package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/pkg/errors"

	"tiktok-mall/app/auth/biz/dal/mysql"
	"tiktok-mall/app/auth/biz/dal/redis"
)

const (
	// Token相关配置
	TokenExpiration        = 24 * time.Hour     // 访问令牌有效期
	RefreshTokenExpiration = 7 * 24 * time.Hour // 刷新令牌有效期
	TokenLength            = 32                 // 令牌长度（字节）
	MaxLoginRetries        = 5                  // 最大登录重试次数
	LoginRetryWindow       = time.Hour          // 登录重试窗口期

	// 用户状态
	UserStatusNormal = 1 // 正常
	UserStatusBanned = 2 // 禁用
)

// AuthService 认证服务
type AuthService struct{}

// NewAuthService 创建认证服务实例
func NewAuthService() *AuthService {
	return &AuthService{}
}

// generateToken 生成随机令牌
func generateToken() (string, error) {
	b := make([]byte, TokenLength)
	if _, err := rand.Read(b); err != nil {
		return "", errors.Wrap(err, "generate random bytes failed")
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// hashPassword 密码加密
func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", errors.Wrap(err, "hash password failed")
	}
	return string(bytes), nil
}

// comparePassword 密码比对
func comparePassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

// validateLoginRetries 验证登录重试次数
func (s *AuthService) validateLoginRetries(ctx context.Context, username string) error {
	count, err := redis.GetLoginRetryCount(ctx, username)
	if err != nil {
		return errors.Wrap(err, "get login retry count failed")
	}

	if count >= MaxLoginRetries {
		return errors.New("too many login attempts, please try again later")
	}

	return nil
}

// createAndCacheTokens 创建并缓存令牌
func (s *AuthService) createAndCacheTokens(ctx context.Context, userID int64) (token, refreshToken string, err error) {
	// 生成访问令牌
	token, err = generateToken()
	if err != nil {
		return "", "", errors.Wrap(err, "generate access token failed")
	}

	// 生成刷新令牌
	refreshToken, err = generateToken()
	if err != nil {
		return "", "", errors.Wrap(err, "generate refresh token failed")
	}

	// 创建Token记录
	tokenRecord := &mysql.Token{
		UserID:       userID,
		Token:        token,
		RefreshToken: refreshToken,
		ExpiredAt:    time.Now().Add(TokenExpiration),
	}

	if err := mysql.CreateToken(tokenRecord); err != nil {
		return "", "", errors.Wrap(err, "create token record failed")
	}

	// 缓存Token
	if err := redis.CacheToken(ctx, token, userID, TokenExpiration); err != nil {
		hlog.CtxWarnf(ctx, "cache token failed: %v", err)
	}

	return token, refreshToken, nil
}

// invalidateTokens 使令牌失效
func (s *AuthService) invalidateTokens(ctx context.Context, token string) error {
	// 从数据库删除令牌
	if err := mysql.DeleteToken(token); err != nil {
		return errors.Wrap(err, "delete token from database failed")
	}

	// 从缓存删除令牌
	if err := redis.DeleteToken(ctx, token); err != nil {
		hlog.CtxWarnf(ctx, "delete token from cache failed: %v", err)
	}

	// 将令牌加入黑名单
	if err := redis.AddToBlacklist(ctx, token, TokenExpiration); err != nil {
		hlog.CtxWarnf(ctx, "add token to blacklist failed: %v", err)
	}

	return nil
}

// Register 用户注册
func (s *AuthService) Register(ctx context.Context, username, password, email, phone string) (int64, string, error) {
	// 检查用户名是否已存在
	exists, err := mysql.CheckUserExists(username)
	if err != nil {
		return 0, "", errors.Wrap(err, "check username exists failed")
	}
	if exists {
		return 0, "", errors.New("username already exists")
	}

	// 检查邮箱是否已被使用
	if email != "" {
		exists, err = mysql.CheckEmailExists(email)
		if err != nil {
			return 0, "", errors.Wrap(err, "check email exists failed")
		}
		if exists {
			return 0, "", errors.New("email already exists")
		}
	}

	// 检查手机号是否已被使用
	if phone != "" {
		exists, err = mysql.CheckPhoneExists(phone)
		if err != nil {
			return 0, "", errors.Wrap(err, "check phone exists failed")
		}
		if exists {
			return 0, "", errors.New("phone already exists")
		}
	}

	// 密码加密
	hashedPassword, err := hashPassword(password)
	if err != nil {
		return 0, "", err
	}

	// 创建用户
	user := &mysql.User{
		Username: username,
		Password: hashedPassword,
		Email:    email,
		Phone:    phone,
		Status:   UserStatusNormal,
	}

	if err := mysql.CreateUser(user); err != nil {
		return 0, "", errors.Wrap(err, "create user failed")
	}

	// 生成访问令牌
	token, _, err := s.createAndCacheTokens(ctx, user.ID)
	if err != nil {
		return 0, "", err
	}

	return user.ID, token, nil
}

// Login 用户登录
func (s *AuthService) Login(ctx context.Context, username, password string) (string, string, error) {
	// 验证登录重试次数
	if err := s.validateLoginRetries(ctx, username); err != nil {
		return "", "", err
	}

	// 获取用户信息
	user, err := mysql.GetUserByUsername(username)
	if err != nil {
		return "", "", errors.Wrap(err, "get user failed")
	}
	if user == nil {
		// 增加登录重试次数
		if _, err := redis.IncrLoginRetry(ctx, username); err != nil {
			hlog.CtxWarnf(ctx, "increment login retry count failed: %v", err)
		}
		return "", "", errors.New("invalid username or password")
	}

	// 检查用户状态
	if user.Status != UserStatusNormal {
		return "", "", errors.New("user is disabled")
	}

	// 验证密码
	if err := comparePassword(user.Password, password); err != nil {
		// 增加登录重试次数
		if _, err := redis.IncrLoginRetry(ctx, username); err != nil {
			hlog.CtxWarnf(ctx, "increment login retry count failed: %v", err)
		}
		return "", "", errors.New("invalid username or password")
	}

	// 重置登录重试次数
	if err := redis.ResetLoginRetry(ctx, username); err != nil {
		hlog.CtxWarnf(ctx, "reset login retry count failed: %v", err)
	}

	// 生成新的访问令牌和刷新令牌
	token, refreshToken, err := s.createAndCacheTokens(ctx, user.ID)
	if err != nil {
		return "", "", err
	}

	return token, refreshToken, nil
}

// RefreshToken 刷新访问令牌
func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (string, string, error) {
	// 获取刷新令牌记录
	tokenRecord, err := mysql.GetTokenByRefreshToken(refreshToken)
	if err != nil {
		return "", "", errors.Wrap(err, "get token record failed")
	}
	if tokenRecord == nil {
		return "", "", errors.New("invalid refresh token")
	}

	// 检查令牌是否过期
	if tokenRecord.ExpiredAt.Before(time.Now()) {
		return "", "", errors.New("refresh token expired")
	}

	// 使旧令牌失效
	if err := s.invalidateTokens(ctx, tokenRecord.Token); err != nil {
		hlog.CtxWarnf(ctx, "invalidate old tokens failed: %v", err)
	}

	// 生成新的访问令牌和刷新令牌
	token, newRefreshToken, err := s.createAndCacheTokens(ctx, tokenRecord.UserID)
	if err != nil {
		return "", "", err
	}

	return token, newRefreshToken, nil
}

// Logout 用户登出
func (s *AuthService) Logout(ctx context.Context, token string) error {
	return s.invalidateTokens(ctx, token)
}

// ValidateToken 验证访问令牌
func (s *AuthService) ValidateToken(ctx context.Context, token string) (int64, string, error) {
	// 检查令牌是否在黑名单中
	inBlacklist, err := redis.IsInBlacklist(ctx, token)
	if err != nil {
		hlog.CtxWarnf(ctx, "check token blacklist failed: %v", err)
	}
	if inBlacklist {
		return 0, "", errors.New("token is invalid")
	}

	// 从缓存获取用户ID
	userID, err := redis.GetCachedUserID(ctx, token)
	if err == nil {
		// 从缓存获取成功，验证用户信息
		user, err := mysql.GetUserByID(userID)
		if err != nil {
			return 0, "", errors.Wrap(err, "get user failed")
		}
		if user == nil || user.Status != UserStatusNormal {
			return 0, "", errors.New("user not found or disabled")
		}
		return user.ID, user.Username, nil
	}

	// 缓存未命中，从数据库查询
	tokenRecord, err := mysql.GetTokenByToken(token)
	if err != nil {
		return 0, "", errors.Wrap(err, "get token record failed")
	}
	if tokenRecord == nil {
		return 0, "", errors.New("token not found")
	}

	// 检查令牌是否过期
	if tokenRecord.ExpiredAt.Before(time.Now()) {
		return 0, "", errors.New("token expired")
	}

	// 获取用户信息
	user, err := mysql.GetUserByID(tokenRecord.UserID)
	if err != nil {
		return 0, "", errors.Wrap(err, "get user failed")
	}
	if user == nil || user.Status != UserStatusNormal {
		return 0, "", errors.New("user not found or disabled")
	}

	// 更新缓存
	if err := redis.CacheToken(ctx, token, user.ID, time.Until(tokenRecord.ExpiredAt)); err != nil {
		hlog.CtxWarnf(ctx, "cache token failed: %v", err)
	}

	return user.ID, user.Username, nil
}

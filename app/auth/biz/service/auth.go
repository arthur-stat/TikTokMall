package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/pkg/errors"

	"TikTokMall/app/auth/biz/dal/mysql"
	"TikTokMall/app/auth/biz/dal/redis"
	"TikTokMall/app/auth/kitex_gen/auth"

	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

const (
	// Token相关配置
	TokenExpiration        = 24 * time.Hour     // 访问令牌有效期
	RefreshTokenExpiration = 7 * 24 * time.Hour // 刷新令牌有效期
	TokenLength            = 32                 // 令牌长度（字节）
	MaxLoginRetries        = 5                  // 最大登录重试次数
	LoginRetryWindow       = time.Hour          // 登录重试窗口期
	maxLoginAttempts       = 5                  // 最大登录尝试次数

	// 用户状态
	UserStatusNormal = 1 // 正常
	UserStatusBanned = 2 // 禁用
)

type authService struct {
	repo AuthRepository
}

func NewAuthService(repo AuthRepository) AuthService {
	return &authService{
		repo: repo,
	}
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
func (s *authService) validateLoginRetries(ctx context.Context, username string, isValidPassword bool) error {
	// 如果密码正确，跳过重试次数验证
	if isValidPassword {
		// 重置重试计数
		if err := redis.ResetLoginRetry(ctx, username); err != nil {
			hlog.CtxWarnf(ctx, "reset login retry count failed: %v", err)
		}
		return nil
	}

	// 增加重试计数
	count, err := redis.IncrLoginRetry(ctx, username)
	if err != nil {
		return errors.Wrap(err, "increment login retry count failed")
	}

	// 检查是否超过重试限制
	if count >= MaxLoginRetries {
		return errors.New("too many login attempts")
	}

	return errors.New("invalid username or password")
}

// createAndCacheTokens 创建并缓存令牌
func (s *authService) createAndCacheTokens(ctx context.Context, userID int64) (token, refreshToken string, err error) {
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

	if err := s.repo.CreateToken(tokenRecord); err != nil {
		return "", "", errors.Wrap(err, "create token record failed")
	}

	// 缓存Token
	if err := redis.CacheToken(ctx, token, userID, TokenExpiration); err != nil {
		hlog.CtxWarnf(ctx, "cache token failed: %v", err)
	}

	return token, refreshToken, nil
}

// invalidateTokens 使令牌失效
func (s *authService) invalidateTokens(ctx context.Context, token string) error {
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

// Register 实现注册功能
func (s *authService) Register(ctx context.Context, req *auth.RegisterRequest) (*auth.RegisterResponse, error) {
	// 检查用户名是否已存在
	exists, err := s.repo.GetUserByUsername(req.Username)
	if err != nil && err != mysql.ErrRecordNotFound {
		return nil, err
	}
	if exists != nil {
		return nil, errors.New("username already exists")
	}

	// 创建用户
	hashedPassword, err := hashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user := &mysql.User{
		Username: req.Username,
		Password: hashedPassword,
		Email:    req.Email,
		Phone:    req.Phone,
		Status:   UserStatusNormal,
	}

	if err := s.repo.CreateUser(user); err != nil {
		return nil, err
	}

	// 生成令牌
	token, err := generateToken()
	if err != nil {
		return nil, err
	}

	// 创建令牌记录
	tokenRecord := &mysql.Token{
		UserID:    user.ID,
		Token:     token,
		ExpiredAt: time.Now().Add(TokenExpiration),
	}

	if err := s.repo.CreateToken(tokenRecord); err != nil {
		return nil, err
	}

	return &auth.RegisterResponse{
		Base: &auth.BaseResp{
			Code:    0,
			Message: "success",
		},
		Data: &auth.RegisterData{
			UserId: user.ID,
			Token:  token,
		},
	}, nil
}

// Login 实现登录功能
func (s *authService) Login(ctx context.Context, req *auth.LoginRequest) (*auth.LoginResponse, error) {
	user, err := s.repo.GetUserByUsername(req.Username)
	if err != nil {
		return nil, err
	}

	if user.Status == UserStatusBanned {
		return nil, errors.New("user is banned")
	}

	// 先验证密码
	isValidPassword := comparePassword(user.Password, req.Password) == nil

	// 然后检查重试次数
	if err := s.validateLoginRetries(ctx, req.Username, isValidPassword); err != nil {
		return nil, err
	}

	// 如果密码验证通过，创建令牌
	token, refreshToken, err := s.createAndCacheTokens(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	return &auth.LoginResponse{
		Base: &auth.BaseResp{
			Code:    0,
			Message: "success",
		},
		Data: &auth.LoginData{
			Token:        token,
			RefreshToken: refreshToken,
		},
	}, nil
}

// RefreshToken 实现刷新令牌功能
func (s *authService) RefreshToken(ctx context.Context, req *auth.RefreshTokenRequest) (*auth.RefreshTokenResponse, error) {
	// 获取令牌记录
	token, err := s.repo.GetTokenByToken(req.RefreshToken)
	if err != nil {
		if err == mysql.ErrRecordNotFound {
			return nil, auth.ErrInvalidToken
		}
		return nil, err
	}

	// 检查令牌是否过期
	if token.ExpiredAt.Before(time.Now()) {
		return nil, auth.ErrTokenExpired
	}

	// 获取用户信息
	user, err := s.repo.GetUserByID(token.UserID)
	if err != nil {
		return nil, err
	}

	// 生成新的令牌
	newToken, newRefreshToken, err := s.createAndCacheTokens(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	return &auth.RefreshTokenResponse{
		Base: &auth.BaseResp{
			Code:    0,
			Message: "success",
		},
		Data: &auth.RefreshTokenData{
			Token:        newToken,
			RefreshToken: newRefreshToken,
		},
	}, nil
}

// ValidateToken 实现验证令牌功能
func (s *authService) ValidateToken(ctx context.Context, req *auth.ValidateTokenRequest) (*auth.ValidateTokenResponse, error) {
	// 获取令牌记录
	token, err := s.repo.GetTokenByToken(req.Token)
	if err != nil {
		if err == mysql.ErrRecordNotFound {
			return &auth.ValidateTokenResponse{
				Base: &auth.BaseResp{
					Code:    int32(consts.StatusUnauthorized),
					Message: "invalid token",
				},
			}, auth.ErrInvalidToken
		}
		return nil, err
	}

	// 检查令牌是否过期
	if token.ExpiredAt.Before(time.Now()) {
		return &auth.ValidateTokenResponse{
			Base: &auth.BaseResp{
				Code:    int32(consts.StatusUnauthorized),
				Message: "token expired",
			},
		}, auth.ErrTokenExpired
	}

	// 获取用户信息
	user, err := s.repo.GetUserByID(token.UserID)
	if err != nil {
		return nil, err
	}

	return &auth.ValidateTokenResponse{
		Base: &auth.BaseResp{
			Code:    int32(consts.StatusOK),
			Message: "success",
		},
		Data: &auth.ValidateTokenData{
			UserId:   user.ID,
			Username: user.Username,
		},
	}, nil
}

// Logout 实现登出功能
func (s *authService) Logout(ctx context.Context, req *auth.LogoutRequest) (*auth.LogoutResponse, error) {
	// 获取令牌记录
	_, err := s.repo.GetTokenByToken(req.Token)
	if err != nil {
		if err == mysql.ErrRecordNotFound {
			return &auth.LogoutResponse{
				Base: &auth.BaseResp{
					Code:    int32(consts.StatusUnauthorized),
					Message: "invalid token",
				},
			}, auth.ErrInvalidToken
		}
		return nil, err
	}

	// 删除令牌
	if err := s.repo.DeleteToken(req.Token); err != nil {
		return nil, err
	}

	// 将令牌加入黑名单
	if err := redis.AddToBlacklist(ctx, req.Token, TokenExpiration); err != nil {
		hlog.CtxWarnf(ctx, "add token to blacklist failed: %v", err)
	}

	return &auth.LogoutResponse{
		Base: &auth.BaseResp{
			Code:    int32(consts.StatusOK),
			Message: "success",
		},
	}, nil
}

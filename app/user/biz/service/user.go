package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/pkg/errors"
	"time"

	"golang.org/x/crypto/bcrypt"

	"TikTokMall/app/user/biz/dal/mysql"
	"TikTokMall/app/user/biz/dal/redis"
)

const (
	TokenExpiration        = 24 * time.Hour
	RefreshTokenExpiration = 7 * 24 * time.Hour
	TokenLength            = 32
	UserCacheExpiration    = 2 * time.Hour
)

type UserService struct{}

func NewUserService() *UserService {
	return &UserService{}
}

// ----------- 独立实现的核心方法 -----------

// generateToken 生成随机令牌
func generateToken() (string, error) {
	b := make([]byte, TokenLength)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// hashPassword 密码加密
func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// validatePassword 密码验证
func validatePassword(hashed, plain string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(plain))
}

// ----------- 业务方法 -----------

// Register 用户注册
func (s *UserService) Register(ctx context.Context, username, password, email, phone string) (int64, string, error) {
	// 检查用户名唯一性
	if exist, _ := mysql.CheckUserExists(username); exist {
		return 0, "", fmt.Errorf("username already exists")
	}

	// 密码加密
	hashedPwd, err := hashPassword(password)
	if err != nil {
		return 0, "", err
	}

	// 创建用户记录
	user := &mysql.User{
		Username: username,
		Password: hashedPwd,
		Email:    email,
		Phone:    phone,
	}
	if err := mysql.CreateUser(user); err != nil {
		return 0, "", err
	}

	// 生成令牌（后续可能要考虑存储令牌？）
	token, err := generateToken()
	if err != nil {
		return 0, "", err
	}

	return user.ID, token, nil
}

// Login 用户登录
func (s *UserService) Login(ctx context.Context, username, password string) (string, string, error) {
	user, err := mysql.GetUserByUsername(username)
	if err != nil || user == nil {
		return "", "", fmt.Errorf("invalid credentials")
	}

	// 验证密码
	if err := validatePassword(user.Password, password); err != nil {
		return "", "", fmt.Errorf("invalid credentials")
	}

	// 生成令牌
	token, _ := generateToken()
	refreshToken, _ := generateToken()

	// 缓存用户信息
	if err := redis.CacheUser(ctx, user, 2*time.Hour); err != nil {
		// 记录日志但不中断流程
	}

	return token, refreshToken, nil
}

// Delete 删除用户
func (s *UserService) Delete(ctx context.Context, userID int64) error {
	// 删除用户记录
	if err := mysql.DeleteUser(userID); err != nil {
		return err
	}

	// 清理缓存
	_ = redis.DeleteUserCache(ctx, userID)
	return nil
}

// Info 获取用户信息
func (s *UserService) Info(ctx context.Context, token string) (int64, string, string, string, error) {
	// 通过Token获取用户ID
	userID, err := s.getUserIDByToken(ctx, token)
	if err != nil {
		return 0, "", "", "", errors.Wrap(err, "token验证失败")
	}

	// 优先从缓存获取
	if user, err := redis.GetUserCache(ctx, userID); err == nil {
		return user.ID, user.Username, user.Email, user.Phone, nil
	}

	// 缓存未命中则查数据库
	user, err := mysql.GetUserByID(userID)
	if err != nil || user == nil {
		return 0, "", "", "", errors.New("用户不存在")
	}

	// 更新缓存
	if err := redis.CacheUser(ctx, user, UserCacheExpiration); err != nil {
		hlog.CtxWarnf(ctx, "缓存用户信息失败: %v", err)
	}

	return user.ID, user.Username, user.Email, user.Phone, nil
}

// Update 更新用户信息
func (s *UserService) Update(ctx context.Context, token string, newUsername, newEmail, newPhone string) error {
	// 通过Token获取用户ID
	userID, err := s.getUserIDByToken(ctx, token)
	if err != nil {
		return errors.Wrap(err, "invalid token")
	}

	// 构建更新字段
	updates := make(map[string]interface{})
	if newUsername != "" {
		// 检查新用户名是否已被占用
		exists, err := mysql.CheckUserExists(newUsername)
		if err != nil {
			return errors.Wrap(err, "check username availability failed")
		}
		if exists {
			return errors.New("username already exists")
		}
		updates["username"] = newUsername
	}
	if newEmail != "" {
		// 邮箱格式验证（可扩展为具体校验逻辑）
		updates["email"] = newEmail
	}
	if newPhone != "" {
		// 手机号格式验证（可扩展为具体校验逻辑）
		updates["phone"] = newPhone
	}

	// 执行数据库更新
	if err := mysql.UpdateUser(userID, updates); err != nil {
		return errors.Wrap(err, "database update failed")
	}

	// 清理用户缓存
	if err := redis.DeleteUserCache(ctx, userID); err != nil {
		hlog.CtxWarnf(ctx, "cache cleanup failed: %v", err)
	}

	return nil
}

// Logout 用户登出
func (s *UserService) Logout(ctx context.Context, token string) error {
	// 从缓存删除Token
	if err := redis.DeleteToken(ctx, token); err != nil {
		hlog.CtxWarnf(ctx, "token cache deletion failed: %v", err)
	}

	// 将Token加入黑名单（复用auth服务逻辑）
	if err := redis.AddToBlacklist(ctx, token, TokenExpiration); err != nil {
		hlog.CtxWarnf(ctx, "blacklist update failed: %v", err)
	}

	// 删除数据库中的Token记录（疑问，有必要吗？）
	if err := mysql.DeleteToken(token); err != nil {
		hlog.CtxWarnf(ctx, "database token cleanup failed: %v", err)
	}

	return nil
}

// getUserIDByToken 通过Token获取用户ID
func (s *UserService) getUserIDByToken(ctx context.Context, token string) (int64, error) {
	// 检查Token是否在黑名单
	if blacklisted, err := redis.IsInBlacklist(ctx, token); err != nil {
		return 0, errors.Wrap(err, "blacklist check failed")
	} else if blacklisted {
		return 0, errors.New("token is invalid")
	}

	// 从缓存获取用户ID
	if userID, err := redis.GetCachedUserID(ctx, token); err == nil {
		return userID, nil
	}

	// 缓存未命中则查询数据库
	tokenRecord, err := mysql.GetTokenByToken(token)
	if err != nil {
		return 0, errors.Wrap(err, "database query failed")
	}
	if tokenRecord == nil {
		return 0, errors.New("token not found")
	}

	// 检查Token是否过期
	if tokenRecord.ExpiredAt.Before(time.Now()) {
		return 0, errors.New("token expired")
	}

	// 缓存结果
	if err := redis.CacheToken(ctx, token, tokenRecord.UserID, time.Until(tokenRecord.ExpiredAt)); err != nil {
		hlog.CtxWarnf(ctx, "cache token failed: %v", err)
	}

	return tokenRecord.UserID, nil
}

package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"

	"TikTokMall/app/user/biz/dal/mysql"
	"TikTokMall/app/user/biz/dal/redis"
)

const (
	TokenExpiration        = 24 * time.Hour
	RefreshTokenExpiration = 7 * 24 * time.Hour
	TokenLength            = 32
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

	// 生成令牌（简化实现，实际需要存储令牌）
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
func (s *UserService) Info(ctx context.Context, userID int64) (*mysql.User, error) {
	// 优先从缓存获取
	if user, err := redis.GetUserCache(ctx, userID); err == nil {
		return user, nil
	}

	// 缓存未命中则查数据库
	return mysql.GetUserByID(userID)
}

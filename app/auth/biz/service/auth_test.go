package service

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"TikTokMall/app/auth/biz/dal"
	"TikTokMall/app/auth/biz/dal/mysql"
	"TikTokMall/app/auth/biz/dal/redis"
)

func TestMain(m *testing.M) {
	// 初始化测试环境
	if err := dal.Init(); err != nil {
		panic(err)
	}

	// 清理测试数据
	if err := mysql.DB.Exec("TRUNCATE TABLE users").Error; err != nil {
		panic(err)
	}
	if err := mysql.DB.Exec("TRUNCATE TABLE tokens").Error; err != nil {
		panic(err)
	}

	// 运行测试
	code := m.Run()

	// 退出
	os.Exit(code)
}

func cleanupTestData() {
	mysql.DB.Exec("TRUNCATE TABLE users")
	mysql.DB.Exec("TRUNCATE TABLE tokens")
	redis.RDB.FlushDB(context.Background())
}

func TestAuthService_Register(t *testing.T) {
	svc := NewAuthService()
	ctx := context.Background()

	tests := []struct {
		name     string
		username string
		password string
		email    string
		phone    string
		wantErr  bool
	}{
		{
			name:     "正常注册",
			username: "testuser1",
			password: "password123",
			email:    "test1@example.com",
			phone:    "13800138001",
			wantErr:  false,
		},
		{
			name:     "用户名已存在",
			username: "testuser1",
			password: "password123",
			email:    "test2@example.com",
			phone:    "13800138002",
			wantErr:  true,
		},
		{
			name:     "邮箱已存在",
			username: "testuser2",
			password: "password123",
			email:    "test1@example.com",
			phone:    "13800138003",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userID, token, err := svc.Register(ctx, tt.username, tt.password, tt.email, tt.phone)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.NotZero(t, userID)
			assert.NotEmpty(t, token)

			// 验证用户是否创建成功
			user, err := mysql.GetUserByUsername(tt.username)
			require.NoError(t, err)
			assert.Equal(t, tt.username, user.Username)
			assert.Equal(t, tt.email, user.Email)
			assert.Equal(t, tt.phone, user.Phone)
		})
	}
}

func TestAuthService_Login(t *testing.T) {
	svc := NewAuthService()
	ctx := context.Background()

	// 创建测试用户
	username := "testlogin"
	password := "password123"
	_, _, err := svc.Register(ctx, username, password, "testlogin@example.com", "13800138004")
	require.NoError(t, err)

	tests := []struct {
		name     string
		username string
		password string
		wantErr  bool
	}{
		{
			name:     "正常登录",
			username: username,
			password: password,
			wantErr:  false,
		},
		{
			name:     "用户名不存在",
			username: "nonexistent",
			password: password,
			wantErr:  true,
		},
		{
			name:     "密码错误",
			username: username,
			password: "wrongpassword",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, refreshToken, err := svc.Login(ctx, tt.username, tt.password)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.NotEmpty(t, token)
			assert.NotEmpty(t, refreshToken)

			// 验证Token是否有效
			userID, _, err := svc.ValidateToken(ctx, token)
			require.NoError(t, err)
			assert.NotZero(t, userID)
		})
	}
}

func TestAuthService_RefreshToken(t *testing.T) {
	svc := NewAuthService()
	ctx := context.Background()

	// 创建测试用户并登录
	username := "testrefresh"
	password := "password123"
	_, _, err := svc.Register(ctx, username, password, "testrefresh@example.com", "13800138005")
	require.NoError(t, err)

	_, refreshToken, err := svc.Login(ctx, username, password)
	require.NoError(t, err)

	tests := []struct {
		name         string
		refreshToken string
		wantErr      bool
	}{
		{
			name:         "正常刷新",
			refreshToken: refreshToken,
			wantErr:      false,
		},
		{
			name:         "无效的刷新令牌",
			refreshToken: "invalid_refresh_token",
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, newRefreshToken, err := svc.RefreshToken(ctx, tt.refreshToken)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.NotEmpty(t, token)
			assert.NotEmpty(t, newRefreshToken)
			assert.NotEqual(t, tt.refreshToken, newRefreshToken)

			// 验证新Token是否有效
			userID, _, err := svc.ValidateToken(ctx, token)
			require.NoError(t, err)
			assert.NotZero(t, userID)
		})
	}
}

func TestAuthService_Logout(t *testing.T) {
	svc := NewAuthService()
	ctx := context.Background()

	// 创建测试用户并登录
	username := "testlogout"
	password := "password123"
	_, _, err := svc.Register(ctx, username, password, "testlogout@example.com", "13800138006")
	require.NoError(t, err)

	token, _, err := svc.Login(ctx, username, password)
	require.NoError(t, err)

	tests := []struct {
		name    string
		token   string
		wantErr bool
	}{
		{
			name:    "正常登出",
			token:   token,
			wantErr: false,
		},
		{
			name:    "重复登出",
			token:   token,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := svc.Logout(ctx, tt.token)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)

			// 验证Token是否已失效
			_, _, err = svc.ValidateToken(ctx, tt.token)
			assert.Error(t, err)
		})
	}
}

func TestAuthService_ValidateToken(t *testing.T) {
	svc := NewAuthService()
	ctx := context.Background()

	// 创建测试用户并登录
	username := "testvalidate"
	password := "password123"
	_, _, err := svc.Register(ctx, username, password, "testvalidate@example.com", "13800138007")
	require.NoError(t, err)

	token, _, err := svc.Login(ctx, username, password)
	require.NoError(t, err)

	tests := []struct {
		name    string
		token   string
		wantErr bool
	}{
		{
			name:    "有效Token",
			token:   token,
			wantErr: false,
		},
		{
			name:    "无效Token",
			token:   "invalid_token",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userID, username, err := svc.ValidateToken(ctx, tt.token)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.NotZero(t, userID)
			assert.NotEmpty(t, username)
		})
	}
}

func TestAuthService_LoginRetryLimit(t *testing.T) {
	svc := NewAuthService()
	ctx := context.Background()

	// 创建测试用户
	username := "testretry"
	password := "password123"
	_, _, err := svc.Register(ctx, username, password, "testretry@example.com", "13800138008")
	require.NoError(t, err)

	// 清除之前的重试记录
	err = redis.ResetLoginRetry(ctx, username)
	require.NoError(t, err)

	// 快速尝试登录直到达到最大重试次数
	for i := 0; i < MaxLoginRetries; i++ {
		_, _, err := svc.Login(ctx, username, "wrongpassword")
		require.Error(t, err)
		if i < MaxLoginRetries-1 {
			require.Contains(t, err.Error(), "invalid username or password")
		}
		time.Sleep(10 * time.Millisecond) // 添加短暂延迟
	}

	// 最后一次尝试应该返回重试次数超限错误
	_, _, err = svc.Login(ctx, username, "wrongpassword")
	require.Error(t, err)
	require.Contains(t, err.Error(), "too many login attempts")
}

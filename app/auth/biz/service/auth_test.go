package service

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"TikTokMall/app/auth/biz/dal/mysql"
	"TikTokMall/app/auth/biz/dal/redis"
	servicemock "TikTokMall/app/auth/biz/service/mock"
	"TikTokMall/app/auth/kitex_gen/auth"

	"gorm.io/gorm"
)

// 添加 mock Redis 客户端
type mockRedis struct {
	retryCount int
}

func (m *mockRedis) CacheToken(ctx context.Context, token string, userID int64, expiration time.Duration) error {
	return nil
}

func (m *mockRedis) DeleteToken(ctx context.Context, token string) error {
	return nil
}

func (m *mockRedis) AddToBlacklist(ctx context.Context, token string, expiration time.Duration) error {
	return nil
}

func (m *mockRedis) GetLoginRetryCount(ctx context.Context, username string) (int, error) {
	return m.retryCount, nil
}

func (m *mockRedis) IncrLoginRetry(ctx context.Context, username string) (int, error) {
	m.retryCount++
	return m.retryCount, nil
}

func (m *mockRedis) ResetLoginRetry(ctx context.Context, username string) error {
	m.retryCount = 0
	return nil
}

func init() {
	// 设置测试用的 DB
	mysql.DB = &gorm.DB{}
	// 设置 mock Redis 客户端
	redis.Client = &mockRedis{}
}

func TestMain(m *testing.M) {
	// 不需要初始化实际的数据库
	code := m.Run()
	os.Exit(code)
}

func TestAuthService_Register(t *testing.T) {
	mockRepo := new(servicemock.MockAuthRepository)
	svc := NewAuthService(mockRepo)

	tests := []struct {
		name    string
		req     *auth.RegisterRequest
		setup   func()
		wantErr bool
	}{
		{
			name: "success",
			req: &auth.RegisterRequest{
				Username: "testuser",
				Password: "password123",
				Email:    "test@example.com",
				Phone:    "13800138000",
			},
			setup: func() {
				mockRepo.On("GetUserByUsername", "testuser").Return(nil, mysql.ErrRecordNotFound)
				mockRepo.On("CreateUser", mock.Anything).Run(func(args mock.Arguments) {
					user := args.Get(0).(*mysql.User)
					user.ID = 1
				}).Return(nil)
				mockRepo.On("CreateToken", mock.Anything).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "username_exists",
			req: &auth.RegisterRequest{
				Username: "existinguser",
				Password: "password123",
				Email:    "test@example.com",
				Phone:    "13800138000",
			},
			setup: func() {
				mockRepo.On("GetUserByUsername", "existinguser").Return(&mysql.User{
					Username: "existinguser",
				}, nil)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 重置 mock
			mockRepo = new(servicemock.MockAuthRepository)
			svc = NewAuthService(mockRepo)

			if tt.setup != nil {
				tt.setup()
			}

			resp, err := svc.Register(context.Background(), tt.req)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, resp)
			assert.NotNil(t, resp.Data)
			assert.NotEmpty(t, resp.Data.Token)
			assert.NotZero(t, resp.Data.UserId)

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestAuthService_Login(t *testing.T) {
	mockRepo := new(servicemock.MockAuthRepository)
	svc := NewAuthService(mockRepo)

	// 设置 mock Redis 客户端
	redisClient := &mockRedis{
		retryCount: 0,
	}
	redis.Client = redisClient

	tests := []struct {
		name    string
		req     *auth.LoginRequest
		setup   func()
		wantErr bool
	}{
		{
			name: "success",
			req: &auth.LoginRequest{
				Username: "testuser",
				Password: "password123",
			},
			setup: func() {
				hashedPassword, _ := hashPassword("password123")
				mockRepo.On("GetUserByUsername", "testuser").Return(&mysql.User{
					ID:       1,
					Username: "testuser",
					Password: hashedPassword,
					Status:   UserStatusNormal,
				}, nil)
				mockRepo.On("CreateToken", mock.Anything).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "user_not_found",
			req: &auth.LoginRequest{
				Username: "nonexistent",
				Password: "password123",
			},
			setup: func() {
				mockRepo.On("GetUserByUsername", "nonexistent").Return(nil, mysql.ErrRecordNotFound)
			},
			wantErr: true,
		},
		{
			name: "wrong_password",
			req: &auth.LoginRequest{
				Username: "testuser",
				Password: "wrongpassword",
			},
			setup: func() {
				hashedPassword, _ := hashPassword("password123")
				mockRepo.On("GetUserByUsername", "testuser").Return(&mysql.User{
					ID:       1,
					Username: "testuser",
					Password: hashedPassword,
					Status:   UserStatusNormal,
				}, nil)
			},
			wantErr: true,
		},
		{
			name: "user_banned",
			req: &auth.LoginRequest{
				Username: "banneduser",
				Password: "password123",
			},
			setup: func() {
				hashedPassword, _ := hashPassword("password123")
				mockRepo.On("GetUserByUsername", "banneduser").Return(&mysql.User{
					ID:       1,
					Username: "banneduser",
					Password: hashedPassword,
					Status:   UserStatusBanned,
				}, nil)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup()
			}

			resp, err := svc.Login(context.Background(), tt.req)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, resp)
			assert.NotNil(t, resp.Data)
			assert.NotEmpty(t, resp.Data.Token)
			assert.NotEmpty(t, resp.Data.RefreshToken)

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestAuthService_RefreshToken(t *testing.T) {
	mockRepo := new(servicemock.MockAuthRepository)
	svc := NewAuthService(mockRepo)

	tests := []struct {
		name    string
		req     *auth.RefreshTokenRequest
		setup   func()
		wantErr bool
	}{
		{
			name: "success",
			req: &auth.RefreshTokenRequest{
				RefreshToken: "valid-refresh-token",
			},
			setup: func() {
				mockRepo.On("GetTokenByToken", "valid-refresh-token").Return(&mysql.Token{
					UserID:    1,
					ExpiredAt: time.Now().Add(24 * time.Hour),
				}, nil)
				mockRepo.On("GetUserByID", int64(1)).Return(&mysql.User{
					ID:       1,
					Username: "testuser",
					Status:   UserStatusNormal,
				}, nil)
				mockRepo.On("CreateToken", mock.Anything).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "invalid_token",
			req: &auth.RefreshTokenRequest{
				RefreshToken: "invalid-token",
			},
			setup: func() {
				mockRepo.On("GetTokenByToken", "invalid-token").Return(nil, mysql.ErrRecordNotFound)
			},
			wantErr: true,
		},
		{
			name: "expired_token",
			req: &auth.RefreshTokenRequest{
				RefreshToken: "expired-token",
			},
			setup: func() {
				mockRepo.On("GetTokenByToken", "expired-token").Return(&mysql.Token{
					UserID:    1,
					ExpiredAt: time.Now().Add(-1 * time.Hour),
				}, nil)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup()
			}

			resp, err := svc.RefreshToken(context.Background(), tt.req)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, resp)
			assert.NotEmpty(t, resp.Data.Token)
			assert.NotEmpty(t, resp.Data.RefreshToken)

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestAuthService_Logout(t *testing.T) {
	mockRepo := new(servicemock.MockAuthRepository)
	svc := NewAuthService(mockRepo)

	tests := []struct {
		name    string
		req     *auth.LogoutRequest
		setup   func()
		wantErr bool
	}{
		{
			name: "success",
			req: &auth.LogoutRequest{
				Token: "valid-token",
			},
			setup: func() {
				mockRepo.On("GetTokenByToken", "valid-token").Return(&mysql.Token{
					UserID: 1,
					Token:  "valid-token",
				}, nil)
				mockRepo.On("DeleteToken", "valid-token").Return(nil)
			},
			wantErr: false,
		},
		{
			name: "invalid_token",
			req: &auth.LogoutRequest{
				Token: "invalid-token",
			},
			setup: func() {
				mockRepo.On("GetTokenByToken", "invalid-token").Return(nil, mysql.ErrRecordNotFound)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 重置 mock
			mockRepo = new(servicemock.MockAuthRepository)
			svc = NewAuthService(mockRepo)

			if tt.setup != nil {
				tt.setup()
			}

			resp, err := svc.Logout(context.Background(), tt.req)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, resp)

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestAuthService_ValidateToken(t *testing.T) {
	mockRepo := new(servicemock.MockAuthRepository)
	svc := NewAuthService(mockRepo)

	tests := []struct {
		name    string
		req     *auth.ValidateTokenRequest
		setup   func()
		wantErr bool
	}{
		{
			name: "success",
			req: &auth.ValidateTokenRequest{
				Token: "valid-token",
			},
			setup: func() {
				// 模拟有效的访问令牌
				mockRepo.On("GetTokenByToken", "valid-token").Return(&mysql.Token{
					UserID:    1,
					ExpiredAt: time.Now().Add(24 * time.Hour),
				}, nil)
				mockRepo.On("GetUserByID", int64(1)).Return(&mysql.User{
					ID:       1,
					Username: "testuser",
					Status:   UserStatusNormal,
				}, nil)
			},
			wantErr: false,
		},
		{
			name: "invalid_token",
			req: &auth.ValidateTokenRequest{
				Token: "invalid-token",
			},
			setup: func() {
				mockRepo.On("GetTokenByToken", "invalid-token").Return(nil, mysql.ErrRecordNotFound)
			},
			wantErr: true,
		},
		{
			name: "expired_token",
			req: &auth.ValidateTokenRequest{
				Token: "expired-token",
			},
			setup: func() {
				mockRepo.On("GetTokenByToken", "expired-token").Return(&mysql.Token{
					UserID:    1,
					ExpiredAt: time.Now().Add(-1 * time.Hour),
				}, nil)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup()
			}

			resp, err := svc.ValidateToken(context.Background(), tt.req)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, resp)
			assert.NotNil(t, resp.Data)
			assert.NotZero(t, resp.Data.UserId)
			assert.NotEmpty(t, resp.Data.Username)

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestAuthService_LoginRetryLimit(t *testing.T) {
	mockRepo := new(servicemock.MockAuthRepository)
	svc := NewAuthService(mockRepo)
	ctx := context.Background()

	// 设置 mock 预期
	hashedPassword, _ := hashPassword("password123")
	// 使用 Times 来指定调用次数
	mockRepo.On("GetUserByUsername", "testretry").Return(&mysql.User{
		ID:       1,
		Username: "testretry",
		Password: hashedPassword,
		Status:   UserStatusNormal,
	}, nil).Times(MaxLoginRetries + 1) // 预期会调用 MaxLoginRetries + 1 次

	// 模拟 Redis 重试计数
	redisClient := &mockRedis{
		retryCount: 0,
	}
	redis.Client = redisClient

	// 快速尝试登录直到达到最大重试次数
	for i := 0; i < MaxLoginRetries; i++ {
		loginReq := &auth.LoginRequest{
			Username: "testretry",
			Password: "wrongpassword",
		}
		_, err := svc.Login(ctx, loginReq)
		require.Error(t, err)
		require.Contains(t, err.Error(), "invalid username or password")
		redisClient.retryCount++ // 增加重试计数
	}

	// 最后一次尝试，此时重试计数已经达到上限
	loginReq := &auth.LoginRequest{
		Username: "testretry",
		Password: "wrongpassword",
	}
	_, err := svc.Login(ctx, loginReq)
	require.Error(t, err)
	require.Contains(t, err.Error(), "too many login attempts")

	mockRepo.AssertExpectations(t)
}

package service

import (
	"context"

	"TikTokMall/app/auth/biz/dal/mysql"
	"TikTokMall/app/auth/kitex_gen/auth"
)

// AuthService 定义认证服务接口
type AuthService interface {
	Register(ctx context.Context, req *auth.RegisterRequest) (*auth.RegisterResponse, error)
	Login(ctx context.Context, req *auth.LoginRequest) (*auth.LoginResponse, error)
	RefreshToken(ctx context.Context, req *auth.RefreshTokenRequest) (*auth.RefreshTokenResponse, error)
	ValidateToken(ctx context.Context, req *auth.ValidateTokenRequest) (*auth.ValidateTokenResponse, error)
	Logout(ctx context.Context, req *auth.LogoutRequest) (*auth.LogoutResponse, error)
}

// AuthRepository 定义数据访问接口
type AuthRepository interface {
	CreateUser(user *mysql.User) error
	GetUserByUsername(username string) (*mysql.User, error)
	GetUserByID(id int64) (*mysql.User, error)
	CreateToken(token *mysql.Token) error
	GetTokenByToken(token string) (*mysql.Token, error)
	DeleteToken(token string) error
}

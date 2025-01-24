package service

import (
	"context"
	auth "TikTokMall/app/auth/kitex_gen/auth"
)

type LoginService struct {
	ctx context.Context
} // NewLoginService new LoginService
func NewLoginService(ctx context.Context) *LoginService {
	return &LoginService{ctx: ctx}
}

// Run create note info
func (s *LoginService) Run(req *auth.LoginRequest) (resp *auth.LoginResponse, err error) {
	// Finish your business logic.

	return
}

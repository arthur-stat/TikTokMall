package service

import (
	"context"
	auth "TikTokMall/app/auth/kitex_gen/auth"
)

type RegisterService struct {
	ctx context.Context
} // NewRegisterService new RegisterService
func NewRegisterService(ctx context.Context) *RegisterService {
	return &RegisterService{ctx: ctx}
}

// Run create note info
func (s *RegisterService) Run(req *auth.RegisterRequest) (resp *auth.RegisterResponse, err error) {
	// Finish your business logic.

	return
}

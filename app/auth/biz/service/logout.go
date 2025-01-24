package service

import (
	"context"
	auth "TikTokMall/app/auth/kitex_gen/auth"
)

type LogoutService struct {
	ctx context.Context
} // NewLogoutService new LogoutService
func NewLogoutService(ctx context.Context) *LogoutService {
	return &LogoutService{ctx: ctx}
}

// Run create note info
func (s *LogoutService) Run(req *auth.LogoutRequest) (resp *auth.LogoutResponse, err error) {
	// Finish your business logic.

	return
}

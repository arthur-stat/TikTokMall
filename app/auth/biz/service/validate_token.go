package service

import (
	"context"
	auth "TikTokMall/app/auth/kitex_gen/auth"
)

type ValidateTokenService struct {
	ctx context.Context
} // NewValidateTokenService new ValidateTokenService
func NewValidateTokenService(ctx context.Context) *ValidateTokenService {
	return &ValidateTokenService{ctx: ctx}
}

// Run create note info
func (s *ValidateTokenService) Run(req *auth.ValidateTokenRequest) (resp *auth.ValidateTokenResponse, err error) {
	// Finish your business logic.

	return
}

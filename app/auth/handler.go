package main

import (
	auth "TikTokMall/app/auth/kitex_gen/auth"
	"context"
	"TikTokMall/app/auth/biz/service"
)

// AuthServiceImpl implements the last service interface defined in the IDL.
type AuthServiceImpl struct{}

// DeliverTokenByRPC implements the AuthServiceImpl interface.
func (s *AuthServiceImpl) DeliverTokenByRPC(ctx context.Context, req *auth.DeliverTokenReq) (resp *auth.DeliveryResp, err error) {
	resp, err = service.NewDeliverTokenByRPCService(ctx).Run(req)

	return resp, err
}

// VerifyTokenByRPC implements the AuthServiceImpl interface.
func (s *AuthServiceImpl) VerifyTokenByRPC(ctx context.Context, req *auth.VerifyTokenReq) (resp *auth.VerifyResp, err error) {
	resp, err = service.NewVerifyTokenByRPCService(ctx).Run(req)

	return resp, err
}

// Register implements the AuthServiceImpl interface.
func (s *AuthServiceImpl) Register(ctx context.Context, req *auth.RegisterRequest) (resp *auth.RegisterResponse, err error) {
	resp, err = service.NewRegisterService(ctx).Run(req)

	return resp, err
}

// Login implements the AuthServiceImpl interface.
func (s *AuthServiceImpl) Login(ctx context.Context, req *auth.LoginRequest) (resp *auth.LoginResponse, err error) {
	resp, err = service.NewLoginService(ctx).Run(req)

	return resp, err
}

// RefreshToken implements the AuthServiceImpl interface.
func (s *AuthServiceImpl) RefreshToken(ctx context.Context, req *auth.RefreshTokenRequest) (resp *auth.RefreshTokenResponse, err error) {
	resp, err = service.NewRefreshTokenService(ctx).Run(req)

	return resp, err
}

// Logout implements the AuthServiceImpl interface.
func (s *AuthServiceImpl) Logout(ctx context.Context, req *auth.LogoutRequest) (resp *auth.LogoutResponse, err error) {
	resp, err = service.NewLogoutService(ctx).Run(req)

	return resp, err
}

// ValidateToken implements the AuthServiceImpl interface.
func (s *AuthServiceImpl) ValidateToken(ctx context.Context, req *auth.ValidateTokenRequest) (resp *auth.ValidateTokenResponse, err error) {
	resp, err = service.NewValidateTokenService(ctx).Run(req)

	return resp, err
}

// Code generated by Kitex v0.9.1. DO NOT EDIT.

package authservice

import (
	auth "TikTokMall/app/auth/kitex_gen/auth"
	"context"
	client "github.com/cloudwego/kitex/client"
	callopt "github.com/cloudwego/kitex/client/callopt"
)

// Client is designed to provide IDL-compatible methods with call-option parameter for kitex framework.
type Client interface {
	Register(ctx context.Context, Req *auth.RegisterRequest, callOptions ...callopt.Option) (r *auth.RegisterResponse, err error)
	Login(ctx context.Context, Req *auth.LoginRequest, callOptions ...callopt.Option) (r *auth.LoginResponse, err error)
	RefreshToken(ctx context.Context, Req *auth.RefreshTokenRequest, callOptions ...callopt.Option) (r *auth.RefreshTokenResponse, err error)
	Logout(ctx context.Context, Req *auth.LogoutRequest, callOptions ...callopt.Option) (r *auth.LogoutResponse, err error)
	ValidateToken(ctx context.Context, Req *auth.ValidateTokenRequest, callOptions ...callopt.Option) (r *auth.ValidateTokenResponse, err error)
}

// NewClient creates a client for the service defined in IDL.
func NewClient(destService string, opts ...client.Option) (Client, error) {
	var options []client.Option
	options = append(options, client.WithDestService(destService))

	options = append(options, opts...)

	kc, err := client.NewClient(serviceInfo(), options...)
	if err != nil {
		return nil, err
	}
	return &kAuthServiceClient{
		kClient: newServiceClient(kc),
	}, nil
}

// MustNewClient creates a client for the service defined in IDL. It panics if any error occurs.
func MustNewClient(destService string, opts ...client.Option) Client {
	kc, err := NewClient(destService, opts...)
	if err != nil {
		panic(err)
	}
	return kc
}

type kAuthServiceClient struct {
	*kClient
}

func (p *kAuthServiceClient) Register(ctx context.Context, Req *auth.RegisterRequest, callOptions ...callopt.Option) (r *auth.RegisterResponse, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.Register(ctx, Req)
}

func (p *kAuthServiceClient) Login(ctx context.Context, Req *auth.LoginRequest, callOptions ...callopt.Option) (r *auth.LoginResponse, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.Login(ctx, Req)
}

func (p *kAuthServiceClient) RefreshToken(ctx context.Context, Req *auth.RefreshTokenRequest, callOptions ...callopt.Option) (r *auth.RefreshTokenResponse, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.RefreshToken(ctx, Req)
}

func (p *kAuthServiceClient) Logout(ctx context.Context, Req *auth.LogoutRequest, callOptions ...callopt.Option) (r *auth.LogoutResponse, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.Logout(ctx, Req)
}

func (p *kAuthServiceClient) ValidateToken(ctx context.Context, Req *auth.ValidateTokenRequest, callOptions ...callopt.Option) (r *auth.ValidateTokenResponse, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.ValidateToken(ctx, Req)
}

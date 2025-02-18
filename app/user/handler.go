package main

import (
	"context"
	"time"

	"TikTokMall/app/user/biz/service"
	"TikTokMall/app/user/kitex_gen/user"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/kitex-contrib/obs-opentelemetry/tracing"
	"go.opentelemetry.io/otel/attribute"
)

// UserServiceImpl implements the last service interface defined in the IDL.
type UserServiceImpl struct{}

// NewUserServiceImpl creates a new instance of UserServiceImpl
func NewUserServiceImpl() *UserServiceImpl {
	return &UserServiceImpl{}
}

// Register implements the UserServiceImpl interface.
func (s *UserServiceImpl) Register(ctx context.Context, req *user.RegisterReq) (resp *user.RegisterResp, err error) {
	// Start span for tracing
	ctx, span := tracing.StartServerSpan(ctx, "UserService.Register")
	defer span.End()

	// Add request attributes to span
	span.SetAttributes(
		attribute.String("email", req.Email),
		attribute.Bool("has_password", len(req.Password) > 0),
	)

	// Set operation timeout
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	// Validate request
	if err := validateRegisterRequest(req); err != nil {
		klog.CtxErrorf(ctx, "Register request validation failed: %v", err)
		span.RecordError(err)
		return nil, err
	}

	// Process request
	startTime := time.Now()
	resp, err = service.NewRegisterService(ctx).Run(req)
	
	// Log result
	if err != nil {
		klog.CtxErrorf(ctx, "Register failed: %v", err)
		span.RecordError(err)
	} else {
		klog.CtxInfof(ctx, "Register success: userId=%d, took=%v", resp.UserId, time.Since(startTime))
	}

	return resp, err
}

// Login implements the UserServiceImpl interface.
func (s *UserServiceImpl) Login(ctx context.Context, req *user.LoginReq) (resp *user.LoginResp, err error) {
	// Start span for tracing
	ctx, span := tracing.StartServerSpan(ctx, "UserService.Login")
	defer span.End()

	// Add request attributes to span
	span.SetAttributes(
		attribute.String("email", req.Email),
		attribute.Bool("has_password", len(req.Password) > 0),
	)

	// Set operation timeout
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	// Validate request
	if err := validateLoginRequest(req); err != nil {
		klog.CtxErrorf(ctx, "Login request validation failed: %v", err)
		span.RecordError(err)
		return nil, err
	}

	// Process request
	startTime := time.Now()
	resp, err = service.NewLoginService(ctx).Run(req)
	
	// Log result
	if err != nil {
		klog.CtxErrorf(ctx, "Login failed: %v", err)
		span.RecordError(err)
	} else {
		klog.CtxInfof(ctx, "Login success: userId=%d, took=%v", resp.UserId, time.Since(startTime))
	}

	return resp, err
}

// Request validation helpers
func validateRegisterRequest(req *user.RegisterReq) error {
	if req == nil {
		return ErrInvalidRequest
	}
	if len(req.Email) == 0 {
		return ErrEmailRequired
	}
	if len(req.Password) == 0 {
		return ErrPasswordRequired
	}
	if len(req.ConfirmPassword) == 0 {
		return ErrConfirmPasswordRequired
	}
	if len(req.Password) < 8 {
		return ErrPasswordTooShort
	}
	return nil
}

func validateLoginRequest(req *user.LoginReq) error {
	if req == nil {
		return ErrInvalidRequest
	}
	if len(req.Email) == 0 {
		return ErrEmailRequired
	}
	if len(req.Password) == 0 {
		return ErrPasswordRequired
	}
	return nil
}
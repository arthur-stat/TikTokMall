package handler

import (
	"context"
	"strings"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"

	"TikTokMall/app/auth/biz/service"
	"TikTokMall/app/auth/kitex_gen/auth"
	"TikTokMall/app/auth/repository/mysql"
)

// AuthHandler 认证服务处理器
type AuthHandler struct {
	svc service.AuthService
}

// NewAuthHandler 创建认证服务处理器
func NewAuthHandler() *AuthHandler {
	repo := mysql.NewAuthRepository()
	svc := service.NewAuthService(repo)
	return &AuthHandler{
		svc: svc,
	}
}

// Register 处理用户注册请求
func (h *AuthHandler) Register(ctx context.Context, c *app.RequestContext) {
	var req auth.RegisterRequest
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(consts.StatusBadRequest, &auth.RegisterResponse{
			Base: &auth.BaseResp{
				Code:    consts.StatusBadRequest,
				Message: err.Error(),
			},
		})
		return
	}

	resp, err := h.svc.Register(ctx, &req)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, &auth.RegisterResponse{
			Base: &auth.BaseResp{
				Code:    consts.StatusInternalServerError,
				Message: err.Error(),
			},
		})
		return
	}

	c.JSON(consts.StatusOK, resp)
}

// Login 处理用户登录请求
func (h *AuthHandler) Login(ctx context.Context, c *app.RequestContext) {
	var req auth.LoginRequest
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(consts.StatusBadRequest, &auth.LoginResponse{
			Base: &auth.BaseResp{
				Code:    consts.StatusBadRequest,
				Message: err.Error(),
			},
		})
		return
	}

	resp, err := h.svc.Login(ctx, &req)
	if err != nil {
		c.JSON(consts.StatusUnauthorized, &auth.LoginResponse{
			Base: &auth.BaseResp{
				Code:    consts.StatusUnauthorized,
				Message: err.Error(),
			},
		})
		return
	}

	c.JSON(consts.StatusOK, resp)
}

// RefreshToken 处理令牌刷新请求
func (h *AuthHandler) RefreshToken(ctx context.Context, c *app.RequestContext) {
	var req auth.RefreshTokenRequest
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(consts.StatusBadRequest, &auth.RefreshTokenResponse{
			Base: &auth.BaseResp{
				Code:    consts.StatusBadRequest,
				Message: err.Error(),
			},
		})
		return
	}

	// 从Authorization头中获取refreshToken
	req.RefreshToken = strings.TrimPrefix(req.RefreshToken, "Bearer ")

	resp, err := h.svc.RefreshToken(ctx, &req)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, &auth.RefreshTokenResponse{
			Base: &auth.BaseResp{
				Code:    consts.StatusInternalServerError,
				Message: err.Error(),
			},
		})
		return
	}

	c.JSON(consts.StatusOK, resp)
}

// Logout 处理登出请求
func (h *AuthHandler) Logout(ctx context.Context, c *app.RequestContext) {
	var req auth.LogoutRequest
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(consts.StatusBadRequest, &auth.LogoutResponse{
			Base: &auth.BaseResp{
				Code:    consts.StatusBadRequest,
				Message: err.Error(),
			},
		})
		return
	}

	req.Token = strings.TrimPrefix(req.Token, "Bearer ")

	resp, err := h.svc.Logout(ctx, &req)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, &auth.LogoutResponse{
			Base: &auth.BaseResp{
				Code:    consts.StatusInternalServerError,
				Message: err.Error(),
			},
		})
		return
	}

	c.JSON(consts.StatusOK, resp)
}

// ValidateToken 处理令牌验证请求
func (h *AuthHandler) ValidateToken(ctx context.Context, c *app.RequestContext) {
	var req auth.ValidateTokenRequest
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(consts.StatusBadRequest, &auth.ValidateTokenResponse{
			Base: &auth.BaseResp{
				Code:    consts.StatusBadRequest,
				Message: err.Error(),
			},
		})
		return
	}

	req.Token = strings.TrimPrefix(req.Token, "Bearer ")

	resp, err := h.svc.ValidateToken(ctx, &req)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, &auth.ValidateTokenResponse{
			Base: &auth.BaseResp{
				Code:    consts.StatusInternalServerError,
				Message: err.Error(),
			},
		})
		return
	}

	c.JSON(consts.StatusOK, resp)
}

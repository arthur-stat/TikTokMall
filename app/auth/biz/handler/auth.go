package handler

import (
	"context"
	"strings"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"

	"tiktok-mall/app/auth/biz/service"
	"tiktok-mall/kitex_gen/auth"
)

// AuthHandler 认证服务处理器
type AuthHandler struct {
	svc *service.AuthService
}

// NewAuthHandler 创建认证服务处理器
func NewAuthHandler() *AuthHandler {
	return &AuthHandler{
		svc: service.NewAuthService(),
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

	// 调用服务层处理注册
	userID, token, err := h.svc.Register(ctx, req.Username, req.Password, req.Email, req.Phone)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, &auth.RegisterResponse{
			Base: &auth.BaseResp{
				Code:    consts.StatusInternalServerError,
				Message: err.Error(),
			},
		})
		return
	}

	c.JSON(consts.StatusOK, &auth.RegisterResponse{
		Base: &auth.BaseResp{
			Code:    consts.StatusOK,
			Message: "success",
		},
		Data: &auth.RegisterData{
			UserId: userID,
			Token:  token,
		},
	})
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

	// 调用服务层处理登录
	token, refreshToken, err := h.svc.Login(ctx, req.Username, req.Password)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, &auth.LoginResponse{
			Base: &auth.BaseResp{
				Code:    consts.StatusInternalServerError,
				Message: err.Error(),
			},
		})
		return
	}

	c.JSON(consts.StatusOK, &auth.LoginResponse{
		Base: &auth.BaseResp{
			Code:    consts.StatusOK,
			Message: "success",
		},
		Data: &auth.LoginData{
			Token:        token,
			RefreshToken: refreshToken,
		},
	})
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
	refreshToken := strings.TrimPrefix(req.RefreshToken, "Bearer ")

	// 调用服务层处理令牌刷新
	token, newRefreshToken, err := h.svc.RefreshToken(ctx, refreshToken)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, &auth.RefreshTokenResponse{
			Base: &auth.BaseResp{
				Code:    consts.StatusInternalServerError,
				Message: err.Error(),
			},
		})
		return
	}

	c.JSON(consts.StatusOK, &auth.RefreshTokenResponse{
		Base: &auth.BaseResp{
			Code:    consts.StatusOK,
			Message: "success",
		},
		Data: &auth.RefreshTokenData{
			Token:        token,
			RefreshToken: newRefreshToken,
		},
	})
}

// Logout 处理用户登出请求
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

	// 从Authorization头中获取token
	token := strings.TrimPrefix(req.Token, "Bearer ")

	// 调用服务层处理登出
	if err := h.svc.Logout(ctx, token); err != nil {
		c.JSON(consts.StatusInternalServerError, &auth.LogoutResponse{
			Base: &auth.BaseResp{
				Code:    consts.StatusInternalServerError,
				Message: err.Error(),
			},
		})
		return
	}

	c.JSON(consts.StatusOK, &auth.LogoutResponse{
		Base: &auth.BaseResp{
			Code:    consts.StatusOK,
			Message: "success",
		},
	})
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

	// 从Authorization头中获取token
	token := strings.TrimPrefix(req.Token, "Bearer ")

	// 调用服务层处理令牌验证
	userID, username, err := h.svc.ValidateToken(ctx, token)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, &auth.ValidateTokenResponse{
			Base: &auth.BaseResp{
				Code:    consts.StatusInternalServerError,
				Message: err.Error(),
			},
		})
		return
	}

	c.JSON(consts.StatusOK, &auth.ValidateTokenResponse{
		Base: &auth.BaseResp{
			Code:    consts.StatusOK,
			Message: "success",
		},
		Data: &auth.ValidateTokenData{
			UserId:   userID,
			Username: username,
		},
	})
}

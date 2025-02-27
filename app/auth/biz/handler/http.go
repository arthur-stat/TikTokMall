package handler

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"

	"TikTokMall/app/auth/biz/service"
	"TikTokMall/app/auth/kitex_gen/auth"
)

// 需要先定义 AuthHTTPHandler
type AuthHTTPHandler struct {
	svc service.AuthService
}

func NewAuthHTTPHandler(svc service.AuthService) *AuthHTTPHandler {
	return &AuthHTTPHandler{
		svc: svc,
	}
}

func (h *AuthHTTPHandler) Register(ctx context.Context, c *app.RequestContext) {
	var req auth.RegisterRequest

	// 调试请求内容
	reqBody := c.Request.Body()
	fmt.Printf("Request body: %s\n", string(reqBody))

	// 尝试直接解析JSON而不是使用Bind
	if err := json.Unmarshal(reqBody, &req); err != nil {
		c.JSON(consts.StatusBadRequest, &auth.RegisterResponse{
			Base: &auth.BaseResp{
				Code:    consts.StatusBadRequest,
				Message: "invalid request format: " + err.Error(),
			},
		})
		return
	}

	// 输出解析后的请求内容进行调试
	fmt.Printf("Parsed request: %+v\n", req)

	// 验证必填字段
	if req.Username == "" || req.Password == "" {
		c.JSON(consts.StatusBadRequest, &auth.RegisterResponse{
			Base: &auth.BaseResp{
				Code:    consts.StatusBadRequest,
				Message: "username and password are required",
			},
		})
		return
	}

	// 调用 service 层的 Register 方法
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

func validateRegisterRequest(req *auth.RegisterRequest) error {
	if req.Username == "" || req.Password == "" {
		return fmt.Errorf("username and password are required")
	}
	return nil
}

func (h *AuthHTTPHandler) Login(ctx context.Context, c *app.RequestContext) {
	var req auth.LoginRequest
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(consts.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	resp, err := h.svc.Login(ctx, &req)
	if err != nil {
		c.JSON(consts.StatusUnauthorized, map[string]interface{}{
			"error": "invalid username or password",
		})
		return
	}

	c.JSON(consts.StatusOK, resp)
}

// RefreshToken 处理令牌刷新请求
func (h *AuthHTTPHandler) RefreshToken(ctx context.Context, c *app.RequestContext) {
	var req auth.RefreshTokenRequest
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(consts.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	resp, err := h.svc.RefreshToken(ctx, &req)
	if err != nil {
		statusCode := consts.StatusInternalServerError
		if err == auth.ErrInvalidToken {
			statusCode = consts.StatusUnauthorized
		}
		c.JSON(statusCode, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	c.JSON(consts.StatusOK, resp)
}

// ValidateToken 处理令牌验证请求
func (h *AuthHTTPHandler) ValidateToken(ctx context.Context, c *app.RequestContext) {
	var req auth.ValidateTokenRequest
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(consts.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	resp, err := h.svc.ValidateToken(ctx, &req)
	if err != nil {
		statusCode := consts.StatusInternalServerError
		if err == auth.ErrInvalidToken {
			statusCode = consts.StatusUnauthorized
		}
		c.JSON(statusCode, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	c.JSON(consts.StatusOK, resp)
}

// Logout 处理登出请求
func (h *AuthHTTPHandler) Logout(ctx context.Context, c *app.RequestContext) {
	var req auth.LogoutRequest
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(consts.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	resp, err := h.svc.Logout(ctx, &req)
	if err != nil {
		statusCode := consts.StatusInternalServerError
		if err == auth.ErrInvalidToken {
			statusCode = consts.StatusUnauthorized
		}
		c.JSON(statusCode, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	c.JSON(consts.StatusOK, resp)
}

// ... 其他处理方法

package handler

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"

	"TikTokMall/app/user/biz/service"
	"TikTokMall/app/user/kitex_gen/user"
	"TikTokMall/app/user/pkg/errno"
)

// UserHandler handles HTTP requests for user service
type UserHandler struct {
	svc *service.UserService
}

// NewUserHandler creates a new UserHandler instance
func NewUserHandler() *UserHandler {
	return &UserHandler{
		svc: service.NewUserService(),
	}
}

// Response represents a common HTTP response structure
type Response struct {
	Code    int32       `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// Register handles user registration request
func (h *UserHandler) Register(ctx context.Context, c *app.RequestContext) {
	var req user.RegisterReq
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(consts.StatusBadRequest, Response{
			Code:    errno.ParamErr.Code,
			Message: err.Error(),
		})
		return
	}

	// Validate passwords match
	if req.Password != req.ConfirmPassword {
		c.JSON(consts.StatusBadRequest, Response{
			Code:    errno.ParamErr.Code,
			Message: "passwords do not match",
		})
		return
	}

	// Call service layer to handle registration
	resp, err := service.NewRegisterService(ctx).Run(&req)
	if err != nil {
		if e, ok := err.(*errno.ErrNo); ok {
			c.JSON(consts.StatusOK, Response{
				Code:    e.Code,
				Message: e.Message,
			})
			return
		}
		c.JSON(consts.StatusInternalServerError, Response{
			Code:    errno.ServiceErr.Code,
			Message: err.Error(),
		})
		return
	}

	c.JSON(consts.StatusOK, Response{
		Code:    errno.Success.Code,
		Message: "success",
		Data: map[string]interface{}{
			"user_id": resp.UserId,
		},
	})
}

// Login handles user login request
func (h *UserHandler) Login(ctx context.Context, c *app.RequestContext) {
	var req user.LoginReq
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(consts.StatusBadRequest, Response{
			Code:    errno.ParamErr.Code,
			Message: err.Error(),
		})
		return
	}

	// Call service layer to handle login
	resp, err := service.NewLoginService(ctx).Run(&req)
	if err != nil {
		if e, ok := err.(*errno.ErrNo); ok {
			c.JSON(consts.StatusOK, Response{
				Code:    e.Code,
				Message: e.Message,
			})
			return
		}
		c.JSON(consts.StatusInternalServerError, Response{
			Code:    errno.ServiceErr.Code,
			Message: err.Error(),
		})
		return
	}

	c.JSON(consts.StatusOK, Response{
		Code:    errno.Success.Code,
		Message: "success",
		Data: map[string]interface{}{
			"user_id": resp.UserId,
		},
	})
}
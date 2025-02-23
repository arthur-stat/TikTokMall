package handler

import (
	"context"
	"strings"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"

	"TikTokMall/app/user/biz/service"
	"TikTokMall/app/user/kitex_gen/user"
)

// UserHandler 用户服务处理器
type UserHandler struct {
	svc *service.UserService
}

// NewUserHandler 创建用户服务处理器
func NewUserHandler() *UserHandler {
	return &UserHandler{
		svc: service.NewUserService(),
	}
}

// Register 处理用户注册请求
func (h *UserHandler) Register(ctx context.Context, c *app.RequestContext) {
	var req user.RegisterReq
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(consts.StatusBadRequest, &user.RegisterResp{
			Base: &user.BaseResp{
				Code:    consts.StatusBadRequest,
				Message: err.Error(),
			},
		})
		return
	}

	// 调用服务层处理注册
	userID, token, err := h.svc.Register(ctx, req.Username, req.Password, req.Email, req.Phone)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, &user.RegisterResp{
			Base: &user.BaseResp{
				Code:    consts.StatusInternalServerError,
				Message: err.Error(),
			},
		})
		return
	}

	c.JSON(consts.StatusOK, &user.RegisterResp{
		Base: &user.BaseResp{
			Code:    consts.StatusOK,
			Message: "success",
		},
		UserId: userID,
		Token:  token,
	})
}

// Login 处理用户登录请求
func (h *UserHandler) Login(ctx context.Context, c *app.RequestContext) {
	var req user.LoginReq
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(consts.StatusBadRequest, &user.LoginResp{
			Base: &user.BaseResp{
				Code:    consts.StatusBadRequest,
				Message: err.Error(),
			},
		})
		return
	}

	// 调用服务层处理登录
	token, refreshToken, err := h.svc.Login(ctx, req.Username, req.Password)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, &user.LoginResp{
			Base: &user.BaseResp{
				Code:    consts.StatusInternalServerError,
				Message: err.Error(),
			},
		})
		return
	}

	c.JSON(consts.StatusOK, &user.LoginResp{
		Base: &user.BaseResp{
			Code:    consts.StatusOK,
			Message: "success",
		},
		Token:        token,
		RefreshToken: refreshToken,
	})
}

// Logout 处理用户登出请求
func (h *UserHandler) Logout(ctx context.Context, c *app.RequestContext) {
	var req user.LogoutReq
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(consts.StatusBadRequest, &user.LogoutResp{
			Base: &user.BaseResp{
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
		c.JSON(consts.StatusInternalServerError, &user.LogoutResp{
			Base: &user.BaseResp{
				Code:    consts.StatusInternalServerError,
				Message: err.Error(),
			},
		})
		return
	}

	c.JSON(consts.StatusOK, &user.LogoutResp{
		Base: &user.BaseResp{
			Code:    consts.StatusOK,
			Message: "success",
		},
	})
}

// Delete 处理用户删除请求
func (h *UserHandler) Delete(ctx context.Context, c *app.RequestContext) {
	var req user.DeleteReq
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(consts.StatusBadRequest, &user.DeleteResp{
			Base: &user.BaseResp{
				Code:    consts.StatusBadRequest,
				Message: err.Error(),
			},
		})
		return
	}

	// 从Authorization头中获取token
	token := strings.TrimPrefix(req.Token, "Bearer ")

	// 调用服务层处理删除
	if err := h.svc.Delete(ctx, token); err != nil {
		c.JSON(consts.StatusInternalServerError, &user.DeleteResp{
			Base: &user.BaseResp{
				Code:    consts.StatusInternalServerError,
				Message: err.Error(),
			},
		})
		return
	}

	c.JSON(consts.StatusOK, &user.DeleteResp{
		Base: &user.BaseResp{
			Code:    consts.StatusOK,
			Message: "success",
		},
	})
}

// Update 处理用户更新请求
func (h *UserHandler) Update(ctx context.Context, c *app.RequestContext) {
	var req user.UpdateReq
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(consts.StatusBadRequest, &user.UpdateResp{
			Base: &user.BaseResp{
				Code:    consts.StatusBadRequest,
				Message: err.Error(),
			},
		})
		return
	}

	// 从Authorization头中获取token
	token := strings.TrimPrefix(req.Token, "Bearer ")

	// 调用服务层处理更新
	if err := h.svc.Update(ctx, token, req.NewUsername, req.NewEmail, req.NewPhone); err != nil {
		c.JSON(consts.StatusInternalServerError, &user.UpdateResp{
			Base: &user.BaseResp{
				Code:    consts.StatusInternalServerError,
				Message: err.Error(),
			},
		})
		return
	}

	c.JSON(consts.StatusOK, &user.UpdateResp{
		Base: &user.BaseResp{
			Code:    consts.StatusOK,
			Message: "success",
		},
	})
}

// Info 处理用户信息请求
func (h *UserHandler) Info(ctx context.Context, c *app.RequestContext) {
	var req user.InfoReq
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(consts.StatusBadRequest, &user.InfoResp{
			Base: &user.BaseResp{
				Code:    consts.StatusBadRequest,
				Message: err.Error(),
			},
		})
		return
	}

	// 从Authorization头中获取token
	token := strings.TrimPrefix(req.Token, "Bearer ")

	// 调用服务层处理信息获取
	userID, username, email, phone, err := h.svc.Info(ctx, token)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, &user.InfoResp{
			Base: &user.BaseResp{
				Code:    consts.StatusInternalServerError,
				Message: err.Error(),
			},
		})
		return
	}

	c.JSON(consts.StatusOK, &user.InfoResp{
		Base: &user.BaseResp{
			Code:    consts.StatusOK,
			Message: "success",
		},
		UserId:   userID,
		Username: username,
		Email:    email,
		Phone:    phone,
	})
}

package handler

import (
	"TikTokMall/app/user/biz/service"
	"context"
)

// UserHandler 用户服务处理器
type UserHandler struct {
	svc *service.UserService
}

// NewAuthHandler 创建用户服务处理器
func NewAuthHandler() *UserHandler {
	return &UserHandler{
		svc: service.NewUserService(),
	}
}

// Delete 处理用户删除请求
func (h *UserHandler) Delete(ctx context.Context, token string) error {

}

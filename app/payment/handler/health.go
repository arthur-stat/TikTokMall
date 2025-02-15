package handler

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"net/http"
)

// HealthHandler 用于返回健康检查结果
func HealthHandler(ctx context.Context, c *app.RequestContext) {
	c.JSON(http.StatusOK, map[string]string{
		"status": "healthy",
	})
}

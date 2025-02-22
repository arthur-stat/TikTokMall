package handler

import (
	"fmt"
	"github.com/cloudwego/hertz/pkg/app"
)

type BaseResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func Success(ctx *app.RequestContext, data any) {
	ctx.JSON(200, BaseResponse{
		Code:    200,
		Message: "success",
		Data:    data,
	})
}

func Error(ctx *app.RequestContext, httpCode int, bizCode int, msg string) {
	ctx.JSON(httpCode, BaseResponse{
		Code:    bizCode,
		Message: msg,
	})
}

// handleError 统一处理错误响应
func handleError(ctx *app.RequestContext, statusCode int, message string, err error) {
	ctx.JSON(statusCode, map[string]string{
		"error": fmt.Sprintf("%s: %v", message, err),
	})
}

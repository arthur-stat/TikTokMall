package middleware

import (
	"context"
	"time"

	"github.com/cloudwego/kitex/pkg/endpoint"
	"github.com/cloudwego/kitex/pkg/klog"
)

// Logging 日志中间件
func Logging() endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, req, resp interface{}) (err error) {
			start := time.Now()
			err = next(ctx, req, resp)
			duration := time.Since(start)

			if err != nil {
				klog.CtxErrorf(ctx, "[Order] method: %s, duration: %v, err: %v", ctx.Value("method"), duration, err)
			} else {
				klog.CtxInfof(ctx, "[Order] method: %s, duration: %v", ctx.Value("method"), duration)
			}

			return err
		}
	}
}

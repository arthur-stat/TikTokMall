package opentracing

import (
	"context"

	"github.com/opentracing/opentracing-go"
)

// StartSpanFromContext 从上下文创建新的span
func StartSpanFromContext(ctx context.Context, operationName string) (opentracing.Span, context.Context) {
	return opentracing.StartSpanFromContext(ctx, operationName)
}

// SpanFromContext 从上下文获取span
func SpanFromContext(ctx context.Context) opentracing.Span {
	return opentracing.SpanFromContext(ctx)
}

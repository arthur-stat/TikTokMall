package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"TikTokMall/app/checkout/kitex_gen/checkout"
)

// TestSimpleResponses 是一个简单的测试，不依赖于特定类型
func TestSimpleResponses(t *testing.T) {
	// 创建一个极简模拟服务
	svc := &mockCheckoutService{}

	req := &checkout.CheckoutReq{
		UserId: 1,
	}

	resp, err := svc.Run(context.Background(), req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "test-order-id", resp.OrderId)
}

// TestEmptyOrderID 测试错误情况
func TestEmptyOrderID(t *testing.T) {
	// 创建模拟服务
	svc := &mockCheckoutService{}

	req := &checkout.CheckoutReq{
		// 空请求
	}

	_, err := svc.Run(context.Background(), req)
	assert.Error(t, err)
}

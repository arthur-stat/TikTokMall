package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"TikTokMall/app/checkout/kitex_gen/checkout"
)

// TestCheckoutWithoutDB 测试无需数据库的功能
func TestCheckoutWithoutDB(t *testing.T) {
	// 使用我们已经实现的模拟服务
	svc := &mockCheckoutService{}

	// 创建测试请求
	req := &checkout.CheckoutReq{
		UserId:    123,
		Firstname: "Test",
		Lastname:  "User",
		Email:     "test@example.com",
	}

	// 调用服务
	resp, err := svc.Run(context.Background(), req)

	// 验证响应
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "test-order-id", resp.OrderId)
	assert.Equal(t, "test-transaction-id", resp.TransactionId)
}

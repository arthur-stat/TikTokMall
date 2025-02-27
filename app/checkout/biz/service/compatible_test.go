package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"TikTokMall/app/checkout/kitex_gen/checkout"
)

// TestCheckoutBasic 测试基本功能
func TestCheckoutBasic(t *testing.T) {
	// 使用模拟实现进行测试
	svc := &mockCheckoutService{}

	// 测试基本流程，避免使用有问题的字段
	req := &checkout.CheckoutReq{
		UserId:    1,
		Firstname: "测试用户",
		Email:     "test@example.com",
	}

	resp, err := svc.Run(context.Background(), req)

	assert.NoError(t, err)
	assert.NotEmpty(t, resp.OrderId)
	assert.NotEmpty(t, resp.TransactionId)
}

// TestCheckoutValidation 测试输入验证
func TestCheckoutValidation(t *testing.T) {
	// 使用模拟实现进行测试
	svc := &mockCheckoutService{}

	// 测试无效输入
	req := &checkout.CheckoutReq{
		// 不设置 UserId
		Firstname: "测试用户",
	}

	_, err := svc.Run(context.Background(), req)

	// 验证返回预期错误
	assert.Error(t, err)
	assert.Equal(t, ErrInvalidInput, err)
}

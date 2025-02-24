package service

import (
	"TikTokMall/app/payment/kitex_gen/payment"
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAlipayChargeService_Run(t *testing.T) {
	// 创建上下文
	ctx := context.Background()

	// 创建AlipayChargeService实例
	service := NewAlipayChargeService(ctx)

	// 创建测试请求
	req := &payment.AlipayChargeReq{
		OrderId:   123456,
		UserId:    123,
		Amount:    100.0,
		ReturnUrl: "https://example.com/return",
	}

	// 执行Run方法
	resp, err := service.Run(req)

	// 验证没有错误
	assert.NoError(t, err)

	// 验证返回的支付URL不为空
	assert.NotEmpty(t, resp.PayUrl)
}

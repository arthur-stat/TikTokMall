package service

import (
	"context"

	"TikTokMall/app/checkout/kitex_gen/checkout"
)

// 模拟结账服务的实现，专门用于测试
type mockCheckoutService struct{}

// Run 处理结账请求并返回响应
func (s *mockCheckoutService) Run(ctx context.Context, req *checkout.CheckoutReq) (*checkout.CheckoutResp, error) {
	// 简单的验证
	if req.UserId == 0 {
		return nil, ErrInvalidInput
	}

	// 这里我们返回一个固定的成功响应，便于测试
	return &checkout.CheckoutResp{
		OrderId:       "test-order-id",
		TransactionId: "test-transaction-id",
	}, nil
}

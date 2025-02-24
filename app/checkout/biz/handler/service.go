package handler

import (
	"TikTokMall/app/checkout/kitex_gen/checkout"
	"context"
)

// CheckoutServiceImpl implements the checkout service interface
type CheckoutServiceImpl struct{}

func NewCheckoutServiceImpl() *CheckoutServiceImpl {
	return &CheckoutServiceImpl{}
}

// Checkout implements the checkout service interface
func (s *CheckoutServiceImpl) Checkout(ctx context.Context, req *checkout.CheckoutReq) (*checkout.CheckoutResp, error) {
	// TODO: 实现结账逻辑
	return &checkout.CheckoutResp{}, nil
}

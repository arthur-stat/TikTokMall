package main

import (
	"context"

	"TikTokMall/app/checkout/biz/service"
	"TikTokMall/app/checkout/kitex_gen/checkout"
)

// CheckoutServiceImpl implements the last service interface defined in the IDL.
type CheckoutServiceImpl struct {
	svc service.CheckoutService
}

func NewCheckoutServiceImpl() *CheckoutServiceImpl {
	return &CheckoutServiceImpl{
		svc: service.NewCheckoutService(),
	}
}

// Checkout implements the CheckoutServiceImpl interface.
func (s *CheckoutServiceImpl) Checkout(ctx context.Context, req *checkout.CheckoutReq) (resp *checkout.CheckoutResp, err error) {
	return s.svc.Run(ctx, req)
}

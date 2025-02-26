package service

import (
	payment "TikTokMall/app/payment/kitex_gen/payment"
	"context"
)

type AlipayRefundService struct {
	ctx context.Context
} // NewAlipayRefundService new AlipayRefundService
func NewAlipayRefundService(ctx context.Context) *AlipayRefundService {
	return &AlipayRefundService{ctx: ctx}
}

// Run create note info
func (s *AlipayRefundService) Run(req *payment.AlipayRefundReq) (resp *payment.AlipayRefundResp, err error) {

	return
}

package service

import (
	"context"
	"testing"
	payment "TikTokMall/app/payment/kitex_gen/payment"
)

func TestAlipayRefund_Run(t *testing.T) {
	ctx := context.Background()
	s := NewAlipayRefundService(ctx)
	// init req and assert value

	req := &payment.AlipayRefundReq{}
	resp, err := s.Run(req)
	t.Logf("err: %v", err)
	t.Logf("resp: %v", resp)

	// todo: edit your unit test

}

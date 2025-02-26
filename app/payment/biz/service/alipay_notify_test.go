package service

import (
	"context"
	"testing"
	payment "TikTokMall/app/payment/kitex_gen/payment"
)

func TestAlipayNotify_Run(t *testing.T) {
	ctx := context.Background()
	s := NewAlipayNotifyService(ctx)
	// init req and assert value

	req := &payment.AlipayNotifyReq{}
	resp, err := s.Run(req)
	t.Logf("err: %v", err)
	t.Logf("resp: %v", resp)

	// todo: edit your unit test

}

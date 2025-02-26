package payment

import (
	payment "TikTokMall/rpc_gen/kitex_gen/payment"
	"context"
	"github.com/cloudwego/kitex/client/callopt"
	"github.com/cloudwego/kitex/pkg/klog"
)

func Charge(ctx context.Context, req *payment.ChargeReq, callOptions ...callopt.Option) (resp *payment.ChargeResp, err error) {
	resp, err = defaultClient.Charge(ctx, req, callOptions...)
	if err != nil {
		klog.CtxErrorf(ctx, "Charge call failed,err =%+v", err)
		return nil, err
	}
	return resp, nil
}

func Refund(ctx context.Context, req *payment.RefundReq, callOptions ...callopt.Option) (resp *payment.RefundResp, err error) {
	resp, err = defaultClient.Refund(ctx, req, callOptions...)
	if err != nil {
		klog.CtxErrorf(ctx, "Refund call failed,err =%+v", err)
		return nil, err
	}
	return resp, nil
}

func AlipayCharge(ctx context.Context, req *payment.AlipayChargeReq, callOptions ...callopt.Option) (resp *payment.AlipayChargeResp, err error) {
	resp, err = defaultClient.AlipayCharge(ctx, req, callOptions...)
	if err != nil {
		klog.CtxErrorf(ctx, "AlipayCharge call failed,err =%+v", err)
		return nil, err
	}
	return resp, nil
}

func AlipayRefund(ctx context.Context, req *payment.AlipayRefundReq, callOptions ...callopt.Option) (resp *payment.AlipayRefundResp, err error) {
	resp, err = defaultClient.AlipayRefund(ctx, req, callOptions...)
	if err != nil {
		klog.CtxErrorf(ctx, "AlipayRefund call failed,err =%+v", err)
		return nil, err
	}
	return resp, nil
}

func AlipayNotify(ctx context.Context, req *payment.AlipayNotifyReq, callOptions ...callopt.Option) (resp *payment.AlipayNotifyResp, err error) {
	resp, err = defaultClient.AlipayNotify(ctx, req, callOptions...)
	if err != nil {
		klog.CtxErrorf(ctx, "AlipayNotify call failed,err =%+v", err)
		return nil, err
	}
	return resp, nil
}

package payment

import (
	payment "TikTokMall/rpc_gen/kitex_gen/payment"
	"context"

	"TikTokMall/rpc_gen/kitex_gen/payment/paymentservice"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/client/callopt"
)

type RPCClient interface {
	KitexClient() paymentservice.Client
	Service() string
	Charge(ctx context.Context, Req *payment.ChargeReq, callOptions ...callopt.Option) (r *payment.ChargeResp, err error)
	Refund(ctx context.Context, Req *payment.RefundReq, callOptions ...callopt.Option) (r *payment.RefundResp, err error)
	AlipayCharge(ctx context.Context, Req *payment.AlipayChargeReq, callOptions ...callopt.Option) (r *payment.AlipayChargeResp, err error)
}

func NewRPCClient(dstService string, opts ...client.Option) (RPCClient, error) {
	kitexClient, err := paymentservice.NewClient(dstService, opts...)
	if err != nil {
		return nil, err
	}
	cli := &clientImpl{
		service:     dstService,
		kitexClient: kitexClient,
	}

	return cli, nil
}

type clientImpl struct {
	service     string
	kitexClient paymentservice.Client
}

func (c *clientImpl) Service() string {
	return c.service
}

func (c *clientImpl) KitexClient() paymentservice.Client {
	return c.kitexClient
}

func (c *clientImpl) Charge(ctx context.Context, Req *payment.ChargeReq, callOptions ...callopt.Option) (r *payment.ChargeResp, err error) {
	return c.kitexClient.Charge(ctx, Req, callOptions...)
}

func (c *clientImpl) Refund(ctx context.Context, Req *payment.RefundReq, callOptions ...callopt.Option) (r *payment.RefundResp, err error) {
	return c.kitexClient.Refund(ctx, Req, callOptions...)
}

func (c *clientImpl) AlipayCharge(ctx context.Context, Req *payment.AlipayChargeReq, callOptions ...callopt.Option) (r *payment.AlipayChargeResp, err error) {
	return c.kitexClient.AlipayCharge(ctx, Req, callOptions...)
}

package service

import (
	"context"

	"TikTokMall/app/order/kitex_gen/order"
)

// OrderService 定义订单服务接口
type OrderService interface {
	PlaceOrder(ctx context.Context, req *order.PlaceOrderReq) (*order.PlaceOrderResp, error)
	ListOrder(ctx context.Context, req *order.ListOrderReq) (*order.ListOrderResp, error)
	MarkOrderPaid(ctx context.Context, req *order.MarkOrderPaidReq) (*order.MarkOrderPaidResp, error)
}

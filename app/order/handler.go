package main

import (
	"context"

	"TikTokMall/app/order/biz/handler"
	"TikTokMall/app/order/biz/service"
	"TikTokMall/app/order/kitex_gen/order"
)

// OrderServiceImpl implements the last service interface defined in the IDL.
type OrderServiceImpl struct{}

// PlaceOrder implements the OrderServiceImpl interface.
func (s *OrderServiceImpl) PlaceOrder(ctx context.Context, req *order.PlaceOrderReq) (resp *order.PlaceOrderResp, err error) {
	svc := service.NewOrderService()
	return svc.PlaceOrder(ctx, req)
}

// ListOrder implements the OrderServiceImpl interface.
func (s *OrderServiceImpl) ListOrder(ctx context.Context, req *order.ListOrderReq) (resp *order.ListOrderResp, err error) {
	svc := service.NewOrderService()
	return svc.ListOrder(ctx, req)
}

// MarkOrderPaid implements the OrderServiceImpl interface.
func (s *OrderServiceImpl) MarkOrderPaid(ctx context.Context, req *order.MarkOrderPaidReq) (resp *order.MarkOrderPaidResp, err error) {
	svc := service.NewOrderService()
	return svc.MarkOrderPaid(ctx, req)
}

// HealthCheck implements the OrderServiceImpl interface.
func (s *OrderServiceImpl) HealthCheck(ctx context.Context) error {
	return handler.HealthCheck(ctx)
}

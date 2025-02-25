package main

import (
	"context"

	"TikTokMall/app/order/biz/dal/mysql"
	"TikTokMall/app/order/biz/dal/redis"
	"TikTokMall/app/order/biz/service"
	"TikTokMall/app/order/kitex_gen/order"
)

// OrderServiceImpl implements the last service interface defined in the IDL.
type OrderServiceImpl struct{}

// PlaceOrder implements the OrderServiceImpl interface.
func (s *OrderServiceImpl) PlaceOrder(ctx context.Context, req *order.PlaceOrderReq) (resp *order.PlaceOrderResp, err error) {
	repo := mysql.NewOrderMySQLRepository()
	svc := service.NewOrderService(repo)
	return svc.PlaceOrder(ctx, req)
}

// ListOrder implements the OrderServiceImpl interface.
func (s *OrderServiceImpl) ListOrder(ctx context.Context, req *order.ListOrderReq) (resp *order.ListOrderResp, err error) {
	repo := mysql.NewOrderMySQLRepository()
	svc := service.NewOrderService(repo)
	return svc.ListOrder(ctx, req)
}

// MarkOrderPaid implements the OrderServiceImpl interface.
func (s *OrderServiceImpl) MarkOrderPaid(ctx context.Context, req *order.MarkOrderPaidReq) (resp *order.MarkOrderPaidResp, err error) {
	repo := mysql.NewOrderMySQLRepository()
	svc := service.NewOrderService(repo)
	return svc.MarkOrderPaid(ctx, req)
}

// HealthCheck implements the OrderServiceImpl interface.
func (s *OrderServiceImpl) HealthCheck(ctx context.Context) error {
	if err := mysql.DB.WithContext(ctx).Raw("SELECT 1").Error; err != nil {
		return err
	}

	if err := redis.RDB.Ping(ctx).Err(); err != nil {
		return err
	}

	return nil
}

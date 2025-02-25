package mysql

import (
	"context"
)

type OrderMySQLRepository struct{}

func NewOrderMySQLRepository() *OrderMySQLRepository {
	return &OrderMySQLRepository{}
}

func (r *OrderMySQLRepository) CreateOrder(ctx context.Context, order *Order, items []*OrderItem) error {
	return CreateOrder(ctx, order, items)
}

func (r *OrderMySQLRepository) ListOrdersByUserID(ctx context.Context, userID uint32) ([]*Order, error) {
	return ListOrdersByUserID(ctx, userID)
}

func (r *OrderMySQLRepository) UpdateOrderStatus(ctx context.Context, orderID int64, status int8) error {
	return UpdateOrderStatus(ctx, orderID, status)
}

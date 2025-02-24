package mysql

import (
	"context"

	"gorm.io/gorm"
)

// CreateOrder 创建订单
func CreateOrder(ctx context.Context, order *Order, items []*OrderItem) error {
	return DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 创建订单
		if err := tx.Create(order).Error; err != nil {
			return err
		}

		// 创建订单项
		for _, item := range items {
			item.OrderID = order.ID
		}
		return tx.Create(items).Error
	})
}

// GetOrderByID 根据ID获取订单
func GetOrderByID(ctx context.Context, orderID int64) (*Order, error) {
	var order Order
	err := DB.WithContext(ctx).First(&order, orderID).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

// ListOrdersByUserID 获取用户订单列表
func ListOrdersByUserID(ctx context.Context, userID uint32) ([]*Order, error) {
	var orders []*Order
	err := DB.WithContext(ctx).Where("user_id = ?", userID).Find(&orders).Error
	return orders, err
}

// UpdateOrderStatus 更新订单状态
func UpdateOrderStatus(ctx context.Context, orderID int64, status int8) error {
	return DB.WithContext(ctx).Model(&Order{}).Where("id = ?", orderID).Update("status", status).Error
}

package mysql

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)

// Cart 相关操作
func AddToCart(ctx context.Context, cart *Cart) error {
	// 检查是否已存在相同商品
	var existingCart Cart
	err := DB.WithContext(ctx).Where("user_id = ? AND product_id = ?", cart.UserID, cart.ProductID).First(&existingCart).Error
	if err == nil {
		// 已存在则更新数量
		existingCart.Quantity += cart.Quantity
		return DB.WithContext(ctx).Save(&existingCart).Error
	}
	// 不存在则创建新记录
	return DB.WithContext(ctx).Create(cart).Error
}

func UpdateCartItem(ctx context.Context, userID int64, cartID int64, quantity int) error {
	result := DB.WithContext(ctx).Model(&Cart{}).
		Where("id = ? AND user_id = ?", cartID, userID).
		Update("quantity", quantity)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("购物车项不存在")
	}
	return nil
}

func DeleteCartItem(ctx context.Context, userID int64, cartID int64) error {
	result := DB.WithContext(ctx).Where("id = ? AND user_id = ?", cartID, userID).Delete(&Cart{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("购物车项不存在")
	}
	return nil
}

func GetUserCart(ctx context.Context, userID int64) ([]*Cart, error) {
	var carts []*Cart
	err := DB.WithContext(ctx).Where("user_id = ?", userID).Find(&carts).Error
	return carts, err
}

// Order 相关操作
func CreateOrder(ctx context.Context, order *Order, items []*OrderItem) error {
	return DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 创建订单
		if err := tx.Create(order).Error; err != nil {
			return err
		}

		// 创建订单项
		for _, item := range items {
			item.OrderID = order.ID
			if err := tx.Create(item).Error; err != nil {
				return err
			}
		}

		// 清除已购买的购物车项
		if err := tx.Where("user_id = ? AND selected = ?", order.UserID, true).Delete(&Cart{}).Error; err != nil {
			return err
		}

		return nil
	})
}

func GetOrderByID(ctx context.Context, orderID int64) (*Order, error) {
	var order Order
	err := DB.WithContext(ctx).First(&order, orderID).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func GetOrderByOrderNo(ctx context.Context, orderNo string) (*Order, error) {
	var order Order
	err := DB.WithContext(ctx).Where("order_no = ?", orderNo).First(&order).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func UpdateOrderStatus(ctx context.Context, orderID int64, status int8) error {
	result := DB.WithContext(ctx).Model(&Order{}).
		Where("id = ?", orderID).
		Update("status", status)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("订单不存在")
	}
	return nil
}

func GetUserOrders(ctx context.Context, userID int64, status *int8) ([]*Order, error) {
	query := DB.WithContext(ctx).Where("user_id = ?", userID)
	if status != nil {
		query = query.Where("status = ?", *status)
	}

	var orders []*Order
	err := query.Order("created_at DESC").Find(&orders).Error
	return orders, err
}

// OrderItem 相关操作
func GetOrderItems(ctx context.Context, orderID int64) ([]*OrderItem, error) {
	var items []*OrderItem
	err := DB.WithContext(ctx).Where("order_id = ?", orderID).Find(&items).Error
	return items, err
}

// 获取订单详情（包含订单项）
type OrderDetail struct {
	Order *Order
	Items []*OrderItem
}

func GetOrderDetail(ctx context.Context, orderID int64) (*OrderDetail, error) {
	order, err := GetOrderByID(ctx, orderID)
	if err != nil {
		return nil, err
	}

	items, err := GetOrderItems(ctx, orderID)
	if err != nil {
		return nil, err
	}

	return &OrderDetail{
		Order: order,
		Items: items,
	}, nil
}

package mysql

import (
	"context"
	"fmt"

	"TikTokMall/app/cart/biz/model"
)

// 这个函数在cart.go中已经定义，所以这里不需要再定义
// func DeleteCartItem(ctx context.Context, userID uint32, productID uint32) error {
//     // 实现...
// }

// 如果需要，可以添加其他不重复的数据库操作函数

// 重命名避免和 cart.go 中的函数冲突
func CreateNewCartItem(ctx context.Context, item *model.CartItem) error {
	return DB.WithContext(ctx).Create(item).Error
}

// UpdateCartItemQuantity 更新购物车项数量
func UpdateCartItemQuantity(ctx context.Context, userID uint32, productID uint32, quantity uint32) error {
	result := DB.WithContext(ctx).Model(&model.CartItem{}).
		Where("user_id = ? AND product_id = ?", userID, productID).
		Update("quantity", quantity)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("购物车项不存在")
	}
	return nil
}

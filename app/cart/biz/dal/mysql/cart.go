package mysql

import (
    "context"
    "TikTokMall/app/cart/biz/model"
)

// AddCartItem 添加购物车商品
func AddCartItem(ctx context.Context, item *model.CartItem) error {
    return DB.WithContext(ctx).Create(item).Error
}

// GetUserCart 获取用户的购物车商品
func GetUserCart(ctx context.Context, userID uint32) ([]*model.CartItem, error) {
    var items []*model.CartItem
    err := DB.WithContext(ctx).Where("user_id = ?", userID).Find(&items).Error
    return items, err
}

// EmptyCart 清空用户的购物车
func EmptyCart(ctx context.Context, userID uint32) error {
    return DB.WithContext(ctx).Where("user_id = ?", userID).Delete(&model.CartItem{}).Error
}

// CreateCartItem 创建购物车商品
func CreateCartItem(ctx context.Context, item *model.CartItem) error {
    return DB.WithContext(ctx).Create(item).Error
}

// GetCartItems 获取用户的购物车商品列表
func GetCartItems(ctx context.Context, userID uint32) ([]*model.CartItem, error) {
    var items []*model.CartItem
    err := DB.WithContext(ctx).Where("user_id = ?", userID).Find(&items).Error
    return items, err
}

// UpdateCartItem 更新购物车商品
func UpdateCartItem(ctx context.Context, item *model.CartItem) error {
    return DB.WithContext(ctx).Save(item).Error
}

// DeleteCartItem 删除购物车商品
func DeleteCartItem(ctx context.Context, userID uint32, itemID uint32) error {
    return DB.WithContext(ctx).Where("user_id = ? AND id = ?", userID, itemID).Delete(&model.CartItem{}).Error
}

// BatchCreateCartItems 批量创建购物车商品
func BatchCreateCartItems(ctx context.Context, items []*model.CartItem) error {
    return DB.WithContext(ctx).Create(items).Error
}

package service

import (
	"context"

	"TikTokMall/app/cart/biz/dal/mysql"
	"TikTokMall/app/cart/biz/model"
)

// 定义一个默认的仓库实现
type defaultCartRepository struct{}

func NewCartRepository() cartRepository {
	return &defaultCartRepository{}
}

// 实现cartRepository接口的各个方法
func (r *defaultCartRepository) AddItem(ctx context.Context, userID uint32, item *model.CartItem) error {
	// 真实实现，调用mysql保存数据
	return mysql.DB.Create(item).Error
}

func (r *defaultCartRepository) GetItems(ctx context.Context, userID uint32) ([]*model.CartItem, error) {
	var items []*model.CartItem
	err := mysql.DB.Where("user_id = ?", userID).Find(&items).Error
	return items, err
}

func (r *defaultCartRepository) RemoveItem(ctx context.Context, userID uint32, productID uint32) error {
	return mysql.DB.Where("user_id = ? AND product_id = ?", userID, productID).Delete(&model.CartItem{}).Error
}

func (r *defaultCartRepository) UpdateItemQuantity(ctx context.Context, userID uint32, productID uint32, quantity uint32) error {
	return mysql.DB.Model(&model.CartItem{}).
		Where("user_id = ? AND product_id = ?", userID, productID).
		Update("quantity", quantity).Error
}

func (r *defaultCartRepository) EmptyCart(ctx context.Context, userID uint32) error {
	return mysql.DB.Where("user_id = ?", userID).Delete(&model.CartItem{}).Error
}

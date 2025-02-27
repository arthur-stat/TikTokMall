package service

import (
	"context"

	"TikTokMall/app/cart/biz/model"
)

// 定义一个默认的仓库实现
type defaultCartRepository struct{}

func NewCartRepository() cartRepository {
	return &defaultCartRepository{}
}

// 实现cartRepository接口的各个方法
func (r *defaultCartRepository) AddItem(ctx context.Context, userID uint32, item *model.CartItem) error {
	// 真实实现会调用mysql或redis存储
	return nil
}

func (r *defaultCartRepository) GetItems(ctx context.Context, userID uint32) ([]*model.CartItem, error) {
	return []*model.CartItem{}, nil
}

func (r *defaultCartRepository) RemoveItem(ctx context.Context, userID uint32, productID uint32) error {
	return nil
}

func (r *defaultCartRepository) UpdateItemQuantity(ctx context.Context, userID uint32, productID uint32, quantity uint32) error {
	return nil
}

func (r *defaultCartRepository) EmptyCart(ctx context.Context, userID uint32) error {
	return nil
}

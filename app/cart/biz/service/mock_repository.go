package service

import (
	"TikTokMall/app/cart/biz/model"
	"context"
)

// MockCartRepository 创建一个内存仓库实现用于测试
type MockCartRepository struct {
	items map[uint32][]*model.CartItem
}

func NewMockCartRepository() cartRepository {
	return &MockCartRepository{
		items: make(map[uint32][]*model.CartItem),
	}
}

func (r *MockCartRepository) AddItem(ctx context.Context, userID uint32, item *model.CartItem) error {
	if r.items[userID] == nil {
		r.items[userID] = []*model.CartItem{}
	}

	// 检查是否已存在相同商品
	for i, existingItem := range r.items[userID] {
		if existingItem.ProductID == item.ProductID {
			// 更新数量
			r.items[userID][i].Quantity += item.Quantity
			return nil
		}
	}

	// 添加新商品
	r.items[userID] = append(r.items[userID], item)
	return nil
}

func (r *MockCartRepository) GetItems(ctx context.Context, userID uint32) ([]*model.CartItem, error) {
	return r.items[userID], nil
}

func (r *MockCartRepository) RemoveItem(ctx context.Context, userID uint32, productID uint32) error {
	if r.items[userID] == nil {
		return nil
	}

	newItems := []*model.CartItem{}
	for _, item := range r.items[userID] {
		if item.ProductID != productID {
			newItems = append(newItems, item)
		}
	}

	r.items[userID] = newItems
	return nil
}

func (r *MockCartRepository) UpdateItemQuantity(ctx context.Context, userID uint32, productID uint32, quantity uint32) error {
	if r.items[userID] == nil {
		return nil
	}

	for i, item := range r.items[userID] {
		if item.ProductID == productID {
			r.items[userID][i].Quantity = quantity
			return nil
		}
	}

	return nil
}

func (r *MockCartRepository) EmptyCart(ctx context.Context, userID uint32) error {
	r.items[userID] = []*model.CartItem{}
	return nil
}

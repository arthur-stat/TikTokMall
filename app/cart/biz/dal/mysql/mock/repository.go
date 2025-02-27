package mock

import (
	"context"

	"github.com/stretchr/testify/mock"

	"TikTokMall/app/cart/biz/model"
)

// MockCartRepository 是 CartRepository 的 mock 实现
type MockCartRepository struct {
	mock.Mock
}

func (m *MockCartRepository) AddItem(ctx context.Context, userID uint32, item *model.CartItem) error {
	args := m.Called(ctx, userID, item)
	return args.Error(0)
}

func (m *MockCartRepository) GetItems(ctx context.Context, userID uint32) ([]*model.CartItem, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.CartItem), args.Error(1)
}

func (m *MockCartRepository) RemoveItem(ctx context.Context, userID uint32, productID uint32) error {
	args := m.Called(ctx, userID, productID)
	return args.Error(0)
}

func (m *MockCartRepository) UpdateItemQuantity(ctx context.Context, userID uint32, productID uint32, quantity uint32) error {
	args := m.Called(ctx, userID, productID, quantity)
	return args.Error(0)
}

func (m *MockCartRepository) EmptyCart(ctx context.Context, userID uint32) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

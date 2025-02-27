package service

import (
	"context"
	"errors"

	"TikTokMall/app/cart/biz/model"
	"TikTokMall/app/cart/kitex_gen/cart"
)

var (
	ErrInvalidQuantity = errors.New("quantity must be greater than 0")
	ErrUserNotFound    = errors.New("user not found")
)

// CartService 定义购物车服务接口
type CartService interface {
	AddItem(ctx context.Context, req *cart.AddItemReq) (*cart.AddItemResp, error)
	GetCart(ctx context.Context, req *cart.GetCartReq) (*cart.GetCartResp, error)
	EmptyCart(ctx context.Context, req *cart.EmptyCartReq) (*cart.EmptyCartResp, error)
}

type cartService struct{}

// NewCartService 创建购物车服务实例
func NewCartService() CartService {
	repo := NewCartRepository()
	return &cartServiceImpl{
		repo: repo,
	}
}

// AddItem adds an item to the user's cart
func (s *cartService) AddItem(ctx context.Context, req *cart.AddItemReq) (*cart.AddItemResp, error) {
	if req.Item.Quantity <= 0 {
		return nil, ErrInvalidQuantity
	}

	// TODO: 实现添加购物车逻辑
	return &cart.AddItemResp{}, nil
}

// GetCart retrieves the user's cart
func (s *cartService) GetCart(ctx context.Context, req *cart.GetCartReq) (*cart.GetCartResp, error) {
	// TODO: 实现获取购物车逻辑
	return &cart.GetCartResp{
		Cart: &cart.Cart{
			UserId: req.UserId,
			Items:  []*cart.CartItem{},
		},
	}, nil
}

// EmptyCart removes all items from the user's cart
func (s *cartService) EmptyCart(ctx context.Context, req *cart.EmptyCartReq) (*cart.EmptyCartResp, error) {
	// TODO: 实现清空购物车逻辑
	return &cart.EmptyCartResp{}, nil
}

// Helper function to convert database items to proto message
func convertToProtoCart(userID uint32, items []*model.CartItem) *cart.GetCartResp {
	protoItems := make([]*cart.CartItem, 0, len(items))
	for _, item := range items {
		protoItems = append(protoItems, &cart.CartItem{
			ProductId: item.ProductID,
			Quantity:  int32(item.Quantity),
		})
	}
	return &cart.GetCartResp{
		Cart: &cart.Cart{
			UserId: userID,
			Items:  protoItems,
		},
	}
}

package service

import (
	"context"
	"fmt"

	"TikTokMall/app/cart/biz/model"
	"TikTokMall/app/cart/kitex_gen/cart"

	"github.com/cloudwego/kitex/pkg/klog"
)

type cartRepository interface {
	AddItem(ctx context.Context, userID uint32, item *model.CartItem) error
	GetItems(ctx context.Context, userID uint32) ([]*model.CartItem, error)
	RemoveItem(ctx context.Context, userID uint32, productID uint32) error
	UpdateItemQuantity(ctx context.Context, userID uint32, productID uint32, quantity uint32) error
	EmptyCart(ctx context.Context, userID uint32) error
}

type cartServiceImpl struct {
	repo cartRepository
}

// NewCartServiceWithRepo 创建购物车服务实例
func NewCartServiceWithRepo(repo cartRepository) CartService {
	return &cartServiceImpl{
		repo: repo,
	}
}

// AddItem 添加商品到购物车
func (s *cartServiceImpl) AddItem(ctx context.Context, req *cart.AddItemReq) (*cart.AddItemResp, error) {
	if req.Item.Quantity <= 0 {
		return nil, ErrInvalidQuantity
	}

	item := &model.CartItem{
		UserID:    req.UserId,
		ProductID: req.Item.ProductId,
		Quantity:  uint32(req.Item.Quantity),
		Selected:  true,
	}

	// 至少返回空响应，而不是 nil
	if s.repo == nil {
		// 避免 nil 指针
		klog.Warnf("Repository is nil, using dummy implementation")
		return &cart.AddItemResp{}, nil
	}

	if err := s.repo.AddItem(ctx, req.UserId, item); err != nil {
		return nil, fmt.Errorf("添加购物车失败: %w", err)
	}

	return &cart.AddItemResp{}, nil
}

// GetCart 获取购物车
func (s *cartServiceImpl) GetCart(ctx context.Context, req *cart.GetCartReq) (*cart.GetCartResp, error) {
	items, err := s.repo.GetItems(ctx, req.UserId)
	if err != nil {
		return nil, fmt.Errorf("获取购物车失败: %w", err)
	}

	cartItems := make([]*cart.CartItem, 0, len(items))
	for _, item := range items {
		cartItems = append(cartItems, &cart.CartItem{
			ProductId: item.ProductID,
			Quantity:  int32(item.Quantity),
		})
	}

	return &cart.GetCartResp{
		Cart: &cart.Cart{
			UserId: req.UserId,
			Items:  cartItems,
		},
	}, nil
}

// EmptyCart 清空购物车
func (s *cartServiceImpl) EmptyCart(ctx context.Context, req *cart.EmptyCartReq) (*cart.EmptyCartResp, error) {
	if err := s.repo.EmptyCart(ctx, req.UserId); err != nil {
		return nil, fmt.Errorf("清空购物车失败: %w", err)
	}

	return &cart.EmptyCartResp{}, nil
}

package handler

import (
	"context"

	"TikTokMall/app/cart/biz/service"
	"TikTokMall/app/cart/kitex_gen/cart"
)

// CartServiceImpl 实现 cart.thrift 中定义的服务接口
type CartServiceImpl struct {
	svc service.CartService
}

// NewCartServiceImpl 创建一个新的 CartServiceImpl
func NewCartServiceImpl() *CartServiceImpl {
	return &CartServiceImpl{
		svc: service.NewCartService(),
	}
}

// AddItem 实现 CartServiceImpl 接口
func (s *CartServiceImpl) AddItem(ctx context.Context, req *cart.AddItemReq) (resp *cart.AddItemResp, err error) {
	return s.svc.AddItem(ctx, req)
}

// GetCart 实现 CartServiceImpl 接口
func (s *CartServiceImpl) GetCart(ctx context.Context, req *cart.GetCartReq) (resp *cart.GetCartResp, err error) {
	return s.svc.GetCart(ctx, req)
}

// EmptyCart 实现 CartServiceImpl 接口
func (s *CartServiceImpl) EmptyCart(ctx context.Context, req *cart.EmptyCartReq) (resp *cart.EmptyCartResp, err error) {
	return s.svc.EmptyCart(ctx, req)
}

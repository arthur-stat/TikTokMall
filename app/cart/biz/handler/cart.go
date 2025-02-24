package handler

import (
	"context"

	"TikTokMall/app/cart/biz/service"
	"TikTokMall/app/cart/kitex_gen/cart"
)

// CartServiceImpl implements the last service interface defined in the IDL.
type CartServiceImpl struct {
	svc service.CartService
}

// NewCartServiceImpl creates a new CartServiceImpl.
func NewCartServiceImpl() *CartServiceImpl {
	return &CartServiceImpl{
		svc: service.NewCartService(),
	}
}

// AddItem implements the CartServiceImpl interface.
func (s *CartServiceImpl) AddItem(ctx context.Context, req *cart.AddItemReq) (resp *cart.AddItemResp, err error) {
	return s.svc.AddItem(ctx, req)
}

// GetCart implements the CartServiceImpl interface.
func (s *CartServiceImpl) GetCart(ctx context.Context, req *cart.GetCartReq) (resp *cart.GetCartResp, err error) {
	return s.svc.GetCart(ctx, req)
}

// EmptyCart implements the CartServiceImpl interface.
func (s *CartServiceImpl) EmptyCart(ctx context.Context, req *cart.EmptyCartReq) (resp *cart.EmptyCartResp, err error) {
	return s.svc.EmptyCart(ctx, req)
}

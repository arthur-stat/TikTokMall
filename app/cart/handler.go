package main

import (
	"context"

	"TikTokMall/app/cart/biz/handler"
	"TikTokMall/app/cart/kitex_gen/cart"
)

// CartServiceImpl implements the last service interface defined in the IDL.
type CartServiceImpl struct{}

// NewCartServiceImpl creates a new CartServiceImpl.
func NewCartServiceImpl() *CartServiceImpl {
	return &CartServiceImpl{}
}

// AddItem implements the CartServiceImpl interface.
func (s *CartServiceImpl) AddItem(ctx context.Context, req *cart.AddItemReq) (resp *cart.AddItemResp, err error) {
	httpHandler := handler.NewCartHTTPHandler()
	return httpHandler.Svc.AddItem(ctx, req)
}

// GetCart implements the CartServiceImpl interface.
func (s *CartServiceImpl) GetCart(ctx context.Context, req *cart.GetCartReq) (resp *cart.GetCartResp, err error) {
	httpHandler := handler.NewCartHTTPHandler()
	return httpHandler.Svc.GetCart(ctx, req)
}

// EmptyCart implements the CartServiceImpl interface.
func (s *CartServiceImpl) EmptyCart(ctx context.Context, req *cart.EmptyCartReq) (resp *cart.EmptyCartResp, err error) {
	httpHandler := handler.NewCartHTTPHandler()
	return httpHandler.Svc.EmptyCart(ctx, req)
}

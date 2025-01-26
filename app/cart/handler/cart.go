package handler

import (
    "context"

    "TikTokMall/app/cart/biz/service"
    "TikTokMall/app/cart/kitex_gen/cart"
)

// CartServiceImpl implements the cart service interface
type CartServiceImpl struct {
    svc *service.CartService
}

// NewCartServiceImpl creates a new cart service implementation
func NewCartServiceImpl() *CartServiceImpl {
    return &CartServiceImpl{
        svc: service.NewCartService(),
    }
}

// AddItem adds an item to the cart
func (s *CartServiceImpl) AddItem(ctx context.Context, req *cart.AddItemReq) (resp *cart.AddItemResp, err error) {
    err = s.svc.AddItem(ctx, req)
    if err != nil {
        return &cart.AddItemResp{}, err
    }
    return &cart.AddItemResp{}, nil
}

// GetCart retrieves all items in the cart
func (s *CartServiceImpl) GetCart(ctx context.Context, req *cart.GetCartReq) (resp *cart.GetCartResp, err error) {
    result, err := s.svc.GetCart(ctx, req.UserId)
    if err != nil {
        return &cart.GetCartResp{}, err
    }
    return &cart.GetCartResp{
        Cart: result,
    }, nil
}

// EmptyCart removes all items from the cart
func (s *CartServiceImpl) EmptyCart(ctx context.Context, req *cart.EmptyCartReq) (resp *cart.EmptyCartResp, err error) {
    err = s.svc.EmptyCart(ctx, req.UserId)
    if err != nil {
        return &cart.EmptyCartResp{}, err
    }
    return &cart.EmptyCartResp{}, nil
}

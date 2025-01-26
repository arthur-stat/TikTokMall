package service

import (
    "context"
    "errors"

    "TikTokMall/app/cart/biz/dal/mysql"
    "TikTokMall/app/cart/biz/dal/redis"
    "TikTokMall/app/cart/biz/model"
    "TikTokMall/app/cart/kitex_gen/cart"
)

var (
    ErrInvalidQuantity = errors.New("quantity must be greater than 0")
    ErrUserNotFound    = errors.New("user not found")
)

type CartService struct{}

// NewCartService creates a new cart service instance
func NewCartService() *CartService {
    return &CartService{}
}

// AddItem adds an item to the user's cart
func (s *CartService) AddItem(ctx context.Context, req *cart.AddItemReq) error {
    if req.Item.Quantity <= 0 {
        return ErrInvalidQuantity
    }

    // Create cart item
    item := &model.CartItem{
        UserID:    req.UserId,
        ProductID: req.Item.ProductId,
        Quantity:  uint32(req.Item.Quantity),
        Selected:  true,
    }

    // Add to database
    if err := mysql.AddCartItem(ctx, item); err != nil {
        return err
    }

    // Invalidate cache to force refresh
    return redis.InvalidateCartCache(ctx, req.UserId)
}

// GetCart retrieves the user's cart
func (s *CartService) GetCart(ctx context.Context, userID uint32) (*cart.Cart, error) {
    // Try to get from cache first
    items, err := redis.GetCachedCart(ctx, userID)
    if err == nil {
        return convertToProtoCart(userID, items), nil
    }

    // If cache miss, get from database
    items, err = mysql.GetUserCart(ctx, userID)
    if err != nil {
        return nil, err
    }

    // Cache the result
    if err := redis.CacheCart(ctx, userID, items); err != nil {
        // Log the error but don't fail the request
        // TODO: Add proper logging
    }

    return convertToProtoCart(userID, items), nil
}

// EmptyCart removes all items from the user's cart
func (s *CartService) EmptyCart(ctx context.Context, userID uint32) error {
    // Empty cart in database
    if err := mysql.EmptyCart(ctx, userID); err != nil {
        return err
    }

    // Invalidate cache
    return redis.InvalidateCartCache(ctx, userID)
}

// Helper function to convert database items to proto message
func convertToProtoCart(userID uint32, items []*model.CartItem) *cart.Cart {
    protoItems := make([]*cart.CartItem, 0, len(items))
    for _, item := range items {
        protoItems = append(protoItems, &cart.CartItem{
            ProductId: item.ProductID,
            Quantity:  int32(item.Quantity),
        })
    }
    return &cart.Cart{
        UserId: userID,
        Items:  protoItems,
    }
}

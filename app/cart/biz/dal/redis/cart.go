package redis

import (
    "context"
    "encoding/json"
    "fmt"

    "TikTokMall/app/cart/biz/model"
)

const (
    cartKeyPrefix = "cart:"
)

// CartCache 购物车缓存结构
type CartCache struct {
    UserID uint32           `json:"user_id"`
    Items  []*model.CartItem `json:"items"`
}

// cartKey 生成购物车的Redis键
func cartKey(userID uint32) string {
    return fmt.Sprintf("%s%d", cartKeyPrefix, userID)
}

// CacheCart 缓存购物车数据
func CacheCart(ctx context.Context, userID uint32, items []*model.CartItem) error {
    cart := &CartCache{
        UserID: userID,
        Items:  items,
    }
    
    data, err := json.Marshal(cart)
    if err != nil {
        return err
    }
    
    return Client.Set(ctx, cartKey(userID), data, cartExpiration).Err()
}

// GetCachedCart 获取缓存的购物车数据
func GetCachedCart(ctx context.Context, userID uint32) ([]*model.CartItem, error) {
    data, err := Client.Get(ctx, cartKey(userID)).Bytes()
    if err != nil {
        return nil, err
    }
    
    var cart CartCache
    if err := json.Unmarshal(data, &cart); err != nil {
        return nil, err
    }
    
    return cart.Items, nil
}

// InvalidateCartCache 使购物车缓存失效
func InvalidateCartCache(ctx context.Context, userID uint32) error {
    return Client.Del(ctx, cartKey(userID)).Err()
}

// UpdateCartCache 更新购物车缓存
func UpdateCartCache(ctx context.Context, userID uint32, items []*model.CartItem) error {
    return CacheCart(ctx, userID, items)
}

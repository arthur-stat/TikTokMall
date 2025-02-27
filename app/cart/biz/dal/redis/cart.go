package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"TikTokMall/app/cart/biz/model"
)

const (
	cartKeyPrefix  = "cart:"
	cartExpiration = 24 * time.Hour
)

// GetCachedCart 从缓存获取购物车
func GetCachedCart(ctx context.Context, userID uint32) ([]*model.CartItem, error) {
	key := fmt.Sprintf("%s%d", cartKeyPrefix, userID)
	data, err := RDB.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	var items []*model.CartItem
	if err := json.Unmarshal([]byte(data), &items); err != nil {
		return nil, err
	}

	return items, nil
}

// CacheCart 缓存购物车
func CacheCart(ctx context.Context, userID uint32, items []*model.CartItem) error {
	key := fmt.Sprintf("%s%d", cartKeyPrefix, userID)
	data, err := json.Marshal(items)
	if err != nil {
		return err
	}

	return RDB.Set(ctx, key, string(data), cartExpiration).Err()
}

// UpdateCartCache 更新购物车缓存
func UpdateCartCache(ctx context.Context, userID uint32, items []*model.CartItem) error {
	return CacheCart(ctx, userID, items)
}

// InvalidateCartCache 使购物车缓存失效
func InvalidateCartCache(ctx context.Context, userID uint32) error {
	key := fmt.Sprintf("%s%d", cartKeyPrefix, userID)
	return RDB.Del(ctx, key).Err()
}

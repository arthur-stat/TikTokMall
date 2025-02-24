package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"TikTokMall/app/order/biz/dal/mysql"
)

const (
	orderKeyPrefix  = "order:"
	orderExpiration = 24 * time.Hour
)

// CacheOrder 缓存订单信息
func CacheOrder(ctx context.Context, order *mysql.Order) error {
	key := fmt.Sprintf("%s%d", orderKeyPrefix, order.ID)
	data, err := json.Marshal(order)
	if err != nil {
		return err
	}
	return RDB.Set(ctx, key, string(data), orderExpiration).Err()
}

// GetCachedOrder 获取缓存的订单信息
func GetCachedOrder(ctx context.Context, orderID int64) (*mysql.Order, error) {
	key := fmt.Sprintf("%s%d", orderKeyPrefix, orderID)
	data, err := RDB.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	var order mysql.Order
	if err := json.Unmarshal([]byte(data), &order); err != nil {
		return nil, err
	}
	return &order, nil
}

// InvalidateOrderCache 使订单缓存失效
func InvalidateOrderCache(ctx context.Context, orderID int64) error {
	key := fmt.Sprintf("%s%d", orderKeyPrefix, orderID)
	return RDB.Del(ctx, key).Err()
}

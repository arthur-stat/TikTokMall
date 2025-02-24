package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"TikTokMall/app/payment/biz/model"
	"github.com/redis/go-redis/v9"
)

// SetPaymentCache 将支付记录缓存到 Redis，过期时间设置为 24 小时
func SetPaymentCache(ctx context.Context, payment *model.Payments) error {
	key := fmt.Sprintf("payment:%d", payment.OrderID)
	data, err := json.Marshal(payment)
	if err != nil {
		return fmt.Errorf("failed to marshal payment data: %v", err)
	}
	err = Client.Set(ctx, key, data, 24*time.Hour).Err()
	if err != nil {
		return fmt.Errorf("failed to set payment cache: %v", err)
	}
	return nil
}

// GetPaymentCache 从 Redis 中获取支付记录
func GetPaymentCache(ctx context.Context, orderID int64) (*model.Payments, error) {
	key := fmt.Sprintf("payment:%d", orderID)
	data, err := Client.Get(ctx, key).Bytes()
	if err == redis.Nil {
		return nil, nil // 缓存中没有记录
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get payment cache: %v", err)
	}

	var payment model.Payments
	if err := json.Unmarshal(data, &payment); err != nil {
		return nil, fmt.Errorf("failed to unmarshal cached payment data: %v", err)
	}
	// 确保 OrderID 正确
	payment.OrderID = orderID
	return &payment, nil
}

// SetPayUrlToCache 将支付url缓存到 Redis，过期时间设置为 10 min
func SetPayUrlToCache(ctx context.Context, orderId int64, payUrl string) error {
	key := strconv.FormatInt(orderId, 10)
	data := payUrl
	err := Client.Set(ctx, key, data, 10*time.Minute).Err()
	if err != nil {
		return fmt.Errorf("failed to set payUrl chahe: %v", err)
	}
	return nil
}

func GetPayUrlFromCache(ctx context.Context, orderId int64) (string, error) {
	key := strconv.FormatInt(orderId, 10)
	data, err := Client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", nil // 缓存中没有记录
	}
	if err != nil {
		return "", fmt.Errorf("failed to get payUrl cache: %v", err)
	}
	return data, nil
}

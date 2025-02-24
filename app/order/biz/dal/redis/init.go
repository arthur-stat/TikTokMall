package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"TikTokMall/app/order/conf"
)

var RDB *redis.Client

// Init 初始化 Redis 连接
func Init() error {
	config := conf.GetConf()
	RDB = redis.NewClient(&redis.Options{
		Addr:     config.Redis.Address,
		Username: config.Redis.Username,
		Password: config.Redis.Password,
		DB:       config.Redis.DB,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := RDB.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("无法连接到Redis: %w", err)
	}

	return nil
}

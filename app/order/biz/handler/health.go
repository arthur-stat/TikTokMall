package handler

import (
	"context"

	"TikTokMall/app/order/biz/dal/mysql"
	"TikTokMall/app/order/biz/dal/redis"
)

// HealthCheck 健康检查
func HealthCheck(ctx context.Context) error {
	// 检查 MySQL 连接
	if err := mysql.DB.WithContext(ctx).Raw("SELECT 1").Error; err != nil {
		return err
	}

	// 检查 Redis 连接
	if err := redis.RDB.Ping(ctx).Err(); err != nil {
		return err
	}

	return nil
}

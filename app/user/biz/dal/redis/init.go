package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// 全局Redis客户端（与auth服务隔离）
var RDB *redis.Client

// Init 初始化Redis连接（参数化配置）
func Init(addr, password string, db int) error {
	RDB = redis.NewClient(&redis.Options{
		Addr:         addr,
		Password:     password,
		DB:           db,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		PoolSize:     10,
	})

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := RDB.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("redis connection failed: %w", err)
	}

	return nil
}

// Close 关闭连接
func Close() error {
	return RDB.Close()
}

// // 旧redis初始化代码
//package redis
//
//import (
//	"context"
//
//	"github.com/redis/go-redis/v9"
//	"TikTokMall/app/user/conf"
//)
//
//var (
//	RedisClient *redis.Client
//)
//
//func Init() {
//	RedisClient = redis.NewClient(&redis.Options{
//		Addr:     conf.GetConf().Redis.Address,
//		Username: conf.GetConf().Redis.Username,
//		Password: conf.GetConf().Redis.Password,
//		DB:       conf.GetConf().Redis.DB,
//	})
//	if err := RedisClient.Ping(context.Background()).Err(); err != nil {
//		panic(err)
//	}
//}

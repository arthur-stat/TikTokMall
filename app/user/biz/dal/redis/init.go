package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	UserRDB *redis.Client
)

// Init 初始化 Redis 连接
func Init(addr, password string, db int) error {
	UserRDB = redis.NewClient(&redis.Options{
		Addr:         addr,
		Password:     password,
		DB:           db,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		PoolSize:     10,
		PoolTimeout:  4 * time.Second,
	})

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := UserRDB.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("ping redis failed: %w", err)
	}

	return nil
}

// Close 关闭 Redis 连接
func Close() error {
	return UserRDB.Close()
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

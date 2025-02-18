package redis

import (
	"TikTokMall/app/user/conf"
	"context"

	"github.com/redis/go-redis/v9"
)

var RDB *redis.Client

// Init initializes the Redis connection
func Init() {
	c := conf.GetConf().Redis
	RDB = redis.NewClient(&redis.Options{
		Addr:     c.Addr,
		Password: c.Password,
		DB:       c.DB,
	})

	// Test connection
	if err := RDB.Ping(context.Background()).Err(); err != nil {
		panic(err)
	}
}
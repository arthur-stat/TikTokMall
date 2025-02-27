package redis

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

var RDB RedisClient

func Init() error {
	// 检查是否在测试环境中
	if os.Getenv("GO_ENV") == "test" || os.Getenv("TESTING") == "1" {
		// 在测试中使用 mock 客户端
		RDB = NewMockRedisClient()
		return nil
	}

	// 实际环境连接
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT")),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 尝试ping，但如果失败也可以继续（部分功能可能不可用）
	if err := client.Ping(ctx).Err(); err != nil {
		fmt.Printf("警告: 无法连接到Redis: %v\n", err)
		// 在测试环境中，我们不把这当作错误
		if os.Getenv("GO_ENV") != "test" && os.Getenv("TESTING") != "1" {
			return err
		}
	}

	RDB = client
	return nil
}

package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

const (
	loginRetryKeyPrefix = "login_retry:"
	retryExpiration     = time.Hour // 重试计数器过期时间
)

// GetLoginRetryCount 获取登录重试次数
func GetLoginRetryCount(ctx context.Context, username string) (int, error) {
	if Client == nil {
		return 0, nil
	}
	key := fmt.Sprintf("%s%s", loginRetryKeyPrefix, username)
	count, err := RDB.Get(ctx, key).Int()
	if err == redis.Nil {
		return 0, nil
	}
	return count, err
}

// IncrLoginRetry 增加登录重试次数
func IncrLoginRetry(ctx context.Context, username string) (int, error) {
	if Client == nil {
		return 0, nil
	}
	return Client.GetLoginRetryCount(ctx, username)
}

// ResetLoginRetry 重置登录重试次数
func ResetLoginRetry(ctx context.Context, username string) error {
	if Client == nil {
		return nil
	}
	return Client.ResetLoginRetry(ctx, username)
}

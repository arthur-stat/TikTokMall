package redis

import (
	"context"
	"fmt"
	"time"
)

const (
	loginRetryKeyPrefix = "login_retry:"
	retryExpiration     = time.Hour // 重试计数器过期时间
)

// GetLoginRetryCount 获取登录重试次数
func GetLoginRetryCount(ctx context.Context, username string) (int64, error) {
	key := fmt.Sprintf("%s%s", loginRetryKeyPrefix, username)
	count, err := RDB.Get(ctx, key).Int64()
	if err != nil && err.Error() == "redis: nil" {
		return 0, nil
	}
	return count, err
}

// IncrLoginRetry 增加登录重试次数
func IncrLoginRetry(ctx context.Context, username string) (int64, error) {
	key := fmt.Sprintf("%s%s", loginRetryKeyPrefix, username)
	pipe := RDB.Pipeline()
	incr := pipe.Incr(ctx, key)
	pipe.Expire(ctx, key, retryExpiration)
	_, err := pipe.Exec(ctx)
	if err != nil {
		return 0, err
	}
	return incr.Val(), nil
}

// ResetLoginRetry 重置登录重试次数
func ResetLoginRetry(ctx context.Context, username string) error {
	key := fmt.Sprintf("%s%s", loginRetryKeyPrefix, username)
	return RDB.Del(ctx, key).Err()
}

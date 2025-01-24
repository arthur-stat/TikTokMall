package redis

import (
	"context"
	"fmt"
	"time"
)

const (
	// Token相关的Key前缀
	tokenKeyPrefix     = "auth:token:"
	blacklistKeyPrefix = "auth:blacklist:"
	retryKeyPrefix     = "auth:retry:"
)

// CacheToken 缓存Token
func CacheToken(ctx context.Context, token string, userID int64, expiration time.Duration) error {
	key := fmt.Sprintf("%s%s", tokenKeyPrefix, token)
	return RDB.Set(ctx, key, userID, expiration).Err()
}

// GetCachedUserID 获取缓存的用户ID
func GetCachedUserID(ctx context.Context, token string) (int64, error) {
	key := fmt.Sprintf("%s%s", tokenKeyPrefix, token)
	val, err := RDB.Get(ctx, key).Int64()
	if err != nil {
		return 0, err
	}
	return val, nil
}

// DeleteToken 删除Token缓存
func DeleteToken(ctx context.Context, token string) error {
	key := fmt.Sprintf("%s%s", tokenKeyPrefix, token)
	return RDB.Del(ctx, key).Err()
}

// AddToBlacklist 将Token加入黑名单
func AddToBlacklist(ctx context.Context, token string, expiration time.Duration) error {
	key := fmt.Sprintf("%s%s", blacklistKeyPrefix, token)
	return RDB.Set(ctx, key, 1, expiration).Err()
}

// IsInBlacklist 检查Token是否在黑名单中
func IsInBlacklist(ctx context.Context, token string) (bool, error) {
	key := fmt.Sprintf("%s%s", blacklistKeyPrefix, token)
	exists, err := RDB.Exists(ctx, key).Result()
	return exists > 0, err
}

// IncrLoginRetry 增加登录重试次数
func IncrLoginRetry(ctx context.Context, username string) (int64, error) {
	key := fmt.Sprintf("%s%s", retryKeyPrefix, username)
	pipe := RDB.TxPipeline()
	incr := pipe.Incr(ctx, key)
	pipe.Expire(ctx, key, time.Hour) // 1小时后重置
	_, err := pipe.Exec(ctx)
	if err != nil {
		return 0, err
	}
	return incr.Val(), nil
}

// ResetLoginRetry 重置登录重试次数
func ResetLoginRetry(ctx context.Context, username string) error {
	key := fmt.Sprintf("%s%s", retryKeyPrefix, username)
	return RDB.Del(ctx, key).Err()
}

// GetLoginRetryCount 获取登录重试次数
func GetLoginRetryCount(ctx context.Context, username string) (int64, error) {
	key := fmt.Sprintf("%s%s", retryKeyPrefix, username)
	count, err := RDB.Get(ctx, key).Int64()
	if err != nil && err.Error() == "redis: nil" {
		return 0, nil
	}
	return count, err
}

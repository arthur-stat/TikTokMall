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

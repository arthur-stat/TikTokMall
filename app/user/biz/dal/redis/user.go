package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"TikTokMall/app/user/biz/dal/mysql"
)

const (
	userCacheKeyPrefix = "user:"
)

// CacheUser 缓存用户信息
func CacheUser(ctx context.Context, user *mysql.User, expiration time.Duration) error {
	key := fmt.Sprintf("%s%d", userCacheKeyPrefix, user.ID)
	data, err := json.Marshal(user)
	if err != nil {
		return err
	}
	return RDB.Set(ctx, key, data, expiration).Err()
}

// GetUserCache 获取缓存的用户信息
func GetUserCache(ctx context.Context, userID int64) (*mysql.User, error) {
	key := fmt.Sprintf("%s%d", userCacheKeyPrefix, userID)
	data, err := RDB.Get(ctx, key).Bytes()
	if err != nil {
		return nil, err
	}
	var user mysql.User
	if err := json.Unmarshal(data, &user); err != nil {
		return nil, err
	}
	return &user, nil
}

// DeleteUserCache 删除用户缓存
func DeleteUserCache(ctx context.Context, userID int64) error {
	key := fmt.Sprintf("%s%d", userCacheKeyPrefix, userID)
	return RDB.Del(ctx, key).Err()
}

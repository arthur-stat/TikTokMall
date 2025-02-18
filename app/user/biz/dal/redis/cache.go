package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"TikTokMall/app/user/biz/dal/mysql"
)

const (
	userKeyPrefix = "user:info:"
	userExpire    = time.Hour * 24
)

// SetUser caches user information
func SetUser(ctx context.Context, user *mysql.User) error {
	key := fmt.Sprintf("%s%d", userKeyPrefix, user.ID)
	data, err := json.Marshal(user)
	if err != nil {
		return err
	}
	return RDB.Set(ctx, key, data, userExpire).Err()
}

// GetUser retrieves user information from cache
func GetUser(ctx context.Context, userID int64) (*mysql.User, error) {
	key := fmt.Sprintf("%s%d", userKeyPrefix, userID)
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

// DeleteUser removes user from cache
func DeleteUser(ctx context.Context, userID int64) error {
	key := fmt.Sprintf("%s%d", userKeyPrefix, userID)
	return RDB.Del(ctx, key).Err()
}
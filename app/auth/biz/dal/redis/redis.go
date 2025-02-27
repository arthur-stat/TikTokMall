package redis

import (
	"context"
	"time"
)

var Client RedisClient

type RedisClient interface {
	CacheToken(ctx context.Context, token string, userID int64, expiration time.Duration) error
	DeleteToken(ctx context.Context, token string) error
	AddToBlacklist(ctx context.Context, token string, expiration time.Duration) error
	GetLoginRetryCount(ctx context.Context, username string) (int, error)
	ResetLoginRetry(ctx context.Context, username string) error
}

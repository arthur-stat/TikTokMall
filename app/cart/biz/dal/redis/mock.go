package redis

import (
	"context"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisClient 是一个抽象的 Redis 客户端接口
type RedisClient interface {
	Get(ctx context.Context, key string) *redis.StringCmd
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd
	Del(ctx context.Context, keys ...string) *redis.IntCmd
}

// 确保实际的 redis.Client 实现了这个接口
var _ RedisClient = (*redis.Client)(nil)

// MockRedisClient 是一个用于测试的 Redis 客户端实现
type MockRedisClient struct {
	storage sync.Map // 用于存储键值对
}

// NewMockRedisClient 创建一个新的模拟Redis客户端
func NewMockRedisClient() *MockRedisClient {
	return &MockRedisClient{
		storage: sync.Map{},
	}
}

func (m *MockRedisClient) Get(ctx context.Context, key string) *redis.StringCmd {
	cmd := redis.NewStringCmd(ctx)
	if val, ok := m.storage.Load(key); ok {
		cmd.SetVal(val.(string))
		return cmd
	}
	cmd.SetErr(redis.Nil)
	return cmd
}

func (m *MockRedisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	cmd := redis.NewStatusCmd(ctx)
	strValue, ok := value.(string)
	if !ok {
		cmd.SetErr(redis.Nil)
		return cmd
	}

	m.storage.Store(key, strValue)
	cmd.SetVal("OK")
	return cmd
}

func (m *MockRedisClient) Del(ctx context.Context, keys ...string) *redis.IntCmd {
	cmd := redis.NewIntCmd(ctx)
	var count int64
	for _, key := range keys {
		if _, ok := m.storage.LoadAndDelete(key); ok {
			count++
		}
	}
	cmd.SetVal(count)
	return cmd
}

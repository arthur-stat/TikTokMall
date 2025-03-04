package redis

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

// 确保实际的 redis.Client 实现了这个接口
// var _ RedisClient = (*redis.Client)(nil)

// MockRedisClient 是一个用于测试的 Redis 客户端实现
type MockRedisClient struct {
	storage     sync.Map // 用于存储键值对
	hashStorage sync.Map // 用于存储哈希表
}

// NewMockRedisClient 创建一个新的模拟Redis客户端
func NewMockRedisClient() *MockRedisClient {
	return &MockRedisClient{
		storage:     sync.Map{},
		hashStorage: sync.Map{},
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

	// 实际删除数据
	for _, key := range keys {
		m.storage.Delete(key)
		m.hashStorage.Delete(key) // 同时删除哈希表
	}

	cmd.SetVal(int64(len(keys))) // 返回删除的键数量
	return cmd
}

// HSet 实现哈希表的设置操作
func (m *MockRedisClient) HSet(ctx context.Context, key string, values ...interface{}) *redis.IntCmd {
	cmd := redis.NewIntCmd(ctx)

	// 获取或创建哈希表
	hashMapInterface, _ := m.hashStorage.LoadOrStore(key, &sync.Map{})
	hashMap := hashMapInterface.(*sync.Map)

	// 确保参数数量是偶数
	if len(values)%2 != 0 {
		cmd.SetErr(fmt.Errorf("wrong number of arguments for HSet"))
		return cmd
	}

	var count int64
	// 处理键值对
	for i := 0; i < len(values); i += 2 {
		field, ok := values[i].(string)
		if !ok {
			continue
		}

		value := values[i+1]
		strValue, ok := value.(string)
		if !ok {
			strValue = fmt.Sprintf("%v", value)
		}

		_, exists := hashMap.Load(field)
		hashMap.Store(field, strValue)

		// 如果字段不存在，计数加1
		if !exists {
			count++
		}
	}

	cmd.SetVal(count)
	return cmd
}

// HGet 实现哈希表的获取操作
func (m *MockRedisClient) HGet(ctx context.Context, key string, field string) *redis.StringCmd {
	cmd := redis.NewStringCmd(ctx)

	// 获取哈希表
	hashMapInterface, ok := m.hashStorage.Load(key)
	if !ok {
		cmd.SetErr(redis.Nil)
		return cmd
	}

	hashMap := hashMapInterface.(*sync.Map)

	// 获取字段值
	value, ok := hashMap.Load(field)
	if !ok {
		cmd.SetErr(redis.Nil)
		return cmd
	}

	cmd.SetVal(value.(string))
	return cmd
}

// HGetAll 实现哈希表的获取所有字段和值的操作
func (m *MockRedisClient) HGetAll(ctx context.Context, key string) *redis.MapStringStringCmd {
	cmd := redis.NewMapStringStringCmd(ctx)

	// 获取哈希表
	hashMapInterface, ok := m.hashStorage.Load(key)
	if !ok {
		// 返回空map而不是错误
		cmd.SetVal(make(map[string]string))
		return cmd
	}

	hashMap := hashMapInterface.(*sync.Map)

	// 构建结果map
	result := make(map[string]string)
	hashMap.Range(func(k, v interface{}) bool {
		result[k.(string)] = v.(string)
		return true
	})

	cmd.SetVal(result)
	return cmd
}

// HDel 实现哈希表的删除字段操作
func (m *MockRedisClient) HDel(ctx context.Context, key string, fields ...string) *redis.IntCmd {
	cmd := redis.NewIntCmd(ctx)

	// 获取哈希表
	hashMapInterface, ok := m.hashStorage.Load(key)
	if !ok {
		cmd.SetVal(0)
		return cmd
	}

	hashMap := hashMapInterface.(*sync.Map)

	// 删除字段
	var count int64
	for _, field := range fields {
		if _, ok := hashMap.LoadAndDelete(field); ok {
			count++
		}
	}

	cmd.SetVal(count)
	return cmd
}

// Ping 实现ping操作
func (m *MockRedisClient) Ping(ctx context.Context) *redis.StatusCmd {
	cmd := redis.NewStatusCmd(ctx)
	cmd.SetVal("PONG")
	return cmd
}

// Close 实现关闭操作
func (m *MockRedisClient) Close() error {
	return nil
}

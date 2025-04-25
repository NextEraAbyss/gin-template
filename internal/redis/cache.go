package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	redisClient "github.com/redis/go-redis/v9"
)

// Cache 缓存接口.
type Cache interface {
	Get(ctx context.Context, key string, value interface{}) error
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Delete(ctx context.Context, key string) error
	Exists(ctx context.Context, key string) (bool, error)
	FlushDB(ctx context.Context) error
}

// RedisCache Redis缓存实现.
type RedisCache struct {
	client *redisClient.Client
}

// NewRedisCache 创建Redis缓存.
func NewRedisCache(client *redisClient.Client) *RedisCache {
	return &RedisCache{
		client: client,
	}
}

// Get 获取缓存.
func (c *RedisCache) Get(ctx context.Context, key string, value interface{}) error {
	var data interface{}
	var err error

	data, err = c.client.Get(ctx, key).Result()
	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(data.(string)), value)
}

// Set 设置缓存.
func (c *RedisCache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return c.client.Set(ctx, key, data, expiration).Err()
}

// Delete 删除缓存.
func (c *RedisCache) Delete(ctx context.Context, key string) error {
	var err error

	err = c.client.Del(ctx, key).Err()
	if err != nil {
		return err
	}

	return nil
}

// Exists 检查缓存是否存在.
func (c *RedisCache) Exists(ctx context.Context, key string) (bool, error) {
	var result int64
	var err error

	result, err = c.client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}

	return result > 0, nil
}

// GetOrSet 获取缓存，如果不存在则设置
func (c *RedisCache) GetOrSet(ctx context.Context, key string, value interface{}, expiration time.Duration, fn func() (interface{}, error)) error {
	// 尝试获取缓存
	err := c.Get(ctx, key, value)
	if err == nil {
		return nil
	}

	// 缓存不存在，执行回调函数
	data, err := fn()
	if err != nil {
		return err
	}

	// 设置缓存
	return c.Set(ctx, key, data, expiration)
}

// ClearByPrefix 清除指定前缀的缓存
func (c *RedisCache) ClearByPrefix(ctx context.Context, prefix string) error {
	pattern := fmt.Sprintf("%s*", prefix)
	keys, err := c.client.Keys(ctx, pattern).Result()
	if err != nil {
		return err
	}

	if len(keys) > 0 {
		return c.client.Del(ctx, keys...).Err()
	}

	return nil
}

// Increment 递增计数器
func (c *RedisCache) Increment(ctx context.Context, key string) (int64, error) {
	return c.client.Incr(ctx, key).Result()
}

// Decrement 递减计数器
func (c *RedisCache) Decrement(ctx context.Context, key string) (int64, error) {
	return c.client.Decr(ctx, key).Result()
}

// SetNX 设置缓存，如果不存在
func (c *RedisCache) SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) (bool, error) {
	data, err := json.Marshal(value)
	if err != nil {
		return false, err
	}

	return c.client.SetNX(ctx, key, data, expiration).Result()
}

// FlushDB 清空缓存.
func (c *RedisCache) FlushDB(ctx context.Context) error {
	return c.client.FlushDB(ctx).Err()
}

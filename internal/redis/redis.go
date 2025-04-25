package redis

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	// Client Redis客户端
	Client *redis.Client
	// Ctx 上下文
	Ctx = context.Background()
)

// InitRedis 初始化Redis连接
func InitRedis(addr, password string, db int) error {
	Client = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	// 测试连接
	ctx := context.Background()
	if err := Client.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("failed to connect to Redis: %v", err)
	}

	log.Println("Redis connected successfully")
	return nil
}

// GetClient 获取Redis客户端
func GetClient() *redis.Client {
	return Client
}

// Close 关闭Redis连接.
func Close() error {
	if Client != nil {
		return Client.Close()
	}

	return nil
}

// Set 设置键值对
func Set(key string, value interface{}, expiration time.Duration) error {
	return Client.Set(Ctx, key, value, expiration).Err()
}

// Get 获取键值
func Get(key string) (string, error) {
	return Client.Get(Ctx, key).Result()
}

// Delete 删除键
func Delete(key string) error {
	return Client.Del(Ctx, key).Err()
}

// SetHash 设置哈希表字段
func SetHash(key, field string, value interface{}) error {
	return Client.HSet(Ctx, key, field, value).Err()
}

// GetHash 获取哈希表字段
func GetHash(key, field string) (string, error) {
	return Client.HGet(Ctx, key, field).Result()
}

// HashExists 检查哈希表字段是否存在
func HashExists(key, field string) (bool, error) {
	return Client.HExists(Ctx, key, field).Result()
}

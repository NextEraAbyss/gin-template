package redis

import (
	"context"
	"fmt"
	"log"
	"time"

	"gitee.com/NextEraAbyss/gin-template/config"
	"github.com/go-redis/redis/v8"
)

var (
	// Client Redis客户端
	Client *redis.Client
	// Ctx 上下文
	Ctx = context.Background()
)

// Init 初始化Redis连接
func Init(config *config.Config) *redis.Client {
	dsn := fmt.Sprintf("%s:%s", config.Redis.Host, config.Redis.Port)

	Client = redis.NewClient(&redis.Options{
		Addr:     dsn,
		Password: config.Redis.Password,
		DB:       config.Redis.DB,
		// 连接池配置
		PoolSize:     10,
		MinIdleConns: 5,
		// 超时配置
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		PoolTimeout:  4 * time.Second,
	})

	// 测试连接
	_, err := Client.Ping(Ctx).Result()
	if err != nil {
		log.Printf("Warning: Failed to connect to Redis: %v", err)
	} else {
		log.Println("Connected to Redis successfully")
	}

	return Client
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

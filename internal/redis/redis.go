package redis

import (
	"context"
	"fmt"
	"time"

	"gitee.com/NextEraAbyss/gin-template/config"
	"github.com/redis/go-redis/v9"
)

var (
	// Client Redis客户端
	Client *redis.Client
	// Ctx 上下文
	Ctx = context.Background()
)

// InitRedis 初始化Redis连接
func InitRedis(config *config.Config) *redis.Client {
	// 创建Redis客户端
	client := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%s", config.Redis.Host, config.Redis.Port),
		Password:     config.Redis.Password,
		DB:           config.Redis.DB,
		PoolSize:     10,              // 默认连接池大小
		MinIdleConns: 5,               // 默认最小空闲连接数
		MaxRetries:   3,               // 最大重试次数
		DialTimeout:  5 * time.Second, // 连接超时时间
		ReadTimeout:  3 * time.Second, // 读取超时时间
		WriteTimeout: 3 * time.Second, // 写入超时时间
	})

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := client.Ping(ctx).Result()
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to Redis: %v", err))
	}

	Client = client
	return client
}

// CloseRedis 关闭Redis连接
func CloseRedis() error {
	if Client != nil {
		err := Client.Close()
		if err != nil {
			return err
		}
		Client = nil
	}
	return nil
}

// GetClient 获取Redis客户端
func GetClient() *redis.Client {
	return Client
}

package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"gitee.com/NextEraAbyss/gin-template/internal/redis"
	"gitee.com/NextEraAbyss/gin-template/utils"
	"github.com/gin-gonic/gin"
)

// RateLimitConfig 限流配置
type RateLimitConfig struct {
	// 时间窗口（秒）
	Window int
	// 最大请求数
	MaxRequests int
	// 是否按IP限流
	ByIP bool
	// 是否按用户限流
	ByUser bool
}

// DefaultRateLimitConfig 默认限流配置
var DefaultRateLimitConfig = RateLimitConfig{
	Window:      60,  // 1分钟
	MaxRequests: 100, // 最多100个请求
	ByIP:        true,
	ByUser:      false,
}

// RateLimitMiddleware 限流中间件
func RateLimitMiddleware(config RateLimitConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取限流键
		key := getRateLimitKey(c, config)
		if key == "" {
			c.Next()
			return
		}

		// 获取当前请求数
		count, err := redis.Client.Incr(context.Background(), key).Result()
		if err != nil {
			utils.Errorf("限流错误: %v", err)
			c.Next()
			return
		}

		// 设置过期时间
		if count == 1 {
			redis.Client.Expire(context.Background(), key, time.Duration(config.Window)*time.Second)
		}

		// 检查是否超过限制
		if count > int64(config.MaxRequests) {
			// 设置响应头
			c.Header("X-RateLimit-Limit", strconv.Itoa(config.MaxRequests))
			c.Header("X-RateLimit-Remaining", "0")
			c.Header("X-RateLimit-Reset", strconv.FormatInt(time.Now().Add(time.Duration(config.Window)*time.Second).Unix(), 10))

			// 返回429状态码
			utils.ResponseErrorWithStatus(c, http.StatusTooManyRequests, utils.CodeServerError, "请求过于频繁，请稍后再试")
			c.Abort()
			return
		}

		// 设置响应头
		c.Header("X-RateLimit-Limit", strconv.Itoa(config.MaxRequests))
		c.Header("X-RateLimit-Remaining", strconv.FormatInt(int64(config.MaxRequests)-count, 10))
		c.Header("X-RateLimit-Reset", strconv.FormatInt(time.Now().Add(time.Duration(config.Window)*time.Second).Unix(), 10))

		c.Next()
	}
}

// getRateLimitKey 获取限流键
func getRateLimitKey(c *gin.Context, config RateLimitConfig) string {
	// 限流键前缀
	prefix := "rate_limit:"

	// 按IP限流
	if config.ByIP {
		ip := c.ClientIP()
		return fmt.Sprintf("%sip:%s", prefix, ip)
	}

	// 按用户限流
	if config.ByUser {
		userID, exists := c.Get("user_id")
		if exists {
			return fmt.Sprintf("%suser:%v", prefix, userID)
		}
	}

	// 按路径限流
	return fmt.Sprintf("%spath:%s", prefix, c.Request.URL.Path)
}

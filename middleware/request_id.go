package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	// RequestIDKey 请求ID的上下文键
	RequestIDKey = "X-Request-ID"
)

// RequestID 请求ID中间件
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 检查请求头中是否已有请求ID
		requestID := c.GetHeader(RequestIDKey)
		if requestID == "" {
			// 如果没有，生成一个新的请求ID
			requestID = uuid.New().String()
		}

		// 将请求ID存储到上下文中
		c.Set(RequestIDKey, requestID)

		// 将请求ID添加到响应头中
		c.Header(RequestIDKey, requestID)

		c.Next()
	}
}

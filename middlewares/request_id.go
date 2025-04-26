package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// RequestIDHeaderName 请求ID的HTTP头名称
const RequestIDHeaderName = "X-Request-ID"

// RequestID 中间件
// 为每个请求生成一个唯一ID，并添加到请求上下文和响应头中
// 可用于跟踪和调试请求
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 检查请求头中是否已经有请求ID
		requestID := c.Request.Header.Get(RequestIDHeaderName)
		if requestID == "" {
			// 如果没有，生成一个新的UUID作为请求ID
			requestID = uuid.New().String()
		}

		// 将请求ID设置到上下文中
		c.Set("request_id", requestID)

		// 将请求ID添加到响应头中
		c.Writer.Header().Set(RequestIDHeaderName, requestID)

		// 继续处理请求
		c.Next()
	}
}

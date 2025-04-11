package middleware

import (
	"bytes"
	"io"
	"time"

	"gitee.com/NextEraAbyss/gin-template/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Logger 日志中间件
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 开始时间
		start := time.Now()

		// 读取请求体
		var requestBody []byte
		if c.Request.Body != nil {
			requestBody, _ = io.ReadAll(c.Request.Body)
			// 恢复请求体
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

		// 处理请求
		c.Next()

		// 结束时间
		end := time.Now()
		latency := end.Sub(start)

		// 获取请求ID
		requestID := c.GetString("request_id")
		if requestID == "" {
			requestID = uuid.New().String()
		}

		// 记录请求日志
		utils.LogInfo("Request completed",
			"request_id", requestID,
			"method", c.Request.Method,
			"uri", c.Request.RequestURI,
			"status", c.Writer.Status(),
			"client_ip", c.ClientIP(),
			"latency", latency,
		)

		// 如果有错误，记录错误日志
		if len(c.Errors) > 0 {
			utils.LogError("Request error",
				"request_id", requestID,
				"method", c.Request.Method,
				"uri", c.Request.RequestURI,
				"status", c.Writer.Status(),
				"errors", c.Errors.String(),
			)
		}
	}
}

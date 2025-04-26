package middlewares

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// Logger 中间件
// 记录请求的处理时间和响应状态
// 输出格式: [请求方法] [状态码] [路径] [客户端IP] [响应时间] [错误信息]
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 请求开始时间
		startTime := time.Now()

		// 处理请求
		c.Next()

		// 请求结束时间
		endTime := time.Now()

		// 计算请求处理时间
		latency := endTime.Sub(startTime)

		// 请求方法和路径
		method := c.Request.Method
		path := c.Request.URL.Path
		if c.Request.URL.RawQuery != "" {
			path = path + "?" + c.Request.URL.RawQuery
		}

		// 请求结果状态码
		statusCode := c.Writer.Status()

		// 客户端IP
		clientIP := c.ClientIP()

		// 构建日志条目
		entry := logrus.WithFields(logrus.Fields{
			"status":     statusCode,
			"latency":    latency,
			"client_ip":  clientIP,
			"method":     method,
			"path":       path,
			"request_id": c.GetString("request_id"),
		})

		// 判断请求是否出错
		if len(c.Errors) > 0 {
			// 记录错误信息
			entry.Error(c.Errors.String())
		} else {
			// 根据状态码选择日志级别
			if statusCode >= 500 {
				entry.Error(fmt.Sprintf("[%s] %d %s %s %s", method, statusCode, path, clientIP, latency))
			} else if statusCode >= 400 {
				entry.Warn(fmt.Sprintf("[%s] %d %s %s %s", method, statusCode, path, clientIP, latency))
			} else {
				entry.Info(fmt.Sprintf("[%s] %d %s %s %s", method, statusCode, path, clientIP, latency))
			}
		}
	}
}

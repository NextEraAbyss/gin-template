package middleware

import (
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/microcosm-cc/bluemonday"
	"golang.org/x/time/rate"
)

// Security 安全中间件
func Security() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 设置安全相关的响应头
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Content-Security-Policy", "default-src 'self'")
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Header("Permissions-Policy", "geolocation=(), microphone=(), camera=()")

		// 检查请求方法
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		// XSS防护
		if c.Request.Method == "POST" || c.Request.Method == "PUT" {
			// 使用bluemonday进行XSS过滤
			p := bluemonday.UGCPolicy()

			// 处理请求体
			if c.Request.Body != nil {
				body, err := c.GetRawData()
				if err == nil {
					// 过滤HTML内容
					cleanBody := p.Sanitize(string(body))
					c.Request.Body = io.NopCloser(strings.NewReader(cleanBody))
				}
			}
		}

		// SQL注入防护
		// 注意：这只是基本的防护，实际项目中应该使用参数化查询
		for _, value := range c.Request.URL.Query() {
			for _, v := range value {
				if containsSQLInjection(v) {
					c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
						"code":    40001,
						"message": "检测到潜在的SQL注入攻击",
					})
					return
				}
			}
		}

		c.Next()
	}
}

// containsSQLInjection 检查是否包含SQL注入特征
func containsSQLInjection(input string) bool {
	// SQL注入特征
	sqlPatterns := []string{
		"--",
		";",
		"/*",
		"*/",
		"@@",
		"@",
		"char",
		"nchar",
		"varchar",
		"nvarchar",
		"alter",
		"begin",
		"cast",
		"create",
		"cursor",
		"declare",
		"delete",
		"drop",
		"end",
		"exec",
		"execute",
		"fetch",
		"insert",
		"kill",
		"select",
		"sys",
		"sysobjects",
		"syscolumns",
		"table",
		"update",
	}

	input = strings.ToLower(input)
	for _, pattern := range sqlPatterns {
		if strings.Contains(input, pattern) {
			return true
		}
	}
	return false
}

// RateLimit 速率限制中间件
func RateLimit(limit int, window time.Duration) gin.HandlerFunc {
	limiter := rate.NewLimiter(rate.Every(window), limit)
	return func(c *gin.Context) {
		if !limiter.Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"code":    42900,
				"message": "请求过于频繁，请稍后再试",
			})
			return
		}
		c.Next()
	}
}

// IPWhitelist IP白名单中间件
func IPWhitelist(whitelist []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIP := c.ClientIP()
		for _, ip := range whitelist {
			if ip == clientIP {
				c.Next()
				return
			}
		}
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"code":    40300,
			"message": "IP不在白名单中",
		})
	}
}

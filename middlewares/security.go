package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// SecurityHeaders 添加安全相关的HTTP头
var SecurityHeaders = map[string]string{
	"X-Content-Type-Options":    "nosniff",
	"X-Frame-Options":           "DENY",
	"X-XSS-Protection":          "1; mode=block",
	"Referrer-Policy":           "strict-origin-when-cross-origin",
	"Permissions-Policy":        "accelerometer=(), camera=(), geolocation=(), gyroscope=(), magnetometer=(), microphone=(), payment=(), usb=()",
	"Strict-Transport-Security": "max-age=31536000; includeSubDomains",
	"Content-Security-Policy":   "default-src 'self'; script-src 'self' 'unsafe-inline'; style-src 'self' 'unsafe-inline'; img-src 'self' data:; font-src 'self'; object-src 'none'; base-uri 'self'; form-action 'self'; frame-ancestors 'none'; upgrade-insecure-requests; block-all-mixed-content",
}

// ContentSecurityPolicy 配置
var ContentSecurityPolicy = "default-src 'self'; " +
	"script-src 'self' 'unsafe-inline' https://cdnjs.cloudflare.com https://cdn.jsdelivr.net; " +
	"style-src 'self' 'unsafe-inline' https://cdn.jsdelivr.net https://fonts.googleapis.com; " +
	"img-src 'self' data: https://cdn.jsdelivr.net; " +
	"font-src 'self' https://fonts.gstatic.com https://cdn.jsdelivr.net; " +
	"object-src 'none'; " +
	"base-uri 'self'; " +
	"form-action 'self'; " +
	"frame-ancestors 'none'; " +
	"upgrade-insecure-requests; " +
	"block-all-mixed-content;"

// Security 中间件设置安全相关的HTTP头
// 添加各种安全头，如内容安全策略、XSS保护等
// 参考: https://owasp.org/www-project-secure-headers/
func Security() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 添加基本安全头
		for key, value := range SecurityHeaders {
			c.Header(key, value)
		}

		// 为Swagger UI设置特殊的CSP头
		if c.Request.URL.Path == "/swagger/index.html" || c.Request.URL.Path == "/swagger/doc.json" {
			c.Header("Content-Security-Policy", "default-src 'self'; script-src 'self' 'unsafe-inline' 'unsafe-eval'; style-src 'self' 'unsafe-inline'; img-src 'self' data:; font-src 'self' data:; connect-src 'self';")
		} else {
			c.Header("Content-Security-Policy", ContentSecurityPolicy)
		}

		// 添加额外的安全措施
		c.Header("X-Permitted-Cross-Domain-Policies", "none")
		c.Header("Cache-Control", "no-store, max-age=0")
		c.Header("Cross-Origin-Embedder-Policy", "require-corp")
		c.Header("Cross-Origin-Opener-Policy", "same-origin")
		c.Header("Cross-Origin-Resource-Policy", "same-origin")

		// 对于API, 禁止内容嗅探
		c.Header("X-Content-Type-Options", "nosniff")

		// 防止点击劫持攻击
		c.Header("X-Frame-Options", "DENY")

		// 设置安全Cookie
		c.SetSameSite(http.SameSiteStrictMode)
		c.SetCookie("secure", "true", 0, "/", "", true, true)

		c.Next()
	}
}

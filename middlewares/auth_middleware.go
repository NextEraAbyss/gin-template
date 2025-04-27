package middlewares

import (
	"net/http"
	"strconv"
	"strings"

	"gitee.com/NextEraAbyss/gin-template/utils"
	"github.com/gin-gonic/gin"
)

// AuthMiddleware 验证JWT令牌的中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header format must be Bearer {token}"})
			c.Abort()
			return
		}

		claims, err := utils.ParseToken(parts[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// 从Subject中提取用户ID
		userID, err := strconv.ParseUint(claims.Subject, 10, 32)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token subject"})
			c.Abort()
			return
		}

		// 将解析的claims保存到上下文中
		c.Set("userID", uint(userID))
		c.Set("username", claims.Issuer) // 使用Issuer存储username
		c.Next()
	}
}

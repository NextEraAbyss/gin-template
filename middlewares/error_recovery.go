package middlewares

import (
	"log"
	"runtime/debug"

	"gitee.com/NextEraAbyss/gin-template/internal/errors"
	"gitee.com/NextEraAbyss/gin-template/utils"
	"github.com/gin-gonic/gin"
)

// Recovery 恢复中间件
// 捕获 panic，记录错误堆栈，并返回统一的错误响应
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// 记录错误堆栈
				log.Printf("Panic: %v\nStack: %s", err, debug.Stack())

				// 返回错误响应
				utils.ResponseError(c, errors.CodeInternalError, "系统内部错误")

				// 中止请求处理
				c.Abort()
			}
		}()

		c.Next()
	}
}

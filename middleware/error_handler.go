package middleware

import (
	"log"
	"runtime/debug"

	"gitee.com/NextEraAbyss/gin-template/internal/errors"
	"gitee.com/NextEraAbyss/gin-template/utils"
	"github.com/gin-gonic/gin"
)

// ErrorHandler 错误处理中间件
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// 检查是否有错误
		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err
			if customErr, ok := err.(*errors.Error); ok {
				utils.ResponseError(c, customErr.Code, customErr.Message)
				return
			}

			// 处理其他错误
			utils.ResponseError(c, errors.CodeInternalError, "服务器内部错误")
		}
	}
}

// Recovery 恢复中间件
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

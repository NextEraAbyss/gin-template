package middleware

import (
	"context"
	"errors"
	"log"
	"runtime/debug"

	appErrors "gitee.com/NextEraAbyss/gin-template/internal/errors"
	"gitee.com/NextEraAbyss/gin-template/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ErrorHandler 错误处理中间件
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// 检查是否有错误
		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err

			// 记录错误堆栈
			log.Printf("Error: %v\nStack: %s", err, debug.Stack())

			// 处理自定义错误
			if customErr, ok := err.(*appErrors.Error); ok {
				utils.ResponseError(c, customErr.Code, customErr.Message)
				return
			}

			// 处理特定类型的错误
			switch {
			case errors.Is(err, gorm.ErrRecordNotFound):
				utils.ResponseError(c, appErrors.CodeArticleNotFound, "资源不存在")
			case errors.Is(err, context.DeadlineExceeded):
				utils.ResponseError(c, appErrors.CodeTimeout, "请求超时")
			case errors.Is(err, context.Canceled):
				utils.ResponseError(c, appErrors.CodeTimeout, "请求已取消")
			case errors.Is(err, gorm.ErrInvalidDB):
				utils.ResponseError(c, appErrors.CodeDBConnectionError, "数据库连接错误")
			default:
				// 处理其他错误
				utils.ResponseError(c, appErrors.CodeInternalError, "系统内部错误")
			}
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
				utils.ResponseError(c, appErrors.CodeInternalError, "系统内部错误")

				// 中止请求处理
				c.Abort()
			}
		}()

		c.Next()
	}
}

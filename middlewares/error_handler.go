package middlewares

import (
	"gitee.com/NextEraAbyss/gin-template/internal/errors"
	"gitee.com/NextEraAbyss/gin-template/utils"
	"github.com/gin-gonic/gin"
)

// ErrorHandler 错误处理中间件
// 捕获请求处理过程中产生的错误，并返回标准格式的错误响应
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// 检查是否有错误
		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err
			if customErr, ok := err.(*errors.Error); ok {
				utils.ResponseError(c, utils.ErrorCode(customErr.Code), customErr.Message)
				return
			}

			// 处理其他错误
			utils.ResponseError(c, utils.CodeInternalError, "服务器内部错误")
		}
	}
}

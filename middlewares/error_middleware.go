package middlewares

import (
	"net/http"
	"strings"

	"gitee.com/NextEraAbyss/gin-template/utils"
	"github.com/gin-gonic/gin"
)

// 自定义错误类型
const (
	ValidationError = "ValidationError"
	DatabaseError   = "DatabaseError"
	AuthError       = "AuthError"
	NotFoundError   = "NotFoundError"
	SystemError     = "SystemError"
)

// ErrorHandlerMiddleware 处理应用程序异常的中间件
func ErrorHandlerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 执行后续中间件和处理函数
		c.Next()

		// 如果存在错误
		if len(c.Errors) > 0 {
			// 获取最后一个错误
			err := c.Errors.Last()
			errMessage := err.Error()

			// 记录错误日志
			utils.Error("API错误: %v", errMessage)

			// 默认HTTP状态码和业务码
			statusCode := http.StatusInternalServerError
			businessCode := utils.CodeFailure

			// 根据错误类型设置响应
			if strings.Contains(errMessage, ValidationError) {
				statusCode = http.StatusBadRequest
				businessCode = utils.CodeInvalidParams
			} else if strings.Contains(errMessage, DatabaseError) {
				statusCode = http.StatusInternalServerError
				businessCode = utils.CodeFailure
			} else if strings.Contains(errMessage, AuthError) {
				statusCode = http.StatusUnauthorized
				businessCode = utils.CodeUnauthorized
			} else if strings.Contains(errMessage, NotFoundError) {
				statusCode = http.StatusNotFound
				businessCode = utils.CodeNotFound
			} else if c.Writer.Status() != http.StatusOK {
				// 如果响应状态码已经被设置，保持一致
				statusCode = c.Writer.Status()
			}

			// 构建统一的错误响应
			c.JSON(statusCode, utils.Response{
				Code:    businessCode,
				Message: errMessage,
				Data:    nil,
				Error:   err.Err.Error(),
			})

			// 中止后续处理
			c.Abort()
		}
	}
}

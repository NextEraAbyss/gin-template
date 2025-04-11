package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 标准API响应结构
type Response struct {
	Code    int         `json:"code"`            // 业务码
	Message string      `json:"message"`         // 提示信息
	Data    interface{} `json:"data"`            // 数据
	Error   string      `json:"error,omitempty"` // 错误信息，只在开发环境显示
}

// ResponseCode 业务码枚举
const (
	CodeSuccess       = 20000 // 成功
	CodeUnauthorized  = 40100 // 未授权
	CodeForbidden     = 40300 // 禁止访问
	CodeNotFound      = 40400 // 资源不存在
	CodeInvalidParams = 40001 // 请求参数错误
	CodeFailure       = 50000 // 系统内部错误
	CodeNoContent     = 20400 // 无内容
)

// CodeMessages 业务码对应消息
var CodeMessages = map[int]string{
	CodeSuccess:       "操作成功",
	CodeUnauthorized:  "未授权",
	CodeForbidden:     "禁止访问",
	CodeNotFound:      "资源不存在",
	CodeInvalidParams: "请求参数错误",
	CodeFailure:       "系统内部错误",
	CodeNoContent:     "无内容",
}

// Success 返回成功响应
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    CodeSuccess,
		Message: CodeMessages[CodeSuccess],
		Data:    data,
	})
}

// SuccessWithMessage 返回带自定义消息的成功响应
func SuccessWithMessage(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    CodeSuccess,
		Message: message,
		Data:    data,
	})
}

// NoContent 返回无内容响应
func NoContent(c *gin.Context) {
	c.JSON(http.StatusOK, Response{
		Code:    CodeNoContent,
		Message: CodeMessages[CodeNoContent],
		Data:    nil,
	})
}

// Unauthorized 返回未授权响应
func Unauthorized(c *gin.Context, message ...string) {
	msg := CodeMessages[CodeUnauthorized]
	if len(message) > 0 {
		msg = message[0]
	}
	c.JSON(http.StatusUnauthorized, Response{
		Code:    CodeUnauthorized,
		Message: msg,
		Data:    nil,
	})
}

// Forbidden 返回禁止访问响应
func Forbidden(c *gin.Context, message ...string) {
	msg := CodeMessages[CodeForbidden]
	if len(message) > 0 {
		msg = message[0]
	}
	c.JSON(http.StatusForbidden, Response{
		Code:    CodeForbidden,
		Message: msg,
		Data:    nil,
	})
}

// NotFound 返回资源不存在响应
func NotFound(c *gin.Context, message ...string) {
	msg := CodeMessages[CodeNotFound]
	if len(message) > 0 {
		msg = message[0]
	}
	c.JSON(http.StatusNotFound, Response{
		Code:    CodeNotFound,
		Message: msg,
		Data:    nil,
	})
}

// InvalidParams 返回参数错误响应
func InvalidParams(c *gin.Context, message ...string) {
	msg := CodeMessages[CodeInvalidParams]
	if len(message) > 0 {
		msg = message[0]
	}
	c.JSON(http.StatusBadRequest, Response{
		Code:    CodeInvalidParams,
		Message: msg,
		Data:    nil,
	})
}

// Failure 返回系统错误响应
func Failure(c *gin.Context, err error) {
	Error("System error: %v", err)

	c.JSON(http.StatusInternalServerError, Response{
		Code:    CodeFailure,
		Message: CodeMessages[CodeFailure],
		Data:    nil,
		Error:   err.Error(),
	})
}

// CustomError 返回自定义错误响应
func CustomError(c *gin.Context, httpStatus, code int, message string, err error) {
	errMsg := ""
	if err != nil {
		Error("Custom error: %v", err)
		errMsg = err.Error()
	}

	c.JSON(httpStatus, Response{
		Code:    code,
		Message: message,
		Data:    nil,
		Error:   errMsg,
	})
}

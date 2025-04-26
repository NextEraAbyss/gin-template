package utils

import (
	"net/http"

	"gitee.com/NextEraAbyss/gin-template/internal/errors"
	"github.com/gin-gonic/gin"
)

// Response 统一响应结构
type Response struct {
	Code    int         `json:"code"`    // 错误码
	Message string      `json:"message"` // 错误信息
	Data    interface{} `json:"data"`    // 数据
}

// ResponseSuccess 成功响应
func ResponseSuccess(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    int(errors.CodeSuccess),
		Message: "success",
		Data:    data,
	})
}

// ResponseError 错误响应
func ResponseError(c *gin.Context, code errors.ErrorCode, message string) {
	// 获取HTTP状态码
	httpCode := errors.GetHTTPCode(code)

	c.JSON(httpCode, Response{
		Code:    int(code),
		Message: message,
		Data:    nil,
	})
}

// ResponseErrorWithDetails 带详细信息的错误响应
func ResponseErrorWithDetails(c *gin.Context, code errors.ErrorCode, message, details string) {
	c.JSON(http.StatusOK, gin.H{
		"code":    code,
		"message": message,
		"details": details,
	})
}

// 系统级状态码
const (
	CodeSuccess       = 0    // 成功
	CodeInvalidParams = 1001 // 无效的参数
	CodeUnauthorized  = 1002 // 未授权
	CodeForbidden     = 1003 // 禁止访问
	CodeNotFound      = 1004 // 资源不存在
	CodeInternalError = 1005 // 内部错误
	CodeServerError   = 1005 // 服务器错误

	// 用户相关错误 (2000-2999)
	CodeUserNotFound  = 2001 // 用户不存在
	CodeUserExists    = 2002 // 用户已存在
	CodePasswordError = 2003 // 密码错误
	CodeTokenExpired  = 2004 // Token过期
	CodeTokenInvalid  = 2005 // Token无效
	CodeUserDisabled  = 2006 // 用户已禁用

	// 评论相关错误 (4000-4999)
	CodeCommentNotFound  = 4001 // 评论不存在
	CodeCommentForbidden = 4002 // 无权操作评论
)

// CodeMessages 状态码对应的消息
var CodeMessages = map[int]string{
	CodeSuccess:       "操作成功",
	CodeUnauthorized:  "未授权",
	CodeForbidden:     "禁止访问",
	CodeNotFound:      "资源不存在",
	CodeInvalidParams: "请求参数错误",
	CodeServerError:   "服务器内部错误",

	// 用户相关错误消息
	CodeUserNotFound:  "用户不存在",
	CodeUserExists:    "用户已存在",
	CodePasswordError: "密码错误",
	CodeTokenExpired:  "Token已过期",
	CodeTokenInvalid:  "Token无效",
	CodeUserDisabled:  "用户已禁用",

	// 评论相关错误消息
	CodeCommentNotFound:  "评论不存在",
	CodeCommentForbidden: "无权操作评论",
}

// ResponseErrorWithData 带数据的错误响应
func ResponseErrorWithData(c *gin.Context, code int, message string, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    code,
		Message: message,
		Data:    data,
	})
}

// ResponseErrorWithStatus 带HTTP状态码的错误响应
func ResponseErrorWithStatus(c *gin.Context, httpStatus, code int, message string) {
	c.JSON(httpStatus, Response{
		Code:    code,
		Message: message,
		Data:    nil,
	})
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
		Code:    CodeSuccess,
		Message: CodeMessages[CodeSuccess],
		Data:    nil,
	})
}

// ResponseUnauthorized 返回未授权响应
func ResponseUnauthorized(c *gin.Context, message string) {
	c.JSON(http.StatusUnauthorized, Response{
		Code:    CodeUnauthorized,
		Message: message,
	})
}

// ResponseForbidden 返回禁止访问响应
func ResponseForbidden(c *gin.Context, message string) {
	c.JSON(http.StatusForbidden, Response{
		Code:    CodeForbidden,
		Message: message,
	})
}

// ResponseNotFound 返回未找到响应
func ResponseNotFound(c *gin.Context, message string) {
	c.JSON(http.StatusNotFound, Response{
		Code:    CodeNotFound,
		Message: message,
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
	Errorf("System error: %v", err)

	c.JSON(http.StatusInternalServerError, Response{
		Code:    CodeServerError,
		Message: CodeMessages[CodeServerError],
		Data:    nil,
	})
}

// CustomError 返回自定义错误响应
func CustomError(c *gin.Context, httpStatus, code int, message string, err error) {
	if err != nil {
		Errorf("Custom error: %v", err)
		message = err.Error()
	}

	c.JSON(httpStatus, Response{
		Code:    code,
		Message: message,
		Data:    nil,
	})
}

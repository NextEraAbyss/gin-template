package utils

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 统一响应结构
type Response struct {
	Code    int         `json:"code"`    // 错误码
	Message string      `json:"message"` // 错误信息
	Data    interface{} `json:"data"`    // 数据
}

// ErrorCode 错误码类型
type ErrorCode int

// AppResponseError 错误响应对象
type AppResponseError struct {
	Code     ErrorCode   `json:"code"`    // 错误码
	Message  string      `json:"message"` // 错误信息
	Details  string      `json:"details"` // 详细信息，可选
	Data     interface{} `json:"data"`    // 数据，可选
	HTTPCode int         `json:"-"`       // HTTP状态码
	Err      error       `json:"-"`       // 原始错误
}

// Error 实现error接口
func (e *AppResponseError) Error() string {
	if e == nil {
		return ""
	}

	if e.Message != "" {
		return e.Message
	}

	if customErr, ok := e.Err.(*AppResponseError); ok {
		return customErr.Message
	}

	return e.Err.Error()
}

// NewAppResponseError 创建新的错误响应
func NewAppResponseError(code ErrorCode, message string) *AppResponseError {
	return &AppResponseError{
		Code:     code,
		Message:  message,
		HTTPCode: GetHTTPStatusCode(code),
	}
}

// WithDetails 添加详细错误信息
func (e *AppResponseError) WithDetails(details string) *AppResponseError {
	e.Details = details
	return e
}

// WithData 添加错误相关数据
func (e *AppResponseError) WithData(data interface{}) *AppResponseError {
	e.Data = data
	return e
}

// WrapAppResponseError 包装已有错误为响应错误
func WrapAppResponseError(err error, code ErrorCode, message string) *AppResponseError {
	if err == nil {
		return nil
	}

	// 如果已经是响应错误，则更新信息
	var customErr *AppResponseError
	if errors.As(err, &customErr) {
		customErr.Code = code
		customErr.Message = message
		customErr.HTTPCode = GetHTTPStatusCode(code)
		return customErr
	}

	// 否则创建新的响应错误
	return &AppResponseError{
		Code:     code,
		Message:  message,
		Details:  err.Error(),
		HTTPCode: GetHTTPStatusCode(code),
		Err:      err,
	}
}

// ToResponse 将AppResponseError转换为Response
func (e *AppResponseError) ToResponse() Response {
	return Response{
		Code:    int(e.Code),
		Message: e.Message,
		Data:    e.Data,
	}
}

// 系统级状态码
const (
	CodeSuccess       ErrorCode = 0    // 成功
	CodeInvalidParams ErrorCode = 1001 // 无效的参数
	CodeUnauthorized  ErrorCode = 1002 // 未授权
	CodeForbidden     ErrorCode = 1003 // 禁止访问
	CodeNotFound      ErrorCode = 1004 // 资源不存在
	CodeInternalError ErrorCode = 1005 // 内部错误
	CodeServerError   ErrorCode = 1006 // 服务器错误

	// CodeUserNotFound 用户相关错误 (2000-2999)
	CodeUserNotFound  ErrorCode = 2001 // 用户不存在
	CodeUserExists    ErrorCode = 2002 // 用户已存在
	CodePasswordError ErrorCode = 2003 // 密码错误
	CodeTokenExpired  ErrorCode = 2004 // Token过期
	CodeTokenInvalid  ErrorCode = 2005 // Token无效
	CodeUserDisabled  ErrorCode = 2006 // 用户已禁用

	// CodeCommentNotFound 评论相关错误 (4000-4999)
	CodeCommentNotFound  ErrorCode = 4001 // 评论不存在
	CodeCommentForbidden ErrorCode = 4002 // 无权操作评论
)

// 预定义错误响应
var (
	// ErrUnknown 系统级错误
	ErrUnknown       = NewAppResponseError(CodeInternalError, "未知错误")
	ErrInternal      = NewAppResponseError(CodeInternalError, "系统内部错误")
	ErrInvalidParams = NewAppResponseError(CodeInvalidParams, "无效的参数")

	// ErrUnauthorized 用户认证错误
	ErrUnauthorized = NewAppResponseError(CodeUnauthorized, "未授权访问")
	ErrForbidden    = NewAppResponseError(CodeForbidden, "禁止访问")
	ErrTokenExpired = NewAppResponseError(CodeTokenExpired, "Token已过期")
	ErrTokenInvalid = NewAppResponseError(CodeTokenInvalid, "Token无效")

	// ErrUserNotFound 用户操作错误
	ErrUserNotFound  = NewAppResponseError(CodeUserNotFound, "用户不存在")
	ErrUserDisabled  = NewAppResponseError(CodeUserDisabled, "用户已被禁用")
	ErrUserExists    = NewAppResponseError(CodeUserExists, "用户已存在")
	ErrPasswordError = NewAppResponseError(CodePasswordError, "密码错误")
)

// CodeMessages 状态码对应的消息
var CodeMessages = map[ErrorCode]string{
	CodeSuccess:       "操作成功",
	CodeUnauthorized:  "未授权",
	CodeForbidden:     "禁止访问",
	CodeNotFound:      "资源不存在",
	CodeInvalidParams: "请求参数错误",
	CodeInternalError: "服务器内部错误",
	CodeServerError:   "服务器错误",

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

// GetHTTPStatusCode 根据错误码获取对应的HTTP状态码
func GetHTTPStatusCode(code ErrorCode) int {
	switch code {
	case CodeSuccess:
		return http.StatusOK
	case CodeInvalidParams:
		return http.StatusBadRequest
	case CodeUnauthorized:
		return http.StatusUnauthorized
	case CodeForbidden:
		return http.StatusForbidden
	case CodeNotFound:
		return http.StatusNotFound
	case CodeInternalError, CodeServerError:
		return http.StatusInternalServerError
	default:
		// 用户相关错误默认使用400状态码
		if code >= 2000 && code < 3000 {
			return http.StatusBadRequest
		}
		return http.StatusInternalServerError
	}
}

// ResponseSuccess 成功响应 - 使用默认成功消息返回数据
func ResponseSuccess(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    int(CodeSuccess),
		Message: CodeMessages[CodeSuccess],
		Data:    data,
	})
}

// ResponseWithMessage 成功响应 - 使用自定义消息返回数据
func ResponseWithMessage(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    int(CodeSuccess),
		Message: message,
		Data:    data,
	})
}

// ResponseWithAppError 响应应用错误对象
func ResponseWithAppError(c *gin.Context, respError *AppResponseError) {
	c.JSON(respError.HTTPCode, respError.ToResponse())
}

// ResponseError 错误响应 - 返回错误码和消息
func ResponseError(c *gin.Context, code ErrorCode, message string) {
	// 如果没有提供自定义消息，使用默认错误消息
	if message == "" {
		message = CodeMessages[code]
	}

	// 获取HTTP状态码
	httpCode := GetHTTPStatusCode(code)

	c.JSON(httpCode, Response{
		Code:    int(code),
		Message: message,
		Data:    nil,
	})
}

// ResponseErrorWithData 错误响应 - 返回错误信息和数据
func ResponseErrorWithData(c *gin.Context, code ErrorCode, message string, data interface{}) {
	if message == "" {
		message = CodeMessages[code]
	}

	httpCode := GetHTTPStatusCode(code)

	c.JSON(httpCode, Response{
		Code:    int(code),
		Message: message,
		Data:    data,
	})
}

// LogAndResponseError 记录错误并返回错误响应
func LogAndResponseError(c *gin.Context, code ErrorCode, err error) {
	if err != nil {
		Errorf("Error: %v", err)
		ResponseError(c, code, err.Error())
		return
	}
	ResponseError(c, code, "")
}

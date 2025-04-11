package utils

import (
	"errors"
	"fmt"
	"net/http"
)

// AppError 应用程序错误
type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Err     error  `json:"-"`
}

// Error 实现error接口
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

// Unwrap 实现errors.Unwrap接口
func (e *AppError) Unwrap() error {
	return e.Err
}

// NewAppError 创建新的应用程序错误
func NewAppError(code int, message string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// 预定义错误类型
var (
	// ErrBadRequest 请求参数错误
	ErrBadRequest = func(err error) *AppError {
		return NewAppError(http.StatusBadRequest, "请求参数错误", err)
	}

	// ErrUnauthorized 未授权
	ErrUnauthorized = func(err error) *AppError {
		return NewAppError(http.StatusUnauthorized, "未授权", err)
	}

	// ErrForbidden 禁止访问
	ErrForbidden = func(err error) *AppError {
		return NewAppError(http.StatusForbidden, "禁止访问", err)
	}

	// ErrNotFound 资源不存在
	ErrNotFound = func(err error) *AppError {
		return NewAppError(http.StatusNotFound, "资源不存在", err)
	}

	// ErrInternalServer 服务器内部错误
	ErrInternalServer = func(err error) *AppError {
		return NewAppError(http.StatusInternalServerError, "服务器内部错误", err)
	}

	// ErrValidation 验证错误
	ErrValidation = func(err error) *AppError {
		return NewAppError(http.StatusBadRequest, "验证错误", err)
	}

	// ErrDatabase 数据库错误
	ErrDatabase = func(err error) *AppError {
		return NewAppError(http.StatusInternalServerError, "数据库错误", err)
	}

	// ErrCache 缓存错误
	ErrCache = func(err error) *AppError {
		return NewAppError(http.StatusInternalServerError, "缓存错误", err)
	}

	// ErrDuplicate 重复数据错误
	ErrDuplicate = func(err error) *AppError {
		return NewAppError(http.StatusConflict, "数据已存在", err)
	}
)

// IsAppError 检查错误是否为AppError类型
func IsAppError(err error) (*AppError, bool) {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr, true
	}
	return nil, false
}

// GetAppError 获取AppError，如果不是则包装为AppError
func GetAppError(err error) *AppError {
	if appErr, ok := IsAppError(err); ok {
		return appErr
	}
	return ErrInternalServer(err)
}

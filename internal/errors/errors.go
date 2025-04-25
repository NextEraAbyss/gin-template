package errors

import (
	"net/http"
)

// ErrorCode 错误码类型.
type ErrorCode int

// Error 自定义错误类型.
type Error struct {
	Code     ErrorCode `json:"code"`    // 错误码.
	Message  string    `json:"message"` // 错误信息.
	Details  string    `json:"details"` // 详细错误信息.
	HTTPCode int       `json:"-"`       // HTTP状态码.
	Err      error     `json:"-"`       // 嵌套错误.
}

// Error 实现error接口.
func (e *Error) Error() string {
	if e == nil {
		return ""
	}

	if e.Message != "" {
		return e.Message
	}

	if customErr, ok := e.Err.(*Error); ok {
		return customErr.Message
	}

	return e.Err.Error()
}

// New 创建新的错误.
func New(code ErrorCode, message string) *Error {
	return &Error{
		Code:     code,
		Message:  message,
		HTTPCode: GetHTTPCode(code),
	}
}

// WithDetails 添加详细错误信息.
func (e *Error) WithDetails(details string) *Error {
	e.Details = details
	return e
}

// Wrap 包装已有错误.
func Wrap(err error, code ErrorCode, message string) *Error {
	if err == nil {
		return nil
	}

	// 如果已经是自定义错误，则更新信息.
	if customErr, ok := err.(*Error); ok {
		customErr.Code = code
		customErr.Message = message
		customErr.HTTPCode = GetHTTPCode(code)
		return customErr
	}

	// 否则创建新的错误.
	return &Error{
		Code:     code,
		Message:  message,
		Details:  err.Error(),
		HTTPCode: GetHTTPCode(code),
		Err:      err,
	}
}

// GetHTTPCode 根据错误码获取对应的HTTP状态码.
func GetHTTPCode(code ErrorCode) int {
	switch {
	case code >= 1000 && code < 2000:
		return http.StatusInternalServerError

	case code >= 2000 && code < 2100:
		return http.StatusBadRequest

	case code >= 2100 && code < 2200:
		return http.StatusUnauthorized

	case code >= 2200 && code < 2300:
		return http.StatusForbidden

	case code >= 2300 && code < 2400:
		return http.StatusNotFound

	case code >= 2400 && code < 2500:
		return http.StatusConflict

	case code >= 2500 && code < 2600:
		return http.StatusTooManyRequests

	default:
		return http.StatusInternalServerError
	}
}

// 系统级错误码 (1000-1999).
const (
	// 通用错误 (1000-1099).
	CodeSuccess            ErrorCode = 1000 // 成功.
	CodeUnknownError       ErrorCode = 1001 // 未知错误.
	CodeInternalError      ErrorCode = 1002 // 内部错误.
	CodeServiceUnavailable ErrorCode = 1003 // 服务不可用.
	CodeNotFound           ErrorCode = 1004 // 资源不存在.
	CodeInvalidParams      ErrorCode = 1005 // 无效参数.
	CodeTimeout            ErrorCode = 1006 // 请求超时.
	CodeTooManyRequests    ErrorCode = 1007 // 请求过于频繁.

	// 数据库错误 (1100-1199).
	CodeDBError            ErrorCode = 1100 // 数据库错误.
	CodeDBConnectionError  ErrorCode = 1101 // 数据库连接错误.
	CodeDBQueryError       ErrorCode = 1102 // 数据库查询错误.
	CodeDBTransactionError ErrorCode = 1103 // 数据库事务错误.

	// 缓存错误 (1200-1299).
	CodeCacheError           ErrorCode = 1200 // 缓存错误.
	CodeCacheConnectionError ErrorCode = 1201 // 缓存连接错误.
	CodeCacheKeyNotFound     ErrorCode = 1202 // 缓存键不存在.

	// 文件操作错误 (1300-1399).
	CodeFileError            ErrorCode = 1300 // 文件错误.
	CodeFileNotFound         ErrorCode = 1301 // 文件不存在.
	CodeFilePermissionDenied ErrorCode = 1302 // 文件权限拒绝.
	CodeFileIOError          ErrorCode = 1303 // 文件IO错误.

	// 网络错误 (1400-1499).
	CodeNetworkError           ErrorCode = 1400 // 网络错误.
	CodeNetworkTimeout         ErrorCode = 1401 // 网络超时.
	CodeNetworkConnectionError ErrorCode = 1402 // 网络连接错误.

	// 第三方服务错误 (1500-1599).
	CodeThirdPartyError       ErrorCode = 1500 // 第三方服务错误.
	CodeThirdPartyTimeout     ErrorCode = 1501 // 第三方服务超时.
	CodeThirdPartyUnavailable ErrorCode = 1502 // 第三方服务不可用.
)

// 用户相关错误码 (2000-2999).
const (
	// 用户认证错误 (2000-2099).
	CodeUnauthorized ErrorCode = 2000 // 未授权.
	CodeForbidden    ErrorCode = 2001 // 禁止访问.
	CodeTokenExpired ErrorCode = 2002 // Token过期.
	CodeTokenInvalid ErrorCode = 2003 // Token无效.
	CodeTokenMissing ErrorCode = 2004 // Token缺失.

	// 用户操作错误 (2100-2199).
	CodeUserNotFound     ErrorCode = 2100 // 用户不存在.
	CodeUserDisabled     ErrorCode = 2101 // 用户已禁用.
	CodeUserExists       ErrorCode = 2102 // 用户已存在.
	CodeUsernameExists   ErrorCode = 2103 // 用户名已存在.
	CodeEmailExists      ErrorCode = 2104 // 邮箱已存在.
	CodePasswordError    ErrorCode = 2105 // 密码错误.
	CodePasswordTooWeak  ErrorCode = 2106 // 密码太弱.
	CodePasswordMismatch ErrorCode = 2107 // 密码不匹配.
	CodePasswordExpired  ErrorCode = 2108 // 密码已过期.
	CodeAccountLocked    ErrorCode = 2109 // 账户已锁定.
	CodeAccountExpired   ErrorCode = 2110 // 账户已过期.
	CodeLoginFailed      ErrorCode = 2111 // 登录失败.
	CodeLogoutFailed     ErrorCode = 2112 // 登出失败.
	CodeRegisterFailed   ErrorCode = 2113 // 注册失败.
	CodeUpdateFailed     ErrorCode = 2114 // 更新失败.
	CodeDeleteFailed     ErrorCode = 2115 // 删除失败.
)

// 文章相关错误码 (3000-3999).
const (
	// 文章操作错误 (3000-3099).
	CodeArticleNotFound         ErrorCode = 3000 // 文章不存在.
	CodeArticleExists           ErrorCode = 3001 // 文章已存在.
	CodeArticleCreateFailed     ErrorCode = 3002 // 文章创建失败.
	CodeArticleUpdateFailed     ErrorCode = 3003 // 文章更新失败.
	CodeArticleDeleteFailed     ErrorCode = 3004 // 文章删除失败.
	CodeArticlePermissionDenied ErrorCode = 3005 // 文章权限拒绝.
)

// 预定义错误.
var (
	// 系统级错误.
	ErrUnknown            = New(CodeUnknownError, "未知错误")
	ErrInternal           = New(CodeInternalError, "系统内部错误")
	ErrServiceUnavailable = New(CodeServiceUnavailable, "服务暂时不可用")
	ErrTimeout            = New(CodeTimeout, "请求超时")
	ErrInvalidParams      = New(CodeInvalidParams, "无效的参数")

	// 数据库错误.
	ErrDB            = New(CodeDBError, "数据库错误")
	ErrDBConnection  = New(CodeDBConnectionError, "数据库连接错误")
	ErrDBQuery       = New(CodeDBQueryError, "数据库查询错误")
	ErrDBTransaction = New(CodeDBTransactionError, "数据库事务错误")

	// 缓存错误.
	ErrCache            = New(CodeCacheError, "缓存错误")
	ErrCacheConnection  = New(CodeCacheConnectionError, "缓存连接错误")
	ErrCacheKeyNotFound = New(CodeCacheKeyNotFound, "缓存键不存在")

	// 文件错误.
	ErrFile                 = New(CodeFileError, "文件错误")
	ErrFileNotFound         = New(CodeFileNotFound, "文件不存在")
	ErrFilePermissionDenied = New(CodeFilePermissionDenied, "文件权限拒绝")
	ErrFileIO               = New(CodeFileIOError, "文件IO错误")

	// 网络错误.
	ErrNetwork           = New(CodeNetworkError, "网络错误")
	ErrNetworkTimeout    = New(CodeNetworkTimeout, "网络超时")
	ErrNetworkConnection = New(CodeNetworkConnectionError, "网络连接错误")

	// 第三方服务错误.
	ErrThirdParty            = New(CodeThirdPartyError, "第三方服务错误")
	ErrThirdPartyTimeout     = New(CodeThirdPartyTimeout, "第三方服务超时")
	ErrThirdPartyUnavailable = New(CodeThirdPartyUnavailable, "第三方服务不可用")

	// 用户认证错误.
	ErrUnauthorized = New(CodeUnauthorized, "未授权访问")
	ErrForbidden    = New(CodeForbidden, "禁止访问")
	ErrTokenExpired = New(CodeTokenExpired, "Token已过期")
	ErrTokenInvalid = New(CodeTokenInvalid, "Token无效")
	ErrTokenMissing = New(CodeTokenMissing, "Token缺失")

	// 用户操作错误.
	ErrUserNotFound     = New(CodeUserNotFound, "用户不存在")
	ErrUserDisabled     = New(CodeUserDisabled, "用户已被禁用")
	ErrUserExists       = New(CodeUserExists, "用户已存在")
	ErrUsernameExists   = New(CodeUsernameExists, "用户名已存在")
	ErrEmailExists      = New(CodeEmailExists, "邮箱已存在")
	ErrPasswordError    = New(CodePasswordError, "密码错误")
	ErrPasswordTooWeak  = New(CodePasswordTooWeak, "密码强度不足")
	ErrPasswordMismatch = New(CodePasswordMismatch, "密码不匹配")
	ErrPasswordExpired  = New(CodePasswordExpired, "密码已过期")
	ErrAccountLocked    = New(CodeAccountLocked, "账户已锁定")
	ErrAccountExpired   = New(CodeAccountExpired, "账户已过期")
	ErrLoginFailed      = New(CodeLoginFailed, "登录失败")
	ErrLogoutFailed     = New(CodeLogoutFailed, "登出失败")
	ErrRegisterFailed   = New(CodeRegisterFailed, "注册失败")
	ErrUpdateFailed     = New(CodeUpdateFailed, "更新失败")
	ErrDeleteFailed     = New(CodeDeleteFailed, "删除失败")

	// 文章操作错误.
	ErrArticleNotFound         = New(CodeArticleNotFound, "文章不存在")
	ErrArticleExists           = New(CodeArticleExists, "文章已存在")
	ErrArticleCreateFailed     = New(CodeArticleCreateFailed, "文章创建失败")
	ErrArticleUpdateFailed     = New(CodeArticleUpdateFailed, "文章更新失败")
	ErrArticleDeleteFailed     = New(CodeArticleDeleteFailed, "文章删除失败")
	ErrArticlePermissionDenied = New(CodeArticlePermissionDenied, "没有文章操作权限")
)

func (e *Error) StatusCode() int {
	if e == nil {
		return http.StatusInternalServerError
	}

	switch e.Code {
	case CodeInvalidParams:
		return http.StatusBadRequest

	case CodeUnauthorized:
		return http.StatusUnauthorized

	case CodeForbidden:
		return http.StatusForbidden

	case CodeNotFound:
		return http.StatusNotFound

	default:
		return http.StatusInternalServerError
	}
}

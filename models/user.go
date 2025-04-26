package models

import (
	"time"
)

// User 用户模型
type User struct {
	ID          uint       `gorm:"primarykey" json:"id"`                                          // 用户ID
	CreatedAt   time.Time  `gorm:"type:datetime;not null" json:"created_at"`                      // 创建时间
	UpdatedAt   time.Time  `gorm:"type:datetime;not null" json:"updated_at"`                      // 更新时间
	DeletedAt   *time.Time `gorm:"index" json:"deleted_at"`                                       // 删除时间
	Username    string     `gorm:"type:varchar(50);uniqueIndex;not null" json:"username"`         // 用户名
	Password    string     `gorm:"type:varchar(255);not null" json:"-"`                           // 密码
	Email       string     `gorm:"type:varchar(100);uniqueIndex;not null" json:"email"`           // 邮箱
	Nickname    string     `gorm:"type:varchar(50)" json:"nickname"`                              // 昵称
	Avatar      string     `gorm:"type:varchar(255)" json:"avatar"`                               // 头像URL
	Status      int        `gorm:"type:tinyint;default:1;index" json:"status"`                    // 用户状态
	LastLoginAt time.Time  `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"last_login_at"` // 最后登录时间
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}

// UserLoginRequest 用户登录请求
type UserLoginRequest struct {
	Username string `json:"username" binding:"required"` // 用户名
	Password string `json:"password" binding:"required"` // 密码
}

// UserRegisterRequest 用户注册请求
type UserRegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"` // 用户名
	Password string `json:"password" binding:"required,min=6,max=20"` // 密码
	Email    string `json:"email" binding:"required,email"`           // 邮箱
	Nickname string `json:"nickname" binding:"required,min=2,max=50"` // 昵称
}

// UserUpdateRequest 用户信息更新请求
type UserUpdateRequest struct {
	Nickname string `json:"nickname" binding:"omitempty,min=2,max=50"` // 昵称
	Avatar   string `json:"avatar" binding:"omitempty,url"`            // 头像URL
	Email    string `json:"email" binding:"omitempty,email"`           // 邮箱
}

// UserResponse 用户信息响应
type UserResponse struct {
	ID          uint      `json:"id"`            // 用户ID
	Username    string    `json:"username"`      // 用户名
	Email       string    `json:"email"`         // 邮箱
	Nickname    string    `json:"nickname"`      // 昵称
	Avatar      string    `json:"avatar"`        // 头像URL
	Status      int       `json:"status"`        // 用户状态
	LastLoginAt time.Time `json:"last_login_at"` // 最后登录时间
	CreatedAt   time.Time `json:"created_at"`    // 创建时间
	UpdatedAt   time.Time `json:"updated_at"`    // 更新时间
}

// ToResponse 将User模型转换为UserResponse
func (u *User) ToResponse() *UserResponse {
	return &UserResponse{
		ID:          u.ID,
		Username:    u.Username,
		Email:       u.Email,
		Nickname:    u.Nickname,
		Avatar:      u.Avatar,
		Status:      u.Status,
		LastLoginAt: u.LastLoginAt,
		CreatedAt:   u.CreatedAt,
		UpdatedAt:   u.UpdatedAt,
	}
}

// UserQueryDTO 用户查询参数
type UserQueryDTO struct {
	Page     int    `form:"page" json:"page"`         // 页码
	PageSize int    `form:"pageSize" json:"pageSize"` // 每页数量
	Keyword  string `form:"keyword" json:"keyword"`   // 搜索关键词
	Status   int    `form:"status" json:"status"`     // 状态
	OrderBy  string `form:"orderBy" json:"orderBy"`   // 排序字段
	Order    string `form:"order" json:"order"`       // 排序方向
}

// ChangePasswordRequest 修改密码请求
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"` // 旧密码
	NewPassword string `json:"new_password" binding:"required"` // 新密码
}

// ResetPasswordRequest 重置密码请求
type ResetPasswordRequest struct {
	Email string `json:"email" binding:"required,email"` // 邮箱地址
}

// UserListResponse 用户列表响应
type UserListResponse struct {
	Total int64  `json:"total"` // 总数
	Items []User `json:"items"` // 用户列表
}

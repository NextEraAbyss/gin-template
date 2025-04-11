package models

import (
	"time"

	"gorm.io/gorm"
)

// User 用户模型
// 包含用户的基本信息、认证信息和时间戳
type User struct {
	gorm.Model
	Username    string    `gorm:"type:varchar(50);uniqueIndex;not null" json:"username"` // 用户名，唯一
	Password    string    `gorm:"type:varchar(255);not null" json:"-"`                   // 密码，不返回给前端
	Email       string    `gorm:"type:varchar(100);uniqueIndex;not null" json:"email"`   // 邮箱，唯一
	Nickname    string    `gorm:"type:varchar(50)" json:"nickname"`                      // 昵称
	Avatar      string    `gorm:"type:varchar(255)" json:"avatar"`                       // 头像URL
	Status      int       `gorm:"type:tinyint;default:1;index" json:"status"`            // 用户状态：1-正常，0-禁用
	LastLoginAt time.Time `gorm:"type:timestamp;index" json:"last_login_at"`             // 最后登录时间
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
	Username string `json:"username" binding:"required,min=3,max=50"` // 用户名，3-50字符
	Password string `json:"password" binding:"required,min=6,max=20"` // 密码，6-20字符
	Email    string `json:"email" binding:"required,email"`           // 邮箱，必须是有效的邮箱格式
	Nickname string `json:"nickname" binding:"required,min=2,max=50"` // 昵称，2-50字符
}

// UserUpdateRequest 用户信息更新请求
type UserUpdateRequest struct {
	Nickname string `json:"nickname" binding:"omitempty,min=2,max=50"` // 昵称，可选，2-50字符
	Avatar   string `json:"avatar" binding:"omitempty,url"`            // 头像URL，可选，必须是有效的URL
	Email    string `json:"email" binding:"omitempty,email"`           // 邮箱，可选，必须是有效的邮箱格式
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
	Page     int    `form:"page" json:"page"`
	PageSize int    `form:"pageSize" json:"pageSize"`
	Keyword  string `form:"keyword" json:"keyword"`
	Status   int    `form:"status" json:"status"`
	OrderBy  string `form:"orderBy" json:"orderBy"`
	Order    string `form:"order" json:"order"`
}

// UserListResponse 用户列表响应
type UserListResponse struct {
	Total int64  `json:"total"`
	Items []User `json:"items"`
}

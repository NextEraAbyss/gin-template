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

package models

import "time"

// User 用户模型
type User struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Username  string    `json:"username" gorm:"size:50;not null;unique"`
	Email     string    `json:"email" gorm:"size:100;not null;unique"`
	Password  string    `json:"-" gorm:"size:100;not null"`
	FullName  string    `json:"full_name" gorm:"size:100"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

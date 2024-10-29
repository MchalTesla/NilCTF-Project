package models

import (
	"gorm.io/gorm"
)

type User struct {
    gorm.Model
    Username string `gorm:"unique;not null"`
    Password string `gorm:"not null"`
    Email    string `gorm:"unique;not null"`
	// enum: active, inactive
    Status   string `gorm:"default:'active;not null'"` // 用户状态，例如：'active', 'inactive'
	// enum: admin, user
    Role     string `gorm:"default:'user';not null"` // 用户角色
}
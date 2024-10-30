package models

import (
	"gorm.io/gorm"
)

type User struct {
    gorm.Model
    Username string `gorm:"unique;not null"`
    Password string `gorm:"not null"`
	Description string `gorm:"type:text"`        // 用户描述，长度小于
    Email    string `gorm:"unique;not null"`
	// enum: active, inactive
    Status   string `gorm:"default:'active';not null"` // 用户状态，例如：'active', 'inactive'
	// enum: admin, user, organizer
    Role     string `gorm:"default:'user';not null"` // 用户角色
	Tag		 string		// 标签
}
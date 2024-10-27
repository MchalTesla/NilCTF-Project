package models

import (
	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	Username string `gorm:"unique;not null" json:"username"`
	Password string `gorm:"notnull" json:"password"`
	Email    string `gorm:"unique;not null" json:"email"` // 添加邮箱字段，确保唯一且不为null
	Status      string `gorm:"default:'active'"` // 组的状态（例如：'active', 'inactive'）
}
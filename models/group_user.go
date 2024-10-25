package models

import (
	"github.com/jinzhu/gorm"
)

type GroupUser struct {
	gorm.Model
	UserID uint `gorm:"not null"`   // 用户ID
	GroupID uint `gorm:"not null"`   // 组ID
}
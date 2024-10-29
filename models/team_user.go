package models

import (
	"gorm.io/gorm"
)

type TeamUser struct {
	gorm.Model
	UserID uint `gorm:"not null"` // 用户ID
	TeamID uint `gorm:"not null"` // 组ID
}

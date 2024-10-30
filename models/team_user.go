package models

import (
	"gorm.io/gorm"
)

type TeamUser struct {
	gorm.Model
	UserID uint `gorm:"not null;uniqueIndex:idx_user_team"` // 用户ID
	TeamID uint `gorm:"not null;uniqueIndex:idx_user_team"` // 组ID
	// enum: leader, member
	Role   string `gorm:"default:'member';not null"`        // 用户在团队里的角色
}
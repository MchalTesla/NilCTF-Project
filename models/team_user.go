package models

import (
	"gorm.io/gorm"
)

type TeamUser struct {
	gorm.Model
	UserID uint   `gorm:"not null;uniqueIndex:idx_user_team"` // 用户ID
	User   User   `gorm:"foreignKey:UserID;references:ID"`    // 外键关联到 User 表
	TeamID uint   `gorm:"not null;uniqueIndex:idx_user_team"` // 组ID
	Team   Team   `gorm:"foreignKey:TeamID;references:ID"`    // 外键关联到 Team 表
	Role   string `gorm:"default:'member';not null"`          // 用户在团队里的角色 (enum: leader, member)
}
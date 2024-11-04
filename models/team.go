package models

import (
	"gorm.io/gorm"
)

type Team struct {
	gorm.Model
	Name        string   `gorm:"not null"`                             // 队伍名
	Description string   `gorm:"type:text"`                            // 队伍描述
	OwnerID     uint     `gorm:"not null"`                             // 创建者用户ID
	Owner       User     `gorm:"foreignKey:OwnerID;references:ID"`     // 外键关联到 User 表
	Status      string   `gorm:"default:'pending'"`                    // 队伍状态 (enum: active, banned, pending)
	Tag         []string `gorm:"type:jsonb"`    // 标签
	Members     []User   `gorm:"many2many:team_users;"`                // 多对多关联，通过 TeamUser 关联
	Competitions []CompetitionTeam `gorm:"foreignKey:TeamID"` // 关联参加的比赛
}

package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username    string   `gorm:"unique;not null"`
	Password    string   `gorm:"not null"`
	Description *string  `gorm:"type:text"`                      // 用户描述，长度小于150
	Email       string   `gorm:"unique;not null"`
	Status      string   `gorm:"default:'pending';not null"`     // 用户状态 (enum: active, banned, pending)
	Role        string   `gorm:"default:'user';not null"`        // 用户角色 (enum: admin, user, organizer)
	Tag         []string `gorm:"type:jsonb"`                    // 标签
	Teams       []Team   `gorm:"many2many:team_users;"`          // 多对多关联，通过 TeamUser 关联
	OwnedTeams  []Team   `gorm:"foreignKey:OwnerID;references:ID"` // 关联创建的队伍
	Competitions []Competition `gorm:"foreignKey:OwnerID;references:ID"` // 关联创建的比赛
}
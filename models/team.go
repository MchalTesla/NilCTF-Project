package models

import (
	"gorm.io/gorm"
)

type Team struct {
	gorm.Model
	Name        string `gorm:"not null"`         // 队伍名
	Description string `gorm:"type:text"`        // 队伍描述
	OwnerID     uint   `gorm:"not null"`         // 创建者用户ID
	// enum: active, inactive
	Status      string `gorm:"default:'active'"` // 队伍的状态（例如：'active', 'inactive'）
	Tag			string		// 标签
}

package models

import (
	"gorm.io/gorm"
)

type Team struct {
	gorm.Model
	Name        string `gorm:"not null"`         // 组名
	Description string `gorm:"type:text"`        // 组描述
	OwnerID     uint   `gorm:"not null"`         // 创建者用户ID
	// enum: active, inactive
	Status      string `gorm:"default:'active'"` // 组的状态（例如：'active', 'inactive'）
	Tag			string		// 标签
}

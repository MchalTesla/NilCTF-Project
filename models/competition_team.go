package models

import (
	"gorm.io/gorm"
)

type CompetitionTeam struct {
	gorm.Model
	TeamID        uint `gorm:"not null"` // 组ID
	CompetitionID uint `gorm:"not null"` // 比赛ID
	// enum: active inactive underobservation disqualified pending
	tatus        string `gorm:"default:'active';not null"` // 状态，默认值为 'active'
}

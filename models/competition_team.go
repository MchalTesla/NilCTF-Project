package models

import (
	"gorm.io/gorm"
)

type CompetitionTeam struct {
	gorm.Model
	TeamID        uint   `gorm:"not null;uniqueIndex:idx_team_competition"` // 组ID
	CompetitionID uint   `gorm:"not null;uniqueIndex:idx_team_competition"` // 比赛ID
	// enum: active inactive underobservation disqualified pending
	Status        string `gorm:"default:'active';not null"`                 // 状态，默认值为 'active'
}

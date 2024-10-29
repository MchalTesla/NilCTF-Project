package models

import (
	"gorm.io/gorm"
)

type CompetitionTeam struct {
	gorm.Model
	TeamID        uint `gorm:"not null"` // 组ID
	CompetitionID uint `gorm:"not null"` // 比赛ID
}

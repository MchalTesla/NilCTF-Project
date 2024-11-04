package models

import (
	"gorm.io/gorm"
)

type CompetitionTeam struct {
	gorm.Model
	TeamID        uint       `gorm:"not null;uniqueIndex:idx_team_competition"` // 组ID
	Team          Team       `gorm:"foreignKey:TeamID;references:ID"`           // 外键关联到 Team 表
	CompetitionID uint       `gorm:"not null;uniqueIndex:idx_team_competition"` // 比赛ID
	Competition   Competition `gorm:"foreignKey:CompetitionID;references:ID"`   // 外键关联到 Competition 表
	Status        string     `gorm:"default:'active';not null"`                 // 状态 (enum: active, inactive, underobservation, disqualified, pending)
}
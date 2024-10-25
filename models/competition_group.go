package models

import (
	"github.com/jinzhu/gorm"
)

type CompetitionUser struct {
	gorm.Model
	GroupID       uint `gorm:"not null"` // 组ID
	CompetitionID uint `gorm:"not null"` // 比赛ID
}
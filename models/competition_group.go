package models

import (
	"github.com/jinzhu/gorm"
)

type CompetitionGroup struct {
	gorm.Model
	UserID        uint `gorm:"not null"` // 用户ID
	CompetitionID uint `gorm:"not null"` // 比赛ID
}

package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Competitions struct {
	gorm.Model
	Name              string    `gorm:"not null"`           // 比赛名称
	Description       string    `gorm:"type:text"`          // 比赛描述
	StartTime         time.Time `gorm:"not null"`           // 比赛开始时间
	EndTime           time.Time `gorm:"not null"`           // 比赛结束时间
	Status            string    `gorm:"default:'upcoming'"` // 比赛状态（例如：'upcoming', 'ongoing', 'completed'）
	OwnerID           uint      `gorm:"not null"`           // 主办方用户ID
	TeamLimit         int       `gorm:"default:0"`          // 参赛队伍/人数限制，0表示不限制
	ParticipationType string    `gorm:"not null"`           // 参赛形式（例如：'individual'，'team'）
	Nature            string    `gorm:"not null"`           // 比赛性质（例如：'qualifying'，'semifinal'，'final'，'test'）
}

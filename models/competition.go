package models

import (
	"time"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Competition struct {
	gorm.Model
	Name          string           `gorm:"unique;not null"`          // 比赛名称
	Description   string           `gorm:"type:text;default:''"`                // 比赛描述
	StartTime     time.Time        `gorm:"not null"`                 // 比赛开始时间
	EndTime       time.Time        `gorm:"not null"`                 // 比赛结束时间
	Status        string           `gorm:"default:'upcoming'"`       // 比赛状态 (enum: upcoming, ongoing, completed)
	OwnerID       uint             `gorm:"not null"`                 // 主办方用户ID
	Owner         User             `gorm:"foreignKey:OwnerID;references:ID"` // 外键关联到 User 表
	TeamLimit     int              `gorm:"default:0"`                // 参赛队伍数限制，0表示不限制
	Nature        string           `gorm:"default:'test';not null"`  // 比赛性质 (enum: qualifying, semifinal, final, test)
	MaxTeamSize   int              `gorm:"default:0"`                // 每支队伍的最大人数限制，0表示不限制，1表示个人赛
	IsHidden      bool             `gorm:"default:false"`            // 比赛是否隐藏
	Public        bool             `gorm:"default:true"`             // 比赛是否公开
	Suspend       bool             `gorm:"default:false"`            // 比赛是否暂停
	JoinLock      bool             `gorm:"default:false"`            // 比赛是否允许加入
	Tag        	  pq.StringArray   `gorm:"type:text[];default:'{}'"` // 标签
	Teams         []Team		   `gorm:"many2many:competition_teams"` // 多对多关联参加的队伍，通过 CompetitionTeam 关联
}


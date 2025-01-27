package models

import (
	"gorm.io/gorm"
	"github.com/lib/pq"
)

type Question struct {
	gorm.Model
	Title           string         `gorm:"not null"`                  // 题目标题
	Description     string         `gorm:"type:text;default:''"`      // 题目描述
	Difficulty      string         `gorm:"default:'medium'"`          // 难度 (enum: easy, medium, hard)
	Tags            pq.StringArray `gorm:"type:text[];default:'{}'"`  // 标签
	CompetitionID   uint           `gorm:"not null"`                  // 关联比赛ID
	Competition     Competition    `gorm:"foreignKey:CompetitionID;references:ID"` // 外键关联到 Competition 表
	Points        int            `gorm:"default:100"`          		  // 题目分值
	IsHidden        bool           `gorm:"default:false"`             // 是否隐藏题目
	SubmissionCount int            `gorm:"default:0"`                 // 提交次数统计
}
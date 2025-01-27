package models

import (
	"gorm.io/gorm"
)

type Submission struct {
	gorm.Model
	UserID     uint      `gorm:"not null"`                  // 提交用户ID
	User       User      `gorm:"foreignKey:UserID;references:ID"` // 外键关联到 User 表
	QuestionID uint      `gorm:"not null"`                  // 提交题目ID
	Question   Question  `gorm:"foreignKey:QuestionID;references:ID"` // 外键关联到 Question 表
	Score      int       `gorm:"default:0"`                 // 提交得分
	Status     string    `gorm:"default:'pending'"`         // 提交状态 (enum: pending, correct, incorrect, partially_correct)
	Answer     string    `gorm:"type:text"`                 // 用户提交的输出
}

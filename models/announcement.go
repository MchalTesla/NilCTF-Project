package models

import (
    "time"

    "gorm.io/gorm"
)
const (
	PriorityLow    = "low"
	PriorityMedium = "medium"
	PriorityHigh   = "high"
)

type Announcement struct {
	gorm.Model
	Title       string    `gorm:"not null"`                  // 公告标题
	Content     string    `gorm:"type:text;not null"`        // 公告内容
	AuthorID    uint      `gorm:"not null"`                  // 作者用户ID
	Author      User      `gorm:"foreignKey:AuthorID;references:ID"` // 外键关联到 User 表
	PublishedAt time.Time `gorm:"not null"`                  // 发布时间
	Status      string    `gorm:"default:'draft';not null"`  // 公告状态 (enum: draft, published, archived)
	Priority    string    `gorm:"default:'medium';not null"` // 公告优先级 (enum: low, medium, high)
}
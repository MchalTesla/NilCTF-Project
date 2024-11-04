package models

import (
	"gorm.io/gorm"
)

type Config struct {
    gorm.Model
    Key   string `gorm:"unique;not null"` // 配置项的名称
    Value string `gorm:"not null"`         // 配置项的值
}

func InitializeConfigs(db *gorm.DB) error {
	defaultConfigs := []Config{
		{Key: "user_status_default", Value: "pending"},
	}

	for _, config := range defaultConfigs {
		var existingConfig Config
		// 检查配置项是否存在
		result := db.Where("key = ?", config.Key).First(&existingConfig)
		if result.RowsAffected == 0 {
			// 不存在，插入
			if err := db.Create(&config).Error; err != nil {
				return err // 返回错误
			}
		}
	}

	return nil // 初始化成功
}
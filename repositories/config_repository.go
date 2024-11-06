package repositories

import (
	"NilCTF/models"

	"gorm.io/gorm"
)

type ConfigRepository struct {
	DB *gorm.DB
}

// 新建 ConfigRepository
func NewConfigRepository(db *gorm.DB) *ConfigRepository {
	return &ConfigRepository{DB: db}
}

// 更新或插入配置项
func (repo *ConfigRepository) Upsert(key, value string) error {
	config := models.Config{Key: key, Value: value}
	return repo.DB.Where("key = ?", key).Assign(config).FirstOrCreate(&config).Error
}

// 获取配置项
func (repo *ConfigRepository) Get(key string) (string, error) {
	var config models.Config
	if err := repo.DB.Where("key = ?", key).First(&config).Error; err != nil {
		return "", err
	}
	return config.Value, nil
}

// 删除配置项
func (repo *ConfigRepository) Delete(key string) error {
	return repo.DB.Where("key = ?", key).Unscoped().Delete(&models.Config{}).Error
}

// ListConfig 根据条件查找配置
func (repo *ConfigRepository) List(condition, value string) ([]models.Config, error) {
	var configs []models.Config

	// 使用条件查询，假设我们要按状态查找
	if err := repo.DB.Where(condition, value).Find(&configs).Error; err != nil {
		return nil, err // 返回错误
	}

	return configs, nil // 返回查询结果
}

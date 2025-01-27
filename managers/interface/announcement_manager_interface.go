package managers_interface

import (
	"NilCTF/models"
)

// AnnouncementManagerInterface 定义了公告管理器的接口
type AnnouncementManagerInterface interface {
	// Create 创建新的公告
	Create(announcement *models.Announcement) error
	// Get 根据ID获取公告
	Get(ID uint) (*models.Announcement, error)
	// Update 更新公告
	Update(announcement *models.Announcement) error
	// Delete 删除公告
	Delete(announcement *models.Announcement) error
	// List 列出公告，支持过滤、分页和模糊查询
	List(filters map[string]interface{}, limit, offset int, isFuzzy bool) ([]models.Announcement, error)
}

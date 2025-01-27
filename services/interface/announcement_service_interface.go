package services_interface

import (
	"NilCTF/dto"
)

// AnnouncementServiceInterface 定义了公告服务的接口
type AnnouncementServiceInterface interface {
	// Create 创建新的公告
	Create(announcementDTO *dto.AnnouncementDTO) error
	// Get 根据ID获取公告
	Get(ID uint) (*dto.AnnouncementDTO, error)
	// Update 更新公告
	Update(announcementDTO *dto.AnnouncementDTO) error
	// Delete 删除公告
	Delete(ID uint) error
	// List 列出公告，支持过滤、分页和模糊查询
	List(filters map[string]interface{}, limit, offset int, isFuzzy bool) ([]dto.AnnouncementDTO, error)
}

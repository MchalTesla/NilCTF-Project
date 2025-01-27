package managers

import (
    "NilCTF/error_code"
    "NilCTF/models"
    repositories_interface "NilCTF/repositories/interface"
)

type AnnouncementManager struct {
    AR repositories_interface.AnnouncementRepositoryInterface
}

// NewAnnouncementManager 创建一个新的 AnnouncementManager 实例
func NewAnnouncementManager(AR repositories_interface.AnnouncementRepositoryInterface) *AnnouncementManager {
    return &AnnouncementManager{AR: AR}
}

// Create 创建新的公告
func (m *AnnouncementManager) Create(announcement *models.Announcement) error {
    if announcement.ID != 0 {
        return error_code.ErrInvalidID
    }

    // 检查公告标题是否为空
    if announcement.Title == "" {
        return error_code.ErrInvalidInput
    }

    // 检查公告内容是否为空
    if announcement.Content == "" {
        return error_code.ErrInvalidInput
    }

    return m.AR.Create(announcement)
}

// Get 根据ID获取公告
func (m *AnnouncementManager) Get(ID uint) (*models.Announcement, error) {
    if ID == 0 {
        return nil, error_code.ErrInvalidID
    }
    return m.AR.Get(ID)
}

// Update 更新公告
func (m *AnnouncementManager) Update(announcement *models.Announcement) error {
    if announcement.ID == 0 {
        return error_code.ErrInvalidID
    }

    // 检查公告标题是否为空
    if announcement.Title == "" {
        return error_code.ErrInvalidInput
    }

    // 检查公告内容是否为空
    if announcement.Content == "" {
        return error_code.ErrInvalidInput
    }

    return m.AR.Update(announcement)
}

// Delete 删除公告
func (m *AnnouncementManager) Delete(announcement *models.Announcement) error {
    if announcement.ID == 0 {
        return error_code.ErrInvalidID
    }
    return m.AR.Delete(announcement)
}

// List 列出公告，支持过滤、分页和模糊查询
func (m *AnnouncementManager) List(filters map[string]interface{}, limit, offset int, isFuzzy bool) ([]models.Announcement, error) {
    return m.AR.List(filters, limit, offset, isFuzzy)
}
package repositories

import (
    "NilCTF/error_code"
    "NilCTF/models"

    "gorm.io/gorm"
)

// AnnouncementRepository 实现了 AnnouncementRepositoryInterface 接口
type AnnouncementRepository struct {
    DB *gorm.DB
}

// NewAnnouncementRepository 创建一个新的 AnnouncementRepository 实例
func NewAnnouncementRepository(DB *gorm.DB) *AnnouncementRepository {
    return &AnnouncementRepository{DB: DB}
}

// Create 创建新的公告
func (r *AnnouncementRepository) Create(announcement *models.Announcement) error {
    if announcement.ID != 0 {
        return error_code.ErrInvalidID
    }
    if err := r.DB.Create(announcement).Error; err != nil {
        return error_code.ErrInternalServer
    }
    return nil
}

// Get 根据ID获取公告
func (r *AnnouncementRepository) Get(ID uint) (*models.Announcement, error) {
    if ID == 0 {
        return nil, error_code.ErrInvalidID
    }

    var announcement models.Announcement
    if err := r.DB.First(&announcement, ID).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            return nil, error_code.ErrNotFound
        }
        return nil, error_code.ErrInternalServer
    }
    return &announcement, nil
}

// Update 更新公告
func (r *AnnouncementRepository) Update(announcement *models.Announcement) error {
    if announcement.ID == 0 {
        return error_code.ErrInvalidID
    }
    if err := r.DB.Updates(announcement).Error; err != nil {
        return error_code.ErrInternalServer
    }
    return nil
}

// Delete 删除公告
func (r *AnnouncementRepository) Delete(announcement *models.Announcement) error {
    if announcement.ID == 0 {
        return error_code.ErrInvalidID
    }
    if err := r.DB.Unscoped().Delete(announcement).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            return error_code.ErrNotFound
        }
        return error_code.ErrInternalServer
    }
    return nil
}

// List 列出公告，支持过滤、分页和模糊查询
func (r *AnnouncementRepository) List(filters map[string]interface{}, limit, offset int, isFuzzy bool) ([]models.Announcement, error) {
    var announcements []models.Announcement
    query := r.DB

    for key, value := range filters {
        if isFuzzy {
            query = query.Where(key+" LIKE ?", "%"+value.(string)+"%")
        } else {
            query = query.Where(key+" = ?", value)
        }
    }

    if limit > 0 {
        query = query.Limit(limit)
    }
    if offset >= 0 {
        query = query.Offset(offset)
    }

    if err := query.Find(&announcements).Error; err != nil {
        return []models.Announcement{}, error_code.ErrInternalServer
    }
    return announcements, nil
}
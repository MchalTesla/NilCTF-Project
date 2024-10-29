package repositories

import (
	"NilCTF/error_code"
	"NilCTF/models"
	"errors"

	"gorm.io/gorm"
)

type TeamRepository struct {
	DB *gorm.DB
}

// NewTeamRepository 返回新的 TeamRepository 实例
func NewTeamRepository(DB *gorm.DB) *TeamRepository {
	return &TeamRepository{DB: DB}
}

// Create 创建Team
func (r *TeamRepository) Create(team *models.Team) error {
	var existingTeams models.Team

	// 检查团队是否已经存在
	if err := r.DB.Where("name = ?", team.Name).First(&existingTeams).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 创建新团队
			if err := r.DB.Create(team).Error; err != nil {
				// 系统错误处理
				return error_code.ErrInternalServer
			}
			return nil
		}
		// 系统错误处理
		return error_code.ErrInternalServer
	}

	// 团队已存在
	return error_code.ErrTeamAlreadyExists
}

// Read 根据ID或者队伍名查找Team
func (r *TeamRepository) Read(ID uint, name string) ([]models.Team, error) {
	var existingTeams []models.Team

	// 根据ID或名称查找团队
	var err error
	switch {
	case ID != 0:
		err = r.DB.Find(&existingTeams, ID).Error
	case name != "":
		err = r.DB.Where("name = ?", name).Find(&existingTeams).Error
	default:
		return nil, error_code.ErrInvalidInput
	}

	if err != nil {
		// 系统错误处理
		return nil, error_code.ErrInternalServer
	} else if len(existingTeams) == 0{
		return nil, error_code.ErrTeamNotFound
	}

	return existingTeams, nil
}

// Update 更新队伍信息
func (r *TeamRepository) Update(team *models.Team) error {
	// 检查团队ID是否存在
	if team.ID == 0 {
		return error_code.ErrInvalidInput
	}

	// 更新团队信息
	if err := r.DB.Save(team).Error; err != nil {
		// 系统错误处理
		return error_code.ErrInternalServer
	}
	return nil
}

// Delete 删除队伍
func (r *TeamRepository) Delete(team *models.Team) error {
	if err := r.DB.Delete(team).Error; err != nil {
		// 系统错误处理
		return error_code.ErrInternalServer
	}
	return nil
}

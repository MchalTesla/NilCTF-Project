package repositories

import (
	"NilCTF/error_code"
	"NilCTF/models"
	"NilCTF/utils"
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

// Create 创建Team，ID必须为0
func (r *TeamRepository) Create(team *models.Team) error {

	// 判断ID是否有效
	if team.ID != 0 {
		return error_code.ErrInvalidID
	}

	// 检查队伍名中是否符合规范
	if !utils.IsValidName(team.Name) {
		return error_code.ErrInvalidInput
	}

	// 检查队伍描述是否符合规范
	if !utils.IsValidDescription(team.Description) {
		return error_code.ErrInvalidDescription
	}

	// 检查队伍是否已经存在
	if err := r.DB.Where("name = ?", team.Name).First(&models.Team{}).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 创建新队伍
			if err := r.DB.Create(team).Error; err != nil {
				// 系统错误处理
				return error_code.ErrInternalServer
			}
			return nil
		}
		// 系统错误处理
		return error_code.ErrInternalServer
	}

	// 队伍已存在
	return error_code.ErrTeamAlreadyExists
}

// Read 根据ID或者队伍名查找Team
func (r *TeamRepository) Read(ID uint, name string) ([]models.Team, error) {
	var existingTeams []models.Team

	// 根据ID或名称查找队伍
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

// Update 更新队伍信息, ID必须存在
func (r *TeamRepository) Update(team *models.Team) error {
	// 检查队伍ID是否有效
	if team.ID == 0 {
		return error_code.ErrInvalidID
	}

	// 检查队伍ID是否存在
	if err := r.DB.Where("id = ?", team.ID).First(&models.Team{}).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return error_code.ErrTeamNotFound
		}
		return error_code.ErrInternalServer
	}

	// 检查队伍名中是否符合规范
	if !utils.IsValidName(team.Name) {
		return error_code.ErrInvalidInput
	}

	// 检查队伍描述是否符合规范
	if !utils.IsValidDescription(team.Description) {
		return error_code.ErrInvalidDescription
	}

	// 更新队伍信息
	if err := r.DB.Model(team).Updates(team).Error; err != nil {
		// 系统错误处理
		return error_code.ErrInternalServer
	}
	return nil
}

// Delete 删除队伍， ID必须存在
func (r *TeamRepository) Delete(team *models.Team) error {
	// 判断ID是否有效
	if team.ID == 0 {
		return error_code.ErrInvalidID
	}

	if err := r.DB.Delete(team).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return error_code.ErrTeamNotFound
		}
		// 系统错误处理
		return error_code.ErrInternalServer
	}
	return nil
}

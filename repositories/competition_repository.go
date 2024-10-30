package repositories

import (
	"NilCTF/error_code"
	"NilCTF/models"
	"errors"

	"gorm.io/gorm"
)

type CompetitionRepository struct {
	DB *gorm.DB
}

// NewCompetitionRepository 返回新的 CompetitionRepository 实例
func NewCompetitionRepository(DB *gorm.DB) *CompetitionRepository {
	return &CompetitionRepository{DB: DB}
}

// Create 创建Competition， ID必须为0
func (r *CompetitionRepository) Create(competition *models.Competition) error {

	// 判断ID是否合规
	if competition.ID != 0 {
		return error_code.ErrInvalidID
	}

	// 检查比赛是否已存在
	if err := r.DB.Where("name = ?", competition.Name).First(&models.Competition{}).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 创建新比赛
			if err := r.DB.Create(competition).Error; err != nil {
				// 系统错误处理
				return error_code.ErrInternalServer
			}
			return nil
		}
		// 系统错误处理
		return error_code.ErrInternalServer
	}

	// 比赛已存在
	return error_code.ErrCompetitionAlreadyExists
}

// Read 根据ID、名称或所有者ID查找Competition
func (r *CompetitionRepository) Read(ID uint, name string, ownerID uint) ([]models.Competition, error) {
	var ExistingCompetitions []models.Competition

	// 根据ID、名称或所有者ID查找比赛
	var err error
	switch {
	case ID != 0:
		err = r.DB.Find(&ExistingCompetitions, ID).Error
	case name != "":
		err = r.DB.Where("name = ?", name).Find(&ExistingCompetitions).Error
	case ownerID != 0:
		err = r.DB.Where("owner_id = ?", ownerID).Find(&ExistingCompetitions).Error
	default:
		return nil, error_code.ErrInvalidInput
	}

	if err != nil {
		// 系统错误处理
		return nil, error_code.ErrInternalServer
	} else if len(ExistingCompetitions) == 0 {
		return nil, error_code.ErrCompetitionNotFound
	}

	return ExistingCompetitions, nil
}

// Update 更新Competition信息, ID必须存在
func (r *CompetitionRepository) Update(competition *models.Competition) error {
	// 检查比赛ID是否有效
	if competition.ID == 0 {
		return error_code.ErrInvalidID
	}

	// 检查比赛ID是否存在
	if err := r.DB.Where("id = ?", competition.ID).First(&models.Competition{}).Error; err == nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return error_code.ErrUserNotFound
		}
		return error_code.ErrInternalServer
	}

	// 检查比赛名字是否存在
	if err := r.DB.Where("name = ? AND id != ?", competition.Name, competition.ID).First(&models.Competition{}).Error; err == nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return error_code.ErrCompetitionAlreadyExists
		}
		return error_code.ErrInternalServer
	}

	// 更新比赛信息
	if err := r.DB.Model(competition).Updates(competition).Error; err != nil {
		// 系统错误处理
		return error_code.ErrInternalServer
	}
	return nil
}

// Delete 删除Competition， ID必须存在
func (r *CompetitionRepository) Delete(competition *models.Competition) error {
	// 判断ID是否有效
	if competition.ID == 0 {
		return error_code.ErrInvalidID
	}
	
	if err := r.DB.Delete(competition).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return error_code.ErrCompetitionNotFound
		}
		// 系统错误处理
		return error_code.ErrInternalServer
	}
	return nil
}

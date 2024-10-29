package repositories

import (
	"NilCTF/error_code"
	"NilCTF/models"

	"github.com/jinzhu/gorm"
)

type CompetitionRepository struct {
	DB *gorm.DB
}

// NewCompetitionRepository 返回新的 CompetitionRepository 实例
func NewCompetitionRepository(DB *gorm.DB) *CompetitionRepository {
	return &CompetitionRepository{DB: DB}
}

// Create 创建Competition
func (r *CompetitionRepository) Create(competition *models.Competition) error {
	var existingCompetition models.Competition

	// 检查比赛是否已存在
	if err := r.DB.Where("name = ?", competition.Name).First(&existingCompetition).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
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
	var sliceExistingCompetition []models.Competition

	// 根据ID、名称或所有者ID查找比赛
	var err error
	switch {
	case ID != 0:
		err = r.DB.Find(&sliceExistingCompetition, ID).Error
	case name != "":
		err = r.DB.Where("name = ?", name).Find(&sliceExistingCompetition).Error
	case ownerID != 0:
		err = r.DB.Where("owner_id = ?", ownerID).Find(&sliceExistingCompetition).Error
	default:
		return nil, error_code.ErrInvalidInput
	}

	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, error_code.ErrCompetitionNotFound
		}
		// 系统错误处理
		return nil, error_code.ErrInternalServer
	}

	return sliceExistingCompetition, nil
}

// Update 更新Competition信息
func (r *CompetitionRepository) Update(competition *models.Competition) error {
	// 检查比赛ID是否存在
	if competition.ID == 0 {
		return error_code.ErrInvalidInput
	}

	// 更新比赛信息
	if err := r.DB.Save(competition).Error; err != nil {
		// 系统错误处理
		return error_code.ErrInternalServer
	}
	return nil
}

// Delete 删除Competition
func (r *CompetitionRepository) Delete(competition *models.Competition) error {
	if err := r.DB.Delete(competition).Error; err != nil {
		// 系统错误处理
		return error_code.ErrInternalServer
	}
	return nil
}
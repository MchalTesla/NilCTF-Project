package repositories

import (
	"NilCTF/error_code"
	"NilCTF/models"

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

	if err := r.DB.Create(competition).Error; err != nil {
		// 系统错误处理
		return error_code.ErrInternalServer
	}
	return nil
}

// Get 根据ID、名称或所有者ID查找Competition
func (r *CompetitionRepository) Get(ID uint) (*models.Competition, error) {
	var ExistingCompetition models.Competition

	// 根据ID查找比赛
	if ID == 0 {
		return nil, error_code.ErrInvalidID
	}

	if err := r.DB.First(&ExistingCompetition, ID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, error_code.ErrCompetitionNotFound
		} else {
			return nil, error_code.ErrInternalServer
		}
	}

	return &ExistingCompetition, nil
}

// Update 更新Competition信息, ID必须存在
func (r *CompetitionRepository) Update(competition *models.Competition) error {
	// 检查比赛ID是否有效
	if competition.ID == 0 {
		return error_code.ErrInvalidID
	}

	// 更新比赛信息
	if err := r.DB.Updates(competition).Error; err != nil {
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

	if err := r.DB.Unscoped().Delete(competition).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return error_code.ErrCompetitionNotFound
		}
		// 系统错误处理
		return error_code.ErrInternalServer
	}
	return nil
}

func (r *CompetitionRepository) List(filters map[string]interface{}, limit, offset int, isFuzzy bool) ([]models.Competition, error) {
	var competitions []models.Competition
	query := r.DB

	// 应用过滤条件
	for key, value := range filters {
		if isFuzzy { // 如果启用模糊搜索
			// 使用 LIKE 查询并在值两端添加 % 通配符
			query = query.Where(key+" LIKE ?", "%"+value.(string)+"%")
		} else {
			// 精确匹配
			query = query.Where(key+" = ?", value)
		}
	}

	// 设置分页
	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset >= 0 {
		query = query.Offset(offset)
	}

	// 执行查询
	if err := query.Find(&competitions).Error; err != nil {
		return []models.Competition{}, error_code.ErrInternalServer
	}
	return competitions, nil
}

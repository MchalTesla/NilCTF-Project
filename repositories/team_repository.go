package repositories

import (
	"NilCTF/error_code"
	"NilCTF/models"

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
	// 创建新队伍
	if err := r.DB.Create(team).Error; err != nil {
		// 系统错误处理
		return error_code.ErrInternalServer
	}
	return nil
}

// Update 更新队伍信息, ID必须存在
func (r *TeamRepository) Update(team *models.Team) error {
	if team.ID == 0 {
		return error_code.ErrInvalidID
	}
	// 更新队伍信息
	if err := r.DB.Model(team).Updates(team).Error; err != nil {
		// 系统错误处理
		return error_code.ErrInternalServer
	}
	return nil
}

// Get 根据ID或者队伍名查找Team
func (r *TeamRepository) Get(ID uint) (*models.Team, error) {
	var existingTeams models.Team

	// 根据ID查找队伍
	if ID == 0 {
		return nil, error_code.ErrInvalidID
	}

	if err := r.DB.First(&existingTeams, ID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, error_code.ErrTeamNotFound
		} else {
			return nil, error_code.ErrInternalServer
		}
	}

	return &existingTeams, nil
}

// Delete 删除队伍， ID必须存在
func (r *TeamRepository) Delete(team *models.Team) error {
	// 判断ID是否有效
	if team.ID == 0 {
		return error_code.ErrInvalidID
	}

	if err := r.DB.Unscoped().Delete(team).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return error_code.ErrTeamNotFound
		}
		// 系统错误处理
		return error_code.ErrInternalServer
	}
	return nil
}

func (r *TeamRepository) List(filters map[string]interface{}, limit, offset int, isFuzzy bool) ([]models.Team, error) {
	var teams []models.Team
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
	if err := query.Find(&teams).Error; err != nil {
		return []models.Team{}, error_code.ErrInternalServer
	}
	return teams, nil
}

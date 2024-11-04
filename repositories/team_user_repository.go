package repositories

import (
	"NilCTF/error_code"
	"NilCTF/models"

	"gorm.io/gorm"
)

type TeamUserRepository struct {
	DB *gorm.DB
}

func NewTeamUserRepository(DB *gorm.DB) *TeamUserRepository {
	return &TeamUserRepository{DB: DB}
}

// Upsert 插入或更新队伍和用户的映射
func (r *TeamUserRepository) Create(teamUser *models.TeamUser) error {
	if teamUser.ID != 0 {
		return error_code.ErrInvalidID
	}
	if err := r.DB.Create(teamUser).Error; err != nil {
		return error_code.ErrInternalServer
	}
	return nil
}
func (r *TeamUserRepository) Update(teamUser *models.TeamUser) error {
	if teamUser.ID == 0 {
		return error_code.ErrInvalidID
	}
	if err := r.DB.Assign(teamUser).Error; err != nil {
		return error_code.ErrInternalServer
	}
	return nil
}

// Get 根据ID、teamID或userID查找队伍和用户的映射
func (r *TeamUserRepository) Get(ID, teamID, userID uint) ([]models.TeamUser, error) {
	var teamUsers []models.TeamUser
	query := r.DB.Model(&teamUsers)

	// 动态构建查询条件
	switch {
	case ID != 0:
		query = query.Where("id = ?", ID)
	case teamID != 0 && userID != 0:
		query = query.Where("teamid = ? AND userid = ?", teamID, userID)
	case teamID != 0:
		query = query.Where("teamid = ?", teamID)
	case userID != 0:
		query = query.Where("userid = ?", userID)
	default:
		return nil, error_code.ErrInvalidInput
	}

	if err := query.Find(&teamUsers).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, error_code.ErrTeamNotInCompetition
		}
		return nil, error_code.ErrInternalServer
	}

	return teamUsers, nil
}

// Delete 删除队伍和用户的映射
func (r *TeamUserRepository) Delete(teamUser *models.TeamUser) error {
	if err := r.DB.Delete(teamUser).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return error_code.ErrUserNotInTeam
		}
		return error_code.ErrInternalServer
	}
	return nil
}

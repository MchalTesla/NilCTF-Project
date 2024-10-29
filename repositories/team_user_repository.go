package repositories

import (
	"NilCTF/error_code"
	"NilCTF/models"
	"errors"

	"gorm.io/gorm"
)

type TeamUserRepository struct {
	DB *gorm.DB
}

// NewTeamUserRepository 返回新的 TeamUserRepository 实例
func NewTeamUserRepository(DB *gorm.DB) *TeamUserRepository {
	return &TeamUserRepository{DB: DB}
}

// Create 创建队伍和用户的映射
func (r *TeamUserRepository) Create(teamUser *models.TeamUser) error {
	var existingTeamUser models.TeamUser

	if err := r.DB.Where("teamid = ? AND userid = ?", teamUser.TeamID, teamUser.UserID).First(&existingTeamUser).Error; err == nil {
		return error_code.ErrUserAlreadyInTeam
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return error_code.ErrInternalServer
	}

	if err := r.DB.Create(teamUser).Error; err != nil {
		return error_code.ErrInternalServer
	}
	return nil
}

// Read 查找队伍和用户的映射
func (r *TeamUserRepository) Read(ID, teamID, userID uint) ([]models.TeamUser, error) {
	var teamUsers []models.TeamUser

	// 根据ID、TeamID或UserID查找
	var err error
	switch {
	case ID != 0:
		err = r.DB.Find(&teamUsers, ID).Error
	case teamID != 0 && userID != 0:
		err = r.DB.Where("teamid = ? AND userid = ?", teamID, userID).Find(&teamUsers).Error
	case teamID != 0:
		err = r.DB.Where("teamid = ?", teamID).Find(&teamUsers).Error
	case userID != 0:
		err = r.DB.Where("userid = ?", userID).Find(&teamUsers).Error
	default:
		return nil, error_code.ErrInvalidInput
	}

	if err != nil {
		return nil, error_code.ErrInternalServer
	} else if len(teamUsers) == 0 {
		return nil, error_code.ErrTeamNotInCompetition
	}

	return teamUsers, nil
}

// Update 更新队伍和用户的映射
func (r *TeamUserRepository) Update(teamUser *models.TeamUser) error {
	if err := r.DB.Save(teamUser).Error; err != nil {
		return error_code.ErrInternalServer
	}
	return nil
}

// Delete 删除队伍和用户的映射
func (r *TeamUserRepository) Delete(teamUser *models.TeamUser) error {
	if err := r.DB.Delete(teamUser).Error; err != nil {
		return error_code.ErrInternalServer
	}
	return nil
}

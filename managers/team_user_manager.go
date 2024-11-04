package managers

import (
	"NilCTF/models"
	"NilCTF/repositories/interface"
	"NilCTF/error_code"
)

type TeamUserManager struct {
	TUR repositories_interface.TeamUserRepositoryInterface
}

func NewTeamUserManager(TUR repositories_interface.TeamUserRepositoryInterface) *TeamUserManager {
	return &TeamUserManager{TUR: TUR}
}

// Create 创建新的队伍和用户映射
func (m *TeamUserManager) Create(teamUser *models.TeamUser) error {
	// 检查ID是否有效
	if teamUser.ID != 0 {
		return error_code.ErrInvalidID
	}
	
	// 检查是否已有该用户和队伍的映射
	existingMappings, err := m.TUR.Get(0, teamUser.TeamID, teamUser.UserID)
	if err != nil {
		return err
	}
	if len(existingMappings) > 0 {
		return error_code.ErrUserAlreadyInTeam
	}

	return m.TUR.Create(teamUser)
}

// Update 更新已有的队伍和用户映射
func (m *TeamUserManager) Update(teamUser *models.TeamUser) error {
	// 检查ID是否有效
	if teamUser.ID == 0 {
		return error_code.ErrInvalidID
	}

	// 使用Get确认记录存在
	existingMappings, err := m.TUR.Get(teamUser.ID, 0, 0)
	if err != nil {
		return err
	}
	if len(existingMappings) == 0 {
		return error_code.ErrUserNotInTeam
	}

	// 如果 TeamID 或 UserID 变更，拒绝更新
	if (teamUser.TeamID != 0 && existingMappings[0].TeamID != teamUser.TeamID) ||
	(teamUser.UserID != 0 && existingMappings[0].UserID != teamUser.UserID) {
		return error_code.ErrInvalidID
	}

	return m.TUR.Update(teamUser)
}

// Get 根据条件获取队伍和用户映射信息
func (m *TeamUserManager) Get(ID, teamID, userID uint) ([]models.TeamUser, error) {
	return m.TUR.Get(ID, teamID, userID)
}

// Delete 删除指定的队伍和用户映射
func (m *TeamUserManager) Delete(teamUser *models.TeamUser) error {
	if teamUser.ID == 0 {
		return error_code.ErrInvalidID
	}
	return m.TUR.Delete(teamUser)
}

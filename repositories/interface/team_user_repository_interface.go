package repositories_interface

import (
	"NilCTF/models"
)

// 队伍用户对照表Repository层接口
type TeamUserRepositoryInterface interface {
	Create(teamUser *models.TeamUser) error
	Read(ID, teamID, userID uint) ([]models.TeamUser, error)
	Update(teamUser *models.TeamUser) error
	Delete(teamUser *models.TeamUser) error
}
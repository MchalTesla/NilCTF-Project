package managers_interface

import (
	"NilCTF/models"
)

// 队伍用户对照表Manager层接口
type TeamUserManagerInterface interface {
	Create(teamUser *models.TeamUser) error
	Get(ID, teamID, userID uint) ([]models.TeamUser, error)
	Update(teamUser *models.TeamUser) error
	Delete(teamUser *models.TeamUser) error
}
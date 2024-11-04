package managers_interface

import (
	"NilCTF/models"
)

type CompetitionTeamManagerInterface interface {
	Create(competitionTeam *models.CompetitionTeam) error
	Get(ID, competitionID, teamID uint) ([]models.CompetitionTeam, error)
	Update(competitionTeam *models.CompetitionTeam) error
	Delete(teamUser *models.CompetitionTeam) error
}
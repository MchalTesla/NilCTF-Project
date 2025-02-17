package repositories_interface

import (
	"NilCTF/models"
)

type CompetitionTeamRepositoryInterface interface {
	Create(competitionTeam *models.CompetitionTeam) error
	Get(ID, competitionID, teamID uint) ([]models.CompetitionTeam, error)
	Update(competitionTeam *models.CompetitionTeam) error
	Delete(teamUser *models.CompetitionTeam) error
	List(filters map[string]interface{}, limit, offset int, isFuzzy bool) ([]models.CompetitionTeam, error)
}

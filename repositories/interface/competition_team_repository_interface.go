package repositories_interface

import (
	"NilCTF/models"
)

type CompetitionTeamRepositoryInterface interface {
	Create(competitionTeam *models.CompetitionTeam) error
	Read(ID, competitionID, teamID uint) ([]models.CompetitionTeam, error)
	Update(competitionTeam *models.CompetitionTeam) error
	Delete(teamUser *models.CompetitionTeam) error
}
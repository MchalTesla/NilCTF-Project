package repositories_interface

import (
	"NilCTF/models"
)

type TeamRepositoryInterface interface {
	Create(team *models.Team) error
	Read(ID uint, username string) ([]models.Team, error)
	Update(team *models.Team) error
	Delete(team *models.Team) error
}
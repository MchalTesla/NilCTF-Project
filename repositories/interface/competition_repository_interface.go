package repositories_interface

import (
	"NilCTF/models"
)

type CompetitionRepositoryInterface interface {
	Create(competition *models.Competition) error
	Read(ID uint, name string, ownerID uint) ([]models.Competition, error)
	Update(competition *models.Competition) error
	Delete(competition *models.Competition) error
}
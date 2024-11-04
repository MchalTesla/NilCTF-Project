package managers_interface

import (
	"NilCTF/models"
)

type CompetitionManagerInterface interface {
	Create(competition *models.Competition) error
	Get(ID uint) (*models.Competition, error)
	Update(competition *models.Competition) error
	Delete(competition *models.Competition) error
	List(filters map[string]interface{}, limit, offset int, isFuzzy bool) ([]models.Competition, error)
}

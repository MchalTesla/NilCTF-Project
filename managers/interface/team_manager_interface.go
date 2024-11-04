package managers_interface

import (
	"NilCTF/models"
)
type TeamManagerInterface interface {
	Create(team *models.Team) error
	Get(ID uint) (*models.Team, error)
	Update(team *models.Team) error
	Delete(team *models.Team) error
	List(filters map[string]interface{}, limit, offset int, isFuzzy bool) ([]models.Team, error)
}

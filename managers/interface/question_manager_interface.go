package managers_interface

import (
	"NilCTF/models"
)

type QuestionManagerInterface interface {
	Create(question *models.Question) error
	Get(ID uint) (*models.Question, error)
	Update(question *models.Question) error
	Delete(question *models.Question) error
	List(filters map[string]interface{}, limit, offset int, isFuzzy bool) ([]models.Question, error)
}

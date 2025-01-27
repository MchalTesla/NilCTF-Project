package repositories_interface

import (
	"NilCTF/models"
)

type SubmissionRepositoryInterface interface {
	Create(submission *models.Submission) error
	Get(ID uint) (*models.Submission, error)
	Update(submission *models.Submission) error
	Delete(submission *models.Submission) error
	List(filters map[string]interface{}, limit, offset int, isFuzzy bool) ([]models.Submission, error)
}

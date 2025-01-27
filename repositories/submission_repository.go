package repositories

import (
	"NilCTF/error_code"
	"NilCTF/models"

	"gorm.io/gorm"
)

type SubmissionRepository struct {
	DB *gorm.DB
}

func NewSubmissionRepository(DB *gorm.DB) *SubmissionRepository {
	return &SubmissionRepository{DB: DB}
}

func (r *SubmissionRepository) Create(submission *models.Submission) error {
	if submission.ID != 0 {
		return error_code.ErrInvalidID
	}
	if err := r.DB.Create(submission).Error; err != nil {
		return error_code.ErrInternalServer
	}
	return nil
}

func (r *SubmissionRepository) Get(ID uint) (*models.Submission, error) {
	if ID == 0 {
		return nil, error_code.ErrInvalidID
	}

	var submission models.Submission
	if err := r.DB.First(&submission, ID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, error_code.ErrSubmissionNotFound
		}
		return nil, error_code.ErrInternalServer
	}
	return &submission, nil
}

func (r *SubmissionRepository) Update(submission *models.Submission) error {
	if submission.ID == 0 {
		return error_code.ErrInvalidID
	}
	if err := r.DB.Updates(submission).Error; err != nil {
		return error_code.ErrInternalServer
	}
	return nil
}

func (r *SubmissionRepository) Delete(submission *models.Submission) error {
	if submission.ID == 0 {
		return error_code.ErrInvalidID
	}
	if err := r.DB.Unscoped().Delete(submission).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return error_code.ErrSubmissionNotFound
		}
		return error_code.ErrInternalServer
	}
	return nil
}

func (r *SubmissionRepository) List(filters map[string]interface{}, limit, offset int, isFuzzy bool) ([]models.Submission, error) {
	var submissions []models.Submission
	query := r.DB

	for key, value := range filters {
		if isFuzzy {
			query = query.Where(key+" LIKE ?", "%"+value.(string)+"%")
		} else {
			query = query.Where(key+" = ?", value)
		}
	}

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset >= 0 {
		query = query.Offset(offset)
	}

	if err := query.Find(&submissions).Error; err != nil {
		return []models.Submission{}, error_code.ErrInternalServer
	}
	return submissions, nil
}

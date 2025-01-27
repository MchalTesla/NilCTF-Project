package repositories

import (
	"NilCTF/error_code"
	"NilCTF/models"

	"gorm.io/gorm"
)

type QuestionRepository struct {
	DB *gorm.DB
}

func NewQuestionRepository(DB *gorm.DB) *QuestionRepository {
	return &QuestionRepository{DB: DB}
}

func (r *QuestionRepository) Create(question *models.Question) error {
	if question.ID != 0 {
		return error_code.ErrInvalidID
	}
	if err := r.DB.Create(question).Error; err != nil {
		return error_code.ErrInternalServer
	}
	return nil
}

func (r *QuestionRepository) Get(ID uint) (*models.Question, error) {
	if ID == 0 {
		return nil, error_code.ErrInvalidID
	}

	var question models.Question
	if err := r.DB.First(&question, ID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, error_code.ErrQuestionNotFound
		}
		return nil, error_code.ErrInternalServer
	}
	return &question, nil
}

func (r *QuestionRepository) Update(question *models.Question) error {
	if question.ID == 0 {
		return error_code.ErrInvalidID
	}
	if err := r.DB.Updates(question).Error; err != nil {
		return error_code.ErrInternalServer
	}
	return nil
}

func (r *QuestionRepository) Delete(question *models.Question) error {
	if question.ID == 0 {
		return error_code.ErrInvalidID
	}
	if err := r.DB.Unscoped().Delete(question).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return error_code.ErrQuestionNotFound
		}
		return error_code.ErrInternalServer
	}
	return nil
}

func (r *QuestionRepository) List(filters map[string]interface{}, limit, offset int, isFuzzy bool) ([]models.Question, error) {
	var questions []models.Question
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

	if err := query.Find(&questions).Error; err != nil {
		return []models.Question{}, error_code.ErrInternalServer
	}
	return questions, nil
}

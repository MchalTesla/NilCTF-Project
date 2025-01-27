package managers

import (
	"NilCTF/error_code"
	"NilCTF/models"
	repositories_interface "NilCTF/repositories/interface"
	"NilCTF/utils"
)

type QuestionManager struct {
	QR repositories_interface.QuestionRepositoryInterface
}

func NewQuestionManager(QR repositories_interface.QuestionRepositoryInterface) *QuestionManager {
	return &QuestionManager{QR: QR}
}

func (m *QuestionManager) Create(question *models.Question) error {
	if question.ID != 0 {
		return error_code.ErrInvalidID
	}

	// 检查题目标题是否符合规范
	if !utils.IsValidName(question.Title) {
		return error_code.ErrInvalidName
	}

	// 检查题目描述是否符合规范
	if !utils.IsValidDescription(question.Description) {
		return error_code.ErrInvalidDescription
	}

	// 检查题目分值是否大于零
	if question.Points < 0 {
		return error_code.ErrInvalidPoints
	}

	return m.QR.Create(question)
}

func (m *QuestionManager) Update(question *models.Question) error {
	if question.ID == 0 {
		return error_code.ErrInvalidID
	}

	if _, err := m.QR.Get(question.ID); err != nil {
		return err
	}

	if !utils.IsValidName(question.Title) {
		return error_code.ErrInvalidName
	}

	if !utils.IsValidDescription(question.Description) {
		return error_code.ErrInvalidDescription
	}

	if question.Points <= 0 {
		return error_code.ErrInvalidPoints
	}

	return m.QR.Update(question)
}

func (m *QuestionManager) Get(ID uint) (*models.Question, error) {
	if ID == 0 {
		return nil, error_code.ErrInvalidID
	}
	return m.QR.Get(ID)
}

func (m *QuestionManager) Delete(question *models.Question) error {
	if question.ID == 0 {
		return error_code.ErrInvalidID
	}
	return m.QR.Delete(question)
}

func (m *QuestionManager) List(filters map[string]interface{}, limit, offset int, isFuzzy bool) ([]models.Question, error) {
	return m.QR.List(filters, limit, offset, isFuzzy)
}

package managers

import (
	"NilCTF/error_code"
	"NilCTF/models"
	repositories_interface "NilCTF/repositories/interface"
)

type SubmissionManager struct {
	SR repositories_interface.SubmissionRepositoryInterface
}

func NewSubmissionManager(SR repositories_interface.SubmissionRepositoryInterface) *SubmissionManager {
	return &SubmissionManager{SR: SR}
}

func (m *SubmissionManager) Create(submission *models.Submission) error {
	if submission.ID != 0 {
		return error_code.ErrInvalidID
	}

	// 检查题目 ID 是否有效
	if submission.QuestionID == 0 {
		return error_code.ErrInvalidQuestionID
	}

	// 检查用户 ID 是否有效
	if submission.UserID == 0 {
		return error_code.ErrInvalidUserID
	}

	// 检查答案是否为空
	if submission.Answer == "" {
		return error_code.ErrEmptyAnswer
	}

	return m.SR.Create(submission)
}

func (m *SubmissionManager) Update(submission *models.Submission) error {
	if submission.ID == 0 {
		return error_code.ErrInvalidID
	}

	if _, err := m.SR.Get(submission.ID); err != nil {
		return err
	}

	if submission.Answer == "" {
		return error_code.ErrEmptyAnswer
	}

	return m.SR.Update(submission)
}

func (m *SubmissionManager) Get(ID uint) (*models.Submission, error) {
	if ID == 0 {
		return nil, error_code.ErrInvalidID
	}
	return m.SR.Get(ID)
}

func (m *SubmissionManager) Delete(submission *models.Submission) error {
	if submission.ID == 0 {
		return error_code.ErrInvalidID
	}
	return m.SR.Delete(submission)
}

func (m *SubmissionManager) List(filters map[string]interface{}, limit, offset int, isFuzzy bool) ([]models.Submission, error) {
	return m.SR.List(filters, limit, offset, isFuzzy)
}
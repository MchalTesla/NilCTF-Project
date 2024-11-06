package managers

import (
	"NilCTF/error_code"
	"NilCTF/models"
	repositories_interface "NilCTF/repositories/interface"
	"NilCTF/utils"
)

type CompetitionManager struct {
	CR repositories_interface.CompetitionRepositoryInterface
}

func NewCompetitionManager(CR repositories_interface.CompetitionRepositoryInterface) *CompetitionManager {
	return &CompetitionManager{CR: CR}
}

func (m *CompetitionManager) Create(competition *models.Competition) error {
	if competition.ID != 0 {
		return error_code.ErrInvalidID
	}

	// 检查比赛名字是否符合规范
	if !utils.IsValidName(competition.Name) {
		return error_code.ErrInvalidName
	}

	// 检查比赛描述是否符合规范
	if !utils.IsValidDescription(competition.Description) {
		return error_code.ErrInvalidDescription
	}
	return m.CR.Create(competition)
}

func (m *CompetitionManager) Update(competition *models.Competition) error {
	if competition.ID == 0 {
		return error_code.ErrInvalidID
	}

	if _, err := m.CR.Get(competition.ID); err != nil {
		return err
	}

	if !utils.IsValidName(competition.Name) {
		return error_code.ErrInvalidName
	}

	if !utils.IsValidDescription(competition.Description) {
		return error_code.ErrInvalidDescription
	}
	return m.Update(competition)
}

func (m *CompetitionManager) Get(ID uint) (*models.Competition, error) {
	if ID == 0 {
		return nil, error_code.ErrInvalidID
	}
	return m.CR.Get(ID)
}

func (m *CompetitionManager) Delete(competition *models.Competition) error {
	if competition.ID == 0 {
		return error_code.ErrInvalidID
	}
	return m.CR.Delete(competition)
}

func (m *CompetitionManager) List(filters map[string]interface{}, limit, offset int, isFuzzy bool) ([]models.Competition, error) {
	return m.CR.List(filters, limit, offset, isFuzzy)
}

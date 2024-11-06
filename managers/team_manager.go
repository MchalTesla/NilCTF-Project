package managers

import (
	"NilCTF/error_code"
	"NilCTF/models"
	repositories_interface "NilCTF/repositories/interface"
	"NilCTF/utils"
)

type TeamManager struct {
	TR repositories_interface.TeamRepositoryInterface
}

func NewTeamManager(TR repositories_interface.TeamRepositoryInterface) *TeamManager {
	return &TeamManager{TR: TR}
}

func (m *TeamManager) Create(team *models.Team) error {
	if team.ID != 0 {
		return error_code.ErrInvalidID
	}

	// 判断队伍名规范性
	if !utils.IsValidName(team.Name) {
		return error_code.ErrInvalidName
	}

	// 检查队伍描述是否符合规范
	if !utils.IsValidDescription(team.Description) {
		return error_code.ErrInvalidDescription
	}

	return m.TR.Create(team)
}

func (m *TeamManager) Update(team *models.Team) error {
	// 检查队伍ID是否有效
	if team.ID == 0 {
		return error_code.ErrInvalidID
	}

	// 检查队伍ID是否存在
	if _, err := m.Get(team.ID); err != nil {
		return err
	}

	// 检查队伍名中是否符合规范
	if !utils.IsValidName(team.Name) {
		return error_code.ErrInvalidName
	}

	// 检查队伍描述是否符合规范
	if !utils.IsValidDescription(team.Description) {
		return error_code.ErrInvalidDescription
	}
	return m.Update(team)
}

func (m *TeamManager) Get(ID uint) (*models.Team, error) {
	if ID == 0 {
		return nil, error_code.ErrInvalidID
	}
	return m.TR.Get(ID)
}

func (m *TeamManager) Delete(team *models.Team) error {
	if team.ID == 0 {
		return error_code.ErrInvalidID
	}
	return m.TR.Delete(team)
}

func (m *TeamManager) List(filters map[string]interface{}, limit, offset int, isFuzzy bool) ([]models.Team, error) {
	return m.TR.List(filters, limit, offset, isFuzzy)
}

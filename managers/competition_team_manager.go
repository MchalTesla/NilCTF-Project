package managers

import (
	"NilCTF/models"
	"NilCTF/repositories/interface"
	"NilCTF/error_code"
)

type CompetitionTeamManager struct {
	CTR repositories_interface.CompetitionTeamRepositoryInterface
}

func NewCompetitionTeamManager(CTR repositories_interface.CompetitionTeamRepositoryInterface) *CompetitionTeamManager {
	return &CompetitionTeamManager{CTR: CTR}
}

func (m *CompetitionTeamManager) Create(competitionTeam *models.CompetitionTeam) error {
	//判断ID是否合规
	if competitionTeam.ID != 0 {
		return error_code.ErrInvalidID
	}

	// 检查是否存在
	existingMappings, err := m.CTR.Get(0, competitionTeam.CompetitionID, competitionTeam.TeamID)
	if err != nil {
		return err
	}
	if len(existingMappings) > 0 {
		return error_code.ErrTeamAlreadyInCompetition
	}

	return m.CTR.Create(competitionTeam)
}

func (m *CompetitionTeamManager) Update(competitionTeam *models.CompetitionTeam) error {
	if competitionTeam.ID == 0 {
		return error_code.ErrInvalidID
	}

	existingMappings, err := m.CTR.Get(competitionTeam.ID, 0, 0)
	if err != nil {
		return err
	}
	if len(existingMappings) == 0 {
		return error_code.ErrTeamNotInCompetition
	}

	// 如果 CompetitionID 或 TeamID 变更，拒绝更新
	if (competitionTeam.CompetitionID != 0 && existingMappings[0].CompetitionID != competitionTeam.CompetitionID) ||
	(competitionTeam.TeamID != 0 && existingMappings[0].TeamID != competitionTeam.TeamID) {
		return error_code.ErrInvalidID
	}
	
	return m.CTR.Update(competitionTeam)
}

func (m *CompetitionTeamManager) Get(ID, competitionID, teamID uint) ([]models.CompetitionTeam, error) {
	return m.CTR.Get(ID, competitionID, teamID)
}

func (m *CompetitionTeamManager) Delete(competitionTeam *models.CompetitionTeam) error {
	if competitionTeam.ID == 0{
		return error_code.ErrInvalidID
	}
	return m.CTR.Delete(competitionTeam)
}
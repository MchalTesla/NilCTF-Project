package repositories

import (
	"NilCTF/error_code"
	"NilCTF/models"
	"errors"

	"gorm.io/gorm"
)

type CompetitionTeamRepository struct {
	DB *gorm.DB
}

// NewCompetitionTeamRepository 返回新的 CompetitionTeamRepository 实例
func NewCompetitionTeamRepository(DB *gorm.DB) *CompetitionTeamRepository {
	return &CompetitionTeamRepository{DB: DB}
}

// Create 创建比赛和队伍的映射 ID必须为0
func (r *CompetitionTeamRepository) Create(competitionTeam *models.CompetitionTeam) error {
	//判断ID是否合规
	if competitionTeam.ID != 0 {
		return error_code.ErrInvalidID
	}

	if err := r.DB.Where("competition = ? AND teamid = ?", competitionTeam.CompetitionID, competitionTeam.TeamID).First(&models.CompetitionTeam{}).Error; err == nil {
		return error_code.ErrTeamAlreadyInCompetition
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return error_code.ErrInternalServer
	}

	if err := r.DB.Create(competitionTeam).Error; err != nil {
		return error_code.ErrInternalServer
	}
	return nil
}

// Read 查找比赛和队伍的映射
func (r *CompetitionTeamRepository) Read(ID, competitionID, teamID uint) ([]models.CompetitionTeam, error) {
	var competitionTeams []models.CompetitionTeam

	// 根据ID、 competitionID或UserID查找
	var err error
	switch {
	case ID != 0:
		err = r.DB.Find(&competitionTeams, ID).Error
	case competitionID != 0 && teamID != 0:
		err = r.DB.Where("competitionid = ? AND teamid = ?", competitionID, teamID).Find(&competitionTeams).Error
	case competitionID != 0:
		err = r.DB.Where("competitionid = ?", competitionID).Find(&competitionTeams).Error
	case teamID != 0:
		err = r.DB.Where("teamid = ?", teamID).Find(&competitionTeams).Error
	default:
		return nil,error_code.ErrInvalidInput
	}

	if err != nil {
		return nil, error_code.ErrInternalServer
	} else if len(competitionTeams) == 0 {
		return nil, error_code.ErrTeamNotInCompetition
	}
	
	return competitionTeams, nil
}

// Update 更新比赛和队伍的映射, 不能更改CompetitionID和TeamID, ID、CompetitionID、TeamID必须存在
func (r *CompetitionTeamRepository) Update(competitionTeam *models.CompetitionTeam) error {
	var existingCompetitionTeam models.CompetitionTeam
	//检查比赛-组ID是否有效
	if competitionTeam.ID == 0 {
		return error_code.ErrInvalidID
	}

	// 检查比赛-组ID是否存在
	if err := r.DB.Where("id = ?", competitionTeam.ID).First(&existingCompetitionTeam).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return error_code.ErrTeamNotInCompetition
		// 判断比赛ID和队伍ID是否和数据库中记录的一样，如果不一样，返回错误
		} else if competitionTeam.CompetitionID != existingCompetitionTeam.CompetitionID || competitionTeam.TeamID != existingCompetitionTeam.TeamID {
			return error_code.ErrInvalidInput
		}
		return error_code.ErrInternalServer
	}

	if err := r.DB.Model(competitionTeam).Updates(competitionTeam).Error; err != nil {
		return error_code.ErrInternalServer
	}
	return nil
}

// Delete 删除队伍和用户的映射, ID必须存在
func (r *CompetitionTeamRepository) Delete(competitionTeam *models.CompetitionTeam) error {
	//判断ID是否有效
	if competitionTeam.ID == 0 {
		return error_code.ErrInvalidID
	}
	if err := r.DB.Delete(competitionTeam).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return error_code.ErrTeamNotInCompetition
		}
		return error_code.ErrInternalServer
	}
	return nil
}
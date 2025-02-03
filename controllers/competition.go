package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type CompetitionController struct {

}

func NewCompetitionController() *CompetitionController {
	return &CompetitionController{}
}


func (r *CompetitionController) ListCompetition(c *gin.Context) {
	// var competitions []models.Competition

	// // 获取所有未隐藏的比赛记录
	// if err := config.DB.Where("is_hidden = ?", false).Find(&competitions).Error; err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": "获取比赛列表失败", "error": err.Error()})
	// 	return
	// }

	// var response []gin.H
	// for _, competition := range competitions {
	// 	var participantCount int64
	// 	// 获取参加当前比赛的队伍数量
	// 	config.DB.Model(&models.CompetitionTeam{}).Where("competition_id = ?", competition.ID).Count(&participantCount)

	// 	response = append(response, gin.H{
	// 		"name":              competition.Name,
	// 		"description":       competition.Description,
	// 		"nature":            competition.Nature,
	// 		"owner":             competition.OwnerID, // 你可能需要在用户表中查找用户名
	// 		"status":            competition.Status,
	// 		"start_time":        competition.StartTime,
	// 		"end_time":          competition.EndTime,
	// 		"participant_count": participantCount,        // 参与队伍数量
	// 		"team_limit":        competition.TeamLimit,   // 限制的队伍数量
	// 		"max_team_size":     competition.MaxTeamSize, // 队伍最大人数限制
	// 	})
	// }

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": gin.H{"name": "test"}})
}

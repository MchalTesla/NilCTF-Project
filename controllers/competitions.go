package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"AWD-Competition-Platform/config"
	"AWD-Competition-Platform/models"
)

func Competitions(c *gin.Context) {
	var competitions []models.Competitions

	// 获取所有比赛记录
	if err := config.DB.Preload("Owner").Find(&competitions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": "获取比赛列表失败", "error": err.Error()})
		return
	}

	var response []gin.H
	for _, competition := range competitions {
		response = append(response, gin.H{
			"name":        competition.Name,
			"description": competition.Description,
			"nature":      competition.Nature,
			"owner":       competition.OwnerID, // 你可能需要在用户表中查找用户名
			"status":      competition.Status,
			"start_time":  competition.StartTime,
			"end_time":    competition.EndTime,
		})
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": response})
}
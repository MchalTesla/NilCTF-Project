package controllers

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"AWD-Competition-Platform/config"
	"AWD-Competition-Platform/models"
	"AWD-Competition-Platform/services"
	"AWD-Competition-Platform/repositories"
)

func Register(c *gin.Context) {
	UR := repositories.NewUserRepository(config.DB)
	US := services.NewUserService(UR)

	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := US.Register(user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "用户注册成功", "redirect": "/login"})
}
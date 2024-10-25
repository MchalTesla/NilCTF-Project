package controllers

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"AWD-Competition-Platform/config"
	"AWD-Competition-Platform/models"
	"AWD-Competition-Platform/middleware"
)

func Login(c *gin.Context) {
	var user models.User
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	if err := config.DB.Where("username = ?", input.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": "用户未注册"})
		return
	}

	if !user.CheckPassword(input.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": "用户名或密码错误"})
		return
	}

	token, err := middleware.GenerateToken(user.ID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": "无法登录"})
        return
    }

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "登录成功", "token": token, "redirect": "/index"})
}
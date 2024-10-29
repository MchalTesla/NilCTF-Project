package controllers

import (
	"NilCTF/config"
	"NilCTF/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func IndexController(c *gin.Context) {
	// middleware获取ID
	userID, existing := c.Get("userID")
	if !existing {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "fail"})
	}
	// 创建一个user变量，用于存储用户信息
	var user models.User
	// 从数据库中获取用户名
	if err := config.DB.Where("ID = ?", userID).Find(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "fail"})
	}

	// 获取当前时间
	currentTime := time.Now().Format("2006-01-02 15:04:05")

	// 返回 JSON 响应
	c.JSON(http.StatusOK, gin.H{
		"message":     "欢迎！",
		"username":    user.Username,
		"currentTime": currentTime,
	})
}

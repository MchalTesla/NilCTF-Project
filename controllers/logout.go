package controllers

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

func Logout(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "退出登录成功", "redirect": "/index"})
}
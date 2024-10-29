package controllers

import (
	"NilCTF/models"
	"NilCTF/services/interface"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Register(c *gin.Context, US services_interface.UserServiceInterface) {

	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := US.Register(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "用户注册成功", "redirect": "/login"})
}

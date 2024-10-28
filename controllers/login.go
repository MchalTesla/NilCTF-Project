package controllers

import (
	"NilCTF/middleware"
	services_interface "NilCTF/services/interface"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context, US services_interface.UserServiceInterface) {

	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "fail"})
		return
	}

	user, err := US.Login(input.Email, input.Username, input.Password)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	token, err := middleware.GenerateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": "无法登录"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "登录成功", "token": token, "redirect": "/index"})
}

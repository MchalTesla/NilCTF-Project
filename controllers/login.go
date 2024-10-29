package controllers

import (
	"NilCTF/error_code"
	"NilCTF/middleware"
	services_interface "NilCTF/services/interface"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context, US services_interface.UserServiceInterface) {

	var input struct {
		LoginIdentifier string `json:"loginidentifier"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": error_code.ErrInvalidInput.Error()})
		return
	}

	user, err := US.Login(input.LoginIdentifier, input.Password)

	if err != nil {
		httpStatus := http.StatusInternalServerError
		if err != error_code.ErrInternalServer{
			httpStatus = http.StatusUnauthorized
			err =  error_code.ErrInvalidCredentials
		}
		c.JSON(httpStatus, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	token, err := middleware.GenerateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": "无法登录"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "登录成功", "token": token, "redirect": "/index"})
}

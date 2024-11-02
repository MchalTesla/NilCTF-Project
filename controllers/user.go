package controllers

import (
	"NilCTF/error_code"
	"NilCTF/middleware"
	"NilCTF/services/interface"
	"NilCTF/models"
	"net/http"
	"errors"

	"github.com/gin-gonic/gin"
)

type UserControllers struct {
}

func (r *UserControllers) Login(c *gin.Context, US services_interface.UserServiceInterface) {

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

	postMiddleware := middleware.NewPostMiddleware()
	token, err := postMiddleware.GenerateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": error_code.ErrInternalServer.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "登录成功", "token": token, "redirect": "/index"})
}

func (r *UserControllers) Logout(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "退出登录成功", "redirect": "/index", "clear_token": true})
}

func (r *UserControllers) Register(c *gin.Context, US services_interface.UserServiceInterface) {

	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": error_code.ErrInvalidInput.Error()})
		return
	}

	if err := US.Register(&user); err != nil {
		if errors.Is(err, error_code.ErrInternalServer) {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": err.Error()})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "用户注册成功", "redirect": "/login"})
}

func (r *UserControllers) VerifyLogin(c *gin.Context) {
	userRole, _ := c.Get("userRole")
	c.JSON(http.StatusOK, gin.H{"status": "success", "user_role": userRole})
}
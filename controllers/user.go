package controllers

import (
	"NilCTF/dto"
	"NilCTF/error_code"
	"NilCTF/middleware"
	"NilCTF/services/interface"
	"NilCTF/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserControllers struct {
	US services_interface.UserServiceInterface
	cookieSecure bool
	jwtTime int
	jwtSecret []byte
	postMiddleware *middleware.PostMiddleware
}

func NewUserControllers(US services_interface.UserServiceInterface, 
		cookieSecure bool,
		jwtTime int,
		jwtSecret []byte,
	) *UserControllers {
	return &UserControllers{
		US: US,
		cookieSecure: cookieSecure,
		jwtTime: jwtTime,
	}
}

func (r *UserControllers) Login(c *gin.Context) {

	var input struct {
		LoginIdentifier string `json:"loginidentifier"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": error_code.ErrInvalidInput.Error()})
		return
	}

	userID, err := r.US.Login(input.LoginIdentifier, input.Password)

	if err != nil {
		httpStatus := http.StatusInternalServerError
		if err != error_code.ErrInternalServer {
			httpStatus = http.StatusUnauthorized
			switch err {
			case error_code.ErrUserBanned, error_code.ErrUserPending:
			default:
				err = error_code.ErrInvalidCredentials
			}
		}
		c.JSON(httpStatus, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	token, err := utils.GenerateToken(userID, r.jwtTime, r.jwtSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": error_code.ErrInternalServer.Error()})
		return
	}

	// Set token as a cookie
	c.SetCookie("auth_token", token, r.jwtTime * 60 * 60, "/", "", r.cookieSecure, true)

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "登录成功", "redirect": "/index"})
}


func (r *UserControllers) Logout(c *gin.Context) {
	c.SetCookie("auth_token", "", -1, "/", "/", r.cookieSecure, true)
	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "退出登录成功", "redirect": "/index", "clear_token": true})
}

func (r *UserControllers) Register(c *gin.Context) {

	var user dto.UserCreate
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": error_code.ErrInvalidInput.Error()})
		return
	}

	if err := r.US.Register(&user); err != nil {
		if err == error_code.ErrInternalServer {
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
	userStatus, _ := c.Get("userStatus")
	c.JSON(http.StatusOK, gin.H{"status": "success", "user_role": userRole, "user_status": userStatus})
	return
}
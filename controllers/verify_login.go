package controllers

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

func VerifyLogin(c *gin.Context) {
	userRole, _ := c.Get("userRole")
	c.JSON(http.StatusOK, gin.H{"status": "success", "user_role": userRole})
}
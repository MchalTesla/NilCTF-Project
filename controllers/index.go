package controllers

import (
	"NilCTF/services/interface"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type IndexControllers struct {
	US services_interface.UserServiceInterface
}

func NewIndexControllers(US services_interface.UserServiceInterface) *IndexControllers {
	return &IndexControllers{US: US}
}

func (ic *IndexControllers) Index(c *gin.Context) {
	// middleware获取ID
	userID, existing := c.Get("userID")
	if !existing {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "fail"})
	}

	user, err := ic.US.GetNow(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": err.Error()})
		return
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

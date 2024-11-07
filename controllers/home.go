package controllers

import (
	"NilCTF/error_code"
	"NilCTF/services/interface"
	"NilCTF/dto"
	"net/http"

	"github.com/gin-gonic/gin"
)

type HomeControllers struct {
    HS services_interface.HomeServiceInterface
}

func NewHomeControllers(HS services_interface.HomeServiceInterface) *HomeControllers {
    return &HomeControllers{HS: HS}
}

// Home 用户主页
func (hc *HomeControllers) Home(c *gin.Context) {
    var userID uint
    {
        id, _ := c.Get("userID")
        userID = id.(uint)
    }

    // 获取用户信息
    updates, err := hc.HS.Info(userID)
    if err != nil {
        if err == error_code.ErrInternalServer {
            c.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": error_code.ErrInvalidInput.Error()})
        } else {
            c.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
        }
        return
    }

    // 返回成功响应和用户信息
    c.JSON(http.StatusOK, gin.H{
        "status": "success",
        "message": gin.H{
            "created_at":   updates.CreatedAt,
            "username":     updates.Username,
            "description":  updates.Description,
            "email":        updates.Email,
            "status":       updates.Status,
            "role":         updates.Role,
            "tag":          updates.Tag,
        },
    })
}

// UpdateUser 更新用户信息的 API
func (r *HomeControllers) Modify(c *gin.Context, US services_interface.UserServiceInterface) {
    var updates dto.UserUpdate
    var userID uint
	{
		id, _ := c.Get("userID")
		userID = id.(uint)
	}

    // 解析请求体中的 JSON 数据
    if err := c.ShouldBindJSON(&updates); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"status":"fail", "message": error_code.ErrInvalidInput.Error()})
        return
    }

    // 调用服务层更新用户信息
    if err := US.Update(userID, &updates); err != nil {
        // 根据错误类型返回不同的状态码
		if err == error_code.ErrInternalServer {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": error_code.ErrInternalServer.Error()})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		}
        return
    }

    c.JSON(http.StatusOK, gin.H{"status": "success", "message": "User updated successfully", "redirect": "/home"})
}
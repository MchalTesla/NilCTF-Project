package controllers

import (
	"NilCTF/dto"
	"NilCTF/error_code"
	"NilCTF/services/interface"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AdminUserController struct {
	MS services_interface.AdminUserService
}

func NewAdminUserController(MS services_interface.AdminUserService) *AdminUserController {
	return &AdminUserController{MS: MS}
}

func (mc *AdminUserController) GetUsersCount(c *gin.Context) {
	userCount, err := mc.MS.GetUsersCount()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": error_code.ErrInvalidPageParameter.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok", "message": userCount})
}

// ListUsers 列出所有用户，可分页并选择每页显示数量
func (mc *AdminUserController) ListUsers(c *gin.Context) {
	// 从查询参数获取分页信息
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	// 转换分页参数为整数
	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": error_code.ErrInvalidPageParameter.Error()})
		return
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": error_code.ErrInvalidLimitParameter.Error()})
		return
	}
	
	usersDTO, err := mc.MS.ListAllUsers(page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": error_code.ErrInternalServer})
		return
	}

	pages, err := mc.MS.GetTotalPages(limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": error_code.ErrInternalServer})
		return
	}

	// 返回用户列表
	c.JSON(http.StatusOK, gin.H{
		"status": "ok", "message": gin.H{
			"pages": pages,
			"page": page,
			"users": usersDTO,
		},
	})
}

// HandleUser 处理用户相关操作
func (mc *AdminUserController) HandleUser(c *gin.Context) {
    var request struct {
        Action string           `json:"action"`
        User   dto.UserUpdateByAdmin `json:"user"`
    }

    if err := c.ShouldBindJSON(&request); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": error_code.ErrInvalidInput.Error()})
        return
    }

    switch request.Action {
    case "create":
        mc.createUser(c, &request.User)
    case "update":
        mc.updateUser(c, &request.User)
    case "delete":
        mc.deleteUser(c, &request.User)
    default:
        c.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "未知操作"})
    }
}

func (mc *AdminUserController) createUser(c *gin.Context, user *dto.UserUpdateByAdmin) {
    if err := mc.MS.CreateUser(user); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (mc *AdminUserController) updateUser(c *gin.Context, user *dto.UserUpdateByAdmin) {
    if err := mc.MS.UpdateUsers(user); err != nil {
        if err == error_code.ErrInternalServer {
            c.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": error_code.ErrInternalServer.Error()})
        } else {
            c.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
        }
        return
    }
    c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (mc *AdminUserController) deleteUser(c *gin.Context, user *dto.UserUpdateByAdmin) {
    if err := mc.MS.DeleteUser(user.ID); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

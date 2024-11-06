package controllers

import (
	"NilCTF/dto"
	"NilCTF/error_code"
	"NilCTF/services/interface"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ManagerController struct {
	MS services_interface.ManagerService
}

func NewManagerController(MS services_interface.ManagerService) *ManagerController {
	return &ManagerController{MS: MS}
}

func (mc *ManagerController) GetUsersCount(c *gin.Context) {
	userCount, err := mc.MS.GetUsersCount()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": error_code.ErrInvalidPageParameter.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok", "message": userCount})
}

// ListUsers 列出所有用户，可分页并选择每页显示数量
func (mc *ManagerController) ListUsers(c *gin.Context) {
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

// UpdateUserByAdmin 管理员更新用户信息
func (mc *ManagerController) UpdateUserByAdmin(c *gin.Context) {
	var updates dto.UserUpdateByAdmin

    // 解析请求体中的 JSON 数据
    if err := c.ShouldBindJSON(&updates); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"status":"fail", "message": error_code.ErrInvalidInput.Error()})
        return
    }

	// 调用服务层更新用户信息
	if err := mc.MS.UpdateUsers(&updates); err != nil {
		// 根据错误类型返回不同的状态码
		if err == error_code.ErrInternalServer {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": error_code.ErrInternalServer.Error()})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
	return
}

// DeleteUserByAdmin 管理员删除用户
func (mc *ManagerController) DeleteUserByAdmin(c *gin.Context) {
	var currentUserID uint
	{
		id, _ := c.Get("userID")
		currentUserID = id.(uint)
	}
	user := struct {ID uint}{0,}
	// 获取指定用户

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status":"fail", "message": error_code.ErrInvalidInput.Error()})
        return
	}

	if currentUserID == user.ID {
		c.JSON(http.StatusBadRequest, gin.H{"status":"fail", "message": error_code.ErrInvalidInput.Error()})
		return 
	}

	// 删除用户
	if err := mc.MS.DeleteUser(user.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

package controllers

import (
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

// // UpdateUserByAdmin 管理员更新用户信息
// func (auc *AdminUserController) UpdateUserByAdmin(c *gin.Context) {
// 	var user models.User
// 	userId := c.Param("id")

// 	// 获取指定用户
// 	if err := auc.DB.First(&user, userId).Error; err != nil {
// 		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
// 		return
// 	}

// 	// 解析更新数据
// 	var updateData struct {
// 		Username string `json:"username,omitempty"`
// 		Email    string `json:"email,omitempty"`
// 		Password string `json:"password,omitempty"`
// 	}

// 	if err := c.ShouldBindJSON(&updateData); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data"})
// 		return
// 	}

// 	// 更新用户名和邮箱
// 	if updateData.Username != "" {
// 		user.Username = updateData.Username
// 	}
// 	if updateData.Email != "" {
// 		user.Email = updateData.Email
// 	}

// 	// 更新密码
// 	if updateData.Password != "" {
// 		if err := user.HashPassword(updateData.Password); err != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": "Password encryption failed"})
// 			return
// 		}
// 	}

// 	// 保存更新
// 	if err := auc.DB.Save(&user).Error; err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
// }

// // DeleteUserByAdmin 管理员删除用户
// func (auc *AdminUserController) DeleteUserByAdmin(c *gin.Context) {
// 	var user models.User
// 	userId := c.Param("id")

// 	// 获取指定用户
// 	if err := auc.DB.First(&user, userId).Error; err != nil {
// 		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
// 		return
// 	}

// 	// 删除用户
// 	if err := auc.DB.Delete(&user).Error; err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
// }

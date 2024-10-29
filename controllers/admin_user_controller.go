package controllers

// import (
// 	"net/http"
// 	"strconv"

// 	"github.com/gin-gonic/gin"
// 	"gorm.io/gorm"
// 	"NilCTF/models"
// )

// type AdminUserController struct {
// 	DB *gorm.DB
// }

// // ListUsers 列出所有用户，可分页并选择每页显示数量
// func (auc *AdminUserController) ListUsers(c *gin.Context) {
// 	// 获取分页参数
// 	pageStr := c.PostForm("page")
// 	if pageStr == "" {
// 		pageStr = "1"
// 	}

// 	perPageStr := c.PostForm("per_page")
// 	if perPageStr == "" {
// 		perPageStr = "10"
// 	}

// 	// 解析分页参数
// 	page, err := strconv.Atoi(pageStr)
// 	if err != nil || page < 1 {
// 		page = 1
// 	}
// 	perPage, err := strconv.Atoi(perPageStr)
// 	if err != nil || perPage < 1 {
// 		perPage = 10
// 	}

// 	// 获取用户数据
// 	var users []models.User
// 	if err := auc.DB.Offset((page - 1) * perPage).Limit(perPage).Find(&users).Error; err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve users"})
// 		return
// 	}

// 	// 获取总用户数
// 	var total int64
// 	auc.DB.Model(&models.User{}).Count(&total)

// 	c.JSON(http.StatusOK, gin.H{
// 		"page":        page,
// 		"per_page":    perPage,
// 		"total":       total,
// 		"total_pages": (total + int64(perPage) - 1) / int64(perPage),
// 		"users":       users,
// 	})
// }

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

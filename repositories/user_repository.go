package repositories

import (
	"NilCTF/error_code"
	"NilCTF/models"
	"NilCTF/utils"
	"errors"

	"gorm.io/gorm"
)

type UserRepository struct {
	DB *gorm.DB
}

// NewUserRepository 返回一个新的 UserRepository 实例
func NewUserRepository(DB *gorm.DB) *UserRepository {
	return &UserRepository{DB: DB}
}

// Create 插入新的用户记录， ID必须为0
func (r *UserRepository) Create(user *models.User) error {

	// 判断ID是否符合规范
	if user.ID != 0 {
		return error_code.ErrInvalidID
	}

	// 检查用户名中是否符合规范
	if !utils.IsValidName(user.Username) {
		return error_code.ErrInvalidUsername
	}

	// 检查邮箱是否符合规范
	if !utils.IsValidEmail(user.Email) {
		return error_code.ErrInvalidEmail
	}

	// 检查描述是否符合规范
	if !utils.IsValidDescription(user.Description) {
		return error_code.ErrInvalidDescription
	}

	// 检查邮箱是否被占用
	if err := r.DB.Where("email = ?", user.Email).First(&models.User{}).Error; err == nil {
		return error_code.ErrEmailExists
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		// 捕获潜在系统错误
		return error_code.ErrInternalServer
	}

	// 检查用户名是否已存在
	if err := r.DB.Where("username = ?", user.Username).First(&models.User{}).Error; err == nil {
		return error_code.ErrUsernameExists
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		// 捕获潜在系统错误
		return error_code.ErrInternalServer
	}

	// 密码哈希化
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return error_code.ErrInternalServer
	}
	user.Password = hashedPassword

	// 插入新用户
	if err := r.DB.Create(user).Error; err != nil {
		return error_code.ErrInternalServer
	}
	return nil
}

// Read 查找单个用户记录
func (r *UserRepository) Read(ID uint, email, username string) (*models.User, error) {
	var user models.User
	query := r.DB.Model(&user)

	// 动态构建查询条件
	if ID != 0 {
		query = query.Where("id = ?", ID)
	} else if email != "" {
		query = query.Where("email = ?", email)
	} else if username != "" {
		query = query.Where("username = ?", username)
	} else {
		return nil, error_code.ErrInvalidInput
	}

	// 执行查询
	if err := query.First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, error_code.ErrUserNotFound
		}
		// 捕获系统错误
		return nil, error_code.ErrInternalServer
	}

	return &user, nil
}

// Update 更新用户记录, ID必须存在
func (r *UserRepository) Update(user *models.User) error {
	// 检查用户ID是否有效
	if user.ID == 0 {
		return error_code.ErrInvalidID
	}

	// 检查用户ID是否存在
	if err := r.DB.Where("id = ?", user.ID).First(&models.User{}).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return error_code.ErrUserNotFound
		}
		return error_code.ErrInternalServer
	}

	// 检查用户名中是否符合规范
	if !utils.IsValidName(user.Username) {
		return error_code.ErrInvalidUsername
	}

	// 检查邮箱是否符合规范
	if !utils.IsValidEmail(user.Email) {
		return error_code.ErrInvalidEmail
	}

	// 检查描述是否符合规范
	if !utils.IsValidDescription(user.Description) {
		return error_code.ErrInvalidDescription
	}

	// 检查邮箱是否被占用
	if err := r.DB.Where("email = ? AND id != ?", user.Email, user.ID).First(&models.User{}).Error; err == nil {
		return error_code.ErrEmailExists
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		// 捕获潜在系统错误
		return error_code.ErrInternalServer
	}

	// 检查用户名是否已存在
	if err := r.DB.Where("username = ? AND id != ?", user.Username, user.ID).First(&models.User{}).Error; err == nil {
		return error_code.ErrUsernameExists
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		// 捕获潜在系统错误
		return error_code.ErrInternalServer
	}

	if err := r.DB.Model(user).Updates(user).Error; err != nil {
		// 捕获系统错误
		return error_code.ErrInternalServer
	}
	return nil
}

// Delete 删除用户记录 ID必须存在
func (r *UserRepository) Delete(user *models.User) error {
	// 判断用户ID是否有效
	if user.ID == 0 {
		return error_code.ErrInvalidID
	}

	if err := r.DB.Delete(user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return error_code.ErrUserNotFound
		}
		// 捕获系统错误
		return error_code.ErrInternalServer
	}
	return nil
}

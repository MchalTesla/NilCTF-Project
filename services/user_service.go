package services

import (
	"NilCTF/error_code"
	"NilCTF/models"
	"NilCTF/repositories/interface"
	"NilCTF/utils"
	"NilCTF/dto"
)

// UserService 提供用户相关服务
type UserService struct {
	UR repositories_interface.UserRepositoryInterface
}

// NewUserService 返回一个新的 UserService 实例
func NewUserService(UR repositories_interface.UserRepositoryInterface) *UserService {
	return &UserService{UR: UR}
}

// Register 注册一个新用户
func (r *UserService) Register(user *models.User) error {
	// 检查用户是否已存在
	existingUser, err := r.UR.Read(0, user.Email, user.Username)
	if err == nil && existingUser != nil {
		return error_code.ErrEmailExists // 返回邮箱已被占用的错误
	}

	if err := r.UR.Create(user); err != nil {
		return err
	}
	return nil
}

// Login 用户登录
func (r *UserService) Login(loginIdentifier string, password string) (*models.User, error) {
	var existingUser *models.User
	var err error

	// 判断 loginIdentifier 是用户名还是邮箱
	if utils.IsValidEmail(loginIdentifier) {
		existingUser, err = r.UR.Read(0, loginIdentifier, "") // 通过邮箱查找用户
	} else if utils.IsValidName(loginIdentifier){
		existingUser, err = r.UR.Read(0, "", loginIdentifier) // 通过用户名查找用户
	} else {
		return nil, error_code.ErrInvalidInput
	}

	if err != nil {
		return nil, err // 处理其他可能的错误
	}

	if existingUser == nil {
		return nil, error_code.ErrUserNotFound // 返回用户未找到的错误
	}

	if !utils.CheckPassword(existingUser.Password, password) {
		return nil, error_code.ErrInvalidCredentials // 返回密码错误的错误
	}

	return existingUser, nil
}

// 修改用户信息
func (r *UserService) Update(userID uint, updates dto.UserUpdate) error{
	var user models.User

	// 根据传入的字段更新值
	user.ID = userID
	if updates.Username != nil {
		user.Username = *updates.Username
	}
	if updates.Password != nil {
		user.Password = *updates.Password
	}
	if updates.Description != nil {
		user.Description = *updates.Description
	}
	if updates.Email != nil {
		user.Email = *updates.Email
	}

	if err := r.UR.Update(&user); err != nil {
		return err
	}
	return nil
}
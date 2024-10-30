package services

import (
	"NilCTF/error_code"
	"NilCTF/models"
	"NilCTF/repositories/interface"
	"NilCTF/utils"
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
	} else {
		existingUser, err = r.UR.Read(0, "", loginIdentifier) // 通过用户名查找用户
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
func (r *UserService) Update(user *models.User) error{
	if err := r.UR.Update(user); err != nil {
		return err
	}
	return nil
}
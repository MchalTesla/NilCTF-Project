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
		return error_code.ErrEmailTaken // 返回邮箱已被占用的错误
	}

	if err := r.UR.Create(user); err != nil {
		return err
	}
	return nil
}

// Login 用户登录
func (r *UserService) Login(email string, username string, password string) (*models.User, error) {
	existingUser, err := r.UR.Read(0, email, username)
	if err != nil {
		if err == error_code.ErrUserNotFound {
			return nil, error_code.ErrUserNotFound // 用户未找到
		}
		return nil, error_code.ErrInternalServer // 处理其他可能的错误
	}

	if !utils.CheckPassword(existingUser.Password, password) {
		return nil, error_code.ErrInvalidInput // 返回密码错误的错误
	}

	return existingUser, nil
}
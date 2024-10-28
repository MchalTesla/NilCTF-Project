package services

import (
	"AWD-Competition-Platform/models"
	"AWD-Competition-Platform/repositories/interface"
	"AWD-Competition-Platform/utils"
	"fmt"
)

type UserService struct {
	UR repositories_interface.UserRepositoryInterface
}

// 返回一个新的Userservice实例
func NewUserService(UR repositories_interface.UserRepositoryInterface) *UserService {
	return &UserService{UR: UR}
}

// 注册一个新用户
func (r *UserService) Register(user models.User) error {
	if err := r.UR.Create(user); err != nil {
		return err
	}
	return nil
}

// 用户登录
func (r *UserService) Login(email string, username string, password string) (models.User, error) {
	existingUser, err := r.UR.Read(0, email, username)
	if err != nil {
		return models.User{}, err
	}

	if !utils.CheckPassword(existingUser.Password, password) {
		return models.User{}, fmt.Errorf("ERR_INCORRECT_PASSWORD")
	}

	return existingUser, nil
}
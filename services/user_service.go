package services

import (
	"NilCTF/dto"
	"NilCTF/error_code"
	"NilCTF/managers/interface"
	"NilCTF/models"
	"NilCTF/utils"
)

// UserService 提供用户相关服务
type UserService struct {
	UM managers_interface.UserManagerInterface
}

// NewUserService 返回一个新的 UserService 实例
func NewUserService(UM managers_interface.UserManagerInterface) *UserService {
	return &UserService{UM: UM}
}

// Register 注册一个新用户
func (r *UserService) Register(creates *dto.UserCreate) error {
	var user models.User

	user.Username = creates.Username
	user.Password = creates.Password
	user.Email = creates.Email
	
	if err := r.UM.Create(&user); err != nil {
		return err
	}
	return nil
}

// Login 用户登录
func (r *UserService) Login(loginIdentifier string, password string) (uint, error) {
	var existingUser *models.User
	var err error

	// 判断 loginIdentifier 是用户名还是邮箱
	if utils.IsValidEmail(loginIdentifier) {
		existingUser, err = r.UM.Get(0, loginIdentifier, "") // 通过邮箱查找用户
	} else if utils.IsValidName(loginIdentifier) {
		existingUser, err = r.UM.Get(0, "", loginIdentifier) // 通过用户名查找用户
	} else {
		return 0, error_code.ErrInvalidInput
	}

	if err != nil {
		return 0, err // 处理其他可能的错误
	}

	if existingUser == nil {
		return 0, error_code.ErrUserNotFound // 返回用户未找到的错误
	}

	if !utils.CheckPassword(existingUser.Password, password) {
		return 0, error_code.ErrInvalidCredentials // 返回密码错误的错误
	}

	switch existingUser.Status {
	case "banned":
		return 0, error_code.ErrUserBanned
	case "pending":
		return 0, error_code.ErrUserPending
	default:
		return existingUser.ID, nil
	}
}

// 修改用户信息
func (r *UserService) Update(userID uint, updates *dto.UserUpdate) error {
	var user models.User

	// 根据传入的字段更新值
	user.ID = userID

	user.Username = updates.Username
	user.Password = updates.Password
	user.Description = updates.Description
	user.Email = updates.Email

	if err := r.UM.Update(&user); err != nil {
		return err
	}
	return nil
}

// 获取当前用户信息
func (r *UserService) GetNow(userID uint) (*dto.UserInfo, error) {
	var userDTO dto.UserInfo

	user, err := r.UM.Get(userID, "", "")
	if err != nil {
		return nil, err
	}
	
	userDTO.Username = user.Username
	userDTO.Email = user.Email
	userDTO.Description = user.Description
	userDTO.Status = user.Status
	userDTO.Role = user.Role
	userDTO.Tag = user.Tag
	userDTO.CreatedAt = user.CreatedAt

	return &userDTO, nil
}

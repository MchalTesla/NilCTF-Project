package managers

import (
	"NilCTF/error_code"
	"NilCTF/models"
	"NilCTF/repositories/interface"
	"NilCTF/utils"
)

type UserManager struct {
	UR repositories_interface.UserRepositoryInterface
	CR repositories_interface.ConfigRepositoryInterface
}

func NewUserManager(UR repositories_interface.UserRepositoryInterface, CR repositories_interface.ConfigRepositoryInterface) *UserManager {
	return &UserManager{UR: UR, CR: CR}
}

func (m *UserManager) Create(user *models.User) error {

	if user.ID != 0 {
		return error_code.ErrInvalidID
	}
	// 校验用户名格式
	if !utils.IsValidName(user.Username) {
		return error_code.ErrInvalidUsername
	}

	// 校验邮箱格式
	if !utils.IsValidEmail(user.Email) {
		return error_code.ErrInvalidEmail
	}

	// 校验描述格式
	if user.Description != nil && !utils.IsValidDescription(*user.Description) {
		return error_code.ErrInvalidDescription
	}

	// 检查用户名和邮箱是否唯一
	if existingUser, err := m.UR.Get(0, "", user.Username); existingUser != nil {
		return error_code.ErrUsernameExists
	} else if err == error_code.ErrInternalServer {
		return err
	}
	if existingUser, err := m.UR.Get(0, user.Email, ""); existingUser != nil {
		return error_code.ErrEmailExists
	} else if err == error_code.ErrInternalServer {
		return err
	}

	// 设置用户状态
	switch user.Status {
	case "active", "banned", "pending":
	case "":
		status, err := m.CR.Get("user_status_default")
		if err != nil {
			return err
		}
		user.Status = status
	default:
		return error_code.ErrInvalidInput
	}

	// 设置用户角色
	switch user.Role {
	case "admin", "user", "organizer":
	case "":
		user.Role = "user"
	default:
		return error_code.ErrInvalidInput
	}

	// 密码哈希化
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return error_code.ErrInternalServer
	}
	user.Password = hashedPassword

	// 将新用户插入数据库
	return m.UR.Create(user)
}

func (m *UserManager) Get(ID uint, email string, username string) (*models.User, error) {
	return m.UR.Get(ID, email, username)
}

func (m *UserManager) Update(user *models.User) error {
	if user.ID == 0 {
		return error_code.ErrInvalidID
	}

	// 校验用户名和邮箱格式
	if user.Username != "" && !utils.IsValidName(user.Username) {
		return error_code.ErrInvalidUsername
	}
	if user.Email != "" && !utils.IsValidEmail(user.Email) {
		return error_code.ErrInvalidEmail
	}

	// 校验描述格式
	if user.Description != nil && !utils.IsValidDescription(*user.Description) {
		return error_code.ErrInvalidDescription
	}

	// 检查用户名和邮箱是否唯一
	if user.Username != "" {
		if existingUser, err := m.UR.Get(0, "", user.Username); existingUser != nil && existingUser.ID != user.ID {
			return error_code.ErrUsernameExists
		} else if err == error_code.ErrInternalServer {
			return err
		}
	}
	if user.Email != "" {
		if existingUser, err := m.UR.Get(0, user.Email, ""); existingUser != nil && existingUser.ID != user.ID {
			return error_code.ErrEmailExists
		} else if err == error_code.ErrInternalServer {
			return err
		}
	}

	// 设置用户状态
	switch user.Status {
	case "active", "banned", "pending", "":
	default:
		return error_code.ErrInvalidInput
	}

	// 设置用户角色
	switch user.Role {
	case "admin", "user", "organizer", "":
	default:
		return error_code.ErrInvalidInput
	}

	// 密码哈希化
	if user.Password != "" {
		hashedPassword, err := utils.HashPassword(user.Password)
		if err != nil {
			return error_code.ErrInternalServer
		}
		user.Password = hashedPassword
	}

	return m.UR.Update(user)
}

func (m *UserManager) Delete(user *models.User) error {
	if user.ID == 0 {
		return error_code.ErrInvalidID
	}
	return m.UR.Delete(user)
}

func (m *UserManager) List(filters map[string]interface{}, limit, offset int, isFuzzy bool) ([]models.User, error) {
	return m.UR.List(filters, limit, offset, isFuzzy)
}

func (m *UserManager) Count() (int64, error) {
	return m.UR.Count()
}
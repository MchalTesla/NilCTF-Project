package services

import (
	"NilCTF/dto"
	managers_interface "NilCTF/managers/interface"
	"NilCTF/models"
)

type AdminUserService struct {
	UM managers_interface.UserManagerInterface
}

// NewAdminUserService 返回一个新的 AdminUserService 实例
func NewAdminUserService(UM managers_interface.UserManagerInterface) *AdminUserService {
	return &AdminUserService{UM: UM}
}

func (s *AdminUserService) GetUsersCount() (int64, error) {
	return s.UM.Count()
}

func (s *AdminUserService) GetTotalPages(limit int) (int64, error) {
	count, err := s.GetUsersCount()
	if err != nil {
		return 0, err
	}
	return (count + int64(limit) - 1) / int64(limit), nil
}

func (s *AdminUserService) ListAllUsers(page int, limit int) ([]dto.UserInfoByAdmin, error) {
	offset := (page - 1) * limit
	users, err := s.UM.List(nil, limit, offset, false)
	if err != nil {
		return nil, err
	}

	var usersDTO []dto.UserInfoByAdmin
	for _, user := range users {
		usersDTO = append(usersDTO, dto.UserInfoByAdmin{
			ID:        user.ID,
			Username:  user.Username,
			Description: user.Description,
			Email:     user.Email,
			Status:    user.Status,
			Role:      user.Role,
			Tag: 	 user.Tag,
			CreatedAt: user.CreatedAt,
		})
	}
	return usersDTO, nil
}

func (s *AdminUserService) CreateUser(user *dto.UserUpdateByAdmin) error {
	var newUser models.User
	newUser.Username = user.Username
	newUser.Password = user.Password
	newUser.Email = user.Email
	newUser.Status = user.Status
	newUser.Role = user.Role
	newUser.Description = user.Description
	newUser.Tag = user.Tag

	return s.UM.Create(&newUser)
}

func (s *AdminUserService) UpdateUsers(updates *dto.UserUpdateByAdmin) error {
	var user models.User

	user.ID = updates.ID
	user.Username = updates.Username
	user.Password = updates.Password
	user.Email = updates.Email
	user.Description = updates.Description
	user.Status = updates.Status
	user.Role = updates.Role
	user.Tag = updates.Tag

	return s.UM.Update(&user)
}

func (s *AdminUserService) DeleteUser(ID uint) error {
	var user models.User
	user.ID = ID
	return s.UM.Delete(&user)
}

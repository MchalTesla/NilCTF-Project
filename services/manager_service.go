package services

import (
	"NilCTF/dto"
	managers_interface "NilCTF/managers/interface"
	"NilCTF/models"
)

type ManagerService struct {
	UM managers_interface.UserManagerInterface
}

// NewManagerService 返回一个新的 ManagerService 实例
func NewManagerService(UM managers_interface.UserManagerInterface) *ManagerService {
	return &ManagerService{UM: UM}
}

func (s *ManagerService) GetUsersCount() (int64, error) {
	return s.UM.Count()
}

func (s *ManagerService) GetTotalPages(limit int) (int64, error) {
	// 获取用户总数
	totalRecords, err := s.GetUsersCount()
	if err != nil {
		return 0, err
	}

	// 计算总页数
	totalPages := (totalRecords + int64(limit) - 1) / int64(limit)
	return totalPages, nil
}

func (s *ManagerService) ListAllUsers(page int, limit int) ([]dto.UserInfoByAdmin, error) {
	var usersDTO []dto.UserInfoByAdmin
	offset := (page - 1) * limit
	users, err := s.UM.List(nil, limit, offset, false)
	if err != nil {
		return nil, err
	}

	// 将users遍历进usersDTO
	for _, user := range users {
		userInfoByAdmin := dto.UserInfoByAdmin{
			ID:          user.ID,
			Username:    user.Username,
			Description: user.Description,
			Email:       user.Email,
			Status:      user.Status,
			Role:        user.Role,
			Tag:         user.Tag,
			CreatedAt:   user.CreatedAt,
		}
		usersDTO = append(usersDTO, userInfoByAdmin)
	}

	return usersDTO, nil
}

func (s *ManagerService) UpdateUsers(updates *dto.UserUpdateByAdmin) error {
	var user models.User

	// 根据传入的字段更新值
	user.ID = updates.ID
	user.Username = updates.Username
	user.Password = updates.Password
	user.Email = updates.Email
	user.Description = updates.Description
	user.Status = updates.Status
	user.Role = updates.Role
	user.Tag = updates.Tag

	if err := s.UM.Update(&user); err != nil {
		return err
	}
	return nil
}

func (s *ManagerService) DeleteUser(ID uint) error {
	var user models.User
	user.ID = ID
	return s.UM.Delete(&user)
}

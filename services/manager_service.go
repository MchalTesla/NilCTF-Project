package services

import (
	"NilCTF/dto"
	"NilCTF/managers/interface"
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

func (s *ManagerService) ListAllUsers(page int, limit int) ([]dto.UserInfo, error) {
	var usersDTO []dto.UserInfo
	offset := (page - 1) * limit
	users, err := s.UM.List(nil, limit, offset, false) 
	if err != nil {
		return nil, err
	}

	// 将users遍历进usersDTO
	for _, user := range users {
		userInfo := dto.UserInfo{
			CreatedAt:	user.CreatedAt,
			Username:	user.Username,
			Description:user.Description,
			Email:		user.Email,
			Status:		user.Status,
			Role:		user.Role,
			Tag:		user.Tag,
		}
		usersDTO = append(usersDTO, userInfo)
	}

	return usersDTO, nil
}
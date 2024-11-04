package services

import (
	"NilCTF/dto"
	"NilCTF/managers/interface"
)

type HomeService struct {
	UM managers_interface.UserManagerInterface
}

// NewUserService 返回一个新的 UserService 实例
func NewHomeService(UM managers_interface.UserManagerInterface) *HomeService {
	return &HomeService{UM: UM}
}

func (r *HomeService) Info(userID uint) (*dto.UserInfo, error) {
	var userInfo dto.UserInfo
	user, err := r.UM.Get(userID, "", "")
	if err != nil {
		return nil, err
	}

	userInfo.CreatedAt = user.CreatedAt
	userInfo.Username = user.Username
	userInfo.Description = user.Description
	userInfo.Email = user.Email
	userInfo.Status = user.Status
	userInfo.Role = user.Role
	userInfo.Tag = user.Tag

	return &userInfo, nil

}

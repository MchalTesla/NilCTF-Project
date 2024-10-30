package services

import (
	"NilCTF/dto"
	"NilCTF/repositories/interface"
)

type HomeService struct {
	UR repositories_interface.UserRepositoryInterface
}

// NewUserService 返回一个新的 UserService 实例
func NewHomeService(UR repositories_interface.UserRepositoryInterface) *HomeService {
	return &HomeService{UR: UR}
}


func (r *HomeService)Info(userID uint) (*dto.UserInfo, error) {
	var userInfo dto.UserInfo
	user, err := r.UR.Read(userID, "", "")
	if err != nil {
		return nil, err
	}

	userInfo.CreatedAt = &user.CreatedAt
	userInfo.Username = &user.Username
	userInfo.Description = &user.Description
	userInfo.Email = &user.Email
	userInfo.Status = &user.Status
	userInfo.Role = &user.Role
	userInfo.Tag = &user.Tag

	return &userInfo, nil

}

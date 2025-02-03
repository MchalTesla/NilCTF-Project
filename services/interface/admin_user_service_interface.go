package services_interface

import (
	"NilCTF/dto"
)



type AdminUserService interface {
	ListAllUsers(page int, limit int) ([]dto.UserInfoByAdmin, error)
	GetUsersCount() (int64, error)
	GetTotalPages(limit int) (int64, error)
	CreateUser(user *dto.UserUpdateByAdmin) error
	UpdateUsers(updates *dto.UserUpdateByAdmin) error
	DeleteUser(ID uint) error
}
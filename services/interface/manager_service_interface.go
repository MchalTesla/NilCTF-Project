package services_interface

import (
	"NilCTF/dto"
)



type ManagerService interface {
	ListAllUsers(page int, limit int) ([]dto.UserInfo, error)
	GetUsersCount() (int64, error)
	GetTotalPages(limit int) (int64, error)
}
package services_interface

import (
	"NilCTF/dto"
)

type UserServiceInterface interface {
	Register(user *dto.UserCreate) error
	Login(loginIdentifier string, password string) (uint, error)
	Update(userID uint, updates *dto.UserUpdate) error
}
package services_interface

import (
	"NilCTF/models"
	"NilCTF/dto"
)

type UserServiceInterface interface {
	Register(user *models.User) error
	Login(loginIdentifier string, password string) (*models.User, error)
	Update(userID uint, updates dto.UserUpdate) error
}
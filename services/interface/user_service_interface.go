package services_interface

import (
	"NilCTF/models"
)

type UserServiceInterface interface {
	Register(user *models.User) error
	Login(loginIdentifier string, password string) (*models.User, error)
	Update(user *models.User) error
}
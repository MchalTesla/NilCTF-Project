package services_interface

import (
	"NilCTF/models"
)

type UserServiceInterface interface {
	Register(user *models.User) error
	Login(email string, username string, password string) (*models.User, error)
}

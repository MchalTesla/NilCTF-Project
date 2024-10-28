package services_interface

import (
	"AWD-Competition-Platform/models"
)

type UserServiceInterface interface {
	Register(user models.User) error
	Login( email string, username string, password string) (models.User, error)
}
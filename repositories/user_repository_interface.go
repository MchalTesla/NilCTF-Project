package repositories

import (
	"AWD-Competition-Platform/models"
)

type UserRepositoryInterface interface {
	Create(user models.User) error
	Read(ID uint, email string, username string) (models.User, error)
	Update(user models.User) error
	Delete(user models.User) error
}
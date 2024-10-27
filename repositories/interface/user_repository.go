package repositories

import (
	"AWD-Competition-Platform/models"
)

type UserRepository interface {
	Create(user models.User) error
	Read(ID uint) (*models.User, error)
	Update(user models.User) error
	Delete(user models.User) error
}
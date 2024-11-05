package managers_interface

import (
	"NilCTF/models"
)

type UserManagerInterface interface {
	Create(user *models.User) error
	Get(ID uint, email string, username string) (*models.User, error)
	Update(user *models.User) error
	Delete(user *models.User) error
	List(filters map[string]interface{}, limit, offset int, isFuzzy bool) ([]models.User, error)
	Count() (int64, error)
}

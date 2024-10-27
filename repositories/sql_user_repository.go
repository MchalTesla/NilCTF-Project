package repositories

import (
	"AWD-Competition-Platform/models"
	"AWD-Competition-Platform/utils"
	"fmt"

	"github.com/jinzhu/gorm"
)

type SQLUserRepository struct {
	db *gorm.DB
}

// NewSQLUserRepository 返回一个新的 SQLUserRepository 实例
func NewSQLUserRepository(db *gorm.DB) *SQLUserRepository {
	return &SQLUserRepository{db: db}
}

func (r *SQLUserRepository) Create(user models.User) error {
	var existingUser models.User
	if err := r.db.Where("username = ?", user.Username).First(&existingUser).Error; err == nil {
		return fmt.Errorf("ERR_USERNAME_EXISTS")
	}
	
	if err := r.db.Where("email = ?", user.Email).First(&existingUser).Error; err == nil {
		return fmt.Errorf("ERR_EMAIL_TAKEN")
	}

	var err error
	user.Password, err = utils.HashPassword(user.Password)
	if err != nil {
		return fmt.Errorf("ERR_INTERNAL_SERVER")
	}

	if err := r.db.Create(&user).Error; err != nil {
		return fmt.Errorf("ERR_INTERNAL_SERVER")
	}
	return nil
}

func (r *SQLUserRepository) Read(ID uint) (*models.User, error) {
	var user models.User

	if err := r.db.First(&user, ID).Error; err != nil {
		return nil, fmt.Errorf("ERR_INTERNAL_SERVER")
	}

	return &user, nil
}

func (r *SQLUserRepository) Update(user models.User) error {
	if err := r.db.Save(&user).Error; err != nil {
		return fmt.Errorf("ERR_INTERNAL_SERVER")
	}
	return nil
}

func (r *SQLUserRepository) Delete(user models.User) error {
	if err := r.db.Delete(user.ID); err != nil {
		return fmt.Errorf("ERR_INTERNAL_SERVER")
	}
	return nil
}
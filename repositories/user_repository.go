package repositories

import (
	"NilCTF/models"
	"NilCTF/utils"
	"fmt"

	"github.com/jinzhu/gorm"
)

type UserRepository struct {
	DB *gorm.DB
}

// NewUserRepository 返回一个新的 UserRepository 实例
func NewUserRepository(DB *gorm.DB) *UserRepository {
	return &UserRepository{DB: DB}
}

func (r *UserRepository) Create(user models.User) error {
	var existingUser models.User

	if err := r.DB.Where("email = ?", user.Email).First(&existingUser).Error; err == nil {
		return fmt.Errorf("ERR_EMAIL_TAKEN")
	}

	if err := r.DB.Where("username = ?", user.Username).First(&existingUser).Error; err == nil {
		return fmt.Errorf("ERR_USERNAME_EXISTS")
	}

	var err error
	user.Password, err = utils.HashPassword(user.Password)
	if err != nil {
		return fmt.Errorf("ERR_INTERNAL_SERVER")
	}

	if err := r.DB.Create(&user).Error; err != nil {
		return fmt.Errorf("ERR_INTERNAL_SERVER")
	}
	return nil
}

func (r *UserRepository) Read(ID uint, email string, username string) (models.User, error) {
	var existingUser models.User

	if ID != 0 {
		if err := r.DB.First(&existingUser, ID).Error; err != nil {
			return models.User{}, fmt.Errorf("USER_NOT_FOUND")
		}
	}

	if username != "" {
		if err := r.DB.Where("username = ?", username).First(&existingUser).Error; err != nil {
			return models.User{}, fmt.Errorf("USER_NOT_FOUND")
		}
	}

	if email != "" {
		if err := r.DB.Where("email = ?", email).First(&existingUser).Error; err != nil {
			return models.User{}, fmt.Errorf("USER_NOT_FOUND")
		}
	}

	return existingUser, nil
}

func (r *UserRepository) Update(user models.User) error {
	if err := r.DB.Save(&user).Error; err != nil {
		return fmt.Errorf("ERR_INTERNAL_SERVER")
	}
	return nil
}

func (r *UserRepository) Delete(user models.User) error {
	if err := r.DB.Delete(user.ID); err != nil {
		return fmt.Errorf("ERR_INTERNAL_SERVER")
	}
	return nil
}

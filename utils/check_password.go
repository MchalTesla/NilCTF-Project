package utils

import (
	"golang.org/x/crypto/bcrypt"
)

// CheckPassword 验证用户密码
func CheckPassword(hashPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(password))
	return err == nil
}

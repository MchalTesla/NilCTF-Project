package utils

import (
	"regexp"
)

// isValidEmail 检查邮箱是否符合格式
func IsValidEmail(email string) bool {
	emailRegex := `^[\w-\.]+@([\w-]+\.)+[\w-]{2,4}$`
	re := regexp.MustCompile(emailRegex)
	return re.MatchString(email)
}
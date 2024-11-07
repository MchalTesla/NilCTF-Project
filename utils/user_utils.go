package utils

import (
	"golang.org/x/crypto/bcrypt"
	"unicode/utf8"
	"regexp"
)

// HashPassword 哈希用户密码
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}


// CheckPassword 验证用户密码
func CheckPassword(hashPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(password))
	return err == nil
}

// IsValidDescription 验证描述的有效性，确保长度不超过 150 个字符
func IsValidDescription(description string) bool {
    // 限制描述长度为 150 个字符
    if len(description) > 150 {
        return false
    }
    return true
}

// isValidEmail 检查邮箱是否符合格式
func IsValidEmail(email string) bool {
	emailRegex := `^\w+(-+.\w+)*@\w+(-.\w+)*.\w+(-.\w+)*$`
	re := regexp.MustCompile(emailRegex)
	return re.MatchString(email)
}

// IsValidName 验证名子的有效性
func IsValidName(username string) bool {
    // 用户名规则：
    // 1. 长度在 3 到 20 个字符之间
    // 2. 只能包含字母、数字和下划线
    // 3. 不能以数字开头
    if utf8.RuneCountInString(username) < 3 || utf8.RuneCountInString(username) > 20 {
        return false
    }

    // 正则表达式匹配
    // 除去电子邮箱格式外，都可以
    return !IsValidEmail(username)
}

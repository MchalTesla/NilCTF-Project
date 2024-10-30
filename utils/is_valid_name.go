package utils

import (
	"regexp"
)

// IsValidName 验证名子的有效性
func IsValidName(username string) bool {
    // 用户名规则：
    // 1. 长度在 3 到 20 个字符之间
    // 2. 只能包含字母、数字和下划线
    // 3. 不能以数字开头
    if len(username) < 3 || len(username) > 20 {
        return false
    }

    // 正则表达式匹配
    // ^[a-zA-Z_][a-zA-Z0-9_]*$ 表示以字母或下划线开头，后面可以跟字母、数字或下划线
    var usernameRegex = regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*$`)
    return usernameRegex.MatchString(username)
}
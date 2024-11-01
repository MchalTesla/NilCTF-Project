package utils

import (
    "unicode/utf8"
)

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
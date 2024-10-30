package utils

// IsValidDescription 验证描述的有效性，确保长度不超过 150 个字符
func IsValidDescription(description string) bool {
    // 限制描述长度为 150 个字符
    if len(description) > 150 {
        return false
    }
    return true
}
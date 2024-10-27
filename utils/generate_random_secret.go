package utils

import (
	"crypto/rand"
)

func GenerateRandomSecret(length int) ([]byte, error) {
	secret := make([]byte, length)
	_, err := rand.Read(secret) // 生成随机字节
	if err != nil {
		return nil, err
	}
	return secret, nil
}

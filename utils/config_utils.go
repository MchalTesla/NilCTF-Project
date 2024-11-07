package utils

import (
	"crypto/rand"
	"NilCTF/error_code"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// Claims 自定义声明结构体，继承 jwt.RegisteredClaims
type Claims struct {
	ID uint `json:"ID"`
	jwt.RegisteredClaims
}


func GenerateRandomSecret(length int) ([]byte, error) {
	secret := make([]byte, length)
	_, err := rand.Read(secret) // 生成随机字节
	if err != nil {
		return nil, err
	}
	return secret, nil
}

// respondWithError 响应错误
func RespondWithError(c *gin.Context, err error) {
	message := err.Error()
	isAPIRequest := strings.HasPrefix(c.FullPath(), "/api/")

	switch err {
	case error_code.ErrPermissionDenied:
		if isAPIRequest {
			c.JSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": message})
		} else {
			c.Redirect(http.StatusFound, "/forbidden")
		}
	case error_code.ErrUserNotLoggedIn, error_code.ErrUserNotFound:
		if isAPIRequest {
			c.JSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": message, "redirect": "/login"})
		} else {
			c.Redirect(http.StatusFound, "/login")
		}
	default:
		if isAPIRequest {
			c.JSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": message})
		} else {
			c.Redirect(http.StatusFound, "/server_error")
		}
	}

	c.Abort()
}

// parseToken 解析 JWT Token 并返回声明
func ParseToken(tokenString string, jwtSecret []byte) (*jwt.Token, *Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil {
		return nil, nil, error_code.ErrInternalServer // 解析错误转换为内部服务器错误
	}

	return token, claims, nil
}
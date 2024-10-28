package middleware

import (
	"NilCTF/config"
	"NilCTF/models"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// Claims 自定义声明结构体，继承 jwt.RegisteredClaims
type Claims struct {
	ID uint `json:"ID"`
	jwt.RegisteredClaims
}

// JWTAuthMiddleware JWT 认证中间件
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" || !strings.HasPrefix(tokenString, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": "未登录，请重新登录", "redirect": "/login"})
			c.Abort()
			return
		}

		tokenString = tokenString[7:] // 去掉 "Bearer " 前缀

		token, claims, err := ParseToken(tokenString)
		if err != nil || !token.Valid || claims.ExpiresAt.Time.Before(time.Now()) {
			c.JSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": "未登录，请重新登录", "redirect": "/login"})
			c.Abort()
			return
		}

		var existingUser models.User
		if err := config.DB.Where("ID = ?", claims.ID).Find(&existingUser).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": "未登录，请重新登录", "redirect": "/login"})
			c.Abort()
			return
		}

		// 将用户信息保存到上下文中
		c.Set("ID", claims.ID)
		c.Next()
	}
}

// ParseToken 解析 JWT Token 并返回声明
func ParseToken(tokenString string) (*jwt.Token, *Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return config.JwtSecret, nil
	})

	if err != nil {
		err = fmt.Errorf("ERR_INTERNAL_SERVER")
	}

	return token, claims, err
}

// 生成 JWT Token
func GenerateToken(ID uint) (string, error) {
	claims := Claims{
		ID: ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 48)), // 设置过期时间
			IssuedAt:  jwt.NewNumericDate(time.Now()),                     // 设置签发时间
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(config.JwtSecret)
	if err != nil {
		return "", fmt.Errorf("ERR_INTERNAL_SERVER")
	}
	return tokenString, nil
}

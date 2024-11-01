package middleware

import (
	"NilCTF/config"
	"NilCTF/error_code"
	"NilCTF/models"
	"NilCTF/repositories"
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

// isValidTokenHeader 检查授权头是否有效
func isValidTokenHeader(header string) bool {
	return header != "" && strings.HasPrefix(header, "Bearer ")
}

// respondWithError 响应错误
func respondWithError(c *gin.Context, err error) {
	message := err.Error()
	c.JSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": message, "redirect": "/login"})
	c.Abort()
}

// parseToken 解析 JWT Token 并返回声明
func parseToken(tokenString string) (*jwt.Token, *Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return config.JwtSecret, nil
	})

	if err != nil {
		return nil, nil, error_code.ErrInternalServer // 解析错误转换为内部服务器错误
	}

	return token, claims, nil
}

type PostMiddleware struct {

}

func NewPostMiddleware() *PostMiddleware {
	return &PostMiddleware{}
}

// JWTAuthMiddleware JWT 认证中间件
// 参数用于验证用户的 JWT 并检查角色权限
// role 参数接受以下值：
// - "all": 允许所有角色访问
// - "admin": 仅允许管理员角色访问
// - "user": 允许用户和管理员角色访问
// - "organizer": 允许比赛创建者访问
func (h *PostMiddleware) JWTAuthMiddleware(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if !isValidTokenHeader(tokenString) {
			respondWithError(c, error_code.ErrInvalidInput)
			return
		}

		tokenString = tokenString[7:] // 去掉 "Bearer " 前缀

		token, claims, err := parseToken(tokenString)
		if err != nil || !token.Valid || claims.ExpiresAt.Time.Before(time.Now()) {
			respondWithError(c, error_code.ErrInvalidInput)
			return
		}

		var user *models.User
		if user, err = repositories.NewUserRepository(config.DB).Read(claims.ID, "", ""); err != nil {
			respondWithError(c, error_code.ErrUserNotFound)
			return
		}

		// 判断用户角色，如果不符合某个角色，就限制访问
		switch role {
		case "all": 
		case "admin":
			if user.Role != "admin" {
				respondWithError(c, error_code.ErrPermissionDenied)
				return
			}
		case "user":
			if user.Role != "user" {
				respondWithError(c, error_code.ErrPermissionDenied)
				return
			}
		case "organizer":
			if user.Role != "organizer" {
				respondWithError(c, error_code.ErrPermissionDenied)
				return
			}
		}

		// 将用户信息保存到上下文中
		c.Set("userID", user.ID)
		c.Set("userName", user.Username)
		c.Set("useremail", user.Email)
		c.Set("userStatus", user.Status)
		c.Set("userRole", user.Role)
		c.Next()
	}
}

// GenerateToken 生成 JWT Token
func (h *PostMiddleware) GenerateToken(ID uint) (string, error) {
	claims := Claims{
		ID: ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(48 * time.Hour)), // 设置过期时间
			IssuedAt:  jwt.NewNumericDate(time.Now()),                     // 设置签发时间
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(config.JwtSecret)
	if err != nil {
		return "", error_code.ErrInternalServer // 生成错误转换为内部服务器错误
	}
	return tokenString, nil
}
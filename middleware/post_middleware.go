package middleware

import (
	"NilCTF/config"
	"NilCTF/error_code"
	managers_interface "NilCTF/managers/interface"
	"NilCTF/models"
	"net/http"
	"time"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// Claims 自定义声明结构体，继承 jwt.RegisteredClaims
type Claims struct {
	ID uint `json:"ID"`
	jwt.RegisteredClaims
}

// respondWithError 响应错误
func respondWithError(c *gin.Context, err error) {
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
	UM managers_interface.UserManagerInterface
}

func NewPostMiddleware(UM managers_interface.UserManagerInterface) *PostMiddleware {
	return &PostMiddleware{UM: UM}
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
		// 从 Cookie 获取 token
		tokenString, err := c.Cookie("auth_token")
		if err != nil {
			respondWithError(c, error_code.ErrUserNotLoggedIn)
			return
		}

		// 解析 token
		token, claims, err := parseToken(tokenString)
		if err != nil || !token.Valid || claims.ExpiresAt.Time.Before(time.Now()) {
			respondWithError(c, error_code.ErrUserNotLoggedIn)
			return
		}

		// 根据 claims.ID 获取用户信息
		var user *models.User
		if user, err = h.UM.Get(claims.ID, "", ""); err != nil {
			respondWithError(c, error_code.ErrUserNotFound)
			return
		}

		// 判断用户角色，如果不符合某个角色，限制访问
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
		c.Set("userEmail", user.Email)
		c.Set("userStatus", user.Status)
		c.Set("userRole", user.Role)
		c.Next()
	}
}

// GenerateToken 生成 JWT Token
func (h *PostMiddleware) GenerateToken(ID uint, jwtTime int) (string, error) {
	claims := Claims{
		ID: ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(jwtTime) * time.Hour)), // 设置过期时间
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

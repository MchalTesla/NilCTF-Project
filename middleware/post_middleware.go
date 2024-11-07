package middleware

import (
	"NilCTF/error_code"
	"NilCTF/managers/interface"
	"NilCTF/models"
	"NilCTF/utils"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)



type PostMiddleware struct {
	UM managers_interface.UserManagerInterface
	jwtSecret []byte
}

func NewPostMiddleware(UM managers_interface.UserManagerInterface, jwtSecret []byte) *PostMiddleware {
	return &PostMiddleware{UM: UM, jwtSecret: jwtSecret}
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
			utils.RespondWithError(c, error_code.ErrUserNotLoggedIn)
			return
		}

		// 解析 token
		token, claims, err := utils.ParseToken(tokenString, h.jwtSecret)
		if err != nil || !token.Valid || claims.ExpiresAt.Time.Before(time.Now()) {
			utils.RespondWithError(c, error_code.ErrUserNotLoggedIn)
			return
		}

		// 根据 claims.ID 获取用户信息
		var user *models.User
		if user, err = h.UM.Get(claims.ID, "", ""); err != nil {
			utils.RespondWithError(c, error_code.ErrUserNotFound)
			return
		}

		// 判断用户角色，如果不符合某个角色，限制访问
		switch role {
		case "all":
		case "admin":
			if user.Role != "admin" {
				utils.RespondWithError(c, error_code.ErrPermissionDenied)
				return
			}
		case "user":
			if user.Role != "user" {
				utils.RespondWithError(c, error_code.ErrPermissionDenied)
				return
			}
		case "organizer":
			if user.Role != "organizer" {
				utils.RespondWithError(c, error_code.ErrPermissionDenied)
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
	claims := utils.Claims{
		ID: ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(jwtTime) * time.Hour)), // 设置过期时间
			IssuedAt:  jwt.NewNumericDate(time.Now()),                     // 设置签发时间
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(h.jwtSecret)
	if err != nil {
		return "", error_code.ErrInternalServer // 生成错误转换为内部服务器错误
	}
	return tokenString, nil
}

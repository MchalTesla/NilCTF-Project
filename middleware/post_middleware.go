package middleware

import (
	"NilCTF/error_code"
	"NilCTF/managers/interface"
	"NilCTF/models"
	"NilCTF/utils"
	"time"
	"strings"
	"net/http"

	"github.com/gin-gonic/gin"
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
			h.RespondWithError(c, error_code.ErrUserNotLoggedIn)
			return
		}

		// 解析 token
		token, claims, err := utils.ParseToken(tokenString, h.jwtSecret)
		if err != nil || !token.Valid || claims.ExpiresAt.Time.Before(time.Now()) {
			h.RespondWithError(c, error_code.ErrUserNotLoggedIn)
			return
		}

		// 根据 claims.ID 获取用户信息
		var user *models.User
		if user, err = h.UM.Get(claims.ID, "", ""); err != nil {
			h.RespondWithError(c, error_code.ErrUserNotFound)
			return
		}

		// 判断用户角色，如果不符合某个角色，限制访问
		switch role {
		case "all":
		case "admin":
			if user.Role != "admin" {
				h.RespondWithError(c, error_code.ErrPermissionDenied)
				return
			}
		case "user":
			if user.Role != "user" {
				h.RespondWithError(c, error_code.ErrPermissionDenied)
				return
			}
		case "organizer":
			if user.Role != "organizer" {
				h.RespondWithError(c, error_code.ErrPermissionDenied)
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

// respondWithError 响应错误
func (h *PostMiddleware) RespondWithError(c *gin.Context, err error) {
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

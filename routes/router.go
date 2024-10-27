package routes

import (
	"AWD-Competition-Platform/controllers"
	"AWD-Competition-Platform/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Setuproutes(r *gin.Engine) {
	r.LoadHTMLGlob("frontend/*.html")

	// 设置静态文件目录
	r.Static("/css", "./frontend/css")
	r.Static("/js", "./frontend/js")

	r.GET("/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", nil)
	})
	r.GET("/register", func(c *gin.Context) {
		c.HTML(http.StatusOK, "register.html", nil)
	})
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})
	r.GET("/index", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	r.POST("/api/register", controllers.Register)
	r.POST("/api/login", controllers.Login)
	r.POST("/api/index", middleware.JWTAuthMiddleware(), controllers.IndexController) // 使用 JWT 中间件保护 Index 路由
	r.POST("/api/logout", controllers.Logout)
	r.POST("api/verify_login", middleware.JWTAuthMiddleware(), controllers.VerifyLogin)
	r.POST("api/competitions", controllers.Competitions)
}

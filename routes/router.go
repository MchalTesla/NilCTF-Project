package routes

import (
	"NilCTF/config"
	"NilCTF/controllers"
	"NilCTF/middleware"
	"NilCTF/repositories"
	"NilCTF/services"
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

	r.POST("/api/register", func(c *gin.Context) {
		controllers.Register(c, services.NewUserService(repositories.NewUserRepository(config.DB)))
	})
	// 登录界面，匿名函数负责初始化一些对象
	r.POST("/api/login", func(c *gin.Context) {
		controllers.Login(c, services.NewUserService(repositories.NewUserRepository(config.DB)))
	})
	r.POST("/api/index", middleware.JWTAuthMiddleware("all"), controllers.IndexController) // 使用 JWT 中间件保护 Index 路由
	r.POST("/api/logout", controllers.Logout)
	r.POST("/api/verify_login", middleware.JWTAuthMiddleware("all"), controllers.VerifyLogin)
	r.POST("/api/competitions", controllers.Competitions)
}

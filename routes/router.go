package routes

import (
	"NilCTF/config"
	"NilCTF/controllers"
	"NilCTF/middleware"
	"NilCTF/repositories"
	"NilCTF/services"
	"net/http"
	"path/filepath"
	"os"

	"github.com/gin-gonic/gin"
)

func Setuproutes(r *gin.Engine) {

	//实例化控制器
	userControllers := &controllers.UserControllers{}
	indexControllers := &controllers.IndexControllers{}
	competitionControllers := &controllers.CompetitionControllers{}
	homeControllers := &controllers.HomeControllers{}

	var files []string
	// 遍历 frontend 目录及其子目录中的所有 .html 文件
	filepath.Walk("frontend", func(path string, info os.FileInfo, err error) error {
		if err == nil && !info.IsDir() && filepath.Ext(path) == ".html" {
			files = append(files, path)
		}
		return nil
	})

	r.LoadHTMLFiles(files...)

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
	r.GET("/home", func(c *gin.Context) {
		c.HTML(http.StatusOK, "home.html", nil)
	})
	r.GET("/home/modify", func(c *gin.Context) {
		c.HTML(http.StatusOK, "modify.html", nil)
	})

	r.POST("/api/register", func(c *gin.Context) {
		userControllers.Register(c, services.NewUserService(repositories.NewUserRepository(config.DB)))
	})
	// 登录界面，匿名函数负责初始化一些对象
	r.POST("/api/login", func(c *gin.Context) {
		userControllers.Login(c, services.NewUserService(repositories.NewUserRepository(config.DB)))
	})
	r.POST("/api/index", middleware.JWTAuthMiddleware("all"), indexControllers.Index) // 使用 JWT 中间件保护 Index 路由
	r.POST("/api/home", middleware.JWTAuthMiddleware("all"), func(c *gin.Context) {
		homeControllers.Home(c, services.NewHomeService(repositories.NewUserRepository(config.DB)))
	})
	r.POST("/api/user/logout", userControllers.Logout)
	r.POST("/api/user/verify", middleware.JWTAuthMiddleware("all"), userControllers.VerifyLogin)
	r.POST("/api/home/modify", middleware.JWTAuthMiddleware("all"), func(c *gin.Context) {
		homeControllers.UpdateUser(c, services.NewUserService(repositories.NewUserRepository(config.DB)))
	})
	r.POST("/api/competition/list_competition", competitionControllers.ListCompetition)
}

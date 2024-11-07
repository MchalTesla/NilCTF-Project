package routes

import (
	"NilCTF/config"
	"NilCTF/controllers"
	"NilCTF/managers"
	"NilCTF/middleware"
	"NilCTF/repositories"
	"NilCTF/services"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// 加载 HTML 文件
func loadHTMLFiles(r *gin.Engine, path string) error {
	var files []string
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err == nil && !info.IsDir() && filepath.Ext(path) == ".html" {
			files = append(files, path)
		}
		return err
	})
	if err != nil {
		return err
	}
	r.LoadHTMLFiles(files...)
	return nil
}

func Setuproutes(r *gin.Engine) {

	// 静态文件目录
	r.Static("/css", "frontend/public/lib/css")
	r.Static("/js", "frontend/public/lib/js")

	// 加载 HTML 模板文件并处理加载错误
	if err := loadHTMLFiles(r, "frontend/public/html"); err != nil {
		panic("Failed to load HTML files: " + err.Error())
	}

	// 实例化服务和存储库
	config := config.NewConfig("config.yaml")
	userRepo := repositories.NewUserRepository(config.Database.DB)
	configRepo := repositories.NewConfigRepository(config.Database.DB)
	userManager := managers.NewUserManager(userRepo, configRepo)
	userService := services.NewUserService(userManager)
	homeService := services.NewHomeService(userRepo)
	managerService := services.NewManagerService(userManager)

	// 控制器初始化
	indexControllers := controllers.NewIndexControllers(userService)
	homeControllers := controllers.NewHomeControllers(homeService)
	competitionControllers := &controllers.CompetitionControllers{}
	managerController := controllers.NewManagerController(managerService)

	// 初始化 Middleware
	preMiddleware := middleware.NewPreMiddleware()
	postMiddleware := middleware.NewPostMiddleware(userManager, config.Jwt.JwtSecret)

	// 初始化user控制器
	userControllers := controllers.NewUserControllers(
		userService, 
		false, 
		config.Jwt.EffectiveDuration, 
		postMiddleware,
	)

	// 配置前置中间件
	r.Use(
		preMiddleware.RateLimitMiddleware(
			rate.Limit(config.Middleware.IPSpeedLimit), 
			config.Middleware.IPSpeedMaxLimit, 
			config.Middleware.IPMaxPlayers,
		),

		preMiddleware.CSPMiddleware(),
	)

	// 页面路由
	r.GET("/login", func(c *gin.Context) { c.HTML(http.StatusOK, "login.html", nil) })
	r.GET("/register", func(c *gin.Context) { c.HTML(http.StatusOK, "register.html", nil) })
	r.GET("/", func(c *gin.Context) { c.HTML(http.StatusOK, "index.html", nil) })
	r.GET("/index", func(c *gin.Context) { c.HTML(http.StatusOK, "index.html", nil) })
	r.GET("/forbidden", func(c *gin.Context) { c.HTML(http.StatusForbidden, "forbidden.html", nil) })
	r.GET("/server_error", func(c *gin.Context) { c.HTML(http.StatusInternalServerError, "server_error.html", nil) })

	homeHTMLGroup := r.Group("/home")
	{
		homeHTMLGroup.Use(postMiddleware.JWTAuthMiddleware("all"))
		{
			homeHTMLGroup.GET("", func(c *gin.Context) { c.HTML(http.StatusOK, "home.html", nil) })
			homeHTMLGroup.GET("/modify", func(c *gin.Context) { c.HTML(http.StatusOK, "modify.html", nil) })
		}
	}

	managerHTMLGroup := r.Group("/manager")
	{
		managerHTMLGroup.Use(postMiddleware.JWTAuthMiddleware("admin"))
		{
			managerHTMLGroup.GET("", func(c *gin.Context) { c.HTML(http.StatusOK, "manager.html", nil) })
			managerHTMLGroup.GET("/users", func(c *gin.Context) { c.HTML(http.StatusOK, "users.html", nil) })
		}
	}

	// 创建 API 路由组并添加预处理中间件
	apiGroup := r.Group("/api")
	{

		// 注册 API 路由
		apiGroup.POST("/register", userControllers.Register)
		apiGroup.POST("/login", userControllers.Login)
		apiGroup.GET("/user/logout", userControllers.Logout)

		// 用户路由组
		userGroup := apiGroup.Group("")
		userGroup.Use(postMiddleware.JWTAuthMiddleware("all"))
		{
			userGroup.GET("/index", indexControllers.Index)
			userGroup.GET("/home", homeControllers.Home)
			userGroup.POST("/home/modify", func(c *gin.Context) { homeControllers.Modify(c, userService) })
			userGroup.GET("/verify", userControllers.VerifyLogin)
		}

		// 管理员路由组
		adminGroup := apiGroup.Group("/manager")
		adminGroup.Use(postMiddleware.JWTAuthMiddleware("admin"))
		{
			adminGroup.GET("/users_count", managerController.GetUsersCount)
			adminGroup.POST("/list_users", managerController.ListUsers)
			adminGroup.POST("/update", managerController.UpdateUserByAdmin)
			adminGroup.POST("/delete", managerController.DeleteUserByAdmin)
		}

		// 不受保护的比赛列表路由
		apiGroup.GET("/competition/list_competition", competitionControllers.ListCompetition)
	}
}

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
	ManagerService := services.NewManagerService(userManager)
	announcementRepo := repositories.NewAnnouncementRepository(config.Database.DB)
	announcementManager := managers.NewAnnouncementManager(announcementRepo)
	announcementService := services.NewAnnouncementService(announcementManager)

	// 公告相关路由
	announcementController := controllers.NewAnnouncementController(announcementService)
	adminAnnouncementController := controllers.NewAdminAnnouncementController(announcementService)

	// 控制器初始化
	indexControllers := controllers.NewIndexControllers(userService)
	userControllers := controllers.NewUserControllers(userService, false, config.Jwt.EffectiveDuration, config.Jwt.JwtSecret)
	homeControllers := controllers.NewHomeControllers(homeService)
	competitionControllers := &controllers.CompetitionControllers{}
	adminManagerController := controllers.NewAdminManagerController(ManagerService)

	// 初始化 Middleware
	preMiddleware := middleware.NewPreMiddleware()
	postMiddleware := middleware.NewPostMiddleware(userManager, config.Jwt.JwtSecret)

	// 部署前置基于IP的速度控制器中间件
	if config.Middleware.StartIPSpeedLimit == true {
		r.Use(
			preMiddleware.RateLimitMiddleware(
				rate.Limit(config.Middleware.IPSpeedLimit),
				config.Middleware.IPSpeedMaxLimit,
				config.Middleware.IPMaxPlayers,
			),
		)
	}
	// 部署前置CSP安全规则中间件
	if config.Middleware.StartCSP == true {
		r.Use(preMiddleware.CSPMiddleware(config.Middleware.CSPValue))
	}

	/*
	页面路由
	*/

	r.GET("/login", func(c *gin.Context) { c.HTML(http.StatusOK, "login.html", nil) })
	r.GET("/register", func(c *gin.Context) { c.HTML(http.StatusOK, "register.html", nil) })
	r.GET("/", func(c *gin.Context) { c.HTML(http.StatusOK, "index.html", nil) })
	r.GET("/index", func(c *gin.Context) { c.HTML(http.StatusOK, "index.html", nil) })
	r.GET("/forbidden", func(c *gin.Context) { c.HTML(http.StatusForbidden, "forbidden.html", nil) })
	r.GET("/server_error", func(c *gin.Context) { c.HTML(http.StatusInternalServerError, "server_error.html", nil) })
	r.GET("/announcements", func(c *gin.Context) { c.HTML(http.StatusOK, "announcement.html", nil) })

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
			managerHTMLGroup.GET("/admin_users", func(c *gin.Context) { c.HTML(http.StatusOK, "admin_users.html", nil) })
			managerHTMLGroup.GET("/admin_announcements", func(c *gin.Context) { c.HTML(http.StatusOK, "admin_announcements.html", nil) })
		}
	}

	// 创建 API 路由组并添加预处理中间件
	apiGroup := r.Group("/api")
	{
		apiGroup.POST("/register", userControllers.Register)
		apiGroup.POST("/login", userControllers.Login)
		apiGroup.GET("/user/logout", userControllers.Logout)
		apiGroup.GET("/announcements", announcementController.ListAnnouncements)

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
			adminGroup.GET("/users_count", adminManagerController.GetUsersCount)
			adminGroup.POST("/list_users", adminManagerController.ListUsers)
			adminGroup.POST("/users", adminManagerController.HandleUser) // 合并用户相关操作
			adminGroup.POST("/announcements", adminAnnouncementController.HandleAnnouncement) // 合并公告相关操作
		}

		// 不受保护的比赛列表路由
		apiGroup.GET("/competition/list_competition", competitionControllers.ListCompetition)
	}
}

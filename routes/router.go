package routes

import (
	"NilCTF/config"
	"NilCTF/controllers"
	"NilCTF/middleware"
	"NilCTF/repositories"
	"NilCTF/managers"
	"NilCTF/services"
	"net/http"
	"path/filepath"
	"os"

	"github.com/gin-gonic/gin"
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
	// 初始化 Middleware
	preMiddleware := middleware.NewPreMiddleware()
	postMiddleware := middleware.NewPostMiddleware()

	// 静态文件目录
	r.Static("/css", "./frontend/css")
	r.Static("/js", "./frontend/js")

	// 加载 HTML 模板文件并处理加载错误
	if err := loadHTMLFiles(r, "frontend"); err != nil {
		panic("Failed to load HTML files: " + err.Error())
	}

	// 实例化服务和存储库
	userRepo := repositories.NewUserRepository(config.DB)
	configRepo := repositories.NewConfigRepository(config.DB)
	userManager := managers.NewUserManager(userRepo, configRepo)
	userService := services.NewUserService(userManager)
	homeService := services.NewHomeService(userRepo)

	// 控制器初始化
	userControllers := &controllers.UserControllers{}
	indexControllers := &controllers.IndexControllers{}
	competitionControllers := &controllers.CompetitionControllers{}
	homeControllers := &controllers.HomeControllers{}

	// 页面路由
	r.GET("/login", func(c *gin.Context) { c.HTML(http.StatusOK, "login.html", nil) })
	r.GET("/register", func(c *gin.Context) { c.HTML(http.StatusOK, "register.html", nil) })
	r.GET("/", func(c *gin.Context) { c.HTML(http.StatusOK, "index.html", nil) })
	r.GET("/index", func(c *gin.Context) { c.HTML(http.StatusOK, "index.html", nil) })
	r.GET("/home", func(c *gin.Context) { c.HTML(http.StatusOK, "home.html", nil) })
	r.GET("/home/modify", func(c *gin.Context) { c.HTML(http.StatusOK, "modify.html", nil) })

	// 创建 API 路由组并添加预处理中间件
	apiGroup := r.Group("/api")
	{
		apiGroup.Use(
			preMiddleware.RateLimitMiddleware(10, 20, 5000),
			preMiddleware.CSPMiddleware(),
		)

		// 注册 API 路由
		apiGroup.POST("/register", func(c *gin.Context) { userControllers.Register(c, userService) })
		apiGroup.POST("/login", func(c *gin.Context) { userControllers.Login(c, userService) })
		apiGroup.GET("/user/logout", userControllers.Logout)

		// 受保护路由组
		protectedGroup := apiGroup.Group("")
		protectedGroup.Use(postMiddleware.JWTAuthMiddleware("all", userManager))
		{
			protectedGroup.GET("/index", indexControllers.Index)
			protectedGroup.GET("/home", func(c *gin.Context) { homeControllers.Home(c, homeService) })
			protectedGroup.GET("/user/verify", userControllers.VerifyLogin)
			protectedGroup.POST("/home/modify", func(c *gin.Context) { homeControllers.Modify(c, userService) })
		}

		// 不受保护的比赛列表路由
		apiGroup.GET("/competition/list_competition", competitionControllers.ListCompetition)
	}
}
package routes

import (
	"NilCTF/controllers"
	"NilCTF/managers"
	"NilCTF/middleware"
	"NilCTF/repositories"
	"NilCTF/services"

	"golang.org/x/time/rate"
)

// 注册repositories
func (app *App) initRepositories() {
    app.Container.Repositories["user"] = repositories.NewUserRepository(app.Config.Database.DB)
    app.Container.Repositories["config"] = repositories.NewConfigRepository(app.Config.Database.DB)
    app.Container.Repositories["announcement"] = repositories.NewAnnouncementRepository(app.Config.Database.DB)
}

// 注册managers
func (app *App) initManagers() {
    app.Container.Managers["user"] = managers.NewUserManager(
        app.Container.Repositories["user"].(*repositories.UserRepository),
        app.Container.Repositories["config"].(*repositories.ConfigRepository),
    )
    app.Container.Managers["announcement"] = managers.NewAnnouncementManager(
        app.Container.Repositories["announcement"].(*repositories.AnnouncementRepository),
    )
}

// 注册services (按需添加)
func (app *App) initServices() {
    app.Container.Services["user"] = services.NewUserService(
        app.Container.Managers["user"].(*managers.UserManager),
    )
	app.Container.Services["home"] = services.NewHomeService(
		app.Container.Repositories["user"].(*repositories.UserRepository),
	)
	app.Container.Services["admin_user"] = services.NewAdminUserService(
		app.Container.Managers["user"].(*managers.UserManager),
	)
    app.Container.Services["announcement"] = services.NewAnnouncementService(
        app.Container.Managers["announcement"].(*managers.AnnouncementManager),
        app.Container.Managers["user"].(*managers.UserManager),
    )
}

// 注册controllers (按需添加)
func (app *App) initControllers() {
	app.Container.Controllers["index"] = controllers.NewIndexControllers(
		app.Container.Services["user"].(*services.UserService),
	)
    app.Container.Controllers["user"] = controllers.NewUserControllers(
        app.Container.Services["user"].(*services.UserService),
        false,
        app.Config.Jwt.EffectiveDuration,
        app.Config.Jwt.JwtSecret,
    )
	app.Container.Controllers["home"] = controllers.NewHomeControllers(
		app.Container.Services["home"].(*services.HomeService),
	)
	app.Container.Controllers["competition"] = controllers.NewCompetitionController()
    app.Container.Controllers["announcement"] = controllers.NewAnnouncementController(
        app.Container.Services["announcement"].(*services.AnnouncementService),
    )
	app.Container.Controllers["admin_user"] = controllers.NewAdminUserController(
		app.Container.Services["admin_user"].(*services.AdminUserService),
	)
	app.Container.Controllers["admin_announcement"] = controllers.NewAdminAnnouncementController(
		app.Container.Services["announcement"].(*services.AnnouncementService),
	)
}

func (app *App) initMiddleware() {
    app.Container.Middleware.Pre = middleware.NewPreMiddleware()
    app.Container.Middleware.Post = middleware.NewPostMiddleware(
        app.Container.Managers["user"].(*managers.UserManager),
        app.Config.Jwt.JwtSecret,
    )
    
    // 配置中间件
    if app.Config.Middleware.StartIPSpeedLimit {
        app.Engine.Use(
            app.Container.Middleware.Pre.RateLimitMiddleware(
                rate.Limit(app.Config.Middleware.IPSpeedLimit),
                app.Config.Middleware.IPSpeedMaxLimit,
                app.Config.Middleware.IPMaxPlayers,
            ),
        )
    }
}
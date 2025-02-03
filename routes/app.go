package routes

import (
	"NilCTF/config"
	"NilCTF/middleware"

	"github.com/gin-gonic/gin"
)

type App struct {
    Config     *config.Config
    Engine     *gin.Engine
    Container  *ServiceContainer
}

type ServiceContainer struct {
    Repositories map[string]interface{}
    Managers     map[string]interface{}
    Services     map[string]interface{}
    Controllers  map[string]interface{}
    Middleware   struct {
        Pre  *middleware.PreMiddleware
        Post *middleware.PostMiddleware
    }
}

func NewApp(r *gin.Engine, configPath string) *App {
    app := &App{
        Config: config.NewConfig(configPath),
        Engine: r,
        Container: &ServiceContainer{
            Repositories: make(map[string]interface{}),
            Managers:    make(map[string]interface{}),
            Services:    make(map[string]interface{}),
            Controllers: make(map[string]interface{}),
        },
    }
    
    app.initRepositories()
    app.initManagers()
    app.initServices()
    app.initControllers()
    app.initMiddleware()
    return app
}
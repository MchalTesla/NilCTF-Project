package routes

import "github.com/gin-gonic/gin"

type Router struct {
	App *App
	htmlRoutes []HTMLRoute
	apiPublicRoutes []APIRoute
	apiUserRoutes []APIRoute
	apiAdminRoutes []APIRoute
	Engine     *gin.Engine
}

// 实例化一个新的Router
func NewRouter(r *gin.Engine) *Router {
	router := &Router{
		App: NewApp(r, "config.yaml"),
		Engine: r,
	}
	router.initHTMLRoutes()
	router.RegisterHTMLRoutes()
	router.initAPIRoutes()
	router.RegisterAPIRoutes()
	router.OtherRoute()
	return router
}
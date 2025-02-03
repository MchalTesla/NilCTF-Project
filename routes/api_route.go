package routes

import (
	"NilCTF/controllers"
	"NilCTF/services"

	"github.com/gin-gonic/gin"
)

type APIRoute struct {
    Path        string
    Method      string
    Handler     gin.HandlerFunc
    Middleware  []gin.HandlerFunc
}

func (r *Router) initAPIRoutes() {
	        // 公共路由
    r.apiPublicRoutes = []APIRoute{
        {"/api/register", "POST", r.App.Container.Controllers["user"].(*controllers.UserControllers).Register, nil},
        {"/api/login", "POST", r.App.Container.Controllers["user"].(*controllers.UserControllers).Login, nil},
        {"/api/user/logout", "GET", r.App.Container.Controllers["user"].(*controllers.UserControllers).Logout, nil},
        {"/api/announcements", "GET", r.App.Container.Controllers["announcement"].(*controllers.AnnouncementController).ListAnnouncements, nil},
		{"/api/competition/list_competition", "GET", r.App.Container.Controllers["competition"].(*controllers.CompetitionController).ListCompetition, nil},
	}

        // 用户路由
	r.apiUserRoutes = []APIRoute{
        {"/api/index", "GET", r.App.Container.Controllers["index"].(*controllers.IndexControllers).Index, nil},
        {"/api/home", "GET", r.App.Container.Controllers["home"].(*controllers.HomeControllers).Home, nil},
        {"/api/home/modify", "POST", func(c *gin.Context) { r.App.Container.Controllers["home"].(*controllers.HomeControllers).Modify(c, r.App.Container.Services["user"].(*services.UserService)) }, nil},
        {"/api/verify", "GET", r.App.Container.Controllers["user"].(*controllers.UserControllers).VerifyLogin, nil},
	}

        // 管理员路由
	r.apiAdminRoutes = []APIRoute{
		{"/api/admin/users_count", "GET", r.App.Container.Controllers["admin_user"].(*controllers.AdminUserController).GetUsersCount, nil},
        {"/api/admin/list_users", "POST", r.App.Container.Controllers["admin_user"].(*controllers.AdminUserController).ListUsers, nil},
        {"/api/admin/users", "POST", r.App.Container.Controllers["admin_user"].(*controllers.AdminUserController).HandleUser, nil},
        {"/api/admin/announcements", "POST", r.App.Container.Controllers["admin_announcement"].(*controllers.AdminAnnouncementController).HandleAnnouncement, nil},
    }
}

// RegisterAPIRoutes 注册API路由
func (r *Router) RegisterAPIRoutes() {
	// 注册公共路由
	for _, route := range r.apiPublicRoutes {
		handlers := []gin.HandlerFunc{}
		if route.Middleware != nil {
			handlers = append(handlers, route.Middleware...)
		}
		handlers = append(handlers, route.Handler)
		switch route.Method {
		case "GET":
			r.Engine.GET(route.Path, handlers...)
		case "POST":
			r.Engine.POST(route.Path, handlers...)
		}
	}

	// 注册用户路由
	for _, route := range r.apiUserRoutes {
		handlers := []gin.HandlerFunc{}
		handlers = append(handlers, r.App.Container.Middleware.Post.JWTAuthMiddleware("all"))
		if route.Middleware != nil {
			handlers = append(handlers, route.Middleware...)
		}
		handlers = append(handlers, route.Handler)
		switch route.Method {
		case "GET":
			r.Engine.GET(route.Path, handlers...)
		case "POST":
			r.Engine.POST(route.Path, handlers...)
		}
	}

	// 注册管理员路由
	for _, route := range r.apiAdminRoutes {
		handlers := []gin.HandlerFunc{}
		handlers = append(handlers, r.App.Container.Middleware.Post.JWTAuthMiddleware("admin"))
		if route.Middleware != nil {
			handlers = append(handlers, route.Middleware...)
		}
		handlers = append(handlers, route.Handler)
		switch route.Method {
		case "GET":
			r.Engine.GET(route.Path, handlers...)
		case "POST":
			r.Engine.POST(route.Path, handlers...)
		}
	}

}
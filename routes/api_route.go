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
	Role        string
}

func (r *Router) initAPIRoutes() {
	// API路由配置
    r.apiRoutes = []APIRoute{
		// 公共路由
        {"/api/register", "POST", r.App.Container.Controllers["user"].(*controllers.UserControllers).Register, nil, "all"},
        {"/api/login", "POST", r.App.Container.Controllers["user"].(*controllers.UserControllers).Login, nil, "all"},
        {"/api/user/logout", "GET", r.App.Container.Controllers["user"].(*controllers.UserControllers).Logout, nil,	"all"},
        {"/api/announcements", "GET", r.App.Container.Controllers["announcement"].(*controllers.AnnouncementController).ListAnnouncements, nil, "all"},
		{"/api/competition/list_competition", "GET", r.App.Container.Controllers["competition"].(*controllers.CompetitionController).ListCompetition, nil, "all"},

    	// 用户路由
        {"/api/index", "GET", r.App.Container.Controllers["index"].(*controllers.IndexControllers).Index, nil, "user"},
        {"/api/home", "GET", r.App.Container.Controllers["home"].(*controllers.HomeControllers).Home, nil, "user"},
        {"/api/home/modify", "POST", func(c *gin.Context) { r.App.Container.Controllers["home"].(*controllers.HomeControllers).Modify(c, r.App.Container.Services["user"].(*services.UserService)); c.Next() }, nil, "user"},
        {"/api/verify", "GET", r.App.Container.Controllers["user"].(*controllers.UserControllers).VerifyLogin, nil, "user"},

    	// 管理员路由
		{"/api/admin/users_count", "GET", r.App.Container.Controllers["admin_user"].(*controllers.AdminUserController).GetUsersCount, nil, "admin"},
        {"/api/admin/list_users", "POST", r.App.Container.Controllers["admin_user"].(*controllers.AdminUserController).ListUsers, nil, "admin"},
        {"/api/admin/users", "POST", r.App.Container.Controllers["admin_user"].(*controllers.AdminUserController).HandleUser, nil, "admin"},
        {"/api/admin/announcements", "POST", r.App.Container.Controllers["admin_announcement"].(*controllers.AdminAnnouncementController).HandleAnnouncement, nil, "admin"},
    }
}

// RegisterAPIRoutes 注册API路由
func (r *Router) RegisterAPIRoutes() {
	// 注册API路由
	for _, route := range r.apiRoutes {
		handlers := []gin.HandlerFunc{}
		handlers = append(handlers, r.App.Container.Middleware.Post.JWTAuthMiddleware(route.Role))
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
package routes

import (
	"github.com/gin-gonic/gin"
)

// HTMLRoute 定义HTML路由配置
type HTMLRoute struct {
    Path        string   // 路由路径
    FilePath   string   // HTML文件路径
    Middleware []gin.HandlerFunc // 中间件
	Role 	 string   // 角色
}

func (r *Router) initHTMLRoutes() {
	// 定义HTML路由配置
	r.htmlRoutes = []HTMLRoute{
		{"/login", "login.html", nil, "all"},
		{"/register", "register.html", nil, "all"},
		{"/", "index.html", nil, "all"},
		{"/index", "index.html", nil, "all"},
		{"/forbidden", "forbidden.html", nil, "all"},
		{"/server_error", "server_error.html", nil, "all"},
		{"/announcements", "announcements.html", nil, "all"},
		
		// 需要用户认证的路由
		{"/home", "home.html", nil, "user"},
		{"/home/modify", "modify.html", nil, "user"},
		
		// 管理员路由
		{"/admin", "", []gin.HandlerFunc{func(c *gin.Context) { c.Redirect(302, "/admin/index") }}, "admin"},
		{"/admin/index", "admin/index.html", nil, "admin"},
		{"/admin/users", "admin/admin_users.html", nil, "admin"},
		{"/admin/announcements", "admin/admin_announcements.html", nil, "admin"},
	}
}

func (r *Router) RegisterHTMLRoutes() {
	for _, route := range r.htmlRoutes {
		handlers := []gin.HandlerFunc{}
		handlers = append(handlers, r.App.Container.Middleware.Post.JWTAuthMiddleware(route.Role))
		if route.Middleware != nil {
			handlers = append(handlers, route.Middleware...)
		}
		if route.FilePath != "" {
			handlers = append(handlers, func(c *gin.Context) {
				c.File("frontend/public/html/" + route.FilePath)
			})
		}
		r.Engine.GET(route.Path, handlers...)
	}
}

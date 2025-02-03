package routes

import (
	"github.com/gin-gonic/gin"
)

// HTMLRoute 定义HTML路由配置
type HTMLRoute struct {
    Path        string   // 路由路径
    FilePath   string   // HTML文件路径
    Middleware []gin.HandlerFunc // 中间件
}

func (r *Router) initHTMLRoutes() {
	// 定义HTML路由配置
	r.htmlRoutes = []HTMLRoute{
		{"/login", "login.html", nil},
		{"/register", "register.html", nil},
		{"/", "index.html", nil},
		{"/index", "index.html", nil},
		{"/forbidden", "forbidden.html", nil},
		{"/server_error", "server_error.html", nil},
		{"/announcements", "announcements.html", nil},
		
		// 需要用户认证的路由
		{"/home", "home.html", []gin.HandlerFunc{r.App.Container.Middleware.Post.JWTAuthMiddleware("all")}},
		{"/home/modify", "modify.html", []gin.HandlerFunc{r.App.Container.Middleware.Post.JWTAuthMiddleware("all")}},
		
		// 管理员路由
		{"/admin", "", []gin.HandlerFunc{func(c *gin.Context) { c.Redirect(302, "/admin/index") }}},
		{"/admin/index", "admin/index.html", []gin.HandlerFunc{r.App.Container.Middleware.Post.JWTAuthMiddleware("all")}},
		{"/admin/users", "admin/admin_users.html", []gin.HandlerFunc{r.App.Container.Middleware.Post.JWTAuthMiddleware("all")}},
		{"/admin/announcements", "admin/admin_announcements.html", []gin.HandlerFunc{r.App.Container.Middleware.Post.JWTAuthMiddleware("all")}},
	}
}

func (r *Router) RegisterHTMLRoutes() {
	for _, route := range r.htmlRoutes {
		handlers := []gin.HandlerFunc{}
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

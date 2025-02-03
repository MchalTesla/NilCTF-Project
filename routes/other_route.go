package routes

func (r *Router) OtherRoute() {
	// 注册静态路由
	r.Engine.Static("/css", "frontend/public/lib/css")
	r.Engine.Static("/js", "frontend/public/lib/js")
}
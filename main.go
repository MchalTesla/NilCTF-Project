package main

import (
	_ "NilCTF/config"
	"NilCTF/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	// 使用gin创建一个新的路由器
	r := gin.Default()

	// 设置路由
	routes.Setuproutes(r)

	// 启动服务器
	r.Run(":8080") // 启动服务，监听8080端口
}

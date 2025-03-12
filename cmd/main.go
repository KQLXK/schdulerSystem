package main

import (
	"schedule/database"
	"schedule/models"
	"schedule/route"
)

func main() {
	// 初始化配置
	//config.LoadConfig()

	// 初始化数据库
	database.InitDB()
	//渲染数据表
	models.InitTables()

	r := route.SetupRoute()

	r.Run(":8080")

	// 创建 Gin 引擎
	//r := gin.Default()

	// 注册中间件
	//r.Use(middleware.LoggerMiddleware()) // 日志中间件（可选）
	//r.Use(middleware.CORSMiddleware())   // CORS 中间件（可选）

	// 注册路由
	//handlers.RegisterRoutes(r)

	// 启动 HTTP 服务
	//if err := r.Run(":8080"); err != nil {
	//	panic("Failed to start server: " + err.Error())
	//}
}

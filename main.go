package main

import (
	"github.com/gin-gonic/gin"
	"schedule/config"
	"schedule/database"
	"schedule/handlers"
	"schedule/middleware"
)

func main() {
	// 初始化配置
	//config.LoadConfig()

	// 初始化数据库
	database.InitDB()

	// 创建 Gin 引擎
	r := gin.Default()

	// 注册中间件
	r.Use(middleware.LoggerMiddleware()) // 日志中间件（可选）
	r.Use(middleware.CORSMiddleware())   // CORS 中间件（可选）

	// 注册路由
	handlers.RegisterRoutes(r)

	// 启动 HTTP 服务
	if err := r.Run(":" + config.ServerPort); err != nil {
		panic("Failed to start server: " + err.Error())
	}
}

package main

import (
	"schedule/database"
	"schedule/models"
	"schedule/route"
)

func main() {

	// 初始化数据库
	database.InitDB()
	//渲染数据表
	models.InitTables()

	r := route.SetupRoute()

	r.Run(":8080")

	//algorithm.GA()

}

package route

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"schedule/handlers"
)

func SetupRoute() *gin.Engine {

	r := gin.Default()

	// 配置 CORS 中间件
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"content-type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	courseGroup := r.Group("/course")
	{
		courseGroup.POST("/create", handlers.AddCourse)
		courseGroup.DELETE("/delete/:course_id", handlers.DeleteCourse)
		courseGroup.PUT("/update", handlers.UpdateCourse)
		courseGroup.GET("/query/:course_id", handlers.GetCourse)
		courseGroup.GET("/queryall", handlers.GetAllCourses)
	}

	return r

}

package route

import (
	"github.com/gin-gonic/gin"
	"schedule/handlers"
)

func SetupRoute() *gin.Engine {

	r := gin.Default()

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

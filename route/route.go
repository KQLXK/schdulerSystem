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
	}

	return r

}

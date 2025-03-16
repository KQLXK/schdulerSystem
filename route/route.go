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
	teacherGroup := r.Group("/teacher")
	{
		teacherGroup.POST("/create", handlers.AddTeacher)
		teacherGroup.DELETE("/delete/:id", handlers.DeleteTeacher)
		teacherGroup.PUT("/update/:id", handlers.UpdateTeacher)
		teacherGroup.GET("/query/:id", handlers.GetTeacherByID)
		teacherGroup.GET("/queryall", handlers.GetTeachers)
	}
	classroomGroup := r.Group("/classroom")
	{
		classroomGroup.POST("/create", handlers.AddClassroom)
		classroomGroup.DELETE("/delete/:id", handlers.DeleteClassroom)
		classroomGroup.PUT("/update/:id", handlers.UpdateClassroom)
		classroomGroup.GET("/query/:id", handlers.GetClassroomByID)
		classroomGroup.GET("/queryall", handlers.GetClassrooms)
	}
	return r

}

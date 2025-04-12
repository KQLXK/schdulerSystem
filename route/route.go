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

	userGroup := r.Group("/user")
	{
		userGroup.POST("/login", handlers.UserLoginHandler)
		//userGroup.POST("/register", handlers.RegisterHandler)
	}

	courseGroup := r.Group("/course")
	{
		courseGroup.POST("/create", handlers.AddCourse)
		courseGroup.POST("/create/file", handlers.AddCourseByExcel)
		courseGroup.DELETE("/delete/:course_id", handlers.DeleteCourse)
		courseGroup.PUT("/update", handlers.UpdateCourse)
		courseGroup.GET("/query/:course_id", handlers.GetCourse)
		courseGroup.GET("/queryall", handlers.GetAllCourses)
		courseGroup.GET("/querybypage", handlers.QueryCourseByPage)
		courseGroup.GET("/search", handlers.SearchCourse)
	}
	teacherGroup := r.Group("/teacher")
	{
		teacherGroup.POST("/create", handlers.AddTeacher)
		teacherGroup.DELETE("/delete/:id", handlers.DeleteTeacher)
		teacherGroup.PUT("/update/:id", handlers.UpdateTeacher)
		teacherGroup.GET("/query/:id", handlers.GetTeacherByID)
		teacherGroup.GET("/queryall", handlers.GetTeachers)
		teacherGroup.GET("/querybypage", handlers.QueryTeacherByPage)
		teacherGroup.GET("/search", handlers.SearchTeacher)
		teacherGroup.POST("/create/file", handlers.AddTeacherByExcel)
	}
	classroomGroup := r.Group("/classroom")
	{
		classroomGroup.POST("/create", handlers.AddClassroom)
		classroomGroup.DELETE("/delete/:id", handlers.DeleteClassroom)
		classroomGroup.PUT("/update/:id", handlers.UpdateClassroom)
		classroomGroup.GET("/query/:id", handlers.GetClassroomByID)
		classroomGroup.GET("/queryall", handlers.GetClassrooms)
		classroomGroup.GET("/querybypage", handlers.QueryClassroomByPage)
		classroomGroup.GET("/search", handlers.SearchClassroom)
		classroomGroup.POST("/create/file", handlers.AddClassroomByExcel)
	}

	schedulegroup := r.Group("/schedule")
	{
		schedulegroup.POST("/create", handlers.CreateSchedule)
		//group.POST("/import", handlers.)
		schedulegroup.PUT("/update/:schedule_id", handlers.UpdateSchedule)
		schedulegroup.DELETE("/delete/:schedule_id", handlers.DeleteSchedule)
		schedulegroup.GET("/query/:schedule_id", handlers.GetSchedule)
		schedulegroup.GET("/queryall", handlers.GetAllSchedules)
		schedulegroup.GET("/querybypage", handlers.QuerySchedulesByPage)
		schedulegroup.POST("/create/file", handlers.AddScheduleByExcel)
		schedulegroup.POST("/ga", handlers.GAHandler)
		//group.GET("/search", handlers.SearchSchedules)
	}

	classGroup := r.Group("/class")
	{
		classGroup.POST("/create", handlers.AddClass)
		classGroup.POST("/create/file", handlers.AddClassByExcel)
		classGroup.DELETE("/delete/:id", handlers.DeleteClass)
		classGroup.PUT("/update/:id", handlers.UpdateClass)
		classGroup.GET("/query/:id", handlers.GetClassByID)
		classGroup.GET("/queryall", handlers.GetAllClasses)
		classGroup.GET("/querybypage", handlers.QueryClassByPage)
		classGroup.GET("/search", handlers.SearchClass)
	}

	tableGroup := r.Group("/table")
	{
		tableGroup.GET("/class", handlers.GetClassTableHandler)
		tableGroup.GET("/teacher", handlers.GetTeacherTableHandler)
		tableGroup.GET("/classroom", handlers.GetClassroomTableHandler)
		tableGroup.POST("/create", handlers.ManualScheduleHandler)
	}

	return r

}

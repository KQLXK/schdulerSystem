package handlers

import (
	"github.com/gin-gonic/gin"
	"schedule/middleware"
)

// RegisterRoutes 注册所有 API 路由
func RegisterRoutes(r *gin.Engine) {
	// 公开路由（无需认证）
	r.POST("/login", LoginHandler) // 登录接口

	// 受保护的路由组（需要认证）
	api := r.Group("/api")
	api.Use(middleware.AuthMiddleware()) // 使用 JWT 认证中间件
	{
		// 教师相关路由
		api.GET("/teachers", GetTeachers)          // 获取所有教师
		api.POST("/teachers", AddTeacher)          // 添加教师
		api.PUT("/teachers/:id", UpdateTeacher)    // 更新教师信息
		api.DELETE("/teachers/:id", DeleteTeacher) // 删除教师

		// 教室相关路由
		api.GET("/classrooms", GetClassrooms)          // 获取所有教室
		api.POST("/classrooms", AddClassroom)          // 添加教室
		api.PUT("/classrooms/:id", UpdateClassroom)    // 更新教室信息
		api.DELETE("/classrooms/:id", DeleteClassroom) // 删除教室

		// 课程相关路由
		api.GET("/courses", GetCourses)          // 获取所有课程
		api.POST("/courses", AddCourse)          // 添加课程
		api.PUT("/courses/:id", UpdateCourse)    // 更新课程信息
		api.DELETE("/courses/:id", DeleteCourse) // 删除课程

		// 班级相关路由
		api.GET("/classes", GetClasses)         // 获取所有班级
		api.POST("/classes", AddClass)          // 添加班级
		api.PUT("/classes/:id", UpdateClass)    // 更新班级信息
		api.DELETE("/classes/:id", DeleteClass) // 删除班级

		// 排课相关路由
		api.GET("/schedules", GetSchedules)          // 获取所有排课
		api.POST("/schedules", AddSchedule)          // 添加排课
		api.PUT("/schedules/:id", UpdateSchedule)    // 更新排课信息
		api.DELETE("/schedules/:id", DeleteSchedule) // 删除排课
	}
}

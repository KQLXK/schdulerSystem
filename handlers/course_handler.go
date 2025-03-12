package handlers

import (
	"github.com/gin-gonic/gin"
	"schedule/commen/result"
	"schedule/dto"
	"schedule/services/course"
)

//// GetCourses 获取所有课程
//func GetCourses(c *gin.Context) {
//	courses, err := services.GetAllCourses()
//	if err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//		return
//	}
//	c.JSON(http.StatusOK, courses)
//}
//
//// GetCourseByID 根据ID获取课程
//func GetCourseByID(c *gin.Context) {
//	id := c.Param("id")
//	course, err := services.GetCourseByID(id)
//	if err != nil {
//		c.JSON(http.StatusNotFound, gin.H{"error": "Course not found"})
//		return
//	}
//	c.JSON(http.StatusOK, course)
//}

// AddCourse 添加课程
func AddCourse(c *gin.Context) {
	var req dto.CourseCreateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		result.Errors(c, err)
		return
	}
	resp, err := course.NewCourseCreateFlow(req).Do()
	if err != nil {
		if err == course.DataExistErr {
			result.Error(c, result.CourseExsitsStatus)
			return
		}
		if err == course.InvalidDataErr {
			result.Error(c, result.CourseDataInvalidStatus)
			return
		}
		result.Errors(c, err)
		return
	}
	result.Sucess(c, resp)
	return
}

func UpdateCourse(c *gin.Context) {
	var req dto.CourseUpdateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		result.Errors(c, err)
		return
	}
	if err := course.NewCourseUpdateFlow(req).Do(); err != nil {
		if err == course.DataNotFoundErr {
			result.Error(c, result.CourseNotFoundStatus)
			return
		}
		if err == course.InvalidDataErr {
			result.Error(c, result.CourseDataInvalidStatus)
			return
		}
		result.Errors(c, err)
		return
	}
	result.Sucess(c, nil) // 更新成功，返回空数据
}

// DeleteCourse 处理删除课程的请求
func DeleteCourse(c *gin.Context) {
	courseID := c.Param("course_id") // 假设课程 ID 是通过 URL 参数传递的
	if courseID == "" {
		result.Error(c, result.CourseIDEmptyStatus)
		return
	}
	if err := course.NewCourseDeleteFlow(courseID).Do(); err != nil {
		if err == course.DataNotFoundErr {
			result.Error(c, result.CourseNotFoundStatus)
			return
		}
		result.Errors(c, err)
		return
	}
	result.Sucess(c, nil) // 删除成功，返回空数据
}

// GetCourse 处理获取单个课程的请求
func GetCourse(c *gin.Context) {
	courseID := c.Param("course_id")
	if courseID == "" {
		result.Error(c, result.CourseIDEmptyStatus)
		return
	}
	data, err := course.NewCourseGetFlow(courseID).Do()
	if err != nil {
		if err == course.DataNotFoundErr {
			result.Error(c, result.CourseNotFoundStatus)
			return
		}
		result.Errors(c, err)
		return
	}
	result.Sucess(c, data)
}

// GetAllCourses 处理获取所有课程的请求
func GetAllCourses(c *gin.Context) {
	courses, err := course.NewCourseGetAllFlow().Do()
	if err != nil {
		result.Errors(c, err)
		return
	}
	result.Sucess(c, courses)
}

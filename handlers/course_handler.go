package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"schedule/models"
	"schedule/services"
)

// GetCourses 获取所有课程
func GetCourses(c *gin.Context) {
	courses, err := services.GetAllCourses()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, courses)
}

// GetCourseByID 根据ID获取课程
func GetCourseByID(c *gin.Context) {
	id := c.Param("id")
	course, err := services.GetCourseByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Course not found"})
		return
	}
	c.JSON(http.StatusOK, course)
}

// AddCourse 添加课程
func AddCourse(c *gin.Context) {
	var course models.Course
	if err := c.ShouldBindJSON(&course); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := services.CreateCourse(&course); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, course)
}

// UpdateCourse 更新课程信息
func UpdateCourse(c *gin.Context) {
	id := c.Param("id")
	var course models.Course
	if err := c.ShouldBindJSON(&course); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := services.UpdateCourse(id, &course); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, course)
}

// DeleteCourse 删除课程
func DeleteCourse(c *gin.Context) {
	id := c.Param("id")
	if err := services.DeleteCourse(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Course deleted successfully"})
}

package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"schedule/services"
)

// GetClassroomUtilization 获取教室利用率
func GetClassroomUtilization(c *gin.Context) {
	utilization, err := services.GetClassroomUtilization()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, utilization)
}

// GetTeacherWorkload 获取教师工作量
func GetTeacherWorkload(c *gin.Context) {
	workload, err := services.GetTeacherWorkload()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, workload)
}

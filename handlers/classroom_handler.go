package handlers

import (
	"github.com/gin-gonic/gin"
	"schedule/commen/result"
	"schedule/dto"
	"schedule/services/classroom"
)

// GetClassrooms 获取所有教室
func GetClassrooms(c *gin.Context) {
	classrooms, err := classroom.GetAllClassrooms()
	if err != nil {
		result.Errors(c, err)
		return
	}
	result.Sucess(c, classrooms)
}

// GetClassroomByID 根据ID获取教室
func GetClassroomByID(c *gin.Context) {
	id := c.Param("id")
	resp, err := classroom.GetClassroomByID(id)
	if err != nil {
		result.Error(c, result.ClassroomNotFoundStatus)
		return
	}
	result.Sucess(c, resp)
}

// AddClassroom 添加教室
func AddClassroom(c *gin.Context) {
	var req dto.ClassroomCreateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		result.Errors(c, err)
		return
	}

	resp, err := classroom.CreateClassroom(req)
	if err != nil {
		if err == classroom.ExistsError {
			result.Error(c, result.ClassroomExistsStatus)
			return
		}
		result.Errors(c, err)
		return
	}
	result.Sucess(c, resp)
}

// UpdateClassroom 更新教室信息
func UpdateClassroom(c *gin.Context) {
	id := c.Param("id")
	var req dto.ClassroomUpdateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		result.Errors(c, err)
		return
	}

	resp, err := classroom.UpdateClassroom(id, req)
	if err != nil {
		if err == classroom.NotFoundError {
			result.Error(c, result.ClassroomNotFoundStatus)
			return
		}
		result.Errors(c, err)
		return
	}
	result.Sucess(c, resp)
}

// DeleteClassroom 删除教室
func DeleteClassroom(c *gin.Context) {
	id := c.Param("id")
	if err := classroom.DeleteClassroom(id); err != nil {
		if err == classroom.NotFoundError {
			result.Error(c, result.ClassroomNotFoundStatus)
			return
		}
		result.Errors(c, err)
		return
	}
	result.Sucess(c, nil)
}

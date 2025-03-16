package handlers

import (
	"github.com/gin-gonic/gin"
	"schedule/commen/result"
	"schedule/dto"
	"schedule/services/teacher"
)

// GetTeachers 获取所有教师
func GetTeachers(c *gin.Context) {
	teachers, err := teacher.GetAllTeachers()
	if err != nil {
		result.Errors(c, err)
		return
	}
	result.Sucess(c, teachers)
}

// GetTeacherByID 根据ID获取教师
func GetTeacherByID(c *gin.Context) {
	id := c.Param("id")
	resp, err := teacher.GetTeacherByID(id)
	if err != nil {
		result.Error(c, result.TeacherNotFoundStatus)
		return
	}
	result.Sucess(c, resp)
}

// AddTeacher 添加教师
func AddTeacher(c *gin.Context) {
	var req dto.TeacherCreateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		result.Errors(c, err)
		return
	}

	resp, err := teacher.CreateTeacher(req)
	if err != nil {
		if err == teacher.ExistsError {
			result.Error(c, result.TeacherExistsStatus)
			return
		}
		result.Errors(c, err)
		return
	}
	result.Sucess(c, resp)
}

// UpdateTeacher 更新教师信息
func UpdateTeacher(c *gin.Context) {
	id := c.Param("id")
	var req dto.TeacherUpdateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		result.Errors(c, err)
		return
	}

	resp, err := teacher.UpdateTeacher(id, req)
	if err != nil {
		if err == teacher.NotFoundError {
			result.Error(c, result.TeacherNotFoundStatus)
			return
		}
		result.Errors(c, err)
		return
	}
	result.Sucess(c, resp)
}

// DeleteTeacher 删除教师
func DeleteTeacher(c *gin.Context) {
	id := c.Param("id")
	if err := teacher.DeleteTeacher(id); err != nil {
		if err == teacher.NotFoundError {
			result.Error(c, result.TeacherNotFoundStatus)
			return
		}
		result.Errors(c, err)
		return
	}
	result.Sucess(c, nil)
}

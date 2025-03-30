package handlers

import (
	"github.com/gin-gonic/gin"
	"log"
	"os"
	"schedule/commen/result"
	"schedule/dto"
	"schedule/services/teacher"
	"strconv"
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

func QueryTeacherByPage(c *gin.Context) {
	// 获取查询参数并转换为 int
	pageStr := c.Query("page")
	page, _ := strconv.Atoi(pageStr)

	pageSizeStr := c.Query("pagesize")
	pageSize, _ := strconv.Atoi(pageSizeStr)

	resp, err := teacher.TeacherQueryByPage(page, pageSize)
	if err != nil {
		if err == teacher.PageNumErr {
			result.Error(c, result.PageDataErrStatus)
			return
		}
		result.Errors(c, err)
		return
	}
	result.Sucess(c, resp)
}

func SearchTeacher(c *gin.Context) {
	searchStr := c.PostForm("search_str")
	resp, err := teacher.TeacherSearch(searchStr)
	if err != nil {
		result.Errors(c, err)
		return
	}
	result.Sucess(c, resp)
}

func AddTeacherByExcel(c *gin.Context) {
	file, err := c.FormFile("teacher_file")
	if err != nil {
		result.Error(c, result.FileNotReceiveStatus)
		return
	}

	tempFilePath := "./tmp/" + file.Filename
	if err = c.SaveUploadedFile(file, tempFilePath); err != nil {
		log.Println("save uploaded file failed, err:", err)
		result.Error(c, result.ServerInteralErrStatus)
		return
	}

	resp, err := teacher.TeacherAddByExcel(file.Filename)
	if err != nil {
		result.Error(c, result.FileFormatErrStatus)
		return
	}

	defer os.Remove(tempFilePath)
	result.Sucess(c, resp)
}

package handlers

import (
	"github.com/gin-gonic/gin"
	"log"
	"os"
	"schedule/commen/result"
	"schedule/dto"
	"schedule/services/classroom"
	"strconv"
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

func QueryClassroomByPage(c *gin.Context) {
	// 获取查询参数并转换为 int
	pageStr := c.Query("page")
	page, _ := strconv.Atoi(pageStr)

	pageSizeStr := c.Query("pagesize")
	pageSize, _ := strconv.Atoi(pageSizeStr)

	resp, err := classroom.ClassroomQueryByPage(page, pageSize)
	if err != nil {
		if err == classroom.PageNumErr {
			result.Error(c, result.PageDataErrStatus)
			return
		}
		result.Errors(c, err)
		return
	}
	result.Sucess(c, resp)
}

func SearchClassroom(c *gin.Context) {
	searchStr := c.Query("search_str")
	resp, err := classroom.ClassroomSearch(searchStr)
	if err != nil {
		result.Errors(c, err)
		return
	}
	result.Sucess(c, resp)
}

func AddClassroomByExcel(c *gin.Context) {
	file, err := c.FormFile("classroom_file")
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

	resp, err := classroom.ClassroomAddByExcel(file.Filename)
	if err != nil {
		result.Error(c, result.FileFormatErrStatus)
		return
	}

	defer os.Remove(tempFilePath)
	result.Sucess(c, resp)
}

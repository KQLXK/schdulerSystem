package handlers

import (
	"github.com/gin-gonic/gin"
	"log"
	"os"
	"schedule/commen/result"
	"schedule/dto"
	"schedule/services/class"
	"strconv"
)

func AddClass(c *gin.Context) {
	var req dto.ClassCreateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		result.Errors(c, err)
		return
	}
	log.Println("get req sucess, req:", req)
	resp, err := class.CreateClass(req)
	if err != nil {
		if err == class.ExistsError {
			result.Error(c, result.ClassExistsStatus)
			return
		}
		result.Errors(c, err)
		return
	}
	result.Sucess(c, resp)
}

func AddClassByExcel(c *gin.Context) {
	file, err := c.FormFile("class_file")
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
	resp, err := class.ClassAddByExcel(file.Filename)
	if err != nil {
		result.Error(c, result.FileFormatErrStatus)
		return
	}
	defer os.Remove(tempFilePath)
	result.Sucess(c, resp)
}

func UpdateClass(c *gin.Context) {
	id := c.Param("id")
	var req dto.ClassUpdateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		result.Errors(c, err)
		return
	}
	req.ID = id
	resp, err := class.UpdateClass(id, req)
	if err != nil {
		if err == class.NotFoundError {
			result.Error(c, result.ClassNotFoundStatus)
			return
		}
		result.Errors(c, err)
		return
	}
	result.Sucess(c, resp)
}

func DeleteClass(c *gin.Context) {
	id := c.Param("id")
	if err := class.DeleteClass(id); err != nil {
		if err == class.NotFoundError {
			result.Error(c, result.ClassNotFoundStatus)
			return
		}
		result.Errors(c, err)
		return
	}
	result.Sucess(c, nil)
}

func GetClassByID(c *gin.Context) {
	id := c.Param("id")
	resp, err := class.GetClassByID(id)
	if err != nil {
		if err == class.NotFoundError {
			result.Error(c, result.ClassNotFoundStatus)
			return
		}
		result.Errors(c, err)
		return
	}
	result.Sucess(c, resp)
}

func GetAllClasses(c *gin.Context) {
	resp, err := class.GetAllClasses()
	if err != nil {
		result.Errors(c, err)
		return
	}
	result.Sucess(c, resp)
}

func QueryClassByPage(c *gin.Context) {
	pageStr := c.Query("page")
	page, _ := strconv.Atoi(pageStr)

	pageSizeStr := c.Query("pagesize")
	pageSize, _ := strconv.Atoi(pageSizeStr)

	resp, err := class.ClassQueryByPage(page, pageSize)
	if err != nil {
		if err == class.PageNumErr {
			result.Error(c, result.PageDataErrStatus)
			return
		}
		result.Errors(c, err)
		return
	}
	result.Sucess(c, resp)
}

func SearchClass(c *gin.Context) {
	searchStr := c.Query("search_str")
	resp, err := class.ClassSearch(searchStr)
	if err != nil {
		result.Errors(c, err)
		return
	}
	result.Sucess(c, resp)
}

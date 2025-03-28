package handlers

import (
	"github.com/gin-gonic/gin"
	"log"
	"os"
	"schedule/commen/result"
	"schedule/dto"
	"schedule/services/course"
	"strconv"
)

// AddCourse 添加课程
func AddCourse(c *gin.Context) {
	var req dto.CourseCreateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		result.Errors(c, err)
		return
	}
	log.Println("get req sucess, req:", req)
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
		if err == course.DataNullErr {
			result.Error(c, result.CourseDataEmptyStatus)
			return
		}
		result.Errors(c, err)
		return
	}
	result.Sucess(c, resp)
	return
}

// 使用excel导入课表
func AddCourseByExcel(c *gin.Context) {
	file, err := c.FormFile("course_file")
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
	resp, err := course.CourseAddByExcel(file.Filename)
	if err != nil {
		result.Error(c, result.FileFormatErrStatus)
		return
	}
	defer os.Remove(tempFilePath)
	result.Sucess(c, resp)
}

func UpdateCourse(c *gin.Context) {
	var req dto.CourseUpdateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		result.Errors(c, err)
		return
	}
	log.Println("get req sucess, req:", req)
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
	resp, err := course.CourseQueryAll()
	if err != nil {
		result.Errors(c, err)
	}
	result.Sucess(c, resp)
}

func QueryCourseByPage(c *gin.Context) {
	// 获取查询参数并转换为 int64
	pageStr := c.Query("page")
	Page, _ := strconv.ParseInt(pageStr, 10, 64)

	pageSizeStr := c.Query("pagesize")
	PageSize, _ := strconv.ParseInt(pageSizeStr, 10, 64)
	log.Println("page:", Page, "pagesize:", PageSize)
	resp, err := course.CourseQueryByPage(int(Page), int(PageSize))
	if err != nil {
		if err == course.PageNumErr {
			result.Error(c, result.PageDataErrStatus)
			return
		}
		result.Errors(c, err)
		return
	}
	result.Sucess(c, resp)
}

func SearchCourse(c *gin.Context) {
	searchStr := c.PostForm("search_str")
	resp, err := course.CourseSearch(searchStr)
	if err != nil {
		result.Errors(c, err)
		return
	}
	result.Sucess(c, resp)
}

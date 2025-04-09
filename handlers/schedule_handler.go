package handlers

import (
	"github.com/gin-gonic/gin"
	"log"
	"os"
	"schedule/commen/result"
	"schedule/dto"
	"schedule/services/schedule"
	"strconv"
)

func CreateSchedule(c *gin.Context) {
	var req dto.ScheduleCreateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println("bind to req failed, err:", err)
		result.Errors(c, err)
		return
	}
	resp, err := schedule.CreateSchedule(req)
	if err != nil {
		if err == schedule.ErrInvalidHours {
			result.Error(c, result.InvalidHoursStatus)
			return
		} else if err == schedule.ErrInvalidSemester {
			result.Error(c, result.InvalidSemesterStatus)
			return
		} else if err == schedule.ErrInvalidTimeFormat {
			result.Error(c, result.InvalidTimeFormatStatus)
			return
		} else {
			result.Errors(c, err)
		}
	}
	result.Sucess(c, resp)
}

// ImportScheduleByExcel Excel导入排课
func AddScheduleByExcel(c *gin.Context) {
	file, err := c.FormFile("schedule_file")
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
	defer os.Remove(tempFilePath)

	resp, err := schedule.AddByExcel(file.Filename)
	if err != nil {
		result.Errors(c, err)
		return
	}
	result.Sucess(c, resp)
}

// UpdateSchedule 更新排课
func UpdateSchedule(c *gin.Context) {
	var req dto.ScheduleUpdateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		result.Errors(c, err)
		return
	}
	log.Println("update schedule req:", req)

	if err := schedule.NewScheduleUpdateFlow(req).Do(); err != nil {
		switch err {
		case schedule.ErrScheduleNotFound:
			result.Error(c, result.ScheduleNotFoundStatus)
		case schedule.ErrInvalidHours:
			result.Error(c, result.ScheduleHoursInvalidStatus)
		//case schedule.ErrScheduleConflict:
		//	result.Error(c, result.ScheduleConflictStatus)
		default:
			result.Errors(c, err)
		}
		return
	}
	result.Sucess(c, nil)
}

// DeleteSchedule 删除排课
func DeleteSchedule(c *gin.Context) {
	scheduleID := c.Param("schedule_id")

	if err := schedule.DeleteSchedule(scheduleID); err != nil {
		if err == schedule.ErrScheduleNotFound {
			result.Error(c, result.ScheduleNotFoundStatus)
			return
		}
		result.Errors(c, err)
		return
	}
	result.Sucess(c, nil)
}

// GetSchedule 获取单个排课
func GetSchedule(c *gin.Context) {
	scheduleID, err := strconv.Atoi(c.Param("schedule_id"))
	if err != nil || scheduleID <= 0 {
		result.Error(c, result.ScheduleIDInvalidStatus)
		return
	}

	data, err := schedule.QuerySchedule(scheduleID)
	if err != nil {
		if err == schedule.ErrScheduleNotFound {
			result.Error(c, result.ScheduleNotFoundStatus)
			return
		}
		result.Errors(c, err)
		return
	}
	result.Sucess(c, data)
}

// GetAllSchedules 获取所有排课
func GetAllSchedules(c *gin.Context) {
	resp, err := schedule.QueryAllSchedules()
	if err != nil {
		result.Errors(c, err)
		return
	}
	result.Sucess(c, resp)
}

// QuerySchedulesByPage 分页查询排课
func QuerySchedulesByPage(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	resp, err := schedule.QuerySchedulesByPage(page, pageSize)
	if err != nil {
		if err == schedule.ErrPageOutOfRange {
			result.Error(c, result.PageDataErrStatus)
			return
		}
		result.Errors(c, err)
		return
	}
	result.Sucess(c, resp)
}

//// SearchSchedules 搜索排课
//func SearchSchedules(c *gin.Context) {
//	searchStr := c.Query("search_str")
//	filter := schedule.QueryParams{
//		Semester:  c.Query("semester"),
//		TeacherID: c.Query("teacher_id"),
//		CourseID:  c.Query("course_id"),
//		Campus:    c.Query("campus"),
//	}
//
//	resp, err := schedule.SearchSchedules(searchStr, filter)
//	if err != nil {
//		result.Errors(c, err)
//		return
//	}
//	result.Success(c, resp)
//}

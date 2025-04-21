package handlers

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"schedule/commen/result"
	"schedule/dto"
	"schedule/models"
	"schedule/services/class"
	"schedule/services/classroom"
	"schedule/services/course"
	"schedule/services/schedule"
	"schedule/services/table"
	"schedule/services/teacher"
	"strings"
)

func AddDataByExcel(c *gin.Context) {
	file, err := c.FormFile("any_file")
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
	switch switchfilename(file.Filename) {
	case "schedule":
		resp, err := schedule.AddByExcel(file.Filename)
		if err != nil {
			result.Errors(c, err)
			return
		}
		result.Sucess(c, resp)
	case "class":
		resp, err := class.ClassAddByExcel(file.Filename)
		if err != nil {
			result.Errors(c, err)
			return
		}
		result.Sucess(c, resp)
	case "teacher":
		resp, err := teacher.TeacherAddByExcel(file.Filename)
		if err != nil {
			result.Errors(c, err)
			return
		}
		result.Sucess(c, resp)
	case "course":
		resp, err := course.CourseAddByExcel(file.Filename)
		if err != nil {
			result.Errors(c, err)
			return
		}
		result.Sucess(c, resp)
	case "classroom":
		resp, err := classroom.ClassroomAddByExcel(file.Filename)
		if err != nil {
			result.Errors(c, err)
			return
		}
		result.Sucess(c, resp)
	default:
		result.Error(c, result.FileNameErrStatus)
	}
}

// GetClassScheduleHandler 处理查询班级课表的请求
func GetClassTableHandler(c *gin.Context) {
	classID := c.Query("class_id")
	semester := c.Query("semester")

	if classID == "" || semester == "" {
		result.Error(c, result.GetQueryStringfailed)
		return
	}

	scheduleResults, err := table.GetClassScheduleBySemester(classID, semester)
	if err != nil {
		result.Errors(c, err)
		return
	}

	var resp dto.GetClassTableResp
	resp.ClassTables = ConvertScheduleResultsToClassTables(scheduleResults)

	c.JSON(http.StatusOK, resp)
}

func GetTeacherTableHandler(c *gin.Context) {
	teacherID := c.Query("teacher_id")
	semester := c.Query("semester")

	if teacherID == "" || semester == "" {
		result.Error(c, result.GetQueryStringfailed)
		return
	}
	log.Println("/table/teacher: get req success, teacher_id:", teacherID, "semester:", semester)

	scheduleResults, err := table.GetTeacherScheduleBySemester(teacherID, semester)
	if err != nil {
		result.Errors(c, err)
		return
	}

	var resp dto.GetClassTableResp
	resp.ClassTables = ConvertScheduleResultsToClassTables(scheduleResults)

	c.JSON(http.StatusOK, resp)
}

func GetClassroomTableHandler(c *gin.Context) {
	ClassroomID := c.Query("classroom_id")
	semester := c.Query("semester")

	if ClassroomID == "" || semester == "" {
		result.Error(c, result.GetQueryStringfailed)
		return
	}
	log.Println("/table/classroom: get req success, classroom_id:", ClassroomID, "semester:", semester)

	scheduleResults, err := table.GetClassroomScheduleBySemester(ClassroomID, semester)
	if err != nil {
		result.Errors(c, err)
		return
	}

	var resp dto.GetClassTableResp
	resp.ClassTables = ConvertScheduleResultsToClassTables(scheduleResults)

	c.JSON(http.StatusOK, resp)
}

func ConvertScheduleResultsToClassTables(scheduleResults []models.ScheduleResult) []dto.ClassTable {
	var classTables []dto.ClassTable

	for _, result := range scheduleResults {
		classTable := dto.ClassTable{
			ID:          int(result.ID),
			Semester:    result.Semester, // 假设 ScheduleResult 中有 Semester 字段
			ScheduleID:  int64(result.ScheduleID),
			CourseID:    result.CourseID,
			CourseName:  result.CourseName,
			TeacherID:   result.TeacherID,
			TeacherName: result.TeacherName,
			ClassroomID: result.ClassroomID,
			ClassIDs:    result.ClassIDs,
			ClassNames:  result.ClassNames,
			Timeslots:   result.TimeSlots, // 假设 TimeSlots 是 JSON 类型
		}
		classTables = append(classTables, classTable)
	}

	return classTables
}

func switchfilename(filename string) string {
	if strings.Contains(filename, "排课任务") {
		return "schedule"
	} else if strings.Contains(filename, "课程") {
		return "course"
	} else if strings.Contains(filename, "教室") {
		return "classroom"
	} else if strings.Contains(filename, "教师") {
		return "teacher"
	} else if strings.Contains(filename, "班级") {
		return "class"
	} else {
		return ""
	}
}

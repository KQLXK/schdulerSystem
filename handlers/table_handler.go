package handlers

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"schedule/commen/result"
	"schedule/dto"
	"schedule/models"
	"schedule/services/table"
)

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

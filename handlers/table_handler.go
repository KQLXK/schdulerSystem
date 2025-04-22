package handlers

import (
	"github.com/gin-gonic/gin"
	"log"
	"mime/multipart"
	"net/http"
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
	// 获取所有文件
	courseFile, err := c.FormFile("course_file")
	if err != nil {
		result.Error(c, result.FileNotReceiveStatus)
		return
	}

	scheduleFile, err := c.FormFile("schedule_file")
	if err != nil {
		result.Error(c, result.FileNotReceiveStatus)
		return
	}

	classroomFile, err := c.FormFile("classroom_file")
	if err != nil {
		result.Error(c, result.FileNotReceiveStatus)
		return
	}

	teacherFile, err := c.FormFile("teacher_file")
	if err != nil {
		result.Error(c, result.FileNotReceiveStatus)
		return
	}

	classFile, err := c.FormFile("class_file")
	if err != nil {
		result.Error(c, result.FileNotReceiveStatus)
		return
	}

	// 创建临时目录
	//tmpDir := "./tmp/"
	//if err := os.MkdirAll(tmpDir, 0755); err != nil {
	//	log.Printf("创建临时目录失败: %v", err)
	//	result.Error(c, result.ServerInteralErrStatus)
	//	return
	//}
	//defer os.RemoveAll(tmpDir)

	// 保存文件到临时目录的辅助函数
	saveFile := func(file *multipart.FileHeader) error {
		filePath := "./tmp/" + file.Filename
		if err := c.SaveUploadedFile(file, filePath); err != nil {
			return err
		}
		return nil
	}

	// 保存所有文件
	if err := saveFile(teacherFile); err != nil {
		log.Printf("保存教师文件失败: %v", err)
		result.Error(c, result.ServerInteralErrStatus)
		return
	}

	if err := saveFile(classFile); err != nil {
		log.Printf("保存班级文件失败: %v", err)
		result.Error(c, result.ServerInteralErrStatus)
		return
	}

	if err := saveFile(classroomFile); err != nil {
		log.Printf("保存教室文件失败: %v", err)
		result.Error(c, result.ServerInteralErrStatus)
		return
	}

	if err := saveFile(courseFile); err != nil {
		log.Printf("保存课程文件失败: %v", err)
		result.Error(c, result.ServerInteralErrStatus)
		return
	}

	if err := saveFile(scheduleFile); err != nil {
		log.Printf("保存课表文件失败: %v", err)
		result.Error(c, result.ServerInteralErrStatus)
		return
	}

	// 按照依赖顺序处理文件
	var response dto.AddAllByExcelResp

	// 处理教师信息
	teacherResp, err := teacher.TeacherAddByExcel(teacherFile.Filename)
	if err != nil {
		result.Errors(c, err)
		return
	}
	response.TeacherResp = teacherResp

	// 处理班级信息
	classResp, err := class.ClassAddByExcel(classFile.Filename)
	if err != nil {
		result.Errors(c, err)
		return
	}
	response.ClassResp = classResp

	// 处理教室信息
	classroomResp, err := classroom.ClassroomAddByExcel(classroomFile.Filename)
	if err != nil {
		result.Errors(c, err)
		return
	}
	response.ClassroomResp = classroomResp

	// 处理课程信息
	courseResp, err := course.CourseAddByExcel(courseFile.Filename)
	if err != nil {
		result.Errors(c, err)
		return
	}
	response.CourseResp = courseResp

	// 处理课表信息
	scheduleResp, err := schedule.AddByExcel(scheduleFile.Filename)
	if err != nil {
		result.Errors(c, err)
		return
	}
	response.ScheduleResp = scheduleResp

	result.Sucess(c, response)
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

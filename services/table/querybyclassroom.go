package table

import (
	"schedule/models"
)

// GetClassScheduleBySemester 查询班级课表
func GetClassroomScheduleBySemester(classroomID string, semester string) ([]models.ScheduleResult, error) {
	// 假设 ScheduleResult 结构体中有一个字段存储学期信息
	scheduleResults, err := models.GetAllScheduleResults()
	if err != nil {
		return nil, err
	}

	// 过滤出符合条件的排课结果
	Results := make([]models.ScheduleResult, 0)
	for _, result := range scheduleResults {
		if result.ClassroomID == classroomID && result.Schedule.Semester == semester {
			Results = append(Results, result)
		}
	}

	return Results, nil
}

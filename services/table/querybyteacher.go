package table

import (
	"schedule/models"
)

// GetTeacherScheduleBySemester 查询教师课表
func GetTeacherScheduleBySemester(teacherID string, semester string) ([]models.ScheduleResult, error) {
	// 假设 ScheduleResult 结构体中有一个字段存储学期信息
	scheduleResults, err := models.GetAllScheduleResults()
	if err != nil {
		return nil, err
	}

	var results []models.ScheduleResult

	// 过滤出符合条件的排课结果
	for _, result := range scheduleResults {
		if result.TeacherID == teacherID && result.Schedule.Semester == semester {
			results = append(results, result)
		}
	}

	return results, nil
}

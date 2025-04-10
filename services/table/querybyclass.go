package table

import (
	"schedule/models"
)

// GetClassScheduleBySemester 查询班级课表
func GetClassScheduleBySemester(classID string, semester string) ([]models.ScheduleResult, error) {
	// 假设 ScheduleResult 结构体中有一个字段存储学期信息
	scheduleResults, err := models.GetAllScheduleResults()
	if err != nil {
		return nil, err
	}

	// 过滤出符合条件的排课结果
	for _, result := range scheduleResults {
		if IsClassInClassList(classID, result.ClassIDs) && result.Schedule.Semester == semester {
			scheduleResults = append(scheduleResults, result)
		}
	}

	return scheduleResults, nil
}

// 判断class是否在[]class中的函数
func IsClassInClassList(class string, classList []string) bool {
	for _, c := range classList {
		if c == class {
			return true
		}
	}
	return false
}

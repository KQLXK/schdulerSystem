package services

import (
	"schedule/database"
)

// GetClassroomUtilization 获取教室利用率
func GetClassroomUtilization() ([]map[string]interface{}, error) {
	var utilization []map[string]interface{}
	if err := database.DB.Table("schedules").
		Select("classroom_id, COUNT(*) * 100.0 / (SELECT COUNT(*) FROM schedules) AS utilization").
		Group("classroom_id").
		Find(&utilization).Error; err != nil {
		return nil, err
	}
	return utilization, nil
}

// GetTeacherWorkload 获取教师工作量
func GetTeacherWorkload() ([]map[string]interface{}, error) {
	var workload []map[string]interface{}
	if err := database.DB.Table("schedules").
		Select("teacher_id, COUNT(*) AS workload").
		Group("teacher_id").
		Find(&workload).Error; err != nil {
		return nil, err
	}
	return workload, nil
}

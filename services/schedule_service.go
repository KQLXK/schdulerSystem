package services

import (
	"schedule/database"
	"schedule/models"
)

// GetAllSchedules 获取所有排课结果
func GetAllSchedules() ([]models.Schedule, error) {
	var schedules []models.Schedule
	if err := database.DB.Find(&schedules).Error; err != nil {
		return nil, err
	}
	return schedules, nil
}

// GetScheduleByID 根据ID获取排课结果
func GetScheduleByID(id string) (*models.Schedule, error) {
	var schedule models.Schedule
	if err := database.DB.Where("id = ?", id).First(&schedule).Error; err != nil {
		return nil, err
	}
	return &schedule, nil
}

// CreateSchedule 创建排课结果
func CreateSchedule(schedule *models.Schedule) error {
	if err := database.DB.Create(schedule).Error; err != nil {
		return err
	}
	return nil
}

// UpdateSchedule 更新排课结果
func UpdateSchedule(id string, schedule *models.Schedule) error {
	if err := database.DB.Model(&models.Schedule{}).Where("id = ?", id).Updates(schedule).Error; err != nil {
		return err
	}
	return nil
}

// DeleteSchedule 删除排课结果
func DeleteSchedule(id string) error {
	if err := database.DB.Where("id = ?", id).Delete(&models.Schedule{}).Error; err != nil {
		return err
	}
	return nil
}

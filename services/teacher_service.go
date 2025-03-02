package services

import (
	"schedule/database"
	"schedule/models"
)

// GetAllTeachers 获取所有教师
func GetAllTeachers() ([]models.Teacher, error) {
	var teachers []models.Teacher
	if err := database.DB.Find(&teachers).Error; err != nil {
		return nil, err
	}
	return teachers, nil
}

// GetTeacherByID 根据ID获取教师
func GetTeacherByID(id string) (*models.Teacher, error) {
	var teacher models.Teacher
	if err := database.DB.Where("id = ?", id).First(&teacher).Error; err != nil {
		return nil, err
	}
	return &teacher, nil
}

// CreateTeacher 创建教师
func CreateTeacher(teacher *models.Teacher) error {
	if err := database.DB.Create(teacher).Error; err != nil {
		return err
	}
	return nil
}

// UpdateTeacher 更新教师信息
func UpdateTeacher(id string, teacher *models.Teacher) error {
	if err := database.DB.Model(&models.Teacher{}).Where("id = ?", id).Updates(teacher).Error; err != nil {
		return err
	}
	return nil
}

// DeleteTeacher 删除教师
func DeleteTeacher(id string) error {
	if err := database.DB.Where("id = ?", id).Delete(&models.Teacher{}).Error; err != nil {
		return err
	}
	return nil
}

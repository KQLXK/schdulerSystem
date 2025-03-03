package services

import (
	"schedule/database"
	"schedule/models"
)

// GetAllCourses 获取所有课程
func GetAllCourses() ([]models.Course, error) {
	var courses []models.Course
	if err := database.DB.Find(&courses).Error; err != nil {
		return nil, err
	}
	return courses, nil
}

// GetCourseByID 根据ID获取课程
func GetCourseByID(id string) (*models.Course, error) {
	var course models.Course
	if err := database.DB.Where("id = ?", id).First(&course).Error; err != nil {
		return nil, err
	}
	return &course, nil
}

// CreateCourse 创建课程
func CreateCourse(course *models.Course) error {
	if err := database.DB.Create(course).Error; err != nil {
		return err
	}
	return nil
}

// UpdateCourse 更新课程信息
func UpdateCourse(id string, course *models.Course) error {
	if err := database.DB.Model(&models.Course{}).Where("id = ?", id).Updates(course).Error; err != nil {
		return err
	}
	return nil
}

// DeleteCourse 删除课程
func DeleteCourse(id string) error {
	if err := database.DB.Where("id = ?", id).Delete(&models.Course{}).Error; err != nil {
		return err
	}
	return nil
}

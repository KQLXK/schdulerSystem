package services

import (
	"schedule/database"
	"schedule/models"
)

// GetAllClassrooms 获取所有教室
func GetAllClassrooms() ([]models.Classroom, error) {
	var classrooms []models.Classroom
	if err := database.DB.Find(&classrooms).Error; err != nil {
		return nil, err
	}
	return classrooms, nil
}

// GetClassroomByID 根据ID获取教室
func GetClassroomByID(id string) (*models.Classroom, error) {
	var classroom models.Classroom
	if err := database.DB.Where("id = ?", id).First(&classroom).Error; err != nil {
		return nil, err
	}
	return &classroom, nil
}

// CreateClassroom 创建教室
func CreateClassroom(classroom *models.Classroom) error {
	if err := database.DB.Create(classroom).Error; err != nil {
		return err
	}
	return nil
}

// UpdateClassroom 更新教室信息
func UpdateClassroom(id string, classroom *models.Classroom) error {
	if err := database.DB.Model(&models.Classroom{}).Where("id = ?", id).Updates(classroom).Error; err != nil {
		return err
	}
	return nil
}

// DeleteClassroom 删除教室
func DeleteClassroom(id string) error {
	if err := database.DB.Where("id = ?", id).Delete(&models.Classroom{}).Error; err != nil {
		return err
	}
	return nil
}

package services

import (
	"schedule/database"
	"schedule/models"
)

// GetAllClasses 获取所有班级
func GetAllClasses() ([]models.Class, error) {
	var classes []models.Class
	if err := database.DB.Find(&classes).Error; err != nil {
		return nil, err
	}
	return classes, nil
}

// GetClassByID 根据ID获取班级
func GetClassByID(id string) (*models.Class, error) {
	var class models.Class
	if err := database.DB.Where("id = ?", id).First(&class).Error; err != nil {
		return nil, err
	}
	return &class, nil
}

// CreateClass 创建班级
func CreateClass(class *models.Class) error {
	if err := database.DB.Create(class).Error; err != nil {
		return err
	}
	return nil
}

// UpdateClass 更新班级信息
func UpdateClass(id string, class *models.Class) error {
	if err := database.DB.Model(&models.Class{}).Where("id = ?", id).Updates(class).Error; err != nil {
		return err
	}
	return nil
}

// DeleteClass 删除班级
func DeleteClass(id string) error {
	if err := database.DB.Where("id = ?", id).Delete(&models.Class{}).Error; err != nil {
		return err
	}
	return nil
}

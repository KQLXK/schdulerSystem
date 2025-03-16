package teacher

import (
	"schedule/models"
)

// GetAllTeachers 获取所有教师
func GetAllTeachers() ([]models.Teacher, error) {
	return models.GetAllTeachers()
}

// GetTeacherByID 根据ID获取教师
func GetTeacherByID(id string) (*models.Teacher, error) {
	teacher, err := models.GetTeacherByID(id)
	if err != nil {
		return nil, NotFoundError
	}
	return teacher, nil
}

package classroom

import (
	"errors"
	"schedule/models"
)

var (
	NotFoundError = errors.New("教室未找到")
)

// GetAllClassrooms 获取所有教室
func GetAllClassrooms() ([]models.Classroom, error) {
	return models.NewClassroomDao().GetAllClassrooms()
}

// GetClassroomByID 根据ID获取教室
func GetClassroomByID(id string) (*models.Classroom, error) {
	classroom, err := models.NewClassroomDao().GetClassroomByID(id)
	if err != nil {
		return nil, NotFoundError
	}
	return classroom, nil
}

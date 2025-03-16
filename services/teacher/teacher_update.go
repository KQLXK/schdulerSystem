package teacher

import (
	"errors"
	"schedule/dto"
	"schedule/models"
)

var (
	NotFoundError = errors.New("教师未找到")
)

// UpdateTeacher 更新教师信息
func UpdateTeacher(id string, req dto.TeacherUpdateReq) (*models.Teacher, error) {
	// 检查教师是否存在
	if _, err := models.GetTeacherByID(id); err != nil {
		return nil, NotFoundError
	}

	teacher := models.Teacher{
		ID:         req.ID,
		Name:       req.Name,
		Gender:     req.Gender,
		Department: req.Department,
		IsExternal: req.IsExternal,
		Status:     req.Status,
	}

	if err := models.UpdateTeacher(id, &teacher); err != nil {
		return nil, err
	}
	return &teacher, nil
}

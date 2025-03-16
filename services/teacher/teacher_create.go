package teacher

import (
	"errors"
	"schedule/dto"
	"schedule/models"
)

var (
	ExistsError = errors.New("教师工号或姓名已存在")
)

// CreateTeacher 创建教师
func CreateTeacher(req dto.TeacherCreateReq) (*models.Teacher, error) {
	// 检查教师是否已存在
	if _, err := models.GetTeacherByID(req.ID); err == nil {
		return nil, ExistsError
	}
	// 设置默认状态
	if req.Status == "" {
		req.Status = "启用"
	}

	teacher := models.Teacher{
		ID:         req.ID,
		Name:       req.Name,
		Gender:     req.Gender,
		Department: req.Department,
		IsExternal: req.IsExternal,
		Status:     req.Status,
	}

	if err := models.CreateTeacher(&teacher); err != nil {
		return nil, err
	}
	return &teacher, nil
}

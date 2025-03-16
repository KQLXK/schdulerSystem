package classroom

import (
	"errors"
	"schedule/dto"
	"schedule/models"
)

var (
	ExistsError = errors.New("教室编号或名称已存在")
)

func CreateClassroom(req dto.ClassroomCreateReq) (*models.Classroom, error) {
	// 检查教室是否已存在
	if _, err := models.NewClassroomDao().GetClassroomByID(req.ID); err == nil {
		return nil, ExistsError
	}

	classroom := models.Classroom{
		ID:          req.ID,
		Name:        req.Name,
		Campus:      req.Campus,
		Building:    req.Building,
		Floor:       req.Floor,
		Capacity:    req.Capacity,
		Type:        req.Type,
		HasAC:       req.HasAC,
		Description: req.Description,
		Department:  req.Department,
		Status:      req.Status,
	}

	if err := models.NewClassroomDao().CreateClassroom(&classroom); err != nil {
		return nil, err
	}
	return &classroom, nil
}

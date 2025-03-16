package classroom

import (
	"schedule/dto"
	"schedule/models"
)

// UpdateClassroom 更新教室信息
func UpdateClassroom(id string, req dto.ClassroomUpdateReq) (*models.Classroom, error) {
	// 检查教室是否存在
	if _, err := models.NewClassroomDao().GetClassroomByID(id); err != nil {
		return nil, NotFoundError
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

	if err := models.NewClassroomDao().UpdateClassroom(id, &classroom); err != nil {
		return nil, err
	}
	return &classroom, nil
}

package class

import (
	"schedule/dto"
	"schedule/models"
)

func UpdateClass(id string, req dto.ClassUpdateReq) (*models.Class, error) {
	// 检查班级是否存在
	if _, err := models.NewClassDaoInstance().GetClassByID(id); err != nil {
		return nil, NotFoundError
	}

	class := models.Class{
		ID:             req.ID,
		Name:           req.Name,
		Academic:       req.Academic,
		Cultivation:    req.Cultivation,
		Type:           req.Type,
		ExpectedYear:   req.ExpectedYear,
		IsGraduation:   req.IsGraduation,
		StudentCount:   req.StudentCount,
		MaxCount:       req.MaxCount,
		Year:           req.Year,
		Department:     req.Department,
		MajorID:        req.MajorID,
		Major:          req.Major,
		Campus:         req.Campus,
		FixedClassroom: req.FixedClassroom,
		IsFixed:        req.IsFixed,
	}

	if err := models.NewClassDaoInstance().UpdateClass(id, &class); err != nil {
		return nil, err
	}
	return &class, nil
}

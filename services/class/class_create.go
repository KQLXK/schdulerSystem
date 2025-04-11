package class

import (
	"errors"
	"schedule/dto"
	"schedule/models"
)

var (
	ExistsError = errors.New("班级编号或名称已存在")
)

func CreateClass(req dto.ClassCreateReq) (*models.Class, error) {
	// 检查班级是否已存在
	if _, err := models.NewClassDaoInstance().GetClassByID(req.ID); err == nil {
		return nil, ExistsError
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
	}

	if err := models.NewClassDaoInstance().CreateClass(&class); err != nil {
		return nil, err
	}
	return &class, nil
}

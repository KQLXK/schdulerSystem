package classroom

import (
	"schedule/dto"
	"schedule/models"
)

type ClassroomSearchFlow struct {
	SearchString string
	Count        int
	Classrooms   []models.Classroom
}

func ClassroomSearch(s string) (*dto.ClassroomSearchResp, error) {
	return NewClassroomSearchFlow(s).Do()
}

func NewClassroomSearchFlow(s string) *ClassroomSearchFlow {
	return &ClassroomSearchFlow{
		SearchString: s,
	}
}

func (f *ClassroomSearchFlow) Do() (*dto.ClassroomSearchResp, error) {
	var resp dto.ClassroomSearchResp
	if err := f.Search(); err != nil {
		return nil, err
	}
	resp.Classrooms = f.Convert()
	resp.TotalCount = int64(f.Count)
	return &resp, nil
}

func (f *ClassroomSearchFlow) Search() error {
	classrooms, err := models.NewClassroomDao().SearchClassroom(f.SearchString)
	if err != nil {
		return err
	}
	f.Classrooms = classrooms
	f.Count = len(classrooms)
	return nil
}

func (f *ClassroomSearchFlow) Convert() []dto.ClassroomGetResp {
	classroomResp := make([]dto.ClassroomGetResp, len(f.Classrooms))
	for i, classroom := range f.Classrooms {
		classroomResp[i] = dto.ClassroomGetResp{
			ID:          classroom.ID,
			Name:        classroom.Name,
			Campus:      classroom.Campus,
			Building:    classroom.Building,
			Floor:       classroom.Floor,
			Capacity:    classroom.Capacity,
			Type:        classroom.Type,
			HasAC:       classroom.HasAC,
			Description: classroom.Description,
			Department:  classroom.Department,
			Status:      classroom.Status,
		}
	}
	return classroomResp
}

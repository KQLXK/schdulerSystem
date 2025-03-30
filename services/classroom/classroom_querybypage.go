package classroom

import (
	"errors"
	"math"
	"schedule/dto"
	"schedule/models"
)

var (
	PageNumErr = errors.New("页码超出范围")
)

type ClassroomQueryByPageFlow struct {
	Page       int
	Pagesize   int
	Classrooms []models.Classroom
	Total      int64
	TotalPage  int64
}

func ClassroomQueryByPage(page int, pagesize int) (*dto.ClassroomQueryByPageResp, error) {
	return NewClassroomQueryByPageFlow(page, pagesize).Do()
}

func NewClassroomQueryByPageFlow(page int, pagesize int) *ClassroomQueryByPageFlow {
	return &ClassroomQueryByPageFlow{
		Page:     page,
		Pagesize: pagesize,
	}
}

func (f *ClassroomQueryByPageFlow) Do() (*dto.ClassroomQueryByPageResp, error) {
	var resp dto.ClassroomQueryByPageResp
	resp.Page = f.Page
	resp.PageSize = f.Pagesize

	if err := f.QueryByPage(); err != nil {
		return nil, err
	}

	resp.Classrooms = f.Convert()

	if err := f.CountTotal(); err != nil {
		return nil, err
	}

	resp.Total = f.Total
	resp.TotalPages = f.TotalPage

	if int64(f.Page) > f.TotalPage {
		return nil, PageNumErr
	}

	return &resp, nil
}

func (f *ClassroomQueryByPageFlow) QueryByPage() error {
	classrooms, err := models.NewClassroomDao().QueryByPage(f.Page, f.Pagesize)
	if err != nil {
		return err
	}
	f.Classrooms = classrooms
	return nil
}

func (f *ClassroomQueryByPageFlow) CountTotal() error {
	total, err := models.NewClassroomDao().CountTotal()
	if err != nil {
		return err
	}
	f.Total = total
	f.TotalPage = int64(math.Ceil(float64(total) / float64(f.Pagesize)))
	return nil
}

func (f *ClassroomQueryByPageFlow) Convert() []dto.ClassroomGetResp {
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

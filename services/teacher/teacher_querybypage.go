package teacher

import (
	"errors"
	"math"
	"schedule/dto"
	"schedule/models"
)

var (
	PageNumErr = errors.New("页码超出范围")
)

type TeacherQueryByPageFlow struct {
	Page      int
	Pagesize  int
	Teachers  []models.Teacher
	Total     int64
	TotalPage int64
}

func TeacherQueryByPage(page int, pagesize int) (*dto.TeacherQueryByPageResp, error) {
	return NewTeacherQueryByPageFlow(page, pagesize).Do()
}

func NewTeacherQueryByPageFlow(page int, pagesize int) *TeacherQueryByPageFlow {
	return &TeacherQueryByPageFlow{
		Page:     page,
		Pagesize: pagesize,
	}
}

func (f *TeacherQueryByPageFlow) Do() (*dto.TeacherQueryByPageResp, error) {
	var resp dto.TeacherQueryByPageResp
	resp.Page = f.Page
	resp.PageSize = f.Pagesize

	if err := f.QueryByPage(); err != nil {
		return nil, err
	}

	resp.Teachers = f.Convert()

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

func (f *TeacherQueryByPageFlow) QueryByPage() error {
	teachers, err := models.QueryTeachersByPage(f.Page, f.Pagesize)
	if err != nil {
		return err
	}
	f.Teachers = teachers
	return nil
}

func (f *TeacherQueryByPageFlow) CountTotal() error {
	total, err := models.CountTeachers()
	if err != nil {
		return err
	}
	f.Total = total
	f.TotalPage = int64(math.Ceil(float64(total) / float64(f.Pagesize)))
	return nil
}

func (f *TeacherQueryByPageFlow) Convert() []dto.TeacherGetResp {
	teacherResp := make([]dto.TeacherGetResp, len(f.Teachers))
	for i, teacher := range f.Teachers {
		teacherResp[i] = dto.TeacherGetResp{
			ID:          teacher.ID,
			Name:        teacher.Name,
			EnglishName: teacher.EnglishName,
			Gender:      teacher.Gender,
			Ethnicity:   teacher.Ethnicity,
			Department:  teacher.Department,
			Title:       teacher.Title,
			Category:    teacher.Category,
			IsExternal:  teacher.IsExternal,
			Status:      teacher.Status,
		}
	}
	return teacherResp
}

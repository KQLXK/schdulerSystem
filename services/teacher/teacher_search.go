package teacher

import (
	"schedule/dto"
	"schedule/models"
)

type TeacherSearchFlow struct {
	SearchString string
	Count        int
	Teachers     []models.Teacher
}

func TeacherSearch(s string) (*dto.TeacherSearchResp, error) {
	return NewTeacherSearchFlow(s).Do()
}

func NewTeacherSearchFlow(s string) *TeacherSearchFlow {
	return &TeacherSearchFlow{
		SearchString: s,
	}
}

func (f *TeacherSearchFlow) Do() (*dto.TeacherSearchResp, error) {
	var resp dto.TeacherSearchResp
	if err := f.Search(); err != nil {
		return nil, err
	}
	resp.Teachers = f.Convert()
	resp.TotalCount = int64(f.Count)
	return &resp, nil
}

func (f *TeacherSearchFlow) Search() error {
	teachers, err := models.SearchTeachers(f.SearchString)
	if err != nil {
		return err
	}
	f.Teachers = teachers
	f.Count = len(teachers)
	return nil
}

func (f *TeacherSearchFlow) Convert() []dto.TeacherGetResp {
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

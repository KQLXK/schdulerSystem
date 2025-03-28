package course

import (
	"schedule/dto"
	"schedule/models"
)

type CourseSearchFlow struct {
	SearchString string
	Count        int
	Courses      []models.Course
}

func CourseSearch(s string) (*dto.CourseSearchResp, error) {
	return NewCourseSearchFlow(s).Do()
}

func NewCourseSearchFlow(s string) *CourseSearchFlow {
	return &CourseSearchFlow{
		SearchString: s,
	}
}

func (f *CourseSearchFlow) Do() (*dto.CourseSearchResp, error) {
	var resp dto.CourseSearchResp
	if err := f.Search(); err != nil {
		return nil, err
	}
	resp.Courses = Convert(f.Courses)
	resp.TotalCount = int64(f.Count)
	return &resp, nil
}

func (f *CourseSearchFlow) Search() error {
	courses, err := models.NewCourseDao().SearchCourse(f.SearchString)
	if err != nil {
		return err
	}
	f.Courses = courses
	f.Count = len(courses)
	return nil
}

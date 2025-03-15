package course

import (
	"errors"
	"math"
	"schedule/dto"
	"schedule/models"
)

var (
	PageNumErr = errors.New("页码超出范围")
)

type CourseQueryByPageFlow struct {
	Page      int
	Pagesize  int
	Courses   []models.Course
	Total     int
	TotalPage int
}

func CourseQueryByPage(page int, pagesize int) (*dto.CourseQueryByPageResp, error) {
	return NewCourseQueryByPageFlow(page, pagesize).Do()
}

func NewCourseQueryByPageFlow(page int, pagesize int) *CourseQueryByPageFlow {
	return &CourseQueryByPageFlow{
		Page:     page,
		Pagesize: pagesize,
	}
}

func (f *CourseQueryByPageFlow) Do() (*dto.CourseQueryByPageResp, error) {
	var resp dto.CourseQueryByPageResp
	resp.Page = f.Page
	resp.PageSize = f.Pagesize
	if err := f.QueryByPage(); err != nil {
		return nil, err
	}
	resp.Courses = f.Convert()
	if err := f.CountTotal(); err != nil {
		return nil, err
	}
	resp.Total = int64(f.Total)
	resp.TotalPages = int64(f.TotalPage)
	if f.TotalPage < f.Page {
		return nil, PageNumErr
	}
	return &resp, nil
}

func (f *CourseQueryByPageFlow) QueryByPage() error {
	courses, err := models.NewCourseDao().QueryByPage(f.Page, f.Pagesize)
	if err != nil {
		return err
	}
	f.Courses = courses
	return nil
}

func (f *CourseQueryByPageFlow) CountTotal() error {
	total, err := models.NewCourseDao().CountTotal()
	if err != nil {
		return err
	}
	f.Total = total
	f.TotalPage = int(math.Ceil(float64(total) / float64(f.Pagesize)))
	return nil
}

func (f *CourseQueryByPageFlow) Convert() []dto.CourseGetResp {
	courseResp := make([]dto.CourseGetResp, len(f.Courses))
	for i, course := range f.Courses {
		courseResp[i] = dto.CourseGetResp{
			CourseID:         course.ID,
			CourseName:       course.Name,
			CourseType:       course.Type,
			CourseProperty:   course.Property,
			CourseCredit:     course.Credit,
			CourseDepartment: course.Department,
			TotalHour:        course.TotalHours,
			TheoryHours:      course.TheoryHours,
			TestHours:        course.TestHours,
			ComputerHours:    course.ComputerHours,
			PracticeHours:    course.PracticeHours,
			OtherHours:       course.OtherHours,
			WeeklyHours:      course.WeeklyHours,
			PurePractice:     course.PurePractice,
		}
	}
	return courseResp
}

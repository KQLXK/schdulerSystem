package course

import (
	"schedule/dto"
	"schedule/models"
)

type CourseQueryAllFlow struct {
	courses []models.Course
	total   int
}

func CourseQueryAll() (*dto.CourseQueryAllResp, error) {
	return NewCourseQueryAllFlow().Do()
}

func NewCourseQueryAllFlow() *CourseQueryAllFlow {
	return &CourseQueryAllFlow{}
}

func (f *CourseQueryAllFlow) Do() (*dto.CourseQueryAllResp, error) {
	var resp dto.CourseQueryAllResp
	if err := f.Queryall(); err != nil {
		return nil, err
	}
	if err := f.CountTotal(); err != nil {
		return nil, err
	}
	resp.Courses = Convert(f.courses)
	resp.Total = int64(f.total)
	return &resp, nil
}

func (f *CourseQueryAllFlow) CountTotal() error {
	total, err := models.NewCourseDao().CountTotal()
	if err != nil {
		return err
	}
	f.total = total
	return nil
}

func (f *CourseQueryAllFlow) Queryall() error {
	courses, err := models.NewCourseDao().GetAllCourses()
	if err != nil {
		return err
	}
	f.courses = courses
	return nil
}

func Convert(courses []models.Course) []dto.CourseGetResp {
	courseResp := make([]dto.CourseGetResp, len(courses))
	for i, course := range courses {
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

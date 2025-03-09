package course

import (
	"schedule/models"
)

// CourseGetFlow 表示获取单个课程的工作流
type CourseGetFlow struct {
	CourseID string
}

func NewCourseGetFlow(courseID string) *CourseGetFlow {
	return &CourseGetFlow{CourseID: courseID}
}

func (f *CourseGetFlow) Do() (*models.Course, error) {
	course, err := models.NewCourseDao().GetCourseByID(f.CourseID)
	if err != nil {
		return nil, DataNotFoundErr
	}
	return course, nil
}

// CourseGetAllFlow 表示获取所有课程的工作流
type CourseGetAllFlow struct{}

func NewCourseGetAllFlow() *CourseGetAllFlow {
	return &CourseGetAllFlow{}
}

func (f *CourseGetAllFlow) Do() ([]models.Course, error) {
	return models.NewCourseDao().GetAllCourses()
}

package course

import (
	"schedule/models"
)

// CourseDeleteFlow 表示删除课程的工作流
type CourseDeleteFlow struct {
	CourseID string
}

func DeleteCourse(courseId string) error {
	return NewCourseDeleteFlow(courseId).Do()
}

func NewCourseDeleteFlow(courseID string) *CourseDeleteFlow {
	return &CourseDeleteFlow{CourseID: courseID}
}

func (f *CourseDeleteFlow) Do() error {
	if err := models.NewCourseDao().DeleteCourse(f.CourseID); err != nil {
		return err
	}
	return nil
}

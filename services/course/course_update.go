package course

import (
	"errors"
	"schedule/dto"
	"schedule/models"
)

var DataNotFoundErr = errors.New("课程未找到")

// CourseUpdateFlow 表示更新课程的工作流
type CourseUpdateFlow struct {
	dto.CourseUpdateReq
}

func NewCourseUpdateFlow(req dto.CourseUpdateReq) *CourseUpdateFlow {
	return &CourseUpdateFlow{req}
}

func (f *CourseUpdateFlow) Do() error {
	if err := f.CheckCourseExists(); err != nil {
		return err
	}
	if err := f.CheckTotalHours(); err != nil {
		return err
	}
	return f.UpdateCourse()
}

func (f *CourseUpdateFlow) CheckCourseExists() error {
	if _, err := models.NewCourseDao().GetCourseByID(f.CourseID); err != nil {
		return DataNotFoundErr
	}
	return nil
}

func (f *CourseUpdateFlow) CheckTotalHours() error {
	if f.TotalHour != (f.TheoryHours + f.TestHours + f.ComputerHours + f.PracticeHours + f.OtherHours) {
		return InvalidDataErr
	}
	return nil
}

func (f *CourseUpdateFlow) UpdateCourse() error {
	course := f.reqToCourse()
	return models.NewCourseDao().UpdateCourse(f.CourseID, course)
}

func (f *CourseUpdateFlow) reqToCourse() *map[string]interface{} {
	return &map[string]interface{}{
		"ID":            f.CourseID,
		"Name":          f.CourseName,
		"Type":          f.CourseType,
		"Property":      f.CourseProperty,
		"Credit":        f.CourseCredit,
		"Department":    f.CourseDepartment,
		"TotalHours":    f.TotalHour,
		"TheoryHours":   f.TheoryHours,
		"TestHours":     f.TestHours,
		"ComputerHours": f.ComputerHours,
		"PracticeHours": f.PracticeHours,
		"OtherHours":    f.OtherHours,
		"WeeklyHours":   f.WeeklyHours,
		"PurePractice":  f.PurePractice,
	}
}

package course

import (
	"errors"
	"log"
	"schedule/dto"
	"schedule/models"
)

var (
	InvalidDataErr = errors.New("学时数据不合法")
	DataExistErr   = errors.New("课程ID或课程名已存在")
)

type CourseCreateFlow struct {
	dto.CourseCreateReq
}

func NewCourseCreateFlow(req dto.CourseCreateReq) *CourseCreateFlow {
	return &CourseCreateFlow{
		req,
	}
}

func (f *CourseCreateFlow) Do() (*dto.CourseCreateResp, error) {
	if err := f.CheckExists(); err != nil {
		return nil, err
	}
	if err := f.CheckTotalHours(); err != nil {
		return nil, err
	}
	if err := f.CreateCourse(); err != nil {
		return nil, err
	}
	return &dto.CourseCreateResp{
		CourseID:   f.CourseID,
		CourseName: f.CourseName,
	}, nil
}

func (f *CourseCreateFlow) CheckExists() error {
	if _, err := models.NewCourseDao().GetCourseByID(f.CourseID); err == nil {
		return DataExistErr
	}
	if _, err := models.NewCourseDao().GetCourseByName(f.CourseName); err == nil {
		return DataExistErr
	}
	return nil
}

func (f *CourseCreateFlow) CheckTotalHours() error {
	if f.TotalHour != (f.TheoryHours + f.TestHours + f.ComputerHours + f.PracticeHours + f.OtherHours) {
		log.Printf("TotalHour check failed: expected %d TheoryHours: %d, TestHours: %d, ComputerHours: %d, PracticeHours: %d, OtherHours: %d", f.TotalHour, f.TheoryHours, f.TestHours, f.ComputerHours, f.PracticeHours, f.OtherHours)
		return InvalidDataErr
	}
	return nil
}

func (f *CourseCreateFlow) CreateCourse() error {
	if err := models.NewCourseDao().CreateCourse(f.reqToCourse()); err != nil {
		return err
	}
	return nil
}

func (f *CourseCreateFlow) reqToCourse() *models.Course {
	return &models.Course{
		ID:            f.CourseID,
		Name:          f.CourseName,
		Type:          f.CourseType,
		Property:      f.CourseProperty,
		Credit:        f.CourseCredit,
		Department:    f.CourseDepartment,
		TotalHours:    f.TotalHour,
		TheoryHours:   f.TheoryHours,
		TestHours:     f.TestHours,
		ComputerHours: f.ComputerHours,
		PracticeHours: f.PracticeHours,
		OtherHours:    f.OtherHours,
		WeeklyHours:   f.WeeklyHours,
		PurePractice:  f.PurePractice,
		// ID 和 gorm.Model 的字段由数据库自动生成
	}
}

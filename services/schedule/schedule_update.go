package schedule

import (
	"errors"
	"schedule/dto"
	"schedule/models"
)

var (
	ErrScheduleNotFound = errors.New("排课记录不存在")
	ErrInvalidHour      = errors.New("学时数据不合法")
)

// ScheduleUpdateFlow 排课更新工作流
type ScheduleUpdateFlow struct {
	dto.ScheduleUpdateReq
}

func NewScheduleUpdateFlow(req dto.ScheduleUpdateReq) *ScheduleUpdateFlow {
	return &ScheduleUpdateFlow{req}
}

func (f *ScheduleUpdateFlow) Do() error {
	if err := f.Validate(); err != nil {
		return err
	}
	if err := f.CheckExists(); err != nil {
		return err
	}
	if err := f.CheckHours(); err != nil {
		return err
	}
	return f.UpdateSchedule()
}

// 基础数据校验
func (f *ScheduleUpdateFlow) Validate() error {
	if f.ID <= 0 {
		return errors.New("无效的排课ID")
	}
	if f.Semester == "" || f.CourseID == "" || f.TeacherID == "" {
		return errors.New("关键字段不能为空")
	}
	if f.TeachingClassSize < 0 {
		return errors.New("教学班人数不能为负数")
	}
	if f.ContinuousPeriods < 1 {
		return errors.New("连排节次至少1节")
	}
	return nil
}

// 检查排课是否存在
func (f *ScheduleUpdateFlow) CheckExists() error {
	if _, err := models.NewScheduleDao().GetScheduleByID(int(f.ID)); err != nil {
		return ErrScheduleNotFound
	}
	return nil
}

// 学时合理性校验
func (f *ScheduleUpdateFlow) CheckHours() error {
	if f.OpeningHours < f.SchedulingHours {
		return ErrInvalidHour
	}
	return nil
}

// 执行更新操作
func (f *ScheduleUpdateFlow) UpdateSchedule() error {
	updateData := f.reqToMap()
	return models.NewScheduleDao().UpdateSchedule(int(f.ID), updateData)
}

// 转换为更新字段映射
func (f *ScheduleUpdateFlow) reqToMap() map[string]interface{} {
	return map[string]interface{}{
		"Semester":               f.Semester,
		"CourseID":               f.CourseID,
		"CourseName":             f.CourseName,
		"TeacherID":              f.TeacherID,
		"TeachingClass":          f.TeachingClass,
		"TeachingClassID":        f.TeachingClassID,
		"TeachingClassName":      f.TeachingClassName,
		"HourType":               f.HourType,
		"OpeningHours":           f.OpeningHours,
		"SchedulingHours":        f.SchedulingHours,
		"TotalHours":             f.TotalHours,
		"SchedulingPriority":     f.SchedulingPriority,
		"TeachingClassSize":      f.TeachingClassSize,
		"OpeningCampus":          f.OpeningCampus,
		"OpeningWeekHours":       f.OpeningWeekHours,
		"ContinuousPeriods":      f.ContinuousPeriods,
		"SpecifiedClassroomType": f.SpecifiedClassroomType,
		"SpecifiedClassroom":     f.SpecifiedClassroom,
		"SpecifiedBuilding":      f.SpecifiedBuilding,
		"SpecifiedTime":          f.SpecifiedTime,
	}
}

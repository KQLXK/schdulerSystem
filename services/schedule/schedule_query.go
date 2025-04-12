package schedule

import (
	"errors"
	"gorm.io/gorm"
	"schedule/dto"
	"schedule/models"
)

// ScheduleGetFlow 排课详情查询工作流
type ScheduleGetFlow struct {
	ScheduleID int
}

// 快速查询入口
func QuerySchedule(scheduleID int) (*dto.ScheduleGetResp, error) {
	return NewScheduleGetFlow(scheduleID).Do()
}

// 创建查询实例
func NewScheduleGetFlow(scheduleID int) *ScheduleGetFlow {
	return &ScheduleGetFlow{ScheduleID: scheduleID}
}

// 执行查询操作
func (f *ScheduleGetFlow) Do() (*dto.ScheduleGetResp, error) {
	// 可选：添加ID有效性校验
	if f.ScheduleID <= 0 {
		return nil, errors.New("无效的排课ID")
	}

	schedule, err := models.NewScheduleDao().GetScheduleByID(f.ScheduleID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrScheduleNotFound
		}
		return nil, err // 返回原始数据库错误
	}
	return ScheduleToResp(schedule), nil
}

func ScheduleToResp(schedule *models.Schedule) *dto.ScheduleGetResp {
	return &dto.ScheduleGetResp{
		ID:                     schedule.ID,
		CreatedAt:              schedule.CreatedAt,
		UpdatedAt:              schedule.UpdatedAt,
		Semester:               schedule.Semester,
		CourseID:               schedule.CourseID,
		CourseName:             schedule.CourseName,
		TeacherID:              schedule.TeacherID,
		TeachingClass:          schedule.TeachingClass,
		TeachingClassID:        schedule.TeachingClassID,
		TeachingClassName:      schedule.TeachingClassName,
		HourType:               schedule.HourType,
		OpeningHours:           schedule.OpeningHours,
		SchedulingHours:        schedule.SchedulingHours,
		TotalHours:             schedule.TotalHours,
		SchedulingPriority:     schedule.SchedulingPriority,
		TeachingClassSize:      schedule.TeachingClassSize,
		OpeningCampus:          schedule.OpeningCampus,
		OpeningWeekHours:       schedule.OpeningWeekHours,
		ContinuousPeriods:      schedule.ContinuousPeriods,
		SpecifiedClassroomType: schedule.SpecifiedClassroomType,
		SpecifiedClassroom:     schedule.SpecifiedClassroom,
		SpecifiedBuilding:      schedule.SpecifiedBuilding,
		SpecifiedTime:          schedule.SpecifiedTime,
	}
}

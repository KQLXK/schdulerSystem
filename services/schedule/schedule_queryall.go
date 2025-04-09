package schedule

import (
	"schedule/dto"
	"schedule/models"
)

type ScheduleQueryAllFlow struct {
	schedules []models.Schedule
	total     int64
}

// 快速查询入口
func QueryAllSchedules() (*dto.ScheduleQueryAllResp, error) {
	return NewScheduleQueryAllFlow().Do()
}

func NewScheduleQueryAllFlow() *ScheduleQueryAllFlow {
	return &ScheduleQueryAllFlow{}
}

func (f *ScheduleQueryAllFlow) Do() (*dto.ScheduleQueryAllResp, error) {
	var resp dto.ScheduleQueryAllResp

	if err := f.QueryAll(); err != nil {
		return nil, err
	}

	if err := f.CountTotal(); err != nil {
		return nil, err
	}

	resp.Schedules = ConvertToDTO(f.schedules)
	resp.Total = int64(f.total)

	return &resp, nil
}

func (f *ScheduleQueryAllFlow) QueryAll() error {
	schedules, err := models.NewScheduleDao().GetAllSchedules()
	if err != nil {
		return err
	}
	f.schedules = schedules
	return nil
}

func (f *ScheduleQueryAllFlow) CountTotal() error {
	total, err := models.NewScheduleDao().CountSchedules()
	if err != nil {
		return err
	}
	f.total = total
	return nil
}

func ConvertToDTO(schedules []models.Schedule) []dto.ScheduleGetResp {
	scheduleResp := make([]dto.ScheduleGetResp, len(schedules))
	for i, s := range schedules {
		scheduleResp[i] = dto.ScheduleGetResp{
			ID:                     s.ID,
			CreatedAt:              s.CreatedAt,
			UpdatedAt:              s.UpdatedAt,
			Semester:               s.Semester,
			CourseID:               s.CourseID,
			CourseName:             s.CourseName,
			TeacherID:              s.TeacherID,
			TeachingClass:          s.TeachingClass,
			TeachingClassID:        s.TeachingClassID,
			TeachingClassName:      s.TeachingClassName,
			HourType:               s.HourType,
			OpeningHours:           int64(s.OpeningHours),
			SchedulingHours:        int64(s.SchedulingHours),
			TotalHours:             int64(s.TotalHours),
			SchedulingPriority:     int64(s.SchedulingPriority),
			TeachingClassSize:      int64(s.TeachingClassSize),
			OpeningCampus:          s.OpeningCampus,
			OpeningWeekHours:       s.OpeningWeekHours,
			ContinuousPeriods:      int64(s.ContinuousPeriods),
			SpecifiedClassroomType: s.SpecifiedClassroomType,
			SpecifiedClassroom:     s.SpecifiedClassroom,
			SpecifiedBuilding:      s.SpecifiedBuilding,
			SpecifiedTime:          s.SpecifiedTime,
		}
	}
	return scheduleResp
}

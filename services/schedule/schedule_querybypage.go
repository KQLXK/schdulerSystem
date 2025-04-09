package schedule

import (
	"errors"
	"math"
	"schedule/dto"
	"schedule/models"
)

var (
	ErrPageOutOfRange = errors.New("页码超出范围")
)

type ScheduleQueryByPageFlow struct {
	Page       int
	PageSize   int
	Schedules  []models.Schedule
	Total      int
	TotalPages int
}

// 快速查询入口
func QuerySchedulesByPage(page, pageSize int) (*dto.ScheduleQueryByPageResp, error) {
	return NewScheduleQueryByPageFlow(page, pageSize).Do()
}

func NewScheduleQueryByPageFlow(page, pageSize int) *ScheduleQueryByPageFlow {
	return &ScheduleQueryByPageFlow{
		Page:     page,
		PageSize: pageSize,
	}
}

func (f *ScheduleQueryByPageFlow) Do() (*dto.ScheduleQueryByPageResp, error) {
	var resp dto.ScheduleQueryByPageResp

	// 设置基础分页信息
	resp.Page = int64(f.Page)
	resp.PageSize = int64(f.PageSize)

	// 执行分页查询
	if err := f.QueryByPage(); err != nil {
		return nil, err
	}

	// 转换数据格式
	resp.Schedules = f.ConvertToDTO()

	// 获取总数和总页数
	if err := f.CountTotal(); err != nil {
		return nil, err
	}

	// 设置总数信息
	resp.Total = int64(f.Total)
	resp.TotalPages = int64(f.TotalPages)

	// 检查页码有效性
	if f.TotalPages > 0 && f.Page > f.TotalPages {
		return nil, ErrPageOutOfRange
	}

	return &resp, nil
}

func (f *ScheduleQueryByPageFlow) QueryByPage() error {
	schedules, err := models.NewScheduleDao().QueryByPage(f.Page, f.PageSize)
	if err != nil {
		return err
	}
	f.Schedules = schedules
	return nil
}

func (f *ScheduleQueryByPageFlow) CountTotal() error {
	total, err := models.NewScheduleDao().CountSchedules()
	if err != nil {
		return err
	}
	f.Total = int(total)
	f.TotalPages = int(math.Ceil(float64(total) / float64(f.PageSize)))
	return nil
}

func (f *ScheduleQueryByPageFlow) ConvertToDTO() []dto.ScheduleGetResp {
	scheduleResp := make([]dto.ScheduleGetResp, len(f.Schedules))
	for i, s := range f.Schedules {
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

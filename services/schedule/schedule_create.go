package schedule

import (
	"errors"
	"log"
	"schedule/dto"
	"schedule/models"
	"strings"
	"time"
)

var (
	ErrInvalidSemester   = errors.New("学期格式错误，示例：2023-2024-秋季")
	ErrInvalidTimeFormat = errors.New("指定时间格式错误，示例：Mon-08:00")
	ErrInvalidHours      = errors.New("开课学时不能小于排课学时")
	//ErrScheduleConflict   = errors.New("排课时间或教室冲突")
	//ErrTeacherUnavailable = errors.New("教师在该时间段已有排课")
)

type ScheduleCreateFlow struct {
	ScheduleId int64
	dto.ScheduleCreateReq
}

func CreateSchedule(req dto.ScheduleCreateReq) (*dto.ScheduleCreateResp, error) {
	return NewScheduleCreateFlow(req).Do()
}

func NewScheduleCreateFlow(req dto.ScheduleCreateReq) *ScheduleCreateFlow {
	return &ScheduleCreateFlow{
		ScheduleCreateReq: req,
	}
}

func (f *ScheduleCreateFlow) Do() (*dto.ScheduleCreateResp, error) {
	if err := f.CheckData(); err != nil {
		return nil, err
	}
	if err := f.ValidateSemester(); err != nil {
		return nil, err
	}
	if err := f.ValidateTimeFormat(); err != nil {
		return nil, err
	}
	if err := f.CheckHours(); err != nil {
		return nil, err
	}
	//if err := f.CheckConflicts(); err != nil {
	//	return nil, err
	//}
	if err := f.CreateSchedule(); err != nil {
		return nil, err
	}

	return &dto.ScheduleCreateResp{
		ScheduleID: f.ScheduleId, // 实际应替换为数据库生成的ID
		Semester:   f.Semester,
	}, nil
}

// 基础数据校验
func (f *ScheduleCreateFlow) CheckData() error {
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

// 学期格式校验
func (f *ScheduleCreateFlow) ValidateSemester() error {
	// 实现实际的学期格式正则校验
	// 示例简单校验：2023-2024-秋季
	parts := strings.Split(f.Semester, "-")
	if len(parts) != 3 {
		return ErrInvalidSemester
	}
	return nil
}

//func isValidSeason(season string) bool {
//	valid := []string{"春季", "夏季", "秋季", "冬季"}
//	for _, v := range valid {
//		if v == season {
//			return true
//		}
//	}
//	return false
//}

// 时间格式校验
func (f *ScheduleCreateFlow) ValidateTimeFormat() error {
	if f.SpecifiedTime == "" {
		return nil
	}
	// 示例格式校验：Mon-08:00
	parts := strings.Split(f.SpecifiedTime, "-")
	if len(parts) != 2 || !isValidWeekday(parts[0]) || !isValidTime(parts[1]) {
		return ErrInvalidTimeFormat
	}
	return nil
}

// 学时合理性校验
func (f *ScheduleCreateFlow) CheckHours() error {
	if f.OpeningHours < f.SchedulingHours {
		log.Printf("学时数据异常 OpeningHours: %d, SchedulingHours: %d",
			f.OpeningHours, f.SchedulingHours)
		return ErrInvalidHours
	}
	return nil
}

//// 冲突检查（需要依赖DAO实现）
//func (f *ScheduleCreateFlow) CheckConflicts() error {
//	dao := models.NewScheduleDao()
//
//	// 检查教师时间冲突
//	if exists, err := dao.CheckTeacherSchedule(f.TeacherID, f.SpecifiedTime); err != nil {
//		return err
//	} else if exists {
//		return ErrTeacherUnavailable
//	}
//
//	// 检查教室冲突
//	if f.SpecifiedClassroom != "" {
//		if exists, err := dao.CheckClassroomSchedule(
//			f.SpecifiedBuilding,
//			f.SpecifiedClassroom,
//			f.SpecifiedTime,
//		); err != nil {
//			return err
//		} else if exists {
//			return ErrScheduleConflict
//		}
//	}
//	return nil
//}

// 创建排课记录
func (f *ScheduleCreateFlow) CreateSchedule() error {
	schedule := f.reqToSchedule()
	if err := models.NewScheduleDao().CreateSchedule(schedule); err != nil {
		log.Printf("创建排课任务失败: %v", err)
		return err
	}
	f.ScheduleId = schedule.ID
	return nil
}

// DTO转模型
func (f *ScheduleCreateFlow) reqToSchedule() *models.Schedule {
	return &models.Schedule{
		Semester:               f.Semester,
		CourseID:               f.CourseID,
		CourseName:             f.CourseName,
		TeacherID:              f.TeacherID,
		TeachingClass:          f.TeachingClass,
		TeachingClassID:        f.TeachingClassID,
		TeachingClassName:      f.TeachingClassName,
		HourType:               f.HourType,
		OpeningHours:           f.OpeningHours,
		SchedulingHours:        f.SchedulingHours,
		TotalHours:             f.TotalHours,
		SchedulingPriority:     f.SchedulingPriority,
		TeachingClassSize:      f.TeachingClassSize,
		OpeningCampus:          f.OpeningCampus,
		OpeningWeekHours:       f.OpeningWeekHours,
		ContinuousPeriods:      f.ContinuousPeriods,
		SpecifiedClassroomType: f.SpecifiedClassroomType,
		SpecifiedClassroom:     f.SpecifiedClassroom,
		SpecifiedBuilding:      f.SpecifiedBuilding,
		SpecifiedTime:          f.SpecifiedTime,
		// 自动填充字段由数据库处理
	}
}

// 辅助校验方法
func isValidWeekday(day string) bool {
	weekdays := []string{"Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"}
	for _, d := range weekdays {
		if d == day {
			return true
		}
	}
	return false
}

func isValidTime(t string) bool {
	_, err := time.Parse("15:04", t)
	return err == nil
}

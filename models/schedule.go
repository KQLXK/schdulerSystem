package models

import (
	"gorm.io/gorm"
	"schedule/database"
	"sync"
	"time"
)

// Schedule 定义了排课结果表的结构
type Schedule struct {
	gorm.Model            // 内嵌 gorm.Model
	CourseID    string    `gorm:"type:varchar(20);not null"` // 课程ID
	TeacherID   string    `gorm:"type:varchar(10);not null"` // 教师ID
	ClassroomID string    `gorm:"type:varchar(20);not null"` // 教室ID
	ClassID     string    `gorm:"type:varchar(20);not null"` // 班级ID
	StartTime   time.Time `gorm:"type:datetime"`             // 开始时间
	EndTime     time.Time `gorm:"type:datetime"`             // 结束时间
	WeekPattern string    `gorm:"type:varchar(20)"`          // 周次模式（单周、双周、全周）
	Status      string    `gorm:"type:varchar(20)"`          // 状态;default:'未排'
}

type ScheduleDao struct{}

var (
	ScheduleOnce sync.Once
	scheduleDao  *ScheduleDao
)

func NewScheduleDao() *ScheduleDao {
	ScheduleOnce.Do(func() {
		scheduleDao = &ScheduleDao{}
	})
	return scheduleDao
}

// GetAllSchedules 获取所有排课结果
func (ScheduleDao) GetAllSchedules() ([]Schedule, error) {
	var schedules []Schedule
	if err := database.DB.Find(&schedules).Error; err != nil {
		return nil, err
	}
	return schedules, nil
}

// GetScheduleByID 根据ID获取排课结果
func (ScheduleDao) GetScheduleByID(id string) (*Schedule, error) {
	var schedule Schedule
	if err := database.DB.Where("id = ?", id).First(&schedule).Error; err != nil {
		return nil, err
	}
	return &schedule, nil
}

// CreateSchedule 创建排课结果
func (ScheduleDao) CreateSchedule(schedule *Schedule) error {
	if err := database.DB.Create(schedule).Error; err != nil {
		return err
	}
	return nil
}

// UpdateSchedule 更新排课结果
func (ScheduleDao) UpdateSchedule(id string, schedule *Schedule) error {
	if err := database.DB.Model(&Schedule{}).Where("id = ?", id).Updates(schedule).Error; err != nil {
		return err
	}
	return nil
}

// DeleteSchedule 删除排课结果
func (ScheduleDao) DeleteSchedule(id string) error {
	if err := database.DB.Where("id = ?", id).Delete(&Schedule{}).Error; err != nil {
		return err
	}
	return nil
}

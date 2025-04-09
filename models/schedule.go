package models

import (
	"gorm.io/gorm"
	"log"
	"schedule/database"
	"sync"
)

// Schedule 定义了排课结果表的结构
type Schedule struct {
	gorm.Model
	ID                     int64  `gorm:"column:id;primaryKey;autoIncrement"`                // 排课任务编号
	Semester               string `gorm:"column:semester; type:varchar(20)"`                 // 学年学期，格式为 "年份-年份-学期"
	CourseID               string `gorm:"column:course_id;type:varchar(20);not null"`        // 课程编号
	Course                 Course `gorm:"foreignKey:CourseID"`                               // 关联的课程
	CourseName             string `gorm:"column:course_name; type:varchar(100)"`             // 课程名称
	TeacherID              string `gorm:"column:teacher_id; type:varchar(10)"`               // 教师工号
	TeachingClass          string `gorm:"column:teaching_class; type:varchar(100)"`          // 教学班组成
	TeachingClassID        string `gorm:"column:teaching_class_number; type:varchar(50)"`    // 教学班编号
	TeachingClassName      string `gorm:"column:teaching_class_name; type:varchar(100)"`     // 教学班名称
	HourType               string `gorm:"column:hour_type; type:varchar(20)"`                // 学时类型
	OpeningHours           int64  `gorm:"column:opening_hours"`                              // 开课学时
	SchedulingHours        int64  `gorm:"column:scheduling_hours"`                           // 排课学时
	TotalHours             int64  `gorm:"column:total_hours"`                                // 总学时
	SchedulingPriority     int64  `gorm:"column:scheduling_priority"`                        // 排课优先级
	TeachingClassSize      int64  `gorm:"column:teaching_class_size"`                        // 教学班人数
	OpeningCampus          string `gorm:"column:opening_campus; type:varchar(50)"`           // 开课校区
	OpeningWeekHours       string `gorm:"column:opening_week_hours; type:varchar(20)"`       // 开课周次学时
	ContinuousPeriods      int64  `gorm:"column:continuous_periods"`                         // 连排节次
	SpecifiedClassroomType string `gorm:"column:specified_classroom_type; type:varchar(50)"` // 指定教室类型
	SpecifiedClassroom     string `gorm:"column:specified_classroom; type:varchar(50)"`      // 指定教室
	SpecifiedBuilding      string `gorm:"column:specified_building; type:varchar(50)"`       // 指定教学楼
	SpecifiedTime          string `gorm:"column:specified_time; type:varchar(50)"`           // 指定时间
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
func (ScheduleDao) GetScheduleByID(id int) (*Schedule, error) {
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
func (ScheduleDao) UpdateSchedule(id int, schedule map[string]interface{}) error {
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

func (ScheduleDao) CountSchedules() (int64, error) {
	var total int64
	if err := database.DB.Model(Schedule{}).Order("created_at DESC").Count(&total).Error; err != nil {
		log.Println("count schedules failed, err:", err)
		return -1, err
	}
	return total, nil
}

func (ScheduleDao) QueryByPage(page int, pagesize int) ([]Schedule, error) {
	var schedules []Schedule
	offset := (page - 1) * pagesize
	if err := database.DB.Model(Schedule{}).Offset(offset).Find(&schedules).Limit(pagesize).Error; err != nil {
		log.Println("query schedule by page failed, err:", err)
		return nil, err
	}
	log.Println("query schedule by page sucess")
	return schedules, nil
}

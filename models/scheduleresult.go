package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
	"schedule/database"
)

type ScheduleResult struct {
	gorm.Model
	ID          uint          `gorm:"primaryKey;autoIncrement"` // 主键ID
	Semester    string        `gorm:"column:semester;varchar(20)"`
	ScheduleID  int           `gorm:"column:schedule_id;varchar(20)"` // 排课ID
	Schedule    Schedule      `gorm:"foreignKey:ScheduleID"`
	CourseID    string        `gorm:"column:course_id"`               // 课程ID，外键
	CourseName  string        `gorm:"column:course_name;varchar(20)"` // 课程名称
	Course      Course        `gorm:"foreignKey:CourseID"`
	ClassroomID string        `gorm:"column:classroom_id;varchar(20)"` // 教室ID，外键
	Classroom   Classroom     `gorm:"foreignKey:ClassroomID"`          // 教室名称
	TeacherID   string        `gorm:"column:teacher_id;varchar(20)"`   // 教师ID，外键
	TeacherName string        `gorm:"column:teacher_name;varchar(20)"`
	Teacher     Teacher       `gorm:"foreignKey:TeacherID"` // 教师名称
	TimeSlots   JSONTimeSlots `gorm:"type:json"`            // 时间槽（JSON类型），可以直接嵌套
}

// 时间段结构
type TimeSlot struct {
	WeekNumbers []int
	Weekday     int // 1-5（周一到周五）
	StartPeriod int // 开始节次（1-12）
	Duration    int // 持续节数
}

// 1. 定义自定义JSON类型
type JSONTimeSlots []TimeSlot

// 实现Valuer接口（写入数据库时自动序列化）
func (ts JSONTimeSlots) Value() (driver.Value, error) {
	return json.Marshal(ts)
}

// 实现Scanner接口（从数据库读取时自动反序列化）
func (ts *JSONTimeSlots) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("无法解析数据库字段：非字节数组类型")
	}
	return json.Unmarshal(bytes, ts)
}

// GetAllScheduleResults 获取所有排课结果
func GetAllScheduleResults() ([]ScheduleResult, error) {
	var scheduleResults []ScheduleResult
	if err := database.DB.Find(&scheduleResults).Error; err != nil {
		return nil, err
	}
	return scheduleResults, nil
}

// GetScheduleResultByID 根据ID获取排课结果
func GetScheduleResultByID(id uint) (*ScheduleResult, error) {
	var scheduleResult ScheduleResult
	if err := database.DB.Where("id = ?", id).First(&scheduleResult).Error; err != nil {
		return nil, err
	}
	return &scheduleResult, nil
}

// CreateScheduleResult 创建排课结果
func CreateScheduleResult(scheduleResult *ScheduleResult) error {
	if err := database.DB.Create(scheduleResult).Error; err != nil {
		return err
	}
	return nil
}

// UpdateScheduleResult 更新排课结果信息
func UpdateScheduleResult(id uint, scheduleResult *ScheduleResult) error {
	if err := database.DB.Model(&ScheduleResult{}).Where("id = ?", id).Updates(scheduleResult).Error; err != nil {
		return err
	}
	return nil
}

// DeleteScheduleResult 删除排课结果
func DeleteScheduleResult(id uint) error {
	if err := database.DB.Where("id = ?", id).Delete(&ScheduleResult{}).Error; err != nil {
		return err
	}
	return nil
}

// QueryScheduleResultsByPage 分页查询排课结果
func QueryScheduleResultsByPage(page int, pageSize int) ([]ScheduleResult, error) {
	var scheduleResults []ScheduleResult
	offset := (page - 1) * pageSize
	if err := database.DB.Model(&ScheduleResult{}).Order("created_at DESC").Limit(pageSize).Offset(offset).Find(&scheduleResults).Error; err != nil {
		return nil, err
	}
	return scheduleResults, nil
}

// CountScheduleResults 计算排课结果总数
func CountScheduleResults() (int64, error) {
	var total int64
	if err := database.DB.Model(&ScheduleResult{}).Count(&total).Error; err != nil {
		return 0, err
	}
	return total, nil
}

// SearchScheduleResults 按课程名称或教师名称模糊查询排课结果
func SearchScheduleResults(s string) ([]ScheduleResult, error) {
	var scheduleResults []ScheduleResult
	if err := database.DB.Model(&ScheduleResult{}).
		Where("course_name LIKE ? OR teacher LIKE ?", "%"+s+"%", "%"+s+"%").
		Find(&scheduleResults).Error; err != nil {
		return nil, err
	}
	return scheduleResults, nil
}

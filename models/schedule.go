package models

import (
	"gorm.io/gorm"
	"time"
)

// Schedule 定义了排课结果表的结构
type Schedule struct {
	gorm.Model            // 内嵌 gorm.Model
	CourseID    string    `gorm:"type:varchar(20);not null"`     // 课程ID
	TeacherID   string    `gorm:"type:varchar(10);not null"`     // 教师ID
	ClassroomID string    `gorm:"type:varchar(20);not null"`     // 教室ID
	ClassID     string    `gorm:"type:varchar(20);not null"`     // 班级ID
	StartTime   time.Time `gorm:"type:datetime"`                 // 开始时间
	EndTime     time.Time `gorm:"type:datetime"`                 // 结束时间
	WeekPattern string    `gorm:"type:varchar(20)"`              // 周次模式（单周、双周、全周）
	Status      string    `gorm:"type:varchar(20);default:'未排'"` // 状态
}

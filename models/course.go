package models

import (
	"gorm.io/gorm"
	"schedule/database"
	"sync"
)

// Course 定义了课程表的结构
type Course struct {
	gorm.Model         // 内嵌 gorm.Model
	ID         string  `gorm:"primaryKey;type:varchar(20)"`   // 课程编号
	Name       string  `gorm:"type:varchar(100);not null"`    // 课程名称
	Type       string  `gorm:"type:varchar(20)"`              // 课程类型（理论、实践、实验）
	Credit     float64 `gorm:"type:float"`                    // 学分
	Department string  `gorm:"type:varchar(50)"`              // 开课院系
	TotalHours int     `gorm:"type:int"`                      // 总学时
	Status     string  `gorm:"type:varchar(20);default:'启用'"` // 状态
}

type CourseDao struct{}

var (
	CourseOnce sync.Once
	courseDao  *CourseDao
)

func NewCourseDao() *CourseDao {
	CourseOnce.Do(func() {
		courseDao = &CourseDao{}
	})
	return courseDao
}

// GetAllCourses 获取所有课程
func (CourseDao) GetAllCourses() ([]Course, error) {
	var courses []Course
	if err := database.DB.Find(&courses).Error; err != nil {
		return nil, err
	}
	return courses, nil
}

// GetCourseByID 根据ID获取课程
func (CourseDao) GetCourseByID(id string) (*Course, error) {
	var course Course
	if err := database.DB.Where("id = ?", id).First(&course).Error; err != nil {
		return nil, err
	}
	return &course, nil
}

// CreateCourse 创建课程
func (CourseDao) CreateCourse(course *Course) error {
	if err := database.DB.Create(course).Error; err != nil {
		return err
	}
	return nil
}

// UpdateCourse 更新课程信息
func (CourseDao) UpdateCourse(id string, course *Course) error {
	if err := database.DB.Model(&Course{}).Where("id = ?", id).Updates(course).Error; err != nil {
		return err
	}
	return nil
}

// DeleteCourse 删除课程
func (CourseDao) DeleteCourse(id string) error {
	if err := database.DB.Where("id = ?", id).Delete(&Course{}).Error; err != nil {
		return err
	}
	return nil
}

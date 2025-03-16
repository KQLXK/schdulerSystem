package models

import (
	"gorm.io/gorm"
	"schedule/database"
	"sync"
)

// Classroom 定义了教室表的结构
type Classroom struct {
	gorm.Model        // 内嵌 gorm.Model
	ID         string `gorm:"primaryKey;type:varchar(20)"` // 教室编号
	Name       string `gorm:"type:varchar(100);not null"`  // 教室名称
	Campus     string `gorm:"type:varchar(50)"`            // 校区
	Building   string `gorm:"type:varchar(50)"`            // 教学楼
	Capacity   int    `gorm:"type:int"`                    // 容量
	Type       string `gorm:"type:varchar(50)"`            // 教室类型（普通教室、多媒体教室等）
	Status     string `gorm:"type:varchar(20)"`            // 状态;default:'启用'
}

type ClassroomDao struct{}

var (
	ClassroomOnce sync.Once
	classroomDao  *ClassroomDao
)

// NewClassroomDao 返回 ClassroomDao 的单例实例
func NewClassroomDao() *ClassroomDao {
	ClassroomOnce.Do(func() {
		classroomDao = &ClassroomDao{}
	})
	return classroomDao
}

// GetAllClassrooms 获取所有教室
func (ClassroomDao) GetAllClassrooms() ([]Classroom, error) {
	var classrooms []Classroom
	if err := database.DB.Find(&classrooms).Error; err != nil {
		return nil, err
	}
	return classrooms, nil
}

// GetClassroomByID 根据ID获取教室
func (ClassroomDao) GetClassroomByID(id string) (*Classroom, error) {
	var classroom Classroom
	if err := database.DB.Where("id = ?", id).First(&classroom).Error; err != nil {
		return nil, err
	}
	return &classroom, nil
}

// CreateClassroom 创建教室
func (ClassroomDao) CreateClassroom(classroom *Classroom) error {
	if err := database.DB.Create(classroom).Error; err != nil {
		return err
	}
	return nil
}

// UpdateClassroom 更新教室信息
func (ClassroomDao) UpdateClassroom(id string, classroom *Classroom) error {
	if err := database.DB.Model(&Classroom{}).Where("id = ?", id).Updates(classroom).Error; err != nil {
		return err
	}
	return nil
}

// DeleteClassroom 删除教室
func (ClassroomDao) DeleteClassroom(id string) error {
	if err := database.DB.Where("id = ?", id).Delete(&Classroom{}).Error; err != nil {
		return err
	}
	return nil
}

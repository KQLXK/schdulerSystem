package models

import (
	"gorm.io/gorm"
	"schedule/database"
	"sync"
)

// Classroom 定义了教室表的结构
type Classroom struct {
	gorm.Model         // 内嵌 gorm.Model
	ID          string `gorm:"primaryKey;type:varchar(20)"` // 教室编号
	Name        string `gorm:"type:varchar(100);not null"`  // 教室名称
	Campus      string `gorm:"type:varchar(50)"`            // 校区
	Building    string `gorm:"type:varchar(50)"`            // 教学楼
	Floor       string `gorm:"type:varchar(20)"`            // 所在楼层
	Capacity    int    `gorm:"type:int"`                    // 容量
	Type        string `gorm:"type:varchar(50)"`            // 教室类型（普通教室、多媒体教室等）
	HasAC       bool   `gorm:"type:boolean;default:false"`  // 是否有空调
	Description string `gorm:"type:text"`                   // 教室描述
	Department  string `gorm:"type:varchar(50)"`            // 管理部门
	Status      string `gorm:"type:varchar(20)"`            // 状态;default:'启用'
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

func (ClassroomDao) QueryByPage(page int, pagesize int) ([]Classroom, error) {
	var classrooms []Classroom
	offset := (page - 1) * pagesize
	if err := database.DB.Model(&Classroom{}).Order("created_at DESC").Limit(pagesize).Offset(offset).Find(&classrooms).Error; err != nil {
		return nil, err
	}
	return classrooms, nil
}

func (ClassroomDao) CountTotal() (int64, error) {
	var total int64
	if err := database.DB.Model(&Classroom{}).Count(&total).Error; err != nil {
		return 0, err
	}
	return total, nil
}

func (ClassroomDao) SearchClassroom(s string) ([]Classroom, error) {
	var classrooms []Classroom
	if err := database.DB.Model(&Classroom{}).
		Where("name LIKE ? OR id LIKE ? OR building LIKE ?",
			"%"+s+"%", "%"+s+"%", "%"+s+"%").
		Find(&classrooms).Error; err != nil {
		return nil, err
	}
	return classrooms, nil
}

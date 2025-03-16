package models

import (
	"gorm.io/gorm"
	"schedule/database"
)

// Teacher 定义了教师表的结构
type Teacher struct {
	gorm.Model         // 内嵌 gorm.Model，包含 ID、CreatedAt、UpdatedAt、DeletedAt 字段
	ID          string `gorm:"primaryKey;type:varchar(10)"` // 教师工号
	Name        string `gorm:"type:varchar(50);not null"`   // 教师姓名
	EnglishName string `gorm:"type:varchar(50)"`            // 英文姓名
	Gender      string `gorm:"type:varchar(10)"`            // 性别
	Ethnicity   string `gorm:"type:varchar(20)"`            // 民族
	Department  string `gorm:"type:varchar(50)"`            // 所属院系
	Title       string `gorm:"type:varchar(50)"`            // 职称
	Category    string `gorm:"type:varchar(50)"`            // 教职工类别
	IsExternal  bool   `gorm:"default:false"`               // 是否外聘
	Status      string `gorm:"type:varchar(20)"`            // 状态
}

// GetAllTeachers 获取所有教师
func GetAllTeachers() ([]Teacher, error) {
	var teachers []Teacher
	if err := database.DB.Find(&teachers).Error; err != nil {
		return nil, err
	}
	return teachers, nil
}

// GetTeacherByID 根据ID获取教师
func GetTeacherByID(id string) (*Teacher, error) {
	var teacher Teacher
	if err := database.DB.Where("id = ?", id).First(&teacher).Error; err != nil {
		return nil, err
	}
	return &teacher, nil
}

// CreateTeacher 创建教师
func CreateTeacher(teacher *Teacher) error {
	if err := database.DB.Create(teacher).Error; err != nil {
		return err
	}
	return nil
}

// UpdateTeacher 更新教师信息
func UpdateTeacher(id string, teacher *Teacher) error {
	if err := database.DB.Model(&Teacher{}).Where("id = ?", id).Updates(teacher).Error; err != nil {
		return err
	}
	return nil
}

// DeleteTeacher 删除教师
func DeleteTeacher(id string) error {
	if err := database.DB.Where("id = ?", id).Delete(&Teacher{}).Error; err != nil {
		return err
	}
	return nil
}

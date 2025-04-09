package models

import (
	"gorm.io/gorm"
	"log"
	"schedule/database"
	"sync"
)

// Class 定义了班级表的结构
type Class struct {
	gorm.Model          // 内嵌 gorm.Model
	ID           string `gorm:"primaryKey;type:varchar(20)"` // 班级编号
	Name         string `gorm:"type:varchar(100);not null"`  // 班级名称
	Department   string `gorm:"type:varchar(50)"`            // 所属院系
	Major        string `gorm:"type:varchar(50)"`            // 专业
	Campus       string `gorm:"type:varchar(50)"`            // 校区
	StudentCount int    `gorm:"type:int"`                    // 班级人数
	Status       string `gorm:"type:varchar(20)"`            // 状态;default:'启用'
}

type ClassDao struct{}

var (
	ClassOnce sync.Once
	classDao  *ClassDao
)

func NewClassDaoInstance() *ClassDao {
	ClassOnce.Do(func() {
		classDao = &ClassDao{}
	})
	return classDao
}

// GetAllClasses 获取所有班级
func (ClassDao) GetAllClasses() ([]Class, error) {
	var classes []Class
	if err := database.DB.Find(&classes).Error; err != nil {
		return nil, err
	}
	return classes, nil
}

// GetClassByID 根据ID获取班级
func (ClassDao) GetClassByID(id string) (*Class, error) {
	var class Class
	if err := database.DB.Where("id = ?", id).First(&class).Error; err != nil {
		return nil, err
	}
	return &class, nil
}

func (ClassDao) GetClassByName(name string) (*Class, error) {
	var class Class
	if err := database.DB.Model(&Class{}).First(&class).Error; err != nil {
		log.Printf("get class by name failed, classname:%s, err:%v", class.Name, err)
		return nil, err
	}
	return &class, nil
}

// CreateClass 创建班级
func (ClassDao) CreateClass(class *Class) error {
	if err := database.DB.Create(class).Error; err != nil {
		return err
	}
	return nil
}

// UpdateClass 更新班级信息
func (ClassDao) UpdateClass(id string, class *Class) error {
	if err := database.DB.Model(&Class{}).Where("id = ?", id).Updates(class).Error; err != nil {
		return err
	}
	return nil
}

// DeleteClass 删除班级
func (ClassDao) DeleteClass(id string) error {
	if err := database.DB.Where("id = ?", id).Delete(&Class{}).Error; err != nil {
		return err
	}
	return nil
}

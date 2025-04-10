package models

import (
	"gorm.io/gorm"
	"log"
	"schedule/database"
	"sync"
)

// Class 定义了班级表的结构
type Class struct {
	gorm.Model            // 内嵌 gorm.Model
	ID             string `gorm:"primaryKey;type:varchar(20)"` // 班级编号
	Name           string `gorm:"type:varchar(100);not null"`  // 班级名称
	Academic       string `gorm:"type:varchar(100)"`           //学制
	Cultivation    string `gorm:"type:varchar(100)"`           //培养层次
	Type           string `gorm:"type:varchar(100)"`           //班级类别
	ExpectedYear   string `gorm:"type:varchar(100)"`           //预计毕业年度
	IsGraduation   string `gorm:"type:varchar(100)"`           //是否毕业
	StudentCount   int    `gorm:"type:int"`                    // 班级人数
	MaxCount       int    `gorm:"type:int"`                    //班级最大人数
	Year           string `gorm:"type:varchar(100)"`           //入学年份
	Department     string `gorm:"type:varchar(50)"`            // 所属院系
	MajorID        string `gorm:"type:varchar(50)"`            //专业编号
	Major          string `gorm:"type:varchar(50)"`            // 专业
	Campus         string `gorm:"type:varchar(50)"`            // 校区
	FixedClassroom string `gorm:"type:varchar(50)"`            //固定教室编号
	IsFixed        string `gorm:"type:varchar(50)"`            //是否固定教室
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
	if err := database.DB.Model(&Class{}).Where("name like ?", "%"+name+"%").First(&class).Error; err != nil {
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

func (ClassDao) QueryByPage(page int, pagesize int) ([]Class, error) {
	var classes []Class
	offset := (page - 1) * pagesize
	if err := database.DB.Model(&Class{}).Order("created_at DESC").Limit(pagesize).Offset(offset).Find(&classes).Error; err != nil {
		return nil, err
	}
	return classes, nil
}

func (ClassDao) CountTotal() (int64, error) {
	var total int64
	if err := database.DB.Model(&Class{}).Count(&total).Error; err != nil {
		return 0, err
	}
	return total, nil
}

func (ClassDao) SearchClass(s string) ([]Class, error) {
	var classes []Class
	if err := database.DB.Model(&Class{}).
		Where("name LIKE ? OR id LIKE ? OR department LIKE ? OR major LIKE ? OR campus LIKE ?",
			"%"+s+"%", "%"+s+"%", "%"+s+"%", "%"+s+"%", "%"+s+"%").
		Find(&classes).Error; err != nil {
		return nil, err
	}
	return classes, nil
}

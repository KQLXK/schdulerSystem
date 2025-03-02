package models

import "gorm.io/gorm"

// Class 定义了班级表的结构
type Class struct {
	gorm.Model          // 内嵌 gorm.Model
	ID           string `gorm:"primaryKey;type:varchar(20)"`   // 班级编号
	Name         string `gorm:"type:varchar(100);not null"`    // 班级名称
	Department   string `gorm:"type:varchar(50)"`              // 所属院系
	Major        string `gorm:"type:varchar(50)"`              // 专业
	Campus       string `gorm:"type:varchar(50)"`              // 校区
	StudentCount int    `gorm:"type:int"`                      // 班级人数
	Status       string `gorm:"type:varchar(20);default:'启用'"` // 状态
}

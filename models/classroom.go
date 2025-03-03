package models

import "gorm.io/gorm"

// Classroom 定义了教室表的结构
type Classroom struct {
	gorm.Model        // 内嵌 gorm.Model
	ID         string `gorm:"primaryKey;type:varchar(20)"`   // 教室编号
	Name       string `gorm:"type:varchar(100);not null"`    // 教室名称
	Campus     string `gorm:"type:varchar(50)"`              // 校区
	Building   string `gorm:"type:varchar(50)"`              // 教学楼
	Capacity   int    `gorm:"type:int"`                      // 容量
	Type       string `gorm:"type:varchar(50)"`              // 教室类型（普通教室、多媒体教室等）
	Status     string `gorm:"type:varchar(20);default:'启用'"` // 状态
}

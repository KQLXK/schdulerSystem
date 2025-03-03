package models

import "gorm.io/gorm"

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

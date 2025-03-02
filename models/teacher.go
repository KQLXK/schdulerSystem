package models

import "gorm.io/gorm"

// Teacher 定义了教师表的结构
type Teacher struct {
	gorm.Model        // 内嵌 gorm.Model，包含 ID、CreatedAt、UpdatedAt、DeletedAt 字段
	ID         string `gorm:"primaryKey;type:varchar(10)"`   // 教师工号
	Name       string `gorm:"type:varchar(50);not null"`     // 教师姓名
	Gender     string `gorm:"type:varchar(10)"`              // 性别
	Department string `gorm:"type:varchar(50)"`              // 所属院系
	IsExternal bool   `gorm:"default:false"`                 // 是否外聘
	Status     string `gorm:"type:varchar(20);default:'启用'"` // 状态
}

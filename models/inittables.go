package models

import (
	"log"
	"schedule/database"
)

func InitTables() error {
	// 自动迁移表结构
	if err := database.DB.AutoMigrate(
		&Teacher{},
		&Classroom{},
		&Course{},
		&Class{},
		&Schedule{},
	); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
		return err
	}

	// 打印成功迁移信息
	log.Println("Database tables migrated successfully!")
	return nil
}

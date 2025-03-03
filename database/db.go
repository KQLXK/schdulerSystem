package database

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"schedule/models"
)

// DB 是全局的数据库连接实例
var DB *gorm.DB

// InitDB 初始化数据库连接
func InitDB() {
	// 直接定义数据库连接信息
	dbHost := "localhost" // 数据库主机地址
	dbPort := "3306"      // 数据库端口
	dbUser := "root"      // 数据库用户名
	dbPassword := "admin" // 数据库密码
	dbName := "schedule"  // 数据库名称

	// 构建 DSN（Data Source Name）
	dsn := dbUser + ":" + dbPassword + "@tcp(" + dbHost + ":" + dbPort + ")/" + dbName + "?charset=utf8mb4&parseTime=True&loc=Local"

	// 连接数据库
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// 打印成功连接信息
	log.Println("Successfully connected to database!")

	// 自动迁移表结构
	if err := DB.AutoMigrate(
		&models.Teacher{},
		&models.Classroom{},
		&models.Course{},
		&models.Class{},
		&models.Schedule{},
	); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// 打印成功迁移信息
	log.Println("Database tables migrated successfully!")
}

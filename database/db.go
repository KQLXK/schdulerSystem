package database

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"schedule/commen/config"
)

// DB 是全局的数据库连接实例
var (
	Rdb *redis.Client
	DB  *gorm.DB
	Ctx = context.Background()
)

// InitDB 初始化数据库连接
func InitDB() {
	// 直接定义数据库连接信息
	dbHost := config.GetConfig().Database.Host         // 数据库主机地址
	dbPort := config.GetConfig().Database.Port         // 数据库端口
	dbUser := config.GetConfig().Database.User         // 数据库用户名
	dbPassword := config.GetConfig().Database.Password // 数据库密码
	dbName := config.GetConfig().Database.Name         // 数据库名称

	// 构建 DSN（Data Source Name）
	dsn := dbUser + ":" + dbPassword + "@tcp(" + dbHost + ":" + dbPort + ")/?charset=utf8mb4&parseTime=True&loc=Local"

	// 连接数据库
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	if err = DB.Exec("CREATE DATABASE IF NOT EXISTS " + dbName).Error; err != nil {
		log.Println("create database failed, err:", err)
		return
	}

	dsn = dbUser + ":" + dbPassword + "@tcp(" + dbHost + ":" + dbPort + ")/" + dbName + "?charset=utf8mb4&parseTime=True&loc=Local"
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Println("connect to database failed, err:", err)
		return
	}

	// 打印成功连接信息
	log.Println("Successfully connected to database!")
}

func InitRedis() {
	Rdb = redis.NewClient(&redis.Options{
		Addr:     config.GetConfig().Redis.Addr,
		Password: "",
		DB:       0,
		PoolSize: 20,
	})
	_, err := Rdb.Ping(Ctx).Result()
	if err != nil {
		log.Fatalf("could not connect to redis: %v", err)
	}
	fmt.Println("connect redis success")
}

package models

import (
	"gorm.io/gorm"
	"log"
	"schedule/database"
	"sync"
)

// Course 定义了课程表的结构
type Course struct {
	gorm.Model            // 内嵌 gorm.Model
	ID            string  `gorm:"primaryKey;type:varchar(20)"` // 课程编号
	Name          string  `gorm:"type:varchar(100);not null"`  // 课程名称
	Type          string  `gorm:"type:varchar(20)"`            // 课程类型（理论、实践、实验）
	Property      string  `gorm:"type:varchar(100)"`           // 课程属性
	Credit        float64 `gorm:"type:float"`                  // 学分
	Department    string  `gorm:"type:varchar(50)"`            // 开课院系
	TotalHours    int64   `gorm:"type:int"`                    // 总学时
	TheoryHours   int64   `gorm:"type:int"`                    //理论学时
	TestHours     int64   `gorm:"type:int"`                    //实验学时
	ComputerHours int64   `gorm:"type:int"`                    //上机学时
	PracticeHours int64   `gorm:"type:int"`                    //实践学时
	OtherHours    int64   `gorm:"type:int"`                    //其他学时
	WeeklyHours   int64   `gorm:"type:int"`                    //周学时
	PurePractice  bool    `gorm:"type:boolean"`                //是否纯实践
}

type CourseDao struct{}

var (
	CourseOnce sync.Once
	courseDao  *CourseDao
)

func NewCourseDao() *CourseDao {
	CourseOnce.Do(func() {
		courseDao = &CourseDao{}
	})
	return courseDao
}

// GetAllCourses 获取所有课程
func (CourseDao) GetAllCourses() ([]Course, error) {
	var courses []Course
	if err := database.DB.Find(&courses).Error; err != nil {
		log.Println("Failed to get all courses, err:", err)
		return nil, err
	}
	log.Println("get all courses sucess")
	return courses, nil
}

// GetCourseByID 根据ID获取课程
func (CourseDao) GetCourseByID(id string) (*Course, error) {
	var course Course
	if err := database.DB.Where("id = ?", id).First(&course).Error; err != nil {
		log.Println("Failed to get course by ID:", id, ", err:", err)
		return nil, err
	}
	log.Println("get course by id sucess")
	return &course, nil
}

func (CourseDao) GetCourseByName(name string) (*Course, error) {
	var course Course
	if err := database.DB.Where("name = ?", name).First(&course).Error; err != nil {
		return nil, err
	}
	return &course, nil
}

// CreateCourse 创建课程
func (CourseDao) CreateCourse(course *Course) error {
	if err := database.DB.Create(course).Error; err != nil {
		log.Println("Database create course failed, err:", err)
		return err
	}
	log.Println("Course created successfully, ID:", course.ID, "Name:", course.Name)
	return nil
}

// UpdateCourse 更新课程信息
func (CourseDao) UpdateCourse(id string, course *Course) error {
	if err := database.DB.Model(&Course{}).Where("id = ?", id).Updates(course).Error; err != nil {
		log.Println("Failed to update course with ID:", id, ", err:", err)
		return err
	}
	log.Println("Course updated successfully, ID:", id)
	return nil
}

// DeleteCourse 删除课程
func (CourseDao) DeleteCourse(id string) error {
	if err := database.DB.Where("id = ?", id).Delete(&Course{}).Error; err != nil {
		log.Println("Failed to delete course with ID:", id, ", err:", err)
		return err
	}
	log.Println("Course deleted successfully, ID:", id)
	return nil
}

func (CourseDao) QueryByPage(page int, pagesize int) ([]Course, error) {
	var courses []Course
	offset := (page - 1) * pagesize
	if err := database.DB.Model(&Course{}).Order("created_at DESC").Limit(pagesize).Offset(offset).Find(&courses).Error; err != nil {
		log.Println("query course by page failed, err:", err)
		return nil, err
	}
	log.Println("query course by page sucess")
	return courses, nil
}

func (CourseDao) CountTotal() (int, error) {
	var total int64
	if err := database.DB.Model(&Course{}).Order("created_at DESC").Count(&total).Error; err != nil {
		log.Println("count total course failed, err:", err)
		return -1, err
	}
	log.Println("count total courses sucess")
	return int(total), nil
}

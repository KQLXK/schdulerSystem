package main

import (
	"schedule/database"
	"schedule/models"
	"schedule/route"
)

func main() {

	// 初始化数据库
	database.InitDB()
	//渲染数据表
	models.InitTables()

	r := route.SetupRoute()

	r.Run(":8080")

	//timeSlots := []models.TimeSlot{
	//	{WeekNumbers: []int{1, 2, 3}, Weekday: 1, StartPeriod: 1, Duration: 2},
	//	{WeekNumbers: []int{4, 5}, Weekday: 2, StartPeriod: 3, Duration: 1},
	//}
	//
	//scheduleResult := models.ScheduleResult{
	//	ScheduleID:  1,
	//	CourseID:    "C001",
	//	CourseName:  "数学",
	//	ClassroomID: "A101",
	//	TeacherID:   "T001",
	//	TeacherName: "张老师",
	//	ClassIDs:    []string{"301", "303"},
	//	TimeSlots:   timeSlots,
	//}
	//
	//// 保存到数据库
	//if err := database.DB.Create(&scheduleResult).Error; err != nil {
	//	log.Fatal(err)
	//}
	//
	//fmt.Println("ScheduleResult created:", scheduleResult)

}

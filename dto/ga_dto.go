package dto

import (
	"schedule/models"
	"time"
)

// API请求参数
type ScheduleGARequest struct {
	ScheduleIDs     []int  `json:"schedule_ids"`
	AdditionalRules Config `json:"rules,omitempty"`
}

type Config struct {
	SemesterWeek               int  `json:"semester_week"`                 //本学期计划周次
	MaxPeriodsPerDay           int  `json:"max_periods_per_day"`           // 每天最大节次
	MorningPeriodEnd           int  `json:"morning_period_end"`            // 上午结束节次（如4）
	AfternoonStartPeriod       int  `json:"afternoon_start_period"`        // 下午开始节次（如5）
	NightStartPeriod           int  `json:"night_start_period"`            // 晚上开始节次（如7）
	SportsAfternoonOnly        bool `json:"sports_afternoon_only"`         //体育课是否只在下午
	LabAtNightOnly             bool `json:"lab_at_night_only"`             //实验课是否只在晚上
	NightClassesAllowed        bool `json:"night_classes_allowed"`         //晚上是否排课
	MultiSessionConsecutive    bool `json:"multi_session_consecutive"`     //多学时类型（理论、实验、上机）的课程学时是否连续排
	TeacherMaxDailyPeriods     int  `json:"teacher_max_daily_periods"`     // 教师每天最大节次
	TeacherMaxWeeklyPeriods    int  `json:"teacher_max_weekly_periods"`    //教师每周最大节次
	TeacherMaxMorningPeriods   int  `json:"teacher_max_morning_periods"`   //教师上午最大节次
	TeacherMaxAfternoonPeriods int  `json:"teacher_max_afternoon_periods"` //教师下午最大节次
}

// 特殊规则配置
//type RuleConfig struct {
//	PriorityCourses  []int `json:"priority_courses"` // 优先课程ID列表
//	ForbiddenPeriods []struct {
//		Weekday int   `json:"weekday"` // 1-7
//		Periods []int `json:"periods"` // 1-12
//	} `json:"forbidden_periods"`
//	ClassroomPreferences map[int]string `json:"classroom_prefs"` // 课程ID -> 偏好教室类型
//}

// API响应结构
type ScheduleGAResponse struct {
	Status      string         `json:"status"`             // processing/success/failed
	Progress    float64        `json:"progress,omitempty"` // 0-1
	GeneratedAt time.Time      `json:"generated_at"`
	SuccessList []SuccessItem  `json:"success_list"`
	FailedList  []FailedItem   `json:"failed_list"`
	Analysis    ResultAnalysis `json:"analysis"`
}

// 成功条目
type SuccessItem struct {
	ScheduleID    int               `json:"schedule_id"`
	CourseID      string            `json:"course_id"`
	CourseName    string            `json:"course_name"`
	TimeSlots     []models.TimeSlot `json:"time_slots"`
	ClassroomID   string            `json:"classroom_id"`
	ClassroomName string            `json:"classroom_name"`
	TeacherID     string            `json:"teacher_id"`
	TeacherName   string            `json:"teacher_name"`
}

// 失败条目
type FailedItem struct {
	ScheduleID       int               `json:"schedule_id"`
	CourseName       string            `json:"course_name"`
	ClassroomID      string            `json:"classroom_id"`
	ClassroomName    string            `json:"classroom_name"`
	TeacherID        string            `json:"teacher_id"`
	TeacherName      string            `json:"teacher_name"`
	TimeSlots        []models.TimeSlot `json:"time_slots"`
	ConflictReasons  []string          `json:"reasons"`
	SuggestedTimes   []models.TimeSlot `json:"suggested_times,omitempty"`
	AlternativeRooms []string          `json:"alternative_rooms,omitempty"`
}

// 时间槽DTO
type SlotDTO struct {
	Weeks       []int `json:"weeks"`        // 周次列表
	Weekday     int   `json:"weekday"`      // 1-7
	StartPeriod int   `json:"start_period"` // 开始节次
	Duration    int   `json:"duration"`     // 持续节数
}

// 结果分析
type ResultAnalysis struct {
	ClassroomUtilization map[string]float64 `json:"classroom_usage"`   // 教室ID -> 使用率
	TeacherWorkload      map[string]int     `json:"teacher_workload"`  // 教师ID -> 总课时
	TimeDistribution     map[string]int     `json:"time_distribution"` // 时间段 -> 课程数
}

// 手动排课请求结构
type ManualScheduleRequest struct {
	//Semester    string            `json:"semester" binding:"required"`
	ScheduleID int `json:"schedule_id" binding:"required"`
	//CourseID    string            `json:"course_id" binding:"required"`
	ClassroomID string `json:"classroom_id" binding:"required"`
	//TeacherID   string            `json:"teacher_id" binding:"required"`
	//ClassIDs    []string          `json:"class_ids" binding:"required"`
	TimeSlots []models.TimeSlot `json:"time_slots" binding:"required"`
}

// 排课结果响应结构
type ScheduleResultResponse struct {
	ID          uint              `json:"id"`
	Semester    string            `json:"semester"`
	CourseID    string            `json:"course_id"`
	CourseName  string            `json:"course_name"`
	ClassroomID string            `json:"classroom_id"`
	TeacherID   string            `json:"teacher_id"`
	TeacherName string            `json:"teacher_name"`
	ClassIDs    []string          `json:"class_ids"`
	TimeSlots   []models.TimeSlot `json:"time_slots"`
}

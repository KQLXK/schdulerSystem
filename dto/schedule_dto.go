package dto

import "time"

// ScheduleCreateReq 创建排课任务的请求
type ScheduleCreateReq struct {
	Semester               string `json:"semester" validate:"required"`         // 学年学期（格式：年份-年份-学期）
	CourseID               string `json:"course_id" validate:"required"`        // 课程编号
	CourseName             string `json:"course_name" validate:"required"`      // 课程名称
	TeacherID              string `json:"teacher_id" validate:"required"`       // 教师工号
	TeachingClass          string `json:"teaching_class"`                       // 教学班组成
	TeachingClassID        string `json:"teaching_class_number"`                // 教学班编号
	TeachingClassName      string `json:"teaching_class_name"`                  // 教学班名称
	HourType               string `json:"hour_type"`                            // 学时类型
	OpeningHours           int64  `json:"opening_hours"`                        // 开课学时
	SchedulingHours        int64  `json:"scheduling_hours"`                     // 排课学时
	TotalHours             int64  `json:"total_hours" validate:"required"`      // 总学时
	SchedulingPriority     int64  `json:"scheduling_priority"`                  // 排课优先级
	TeachingClassSize      int64  `json:"teaching_class_size" validate:"min=0"` // 教学班人数
	OpeningCampus          string `json:"opening_campus" validate:"required"`   // 开课校区
	OpeningWeekHours       string `json:"opening_week_hours"`                   // 开课周次学时
	ContinuousPeriods      int64  `json:"continuous_periods" validate:"min=1"`  // 连排节次
	SpecifiedClassroomType string `json:"specified_classroom_type"`             // 指定教室类型
	SpecifiedClassroom     string `json:"specified_classroom"`                  // 指定教室
	SpecifiedBuilding      string `json:"specified_building"`                   // 指定教学楼
	SpecifiedTime          string `json:"specified_time"`                       // 指定时间（格式需约定，如：Mon-08:00）
}

// ScheduleCreateResp 创建排课任务的响应
type ScheduleCreateResp struct {
	ScheduleID  int64  `json:"schedule_id"` // 排课任务ID
	CourseName  string `json:"course_name"`
	TeacherName string `json:"teacher_name"`
	Semester    string `json:"semester"`      // 回显关键信息
	Err         string `json:"err,omitempty"` // 错误信息（成功时为空）
}

// ScheduleUpdateReq 更新排课任务的请求
type ScheduleUpdateReq struct {
	ID                     int64  `json:"id" validate:"required"` // 排课任务ID
	Semester               string `json:"semester"`
	CourseID               string `json:"course_id"`
	CourseName             string `json:"course_name"`
	TeacherID              string `json:"teacher_id"`
	TeachingClass          string `json:"teaching_class"`
	TeachingClassID        string `json:"teaching_class_number"`
	TeachingClassName      string `json:"teaching_class_name"`
	HourType               string `json:"hour_type"`
	OpeningHours           int64  `json:"opening_hours"`
	SchedulingHours        int64  `json:"scheduling_hours"`
	TotalHours             int64  `json:"total_hours"`
	SchedulingPriority     int64  `json:"scheduling_priority"`
	TeachingClassSize      int64  `json:"teaching_class_size" validate:"min=0"`
	OpeningCampus          string `json:"opening_campus"`
	OpeningWeekHours       string `json:"opening_week_hours"`
	ContinuousPeriods      int64  `json:"continuous_periods" validate:"min=1"`
	SpecifiedClassroomType string `json:"specified_classroom_type"`
	SpecifiedClassroom     string `json:"specified_classroom"`
	SpecifiedBuilding      string `json:"specified_building"`
	SpecifiedTime          string `json:"specified_time"`
}

// ScheduleGetResp 获取排课详情的响应
type ScheduleGetResp struct {
	ID                     int64     `json:"id"`
	CreatedAt              time.Time `json:"created_at"`
	UpdatedAt              time.Time `json:"updated_at"`
	Semester               string    `json:"semester"`
	CourseID               string    `json:"course_id"`
	CourseName             string    `json:"course_name"`
	TeacherID              string    `json:"teacher_id"`
	TeachingClass          string    `json:"teaching_class"`
	TeachingClassID        string    `json:"teaching_class_number"`
	TeachingClassName      string    `json:"teaching_class_name"`
	HourType               string    `json:"hour_type"`
	OpeningHours           int64     `json:"opening_hours"`
	SchedulingHours        int64     `json:"scheduling_hours"`
	TotalHours             int64     `json:"total_hours"`
	SchedulingPriority     int64     `json:"scheduling_priority"`
	TeachingClassSize      int64     `json:"teaching_class_size"`
	OpeningCampus          string    `json:"opening_campus"`
	OpeningWeekHours       string    `json:"opening_week_hours"`
	ContinuousPeriods      int64     `json:"continuous_periods"`
	SpecifiedClassroomType string    `json:"specified_classroom_type"`
	SpecifiedClassroom     string    `json:"specified_classroom"`
	SpecifiedBuilding      string    `json:"specified_building"`
	SpecifiedTime          string    `json:"specified_time"`
}

// ScheduleQueryByPageResp 分页查询响应
type ScheduleQueryByPageResp struct {
	Total      int64             `json:"total"`      // 总记录数
	TotalPages int64             `json:"totalPages"` // 总页数
	Page       int64             `json:"page"`       // 当前页码
	PageSize   int64             `json:"pageSize"`   // 每页数量
	Schedules  []ScheduleGetResp `json:"schedules"`  // 排课列表
}

// ScheduleSearchResp 搜索排课响应
type ScheduleSearchResp struct {
	TotalCount int64             `json:"total_count"` // 匹配总数
	Schedules  []ScheduleGetResp `json:"schedules"`   // 排课列表
}

// ScheduleAddByExcelResp Excel批量导入响应
type ScheduleAddByExcelResp struct {
	AddSuccess int64                 `json:"add_success"` // 成功数量
	AddFail    int64                 `json:"add_fail"`    // 失败数量
	FailList   []*ScheduleCreateResp `json:"fail_list"`   // 失败明细
}

type ScheduleQueryAllResp struct {
	Total     int64
	Schedules []ScheduleGetResp
}

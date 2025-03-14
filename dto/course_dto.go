package dto

type CourseCreateReq struct {
	// 上机学时
	ComputerHours int64 `json:"computer_hours"`
	// 课程学分
	CourseCredit float64 `json:"course_credit"`
	// 开课院系
	CourseDepartment string `json:"course_department"`
	// 课程编号
	CourseID string `json:"course_id"`
	// 课程名称
	CourseName string `json:"course_name"`
	// 课程属性
	CourseProperty string `json:"course_property"`
	// 课程类型
	CourseType string `json:"course_type"`
	// 其他学时
	OtherHours int64 `json:"other_hours"`
	// 实践学时
	PracticeHours int64 `json:"practice_hours"`
	// 是否纯实践
	PurePractice bool `json:"pure_practice"`
	// 实验学时
	TestHours int64 `json:"test_hours"`
	// 理论学时
	TheoryHours int64 `json:"theory_hours"`
	// 总学时
	TotalHour int64 `json:"total_hours"`
	// 周学时
	WeeklyHours int64 `json:"weekly_hours"`
}

type CourseCreateResp struct {
	CourseID   string `json:"course_id"`
	CourseName string `json:"course_name"`
}

type CourseUpdateReq struct {
	CourseID         string  `json:"course_id" validate:"required"`   // 课程编号
	CourseName       string  `json:"course_name" validate:"required"` // 课程名称
	CourseType       string  `json:"course_type"`                     // 课程类型
	CourseProperty   string  `json:"course_property"`                 // 课程属性
	CourseCredit     float64 `json:"course_credit"`                   // 学分
	CourseDepartment string  `json:"course_department"`               // 开课院系
	TotalHour        int64   `json:"total_hour" validate:"required"`  // 总学时
	TheoryHours      int64   `json:"theory_hours"`                    // 理论学时
	TestHours        int64   `json:"test_hours"`                      // 实验学时
	ComputerHours    int64   `json:"computer_hours"`                  // 上机学时
	PracticeHours    int64   `json:"practice_hours"`                  // 实践学时
	OtherHours       int64   `json:"other_hours"`                     // 其他学时
	WeeklyHours      int64   `json:"weekly_hours"`                    // 周学时
	PurePractice     bool    `json:"pure_practice"`                   // 是否纯实践
}

type CourseGetResp struct {
	CourseID         string  `json:"course_id"`         // 课程编号
	CourseName       string  `json:"course_name"`       // 课程名称
	CourseType       string  `json:"course_type"`       // 课程类型
	CourseProperty   string  `json:"course_property"`   // 课程属性
	CourseCredit     float64 `json:"course_credit"`     // 学分
	CourseDepartment string  `json:"course_department"` // 开课院系
	TotalHour        int     `json:"total_hour"`        // 总学时
	TheoryHours      int     `json:"theory_hours"`      // 理论学时
	TestHours        int     `json:"test_hours"`        // 实验学时
	ComputerHours    int     `json:"computer_hours"`    // 上机学时
	PracticeHours    int     `json:"practice_hours"`    // 实践学时
	OtherHours       int     `json:"other_hours"`       // 其他学时
	WeeklyHours      int     `json:"weekly_hours"`      // 周学时
	PurePractice     bool    `json:"pure_practice"`     // 是否纯实践
}

type CourseGetAllResp struct {
	Courses []CourseGetResp `json:"courses"` // 课程列表
}

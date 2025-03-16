package dto

// TeacherCreateReq 定义了创建教师的请求结构
type TeacherCreateReq struct {
	ID         string `json:"id" binding:"required"`   // 教师工号
	Name       string `json:"name" binding:"required"` // 教师姓名
	Gender     string `json:"gender"`                  // 性别
	Department string `json:"department"`              // 所属院系
	IsExternal bool   `json:"is_external"`             // 是否外聘
	Status     string `json:"status"`                  // 状态
}

// TeacherUpdateReq 定义了更新教师的请求结构
type TeacherUpdateReq struct {
	ID         string `json:"id" binding:"required"`   // 教师工号
	Name       string `json:"name" binding:"required"` // 教师姓名
	Gender     string `json:"gender"`                  // 性别
	Department string `json:"department"`              // 所属院系
	IsExternal bool   `json:"is_external"`             // 是否外聘
	Status     string `json:"status"`                  // 状态
}

// TeacherGetResp 定义了获取教师的响应结构
type TeacherGetResp struct {
	ID         string `json:"id"`          // 教师工号
	Name       string `json:"name"`        // 教师姓名
	Gender     string `json:"gender"`      // 性别
	Department string `json:"department"`  // 所属院系
	IsExternal bool   `json:"is_external"` // 是否外聘
	Status     string `json:"status"`      // 状态
}

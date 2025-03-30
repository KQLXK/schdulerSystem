package dto

// TeacherCreateReq 定义了创建教师的请求结构
type TeacherCreateReq struct {
	ID          string `json:"id" binding:"required"`   // 教师工号
	Name        string `json:"name" binding:"required"` // 教师姓名
	EnglishName string `json:"english_name"`            // 英文姓名
	Gender      string `json:"gender"`                  // 性别
	Ethnicity   string `json:"ethnicity"`               // 民族
	Department  string `json:"department"`              // 所属院系
	Title       string `json:"title"`                   // 职称
	Category    string `json:"category"`                // 教职工类别
	IsExternal  bool   `json:"is_external"`             // 是否外聘
	Status      string `json:"status"`                  // 状态
}

type TeacherCreateResp struct {
	ID   string `json:"id"`            // 教师工号
	Name string `json:"name"`          // 教师姓名
	Err  string `json:"err,omitempty"` // 错误信息(可选)
}

// TeacherUpdateReq 定义了更新教师的请求结构
type TeacherUpdateReq struct {
	ID          string `json:"id" binding:"required"`   // 教师工号
	Name        string `json:"name" binding:"required"` // 教师姓名
	EnglishName string `json:"english_name"`            // 英文姓名
	Gender      string `json:"gender"`                  // 性别
	Ethnicity   string `json:"ethnicity"`               // 民族
	Department  string `json:"department"`              // 所属院系
	Title       string `json:"title"`                   // 职称
	Category    string `json:"category"`                // 教职工类别
	IsExternal  bool   `json:"is_external"`             // 是否外聘
	Status      string `json:"status"`                  // 状态
}

// TeacherGetResp 定义了获取教师的响应结构
type TeacherGetResp struct {
	ID          string `json:"id"`           // 教师工号
	Name        string `json:"name"`         // 教师姓名
	EnglishName string `json:"english_name"` // 英文姓名
	Gender      string `json:"gender"`       // 性别
	Ethnicity   string `json:"ethnicity"`    // 民族
	Department  string `json:"department"`   // 所属院系
	Title       string `json:"title"`        // 职称
	Category    string `json:"category"`     // 教职工类别
	IsExternal  bool   `json:"is_external"`  // 是否外聘
	Status      string `json:"status"`       // 状态
}

type TeacherQueryByPageResp struct {
	Total      int64            `json:"total"`      // 总记录数
	TotalPages int64            `json:"totalPages"` // 总页数
	Page       int              `json:"page"`       // 当前页码
	PageSize   int              `json:"pageSize"`   // 每页数量
	Teachers   []TeacherGetResp `json:"teachers"`   // 教师列表
}

type TeacherSearchResp struct {
	TotalCount int64            `json:"total_count"`
	Teachers   []TeacherGetResp `json:"teachers"`
}

type TeacherAddByExcelResp struct {
	AddSuccess int                  `json:"add_success"`
	AddFail    int                  `json:"add_fail"`
	FailList   []*TeacherCreateResp `json:"fail_list"`
}

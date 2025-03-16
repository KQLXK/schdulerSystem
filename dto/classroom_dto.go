package dto

// ClassroomCreateReq 定义了创建教室的请求结构
type ClassroomCreateReq struct {
	ID          string `json:"id" binding:"required"`       // 教室编号
	Name        string `json:"name" binding:"required"`     // 教室名称
	Campus      string `json:"campus"`                      // 校区
	Building    string `json:"building"`                    // 教学楼
	Floor       string `json:"floor"`                       // 所在楼层
	Capacity    int    `json:"capacity" binding:"required"` // 容量
	Type        string `json:"type"`                        // 教室类型（普通教室、多媒体教室等）
	HasAC       bool   `json:"has_ac"`                      // 是否有空调
	Description string `json:"description"`                 // 教室描述
	Department  string `json:"department"`                  // 管理部门
	Status      string `json:"status"`                      // 状态
}

// ClassroomUpdateReq 定义了更新教室的请求结构
type ClassroomUpdateReq struct {
	ID          string `json:"id" binding:"required"`       // 教室编号
	Name        string `json:"name" binding:"required"`     // 教室名称
	Campus      string `json:"campus"`                      // 校区
	Building    string `json:"building"`                    // 教学楼
	Floor       string `json:"floor"`                       // 所在楼层
	Capacity    int    `json:"capacity" binding:"required"` // 容量
	Type        string `json:"type"`                        // 教室类型（普通教室、多媒体教室等）
	HasAC       bool   `json:"has_ac"`                      // 是否有空调
	Description string `json:"description"`                 // 教室描述
	Department  string `json:"department"`                  // 管理部门
	Status      string `json:"status"`                      // 状态
}

// ClassroomGetResp 定义了获取教室的响应结构
type ClassroomGetResp struct {
	ID          string `json:"id"`          // 教室编号
	Name        string `json:"name"`        // 教室名称
	Campus      string `json:"campus"`      // 校区
	Building    string `json:"building"`    // 教学楼
	Floor       string `json:"floor"`       // 所在楼层
	Capacity    int    `json:"capacity"`    // 容量
	Type        string `json:"type"`        // 教室类型
	HasAC       bool   `json:"has_ac"`      // 是否有空调
	Description string `json:"description"` // 教室描述
	Department  string `json:"department"`  // 管理部门
	Status      string `json:"status"`      // 状态
}

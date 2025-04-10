package dto

import "time"

type ClassCreateReq struct {
	ID             string `json:"id" binding:"required"`
	Name           string `json:"name" binding:"required"`
	Academic       string `json:"academic"`
	Cultivation    string `json:"cultivation"`
	Type           string `json:"type"`
	ExpectedYear   string `json:"expected_year"`
	IsGraduation   string `json:"is_graduation"`
	StudentCount   int    `json:"student_count"`
	MaxCount       int    `json:"max_count"`
	Year           string `json:"year"`
	Department     string `json:"department"`
	MajorID        string `json:"major_id"`
	Major          string `json:"major"`
	Campus         string `json:"campus"`
	FixedClassroom string `json:"fixed_classroom"`
	IsFixed        string `json:"is_fixed"`
}

type ClassCreateResp struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Err  string `json:"err,omitempty"`
}

type ClassUpdateReq struct {
	ID             string `json:"id" binding:"required"`
	Name           string `json:"name" binding:"required"`
	Academic       string `json:"academic"`
	Cultivation    string `json:"cultivation"`
	Type           string `json:"type"`
	ExpectedYear   string `json:"expected_year"`
	IsGraduation   string `json:"is_graduation"`
	StudentCount   int    `json:"student_count"`
	MaxCount       int    `json:"max_count"`
	Year           string `json:"year"`
	Department     string `json:"department"`
	MajorID        string `json:"major_id"`
	Major          string `json:"major"`
	Campus         string `json:"campus"`
	FixedClassroom string `json:"fixed_classroom"`
	IsFixed        string `json:"is_fixed"`
}

type ClassGetResp struct {
	ID             string    `json:"id"`
	Name           string    `json:"name"`
	Academic       string    `json:"academic"`
	Cultivation    string    `json:"cultivation"`
	Type           string    `json:"type"`
	ExpectedYear   string    `json:"expected_year"`
	IsGraduation   string    `json:"is_graduation"`
	StudentCount   int       `json:"student_count"`
	MaxCount       int       `json:"max_count"`
	Year           string    `json:"year"`
	Department     string    `json:"department"`
	MajorID        string    `json:"major_id"`
	Major          string    `json:"major"`
	Campus         string    `json:"campus"`
	FixedClassroom string    `json:"fixed_classroom"`
	IsFixed        string    `json:"is_fixed"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type ClassQueryByPageResp struct {
	Total      int64          `json:"total"`
	TotalPages int64          `json:"totalPages"`
	Page       int            `json:"page"`
	PageSize   int            `json:"pageSize"`
	Classes    []ClassGetResp `json:"classes"`
}

type ClassSearchResp struct {
	TotalCount int64          `json:"total_count"`
	Classes    []ClassGetResp `json:"classes"`
}

type ClassAddByExcelResp struct {
	AddSuccess int                `json:"add_success"`
	AddFail    int                `json:"add_fail"`
	FailList   []*ClassCreateResp `json:"fail_list"`
}

package teacher

import (
	"errors"
	"fmt"
	"log"
	"schedule/commen/utils"
	"schedule/dto"
	"schedule/models"
	"strings"
)

var (
	FileFormatErr        = errors.New("文件格式有误")
	MissingColumnErr     = errors.New("缺少必要的列")
	InvalidDataFormatErr = errors.New("数据格式错误")
)

var teacherFieldToColumn = map[string]string{
	"ID":          "工号",
	"Name":        "姓名",
	"Gender":      "性别",
	"EnglishName": "英文姓名",
	"Ethnicity":   "民族",
	"Title":       "职称",
	"Department":  "单位",
	"IsExternal":  "是否外聘",
	"Category":    "教职工类别",
}

type TeacherAddByExcelFlow struct {
	Filename    string
	TeacherList [][]string
}

func TeacherAddByExcel(filename string) (*dto.TeacherAddByExcelResp, error) {
	return NewTeacherAddByExcelFlow(filename).Do()
}

func NewTeacherAddByExcelFlow(filename string) *TeacherAddByExcelFlow {
	return &TeacherAddByExcelFlow{
		Filename: filename,
	}
}

func (f *TeacherAddByExcelFlow) Do() (*dto.TeacherAddByExcelResp, error) {
	teacherList, err := f.ReadFile()
	if err != nil {
		return nil, err
	}
	f.TeacherList = teacherList

	if len(f.TeacherList) < 2 {
		return nil, errors.New("文件内容为空")
	}
	headers := f.TeacherList[1]
	dataRows := f.TeacherList[2:]

	columnMap := make(map[string]int)
	for i, header := range headers {
		columnMap[header] = i
	}

	fieldMap := make(map[string]int)
	for field, colName := range teacherFieldToColumn {
		idx, ok := columnMap[colName]
		if !ok {
			err = fmt.Errorf("%w: %s", MissingColumnErr, colName)
			log.Println(err)
			return nil, err
		}
		fieldMap[field] = idx
	}

	return f.AddTeacher(fieldMap, dataRows), nil
}

func (f *TeacherAddByExcelFlow) ReadFile() ([][]string, error) {
	res, err := utils.ReadExcel(f.Filename)
	if err != nil {
		if err == utils.ExcelFormatErr {
			return nil, FileFormatErr
		}
		return nil, err
	}
	return res, nil
}

func (f *TeacherAddByExcelFlow) AddTeacher(fieldMap map[string]int, dataRows [][]string) *dto.TeacherAddByExcelResp {
	faillist := make([]*dto.TeacherCreateResp, 0)
	var successcount, failcount int

	for _, row := range dataRows {
		teacher, err := parseTeacher(row, fieldMap)
		if err != nil {
			failcount++
			faillist = append(faillist, &dto.TeacherCreateResp{
				ID:   getValueSafely(row, fieldMap["ID"]),
				Name: getValueSafely(row, fieldMap["Name"]),
				Err:  err.Error(),
			})
			continue
		}

		if err = models.CreateTeacher(teacher); err != nil {
			failcount++
			faillist = append(faillist, &dto.TeacherCreateResp{
				ID:   teacher.ID,
				Name: teacher.Name,
				Err:  err.Error(),
			})
		} else {
			successcount++
		}
	}

	return &dto.TeacherAddByExcelResp{
		AddSuccess: successcount,
		AddFail:    failcount,
		FailList:   faillist,
	}
}

func parseTeacher(row []string, fieldMap map[string]int) (*models.Teacher, error) {
	teacher := &models.Teacher{
		ID:          getValueSafely(row, fieldMap["ID"]),
		Name:        getValueSafely(row, fieldMap["Name"]),
		EnglishName: getValueSafely(row, fieldMap["EnglishName"]),
		Gender:      getValueSafely(row, fieldMap["Gender"]),
		Ethnicity:   getValueSafely(row, fieldMap["Ethnicity"]),
		Title:       getValueSafely(row, fieldMap["Title"]),
		Department:  getValueSafely(row, fieldMap["Department"]),
		Category:    getValueSafely(row, fieldMap["Category"]),
		Status:      "启用", // 默认状态
	}

	// 处理是否外聘字段
	isExternalStr := getValueSafely(row, fieldMap["IsExternal"])
	if isExternalStr == "是" {
		teacher.IsExternal = true
	} else {
		teacher.IsExternal = false
	}

	if teacher.ID == "" || teacher.Name == "" {
		return nil, fmt.Errorf("教师工号或姓名不能为空")
	}

	return teacher, nil
}

func getValueSafely(row []string, index int) string {
	if index < 0 || index >= len(row) {
		return ""
	}
	return strings.TrimSpace(row[index])
}

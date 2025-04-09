package classroom

import (
	"errors"
	"fmt"
	"log"
	"schedule/commen/utils"
	"schedule/dto"
	"schedule/models"
	"strconv"
	"strings"
)

var (
	FileFormatErr        = errors.New("文件格式有误")
	MissingColumnErr     = errors.New("缺少必要的列")
	InvalidDataFormatErr = errors.New("数据格式错误")
)

var classroomFieldToColumn = map[string]string{
	"ID":          "教室编号",
	"Name":        "教室名称",
	"Campus":      "校区",
	"Building":    "教学楼",
	"Floor":       "所在楼层",
	"Type":        "教室类型",
	"Capacity":    "最大上课容纳人数",
	"HasAC":       "是否有空调",
	"Status":      "是否启用",
	"Description": "教室描述",
	"Department":  "管理部门",
}

type ClassroomAddByExcelFlow struct {
	Filename      string
	ClassroomList [][]string
}

func ClassroomAddByExcel(filename string) (*dto.ClassroomAddByExcelResp, error) {
	return NewClassroomAddByExcelFlow(filename).Do()
}

func NewClassroomAddByExcelFlow(filename string) *ClassroomAddByExcelFlow {
	return &ClassroomAddByExcelFlow{
		Filename: filename,
	}
}

func (f *ClassroomAddByExcelFlow) Do() (*dto.ClassroomAddByExcelResp, error) {
	classroomList, err := f.ReadFile()
	if err != nil {
		return nil, err
	}
	f.ClassroomList = classroomList

	if len(f.ClassroomList) < 1 {
		return nil, errors.New("文件内容为空")
	}
	headers := f.ClassroomList[1]
	dataRows := f.ClassroomList[2:]

	columnMap := make(map[string]int)
	for i, header := range headers {
		columnMap[header] = i
	}

	fieldMap := make(map[string]int)
	for field, colName := range classroomFieldToColumn {
		idx, ok := columnMap[colName]
		if !ok {
			err = fmt.Errorf("%w: %s", MissingColumnErr, colName)
			log.Println(err)
			return nil, err
		}
		fieldMap[field] = idx
	}

	return f.AddClassroom(fieldMap, dataRows), nil
}

func (f *ClassroomAddByExcelFlow) ReadFile() ([][]string, error) {
	res, err := utils.ReadExcel(f.Filename)
	if err != nil {
		if err == utils.ExcelFormatErr {
			return nil, FileFormatErr
		}
		return nil, err
	}
	return res, nil
}

func (f *ClassroomAddByExcelFlow) AddClassroom(fieldMap map[string]int, dataRows [][]string) *dto.ClassroomAddByExcelResp {
	faillist := make([]*dto.ClassroomCreateResp, 0)
	var successcount, failcount int

	for _, row := range dataRows {
		classroom, err := parseClassroom(row, fieldMap)
		if err != nil {
			failcount++
			faillist = append(faillist, &dto.ClassroomCreateResp{
				ID:   getValueSafely(row, fieldMap["ID"]),
				Name: getValueSafely(row, fieldMap["Name"]),
				Err:  err.Error(),
			})
			continue
		}

		if err = models.NewClassroomDao().CreateClassroom(classroom); err != nil {
			failcount++
			faillist = append(faillist, &dto.ClassroomCreateResp{
				ID:   classroom.ID,
				Name: classroom.Name,
				Err:  err.Error(),
			})
		} else {
			successcount++
		}
	}

	return &dto.ClassroomAddByExcelResp{
		AddSuccess: successcount,
		AddFail:    failcount,
		FailList:   faillist,
	}
}

func parseClassroom(row []string, fieldMap map[string]int) (*models.Classroom, error) {
	var err error
	classroom := &models.Classroom{
		ID:          getValueSafely(row, fieldMap["ID"]),
		Name:        getValueSafely(row, fieldMap["Name"]),
		Campus:      getValueSafely(row, fieldMap["Campus"]),
		Building:    getValueSafely(row, fieldMap["Building"]),
		Floor:       getValueSafely(row, fieldMap["Floor"]),
		Type:        getValueSafely(row, fieldMap["Type"]),
		Description: getValueSafely(row, fieldMap["Description"]),
		Department:  getValueSafely(row, fieldMap["Department"]),
	}

	// 处理容量
	capacityStr := getValueSafely(row, fieldMap["Capacity"])
	if classroom.Capacity, err = strconv.Atoi(capacityStr); err != nil {
		return nil, fmt.Errorf("%w: 容量", InvalidDataFormatErr)
	}

	// 处理是否有空调
	hasACStr := getValueSafely(row, fieldMap["HasAC"])
	if hasACStr == "是" {
		classroom.HasAC = true
	} else {
		classroom.HasAC = false
	}

	// 处理是否启用
	statusStr := getValueSafely(row, fieldMap["Status"])
	if statusStr == "是" {
		classroom.Status = "启用"
	} else {
		classroom.Status = "停用"
	}

	if classroom.ID == "" || classroom.Name == "" {
		return nil, fmt.Errorf("教室编号或名称不能为空")
	}

	return classroom, nil
}

func getValueSafely(row []string, index int) string {
	if index < 0 || index >= len(row) {
		return ""
	}
	return strings.TrimSpace(row[index])
}

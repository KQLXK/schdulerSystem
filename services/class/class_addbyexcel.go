package class

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

var classFieldToColumn = map[string]string{
	"ID":             "班级编号",
	"Name":           "班级名称",
	"Academic":       "学制",
	"Cultivation":    "培养层次",
	"Type":           "班级类别",
	"ExpectedYear":   "预计毕业年度",
	"IsGraduation":   "是否毕业",
	"StudentCount":   "班级人数",
	"MaxCount":       "班级最大人数",
	"Year":           "入学年份",
	"Department":     "所属院系",
	"MajorID":        "专业编号",
	"Major":          "专业",
	"Campus":         "校区",
	"FixedClassroom": "固定教室", // 只映射第一个"固定教室"列
}

type ClassAddByExcelFlow struct {
	Filename  string
	ClassList [][]string
}

func ClassAddByExcel(filename string) (*dto.ClassAddByExcelResp, error) {
	return NewClassAddByExcelFlow(filename).Do()
}

func NewClassAddByExcelFlow(filename string) *ClassAddByExcelFlow {
	return &ClassAddByExcelFlow{
		Filename: filename,
	}
}

func (f *ClassAddByExcelFlow) Do() (*dto.ClassAddByExcelResp, error) {
	classList, err := f.ReadFile()
	if err != nil {
		return nil, err
	}
	f.ClassList = classList

	if len(f.ClassList) < 2 {
		return nil, errors.New("文件内容为空")
	}
	headers := f.ClassList[1]
	dataRows := f.ClassList[2:]

	// 创建列名到索引的映射，确保只取第一个"固定教室"
	columnMap := make(map[string]int)
	seenColumns := make(map[string]bool) // 用于跟踪已处理的列名

	for i, header := range headers {
		trimmedHeader := strings.TrimSpace(header)
		if !seenColumns[trimmedHeader] {
			columnMap[trimmedHeader] = i
			seenColumns[trimmedHeader] = true
		}
	}

	fieldMap := make(map[string]int)
	for field, colName := range classFieldToColumn {
		idx, ok := columnMap[colName]
		if !ok {
			err = fmt.Errorf("%w: %s", MissingColumnErr, colName)
			log.Println(err)
			return nil, err
		}
		fieldMap[field] = idx
	}

	return f.AddClass(fieldMap, dataRows), nil
}

func (f *ClassAddByExcelFlow) ReadFile() ([][]string, error) {
	res, err := utils.ReadExcel(f.Filename)
	if err != nil {
		if err == utils.ExcelFormatErr {
			return nil, FileFormatErr
		}
		return nil, err
	}
	return res, nil
}

func (f *ClassAddByExcelFlow) AddClass(fieldMap map[string]int, dataRows [][]string) *dto.ClassAddByExcelResp {
	faillist := make([]*dto.ClassCreateResp, 0)
	var successcount, failcount int

	for _, row := range dataRows {
		class, err := parseClass(row, fieldMap)
		if err != nil {
			failcount++
			faillist = append(faillist, &dto.ClassCreateResp{
				ID:   getValueSafely(row, fieldMap["ID"]),
				Name: getValueSafely(row, fieldMap["Name"]),
				Err:  err.Error(),
			})
			continue
		}

		if err = models.NewClassDaoInstance().CreateClass(class); err != nil {
			failcount++
			faillist = append(faillist, &dto.ClassCreateResp{
				ID:   class.ID,
				Name: class.Name,
				Err:  err.Error(),
			})
		} else {
			successcount++
		}
	}

	return &dto.ClassAddByExcelResp{
		AddSuccess: successcount,
		AddFail:    failcount,
		FailList:   faillist,
	}
}

func parseClass(row []string, fieldMap map[string]int) (*models.Class, error) {
	//var err error
	class := &models.Class{
		ID:             getValueSafely(row, fieldMap["ID"]),
		Name:           getValueSafely(row, fieldMap["Name"]),
		Academic:       getValueSafely(row, fieldMap["Academic"]),
		Cultivation:    getValueSafely(row, fieldMap["Cultivation"]),
		Type:           getValueSafely(row, fieldMap["Type"]),
		ExpectedYear:   getValueSafely(row, fieldMap["ExpectedYear"]),
		IsGraduation:   getValueSafely(row, fieldMap["IsGraduation"]),
		StudentCount:   getValueSafely(row, fieldMap["StudentCount"]),
		MaxCount:       getValueSafely(row, fieldMap["MaxCount"]),
		Year:           getValueSafely(row, fieldMap["Year"]),
		Department:     getValueSafely(row, fieldMap["Department"]),
		MajorID:        getValueSafely(row, fieldMap["MajorID"]),
		Major:          getValueSafely(row, fieldMap["Major"]),
		Campus:         getValueSafely(row, fieldMap["Campus"]),
		FixedClassroom: getValueSafely(row, fieldMap["FixedClassroom"]),
	}

	// 处理班级人数
	//studentCountStr := getValueSafely(row, fieldMap["StudentCount"])
	//if class.StudentCount, err = strconv.Atoi(studentCountStr); err != nil {
	//	return nil, fmt.Errorf("%w: 班级人数", InvalidDataFormatErr)
	//}

	// 处理班级最大人数
	//maxCountStr := getValueSafely(row, fieldMap["MaxCount"])
	//if class.MaxCount, err = strconv.Atoi(maxCountStr); err != nil {
	//	return nil, fmt.Errorf("%w: 班级最大人数", InvalidDataFormatErr)
	//}

	if class.ID == "" || class.Name == "" {
		return nil, fmt.Errorf("班级编号或名称不能为空")
	}

	return class, nil
}

func getValueSafely(row []string, index int) string {
	if index < 0 || index >= len(row) {
		return ""
	}
	return strings.TrimSpace(row[index])
}

package course

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

// 定义结构体字段与Excel列名的映射
var fieldToColumn = map[string]string{
	"ID":            "课程编号",
	"Name":          "课程名称",
	"Type":          "课程类型",
	"Property":      "课程属性",
	"Credit":        "学分",
	"Department":    "开课院系",
	"TotalHours":    "总学时",
	"TheoryHours":   "理论学时",
	"TestHours":     "实验学时",
	"ComputerHours": "上机学时",
	"PracticeHours": "实践学时",
	"OtherHours":    "其他学时",
	"WeeklyHours":   "周学时",
	"PurePractice":  "是否纯实践环节",
}

type CourseAddByExcelFlow struct {
	Filename   string
	CourseList [][]string
}

func CourseAddByExcel(filename string) (*dto.CourseAddByExcelResp, error) {
	return NewCourseAddByExcelFlow(filename).Do()
}

func NewCourseAddByExcelFlow(filename string) *CourseAddByExcelFlow {
	return &CourseAddByExcelFlow{
		Filename: filename,
	}
}

func (f *CourseAddByExcelFlow) Do() (*dto.CourseAddByExcelResp, error) {
	courseList, err := f.ReadFile()
	if err != nil {
		return nil, err
	}
	f.CourseList = courseList

	// 检查是否存在标题行和数据行
	if len(f.CourseList) < 1 {
		return nil, errors.New("文件内容为空")
	}
	headers := f.CourseList[0]
	dataRows := f.CourseList[1:]

	// 建立列名到索引的映射
	columnMap := make(map[string]int)
	for i, header := range headers {
		columnMap[header] = i
	}

	// 验证必要列是否存在，并建立字段索引映射
	fieldMap := make(map[string]int)
	for field, colName := range fieldToColumn {
		idx, ok := columnMap[colName]
		if !ok {
			err = fmt.Errorf("%w: %s", MissingColumnErr, colName)
			log.Println(err)
			return nil, err
		}
		fieldMap[field] = idx
	}

	return f.AddCourse(fieldMap, dataRows), nil
}

func (f *CourseAddByExcelFlow) ReadFile() ([][]string, error) {
	res, err := utils.ReadExcel(f.Filename)
	if err != nil {
		if err == utils.ExcelFormatErr {
			return nil, FileFormatErr
		} else {
			return nil, err
		}
	}
	return res, nil
}

func (f *CourseAddByExcelFlow) AddCourse(fieldMap map[string]int, dataRows [][]string) *dto.CourseAddByExcelResp {
	faillist := make([]*dto.CourseCreateResp, 0)
	var successcount, failcount int64

	for _, row := range dataRows {
		//if i > 2 {
		//	break
		//}
		course, err := parseCourse(row, fieldMap)
		if err != nil {
			failcount++
			faillist = append(faillist, &dto.CourseCreateResp{
				CourseID:   getValueSafely(row, fieldMap["ID"]),
				CourseName: getValueSafely(row, fieldMap["Name"]),
				Err:        err.Error(),
			})
			continue
		}

		if err = models.NewCourseDao().CreateCourse(course); err != nil {
			failcount++
			faillist = append(faillist, &dto.CourseCreateResp{
				CourseID:   course.ID,
				CourseName: course.Name,
				Err:        err.Error(),
			})
		} else {
			successcount++
		}
	}

	return &dto.CourseAddByExcelResp{
		AddSuccess: successcount,
		AddFail:    failcount,
		FailList:   faillist,
	}
}

// parseCourse 根据映射关系解析课程数据
func parseCourse(row []string, fieldMap map[string]int) (*models.Course, error) {
	var err error
	course := &models.Course{}

	// 解析各字段，处理可能的转换错误
	course.ID = getValueSafely(row, fieldMap["ID"])
	course.Name = getValueSafely(row, fieldMap["Name"])
	if course.ID == "" || course.Name == "" {
		return nil, fmt.Errorf("课程名或课程ID不能为空")
	}
	course.Type = getValueSafely(row, fieldMap["Type"])
	course.Property = getValueSafely(row, fieldMap["Property"])

	if course.Credit, err = utils.ParseFloat(getValueSafely(row, fieldMap["Credit"])); err != nil {
		return nil, fmt.Errorf("%w: 学分", InvalidDataFormatErr)
	}

	course.Department = getValueSafely(row, fieldMap["Department"])

	if course.TotalHours, err = utils.ParseInt(getValueSafely(row, fieldMap["TotalHours"])); err != nil {
		return nil, fmt.Errorf("%w: 总学时", InvalidDataFormatErr)
	}

	// 类似处理其他数值字段...
	// 示例：解析TheoryHours
	if course.TheoryHours, err = utils.ParseInt(getValueSafely(row, fieldMap["TheoryHours"])); err != nil {
		return nil, fmt.Errorf("%w: 理论学时", InvalidDataFormatErr)
	}

	if course.TestHours, err = utils.ParseInt(getValueSafely(row, fieldMap["TestHours"])); err != nil {
		return nil, fmt.Errorf("%w: 实验学时", InvalidDataFormatErr)
	}

	if course.PracticeHours, err = utils.ParseInt(getValueSafely(row, fieldMap["PracticeHours"])); err != nil {
		return nil, fmt.Errorf("%w: 实践学时", InvalidDataFormatErr)
	}

	if course.ComputerHours, err = utils.ParseInt(getValueSafely(row, fieldMap["ComputerHours"])); err != nil {
		return nil, fmt.Errorf("%w: 上机学时", InvalidDataFormatErr)
	}

	if course.OtherHours, err = utils.ParseInt(getValueSafely(row, fieldMap["OtherHours"])); err != nil {
		return nil, fmt.Errorf("%w: 其它学时", InvalidDataFormatErr)
	}

	if course.WeeklyHours, err = utils.ParseInt(getValueSafely(row, fieldMap["WeeklyHours"])); err != nil {
		return nil, fmt.Errorf("%w: 周学时", InvalidDataFormatErr)
	}

	// 解析布尔值
	purePracticeStr := getValueSafely(row, fieldMap["PurePractice"])
	if course.PurePractice, err = utils.ParseBool(purePracticeStr); err != nil {
		return nil, fmt.Errorf("%w: 是否纯实践", InvalidDataFormatErr)
	}

	return course, nil
}

// getValueSafely 安全获取行中的值，防止索引越界
func getValueSafely(row []string, index int) string {
	if index < 0 || index >= len(row) {
		return ""
	}
	return strings.TrimSpace(row[index])
}

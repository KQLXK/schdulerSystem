package course

import (
	"errors"
	"schedule/commen/utils"
	"schedule/dto"
	"schedule/models"
	"strconv"
	"strings"
)

var (
	FileFormatErr = errors.New("文件格式有误")
)

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
	courselist, err := f.ReadFile()
	if err != nil {
		return nil, err
	} else {
		f.CourseList = courselist
	}
	return f.AddCourse(), nil
}

func (f *CourseAddByExcelFlow) ReadFile() ([][]string, error) {
	data := strings.Split(f.Filename, ".")
	if data[1] == "xlsx" {
		res, err := utils.ReadXlsx(f.Filename)
		if err != nil {
			return nil, err
		}
		return res, nil
	} else if data[1] == "xls" {
		res, err := utils.ReadXls(f.Filename)
		if err != nil {
			return nil, err
		}
		return res, nil
	}
	return nil, FileFormatErr
}

func (f *CourseAddByExcelFlow) AddCourse() *dto.CourseAddByExcelResp {
	faillist := make([]*dto.CourseCreateResp, 0)
	var successcount, failcount int
	for _, row := range f.CourseList {
		purepractice, _ := strconv.ParseBool(row[13])
		course := &models.Course{
			ID:            row[0],
			Name:          row[1],
			Type:          row[2],
			Property:      row[3],
			Credit:        utils.ParseFloat(row[4]),
			Department:    row[5],
			TotalHours:    utils.ParseInt(row[6]),
			TheoryHours:   utils.ParseInt(row[7]),
			TestHours:     utils.ParseInt(row[8]),
			ComputerHours: utils.ParseInt(row[9]),
			PracticeHours: utils.ParseInt(row[10]),
			OtherHours:    utils.ParseInt(row[11]),
			WeeklyHours:   utils.ParseInt(row[12]),
			PurePractice:  purepractice,
		}
		if err := models.NewCourseDao().CreateCourse(course); err != nil {
			faillist = append(faillist, &dto.CourseCreateResp{
				CourseID:   course.ID,
				CourseName: course.Name,
			})
			failcount++
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

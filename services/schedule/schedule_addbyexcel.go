package schedule

import (
	"errors"
	"fmt"
	"schedule/commen/utils"
	"schedule/dto"
	"schedule/models"
	"strings"
)

var (
	FileFormatErr        = errors.New("不支持的文件格式")
	FileEmptyErr         = errors.New("文件内容为空")
	MissingColumnErr     = errors.New("缺少必要的列")
	InvalidDataFormatErr = errors.New("数据格式错误")
)

var filedToColumn map[string]string = map[string]string{
	"Semester":               "学年学期",
	"CourseID":               "课程编号",
	"CourseName":             "课程名称",
	"TeacherID":              "教师工号",
	"TeachingClass":          "教学班组成",
	"TeachingClassID":        "教学班编号",
	"TeachingClassName":      "教学班名称",
	"HourType":               "学时类型",
	"OpeningHours":           "开课学时",
	"SchedulingHours":        "排课学时",
	"TotalHours":             "总学时",
	"SchedulingPriority":     "排课优先级",
	"TeachingClassSize":      "教学班人数",
	"OpeningCampus":          "开课校区",
	"OpeningWeekHours":       "开课周次学时",
	"ContinuousPeriods":      "连排节次",
	"SpecifiedClassroomType": "指定教室类型",
	"SpecifiedClassroom":     "指定教室",
	"SpecifiedBuilding":      "指定教学楼",
	"SpecifiedTime":          "指定时间",
}

type ScheduleAddByExcelFlow struct {
	Filename     string
	ScheduleList [][]string
}

func AddByExcel(filename string) (*dto.ScheduleAddByExcelResp, error) {
	return NewScheduleAddByExcelFlow(filename).Do()
}

func NewScheduleAddByExcelFlow(filename string) *ScheduleAddByExcelFlow {
	return &ScheduleAddByExcelFlow{
		Filename: filename,
	}
}

func (f *ScheduleAddByExcelFlow) Do() (*dto.ScheduleAddByExcelResp, error) {
	if err := f.ReadFile(); err != nil {
		return nil, err
	}
	if len(f.ScheduleList) < 1 {
		return nil, FileEmptyErr
	}
	headers := f.ScheduleList[0]
	dataRows := f.ScheduleList[1:]

	columnMap := make(map[string]int)
	for idx, header := range headers {
		columnMap[header] = idx
	}

	fieldMap := make(map[string]int)
	for filed, colName := range filedToColumn {
		idx, ok := columnMap[colName]
		if !ok {
			return nil, fmt.Errorf("%v:%s", MissingColumnErr, colName)
		}
		fieldMap[filed] = idx
	}

	resp := f.AddSchedule(fieldMap, dataRows)
	return resp, nil

}

func (f *ScheduleAddByExcelFlow) ReadFile() error {
	schedulelist, err := utils.ReadExcel(f.Filename)
	if err != nil {
		if err == utils.ExcelFormatErr {
			return FileFormatErr
		} else {
			return err
		}
	}
	f.ScheduleList = schedulelist
	return nil
}

func (f *ScheduleAddByExcelFlow) AddSchedule(fieldMap map[string]int, dataRows [][]string) *dto.ScheduleAddByExcelResp {
	failList := make([]*dto.ScheduleCreateResp, 0)
	var successCount, failCount int64

	for _, row := range dataRows {
		schedule, err := parseSchedule(row, fieldMap)
		if err != nil {
			failCount++
			id, _ := utils.ParseInt(getValueSafely(row, fieldMap["ScheduleID"]))
			failList = append(failList, &dto.ScheduleCreateResp{
				ScheduleID:  id,
				CourseName:  getValueSafely(row, fieldMap["CourseName"]),
				TeacherName: getValueSafely(row, fieldMap["TeacherName"]),
				Semester:    getValueSafely(row, fieldMap["Semester"]),
				Err:         err.Error(),
			})
			continue
		}

		if err := models.NewScheduleDao().CreateSchedule(schedule); err != nil {
			failCount++
			failList = append(failList, &dto.ScheduleCreateResp{
				ScheduleID: schedule.ID,
				Semester:   schedule.Semester,
				Err:        err.Error(),
			})
		} else {
			successCount++
		}
	}

	return &dto.ScheduleAddByExcelResp{
		AddSuccess: successCount,
		AddFail:    failCount,
		FailList:   failList,
	}
}

func parseSchedule(row []string, fieldMap map[string]int) (*models.Schedule, error) {
	var err error
	schedule := &models.Schedule{}

	// 基础字段
	schedule.Semester = getValueSafely(row, fieldMap["Semester"])
	schedule.CourseID = getValueSafely(row, fieldMap["CourseID"])
	schedule.CourseName = getValueSafely(row, fieldMap["CourseName"])
	schedule.TeacherID = getValueSafely(row, fieldMap["TeacherID"])
	schedule.TeachingClassName = getValueSafely(row, fieldMap["TeachingClassName"])
	schedule.TeachingClass = getValueSafely(row, fieldMap["TeachingClass"])
	schedule.TeachingClassID = getValueSafely(row, fieldMap["TeachingClassID"])
	schedule.HourType = getValueSafely(row, fieldMap["HourType"])
	schedule.OpeningCampus = getValueSafely(row, fieldMap["OpeningCampus"])
	schedule.OpeningWeekHours = getValueSafely(row, fieldMap["OpeningWeekHours"])
	schedule.SpecifiedBuilding = getValueSafely(row, fieldMap["SpecifiedBuilding"])
	schedule.SpecifiedClassroomType = getValueSafely(row, fieldMap["SpecifiedClassroomType"])
	schedule.SpecifiedClassroom = getValueSafely(row, fieldMap["SpecifiedClassroom"])
	schedule.SpecifiedTime = getValueSafely(row, fieldMap["SpecifiedTime"])

	// 数值型字段
	if schedule.OpeningHours, err = utils.ParseInt(getValueSafely(row, fieldMap["OpeningHours"])); err != nil {
		return nil, fmt.Errorf("%w: 开课学时", InvalidDataFormatErr)
	}

	if schedule.SchedulingPriority, err = utils.ParseInt(getValueSafely(row, fieldMap["SchedulingPriority"])); err != nil {
		return nil, fmt.Errorf("%w: 排课优先级", InvalidDataFormatErr)
	}

	if schedule.TeachingClassSize, err = utils.ParseInt(getValueSafely(row, fieldMap["TeachingClassSize"])); err != nil {
		return nil, fmt.Errorf("%w: 教学班人数", InvalidDataFormatErr)
	}

	if schedule.SchedulingHours, err = utils.ParseInt(getValueSafely(row, fieldMap["SchedulingHours"])); err != nil {
		return nil, fmt.Errorf("%w: 排课学时", InvalidDataFormatErr)
	}

	if schedule.ContinuousPeriods, err = utils.ParseInt(getValueSafely(row, fieldMap["ContinuousPeriods"])); err != nil {
		return nil, fmt.Errorf("%w: 连排节次", InvalidDataFormatErr)
	}

	if schedule.TotalHours, err = utils.ParseInt(getValueSafely(row, fieldMap["TotalHours"])); err != nil {
		return nil, fmt.Errorf("%w: 总学时", InvalidDataFormatErr)
	}

	if schedule.OpeningHours, err = utils.ParseInt(getValueSafely(row, fieldMap["OpeningHours"])); err != nil {
		return nil, fmt.Errorf("%w: 开课学时", InvalidDataFormatErr)
	}

	return schedule, nil
}

func getValueSafely(row []string, idx int) string {
	if idx < 0 || idx >= len(row) {
		return ""
	}
	return strings.TrimSpace(row[idx])
}

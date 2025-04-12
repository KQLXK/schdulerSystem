package table

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"schedule/database"
	"schedule/dto"
	"schedule/models"
	"strings"
)

type ManualScheduleFlow struct {
	*dto.ManualScheduleRequest
	timeslots []models.TimeSlot
	schedule  *models.Schedule
	course    *models.Course
	classroom *models.Classroom
	teacher   *models.Teacher
	classes   []models.Class
}

var (
	CourseExistsErr    = errors.New("课程不存在")
	ScheduleExistsErr  = errors.New("排课任务不存在")
	ClassExistsErr     = errors.New("班级不存在")
	TeacherExistsErr   = errors.New("教师不存在")
	ClassroomExistsErr = errors.New("教室不存在")
	TimeSlotExistsErr  = errors.New("时间段冲突")
)

func ManualSchedule(req *dto.ManualScheduleRequest) error {
	return NewManualScheduleFlow(req).Do()
}

func NewManualScheduleFlow(req *dto.ManualScheduleRequest) *ManualScheduleFlow {
	return &ManualScheduleFlow{ManualScheduleRequest: req}
}

// Do 执行手动排课流程
func (m *ManualScheduleFlow) Do() error {
	// 使用事务保证数据一致性
	return database.DB.Transaction(func(tx *gorm.DB) error {
		if err := m.checkSchedule(); err != nil {
			return err
		}
		if err := m.checkClassroom(); err != nil {
			return err
		}
		if err := m.prepareData(); err != nil {
			return err
		}
		if err := m.checkTimeSlot(); err != nil {
			return err
		}
		// 6. 创建排课结果
		return m.createSchedule()
	})
}

func (m *ManualScheduleFlow) checkSchedule() error {
	var schedule *models.Schedule
	schedule, err := models.NewScheduleDao().GetScheduleByID(m.ScheduleID)
	if err != nil {
		return fmt.Errorf("%w: %d", ScheduleExistsErr, m.ScheduleID)
	}
	m.schedule = schedule
	return nil
}

// prepareData 准备数据
func (m *ManualScheduleFlow) prepareData() error {
	course, _ := models.NewCourseDao().GetCourseByID(m.schedule.CourseID)
	classes := parseTeachingClasses(m.schedule.TeachingClass)
	teacher, _ := models.GetTeacherByID(m.schedule.TeacherID)
	timeslots := dto.ConvertSlotsToModel(m.TimeSlots)
	m.timeslots = timeslots
	m.course = course
	m.teacher = teacher
	m.classes = classes
	return nil
}

func (m *ManualScheduleFlow) checkClassroom() error {
	classroom, err := models.NewClassroomDao().GetClassroomByID(m.ClassroomID)
	if err != nil {
		return fmt.Errorf("%w: %s", ClassroomExistsErr, m.ClassroomID)
	}
	m.classroom = classroom
	if m.classroom.Type != m.schedule.SpecifiedClassroomType {
		return fmt.Errorf("教室类型不符合要求, 期望: %s, 实际: %s", m.schedule.SpecifiedClassroomType, m.classroom.Type)
	}
	return nil
}

// checkTimeSlot 检查时间冲突
func (m *ManualScheduleFlow) checkTimeSlot() error {
	// 检查教室时间冲突
	if conflict, err := m.checkClassroomConflict(); err != nil || conflict {
		return fmt.Errorf("教室时间冲突: %w", err)
	}
	// 检查教师时间冲突
	if conflict, err := m.checkTeacherConflict(); err != nil || conflict {
		return fmt.Errorf("教师时间冲突: %w", err)
	}
	// 检查班级时间冲突
	if conflict, err := m.checkClassConflict(); err != nil || conflict {
		return fmt.Errorf("班级时间冲突: %w", err)
	}
	return nil
}

// checkClassroomConflict 检查教室时间冲突
func (m *ManualScheduleFlow) checkClassroomConflict() (bool, error) {
	var existing []models.ScheduleResult
	tx := database.DB
	if err := tx.Where("classroom_id = ?", m.ClassroomID).Find(&existing).Error; err != nil {
		return false, err
	}
	return m.hasTimeConflict(existing)
}

// checkTeacherConflict 检查教师时间冲突
func (m *ManualScheduleFlow) checkTeacherConflict() (bool, error) {
	var existing []models.ScheduleResult
	tx := database.DB
	if err := tx.Where("teacher_id = ?", m.teacher.ID).Find(&existing).Error; err != nil {
		return false, err
	}
	return m.hasTimeConflict(existing)
}

// checkClassConflict 检查班级时间冲突
func (m *ManualScheduleFlow) checkClassConflict() (bool, error) {
	for _, class := range m.classes {
		var existing []models.ScheduleResult
		// 使用JSON_CONTAINS查询包含该班级的排课
		existing, err := GetClassScheduleBySemester(class.ID, m.schedule.Semester)
		if err != nil {
			return false, err
		}
		if existing == nil {
			return false, nil
		}
		if conflict, err := m.hasTimeConflict(existing); err != nil || conflict {
			return true, err
		}
	}
	return false, nil
}

// hasTimeConflict 检查时间冲突核心逻辑
func (m *ManualScheduleFlow) hasTimeConflict(existing []models.ScheduleResult) (bool, error) {
	for _, result := range existing {
		for _, existSlot := range result.TimeSlots {
			for _, newSlot := range m.timeslots {
				if isTimeOverlap(existSlot, newSlot) {
					return true, TimeSlotExistsErr
				}
			}
		}
	}
	return false, nil
}

// isTimeOverlap 判断时间段是否重叠
func isTimeOverlap(a, b models.TimeSlot) bool {
	// 周次交集检查
	weekSet := make(map[int]bool)
	for _, week := range a.WeekNumbers {
		weekSet[week] = true
	}
	hasCommonWeek := false
	for _, week := range b.WeekNumbers {
		if weekSet[week] {
			hasCommonWeek = true
			break
		}
	}
	if !hasCommonWeek {
		return false
	}

	// 星期几检查
	if a.Weekday != b.Weekday {
		return false
	}

	// 节次范围检查
	aStart := a.StartPeriod
	aEnd := a.StartPeriod + a.Duration - 1
	bStart := b.StartPeriod
	bEnd := b.StartPeriod + b.Duration - 1

	return aStart <= bEnd && bStart <= aEnd
}

// createSchedule 创建排课记录
func (m *ManualScheduleFlow) createSchedule() error {
	// 收集班级名称
	classNames := make([]string, len(m.classes))
	classIDs := make([]string, len(m.classes))
	for i, class := range m.classes {
		classIDs[i] = class.ID
		classNames[i] = class.Name
	}

	scheduleResult := &models.ScheduleResult{
		Semester:    m.schedule.Semester,
		ScheduleID:  m.ScheduleID,
		CourseID:    m.schedule.CourseID,
		CourseName:  m.course.Name,
		ClassroomID: m.ClassroomID,
		TeacherID:   m.schedule.TeacherID,
		TeacherName: m.teacher.Name,
		ClassIDs:    models.JSONStrings(classIDs),
		ClassNames:  models.JSONStrings(classNames),
		TimeSlots:   models.JSONTimeSlots(m.timeslots),
	}

	return models.CreateScheduleResult(scheduleResult)
}

func parseTeachingClasses(classes string) []models.Class {
	classlist := strings.Split(classes, ",")
	res := make([]models.Class, 0)
	for _, val := range classlist {
		class, _ := models.NewClassDaoInstance().GetClassByName(val)
		res = append(res, *class)
	}
	return res
}

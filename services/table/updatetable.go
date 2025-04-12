package table

import (
	"errors"
	"fmt"
	"schedule/database"
	"schedule/dto"
	"schedule/models"
)

var (
	ErrInvalidID         = errors.New("排课结果ID不存在")
	ErrInvalidHours      = errors.New("排课结果学时无效")
	ErrScheduleConflict  = errors.New("排课结果冲突")
	ErrClassroomNotFound = errors.New("教室不存在")
	ErrClassroomConflict = errors.New("教室冲突")
)

type TableUpdateFlow struct {
	*dto.ScheduleResultUpdateReq
	timeslots []models.TimeSlot
	result    *models.ScheduleResult
}

func UpdateTable(req *dto.ScheduleResultUpdateReq) error {
	return NewTableUpdateFlow(req).Do()
}

func NewTableUpdateFlow(req *dto.ScheduleResultUpdateReq) *TableUpdateFlow {
	return &TableUpdateFlow{
		ScheduleResultUpdateReq: req,
	}
}

func (f *TableUpdateFlow) Do() error {
	if err := f.Validate(); err != nil {
		return err
	}
	//if err := f.CheckExists(); err != nil {
	//	return err
	//}
	if err := f.CheckTimeslots(); err != nil {
		return err
	}
	return f.UpdateTable()
}

func (f *TableUpdateFlow) Validate() error {
	result, err := models.GetScheduleResultByID(f.ScheduleResultID)
	if err != nil {
		return ErrInvalidID
	}
	if _, err := models.NewClassroomDao().GetClassroomByID(f.ClassroomID); err != nil {
		return ErrClassroomNotFound
	}
	timeslots := dto.ConvertSlotsToModel(f.TimeSlots)
	f.result = result
	f.timeslots = timeslots
	return nil
}

func (f *TableUpdateFlow) CheckTimeslots() error {
	// 检查教室时间冲突
	if conflict, err := f.checkClassroomConflict(); err != nil || conflict {
		return fmt.Errorf("教室时间冲突: %w", err)
	}
	// 检查教师时间冲突
	if conflict, err := f.checkTeacherConflict(); err != nil || conflict {
		return fmt.Errorf("教师时间冲突: %w", err)
	}
	// 检查班级时间冲突
	if conflict, err := f.checkClassConflict(); err != nil || conflict {
		return fmt.Errorf("班级时间冲突: %w", err)
	}
	return nil
}

func (f *TableUpdateFlow) checkClassroomConflict() (bool, error) {
	var existing []models.ScheduleResult
	tx := database.DB
	if err := tx.Where("classroom_id = ?", f.ClassroomID).Find(&existing).Error; err != nil {
		return false, err
	}
	return f.hasTimeConflict(existing)
}

func (f *TableUpdateFlow) checkTeacherConflict() (bool, error) {
	var existing []models.ScheduleResult
	tx := database.DB

	if err := tx.Where("teacher_id =?", f.result.TeacherID).Find(&existing).Error; err != nil {
		return false, err
	}
	return f.hasTimeConflict(existing)
}

func (f *TableUpdateFlow) checkClassConflict() (bool, error) {
	for _, classID := range f.result.ClassIDs {
		var existing []models.ScheduleResult
		existing, err := GetClassScheduleBySemester(classID, f.result.Semester)
		if err != nil {
			return false, err
		}
		if conflict, err := f.hasTimeConflict(existing); err != nil || conflict {
			return true, nil
		}
	}
	return false, nil
}

func (f *TableUpdateFlow) UpdateTable() error {
	newresult := models.ScheduleResult{
		Semester:    f.result.Semester,
		ScheduleID:  f.result.ScheduleID,
		CourseID:    f.result.CourseID,
		CourseName:  f.result.CourseName,
		ClassroomID: f.ClassroomID,
		TeacherID:   f.result.TeacherID,
		TeacherName: f.result.TeacherName,
		TimeSlots:   f.timeslots,
	}
	return models.UpdateScheduleResult(f.ScheduleResultID, &newresult)
}

func (f *TableUpdateFlow) hasTimeConflict(existing []models.ScheduleResult) (bool, error) {
	for _, result := range existing {
		for _, timeslot := range result.TimeSlots {
			for _, newtimeslot := range f.timeslots {
				if isTimeOverlap(timeslot, newtimeslot) {
					return true, TimeSlotExistsErr
				}
			}
		}
	}
	return false, nil
}

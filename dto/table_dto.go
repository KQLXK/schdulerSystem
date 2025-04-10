package dto

import "schedule/models"

type GetClassTableResp struct {
	ClassTables []ClassTable
}

type ClassTable struct {
	ID          int                  `json:"id"`
	Semester    string               `json:"semester"`
	ScheduleID  int64                `json:"schedule_id"`
	CourseID    string               `json:"course_id"`
	CourseName  string               `json:"course_name"`
	TeacherID   string               `json:"teacher_id"`
	TeacherName string               `json:"teacher_name"`
	ClassroomID string               `json:"classroom_id"`
	ClassIDs    []string             `json:"class_ids"`
	ClassNames  []string             `json:"class_name"`
	Timeslots   models.JSONTimeSlots `json:"time_slots"`
}

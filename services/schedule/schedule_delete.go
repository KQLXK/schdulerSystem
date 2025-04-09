package schedule

import (
	"schedule/models"
)

// ScheduleDeleteFlow 排课删除工作流
type ScheduleDeleteFlow struct {
	ScheduleID string
}

// 快速删除入口
func DeleteSchedule(scheduleID string) error {
	return NewScheduleDeleteFlow(scheduleID).Do()
}

// 创建删除工作流实例
func NewScheduleDeleteFlow(scheduleID string) *ScheduleDeleteFlow {
	return &ScheduleDeleteFlow{ScheduleID: scheduleID}
}

// 执行删除操作
func (f *ScheduleDeleteFlow) Do() error {
	// 直接调用DAO层删除（依赖DAO处理软删除/硬删除）
	if err := models.NewScheduleDao().DeleteSchedule(f.ScheduleID); err != nil {
		return err
	}
	return nil
}

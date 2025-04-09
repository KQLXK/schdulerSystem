package schedule

import (
	"errors"
	"gorm.io/gorm"
	"schedule/models"
)

// ScheduleGetFlow 排课详情查询工作流
type ScheduleGetFlow struct {
	ScheduleID int
}

// 快速查询入口
func QuerySchedule(scheduleID int) (*models.Schedule, error) {
	return NewScheduleGetFlow(scheduleID).Do()
}

// 创建查询实例
func NewScheduleGetFlow(scheduleID int) *ScheduleGetFlow {
	return &ScheduleGetFlow{ScheduleID: scheduleID}
}

// 执行查询操作
func (f *ScheduleGetFlow) Do() (*models.Schedule, error) {
	// 可选：添加ID有效性校验
	if f.ScheduleID <= 0 {
		return nil, errors.New("无效的排课ID")
	}

	schedule, err := models.NewScheduleDao().GetScheduleByID(f.ScheduleID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrScheduleNotFound
		}
		return nil, err // 返回原始数据库错误
	}
	return schedule, nil
}

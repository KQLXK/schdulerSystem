package algorithm

import (
	"fmt"
	"schedule/dto"
	"schedule/models"
	"sync"
	"time"
)

func GASchedule(conf *dto.Config, ScheduleList []int) *dto.ScheduleGAResponse {
	scheduler := NewScheduler(conf, ScheduleList)
	chromosome := scheduler.Run()
	resp := scheduler.GenerateScheduleResponse(chromosome)
	return resp
}

// 生成排课结果响应（包含成功/失败列表和统计分析）
func (s *Scheduler) GenerateScheduleResponse(chromosome Chromosome) *dto.ScheduleGAResponse {
	response := &dto.ScheduleGAResponse{
		Status:      "processing",
		GeneratedAt: time.Now().UTC(),
	}

	// Step 1: 执行排课验证
	validation := s.ValidateSchedule(chromosome)

	// Step 2: 转换验证结果
	var wg sync.WaitGroup
	wg.Add(3)

	// 处理成功条目
	go func() {
		defer wg.Done()
		successItems := make([]dto.SuccessItem, 0)
		for _, gene := range chromosome {
			if !isGeneFailed(gene.ScheduleID, validation.FailedGenes) {
				successItems = append(successItems, convertToSuccessItem(gene, s))
				schedule := convertToScheduleResult(gene, s)
				_ = models.CreateScheduleResult(&schedule)
			}
		}
		response.SuccessList = successItems
	}()

	// 处理失败条目
	go func() {
		defer wg.Done()
		failedItems := make([]dto.FailedItem, len(validation.FailedGenes))
		for i, fg := range validation.FailedGenes {
			failedItems[i] = dto.FailedItem{
				ScheduleID:      fg.ScheduleID,
				CourseName:      s.scheduleMap[int64(fg.ScheduleID)].CourseName,
				ClassroomID:     fg.ClassroomID,
				TeacherID:       fg.TeacherID,
				TimeSlots:       fg.TimeSlots,
				TeacherName:     s.TeacherMap[fg.TeacherID].Name,
				ClassroomName:   s.ClassroomMap[fg.ClassroomID].Name,
				ConflictReasons: fg.ConflictReasons,
				//SuggestedTimes:  s.suggestAlternativeTimes(fg),
				//AlternativeRooms: s.findAlternativeClassrooms(fg),
			}
		}
		response.FailedList = failedItems
	}()

	// 生成分析数据
	go func() {
		defer wg.Done()
		analysis := dto.ResultAnalysis{
			ClassroomUtilization: s.calculateClassroomUtilization(chromosome),
			TeacherWorkload:      s.calculateTeacherWorkload(chromosome),
			TimeDistribution:     s.calculateTimeDistribution(chromosome),
		}
		response.Analysis = analysis
	}()

	wg.Wait()

	// 更新最终状态
	if validation.IsValid && len(response.FailedList) == 0 {
		response.Status = "success"
		response.Progress = 1.0
	} else {
		response.Status = "failed"
		response.Progress = 1.0
	}

	return response
}

// --- 辅助函数 ---

// 检查基因是否在失败列表中
func isGeneFailed(scheduleID int64, failedGenes []FailedGene) bool {
	for _, fg := range failedGenes {
		if int64(fg.ScheduleID) == scheduleID {
			return true
		}
	}
	return false
}

// 转换成功条目
func convertToSuccessItem(gene ScheduleGene, s *Scheduler) dto.SuccessItem {
	schedule := s.scheduleMap[gene.ScheduleID]
	classroom := s.ClassroomMap[gene.ClassroomID]
	teacher := s.TeacherMap[gene.TeacherID]

	return dto.SuccessItem{
		ScheduleID:    int(gene.ScheduleID),
		CourseID:      schedule.CourseID,
		CourseName:    schedule.CourseName,
		TimeSlots:     gene.TimeSlots,
		ClassroomID:   classroom.ID,
		ClassroomName: classroom.Name,
		TeacherID:     teacher.ID,
		TeacherName:   teacher.Name,
		//ClassIDs:      schedule.TeachingClass,
	}
}

func convertToScheduleResult(gene ScheduleGene, s *Scheduler) models.ScheduleResult {
	schedule := s.scheduleMap[gene.ScheduleID]
	teacher := s.TeacherMap[gene.TeacherID]
	class := s.parseTeachingClasses(schedule.TeachingClass)
	var classIDs, classNames []string
	for _, c := range class {
		classIDs = append(classIDs, c.ID)
		classNames = append(classNames, c.Name)
	}
	return models.ScheduleResult{
		Semester:    schedule.Semester,
		ScheduleID:  int(schedule.ID),
		CourseID:    schedule.CourseID,
		CourseName:  schedule.CourseName,
		ClassroomID: gene.ClassroomID,
		TeacherID:   gene.TeacherID,
		TeacherName: teacher.Name,
		TimeSlots:   gene.TimeSlots,
		ClassIDs:    classIDs,
		ClassNames:  classNames,
	}
}

// 计算教室使用率
func (s *Scheduler) calculateClassroomUtilization(chromosome Chromosome) map[string]float64 {
	utilization := make(map[string]float64)
	totalPeriods := s.Config.SemesterWeek * 5 * s.Config.MaxPeriodsPerDay

	// 统计每个教室的总使用节次
	classroomUsage := make(map[string]int)
	for _, gene := range chromosome {
		for _, slot := range gene.TimeSlots {
			periods := s.GetAbsolutePeriods(slot)
			classroomUsage[gene.ClassroomID] += len(periods)
		}
	}

	// 计算使用率
	for classroomID, used := range classroomUsage {
		utilization[classroomID] = float64(used) / float64(totalPeriods)
	}
	return utilization
}

// 计算教师工作量
func (s *Scheduler) calculateTeacherWorkload(chromosome Chromosome) map[string]int {
	workload := make(map[string]int)
	for _, gene := range chromosome {
		for _, slot := range gene.TimeSlots {
			workload[gene.TeacherID] += len(slot.WeekNumbers) * slot.Duration
		}
	}
	return workload
}

// 计算时间分布
func (s *Scheduler) calculateTimeDistribution(chromosome Chromosome) map[string]int {
	distribution := make(map[string]int)
	for _, gene := range chromosome {
		for _, slot := range gene.TimeSlots {
			key := fmt.Sprintf("周%d-星期%d-%d节",
				slot.WeekNumbers[0], // 取第一个周次代表
				slot.Weekday,
				slot.StartPeriod,
			)
			distribution[key] += 1
		}
	}
	return distribution
}

// 建议替代时间
func (s *Scheduler) suggestAlternativeTimes(fg FailedGene) []models.TimeSlot {
	// 实现基于时间窗口的推荐算法（示例逻辑）
	return []models.TimeSlot{
		{
			WeekNumbers: []int{5, 6, 7},
			Weekday:     2,
			StartPeriod: 3,
			Duration:    2,
		},
	}
}

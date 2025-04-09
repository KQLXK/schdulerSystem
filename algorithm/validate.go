package algorithm

import (
	"fmt"
	"log"
	"schedule/models"
	"strings"
	"sync"
	"time"
)

// 更新ScheduleValidation结构体
type ScheduleValidation struct {
	IsValid       bool
	SuccessCount  int // 新增字段
	FailedGenes   []FailedGene
	ConflictStats map[string]int
}

// 失败基因详情
type FailedGene struct {
	ScheduleID      int
	ClassroomID     string
	TeacherID       string
	TimeSlots       []models.TimeSlot
	ConflictReasons []string
	SuggestedFix    string
}

// 全局冲突记录器
type ConflictRecorder struct {
	mu     sync.Mutex
	report map[int][]Conflict // key: ScheduleID
}

type Conflict struct {
	Type       string // "教室冲突"/"教师冲突"/"班级冲突"
	Timestamp  time.Time
	Details    string
	RelatedIDs []string // 相关资源ID
}

// 新增结构体
type ClassroomSlotKey struct {
	ClassroomID string
	Week        int
	Weekday     int
	Period      int
}

type ResourceTracker struct {
	// classroomID -> 绝对节次集合
	classroomOccupancy map[string]map[int]bool
	// teacherID -> 绝对节次集合
	teacherOccupancy map[string]map[int]bool
	// classID -> 绝对节次集合
	classOccupancy map[string]map[int]bool
	mu             sync.Mutex
}

func (s *Scheduler) ValidateSchedule(chromosome Chromosome) ScheduleValidation {
	validation := ScheduleValidation{
		IsValid:       true,
		ConflictStats: make(map[string]int),
	}

	// 初始化资源占用记录
	resourceTracker := NewResourceTracker()

	// 冲突记录器初始化
	conflictRecorder := NewConflictRecorder()

	for _, gene := range chromosome {
		schedule := s.findScheduleByID(gene.ScheduleID)
		if schedule == nil {
			recordInvalidGene(&validation, gene, "无效教学班ID")
			continue
		}

		// 并行检查多个约束
		var wg sync.WaitGroup
		checkers := []func(){
			func() { s.checkClassroomConflicts(gene, resourceTracker, conflictRecorder) },
			func() { s.checkTeacherConflicts(gene, resourceTracker, conflictRecorder) },
			func() { s.checkClassConflicts(gene, resourceTracker, conflictRecorder) },
			func() { s.checkTimeValidity(gene, conflictRecorder) },
		}

		for _, checker := range checkers {
			wg.Add(1)
			go func(f func()) {
				defer wg.Done()
				f()
			}(checker)
		}
		wg.Wait()

		// 收集当前基因的冲突
		if conflicts := conflictRecorder.GetConflicts(int(gene.ScheduleID)); len(conflicts) > 0 {
			validation.IsValid = false
			failedGene := buildFailedGene(gene, conflicts)
			validation.FailedGenes = append(validation.FailedGenes, failedGene)

			for _, c := range conflicts {
				validation.ConflictStats[c.Type]++
			}
		}
	}

	// 补充检查全局约束
	s.checkGlobalConstraints(resourceTracker, conflictRecorder)

	return validation
}

// 示例：详细教室冲突检查
// 检查教室冲突（示例）
func (s *Scheduler) checkClassroomConflicts(gene ScheduleGene, tracker *ResourceTracker, recorder *ConflictRecorder) {
	classroomID := gene.ClassroomID

	for _, slot := range gene.TimeSlots {
		// 获取该时间段对应的所有绝对节次
		periods := s.GetAbsolutePeriods(slot)

		tracker.mu.Lock()
		// 初始化教室记录
		if tracker.classroomOccupancy[classroomID] == nil {
			tracker.classroomOccupancy[classroomID] = make(map[int]bool)
		}

		// 检查每个节次是否已占用
		for _, p := range periods {
			if tracker.classroomOccupancy[classroomID][p] {
				recorder.RecordConflict(int(gene.ScheduleID), Conflict{
					Type:    "教室冲突",
					Details: fmt.Sprintf("教室 %s 在第%d节已被占用", classroomID, p),
				})
			} else {
				tracker.classroomOccupancy[classroomID][p] = true
			}
		}
		tracker.mu.Unlock()
	}
}

// 教师冲突检查
func (s *Scheduler) checkTeacherConflicts(gene ScheduleGene, tracker *ResourceTracker, recorder *ConflictRecorder) {
	teacherID := gene.TeacherID
	for _, slot := range gene.TimeSlots {

		periods := s.GetAbsolutePeriods(slot)

		// 创建教师时间槽唯一标识
		tracker.mu.Lock()
		if tracker.classroomOccupancy[teacherID] == nil {
			tracker.classroomOccupancy[teacherID] = make(map[int]bool)
		}
		for _, p := range periods {
			if tracker.classroomOccupancy[teacherID][p] {
				recorder.RecordConflict(int(gene.ScheduleID), Conflict{
					Type:    "教师冲突",
					Details: fmt.Sprintf("教师 %s 在第%d节已被占用", teacherID, p),
				})
			} else {
				tracker.classroomOccupancy[teacherID][p] = true
			}
		}

		tracker.mu.Unlock()
	}
}

// 班级冲突检查
func (s *Scheduler) checkClassConflicts(gene ScheduleGene, tracker *ResourceTracker, recorder *ConflictRecorder) {
	schedule := s.findScheduleByID(gene.ScheduleID)
	classes := s.parseTeachingClasses(schedule.TeachingClass)

	for _, class := range classes {
		for _, slot := range gene.TimeSlots {

			periods := s.GetAbsolutePeriods(slot)

			// 创建教师时间槽唯一标识
			tracker.mu.Lock()
			if tracker.classroomOccupancy[class.ID] == nil {
				tracker.classroomOccupancy[class.ID] = make(map[int]bool)
			}
			for _, p := range periods {
				if tracker.classroomOccupancy[class.ID][p] {
					recorder.RecordConflict(int(gene.ScheduleID), Conflict{
						Type:    "班级冲突",
						Details: fmt.Sprintf("班级 %s 在第%d节已被占用", class.ID, p),
					})
				} else {
					tracker.classroomOccupancy[class.ID][p] = true
				}
			}
			tracker.mu.Unlock()
		}
	}
}

// 时间有效性检查
func (s *Scheduler) checkTimeValidity(gene ScheduleGene, recorder *ConflictRecorder) {
	for _, slot := range gene.TimeSlots {
		if slot.StartPeriod < 1 || slot.StartPeriod+slot.Duration-1 > s.Config.MaxPeriodsPerDay {
			recorder.RecordConflict(int(gene.ScheduleID), Conflict{
				Type:    "时间无效",
				Details: fmt.Sprintf("无效时间段：%d-%d节", slot.StartPeriod, slot.StartPeriod+slot.Duration),
			})
		}
	}
}

// 全局约束检查
func (s *Scheduler) checkGlobalConstraints(tracker *ResourceTracker, recorder *ConflictRecorder) {
	// 示例：检查教师周课时限制
	// 需要实现具体逻辑遍历教师课时统计
}

// 生成可读性报告
func GenerateReport(validation ScheduleValidation) string {
	builder := &strings.Builder{}

	// 汇总统计
	fmt.Fprintf(builder, "排课验证结果：\n")
	fmt.Fprintf(builder, "%-20s %d\n", "总排课数量", len(validation.FailedGenes)+validation.SuccessCount)
	fmt.Fprintf(builder, "%-20s %t\n", "整体有效性", validation.IsValid)
	fmt.Fprintf(builder, "\n冲突统计：\n")
	for t, c := range validation.ConflictStats {
		fmt.Fprintf(builder, "%-15s %d\n", t+":", c)
	}

	// 详细失败列表
	if len(validation.FailedGenes) > 0 {
		fmt.Fprintf(builder, "\n失败详情：\n")
		for i, fg := range validation.FailedGenes {
			fmt.Fprintf(builder, "%d. 教学班ID：%d\n", i+1, fg.ScheduleID)
			fmt.Fprintf(builder, "   冲突原因：\n")
			for _, r := range fg.ConflictReasons {
				fmt.Fprintf(builder, "   - %s\n", r)
			}
			fmt.Fprintf(builder, "   建议解决方案：%s\n", fg.SuggestedFix)
		}
	}

	// 智能建议
	if !validation.IsValid {
		builder.WriteString("\n全局建议：\n")
		if validation.ConflictStats["教室冲突"] > 0 {
			builder.WriteString(" - 考虑增加特殊时段（如晚上）的教室开放\n")
		}
		if validation.ConflictStats["教师冲突"] > 0 {
			builder.WriteString(" - 检查教师跨校区授课的可行性\n")
		}
	}

	return builder.String()
}

// 根据冲突生成解决建议
func (s *Scheduler) generateSuggestions(gene ScheduleGene, conflicts []Conflict) string {
	var suggestions []string

	for _, c := range conflicts {
		switch c.Type {
		case "教室冲突":
			// 寻找替代教室建议
			altRooms := s.findAlternativeClassrooms(gene)
			if len(altRooms) > 0 {
				suggestions = append(suggestions, fmt.Sprintf(
					"可尝试更换教室至：%v", altRooms))
			} else {
				suggestions = append(suggestions, "建议调整排课时间")
			}

		case "教师冲突":
			// 检查教师其他可用时间
			availableSlots := s.findTeacherAvailability(gene.TeacherID)
			suggestions = append(suggestions, fmt.Sprintf(
				"该教师其他可用时段：%v", availableSlots))

		case "班级冲突":
			suggestions = append(suggestions,
				"考虑拆分班级或调整课程时间")
		}
	}

	// 去重建议
	return strings.Join(unique(suggestions), "；")
}

// 示例：寻找替代教室
func (s *Scheduler) findAlternativeClassrooms(gene ScheduleGene) []string {
	original := s.ClassroomMap[gene.ClassroomID]
	var candidates []string

	for _, c := range s.Classrooms {
		if c.ID == gene.ClassroomID || c.Type != original.Type {
			continue
		}

		// 检查容量是否足够
		schedule := s.findScheduleByID(gene.ScheduleID)
		if c.Capacity >= int(schedule.TeachingClassSize) {
			candidates = append(candidates, c.ID)
		}
	}
	return candidates
}

func NewResourceTracker() *ResourceTracker {
	return &ResourceTracker{
		classroomOccupancy: make(map[string]map[int]bool),
		teacherOccupancy:   make(map[string]map[int]bool),
		classOccupancy:     make(map[string]map[int]bool),
	}
}

//func (rt *ResourceTracker) IsClassroomOccupied(key ClassroomSlotKey, duration int) bool {
//	rt.mu.Lock()
//	defer rt.mu.Unlock()
//	if maxDur, exists := rt.classroomOccupancy[key]; exists {
//		return key.Period < maxDur
//	}
//	return false
//}

//func (rt *ResourceTracker) MarkClassroom(key ClassroomSlotKey, duration int) {
//	rt.mu.Lock()
//	defer rt.mu.Unlock()
//	end := key.Period + duration
//	if current, exists := rt.classroomOccupancy[key]; exists {
//		if end > current {
//			rt.classroomOccupancy[key] = end
//		}
//	} else {
//		rt.classroomOccupancy[key] = end
//	}
//}

func NewConflictRecorder() *ConflictRecorder {
	return &ConflictRecorder{
		report: make(map[int][]Conflict),
	}
}

func (cr *ConflictRecorder) RecordConflict(scheduleID int, conflict Conflict) {
	cr.mu.Lock()
	defer cr.mu.Unlock()
	conflict.Timestamp = time.Now()
	cr.report[scheduleID] = append(cr.report[scheduleID], conflict)
}

func (cr *ConflictRecorder) GetConflicts(scheduleID int) []Conflict {
	cr.mu.Lock()
	defer cr.mu.Unlock()
	return cr.report[scheduleID]
}

func recordInvalidGene(validation *ScheduleValidation, gene ScheduleGene, reason string) {
	validation.IsValid = false
	validation.FailedGenes = append(validation.FailedGenes, FailedGene{
		ScheduleID:      int(gene.ScheduleID),
		ConflictReasons: []string{reason},
	})
}

func buildFailedGene(gene ScheduleGene, conflicts []Conflict) FailedGene {
	reasons := make([]string, 0, len(conflicts))
	for _, c := range conflicts {
		reasons = append(reasons, c.Details)
	}
	return FailedGene{
		ScheduleID:      int(gene.ScheduleID),
		ClassroomID:     gene.ClassroomID,
		TeacherID:       gene.TeacherID,
		TimeSlots:       gene.TimeSlots,
		ConflictReasons: reasons,
		SuggestedFix:    "自动生成建议逻辑需完善",
	}
}

// 去重辅助函数
func unique(input []string) []string {
	keys := make(map[string]bool)
	result := make([]string, 0)
	for _, item := range input {
		if _, value := keys[item]; !value {
			keys[item] = true
			result = append(result, item)
		}
	}
	return result
}

// 查找教师可用时间
func (s *Scheduler) findTeacherAvailability(teacherID string) []models.TimeSlot {
	// 实现需要结合具体排课数据
	return []models.TimeSlot{}
}

// 计算学期内的绝对节次编号
func (s *Scheduler) GetAbsolutePeriods(ts models.TimeSlot) []int {
	periods := make([]int, 0)

	// 获取配置参数
	maxPeriodsPerDay := s.Config.MaxPeriodsPerDay // 每天总节次（如12）
	daysPerWeek := 5                              // 每周5天（周一到周五）

	// 遍历所有周次
	for _, week := range ts.WeekNumbers {
		// 计算周偏移量（周数从1开始）
		weekOffset := (week - 1) * daysPerWeek * maxPeriodsPerDay

		// 计算天偏移量（星期几从1开始）
		dayOffset := (ts.Weekday - 1) * maxPeriodsPerDay

		// 处理持续节次（如从第4节开始，持续2节）
		for p := ts.StartPeriod; p < ts.StartPeriod+ts.Duration; p++ {
			// 检查节次有效性
			if p < 1 || p > maxPeriodsPerDay {
				log.Printf("无效节次：周%d 星期%d 第%d节", week, ts.Weekday, p)
				continue
			}

			// 生成绝对节次编号
			absolutePeriod := weekOffset + dayOffset + p
			periods = append(periods, absolutePeriod)
		}
	}

	return periods
}

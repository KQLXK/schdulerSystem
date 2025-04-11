package algorithm

import (
	"fmt"
	"log"
	"math/rand"
	"runtime"
	"schedule/commen/config"
	"schedule/dto"
	"schedule/models"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	sameRoomScore       = 5.0  // 同一教室
	sameFloorScore      = 3.0  // 同楼层不同教室
	sameBuildingScore   = 1.0  // 同楼栋不同层
	diffBuildingPenalty = -2.0 // 不同楼栋
)

//TODO
//软约束：
//1.排课优先级
//2.班级排课教室尽量集中、教师排课教室尽量集中，同一课程使用相同教室
//3.教师每天、每周排课的最大节数，教师上午、下午的最大节次
//硬约束：
//1.排课的班级人数须小于教室的容量
//可选约束：
//1.体育课安排是否只能在下午，体育课后是否安排课程
//2.学校晚上是否上课
//3.实验课程是否只能安排在晚上
//4.多学时类型（理论、实验、上机）的课程学时是否连续排

// 排课基因（教学班安排）
type ScheduleGene struct {
	ScheduleID  int64
	ClassroomID string
	TeacherID   string
	TimeSlots   []models.TimeSlot //一节课可能有多个时间段
}

// 染色体（排课方案）
type Chromosome []ScheduleGene

// 适应度评估结果
type Fitness struct {
	HardViolations int
	SoftScores     int
}

// 遗传算法参数
type GAParams struct {
	PopulationSize int
	CrossoverRate  float64
	MutationRate   float64
	MaxGenerations int
	TournamentSize int
	ElitismCount   int // 保留的精英数量
}

// 排课系统
type Scheduler struct {
	Schedules    []models.Schedule
	scheduleMap  map[int64]models.Schedule
	Classrooms   []models.Classroom
	ClassroomMap map[string]models.Classroom
	Teachers     []models.Teacher
	TeacherMap   map[string]models.Teacher
	Classes      []models.Class
	ClassMap     map[string]models.Class
	Config       dto.Config
	GAParams     GAParams
	Rng          *rand.Rand
}

// 辅助结构用于返回最佳染色体及其相关信息
type BestChromosome struct {
	Chromosome Chromosome
	Fitness    Fitness
	Index      int
}

// 辅助结构用于统计教师的排课信息
type TeacherScheduleStats struct {
	Daily      map[int]int // 每天总课时
	Weekly     int         // 周总课时
	Morning    map[int]int // 每天上午课时
	Afternoon  map[int]int // 每天下午课时
	Classrooms map[string]bool
}

// 辅助结构，周次区间结构
type WeekBlock struct {
	StartWeek    int // 起始周
	EndWeek      int // 结束周
	HoursPerWeek int // 每周学时
}

func NewScheduler(conf *dto.Config, ScheduleList []int) *Scheduler {
	// 1. 获取基础数据
	classroomDao := models.NewClassroomDao()
	classrooms, err := classroomDao.GetAllClassrooms()
	if err != nil {
		log.Fatalf("Failed to load classrooms: %v", err)
	}

	scheduleDao := models.NewScheduleDao()
	schedules := make([]models.Schedule, 0)
	if len(ScheduleList) == 0 {
		schedules, err = scheduleDao.GetAllSchedules()
		if err != nil {
			log.Fatalf("Failed to load schedules: %v", err)
		}
	} else {
		for _, s := range ScheduleList {
			schedule, err := scheduleDao.GetScheduleByID(s)
			if err != nil {
				log.Fatalf("Failed to load schedules: %v", err)
			}
			schedules = append(schedules, *schedule)
		}
	}
	// 按排课优先级排序
	sort.Slice(schedules, func(i, j int) bool {
		return schedules[i].SchedulingPriority > schedules[j].SchedulingPriority
	})

	teachers, err := models.GetAllTeachers()
	if err != nil {
		log.Fatalf("Failed to load teachers: %v", err)
	}

	classDao := models.NewClassDaoInstance()
	classes, err := classDao.GetAllClasses()
	if err != nil {
		log.Fatalf("Failed to load classes: %v", err)
	}

	// 2. 创建教室快速查询映射
	classroomMap := make(map[string]models.Classroom)
	for _, c := range classrooms {
		classroomMap[c.ID] = c
	}

	scheduleMap := make(map[int64]models.Schedule)
	for _, c := range schedules {
		scheduleMap[c.ID] = c
	}

	//帮我写一下
	teacherMap := make(map[string]models.Teacher)
	for _, c := range teachers {
		teacherMap[c.ID] = c
	}

	classMap := make(map[string]models.Class)
	for _, c := range classes {
		classMap[c.Name] = c
	}

	// 3. 加载配置
	gaConfig := config.GetConfig().GA

	// 4. 初始化随机数生成器
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	// 5. 构建Scheduler实例
	return &Scheduler{
		Schedules:    schedules,
		Classrooms:   classrooms,
		Teachers:     teachers,
		Classes:      classes,
		ClassroomMap: classroomMap, // 新增的快速查询映射
		TeacherMap:   teacherMap,
		ClassMap:     classMap,
		scheduleMap:  scheduleMap,
		Rng:          rng,
		Config:       *conf,
		GAParams: GAParams{
			PopulationSize: gaConfig.PopulationSize,
			CrossoverRate:  gaConfig.CrossoverRate,
			MutationRate:   gaConfig.MutationRate,
			MaxGenerations: gaConfig.MaxGenerations,
			TournamentSize: gaConfig.TournamentSize,
			ElitismCount:   gaConfig.ElitismCount,
		},
	}
}

//func GA(conf *Config) {
//	scheduler := NewScheduler(conf)
//	chromosome := scheduler.Run()
//	scheduler.ValidateSchedule(chromosome)
//	// 结果处理
//	//for _, gene := range chromosome {
//	//	fmt.Printf("ScheduleID: %d, Classroom: %s, TimeSlots: %v\n",
//	//		gene.ScheduleID,
//	//		gene.ClassroomID,
//	//		gene.TimeSlots)
//	//}
//
//	// 可以添加数据库存储逻辑
//	// saveScheduleResults(chromosome)
//}

// 主算法流程
func (s *Scheduler) Run() Chromosome {
	population := s.InitializePopulation()
	bestFitness := Fitness{HardViolations: 1e9, SoftScores: -1e9}

	for gen := 0; gen < s.GAParams.MaxGenerations; gen++ {
		// 计算适应度
		fitnesses := make([]Fitness, len(population))
		fitnesses = s.CalculateAllFitness(population)

		// 找到最佳个体
		currentBest := s.findBestChromosome(population, fitnesses)
		fmt.Printf("Generation %d: Best fitness %+v\n", gen, bestFitness)
		if s.compareFitness(currentBest.Fitness, bestFitness) {
			bestFitness = currentBest.Fitness
		}

		// 生成新一代
		newPopulation := make([]Chromosome, 0, s.GAParams.PopulationSize)

		// 精英保留
		for i := 0; i < s.GAParams.ElitismCount; i++ {
			newPopulation = append(newPopulation, population[currentBest.Index])
		}

		// 填充剩余个体
		for len(newPopulation) < s.GAParams.PopulationSize {
			parent1 := s.TournamentSelection(population, fitnesses)
			parent2 := s.TournamentSelection(population, fitnesses)
			child1, child2 := s.Crossover(parent1, parent2)
			newPopulation = append(newPopulation, s.Mutate(child1))
			if len(newPopulation) < s.GAParams.PopulationSize {
				newPopulation = append(newPopulation, s.Mutate(child2))
			}
		}

		population = newPopulation
	}

	return s.findBestChromosome(population, nil).Chromosome
}

// 初始化种群
func (s *Scheduler) InitializePopulation() []Chromosome {
	// 按优先级降序排序
	sortedSchedules := make([]models.Schedule, len(s.Schedules))
	copy(sortedSchedules, s.Schedules)
	//sort.Slice(sortedSchedules, func(i, j int) bool {
	//	return sortedSchedules[i].SchedulingPriority > sortedSchedules[j].SchedulingPriority
	//})

	population := make([]Chromosome, s.GAParams.PopulationSize)
	for i := 0; i < s.GAParams.PopulationSize; i++ {
		chromosome := make(Chromosome, len(sortedSchedules))
		// 按优先级顺序生成基因
		for idx, schedule := range sortedSchedules {
			chromosome[idx] = s.GenerateRandomGene(&schedule)
		}
		population[i] = chromosome
	}
	return population
}

// 随机生成教学班基因（需优化约束）
func (s *Scheduler) GenerateRandomGene(schedule *models.Schedule) ScheduleGene {
	// 步骤1: 选择符合类型的教室
	var classroom models.Classroom
	for {
		classroom = s.Classrooms[s.Rng.Intn(len(s.Classrooms))]
		if classroom.Type == schedule.SpecifiedClassroomType && classroom.Campus == schedule.OpeningCampus {
			break
		}
	}

	// 步骤2: 解析开课周次配置
	var timeSlots []models.TimeSlot
	if weekBlocks, err := s.parseWeekBlocks(schedule.OpeningWeekHours); err == nil {
		// 步骤3: 为每个周次块生成时间段
		for _, block := range weekBlocks {
			// 生成该周次块的所有周次
			weeks := make([]int, 0, block.EndWeek-block.StartWeek+1)
			for week := block.StartWeek; week <= block.EndWeek; week++ {
				weeks = append(weeks, week)
			}

			// 生成随机时间段（确保同一周次块的时间连续）
			timeSlot := models.TimeSlot{
				WeekNumbers: weeks,
				Weekday:     s.Rng.Intn(5) + 1,                             // 1-5（周一到周五）
				StartPeriod: s.Rng.Intn(s.Config.MaxPeriodsPerDay/2)*2 + 1, // 1,3,5,7,9,11 节开始
				Duration:    int(schedule.ContinuousPeriods),               // 持续时间=每周学时数
			}

			// 调整时间段确保不超过最大节次（假设每天12节课）
			if timeSlot.StartPeriod+timeSlot.Duration > 12 {
				timeSlot.StartPeriod = 12 - timeSlot.Duration
				if timeSlot.StartPeriod < 1 {
					timeSlot.StartPeriod = 1
					timeSlot.Duration = 11 // 最大持续11节（极端情况）
				}
			}

			timeSlots = append(timeSlots, timeSlot)
		}
	} else {
		// 错误处理：使用默认时间段（示例）
		timeSlots = append(timeSlots, models.TimeSlot{
			WeekNumbers: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16},
			Weekday:     s.Rng.Intn(5) + 1,
			StartPeriod: s.Rng.Intn(4)*2 + 1,
			Duration:    2,
		})
	}

	return ScheduleGene{
		ScheduleID:  schedule.ID,
		ClassroomID: classroom.ID,
		TeacherID:   schedule.TeacherID,
		TimeSlots:   timeSlots,
	}
}

func (s *Scheduler) CalculateAllFitness(population []Chromosome) []Fitness {
	result := make([]Fitness, len(population))
	var wg sync.WaitGroup
	sem := make(chan struct{}, runtime.NumCPU()) // 并发控制

	for i := range population {
		wg.Add(1)
		sem <- struct{}{}
		go func(idx int) {
			defer wg.Done()
			result[idx] = s.CalculateFitness(population[idx])
			<-sem
		}(i)
	}
	wg.Wait()
	return result
}

func (s *Scheduler) CalculateFitness(chromosome Chromosome) Fitness {
	// teacherClassrooms,courseClassrooms 初始化
	fitness := Fitness{}

	// 初始化所有跟踪数据结构
	classroomOccupancy := make(map[string]map[int]map[int]bool)
	teacherOccupancy := make(map[string]map[int]map[int]bool)
	classOccupancy := make(map[string]map[int]map[int]bool)

	// 用于软约束跟踪
	teacherStats := make(map[string]*TeacherScheduleStats)
	classClassrooms := make(map[string]map[string]bool)
	teacherClassrooms := make(map[string]map[string]bool)
	courseClassrooms := make(map[string]map[string]bool)
	multiSessionTracking := make(map[string][]models.TimeSlot)

	// 第一遍遍历：收集基础数据和硬约束
	for _, gene := range chromosome {
		schedule := s.scheduleMap[gene.ScheduleID]
		classes := s.parseTeachingClasses(schedule.TeachingClass)
		courseID := schedule.CourseID

		//=== 教室使用记录 ===//
		//初始化教师和教师跟踪
		if teacherClassrooms[schedule.TeacherID] == nil {
			teacherClassrooms[schedule.TeacherID] = make(map[string]bool)
		}
		teacherClassrooms[schedule.TeacherID][gene.ClassroomID] = true

		//初始化班级教室跟踪
		for _, class := range classes {
			if classClassrooms[class.ID] == nil {
				classClassrooms[class.ID] = make(map[string]bool)
			}
			classClassrooms[class.ID][gene.ClassroomID] = true
		}

		// 初始化课程教室跟踪
		if courseClassrooms[courseID] == nil {
			courseClassrooms[courseID] = make(map[string]bool)
		}
		courseClassrooms[courseID][gene.ClassroomID] = true

		// 初始化教师统计
		if teacherStats[gene.TeacherID] == nil {
			teacherStats[gene.TeacherID] = &TeacherScheduleStats{
				Daily:      make(map[int]int),
				Weekly:     0,
				Morning:    make(map[int]int),
				Afternoon:  make(map[int]int),
				Classrooms: make(map[string]bool),
			}
		}

		// 处理每个时间段
		for _, ts := range gene.TimeSlots {
			//=== 硬约束检查 ===//
			// 1. 资源冲突检查
			if s.isResourceOccupied(gene.ClassroomID, ts, classroomOccupancy) {
				fitness.HardViolations++
			}
			if s.isResourceOccupied(gene.TeacherID, ts, teacherOccupancy) {
				fitness.HardViolations++
			}
			for _, class := range classes {
				if s.isResourceOccupied(class.ID, ts, classOccupancy) {
					fitness.HardViolations++
				}
			}

			// 2. 教室容量检查
			classroom := s.ClassroomMap[gene.ClassroomID]
			if classroom.Capacity < int(schedule.TeachingClassSize) {
				fitness.HardViolations += 3 // 严重违规
			}

			// 3. 时间有效性检查
			if ts.StartPeriod < 1 || (ts.StartPeriod+ts.Duration-1) > s.Config.MaxPeriodsPerDay {
				fitness.HardViolations += 2
			}

			//=== 记录占用 ===//
			s.markResourceOccupancy(gene.ClassroomID, ts, classroomOccupancy)
			s.markResourceOccupancy(gene.TeacherID, ts, teacherOccupancy)
			for _, class := range classes {
				s.markResourceOccupancy(class.ID, ts, classOccupancy)
			}

			//=== 收集统计信息 ===//
			// 教师课时统计
			teacherStats[gene.TeacherID].Daily[ts.Weekday] += ts.Duration
			teacherStats[gene.TeacherID].Weekly += ts.Duration
			for p := ts.StartPeriod; p < ts.StartPeriod+ts.Duration; p++ {
				if p <= s.Config.MorningPeriodEnd {
					teacherStats[gene.TeacherID].Morning[ts.Weekday]++
				} else {
					teacherStats[gene.TeacherID].Afternoon[ts.Weekday]++
				}
			}

			// 多学时课程跟踪
			if s.Config.MultiSessionConsecutive {
				multiSessionTracking[courseID] = append(multiSessionTracking[courseID], ts)
			}
		}
	}

	//=== 软约束检查 ===//
	// 1. 教室集中度
	s.checkClassroomConcentration(&fitness, classClassrooms, teacherClassrooms, courseClassrooms)

	// 2. 教师课时限制
	s.checkTeacherConstraints(&fitness, teacherStats)

	// 3. 课程类型约束
	//s.checkCourseTypeConstraints(&fitness, chromosome)

	// 4. 多学时连续性检查
	//if s.Config.MultiSessionConsecutive {
	//	s.checkConsecutiveSessions(&fitness, multiSessionTracking)
	//}

	// 5. 优先级加分
	// todo: 处理优先级
	for _, gene := range chromosome {
		schedule := s.scheduleMap[gene.ScheduleID]
		fitness.SoftScores += int(schedule.SchedulingPriority) * 10
	}

	return fitness
}

/* ------------------------------------------ */
// 辅助函数

// 解析 opening_week_hours 字符串（示例："5-8:2,13-16:2"）
func (s *Scheduler) parseWeekBlocks(openingHours string) ([]WeekBlock, error) {
	blocks := strings.Split(openingHours, ",")
	var weekBlocks []WeekBlock

	for _, block := range blocks {
		// 分割周次和学时部分（如 "5-8:2"）
		parts := strings.Split(block, ":")
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid block format: %s", block)
		}

		// 解析周次范围（如 "5-8"）
		weeks := strings.Split(parts[0], "-")
		if len(weeks) != 2 {
			return nil, fmt.Errorf("invalid weeks format: %s", parts[0])
		}

		startWeek, err := strconv.Atoi(weeks[0])
		if err != nil {
			return nil, err
		}
		endWeek, err := strconv.Atoi(weeks[1])
		if err != nil {
			return nil, err
		}

		// 解析每周学时
		hours, err := strconv.Atoi(parts[1])
		if err != nil {
			return nil, err
		}

		weekBlocks = append(weekBlocks, WeekBlock{
			StartWeek:    startWeek,
			EndWeek:      endWeek,
			HoursPerWeek: hours,
		})
	}
	return weekBlocks, nil
}

// 检查教室集中度
func (s *Scheduler) checkClassroomConcentration(fitness *Fitness, classClassrooms, teacherClassrooms, courseClassrooms map[string]map[string]bool) {
	// 计算三类集中度并加权
	classScore := s.calculateConcentrationScore(classClassrooms, 1.0)     // 班级权重
	teacherScore := s.calculateConcentrationScore(teacherClassrooms, 0.8) // 教师权重
	courseScore := s.calculateConcentrationScore(courseClassrooms, 0.5)   // 课程权重

	totalScore := classScore + teacherScore + courseScore
	fitness.SoftScores += int(totalScore) // 假设SoftScores是总得分
}

// 通用集中度计算逻辑
func (s *Scheduler) calculateConcentrationScore(entityMap map[string]map[string]bool, weight float64) float64 {
	totalScore := 0.0

	for _, classrooms := range entityMap {
		// 跳过空实体
		if len(classrooms) == 0 {
			continue
		}

		// 统计位置分布
		locationStats := struct {
			buildingCounter map[string]int
			floorCounter    map[string]int // building-floor为key
			roomCounter     map[string]int
		}{
			buildingCounter: make(map[string]int),
			floorCounter:    make(map[string]int),
			roomCounter:     make(map[string]int),
		}
		// 收集统计数据
		for classroomID := range classrooms {
			// 从缓存获取教室信息
			classroom, exists := s.ClassroomMap[classroomID]
			if !exists {
				continue
			}
			// 更新房间计数
			locationStats.roomCounter[classroomID]++
			// 更新楼层计数
			floorKey := fmt.Sprintf("%s-%s", classroom.Building, classroom.Floor)
			locationStats.floorCounter[floorKey]++
			// 更新楼栋计数
			locationStats.buildingCounter[classroom.Building]++
		}
		// 计算当前实体得分
		entityScore := 0.0
		// 1. 同一教室得分
		for _, count := range locationStats.roomCounter {
			if count > 1 {
				entityScore += sameRoomScore * float64(count-1)
			}
		}
		// 2. 同楼层不同教室
		for _, count := range locationStats.floorCounter {
			if count > 1 {
				entityScore += sameFloorScore * float64(count-1)
			}
		}
		// 3. 同楼栋不同层
		for building, bCount := range locationStats.buildingCounter {
			// 计算本楼栋的总楼层数
			floorCount := 0
			for fKey := range locationStats.floorCounter {
				if strings.HasPrefix(fKey, building+"-") {
					floorCount++
				}
			}
			if bCount > floorCount {
				entityScore += sameBuildingScore * float64(bCount-floorCount)
			}
		}
		// 4. 跨楼栋惩罚
		if len(locationStats.buildingCounter) > 1 {
			entityScore += diffBuildingPenalty * float64(len(locationStats.buildingCounter)-1)
		}
		// 加权后加入总分
		totalScore += entityScore * weight
	}

	return totalScore
}

// 检查教师约束
func (s *Scheduler) checkTeacherConstraints(fitness *Fitness, stats map[string]*TeacherScheduleStats) {
	for _, stat := range stats {
		// 周总课时限制
		if stat.Weekly > s.Config.TeacherMaxWeeklyPeriods {
			fitness.SoftScores -= (stat.Weekly - s.Config.TeacherMaxWeeklyPeriods) * 5
		}

		// 每日限制
		for weekday, count := range stat.Daily {
			// 日总课时
			if count > s.Config.TeacherMaxDailyPeriods {
				fitness.SoftScores -= (count - s.Config.TeacherMaxDailyPeriods) * 10
			}

			// 上午课时
			if stat.Morning[weekday] > s.Config.TeacherMaxMorningPeriods {
				fitness.SoftScores -= (stat.Morning[weekday] - s.Config.TeacherMaxMorningPeriods) * 5
			}

			// 下午课时
			if stat.Afternoon[weekday] > s.Config.TeacherMaxAfternoonPeriods {
				fitness.SoftScores -= (stat.Afternoon[weekday] - s.Config.TeacherMaxAfternoonPeriods) * 5
			}
		}
	}
}

// 检查课程类型约束
// todo:体育，实验等类型的区分
func (s *Scheduler) checkCourseTypeConstraints(fitness *Fitness, chromosome Chromosome) {
	for _, gene := range chromosome {
		schedule := s.scheduleMap[gene.ScheduleID]
		courseType := schedule.Course.Type

		for _, ts := range gene.TimeSlots {
			// 体育课必须在下午
			if courseType == "体育" && s.Config.SportsAfternoonOnly {
				if ts.StartPeriod < s.Config.AfternoonStartPeriod {
					fitness.SoftScores -= 20
				}
			}

			// 实验课在晚上
			if courseType == "C类(纯实践课)" && s.Config.LabAtNightOnly {
				if ts.StartPeriod < s.Config.NightStartPeriod {
					fitness.HardViolations++ // 作为硬约束
				}
			}

			// 晚上课程检查
			if !s.Config.NightClassesAllowed && ts.StartPeriod >= s.Config.NightStartPeriod {
				fitness.HardViolations++
			}
		}
	}
}

// 检查连续排课
// todo: 实验课，上机课是否连排
func (s *Scheduler) checkConsecutiveSessions(fitness *Fitness, sessions map[string][]models.TimeSlot) {
	for _, slots := range sessions {
		// 按时间排序
		sort.Slice(slots, func(i, j int) bool {
			if slots[i].Weekday == slots[j].Weekday {
				return slots[i].StartPeriod < slots[j].StartPeriod
			}
			return slots[i].Weekday < slots[j].Weekday
		})

		// 检查连续性
		for i := 1; i < len(slots); i++ {
			prev := slots[i-1]
			current := slots[i]

			// 同一天检查
			if prev.Weekday == current.Weekday {
				expectedStart := prev.StartPeriod + prev.Duration
				if current.StartPeriod != expectedStart {
					fitness.SoftScores -= 15
				}
			} else {
				// 跨天检查（允许隔天但必须从第一节开始）
				if current.Weekday != prev.Weekday+1 || current.StartPeriod != 1 {
					fitness.SoftScores -= 20
				}
			}
		}
	}
}

// 检查每个时间段是否有教室老师或班级冲突
func (s *Scheduler) calculateResource(fitness *Fitness, chromosome *Chromosome) {
	// 记录各资源的时间占用
	classroomOccupancy := make(map[string]map[int]map[int]bool)
	teacherOccupancy := make(map[string]map[int]map[int]bool)
	classOccupancy := make(map[string]map[int]map[int]bool)

	// 遍历每个基因（教学班）
	for _, gene := range *chromosome {
		schedule := s.scheduleMap[gene.ScheduleID]
		classes := s.parseTeachingClasses(schedule.TeachingClass)

		// 检查每个时间段
		for _, ts := range gene.TimeSlots {
			// 硬约束检查
			// 1. 教室冲突
			if s.isResourceOccupied(gene.ClassroomID, ts, classroomOccupancy) {
				fitness.HardViolations++
			}
			// 2. 教师冲突
			if s.isResourceOccupied(gene.TeacherID, ts, teacherOccupancy) {
				fitness.HardViolations++
			}
			// 3. 班级冲突
			for _, class := range classes {
				if s.isResourceOccupied(class.ID, ts, classOccupancy) {
					fitness.HardViolations++
				}
			}
			// 记录占用
			s.markResourceOccupancy(gene.ClassroomID, ts, classroomOccupancy)
			s.markResourceOccupancy(gene.TeacherID, ts, teacherOccupancy)
			for _, class := range classes {
				s.markResourceOccupancy(class.ID, ts, classOccupancy)
			}
			// 示例：教室容量是否足够
			classroom := s.ClassroomMap[gene.ClassroomID]
			if classroom.Capacity < int(schedule.TeachingClassSize) {
				fitness.HardViolations += 1
			}
		}
	}
}

// 辅助函数：检查资源占用
func (s *Scheduler) isResourceOccupied(resourceID string, ts models.TimeSlot, occupancy map[string]map[int]map[int]bool) bool {
	if _, ok := occupancy[resourceID]; !ok {
		return false
	}
	for p := ts.StartPeriod; p < ts.StartPeriod+ts.Duration; p++ {
		if p > s.Config.MaxPeriodsPerDay {
			return true // 超出时间范围
		}
		if _, ok := occupancy[resourceID][ts.Weekday]; ok {
			if occupancy[resourceID][ts.Weekday][p] {
				return true
			}
		}
	}
	return false
}

// 标记资源占用
func (s *Scheduler) markResourceOccupancy(resourceID string, ts models.TimeSlot, occupancy map[string]map[int]map[int]bool) {
	if _, ok := occupancy[resourceID]; !ok {
		occupancy[resourceID] = make(map[int]map[int]bool)
	}
	if _, ok := occupancy[resourceID][ts.Weekday]; !ok {
		occupancy[resourceID][ts.Weekday] = make(map[int]bool)
	}
	for p := ts.StartPeriod; p < ts.StartPeriod+ts.Duration; p++ {
		if p <= s.Config.MaxPeriodsPerDay {
			occupancy[resourceID][ts.Weekday][p] = true
		}
	}
}

// 锦标赛选择
func (s *Scheduler) TournamentSelection(population []Chromosome, fitnesses []Fitness) Chromosome {
	tournament := make([]int, s.GAParams.TournamentSize)
	for i := 0; i < s.GAParams.TournamentSize; i++ {
		tournament[i] = s.Rng.Intn(len(population))
	}
	bestIdx := tournament[0]
	for _, idx := range tournament[1:] {
		if fitnesses[idx].HardViolations < fitnesses[bestIdx].HardViolations ||
			(fitnesses[idx].HardViolations == fitnesses[bestIdx].HardViolations &&
				fitnesses[idx].SoftScores > fitnesses[bestIdx].SoftScores) {
			bestIdx = idx
		}
	}
	return population[bestIdx]
}

// 交叉操作（单点交叉）
func (s *Scheduler) Crossover(parent1, parent2 Chromosome) (Chromosome, Chromosome) {
	if s.Rng.Float64() > s.GAParams.CrossoverRate {
		return parent1, parent2
	}
	crossoverPoint := s.Rng.Intn(len(parent1))
	child1 := append(parent1[:crossoverPoint], parent2[crossoverPoint:]...)
	child2 := append(parent2[:crossoverPoint], parent1[crossoverPoint:]...)
	return child1, child2
}

// 变异操作
func (s *Scheduler) Mutate(chromosome Chromosome) Chromosome {
	for i := range chromosome {
		if s.Rng.Float64() < s.GAParams.MutationRate {
			// 随机改变时间或教室
			schedule := s.scheduleMap[chromosome[i].ScheduleID]
			newGene := s.GenerateRandomGene(&schedule)
			chromosome[i] = newGene
		}
	}
	return chromosome
}

func (s *Scheduler) findScheduleByID(id int64) *models.Schedule {
	res, _ := models.NewScheduleDao().GetScheduleByID(int(id))
	return res
}

func (s *Scheduler) findClassroomByID(id string) *models.Classroom {
	res, _ := models.NewClassroomDao().GetClassroomByID(id)
	return res
}

func (s *Scheduler) findTeacherByID(id string) *models.Teacher {
	res, _ := models.GetTeacherByID(id)
	return res
}

func (s *Scheduler) findCourseByID(id string) *models.Course {
	res, _ := models.NewCourseDao().GetCourseByID(id)
	return res
}

func (s *Scheduler) parseTeachingClasses(classes string) []models.Class {
	classlist := strings.Split(classes, ",")
	res := make([]models.Class, 0)
	for _, val := range classlist {
		class, _ := s.ClassMap[val]
		res = append(res, class)
	}
	return res
}

func (s *Scheduler) findBestChromosome(population []Chromosome, fitnesses []Fitness) BestChromosome {
	// 自动计算适应度（当传入的 fitnesses 为 nil 时）
	effectiveFitness := fitnesses
	if effectiveFitness == nil {
		effectiveFitness = make([]Fitness, len(population))
		for i := range population {
			effectiveFitness[i] = s.CalculateFitness(population[i])
		}
	}

	bestIndex := 0
	bestFit := effectiveFitness[0]

	// 遍历所有个体寻找最优解
	for i := 1; i < len(effectiveFitness); i++ {
		currentFit := effectiveFitness[i]

		// 使用 compareFitness 进行比较
		if s.compareFitness(currentFit, bestFit) {
			bestIndex = i
			bestFit = currentFit
		}
	}

	return BestChromosome{
		Chromosome: population[bestIndex],
		Fitness:    bestFit,
		Index:      bestIndex,
	}
}

func (s *Scheduler) compareFitness(a, b Fitness) bool {
	// 优先比较硬约束违反次数，次数越少越好
	if a.HardViolations < b.HardViolations {
		return true
	}
	// 硬约束违反次数相同的情况下，比较软约束得分，得分越高越好
	if a.HardViolations == b.HardViolations {
		return a.SoftScores > b.SoftScores
	}
	return false
}

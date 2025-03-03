package algorithm

import (
	"math/rand"
	"schedule/models"
	"time"
)

// Chromosome 表示一个排课方案
type Chromosome struct {
	Genes   []Gene  // 基因（课程安排）
	Fitness float64 // 适应度
}

// Gene 表示一个课程安排
type Gene struct {
	CourseID    string    // 课程ID
	TeacherID   string    // 教师ID
	ClassroomID string    // 教室ID
	ClassID     string    // 班级ID
	StartTime   time.Time // 开始时间
	EndTime     time.Time // 结束时间
}

// GeneticAlgorithm 是遗传算法的核心实现
type GeneticAlgorithm struct {
	PopulationSize int                // 种群大小
	MaxGenerations int                // 最大迭代次数
	MutationRate   float64            // 变异率
	Courses        []models.Course    // 所有课程
	Teachers       []models.Teacher   // 所有教师
	Classrooms     []models.Classroom // 所有教室
	Classes        []models.Class     // 所有班级
}

// NewGeneticAlgorithm 创建一个新的遗传算法实例
func NewGeneticAlgorithm(populationSize, maxGenerations int, mutationRate float64, courses []models.Course, teachers []models.Teacher, classrooms []models.Classroom, classes []models.Class) *GeneticAlgorithm {
	return &GeneticAlgorithm{
		PopulationSize: populationSize,
		MaxGenerations: maxGenerations,
		MutationRate:   mutationRate,
		Courses:        courses,
		Teachers:       teachers,
		Classrooms:     classrooms,
		Classes:        classes,
	}
}

// Run 运行遗传算法，返回最优的排课方案
func (ga *GeneticAlgorithm) Run() Chromosome {
	// 初始化种群
	population := ga.initializePopulation()

	// 迭代优化
	for generation := 0; generation < ga.MaxGenerations; generation++ {
		// 计算适应度
		for i := range population {
			population[i].Fitness = ga.calculateFitness(population[i])
		}

		// 选择
		newPopulation := ga.selectPopulation(population)

		// 交叉
		newPopulation = ga.crossoverPopulation(newPopulation)

		// 变异
		newPopulation = ga.mutatePopulation(newPopulation)

		// 更新种群
		population = newPopulation
	}

	// 返回最优解
	bestChromosome := population[0]
	for _, chromosome := range population {
		if chromosome.Fitness > bestChromosome.Fitness {
			bestChromosome = chromosome
		}
	}
	return bestChromosome
}

// initializePopulation 初始化种群
func (ga *GeneticAlgorithm) initializePopulation() []Chromosome {
	population := make([]Chromosome, ga.PopulationSize)
	for i := range population {
		population[i].Genes = ga.generateRandomGenes()
	}
	return population
}

// generateRandomGenes 随机生成一个基因序列
func (ga *GeneticAlgorithm) generateRandomGenes() []Gene {
	genes := make([]Gene, len(ga.Courses))
	for i, course := range ga.Courses {
		teacher := ga.Teachers[rand.Intn(len(ga.Teachers))]
		classroom := ga.Classrooms[rand.Intn(len(ga.Classrooms))]
		class := ga.Classes[rand.Intn(len(ga.Classes))]
		startTime := time.Now().Add(time.Duration(rand.Intn(24)) * time.Hour)
		endTime := startTime.Add(time.Duration(course.TotalHours) * time.Hour)

		genes[i] = Gene{
			CourseID:    course.ID,
			TeacherID:   teacher.ID,
			ClassroomID: classroom.ID,
			ClassID:     class.ID,
			StartTime:   startTime,
			EndTime:     endTime,
		}
	}
	return genes
}

// calculateFitness 计算染色体的适应度
func (ga *GeneticAlgorithm) calculateFitness(chromosome Chromosome) float64 {
	// 适应度计算逻辑（示例：冲突越少，适应度越高）
	conflicts := 0
	for i, gene1 := range chromosome.Genes {
		for j, gene2 := range chromosome.Genes {
			if i != j && gene1.TeacherID == gene2.TeacherID && gene1.StartTime.Before(gene2.EndTime) && gene1.EndTime.After(gene2.StartTime) {
				conflicts++ // 教师时间冲突
			}
			if i != j && gene1.ClassroomID == gene2.ClassroomID && gene1.StartTime.Before(gene2.EndTime) && gene1.EndTime.After(gene2.StartTime) {
				conflicts++ // 教室时间冲突
			}
		}
	}
	return 1.0 / (float64(conflicts) + 1.0) // 冲突越少，适应度越高
}

// selectPopulation 选择优质个体
func (ga *GeneticAlgorithm) selectPopulation(population []Chromosome) []Chromosome {
	newPopulation := make([]Chromosome, ga.PopulationSize)
	for i := range newPopulation {
		// 随机选择两个个体，保留适应度较高的一个
		a := population[rand.Intn(len(population))]
		b := population[rand.Intn(len(population))]
		if a.Fitness > b.Fitness {
			newPopulation[i] = a
		} else {
			newPopulation[i] = b
		}
	}
	return newPopulation
}

// crossoverPopulation 交叉生成新种群
func (ga *GeneticAlgorithm) crossoverPopulation(population []Chromosome) []Chromosome {
	newPopulation := make([]Chromosome, ga.PopulationSize)
	for i := 0; i < ga.PopulationSize; i += 2 {
		parent1 := population[i]
		parent2 := population[i+1]
		child1, child2 := ga.crossover(parent1, parent2)
		newPopulation[i] = child1
		newPopulation[i+1] = child2
	}
	return newPopulation
}

// crossover 交叉两个染色体
func (ga *GeneticAlgorithm) crossover(parent1, parent2 Chromosome) (Chromosome, Chromosome) {
	// 单点交叉
	point := rand.Intn(len(parent1.Genes))
	child1 := Chromosome{Genes: append([]Gene{}, parent1.Genes[:point]...)}
	child2 := Chromosome{Genes: append([]Gene{}, parent2.Genes[:point]...)}
	child1.Genes = append(child1.Genes, parent2.Genes[point:]...)
	child2.Genes = append(child2.Genes, parent1.Genes[point:]...)
	return child1, child2
}

// mutatePopulation 变异种群
func (ga *GeneticAlgorithm) mutatePopulation(population []Chromosome) []Chromosome {
	for i := range population {
		if rand.Float64() < ga.MutationRate {
			population[i] = ga.mutate(population[i])
		}
	}
	return population
}

// mutate 变异染色体
func (ga *GeneticAlgorithm) mutate(chromosome Chromosome) Chromosome {
	// 随机选择一个基因进行变异
	index := rand.Intn(len(chromosome.Genes))
	chromosome.Genes[index] = ga.generateRandomGenes()[0]
	return chromosome
}

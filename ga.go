package main

import (
	"fmt"
	"math/rand"
	"time"
)

// 全局常量默认参数值
const (
	IDV_NUM        int     = 3000
	MAX_NUM        int     = 100000
	CROSS_RATE     float64 = 0.4
	VARIATION_RATE float64 = 0.1
	INTERVAL       int     = 200
	FRESH          int     = 2000
)

// 初始化函数
func init() {
	rand.Seed(time.Now().UnixNano()) // 设定随机种子
}

// 染色体结构体
type Individual struct {
	grid  Board // 生成解
	vis   Board // 原题
	score int   // 适应度
}

// 填充随机解
func RandFill(grid *Board, r, c int, fixed Board) {
	vis := [10]bool{}
	for i := r; i < r+3; i++ {
		for j := c; j < c+3; j++ {
			x := fixed[i][j]
			if x != 0 {
				vis[x] = true
				grid[i][j] = x
			}
		}
	}
	for i := r; i < r+3; i++ {
		for j := c; j < c+3; j++ {
			if fixed[i][j] == 0 {
				x := rand.Intn(9) + 1
				for vis[x] {
					x = rand.Intn(9) + 1
				}
				grid[i][j] = x
				vis[x] = true
			}
		}
	}
}

// Individual构造函数
func NewIndividual(gridInput Board) Individual {
	var grid Board
	for i := 0; i < 9; i += 3 {
		for j := 0; j < 9; j += 3 {
			RandFill(&grid, i, j, gridInput)
		}
	}
	res := Individual{grid, gridInput, 0}
	res.GetFitness()
	return res
}

// 计算适应度
func (this *Individual) GetFitness() int {
	res := 100
	checkRow := func(r int) int {
		var cnt [10]int
		rs := 0
		for i := 0; i < 9; i++ {
			cnt[this.grid[r][i]]++
		}
		for i := 1; i <= 9; i++ {
			if cnt[i] > 1 {
				rs += cnt[i] - 1
			}
		}
		return rs
	}
	checkCol := func(c int) int {
		var cnt [10]int
		rs := 0
		for i := 0; i < 9; i++ {
			cnt[this.grid[i][c]]++
		}
		for i := 1; i <= 9; i++ {
			if cnt[i] > 1 {
				rs += cnt[i] - 1
			}
		}
		return rs
	}

	for i := 0; i < 9; i++ {
		res -= checkRow(i)
		res -= checkCol(i)
	}
	this.score = res
	return res
}

// 深拷贝
func (this *Individual) Copy(other Individual) {
	this.score = other.score
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			this.grid[i][j] = other.grid[i][j]
			this.vis[i][j] = other.vis[i][j]
		}
	}
}

// 显示解状态
func (this Individual) Show() {
	this.grid.Show()
}

// GA遗传模型
type GA struct {
	idvNum        int          // 种群个体数目
	maxNum        int          // 最大迭代次数
	crossRate     float64      // 交叉率
	variationRate float64      // 变异率
	bestIdv       Individual   // 最优个体
	population    []Individual // 种群
	vis           Board        // 原题
}

// GA结构体构造函数:(种群大小,迭代次数,交叉率,变异率)
func NewGAModel(par ...interface{}) GA {
	n := len(par)
	idvNum := IDV_NUM
	maxNum := MAX_NUM
	crossRate := CROSS_RATE
	variationRate := VARIATION_RATE
	if n > 0 {
		idvNum = par[0].(int)
		if idvNum&1 == 1 {
			idvNum--
		}
	}
	if n > 1 {
		maxNum = par[1].(int)
	}
	if n > 2 {
		crossRate = par[2].(float64)
	}
	if n > 3 {
		variationRate = par[3].(float64)
	}
	return GA{idvNum: idvNum, maxNum: maxNum, crossRate: crossRate, variationRate: variationRate}
}

// 导入数独习题
func (this *GA) ModelInit(gridInput Board) {
	this.vis = gridInput
	maxFit := 0
	this.population = make([]Individual, this.idvNum)
	for i := range this.population {
		this.population[i] = NewIndividual(gridInput)
		this.population[i].GetFitness()
		if maxFit < this.population[i].score {
			maxFit = this.population[i].score
			this.bestIdv = this.population[i]
		}
	}
}

// 更新最优染色体
func (this *GA) GetBestIdv() Individual {
	maxFit := 0
	for i := range this.population {
		if maxFit < this.population[i].score {
			maxFit = this.population[i].score
			this.bestIdv = this.population[i]
		}
	}
	return this.bestIdv
}

// 交叉互换
func (this *GA) Cross(idv1, idv2 Individual) (Individual, Individual) {
	var idv3, idv4 Individual
	idv3.Copy(idv1)
	idv4.Copy(idv2)
	partSwap := func(chd1, chd2 *Individual, r, c int, vis Board) {
		for i := r; i < r+3; i++ {
			for j := c; j < c+3; j++ {
				if vis[i][j] == 0 {
					chd1.grid[i][j], chd2.grid[i][j] = chd2.grid[i][j], chd1.grid[i][j]
				}
			}
		}
	}
	for i := 0; i < 9; i += 3 {
		for j := 0; j < 9; j += 3 {
			if rand.Float64() > this.crossRate {
				partSwap(&idv3, &idv4, i, j, this.vis)
			}
		}
	}
	idv3.GetFitness()
	idv4.GetFitness()
	return idv3, idv4
}

// 锦标赛选择
func (this *GA) Select() {
	rand.Shuffle(len(this.population), func(i, j int) {
		this.population[i], this.population[j] = this.population[j], this.population[i]
	})
	newPopulation := make([]Individual, this.idvNum)
	better := func(a, b Individual) Individual {
		if a.score > b.score {
			return a
		}
		return b
	}
	for i, j := 0, 0; i < len(this.population); i += 2 {
		newPopulation[j] = better(this.population[i], this.population[i+1])
		j++
	}
	this.population = newPopulation
}

// 变异操作
func (this *GA) Variate(Idx *Individual) {
	for i := 0; i < 9; i += 3 {
		for j := 0; j < 9; j += 3 {
			if rand.Float64() < this.variationRate {
				RandFill(&Idx.grid, i, j, this.vis)
			}
		}
	}
}

// 模型训练
func (this *GA) Train() {
	for t := 0; t < this.maxNum; t++ {
		// 父代母代进行染色体交配
		for i := 0; i < this.idvNum; i += 2 {
			father, mother := this.population[i], this.population[i+1]
			newIdv1, newIdv2 := this.Cross(father, mother)
			newIdv1.GetFitness()
			newIdv2.GetFitness()
			this.population = append(this.population, newIdv1, newIdv2)
		}
		// 种群选择操作优胜劣汰
		this.Select()
		// 种群变异操作
		for i := 0; i < this.idvNum; i++ {
			if rand.Float64() < this.variationRate {
				this.Variate(&(this.population[i]))
				this.population[i].GetFitness()
			}
		}
		// 获取当前种群最优染色体解
		this.GetBestIdv()
		if this.bestIdv.score == 100 {
			fmt.Printf("-------------NO.%d--------------\n", t)
			fmt.Println("Result:")
			this.bestIdv.Show()
			fmt.Println()
			return
		}
		if t%INTERVAL == 0 {
			fmt.Printf("-------------%d--------------\n", t)
			fmt.Printf("Score:%d\n", this.bestIdv.score)
			fmt.Println("Best Individual:")
			this.bestIdv.Show()
			fmt.Println()
			time.Sleep(time.Second)
		}
		if t > 0 && t%FRESH == 0 {
			if this.variationRate < 0.8 {
				this.variationRate *= 1.1
			}
			fmt.Printf("variation rate:%.2f\n", this.variationRate)
			for i := 0; i < this.idvNum; i++ {
				this.population[i] = NewIndividual(this.vis)
				this.population[i].GetFitness()
			}
		}
	}
	fmt.Println("failed to find!")
}

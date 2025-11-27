package models

import (
	"math/rand"
	"time"
)

// Direction 蛇的移动方向
type Direction int

const (
	Up Direction = iota
	Down
	Left
	Right
)

// Position 位置坐标
type Position struct {
	X int `json:"x"`
	Y int `json:"y"`
}

// Snake 蛇的结构
type Snake struct {
	Body      []Position `json:"body"`
	Direction Direction  `json:"direction"`
}

// Food 食物结构
type Food struct {
	Position Position `json:"position"`
}

// Wall 墙体结构
type Wall struct {
	Position  Position  `json:"position"`
	CreatedAt time.Time `json:"created_at"`
	Lifetime  int       // 存活时间（秒）
}

// 游戏状态常量
const (
	GameStatusRunning = "running" // 游戏运行中
	GameStatusEnded   = "ended"   // 游戏结束
)

// Game 游戏结构体
type Game struct {
	ID             string    `json:"id"`
	Snake          Snake     `json:"snake"`
	Food           Food      `json:"food"`
	Walls          []Wall    `json:"walls"`
	Width          int       `json:"width"`
	Height         int       `json:"height"`
	Status         string    `json:"status"`
	Score          int       `json:"score"`
	FoodCount      int       `json:"foodCount"`
	Time           int       `json:"time"`
	LastFoodTime   time.Time `json:"lastFoodTime"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
	LastUpdateTime time.Time `json:"lastUpdateTime"` // 上一次更新时间，用于计算时间差
	MaxWalls       int       `json:"maxWalls"`
}

// NewGame 创建一个新游戏
func NewGame(id string, width, height, maxWalls int) *Game {
	snake := Snake{
		Body: []Position{
			{X: width / 2, Y: height / 2},
			{X: width/2 - 1, Y: height / 2},
			{X: width/2 - 2, Y: height / 2},
		},
		Direction: Right,
	}

	now := time.Now()
	game := &Game{
		ID:             id,
		Snake:          snake,
		Width:          width,
		Height:         height,
		Status:         GameStatusRunning,
		Score:          0,
		FoodCount:      0,
		Time:           0,
		LastFoodTime:   now,
		CreatedAt:      now,
		UpdatedAt:      now,
		LastUpdateTime: now,
		MaxWalls:       maxWalls,
	}

	// 生成初始食物
	game.GenerateFood()

	return game
}

// GenerateFood 生成新食物
func (g *Game) GenerateFood() {
	rand.Seed(time.Now().UnixNano())

	// 找到所有可用的位置
	usedPositions := make(map[Position]bool)

	// 添加蛇的位置
	for _, pos := range g.Snake.Body {
		usedPositions[pos] = true
	}

	// 添加墙体的位置
	for _, wall := range g.Walls {
		usedPositions[wall.Position] = true
	}

	// 收集所有可用位置
	availablePositions := []Position{}
	for x := 0; x < g.Width; x++ {
		for y := 0; y < g.Height; y++ {
			pos := Position{X: x, Y: y}
			if !usedPositions[pos] {
				availablePositions = append(availablePositions, pos)
			}
		}
	}

	if len(availablePositions) > 0 {
		// 随机选择一个可用位置
		randomIndex := rand.Intn(len(availablePositions))
		g.Food.Position = availablePositions[randomIndex]
		g.LastFoodTime = time.Now()
	}
}

// GenerateWall 生成新墙体
func (g *Game) GenerateWall() {
	if len(g.Walls) >= g.MaxWalls {
		return
	}

	rand.Seed(time.Now().UnixNano())

	// 找到所有可用的位置
	usedPositions := make(map[Position]bool)

	// 添加蛇的位置
	for _, pos := range g.Snake.Body {
		usedPositions[pos] = true
	}

	// 添加食物的位置
	usedPositions[g.Food.Position] = true

	// 添加现有墙体的位置
	for _, wall := range g.Walls {
		usedPositions[wall.Position] = true
	}

	// 收集所有可用位置，同时确保墙体不会阻挡关键路径
	availablePositions := []Position{}
	for x := 0; x < g.Width; x++ {
		for y := 0; y < g.Height; y++ {
			pos := Position{X: x, Y: y}
			if !usedPositions[pos] {
				// 确保墙体不会完全阻挡蛇的移动路径
				// 1. 避免在蛇头周围形成包围圈
				isNearHead := false
				head := g.Snake.Body[0]
				// 检查是否在蛇头的2格范围内
				if abs(pos.X-head.X) <= 2 && abs(pos.Y-head.Y) <= 2 {
					isNearHead = true
				}

				// 2. 避免在边界形成无法通过的区域
				// 允许大部分位置，但确保不会完全阻挡到食物的路径
				if !isNearHead || rand.Float64() < 0.3 { // 有30%的概率即使在头部附近也允许生成
					availablePositions = append(availablePositions, pos)
				}
			}
		}
	}

	if len(availablePositions) > 0 {
		// 随机选择一个可用位置
		randomIndex := rand.Intn(len(availablePositions))
		wall := Wall{
			Position:  availablePositions[randomIndex],
			CreatedAt: time.Now(),
			Lifetime:  5 + rand.Intn(10), // 5-14秒的随机生命周期，符合墙体消失的需求
		}
		g.Walls = append(g.Walls, wall)
	}
}

// abs 计算绝对值
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// Update 更新游戏状态
func (g *Game) Update() {
	if g.Status != GameStatusRunning {
		return
	}

	// 移动蛇
	g.MoveSnake()

	// 检查碰撞
	if g.CheckCollision() {
		g.Status = GameStatusEnded
		return
	}

	// 检查是否吃到食物 - 根据返回值决定是否移除尾部
	ateFood := g.CheckFoodCollision()

	// 如果没有吃到食物，移除尾部（保持长度不变）
	if !ateFood && len(g.Snake.Body) > 0 {
		g.Snake.Body = g.Snake.Body[:len(g.Snake.Body)-1]
	}

	// 基于实际时间差更新游戏时间
	now := time.Now()
	elapsedMilliseconds := now.Sub(g.LastUpdateTime).Milliseconds()
	// 如果经过了至少1000毫秒，更新游戏时间
	if elapsedMilliseconds >= 1000 {
		g.Time += int(elapsedMilliseconds / 1000) // 转换为秒
		g.LastUpdateTime = now
	}

	// 检查是否长时间没有吃到食物（10秒规则）
	if time.Since(g.LastFoodTime).Seconds() > 10 {
		g.Status = GameStatusEnded
		return
	}

	// 更新分数：时间(秒)*1 + 豆子数量*10
	g.Score = g.Time + g.FoodCount*10

	// 更新最后更新时间
	g.UpdatedAt = now // 使用之前定义的now变量

	// 移除过期的墙体
	now = time.Now() // 重用之前定义的now变量
	var validWalls []Wall
	for _, wall := range g.Walls {
		if now.Sub(wall.CreatedAt).Seconds() < float64(wall.Lifetime) {
			validWalls = append(validWalls, wall)
		}
	}
	g.Walls = validWalls

	// 随机生成新墙体（每500ms约20%的概率，符合随时间生成的需求）
	if rand.Float64() < 0.2 {
		g.GenerateWall()
	}
}

// MoveSnake 移动蛇 - 只有吃到豆子时才增加长度
func (g *Game) MoveSnake() {
	head := g.Snake.Body[0]
	newHead := Position{X: head.X, Y: head.Y}

	switch g.Snake.Direction {
	case Up:
		newHead.Y--
	case Down:
		newHead.Y++
	case Left:
		newHead.X--
	case Right:
		newHead.X++
	}

	// 将新头部添加到蛇身
	g.Snake.Body = append([]Position{newHead}, g.Snake.Body...)

	// 默认情况下移除尾部，只有吃到豆子时才保留
	// 注意：在Update方法中会根据是否吃到豆子来决定是否保留尾部
}

// CheckCollision 检查碰撞
func (g *Game) CheckCollision() bool {
	head := g.Snake.Body[0]

	// 检查边界碰撞
	if head.X < 0 || head.X >= g.Width || head.Y < 0 || head.Y >= g.Height {
		return true
	}

	// 检查自身碰撞
	for i := 1; i < len(g.Snake.Body); i++ {
		if head.X == g.Snake.Body[i].X && head.Y == g.Snake.Body[i].Y {
			return true
		}
	}

	// 检查墙体碰撞
	for _, wall := range g.Walls {
		if head.X == wall.Position.X && head.Y == wall.Position.Y {
			return true
		}
	}

	return false
}

// CheckFoodCollision 检查是否吃到食物
func (g *Game) CheckFoodCollision() bool {
	head := g.Snake.Body[0]
	if head.X == g.Food.Position.X && head.Y == g.Food.Position.Y {
		// 吃到食物，增加分数和长度
		g.FoodCount++
		g.LastFoodTime = time.Now()

		// 生成新食物
		g.GenerateFood()
		return true
	}

	return false
}

// ChangeDirection 改变蛇的移动方向
func (g *Game) ChangeDirection(newDirection Direction) {
	// 防止180度转向
	if (g.Snake.Direction == Up && newDirection == Down) ||
		(g.Snake.Direction == Down && newDirection == Up) ||
		(g.Snake.Direction == Left && newDirection == Right) ||
		(g.Snake.Direction == Right && newDirection == Left) {
		return
	}

	if g.Status == GameStatusRunning {
		g.Snake.Direction = newDirection
	}
}

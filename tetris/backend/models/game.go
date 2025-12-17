package models

import (
	"math/rand"
	"time"
)

// BlockType 方块类型
type BlockType int

const (
	BlockTypeS BlockType = iota
	BlockTypeZ
	BlockTypeL
	BlockTypeJ
	BlockTypeI
	BlockTypeO
	BlockTypeT
)

// Position 位置结构体
type Position struct {
	Row int `json:"row"`
	Col int `json:"col"`
}

// Block 方块结构体
type Block struct {
	Type     BlockType `json:"type"`
	Shape    [][]int   `json:"shape"`
	Position Position  `json:"position"`
	Color    string    `json:"color"`
	Rotation int       `json:"rotation"`
}

// GameBoard 游戏面板结构体
type GameBoard struct {
	Width  int     `json:"width"`
	Height int     `json:"height"`
	Grid   [][]int `json:"grid"`
	Locked []bool  `json:"locked"`
}

// GameState 游戏状态枚举
type GameState int

const (
	GameStateStopped GameState = iota
	GameStateRunning
	GameStatePaused
	GameStateGameOver
)

// Game 游戏结构体
type Game struct {
	ID           string    `json:"gameId"`
	Board        GameBoard `json:"board"`
	CurrentPiece Block     `json:"currentPiece"`
	NextPiece    Block     `json:"nextPiece"`
	Score        int       `json:"score"`
	Lines        int       `json:"lines"`
	Level        int       `json:"level"`
	State        GameState `json:"state"`
	GameSpeed    int       `json:"gameSpeed"`
	LastUpdate   time.Time `json:"lastUpdate"`
	IsGameOver   bool      `json:"gameOver"`
	GameStarted  bool      `json:"gameStarted"`
	Paused       bool      `json:"paused"`
}

// BlockShapes 方块形状定义
var BlockShapes = map[BlockType][][][]int{
	BlockTypeS: {
		{{0, 1, 1}, {1, 1, 0}},   // 第一个旋转状态
		{{1, 0}, {1, 1}, {0, 1}}, // 第二个旋转状态
		{{0, 1, 1}, {1, 1, 0}},   // 第三个旋转状态
		{{1, 0}, {1, 1}, {0, 1}}, // 第四个旋转状态
	},
	BlockTypeZ: {
		{{1, 1, 0}, {0, 1, 1}},   // 第一个旋转状态
		{{0, 1}, {1, 1}, {1, 0}}, // 第二个旋转状态
		{{1, 1, 0}, {0, 1, 1}},   // 第三个旋转状态
		{{0, 1}, {1, 1}, {1, 0}}, // 第四个旋转状态
	},
	BlockTypeL: {
		{{0, 0, 1}, {1, 1, 1}},   // 第一个旋转状态
		{{1, 0}, {1, 0}, {1, 1}}, // 第二个旋转状态
		{{1, 1, 1}, {1, 0, 0}},   // 第三个旋转状态
		{{1, 1}, {0, 1}, {0, 1}}, // 第四个旋转状态
	},
	BlockTypeJ: {
		{{1, 0, 0}, {1, 1, 1}},   // 第一个旋转状态
		{{1, 1}, {1, 0}, {1, 0}}, // 第二个旋转状态
		{{1, 1, 1}, {0, 0, 1}},   // 第三个旋转状态
		{{0, 1}, {0, 1}, {1, 1}}, // 第四个旋转状态
	},
	BlockTypeI: {
		{{1, 1, 1, 1}},       // 第一个旋转状态
		{{1}, {1}, {1}, {1}}, // 第二个旋转状态
		{{1, 1, 1, 1}},       // 第三个旋转状态
		{{1}, {1}, {1}, {1}}, // 第四个旋转状态
	},
	BlockTypeO: {
		{{1, 1}, {1, 1}}, // 第一个旋转状态
	},
	BlockTypeT: {
		{{0, 1, 0}, {1, 1, 1}},   // 第一个旋转状态
		{{1, 0}, {1, 1}, {1, 0}}, // 第二个旋转状态
		{{1, 1, 1}, {0, 1, 0}},   // 第三个旋转状态
		{{0, 1}, {1, 1}, {0, 1}}, // 第四个旋转状态
	},
}

// 方块颜色定义
var BlockColors = map[BlockType]string{
	BlockTypeS: "#00FF00", // 绿色
	BlockTypeZ: "#FF0000", // 红色
	BlockTypeL: "#FFA500", // 橙色
	BlockTypeJ: "#0000FF", // 蓝色
	BlockTypeI: "#00FFFF", // 青色
	BlockTypeO: "#FFFF00", // 黄色
	BlockTypeT: "#800080", // 紫色
}

// NewBlock 创建新方块
func NewBlock(blockType BlockType) Block {
	shapes := BlockShapes[blockType]
	shape := shapes[0]
	color := BlockColors[blockType]

	return Block{
		Type:     blockType,
		Shape:    shape,
		Position: Position{Row: 0, Col: 4}, // 从顶部中间开始
		Color:    color,
		Rotation: 0,
	}
}

// NewRandomBlock 创建随机方块
func NewRandomBlock() Block {
	blockTypes := []BlockType{BlockTypeS, BlockTypeZ, BlockTypeL, BlockTypeJ, BlockTypeI, BlockTypeO, BlockTypeT}
	blockType := blockTypes[rand.Intn(len(blockTypes))]
	return NewBlock(blockType)
}

// NewGameBoard 创建新游戏面板
func NewGameBoard(width, height int) GameBoard {
	grid := make([][]int, height)
	locked := make([]bool, height)

	for i := range grid {
		grid[i] = make([]int, width)
		locked[i] = false
	}

	return GameBoard{
		Width:  width,
		Height: height,
		Grid:   grid,
		Locked: locked,
	}
}

// NewGame 创建新游戏
func NewGame(id string) *Game {
	board := NewGameBoard(10, 20) // 标准俄罗斯方块尺寸：宽10，高20
	currentPiece := NewRandomBlock()
	nextPiece := NewRandomBlock()

	return &Game{
		ID:           id,
		Board:        board,
		CurrentPiece: currentPiece,
		NextPiece:    nextPiece,
		Score:        0,
		Lines:        0,
		Level:        1,
		State:        GameStateRunning,
		GameSpeed:    1000, // 初始速度：1秒下落一格
		LastUpdate:   time.Now(),
		IsGameOver:   false,
		GameStarted:  true,
		Paused:       false,
	}
}

// RotateBlock 旋转方块
func (g *Game) RotateBlock() {
	if g.State != GameStateRunning || g.IsGameOver {
		return
	}

	shapes := BlockShapes[g.CurrentPiece.Type]
	nextRotation := (g.CurrentPiece.Rotation + 1) % len(shapes)
	nextShape := shapes[nextRotation]

	// 检查旋转是否合法
	if g.CanPlaceBlock(g.CurrentPiece.Position, nextShape) {
		g.CurrentPiece.Shape = nextShape
		g.CurrentPiece.Rotation = nextRotation
	}
}

// MoveBlockLeft 向左移动方块
func (g *Game) MoveBlockLeft() {
	if g.State != GameStateRunning || g.IsGameOver {
		return
	}

	newPosition := Position{Row: g.CurrentPiece.Position.Row, Col: g.CurrentPiece.Position.Col - 1}
	if g.CanPlaceBlock(newPosition, g.CurrentPiece.Shape) {
		g.CurrentPiece.Position = newPosition
	}
}

// MoveBlockRight 向右移动方块
func (g *Game) MoveBlockRight() {
	if g.State != GameStateRunning || g.IsGameOver {
		return
	}

	newPosition := Position{Row: g.CurrentPiece.Position.Row, Col: g.CurrentPiece.Position.Col + 1}
	if g.CanPlaceBlock(newPosition, g.CurrentPiece.Shape) {
		g.CurrentPiece.Position = newPosition
	}
}

// MoveBlockDown 向下移动方块
func (g *Game) MoveBlockDown() bool {
	if g.State != GameStateRunning || g.IsGameOver {
		return false
	}

	newPosition := Position{Row: g.CurrentPiece.Position.Row + 1, Col: g.CurrentPiece.Position.Col}
	if g.CanPlaceBlock(newPosition, g.CurrentPiece.Shape) {
		g.CurrentPiece.Position = newPosition
		return true
	}

	// 方块无法继续下移，锁定到面板
	g.LockBlock()
	return false
}

// DropBlock 快速下落方块
func (g *Game) DropBlock() {
	if g.State != GameStateRunning || g.IsGameOver {
		return
	}

	for g.MoveBlockDown() {
		// 持续下移直到无法移动
	}
}

// CanPlaceBlock 检查方块是否可以放置在指定位置
func (g *Game) CanPlaceBlock(position Position, shape [][]int) bool {
	shapeHeight := len(shape)
	if shapeHeight == 0 {
		return false
	}
	shapeWidth := len(shape[0])

	for y := 0; y < shapeHeight; y++ {
		for x := 0; x < shapeWidth; x++ {
			if shape[y][x] == 1 {
				boardCol := position.Col + x
				boardRow := position.Row + y

				// 检查边界
				if boardCol < 0 || boardCol >= g.Board.Width || boardRow >= g.Board.Height {
					return false
				}

				// 检查是否与已有方块重叠（忽略Row < 0的情况，因为方块可以从顶部上方开始）
				if boardRow >= 0 && g.Board.Grid[boardRow][boardCol] != 0 {
					return false
				}
			}
		}
	}

	return true
}

// LockBlock 锁定方块到游戏面板
func (g *Game) LockBlock() {
	shape := g.CurrentPiece.Shape
	shapeHeight := len(shape)
	if shapeHeight == 0 {
		return
	}
	shapeWidth := len(shape[0])

	// 将方块锁定到游戏面板
	for y := 0; y < shapeHeight; y++ {
		for x := 0; x < shapeWidth; x++ {
			if shape[y][x] == 1 {
				boardRow := g.CurrentPiece.Position.Row + y
				boardCol := g.CurrentPiece.Position.Col + x

				// 确保行号在合法范围内
				if boardRow >= 0 && boardRow < g.Board.Height {
					// 使用方块类型作为单元格值
					g.Board.Grid[boardRow][boardCol] = int(g.CurrentPiece.Type)
				}
			}
		}
	}

	// 检查是否有行可以消除
	g.ClearCompletedLines()

	// 生成新方块
	g.CurrentPiece = g.NextPiece
	g.NextPiece = NewRandomBlock()

	// 检查游戏是否结束
	if !g.CanPlaceBlock(g.CurrentPiece.Position, g.CurrentPiece.Shape) {
		g.GameOver()
	}
}

// ClearCompletedLines 消除已完成的行
func (g *Game) ClearCompletedLines() {
	linesCleared := 0

	// 从底部向上检查每一行
	for y := g.Board.Height - 1; y >= 0; y-- {
		isComplete := true
		for x := 0; x < g.Board.Width; x++ {
			if g.Board.Grid[y][x] == 0 {
				isComplete = false
				break
			}
		}

		if isComplete {
			// 消除当前行
			for i := y; i > 0; i-- {
				copy(g.Board.Grid[i], g.Board.Grid[i-1])
			}

			// 清空第一行
			for x := 0; x < g.Board.Width; x++ {
				g.Board.Grid[0][x] = 0
			}

			linesCleared++
			// 检查当前行（现在是新的一行）
			y++
		}
	}

	if linesCleared > 0 {
		g.Lines += linesCleared
		g.UpdateScore(linesCleared)
		g.UpdateLevel()
	}
}

// UpdateScore 更新分数
func (g *Game) UpdateScore(linesCleared int) {
	// 根据消除的行数计算分数
	scoreMap := map[int]int{
		1: 100,
		2: 300,
		3: 500,
		4: 800,
	}

	score := scoreMap[linesCleared] * g.Level
	g.Score += score
}

// UpdateLevel 更新等级
func (g *Game) UpdateLevel() {
	// 每消除10行升级一次
	newLevel := g.Lines/10 + 1
	if newLevel > g.Level {
		g.Level = newLevel
		// 提高游戏速度
		g.GameSpeed = 1000 - (g.Level-1)*100
		if g.GameSpeed < 100 {
			g.GameSpeed = 100 // 最小速度100ms
		}
	}
}

// GameOver 游戏结束
func (g *Game) GameOver() {
	g.State = GameStateGameOver
	g.IsGameOver = true
}

// StartGame 开始游戏
func (g *Game) StartGame() {
	if g.State == GameStateGameOver {
		// 游戏结束后重新开始，创建新游戏
		*g = *NewGame(g.ID)
	}

	g.State = GameStateRunning
	g.LastUpdate = time.Now()
}

// PauseGame 暂停游戏
func (g *Game) PauseGame() {
	if g.State == GameStateRunning {
		g.State = GameStatePaused
	}
}

// ResumeGame 继续游戏
func (g *Game) ResumeGame() {
	if g.State == GameStatePaused {
		g.State = GameStateRunning
		g.LastUpdate = time.Now()
	}
}

// UpdateGame 更新游戏状态
func (g *Game) UpdateGame() {
	if g.State != GameStateRunning || g.IsGameOver {
		return
	}

	now := time.Now()
	if now.Sub(g.LastUpdate).Milliseconds() >= int64(g.GameSpeed) {
		g.MoveBlockDown()
		g.LastUpdate = now
	}
}

// SendAction 处理游戏动作
func (g *Game) SendAction(action string) {
	switch action {
	case "left":
		g.MoveBlockLeft()
	case "right":
		g.MoveBlockRight()
	case "down":
		g.MoveBlockDown()
	case "rotate":
		g.RotateBlock()
	case "drop":
		g.DropBlock()
	case "start":
		g.StartGame()
	case "pause":
		g.PauseGame()
	case "resume":
		g.ResumeGame()
	}
}

package utils

import (
	"blockcade/models"
	"math/rand"
	"sync"
	"time"

	"github.com/astaxie/beego"
)

// GameManager 游戏管理器
type GameManager struct {
	games          map[string]*models.Game
	mu             sync.RWMutex
	updateInterval time.Duration
}

var gameManager *GameManager
var once sync.Once

// GetGameManager 获取游戏管理器单例
func GetGameManager() *GameManager {
	once.Do(func() {
		// 初始化随机数生成器
		rand.Seed(time.Now().UnixNano())

		// 设置游戏更新间隔（100ms）
		updateInterval := 100 * time.Millisecond

		// 从配置文件读取配置
		if ui := beego.AppConfig.String("game.updateInterval"); ui != "" {
			interval, err := beego.AppConfig.Int("game.updateInterval")
			if err == nil && interval > 0 {
				updateInterval = time.Duration(interval) * time.Millisecond
			}
		}

		gameManager = &GameManager{
			games:          make(map[string]*models.Game),
			updateInterval: updateInterval,
		}

		// 启动游戏更新循环
		go gameManager.startUpdateLoop()
	})

	return gameManager
}

// CreateGame 创建新游戏
func (gm *GameManager) CreateGame(gameID string) *models.Game {
	gm.mu.Lock()
	defer gm.mu.Unlock()

	game := models.NewGame(gameID)
	gm.games[gameID] = game

	return game
}

// GetGame 获取游戏实例
func (gm *GameManager) GetGame(gameID string) (*models.Game, bool) {
	gm.mu.RLock()
	defer gm.mu.RUnlock()

	game, exists := gm.games[gameID]
	return game, exists
}

// RemoveGame 移除游戏实例
func (gm *GameManager) RemoveGame(gameID string) {
	gm.mu.Lock()
	defer gm.mu.Unlock()

	delete(gm.games, gameID)
}

// startUpdateLoop 启动游戏更新循环
func (gm *GameManager) startUpdateLoop() {
	ticker := time.NewTicker(gm.updateInterval)
	defer ticker.Stop()

	for range ticker.C {
		gm.updateGames()
	}
}

// updateGames 更新所有游戏状态
func (gm *GameManager) updateGames() {
	gm.mu.RLock()
	// 创建游戏列表副本，避免在更新时锁定整个映射
	games := make([]*models.Game, 0, len(gm.games))
	for _, game := range gm.games {
		games = append(games, game)
	}
	gm.mu.RUnlock()

	// 更新每个游戏状态
	for _, game := range games {
		game.UpdateGame()
	}
}

// MoveLeft 向左移动方块
func (gm *GameManager) MoveLeft(gameID string) bool {
	gm.mu.Lock()
	defer gm.mu.Unlock()

	game, exists := gm.games[gameID]
	if !exists || game.State != models.GameStateRunning {
		return false
	}

	game.MoveBlockLeft()
	return true
}

// MoveRight 向右移动方块
func (gm *GameManager) MoveRight(gameID string) bool {
	gm.mu.Lock()
	defer gm.mu.Unlock()

	game, exists := gm.games[gameID]
	if !exists || game.State != models.GameStateRunning {
		return false
	}

	game.MoveBlockRight()
	return true
}

// MoveDown 向下移动方块
func (gm *GameManager) MoveDown(gameID string) bool {
	gm.mu.Lock()
	defer gm.mu.Unlock()

	game, exists := gm.games[gameID]
	if !exists || game.State != models.GameStateRunning {
		return false
	}

	game.MoveBlockDown()
	return true
}

// HardDrop 直接落底
func (gm *GameManager) HardDrop(gameID string) bool {
	gm.mu.Lock()
	defer gm.mu.Unlock()

	game, exists := gm.games[gameID]
	if !exists || game.State != models.GameStateRunning {
		return false
	}

	game.DropBlock()
	return true
}

// Rotate 旋转方块
func (gm *GameManager) Rotate(gameID string) bool {
	gm.mu.Lock()
	defer gm.mu.Unlock()

	game, exists := gm.games[gameID]
	if !exists || game.State != models.GameStateRunning {
		return false
	}

	game.RotateBlock()
	return true
}

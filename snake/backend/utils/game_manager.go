package utils

import (
	"blockcade/models"
	"sync"
	"time"

	"github.com/astaxie/beego"
)

// GameManager 游戏管理器
type GameManager struct {
	games    map[string]*models.Game
	mutex    sync.RWMutex
	width    int
	height   int
	maxWalls int
	speed    int
}

var gameManager *GameManager
var once sync.Once

// GetGameManager 获取游戏管理器单例
func GetGameManager() *GameManager {
	once.Do(func() {
		width := 15  // 固定为15，不允许修改
		height := 15 // 固定为15，不允许修改
		maxWalls := 6 // 最多6道墙，符合需求
		speed := 200

		// 从配置文件读取配置（除了width和height）
		if mw := beego.AppConfig.String("wall.max"); mw != "" {
			maxWalls, _ = beego.AppConfig.Int("wall.max")
		}
		if s := beego.AppConfig.String("game.speed"); s != "" {
			speed, _ = beego.AppConfig.Int("game.speed")
		}

		gameManager = &GameManager{
			games:    make(map[string]*models.Game),
			width:    width,
			height:   height,
			maxWalls: maxWalls,
			speed:    speed,
		}

		// 启动游戏更新循环
		go gameManager.startUpdateLoop()
	})

	return gameManager
}

// CreateGame 创建新游戏
func (gm *GameManager) CreateGame(gameID string) *models.Game {
	gm.mutex.Lock()
	defer gm.mutex.Unlock()

	game := models.NewGame(gameID, gm.width, gm.height, gm.maxWalls)
	gm.games[gameID] = game

	return game
}

// GetGame 获取游戏实例
func (gm *GameManager) GetGame(gameID string) (*models.Game, bool) {
	gm.mutex.RLock()
	defer gm.mutex.RUnlock()

	game, exists := gm.games[gameID]
	return game, exists
}

// RemoveGame 移除游戏实例
func (gm *GameManager) RemoveGame(gameID string) {
	gm.mutex.Lock()
	defer gm.mutex.Unlock()

	delete(gm.games, gameID)
}

// UpdateGame 更新游戏方向
func (gm *GameManager) UpdateGameDirection(gameID string, direction models.Direction) bool {
	gm.mutex.Lock()
	defer gm.mutex.Unlock()

	game, exists := gm.games[gameID]
	if !exists || game.Status != models.GameStatusRunning {
		return false
	}

	game.ChangeDirection(direction)
	return true
}

// startUpdateLoop 启动游戏更新循环
func (gm *GameManager) startUpdateLoop() {
	ticker := time.NewTicker(time.Duration(gm.speed) * time.Millisecond)
	defer ticker.Stop()

	for range ticker.C {
		gm.updateAllGames()
	}
}

// updateAllGames 更新所有游戏
func (gm *GameManager) updateAllGames() {
	gm.mutex.Lock()
	defer gm.mutex.Unlock()

	// 清理超过一定时间的游戏记录
	now := time.Now()
	for gameID, game := range gm.games {
		// 更新游戏状态
		game.Update()

		// 如果超过10分钟没有活动，移除游戏
		if now.Sub(game.CreatedAt).Minutes() > 10 {
			delete(gm.games, gameID)
		}
	}
}

package controllers

import (
	"blockcade/utils"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/astaxie/beego"
)

// GameController 游戏控制器
type GameController struct {
	beego.Controller
}

// GameActionRequest 游戏操作请求结构
type GameActionRequest struct {
	Action string `json:"action"`
}

// SaveRecordRequest 保存记录请求结构
type SaveRecordRequest struct {
	PlayerName string `json:"playerName"`
	Score      int    `json:"score"`
}

// NewGame 创建新游戏
// @Title 创建新游戏
// @Description 创建一个新的游戏实例
// @Success 200 {object} models.Game
// @router /api/game [post]
func (c *GameController) NewGame() {
	// 生成游戏ID（可以使用UUID，但为了简化，这里使用时间戳）
	gameID := time.Now().Format("20060102150405") + "-" + c.Ctx.Request.RemoteAddr

	// 获取游戏管理器并创建游戏
	gameManager := utils.GetGameManager()
	game := gameManager.CreateGame(gameID)

	// 返回游戏信息
	c.Data["json"] = game
	c.ServeJSON()
}

// GetGame 获取游戏状态
// @Title 获取游戏状态
// @Description 根据游戏ID获取游戏状态
// @Param id path string true "游戏ID"
// @Success 200 {object} models.Game
// @Failure 404 {object} ErrorResponse
// @router /api/game/:id [get]
func (c *GameController) GetGame() {
	gameID := c.Ctx.Input.Param(":id")

	// 获取游戏管理器
	gameManager := utils.GetGameManager()
	game, exists := gameManager.GetGame(gameID)

	if !exists {
		c.Data["json"] = map[string]string{"error": "游戏不存在"}
		c.Ctx.Output.Status = http.StatusNotFound
	} else {
		c.Data["json"] = game
	}

	c.ServeJSON()
}

// UpdateGame 处理游戏操作
// @Title 处理游戏操作
// @Description 根据游戏ID和操作类型更新游戏状态
// @Param id path string true "游戏ID"
// @Param action body GameActionRequest true "操作请求"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @router /api/game/:id/action [post]
func (c *GameController) UpdateGame() {
	gameID := c.Ctx.Input.Param(":id")

	// 解析请求体
	var request GameActionRequest
	if err := json.NewDecoder(c.Ctx.Request.Body).Decode(&request); err != nil {
		c.Data["json"] = map[string]string{"error": "无效的请求格式"}
		c.Ctx.Output.Status = http.StatusBadRequest
		c.ServeJSON()
		return
	}

	// 获取游戏管理器
	gameManager := utils.GetGameManager()
	_, exists := gameManager.GetGame(gameID)

	if !exists {
		c.Data["json"] = map[string]string{"error": "游戏不存在"}
		c.Ctx.Output.Status = http.StatusNotFound
		c.ServeJSON()
		return
	}

	// 处理操作
	success := false
	switch request.Action {
	case "left":
		success = gameManager.MoveLeft(gameID)
	case "right":
		success = gameManager.MoveRight(gameID)
	case "down":
		success = gameManager.MoveDown(gameID)
	case "harddrop":
		success = gameManager.HardDrop(gameID)
	case "hardDrop":
		success = gameManager.HardDrop(gameID)
	case "rotate":
		success = gameManager.Rotate(gameID)
	default:
		c.Data["json"] = map[string]string{"error": "无效的操作类型"}
		c.Ctx.Output.Status = http.StatusBadRequest
		c.ServeJSON()
		return
	}

	if success {
		c.Data["json"] = map[string]bool{"success": true}
	} else {
		c.Data["json"] = map[string]string{"error": "操作失败"}
		c.Ctx.Output.Status = http.StatusBadRequest
	}

	c.ServeJSON()
}

// RemoveGame 移除游戏
// @Title 移除游戏
// @Description 根据游戏ID移除游戏实例
// @Param id path string true "游戏ID"
// @Success 200 {object} SuccessResponse
// @Failure 404 {object} ErrorResponse
// @router /api/game/:id [delete]
func (c *GameController) RemoveGame() {
	gameID := c.Ctx.Input.Param(":id")

	// 获取游戏管理器
	gameManager := utils.GetGameManager()
	_, exists := gameManager.GetGame(gameID)

	if !exists {
		c.Data["json"] = map[string]string{"error": "游戏不存在"}
		c.Ctx.Output.Status = http.StatusNotFound
	} else {
		gameManager.RemoveGame(gameID)
		c.Data["json"] = map[string]bool{"success": true}
	}

	c.ServeJSON()
}

// SaveRecord 保存游戏记录
// @Title 保存游戏记录
// @Description 保存游戏得分记录
// @Param id path string true "游戏ID"
// @Param request body SaveRecordRequest true "记录请求"
// @Success 200 {object} SuccessResponse
// @Failure 404 {object} ErrorResponse
// @router /api/game/:id/record [post]
func (c *GameController) SaveRecord() {
	gameID := c.Ctx.Input.Param(":id")

	// 解析请求体
	requestBody, err := io.ReadAll(c.Ctx.Request.Body)
	if err != nil {
		beego.Error("读取请求体失败:", err)
		c.Data["json"] = map[string]string{"error": "读取请求体失败"}
		c.Ctx.Output.Status = http.StatusBadRequest
		c.ServeJSON()
		return
	}

	var req SaveRecordRequest
	if err := json.Unmarshal(requestBody, &req); err != nil {
		c.Data["json"] = map[string]string{"error": "无效的请求格式"}
		c.Ctx.Output.Status = http.StatusBadRequest
		c.ServeJSON()
		return
	}

	// 获取游戏管理器
	gameManager := utils.GetGameManager()
	game, exists := gameManager.GetGame(gameID)

	if !exists {
		c.Data["json"] = map[string]string{"error": "游戏不存在"}
		c.Ctx.Output.Status = http.StatusNotFound
		c.ServeJSON()
		return
	}

	// 保存记录到数据库
	if utils.DB != nil {
		// 将游戏时间转换为整数秒
		timePlayed := int(time.Since(game.LastUpdate).Seconds())
		_, err := utils.DB.Exec(
			"INSERT INTO game_records (score, lines, time_played, player_name) VALUES ($1, $2, $3, $4)",
			req.Score, game.Lines, timePlayed, req.PlayerName,
		)
		if err != nil {
			beego.Error("Failed to save game record:", err)
			c.Data["json"] = map[string]string{"error": "保存记录失败"}
			c.Ctx.Output.Status = http.StatusInternalServerError
			c.ServeJSON()
			return
		}
	}

	c.Data["json"] = map[string]string{"success": "记录保存成功"}
	c.ServeJSON()
}

// GetLeaderboard 获取排行榜
// @Title 获取排行榜
// @Description 获取游戏得分排行榜
// @Param limit query int false "限制数量" default(10)
// @Success 200 {array} LeaderboardItem
// @router /api/leaderboard [get]
func (c *GameController) GetLeaderboard() {
	// 获取限制参数
	limit := c.GetString("limit")
	if limit == "" {
		limit = "10"
	}

	// 从数据库查询排行榜
	type LeaderboardItem struct {
		Score      int       `json:"score"`
		Lines      int       `json:"lines"`
		TimePlayed int       `json:"time_played"`
		PlayerName string    `json:"player_name"`
		CreatedAt  time.Time `json:"created_at"`
	}

	var items []LeaderboardItem

	if utils.DB != nil {
		rows, err := utils.DB.Query(
			"SELECT score, lines, time_played, player_name, created_at FROM game_records ORDER BY score DESC LIMIT $1",
			limit,
		)
		if err != nil {
			beego.Error("Failed to get leaderboard:", err)
			c.Data["json"] = []LeaderboardItem{}
			c.ServeJSON()
			return
		}
		defer rows.Close()

		for rows.Next() {
			var item LeaderboardItem
			if err := rows.Scan(&item.Score, &item.Lines, &item.TimePlayed, &item.PlayerName, &item.CreatedAt); err != nil {
				beego.Error("Failed to scan leaderboard row:", err)
				continue
			}
			items = append(items, item)
		}
	}

	c.Data["json"] = items
	c.ServeJSON()
}

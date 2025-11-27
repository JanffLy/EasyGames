package controllers

import (
	"blockcade/models"
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

// GameRequest 游戏请求结构
type GameRequest struct {
	Direction string `json:"direction"`
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

// UpdateDirection 更新游戏方向
// @Title 更新游戏方向
// @Description 更新蛇的移动方向
// @Param id path string true "游戏ID"
// @Param direction body object true "方向请求"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @router /api/game/:id/direction [post]
func (c *GameController) UpdateDirection() {
	gameID := c.Ctx.Input.Param(":id")
	beego.Info("接收到更新方向请求，游戏ID:", gameID)

	// 记录请求体原始内容
	// TODO: 这里接收不到具体的数据
	requestBody, err := io.ReadAll(c.Ctx.Request.Body)
	if err != nil {
		beego.Error("读取请求体失败:", err)
		c.Data["json"] = map[string]string{"error": "读取请求体失败"}
		c.Ctx.Output.Status = http.StatusBadRequest
		c.ServeJSON()
		return
	}
	beego.Info("请求体原始内容:", string(requestBody))

	// 使用map解析JSON，更加灵活
	var data map[string]interface{}
	if err := json.Unmarshal(requestBody, &data); err != nil {
		beego.Error("解析请求体失败:", err)
		c.Data["json"] = map[string]string{"error": "无效的请求格式"}
		c.Ctx.Output.Status = http.StatusBadRequest
		c.ServeJSON()
		return
	}

	// 从map中获取direction字段，支持大小写
	directionStr := ""
	if val, ok := data["direction"]; ok {
		directionStr = val.(string)
	} else if val, ok := data["Direction"]; ok {
		directionStr = val.(string)
	}

	beego.Info("解析到的方向值:", directionStr)
	if directionStr == "" {
		beego.Error("请求体中没有找到direction字段")
		c.Data["json"] = map[string]string{"error": "请求体中缺少direction字段"}
		c.Ctx.Output.Status = http.StatusBadRequest
		c.ServeJSON()
		return
	}

	// 转换方向字符串为枚举值
	var direction models.Direction
	switch directionStr {
	case "up":
		direction = models.Up
	case "down":
		direction = models.Down
	case "left":
		direction = models.Left
	case "right":
		direction = models.Right
	default:
		beego.Error("无效的方向值:", directionStr)
		c.Data["json"] = map[string]string{"error": "无效的方向"}
		c.Ctx.Output.Status = http.StatusBadRequest
		c.ServeJSON()
		return
	}
	beego.Info("方向值转换成功:", direction)

	// 获取游戏管理器并更新方向
	gameManager := utils.GetGameManager()
	if success := gameManager.UpdateGameDirection(gameID, direction); !success {
		c.Data["json"] = map[string]string{"error": "游戏不存在或已结束"}
		c.Ctx.Output.Status = http.StatusNotFound
	} else {
		c.Data["json"] = map[string]string{"success": "方向更新成功"}
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

// SaveRecordRequest 保存记录请求结构
type SaveRecordRequest struct {
	PlayerName string `json:"playerName"`
	Score      int    `json:"score"`
}

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
	beego.Info("请求体原始内容:", string(requestBody))

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
		_, err := utils.DB.Exec(
			"INSERT INTO game_records (score, time_played, food_count, player_name) VALUES ($1, $2, $3, $4)",
			req.Score, game.Time, game.FoodCount, req.PlayerName,
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
		TimePlayed int       `json:"time_played"`
		FoodCount  int       `json:"food_count"`
		CreatedAt  time.Time `json:"created_at"`
	}

	var items []LeaderboardItem

	if utils.DB != nil {
		rows, err := utils.DB.Query(
			"SELECT score, time_played, food_count, created_at FROM game_records ORDER BY score DESC LIMIT $1",
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
			if err := rows.Scan(&item.Score, &item.TimePlayed, &item.FoodCount, &item.CreatedAt); err != nil {
				beego.Error("Failed to scan leaderboard row:", err)
				continue
			}
			items = append(items, item)
		}
	}

	c.Data["json"] = items
	c.ServeJSON()
}

package controllers

import (
	"blockcade/utils"
	"encoding/json"
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

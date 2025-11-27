package routers

import (
	"blockcade/controllers"

	"github.com/astaxie/beego"
	beegoCtx "github.com/astaxie/beego/context"
)

// InitRouter 初始化路由
func InitRouter() {
	// 创建控制器实例
	gameController := &controllers.GameController{}

	// 设置CORS中间件
	beego.InsertFilter("*", beego.BeforeRouter, corsHandler())

	// 注册路由
	beego.Router("/api/game", gameController, "post:NewGame")
	beego.Router("/api/game/:id", gameController, "get:GetGame")
	beego.Router("/api/game/:id/direction", gameController, "post:UpdateDirection")
	beego.Router("/api/game/:id/record", gameController, "post:SaveRecord")
	beego.Router("/api/leaderboard", gameController, "get:GetLeaderboard")
}

// corsHandler CORS中间件
func corsHandler() beego.FilterFunc {
	return func(ctx *beegoCtx.Context) {
		// 设置CORS头信息
		ctx.Output.Header("Access-Control-Allow-Origin", "*")
		ctx.Output.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		ctx.Output.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
		ctx.Output.Header("Access-Control-Allow-Credentials", "true")

		// 处理预检请求
		if ctx.Request.Method == "OPTIONS" {
			ctx.Abort(204, "")
			return
		}
	}
}

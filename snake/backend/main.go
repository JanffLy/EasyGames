package main

import (
	"blockcade/routers"
	"blockcade/utils"
	"log"

	"github.com/astaxie/beego"
)

func main() {
	// 初始化数据库连接
	utils.InitDB()
	defer utils.CloseDB()

	// 初始化Redis连接
	utils.InitRedis()
	defer utils.CloseRedis()

	// 初始化游戏管理器
	gameManager := utils.GetGameManager()
	if gameManager == nil {
		log.Fatal("Failed to initialize game manager")
	}

	// 初始化路由
	routers.InitRouter()

	// 启动服务器
	log.Println("Server starting on port", beego.AppConfig.String("httpport"))
	beego.Run()
}

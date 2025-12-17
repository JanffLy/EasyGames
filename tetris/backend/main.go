package main

import (
	"fmt"
	"log"

	"blockcade/routers"
	"blockcade/utils"

	"github.com/astaxie/beego"
)

func main() {
	// 初始化游戏管理器
	gameManager := utils.GetGameManager()
	if gameManager == nil {
		log.Fatal("Failed to initialize game manager")
	}

	// 初始化路由
	routers.InitRouter()

	// 打印启动信息
	port := beego.AppConfig.String("httpport")
	if port == "" {
		port = "8080"
	}
	fmt.Printf("Starting Tetris game server on port %s...\n", port)

	// 启动beego应用
	beego.Run()
}

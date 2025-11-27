package utils

import (
	"context"
	"fmt"
	"log"

	"github.com/astaxie/beego"
	"github.com/go-redis/redis/v8"
)

var RedisClient *redis.Client
var Ctx = context.Background()

// InitRedis 初始化Redis连接
func InitRedis() {
	host := beego.AppConfig.String("redis.host")
	port := beego.AppConfig.String("redis.port")
	password := beego.AppConfig.String("redis.password")
	dbStr := beego.AppConfig.String("redis.db")

	db := 0
	_, err := fmt.Sscanf(dbStr, "%d", &db)
	if err != nil {
		log.Printf("Failed to parse Redis DB index, using default 0: %v\n", err)
	}

	RedisClient = redis.NewClient(&redis.Options{
		Addr:     host + ":" + port,
		Password: password,
		DB:       db,
	})

	// 测试连接
	_, err = RedisClient.Ping(Ctx).Result()
	if err != nil {
		log.Printf("Failed to connect to Redis: %v\n", err)
		// Redis连接失败不应该导致应用崩溃，因为我们可以在内存中管理游戏状态
	} else {
		log.Println("Redis connected successfully")
	}
}

// CloseRedis 关闭Redis连接
func CloseRedis() {
	if RedisClient != nil {
		RedisClient.Close()
	}
}

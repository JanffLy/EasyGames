package utils

import (
	"strconv"

	"github.com/astaxie/beego"
	"github.com/go-redis/redis"
)

var RedisClient *redis.Client

// InitRedis 初始化Redis连接
func InitRedis() {
	host := beego.AppConfig.String("redis.host")
	port := beego.AppConfig.String("redis.port")
	password := beego.AppConfig.String("redis.password")
	dbStr := beego.AppConfig.String("redis.db")

	db, err := strconv.Atoi(dbStr)
	if err != nil {
		db = 0
	}

	// 连接Redis
	client := redis.NewClient(&redis.Options{
		Addr:     host + ":" + port,
		Password: password,
		DB:       db,
	})

	// 测试连接
	pong, err := client.Ping().Result()
	if err != nil {
		beego.Error("Failed to connect to Redis:", err)
		// 不崩溃，只记录错误
		return
	}

	RedisClient = client
	beego.Info("Redis connected successfully:", pong)
}

// CloseRedis 关闭Redis连接
func CloseRedis() {
	if RedisClient != nil {
		RedisClient.Close()
	}
}

package utils

import (
	"database/sql"
	"log"

	"github.com/astaxie/beego"
	_ "github.com/lib/pq"
)

var DB *sql.DB

// InitDB 初始化数据库连接
func InitDB() {
	host := beego.AppConfig.String("db.host")
	port := beego.AppConfig.String("db.port")
	user := beego.AppConfig.String("db.user")
	password := beego.AppConfig.String("db.password")
	dbname := beego.AppConfig.String("db.name")

	psqlInfo := "host=" + host + " port=" + port + " user=" + user + " password=" + password + " dbname=" + dbname + " sslmode=disable"

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v\n", err)
	}

	// 测试连接
	err = db.Ping()
	if err != nil {
		log.Fatalf("Failed to ping database: %v\n", err)
	}

	DB = db
	log.Println("Database connected successfully")

	// 创建表
	createTables()
}

// createTables 创建必要的数据表
func createTables() {
	// 创建游戏记录表
	_, err := DB.Exec(`
	CREATE TABLE IF NOT EXISTS game_records (
		id SERIAL PRIMARY KEY,
		score INTEGER NOT NULL,
		time_played INTEGER NOT NULL,
		food_count INTEGER NOT NULL,
		player_name VARCHAR(100),
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)
	`)
	if err != nil {
		log.Printf("Failed to create game_records table: %v\n", err)
	}
}

// CloseDB 关闭数据库连接
func CloseDB() {
	if DB != nil {
		DB.Close()
	}
}

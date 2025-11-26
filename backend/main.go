package main

import (
	"log"
	"text-wow/api"
	"text-wow/game"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化游戏引擎
	engine := game.NewGameEngine()

	// 创建Gin路由
	r := gin.Default()

	// 配置CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	// 注册API路由
	api.RegisterRoutes(r, engine)

	log.Println("🎮 Text WoW Server starting on :8080")
	log.Println("⚔️  Battle Engine initialized")

	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

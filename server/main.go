package main

import (
	"log"
	"net/http"
	"text-wow/game"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// CORS 配置
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	// 初始化游戏引擎
	engine := game.NewEngine()

	// API 路由
	api := r.Group("/api")
	{
		// 创建新角色
		api.POST("/character", func(c *gin.Context) {
			var req struct {
				Name      string `json:"name"`
				Race      string `json:"race"`
				ClassName string `json:"class"`
			}
			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			char := engine.CreateCharacter(req.Name, req.Race, req.ClassName)
			c.JSON(http.StatusOK, char)
		})

		// 获取角色信息
		api.GET("/character", func(c *gin.Context) {
			char := engine.GetCharacter()
			if char == nil {
				c.JSON(http.StatusNotFound, gin.H{"error": "no character"})
				return
			}
			c.JSON(http.StatusOK, char)
		})

		// 开始/停止挂机
		api.POST("/battle/toggle", func(c *gin.Context) {
			isRunning := engine.ToggleBattle()
			c.JSON(http.StatusOK, gin.H{"running": isRunning})
		})

		// 获取战斗状态
		api.GET("/battle/status", func(c *gin.Context) {
			status := engine.GetBattleStatus()
			c.JSON(http.StatusOK, status)
		})

		// 获取战斗日志
		api.GET("/battle/logs", func(c *gin.Context) {
			logs := engine.GetBattleLogs()
			c.JSON(http.StatusOK, gin.H{"logs": logs})
		})

		// 执行一次战斗回合（用于原型测试）
		api.POST("/battle/tick", func(c *gin.Context) {
			result := engine.BattleTick()
			c.JSON(http.StatusOK, result)
		})

		// 获取可用区域
		api.GET("/zones", func(c *gin.Context) {
			zones := engine.GetZones()
			c.JSON(http.StatusOK, gin.H{"zones": zones})
		})

		// 切换区域
		api.POST("/zone/change", func(c *gin.Context) {
			var req struct {
				ZoneID string `json:"zone_id"`
			}
			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			if err := engine.ChangeZone(req.ZoneID); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, gin.H{"success": true})
		})
	}

	log.Println("🎮 Text WoW Server starting on :8080")
	r.Run(":8080")
}

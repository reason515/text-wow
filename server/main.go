package main

import (
	"log"
	"net/http"

	"text-wow/internal/api"
	"text-wow/internal/database"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化数据库
	if err := database.Init(); err != nil {
		log.Fatalf("❌ Failed to initialize database: %v", err)
	}
	defer database.Close()

	// 创建Gin实例
	r := gin.Default()

	// CORS 配置
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://localhost:5174", "http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	// 创建API处理器
	h := api.NewHandler()
	chatHandler := api.NewChatHandler()
	battleHandler := api.NewBattleHandler()
	strategyHandler := api.NewStrategyHandlers()

	// API 路由
	apiGroup := r.Group("/api")
	{
		// ═══════════════════════════════════════════════════════════
		// 公开API（无需认证）
		// ═══════════════════════════════════════════════════════════

		// 健康检查
		apiGroup.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "ok", "version": "0.1.0"})
		})

		// 认证
		auth := apiGroup.Group("/auth")
		{
			auth.POST("/register", h.Register)
			auth.POST("/login", h.Login)
		}

		// 游戏配置（公开）
		apiGroup.GET("/races", h.GetRaces)
		apiGroup.GET("/classes", h.GetClasses)
		apiGroup.GET("/zones", h.GetZones)

		// ═══════════════════════════════════════════════════════════
		// 需要认证的API
		// ═══════════════════════════════════════════════════════════

		protected := apiGroup.Group("")
		protected.Use(h.AuthMiddleware())
		{
			// 用户
			protected.GET("/user", h.GetCurrentUser)

			// 角色
			protected.GET("/character", h.GetCharacter) // 获取当前活跃角色（单数）
			protected.GET("/characters", h.GetCharacters)
			protected.POST("/characters", h.CreateCharacter)
			protected.PUT("/characters/active", h.SetCharacterActive)
			protected.POST("/characters/:characterId/allocate", h.AllocateAttributePoints)

			// 技能选择
			protected.GET("/characters/:characterId/skills/initial", h.GetInitialSkillSelection)
			protected.POST("/characters/:characterId/skills/initial", h.SelectInitialSkills)
			protected.GET("/characters/:characterId/skills/selection", h.GetSkillSelection)
			protected.POST("/characters/:characterId/skills/select", h.SelectSkill)
			protected.GET("/characters/:characterId/skills", h.GetCharacterSkills)

			// 小队
			protected.GET("/team", h.GetTeam)

			// 聊天
			chat := protected.Group("/chat")
			{
				chat.GET("/messages", chatHandler.GetMessages)
				chat.POST("/send", chatHandler.SendMessage)
				chat.GET("/online", chatHandler.GetOnlineUsers)
				chat.POST("/block", chatHandler.BlockUser)
				chat.POST("/unblock", chatHandler.UnblockUser)
				chat.POST("/online", chatHandler.SetOnline)
				chat.POST("/offline", chatHandler.SetOffline)
				chat.POST("/heartbeat", chatHandler.Heartbeat)
			}

			// 战斗
			battle := protected.Group("/battle")
			{
				battle.POST("/start", battleHandler.StartBattle)
				battle.POST("/stop", battleHandler.StopBattle)
				battle.POST("/toggle", battleHandler.ToggleBattle)
				battle.POST("/tick", battleHandler.BattleTick)
				battle.GET("/status", battleHandler.GetBattleStatus)
				battle.GET("/logs", battleHandler.GetBattleLogs)
				battle.POST("/zone", battleHandler.ChangeZone)
			}

			// 策略
			protected.GET("/characters/:characterId/strategies", strategyHandler.GetStrategies)
			protected.POST("/characters/:characterId/strategies", strategyHandler.CreateStrategy)
			protected.GET("/strategies/:strategyId", strategyHandler.GetStrategy)
			protected.PUT("/strategies/:strategyId", strategyHandler.UpdateStrategy)
			protected.DELETE("/strategies/:strategyId", strategyHandler.DeleteStrategy)
			protected.POST("/strategies/:strategyId/activate", strategyHandler.SetActiveStrategy)
			protected.GET("/strategy-templates", strategyHandler.GetStrategyTemplates)
			protected.GET("/strategy-condition-types", strategyHandler.GetConditionTypes)

			// 战斗统计
			stats := protected.Group("/stats")
			{
				stats.GET("/overview", h.GetBattleStatsOverview)
				stats.GET("/daily", h.GetDailyStats)
				stats.GET("/battles", h.GetRecentBattles)
				stats.GET("/battles/:battleId", h.GetBattleDetail)
				stats.POST("/session/start", h.StartStatsSession)
				stats.POST("/session/reset", h.ResetStatsSession)
				stats.GET("/session/status", h.GetStatsSessionStatus)
				stats.GET("/cumulative/dps", h.GetCumulativeDPSAnalysis)
				stats.GET("/battles/:battleId/dps", h.GetBattleDPSAnalysis)
			}
			protected.GET("/characters/:characterId/stats", h.GetCharacterLifetimeStats)
			protected.GET("/characters/:characterId/stats/summary", h.GetCharacterBattleSummary)
		}
	}

	log.Println("🎮 Text WoW Server starting on :8080")
	log.Println("📌 API Documentation:")
	log.Println("   POST /api/auth/register    - 用户注册")
	log.Println("   POST /api/auth/login       - 用户登录")
	log.Println("   GET  /api/races            - 获取种族列表")
	log.Println("   GET  /api/classes          - 获取职业列表")
	log.Println("   GET  /api/characters       - 获取角色列表 (需认证)")
	log.Println("   POST /api/characters       - 创建角色 (需认证)")
	log.Println("   GET  /api/team             - 获取小队 (需认证)")
	log.Println("   POST /api/battle/start     - 开始战斗 (需认证)")
	log.Println("   POST /api/battle/stop      - 停止战斗 (需认证)")
	log.Println("   POST /api/battle/toggle    - 切换战斗 (需认证)")
	log.Println("   POST /api/battle/tick      - 战斗回合 (需认证)")
	log.Println("   GET  /api/battle/status    - 战斗状态 (需认证)")
	log.Println("   GET  /api/battle/logs      - 战斗日志 (需认证)")
	log.Println("   POST /api/battle/zone      - 切换区域 (需认证)")
	log.Println("   GET  /api/characters/:id/strategies - 获取策略列表 (需认证)")
	log.Println("   POST /api/characters/:id/strategies - 创建策略 (需认证)")
	log.Println("   PUT  /api/strategies/:id   - 更新策略 (需认证)")
	log.Println("   DELETE /api/strategies/:id - 删除策略 (需认证)")

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("❌ Failed to start server: %v", err)
	}
}

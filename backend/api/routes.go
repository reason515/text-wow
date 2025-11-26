package api

import (
	"net/http"
	"text-wow/game"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册API路由
func RegisterRoutes(r *gin.Engine, engine *game.GameEngine) {
	api := r.Group("/api")
	{
		// 获取游戏状态
		api.GET("/state", func(c *gin.Context) {
			state := engine.GetState()
			c.JSON(http.StatusOK, gin.H{
				"success": true,
				"data":    state,
			})
		})

		// 获取所有区域
		api.GET("/zones", func(c *gin.Context) {
			zones := engine.GetZones()
			c.JSON(http.StatusOK, gin.H{
				"success": true,
				"data":    zones,
			})
		})

		// 切换区域
		api.POST("/zone/:id", func(c *gin.Context) {
			zoneID := c.Param("id")
			success := engine.SetZone(zoneID)
			if success {
				c.JSON(http.StatusOK, gin.H{
					"success": true,
					"message": "已切换到新区域",
				})
			} else {
				c.JSON(http.StatusBadRequest, gin.H{
					"success": false,
					"message": "无法切换区域（可能等级不足或区域不存在）",
				})
			}
		})

		// 执行一次战斗
		api.POST("/battle", func(c *gin.Context) {
			result := engine.DoBattle()
			if result == nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"success": false,
					"message": "战斗失败",
				})
				return
			}
			c.JSON(http.StatusOK, gin.H{
				"success": true,
				"data":    result,
				"state":   engine.GetState(),
			})
		})

		// 开始自动战斗
		api.POST("/auto/start", func(c *gin.Context) {
			engine.StartAutoFight()
			c.JSON(http.StatusOK, gin.H{
				"success": true,
				"message": "自动战斗已开启",
			})
		})

		// 停止自动战斗
		api.POST("/auto/stop", func(c *gin.Context) {
			engine.StopAutoFight()
			c.JSON(http.StatusOK, gin.H{
				"success": false,
				"message": "自动战斗已停止",
			})
		})

		// 获取角色信息
		api.GET("/character", func(c *gin.Context) {
			state := engine.GetState()
			c.JSON(http.StatusOK, gin.H{
				"success": true,
				"data":    state.Character,
			})
		})
	}
}

package api

import (
	"fmt"
	"net/http"

	"text-wow/internal/game"
	"text-wow/internal/models"
	"text-wow/internal/repository"

	"github.com/gin-gonic/gin"
)

// BattleHandler 战斗API处理器
type BattleHandler struct {
	battleMgr *game.BattleManager
	charRepo  *repository.CharacterRepository
	gameRepo  *repository.GameRepository
}

// NewBattleHandler 创建战斗处理器
func NewBattleHandler() *BattleHandler {
	return &BattleHandler{
		battleMgr: game.GetBattleManager(),
		charRepo:  repository.NewCharacterRepository(),
		gameRepo:  repository.NewGameRepository(),
	}
}

// ═══════════════════════════════════════════════════════════
// 战斗控制 API
// ═══════════════════════════════════════════════════════════

// StartBattle 开始战斗
func (h *BattleHandler) StartBattle(c *gin.Context) {
	userID := c.GetInt("userID")

	isRunning, err := h.battleMgr.StartBattle(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data: gin.H{
			"isRunning": isRunning,
		},
	})
}

// StopBattle 停止战斗
func (h *BattleHandler) StopBattle(c *gin.Context) {
	userID := c.GetInt("userID")

	err := h.battleMgr.StopBattle(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data: gin.H{
			"isRunning": false,
		},
	})
}

// ToggleBattle 切换战斗状态
func (h *BattleHandler) ToggleBattle(c *gin.Context) {
	userID := c.GetInt("userID")

	isRunning, err := h.battleMgr.ToggleBattle(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data: gin.H{
			"isRunning": isRunning,
		},
	})
}

// BattleTick 执行战斗回合
func (h *BattleHandler) BattleTick(c *gin.Context) {
	// 添加 panic 恢复
	defer func() {
		if r := recover(); r != nil {
			c.JSON(http.StatusInternalServerError, models.APIResponse{
				Success: false,
				Error:   fmt.Sprintf("战斗处理发生错误: %v", r),
			})
		}
	}()

	userID := c.GetInt("userID")

	// 获取用户的所有角色（所有角色都参与战斗）
	characters, err := h.charRepo.GetByUserID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   fmt.Sprintf("获取角色失败: %v", err),
		})
		return
	}

	if len(characters) == 0 {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "no characters",
		})
		return
	}

	// 执行战斗回合
	result, err := h.battleMgr.ExecuteBattleTick(userID, characters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   fmt.Sprintf("执行战斗回合失败: %v", err),
		})
		return
	}

	if result == nil {
		// 战斗未运行
		c.JSON(http.StatusOK, models.APIResponse{
			Success: true,
			Data: gin.H{
				"isRunning": false,
			},
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    result,
	})
}

// GetBattleStatus 获取战斗状态
func (h *BattleHandler) GetBattleStatus(c *gin.Context) {
	userID := c.GetInt("userID")
	
	// 获取所有角色（所有角色都参与战斗）
	characters, _ := h.charRepo.GetByUserID(userID)
	
	// 确保会话已创建，并根据角色阵营设置默认地图
	if len(characters) > 0 {
		// 先获取当前状态，检查是否已有地图
		status := h.battleMgr.GetBattleStatus(userID)
		faction := characters[0].Faction
		playerLevel := characters[0].Level
		
		// 根据角色阵营确定应该使用的默认地图
		defaultZoneID := "elwynn" // 默认联盟地图
		if faction == "horde" {
			defaultZoneID = "durotar" // 部落默认地图
		}
		
		// 如果没有地图，或者当前地图是 elwynn 但角色是部落（需要设置为 durotar）
		// 或者当前地图是 durotar 但角色是联盟（需要设置为 elwynn）
		needsUpdate := false
		if status.CurrentZoneID == "" {
			needsUpdate = true
		} else if status.CurrentZoneID == "elwynn" && faction == "horde" {
			needsUpdate = true
		} else if status.CurrentZoneID == "durotar" && faction == "alliance" {
			needsUpdate = true
		}
		
		if needsUpdate {
			// 使用 ChangeZone 方法来设置默认地图
			// 这会自动处理锁和验证
			err := h.battleMgr.ChangeZone(userID, defaultZoneID, playerLevel, faction)
			// 如果设置失败（比如等级不够），尝试使用更基础的默认地图
			if err != nil && defaultZoneID != "elwynn" {
				h.battleMgr.ChangeZone(userID, "elwynn", playerLevel, faction)
			}
		}
	}
	
	status := h.battleMgr.GetBattleStatus(userID)
	
	// 为每个角色添加buff信息
	for _, char := range characters {
		buffs := h.battleMgr.GetCharacterBuffs(char.ID)
		char.Buffs = buffs
	}
	
	status.Team = characters

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    status,
	})
}

// GetBattleLogs 获取战斗日志
func (h *BattleHandler) GetBattleLogs(c *gin.Context) {
	userID := c.GetInt("userID")

	logs := h.battleMgr.GetBattleLogs(userID, 100)

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data: gin.H{
			"logs": logs,
		},
	})
}

// ═══════════════════════════════════════════════════════════
// 区域 API
// ═══════════════════════════════════════════════════════════

// ChangeZone 切换区域
func (h *BattleHandler) ChangeZone(c *gin.Context) {
	userID := c.GetInt("userID")

	var req struct {
		ZoneID string `json:"zoneId" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "invalid request",
		})
		return
	}

	// 获取玩家等级和阵营（使用第一个角色）
	characters, err := h.charRepo.GetByUserID(userID)
	if err != nil || len(characters) == 0 {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "no characters",
		})
		return
	}

	playerLevel := characters[0].Level
	playerFaction := characters[0].Faction

	// 切换区域
	err = h.battleMgr.ChangeZone(userID, req.ZoneID, playerLevel, playerFaction)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	// 返回新状态
	status := h.battleMgr.GetBattleStatus(userID)
	logs := h.battleMgr.GetBattleLogs(userID, 10)

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data: gin.H{
			"status": status,
			"logs":   logs,
		},
	})
}

// GetZonesWithMonsters 获取区域列表（包含怪物信息）
func (h *BattleHandler) GetZonesWithMonsters(c *gin.Context) {
	userID := c.GetInt("userID")
	
	// 获取玩家阵营
	characters, err := h.charRepo.GetByUserID(userID)
	var playerFaction string
	if err == nil && len(characters) > 0 {
		playerFaction = characters[0].Faction
	}

	zones, err := h.gameRepo.GetZones()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "failed to get zones",
		})
		return
	}

	// 过滤区域：只返回玩家可以进入的区域（阵营匹配或中立区域）
	var availableZones []models.Zone
	for _, zone := range zones {
		// 中立区域或阵营匹配的区域
		if zone.Faction == "" || zone.Faction == playerFaction {
			availableZones = append(availableZones, zone)
		}
	}

	// 获取玩家所有地图的探索度
	explorationRepo := repository.NewExplorationRepository()
	explorations, _ := explorationRepo.GetAllExplorations(userID)

	// 为每个区域加载怪物
	for i := range availableZones {
		monsters, err := h.gameRepo.GetMonstersByZone(availableZones[i].ID)
		if err == nil {
			availableZones[i].Monsters = monsters
		}
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data: gin.H{
			"zones": availableZones,
			"explorations": explorations,
		},
	})
}

// GetExplorations 获取玩家所有地图的探索度
func (h *BattleHandler) GetExplorations(c *gin.Context) {
	userID := c.GetInt("userID")
	
	explorationRepo := repository.NewExplorationRepository()
	explorations, err := explorationRepo.GetAllExplorations(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "failed to get explorations",
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    explorations,
	})
}











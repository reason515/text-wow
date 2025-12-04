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

	status := h.battleMgr.GetBattleStatus(userID)

	// 获取所有角色（所有角色都参与战斗）
	characters, _ := h.charRepo.GetByUserID(userID)
	
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

	// 获取玩家等级（使用第一个角色）
	characters, err := h.charRepo.GetByUserID(userID)
	if err != nil || len(characters) == 0 {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "no characters",
		})
		return
	}

	playerLevel := characters[0].Level

	// 切换区域
	err = h.battleMgr.ChangeZone(userID, req.ZoneID, playerLevel)
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
	zones, err := h.gameRepo.GetZones()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "failed to get zones",
		})
		return
	}

	// 为每个区域加载怪物
	for i := range zones {
		monsters, err := h.gameRepo.GetMonstersByZone(zones[i].ID)
		if err == nil {
			zones[i].Monsters = monsters
		}
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    zones,
	})
}











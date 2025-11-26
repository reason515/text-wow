package api

import (
	"net/http"
	"strconv"
	"text-wow/internal/game"
	"text-wow/internal/models"

	"github.com/gin-gonic/gin"
)

// === 角色相关 ===

// GetCharacter 获取角色信息
func GetCharacter(c *gin.Context) {
	engine := game.GetEngine()
	char := engine.GetCharacter()

	if char == nil {
		c.JSON(http.StatusOK, gin.H{
			"exists": false,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"exists":    true,
		"character": char,
	})
}

// CreateCharacterRequest 创建角色请求
type CreateCharacterRequest struct {
	Name    string `json:"name" binding:"required"`
	Faction string `json:"faction" binding:"required"`
	Race    string `json:"race" binding:"required"`
	Class   string `json:"class" binding:"required"`
}

// CreateCharacter 创建角色
func CreateCharacter(c *gin.Context) {
	var req CreateCharacterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 根据职业设置初始属性
	baseStats := getClassBaseStats(req.Class)

	char := &models.Character{
		ID:          1,
		Name:        req.Name,
		Faction:     req.Faction,
		Race:        req.Race,
		Class:       req.Class,
		Level:       1,
		Exp:         0,
		ExpToNext:   100,
		HP:          baseStats.HP,
		MaxHP:       baseStats.HP,
		MP:          baseStats.MP,
		MaxMP:       baseStats.MP,
		Strength:    baseStats.Strength,
		Agility:     baseStats.Agility,
		Intellect:   baseStats.Intellect,
		Stamina:     baseStats.Stamina,
		Spirit:      baseStats.Spirit,
		Gold:        10,
		CurrentZone: getDefaultZone(req.Faction),
	}

	engine := game.GetEngine()
	engine.SetCharacter(char)

	// 设置默认策略
	strategy := &models.Strategy{
		ID:                1,
		CharacterID:       1,
		SkillPriority:     getDefaultSkillPriority(req.Class),
		HPPotionThreshold: 30,
		MPPotionThreshold: 20,
		TargetPriority:    "lowest_hp",
		AutoLoot:          true,
	}
	engine.SetStrategy(strategy)

	// 设置默认区域
	engine.SetZone(char.CurrentZone)

	c.JSON(http.StatusOK, gin.H{
		"success":   true,
		"character": char,
	})
}

type baseStatsType struct {
	HP, MP, Strength, Agility, Intellect, Stamina, Spirit int
}

func getClassBaseStats(class string) baseStatsType {
	stats := map[string]baseStatsType{
		"warrior": {HP: 120, MP: 30, Strength: 15, Agility: 8, Intellect: 5, Stamina: 14, Spirit: 5},
		"mage":    {HP: 70, MP: 100, Strength: 5, Agility: 8, Intellect: 18, Stamina: 6, Spirit: 12},
		"hunter":  {HP: 90, MP: 50, Strength: 8, Agility: 16, Intellect: 8, Stamina: 10, Spirit: 8},
		"rogue":   {HP: 85, MP: 40, Strength: 10, Agility: 18, Intellect: 6, Stamina: 8, Spirit: 6},
		"priest":  {HP: 65, MP: 120, Strength: 4, Agility: 6, Intellect: 16, Stamina: 5, Spirit: 18},
	}

	if s, ok := stats[class]; ok {
		return s
	}
	return stats["warrior"]
}

func getDefaultZone(faction string) string {
	if faction == "horde" {
		return "durotar"
	}
	return "elwynn_forest"
}

func getDefaultSkillPriority(class string) []string {
	priorities := map[string][]string{
		"warrior": {"heroic_strike", "thunder_clap", "execute", "attack"},
		"mage":    {"fireball", "frostbolt", "arcane_missiles", "attack"},
		"hunter":  {"aimed_shot", "multi_shot", "kill_shot", "attack"},
		"rogue":   {"backstab", "sinister_strike", "eviscerate", "attack"},
		"priest":  {"smite", "shadow_word_pain", "mind_blast", "attack"},
	}

	if p, ok := priorities[class]; ok {
		return p
	}
	return []string{"attack"}
}

// === 战斗相关 ===

// GetBattleStatus 获取战斗状态
func GetBattleStatus(c *gin.Context) {
	engine := game.GetEngine()
	status := engine.GetBattleStatus()
	char := engine.GetCharacter()
	zone := engine.GetCurrentZone()

	c.JSON(http.StatusOK, gin.H{
		"status":    status,
		"character": char,
		"zone":      zone,
	})
}

// StartBattle 开始战斗
func StartBattle(c *gin.Context) {
	engine := game.GetEngine()

	if engine.GetCharacter() == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请先创建角色"})
		return
	}

	if engine.StartBattle() {
		c.JSON(http.StatusOK, gin.H{"success": true, "message": "战斗开始"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "战斗已在进行中"})
	}
}

// StopBattle 停止战斗
func StopBattle(c *gin.Context) {
	engine := game.GetEngine()
	engine.StopBattle()
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "战斗已停止"})
}

// GetBattleLogs 获取战斗日志
func GetBattleLogs(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "50")
	limit, _ := strconv.Atoi(limitStr)

	engine := game.GetEngine()
	logs := engine.GetBattleLogs(limit)

	c.JSON(http.StatusOK, gin.H{"logs": logs})
}

// === 策略相关 ===

// GetStrategy 获取策略
func GetStrategy(c *gin.Context) {
	engine := game.GetEngine()
	strategy := engine.GetStrategy()
	char := engine.GetCharacter()

	var skills []models.Skill
	if char != nil {
		if s, ok := models.ClassSkills[char.Class]; ok {
			skills = s
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"strategy":        strategy,
		"availableSkills": skills,
	})
}

// UpdateStrategyRequest 更新策略请求
type UpdateStrategyRequest struct {
	SkillPriority     []string `json:"skillPriority"`
	HPPotionThreshold int      `json:"hpPotionThreshold"`
	MPPotionThreshold int      `json:"mpPotionThreshold"`
	TargetPriority    string   `json:"targetPriority"`
	AutoLoot          bool     `json:"autoLoot"`
}

// UpdateStrategy 更新策略
func UpdateStrategy(c *gin.Context) {
	var req UpdateStrategyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	engine := game.GetEngine()
	strategy := engine.GetStrategy()

	if strategy == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "策略不存在"})
		return
	}

	strategy.SkillPriority = req.SkillPriority
	strategy.HPPotionThreshold = req.HPPotionThreshold
	strategy.MPPotionThreshold = req.MPPotionThreshold
	strategy.TargetPriority = req.TargetPriority
	strategy.AutoLoot = req.AutoLoot

	engine.SetStrategy(strategy)

	c.JSON(http.StatusOK, gin.H{"success": true, "strategy": strategy})
}

// === 区域相关 ===

// GetZones 获取区域列表
func GetZones(c *gin.Context) {
	engine := game.GetEngine()
	char := engine.GetCharacter()

	zones := models.Zones

	// 标记可用区域
	var result []gin.H
	for _, z := range zones {
		available := true
		if char != nil && char.Level < z.MinLevel {
			available = false
		}

		result = append(result, gin.H{
			"zone":      z,
			"available": available,
		})
	}

	c.JSON(http.StatusOK, gin.H{"zones": result})
}

// SelectZoneRequest 选择区域请求
type SelectZoneRequest struct {
	ZoneID string `json:"zoneId" binding:"required"`
}

// SelectZone 选择区域
func SelectZone(c *gin.Context) {
	var req SelectZoneRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	engine := game.GetEngine()

	// 如果正在战斗，先停止
	status := engine.GetBattleStatus()
	if status.IsRunning {
		engine.StopBattle()
	}

	if engine.SetZone(req.ZoneID) {
		zone := engine.GetCurrentZone()
		c.JSON(http.StatusOK, gin.H{"success": true, "zone": zone})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "区域不存在"})
	}
}


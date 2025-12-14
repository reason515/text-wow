package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"text-wow/internal/models"
	"text-wow/internal/repository"
)

// StrategyHandlers 策略相关的处理器
type StrategyHandlers struct {
	strategyRepo  *repository.StrategyRepository
	characterRepo *repository.CharacterRepository
}

// NewStrategyHandlers 创建策略处理器
func NewStrategyHandlers() *StrategyHandlers {
	return &StrategyHandlers{
		strategyRepo:  repository.NewStrategyRepository(),
		characterRepo: repository.NewCharacterRepository(),
	}
}

// GetStrategies 获取角色的所有策略
// GET /api/characters/:characterId/strategies
func (h *StrategyHandlers) GetStrategies(c *gin.Context) {
	userID := c.GetInt("userID")
	characterID, err := strconv.Atoi(c.Param("characterId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "invalid character id",
		})
		return
	}

	// 验证角色归属
	char, err := h.characterRepo.GetByID(characterID)
	if err != nil {
		c.JSON(http.StatusNotFound, models.APIResponse{
			Success: false,
			Error:   "character not found",
		})
		return
	}
	if char.UserID != userID {
		c.JSON(http.StatusForbidden, models.APIResponse{
			Success: false,
			Error:   "character does not belong to user",
		})
		return
	}

	strategies, err := h.strategyRepo.GetByCharacterID(characterID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "failed to get strategies: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    strategies,
	})
}

// GetStrategy 获取单个策略
// GET /api/strategies/:strategyId
func (h *StrategyHandlers) GetStrategy(c *gin.Context) {
	userID := c.GetInt("userID")
	strategyID, err := strconv.Atoi(c.Param("strategyId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "invalid strategy id",
		})
		return
	}

	strategy, err := h.strategyRepo.GetByID(strategyID)
	if err != nil {
		c.JSON(http.StatusNotFound, models.APIResponse{
			Success: false,
			Error:   "strategy not found",
		})
		return
	}

	// 验证角色归属
	char, err := h.characterRepo.GetByID(strategy.CharacterID)
	if err != nil || char.UserID != userID {
		c.JSON(http.StatusForbidden, models.APIResponse{
			Success: false,
			Error:   "strategy does not belong to user",
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    strategy,
	})
}

// CreateStrategy 创建策略
// POST /api/characters/:characterId/strategies
func (h *StrategyHandlers) CreateStrategy(c *gin.Context) {
	userID := c.GetInt("userID")
	characterID, err := strconv.Atoi(c.Param("characterId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "invalid character id",
		})
		return
	}

	// 验证角色归属
	char, err := h.characterRepo.GetByID(characterID)
	if err != nil {
		c.JSON(http.StatusNotFound, models.APIResponse{
			Success: false,
			Error:   "character not found",
		})
		return
	}
	if char.UserID != userID {
		c.JSON(http.StatusForbidden, models.APIResponse{
			Success: false,
			Error:   "character does not belong to user",
		})
		return
	}

	// 检查策略数量限制 (最多5个)
	count, err := h.strategyRepo.CountByCharacterID(characterID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "failed to count strategies: " + err.Error(),
		})
		return
	}
	if count >= 5 {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "maximum 5 strategies per character",
		})
		return
	}

	var req models.StrategyCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "invalid request: " + err.Error(),
		})
		return
	}

	// 创建策略
	var strategy *models.BattleStrategy
	if req.FromTemplate != "" {
		// 从模板创建
		templates := repository.GetStrategyTemplates()
		if tmpl, ok := templates[req.FromTemplate]; ok {
			strategy = &models.BattleStrategy{
				CharacterID:          characterID,
				Name:                 req.Name,
				IsActive:             false,
				SkillPriority:        tmpl.SkillPriority,
				ConditionalRules:     tmpl.ConditionalRules,
				TargetPriority:       tmpl.TargetPriority,
				SkillTargetOverrides: tmpl.SkillTargetOverrides,
				ResourceThreshold:    tmpl.ResourceThreshold,
				ReservedSkills:       tmpl.ReservedSkills,
				AutoTargetSettings:   tmpl.AutoTargetSettings,
			}
		} else {
			c.JSON(http.StatusBadRequest, models.APIResponse{
				Success: false,
				Error:   "template not found: " + req.FromTemplate,
			})
			return
		}
	} else {
		// 创建默认策略
		strategy = repository.GetDefaultStrategy(characterID, req.Name)
	}

	strategy, err = h.strategyRepo.Create(strategy)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "failed to create strategy: " + err.Error(),
		})
		return
	}

	// 如果是第一个策略，自动设为激活
	if count == 0 {
		h.strategyRepo.SetActive(strategy.ID, characterID)
		strategy.IsActive = true
	}

	c.JSON(http.StatusCreated, models.APIResponse{
		Success: true,
		Message: "strategy created",
		Data:    strategy,
	})
}

// UpdateStrategy 更新策略
// PUT /api/strategies/:strategyId
func (h *StrategyHandlers) UpdateStrategy(c *gin.Context) {
	userID := c.GetInt("userID")
	strategyID, err := strconv.Atoi(c.Param("strategyId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "invalid strategy id",
		})
		return
	}

	strategy, err := h.strategyRepo.GetByID(strategyID)
	if err != nil {
		c.JSON(http.StatusNotFound, models.APIResponse{
			Success: false,
			Error:   "strategy not found",
		})
		return
	}

	// 验证角色归属
	char, err := h.characterRepo.GetByID(strategy.CharacterID)
	if err != nil || char.UserID != userID {
		c.JSON(http.StatusForbidden, models.APIResponse{
			Success: false,
			Error:   "strategy does not belong to user",
		})
		return
	}

	var req models.StrategyUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "invalid request: " + err.Error(),
		})
		return
	}

	// 更新字段
	if req.Name != nil {
		strategy.Name = *req.Name
	}
	if req.SkillPriority != nil {
		strategy.SkillPriority = req.SkillPriority
	}
	if req.ConditionalRules != nil {
		strategy.ConditionalRules = req.ConditionalRules
	}
	if req.TargetPriority != nil {
		strategy.TargetPriority = *req.TargetPriority
	}
	if req.SkillTargetOverrides != nil {
		strategy.SkillTargetOverrides = req.SkillTargetOverrides
	}
	if req.ResourceThreshold != nil {
		strategy.ResourceThreshold = *req.ResourceThreshold
	}
	if req.ReservedSkills != nil {
		strategy.ReservedSkills = req.ReservedSkills
	}
	if req.AutoTargetSettings != nil {
		strategy.AutoTargetSettings = *req.AutoTargetSettings
	}

	// 处理激活状态
	if req.IsActive != nil && *req.IsActive {
		if err := h.strategyRepo.SetActive(strategy.ID, strategy.CharacterID); err != nil {
			c.JSON(http.StatusInternalServerError, models.APIResponse{
				Success: false,
				Error:   "failed to set strategy active: " + err.Error(),
			})
			return
		}
		strategy.IsActive = true
	}

	if err := h.strategyRepo.Update(strategy); err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "failed to update strategy: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "strategy updated",
		Data:    strategy,
	})
}

// DeleteStrategy 删除策略
// DELETE /api/strategies/:strategyId
func (h *StrategyHandlers) DeleteStrategy(c *gin.Context) {
	userID := c.GetInt("userID")
	strategyID, err := strconv.Atoi(c.Param("strategyId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "invalid strategy id",
		})
		return
	}

	strategy, err := h.strategyRepo.GetByID(strategyID)
	if err != nil {
		c.JSON(http.StatusNotFound, models.APIResponse{
			Success: false,
			Error:   "strategy not found",
		})
		return
	}

	// 验证角色归属
	char, err := h.characterRepo.GetByID(strategy.CharacterID)
	if err != nil || char.UserID != userID {
		c.JSON(http.StatusForbidden, models.APIResponse{
			Success: false,
			Error:   "strategy does not belong to user",
		})
		return
	}

	if err := h.strategyRepo.Delete(strategyID); err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "failed to delete strategy: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "strategy deleted",
	})
}

// SetActiveStrategy 设置激活策略
// POST /api/strategies/:strategyId/activate
func (h *StrategyHandlers) SetActiveStrategy(c *gin.Context) {
	userID := c.GetInt("userID")
	strategyID, err := strconv.Atoi(c.Param("strategyId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "invalid strategy id",
		})
		return
	}

	strategy, err := h.strategyRepo.GetByID(strategyID)
	if err != nil {
		c.JSON(http.StatusNotFound, models.APIResponse{
			Success: false,
			Error:   "strategy not found",
		})
		return
	}

	// 验证角色归属
	char, err := h.characterRepo.GetByID(strategy.CharacterID)
	if err != nil || char.UserID != userID {
		c.JSON(http.StatusForbidden, models.APIResponse{
			Success: false,
			Error:   "strategy does not belong to user",
		})
		return
	}

	if err := h.strategyRepo.SetActive(strategyID, strategy.CharacterID); err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "failed to activate strategy: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "strategy activated",
	})
}

// GetStrategyTemplates 获取策略模板列表
// GET /api/strategy-templates
func (h *StrategyHandlers) GetStrategyTemplates(c *gin.Context) {
	templates := repository.GetStrategyTemplates()

	// 转换为列表格式
	type TemplateInfo struct {
		ID          string `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	templateList := []TemplateInfo{
		{ID: "aggressive", Name: "激进输出", Description: "优先高伤害技能，适合低级区刷怪"},
		{ID: "defensive", Name: "稳健生存", Description: "HP低时防御优先，适合高级区探索"},
		{ID: "aoe", Name: "AOE清怪", Description: "优先AOE技能，适合多敌人战斗"},
		{ID: "tank", Name: "坦克", Description: "优先嘲讽和减伤技能，适合坦克角色"},
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data: map[string]interface{}{
			"templates":    templates,
			"templateList": templateList,
		},
	})
}

// GetConditionTypes 获取条件类型列表
// GET /api/strategy-condition-types
func (h *StrategyHandlers) GetConditionTypes(c *gin.Context) {
	conditionTypes := []map[string]interface{}{
		// 自身状态
		{"type": "self_hp_percent", "name": "自身HP百分比", "category": "self", "operators": []string{"<", ">", "<=", ">=", "="}, "valueType": "percent"},
		{"type": "self_resource_percent", "name": "自身资源百分比", "category": "self", "operators": []string{"<", ">", "<=", ">=", "="}, "valueType": "percent"},
		{"type": "self_resource", "name": "自身资源值", "category": "self", "operators": []string{"<", ">", "<=", ">=", "="}, "valueType": "number"},
		{"type": "self_has_buff", "name": "自身有Buff", "category": "self", "operators": []string{"=", "!="}, "valueType": "buff_id"},
		{"type": "self_missing_buff", "name": "自身无Buff", "category": "self", "operators": []string{"="}, "valueType": "buff_id"},

		// 敌人状态
		{"type": "alive_enemy_count", "name": "存活敌人数量", "category": "enemy", "operators": []string{"<", ">", "<=", ">=", "="}, "valueType": "number"},
		{"type": "target_hp_percent", "name": "目标HP百分比", "category": "enemy", "operators": []string{"<", ">", "<=", ">=", "="}, "valueType": "percent"},
		{"type": "lowest_enemy_hp_percent", "name": "最低敌人HP", "category": "enemy", "operators": []string{"<", ">"}, "valueType": "percent"},
		{"type": "highest_enemy_hp_percent", "name": "最高敌人HP", "category": "enemy", "operators": []string{"<", ">"}, "valueType": "percent"},
		{"type": "total_enemy_hp_percent", "name": "敌人总HP百分比", "category": "enemy", "operators": []string{"<", ">"}, "valueType": "percent"},

		// 队友状态
		{"type": "alive_ally_count", "name": "存活队友数量", "category": "ally", "operators": []string{"<", ">", "<=", ">=", "="}, "valueType": "number"},
		{"type": "any_ally_hp_percent", "name": "任意队友HP", "category": "ally", "operators": []string{"<", ">"}, "valueType": "percent"},
		{"type": "lowest_ally_hp_percent", "name": "最低队友HP", "category": "ally", "operators": []string{"<", ">"}, "valueType": "percent"},

		// 战斗状态
		{"type": "battle_round", "name": "战斗回合数", "category": "battle", "operators": []string{"<", ">", "<=", ">=", "="}, "valueType": "number"},
		{"type": "skill_ready", "name": "技能可用", "category": "battle", "operators": []string{"="}, "valueType": "skill_id"},
		{"type": "skill_on_cooldown", "name": "技能冷却中", "category": "battle", "operators": []string{"="}, "valueType": "skill_id"},
		{"type": "always", "name": "始终", "category": "battle", "operators": []string{}, "valueType": "none"},
	}

	targetPriorities := []map[string]string{
		{"value": "lowest_hp", "label": "血量最低"},
		{"value": "highest_hp", "label": "血量最高"},
		{"value": "highest_threat", "label": "威胁最高"},
		{"value": "random", "label": "随机"},
		{"value": "max_adjacent", "label": "最大波及数"},
		{"value": "max_splash_damage", "label": "最大波及伤害"},
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data: map[string]interface{}{
			"conditionTypes":   conditionTypes,
			"targetPriorities": targetPriorities,
		},
	})
}




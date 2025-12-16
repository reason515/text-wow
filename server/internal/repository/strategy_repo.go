package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"text-wow/internal/database"
	"text-wow/internal/models"
)

// StrategyRepository 策略数据仓库
type StrategyRepository struct{}

// NewStrategyRepository 创建策略仓库
func NewStrategyRepository() *StrategyRepository {
	return &StrategyRepository{}
}

// Create 创建策略
func (r *StrategyRepository) Create(strategy *models.BattleStrategy) (*models.BattleStrategy, error) {
	skillPriorityJSON, err := json.Marshal(strategy.SkillPriority)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal skill_priority: %w", err)
	}

	conditionalRulesJSON, err := json.Marshal(strategy.ConditionalRules)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal conditional_rules: %w", err)
	}

	skillTargetOverridesJSON, err := json.Marshal(strategy.SkillTargetOverrides)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal skill_target_overrides: %w", err)
	}

	reservedSkillsJSON, err := json.Marshal(strategy.ReservedSkills)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal reserved_skills: %w", err)
	}

	autoTargetSettingsJSON, err := json.Marshal(strategy.AutoTargetSettings)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal auto_target_settings: %w", err)
	}

	result, err := database.DB.Exec(`
		INSERT INTO battle_strategies (
			character_id, name, is_active,
			skill_priority, conditional_rules, target_priority,
			skill_target_overrides, resource_threshold, reserved_skills,
			auto_target_settings, created_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		strategy.CharacterID, strategy.Name, boolToInt(strategy.IsActive),
		string(skillPriorityJSON), string(conditionalRulesJSON), strategy.TargetPriority,
		string(skillTargetOverridesJSON), strategy.ResourceThreshold, string(reservedSkillsJSON),
		string(autoTargetSettingsJSON), time.Now(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to insert strategy: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get last insert id: %w", err)
	}

	strategy.ID = int(id)
	strategy.CreatedAt = time.Now()
	return strategy, nil
}

// GetByID 根据ID获取策略
func (r *StrategyRepository) GetByID(id int) (*models.BattleStrategy, error) {
	strategy := &models.BattleStrategy{}
	var isActive int
	var skillPriorityJSON, conditionalRulesJSON, skillTargetOverridesJSON sql.NullString
	var reservedSkillsJSON, autoTargetSettingsJSON sql.NullString
	var updatedAt sql.NullTime

	err := database.DB.QueryRow(`
		SELECT id, character_id, name, is_active,
		       skill_priority, conditional_rules, target_priority,
		       skill_target_overrides, resource_threshold, reserved_skills,
		       auto_target_settings, created_at, updated_at
		FROM battle_strategies WHERE id = ?`, id,
	).Scan(
		&strategy.ID, &strategy.CharacterID, &strategy.Name, &isActive,
		&skillPriorityJSON, &conditionalRulesJSON, &strategy.TargetPriority,
		&skillTargetOverridesJSON, &strategy.ResourceThreshold, &reservedSkillsJSON,
		&autoTargetSettingsJSON, &strategy.CreatedAt, &updatedAt,
	)
	if err != nil {
		return nil, err
	}

	strategy.IsActive = isActive == 1
	if updatedAt.Valid {
		strategy.UpdatedAt = &updatedAt.Time
	}

	// 解析JSON字段
	if err := r.parseJSONFields(strategy, skillPriorityJSON, conditionalRulesJSON,
		skillTargetOverridesJSON, reservedSkillsJSON, autoTargetSettingsJSON); err != nil {
		return nil, err
	}

	return strategy, nil
}

// GetByCharacterID 获取角色的所有策略
func (r *StrategyRepository) GetByCharacterID(characterID int) ([]*models.BattleStrategy, error) {
	rows, err := database.DB.Query(`
		SELECT id, character_id, name, is_active,
		       skill_priority, conditional_rules, target_priority,
		       skill_target_overrides, resource_threshold, reserved_skills,
		       auto_target_settings, created_at, updated_at
		FROM battle_strategies 
		WHERE character_id = ?
		ORDER BY is_active DESC, created_at DESC`, characterID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var strategies []*models.BattleStrategy
	for rows.Next() {
		strategy := &models.BattleStrategy{}
		var isActive int
		var skillPriorityJSON, conditionalRulesJSON, skillTargetOverridesJSON sql.NullString
		var reservedSkillsJSON, autoTargetSettingsJSON sql.NullString
		var updatedAt sql.NullTime

		err := rows.Scan(
			&strategy.ID, &strategy.CharacterID, &strategy.Name, &isActive,
			&skillPriorityJSON, &conditionalRulesJSON, &strategy.TargetPriority,
			&skillTargetOverridesJSON, &strategy.ResourceThreshold, &reservedSkillsJSON,
			&autoTargetSettingsJSON, &strategy.CreatedAt, &updatedAt,
		)
		if err != nil {
			return nil, err
		}

		strategy.IsActive = isActive == 1
		if updatedAt.Valid {
			strategy.UpdatedAt = &updatedAt.Time
		}

		if err := r.parseJSONFields(strategy, skillPriorityJSON, conditionalRulesJSON,
			skillTargetOverridesJSON, reservedSkillsJSON, autoTargetSettingsJSON); err != nil {
			return nil, err
		}

		strategies = append(strategies, strategy)
	}

	return strategies, nil
}

// GetActiveByCharacterID 获取角色当前激活的策略
func (r *StrategyRepository) GetActiveByCharacterID(characterID int) (*models.BattleStrategy, error) {
	strategy := &models.BattleStrategy{}
	var isActive int
	var skillPriorityJSON, conditionalRulesJSON, skillTargetOverridesJSON sql.NullString
	var reservedSkillsJSON, autoTargetSettingsJSON sql.NullString
	var updatedAt sql.NullTime

	err := database.DB.QueryRow(`
		SELECT id, character_id, name, is_active,
		       skill_priority, conditional_rules, target_priority,
		       skill_target_overrides, resource_threshold, reserved_skills,
		       auto_target_settings, created_at, updated_at
		FROM battle_strategies 
		WHERE character_id = ? AND is_active = 1`, characterID,
	).Scan(
		&strategy.ID, &strategy.CharacterID, &strategy.Name, &isActive,
		&skillPriorityJSON, &conditionalRulesJSON, &strategy.TargetPriority,
		&skillTargetOverridesJSON, &strategy.ResourceThreshold, &reservedSkillsJSON,
		&autoTargetSettingsJSON, &strategy.CreatedAt, &updatedAt,
	)
	if err != nil {
		return nil, err
	}

	strategy.IsActive = isActive == 1
	if updatedAt.Valid {
		strategy.UpdatedAt = &updatedAt.Time
	}

	if err := r.parseJSONFields(strategy, skillPriorityJSON, conditionalRulesJSON,
		skillTargetOverridesJSON, reservedSkillsJSON, autoTargetSettingsJSON); err != nil {
		return nil, err
	}

	return strategy, nil
}

// Update 更新策略
func (r *StrategyRepository) Update(strategy *models.BattleStrategy) error {
	skillPriorityJSON, err := json.Marshal(strategy.SkillPriority)
	if err != nil {
		return fmt.Errorf("failed to marshal skill_priority: %w", err)
	}

	conditionalRulesJSON, err := json.Marshal(strategy.ConditionalRules)
	if err != nil {
		return fmt.Errorf("failed to marshal conditional_rules: %w", err)
	}

	skillTargetOverridesJSON, err := json.Marshal(strategy.SkillTargetOverrides)
	if err != nil {
		return fmt.Errorf("failed to marshal skill_target_overrides: %w", err)
	}

	reservedSkillsJSON, err := json.Marshal(strategy.ReservedSkills)
	if err != nil {
		return fmt.Errorf("failed to marshal reserved_skills: %w", err)
	}

	autoTargetSettingsJSON, err := json.Marshal(strategy.AutoTargetSettings)
	if err != nil {
		return fmt.Errorf("failed to marshal auto_target_settings: %w", err)
	}

	_, err = database.DB.Exec(`
		UPDATE battle_strategies SET
			name = ?, is_active = ?,
			skill_priority = ?, conditional_rules = ?, target_priority = ?,
			skill_target_overrides = ?, resource_threshold = ?, reserved_skills = ?,
			auto_target_settings = ?, updated_at = ?
		WHERE id = ?`,
		strategy.Name, boolToInt(strategy.IsActive),
		string(skillPriorityJSON), string(conditionalRulesJSON), strategy.TargetPriority,
		string(skillTargetOverridesJSON), strategy.ResourceThreshold, string(reservedSkillsJSON),
		string(autoTargetSettingsJSON), time.Now(),
		strategy.ID,
	)
	return err
}

// SetActive 设置策略为激活状态（同时取消同角色其他策略的激活状态）
func (r *StrategyRepository) SetActive(strategyID int, characterID int) error {
	tx, err := database.DB.Begin()
	if err != nil {
		return err
	}

	// 先取消该角色所有策略的激活状态
	_, err = tx.Exec(`UPDATE battle_strategies SET is_active = 0 WHERE character_id = ?`, characterID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// 设置指定策略为激活状态
	_, err = tx.Exec(`UPDATE battle_strategies SET is_active = 1, updated_at = ? WHERE id = ?`, time.Now(), strategyID)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

// Delete 删除策略
func (r *StrategyRepository) Delete(id int) error {
	_, err := database.DB.Exec(`DELETE FROM battle_strategies WHERE id = ?`, id)
	return err
}

// CountByCharacterID 统计角色的策略数量
func (r *StrategyRepository) CountByCharacterID(characterID int) (int, error) {
	var count int
	err := database.DB.QueryRow(`SELECT COUNT(*) FROM battle_strategies WHERE character_id = ?`, characterID).Scan(&count)
	return count, err
}

// parseJSONFields 解析JSON字段
func (r *StrategyRepository) parseJSONFields(strategy *models.BattleStrategy,
	skillPriorityJSON, conditionalRulesJSON, skillTargetOverridesJSON,
	reservedSkillsJSON, autoTargetSettingsJSON sql.NullString) error {

	// 初始化默认值
	strategy.SkillPriority = []string{}
	strategy.ConditionalRules = []models.ConditionalRule{}
	strategy.SkillTargetOverrides = make(map[string]string)
	strategy.ReservedSkills = []models.ReservedSkill{}
	strategy.AutoTargetSettings = models.AutoTargetSettings{
		PositionalAutoOptimize: true,
		ExecuteAutoTarget:      true,
		HealAutoTarget:         true,
	}

	if skillPriorityJSON.Valid && skillPriorityJSON.String != "" {
		if err := json.Unmarshal([]byte(skillPriorityJSON.String), &strategy.SkillPriority); err != nil {
			return fmt.Errorf("failed to unmarshal skill_priority: %w", err)
		}
	}

	if conditionalRulesJSON.Valid && conditionalRulesJSON.String != "" {
		if err := json.Unmarshal([]byte(conditionalRulesJSON.String), &strategy.ConditionalRules); err != nil {
			return fmt.Errorf("failed to unmarshal conditional_rules: %w", err)
		}
	}

	if skillTargetOverridesJSON.Valid && skillTargetOverridesJSON.String != "" {
		if err := json.Unmarshal([]byte(skillTargetOverridesJSON.String), &strategy.SkillTargetOverrides); err != nil {
			return fmt.Errorf("failed to unmarshal skill_target_overrides: %w", err)
		}
	}

	if reservedSkillsJSON.Valid && reservedSkillsJSON.String != "" {
		if err := json.Unmarshal([]byte(reservedSkillsJSON.String), &strategy.ReservedSkills); err != nil {
			return fmt.Errorf("failed to unmarshal reserved_skills: %w", err)
		}
	}

	if autoTargetSettingsJSON.Valid && autoTargetSettingsJSON.String != "" {
		if err := json.Unmarshal([]byte(autoTargetSettingsJSON.String), &strategy.AutoTargetSettings); err != nil {
			return fmt.Errorf("failed to unmarshal auto_target_settings: %w", err)
		}
	}

	return nil
}

// GetDefaultStrategy 获取默认策略模板
func GetDefaultStrategy(characterID int, name string) *models.BattleStrategy {
	return &models.BattleStrategy{
		CharacterID:   characterID,
		Name:          name,
		IsActive:      false,
		SkillPriority: []string{},
		ConditionalRules: []models.ConditionalRule{
			{
				ID:       "rule_1",
				Priority: 1,
				Enabled:  true,
				Condition: models.RuleCondition{
					Type:     "self_hp_percent",
					Operator: "<",
					Value:    30,
				},
				Action: models.RuleAction{
					Type:    "use_skill",
					SkillID: "shield_wall",
					Comment: "低血量时使用盾墙",
				},
			},
		},
		TargetPriority:       "lowest_hp",
		SkillTargetOverrides: make(map[string]string),
		ResourceThreshold:    20,
		ReservedSkills:       []models.ReservedSkill{},
		AutoTargetSettings: models.AutoTargetSettings{
			PositionalAutoOptimize: true,
			ExecuteAutoTarget:      true,
			HealAutoTarget:         true,
		},
	}
}

// GetStrategyTemplates 获取策略模板列表
func GetStrategyTemplates() map[string]*models.BattleStrategy {
	return map[string]*models.BattleStrategy{
		"aggressive": {
			Name:                 "激进输出",
			TargetPriority:       "lowest_hp",
			ResourceThreshold:    10,
			SkillPriority:        []string{},               // 初始化为空数组
			SkillTargetOverrides: make(map[string]string),  // 初始化为空map
			ReservedSkills:       []models.ReservedSkill{}, // 初始化为空数组
			ConditionalRules: []models.ConditionalRule{
				{
					ID: "rule_1", Priority: 1, Enabled: true,
					Condition: models.RuleCondition{Type: "alive_enemy_count", Operator: ">=", Value: 3},
					Action:    models.RuleAction{Type: "use_skill", SkillID: "whirlwind"},
				},
				{
					ID: "rule_2", Priority: 2, Enabled: true,
					Condition: models.RuleCondition{Type: "target_hp_percent", Operator: "<", Value: 20},
					Action:    models.RuleAction{Type: "use_skill", SkillID: "execute"},
				},
			},
			AutoTargetSettings: models.AutoTargetSettings{
				PositionalAutoOptimize: true,
				ExecuteAutoTarget:      true,
				HealAutoTarget:         true,
			},
		},
		"defensive": {
			Name:                 "稳健生存",
			TargetPriority:       "lowest_hp",
			ResourceThreshold:    20,
			SkillPriority:        []string{},
			SkillTargetOverrides: make(map[string]string),
			ReservedSkills:       []models.ReservedSkill{},
			ConditionalRules: []models.ConditionalRule{
				{
					ID: "rule_1", Priority: 1, Enabled: true,
					Condition: models.RuleCondition{Type: "self_hp_percent", Operator: "<", Value: 20},
					Action:    models.RuleAction{Type: "use_skill", SkillID: "last_stand"},
				},
				{
					ID: "rule_2", Priority: 2, Enabled: true,
					Condition: models.RuleCondition{Type: "self_hp_percent", Operator: "<", Value: 30},
					Action:    models.RuleAction{Type: "use_skill", SkillID: "shield_wall"},
				},
				{
					ID: "rule_3", Priority: 3, Enabled: true,
					Condition: models.RuleCondition{Type: "alive_enemy_count", Operator: ">=", Value: 3},
					Action:    models.RuleAction{Type: "use_skill", SkillID: "whirlwind"},
				},
				{
					ID: "rule_4", Priority: 4, Enabled: true,
					Condition: models.RuleCondition{Type: "target_hp_percent", Operator: "<", Value: 20},
					Action:    models.RuleAction{Type: "use_skill", SkillID: "execute"},
				},
				{
					ID: "rule_5", Priority: 5, Enabled: true,
					Condition: models.RuleCondition{Type: "target_hp_percent", Operator: "<", Value: 10},
					Action:    models.RuleAction{Type: "normal_attack", Comment: "残血普攻节省资源"},
				},
			},
			AutoTargetSettings: models.AutoTargetSettings{
				PositionalAutoOptimize: true,
				ExecuteAutoTarget:      true,
				HealAutoTarget:         true,
			},
		},
		"aoe": {
			Name:                 "AOE清怪",
			TargetPriority:       "highest_hp",
			ResourceThreshold:    15,
			SkillPriority:        []string{},
			SkillTargetOverrides: make(map[string]string),
			ReservedSkills:       []models.ReservedSkill{},
			ConditionalRules: []models.ConditionalRule{
				{
					ID: "rule_1", Priority: 1, Enabled: true,
					Condition: models.RuleCondition{Type: "self_hp_percent", Operator: "<", Value: 30},
					Action:    models.RuleAction{Type: "use_skill", SkillID: "shield_wall"},
				},
				{
					ID: "rule_2", Priority: 2, Enabled: true,
					Condition: models.RuleCondition{Type: "alive_enemy_count", Operator: ">=", Value: 3},
					Action:    models.RuleAction{Type: "use_skill", SkillID: "whirlwind"},
				},
				{
					ID: "rule_3", Priority: 3, Enabled: true,
					Condition: models.RuleCondition{Type: "alive_enemy_count", Operator: ">=", Value: 2},
					Action:    models.RuleAction{Type: "use_skill", SkillID: "cleave"},
				},
			},
			AutoTargetSettings: models.AutoTargetSettings{
				PositionalAutoOptimize: true,
				ExecuteAutoTarget:      true,
				HealAutoTarget:         true,
			},
		},
		"tank": {
			Name:                 "坦克",
			TargetPriority:       "highest_threat",
			ResourceThreshold:    10,
			SkillPriority:        []string{},
			SkillTargetOverrides: make(map[string]string),
			ReservedSkills:       []models.ReservedSkill{},
			ConditionalRules: []models.ConditionalRule{
				{
					ID: "rule_1", Priority: 1, Enabled: true,
					Condition: models.RuleCondition{Type: "self_hp_percent", Operator: "<", Value: 20},
					Action:    models.RuleAction{Type: "use_skill", SkillID: "last_stand"},
				},
				{
					ID: "rule_2", Priority: 2, Enabled: true,
					Condition: models.RuleCondition{Type: "self_hp_percent", Operator: "<", Value: 40},
					Action:    models.RuleAction{Type: "use_skill", SkillID: "shield_wall"},
				},
				{
					ID: "rule_3", Priority: 3, Enabled: true,
					Condition: models.RuleCondition{Type: "battle_round", Operator: "=", Value: 1},
					Action:    models.RuleAction{Type: "use_skill", SkillID: "taunt"},
				},
				{
					ID: "rule_4", Priority: 4, Enabled: true,
					Condition: models.RuleCondition{Type: "alive_enemy_count", Operator: ">=", Value: 3},
					Action:    models.RuleAction{Type: "use_skill", SkillID: "challenging_shout"},
				},
			},
			AutoTargetSettings: models.AutoTargetSettings{
				PositionalAutoOptimize: true,
				ExecuteAutoTarget:      false,
				HealAutoTarget:         true,
			},
		},
	}
}








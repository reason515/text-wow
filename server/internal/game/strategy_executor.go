package game

import (
	"database/sql"

	"text-wow/internal/models"
	"text-wow/internal/repository"
)

// StrategyExecutor 策略执行器
type StrategyExecutor struct {
	strategyRepo *repository.StrategyRepository
}

// BattleContext 战斗上下文（用于条件评估）
type BattleContext struct {
	Character    *models.Character
	Enemies      []*models.Monster
	Allies       []*models.Character
	Target       *models.Monster
	CurrentRound int
	SkillManager *SkillManager
	BuffManager  *BuffManager
}

// SkillDecision 技能决策结果
type SkillDecision struct {
	SkillID        string
	IsNormalAttack bool
	TargetIndex    int // 目标索引（敌人）
	Reason         string
}

// NewStrategyExecutor 创建策略执行器
func NewStrategyExecutor() *StrategyExecutor {
	return &StrategyExecutor{
		strategyRepo: repository.NewStrategyRepository(),
	}
}

// GetActiveStrategy 获取角色的激活策略
func (e *StrategyExecutor) GetActiveStrategy(characterID int) *models.BattleStrategy {
	strategy, err := e.strategyRepo.GetActiveByCharacterID(characterID)
	if err == sql.ErrNoRows || strategy == nil {
		return nil
	}
	if err != nil {
		return nil
	}
	return strategy
}

// ExecuteStrategy 执行策略，返回技能决策
func (e *StrategyExecutor) ExecuteStrategy(strategy *models.BattleStrategy, ctx *BattleContext) *SkillDecision {
	if strategy == nil || ctx == nil {
		return nil
	}

	// 1. 检查资源阈值 - 如果资源低于阈值，优先普通攻击积攒资源
	if ctx.Character.Resource < strategy.ResourceThreshold {
		// 除非有紧急条件规则触发
		urgentDecision := e.checkUrgentRules(strategy, ctx)
		if urgentDecision != nil {
			return urgentDecision
		}
		return &SkillDecision{
			IsNormalAttack: true,
			Reason:         "资源低于阈值，使用普通攻击积攒资源",
		}
	}

	// 2. 按优先级检查条件规则
	for _, rule := range strategy.ConditionalRules {
		if !rule.Enabled {
			continue
		}

		if e.evaluateCondition(&rule.Condition, ctx) {
			// 检查技能是否可用
			if rule.Action.Type == "normal_attack" {
				return &SkillDecision{
					IsNormalAttack: true,
					Reason:         "条件规则触发: " + rule.Action.Comment,
				}
			}

			if rule.Action.SkillID != "" {
				// 检查技能是否可用（冷却、资源）
				if e.isSkillAvailable(rule.Action.SkillID, ctx) {
					return &SkillDecision{
						SkillID:     rule.Action.SkillID,
						TargetIndex: e.selectTarget(strategy, ctx, rule.Action.SkillID),
						Reason:      "条件规则触发",
					}
				}
			}
		}
	}

	// 3. 按技能优先级使用技能
	for _, skillID := range strategy.SkillPriority {
		if e.isSkillAvailable(skillID, ctx) {
			// 检查是否是保留技能
			if e.isReservedSkill(strategy, skillID, ctx) {
				continue
			}
			return &SkillDecision{
				SkillID:     skillID,
				TargetIndex: e.selectTarget(strategy, ctx, skillID),
				Reason:      "技能优先级",
			}
		}
	}

	// 4. 没有可用技能，返回 nil（将使用默认逻辑或普通攻击）
	return nil
}

// checkUrgentRules 检查紧急规则（低血量时的保命技能）
func (e *StrategyExecutor) checkUrgentRules(strategy *models.BattleStrategy, ctx *BattleContext) *SkillDecision {
	hpPercent := float64(ctx.Character.HP) / float64(ctx.Character.MaxHP) * 100

	for _, rule := range strategy.ConditionalRules {
		if !rule.Enabled {
			continue
		}

		// 只检查 self_hp_percent < 某值的紧急规则
		if rule.Condition.Type == "self_hp_percent" && rule.Condition.Operator == "<" {
			if hpPercent < rule.Condition.Value {
				if rule.Action.Type == "use_skill" && rule.Action.SkillID != "" {
					if e.isSkillAvailable(rule.Action.SkillID, ctx) {
						return &SkillDecision{
							SkillID:     rule.Action.SkillID,
							TargetIndex: e.selectTarget(strategy, ctx, rule.Action.SkillID),
							Reason:      "紧急规则触发: HP低",
						}
					}
				}
			}
		}
	}
	return nil
}

// evaluateCondition 评估条件是否满足
func (e *StrategyExecutor) evaluateCondition(cond *models.RuleCondition, ctx *BattleContext) bool {
	var currentValue float64

	switch cond.Type {
	case "self_hp_percent":
		currentValue = float64(ctx.Character.HP) / float64(ctx.Character.MaxHP) * 100

	case "self_resource_percent":
		if ctx.Character.MaxResource > 0 {
			currentValue = float64(ctx.Character.Resource) / float64(ctx.Character.MaxResource) * 100
		}

	case "self_resource":
		currentValue = float64(ctx.Character.Resource)

	case "alive_enemy_count":
		count := 0
		for _, enemy := range ctx.Enemies {
			if enemy.HP > 0 {
				count++
			}
		}
		currentValue = float64(count)

	case "target_hp_percent":
		if ctx.Target != nil && ctx.Target.MaxHP > 0 {
			currentValue = float64(ctx.Target.HP) / float64(ctx.Target.MaxHP) * 100
		}

	case "lowest_enemy_hp_percent":
		lowestPercent := 100.0
		for _, enemy := range ctx.Enemies {
			if enemy.HP > 0 && enemy.MaxHP > 0 {
				percent := float64(enemy.HP) / float64(enemy.MaxHP) * 100
				if percent < lowestPercent {
					lowestPercent = percent
				}
			}
		}
		currentValue = lowestPercent

	case "highest_enemy_hp_percent":
		highestPercent := 0.0
		for _, enemy := range ctx.Enemies {
			if enemy.HP > 0 && enemy.MaxHP > 0 {
				percent := float64(enemy.HP) / float64(enemy.MaxHP) * 100
				if percent > highestPercent {
					highestPercent = percent
				}
			}
		}
		currentValue = highestPercent

	case "alive_ally_count":
		count := 0
		for _, ally := range ctx.Allies {
			if ally.HP > 0 {
				count++
			}
		}
		currentValue = float64(count)

	case "lowest_ally_hp_percent":
		lowestPercent := 100.0
		for _, ally := range ctx.Allies {
			if ally.HP > 0 && ally.MaxHP > 0 {
				percent := float64(ally.HP) / float64(ally.MaxHP) * 100
				if percent < lowestPercent {
					lowestPercent = percent
				}
			}
		}
		currentValue = lowestPercent

	case "battle_round":
		currentValue = float64(ctx.CurrentRound)

	case "skill_ready":
		if cond.SkillID != "" && e.isSkillAvailable(cond.SkillID, ctx) {
			return true
		}
		return false

	case "skill_on_cooldown":
		if cond.SkillID != "" && !e.isSkillAvailable(cond.SkillID, ctx) {
			return true
		}
		return false

	case "self_has_buff":
		if cond.BuffID != "" && ctx.BuffManager != nil {
			return ctx.BuffManager.HasBuff(ctx.Character.ID, cond.BuffID)
		}
		return false

	case "self_missing_buff":
		if cond.BuffID != "" && ctx.BuffManager != nil {
			return !ctx.BuffManager.HasBuff(ctx.Character.ID, cond.BuffID)
		}
		return true

	case "always":
		return true

	default:
		return false
	}

	// 比较运算
	return e.compareValues(currentValue, cond.Operator, cond.Value)
}

// compareValues 比较数值
func (e *StrategyExecutor) compareValues(current float64, operator string, target float64) bool {
	switch operator {
	case "<":
		return current < target
	case ">":
		return current > target
	case "<=":
		return current <= target
	case ">=":
		return current >= target
	case "=", "==":
		return current == target
	case "!=":
		return current != target
	default:
		return false
	}
}

// isSkillAvailable 检查技能是否可用
func (e *StrategyExecutor) isSkillAvailable(skillID string, ctx *BattleContext) bool {
	if ctx.SkillManager == nil {
		return false
	}

	// 获取可用技能列表
	available := ctx.SkillManager.GetAvailableSkills(ctx.Character.ID, ctx.Character.Resource, ctx.BuffManager)
	for _, skill := range available {
		if skill.SkillID == skillID || skill.SkillID == "warrior_"+skillID {
			return true
		}
	}
	return false
}

// isReservedSkill 检查是否是保留技能（且条件未满足）
func (e *StrategyExecutor) isReservedSkill(strategy *models.BattleStrategy, skillID string, ctx *BattleContext) bool {
	for _, reserved := range strategy.ReservedSkills {
		if reserved.SkillID == skillID {
			// 如果保留条件未满足，则该技能被保留
			if !e.evaluateCondition(&reserved.Condition, ctx) {
				return true
			}
		}
	}
	return false
}

// selectTarget 选择目标
func (e *StrategyExecutor) selectTarget(strategy *models.BattleStrategy, ctx *BattleContext, skillID string) int {
	// 检查是否有特定技能的目标覆盖
	if override, ok := strategy.SkillTargetOverrides[skillID]; ok {
		return e.selectTargetByPriority(override, ctx, skillID)
	}

	// 使用默认目标策略
	return e.selectTargetByPriority(strategy.TargetPriority, ctx, skillID)
}

// selectTargetByPriority 根据优先级选择目标
func (e *StrategyExecutor) selectTargetByPriority(priority string, ctx *BattleContext, skillID string) int {
	aliveEnemies := make([]int, 0)
	for i, enemy := range ctx.Enemies {
		if enemy.HP > 0 {
			aliveEnemies = append(aliveEnemies, i)
		}
	}

	if len(aliveEnemies) == 0 {
		return 0
	}

	// 智能目标选择
	if ctx.SkillManager != nil {
		skill := e.getSkillByID(skillID, ctx)
		if skill != nil {
			// 检查技能标签，应用智能目标
			// TODO: 根据技能标签自动选择目标
		}
	}

	switch priority {
	case "lowest_hp":
		return e.findLowestHPEnemy(ctx.Enemies, aliveEnemies)
	case "highest_hp":
		return e.findHighestHPEnemy(ctx.Enemies, aliveEnemies)
	case "highest_threat":
		// TODO: 实现威胁值选择
		return aliveEnemies[0]
	case "random":
		if len(aliveEnemies) > 0 {
			return aliveEnemies[len(aliveEnemies)/2] // 简单实现
		}
	case "max_adjacent":
		// TODO: 实现位置技能的最优目标选择
		return e.findCenterEnemy(aliveEnemies)
	}

	return aliveEnemies[0]
}

// findLowestHPEnemy 找到HP最低的敌人
func (e *StrategyExecutor) findLowestHPEnemy(enemies []*models.Monster, aliveIndices []int) int {
	lowestIndex := aliveIndices[0]
	lowestHP := enemies[lowestIndex].HP

	for _, i := range aliveIndices {
		if enemies[i].HP < lowestHP {
			lowestHP = enemies[i].HP
			lowestIndex = i
		}
	}
	return lowestIndex
}

// findHighestHPEnemy 找到HP最高的敌人
func (e *StrategyExecutor) findHighestHPEnemy(enemies []*models.Monster, aliveIndices []int) int {
	highestIndex := aliveIndices[0]
	highestHP := enemies[highestIndex].HP

	for _, i := range aliveIndices {
		if enemies[i].HP > highestHP {
			highestHP = enemies[i].HP
			highestIndex = i
		}
	}
	return highestIndex
}

// findCenterEnemy 找到中间位置的敌人（用于位置技能）
func (e *StrategyExecutor) findCenterEnemy(aliveIndices []int) int {
	if len(aliveIndices) == 0 {
		return 0
	}
	return aliveIndices[len(aliveIndices)/2]
}

// getSkillByID 获取技能信息
func (e *StrategyExecutor) getSkillByID(skillID string, ctx *BattleContext) *models.Skill {
	if ctx.SkillManager == nil {
		return nil
	}

	available := ctx.SkillManager.GetAvailableSkills(ctx.Character.ID, ctx.Character.Resource, ctx.BuffManager)
	for _, skill := range available {
		if skill.SkillID == skillID || skill.SkillID == "warrior_"+skillID {
			return skill.Skill
		}
	}
	return nil
}

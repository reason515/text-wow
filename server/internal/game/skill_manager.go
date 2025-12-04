package game

import (
	"fmt"
	"math/rand"
	"sync"

	"text-wow/internal/models"
	"text-wow/internal/repository"
	"text-wow/internal/service"
)

// SkillManager 技能管理器 - 管理战斗中的技能使用
type SkillManager struct {
	mu              sync.RWMutex
	characterSkills map[int][]*CharacterSkillState // key: characterID
	skillService    *service.SkillService
	skillRepo       *repository.SkillRepository
}

// CharacterSkillState 角色技能状态（包含冷却时间等）
type CharacterSkillState struct {
	SkillID      string
	SkillLevel   int
	CooldownLeft int // 剩余冷却时间（回合数）
	Skill        *models.Skill
	Effect       map[string]interface{} // 计算后的技能效果
}

// NewSkillManager 创建技能管理器
func NewSkillManager() *SkillManager {
	skillRepo := repository.NewSkillRepository()
	skillService := service.NewSkillService(skillRepo, repository.NewCharacterRepository())
	return &SkillManager{
		characterSkills: make(map[int][]*CharacterSkillState),
		skillService:    skillService,
		skillRepo:       skillRepo,
	}
}

// LoadCharacterSkills 加载角色的技能
func (sm *SkillManager) LoadCharacterSkills(characterID int) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	characterSkills, err := sm.skillRepo.GetCharacterSkills(characterID)
	if err != nil {
		return err
	}

	skillStates := make([]*CharacterSkillState, 0)
	for _, cs := range characterSkills {
		skill, err := sm.skillRepo.GetSkillByID(cs.SkillID)
		if err != nil {
			continue // 跳过不存在的技能
		}

		effect := sm.skillService.CalculateSkillEffect(skill, cs.SkillLevel)

		skillState := &CharacterSkillState{
			SkillID:      cs.SkillID,
			SkillLevel:   cs.SkillLevel,
			CooldownLeft: 0,
			Skill:        skill,
			Effect:       effect,
		}
		skillStates = append(skillStates, skillState)
	}

	sm.characterSkills[characterID] = skillStates
	return nil
}

// GetAvailableSkills 获取可用的技能（冷却完成且资源足够）
func (sm *SkillManager) GetAvailableSkills(characterID int, currentResource int) []*CharacterSkillState {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	skillStates, exists := sm.characterSkills[characterID]
	if !exists {
		return []*CharacterSkillState{}
	}

	available := make([]*CharacterSkillState, 0)
	for _, state := range skillStates {
		// 检查冷却时间
		if state.CooldownLeft > 0 {
			continue
		}

		// 检查资源消耗
		resourceCost := state.Skill.ResourceCost
		if resourceCost > currentResource {
			continue
		}

		// 检查特殊条件（如斩杀需要HP<20%）
		if !sm.checkSkillConditions(state, nil) {
			continue
		}

		available = append(available, state)
	}

	return available
}

// SelectBestSkill 选择最佳技能（根据当前情况）
func (sm *SkillManager) SelectBestSkill(characterID int, currentResource int, targetHPPercent float64, hasMultipleEnemies bool) *CharacterSkillState {
	available := sm.GetAvailableSkills(characterID, currentResource)
	if len(available) == 0 {
		return nil
	}

	// 简单的技能选择逻辑
	// 1. 如果有多个敌人，优先选择AOE技能
	if hasMultipleEnemies {
		for _, skill := range available {
			if skill.Skill.TargetType == "enemy_all" {
				return skill
			}
		}
	}

	// 2. 如果目标HP<20%，优先使用斩杀
	if targetHPPercent < 0.20 {
		for _, skill := range available {
			if skill.SkillID == "warrior_execute" {
				return skill
			}
		}
	}

	// 3. 优先使用高伤害技能
	bestSkill := available[0]
	maxDamage := 0.0
	for _, skill := range available {
		if skill.Skill.Type == "attack" {
			multiplier := 1.0
			if m, ok := skill.Effect["damageMultiplier"].(float64); ok {
				multiplier = m
			}
			if multiplier > maxDamage {
				maxDamage = multiplier
				bestSkill = skill
			}
		}
	}

	return bestSkill
}

// UseSkill 使用技能
func (sm *SkillManager) UseSkill(characterID int, skillID string) (*CharacterSkillState, error) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	skillStates, exists := sm.characterSkills[characterID]
	if !exists {
		return nil, fmt.Errorf("character skills not loaded")
	}

	for _, state := range skillStates {
		if state.SkillID == skillID {
			// 设置冷却时间
			if cooldown, ok := state.Effect["cooldown"].(int); ok {
				state.CooldownLeft = cooldown
			} else {
				state.CooldownLeft = state.Skill.Cooldown
			}
			return state, nil
		}
	}

	return nil, fmt.Errorf("skill not found")
}

// TickCooldowns 减少所有技能的冷却时间（每回合调用）
func (sm *SkillManager) TickCooldowns(characterID int) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	skillStates, exists := sm.characterSkills[characterID]
	if !exists {
		return
	}

	for _, state := range skillStates {
		if state.CooldownLeft > 0 {
			state.CooldownLeft--
		}
	}
}

// checkSkillConditions 检查技能使用条件
func (sm *SkillManager) checkSkillConditions(state *CharacterSkillState, target *models.Monster) bool {
	switch state.SkillID {
	case "warrior_execute":
		// 斩杀：目标HP必须<20%
		if target == nil {
			return false
		}
		hpPercent := float64(target.HP) / float64(target.MaxHP)
		return hpPercent < 0.20
	default:
		return true
	}
}

// CalculateSkillDamage 计算技能伤害
func (sm *SkillManager) CalculateSkillDamage(skillState *CharacterSkillState, character *models.Character, target *models.Monster, passiveSkillManager *PassiveSkillManager, buffManager *BuffManager) int {
	skill := skillState.Skill
	effect := skillState.Effect

	// 计算实际攻击力（应用被动技能加成）
	actualAttack := float64(character.Attack)
	if passiveSkillManager != nil {
		// 应用被动技能的攻击力加成（百分比）
		attackModifier := passiveSkillManager.GetPassiveModifier(character.ID, "attack")
		actualAttack = actualAttack * (1.0 + attackModifier/100.0)
		
		// 处理低血量时的攻击力加成（狂暴之心）
		hpPercent := float64(character.HP) / float64(character.MaxHP)
		passives := passiveSkillManager.GetPassiveSkills(character.ID)
		for _, passive := range passives {
			if passive.Passive.EffectType == "stat_mod" && passive.Passive.ID == "warrior_passive_berserker_heart" {
				// 根据等级计算触发阈值（1级50%，5级30%）
				threshold := 0.50 - float64(passive.Level-1)*0.05
				if hpPercent < threshold {
					// 根据等级计算攻击力加成（1级20%，5级60%）
					attackBonus := 20.0 + float64(passive.Level-1)*10.0
					actualAttack = actualAttack * (1.0 + attackBonus/100.0)
				}
			}
		}
	}

	// 应用Buff的攻击力加成（战斗怒吼、狂暴之怒、天神下凡等）
	if buffManager != nil {
		attackBuffValue := buffManager.GetBuffValue(character.ID, "attack")
		if attackBuffValue > 0 {
			// Buff值是百分比加成
			actualAttack = actualAttack * (1.0 + attackBuffValue/100.0)
		}
	}

	// 基础伤害计算
	var baseDamage float64

	switch skill.ID {
	case "warrior_shield_slam":
		// 盾牌猛击：基于攻击力和防御力
		attackMult := 1.0
		defenseMult := 0.5
		if m, ok := effect["attackMultiplier"].(float64); ok {
			attackMult = m
		}
		if m, ok := effect["defenseMultiplier"].(float64); ok {
			defenseMult = m
		}
		// 计算实际防御力（应用被动技能加成）
		actualDefense := float64(character.Defense)
		if passiveSkillManager != nil {
			defenseModifier := passiveSkillManager.GetPassiveModifier(character.ID, "defense")
			actualDefense = actualDefense * (1.0 + defenseModifier/100.0)
		}
		baseDamage = actualAttack*attackMult + actualDefense*defenseMult
	default:
		// 默认：基于攻击力
		multiplier := skill.ScalingRatio
		if m, ok := effect["damageMultiplier"].(float64); ok {
			multiplier = m
		}
		baseDamage = actualAttack * multiplier
	}

	// 应用被动技能的伤害加成（如战斗大师、战争领主等）
	if passiveSkillManager != nil {
		damageModifier := passiveSkillManager.GetPassiveModifier(character.ID, "damage")
		baseDamage = baseDamage * (1.0 + damageModifier/100.0)
	}

	// 计算目标实际防御力（应用Debuff效果）
	actualDefense := float64(target.Defense)
	if buffManager != nil {
		defenseDebuffValue := buffManager.GetEnemyDebuffValue(target.ID, "defense")
		if defenseDebuffValue > 0 {
			// Debuff值是负数（降低防御）
			actualDefense = actualDefense * (1.0 - defenseDebuffValue/100.0)
			if actualDefense < 0 {
				actualDefense = 0
			}
		}
	}
	
	// 应用目标防御
	finalDamage := baseDamage - actualDefense/2.0
	if finalDamage < 1 {
		finalDamage = 1
	}

	// 添加随机波动 ±20%
	variance := finalDamage * 0.2
	finalDamage = finalDamage + (rand.Float64()*2-1)*variance

	return int(finalDamage)
}

// ApplySkillEffects 应用技能效果（buff/debuff等）
func (sm *SkillManager) ApplySkillEffects(skillState *CharacterSkillState, character *models.Character, target *models.Monster) map[string]interface{} {
	effects := make(map[string]interface{})
	skill := skillState.Skill
	effect := skillState.Effect

	switch skill.ID {
	case "warrior_charge":
		// 冲锋：获得怒气，可能眩晕
		if rageGain, ok := effect["rageGain"].(int); ok {
			character.Resource += rageGain
			if character.Resource > character.MaxResource {
				character.Resource = character.MaxResource
			}
			effects["rageGain"] = rageGain
		}
		if stunChance, ok := effect["stunChance"].(float64); ok {
			if rand.Float64() < stunChance {
				effects["stun"] = true
				effects["stunDuration"] = 1
			}
		}
	case "warrior_bloodthirst":
		// 嗜血：恢复生命值
		if healPercent, ok := effect["healPercent"].(float64); ok {
			effects["healPercent"] = healPercent
		}
	case "warrior_mortal_strike":
		// 致死打击：降低治疗效果
		if healingReduction, ok := effect["healingReduction"].(float64); ok {
			effects["healingReduction"] = healingReduction
			if duration, ok := effect["debuffDuration"].(float64); ok {
				effects["debuffDuration"] = duration
			}
		}
	}

	return effects
}

// GetSkillResourceCost 获取技能资源消耗
func (sm *SkillManager) GetSkillResourceCost(skillState *CharacterSkillState) int {
	return skillState.Skill.ResourceCost
}

// ClearCharacterSkills 清除角色的技能状态（战斗结束时）
func (sm *SkillManager) ClearCharacterSkills(characterID int) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	delete(sm.characterSkills, characterID)
}


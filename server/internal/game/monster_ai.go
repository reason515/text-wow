package game

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"sort"

	"text-wow/internal/models"
)

// MonsterAI 怪物AI系统
type MonsterAI struct {
	AIType     string
	Behavior   *AIBehavior
	Monster    *models.Monster
	SkillManager *SkillManager
}

// AIBehavior AI行为配置
type AIBehavior struct {
	TargetPriority  []string `json:"target_priority"`  // 目标优先级: lowest_hp, highest_threat, random等
	SkillPriority   []string `json:"skill_priority"`   // 技能优先级: high_damage, execute, defense, heal等
	DefenseThreshold float64 `json:"defense_threshold"` // 防御阈值（HP百分比）
	RandomFactor   float64 `json:"random_factor"`       // 随机因子（0-1）
	Phases          []AIPhase `json:"phases,omitempty"` // 阶段配置（Boss用）
}

// AIPhase AI阶段配置
type AIPhase struct {
	HPThreshold float64 `json:"hp_threshold"` // HP阈值（0-1）
	Behavior    string  `json:"behavior"`      // 行为类型
	Skills      []string `json:"skills"`       // 可用技能列表
}

// NewMonsterAI 创建怪物AI
func NewMonsterAI(monster *models.Monster, skillManager *SkillManager) (*MonsterAI, error) {
	ai := &MonsterAI{
		AIType:      monster.AIType,
		Monster:     monster,
		SkillManager: skillManager,
	}

	// 解析AI行为配置
	if monster.AIBehavior != "" {
		var behavior AIBehavior
		if err := json.Unmarshal([]byte(monster.AIBehavior), &behavior); err != nil {
			// 如果解析失败，使用默认行为
			ai.Behavior = getDefaultBehavior(monster.Type, monster.AIType)
		} else {
			ai.Behavior = &behavior
		}
	} else {
		// 如果没有配置，使用默认行为
		ai.Behavior = getDefaultBehavior(monster.Type, monster.AIType)
	}

	return ai, nil
}

// getDefaultBehavior 获取默认AI行为
func getDefaultBehavior(monsterType, aiType string) *AIBehavior {
	switch aiType {
	case "aggressive":
		return &AIBehavior{
			TargetPriority:  []string{"lowest_hp", "lowest_defense"},
			SkillPriority:   []string{"high_damage", "execute"},
			DefenseThreshold: 0.2,
			RandomFactor:    0.1,
		}
	case "defensive":
		return &AIBehavior{
			TargetPriority:  []string{"highest_threat"},
			SkillPriority:   []string{"defense", "heal", "attack"},
			DefenseThreshold: 0.5,
			RandomFactor:    0.1,
		}
	case "special":
		// Boss和特殊怪物可能有阶段
		if monsterType == "boss" {
			return &AIBehavior{
				TargetPriority:  []string{"highest_threat", "lowest_hp"},
				SkillPriority:   []string{"special", "high_damage"},
				DefenseThreshold: 0.3,
				RandomFactor:    0.05,
				Phases: []AIPhase{
					{HPThreshold: 1.0, Behavior: "aggressive", Skills: []string{}},
					{HPThreshold: 0.5, Behavior: "defensive", Skills: []string{}},
				},
			}
		}
		return &AIBehavior{
			TargetPriority:  []string{"random", "lowest_hp"},
			SkillPriority:   []string{"balanced"},
			DefenseThreshold: 0.3,
			RandomFactor:    0.2,
		}
	default: // balanced
		return &AIBehavior{
			TargetPriority:  []string{"random", "lowest_hp"},
			SkillPriority:   []string{"balanced"},
			DefenseThreshold: 0.3,
			RandomFactor:    0.3,
		}
	}
}

// SelectTarget 选择目标
// threatTable: 威胁值表（角色ID -> 威胁值），由BattleManager传入
func (ai *MonsterAI) SelectTarget(enemies []*models.Character, threatTable map[int]int) *models.Character {
	if len(enemies) == 0 {
		return nil
	}
	if len(enemies) == 1 {
		return enemies[0]
	}

	// 根据目标优先级选择
	for _, priority := range ai.Behavior.TargetPriority {
		switch priority {
		case "lowest_hp":
			target := ai.selectLowestHP(enemies)
			if target != nil {
				return target
			}
		case "highest_threat":
			target := ai.selectHighestThreat(enemies, threatTable)
			if target != nil {
				return target
			}
		case "lowest_defense":
			target := ai.selectLowestDefense(enemies)
			if target != nil {
				return target
			}
		case "random":
			if rand.Float64() < ai.Behavior.RandomFactor {
				return enemies[rand.Intn(len(enemies))]
			}
		}
	}

	// 默认返回第一个
	return enemies[0]
}

// selectLowestHP 选择HP最低的目标
func (ai *MonsterAI) selectLowestHP(enemies []*models.Character) *models.Character {
	if len(enemies) == 0 {
		return nil
	}

	aliveEnemies := make([]*models.Character, 0)
	for _, e := range enemies {
		if e.HP > 0 {
			aliveEnemies = append(aliveEnemies, e)
		}
	}

	if len(aliveEnemies) == 0 {
		return nil
	}

	target := aliveEnemies[0]
	for _, e := range aliveEnemies {
		if e.HP < target.HP {
			target = e
		}
	}
	return target
}

// selectLowestDefense 选择防御最低的目标
func (ai *MonsterAI) selectLowestDefense(enemies []*models.Character) *models.Character {
	if len(enemies) == 0 {
		return nil
	}

	aliveEnemies := make([]*models.Character, 0)
	for _, e := range enemies {
		if e.HP > 0 {
			aliveEnemies = append(aliveEnemies, e)
		}
	}

	if len(aliveEnemies) == 0 {
		return nil
	}

	target := aliveEnemies[0]
	minDefense := target.PhysicalDefense + target.MagicDefense
	for _, e := range aliveEnemies {
		defense := e.PhysicalDefense + e.MagicDefense
		if defense < minDefense {
			minDefense = defense
			target = e
		}
	}
	return target
}

// selectHighestThreat 选择威胁值最高的目标
// 注意：威胁值表应该由BattleManager维护并传递给AI
func (ai *MonsterAI) selectHighestThreat(enemies []*models.Character, threatTable map[int]int) *models.Character {
	if len(enemies) == 0 {
		return nil
	}

	aliveEnemies := make([]*models.Character, 0)
	for _, e := range enemies {
		if e.HP > 0 {
			aliveEnemies = append(aliveEnemies, e)
		}
	}

	if len(aliveEnemies) == 0 {
		return nil
	}

	// 如果没有威胁值数据，使用默认选择（第一个）
	if threatTable == nil || len(threatTable) == 0 {
		return aliveEnemies[0]
	}

	// 选择威胁值最高的目标
	maxThreat := -1
	var target *models.Character
	for _, enemy := range aliveEnemies {
		threat := threatTable[enemy.ID]
		if threat > maxThreat {
			maxThreat = threat
			target = enemy
		}
	}

	if target == nil {
		// 如果所有目标威胁值都是0或负数，选择第一个
		return aliveEnemies[0]
	}

	return target
}

// SelectSkill 选择技能
func (ai *MonsterAI) SelectSkill(target *models.Character, buffManager *BuffManager) *models.MonsterSkill {
	if len(ai.Monster.MonsterSkills) == 0 {
		return nil
	}

	// 获取当前HP百分比
	hpPercent := float64(ai.Monster.HP) / float64(ai.Monster.MaxHP)

	// 检查是否需要防御
	if hpPercent < ai.Behavior.DefenseThreshold {
		// 优先使用防御或治疗技能
		for _, skill := range ai.Monster.MonsterSkills {
			if skill.CooldownLeft > 0 {
				continue
			}
			if skill.SkillType == "defense" || skill.SkillType == "heal" {
				return skill
			}
		}
	}

	// 根据技能优先级选择
	availableSkills := make([]*models.MonsterSkill, 0)
	for _, skill := range ai.Monster.MonsterSkills {
		if skill.CooldownLeft > 0 {
			continue
		}
		// 检查资源消耗
		if skill.Skill != nil && skill.Skill.ResourceCost > ai.Monster.MP {
			continue
		}
		// 检查使用条件
		if !ai.checkSkillCondition(skill, target) {
			continue
		}
		availableSkills = append(availableSkills, skill)
	}

	if len(availableSkills) == 0 {
		return nil
	}

	// 根据优先级排序
	sort.Slice(availableSkills, func(i, j int) bool {
		return availableSkills[i].Priority > availableSkills[j].Priority
	})

	// 根据技能优先级选择
	for _, priority := range ai.Behavior.SkillPriority {
		for _, skill := range availableSkills {
			if ai.matchesSkillPriority(skill, priority, target) {
				return skill
			}
		}
	}

	// 默认返回优先级最高的可用技能
	return availableSkills[0]
}

// matchesSkillPriority 检查技能是否匹配优先级
func (ai *MonsterAI) matchesSkillPriority(skill *models.MonsterSkill, priority string, target *models.Character) bool {
	if skill.Skill == nil {
		return false
	}

	switch priority {
	case "high_damage":
		return skill.Skill.Type == "attack" && skill.Skill.BaseValue > 50
	case "execute":
		if target != nil {
			hpPercent := float64(target.HP) / float64(target.MaxHP)
			return hpPercent < 0.2 && skill.Skill.Type == "attack"
		}
		return false
	case "defense":
		return skill.SkillType == "defense"
	case "heal":
		return skill.SkillType == "heal"
	case "control":
		return skill.SkillType == "control"
	case "special":
		return skill.SkillType == "special"
	case "balanced":
		return true
	default:
		return false
	}
}

// checkSkillCondition 检查技能使用条件
func (ai *MonsterAI) checkSkillCondition(skill *models.MonsterSkill, target *models.Character) bool {
	if skill.UseCondition == "" {
		return true
	}

	var condition map[string]interface{}
	if err := json.Unmarshal([]byte(skill.UseCondition), &condition); err != nil {
		return true
	}

	// 检查HP条件
	if hpMin, ok := condition["hp_min"].(float64); ok {
		hpPercent := float64(ai.Monster.HP) / float64(ai.Monster.MaxHP)
		if hpPercent < hpMin {
			return false
		}
	}
	if hpMax, ok := condition["hp_max"].(float64); ok {
		hpPercent := float64(ai.Monster.HP) / float64(ai.Monster.MaxHP)
		if hpPercent > hpMax {
			return false
		}
	}

	// 检查目标HP条件
	if target != nil {
		if targetHpMin, ok := condition["target_hp_min"].(float64); ok {
			targetHpPercent := float64(target.HP) / float64(target.MaxHP)
			if targetHpPercent < targetHpMin {
				return false
			}
		}
		if targetHpMax, ok := condition["target_hp_max"].(float64); ok {
			targetHpPercent := float64(target.HP) / float64(target.MaxHP)
			if targetHpPercent > targetHpMax {
				return false
			}
		}
	}

	return true
}

// GetCurrentPhase 获取当前阶段（Boss用）
func (ai *MonsterAI) GetCurrentPhase() *AIPhase {
	if len(ai.Behavior.Phases) == 0 {
		return nil
	}

	hpPercent := float64(ai.Monster.HP) / float64(ai.Monster.MaxHP)

	// 找到当前HP对应的阶段
	for _, phase := range ai.Behavior.Phases {
		if hpPercent >= phase.HPThreshold {
			return &phase
		}
	}

	// 返回最后一个阶段
	return &ai.Behavior.Phases[len(ai.Behavior.Phases)-1]
}

// TickCooldowns 减少技能冷却时间
func (ai *MonsterAI) TickCooldowns() {
	for _, skill := range ai.Monster.MonsterSkills {
		if skill.CooldownLeft > 0 {
			skill.CooldownLeft--
		}
	}
}

// UseSkill 使用技能（设置冷却时间）
func (ai *MonsterAI) UseSkill(skill *models.MonsterSkill) {
	skill.CooldownLeft = skill.Cooldown
	if skill.Skill != nil {
		// 消耗资源
		ai.Monster.MP -= skill.Skill.ResourceCost
		if ai.Monster.MP < 0 {
			ai.Monster.MP = 0
		}
	}
}

// String 返回AI类型字符串
func (ai *MonsterAI) String() string {
	return fmt.Sprintf("MonsterAI{Type: %s, AIType: %s}", ai.Monster.Type, ai.AIType)
}


package game

import (
	"sync"

	"text-wow/internal/models"
	"text-wow/internal/repository"
)

// PassiveSkillManager 被动技能管理器
type PassiveSkillManager struct {
	mu                  sync.RWMutex
	characterPassives   map[int][]*CharacterPassiveState // key: characterID
	passiveSkillRepo    *repository.SkillRepository
}

// CharacterPassiveState 角色被动技能状态
type CharacterPassiveState struct {
	PassiveID   string
	Level       int
	Passive     *models.PassiveSkill
	EffectValue float64 // 计算后的效果数值
}

// NewPassiveSkillManager 创建被动技能管理器
func NewPassiveSkillManager() *PassiveSkillManager {
	return &PassiveSkillManager{
		characterPassives: make(map[int][]*CharacterPassiveState),
		passiveSkillRepo:  repository.NewSkillRepository(),
	}
}

// LoadCharacterPassiveSkills 加载角色的被动技能
func (psm *PassiveSkillManager) LoadCharacterPassiveSkills(characterID int) error {
	psm.mu.Lock()
	defer psm.mu.Unlock()

	characterPassives, err := psm.passiveSkillRepo.GetCharacterPassiveSkills(characterID)
	if err != nil {
		return err
	}

	passiveStates := make([]*CharacterPassiveState, 0)
	for _, cp := range characterPassives {
		passive, err := psm.passiveSkillRepo.GetPassiveSkillByID(cp.PassiveID)
		if err != nil {
			continue // 跳过不存在的被动技能
		}

		// 计算被动技能效果数值（根据等级）
		effectValue := psm.calculatePassiveEffectValue(passive, cp.Level)

		passiveState := &CharacterPassiveState{
			PassiveID:   cp.PassiveID,
			Level:       cp.Level,
			Passive:     passive,
			EffectValue: effectValue,
		}
		passiveStates = append(passiveStates, passiveState)
	}

	psm.characterPassives[characterID] = passiveStates
	return nil
}

// calculatePassiveEffectValue 计算被动技能效果数值
func (psm *PassiveSkillManager) calculatePassiveEffectValue(passive *models.PassiveSkill, level int) float64 {
	// 基础效果值 + (等级-1) * 等级缩放
	return passive.EffectValue + float64(level-1)*passive.LevelScaling
}

// GetPassiveSkills 获取角色的所有被动技能
func (psm *PassiveSkillManager) GetPassiveSkills(characterID int) []*CharacterPassiveState {
	psm.mu.RLock()
	defer psm.mu.RUnlock()

	if passives, exists := psm.characterPassives[characterID]; exists {
		return passives
	}
	return []*CharacterPassiveState{}
}

// GetPassiveEffectValue 获取特定类型的被动技能效果值总和
func (psm *PassiveSkillManager) GetPassiveEffectValue(characterID int, effectType, effectStat string) float64 {
	psm.mu.RLock()
	defer psm.mu.RUnlock()

	totalValue := 0.0
	if passives, exists := psm.characterPassives[characterID]; exists {
		for _, passive := range passives {
			if passive.Passive.EffectType == effectType {
				// 检查effect_stat是否匹配
				if effectStat == "" || passive.Passive.EffectStat == effectStat {
					totalValue += passive.EffectValue
				} else if psm.matchesEffectStat(passive.Passive.EffectStat, effectStat) {
					// 处理多属性被动技能（如"threat_and_defense"）
					totalValue += passive.EffectValue
				}
			}
		}
	}
	return totalValue
}

// matchesEffectStat 检查effect_stat是否匹配（处理多属性情况）
func (psm *PassiveSkillManager) matchesEffectStat(passiveStat, targetStat string) bool {
	// 处理多属性被动技能
	switch passiveStat {
	case "threat_and_defense":
		return targetStat == "threat" || targetStat == "defense"
	case "hp_defense_resistance":
		return targetStat == "hp" || targetStat == "defense" || targetStat == "resistance"
	default:
		return passiveStat == targetStat
	}
}

// GetPassiveModifier 获取被动技能的属性修正值（百分比）
func (psm *PassiveSkillManager) GetPassiveModifier(characterID int, statType string) float64 {
	passives := psm.GetPassiveSkills(characterID)
	
	totalModifier := 0.0
	for _, passive := range passives {
		if passive.Passive.EffectType == "stat_mod" {
			// 检查是否影响该属性
			if psm.matchesEffectStat(passive.Passive.EffectStat, statType) {
				totalModifier += passive.EffectValue
			}
		}
	}
	return totalModifier
}

// ApplyPassiveEffects 应用被动技能效果到角色属性
func (psm *PassiveSkillManager) ApplyPassiveEffects(character *models.Character) {
	passives := psm.GetPassiveSkills(character.ID)

	// 重置临时属性（如果需要）
	// 注意：被动技能的效果是永久的，应该在角色创建/加载时应用
	// 这里主要用于战斗中的动态计算

	for _, passive := range passives {
		switch passive.Passive.EffectType {
		case "stat_mod":
			// 属性修正在计算时应用，不直接修改角色属性
		case "rage_generation":
			// 怒气生成加成在获得怒气时应用
		case "on_hit_heal":
			// 攻击回血在攻击时应用
		case "counter_attack":
			// 反击效果在受到攻击时应用
		case "survival":
			// 生存效果在特定条件下应用
		case "reflect":
			// 反射效果在受到攻击时应用
		case "resistance":
			// 抗性效果在受到控制效果时应用
		}
	}
}

// HasPassiveSkill 检查角色是否有特定的被动技能
func (psm *PassiveSkillManager) HasPassiveSkill(characterID int, passiveID string) bool {
	psm.mu.RLock()
	defer psm.mu.RUnlock()

	if passives, exists := psm.characterPassives[characterID]; exists {
		for _, passive := range passives {
			if passive.PassiveID == passiveID {
				return true
			}
		}
	}
	return false
}

// GetPassiveSkillLevel 获取被动技能的等级
func (psm *PassiveSkillManager) GetPassiveSkillLevel(characterID int, passiveID string) int {
	psm.mu.RLock()
	defer psm.mu.RUnlock()

	if passives, exists := psm.characterPassives[characterID]; exists {
		for _, passive := range passives {
			if passive.PassiveID == passiveID {
				return passive.Level
			}
		}
	}
	return 0
}

// ClearCharacterPassives 清除角色的被动技能状态（战斗结束时）
func (psm *PassiveSkillManager) ClearCharacterPassives(characterID int) {
	psm.mu.Lock()
	defer psm.mu.Unlock()
	delete(psm.characterPassives, characterID)
}














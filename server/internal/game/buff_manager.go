package game

import (
	"sync"
	"time"

	"text-wow/internal/models"
)

// BuffManager Buff/Debuff管理器
type BuffManager struct {
	mu             sync.RWMutex
	characterBuffs map[int]map[string]*BuffInstance   // key: characterID, value: map[buffID]*BuffInstance
	enemyBuffs     map[string]map[string]*BuffInstance // key: enemyID, value: map[buffID]*BuffInstance
}

// BuffInstance Buff/Debuff实例
type BuffInstance struct {
	EffectID    string
	Name        string
	Type        string
	IsBuff      bool
	Duration    int    // 剩余回合数
	Value       float64 // 效果数值
	StatAffected string // 影响的属性
	DamageType  string // DOT伤害类型
	CreatedAt   time.Time
	// DOT/HOT相关字段
	IsDOT       bool    // 是否为持续伤害
	IsHOT       bool    // 是否为持续治疗
	Interval    int    // 触发间隔（回合数，0表示每回合触发）
	LastTick    int    // 上次触发的回合数
}

// NewBuffManager 创建Buff管理器
func NewBuffManager() *BuffManager {
	return &BuffManager{
		characterBuffs: make(map[int]map[string]*BuffInstance),
		enemyBuffs:     make(map[string]map[string]*BuffInstance),
	}
}

// ApplyBuff 应用Buff/Debuff
func (bm *BuffManager) ApplyBuff(characterID int, effectID, name, effectType string, isBuff bool, duration int, value float64, statAffected, damageType string) {
	bm.ApplyBuffWithDOT(characterID, effectID, name, effectType, isBuff, duration, value, statAffected, damageType, false, false, 0)
}

// ApplyBuffWithDOT 应用Buff/Debuff（支持DOT/HOT）
func (bm *BuffManager) ApplyBuffWithDOT(characterID int, effectID, name, effectType string, isBuff bool, duration int, value float64, statAffected, damageType string, isDOT, isHOT bool, interval int) {
	bm.mu.Lock()
	defer bm.mu.Unlock()

	if bm.characterBuffs[characterID] == nil {
		bm.characterBuffs[characterID] = make(map[string]*BuffInstance)
	}

	// 如果已存在相同类型的buff，根据叠加规则处理
	if existing, exists := bm.characterBuffs[characterID][effectID]; exists {
		// 根据叠加规则处理
		stackingRule := bm.getStackingRule(effectID)
		switch stackingRule {
		case "refresh":
			// 刷新持续时间（不叠加）
			if existing.Duration < duration {
				existing.Duration = duration
			}
		case "stack":
			// 叠加：增加层数或刷新持续时间
			existing.Duration = duration
			existing.Value += value // 叠加数值
		case "replace":
			// 替换：完全替换
			existing.Duration = duration
			existing.Value = value
		default:
			// 默认：刷新持续时间
			if existing.Duration < duration {
				existing.Duration = duration
			}
		}
		return
	}

	bm.characterBuffs[characterID][effectID] = &BuffInstance{
		EffectID:     effectID,
		Name:         name,
		Type:         effectType,
		IsBuff:       isBuff,
		Duration:     duration,
		Value:        value,
		StatAffected: statAffected,
		DamageType:   damageType,
		CreatedAt:    time.Now(),
		IsDOT:        isDOT,
		IsHOT:        isHOT,
		Interval:     interval,
		LastTick:     0,
	}
}

// ApplyEnemyDebuff 应用Debuff到敌人
func (bm *BuffManager) ApplyEnemyDebuff(enemyID string, effectID, name, effectType string, duration int, value float64, statAffected, damageType string) {
	bm.ApplyEnemyDebuffWithDOT(enemyID, effectID, name, effectType, duration, value, statAffected, damageType, false, 0)
}

// ApplyEnemyDebuffWithDOT 应用Debuff到敌人（支持DOT）
func (bm *BuffManager) ApplyEnemyDebuffWithDOT(enemyID string, effectID, name, effectType string, duration int, value float64, statAffected, damageType string, isDOT bool, interval int) {
	bm.mu.Lock()
	defer bm.mu.Unlock()

	if bm.enemyBuffs[enemyID] == nil {
		bm.enemyBuffs[enemyID] = make(map[string]*BuffInstance)
	}

	// 如果已存在相同类型的debuff，根据叠加规则处理
	if existing, exists := bm.enemyBuffs[enemyID][effectID]; exists {
		// 根据叠加规则处理
		stackingRule := bm.getStackingRule(effectID)
		switch stackingRule {
		case "refresh":
			// 刷新持续时间（不叠加）
			if existing.Duration < duration {
				existing.Duration = duration
			}
		case "stack":
			// 叠加：增加层数或刷新持续时间
			existing.Duration = duration
			existing.Value += value // 叠加数值
		case "replace":
			// 替换：完全替换
			existing.Duration = duration
			existing.Value = value
		default:
			// 默认：刷新持续时间
			if existing.Duration < duration {
				existing.Duration = duration
			}
		}
		return
	}

	bm.enemyBuffs[enemyID][effectID] = &BuffInstance{
		EffectID:     effectID,
		Name:         name,
		Type:         effectType,
		IsBuff:       false, // 敌人debuff都是debuff
		Duration:     duration,
		Value:        value,
		StatAffected: statAffected,
		DamageType:   damageType,
		CreatedAt:    time.Now(),
		IsDOT:        isDOT,
		IsHOT:        false, // 敌人不会有HOT
		Interval:     interval,
		LastTick:     0,
	}
}

// RemoveBuff 移除Buff/Debuff
func (bm *BuffManager) RemoveBuff(characterID int, effectID string) {
	bm.mu.Lock()
	defer bm.mu.Unlock()

	if buffs, exists := bm.characterBuffs[characterID]; exists {
		delete(buffs, effectID)
	}
}

// ExpiredBuff 过期的Buff信息
type ExpiredBuff struct {
	EffectID string
	Name     string
}

// TickBuffs 减少所有Buff/Debuff的持续时间（每回合调用）
func (bm *BuffManager) TickBuffs(characterID int) []ExpiredBuff {
	bm.mu.Lock()
	defer bm.mu.Unlock()

	expired := make([]ExpiredBuff, 0)
	if buffs, exists := bm.characterBuffs[characterID]; exists {
		for effectID, buff := range buffs {
			buff.Duration--
			if buff.Duration <= 0 {
				expired = append(expired, ExpiredBuff{EffectID: effectID, Name: buff.Name})
				delete(buffs, effectID)
			}
		}
	}

	return expired
}

// TickEnemyDebuffs 减少所有敌人Debuff的持续时间（每回合调用）
func (bm *BuffManager) TickEnemyDebuffs(enemyID string) []string {
	bm.mu.Lock()
	defer bm.mu.Unlock()

	expired := make([]string, 0)
	if debuffs, exists := bm.enemyBuffs[enemyID]; exists {
		for effectID, debuff := range debuffs {
			debuff.Duration--
			if debuff.Duration <= 0 {
				expired = append(expired, effectID)
				delete(debuffs, effectID)
			}
		}
	}

	return expired
}

// GetBuffs 获取角色的所有Buff/Debuff
func (bm *BuffManager) GetBuffs(characterID int) map[string]*BuffInstance {
	bm.mu.RLock()
	defer bm.mu.RUnlock()

	if buffs, exists := bm.characterBuffs[characterID]; exists {
		result := make(map[string]*BuffInstance)
		for k, v := range buffs {
			result[k] = v
		}
		return result
	}
	return make(map[string]*BuffInstance)
}

// GetBuffValue 获取Buff/Debuff的数值
func (bm *BuffManager) GetBuffValue(characterID int, statAffected string) float64 {
	bm.mu.RLock()
	defer bm.mu.RUnlock()

	totalValue := 0.0
	if buffs, exists := bm.characterBuffs[characterID]; exists {
		for _, buff := range buffs {
			if buff.StatAffected == statAffected {
				totalValue += buff.Value
			}
		}
	}
	return totalValue
}

// GetEnemyDebuffValue 获取敌人Debuff的数值
func (bm *BuffManager) GetEnemyDebuffValue(enemyID string, statAffected string) float64 {
	bm.mu.RLock()
	defer bm.mu.RUnlock()

	totalValue := 0.0
	if debuffs, exists := bm.enemyBuffs[enemyID]; exists {
		for _, debuff := range debuffs {
			if debuff.StatAffected == statAffected {
				totalValue += debuff.Value
			}
		}
	}
	return totalValue
}

// GetEnemyDebuffs 获取敌人的所有Debuff
func (bm *BuffManager) GetEnemyDebuffs(enemyID string) map[string]*BuffInstance {
	bm.mu.RLock()
	defer bm.mu.RUnlock()

	if debuffs, exists := bm.enemyBuffs[enemyID]; exists {
		result := make(map[string]*BuffInstance)
		for k, v := range debuffs {
			result[k] = v
		}
		return result
	}
	return make(map[string]*BuffInstance)
}

// HasBuff 检查是否有特定的Buff/Debuff
func (bm *BuffManager) HasBuff(characterID int, effectID string) bool {
	bm.mu.RLock()
	defer bm.mu.RUnlock()

	if buffs, exists := bm.characterBuffs[characterID]; exists {
		_, has := buffs[effectID]
		return has
	}
	return false
}

// ClearBuffs 清除角色的所有Buff/Debuff（战斗结束时）
func (bm *BuffManager) ClearBuffs(characterID int) {
	bm.mu.Lock()
	defer bm.mu.Unlock()
	delete(bm.characterBuffs, characterID)
}

// ClearEnemyDebuffs 清除敌人的所有Debuff（战斗结束时）
func (bm *BuffManager) ClearEnemyDebuffs(enemyID string) {
	bm.mu.Lock()
	defer bm.mu.Unlock()
	delete(bm.enemyBuffs, enemyID)
}

// ClearAllEnemyDebuffs 清除所有敌人的Debuff（战斗结束时）
func (bm *BuffManager) ClearAllEnemyDebuffs() {
	bm.mu.Lock()
	defer bm.mu.Unlock()
	bm.enemyBuffs = make(map[string]map[string]*BuffInstance)
}

// ApplyBuffToCharacter 应用Buff效果到角色属性
func (bm *BuffManager) ApplyBuffToCharacter(character *models.Character) {
	buffs := bm.GetBuffs(character.ID)

	// 重置临时属性加成
	// 注意：这里需要根据实际的buff类型来应用效果
	// 例如：攻击力加成、防御力加成、减伤等

	for _, buff := range buffs {
		switch buff.StatAffected {
		case "attack":
			// 攻击力加成在伤害计算时应用
		case "defense":
			// 防御力加成在伤害计算时应用
		case "damage_taken":
			// 减伤效果在受到伤害时应用
		case "crit_rate":
			// 通用暴击率加成在计算暴击时应用（同时影响物理和法术）
		case "phys_crit_rate":
			// 物理暴击率加成在计算暴击时应用
		case "spell_crit_rate":
			// 法术暴击率加成在计算暴击时应用
		}
	}
}

// CalculateDamageWithBuffs 计算带Buff的伤害
func (bm *BuffManager) CalculateDamageWithBuffs(baseDamage int, characterID int, isPhysical bool) int {
	buffs := bm.GetBuffs(characterID)
	
	damageMultiplier := 1.0
	for _, buff := range buffs {
		if buff.StatAffected == "attack" && buff.IsBuff {
			damageMultiplier += buff.Value / 100.0
		}
	}

	return int(float64(baseDamage) * damageMultiplier)
}

// CalculateDamageTakenWithBuffs 计算带Buff的承受伤害
func (bm *BuffManager) CalculateDamageTakenWithBuffs(baseDamage int, characterID int, isPhysical bool) int {
	buffs := bm.GetBuffs(characterID)
	
	damageReduction := 0.0
	for _, buff := range buffs {
		if buff.StatAffected == "damage_taken" || buff.StatAffected == "physical_damage_taken" {
			if isPhysical && buff.StatAffected == "physical_damage_taken" {
				damageReduction += buff.Value
			} else if buff.StatAffected == "damage_taken" {
				damageReduction += buff.Value
			}
		}
	}

	// 减伤是百分比，负数表示减少
	if damageReduction < 0 {
		reductionMultiplier := 1.0 + (damageReduction / 100.0)
		return int(float64(baseDamage) * reductionMultiplier)
	}

	return baseDamage
}

// getStackingRule 获取Buff的叠加规则
func (bm *BuffManager) getStackingRule(effectID string) string {
	// 定义叠加规则
	// "refresh": 刷新持续时间（不叠加）
	// "stack": 叠加（增加层数或数值）
	// "replace": 替换（完全替换）
	stackingRules := map[string]string{
		"battle_shout":      "refresh", // 战斗怒吼：刷新
		"shield_block":     "refresh", // 盾牌格挡：刷新
		"recklessness_crit": "refresh", // 鲁莽：刷新
		"berserker_rage":   "refresh", // 狂暴之怒：刷新
		"avatar":            "refresh", // 天神下凡：刷新
		"shield_wall":       "refresh", // 盾墙：刷新
		"unbreakable_barrier": "refresh", // 不破壁垒：刷新
		"shield_reflection": "refresh", // 盾牌反射：刷新
		"retaliation":       "refresh", // 反击：刷新
		// DOT/HOT效果可以叠加
		"dot_poison":        "stack",   // 毒药DOT：叠加
		"dot_bleed":         "stack",   // 流血DOT：叠加
		"hot_regen":         "stack",   // 恢复HOT：叠加
	}
	
	if rule, exists := stackingRules[effectID]; exists {
		return rule
	}
	return "refresh" // 默认刷新
}

// ProcessDOTEffects 处理DOT效果（每回合调用）
// 返回: (damageMap, healingMap) damageMap: map[characterID]damage, healingMap: map[characterID]healing
func (bm *BuffManager) ProcessDOTEffects(characterID int, currentRound int) (int, int) {
	bm.mu.Lock()
	defer bm.mu.Unlock()

	damage := 0
	healing := 0

	if buffs, exists := bm.characterBuffs[characterID]; exists {
		for _, buff := range buffs {
			// 检查是否应该触发（根据间隔）
			shouldTick := false
			if buff.Interval == 0 {
				// 每回合触发
				shouldTick = true
			} else {
				// 根据间隔触发
				if buff.LastTick == 0 || (currentRound-buff.LastTick) >= buff.Interval {
					shouldTick = true
					buff.LastTick = currentRound
				}
			}

			if shouldTick {
				if buff.IsDOT {
					// DOT：造成伤害
					damage += int(buff.Value)
				} else if buff.IsHOT {
					// HOT：恢复生命值
					healing += int(buff.Value)
				}
			}
		}
	}

	return damage, healing
}

// ProcessEnemyDOTEffects 处理敌人的DOT效果（每回合调用）
func (bm *BuffManager) ProcessEnemyDOTEffects(enemyID string, currentRound int) int {
	bm.mu.Lock()
	defer bm.mu.Unlock()

	damage := 0

	if debuffs, exists := bm.enemyBuffs[enemyID]; exists {
		for _, debuff := range debuffs {
			// 检查是否应该触发（根据间隔）
			shouldTick := false
			if debuff.Interval == 0 {
				// 每回合触发
				shouldTick = true
			} else {
				// 根据间隔触发
				if debuff.LastTick == 0 || (currentRound-debuff.LastTick) >= debuff.Interval {
					shouldTick = true
					debuff.LastTick = currentRound
				}
			}

			if shouldTick && debuff.IsDOT {
				// DOT：造成伤害
				damage += int(debuff.Value)
			}
		}
	}

	return damage
}


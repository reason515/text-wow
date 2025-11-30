package battle

import (
	"math"
)

// ═══════════════════════════════════════════════════════════
// 战斗单位管理
// ═══════════════════════════════════════════════════════════

// NewBattleUnit 创建战斗单位
func NewBattleUnit(id, name string, level int, isPlayer bool, classID string, baseStats BaseStats) *BattleUnit {
	unit := &BattleUnit{
		ID:           id,
		Name:         name,
		Level:        level,
		IsPlayer:     isPlayer,
		ClassID:      classID,
		BaseStats:    baseStats,
		ResourceType: GetResourceType(classID),
		Skills:       make([]*Skill, 0),
		SkillStates:  make([]*SkillState, 0),
		ActiveEffects: make([]*ActiveEffect, 0),
	}

	// 计算战斗属性
	unit.CalculateCombatStats()
	
	// 初始化HP和资源
	unit.CurrentHP = unit.CombatStats.MaxHP
	unit.CurrentResource = unit.GetInitialResource()

	return unit
}

// GetResourceType 根据职业获取资源类型
func GetResourceType(classID string) ResourceType {
	switch classID {
	case "warrior":
		return ResourceRage
	case "rogue":
		return ResourceEnergy
	default:
		return ResourceMana
	}
}

// GetInitialResource 获取初始资源值
func (u *BattleUnit) GetInitialResource() int {
	switch u.ResourceType {
	case ResourceRage:
		return 0 // 怒气从0开始
	case ResourceEnergy:
		return u.CombatStats.MaxResource // 能量满值开始
	case ResourceMana:
		return u.CombatStats.MaxResource // 法力满值开始
	default:
		return u.CombatStats.MaxResource
	}
}

// CalculateCombatStats 计算战斗属性 (基于数据库设计文档的公式)
func (u *BattleUnit) CalculateCombatStats() {
	stats := &u.CombatStats

	// 最大生命值 = 基础HP + 耐力 × 10
	baseHP := 20 + u.Level*3 // 基础随等级成长
	stats.MaxHP = baseHP + u.BaseStats.Stamina*10

	// 最大资源
	switch u.ResourceType {
	case ResourceRage:
		stats.MaxResource = 100 // 怒气固定100
	case ResourceEnergy:
		stats.MaxResource = 100 // 能量固定100
	case ResourceMana:
		stats.MaxResource = 20 + u.BaseStats.Intellect*2 // 法力基于智力
	}

	// 攻击力 = 基础攻击 + 主属性 × 系数
	baseAttack := 5 + u.Level
	switch u.ClassID {
	case "warrior", "paladin":
		stats.Attack = baseAttack + u.BaseStats.Strength
	case "rogue", "hunter":
		stats.Attack = baseAttack + u.BaseStats.Agility
	case "mage", "warlock", "priest":
		stats.Attack = baseAttack + u.BaseStats.Intellect
	case "druid", "shaman":
		stats.Attack = baseAttack + (u.BaseStats.Intellect+u.BaseStats.Agility)/2
	default:
		stats.Attack = baseAttack + u.BaseStats.Strength
	}

	// 防御力 = 基础防御 + 耐力 × 0.5 + 敏捷 × 0.3
	stats.Defense = 2 + u.Level/2 + int(float64(u.BaseStats.Stamina)*0.5+float64(u.BaseStats.Agility)*0.3)

	// 暴击率 = 敏捷 × 0.05% + 智力 × 0.02% (上限40%)
	stats.CritRate = math.Min(0.4, float64(u.BaseStats.Agility)*0.0005+float64(u.BaseStats.Intellect)*0.0002)

	// 暴击伤害默认1.5倍
	stats.CritDamage = 1.5

	// 闪避率 = 敏捷 × 0.03% (上限30%)
	stats.DodgeRate = math.Min(0.3, float64(u.BaseStats.Agility)*0.0003)

	// 命中率默认100%
	stats.HitRate = 1.0

	// 加速默认0
	stats.HastePct = 0
}

// ═══════════════════════════════════════════════════════════
// 技能管理
// ═══════════════════════════════════════════════════════════

// AddSkill 添加技能
func (u *BattleUnit) AddSkill(skill *Skill) {
	u.Skills = append(u.Skills, skill)
	u.SkillStates = append(u.SkillStates, &SkillState{
		Skill:           skill,
		CurrentCooldown: 0,
	})
}

// GetAvailableSkills 获取可用技能 (未在冷却且资源足够)
func (u *BattleUnit) GetAvailableSkills() []*SkillState {
	available := make([]*SkillState, 0)
	for _, ss := range u.SkillStates {
		if ss.CurrentCooldown <= 0 && u.CurrentResource >= ss.Skill.ResourceCost {
			// 检查沉默状态 - 沉默时不能使用法术类技能
			if u.IsSilenced && isMagicSkill(ss.Skill) {
				continue
			}
			available = append(available, ss)
		}
	}
	return available
}

// isMagicSkill 判断是否为法术技能
func isMagicSkill(skill *Skill) bool {
	// 物理伤害技能不受沉默影响
	return skill.DamageType != DamagePhysical
}

// UseSkill 使用技能
func (u *BattleUnit) UseSkill(skillState *SkillState) bool {
	if skillState.CurrentCooldown > 0 {
		return false
	}
	if u.CurrentResource < skillState.Skill.ResourceCost {
		return false
	}

	// 消耗资源
	u.CurrentResource -= skillState.Skill.ResourceCost

	// 设置冷却
	skillState.CurrentCooldown = skillState.Skill.Cooldown

	return true
}

// TickCooldowns 回合结束时减少所有技能冷却
func (u *BattleUnit) TickCooldowns() {
	for _, ss := range u.SkillStates {
		if ss.CurrentCooldown > 0 {
			ss.CurrentCooldown--
		}
	}
}

// ═══════════════════════════════════════════════════════════
// 资源管理
// ═══════════════════════════════════════════════════════════

// RegenerateResource 回合开始时资源恢复
func (u *BattleUnit) RegenerateResource() int {
	var regen int

	switch u.ResourceType {
	case ResourceRage:
		// 怒气: 攻击获得，这里只处理自然衰减
		// 战斗中不衰减，但也不自然恢复
		return 0
	case ResourceEnergy:
		// 能量: 每回合恢复20点
		regen = 20
	case ResourceMana:
		// 法力: 每回合恢复 精神 × 0.5% 的最大法力
		regen = int(float64(u.BaseStats.Spirit) * 0.005 * float64(u.CombatStats.MaxResource))
		if regen < 1 {
			regen = 1
		}
	}

	oldResource := u.CurrentResource
	u.CurrentResource = min(u.CurrentResource+regen, u.CombatStats.MaxResource)
	return u.CurrentResource - oldResource
}

// GainRage 获得怒气 (战士专用)
func (u *BattleUnit) GainRage(amount int) int {
	if u.ResourceType != ResourceRage {
		return 0
	}
	oldRage := u.CurrentResource
	u.CurrentResource = min(u.CurrentResource+amount, u.CombatStats.MaxResource)
	return u.CurrentResource - oldRage
}

// GainRageFromDamage 受到伤害时获得怒气
func (u *BattleUnit) GainRageFromDamage(damage int) int {
	if u.ResourceType != ResourceRage {
		return 0
	}
	// 受到伤害获得怒气: 伤害/最大HP × 50
	rageGain := int(float64(damage) / float64(u.CombatStats.MaxHP) * 50)
	if rageGain < 1 {
		rageGain = 1
	}
	return u.GainRage(rageGain)
}

// GainRageFromAttack 攻击时获得怒气
func (u *BattleUnit) GainRageFromAttack() int {
	if u.ResourceType != ResourceRage {
		return 0
	}
	// 每次攻击获得5点怒气
	return u.GainRage(5)
}

// ═══════════════════════════════════════════════════════════
// 生命值管理
// ═══════════════════════════════════════════════════════════

// TakeDamage 受到伤害 (返回实际伤害)
func (u *BattleUnit) TakeDamage(damage int) int {
	if damage <= 0 {
		return 0
	}

	// 先检查护盾
	shieldAbsorbed := u.AbsorbDamageWithShield(damage)
	actualDamage := damage - shieldAbsorbed

	if actualDamage <= 0 {
		return 0
	}

	// 扣除生命值
	u.CurrentHP -= actualDamage
	u.DamageTaken += actualDamage

	// 战士受伤获得怒气
	u.GainRageFromDamage(actualDamage)

	// 检查死亡
	if u.CurrentHP <= 0 {
		u.CurrentHP = 0
		u.IsDead = true
	}

	return actualDamage
}

// Heal 治疗 (返回实际治疗量)
func (u *BattleUnit) Heal(amount int) int {
	if amount <= 0 || u.IsDead {
		return 0
	}

	oldHP := u.CurrentHP
	u.CurrentHP = min(u.CurrentHP+amount, u.CombatStats.MaxHP)
	actualHeal := u.CurrentHP - oldHP
	u.HealingTaken += actualHeal

	return actualHeal
}

// AbsorbDamageWithShield 护盾吸收伤害
func (u *BattleUnit) AbsorbDamageWithShield(damage int) int {
	absorbed := 0
	for _, effect := range u.ActiveEffects {
		if effect.Effect.Type == EffectShield && effect.ShieldAmount > 0 {
			if effect.ShieldAmount >= damage-absorbed {
				effect.ShieldAmount -= (damage - absorbed)
				absorbed = damage
				break
			} else {
				absorbed += effect.ShieldAmount
				effect.ShieldAmount = 0
			}
		}
	}
	return absorbed
}

// ═══════════════════════════════════════════════════════════
// 效果管理
// ═══════════════════════════════════════════════════════════

// ApplyEffect 应用效果
func (u *BattleUnit) ApplyEffect(effect *Effect, sourceID string) *ActiveEffect {
	// 检查是否已存在该效果
	for _, ae := range u.ActiveEffects {
		if ae.Effect.ID == effect.ID {
			if effect.IsStackable && ae.Stacks < effect.MaxStacks {
				ae.Stacks++
				ae.RemainingTurns = effect.Duration // 刷新持续时间
			} else {
				ae.RemainingTurns = effect.Duration // 仅刷新持续时间
			}
			return ae
		}
	}

	// 创建新效果
	activeEffect := &ActiveEffect{
		Effect:         effect,
		SourceID:       sourceID,
		RemainingTurns: effect.Duration,
		Stacks:         1,
	}

	// 护盾效果设置初始护盾量
	if effect.Type == EffectShield {
		activeEffect.ShieldAmount = int(effect.Value)
	}

	u.ActiveEffects = append(u.ActiveEffects, activeEffect)

	// 更新控制状态
	u.UpdateControlStates()

	return activeEffect
}

// RemoveEffect 移除效果
func (u *BattleUnit) RemoveEffect(effectID string) {
	for i, ae := range u.ActiveEffects {
		if ae.Effect.ID == effectID {
			u.ActiveEffects = append(u.ActiveEffects[:i], u.ActiveEffects[i+1:]...)
			break
		}
	}
	u.UpdateControlStates()
}

// TickEffects 回合开始时处理效果 (返回DOT/HOT结果)
func (u *BattleUnit) TickEffects() []ActionResult {
	results := make([]ActionResult, 0)
	expiredEffects := make([]string, 0)

	for _, ae := range u.ActiveEffects {
		// 处理DOT
		if ae.Effect.Type == EffectDOT {
			damage := int(ae.Effect.Value) * ae.Stacks
			u.CurrentHP -= damage
			u.DamageTaken += damage
			results = append(results, ActionResult{
				Type:       ActionDOTTick,
				SourceID:   ae.SourceID,
				TargetID:   u.ID,
				TargetName: u.Name,
				EffectID:   ae.Effect.ID,
				EffectName: ae.Effect.Name,
				Value:      damage,
			})

			if u.CurrentHP <= 0 {
				u.CurrentHP = 0
				u.IsDead = true
			}
		}

		// 处理HOT
		if ae.Effect.Type == EffectHOT {
			heal := int(ae.Effect.Value) * ae.Stacks
			actualHeal := u.Heal(heal)
			results = append(results, ActionResult{
				Type:       ActionHOTTick,
				SourceID:   ae.SourceID,
				TargetID:   u.ID,
				TargetName: u.Name,
				EffectID:   ae.Effect.ID,
				EffectName: ae.Effect.Name,
				Value:      actualHeal,
			})
		}

		// 减少持续时间
		ae.RemainingTurns--
		if ae.RemainingTurns <= 0 {
			expiredEffects = append(expiredEffects, ae.Effect.ID)
			results = append(results, ActionResult{
				Type:       ActionEffectExpire,
				TargetID:   u.ID,
				TargetName: u.Name,
				EffectID:   ae.Effect.ID,
				EffectName: ae.Effect.Name,
			})
		}
	}

	// 移除过期效果
	for _, effectID := range expiredEffects {
		u.RemoveEffect(effectID)
	}

	return results
}

// UpdateControlStates 更新控制状态
func (u *BattleUnit) UpdateControlStates() {
	u.IsStunned = false
	u.IsSilenced = false
	u.TauntedBy = ""

	for _, ae := range u.ActiveEffects {
		switch ae.Effect.Type {
		case EffectStun, EffectRoot:
			u.IsStunned = true
		case EffectSilence:
			u.IsSilenced = true
		case EffectTaunt:
			u.TauntedBy = ae.SourceID
		}
	}
}

// GetStatModifier 获取属性修正值 (来自所有Buff/Debuff)
func (u *BattleUnit) GetStatModifier(stat string) float64 {
	modifier := 0.0
	for _, ae := range u.ActiveEffects {
		if ae.Effect.Type == EffectStatMod && ae.Effect.StatAffected == stat {
			if ae.Effect.ValueType == "percent" {
				modifier += ae.Effect.Value * float64(ae.Stacks) / 100.0
			}
		}
	}
	return modifier
}

// GetModifiedStat 获取修正后的属性值
func (u *BattleUnit) GetModifiedAttack() int {
	modifier := u.GetStatModifier("attack")
	return int(float64(u.CombatStats.Attack) * (1.0 + modifier))
}

func (u *BattleUnit) GetModifiedDefense() int {
	modifier := u.GetStatModifier("defense")
	return int(float64(u.CombatStats.Defense) * (1.0 + modifier))
}

func (u *BattleUnit) GetModifiedCritRate() float64 {
	modifier := u.GetStatModifier("crit_rate")
	return math.Min(0.4, u.CombatStats.CritRate+modifier)
}

func (u *BattleUnit) GetModifiedDodgeRate() float64 {
	modifier := u.GetStatModifier("dodge_rate")
	return math.Min(0.5, u.CombatStats.DodgeRate+modifier)
}

// ═══════════════════════════════════════════════════════════
// 辅助函数
// ═══════════════════════════════════════════════════════════

// IsAlive 是否存活
func (u *BattleUnit) IsAlive() bool {
	return !u.IsDead && u.CurrentHP > 0
}

// CanAct 是否可以行动
func (u *BattleUnit) CanAct() bool {
	return u.IsAlive() && !u.IsStunned
}

// GetHPPercent 获取HP百分比
func (u *BattleUnit) GetHPPercent() float64 {
	return float64(u.CurrentHP) / float64(u.CombatStats.MaxHP)
}

// min 辅助函数
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}



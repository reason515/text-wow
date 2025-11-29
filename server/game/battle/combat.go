package battle

import (
	"fmt"
	"math/rand"
	"sort"
	"time"
)

// ═══════════════════════════════════════════════════════════
// 战斗引擎
// ═══════════════════════════════════════════════════════════

// CombatEngine 战斗引擎
type CombatEngine struct {
	battle    *Battle
	rng       *rand.Rand
	turnLimit int
}

// NewCombatEngine 创建战斗引擎
func NewCombatEngine() *CombatEngine {
	return &CombatEngine{
		rng:       rand.New(rand.NewSource(time.Now().UnixNano())),
		turnLimit: 100, // 默认100回合限制
	}
}

// StartBattle 开始战斗
func (e *CombatEngine) StartBattle(playerTeam, enemyTeam *Team) *Battle {
	e.battle = &Battle{
		ID:          fmt.Sprintf("battle_%d", time.Now().UnixNano()),
		State:       BattleStateInProgress,
		CurrentTurn: 0,
		MaxTurns:    e.turnLimit,
		PlayerTeam:  playerTeam,
		EnemyTeam:   enemyTeam,
		ActionLog:   make([]ActionLog, 0),
		StartTime:   time.Now(),
	}

	// 确定行动顺序 (基于敏捷)
	e.DetermineInitiative()

	return e.battle
}

// DetermineInitiative 确定行动顺序
func (e *CombatEngine) DetermineInitiative() {
	allUnits := make([]*BattleUnit, 0)
	allUnits = append(allUnits, e.battle.PlayerTeam.Units...)
	allUnits = append(allUnits, e.battle.EnemyTeam.Units...)

	// 按敏捷排序，敏捷高的先行动
	sort.Slice(allUnits, func(i, j int) bool {
		return allUnits[i].BaseStats.Agility > allUnits[j].BaseStats.Agility
	})

	e.battle.TurnOrder = allUnits
}

// ExecuteTurn 执行一个完整回合
func (e *CombatEngine) ExecuteTurn() *ActionLog {
	if e.battle.State != BattleStateInProgress {
		return nil
	}

	e.battle.CurrentTurn++
	turnLog := &ActionLog{
		Turn:      e.battle.CurrentTurn,
		Timestamp: time.Now(),
		Results:   make([]ActionResult, 0),
	}

	// 回合开始: 处理所有单位的效果
	for _, unit := range e.battle.TurnOrder {
		if unit.IsAlive() {
			// 资源恢复
			unit.RegenerateResource()

			// 效果Tick (DOT/HOT)
			effectResults := unit.TickEffects()
			turnLog.Results = append(turnLog.Results, effectResults...)

			// 检查是否因DOT死亡
			if unit.IsDead {
				turnLog.Results = append(turnLog.Results, ActionResult{
					Type:       ActionDeath,
					TargetID:   unit.ID,
					TargetName: unit.Name,
				})
			}
		}
	}

	// 检查战斗是否结束
	if e.CheckBattleEnd() {
		e.battle.ActionLog = append(e.battle.ActionLog, *turnLog)
		return turnLog
	}

	// 每个单位执行行动
	for _, unit := range e.battle.TurnOrder {
		if !unit.CanAct() {
			if unit.IsAlive() && unit.IsStunned {
				turnLog.Results = append(turnLog.Results, ActionResult{
					Type:       ActionSkipped,
					SourceID:   unit.ID,
					SourceName: unit.Name,
				})
			}
			continue
		}

		// 获取目标队伍
		var enemyTeam, allyTeam *Team
		if unit.IsPlayer {
			enemyTeam = e.battle.EnemyTeam
			allyTeam = e.battle.PlayerTeam
		} else {
			enemyTeam = e.battle.PlayerTeam
			allyTeam = e.battle.EnemyTeam
		}

		// AI选择技能和目标
		action := e.SelectAction(unit, allyTeam, enemyTeam)

		// 执行行动
		if action != nil {
			results := e.ExecuteAction(unit, action, allyTeam, enemyTeam)
			turnLog.Results = append(turnLog.Results, results...)
		}

		// 减少技能冷却
		unit.TickCooldowns()

		// 检查战斗是否结束
		if e.CheckBattleEnd() {
			break
		}
	}

	// 检查回合上限
	if e.battle.CurrentTurn >= e.battle.MaxTurns {
		e.battle.State = BattleStateDraw
		e.battle.EndTime = time.Now()
	}

	e.battle.ActionLog = append(e.battle.ActionLog, *turnLog)
	return turnLog
}

// ═══════════════════════════════════════════════════════════
// AI 决策
// ═══════════════════════════════════════════════════════════

// ActionChoice AI选择的行动
type ActionChoice struct {
	SkillState *SkillState
	Targets    []*BattleUnit
}

// SelectAction AI选择行动
func (e *CombatEngine) SelectAction(unit *BattleUnit, allyTeam, enemyTeam *Team) *ActionChoice {
	availableSkills := unit.GetAvailableSkills()

	// 如果没有可用技能，使用普通攻击
	if len(availableSkills) == 0 {
		return e.CreateBasicAttackAction(unit, enemyTeam)
	}

	// 简单AI策略:
	// 1. 如果是治疗职业且有队友血量低于50%，优先治疗
	// 2. 否则选择伤害最高的可用技能
	// 3. 如果被嘲讽，必须攻击嘲讽者

	// 检查是否需要治疗
	lowHPAlly := e.FindLowHPAlly(allyTeam, 0.5)
	if lowHPAlly != nil {
		healSkill := e.FindHealSkill(availableSkills)
		if healSkill != nil {
			return &ActionChoice{
				SkillState: healSkill,
				Targets:    []*BattleUnit{lowHPAlly},
			}
		}
	}

	// 选择攻击技能
	attackSkill := e.FindBestAttackSkill(availableSkills)
	if attackSkill != nil {
		targets := e.SelectTargets(unit, attackSkill.Skill, allyTeam, enemyTeam)
		if len(targets) > 0 {
			return &ActionChoice{
				SkillState: attackSkill,
				Targets:    targets,
			}
		}
	}

	// 后备: 普通攻击
	return e.CreateBasicAttackAction(unit, enemyTeam)
}

// CreateBasicAttackAction 创建普通攻击行动
func (e *CombatEngine) CreateBasicAttackAction(unit *BattleUnit, enemyTeam *Team) *ActionChoice {
	basicAttack := &Skill{
		ID:           "basic_attack",
		Name:         "普通攻击",
		Type:         SkillAttack,
		TargetType:   TargetEnemy,
		DamageType:   DamagePhysical,
		BaseValue:    0,
		ScalingStat:  "strength",
		ScalingRatio: 0.5,
		ResourceCost: 0,
		Cooldown:     0,
	}

	target := e.SelectSingleTarget(unit, enemyTeam)
	if target == nil {
		return nil
	}

	return &ActionChoice{
		SkillState: &SkillState{Skill: basicAttack, CurrentCooldown: 0},
		Targets:    []*BattleUnit{target},
	}
}

// FindLowHPAlly 找到血量低的队友
func (e *CombatEngine) FindLowHPAlly(team *Team, threshold float64) *BattleUnit {
	var lowest *BattleUnit
	lowestPct := 1.0

	for _, unit := range team.Units {
		if unit.IsAlive() {
			pct := unit.GetHPPercent()
			if pct < threshold && pct < lowestPct {
				lowest = unit
				lowestPct = pct
			}
		}
	}
	return lowest
}

// FindHealSkill 找到治疗技能
func (e *CombatEngine) FindHealSkill(skills []*SkillState) *SkillState {
	for _, ss := range skills {
		if ss.Skill.Type == SkillHeal || ss.Skill.Type == SkillHOT {
			return ss
		}
	}
	return nil
}

// FindBestAttackSkill 找到最佳攻击技能
func (e *CombatEngine) FindBestAttackSkill(skills []*SkillState) *SkillState {
	var best *SkillState
	bestValue := 0

	for _, ss := range skills {
		if ss.Skill.Type == SkillAttack || ss.Skill.Type == SkillDOT {
			if ss.Skill.BaseValue > bestValue {
				best = ss
				bestValue = ss.Skill.BaseValue
			}
		}
	}

	// 如果没有攻击技能，返回任意可用技能
	if best == nil && len(skills) > 0 {
		return skills[0]
	}

	return best
}

// SelectSingleTarget 选择单个目标
func (e *CombatEngine) SelectSingleTarget(unit *BattleUnit, enemyTeam *Team) *BattleUnit {
	// 如果被嘲讽，必须攻击嘲讽者
	if unit.TauntedBy != "" {
		for _, enemy := range enemyTeam.Units {
			if enemy.ID == unit.TauntedBy && enemy.IsAlive() {
				return enemy
			}
		}
	}

	// 默认攻击第一个存活的敌人
	for _, enemy := range enemyTeam.Units {
		if enemy.IsAlive() {
			return enemy
		}
	}
	return nil
}

// SelectTargets 根据技能目标类型选择目标
func (e *CombatEngine) SelectTargets(unit *BattleUnit, skill *Skill, allyTeam, enemyTeam *Team) []*BattleUnit {
	targets := make([]*BattleUnit, 0)

	switch skill.TargetType {
	case TargetSelf:
		targets = append(targets, unit)

	case TargetAlly:
		for _, ally := range allyTeam.Units {
			if ally.IsAlive() && ally.ID != unit.ID {
				targets = append(targets, ally)
				break
			}
		}
		if len(targets) == 0 {
			targets = append(targets, unit) // 没有其他队友就选自己
		}

	case TargetAllyAll:
		for _, ally := range allyTeam.Units {
			if ally.IsAlive() {
				targets = append(targets, ally)
			}
		}

	case TargetAllyLowest:
		var lowest *BattleUnit
		lowestHP := 1.0
		for _, ally := range allyTeam.Units {
			if ally.IsAlive() && ally.GetHPPercent() < lowestHP {
				lowest = ally
				lowestHP = ally.GetHPPercent()
			}
		}
		if lowest != nil {
			targets = append(targets, lowest)
		}

	case TargetEnemy:
		target := e.SelectSingleTarget(unit, enemyTeam)
		if target != nil {
			targets = append(targets, target)
		}

	case TargetEnemyAll:
		for _, enemy := range enemyTeam.Units {
			if enemy.IsAlive() {
				targets = append(targets, enemy)
			}
		}

	case TargetEnemyRandom:
		aliveEnemies := make([]*BattleUnit, 0)
		for _, enemy := range enemyTeam.Units {
			if enemy.IsAlive() {
				aliveEnemies = append(aliveEnemies, enemy)
			}
		}
		if len(aliveEnemies) > 0 {
			targets = append(targets, aliveEnemies[e.rng.Intn(len(aliveEnemies))])
		}

	case TargetEnemyLowest:
		var lowest *BattleUnit
		lowestHP := 1.0
		for _, enemy := range enemyTeam.Units {
			if enemy.IsAlive() && enemy.GetHPPercent() < lowestHP {
				lowest = enemy
				lowestHP = enemy.GetHPPercent()
			}
		}
		if lowest != nil {
			targets = append(targets, lowest)
		}
	}

	return targets
}

// ═══════════════════════════════════════════════════════════
// 行动执行
// ═══════════════════════════════════════════════════════════

// ExecuteAction 执行行动
func (e *CombatEngine) ExecuteAction(unit *BattleUnit, action *ActionChoice, allyTeam, enemyTeam *Team) []ActionResult {
	results := make([]ActionResult, 0)
	skill := action.SkillState.Skill

	// 使用技能 (消耗资源，设置冷却)
	if action.SkillState.Skill.ID != "basic_attack" {
		unit.UseSkill(action.SkillState)
	}

	// 战士攻击获得怒气
	if skill.Type == SkillAttack {
		unit.GainRageFromAttack()
	}

	for _, target := range action.Targets {
		result := e.ExecuteSkillOnTarget(unit, skill, target)
		results = append(results, result)

		// 检查目标是否死亡
		if target.IsDead {
			results = append(results, ActionResult{
				Type:       ActionDeath,
				TargetID:   target.ID,
				TargetName: target.Name,
			})
		}
	}

	return results
}

// ExecuteSkillOnTarget 对目标执行技能
func (e *CombatEngine) ExecuteSkillOnTarget(source *BattleUnit, skill *Skill, target *BattleUnit) ActionResult {
	result := ActionResult{
		SourceID:   source.ID,
		SourceName: source.Name,
		TargetID:   target.ID,
		TargetName: target.Name,
		SkillID:    skill.ID,
		SkillName:  skill.Name,
	}

	switch skill.Type {
	case SkillAttack:
		result.Type = ActionSkill
		damage, isCrit, isDodged := e.CalculateDamage(source, target, skill)
		result.Value = damage
		result.IsCrit = isCrit
		result.IsDodged = isDodged
		if !isDodged {
			target.TakeDamage(damage)
			source.DamageDealt += damage
		}

	case SkillHeal:
		result.Type = ActionSkill
		healAmount := e.CalculateHeal(source, skill)
		actualHeal := target.Heal(healAmount)
		result.Value = actualHeal
		source.HealingDone += actualHeal

	case SkillBuff, SkillDebuff:
		result.Type = ActionEffectApply
		if skill.EffectID != "" {
			effect := e.GetEffect(skill.EffectID)
			if effect != nil && e.rng.Float64() < skill.EffectChance {
				target.ApplyEffect(effect, source.ID)
				result.EffectID = effect.ID
				result.EffectName = effect.Name
			}
		}

	case SkillDOT:
		result.Type = ActionSkill
		// DOT技能先造成初始伤害
		damage, isCrit, isDodged := e.CalculateDamage(source, target, skill)
		result.Value = damage
		result.IsCrit = isCrit
		result.IsDodged = isDodged
		if !isDodged {
			target.TakeDamage(damage)
			source.DamageDealt += damage
			// 然后应用DOT效果
			if skill.EffectID != "" {
				effect := e.GetEffect(skill.EffectID)
				if effect != nil {
					target.ApplyEffect(effect, source.ID)
					result.EffectID = effect.ID
					result.EffectName = effect.Name
				}
			}
		}

	case SkillHOT:
		result.Type = ActionSkill
		// HOT技能先进行初始治疗
		healAmount := e.CalculateHeal(source, skill)
		actualHeal := target.Heal(healAmount)
		result.Value = actualHeal
		source.HealingDone += actualHeal
		// 然后应用HOT效果
		if skill.EffectID != "" {
			effect := e.GetEffect(skill.EffectID)
			if effect != nil {
				target.ApplyEffect(effect, source.ID)
				result.EffectID = effect.ID
				result.EffectName = effect.Name
			}
		}

	case SkillShield:
		result.Type = ActionEffectApply
		if skill.EffectID != "" {
			effect := e.GetEffect(skill.EffectID)
			if effect != nil {
				// 护盾值可能受属性加成
				shieldAmount := skill.BaseValue + int(float64(source.BaseStats.Intellect)*skill.ScalingRatio)
				modifiedEffect := *effect
				modifiedEffect.Value = float64(shieldAmount)
				target.ApplyEffect(&modifiedEffect, source.ID)
				result.EffectID = effect.ID
				result.EffectName = effect.Name
				result.Value = shieldAmount
			}
		}

	case SkillControl:
		result.Type = ActionEffectApply
		if skill.EffectID != "" && e.rng.Float64() < skill.EffectChance {
			effect := e.GetEffect(skill.EffectID)
			if effect != nil {
				target.ApplyEffect(effect, source.ID)
				result.EffectID = effect.ID
				result.EffectName = effect.Name
			}
		}

	case SkillDispel:
		result.Type = ActionEffectExpire
		// 驱散敌方增益或友方减益
		e.DispelEffect(target, source.IsPlayer != target.IsPlayer)

	case SkillInterrupt:
		result.Type = ActionEffectApply
		// 打断目前由简单的沉默效果实现
		if skill.EffectID != "" {
			effect := e.GetEffect(skill.EffectID)
			if effect != nil && e.rng.Float64() < skill.EffectChance {
				target.ApplyEffect(effect, source.ID)
				result.EffectID = effect.ID
				result.EffectName = effect.Name
			}
		}
	}

	return result
}

// ═══════════════════════════════════════════════════════════
// 伤害计算
// ═══════════════════════════════════════════════════════════

// CalculateDamage 计算伤害
func (e *CombatEngine) CalculateDamage(source, target *BattleUnit, skill *Skill) (damage int, isCrit, isDodged bool) {
	// 检查闪避
	if e.rng.Float64() < target.GetModifiedDodgeRate() {
		return 0, false, true
	}

	// 基础技能伤害 = 技能基础值 + 成长属性 × 成长系数
	scalingStat := e.GetScalingStat(source, skill.ScalingStat)
	baseDamage := float64(skill.BaseValue) + float64(scalingStat)*skill.ScalingRatio

	// 如果是普通攻击，使用攻击力
	if skill.ID == "basic_attack" {
		baseDamage = float64(source.GetModifiedAttack())
	}

	// 应用攻击力加成
	baseDamage *= (1.0 + source.GetStatModifier("attack"))

	// 计算暴击
	critRate := source.GetModifiedCritRate()
	if e.rng.Float64() < critRate {
		isCrit = true
		baseDamage *= source.CombatStats.CritDamage
	}

	// 计算防御减伤
	defense := float64(target.GetModifiedDefense())
	reduction := defense / (defense + 50.0) // 减伤公式，50防御=50%减伤
	baseDamage *= (1.0 - reduction)

	// 应用伤害减免效果
	damageReduction := target.GetStatModifier("damage_taken")
	baseDamage *= (1.0 + damageReduction) // damage_taken是负数表示减伤

	// 添加随机波动 ±10%
	variance := baseDamage * 0.1
	baseDamage += (e.rng.Float64()*2 - 1) * variance

	// 最小伤害为1
	damage = int(baseDamage)
	if damage < 1 {
		damage = 1
	}

	return damage, isCrit, false
}

// CalculateHeal 计算治疗量
func (e *CombatEngine) CalculateHeal(source *BattleUnit, skill *Skill) int {
	scalingStat := e.GetScalingStat(source, skill.ScalingStat)
	baseHeal := float64(skill.BaseValue) + float64(scalingStat)*skill.ScalingRatio

	// 应用治疗加成
	healBonus := source.GetStatModifier("healing_done")
	baseHeal *= (1.0 + healBonus)

	// 添加随机波动 ±10%
	variance := baseHeal * 0.1
	baseHeal += (e.rng.Float64()*2 - 1) * variance

	return int(baseHeal)
}

// GetScalingStat 获取成长属性值
func (e *CombatEngine) GetScalingStat(unit *BattleUnit, stat string) int {
	switch stat {
	case "strength":
		return unit.BaseStats.Strength
	case "agility":
		return unit.BaseStats.Agility
	case "intellect":
		return unit.BaseStats.Intellect
	case "spirit":
		return unit.BaseStats.Spirit
	case "stamina":
		return unit.BaseStats.Stamina
	default:
		return unit.BaseStats.Strength
	}
}

// ═══════════════════════════════════════════════════════════
// 效果处理
// ═══════════════════════════════════════════════════════════

// 效果定义缓存 (实际应从数据库加载)
var effectsCache = map[string]*Effect{
	"eff_stun": {
		ID: "eff_stun", Name: "眩晕", Type: EffectStun,
		IsBuff: false, Duration: 1, CanDispel: true,
	},
	"eff_silence": {
		ID: "eff_silence", Name: "沉默", Type: EffectSilence,
		IsBuff: false, Duration: 2, CanDispel: true,
	},
	"eff_slow": {
		ID: "eff_slow", Name: "减速", Type: EffectStatMod,
		IsBuff: false, Duration: 3, ValueType: "percent", Value: -30, StatAffected: "attack_speed", CanDispel: true,
	},
	"eff_rend": {
		ID: "eff_rend", Name: "撕裂", Type: EffectDOT,
		IsBuff: false, IsStackable: true, MaxStacks: 3, Duration: 3, Value: 2, DamageType: DamagePhysical, CanDispel: true,
	},
	"eff_ignite": {
		ID: "eff_ignite", Name: "点燃", Type: EffectDOT,
		IsBuff: false, IsStackable: true, MaxStacks: 5, Duration: 3, Value: 2, DamageType: DamageFire, CanDispel: true,
	},
	"eff_poison": {
		ID: "eff_poison", Name: "中毒", Type: EffectDOT,
		IsBuff: false, IsStackable: true, MaxStacks: 5, Duration: 5, Value: 2, DamageType: DamageNature, CanDispel: true,
	},
	"eff_renew": {
		ID: "eff_renew", Name: "恢复", Type: EffectHOT,
		IsBuff: true, Duration: 4, Value: 3, CanDispel: true,
	},
	"eff_battle_shout": {
		ID: "eff_battle_shout", Name: "战斗怒吼", Type: EffectStatMod,
		IsBuff: true, Duration: 5, ValueType: "percent", Value: 10, StatAffected: "attack", CanDispel: true,
	},
	"eff_shield_wall": {
		ID: "eff_shield_wall", Name: "盾墙", Type: EffectStatMod,
		IsBuff: true, Duration: 3, ValueType: "percent", Value: -50, StatAffected: "damage_taken", CanDispel: false,
	},
	"eff_pw_shield": {
		ID: "eff_pw_shield", Name: "真言术:盾", Type: EffectShield,
		IsBuff: true, Duration: 4, Value: 12, CanDispel: true,
	},
}

// GetEffect 获取效果定义
func (e *CombatEngine) GetEffect(effectID string) *Effect {
	if effect, ok := effectsCache[effectID]; ok {
		return effect
	}
	return nil
}

// DispelEffect 驱散效果
func (e *CombatEngine) DispelEffect(target *BattleUnit, dispelBuff bool) {
	for i := len(target.ActiveEffects) - 1; i >= 0; i-- {
		effect := target.ActiveEffects[i]
		if effect.Effect.CanDispel {
			if (dispelBuff && effect.Effect.IsBuff) || (!dispelBuff && !effect.Effect.IsBuff) {
				target.ActiveEffects = append(target.ActiveEffects[:i], target.ActiveEffects[i+1:]...)
				break // 只驱散一个
			}
		}
	}
	target.UpdateControlStates()
}

// ═══════════════════════════════════════════════════════════
// 战斗状态检查
// ═══════════════════════════════════════════════════════════

// CheckBattleEnd 检查战斗是否结束
func (e *CombatEngine) CheckBattleEnd() bool {
	playerAlive := false
	enemyAlive := false

	for _, unit := range e.battle.PlayerTeam.Units {
		if unit.IsAlive() {
			playerAlive = true
			break
		}
	}

	for _, unit := range e.battle.EnemyTeam.Units {
		if unit.IsAlive() {
			enemyAlive = true
			break
		}
	}

	if !playerAlive {
		e.battle.State = BattleStateDefeat
		e.battle.EndTime = time.Now()
		return true
	}

	if !enemyAlive {
		e.battle.State = BattleStateVictory
		e.battle.EndTime = time.Now()
		return true
	}

	return false
}

// GetBattle 获取当前战斗
func (e *CombatEngine) GetBattle() *Battle {
	return e.battle
}

// RunFullBattle 运行完整战斗直到结束
func (e *CombatEngine) RunFullBattle() *Battle {
	for e.battle.State == BattleStateInProgress {
		e.ExecuteTurn()
	}
	return e.battle
}


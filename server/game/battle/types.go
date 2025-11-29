package battle

import "time"

// ═══════════════════════════════════════════════════════════
// 战斗系统核心类型定义
// ═══════════════════════════════════════════════════════════

// ResourceType 资源类型
type ResourceType string

const (
	ResourceRage   ResourceType = "rage"   // 怒气 (战士)
	ResourceEnergy ResourceType = "energy" // 能量 (盗贼/猎人)
	ResourceMana   ResourceType = "mana"   // 法力 (法系职业)
)

// DamageType 伤害类型
type DamageType string

const (
	DamagePhysical  DamageType = "physical"
	DamageFire      DamageType = "fire"
	DamageFrost     DamageType = "frost"
	DamageLightning DamageType = "lightning"
	DamageHoly      DamageType = "holy"
	DamageShadow    DamageType = "shadow"
	DamageNature    DamageType = "nature"
	DamageMagic     DamageType = "magic" // 通用魔法
)

// TargetType 目标类型
type TargetType string

const (
	TargetSelf        TargetType = "self"
	TargetAlly        TargetType = "ally"
	TargetAllyAll     TargetType = "ally_all"
	TargetAllyLowest  TargetType = "ally_lowest_hp"
	TargetEnemy       TargetType = "enemy"
	TargetEnemyAll    TargetType = "enemy_all"
	TargetEnemyRandom TargetType = "enemy_random"
	TargetEnemyLowest TargetType = "enemy_lowest_hp"
)

// SkillType 技能类型
type SkillType string

const (
	SkillAttack    SkillType = "attack"
	SkillHeal      SkillType = "heal"
	SkillBuff      SkillType = "buff"
	SkillDebuff    SkillType = "debuff"
	SkillDOT       SkillType = "dot"
	SkillHOT       SkillType = "hot"
	SkillShield    SkillType = "shield"
	SkillControl   SkillType = "control"
	SkillDispel    SkillType = "dispel"
	SkillInterrupt SkillType = "interrupt"
)

// EffectType 效果类型
type EffectType string

const (
	EffectStatMod     EffectType = "stat_mod"     // 属性修改
	EffectDOT         EffectType = "dot"          // 持续伤害
	EffectHOT         EffectType = "hot"          // 持续治疗
	EffectShield      EffectType = "shield"       // 护盾
	EffectStun        EffectType = "stun"         // 眩晕
	EffectSilence     EffectType = "silence"      // 沉默
	EffectSlow        EffectType = "slow"         // 减速
	EffectRoot        EffectType = "root"         // 定身
	EffectTaunt       EffectType = "taunt"        // 嘲讽
	EffectInvulnerable EffectType = "invulnerable" // 无敌
	EffectLifesteal   EffectType = "lifesteal"    // 吸血
	EffectProc        EffectType = "proc"         // 触发效果
)

// ═══════════════════════════════════════════════════════════
// 基础属性结构
// ═══════════════════════════════════════════════════════════

// BaseStats 基础属性
type BaseStats struct {
	Strength  int `json:"strength"`  // 力量
	Agility   int `json:"agility"`   // 敏捷
	Intellect int `json:"intellect"` // 智力
	Stamina   int `json:"stamina"`   // 耐力
	Spirit    int `json:"spirit"`    // 精神
}

// CombatStats 战斗属性 (从基础属性计算得出)
type CombatStats struct {
	MaxHP       int     `json:"max_hp"`
	MaxResource int     `json:"max_resource"`
	Attack      int     `json:"attack"`
	Defense     int     `json:"defense"`
	CritRate    float64 `json:"crit_rate"`    // 0-1
	CritDamage  float64 `json:"crit_damage"`  // 默认1.5
	DodgeRate   float64 `json:"dodge_rate"`   // 0-1
	HitRate     float64 `json:"hit_rate"`     // 默认1.0
	HastePct    float64 `json:"haste_pct"`    // 加速百分比
}

// ═══════════════════════════════════════════════════════════
// 技能定义
// ═══════════════════════════════════════════════════════════

// Skill 技能定义
type Skill struct {
	ID           string     `json:"id"`
	Name         string     `json:"name"`
	Description  string     `json:"description"`
	ClassID      string     `json:"class_id"`      // 空为通用技能
	Type         SkillType  `json:"type"`
	TargetType   TargetType `json:"target_type"`
	DamageType   DamageType `json:"damage_type"`
	BaseValue    int        `json:"base_value"`     // 基础数值
	ScalingStat  string     `json:"scaling_stat"`   // 成长属性
	ScalingRatio float64    `json:"scaling_ratio"`  // 成长系数
	ResourceCost int        `json:"resource_cost"`
	Cooldown     int        `json:"cooldown"`       // 冷却回合数
	EffectID     string     `json:"effect_id"`      // 附加效果
	EffectChance float64    `json:"effect_chance"`  // 效果触发概率
}

// SkillState 技能状态 (运行时)
type SkillState struct {
	Skill           *Skill `json:"skill"`
	CurrentCooldown int    `json:"current_cooldown"` // 当前冷却剩余
}

// ═══════════════════════════════════════════════════════════
// 效果定义
// ═══════════════════════════════════════════════════════════

// Effect 效果定义 (Buff/Debuff模板)
type Effect struct {
	ID           string     `json:"id"`
	Name         string     `json:"name"`
	Description  string     `json:"description"`
	Type         EffectType `json:"type"`
	IsBuff       bool       `json:"is_buff"`
	IsStackable  bool       `json:"is_stackable"`
	MaxStacks    int        `json:"max_stacks"`
	Duration     int        `json:"duration"`      // 持续回合
	ValueType    string     `json:"value_type"`    // flat/percent
	Value        float64    `json:"value"`
	StatAffected string     `json:"stat_affected"` // 影响的属性
	DamageType   DamageType `json:"damage_type"`   // DOT伤害类型
	CanDispel    bool       `json:"can_dispel"`
}

// ActiveEffect 活跃效果 (运行时)
type ActiveEffect struct {
	Effect         *Effect `json:"effect"`
	SourceID       string  `json:"source_id"`       // 施加者ID
	RemainingTurns int     `json:"remaining_turns"` // 剩余回合
	Stacks         int     `json:"stacks"`          // 当前层数
	ShieldAmount   int     `json:"shield_amount"`   // 护盾剩余量
}

// ═══════════════════════════════════════════════════════════
// 战斗单位
// ═══════════════════════════════════════════════════════════

// BattleUnit 战斗单位 (角色或怪物)
type BattleUnit struct {
	ID           string       `json:"id"`
	Name         string       `json:"name"`
	Level        int          `json:"level"`
	IsPlayer     bool         `json:"is_player"`     // true=玩家角色, false=怪物
	ClassID      string       `json:"class_id"`
	ResourceType ResourceType `json:"resource_type"`

	// 属性
	BaseStats   BaseStats   `json:"base_stats"`
	CombatStats CombatStats `json:"combat_stats"`

	// 当前状态
	CurrentHP       int `json:"current_hp"`
	CurrentResource int `json:"current_resource"`

	// 技能
	Skills      []*Skill       `json:"skills"`
	SkillStates []*SkillState  `json:"skill_states"`

	// 效果
	ActiveEffects []*ActiveEffect `json:"active_effects"`

	// 战斗状态
	IsDead       bool `json:"is_dead"`
	IsStunned    bool `json:"is_stunned"`
	IsSilenced   bool `json:"is_silenced"`
	TauntedBy    string `json:"taunted_by"` // 被嘲讽时指向施法者ID

	// 战斗统计 (本场战斗)
	DamageDealt   int `json:"damage_dealt"`
	DamageTaken   int `json:"damage_taken"`
	HealingDone   int `json:"healing_done"`
	HealingTaken  int `json:"healing_taken"`
}

// ═══════════════════════════════════════════════════════════
// 战斗相关
// ═══════════════════════════════════════════════════════════

// Team 队伍
type Team struct {
	Units    []*BattleUnit `json:"units"`
	IsPlayer bool          `json:"is_player"` // 是否玩家队伍
	Faction  string        `json:"faction"`   // alliance/horde
}

// BattleState 战斗状态
type BattleState string

const (
	BattleStateNotStarted BattleState = "not_started"
	BattleStateInProgress BattleState = "in_progress"
	BattleStateVictory    BattleState = "victory"
	BattleStateDefeat     BattleState = "defeat"
	BattleStateDraw       BattleState = "draw"
)

// Battle 战斗实例
type Battle struct {
	ID          string      `json:"id"`
	State       BattleState `json:"state"`
	CurrentTurn int         `json:"current_turn"`
	MaxTurns    int         `json:"max_turns"` // 最大回合数限制

	PlayerTeam *Team `json:"player_team"`
	EnemyTeam  *Team `json:"enemy_team"`

	TurnOrder []*BattleUnit `json:"turn_order"` // 行动顺序

	ActionLog []ActionLog `json:"action_log"` // 战斗日志

	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
}

// ═══════════════════════════════════════════════════════════
// 战斗行动
// ═══════════════════════════════════════════════════════════

// ActionType 行动类型
type ActionType string

const (
	ActionSkill       ActionType = "skill"
	ActionBasicAttack ActionType = "basic_attack"
	ActionSkipped     ActionType = "skipped" // 被控制跳过
	ActionDOTTick     ActionType = "dot_tick"
	ActionHOTTick     ActionType = "hot_tick"
	ActionEffectApply ActionType = "effect_apply"
	ActionEffectExpire ActionType = "effect_expire"
	ActionDeath       ActionType = "death"
	ActionRevive      ActionType = "revive"
)

// ActionResult 行动结果
type ActionResult struct {
	Type        ActionType `json:"type"`
	SourceID    string     `json:"source_id"`
	SourceName  string     `json:"source_name"`
	TargetID    string     `json:"target_id"`
	TargetName  string     `json:"target_name"`
	SkillID     string     `json:"skill_id"`
	SkillName   string     `json:"skill_name"`
	Value       int        `json:"value"`       // 伤害/治疗量
	IsCrit      bool       `json:"is_crit"`
	IsDodged    bool       `json:"is_dodged"`
	EffectID    string     `json:"effect_id"`   // 附加效果
	EffectName  string     `json:"effect_name"`
}

// ActionLog 行动日志
type ActionLog struct {
	Turn      int            `json:"turn"`
	Timestamp time.Time      `json:"timestamp"`
	Results   []ActionResult `json:"results"`
}

// ═══════════════════════════════════════════════════════════
// 战斗奖励
// ═══════════════════════════════════════════════════════════

// BattleReward 战斗奖励
type BattleReward struct {
	Exp         int           `json:"exp"`
	Gold        int           `json:"gold"`
	Items       []DroppedItem `json:"items"`
	IsMiracle   bool          `json:"is_miracle"` // 是否奇迹掉落
}

// DroppedItem 掉落物品
type DroppedItem struct {
	ItemID   string `json:"item_id"`
	Name     string `json:"name"`
	Quality  string `json:"quality"`
	Quantity int    `json:"quantity"`
}


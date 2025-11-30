package battle

import (
	"fmt"
)

// ═══════════════════════════════════════════════════════════
// 战斗单位工厂 - 创建战斗实体
// ═══════════════════════════════════════════════════════════

// CharacterData 角色数据 (从游戏系统传入)
type CharacterData struct {
	ID        string
	Name      string
	Level     int
	ClassID   string
	RaceID    string
	
	// 属性
	Strength  int
	Agility   int
	Intellect int
	Stamina   int
	Spirit    int
	
	// 当前状态
	CurrentHP       int
	CurrentResource int
}

// MonsterData 怪物数据
type MonsterData struct {
	ID      string
	Name    string
	Level   int
	MaxHP   int
	Attack  int
	Defense int
	ExpDrop int
	GoldMin int
	GoldMax int
}

// CreatePlayerUnit 从角色数据创建战斗单位
func CreatePlayerUnit(data *CharacterData) *BattleUnit {
	baseStats := BaseStats{
		Strength:  data.Strength,
		Agility:   data.Agility,
		Intellect: data.Intellect,
		Stamina:   data.Stamina,
		Spirit:    data.Spirit,
	}

	unit := NewBattleUnit(
		data.ID,
		data.Name,
		data.Level,
		true,
		data.ClassID,
		baseStats,
	)

	// 如果有当前HP/资源状态，恢复它
	if data.CurrentHP > 0 {
		unit.CurrentHP = data.CurrentHP
	}
	if data.CurrentResource >= 0 {
		unit.CurrentResource = data.CurrentResource
	}

	// 添加职业技能
	skills := GetSkillsForClass(data.ClassID)
	for _, skill := range skills {
		unit.AddSkill(skill)
	}

	// 添加普通攻击
	unit.AddSkill(GetBasicAttack())

	return unit
}

// CreateMonsterUnit 从怪物数据创建战斗单位
func CreateMonsterUnit(data *MonsterData) *BattleUnit {
	// 根据怪物属性反推基础属性
	baseStats := BaseStats{
		Strength:  data.Attack / 2,
		Agility:   data.Level * 2,
		Intellect: data.Level,
		Stamina:   data.MaxHP / 10,
		Spirit:    data.Level,
	}

	unit := NewBattleUnit(
		data.ID,
		data.Name,
		data.Level,
		false,
		"monster",
		baseStats,
	)

	// 覆盖计算的属性，使用怪物定义的数值
	unit.CombatStats.MaxHP = data.MaxHP
	unit.CombatStats.Attack = data.Attack
	unit.CombatStats.Defense = data.Defense
	unit.CurrentHP = data.MaxHP

	// 怪物使用简单技能
	monsterSkills := CreateMonsterSkills(data.Level)
	for _, skill := range monsterSkills {
		unit.AddSkill(skill)
	}

	return unit
}

// CreateTeam 创建队伍
func CreateTeam(units []*BattleUnit, isPlayer bool, faction string) *Team {
	return &Team{
		Units:    units,
		IsPlayer: isPlayer,
		Faction:  faction,
	}
}

// ═══════════════════════════════════════════════════════════
// 快速战斗接口
// ═══════════════════════════════════════════════════════════

// QuickBattleResult 快速战斗结果
type QuickBattleResult struct {
	Victory      bool            `json:"victory"`
	TotalTurns   int             `json:"total_turns"`
	PlayerStats  []UnitBattleStats `json:"player_stats"`
	EnemyStats   []UnitBattleStats `json:"enemy_stats"`
	ExpGained    int             `json:"exp_gained"`
	GoldGained   int             `json:"gold_gained"`
	Logs         []ActionLog     `json:"logs"`
}

// UnitBattleStats 单位战斗统计
type UnitBattleStats struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	FinalHP      int    `json:"final_hp"`
	MaxHP        int    `json:"max_hp"`
	DamageDealt  int    `json:"damage_dealt"`
	DamageTaken  int    `json:"damage_taken"`
	HealingDone  int    `json:"healing_done"`
	IsDead       bool   `json:"is_dead"`
}

// RunQuickBattle 运行快速战斗
func RunQuickBattle(playerUnits []*BattleUnit, enemyUnits []*BattleUnit) *QuickBattleResult {
	// 创建队伍
	playerTeam := CreateTeam(playerUnits, true, "alliance")
	enemyTeam := CreateTeam(enemyUnits, false, "horde")

	// 创建战斗引擎
	engine := NewCombatEngine()
	engine.StartBattle(playerTeam, enemyTeam)

	// 运行战斗
	battle := engine.RunFullBattle()

	// 收集结果
	result := &QuickBattleResult{
		Victory:    battle.State == BattleStateVictory,
		TotalTurns: battle.CurrentTurn,
		Logs:       battle.ActionLog,
	}

	// 玩家统计
	for _, unit := range playerTeam.Units {
		result.PlayerStats = append(result.PlayerStats, UnitBattleStats{
			ID:           unit.ID,
			Name:         unit.Name,
			FinalHP:      unit.CurrentHP,
			MaxHP:        unit.CombatStats.MaxHP,
			DamageDealt:  unit.DamageDealt,
			DamageTaken:  unit.DamageTaken,
			HealingDone:  unit.HealingDone,
			IsDead:       unit.IsDead,
		})
	}

	// 敌人统计
	for _, unit := range enemyTeam.Units {
		result.EnemyStats = append(result.EnemyStats, UnitBattleStats{
			ID:           unit.ID,
			Name:         unit.Name,
			FinalHP:      unit.CurrentHP,
			MaxHP:        unit.CombatStats.MaxHP,
			DamageDealt:  unit.DamageDealt,
			DamageTaken:  unit.DamageTaken,
			HealingDone:  unit.HealingDone,
			IsDead:       unit.IsDead,
		})
	}

	return result
}

// ═══════════════════════════════════════════════════════════
// 战斗日志格式化
// ═══════════════════════════════════════════════════════════

// FormatActionResult 格式化行动结果为文本
func FormatActionResult(result ActionResult) string {
	switch result.Type {
	case ActionSkill:
		if result.IsDodged {
			return fmt.Sprintf("%s 使用 [%s]，但被 %s 闪避了！",
				result.SourceName, result.SkillName, result.TargetName)
		}
		if result.IsCrit {
			return fmt.Sprintf("%s 使用 [%s] 对 %s 造成 %d 点暴击伤害！",
				result.SourceName, result.SkillName, result.TargetName, result.Value)
		}
		if result.Value > 0 {
			return fmt.Sprintf("%s 使用 [%s] 对 %s 造成 %d 点伤害",
				result.SourceName, result.SkillName, result.TargetName, result.Value)
		}
		return fmt.Sprintf("%s 使用 [%s]", result.SourceName, result.SkillName)

	case ActionDOTTick:
		return fmt.Sprintf("%s 受到 [%s] 效果 %d 点伤害",
			result.TargetName, result.EffectName, result.Value)

	case ActionHOTTick:
		return fmt.Sprintf("%s 受到 [%s] 效果恢复 %d 点生命",
			result.TargetName, result.EffectName, result.Value)

	case ActionEffectApply:
		return fmt.Sprintf("%s 对 %s 施加了 [%s] 效果",
			result.SourceName, result.TargetName, result.EffectName)

	case ActionEffectExpire:
		return fmt.Sprintf("%s 的 [%s] 效果消失了",
			result.TargetName, result.EffectName)

	case ActionDeath:
		return fmt.Sprintf("%s 被击败了！", result.TargetName)

	case ActionSkipped:
		return fmt.Sprintf("%s 被控制无法行动", result.SourceName)

	default:
		return ""
	}
}

// GetActionColor 获取行动日志颜色
func GetActionColor(result ActionResult) string {
	switch result.Type {
	case ActionSkill:
		if result.IsDodged {
			return "#888888"
		}
		if result.IsCrit {
			return "#FF6B6B"
		}
		if result.Value > 0 {
			return "#FFAA00"
		}
		return "#00FF00"
	case ActionDOTTick:
		return "#FF4444"
	case ActionHOTTick:
		return "#00FF00"
	case ActionEffectApply:
		return "#00FFFF"
	case ActionEffectExpire:
		return "#888888"
	case ActionDeath:
		return "#FF0000"
	case ActionSkipped:
		return "#FFFF00"
	default:
		return "#FFFFFF"
	}
}



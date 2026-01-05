package game

import (
	"math"
	"math/rand"
	"strings"

	"text-wow/internal/models"
)

// Calculator 数值计算器 - 统一管理所有数值计算逻辑
type Calculator struct {
	rng *rand.Rand
}

// NewCalculator 创建数值计算器
func NewCalculator() *Calculator {
	return &Calculator{
		rng: rand.New(rand.NewSource(rand.Int63())),
	}
}

// ═══════════════════════════════════════════════════════════
// 属性转换计算
// ═══════════════════════════════════════════════════════════

// CalculatePhysicalAttack 计算物理攻击力
// 公式: (力量 × 0.4) + (敏捷 × 0.2) + 武器伤害 + 装备加成
// 边界处理: 确保返回值至少为1
func (c *Calculator) CalculatePhysicalAttack(char *models.Character) int {
	if char == nil {
		return 1
	}
	baseAttack := float64(char.Strength)*0.4 + float64(char.Agility)*0.2
	result := int(math.Round(baseAttack))
	if result < 1 {
		result = 1
	}
	return result
}

// CalculateMagicAttack 计算法术攻击力
// 公式: (智力 × 1.0) + (精神 × 0.2) + 装备加成
// 边界处理: 确保返回值至少为1
func (c *Calculator) CalculateMagicAttack(char *models.Character) int {
	if char == nil {
		return 1
	}
	baseAttack := float64(char.Intellect)*1.0 + float64(char.Spirit)*0.2
	result := int(math.Round(baseAttack))
	if result < 1 {
		result = 1
	}
	return result
}

// CalculateHP 计算最大生命值
// 公式: 职业基础HP + (耐力 × 2) + 装备加成
// 边界处理: 确保返回值至少为1
func (c *Calculator) CalculateHP(char *models.Character, baseHP int) int {
	if char == nil {
		return 1
	}
	if baseHP < 0 {
		baseHP = 0
	}
	result := baseHP + char.Stamina*2
	if result < 1 {
		result = 1
	}
	return result
}

// CalculateMP 计算最大法力值
// 公式: 职业基础MP + (精神 × 2) + 装备加成
// 边界处理: 确保返回值至少为0
func (c *Calculator) CalculateMP(char *models.Character, baseMP int) int {
	if char == nil {
		return 0
	}
	if baseMP < 0 {
		baseMP = 0
	}
	result := baseMP + char.Spirit*2
	if result < 0 {
		result = 0
	}
	return result
}

// CalculatePhysCritRate 计算物理暴击率
// 公式: 5% (基础) + (敏捷 / 20) + 装备加成
// 上限: 50%
// 边界处理: 确保返回值在0-0.5之间
func (c *Calculator) CalculatePhysCritRate(char *models.Character) float64 {
	if char == nil {
		return 0.05
	}
	baseRate := 0.05
	agilityBonus := float64(char.Agility) / 20.0 / 100.0 // 转换为小数
	rate := baseRate + agilityBonus
	// 限制范围在0-50%之间
	if rate < 0 {
		rate = 0
	}
	if rate > 0.5 {
		rate = 0.5
	}
	return rate
}

// CalculatePhysCritDamage 计算物理暴击伤害倍率
// 公式: 150% (基础) + (力量 × 0.3%) + 装备加成
// 边界处理: 确保返回值至少为1.0
func (c *Calculator) CalculatePhysCritDamage(char *models.Character) float64 {
	if char == nil {
		return 1.5
	}
	baseRate := 1.5
	strengthBonus := float64(char.Strength) * 0.003 // 0.3% = 0.003
	result := baseRate + strengthBonus
	if result < 1.0 {
		result = 1.0
	}
	return result
}

// CalculateSpellCritRate 计算法术暴击率
// 公式: 5% (基础) + (精神 / 20) + 装备加成
func (c *Calculator) CalculateSpellCritRate(char *models.Character) float64 {
	baseRate := 0.05
	spiritBonus := float64(char.Spirit) / 20.0 / 100.0 // 转换为小数
	return baseRate + spiritBonus
}

// CalculateSpellCritDamage 计算法术暴击伤害倍率
// 公式: 150% (基础) + (智力 × 0.3%) + 装备加成
// 边界处理: 确保返回值至少为1.0
func (c *Calculator) CalculateSpellCritDamage(char *models.Character) float64 {
	if char == nil {
		return 1.5
	}
	baseRate := 1.5
	intellectBonus := float64(char.Intellect) * 0.003 // 0.3% = 0.003
	result := baseRate + intellectBonus
	if result < 1.0 {
		result = 1.0
	}
	return result
}

// CalculateDodgeRate 计算闪避率
// 公式: 5% (基础) + (敏捷 / 20) + 装备加成
// 上限: 50%
// 边界处理: 确保返回值在0-0.5之间
func (c *Calculator) CalculateDodgeRate(char *models.Character) float64 {
	if char == nil {
		return 0.05
	}
	baseRate := 0.05
	agilityBonus := float64(char.Agility) / 20.0 / 100.0 // 转换为小数
	rate := baseRate + agilityBonus
	// 限制范围在0-50%之间
	if rate < 0 {
		rate = 0
	}
	if rate > 0.5 {
		rate = 0.5
	}
	return rate
}

// ═══════════════════════════════════════════════════════════
// 伤害计算
// ═══════════════════════════════════════════════════════════

// DamageCalculationResult 伤害计算结果
type DamageCalculationResult struct {
	BaseDamage    float64 // 基础伤害
	DefenseReduction float64 // 防御减伤
	DamageAfterDefense float64 // 减伤后伤害
	IsCrit        bool    // 是否暴击
	CritMultiplier float64 // 暴击倍率
	FinalDamage   int     // 最终伤害
	IsDodged      bool    // 是否闪避
}

// CalculateDamage 计算伤害
// 完整流程: 基础伤害 → 防御减伤 → 暴击判定 → 闪避判定 → 最终伤害
// 支持伤害类型: physical, magic, fire, frost, shadow, holy, nature
// 边界处理: 处理nil指针、负数输入、无效参数
func (c *Calculator) CalculateDamage(
	attacker *models.Character,
	defender *models.Character,
	baseAttack int,
	skillMultiplier float64,
	damageType string, // physical/magic/fire/frost/shadow/holy/nature
	ignoreDodge bool, // 是否无视闪避
) *DamageCalculationResult {
	result := &DamageCalculationResult{}

	// 边界检查
	if attacker == nil || defender == nil {
		result.FinalDamage = 1
		return result
	}
	if baseAttack < 0 {
		baseAttack = 0
	}
	if skillMultiplier < 0 {
		skillMultiplier = 0
	}

	// 规范化伤害类型（转换为小写并验证）
	damageType = strings.ToLower(damageType)
	validTypes := map[string]bool{
		"physical": true,
		"magic":    true,
		"fire":     true,
		"frost":    true,
		"shadow":   true,
		"holy":     true,
		"nature":   true,
	}
	if !validTypes[damageType] {
		damageType = "physical" // 默认物理伤害
	}

	// 1. 基础伤害计算
	result.BaseDamage = float64(baseAttack) * skillMultiplier
	if result.BaseDamage < 0 {
		result.BaseDamage = 0
	}

	// 2. 防御减伤计算
	var defense int
	// 物理伤害使用物理防御
	if damageType == "physical" {
		defense = defender.PhysicalDefense
	} else {
		// 所有法术伤害（magic/fire/frost/shadow/holy/nature）使用魔法防御
		defense = defender.MagicDefense
	}
	if defense < 0 {
		defense = 0
	}

	// 防御减伤公式: 伤害 = 攻击 - 防御（减法公式）
	// 最低伤害: 1点
	damageAfterDefense := result.BaseDamage - float64(defense)
	if damageAfterDefense < 1 {
		damageAfterDefense = 1
	}
	
	// DefenseReduction 字段保留用于兼容性，但不再使用百分比
	// 计算实际的减伤率（用于显示）
	if result.BaseDamage > 0 {
		result.DefenseReduction = float64(defense) / result.BaseDamage
	} else {
		result.DefenseReduction = 0
	}
	
	result.DamageAfterDefense = damageAfterDefense

	// 3. 暴击判定
	var critRate float64
	var critDamage float64
	if damageType == "physical" {
		critRate = attacker.PhysCritRate
		critDamage = attacker.PhysCritDamage
	} else {
		// 所有法术伤害使用法术暴击
		critRate = attacker.SpellCritRate
		critDamage = attacker.SpellCritDamage
	}

	// 限制暴击率上限
	if critRate > 0.5 {
		critRate = 0.5
	}

	roll := c.rng.Float64()
	if roll < critRate {
		result.IsCrit = true
		result.CritMultiplier = critDamage
		result.DamageAfterDefense = result.DamageAfterDefense * critDamage
	}

	// 4. 闪避判定
	// 只有物理伤害可以闪避
	if !ignoreDodge && damageType == "physical" {
		dodgeRate := defender.DodgeRate
		// 限制闪避率上限
		if dodgeRate > 0.5 {
			dodgeRate = 0.5
		}

		roll := c.rng.Float64()
		if roll < dodgeRate {
			result.IsDodged = true
			result.FinalDamage = 0
			return result
		}
	}

	// 5. 最终伤害（至少1点）
	result.FinalDamage = int(math.Round(result.DamageAfterDefense))
	if result.FinalDamage < 1 {
		result.FinalDamage = 1
	}

	return result
}


// CalculateDefenseReduction 计算防御减伤（已废弃，保留用于兼容性）
// 注意：防御现在使用减法公式（伤害 = 攻击 - 防御），不再使用百分比减伤
// 此函数保留用于向后兼容，但实际计算已改为直接减法
// 元素抗性才使用百分比减伤（未来实现）
// 返回值：0（表示不使用百分比减伤）
func (c *Calculator) CalculateDefenseReduction(defense int) float64 {
	// 此函数已废弃，防御现在使用减法公式
	// 返回0表示不使用百分比减伤
	return 0
}

// ShouldCrit 判断是否暴击
// 边界处理: 确保critRate在0-0.5之间
func (c *Calculator) ShouldCrit(critRate float64) bool {
	if critRate < 0 {
		critRate = 0
	}
	if critRate > 0.5 {
		critRate = 0.5
	}
	return c.rng.Float64() < critRate
}

// ShouldDodge 判断是否闪避
// 边界处理: 确保dodgeRate在0-0.5之间
func (c *Calculator) ShouldDodge(dodgeRate float64) bool {
	if dodgeRate < 0 {
		dodgeRate = 0
	}
	if dodgeRate > 0.5 {
		dodgeRate = 0.5
	}
	return c.rng.Float64() < dodgeRate
}

// ═══════════════════════════════════════════════════════════
// 治疗计算
// ═══════════════════════════════════════════════════════════

// HealingCalculationResult 治疗计算结果
type HealingCalculationResult struct {
	BaseHealing    float64 // 基础治疗
	HealingBonus   float64 // 治疗加成
	FinalHealing   int     // 最终治疗
	Overhealing    int     // 过量治疗
	ActualHealing  int     // 实际治疗（考虑最大HP）
}

// CalculateHealing 计算治疗
// 公式: 基础治疗 = 治疗力 × 技能系数
//       最终治疗 = 基础治疗 × (1 + 治疗加成%)
// 边界处理: 处理nil指针、负数输入
func (c *Calculator) CalculateHealing(
	healer *models.Character,
	target *models.Character,
	baseHealing int,
	skillMultiplier float64,
	healingBonusPercent float64, // 治疗加成百分比
) *HealingCalculationResult {
	result := &HealingCalculationResult{}

	// 边界检查
	if healer == nil || target == nil {
		result.FinalHealing = 0
		result.ActualHealing = 0
		return result
	}
	if baseHealing < 0 {
		baseHealing = 0
	}
	if skillMultiplier < 0 {
		skillMultiplier = 0
	}
	if healingBonusPercent < 0 {
		healingBonusPercent = 0
	}

	// 1. 基础治疗
	result.BaseHealing = float64(baseHealing) * skillMultiplier
	if result.BaseHealing < 0 {
		result.BaseHealing = 0
	}

	// 2. 治疗加成
	result.HealingBonus = healingBonusPercent / 100.0
	finalHealing := result.BaseHealing * (1.0 + result.HealingBonus)
	result.FinalHealing = int(math.Round(finalHealing))

	// 3. 计算实际治疗（考虑最大HP）
	if target.MaxHP < target.HP {
		target.MaxHP = target.HP // 修复数据不一致
	}
	maxHealable := target.MaxHP - target.HP
	if maxHealable < 0 {
		maxHealable = 0
	}
	if result.FinalHealing > maxHealable {
		result.Overhealing = result.FinalHealing - maxHealable
		result.ActualHealing = maxHealable
	} else {
		result.Overhealing = 0
		result.ActualHealing = result.FinalHealing
	}
	if result.ActualHealing < 0 {
		result.ActualHealing = 0
	}

	return result
}

// ═══════════════════════════════════════════════════════════
// 资源计算
// ═══════════════════════════════════════════════════════════

// CalculateManaRegen 计算法力恢复
// 公式: 基础恢复 + (精神 × 0.1) + 装备加成
// 边界处理: 确保返回值至少为0
func (c *Calculator) CalculateManaRegen(char *models.Character, baseRegen int) int {
	if char == nil {
		return 0
	}
	if baseRegen < 0 {
		baseRegen = 0
	}
	spiritBonus := float64(char.Spirit) * 0.1
	result := baseRegen + int(math.Round(spiritBonus))
	if result < 0 {
		result = 0
	}
	return result
}

// CalculateRageGain 计算怒气获取
// 攻击获得: 基础获得 + 技能加成
// 受击获得: 基础获得 + 技能加成
// 边界处理: 确保返回值至少为0
func (c *Calculator) CalculateRageGain(baseGain int, bonusPercent float64) int {
	if baseGain < 0 {
		baseGain = 0
	}
	if bonusPercent < 0 {
		bonusPercent = 0
	}
	result := int(math.Round(float64(baseGain) * (1.0 + bonusPercent/100.0)))
	if result < 0 {
		result = 0
	}
	return result
}

// CalculateEnergyRegen 计算能量恢复
// 公式: 基础恢复 + 装备加成
// 边界处理: 确保返回值至少为0
func (c *Calculator) CalculateEnergyRegen(char *models.Character, baseRegen int) int {
	if baseRegen < 0 {
		baseRegen = 0
	}
	return baseRegen // 能量恢复主要依赖基础值，精神影响较小
}

// CalculateSpeed 计算角色速度
// 公式: 速度 = 敏捷
// 边界处理: 确保返回值至少为1
func (c *Calculator) CalculateSpeed(char *models.Character) int {
	if char == nil {
		return 10 // 默认速度
	}
	result := char.Agility
	if result < 1 {
		result = 1
	}
	return result
}

// ═══════════════════════════════════════════════════════════
// 辅助计算
// ═══════════════════════════════════════════════════════════

// Clamp 限制数值范围
func (c *Calculator) Clamp(value, min, max float64) float64 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

// ClampInt 限制整数范围
func (c *Calculator) ClampInt(value, min, max int) int {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}



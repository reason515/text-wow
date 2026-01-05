package game

import (
	"testing"

	"text-wow/internal/models"

	"github.com/stretchr/testify/assert"
)

// ═══════════════════════════════════════════════════════════
// 属性转换计算测试
// ═══════════════════════════════════════════════════════════
// 注意：属性转换的详细测试用例已迁移到 YAML 格式
// 文件位置：server/internal/test/cases/calculator/attribute_conversion.yaml
// 根据 PRINCIPLES.md，测试用例应使用 YAML 自然语言描述
// 以下Go测试用例保留用于快速验证和调试（特殊情况）

func TestCalculatePhysicalAttack(t *testing.T) {
	calc := NewCalculator()

	tests := []struct {
		name     string
		char     *models.Character
		expected int
	}{
		{
			name: "正常计算",
			char: &models.Character{
				Strength: 20,
				Agility:  10,
			},
			expected: 10, // 20*0.4 + 10*0.2 = 8 + 2 = 10
		},
		{
			name: "零属性",
			char: &models.Character{
				Strength: 0,
				Agility:  0,
			},
			expected: 1, // 最小值
		},
		{
			name: "高属性",
			char: &models.Character{
				Strength: 100,
				Agility:  50,
			},
			expected: 50, // 100*0.4 + 50*0.2 = 40 + 10 = 50
		},
		{
			name: "nil角色",
			char: nil,
			expected: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calc.CalculatePhysicalAttack(tt.char)
			assert.Equal(t, tt.expected, result, "物理攻击力计算不正确")
			assert.GreaterOrEqual(t, result, 1, "物理攻击力应该至少为1")
		})
	}
}

func TestCalculateMagicAttack(t *testing.T) {
	calc := NewCalculator()

	tests := []struct {
		name     string
		char     *models.Character
		expected int
	}{
		{
			name: "正常计算",
			char: &models.Character{
				Intellect: 15,
				Spirit:    10,
			},
			expected: 17, // 15*1.0 + 10*0.2 = 15 + 2 = 17
		},
		{
			name: "零属性",
			char: &models.Character{
				Intellect: 0,
				Spirit:    0,
			},
			expected: 1,
		},
		{
			name: "高属性",
			char: &models.Character{
				Intellect: 50,
				Spirit:    30,
			},
			expected: 56, // 50*1.0 + 30*0.2 = 50 + 6 = 56
		},
		{
			name: "nil角色",
			char: nil,
			expected: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calc.CalculateMagicAttack(tt.char)
			assert.Equal(t, tt.expected, result, "法术攻击力计算不正确")
			assert.GreaterOrEqual(t, result, 1, "法术攻击力应该至少为1")
		})
	}
}

func TestCalculateHP(t *testing.T) {
	calc := NewCalculator()

	tests := []struct {
		name     string
		char     *models.Character
		baseHP   int
		expected int
	}{
		{
			name: "正常计算",
			char: &models.Character{
				Stamina: 15,
			},
			baseHP:   35,
			expected: 65, // 35 + 15*2 = 65
		},
		{
			name: "零耐力",
			char: &models.Character{
				Stamina: 0,
			},
			baseHP:   10,
			expected: 10,
		},
		{
			name: "负数基础HP",
			char: &models.Character{
				Stamina: 10,
			},
			baseHP:   -5,
			expected: 20, // 负数被处理为0，0 + 10*2 = 20
		},
		{
			name: "nil角色",
			char: nil,
			baseHP:   10,
			expected: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calc.CalculateHP(tt.char, tt.baseHP)
			assert.Equal(t, tt.expected, result, "最大生命值计算不正确")
			assert.GreaterOrEqual(t, result, 1, "最大生命值应该至少为1")
		})
	}
}

func TestCalculatePhysCritRate(t *testing.T) {
	calc := NewCalculator()

	tests := []struct {
		name     string
		char     *models.Character
		expected float64
	}{
		{
			name: "正常计算",
			char: &models.Character{
				Agility: 20,
			},
			expected: 0.06, // 0.05 + 20/20/100 = 0.05 + 0.01 = 0.06
		},
		{
			name: "零敏捷",
			char: &models.Character{
				Agility: 0,
			},
			expected: 0.05, // 基础5%
		},
		{
			name: "高敏捷（超过上限）",
			char: &models.Character{
				Agility: 1000,
			},
			expected: 0.5, // 上限50%
		},
		{
			name: "nil角色",
			char: nil,
			expected: 0.05,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calc.CalculatePhysCritRate(tt.char)
			assert.InDelta(t, tt.expected, result, 0.001, "物理暴击率计算不正确")
			assert.GreaterOrEqual(t, result, 0.0, "物理暴击率应该至少为0")
			assert.LessOrEqual(t, result, 0.5, "物理暴击率应该不超过50%")
		})
	}
}

func TestCalculateDodgeRate(t *testing.T) {
	calc := NewCalculator()

	tests := []struct {
		name     string
		char     *models.Character
		expected float64
	}{
		{
			name: "正常计算",
			char: &models.Character{
				Agility: 30,
			},
			expected: 0.065, // 0.05 + 30/20/100 = 0.05 + 0.015 = 0.065
		},
		{
			name: "零敏捷",
			char: &models.Character{
				Agility: 0,
			},
			expected: 0.05,
		},
		{
			name: "高敏捷（超过上限）",
			char: &models.Character{
				Agility: 1000,
			},
			expected: 0.5, // 上限50%
		},
		{
			name: "nil角色",
			char: nil,
			expected: 0.05,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calc.CalculateDodgeRate(tt.char)
			assert.InDelta(t, tt.expected, result, 0.001, "闪避率计算不正确")
			assert.GreaterOrEqual(t, result, 0.0, "闪避率应该至少为0")
			assert.LessOrEqual(t, result, 0.5, "闪避率应该不超过50%")
		})
	}
}

// ═══════════════════════════════════════════════════════════
// 伤害计算测试
// ═══════════════════════════════════════════════════════════
// 注意：伤害计算的详细测试用例已迁移到 YAML 格式
// 文件位置：server/internal/test/cases/calculator/damage_calculation.yaml
// 根据 PRINCIPLES.md，测试用例应使用 YAML 自然语言描述
// 以下Go测试用例保留用于快速验证和调试（特殊情况）

func TestCalculateDamage(t *testing.T) {
	calc := NewCalculator()

	attacker := &models.Character{
		PhysicalAttack:   20,
		MagicAttack:      15,
		PhysCritRate:     0.1,
		PhysCritDamage:   1.5,
		SpellCritRate:    0.1,
		SpellCritDamage:  1.5,
	}

	defender := &models.Character{
		PhysicalDefense: 10,
		MagicDefense:    5,
		DodgeRate:       0.1,
	}

	tests := []struct {
		name          string
		attacker      *models.Character
		defender      *models.Character
		baseAttack    int
		skillMulti    float64
		damageType    string
		ignoreDodge   bool
		checkMin      bool
		checkMax      bool
		minDamage     int
		maxDamage     int
	}{
		{
			name:        "基础物理伤害",
			attacker:     attacker,
			defender:     defender,
			baseAttack:   20,
			skillMulti:   1.0,
			damageType:   "physical",
			ignoreDodge:  false,
			checkMin:     true,
			checkMax:     false,
			minDamage:    1,
		},
		{
			name:        "基础法术伤害",
			attacker:     attacker,
			defender:     defender,
			baseAttack:   15,
			skillMulti:   1.0,
			damageType:   "magic",
			ignoreDodge:  false,
			checkMin:     true,
			checkMax:     false,
			minDamage:    1,
		},
		{
			name:        "nil攻击者",
			attacker:     nil,
			defender:     defender,
			baseAttack:   20,
			skillMulti:   1.0,
			damageType:   "physical",
			ignoreDodge:  false,
			checkMin:     true,
			checkMax:     false,
			minDamage:    1,
		},
		{
			name:        "负数攻击力",
			attacker:     attacker,
			defender:     defender,
			baseAttack:   -10,
			skillMulti:   1.0,
			damageType:   "physical",
			ignoreDodge:  false,
			checkMin:     true,
			checkMax:     false,
			minDamage:    1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calc.CalculateDamage(tt.attacker, tt.defender, tt.baseAttack, tt.skillMulti, tt.damageType, tt.ignoreDodge)
			
			assert.NotNil(t, result, "结果不应该为nil")
			
			if tt.checkMin {
				assert.GreaterOrEqual(t, result.FinalDamage, tt.minDamage, "最终伤害应该至少为%d", tt.minDamage)
			}
			
			if tt.checkMax {
				assert.LessOrEqual(t, result.FinalDamage, tt.maxDamage, "最终伤害应该不超过%d", tt.maxDamage)
			}
			
			// 如果闪避，伤害应该为0
			if result.IsDodged {
				assert.Equal(t, 0, result.FinalDamage, "闪避时伤害应该为0")
			}
			
		})
	}
}

// TestCalculateDefenseReduction 已在 battle_system_core_test.go 中定义

// ═══════════════════════════════════════════════════════════
// 治疗计算测试
// ═══════════════════════════════════════════════════════════
// 注意：治疗计算的详细测试用例已迁移到 YAML 格式
// 文件位置：server/internal/test/cases/calculator/（待补充）
// 根据 PRINCIPLES.md，测试用例应使用 YAML 自然语言描述
// 以下Go测试用例保留用于快速验证和调试（特殊情况）

func TestCalculateHealing(t *testing.T) {
	calc := NewCalculator()

	healer := &models.Character{
		MagicAttack: 20,
	}

	tests := []struct {
		name                string
		healer              *models.Character
		target              *models.Character
		baseHealing         int
		skillMulti          float64
		healingBonusPercent float64
		expectedMin         int
		expectedMax         int
	}{
		{
			name:                "正常治疗",
			healer:              healer,
			target:              &models.Character{HP: 50, MaxHP: 100},
			baseHealing:         20,
			skillMulti:          1.0,
			healingBonusPercent: 0,
			expectedMin:         20,
			expectedMax:         20,
		},
		{
			name:                "治疗加成",
			healer:              healer,
			target:              &models.Character{HP: 50, MaxHP: 100},
			baseHealing:         20,
			skillMulti:          1.0,
			healingBonusPercent: 10, // 10%加成
			expectedMin:         22,
			expectedMax:         22,
		},
		{
			name:                "过量治疗",
			healer:              healer,
			target:              &models.Character{HP: 95, MaxHP: 100},
			baseHealing:         20,
			skillMulti:          1.0,
			healingBonusPercent: 0,
			expectedMin:         5, // 只能治疗5点
			expectedMax:         5,
		},
		{
			name:                "nil治疗者",
			healer:              nil,
			target:              &models.Character{HP: 50, MaxHP: 100},
			baseHealing:         20,
			skillMulti:          1.0,
			healingBonusPercent: 0,
			expectedMin:         0,
			expectedMax:         0,
		},
		{
			name:                "负数基础治疗",
			healer:              healer,
			target:              &models.Character{HP: 50, MaxHP: 100},
			baseHealing:         -10,
			skillMulti:          1.0,
			healingBonusPercent: 0,
			expectedMin:         0,
			expectedMax:         0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calc.CalculateHealing(tt.healer, tt.target, tt.baseHealing, tt.skillMulti, tt.healingBonusPercent)
			
			assert.NotNil(t, result, "结果不应该为nil")
			assert.GreaterOrEqual(t, result.ActualHealing, tt.expectedMin, "实际治疗应该至少为%d", tt.expectedMin)
			assert.LessOrEqual(t, result.ActualHealing, tt.expectedMax, "实际治疗应该不超过%d", tt.expectedMax)
			assert.GreaterOrEqual(t, result.ActualHealing, 0, "实际治疗应该至少为0")
			
			// 过量治疗检查
			if result.Overhealing > 0 {
				assert.Equal(t, result.FinalHealing, result.ActualHealing+result.Overhealing, "最终治疗应该等于实际治疗+过量治疗")
			}
		})
	}
}

// ═══════════════════════════════════════════════════════════
// 辅助函数测试
// ═══════════════════════════════════════════════════════════

func TestClamp(t *testing.T) {
	calc := NewCalculator()

	tests := []struct {
		name     string
		value    float64
		min      float64
		max      float64
		expected float64
	}{
		{
			name:     "正常值",
			value:    5.0,
			min:      0.0,
			max:      10.0,
			expected: 5.0,
		},
		{
			name:     "低于最小值",
			value:    -5.0,
			min:      0.0,
			max:      10.0,
			expected: 0.0,
		},
		{
			name:     "高于最大值",
			value:    15.0,
			min:      0.0,
			max:      10.0,
			expected: 10.0,
		},
		{
			name:     "等于最小值",
			value:    0.0,
			min:      0.0,
			max:      10.0,
			expected: 0.0,
		},
		{
			name:     "等于最大值",
			value:    10.0,
			min:      0.0,
			max:      10.0,
			expected: 10.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calc.Clamp(tt.value, tt.min, tt.max)
			assert.Equal(t, tt.expected, result, "Clamp结果不正确")
			assert.GreaterOrEqual(t, result, tt.min, "结果应该至少为最小值")
			assert.LessOrEqual(t, result, tt.max, "结果应该不超过最大值")
		})
	}
}

func TestClampInt(t *testing.T) {
	calc := NewCalculator()

	tests := []struct {
		name     string
		value    int
		min      int
		max      int
		expected int
	}{
		{
			name:     "正常值",
			value:    5,
			min:      0,
			max:      10,
			expected: 5,
		},
		{
			name:     "低于最小值",
			value:    -5,
			min:      0,
			max:      10,
			expected: 0,
		},
		{
			name:     "高于最大值",
			value:    15,
			min:      0,
			max:      10,
			expected: 10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calc.ClampInt(tt.value, tt.min, tt.max)
		assert.Equal(t, tt.expected, result, "ClampInt结果不正确")
		assert.GreaterOrEqual(t, result, tt.min, "结果应该至少为最小值")
		assert.LessOrEqual(t, result, tt.max, "结果应该不超过最大值")
		})
	}
}

// TestCalculateSpeed 已在 battle_system_core_test.go 中定义

// ═══════════════════════════════════════════════════════════
// 资源计算测试
// ═══════════════════════════════════════════════════════════
// 注意：资源计算的详细测试用例已迁移到 YAML 格式
// 文件位置：server/internal/test/cases/calculator/resource_calculation.yaml
// 根据 PRINCIPLES.md，测试用例应使用 YAML 自然语言描述

// ═══════════════════════════════════════════════════════════
// 边界情况测试
// ═══════════════════════════════════════════════════════════
// 注意：边界情况的详细测试用例已迁移到 YAML 格式
// 文件位置：server/internal/test/cases/calculator/edge_cases.yaml
// 根据 PRINCIPLES.md，测试用例应使用 YAML 自然语言描述
// 以下Go测试用例保留用于快速验证和调试（特殊情况）

func TestCalculator_EdgeCases(t *testing.T) {
	calc := NewCalculator()

	t.Run("极端属性值", func(t *testing.T) {
		char := &models.Character{
			Strength:   10000,
			Agility:    10000,
			Intellect:  10000,
			Spirit:     10000,
			Stamina:    10000,
		}

		// 测试所有计算函数不会崩溃
		_ = calc.CalculatePhysicalAttack(char)
		_ = calc.CalculateMagicAttack(char)
		_ = calc.CalculateHP(char, 100)
		_ = calc.CalculateMP(char, 100)
		_ = calc.CalculatePhysCritRate(char)
		_ = calc.CalculateSpellCritRate(char)
		_ = calc.CalculateDodgeRate(char)
		_ = calc.CalculateSpeed(char)

		// 验证上限
		critRate := calc.CalculatePhysCritRate(char)
		assert.LessOrEqual(t, critRate, 0.5, "暴击率应该不超过50%")

		dodgeRate := calc.CalculateDodgeRate(char)
		assert.LessOrEqual(t, dodgeRate, 0.5, "闪避率应该不超过50%")
	})

	t.Run("负数属性值", func(t *testing.T) {
		char := &models.Character{
			Strength:   -10,
			Agility:    -10,
			Intellect:  -10,
			Spirit:     -10,
			Stamina:    -10,
		}

		// 测试所有计算函数不会崩溃
		physAttack := calc.CalculatePhysicalAttack(char)
		assert.GreaterOrEqual(t, physAttack, 1, "物理攻击应该至少为1")

		magicAttack := calc.CalculateMagicAttack(char)
		assert.GreaterOrEqual(t, magicAttack, 1, "法术攻击应该至少为1")

		hp := calc.CalculateHP(char, 10)
		assert.GreaterOrEqual(t, hp, 1, "最大HP应该至少为1")
	})

	t.Run("伤害计算边界情况", func(t *testing.T) {
		attacker := &models.Character{
			PhysicalAttack: 1,
			PhysCritRate:   0.0,
		}
		defender := &models.Character{
			PhysicalDefense: 1000, // 极高防御
			DodgeRate:       0.0,
		}

		result := calc.CalculateDamage(attacker, defender, 1, 1.0, "physical", false)
		assert.GreaterOrEqual(t, result.FinalDamage, 1, "即使防御很高，也应该至少造成1点伤害")

	})
}


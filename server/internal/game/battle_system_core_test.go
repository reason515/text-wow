package game

import (
	"fmt"
	"testing"

	"text-wow/internal/models"

	"github.com/stretchr/testify/assert"
)

// ═══════════════════════════════════════════════════════════
// 伤害计算测试
// ═══════════════════════════════════════════════════════════

func TestCalculateDamage_PhysicalDamage(t *testing.T) {
	calc := NewCalculator()
	
	attacker := &models.Character{
		PhysicalAttack: 100,
		PhysCritRate:   0.0,
		PhysCritDamage: 2.0,
	}
	
	defender := &models.Character{
		PhysicalDefense: 50,
		DodgeRate:       0.0,
	}
	
	result := calc.CalculateDamage(attacker, defender, 100, 1.0, "physical", false)
	
	assert.NotNil(t, result)
	assert.Greater(t, result.FinalDamage, 0)
	assert.Equal(t, false, result.IsDodged)
	// 防御减伤（减法公式）: 100 - 50 = 50
	assert.Equal(t, 50, result.FinalDamage)
}

func TestCalculateDamage_MagicDamage(t *testing.T) {
	calc := NewCalculator()
	
	attacker := &models.Character{
		MagicAttack:    100,
		SpellCritRate:  0.0,
		SpellCritDamage: 2.0,
	}
	
	defender := &models.Character{
		MagicDefense: 50,
		DodgeRate:    0.0,
	}
	
	result := calc.CalculateDamage(attacker, defender, 100, 1.0, "magic", false)
	
	assert.NotNil(t, result)
	assert.Greater(t, result.FinalDamage, 0)
	assert.Equal(t, false, result.IsDodged)
	// 防御减伤（减法公式）: 100 - 50 = 50
	assert.Equal(t, 50, result.FinalDamage)
}

func TestCalculateDamage_FireDamage(t *testing.T) {
	calc := NewCalculator()
	
	attacker := &models.Character{
		MagicAttack:    100,
		SpellCritRate:  0.0,
		SpellCritDamage: 2.0,
	}
	
	defender := &models.Character{
		MagicDefense: 50,
		DodgeRate:    0.0,
	}
	
	result := calc.CalculateDamage(attacker, defender, 100, 1.0, "fire", false)
	
	assert.NotNil(t, result)
	assert.Greater(t, result.FinalDamage, 0)
	// 元素伤害使用魔法防御（减法公式）: 100 - 50 = 50
	assert.Equal(t, 50, result.FinalDamage)
}

func TestCalculateDamage_Dodge(t *testing.T) {
	calc := NewCalculator()
	
	attacker := &models.Character{
		PhysicalAttack: 100,
		PhysCritRate:   0.0,
	}
	
	defender := &models.Character{
		PhysicalDefense: 0,
		DodgeRate:       1.0, // 100%闪避率（会被限制到50%）
	}
	
	// 由于闪避是随机的，我们需要多次测试来验证闪避机制
	dodgedCount := 0
	totalTests := 100
	
	for i := 0; i < totalTests; i++ {
		result := calc.CalculateDamage(attacker, defender, 100, 1.0, "physical", false)
		if result.IsDodged {
			dodgedCount++
			assert.Equal(t, 0, result.FinalDamage)
		}
	}
	
	// 由于闪避率被限制到50%，应该有大约50%的闪避率
	// 允许一些误差（30%-70%）
	assert.Greater(t, dodgedCount, totalTests*30/100)
	assert.Less(t, dodgedCount, totalTests*70/100)
}

func TestCalculateDamage_Crit(t *testing.T) {
	calc := NewCalculator()
	
	attacker := &models.Character{
		PhysicalAttack: 100,
		PhysCritRate:   1.0, // 100%暴击率（会被限制到50%）
		PhysCritDamage: 2.0,
	}
	
	defender := &models.Character{
		PhysicalDefense: 0,
		DodgeRate:       0.0,
	}
	
	// 由于暴击是随机的，我们需要多次测试来验证暴击机制
	critCount := 0
	totalTests := 100
	
	for i := 0; i < totalTests; i++ {
		result := calc.CalculateDamage(attacker, defender, 100, 1.0, "physical", false)
		if result.IsCrit {
			critCount++
			assert.Greater(t, result.FinalDamage, 100) // 暴击伤害应该更高
		}
	}
	
	// 由于暴击率被限制到50%，应该有大约50%的暴击率
	// 允许一些误差（30%-70%）
	assert.Greater(t, critCount, totalTests*30/100)
	assert.Less(t, critCount, totalTests*70/100)
}

func TestCalculateDamage_AllElementTypes(t *testing.T) {
	calc := NewCalculator()
	
	attacker := &models.Character{
		MagicAttack:    100,
		SpellCritRate:  0.0,
		SpellCritDamage: 2.0,
	}
	
	defender := &models.Character{
		MagicDefense: 50,
		DodgeRate:    0.0,
	}
	
	elementTypes := []string{"fire", "frost", "shadow", "holy", "nature"}
	
	for _, damageType := range elementTypes {
		result := calc.CalculateDamage(attacker, defender, 100, 1.0, damageType, false)
		assert.NotNil(t, result)
		assert.Greater(t, result.FinalDamage, 0)
		// 所有元素伤害都使用魔法防御（减法公式）: 100 - 50 = 50
		assert.Equal(t, 50, result.FinalDamage)
	}
}

func TestCalculateDamage_InvalidType(t *testing.T) {
	calc := NewCalculator()
	
	attacker := &models.Character{
		PhysicalAttack: 100,
	}
	
	defender := &models.Character{
		PhysicalDefense: 0,
	}
	
	// 无效的伤害类型应该默认为物理伤害
	result := calc.CalculateDamage(attacker, defender, 100, 1.0, "invalid_type", false)
	assert.NotNil(t, result)
	assert.Greater(t, result.FinalDamage, 0)
}

func TestCalculateDamage_NilPointers(t *testing.T) {
	calc := NewCalculator()
	
	// nil指针应该返回最小伤害
	result := calc.CalculateDamage(nil, nil, 100, 1.0, "physical", false)
	assert.NotNil(t, result)
	assert.Equal(t, 1, result.FinalDamage) // 至少1点伤害
}

func TestCalculateDamage_NegativeInputs(t *testing.T) {
	calc := NewCalculator()
	
	attacker := &models.Character{
		PhysicalAttack: 100,
	}
	
	defender := &models.Character{
		PhysicalDefense: 0,
	}
	
	// 负数输入应该被处理为0
	result := calc.CalculateDamage(attacker, defender, -100, -1.0, "physical", false)
	assert.NotNil(t, result)
	assert.Equal(t, 1, result.FinalDamage) // 至少1点伤害
}

// ═══════════════════════════════════════════════════════════
// 治疗计算测试
// ═══════════════════════════════════════════════════════════

func TestCalculateHealing_BasicHealing(t *testing.T) {
	calc := NewCalculator()
	
	healer := &models.Character{
		MagicAttack: 100, // 治疗力
	}
	
	target := &models.Character{
		HP:    50,
		MaxHP: 100,
	}
	
	result := calc.CalculateHealing(healer, target, 50, 1.0, 0.0)
	
	assert.NotNil(t, result)
	assert.Equal(t, 50, result.FinalHealing)
	assert.Equal(t, 50, result.ActualHealing) // 实际治疗 = 最终治疗（没有过量治疗）
	assert.Equal(t, 0, result.Overhealing)
}

func TestCalculateHealing_WithHealingBonus(t *testing.T) {
	calc := NewCalculator()
	
	healer := &models.Character{
		MagicAttack: 100,
	}
	
	target := &models.Character{
		HP:    50,
		MaxHP: 100,
	}
	
	// 50%治疗加成
	result := calc.CalculateHealing(healer, target, 50, 1.0, 50.0)
	
	assert.NotNil(t, result)
	assert.Equal(t, 75, result.FinalHealing) // 50 * 1.5 = 75
	assert.Equal(t, 50, result.ActualHealing) // 实际治疗 = 50（因为目标只能恢复到100HP）
	assert.Equal(t, 25, result.Overhealing)   // 过量治疗 = 25
}

func TestCalculateHealing_Overhealing(t *testing.T) {
	calc := NewCalculator()
	
	healer := &models.Character{
		MagicAttack: 100,
	}
	
	target := &models.Character{
		HP:    90,
		MaxHP: 100,
	}
	
	// 治疗100点，但目标只能恢复10点
	result := calc.CalculateHealing(healer, target, 100, 1.0, 0.0)
	
	assert.NotNil(t, result)
	assert.Equal(t, 100, result.FinalHealing)
	assert.Equal(t, 10, result.ActualHealing) // 实际治疗 = 10
	assert.Equal(t, 90, result.Overhealing)   // 过量治疗 = 90
}

func TestCalculateHealing_FullHP(t *testing.T) {
	calc := NewCalculator()
	
	healer := &models.Character{
		MagicAttack: 100,
	}
	
	target := &models.Character{
		HP:    100,
		MaxHP: 100,
	}
	
	result := calc.CalculateHealing(healer, target, 50, 1.0, 0.0)
	
	assert.NotNil(t, result)
	assert.Equal(t, 50, result.FinalHealing)
	assert.Equal(t, 0, result.ActualHealing) // 实际治疗 = 0（目标已满血）
	assert.Equal(t, 50, result.Overhealing) // 全部过量治疗
}

func TestCalculateHealing_NilPointers(t *testing.T) {
	calc := NewCalculator()
	
	result := calc.CalculateHealing(nil, nil, 50, 1.0, 0.0)
	assert.NotNil(t, result)
	assert.Equal(t, 0, result.FinalHealing)
	assert.Equal(t, 0, result.ActualHealing)
}

func TestCalculateHealing_NegativeInputs(t *testing.T) {
	calc := NewCalculator()
	
	healer := &models.Character{
		MagicAttack: 100,
	}
	
	target := &models.Character{
		HP:    50,
		MaxHP: 100,
	}
	
	// 负数输入应该被处理为0
	result := calc.CalculateHealing(healer, target, -50, -1.0, -10.0)
	assert.NotNil(t, result)
	assert.Equal(t, 0, result.FinalHealing)
	assert.Equal(t, 0, result.ActualHealing)
}

// ═══════════════════════════════════════════════════════════
// 速度计算和排序测试
// ═══════════════════════════════════════════════════════════

func TestCalculateSpeed(t *testing.T) {
	calc := NewCalculator()
	
	char := &models.Character{
		Agility: 50,
		Level:   1,
	}
	
	speed := calc.CalculateSpeed(char)
	
	assert.Greater(t, speed, 0)
	// 速度 = 敏捷
	// = 50
	assert.Equal(t, 50, speed)
}

func TestCalculateSpeed_ZeroAgility(t *testing.T) {
	calc := NewCalculator()
	
	char := &models.Character{
		Agility: 0,
		Level:   1,
	}
	
	speed := calc.CalculateSpeed(char)
	// 敏捷为0时，速度最小值为1
	assert.Equal(t, 1, speed)
}

func TestCalculateSpeed_NilPointer(t *testing.T) {
	calc := NewCalculator()
	
	speed := calc.CalculateSpeed(nil)
	// nil指针应该返回基础速度10
	assert.Equal(t, 10, speed)
}

func TestBuildTurnOrder(t *testing.T) {
	manager := NewBattleManager()
	session := manager.GetOrCreateSession(1)
	
	characters := []*models.Character{
		{
			ID:      1,
			Name:    "Char1",
			HP:      100,
			MaxHP:   100,
			Agility: 30, // 速度 = 60
		},
		{
			ID:      2,
			Name:    "Char2",
			HP:      100,
			MaxHP:   100,
			Agility: 50, // 速度 = 100
		},
	}
	
	enemies := []*models.Monster{
		{
			ID:     "enemy1",
			Name:   "Enemy1",
			HP:     100,
			MaxHP:  100,
			Speed:  80,
		},
		{
			ID:     "enemy2",
			Name:   "Enemy2",
			HP:     100,
			MaxHP:  100,
			Speed:  40,
		},
	}
	
	manager.buildTurnOrder(session, characters, enemies)
	
	assert.NotNil(t, session.TurnOrder)
	assert.Equal(t, 4, len(session.TurnOrder))
	
	// 验证排序：速度从高到低（速度相同时随机，所以只验证顺序性）
	// 验证所有参与者都在队列中
	participantIDs := make(map[string]bool)
	speeds := make([]int, 0, 4)
	
	for _, p := range session.TurnOrder {
		assert.NotNil(t, p)
		if p.Type == "character" {
			assert.NotNil(t, p.Character)
			participantIDs[fmt.Sprintf("char_%d", p.Character.ID)] = true
			speeds = append(speeds, p.Speed)
		} else {
			assert.NotNil(t, p.Monster)
			participantIDs[fmt.Sprintf("enemy_%s", p.Monster.ID)] = true
			speeds = append(speeds, p.Speed)
		}
	}
	
	// 验证所有参与者都在队列中
	assert.True(t, participantIDs["char_1"])
	assert.True(t, participantIDs["char_2"])
	assert.True(t, participantIDs["enemy_enemy1"])
	assert.True(t, participantIDs["enemy_enemy2"])
	
	// 验证速度排序（从高到低，允许速度相同时的随机性）
	for i := 0; i < len(speeds)-1; i++ {
		assert.GreaterOrEqual(t, speeds[i], speeds[i+1], "速度应该从高到低排序")
	}
}

func TestBuildTurnOrder_DeadParticipants(t *testing.T) {
	manager := NewBattleManager()
	session := manager.GetOrCreateSession(1)
	
	characters := []*models.Character{
		{
			ID:      1,
			Name:    "Char1",
			HP:      0, // 死亡
			MaxHP:   100,
			Agility: 30,
		},
		{
			ID:      2,
			Name:    "Char2",
			HP:      100,
			MaxHP:   100,
			Agility: 50,
		},
	}
	
	enemies := []*models.Monster{
		{
			ID:     "enemy1",
			Name:   "Enemy1",
			HP:     0, // 死亡
			MaxHP:  100,
			Speed:  80,
		},
		{
			ID:     "enemy2",
			Name:   "Enemy2",
			HP:     100,
			MaxHP:  100,
			Speed:  40,
		},
	}
	
	manager.buildTurnOrder(session, characters, enemies)
	
	// 只有存活的参与者应该在队列中
	assert.Equal(t, 2, len(session.TurnOrder))
	
	// 验证只有存活的参与者在队列中
	for _, p := range session.TurnOrder {
		if p.Type == "character" {
			assert.Greater(t, p.Character.HP, 0)
		} else {
			assert.Greater(t, p.Monster.HP, 0)
		}
	}
}

func TestBuildTurnOrder_EmptyLists(t *testing.T) {
	manager := NewBattleManager()
	session := manager.GetOrCreateSession(1)
	
	manager.buildTurnOrder(session, []*models.Character{}, []*models.Monster{})
	
	assert.NotNil(t, session.TurnOrder)
	assert.Equal(t, 0, len(session.TurnOrder))
}

// ═══════════════════════════════════════════════════════════
// 防御减伤测试
// ═══════════════════════════════════════════════════════════

// TestCalculateDefenseReduction 测试防御减伤计算（已废弃）
// 注意：防御现在使用减法公式（伤害 = 攻击 - 防御），不再使用百分比减伤
// 此函数已废弃，保留测试用于向后兼容
func TestCalculateDefenseReduction(t *testing.T) {
	calc := NewCalculator()
	
	// CalculateDefenseReduction 现在返回 0（已废弃）
	reduction := calc.CalculateDefenseReduction(0)
	assert.Equal(t, 0.0, reduction)
	
	reduction = calc.CalculateDefenseReduction(100)
	assert.Equal(t, 0.0, reduction)
	
	reduction = calc.CalculateDefenseReduction(300)
	assert.Equal(t, 0.0, reduction)
	
	reduction = calc.CalculateDefenseReduction(1000)
	assert.Equal(t, 0.0, reduction)
}

func TestCalculateDefenseReduction_Negative(t *testing.T) {
	calc := NewCalculator()
	
	// 负数防御应该被处理为0
	reduction := calc.CalculateDefenseReduction(-100)
	assert.Equal(t, 0.0, reduction)
}

// TestDefenseSubtraction 测试防御减法公式
func TestDefenseSubtraction(t *testing.T) {
	calc := NewCalculator()
	
	// 创建测试角色
	attacker := &models.Character{
		ID:             1,
		PhysicalAttack: 100,
		MagicAttack:    100,
	}
	
	defender := &models.Character{
		ID:              2,
		PhysicalDefense: 30,
		MagicDefense:    40,
	}
	
	// 测试物理伤害：100 攻击 - 30 防御 = 70 伤害
	result := calc.CalculateDamage(attacker, defender, 100, 1.0, "physical", false)
	assert.Equal(t, 70, result.FinalDamage)
	assert.Equal(t, 70.0, result.DamageAfterDefense)
	
	// 测试魔法伤害：100 攻击 - 40 防御 = 60 伤害
	result = calc.CalculateDamage(attacker, defender, 100, 1.0, "magic", false)
	assert.Equal(t, 60, result.FinalDamage)
	assert.Equal(t, 60.0, result.DamageAfterDefense)
	
	// 测试防御大于攻击：50 攻击 - 100 防御 = 1 伤害（最低）
	highDefenseDefender := &models.Character{
		ID:              3,
		PhysicalDefense: 100,
		MagicDefense:    40,
	}
	result = calc.CalculateDamage(attacker, highDefenseDefender, 50, 1.0, "physical", false)
	assert.Equal(t, 1, result.FinalDamage)
	assert.Equal(t, 1.0, result.DamageAfterDefense)
	
	// 测试元素伤害（使用魔法防御）：100 攻击 - 40 防御 = 60 伤害
	result = calc.CalculateDamage(attacker, defender, 100, 1.0, "fire", false)
	assert.Equal(t, 60, result.FinalDamage)
	assert.Equal(t, 60.0, result.DamageAfterDefense)
}


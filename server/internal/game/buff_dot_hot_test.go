package game

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// ═══════════════════════════════════════════════════════════
// DOT/HOT效果测试
// ═══════════════════════════════════════════════════════════

func TestBuffManager_ProcessDOTEffects(t *testing.T) {
	bm := NewBuffManager()
	characterID := 1

	// 应用DOT效果（每回合造成10点伤害）
	bm.ApplyBuffWithDOT(characterID, "dot_poison", "毒药", "debuff", false, 5, 10.0, "damage", "nature", true, false, 0)

	// 第一回合：应该造成10点伤害
	damage, _ := bm.ProcessDOTEffects(characterID, 1)
	assert.Equal(t, 10, damage)

	// 第二回合：应该再次造成10点伤害
	damage, _ = bm.ProcessDOTEffects(characterID, 2)
	assert.Equal(t, 10, damage)

	// 减少持续时间
	bm.TickBuffs(characterID)
	bm.TickBuffs(characterID)
	bm.TickBuffs(characterID)
	bm.TickBuffs(characterID)
	bm.TickBuffs(characterID)

	// DOT应该过期，不再造成伤害
	damage, _ = bm.ProcessDOTEffects(characterID, 6)
	assert.Equal(t, 0, damage)
}

func TestBuffManager_ProcessHOTEffects(t *testing.T) {
	bm := NewBuffManager()
	characterID := 1

	// 应用HOT效果（每回合恢复15点生命值）
	bm.ApplyBuffWithDOT(characterID, "hot_regen", "恢复", "buff", true, 4, 15.0, "healing", "", false, true, 0)

	// 第一回合：应该恢复15点生命值
	damage, healing := bm.ProcessDOTEffects(characterID, 1)
	assert.Equal(t, 0, damage)
	assert.Equal(t, 15, healing)

	// 第二回合：应该再次恢复15点生命值
	_, healing = bm.ProcessDOTEffects(characterID, 2)
	assert.Equal(t, 15, healing)

	// 减少持续时间
	bm.TickBuffs(characterID)
	bm.TickBuffs(characterID)
	bm.TickBuffs(characterID)
	bm.TickBuffs(characterID)

	// HOT应该过期，不再恢复生命值
	_, healing = bm.ProcessDOTEffects(characterID, 5)
	assert.Equal(t, 0, healing)
}

func TestBuffManager_ProcessDOTEffects_WithInterval(t *testing.T) {
	bm := NewBuffManager()
	characterID := 1

	// 应用DOT效果（每2回合触发一次）
	bm.ApplyBuffWithDOT(characterID, "dot_bleed", "流血", "debuff", false, 6, 20.0, "damage", "physical", true, false, 2)

	// 第一回合：应该触发（LastTick=0，间隔2）
	damage, _ := bm.ProcessDOTEffects(characterID, 1)
	assert.Equal(t, 20, damage)

	// 第二回合：不应该触发（LastTick=1，间隔2，还没到）
	damage, _ = bm.ProcessDOTEffects(characterID, 2)
	assert.Equal(t, 0, damage)

	// 第三回合：应该触发（LastTick=1，间隔2，3-1=2>=2）
	damage, _ = bm.ProcessDOTEffects(characterID, 3)
	assert.Equal(t, 20, damage)

	// 第四回合：不应该触发
	damage, _ = bm.ProcessDOTEffects(characterID, 4)
	assert.Equal(t, 0, damage)

	// 第五回合：应该触发
	damage, _ = bm.ProcessDOTEffects(characterID, 5)
	assert.Equal(t, 20, damage)
}

func TestBuffManager_ProcessEnemyDOTEffects(t *testing.T) {
	bm := NewBuffManager()
	enemyID := "enemy_1"

	// 应用敌人DOT效果
	bm.ApplyEnemyDebuffWithDOT(enemyID, "dot_poison", "毒药", "debuff", 3, 12.0, "damage", "nature", true, 0)

	// 第一回合：应该造成12点伤害
	damage := bm.ProcessEnemyDOTEffects(enemyID, 1)
	assert.Equal(t, 12, damage)

	// 第二回合：应该再次造成12点伤害
	damage = bm.ProcessEnemyDOTEffects(enemyID, 2)
	assert.Equal(t, 12, damage)

	// 减少持续时间
	bm.TickEnemyDebuffs(enemyID)
	bm.TickEnemyDebuffs(enemyID)
	bm.TickEnemyDebuffs(enemyID)

	// DOT应该过期，不再造成伤害
	damage = bm.ProcessEnemyDOTEffects(enemyID, 4)
	assert.Equal(t, 0, damage)
}

func TestBuffManager_DOTAndHOTTogether(t *testing.T) {
	bm := NewBuffManager()
	characterID := 1

	// 同时应用DOT和HOT
	bm.ApplyBuffWithDOT(characterID, "dot_poison", "毒药", "debuff", false, 3, 10.0, "damage", "nature", true, false, 0)
	bm.ApplyBuffWithDOT(characterID, "hot_regen", "恢复", "buff", true, 3, 15.0, "healing", "", false, true, 0)

	// 应该同时造成伤害和恢复生命值
	damage, healing := bm.ProcessDOTEffects(characterID, 1)
	assert.Equal(t, 10, damage)
	assert.Equal(t, 15, healing)
	
	// 验证两个效果都生效
	assert.Greater(t, damage, 0)
	assert.Greater(t, healing, 0)
}

// ═══════════════════════════════════════════════════════════
// Buff叠加规则测试
// ═══════════════════════════════════════════════════════════

func TestBuffManager_StackingRule_Refresh(t *testing.T) {
	bm := NewBuffManager()
	characterID := 1

	// 应用战斗怒吼（refresh规则）
	bm.ApplyBuff(characterID, "battle_shout", "战斗怒吼", "buff", true, 5, 20.0, "attack", "")
	
	buffs := bm.GetBuffs(characterID)
	buff := buffs["battle_shout"]
	assert.Equal(t, 5, buff.Duration)
	assert.Equal(t, 20.0, buff.Value)

	// 再次应用相同Buff（应该刷新持续时间，不叠加数值）
	bm.ApplyBuff(characterID, "battle_shout", "战斗怒吼", "buff", true, 8, 20.0, "attack", "")
	
	buffs = bm.GetBuffs(characterID)
	buff = buffs["battle_shout"]
	assert.Equal(t, 8, buff.Duration) // 持续时间刷新为8
	assert.Equal(t, 20.0, buff.Value)  // 数值不变

	// 如果新持续时间更短，应该保持原持续时间
	bm.ApplyBuff(characterID, "battle_shout", "战斗怒吼", "buff", true, 3, 20.0, "attack", "")
	
	buffs = bm.GetBuffs(characterID)
	buff = buffs["battle_shout"]
	assert.Equal(t, 8, buff.Duration) // 保持较长的持续时间
}

func TestBuffManager_StackingRule_Stack(t *testing.T) {
	bm := NewBuffManager()
	characterID := 1

	// 应用DOT效果（stack规则）
	bm.ApplyBuffWithDOT(characterID, "dot_poison", "毒药", "debuff", false, 5, 10.0, "damage", "nature", true, false, 0)
	
	buffs := bm.GetBuffs(characterID)
	buff := buffs["dot_poison"]
	assert.Equal(t, 5, buff.Duration)
	assert.Equal(t, 10.0, buff.Value)

	// 再次应用相同DOT（应该叠加数值）
	bm.ApplyBuffWithDOT(characterID, "dot_poison", "毒药", "debuff", false, 5, 10.0, "damage", "nature", true, false, 0)
	
	buffs = bm.GetBuffs(characterID)
	buff = buffs["dot_poison"]
	assert.Equal(t, 5, buff.Duration) // 持续时间刷新
	assert.Equal(t, 20.0, buff.Value) // 数值叠加：10 + 10 = 20
}

func TestBuffManager_StackingRule_Replace(t *testing.T) {
	bm := NewBuffManager()
	characterID := 1

	// 应用Buff（replace规则）
	bm.ApplyBuff(characterID, "test_replace", "测试替换", "buff", true, 5, 15.0, "attack", "")
	
	buffs := bm.GetBuffs(characterID)
	buff := buffs["test_replace"]
	assert.Equal(t, 15.0, buff.Value)

	// 再次应用（refresh规则：保持较长的持续时间，数值不变）
	bm.ApplyBuff(characterID, "test_replace", "测试替换", "buff", true, 3, 25.0, "attack", "")
	
	buffs = bm.GetBuffs(characterID)
	buff = buffs["test_replace"]
	// refresh规则：保持较长的持续时间，数值不变
	assert.Equal(t, 5, buff.Duration) // 保持较长的持续时间
	assert.Equal(t, 15.0, buff.Value) // 数值不变（refresh规则）
}

func TestBuffManager_GetStackingRule(t *testing.T) {
	bm := NewBuffManager()

	// 测试不同Buff的叠加规则
	testCases := []struct {
		effectID     string
		expectedRule string
	}{
		{"battle_shout", "refresh"},
		{"shield_block", "refresh"},
		{"dot_poison", "stack"},
		{"dot_bleed", "stack"},
		{"hot_regen", "stack"},
		{"unknown_buff", "refresh"}, // 默认规则
	}

	for _, tc := range testCases {
		rule := bm.getStackingRule(tc.effectID)
		assert.Equal(t, tc.expectedRule, rule, "EffectID: %s should have rule: %s", tc.effectID, tc.expectedRule)
	}
}


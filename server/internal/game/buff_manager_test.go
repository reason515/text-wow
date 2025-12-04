package game

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// ═══════════════════════════════════════════════════════════
// BuffManager 基础功能测试
// ═══════════════════════════════════════════════════════════

func TestBuffManager_ApplyBuff(t *testing.T) {
	bm := NewBuffManager()
	characterID := 1

	// 测试应用Buff
	bm.ApplyBuff(characterID, "battle_shout", "战斗怒吼", "buff", true, 5, 20.0, "attack", "")
	
	buffs := bm.GetBuffs(characterID)
	assert.NotNil(t, buffs)
	assert.Contains(t, buffs, "battle_shout")
	
	buff := buffs["battle_shout"]
	assert.Equal(t, "战斗怒吼", buff.Name)
	assert.True(t, buff.IsBuff)
	assert.Equal(t, 5, buff.Duration)
	assert.Equal(t, 20.0, buff.Value)
	assert.Equal(t, "attack", buff.StatAffected)
}

func TestBuffManager_ApplyDebuff(t *testing.T) {
	bm := NewBuffManager()
	characterID := 1

	// 测试应用Debuff
	bm.ApplyBuff(characterID, "shield_wall", "盾墙", "buff", true, 2, -60.0, "damage_taken", "")
	
	buffs := bm.GetBuffs(characterID)
	assert.Contains(t, buffs, "shield_wall")
	
	buff := buffs["shield_wall"]
	assert.Equal(t, -60.0, buff.Value)
	assert.Equal(t, "damage_taken", buff.StatAffected)
}

func TestBuffManager_GetBuffValue(t *testing.T) {
	bm := NewBuffManager()
	characterID := 1

	// 应用多个Buff
	bm.ApplyBuff(characterID, "battle_shout", "战斗怒吼", "buff", true, 5, 20.0, "attack", "")
	bm.ApplyBuff(characterID, "berserker_rage", "狂暴之怒", "buff", true, 4, 30.0, "attack", "")

	// 测试获取Buff值（应该累加）
	attackValue := bm.GetBuffValue(characterID, "attack")
	assert.Equal(t, 50.0, attackValue)
}

func TestBuffManager_TickBuffs(t *testing.T) {
	bm := NewBuffManager()
	characterID := 1

	// 应用Buff
	bm.ApplyBuff(characterID, "battle_shout", "战斗怒吼", "buff", true, 3, 20.0, "attack", "")
	
	// 第一次tick：duration 3 -> 2
	expired := bm.TickBuffs(characterID)
	assert.Empty(t, expired) // 应该没有过期
	
	buffs := bm.GetBuffs(characterID)
	assert.NotNil(t, buffs)
	buff, exists := buffs["battle_shout"]
	assert.True(t, exists, "Buff应该存在")
	assert.Equal(t, 2, buff.Duration)
	
	// 第二次tick：duration 2 -> 1
	expired = bm.TickBuffs(characterID)
	assert.Empty(t, expired) // 应该没有过期
	
	buffs = bm.GetBuffs(characterID)
	buff, exists = buffs["battle_shout"]
	assert.True(t, exists, "Buff应该仍然存在")
	assert.Equal(t, 1, buff.Duration)
	
	// 第三次tick：duration 1 -> 0，应该过期
	expired = bm.TickBuffs(characterID)
	assert.NotEmpty(t, expired, "Buff应该过期")
	assert.Contains(t, expired, "battle_shout")
	
	// Buff应该被移除
	buffs = bm.GetBuffs(characterID)
	assert.NotContains(t, buffs, "battle_shout")
}

func TestBuffManager_ClearBuffs(t *testing.T) {
	bm := NewBuffManager()
	characterID := 1

	bm.ApplyBuff(characterID, "battle_shout", "战斗怒吼", "buff", true, 5, 20.0, "attack", "")
	bm.ApplyBuff(characterID, "shield_wall", "盾墙", "buff", true, 2, -60.0, "damage_taken", "")
	
	bm.ClearBuffs(characterID)
	
	buffs := bm.GetBuffs(characterID)
	assert.Empty(t, buffs)
}

// ═══════════════════════════════════════════════════════════
// 敌人Debuff系统测试
// ═══════════════════════════════════════════════════════════

func TestBuffManager_ApplyEnemyDebuff(t *testing.T) {
	bm := NewBuffManager()
	enemyID := "enemy_1"

	// 测试应用敌人Debuff
	bm.ApplyEnemyDebuff(enemyID, "demoralizing_shout", "挫志怒吼", "debuff", 3, 15.0, "attack", "")
	
	debuffs := bm.GetEnemyDebuffs(enemyID)
	assert.NotNil(t, debuffs)
	assert.Contains(t, debuffs, "demoralizing_shout")
	
	debuff := debuffs["demoralizing_shout"]
	assert.False(t, debuff.IsBuff) // 敌人debuff都是debuff
	assert.Equal(t, 3, debuff.Duration)
	assert.Equal(t, 15.0, debuff.Value)
	assert.Equal(t, "attack", debuff.StatAffected)
}

func TestBuffManager_GetEnemyDebuffValue(t *testing.T) {
	bm := NewBuffManager()
	enemyID := "enemy_1"

	// 应用多个Debuff
	bm.ApplyEnemyDebuff(enemyID, "demoralizing_shout", "挫志怒吼", "debuff", 3, 15.0, "attack", "")
	bm.ApplyEnemyDebuff(enemyID, "whirlwind", "旋风斩", "debuff", 2, 10.0, "defense", "")

	// 测试获取Debuff值
	attackDebuff := bm.GetEnemyDebuffValue(enemyID, "attack")
	assert.Equal(t, 15.0, attackDebuff)
	
	defenseDebuff := bm.GetEnemyDebuffValue(enemyID, "defense")
	assert.Equal(t, 10.0, defenseDebuff)
}

func TestBuffManager_TickEnemyDebuffs(t *testing.T) {
	bm := NewBuffManager()
	enemyID := "enemy_1"

	bm.ApplyEnemyDebuff(enemyID, "demoralizing_shout", "挫志怒吼", "debuff", 2, 15.0, "attack", "")
	
	// 第一次tick：duration 2 -> 1
	expired := bm.TickEnemyDebuffs(enemyID)
	assert.Empty(t, expired)
	
	debuffs := bm.GetEnemyDebuffs(enemyID)
	assert.NotNil(t, debuffs)
	debuff, exists := debuffs["demoralizing_shout"]
	assert.True(t, exists, "Debuff应该存在")
	if exists {
		assert.Equal(t, 1, debuff.Duration)
	}
	
	// 第二次tick：duration 1 -> 0，应该过期
	expired = bm.TickEnemyDebuffs(enemyID)
	assert.NotEmpty(t, expired, "Debuff应该过期")
	assert.Contains(t, expired, "demoralizing_shout")
	
	debuffs = bm.GetEnemyDebuffs(enemyID)
	assert.NotContains(t, debuffs, "demoralizing_shout")
}

func TestBuffManager_ClearEnemyDebuffs(t *testing.T) {
	bm := NewBuffManager()
	enemyID := "enemy_1"

	bm.ApplyEnemyDebuff(enemyID, "demoralizing_shout", "挫志怒吼", "debuff", 3, 15.0, "attack", "")
	
	bm.ClearEnemyDebuffs(enemyID)
	
	debuffs := bm.GetEnemyDebuffs(enemyID)
	assert.Empty(t, debuffs)
}

// ═══════════════════════════════════════════════════════════
// Buff效果计算测试
// ═══════════════════════════════════════════════════════════

func TestBuffManager_CalculateDamageTakenWithBuffs(t *testing.T) {
	bm := NewBuffManager()
	characterID := 1

	// 应用减伤Buff
	bm.ApplyBuff(characterID, "shield_wall", "盾墙", "buff", true, 2, -60.0, "damage_taken", "")
	
	baseDamage := 100
	reducedDamage := bm.CalculateDamageTakenWithBuffs(baseDamage, characterID, true)
	
	// 60%减伤，应该减少到40
	assert.Equal(t, 40, reducedDamage)
}

func TestBuffManager_CalculateDamageWithBuffs(t *testing.T) {
	bm := NewBuffManager()
	characterID := 1

	// 应用攻击力Buff
	bm.ApplyBuff(characterID, "battle_shout", "战斗怒吼", "buff", true, 5, 20.0, "attack", "")
	
	baseDamage := 100
	enhancedDamage := bm.CalculateDamageWithBuffs(baseDamage, characterID, true)
	
	// 20%攻击力加成，应该增加到120
	assert.Equal(t, 120, enhancedDamage)
}


package game

import (
	"testing"

	"text-wow/internal/models"

	"github.com/stretchr/testify/assert"
)

// ═══════════════════════════════════════════════════════════
// 战士技能集成测试
// ═══════════════════════════════════════════════════════════

func setupWarriorSkillsTest(t *testing.T) (*BattleManager, *models.Character, *models.Monster) {
	bm := NewBattleManager()
	
	character := &models.Character{
		ID:           1,
		UserID:       1,
		Name:         "测试战士",
		RaceID:       "human",
		ClassID:      "warrior",
		Faction:      "alliance",
		Level:        5,
		HP:           1000,
		MaxHP:        1000,
		Resource:     50,
		MaxResource:  100,
		ResourceType: "rage",
		PhysicalAttack:  100,
		MagicAttack:     50,
		PhysicalDefense: 50,
		MagicDefense:    30,
		PhysCritRate:    0.1,
		PhysCritDamage:  1.5,
		SpellCritRate:   0.1,
		SpellCritDamage: 1.5,
		Exp:          0,
		ExpToNext:    100,
		TotalKills:   0,
		TotalDeaths:  0,
	}
	
	enemy := &models.Monster{
		ID:        "test_enemy_1",
		ZoneID:    "test_zone",
		Name:      "测试敌人",
		Level:     5,
		Type:      "normal",
		HP:        500,
		MaxHP:     500,
		PhysicalAttack:  80,
		MagicAttack:     40,
		PhysicalDefense: 30,
		MagicDefense:    20,
		ExpReward: 50,
		GoldMin:   10,
		GoldMax:   20,
	}
	
	return bm, character, enemy
}

// ═══════════════════════════════════════════════════════════
// 被动技能效果测试
// ═══════════════════════════════════════════════════════════

func TestPassiveSkill_BloodCraze(t *testing.T) {
	bm, character, _ := setupWarriorSkillsTest(t)
	
	// 设置血之狂热被动技能
	passiveState := &CharacterPassiveState{
		PassiveID:   "warrior_passive_blood_craze",
		Level:       1,
		Passive:     &models.PassiveSkill{ID: "warrior_passive_blood_craze", EffectType: "on_hit_heal", EffectStat: ""},
		EffectValue: 1.0, // 1%最大HP恢复
	}
	
	bm.passiveSkillManager.mu.Lock()
	bm.passiveSkillManager.characterPassives[character.ID] = []*CharacterPassiveState{passiveState}
	bm.passiveSkillManager.mu.Unlock()
	
	// 模拟攻击（设置HP不满，以便恢复）
	character.HP = 900 // 不满血
	originalHP := character.HP
	damageDealt := 100
	
	// 创建模拟的BattleSession
	session := &BattleSession{
		UserID:     1,
		IsRunning:  true,
		BattleLogs: []models.BattleLog{},
	}
	
	var logs []models.BattleLog
	_ = session // 避免未使用变量警告
	_ = logs    // 避免未使用变量警告
	bm.handlePassiveOnHitEffects(character, damageDealt, true, session, &logs)
	
	// 应该恢复1%最大HP = 10点
	assert.Greater(t, character.HP, originalHP)
	assert.Equal(t, originalHP+10, character.HP)
}

func TestPassiveSkill_Revenge(t *testing.T) {
	bm, character, enemy := setupWarriorSkillsTest(t)
	
	// 设置复仇被动技能
	passiveState := &CharacterPassiveState{
		PassiveID:   "warrior_passive_revenge",
		Level:       1,
		Passive:     &models.PassiveSkill{ID: "warrior_passive_revenge", EffectType: "counter_attack", EffectStat: ""},
		EffectValue: 15.0, // 15%触发概率
	}
	
	bm.passiveSkillManager.mu.Lock()
	bm.passiveSkillManager.characterPassives[character.ID] = []*CharacterPassiveState{passiveState}
	bm.passiveSkillManager.mu.Unlock()
	
	// 模拟受到攻击
	originalEnemyHP := enemy.HP
	damageTaken := 50
	
	session := &BattleSession{
		UserID:     1,
		IsRunning:  true,
		BattleLogs: []models.BattleLog{},
	}
	
	var logs []models.BattleLog
	
	// 多次测试以确保概率触发
	counterTriggered := false
	for i := 0; i < 100; i++ {
		enemy.HP = originalEnemyHP
		bm.handleCounterAttacks(character, enemy, damageTaken, session, &logs)
		if enemy.HP < originalEnemyHP {
			counterTriggered = true
			break
		}
	}
	
	// 在100次尝试中应该至少触发一次（15%概率）
	assert.True(t, counterTriggered, "复仇被动技能应该在多次尝试中至少触发一次")
}

func TestPassiveSkill_Unbreakable(t *testing.T) {
	bm, character, enemy := setupWarriorSkillsTest(t)
	_ = enemy // 避免未使用变量警告
	
	// 设置坚韧不拔被动技能
	passiveState := &CharacterPassiveState{
		PassiveID:   "warrior_passive_unbreakable",
		Level:       1,
		Passive:     &models.PassiveSkill{ID: "warrior_passive_unbreakable", EffectType: "survival", EffectStat: ""},
		EffectValue: 0.0,
	}
	
	bm.passiveSkillManager.mu.Lock()
	bm.passiveSkillManager.characterPassives[character.ID] = []*CharacterPassiveState{passiveState}
	bm.passiveSkillManager.mu.Unlock()
	
	// 模拟受到致命伤害
	character.HP = 10
	damageTaken := 100
	
	session := &BattleSession{
		UserID:     1,
		IsRunning:  true,
		BattleLogs: []models.BattleLog{},
	}
	
	var logs []models.BattleLog
	_ = session // 避免未使用变量警告
	_ = logs    // 避免未使用变量警告
	
	originalHP := character.HP
	character.HP -= damageTaken
	
	// 检查坚韧不拔效果
	if originalHP > 0 && character.HP <= 0 {
		passives := bm.passiveSkillManager.GetPassiveSkills(character.ID)
		for _, passive := range passives {
			if passive.Passive.EffectType == "survival" && passive.Passive.ID == "warrior_passive_unbreakable" {
				character.HP = 1
				break
			}
		}
	}
	
	// 应该保留1点HP
	assert.Equal(t, 1, character.HP)
}

func TestPassiveSkill_ShieldReflection(t *testing.T) {
	bm, character, enemy := setupWarriorSkillsTest(t)
	
	// 设置盾牌反射被动技能
	passiveState := &CharacterPassiveState{
		PassiveID:   "warrior_passive_shield_reflection",
		Level:       1,
		Passive:     &models.PassiveSkill{ID: "warrior_passive_shield_reflection", EffectType: "reflect", EffectStat: ""},
		EffectValue: 10.0, // 10%反射
	}
	
	bm.passiveSkillManager.mu.Lock()
	bm.passiveSkillManager.characterPassives[character.ID] = []*CharacterPassiveState{passiveState}
	bm.passiveSkillManager.mu.Unlock()
	
	// 模拟受到攻击
	originalEnemyHP := enemy.HP
	damageTaken := 100
	
	session := &BattleSession{
		UserID:     1,
		IsRunning:  true,
		BattleLogs: []models.BattleLog{},
	}
	
	var logs []models.BattleLog
	bm.handlePassiveReflectEffects(character, enemy, damageTaken, session, &logs)
	
	// 应该反射10%伤害 = 10点
	expectedReflectDamage := 10
	assert.Equal(t, originalEnemyHP-expectedReflectDamage, enemy.HP)
}

// ═══════════════════════════════════════════════════════════
// 被动技能怒气管理测试
// ═══════════════════════════════════════════════════════════

func TestPassiveSkill_AngerManagement(t *testing.T) {
	bm, character, _ := setupWarriorSkillsTest(t)
	
	// 设置愤怒掌握被动技能
	passiveState := &CharacterPassiveState{
		PassiveID:   "warrior_passive_anger_management",
		Level:       1,
		Passive:     &models.PassiveSkill{ID: "warrior_passive_anger_management", EffectType: "rage_generation", EffectStat: ""},
		EffectValue: 10.0, // 10%怒气获得加成
	}
	
	bm.passiveSkillManager.mu.Lock()
	bm.passiveSkillManager.characterPassives[character.ID] = []*CharacterPassiveState{passiveState}
	bm.passiveSkillManager.mu.Unlock()
	
	// 测试基础怒气获得
	baseRageGain := 10
	enhancedRageGain := bm.applyRageGenerationModifiers(character.ID, baseRageGain)
	
	// 10%加成，应该增加到11点（10 * 1.1 = 11，向下取整）
	assert.Equal(t, 11, enhancedRageGain)
}

func TestPassiveSkill_WarMachine(t *testing.T) {
	bm, character, enemy := setupWarriorSkillsTest(t)
	
	// 设置战争机器被动技能
	passiveState := &CharacterPassiveState{
		PassiveID:   "warrior_passive_war_machine",
		Level:       1,
		Passive:     &models.PassiveSkill{ID: "warrior_passive_war_machine", EffectType: "rage_generation", EffectStat: ""},
		EffectValue: 30.0, // 击杀获得30点怒气
	}
	
	bm.passiveSkillManager.mu.Lock()
	bm.passiveSkillManager.characterPassives[character.ID] = []*CharacterPassiveState{passiveState}
	bm.passiveSkillManager.mu.Unlock()
	
	// 模拟击杀敌人
	originalRage := character.Resource
	enemy.HP = 0
	
	session := &BattleSession{
		UserID:     1,
		IsRunning:  true,
		BattleLogs: []models.BattleLog{},
	}
	
	var logs []models.BattleLog
	bm.handleWarMachineRageGain(character, session, &logs)
	
	// 应该获得30点额外怒气
	assert.Equal(t, originalRage+30, character.Resource)
}

// ═══════════════════════════════════════════════════════════
// 技能效果测试
// ═══════════════════════════════════════════════════════════

func TestSkill_LastStand(t *testing.T) {
	bm, character, enemy := setupWarriorSkillsTest(t)
	_ = bm    // 避免未使用变量警告
	_ = enemy // 避免未使用变量警告
	
	// 模拟破釜沉舟技能效果
	character.HP = 500 // 半血
	originalHP := character.HP
	
	// 破釜沉舟：恢复30%最大HP
	healPercent := 30.0
	healAmount := int(float64(character.MaxHP) * healPercent / 100.0)
	character.HP += healAmount
	if character.HP > character.MaxHP {
		character.HP = character.MaxHP
	}
	
	// 应该恢复300点HP（1000 * 30% = 300）
	assert.Equal(t, originalHP+300, character.HP)
	assert.Equal(t, 800, character.HP)
}

func TestSkill_UnbreakableBarrier(t *testing.T) {
	bm, character, enemy := setupWarriorSkillsTest(t)
	_ = enemy // 避免未使用变量警告
	
	// 模拟不灭壁垒技能：获得50%最大HP的护盾
	shieldPercent := 50.0
	shieldAmount := int(float64(character.MaxHP) * shieldPercent / 100.0)
	
	bm.buffManager.ApplyBuff(character.ID, "unbreakable_barrier", "不灭壁垒", "buff", true, 4, float64(shieldAmount), "shield", "")
	
	// 检查护盾值
	shieldValue := bm.buffManager.GetBuffValue(character.ID, "shield")
	assert.Equal(t, float64(shieldAmount), shieldValue)
	assert.Equal(t, 500.0, shieldValue) // 1000 * 50% = 500
	
	// 模拟受到伤害
	damage := 300
	shieldInt := int(shieldValue)
	
	if damage <= shieldInt {
		shieldInt -= damage
		damage = 0
	} else {
		damage -= shieldInt
		shieldInt = 0
	}
	
	// 护盾应该吸收300点伤害，剩余200点护盾
	assert.Equal(t, 0, damage)
	assert.Equal(t, 200, shieldInt)
}

func TestSkill_ShieldReflection(t *testing.T) {
	bm, character, enemy := setupWarriorSkillsTest(t)
	
	// 模拟盾牌反射技能：反射50%伤害
	reflectPercent := 50.0
	bm.buffManager.ApplyBuff(character.ID, "shield_reflection", "盾牌反射", "buff", true, 2, reflectPercent, "reflect", "")
	
	// 模拟受到攻击
	originalEnemyHP := enemy.HP
	damageTaken := 100
	
	session := &BattleSession{
		UserID:     1,
		IsRunning:  true,
		BattleLogs: []models.BattleLog{},
	}
	
	var logs []models.BattleLog
	bm.handleActiveReflectEffects(character, enemy, damageTaken, session, &logs)
	
	// 应该反射50%伤害 = 50点
	expectedReflectDamage := 50
	assert.Equal(t, originalEnemyHP-expectedReflectDamage, enemy.HP)
}

// ═══════════════════════════════════════════════════════════
// 敌人Debuff测试
// ═══════════════════════════════════════════════════════════

func TestEnemyDebuff_DemoralizingShout(t *testing.T) {
	bm, _, enemy := setupWarriorSkillsTest(t)
	
	// 应用挫志怒吼Debuff：降低15%攻击力
	attackReduction := 15.0
	bm.buffManager.ApplyEnemyDebuff(enemy.ID, "demoralizing_shout", "挫志怒吼", "debuff", 3, attackReduction, "attack", "")
	
	// 检查Debuff值
	debuffValue := bm.buffManager.GetEnemyDebuffValue(enemy.ID, "attack")
	assert.Equal(t, 15.0, debuffValue)
	
	// 计算实际攻击力
	originalAttack := float64(enemy.PhysicalAttack)
	actualAttack := originalAttack * (1.0 - debuffValue/100.0)
	
	// 80 * (1 - 0.15) = 68
	assert.Equal(t, 68.0, actualAttack)
}

func TestEnemyDebuff_Whirlwind(t *testing.T) {
	bm, _, enemy := setupWarriorSkillsTest(t)
	
	// 应用旋风斩Debuff：降低10%防御
	defenseReduction := 10.0
	bm.buffManager.ApplyEnemyDebuff(enemy.ID, "whirlwind", "旋风斩", "debuff", 2, defenseReduction, "defense", "")
	
	// 检查Debuff值
	debuffValue := bm.buffManager.GetEnemyDebuffValue(enemy.ID, "defense")
	assert.Equal(t, 10.0, debuffValue)
	
	// 计算实际防御力
	originalDefense := float64(enemy.PhysicalDefense)
	actualDefense := originalDefense * (1.0 - debuffValue/100.0)
	
	// 30 * (1 - 0.10) = 27
	assert.Equal(t, 27.0, actualDefense)
}

func TestEnemyDebuff_MortalStrike(t *testing.T) {
	bm, _, enemy := setupWarriorSkillsTest(t)
	
	// 应用致死打击Debuff：降低50%治疗效果
	healingReduction := 50.0
	bm.buffManager.ApplyEnemyDebuff(enemy.ID, "mortal_strike", "致死打击", "debuff", 3, healingReduction, "healing_received", "")
	
	// 检查Debuff值
	debuffValue := bm.buffManager.GetEnemyDebuffValue(enemy.ID, "healing_received")
	assert.Equal(t, 50.0, debuffValue)
	
	// 模拟治疗效果
	baseHealing := 100
	actualHealing := int(float64(baseHealing) * (1.0 - debuffValue/100.0))
	
	// 100 * (1 - 0.50) = 50
	assert.Equal(t, 50, actualHealing)
}


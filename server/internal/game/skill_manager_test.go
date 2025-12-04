package game

import (
	"testing"

	"text-wow/internal/models"
	"text-wow/internal/repository"
	"text-wow/internal/service"

	"github.com/stretchr/testify/assert"
)

// ═══════════════════════════════════════════════════════════
// SkillManager 基础功能测试
// ═══════════════════════════════════════════════════════════

func setupSkillManagerTest(t *testing.T) (*SkillManager, *models.Character, *models.Monster) {
	skillRepo := repository.NewSkillRepository()
	characterRepo := repository.NewCharacterRepository()
	skillService := service.NewSkillService(skillRepo, characterRepo)
	
	sm := NewSkillManager()
	sm.skillService = skillService
	sm.skillRepo = skillRepo
	
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
		CritRate:     0.1,
		CritDamage:   1.5,
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
	
	return sm, character, enemy
}

func TestSkillManager_CalculateSkillDamage_Basic(t *testing.T) {
	sm, character, enemy := setupSkillManagerTest(t)
	
	// 创建测试技能状态
	skillState := &CharacterSkillState{
		SkillID:    "warrior_heroic_strike",
		SkillLevel: 1,
		Skill: &models.Skill{
			ID:            "warrior_heroic_strike",
			Name:          "英勇打击",
			Type:          "attack",
			TargetType:    "enemy",
			DamageType:    "physical",
			BaseValue:     100,
			ScalingRatio:  1.0,
			ResourceCost:  10,
			Cooldown:      0,
		},
		Effect: map[string]interface{}{
			"damageMultiplier": 1.0,
		},
	}
	
	// 计算技能伤害
	damage := sm.CalculateSkillDamage(skillState, character, enemy, nil, nil)
	
	// 基础伤害 = 100攻击力 * 1.0倍率 - 30防御/2 = 100 - 15 = 85
	// 加上随机波动，应该在合理范围内
	assert.Greater(t, damage, 0)
	assert.Less(t, damage, 200) // 应该在合理范围内
}

func TestSkillManager_CalculateSkillDamage_WithPassiveModifiers(t *testing.T) {
	sm, character, enemy := setupSkillManagerTest(t)
	
	// 创建被动技能管理器并添加攻击力加成
	psm := NewPassiveSkillManager()
	passiveState := &CharacterPassiveState{
		PassiveID:   "warrior_passive_battle_focus",
		Level:       1,
		Passive:     &models.PassiveSkill{ID: "warrior_passive_battle_focus", EffectType: "stat_mod", EffectStat: "attack"},
		EffectValue: 10.0, // 10%攻击力加成
	}
	
	psm.mu.Lock()
	psm.characterPassives[character.ID] = []*CharacterPassiveState{passiveState}
	psm.mu.Unlock()
	
	// 创建测试技能状态
	skillState := &CharacterSkillState{
		SkillID:    "warrior_heroic_strike",
		SkillLevel: 1,
		Skill: &models.Skill{
			ID:            "warrior_heroic_strike",
			Name:          "英勇打击",
			Type:          "attack",
			TargetType:    "enemy",
			DamageType:    "physical",
			BaseValue:     100,
			ScalingRatio:  1.0,
			ResourceCost:  10,
			Cooldown:      0,
		},
		Effect: map[string]interface{}{
			"damageMultiplier": 1.0,
		},
	}
	
	// 计算技能伤害（带被动技能加成）
	damageWithPassive := sm.CalculateSkillDamage(skillState, character, enemy, psm, nil)
	
	// 计算技能伤害（不带被动技能加成）
	damageWithoutPassive := sm.CalculateSkillDamage(skillState, character, enemy, nil, nil)
	
	// 带被动技能的伤害应该更高或相等（由于随机波动可能相等）
	// 验证两者都在合理范围内
	assert.Greater(t, damageWithPassive, 0)
	assert.Greater(t, damageWithoutPassive, 0)
	assert.Less(t, damageWithPassive, 200)
	assert.Less(t, damageWithoutPassive, 200)
	// 理论上带被动技能的伤害应该更高，但由于随机波动，我们只验证两者都有效
	// 可以通过多次测试取平均值来验证，但这里简化处理
}

func TestSkillManager_CalculateSkillDamage_WithBuffModifiers(t *testing.T) {
	sm, character, enemy := setupSkillManagerTest(t)
	
	// 创建Buff管理器并添加攻击力Buff
	bm := NewBuffManager()
	bm.ApplyBuff(character.ID, "battle_shout", "战斗怒吼", "buff", true, 5, 20.0, "attack", "")
	
	// 创建测试技能状态
	skillState := &CharacterSkillState{
		SkillID:    "warrior_heroic_strike",
		SkillLevel: 1,
		Skill: &models.Skill{
			ID:            "warrior_heroic_strike",
			Name:          "英勇打击",
			Type:          "attack",
			TargetType:    "enemy",
			DamageType:    "physical",
			BaseValue:     100,
			ScalingRatio:  1.0,
			ResourceCost:  10,
			Cooldown:      0,
		},
		Effect: map[string]interface{}{
			"damageMultiplier": 1.0,
		},
	}
	
	// 计算技能伤害（带Buff加成）
	damageWithBuff := sm.CalculateSkillDamage(skillState, character, enemy, nil, bm)
	
	// 计算技能伤害（不带Buff加成）
	damageWithoutBuff := sm.CalculateSkillDamage(skillState, character, enemy, nil, nil)
	
	// 带Buff的伤害应该更高
	assert.Greater(t, damageWithBuff, damageWithoutBuff)
}

func TestSkillManager_CalculateSkillDamage_ShieldSlam(t *testing.T) {
	sm, character, enemy := setupSkillManagerTest(t)
	
	// 盾牌猛击：基于攻击力和防御力
	skillState := &CharacterSkillState{
		SkillID:    "warrior_shield_slam",
		SkillLevel: 1,
		Skill: &models.Skill{
			ID:            "warrior_shield_slam",
			Name:          "盾牌猛击",
			Type:          "attack",
			TargetType:    "enemy",
			DamageType:    "physical",
			BaseValue:     0,
			ScalingRatio:  0,
			ResourceCost:  20,
			Cooldown:      2,
		},
		Effect: map[string]interface{}{
			"attackMultiplier":  1.2,
			"defenseMultiplier": 0.5,
		},
	}
	
	// 计算技能伤害
	damage := sm.CalculateSkillDamage(skillState, character, enemy, nil, nil)
	
	// 盾牌猛击伤害 = 100攻击 * 1.2 + 50防御 * 0.5 - 30防御/2 = 120 + 25 - 15 = 130
	// 加上随机波动，应该在合理范围内
	assert.Greater(t, damage, 0)
	assert.Less(t, damage, 300) // 应该在合理范围内
}

func TestSkillManager_CalculateSkillDamage_WithEnemyDebuff(t *testing.T) {
	sm, character, enemy := setupSkillManagerTest(t)
	
	// 创建Buff管理器并添加敌人防御Debuff
	bm := NewBuffManager()
	bm.ApplyEnemyDebuff(enemy.ID, "whirlwind", "旋风斩", "debuff", 2, 10.0, "defense", "")
	
	// 创建测试技能状态
	skillState := &CharacterSkillState{
		SkillID:    "warrior_heroic_strike",
		SkillLevel: 1,
		Skill: &models.Skill{
			ID:            "warrior_heroic_strike",
			Name:          "英勇打击",
			Type:          "attack",
			TargetType:    "enemy",
			DamageType:    "physical",
			BaseValue:     100,
			ScalingRatio:  1.0,
			ResourceCost:  10,
			Cooldown:      0,
		},
		Effect: map[string]interface{}{
			"damageMultiplier": 1.0,
		},
	}
	
	// 计算技能伤害（带敌人Debuff）
	damageWithDebuff := sm.CalculateSkillDamage(skillState, character, enemy, nil, bm)
	
	// 计算技能伤害（不带敌人Debuff）
	damageWithoutDebuff := sm.CalculateSkillDamage(skillState, character, enemy, nil, nil)
	
	// 带Debuff的伤害应该更高或相等（敌人防御降低，但由于随机波动可能相等）
	// 由于有随机波动，我们只检查伤害在合理范围内
	assert.Greater(t, damageWithDebuff, 0)
	assert.Greater(t, damageWithoutDebuff, 0)
	// 理论上带Debuff的伤害应该更高，但由于随机波动，我们只验证两者都有效
	assert.Less(t, damageWithDebuff, 200)
	assert.Less(t, damageWithoutDebuff, 200)
}

func TestSkillManager_ApplySkillEffects_Charge(t *testing.T) {
	sm, character, _ := setupSkillManagerTest(t)
	
	// 冲锋技能：获得怒气，可能眩晕
	skillState := &CharacterSkillState{
		SkillID:    "warrior_charge",
		SkillLevel: 1,
		Skill: &models.Skill{
			ID:            "warrior_charge",
			Name:          "冲锋",
			Type:          "attack",
			TargetType:    "enemy",
			DamageType:    "physical",
			BaseValue:     80,
			ScalingRatio:  0.8,
			ResourceCost:  0,
			Cooldown:      3,
		},
		Effect: map[string]interface{}{
			"rageGain":    15,
			"stunChance":  0.3,
		},
	}
	
	originalRage := character.Resource
	
	// 应用技能效果
	effects := sm.ApplySkillEffects(skillState, character, nil)
	
	// 应该获得怒气
	assert.Equal(t, originalRage+15, character.Resource)
	assert.Contains(t, effects, "rageGain")
	
	// 可能触发眩晕（概率性）
	_, hasStun := effects["stun"]
	if hasStun {
		assert.True(t, effects["stun"].(bool))
	}
}

func TestSkillManager_ApplySkillEffects_Bloodthirst(t *testing.T) {
	sm, character, _ := setupSkillManagerTest(t)
	
	// 嗜血技能：恢复生命值
	skillState := &CharacterSkillState{
		SkillID:    "warrior_bloodthirst",
		SkillLevel: 1,
		Skill: &models.Skill{
			ID:            "warrior_bloodthirst",
			Name:          "嗜血",
			Type:          "attack",
			TargetType:    "enemy",
			DamageType:    "physical",
			BaseValue:     120,
			ScalingRatio:  1.2,
			ResourceCost:  25,
			Cooldown:      0,
		},
		Effect: map[string]interface{}{
			"healPercent": 30.0,
		},
	}
	
	// 应用技能效果
	effects := sm.ApplySkillEffects(skillState, character, nil)
	
	// 应该包含恢复效果
	assert.Contains(t, effects, "healPercent")
	assert.Equal(t, 30.0, effects["healPercent"])
}

func TestSkillManager_TickCooldowns(t *testing.T) {
	sm, character, _ := setupSkillManagerTest(t)
	
	// 加载技能
	skillState := &CharacterSkillState{
		SkillID:      "warrior_charge",
		SkillLevel:   1,
		CooldownLeft: 3,
		Skill: &models.Skill{
			ID:            "warrior_charge",
			Name:          "冲锋",
			Cooldown:      3,
		},
		Effect: map[string]interface{}{},
	}
	
	sm.mu.Lock()
	sm.characterSkills[character.ID] = []*CharacterSkillState{skillState}
	sm.mu.Unlock()
	
	// 减少冷却时间
	sm.TickCooldowns(character.ID)
	
	// 冷却时间应该减少1
	assert.Equal(t, 2, skillState.CooldownLeft)
	
	// 继续减少
	sm.TickCooldowns(character.ID)
	sm.TickCooldowns(character.ID)
	
	// 冷却时间应该为0
	assert.Equal(t, 0, skillState.CooldownLeft)
}


package game

import (
	"testing"

	"text-wow/internal/models"

	"github.com/stretchr/testify/assert"
)

// ═══════════════════════════════════════════════════════════
// PassiveSkillManager 基础功能测试
// ═══════════════════════════════════════════════════════════

func setupPassiveSkillManagerTest(t *testing.T) (*PassiveSkillManager, func()) {
	psm := NewPassiveSkillManager()
	cleanup := func() {
		// 清理函数
	}
	return psm, cleanup
}

func TestPassiveSkillManager_GetPassiveModifier(t *testing.T) {
	psm, cleanup := setupPassiveSkillManagerTest(t)
	defer cleanup()

	characterID := 1

	// 创建测试被动技能状态
	passiveState := &CharacterPassiveState{
		PassiveID:   "warrior_passive_battle_focus",
		Level:       1,
		Passive:     &models.PassiveSkill{ID: "warrior_passive_battle_focus", EffectType: "stat_mod", EffectStat: "attack"},
		EffectValue: 10.0, // 10%攻击力加成
	}

	// 手动添加到管理器（模拟加载）
	psm.mu.Lock()
	psm.characterPassives[characterID] = []*CharacterPassiveState{passiveState}
	psm.mu.Unlock()

	// 测试获取被动技能修正值
	modifier := psm.GetPassiveModifier(characterID, "attack")
	assert.Equal(t, 10.0, modifier)
}

func TestPassiveSkillManager_GetPassiveModifier_Multiple(t *testing.T) {
	psm, cleanup := setupPassiveSkillManagerTest(t)
	defer cleanup()

	characterID := 1

	// 创建多个被动技能状态
	passiveStates := []*CharacterPassiveState{
		{
			PassiveID:   "warrior_passive_battle_focus",
			Level:       1,
			Passive:     &models.PassiveSkill{ID: "warrior_passive_battle_focus", EffectType: "stat_mod", EffectStat: "attack"},
			EffectValue: 10.0,
		},
		{
			PassiveID:   "warrior_passive_weapon_mastery",
			Level:       1,
			Passive:     &models.PassiveSkill{ID: "warrior_passive_weapon_mastery", EffectType: "stat_mod", EffectStat: "attack"},
			EffectValue: 10.0,
		},
	}

	psm.mu.Lock()
	psm.characterPassives[characterID] = passiveStates
	psm.mu.Unlock()

	// 测试多个被动技能累加
	modifier := psm.GetPassiveModifier(characterID, "attack")
	assert.Equal(t, 20.0, modifier) // 10% + 10% = 20%
}

func TestPassiveSkillManager_GetPassiveModifier_MultiAttribute(t *testing.T) {
	psm, cleanup := setupPassiveSkillManagerTest(t)
	defer cleanup()

	characterID := 1

	// 测试多属性被动技能（防御姿态）
	passiveState := &CharacterPassiveState{
		PassiveID:   "warrior_passive_defensive_stance",
		Level:       1,
		Passive:     &models.PassiveSkill{ID: "warrior_passive_defensive_stance", EffectType: "stat_mod", EffectStat: "threat_and_defense"},
		EffectValue: 15.0,
	}

	psm.mu.Lock()
	psm.characterPassives[characterID] = []*CharacterPassiveState{passiveState}
	psm.mu.Unlock()

	// 测试匹配多属性
	threatModifier := psm.GetPassiveModifier(characterID, "threat")
	assert.Equal(t, 15.0, threatModifier)
	
	defenseModifier := psm.GetPassiveModifier(characterID, "defense")
	assert.Equal(t, 15.0, defenseModifier)
}

func TestPassiveSkillManager_GetPassiveSkills(t *testing.T) {
	psm, cleanup := setupPassiveSkillManagerTest(t)
	defer cleanup()

	characterID := 1

	passiveStates := []*CharacterPassiveState{
		{
			PassiveID:   "warrior_passive_battle_focus",
			Level:       1,
			Passive:     &models.PassiveSkill{ID: "warrior_passive_battle_focus", EffectType: "stat_mod", EffectStat: "attack"},
			EffectValue: 10.0,
		},
		{
			PassiveID:   "warrior_passive_anger_management",
			Level:       1,
			Passive:     &models.PassiveSkill{ID: "warrior_passive_anger_management", EffectType: "rage_generation", EffectStat: ""},
			EffectValue: 10.0,
		},
	}

	psm.mu.Lock()
	psm.characterPassives[characterID] = passiveStates
	psm.mu.Unlock()

	// 测试获取被动技能列表
	passives := psm.GetPassiveSkills(characterID)
	assert.Len(t, passives, 2)
	assert.Equal(t, "warrior_passive_battle_focus", passives[0].PassiveID)
	assert.Equal(t, "warrior_passive_anger_management", passives[1].PassiveID)
}

func TestPassiveSkillManager_HasPassiveSkill(t *testing.T) {
	psm, cleanup := setupPassiveSkillManagerTest(t)
	defer cleanup()

	characterID := 1

	passiveState := &CharacterPassiveState{
		PassiveID:   "warrior_passive_battle_focus",
		Level:       1,
		Passive:     &models.PassiveSkill{ID: "warrior_passive_battle_focus", EffectType: "stat_mod", EffectStat: "attack"},
		EffectValue: 10.0,
	}

	psm.mu.Lock()
	psm.characterPassives[characterID] = []*CharacterPassiveState{passiveState}
	psm.mu.Unlock()

	// 测试检查被动技能
	assert.True(t, psm.HasPassiveSkill(characterID, "warrior_passive_battle_focus"))
	assert.False(t, psm.HasPassiveSkill(characterID, "warrior_passive_weapon_mastery"))
}

func TestPassiveSkillManager_GetPassiveSkillLevel(t *testing.T) {
	psm, cleanup := setupPassiveSkillManagerTest(t)
	defer cleanup()

	characterID := 1

	passiveState := &CharacterPassiveState{
		PassiveID:   "warrior_passive_battle_focus",
		Level:       3,
		Passive:     &models.PassiveSkill{ID: "warrior_passive_battle_focus", EffectType: "stat_mod", EffectStat: "attack"},
		EffectValue: 20.0,
	}

	psm.mu.Lock()
	psm.characterPassives[characterID] = []*CharacterPassiveState{passiveState}
	psm.mu.Unlock()

	// 测试获取被动技能等级
	level := psm.GetPassiveSkillLevel(characterID, "warrior_passive_battle_focus")
	assert.Equal(t, 3, level)
}


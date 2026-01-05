package game

import (
	"testing"

	"text-wow/internal/database"
	"text-wow/internal/models"
	"text-wow/internal/repository"

	"github.com/stretchr/testify/assert"
)

// ═══════════════════════════════════════════════════════════
// 被动技能触发条件测试
// ═══════════════════════════════════════════════════════════

func setupPassiveTriggerTest(t *testing.T) (*BattleManager, *models.Character, func()) {
	testDB, err := database.SetupTestDB()
	if err != nil {
		t.Fatalf("Failed to setup test database: %v", err)
	}

	manager := &BattleManager{
		sessions:            make(map[int]*BattleSession),
		gameRepo:            repository.NewGameRepository(),
		charRepo:            repository.NewCharacterRepository(),
		skillManager:        NewSkillManager(),
		buffManager:         NewBuffManager(),
		passiveSkillManager: NewPassiveSkillManager(),
		calculator:          NewCalculator(),
		monsterManager:      NewMonsterManager(),
	}

	char := &models.Character{
		ID:           1,
		UserID:       1,
		Name:         "测试角色",
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
	}

	cleanup := func() {
		database.TeardownTestDB(testDB)
	}

	return manager, char, cleanup
}

func TestPassiveTrigger_OnHit(t *testing.T) {
	manager, char, cleanup := setupPassiveTriggerTest(t)
	defer cleanup()

	// 添加on_hit_heal被动技能
	passiveSkill := &models.PassiveSkill{
		ID:          "test_on_hit_heal",
		Name:        "血之狂热",
		EffectType:  "on_hit_heal",
		EffectValue: 1.0, // 1%最大HP
		MaxLevel:    5,
		LevelScaling: 0.2,
	}

	passiveState := &CharacterPassiveState{
		PassiveID:   "test_on_hit_heal",
		Level:       1,
		Passive:     passiveSkill,
		EffectValue: 1.0, // 1%最大HP
	}
	
	// 确保Passive不为nil
	assert.NotNil(t, passiveState.Passive, "Passive不应该为nil")
	assert.Equal(t, "on_hit_heal", passiveState.Passive.EffectType, "EffectType应该匹配")

	// 手动添加到被动技能管理器（直接访问内部字段，测试用）
	manager.passiveSkillManager.mu.Lock()
	if manager.passiveSkillManager.characterPassives == nil {
		manager.passiveSkillManager.characterPassives = make(map[int][]*CharacterPassiveState)
	}
	manager.passiveSkillManager.characterPassives[char.ID] = []*CharacterPassiveState{passiveState}
	manager.passiveSkillManager.mu.Unlock()

	// 验证被动技能已添加
	passives := manager.passiveSkillManager.GetPassiveSkills(char.ID)
	assert.NotEmpty(t, passives, "被动技能应该已添加")
	if len(passives) == 0 {
		t.Fatalf("被动技能列表为空，无法继续测试")
	}
	assert.NotNil(t, passives[0].Passive, "Passive不应该为nil")
	if passives[0].Passive == nil {
		t.Fatalf("Passive为nil，无法继续测试")
	}
	assert.Equal(t, "on_hit_heal", passives[0].Passive.EffectType, "被动技能类型应该匹配")
	
	// 调试信息
	t.Logf("被动技能验证：passives数量=%d, EffectType=%s, EffectValue=%.2f", 
		len(passives), passives[0].Passive.EffectType, passives[0].EffectValue)

	// 创建战斗会话
	session := &BattleSession{
		UserID:     1,
		BattleLogs: make([]models.BattleLog, 0),
	}
	manager.sessions[1] = session

	logs := make([]models.BattleLog, 0)
	// 先将HP降低，以便能看到治疗效果
	char.HP = 900 // 降低到900，这样治疗10点后可以观察到
	originalHP := char.HP

	// 触发on_hit效果
	manager.handlePassiveOnHitEffects(char, 100, true, session, &logs)

	// 应该恢复1%最大HP = 10点
	expectedHealAmount := int(float64(char.MaxHP) * passiveState.EffectValue / 100.0)
	t.Logf("触发后：originalHP=%d, currentHP=%d, expectedHeal=%d, EffectValue=%.2f, logs数量=%d", 
		originalHP, char.HP, expectedHealAmount, passiveState.EffectValue, len(logs))
	
	if expectedHealAmount > 0 {
		if char.HP > originalHP {
			// HP增加了，说明被动技能触发了
			actualHealAmount := char.HP - originalHP
			assert.Equal(t, expectedHealAmount, actualHealAmount, "HP应该增加%d点", expectedHealAmount)
			assert.NotEmpty(t, logs, "应该有日志记录")
		} else {
			// HP没有增加，说明被动技能没有触发
			// 再次检查被动技能
			checkPassives := manager.passiveSkillManager.GetPassiveSkills(char.ID)
			t.Errorf("被动技能未触发：HP没有增加。passives数量=%d, EffectType=%s, EffectValue=%.2f, healAmount计算=%d", 
				len(checkPassives), 
				func() string {
					if len(checkPassives) > 0 && checkPassives[0].Passive != nil {
						return checkPassives[0].Passive.EffectType
					}
					return "nil"
				}(),
				passiveState.EffectValue,
				expectedHealAmount)
		}
	} else {
		t.Errorf("expectedHealAmount为0，无法触发治疗效果。MaxHP=%d, EffectValue=%.2f", 
			char.MaxHP, passiveState.EffectValue)
	}
}

func TestPassiveTrigger_OnCrit(t *testing.T) {
	manager, char, cleanup := setupPassiveTriggerTest(t)
	defer cleanup()

	// 添加on_crit_heal被动技能
	passiveSkill := &models.PassiveSkill{
		ID:          "test_on_crit_heal",
		Name:        "暴击恢复",
		EffectType:  "on_crit_heal",
		EffectValue: 5.0, // 5%暴击伤害
		MaxLevel:    5,
		LevelScaling: 1.0,
	}

	passiveState := &CharacterPassiveState{
		PassiveID:   "test_on_crit_heal",
		Level:       1,
		Passive:     passiveSkill,
		EffectValue: 5.0,
	}

	manager.passiveSkillManager.mu.Lock()
	if manager.passiveSkillManager.characterPassives == nil {
		manager.passiveSkillManager.characterPassives = make(map[int][]*CharacterPassiveState)
	}
	manager.passiveSkillManager.characterPassives[char.ID] = []*CharacterPassiveState{passiveState}
	manager.passiveSkillManager.mu.Unlock()

	session := &BattleSession{
		UserID:     1,
		BattleLogs: make([]models.BattleLog, 0),
	}
	manager.sessions[1] = session

	logs := make([]models.BattleLog, 0)
	// 先将HP降低，以便能看到治疗效果
	char.HP = 900 // 降低到900，这样治疗10点后可以观察到
	originalHP := char.HP
	critDamage := 200

	// 触发on_crit效果
	manager.handlePassiveOnCritEffects(char, critDamage, true, session, &logs)

	// 应该恢复5%暴击伤害 = 10点
	assert.Greater(t, char.HP, originalHP)
	assert.Equal(t, originalHP+10, char.HP)
	assert.NotEmpty(t, logs)
}

func TestPassiveTrigger_OnKill(t *testing.T) {
	manager, char, cleanup := setupPassiveTriggerTest(t)
	defer cleanup()
	
	enemy := &models.Monster{
		ID:   "test_enemy",
		Name: "测试敌人",
		HP:   0,
		MaxHP: 100,
	}

	// 添加on_kill_heal被动技能
	passiveSkill := &models.PassiveSkill{
		ID:          "test_on_kill_heal",
		Name:        "击杀恢复",
		EffectType:  "on_kill_heal",
		EffectValue: 10.0, // 10%最大HP
		MaxLevel:    5,
		LevelScaling: 2.0,
	}

	passiveState := &CharacterPassiveState{
		PassiveID:   "test_on_kill_heal",
		Level:       1,
		Passive:     passiveSkill,
		EffectValue: 10.0,
	}

	manager.passiveSkillManager.mu.Lock()
	if manager.passiveSkillManager.characterPassives == nil {
		manager.passiveSkillManager.characterPassives = make(map[int][]*CharacterPassiveState)
	}
	manager.passiveSkillManager.characterPassives[char.ID] = []*CharacterPassiveState{passiveState}
	manager.passiveSkillManager.mu.Unlock()

	session := &BattleSession{
		UserID:     1,
		BattleLogs: make([]models.BattleLog, 0),
	}
	manager.sessions[1] = session

	logs := make([]models.BattleLog, 0)
	// 先将HP降低，以便能看到治疗效果
	char.HP = 800 // 降低到800，这样治疗100点后可以观察到
	originalHP := char.HP

	// 触发on_kill效果
	manager.handlePassiveOnKillEffects(char, enemy, session, &logs)

	// 应该恢复10%最大HP = 100点
	assert.Greater(t, char.HP, originalHP)
	assert.Equal(t, originalHP+100, char.HP)
	assert.NotEmpty(t, logs)
}

func TestPassiveTrigger_OnDamageTaken(t *testing.T) {
	manager, char, cleanup := setupPassiveTriggerTest(t)
	defer cleanup()

	// 添加on_damage_taken_resource被动技能
	passiveSkill := &models.PassiveSkill{
		ID:          "test_on_damage_taken_resource",
		Name:        "受击回怒",
		EffectType:  "on_damage_taken_resource",
		EffectValue: 1.0, // 1%伤害转换为怒气
		MaxLevel:    5,
		LevelScaling: 0.2,
	}

	passiveState := &CharacterPassiveState{
		PassiveID:   "test_on_damage_taken_resource",
		Level:       1,
		Passive:     passiveSkill,
		EffectValue: 1.0,
	}

	manager.passiveSkillManager.mu.Lock()
	if manager.passiveSkillManager.characterPassives == nil {
		manager.passiveSkillManager.characterPassives = make(map[int][]*CharacterPassiveState)
	}
	manager.passiveSkillManager.characterPassives[char.ID] = []*CharacterPassiveState{passiveState}
	manager.passiveSkillManager.mu.Unlock()

	session := &BattleSession{
		UserID:     1,
		BattleLogs: make([]models.BattleLog, 0),
	}
	manager.sessions[1] = session

	logs := make([]models.BattleLog, 0)
	originalResource := char.Resource
	damageTaken := 100

	// 触发on_damage_taken效果
	manager.handlePassiveOnDamageTakenEffects(char, damageTaken, session, &logs)

	// 应该获得1%伤害 = 1点怒气
	assert.Greater(t, char.Resource, originalResource)
	assert.Equal(t, originalResource+1, char.Resource)
	assert.NotEmpty(t, logs)
}

func TestPassiveTrigger_OnSkillUse(t *testing.T) {
	manager, char, cleanup := setupPassiveTriggerTest(t)
	defer cleanup()

	// 添加on_skill_use_resource被动技能
	passiveSkill := &models.PassiveSkill{
		ID:          "test_on_skill_use_resource",
		Name:        "技能回怒",
		EffectType:  "on_skill_use_resource",
		EffectValue: 5.0, // 5点怒气
		MaxLevel:    5,
		LevelScaling: 1.0,
	}

	passiveState := &CharacterPassiveState{
		PassiveID:   "test_on_skill_use_resource",
		Level:       1,
		Passive:     passiveSkill,
		EffectValue: 5.0,
	}

	manager.passiveSkillManager.mu.Lock()
	if manager.passiveSkillManager.characterPassives == nil {
		manager.passiveSkillManager.characterPassives = make(map[int][]*CharacterPassiveState)
	}
	manager.passiveSkillManager.characterPassives[char.ID] = []*CharacterPassiveState{passiveState}
	manager.passiveSkillManager.mu.Unlock()

	session := &BattleSession{
		UserID:     1,
		BattleLogs: make([]models.BattleLog, 0),
	}
	manager.sessions[1] = session

	logs := make([]models.BattleLog, 0)
	originalResource := char.Resource

	// 触发on_skill_use效果
	manager.handlePassiveOnSkillUseEffects(char, "test_skill", session, &logs)

	// 应该获得5点怒气
	assert.Greater(t, char.Resource, originalResource)
	assert.Equal(t, originalResource+5, char.Resource)
	assert.NotEmpty(t, logs)
}

func TestPassiveTrigger_MultipleTriggers(t *testing.T) {
	manager, char, cleanup := setupPassiveTriggerTest(t)
	defer cleanup()

	// 添加多个被动技能
	passiveStates := []*CharacterPassiveState{
		{
			PassiveID:   "test_on_hit_heal",
			Level:       1,
			Passive:     &models.PassiveSkill{ID: "test_on_hit_heal", EffectType: "on_hit_heal", EffectValue: 1.0},
			EffectValue: 1.0,
		},
		{
			PassiveID:   "test_on_hit_resource",
			Level:       1,
			Passive:     &models.PassiveSkill{ID: "test_on_hit_resource", EffectType: "on_hit_resource", EffectValue: 2.0},
			EffectValue: 2.0,
		},
	}

	manager.passiveSkillManager.mu.Lock()
	if manager.passiveSkillManager.characterPassives == nil {
		manager.passiveSkillManager.characterPassives = make(map[int][]*CharacterPassiveState)
	}
	manager.passiveSkillManager.characterPassives[char.ID] = passiveStates
	manager.passiveSkillManager.mu.Unlock()

	session := &BattleSession{
		UserID:     1,
		BattleLogs: make([]models.BattleLog, 0),
	}
	manager.sessions[1] = session

	logs := make([]models.BattleLog, 0)
	// 先将HP降低，以便能看到治疗效果
	char.HP = 900 // 降低到900，这样治疗10点后可以观察到
	originalHP := char.HP
	originalResource := char.Resource

	// 触发on_hit效果（应该同时触发两个被动技能）
	manager.handlePassiveOnHitEffects(char, 100, true, session, &logs)

	// 应该同时恢复生命值和获得资源
	assert.Greater(t, char.HP, originalHP)
	assert.Greater(t, char.Resource, originalResource)
	assert.Equal(t, 2, len(logs)) // 应该产生2条日志
}


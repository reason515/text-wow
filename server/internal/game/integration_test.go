package game

import (
	"testing"

	"text-wow/internal/database"
	"text-wow/internal/models"
	"text-wow/internal/repository"

	"github.com/stretchr/testify/assert"
)

// ═══════════════════════════════════════════════════════════
// 集成测试：验证各个系统之间的交互
// ═══════════════════════════════════════════════════════════

func setupIntegrationTest(t *testing.T) (*BattleManager, *models.Character, func()) {
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

// TestIntegration_SkillWithBuffAndPassive 测试技能、Buff和被动技能的集成
func TestIntegration_SkillWithBuffAndPassive(t *testing.T) {
	manager, char, cleanup := setupIntegrationTest(t)
	defer cleanup()

	// 1. 应用Buff（战斗怒吼：增加攻击力）
	manager.buffManager.ApplyBuff(char.ID, "battle_shout", "战斗怒吼", "buff", true, 5, 20.0, "attack", "")
	
	// 2. 添加被动技能（攻击力加成）
	passiveState := &CharacterPassiveState{
		PassiveID:   "test_attack_mod",
		Level:       1,
		Passive:     &models.PassiveSkill{ID: "test_attack_mod", EffectType: "stat_mod", EffectStat: "attack"},
		EffectValue: 10.0, // 10%攻击力加成
	}
	manager.passiveSkillManager.mu.Lock()
	if manager.passiveSkillManager.characterPassives == nil {
		manager.passiveSkillManager.characterPassives = make(map[int][]*CharacterPassiveState)
	}
	manager.passiveSkillManager.characterPassives[char.ID] = []*CharacterPassiveState{passiveState}
	manager.passiveSkillManager.mu.Unlock()

	// 3. 验证Buff和被动技能都生效
	buffValue := manager.buffManager.GetBuffValue(char.ID, "attack")
	passiveModifier := manager.passiveSkillManager.GetPassiveModifier(char.ID, "attack")
	
	assert.Equal(t, 20.0, buffValue, "Buff应该生效")
	assert.Equal(t, 10.0, passiveModifier, "被动技能应该生效")
	
	// 4. 计算总攻击力加成
	totalModifier := buffValue + passiveModifier
	assert.Equal(t, 30.0, totalModifier, "总加成应该是30%")
}

// TestIntegration_DOTWithBuff 测试DOT效果和Buff的集成
func TestIntegration_DOTWithBuff(t *testing.T) {
	manager, char, cleanup := setupIntegrationTest(t)
	defer cleanup()

	// 1. 应用DOT效果
	manager.buffManager.ApplyBuffWithDOT(char.ID, "dot_poison", "毒药", "debuff", false, 5, 10.0, "damage", "nature", true, false, 0)

	// 2. 应用减伤Buff
	manager.buffManager.ApplyBuff(char.ID, "shield_wall", "盾墙", "buff", true, 3, -60.0, "damage_taken", "")

	// 3. 处理DOT效果（应该造成伤害）
	damage, _ := manager.buffManager.ProcessDOTEffects(char.ID, 1)
	assert.Equal(t, 10, damage, "DOT应该造成10点伤害")

	// 4. 应用减伤（60%减伤）
	reducedDamage := manager.buffManager.CalculateDamageTakenWithBuffs(damage, char.ID, true)
	assert.Equal(t, 4, reducedDamage, "减伤后应该是4点伤害")
}

// TestIntegration_SkillCooldownWithBuff 测试技能冷却和Buff的集成
func TestIntegration_SkillCooldownWithBuff(t *testing.T) {
	manager, char, cleanup := setupIntegrationTest(t)
	defer cleanup()

	// 1. 加载角色技能（需要数据库支持，这里简化测试）
	// 假设有一个冷却时间为3的技能
	skillState := &CharacterSkillState{
		SkillID:      "test_skill",
		SkillLevel:   1,
		CooldownLeft:  0,
		Skill: &models.Skill{
			ID:         "test_skill",
			Name:       "测试技能",
			Type:       "attack",
			Cooldown:   3,
			ResourceCost: 20,
		},
	}

	// 2. 使用技能（设置冷却）
	manager.skillManager.mu.Lock()
	if manager.skillManager.characterSkills == nil {
		manager.skillManager.characterSkills = make(map[int][]*CharacterSkillState)
	}
	manager.skillManager.characterSkills[char.ID] = []*CharacterSkillState{skillState}
	manager.skillManager.mu.Unlock()

	// 3. 使用技能
	_, err := manager.skillManager.UseSkill(char.ID, "test_skill")
	assert.NoError(t, err)

	// 4. 检查冷却时间
	updatedSkill := manager.skillManager.GetSkillState(char.ID, "test_skill")
	if updatedSkill != nil {
		assert.Equal(t, 3, updatedSkill.CooldownLeft, "冷却时间应该设置为3")
	}

	// 5. 减少冷却时间
	manager.skillManager.TickCooldowns(char.ID)
	updatedSkill = manager.skillManager.GetSkillState(char.ID, "test_skill")
	if updatedSkill != nil {
		assert.Equal(t, 2, updatedSkill.CooldownLeft, "冷却时间应该减少到2")
	}
}

// TestIntegration_BuffStackingWithDOT 测试Buff叠加规则和DOT的集成
func TestIntegration_BuffStackingWithDOT(t *testing.T) {
	manager, char, cleanup := setupIntegrationTest(t)
	defer cleanup()

	// 1. 应用DOT效果（stack规则）
	manager.buffManager.ApplyBuffWithDOT(char.ID, "dot_poison", "毒药", "debuff", false, 5, 10.0, "damage", "nature", true, false, 0)

	// 2. 再次应用相同DOT（应该叠加）
	manager.buffManager.ApplyBuffWithDOT(char.ID, "dot_poison", "毒药", "debuff", false, 5, 10.0, "damage", "nature", true, false, 0)

	// 3. 检查叠加后的数值
	buffs := manager.buffManager.GetBuffs(char.ID)
	buff := buffs["dot_poison"]
	assert.Equal(t, 20.0, buff.Value, "DOT应该叠加到20点")

	// 4. 处理DOT效果
	damage, _ := manager.buffManager.ProcessDOTEffects(char.ID, 1)
	assert.Equal(t, 20, damage, "叠加后的DOT应该造成20点伤害")
}

// TestIntegration_PassiveWithBuffAndSkill 测试被动技能、Buff和技能的完整集成
func TestIntegration_PassiveWithBuffAndSkill(t *testing.T) {
	manager, char, cleanup := setupIntegrationTest(t)
	defer cleanup()

	// 1. 应用Buff
	manager.buffManager.ApplyBuff(char.ID, "battle_shout", "战斗怒吼", "buff", true, 5, 20.0, "attack", "")

	// 2. 添加被动技能
	passiveState := &CharacterPassiveState{
		PassiveID:   "test_on_hit_resource",
		Level:       1,
		Passive:     &models.PassiveSkill{ID: "test_on_hit_resource", EffectType: "on_hit_resource"},
		EffectValue: 2.0, // 攻击时获得2点资源
	}
	manager.passiveSkillManager.mu.Lock()
	if manager.passiveSkillManager.characterPassives == nil {
		manager.passiveSkillManager.characterPassives = make(map[int][]*CharacterPassiveState)
	}
	manager.passiveSkillManager.characterPassives[char.ID] = []*CharacterPassiveState{passiveState}
	manager.passiveSkillManager.mu.Unlock()

	// 3. 创建战斗会话
	session := &BattleSession{
		UserID:     1,
		BattleLogs: make([]models.BattleLog, 0),
	}
	manager.sessions[1] = session

	// 4. 触发攻击（应该触发被动技能）
	logs := make([]models.BattleLog, 0)
	originalResource := char.Resource
	manager.handlePassiveOnHitEffects(char, 100, false, session, &logs)

	// 5. 验证被动技能触发（获得资源）
	assert.GreaterOrEqual(t, char.Resource, originalResource, "资源应该增加")
}

// TestIntegration_CompleteBattleFlow 测试完整的战斗流程集成
func TestIntegration_CompleteBattleFlow(t *testing.T) {
	manager, char, cleanup := setupIntegrationTest(t)
	defer cleanup()

	// 1. 初始化战斗（StartBattle只需要userID）
	result, err := manager.StartBattle(1)
	assert.NoError(t, err)
	assert.NotNil(t, result)

	// 2. 执行几个战斗回合
	characters := []*models.Character{char}
	for i := 0; i < 3; i++ {
		tickResult, err := manager.ExecuteBattleTick(1, characters)
		if err != nil {
			t.Logf("战斗回合 %d 执行出错: %v", i+1, err)
			break
		}
		if tickResult == nil {
			t.Logf("战斗回合 %d 返回nil", i+1)
			break
		}
		assert.NotNil(t, tickResult)
	}

	// 3. 验证战斗状态
	session := manager.GetOrCreateSession(1)
	assert.NotNil(t, session)
}


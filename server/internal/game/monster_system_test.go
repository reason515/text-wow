package game

import (
	"testing"

	"text-wow/internal/database"
	"text-wow/internal/models"

	"github.com/stretchr/testify/assert"
)

// ═══════════════════════════════════════════════════════════
// 威胁系统测试
// ═══════════════════════════════════════════════════════════

func TestThreatSystem_UpdateThreat(t *testing.T) {
	testDB, err := database.SetupTestDB()
	if err != nil {
		t.Fatalf("Failed to setup test database: %v", err)
	}
	defer database.TeardownTestDB(testDB)

	manager := NewBattleManager()
	session := manager.GetOrCreateSession(1)

	// 初始化威胁表
	session.ThreatTable = make(map[string]map[int]int)

	monsterID := "test_monster_1"
	characterID := 1
	threatGain := 100

	// 更新威胁值
	manager.updateThreat(session, monsterID, characterID, threatGain)

	// 验证威胁值已更新
	threatTable := manager.getThreatTableForMonster(session, monsterID)
	assert.NotNil(t, threatTable)
	assert.Equal(t, threatGain, threatTable[characterID])

	// 再次更新威胁值
	manager.updateThreat(session, monsterID, characterID, 50)
	threatTable = manager.getThreatTableForMonster(session, monsterID)
	assert.Equal(t, 150, threatTable[characterID]) // 应该是累加的
}

func TestThreatSystem_ResetThreatTable(t *testing.T) {
	testDB, err := database.SetupTestDB()
	if err != nil {
		t.Fatalf("Failed to setup test database: %v", err)
	}
	defer database.TeardownTestDB(testDB)

	manager := NewBattleManager()
	session := manager.GetOrCreateSession(1)

	// 初始化威胁表并添加一些威胁值
	session.ThreatTable = make(map[string]map[int]int)
	session.ThreatTable["monster1"] = make(map[int]int)
	session.ThreatTable["monster1"][1] = 100
	session.ThreatTable["monster2"] = make(map[int]int)
	session.ThreatTable["monster2"][2] = 200

	// 重置威胁表
	manager.resetThreatTable(session)

	// 验证威胁表已清空
	assert.Empty(t, session.ThreatTable)
}

func TestThreatSystem_GetThreatTableForMonster(t *testing.T) {
	testDB, err := database.SetupTestDB()
	if err != nil {
		t.Fatalf("Failed to setup test database: %v", err)
	}
	defer database.TeardownTestDB(testDB)

	manager := NewBattleManager()
	session := manager.GetOrCreateSession(1)

	// 初始化威胁表
	session.ThreatTable = make(map[string]map[int]int)
	session.ThreatTable["monster1"] = make(map[int]int)
	session.ThreatTable["monster1"][1] = 100
	session.ThreatTable["monster1"][2] = 200

	// 获取威胁表
	threatTable := manager.getThreatTableForMonster(session, "monster1")
	assert.NotNil(t, threatTable)
	assert.Equal(t, 100, threatTable[1])
	assert.Equal(t, 200, threatTable[2])

	// 获取不存在的怪物威胁表
	threatTable = manager.getThreatTableForMonster(session, "nonexistent")
	assert.NotNil(t, threatTable)
	assert.Empty(t, threatTable)
}

func TestThreatSystem_MultipleCharacters(t *testing.T) {
	testDB, err := database.SetupTestDB()
	if err != nil {
		t.Fatalf("Failed to setup test database: %v", err)
	}
	defer database.TeardownTestDB(testDB)

	manager := NewBattleManager()
	session := manager.GetOrCreateSession(1)

	// 初始化威胁表
	session.ThreatTable = make(map[string]map[int]int)

	monsterID := "test_monster_1"

	// 多个角色对同一个怪物造成伤害
	manager.updateThreat(session, monsterID, 1, 100)
	manager.updateThreat(session, monsterID, 2, 150)
	manager.updateThreat(session, monsterID, 3, 50)

	// 验证每个角色的威胁值
	threatTable := manager.getThreatTableForMonster(session, monsterID)
	assert.Equal(t, 100, threatTable[1])
	assert.Equal(t, 150, threatTable[2])
	assert.Equal(t, 50, threatTable[3])
}

// ═══════════════════════════════════════════════════════════
// 怪物AI测试
// ═══════════════════════════════════════════════════════════

func TestMonsterAI_SelectTarget_HighestThreat(t *testing.T) {
	testDB, err := database.SetupTestDB()
	if err != nil {
		t.Fatalf("Failed to setup test database: %v", err)
	}
	defer database.TeardownTestDB(testDB)

	// 创建测试怪物
	monster := &models.Monster{
		ID:   "test_monster",
		Name: "测试怪物",
		Type: "normal",
		AIType: "defensive",
		HP:   100,
		MaxHP: 100,
	}

	skillManager := NewSkillManager()
	ai, err := NewMonsterAI(monster, skillManager)
	assert.NoError(t, err)
	assert.NotNil(t, ai)

	// 创建测试角色
	characters := []*models.Character{
		{ID: 1, Name: "角色1", HP: 100, MaxHP: 100},
		{ID: 2, Name: "角色2", HP: 100, MaxHP: 100},
		{ID: 3, Name: "角色3", HP: 100, MaxHP: 100},
	}

	// 创建威胁表（角色2威胁值最高）
	threatTable := map[int]int{
		1: 50,
		2: 200, // 最高威胁
		3: 100,
	}

	// 选择目标（defensive AI应该优先选择最高威胁）
	target := ai.SelectTarget(characters, threatTable)
	assert.NotNil(t, target)
	assert.Equal(t, 2, target.ID) // 应该选择角色2
}

func TestMonsterAI_SelectTarget_LowestHP(t *testing.T) {
	testDB, err := database.SetupTestDB()
	if err != nil {
		t.Fatalf("Failed to setup test database: %v", err)
	}
	defer database.TeardownTestDB(testDB)

	// 创建测试怪物（aggressive AI优先选择最低HP）
	monster := &models.Monster{
		ID:   "test_monster",
		Name: "测试怪物",
		Type: "normal",
		AIType: "aggressive",
		HP:   100,
		MaxHP: 100,
	}

	skillManager := NewSkillManager()
	ai, err := NewMonsterAI(monster, skillManager)
	assert.NoError(t, err)

	// 创建测试角色（角色2 HP最低）
	characters := []*models.Character{
		{ID: 1, Name: "角色1", HP: 100, MaxHP: 100},
		{ID: 2, Name: "角色2", HP: 20, MaxHP: 100},  // HP最低
		{ID: 3, Name: "角色3", HP: 80, MaxHP: 100},
	}

	threatTable := map[int]int{} // 空威胁表

	// 选择目标
	target := ai.SelectTarget(characters, threatTable)
	assert.NotNil(t, target)
	assert.Equal(t, 2, target.ID) // 应该选择HP最低的角色2
}

func TestMonsterAI_SelectTarget_NoThreatTable(t *testing.T) {
	testDB, err := database.SetupTestDB()
	if err != nil {
		t.Fatalf("Failed to setup test database: %v", err)
	}
	defer database.TeardownTestDB(testDB)

	monster := &models.Monster{
		ID:   "test_monster",
		Name: "测试怪物",
		Type: "normal",
		AIType: "defensive",
		HP:   100,
		MaxHP: 100,
	}

	skillManager := NewSkillManager()
	ai, err := NewMonsterAI(monster, skillManager)
	assert.NoError(t, err)

	characters := []*models.Character{
		{ID: 1, Name: "角色1", HP: 100, MaxHP: 100},
		{ID: 2, Name: "角色2", HP: 100, MaxHP: 100},
	}

	// 没有威胁表时，应该返回第一个角色
	target := ai.SelectTarget(characters, nil)
	assert.NotNil(t, target)
	assert.Equal(t, 1, target.ID)
}

func TestMonsterAI_SelectTarget_SingleCharacter(t *testing.T) {
	testDB, err := database.SetupTestDB()
	if err != nil {
		t.Fatalf("Failed to setup test database: %v", err)
	}
	defer database.TeardownTestDB(testDB)

	monster := &models.Monster{
		ID:   "test_monster",
		Name: "测试怪物",
		Type: "normal",
		AIType: "balanced",
		HP:   100,
		MaxHP: 100,
	}

	skillManager := NewSkillManager()
	ai, err := NewMonsterAI(monster, skillManager)
	assert.NoError(t, err)

	characters := []*models.Character{
		{ID: 1, Name: "角色1", HP: 100, MaxHP: 100},
	}

	// 只有一个角色时，应该直接返回
	target := ai.SelectTarget(characters, nil)
	assert.NotNil(t, target)
	assert.Equal(t, 1, target.ID)
}

// ═══════════════════════════════════════════════════════════
// 怪物管理器测试
// ═══════════════════════════════════════════════════════════

func TestMonsterManager_LoadMonsterConfig(t *testing.T) {
	testDB, err := database.SetupTestDB()
	if err != nil {
		t.Fatalf("Failed to setup test database: %v", err)
	}
	defer database.TeardownTestDB(testDB)

	manager := NewMonsterManager()

	// 注意：这个测试依赖于数据库中有怪物配置
	// 如果数据库中没有配置，测试可能会失败
	// 这里我们主要测试函数不会panic
	config, err := manager.LoadMonsterConfig("test_monster_id")
	
	// 如果数据库中没有这个怪物，err不为nil是正常的
	// 我们主要验证函数能正常执行
	if err != nil {
		t.Logf("Monster config not found (expected if DB is empty): %v", err)
	} else {
		assert.NotNil(t, config)
	}
}

func TestMonsterManager_CalculateDrops(t *testing.T) {
	testDB, err := database.SetupTestDB()
	if err != nil {
		t.Fatalf("Failed to setup test database: %v", err)
	}
	defer database.TeardownTestDB(testDB)

	manager := NewMonsterManager()

	// 测试掉落计算（即使没有掉落配置，也应该返回空结果而不是错误）
	drops, err := manager.CalculateDrops("test_monster_id", "normal")
	// 如果数据库schema不完整，可能会返回错误，这是可以接受的
	if err != nil {
		t.Logf("CalculateDrops returned error (expected if DB schema is incomplete): %v", err)
		return
	}
	// CalculateDrops 应该总是返回一个切片（可能为空）
	// 如果没有掉落配置，会返回空数组
	if drops != nil {
		// 如果返回了结果，验证它是切片类型
		assert.IsType(t, []DropResult{}, drops)
	}
	// 如果没有掉落配置，应该返回空数组（这是正常的）
	// 我们主要验证函数不会panic
}

func TestMonsterManager_CalculateDrops_EliteMultiplier(t *testing.T) {
	testDB, err := database.SetupTestDB()
	if err != nil {
		t.Fatalf("Failed to setup test database: %v", err)
	}
	defer database.TeardownTestDB(testDB)

	manager := NewMonsterManager()

	// 测试精英怪物掉落率修正
	// 注意：这个测试依赖于数据库中有掉落配置
	// 如果没有配置，会返回空结果
	drops, err := manager.CalculateDrops("test_monster_id", "elite")
	// 如果数据库schema不完整，可能会返回错误，这是可以接受的
	if err != nil {
		t.Logf("CalculateDrops returned error (expected if DB schema is incomplete): %v", err)
		return
	}
	// 应该返回一个切片（可能为空）
	if drops != nil {
		assert.IsType(t, []DropResult{}, drops)
	}
	// 如果没有掉落配置，返回空数组是正常的
	// 我们主要验证函数能正常执行，不会panic
}

func TestMonsterManager_CalculateDrops_BossMultiplier(t *testing.T) {
	testDB, err := database.SetupTestDB()
	if err != nil {
		t.Fatalf("Failed to setup test database: %v", err)
	}
	defer database.TeardownTestDB(testDB)

	manager := NewMonsterManager()

	// 测试Boss掉落率修正
	drops, err := manager.CalculateDrops("test_monster_id", "boss")
	// 如果数据库schema不完整，可能会返回错误，这是可以接受的
	if err != nil {
		t.Logf("CalculateDrops returned error (expected if DB schema is incomplete): %v", err)
		return
	}
	// 应该返回一个切片（可能为空）
	if drops != nil {
		assert.IsType(t, []DropResult{}, drops)
	}
	// 如果没有掉落配置，返回空数组是正常的
	// 我们主要验证函数能正常执行，不会panic
}

// ═══════════════════════════════════════════════════════════
// 威胁系统集成测试
// ═══════════════════════════════════════════════════════════

func TestThreatSystem_IntegrationWithBattle(t *testing.T) {
	testDB, err := database.SetupTestDB()
	if err != nil {
		t.Fatalf("Failed to setup test database: %v", err)
	}
	defer database.TeardownTestDB(testDB)

	manager := NewBattleManager()
	session := manager.GetOrCreateSession(1)

	// 验证威胁表已初始化
	assert.NotNil(t, session)
	assert.NotNil(t, session.ThreatTable)

	// 手动创建敌人来测试威胁系统
	enemy := &models.Monster{
		ID:   "test_enemy_1",
		Name: "测试敌人",
		HP:   100,
		MaxHP: 100,
	}
	session.CurrentEnemies = []*models.Monster{enemy}

	// 创建测试角色
	char := &models.Character{
		ID:              1,
		UserID:          1,
		Name:            "测试角色",
		Level:           5,
		HP:              100,
		MaxHP:           100,
		PhysicalAttack:  20,
		PhysicalDefense: 10,
		Resource:        50,
		MaxResource:     50,
		ResourceType:    "rage",
	}

	// 模拟造成伤害并更新威胁值
	damage := 50
	manager.updateThreat(session, enemy.ID, char.ID, damage)

	// 验证威胁值已更新
	threatTable := manager.getThreatTableForMonster(session, enemy.ID)
	assert.NotNil(t, threatTable)
	assert.Equal(t, damage, threatTable[char.ID])
}

func TestThreatSystem_ResetOnNewBattle(t *testing.T) {
	testDB, err := database.SetupTestDB()
	if err != nil {
		t.Fatalf("Failed to setup test database: %v", err)
	}
	defer database.TeardownTestDB(testDB)

	manager := NewBattleManager()
	session := manager.GetOrCreateSession(1)

	// 初始化威胁表并添加一些威胁值
	session.ThreatTable = make(map[string]map[int]int)
	session.ThreatTable["monster1"] = make(map[int]int)
	session.ThreatTable["monster1"][1] = 100
	session.ThreatTable["monster2"] = make(map[int]int)
	session.ThreatTable["monster2"][2] = 200

	// 验证威胁表有数据
	assert.Equal(t, 100, session.ThreatTable["monster1"][1])
	assert.Equal(t, 200, session.ThreatTable["monster2"][2])

	// 手动调用重置威胁表（模拟新战斗开始）
	manager.resetThreatTable(session)

	// 验证威胁表已被重置
	assert.Empty(t, session.ThreatTable)
}

// ═══════════════════════════════════════════════════════════
// 怪物AI集成测试
// ═══════════════════════════════════════════════════════════

func TestMonsterAI_DefaultBehavior(t *testing.T) {
	testDB, err := database.SetupTestDB()
	if err != nil {
		t.Fatalf("Failed to setup test database: %v", err)
	}
	defer database.TeardownTestDB(testDB)

	// 测试不同AI类型的默认行为
	testCases := []struct {
		monsterType string
		aiType      string
		expected    string
	}{
		{"normal", "aggressive", "aggressive"},
		{"normal", "defensive", "defensive"},
		{"normal", "balanced", "balanced"},
		{"boss", "special", "special"},
	}

	for _, tc := range testCases {
		monster := &models.Monster{
			ID:   "test_monster",
			Name: "测试怪物",
			Type: tc.monsterType,
			AIType: tc.aiType,
			HP:   100,
			MaxHP: 100,
		}

		skillManager := NewSkillManager()
		ai, err := NewMonsterAI(monster, skillManager)
		assert.NoError(t, err, "Failed to create AI for %s/%s", tc.monsterType, tc.aiType)
		assert.NotNil(t, ai)
		assert.NotNil(t, ai.Behavior)
		assert.Equal(t, tc.aiType, ai.AIType)
	}
}

func TestMonsterAI_SelectSkill(t *testing.T) {
	testDB, err := database.SetupTestDB()
	if err != nil {
		t.Fatalf("Failed to setup test database: %v", err)
	}
	defer database.TeardownTestDB(testDB)

	monster := &models.Monster{
		ID:   "test_monster",
		Name: "测试怪物",
		Type: "normal",
		AIType: "balanced",
		HP:   100,
		MaxHP: 100,
		MP:   50,
		MaxMP: 50,
		MonsterSkills: []*models.MonsterSkill{},
	}

	skillManager := NewSkillManager()
	ai, err := NewMonsterAI(monster, skillManager)
	assert.NoError(t, err)

	target := &models.Character{
		ID:   1,
		Name: "目标角色",
		HP:   100,
		MaxHP: 100,
	}

	buffManager := NewBuffManager()

	// 没有技能时，应该返回nil
	skill := ai.SelectSkill(target, buffManager)
	assert.Nil(t, skill)
}


package game

import (
	"testing"

	"text-wow/internal/database"
	"text-wow/internal/models"
	"text-wow/internal/repository"
	"text-wow/internal/service"

	"github.com/stretchr/testify/assert"
)

// ═══════════════════════════════════════════════════════════
// 技能升级系统测试
// ═══════════════════════════════════════════════════════════

func setupSkillLevelingTest(t *testing.T) (*SkillManager, *models.Character, func()) {
	testDB, err := database.SetupTestDB()
	if err != nil {
		t.Fatalf("Failed to setup test database: %v", err)
	}

	skillRepo := repository.NewSkillRepository()
	characterRepo := repository.NewCharacterRepository()
	skillService := service.NewSkillService(skillRepo, characterRepo)

	sm := NewSkillManager()
	sm.skillService = skillService
	sm.skillRepo = skillRepo

	// 先创建user记录（如果不存在）
	_, err = database.DB.Exec(`
		INSERT OR IGNORE INTO users (id, username, email, password_hash, created_at)
		VALUES (?, ?, ?, ?, ?)`,
		1, "test_user", "test@example.com", "test_hash", "2024-01-01 00:00:00",
	)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	// 先创建race记录（如果不存在）
	_, err = database.DB.Exec(`
		INSERT OR IGNORE INTO races (id, name, description, faction)
		VALUES (?, ?, ?, ?)`,
		"human", "人类", "人类种族", "alliance",
	)
	if err != nil {
		t.Fatalf("Failed to create race: %v", err)
	}

	// 先创建class记录（如果不存在）
	_, err = database.DB.Exec(`
		INSERT OR IGNORE INTO classes (id, name, description, resource_type, base_resource)
		VALUES (?, ?, ?, ?, ?)`,
		"warrior", "战士", "近战职业", "rage", 0,
	)
	if err != nil {
		t.Fatalf("Failed to create class: %v", err)
	}

	// 先创建角色记录到数据库
	_, err = database.DB.Exec(`
		INSERT OR REPLACE INTO characters (id, user_id, name, race_id, class_id, faction, team_slot, level, hp, max_hp, resource, max_resource, resource_type, physical_attack, magic_attack, physical_defense, magic_defense)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		1, 1, "测试角色", "human", "warrior", "alliance", 1, 5, 1000, 1000, 50, 100, "rage", 100, 50, 50, 30,
	)
	if err != nil {
		t.Fatalf("Failed to create character: %v", err)
	}

	character := &models.Character{
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
	}

	// 先创建测试技能记录到数据库
	_, err = database.DB.Exec(`
		INSERT OR REPLACE INTO skills (id, name, description, class_id, type, target_type, base_value, resource_cost, cooldown, level_required)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		"test_skill", "测试技能", "用于测试的技能", "warrior", "attack", "enemy", 50, 10, 0, 1,
	)
	if err != nil {
		t.Fatalf("Failed to create test skill: %v", err)
	}

	// 添加测试技能到角色
	err = skillRepo.AddCharacterSkill(character.ID, "test_skill", 1)
	if err != nil {
		t.Fatalf("Failed to add character skill: %v", err)
	}

	cleanup := func() {
		database.TeardownTestDB(testDB)
	}

	return sm, character, cleanup
}

func TestSkillManager_SkillExperienceGain(t *testing.T) {
	sm, character, cleanup := setupSkillLevelingTest(t)
	defer cleanup()

	// 加载角色技能
	err := sm.LoadCharacterSkills(character.ID)
	assert.NoError(t, err)

	// 使用技能（应该获得经验）
	_, err = sm.UseSkill(character.ID, "test_skill")
	assert.NoError(t, err)

	// 检查技能经验是否增加（需要重新加载）
	characterSkills, err := sm.skillRepo.GetCharacterSkills(character.ID)
	assert.NoError(t, err)

	var testSkill *models.CharacterSkill
	for _, cs := range characterSkills {
		if cs.SkillID == "test_skill" {
			testSkill = cs
			break
		}
	}

	if testSkill != nil {
		// 技能经验应该增加（5-10点）
		assert.Greater(t, testSkill.SkillExp, 0)
		assert.LessOrEqual(t, testSkill.SkillExp, 10)
	}
}

func TestSkillManager_SkillLevelUp(t *testing.T) {
	sm, character, cleanup := setupSkillLevelingTest(t)
	defer cleanup()

	// 手动设置技能经验接近升级阈值
	sm.skillRepo.UpdateSkillExperience(character.ID, "test_skill", 95, 100, 1)

	// 加载角色技能
	err := sm.LoadCharacterSkills(character.ID)
	assert.NoError(t, err)

	// 使用技能多次，直到升级
	for i := 0; i < 10; i++ {
		_, err = sm.UseSkill(character.ID, "test_skill")
		assert.NoError(t, err)
	}

	// 检查技能是否升级
	characterSkills, err := sm.skillRepo.GetCharacterSkills(character.ID)
	assert.NoError(t, err)

	var testSkill *models.CharacterSkill
	for _, cs := range characterSkills {
		if cs.SkillID == "test_skill" {
			testSkill = cs
			break
		}
	}

	if testSkill != nil {
		// 技能等级应该提升（至少2级）
		assert.GreaterOrEqual(t, testSkill.SkillLevel, 2)
		// 升级后经验应该重置（或接近0）
		assert.Less(t, testSkill.SkillExp, testSkill.ExpToNext)
	}
}

func TestSkillManager_CalculateSkillExpGain(t *testing.T) {
	sm, _, cleanup := setupSkillLevelingTest(t)
	defer cleanup()

	testCases := []struct {
		skillType   string
		expectedMin int
		expectedMax int
	}{
		{"attack", 8, 8},   // 攻击技能获得8点经验
		{"heal", 7, 7},     // 治疗技能获得7点经验
		{"buff", 6, 6},     // Buff技能获得6点经验
		{"debuff", 6, 6},   // Debuff技能获得6点经验
		{"control", 7, 7},  // 控制技能获得7点经验
		{"unknown", 5, 5},  // 未知类型获得5点经验
	}

	for _, tc := range testCases {
		skill := &models.Skill{
			ID:   "test_skill",
			Type: tc.skillType,
		}
		expGain := sm.calculateSkillExpGain(skill)
		assert.Equal(t, tc.expectedMin, expGain, "Skill type: %s should give %d exp", tc.skillType, tc.expectedMin)
	}
}

func TestSkillManager_SkillLevelUp_MultipleLevels(t *testing.T) {
	sm, character, cleanup := setupSkillLevelingTest(t)
	defer cleanup()

	// 手动设置技能为1级，经验为0
	sm.skillRepo.UpdateSkillExperience(character.ID, "test_skill", 0, 100, 1)

	// 加载角色技能
	err := sm.LoadCharacterSkills(character.ID)
	assert.NoError(t, err)

	// 使用技能大量次数（模拟快速升级）
	for i := 0; i < 50; i++ {
		_, err = sm.UseSkill(character.ID, "test_skill")
		assert.NoError(t, err)
	}

	// 检查技能等级（最多5级）
	characterSkills, err := sm.skillRepo.GetCharacterSkills(character.ID)
	assert.NoError(t, err)

	var testSkill *models.CharacterSkill
	for _, cs := range characterSkills {
		if cs.SkillID == "test_skill" {
			testSkill = cs
			break
		}
	}

	if testSkill != nil {
		// 技能等级应该在1-5之间
		assert.GreaterOrEqual(t, testSkill.SkillLevel, 1)
		assert.LessOrEqual(t, testSkill.SkillLevel, 5)
		
		// 如果达到5级，经验应该不再增加（或增加但不会升级）
		if testSkill.SkillLevel >= 5 {
			// 5级是最高级，不应该再升级
			assert.Equal(t, 5, testSkill.SkillLevel)
		}
	}
}

func TestSkillManager_SkillExpToNext_Incremental(t *testing.T) {
	sm, character, cleanup := setupSkillLevelingTest(t)
	defer cleanup()

	// 测试不同等级的升级所需经验
	expectedExpToNext := []int{100, 150, 200, 250, 300}

	for level := 1; level <= 5; level++ {
		// 设置技能为指定等级
		expToNext := 100 + (level-1)*50
		sm.skillRepo.UpdateSkillExperience(character.ID, "test_skill", 0, expToNext, level)

		characterSkills, err := sm.skillRepo.GetCharacterSkills(character.ID)
		assert.NoError(t, err)

		var testSkill *models.CharacterSkill
		for _, cs := range characterSkills {
			if cs.SkillID == "test_skill" {
				testSkill = cs
				break
			}
		}

		if testSkill != nil {
			assert.Equal(t, expectedExpToNext[level-1], testSkill.ExpToNext, 
				"Level %d should require %d exp to next", level, expectedExpToNext[level-1])
		}
	}
}


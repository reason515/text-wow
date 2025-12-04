package repository

import (
	"testing"

	"text-wow/internal/database"

	"github.com/stretchr/testify/assert"
)

// ═══════════════════════════════════════════════════════════
// 区域查询测试 - 验证列名正确性
// ═══════════════════════════════════════════════════════════

func TestGameRepository_GetZoneByID_ColumnNames(t *testing.T) {
	// 测试：验证区域查询使用正确的列名（exp_modifier, gold_modifier）
	testDB, err := database.SetupTestDB()
	if err != nil {
		t.Fatalf("Failed to setup test database: %v", err)
	}
	defer database.TeardownTestDB(testDB)

	// 确保区域数据存在（testdb.go应该已经插入了，但为了保险，我们检查一下）
	_, err = testDB.Exec(`INSERT OR IGNORE INTO zones (id, name, description, min_level, max_level, faction, exp_modifier, gold_modifier)
		VALUES ('elwynn', '艾尔文森林', '测试区域', 1, 10, 'alliance', 1.0, 1.0)`)
	if err != nil {
		t.Fatalf("Failed to insert test zone: %v", err)
	}

	repo := NewGameRepository()

	// 查询存在的区域
	zone, err := repo.GetZoneByID("elwynn")
	assert.NoError(t, err, "Should successfully query zone with correct column names")
	assert.NotNil(t, zone, "Zone should not be nil")
	if zone != nil {
		assert.Equal(t, "elwynn", zone.ID, "Zone ID should match")
		assert.Equal(t, 1.0, zone.ExpMulti, "ExpMulti should be 1.0")
		assert.Equal(t, 1.0, zone.GoldMulti, "GoldMulti should be 1.0")
	}
}

func TestGameRepository_GetZones_ColumnNames(t *testing.T) {
	// 测试：验证获取所有区域时列名正确
	testDB, err := database.SetupTestDB()
	if err != nil {
		t.Fatalf("Failed to setup test database: %v", err)
	}
	defer database.TeardownTestDB(testDB)

	// 确保区域数据存在
	_, err = testDB.Exec(`INSERT OR IGNORE INTO zones (id, name, description, min_level, max_level, faction, exp_modifier, gold_modifier)
		VALUES ('elwynn', '艾尔文森林', '测试区域', 1, 10, 'alliance', 1.0, 1.0)`)
	if err != nil {
		t.Fatalf("Failed to insert test zone: %v", err)
	}

	repo := NewGameRepository()

	zones, err := repo.GetZones()
	assert.NoError(t, err, "Should successfully query zones with correct column names")
	assert.Greater(t, len(zones), 0, "Should have at least one zone")

	// 验证每个区域都有正确的ExpMulti和GoldMulti值
	for _, zone := range zones {
		assert.GreaterOrEqual(t, zone.ExpMulti, 0.0, "ExpMulti should be non-negative")
		assert.GreaterOrEqual(t, zone.GoldMulti, 0.0, "GoldMulti should be non-negative")
	}
}

func TestGameRepository_GetZoneByID_NotFound(t *testing.T) {
	// 测试：查询不存在的区域
	testDB, err := database.SetupTestDB()
	if err != nil {
		t.Fatalf("Failed to setup test database: %v", err)
	}
	defer database.TeardownTestDB(testDB)

	repo := NewGameRepository()

	// 查询不存在的区域
	zone, err := repo.GetZoneByID("nonexistent_zone")
	assert.Error(t, err, "Should return error for non-existent zone")
	assert.Nil(t, zone, "Zone should be nil when not found")
}

func TestGameRepository_GetZoneByID_WrongIDFormat(t *testing.T) {
	// 测试：使用错误的区域ID格式（如 elwynn_forest 而不是 elwynn）
	testDB, err := database.SetupTestDB()
	if err != nil {
		t.Fatalf("Failed to setup test database: %v", err)
	}
	defer database.TeardownTestDB(testDB)

	repo := NewGameRepository()

	// 使用错误的ID格式
	zone, err := repo.GetZoneByID("elwynn_forest")
	assert.Error(t, err, "Should return error for wrong zone ID format")
	assert.Nil(t, zone, "Zone should be nil when ID format is wrong")
}

func TestGameRepository_GetMonstersByZone_ValidZone(t *testing.T) {
	// 测试：获取有效区域的怪物
	testDB, err := database.SetupTestDB()
	if err != nil {
		t.Fatalf("Failed to setup test database: %v", err)
	}
	defer database.TeardownTestDB(testDB)

	// 确保区域和怪物数据存在
	_, err = testDB.Exec(`INSERT OR IGNORE INTO zones (id, name, description, min_level, max_level, faction, exp_modifier, gold_modifier)
		VALUES ('elwynn', '艾尔文森林', '测试区域', 1, 10, 'alliance', 1.0, 1.0)`)
	if err != nil {
		t.Fatalf("Failed to insert test zone: %v", err)
	}
	_, err = testDB.Exec(`INSERT OR IGNORE INTO monsters (id, zone_id, name, level, type, hp, physical_attack, magic_attack, physical_defense, magic_defense, exp_reward, gold_min, gold_max, spawn_weight)
		VALUES ('wolf', 'elwynn', '森林狼', 2, 'normal', 30, 8, 4, 2, 1, 15, 1, 5, 100)`)
	if err != nil {
		t.Fatalf("Failed to insert test monster: %v", err)
	}

	repo := NewGameRepository()

	monsters, err := repo.GetMonstersByZone("elwynn")
	assert.NoError(t, err, "Should successfully get monsters for valid zone")
	assert.Greater(t, len(monsters), 0, "Should have at least one monster in elwynn zone")

	// 验证怪物数据
	for _, monster := range monsters {
		assert.Equal(t, "elwynn", monster.ZoneID, "Monster zone ID should match")
		assert.Greater(t, monster.HP, 0, "Monster HP should be positive")
		assert.Greater(t, monster.ExpReward, 0, "Monster exp reward should be positive")
	}
}

func TestGameRepository_GetMonstersByZone_InvalidZone(t *testing.T) {
	// 测试：获取不存在区域的怪物
	testDB, err := database.SetupTestDB()
	if err != nil {
		t.Fatalf("Failed to setup test database: %v", err)
	}
	defer database.TeardownTestDB(testDB)

	repo := NewGameRepository()

	monsters, err := repo.GetMonstersByZone("nonexistent_zone")
	assert.NoError(t, err, "Should not error for non-existent zone (returns empty list)")
	assert.Equal(t, 0, len(monsters), "Should return empty list for non-existent zone")
}

func TestGameRepository_GetMonstersByZone_FieldNames(t *testing.T) {
	// 测试：验证怪物查询使用正确的字段名（physical_attack, magic_attack等）
	// 这个测试确保不会因为字段名不匹配而失败
	testDB, err := database.SetupTestDB()
	if err != nil {
		t.Fatalf("Failed to setup test database: %v", err)
	}
	defer database.TeardownTestDB(testDB)

	// 确保有怪物数据（使用正确的字段名）
	_, err = testDB.Exec(`INSERT OR IGNORE INTO zones (id, name, description, min_level, max_level, faction, exp_modifier, gold_modifier)
		VALUES ('test_zone', '测试区域', '测试', 1, 10, NULL, 1.0, 1.0)`)
	if err != nil {
		t.Fatalf("Failed to insert test zone: %v", err)
	}
	
	_, err = testDB.Exec(`INSERT OR IGNORE INTO monsters (id, zone_id, name, level, type, hp, physical_attack, magic_attack, physical_defense, magic_defense, exp_reward, gold_min, gold_max, spawn_weight)
		VALUES ('test_monster', 'test_zone', '测试怪物', 1, 'normal', 20, 5, 2, 3, 1, 10, 1, 3, 100)`)
	if err != nil {
		t.Fatalf("Failed to insert test monster: %v", err)
	}

	repo := NewGameRepository()

	monsters, err := repo.GetMonstersByZone("test_zone")
	assert.NoError(t, err, "Should successfully query monsters with correct field names")
	assert.Greater(t, len(monsters), 0, "Should have at least one monster")

	// 验证怪物数据字段都正确加载
	for _, monster := range monsters {
		assert.Greater(t, monster.PhysicalAttack, 0, "PhysicalAttack should be loaded correctly")
		assert.GreaterOrEqual(t, monster.MagicAttack, 0, "MagicAttack should be loaded correctly")
		assert.Greater(t, monster.PhysicalDefense, 0, "PhysicalDefense should be loaded correctly")
		assert.Greater(t, monster.MagicDefense, 0, "MagicDefense should be loaded correctly")
		assert.Equal(t, monster.HP, monster.MaxHP, "MaxHP should be set to HP")
	}
}

func TestGameRepository_GetMonstersByZone_EmptyZone(t *testing.T) {
	// 测试：区域存在但没有怪物（应该返回空列表，不应该报错）
	testDB, err := database.SetupTestDB()
	if err != nil {
		t.Fatalf("Failed to setup test database: %v", err)
	}
	defer database.TeardownTestDB(testDB)

	// 创建一个没有怪物的区域
	_, err = testDB.Exec(`INSERT OR IGNORE INTO zones (id, name, description, min_level, max_level, faction, exp_modifier, gold_modifier)
		VALUES ('empty_zone', '空区域', '没有怪物的区域', 1, 10, NULL, 1.0, 1.0)`)
	if err != nil {
		t.Fatalf("Failed to insert test zone: %v", err)
	}

	repo := NewGameRepository()

	monsters, err := repo.GetMonstersByZone("empty_zone")
	assert.NoError(t, err, "Should not error for zone with no monsters")
	assert.Equal(t, 0, len(monsters), "Should return empty list for zone with no monsters")
}


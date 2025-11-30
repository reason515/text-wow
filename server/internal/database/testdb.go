package database

import (
	"database/sql"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

// SetupTestDB 创建一个用于测试的内存数据库
func SetupTestDB() (*sql.DB, error) {
	testDB, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		return nil, err
	}

	// 读取 schema.sql 初始化表结构
	schemaPath := filepath.Join(getProjectRoot(), "database", "schema.sql")
	schema, err := os.ReadFile(schemaPath)
	if err != nil {
		// 如果找不到文件，使用内联的基础schema（包含测试数据）
		schema = []byte(getBaseSchema())
		_, err = testDB.Exec(string(schema))
		if err != nil {
			testDB.Close()
			return nil, err
		}
	} else {
		// 如果找到了schema.sql，先执行表结构
		_, err = testDB.Exec(string(schema))
		if err != nil {
			testDB.Close()
			return nil, err
		}
		// 然后执行测试数据插入（从getBaseSchema中提取INSERT语句）
		testData := getTestDataInserts()
		_, err = testDB.Exec(testData)
		if err != nil {
			testDB.Close()
			return nil, err
		}
	}

	// 设置全局DB为测试数据库
	DB = testDB

	return testDB, nil
}

// TeardownTestDB 清理测试数据库
func TeardownTestDB(testDB *sql.DB) {
	if testDB != nil {
		testDB.Close()
	}
	DB = nil
}

// getProjectRoot 获取项目根目录
func getProjectRoot() string {
	// 尝试从当前目录向上查找
	dir, _ := os.Getwd()
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return "."
}

// getTestDataInserts 返回测试数据INSERT语句
func getTestDataInserts() string {
	return `
-- 插入测试数据
INSERT OR IGNORE INTO races (id, name, faction, description, strength_base, agility_base, intellect_base, stamina_base, spirit_base, strength_pct, agility_pct, intellect_pct, stamina_pct, spirit_pct)
VALUES 
    ('human', '人类', 'alliance', '艾泽拉斯的人类', 2, 0, 0, 0, 2, 0, 0, 0, 0, 0),
    ('orc', '兽人', 'horde', '德拉诺的兽人', 3, 0, 0, 1, 0, 0, 0, 0, 0, 0);

INSERT OR IGNORE INTO classes (id, name, role, description, primary_stat, resource_type, base_hp, base_resource, hp_per_level, resource_per_level, resource_regen, resource_regen_pct, base_strength, base_agility, base_intellect, base_stamina, base_spirit, base_threat_modifier, combat_role, is_ranged)
VALUES 
    ('warrior', '战士', 'tank', '护甲坚固的近战职业', 'strength', 'rage', 120, 100, 15, 10, 5.0, 0.0, 15, 10, 5, 12, 8, 1.0, 'tank', 0),
    ('mage', '法师', 'dps', '强大的奥术施法者', 'intellect', 'mana', 80, 200, 8, 15, 2.0, 0.05, 5, 8, 18, 8, 12, 0.7, 'dps', 1);

INSERT OR IGNORE INTO zones (id, name, description, min_level, max_level, faction, exp_modifier, gold_modifier)
VALUES 
    ('elwynn', '艾尔文森林', '暴风城周边的和平森林', 1, 10, 'alliance', 1.0, 1.0),
    ('durotar', '杜隆塔尔', '兽人的荒芜家园', 1, 10, 'horde', 1.0, 1.0);

INSERT OR IGNORE INTO monsters (id, zone_id, name, level, type, hp, attack, defense, exp_reward, gold_min, gold_max, spawn_weight)
VALUES 
    ('wolf', 'elwynn', '森林狼', 2, 'normal', 30, 8, 2, 15, 1, 5, 100),
    ('kobold', 'elwynn', '狗头人', 3, 'normal', 40, 10, 3, 20, 2, 8, 80),
    ('boar', 'durotar', '野猪', 2, 'normal', 35, 9, 3, 18, 2, 6, 100);
`
}

// getBaseSchema 返回基础的数据库schema
func getBaseSchema() string {
	return `
-- 用户表
CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    email TEXT,
    max_team_size INTEGER DEFAULT 5,
    unlocked_slots INTEGER DEFAULT 1,
    gold INTEGER DEFAULT 0,
    current_zone_id TEXT DEFAULT 'elwynn',
    total_kills INTEGER DEFAULT 0,
    total_gold_gained INTEGER DEFAULT 0,
    play_time INTEGER DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    last_login_at DATETIME,
    status INTEGER DEFAULT 1
);

-- 种族表 (匹配真实schema)
CREATE TABLE IF NOT EXISTS races (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    faction TEXT NOT NULL,
    description TEXT,
    strength_base INTEGER DEFAULT 0,
    agility_base INTEGER DEFAULT 0,
    intellect_base INTEGER DEFAULT 0,
    stamina_base INTEGER DEFAULT 0,
    spirit_base INTEGER DEFAULT 0,
    strength_pct REAL DEFAULT 0,
    agility_pct REAL DEFAULT 0,
    intellect_pct REAL DEFAULT 0,
    stamina_pct REAL DEFAULT 0,
    spirit_pct REAL DEFAULT 0
);

-- 职业表 (匹配真实schema)
CREATE TABLE IF NOT EXISTS classes (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT,
    role TEXT NOT NULL,
    primary_stat TEXT DEFAULT 'strength',
    resource_type TEXT DEFAULT 'mana',
    base_hp INTEGER DEFAULT 100,
    base_resource INTEGER DEFAULT 100,
    hp_per_level INTEGER DEFAULT 10,
    resource_per_level INTEGER DEFAULT 5,
    resource_regen REAL DEFAULT 1.0,
    resource_regen_pct REAL DEFAULT 0.0,
    base_strength INTEGER DEFAULT 10,
    base_agility INTEGER DEFAULT 10,
    base_intellect INTEGER DEFAULT 10,
    base_stamina INTEGER DEFAULT 10,
    base_spirit INTEGER DEFAULT 10,
    base_threat_modifier REAL DEFAULT 1.0,
    combat_role TEXT DEFAULT 'dps',
    is_ranged INTEGER DEFAULT 0
);

-- 角色表
CREATE TABLE IF NOT EXISTS characters (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    name TEXT UNIQUE NOT NULL,
    race_id TEXT NOT NULL,
    class_id TEXT NOT NULL,
    faction TEXT NOT NULL,
    team_slot INTEGER NOT NULL,
    is_active INTEGER DEFAULT 1,
    is_dead INTEGER DEFAULT 0,
    revive_at DATETIME,
    level INTEGER DEFAULT 1,
    exp INTEGER DEFAULT 0,
    exp_to_next INTEGER DEFAULT 100,
    hp INTEGER DEFAULT 100,
    max_hp INTEGER DEFAULT 100,
    resource INTEGER DEFAULT 100,
    max_resource INTEGER DEFAULT 100,
    resource_type TEXT DEFAULT 'mana',
    strength INTEGER DEFAULT 10,
    agility INTEGER DEFAULT 10,
    intellect INTEGER DEFAULT 10,
    stamina INTEGER DEFAULT 10,
    spirit INTEGER DEFAULT 10,
    attack INTEGER DEFAULT 5,
    defense INTEGER DEFAULT 5,
    crit_rate REAL DEFAULT 0.05,
    crit_damage REAL DEFAULT 1.5,
    total_kills INTEGER DEFAULT 0,
    total_deaths INTEGER DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (race_id) REFERENCES races(id),
    FOREIGN KEY (class_id) REFERENCES classes(id)
);

-- 区域表 (使用正确的列名: exp_modifier, gold_modifier)
CREATE TABLE IF NOT EXISTS zones (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT,
    min_level INTEGER DEFAULT 1,
    max_level INTEGER DEFAULT 60,
    faction TEXT,
    exp_modifier REAL DEFAULT 1.0,
    gold_modifier REAL DEFAULT 1.0
);

-- 怪物表
CREATE TABLE IF NOT EXISTS monsters (
    id TEXT PRIMARY KEY,
    zone_id TEXT NOT NULL,
    name TEXT NOT NULL,
    level INTEGER DEFAULT 1,
    type TEXT DEFAULT 'normal',
    hp INTEGER DEFAULT 50,
    attack INTEGER DEFAULT 10,
    defense INTEGER DEFAULT 5,
    exp_reward INTEGER DEFAULT 20,
    gold_min INTEGER DEFAULT 1,
    gold_max INTEGER DEFAULT 10,
    spawn_weight INTEGER DEFAULT 100,
    FOREIGN KEY (zone_id) REFERENCES zones(id)
);

-- 插入测试数据
INSERT INTO races (id, name, faction, description, strength_base, agility_base, intellect_base, stamina_base, spirit_base, strength_pct, agility_pct, intellect_pct, stamina_pct, spirit_pct)
VALUES 
    ('human', '人类', 'alliance', '艾泽拉斯的人类', 2, 0, 0, 0, 2, 0, 0, 0, 0, 0),
    ('orc', '兽人', 'horde', '德拉诺的兽人', 3, 0, 0, 1, 0, 0, 0, 0, 0, 0);

INSERT INTO classes (id, name, role, description, primary_stat, resource_type, base_hp, base_resource, hp_per_level, resource_per_level, resource_regen, resource_regen_pct, base_strength, base_agility, base_intellect, base_stamina, base_spirit, base_threat_modifier, combat_role, is_ranged)
VALUES 
    ('warrior', '战士', 'tank', '护甲坚固的近战职业', 'strength', 'rage', 120, 100, 15, 10, 5.0, 0.0, 15, 10, 5, 12, 8, 1.0, 'tank', 0),
    ('mage', '法师', 'dps', '强大的奥术施法者', 'intellect', 'mana', 80, 200, 8, 15, 2.0, 0.05, 5, 8, 18, 8, 12, 0.7, 'dps', 1);

INSERT INTO zones (id, name, description, min_level, max_level, faction, exp_modifier, gold_modifier)
VALUES 
    ('elwynn', '艾尔文森林', '暴风城周边的和平森林', 1, 10, 'alliance', 1.0, 1.0),
    ('durotar', '杜隆塔尔', '兽人的荒芜家园', 1, 10, 'horde', 1.0, 1.0);

INSERT INTO monsters (id, zone_id, name, level, type, hp, attack, defense, exp_reward, gold_min, gold_max, spawn_weight)
VALUES 
    ('wolf', 'elwynn', '森林狼', 2, 'normal', 30, 8, 2, 15, 1, 5, 100),
    ('kobold', 'elwynn', '狗头人', 3, 'normal', 40, 10, 3, 20, 2, 8, 80),
    ('boar', 'durotar', '野猪', 2, 'normal', 35, 9, 3, 18, 2, 6, 100);
`
}


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
		// 如果找不到文件，使用内联的基础schema
		schema = []byte(getBaseSchema())
	}

	_, err = testDB.Exec(string(schema))
	if err != nil {
		testDB.Close()
		return nil, err
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

-- 种族表
CREATE TABLE IF NOT EXISTS races (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    name_en TEXT NOT NULL,
    faction TEXT NOT NULL,
    description TEXT,
    racial_ability TEXT,
    strength_base INTEGER DEFAULT 0,
    agility_base INTEGER DEFAULT 0,
    intellect_base INTEGER DEFAULT 0,
    stamina_base INTEGER DEFAULT 0,
    spirit_base INTEGER DEFAULT 0
);

-- 职业表
CREATE TABLE IF NOT EXISTS classes (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    name_en TEXT NOT NULL,
    role TEXT NOT NULL,
    description TEXT,
    resource_type TEXT DEFAULT 'mana',
    base_hp INTEGER DEFAULT 100,
    base_resource INTEGER DEFAULT 100,
    base_strength INTEGER DEFAULT 10,
    base_agility INTEGER DEFAULT 10,
    base_intellect INTEGER DEFAULT 10,
    base_stamina INTEGER DEFAULT 10,
    base_spirit INTEGER DEFAULT 10
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
    is_active BOOLEAN DEFAULT 1,
    is_dead BOOLEAN DEFAULT 0,
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
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (race_id) REFERENCES races(id),
    FOREIGN KEY (class_id) REFERENCES classes(id)
);

-- 区域表
CREATE TABLE IF NOT EXISTS zones (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    name_en TEXT NOT NULL,
    description TEXT,
    min_level INTEGER DEFAULT 1,
    max_level INTEGER DEFAULT 60,
    faction TEXT
);

-- 插入测试数据
INSERT INTO races (id, name, name_en, faction, description, strength_base, agility_base, intellect_base, stamina_base, spirit_base)
VALUES 
    ('human', '人类', 'Human', 'alliance', '艾泽拉斯的人类', 2, 0, 0, 0, 2),
    ('orc', '兽人', 'Orc', 'horde', '德拉诺的兽人', 3, 0, 0, 1, 0);

INSERT INTO classes (id, name, name_en, role, description, resource_type, base_hp, base_resource, base_strength, base_agility, base_intellect, base_stamina, base_spirit)
VALUES 
    ('warrior', '战士', 'Warrior', 'tank', '护甲坚固的近战职业', 'rage', 120, 100, 15, 10, 5, 12, 8),
    ('mage', '法师', 'Mage', 'dps', '强大的奥术施法者', 'mana', 80, 200, 5, 8, 18, 8, 12);

INSERT INTO zones (id, name, name_en, description, min_level, max_level, faction)
VALUES 
    ('elwynn_forest', '艾尔文森林', 'Elwynn Forest', '暴风城周边的和平森林', 1, 10, 'alliance'),
    ('durotar', '杜隆塔尔', 'Durotar', '兽人的荒芜家园', 1, 10, 'horde');
`
}


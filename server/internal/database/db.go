package database

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

// debugLog 只在 TEST_DEBUG=1 时输出日志，避免序列化错误
func debugLog(format string, args ...interface{}) {
	if os.Getenv("TEST_DEBUG") == "1" || os.Getenv("TEST_DEBUG") == "true" {
		log.Printf(format, args...)
	}
}

// Init 初始化数据库连接
func Init() error {
	var err error

	// 打开数据库连接，启用WAL模式
	DB, err = sql.Open("sqlite3", "./game.db?_journal_mode=WAL&_busy_timeout=5000")
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	// 设置连接池
	DB.SetMaxOpenConns(25)
	DB.SetMaxIdleConns(5)

	// 启用外键约束
	if _, err := DB.Exec("PRAGMA foreign_keys = ON"); err != nil {
		return fmt.Errorf("failed to enable foreign keys: %w", err)
	}

	// 创建表结构
	if err := initSchema(); err != nil {
		return fmt.Errorf("failed to init schema: %w", err)
	}

	// 运行数据库迁移
	if err := runMigrations(); err != nil {
		debugLog("Migration warning: %v", err)
	}

	// 导入种子数据（如果需要）
	if err := seedData(); err != nil {
		debugLog("Seed data warning: %v", err)
	}

	debugLog("Database initialized successfully")
	return nil
}

// Close 关闭数据库连接
func Close() {
	if DB != nil {
		DB.Close()
	}
}

// initSchema 初始化数据库表结构
func initSchema() error {
	// 尝试从文件加载schema
	// 尝试多个可能的路径
	possiblePaths := []string{
		filepath.Join("database", "schema.sql"),
		filepath.Join("server", "database", "schema.sql"),
		filepath.Join("..", "database", "schema.sql"),
		filepath.Join("..", "..", "database", "schema.sql"),
	}
	
	var schemaPath string
	var found bool
	for _, path := range possiblePaths {
		if _, err := os.Stat(path); err == nil {
			schemaPath = path
			found = true
			break
		}
	}
	
	if found {
		content, err := ioutil.ReadFile(schemaPath)
		if err != nil {
			return fmt.Errorf("failed to read schema file: %w", err)
		}
		if _, err := DB.Exec(string(content)); err != nil {
			return fmt.Errorf("failed to execute schema: %w", err)
		}
		debugLog("Schema loaded from file")
		return nil
	}

	// 如果没有schema文件，使用内置的基础schema
	debugLog("Using embedded schema")
	return createBasicTables()
}

// createBasicTables 创建基础表（备用方案）
func createBasicTables() error {
	schema := `
	-- 用户表
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username VARCHAR(32) UNIQUE NOT NULL,
		password_hash VARCHAR(256) NOT NULL,
		email VARCHAR(128) UNIQUE,
		max_team_size INTEGER DEFAULT 5,
		unlocked_slots INTEGER DEFAULT 1,
		gold INTEGER DEFAULT 0,
		current_zone_id VARCHAR(32) DEFAULT 'elwynn',
		total_kills INTEGER DEFAULT 0,
		total_gold_gained INTEGER DEFAULT 0,
		play_time INTEGER DEFAULT 0,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		last_login_at DATETIME,
		status INTEGER DEFAULT 1
	);

	CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);

	-- 角色表
	CREATE TABLE IF NOT EXISTS characters (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		name VARCHAR(32) NOT NULL,
		race_id VARCHAR(32) NOT NULL,
		class_id VARCHAR(32) NOT NULL,
		faction VARCHAR(16) NOT NULL,
		team_slot INTEGER NOT NULL,
		is_active INTEGER DEFAULT 1,
		is_dead INTEGER DEFAULT 0,
		revive_at DATETIME,
		level INTEGER DEFAULT 1,
		exp INTEGER DEFAULT 0,
		exp_to_next INTEGER DEFAULT 100,
		hp INTEGER NOT NULL,
		max_hp INTEGER NOT NULL,
		resource INTEGER NOT NULL,
		max_resource INTEGER NOT NULL,
		resource_type VARCHAR(16) NOT NULL,
		strength INTEGER DEFAULT 10,
		agility INTEGER DEFAULT 10,
		intellect INTEGER DEFAULT 10,
		stamina INTEGER DEFAULT 10,
		spirit INTEGER DEFAULT 10,
		unspent_points INTEGER DEFAULT 0,
		physical_attack INTEGER DEFAULT 10,
		magic_attack INTEGER DEFAULT 10,
		physical_defense INTEGER DEFAULT 5,
		magic_defense INTEGER DEFAULT 5,
		phys_crit_rate REAL DEFAULT 0.05,
		phys_crit_damage REAL DEFAULT 1.5,
		spell_crit_rate REAL DEFAULT 0.05,
		spell_crit_damage REAL DEFAULT 1.5,
		dodge_rate REAL DEFAULT 0.05,
		total_kills INTEGER DEFAULT 0,
		total_deaths INTEGER DEFAULT 0,
		total_damage_dealt INTEGER DEFAULT 0,
		total_healing_done INTEGER DEFAULT 0,
		offline_time DATETIME,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
		UNIQUE(user_id, team_slot)
	);

	CREATE INDEX IF NOT EXISTS idx_characters_user_id ON characters(user_id);
	CREATE INDEX IF NOT EXISTS idx_characters_level ON characters(level);
	`
	_, err := DB.Exec(schema)
	return err
}

// seedData 导入种子数据
func seedData() error {
	// 检查是否已有数据
	var count int
	err := DB.QueryRow("SELECT COUNT(*) FROM races").Scan(&count)
	if err == nil && count > 0 {
		debugLog("Seed data already exists")
		// 即使已有基础数据，也检查是否需要加载战士技能数据
		var skillCount int
		err := DB.QueryRow("SELECT COUNT(*) FROM skills WHERE class_id = 'warrior'").Scan(&skillCount)
		if err == nil && skillCount == 0 {
			debugLog("Warrior skills not found, loading...")
			if err := loadWarriorSkills(); err != nil {
				debugLog("Failed to load warrior skills: %v", err)
			}
		}
		return nil
	}

	// 尝试从文件加载seed
	seedPath := filepath.Join("database", "seed.sql")
	if _, err := os.Stat(seedPath); err == nil {
		content, err := ioutil.ReadFile(seedPath)
		if err != nil {
			return fmt.Errorf("failed to read seed file: %w", err)
		}
		if _, err := DB.Exec(string(content)); err != nil {
			return fmt.Errorf("failed to execute seed: %w", err)
		}
		debugLog("Seed data loaded from file")
	}

	// 加载战士技能数据
	if err := loadWarriorSkills(); err != nil {
		debugLog("Failed to load warrior skills: %v", err)
		// 不返回错误，因为技能数据可能已经存在
	}

	return nil
}

// loadWarriorSkills 加载战士技能数据
func loadWarriorSkills() error {
	warriorSkillsPath := filepath.Join("database", "warrior_skills.sql")
	if _, err := os.Stat(warriorSkillsPath); err != nil {
		return fmt.Errorf("warrior_skills.sql not found: %w", err)
	}

	content, err := ioutil.ReadFile(warriorSkillsPath)
	if err != nil {
		return fmt.Errorf("failed to read warrior_skills.sql: %w", err)
	}

	// 执行SQL文件
	if _, err := DB.Exec(string(content)); err != nil {
		return fmt.Errorf("failed to execute warrior_skills.sql: %w", err)
	}

	debugLog("Warrior skills loaded")
	return nil
}

// Transaction 执行事务
func Transaction(fn func(*sql.Tx) error) error {
	tx, err := DB.Begin()
	if err != nil {
		return err
	}

	if err := fn(tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

// runMigrations 运行数据库迁移
func runMigrations() error {
	// 迁移0: 添加dodge_rate列到characters表
	if err := migrateDodgeRate(); err != nil {
		return fmt.Errorf("failed to migrate dodge_rate: %w", err)
	}
	// 迁移1: 更新 battle_strategies 表结构
	if err := migrateBattleStrategies(); err != nil {
		return fmt.Errorf("failed to migrate battle_strategies: %w", err)
	}
	// 迁移2: 创建怪物强度配置表
	if err := migrateMonsterStrengthConfig(); err != nil {
		return fmt.Errorf("failed to migrate monster_strength_config: %w", err)
	}
	// 迁移3: 创建配置版本表
	if err := migrateConfigVersions(); err != nil {
		return fmt.Errorf("failed to migrate config_versions: %w", err)
	}
	return nil
}

// migrateDodgeRate 添加dodge_rate列到characters表
func migrateDodgeRate() error {
	// 检查列是否已存在
	rows, err := DB.Query("PRAGMA table_info(characters)")
	if err != nil {
		return err
	}
	defer rows.Close()

	hasDodgeRate := false
	for rows.Next() {
		var cid int
		var name, colType string
		var notNull, pk int
		var dfltValue sql.NullString
		if err := rows.Scan(&cid, &name, &colType, &notNull, &dfltValue, &pk); err != nil {
			return err
		}
		if name == "dodge_rate" {
			hasDodgeRate = true
			break
		}
	}

	if !hasDodgeRate {
		debugLog("Adding dodge_rate column to characters table...")
		_, err := DB.Exec("ALTER TABLE characters ADD COLUMN dodge_rate REAL DEFAULT 0.05")
		if err != nil {
			return fmt.Errorf("failed to add dodge_rate column: %w", err)
		}
		debugLog("dodge_rate column added successfully")
	}

	return nil
}

// migrateBattleStrategies 迁移 battle_strategies 表
func migrateBattleStrategies() error {
	// 检查表是否存在
	var tableName string
	err := DB.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name='battle_strategies'").Scan(&tableName)
	if err == sql.ErrNoRows {
		// 表不存在，将由 schema.sql 创建
		return nil
	}

	// 检查是否有 skill_priority 列
	rows, err := DB.Query("PRAGMA table_info(battle_strategies)")
	if err != nil {
		return err
	}
	defer rows.Close()

	hasSkillPriority := false
	hasSkillTargetOverrides := false
	hasAutoTargetSettings := false

	for rows.Next() {
		var cid int
		var name, colType string
		var notNull, pk int
		var dfltValue sql.NullString
		if err := rows.Scan(&cid, &name, &colType, &notNull, &dfltValue, &pk); err != nil {
			return err
		}
		switch name {
		case "skill_priority":
			hasSkillPriority = true
		case "skill_target_overrides":
			hasSkillTargetOverrides = true
		case "auto_target_settings":
			hasAutoTargetSettings = true
		}
	}

	// 如果缺少新列，需要重建表
	if !hasSkillPriority || !hasSkillTargetOverrides || !hasAutoTargetSettings {
		debugLog("Migrating battle_strategies table...")

		// SQLite 不支持 DROP COLUMN，需要重建表
		// 1. 重命名旧表
		_, err := DB.Exec("ALTER TABLE battle_strategies RENAME TO battle_strategies_old")
		if err != nil {
			return fmt.Errorf("failed to rename old table: %w", err)
		}

		// 2. 创建新表
		_, err = DB.Exec(`
			CREATE TABLE battle_strategies (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				character_id INTEGER NOT NULL,
				name VARCHAR(32) NOT NULL,
				is_active INTEGER DEFAULT 0,
				skill_priority TEXT,
				conditional_rules TEXT,
				target_priority VARCHAR(32) DEFAULT 'lowest_hp',
				skill_target_overrides TEXT,
				resource_threshold INTEGER DEFAULT 0,
				reserved_skills TEXT,
				auto_target_settings TEXT,
				created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
				updated_at DATETIME,
				FOREIGN KEY (character_id) REFERENCES characters(id) ON DELETE CASCADE
			)
		`)
		if err != nil {
			// 回滚：重命名回来
			DB.Exec("ALTER TABLE battle_strategies_old RENAME TO battle_strategies")
			return fmt.Errorf("failed to create new table: %w", err)
		}

		// 3. 迁移数据（如果有）
		_, err = DB.Exec(`
			INSERT INTO battle_strategies (id, character_id, name, is_active, created_at)
			SELECT id, character_id, name, is_active, created_at 
			FROM battle_strategies_old
		`)
		if err != nil {
			debugLog("Failed to migrate data: %v", err)
			// 不是致命错误，继续
		}

		// 4. 删除旧表
		_, err = DB.Exec("DROP TABLE battle_strategies_old")
		if err != nil {
			debugLog("Failed to drop old table: %v", err)
		}

		// 5. 创建索引
		DB.Exec("CREATE INDEX IF NOT EXISTS idx_battle_strategies_character ON battle_strategies(character_id)")
		DB.Exec("CREATE INDEX IF NOT EXISTS idx_battle_strategies_active ON battle_strategies(character_id, is_active)")

		debugLog("battle_strategies table migrated successfully")
	}

	return nil
}

// migrateMonsterStrengthConfig 迁移怪物强度配置表
func migrateMonsterStrengthConfig() error {
	// 检查表是否已存在
	var tableName string
	err := DB.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name='monster_strength_config'").Scan(&tableName)
	if err == nil {
		// 表已存在，跳过
		return nil
	}
	if err != sql.ErrNoRows {
		return err
	}

	// 创建配置表
	debugLog("Creating monster_strength_config table...")
	_, err = DB.Exec(`
		CREATE TABLE IF NOT EXISTS monster_strength_config (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			level_min INTEGER NOT NULL,
			level_max INTEGER NOT NULL,
			hp_multiplier REAL DEFAULT 1.0,
			attack_multiplier REAL DEFAULT 1.0,
			defense_multiplier REAL DEFAULT 1.0,
			crit_rate_bonus REAL DEFAULT 0.0,
			description TEXT,
			is_active INTEGER DEFAULT 1,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			UNIQUE(level_min, level_max)
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create monster_strength_config table: %w", err)
	}

	// 创建索引
	DB.Exec("CREATE INDEX IF NOT EXISTS idx_strength_config_level ON monster_strength_config(level_min, level_max)")
	DB.Exec("CREATE INDEX IF NOT EXISTS idx_strength_config_active ON monster_strength_config(is_active)")

	// 插入默认配置
	_, err = DB.Exec(`
		INSERT OR IGNORE INTO monster_strength_config 
		(level_min, level_max, hp_multiplier, attack_multiplier, defense_multiplier, crit_rate_bonus, description) 
		VALUES
		(1, 10, 1.5, 1.4, 1.4, 0.02, '1-10级：HP +50%, 攻击 +40%, 防御 +40%, 暴击率 +2%'),
		(11, 20, 1.45, 1.35, 1.35, 0.02, '10-20级：HP +45%, 攻击 +35%, 防御 +35%, 暴击率 +2%'),
		(21, 30, 1.4, 1.35, 1.35, 0.02, '20-30级：HP +40%, 攻击 +35%, 防御 +35%, 暴击率 +2%'),
		(31, 40, 1.4, 1.35, 1.35, 0.03, '30-40级：HP +40%, 攻击 +35%, 防御 +35%, 暴击率 +3%'),
		(41, 50, 1.35, 1.3, 1.3, 0.03, '40-50级：HP +35%, 攻击 +30%, 防御 +30%, 暴击率 +3%'),
		(51, 60, 1.35, 1.3, 1.3, 0.03, '50-60级：HP +35%, 攻击 +30%, 防御 +30%, 暴击率 +3%')
	`)
	if err != nil {
		debugLog("Failed to insert default configs: %v", err)
	}

	debugLog("monster_strength_config table created successfully")
	return nil
}

// migrateConfigVersions 迁移配置版本表
func migrateConfigVersions() error {
	// 检查表是否已存在
	var tableName string
	err := DB.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name='config_versions'").Scan(&tableName)
	if err == nil {
		// 表已存在，跳过
		return nil
	}
	if err != sql.ErrNoRows {
		return err
	}

	// 创建配置版本表
	debugLog("Creating config_versions table...")
	_, err = DB.Exec(`
		CREATE TABLE IF NOT EXISTS config_versions (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			config_type VARCHAR(32) NOT NULL,
			version INTEGER NOT NULL,
			config_data TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			created_by VARCHAR(32),
			description TEXT,
			UNIQUE(config_type, version)
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create config_versions table: %w", err)
	}

	// 创建索引
	DB.Exec("CREATE INDEX IF NOT EXISTS idx_config_versions_type ON config_versions(config_type)")
	DB.Exec("CREATE INDEX IF NOT EXISTS idx_config_versions_version ON config_versions(config_type, version DESC)")

	debugLog("config_versions table created successfully")
	return nil
}

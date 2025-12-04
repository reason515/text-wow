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

	// 导入种子数据（如果需要）
	if err := seedData(); err != nil {
		log.Printf("⚠️ Seed data warning: %v", err)
	}

	log.Println("✅ 数据库初始化完成")
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
	schemaPath := filepath.Join("database", "schema.sql")
	if _, err := os.Stat(schemaPath); err == nil {
		content, err := ioutil.ReadFile(schemaPath)
		if err != nil {
			return fmt.Errorf("failed to read schema file: %w", err)
		}
		if _, err := DB.Exec(string(content)); err != nil {
			return fmt.Errorf("failed to execute schema: %w", err)
		}
		log.Println("📄 Schema loaded from file")
		return nil
	}

	// 如果没有schema文件，使用内置的基础schema
	log.Println("📄 Using embedded schema")
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
		physical_attack INTEGER DEFAULT 10,
		magic_attack INTEGER DEFAULT 10,
		physical_defense INTEGER DEFAULT 5,
		magic_defense INTEGER DEFAULT 5,
		crit_rate REAL DEFAULT 0.05,
		crit_damage REAL DEFAULT 1.5,
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
		log.Println("📊 Seed data already exists")
		// 即使已有基础数据，也检查是否需要加载战士技能数据
		var skillCount int
		err := DB.QueryRow("SELECT COUNT(*) FROM skills WHERE class_id = 'warrior'").Scan(&skillCount)
		if err == nil && skillCount == 0 {
			log.Println("⚠️ Warrior skills not found, loading...")
			if err := loadWarriorSkills(); err != nil {
				log.Printf("⚠️ Failed to load warrior skills: %v", err)
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
		log.Println("🌱 Seed data loaded from file")
	}

	// 加载战士技能数据
	if err := loadWarriorSkills(); err != nil {
		log.Printf("⚠️ Failed to load warrior skills: %v", err)
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

	log.Println("⚔️ Warrior skills loaded")
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

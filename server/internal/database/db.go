package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func Init() error {
	var err error
	DB, err = sql.Open("sqlite3", "./game.db?_journal_mode=WAL")
	if err != nil {
		return err
	}

	// 创建表
	if err := createTables(); err != nil {
		return err
	}

	log.Println("✅ 数据库初始化完成")
	return nil
}

func Close() {
	if DB != nil {
		DB.Close()
	}
}

func createTables() error {
	schema := `
	CREATE TABLE IF NOT EXISTS characters (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		faction TEXT NOT NULL,
		race TEXT NOT NULL,
		class TEXT NOT NULL,
		level INTEGER DEFAULT 1,
		exp INTEGER DEFAULT 0,
		exp_to_next INTEGER DEFAULT 100,
		hp INTEGER DEFAULT 100,
		max_hp INTEGER DEFAULT 100,
		mp INTEGER DEFAULT 50,
		max_mp INTEGER DEFAULT 50,
		strength INTEGER DEFAULT 10,
		agility INTEGER DEFAULT 10,
		intellect INTEGER DEFAULT 10,
		stamina INTEGER DEFAULT 10,
		spirit INTEGER DEFAULT 10,
		gold INTEGER DEFAULT 0,
		current_zone TEXT DEFAULT 'elwynn_forest',
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS strategies (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		character_id INTEGER NOT NULL,
		skill_priority TEXT DEFAULT '["attack"]',
		hp_potion_threshold INTEGER DEFAULT 30,
		mp_potion_threshold INTEGER DEFAULT 20,
		target_priority TEXT DEFAULT 'lowest_hp',
		auto_loot INTEGER DEFAULT 1,
		FOREIGN KEY (character_id) REFERENCES characters(id)
	);

	CREATE TABLE IF NOT EXISTS battle_logs (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		character_id INTEGER NOT NULL,
		message TEXT NOT NULL,
		log_type TEXT DEFAULT 'combat',
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (character_id) REFERENCES characters(id)
	);
	`
	_, err := DB.Exec(schema)
	return err
}

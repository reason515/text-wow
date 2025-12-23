package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// 获取当前工作目录，然后找到 server 目录
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get working directory: %v", err)
	}

	// 如果从 cmd/init_database 运行，需要回到 server 目录
	serverDir := wd
	if filepath.Base(wd) == "init_database" {
		serverDir = filepath.Join(wd, "..", "..")
	} else if filepath.Base(filepath.Dir(wd)) == "cmd" {
		serverDir = filepath.Join(wd, "..")
	}

	// 打开数据库
	dbPath := filepath.Join(serverDir, "game.db")
	db, err := sql.Open("sqlite3", dbPath+"?_journal_mode=WAL&_busy_timeout=5000")
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	// 启用外键约束
	if _, err := db.Exec("PRAGMA foreign_keys = ON"); err != nil {
		log.Fatalf("Failed to enable foreign keys: %v", err)
	}

	// 1. 加载 schema.sql
	fmt.Println("Loading schema.sql...")
	schemaPath := filepath.Join(serverDir, "database", "schema.sql")
	schemaContent, err := ioutil.ReadFile(schemaPath)
	if err != nil {
		log.Fatalf("Failed to read schema file: %v", err)
	}
	if _, err := db.Exec(string(schemaContent)); err != nil {
		log.Printf("Warning: Schema execution had errors (tables may already exist): %v", err)
	} else {
		fmt.Println("✓ Schema loaded")
	}

	// 2. 加载 seed.sql（基础数据）
	fmt.Println("\nLoading seed.sql...")
	seedPath := filepath.Join(serverDir, "database", "seed.sql")
	seedContent, err := ioutil.ReadFile(seedPath)
	if err != nil {
		log.Fatalf("Failed to read seed file: %v", err)
	}
	if _, err := db.Exec(string(seedContent)); err != nil {
		log.Printf("Warning: Seed execution had errors: %v", err)
		log.Println("Continuing with other data...")
	} else {
		fmt.Println("✓ Seed data loaded")
	}

	// 3. 加载战士技能数据
	fmt.Println("\nLoading warrior_skills.sql...")
	warriorSkillsPath := filepath.Join(serverDir, "database", "warrior_skills.sql")
	warriorSkillsContent, err := ioutil.ReadFile(warriorSkillsPath)
	if err != nil {
		log.Printf("Warning: warrior_skills.sql not found: %v", err)
	} else {
		if _, err := db.Exec(string(warriorSkillsContent)); err != nil {
			log.Printf("Warning: Warrior skills execution had errors: %v", err)
		} else {
			fmt.Println("✓ Warrior skills loaded")
		}
	}

	// 4. 添加杜隆塔尔怪物数据（因为seed.sql可能有语法问题）
	fmt.Println("\nAdding Durotar monsters...")
	durotarMonsters := []struct {
		id, zoneID, name, monsterType, attackType string
		level, hp, physAtk, magicAtk, physDef, magicDef int
		physCritRate, physCritDmg, spellCritRate, spellCritDmg float64
		expReward, goldMin, goldMax, spawnWeight int
	}{
		{"scorpid_durotar", "durotar", "蝎子", "normal", "physical", 1, 22, 4, 0, 1, 1, 0.07, 1.5, 0.07, 1.5, 5, 1, 2, 100},
		{"dire_wolf_durotar", "durotar", "恐狼", "normal", "physical", 1, 18, 4, 0, 1, 1, 0.07, 1.5, 0.07, 1.5, 4, 1, 1, 100},
		{"razormane_scout_durotar", "durotar", "钢鬃斥候", "normal", "physical", 2, 27, 6, 0, 1, 1, 0.07, 1.5, 0.07, 1.5, 6, 1, 2, 80},
		{"razormane_warrior_durotar", "durotar", "钢鬃战士", "normal", "physical", 3, 33, 7, 0, 3, 3, 0.07, 1.5, 0.07, 1.5, 7, 1, 3, 60},
		{"razormane_thornweaver_durotar", "durotar", "钢鬃织棘者", "normal", "physical", 4, 39, 8, 0, 3, 3, 0.07, 1.5, 0.07, 1.5, 8, 2, 3, 50},
		{"razormane_battleboar_durotar", "durotar", "钢鬃战猪", "normal", "physical", 5, 45, 10, 0, 4, 4, 0.07, 1.5, 0.07, 1.5, 10, 2, 4, 40},
		{"venomtail_scorpid_durotar", "durotar", "毒尾蝎", "normal", "physical", 3, 30, 7, 0, 3, 3, 0.07, 1.5, 0.07, 1.5, 7, 1, 3, 70},
		{"elder_scorpid_durotar", "durotar", "老蝎子", "normal", "physical", 5, 42, 10, 0, 4, 4, 0.07, 1.5, 0.07, 1.5, 9, 2, 4, 30},
		{"razormane_geomancer_durotar", "durotar", "钢鬃地卜师", "normal", "magic", 6, 45, 6, 14, 3, 4, 0.07, 1.5, 0.10, 1.6, 12, 2, 5, 35},
		{"razormane_champion_durotar", "durotar", "钢鬃勇士", "normal", "physical", 6, 48, 11, 0, 4, 4, 0.07, 1.5, 0.07, 1.5, 11, 2, 5, 25},
		{"captain_flat_tusk", "durotar", "平牙队长", "elite", "physical", 8, 120, 17, 0, 7, 7, 0.10, 1.6, 0.07, 1.5, 35, 5, 12, 5},
	}

	for _, m := range durotarMonsters {
		_, err := db.Exec(`
			INSERT OR REPLACE INTO monsters (
				id, zone_id, name, level, type, hp, physical_attack, magic_attack, 
				physical_defense, magic_defense, attack_type, 
				phys_crit_rate, phys_crit_damage, spell_crit_rate, spell_crit_damage,
				exp_reward, gold_min, gold_max, spawn_weight
			) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		`, m.id, m.zoneID, m.name, m.level, m.monsterType, m.hp, m.physAtk, m.magicAtk,
			m.physDef, m.magicDef, m.attackType,
			m.physCritRate, m.physCritDmg, m.spellCritRate, m.spellCritDmg,
			m.expReward, m.goldMin, m.goldMax, m.spawnWeight)
		if err != nil {
			log.Printf("Failed to insert %s: %v", m.id, err)
		}
	}
	fmt.Println("✓ Durotar monsters added")

	// 验证
	var monsterCount int
	err = db.QueryRow("SELECT COUNT(*) FROM monsters WHERE zone_id = 'durotar'").Scan(&monsterCount)
	if err == nil {
		fmt.Printf("\n✅ Database initialized successfully!\n")
		fmt.Printf("   Durotar monsters: %d\n", monsterCount)
		
		var skillCount int
		err = db.QueryRow("SELECT COUNT(*) FROM skills WHERE class_id = 'warrior'").Scan(&skillCount)
		if err == nil {
			fmt.Printf("   Warrior skills: %d\n", skillCount)
		}
	}
}




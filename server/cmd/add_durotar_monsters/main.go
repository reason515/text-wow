package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// 打开数据库
	db, err := sql.Open("sqlite3", "../../game.db?_journal_mode=WAL&_busy_timeout=5000")
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	// 启用外键约束
	if _, err := db.Exec("PRAGMA foreign_keys = ON"); err != nil {
		log.Fatalf("Failed to enable foreign keys: %v", err)
	}

	// 插入杜隆塔尔怪物数据
	monsters := []struct {
		id, zoneID, name, monsterType, attackType string
		level, hp, physAtk, magicAtk, physDef, magicDef int
		physCritRate, physCritDmg, spellCritRate, spellCritDmg float64
		expReward, goldMin, goldMax, spawnWeight int
	}{
		{"scorpid", "durotar", "蝎子", "normal", "physical", 1, 22, 4, 0, 1, 1, 0.07, 1.5, 0.07, 1.5, 5, 1, 2, 100},
		{"dire_wolf", "durotar", "恐狼", "normal", "physical", 1, 18, 4, 0, 1, 1, 0.07, 1.5, 0.07, 1.5, 4, 1, 1, 100},
		{"razormane_scout", "durotar", "钢鬃斥候", "normal", "physical", 2, 27, 6, 0, 1, 1, 0.07, 1.5, 0.07, 1.5, 6, 1, 2, 80},
		{"razormane_warrior", "durotar", "钢鬃战士", "normal", "physical", 3, 33, 7, 0, 3, 3, 0.07, 1.5, 0.07, 1.5, 7, 1, 3, 60},
		{"razormane_thornweaver", "durotar", "钢鬃织棘者", "normal", "physical", 4, 39, 8, 0, 3, 3, 0.07, 1.5, 0.07, 1.5, 8, 2, 3, 50},
		{"razormane_battleboar", "durotar", "钢鬃战猪", "normal", "physical", 5, 45, 10, 0, 4, 4, 0.07, 1.5, 0.07, 1.5, 10, 2, 4, 40},
		{"venomtail_scorpid", "durotar", "毒尾蝎", "normal", "physical", 3, 30, 7, 0, 3, 3, 0.07, 1.5, 0.07, 1.5, 7, 1, 3, 70},
		{"elder_scorpid", "durotar", "老蝎子", "normal", "physical", 5, 42, 10, 0, 4, 4, 0.07, 1.5, 0.07, 1.5, 9, 2, 4, 30},
		{"razormane_geomancer", "durotar", "钢鬃地卜师", "normal", "magic", 6, 45, 6, 14, 3, 4, 0.07, 1.5, 0.10, 1.6, 12, 2, 5, 35},
		{"razormane_champion", "durotar", "钢鬃勇士", "normal", "physical", 6, 48, 11, 0, 4, 4, 0.07, 1.5, 0.07, 1.5, 11, 2, 5, 25},
		{"captain_flat_tusk", "durotar", "平牙队长", "elite", "physical", 8, 120, 17, 0, 7, 7, 0.10, 1.6, 0.07, 1.5, 35, 5, 12, 5},
	}

	fmt.Println("正在插入杜隆塔尔怪物数据...")
	for _, m := range monsters {
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
		} else {
			fmt.Printf("  ✓ %s\n", m.name)
		}
	}

	// 验证
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM monsters WHERE zone_id = 'durotar'").Scan(&count)
	if err != nil {
		log.Fatalf("Failed to verify: %v", err)
	}
	fmt.Printf("\n✅ 完成！杜隆塔尔区域现在有 %d 个怪物\n", count)
}


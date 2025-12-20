package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// 打开数据库（从server目录查找）
	db, err := sql.Open("sqlite3", "../../game.db?_journal_mode=WAL&_busy_timeout=5000")
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	// 启用外键约束
	if _, err := db.Exec("PRAGMA foreign_keys = ON"); err != nil {
		log.Fatalf("Failed to enable foreign keys: %v", err)
	}

	// 更新冲锋技能的冷却时间
	fmt.Println("Updating warrior_charge cooldown from 3 to 5...")
	result, err := db.Exec(`
		UPDATE skills
		SET cooldown = 5,
		    description = '战士的招牌技能！快速冲向敌人并造成伤害，获得怒气，有概率眩晕目标。1级：80%伤害，+15怒气，30%眩晕，冷却5回合，每级+10%伤害，+3怒气，+5%眩晕概率，-0.5回合冷却'
		WHERE id = 'warrior_charge'
	`)
	if err != nil {
		log.Fatalf("Failed to update warrior_charge: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Fatalf("Failed to get rows affected: %v", err)
	}

	if rowsAffected > 0 {
		fmt.Printf("✅ Successfully updated warrior_charge cooldown (affected %d row(s))\n", rowsAffected)
	} else {
		fmt.Println("⚠️ No rows updated. The skill 'warrior_charge' may not exist in the database.")
	}
}


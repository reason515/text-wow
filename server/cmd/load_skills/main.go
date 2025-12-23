package main

import (
	"database/sql"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	var (
		dbPath = flag.String("db", "../../game.db", "数据库路径")
		check  = flag.Bool("check", false, "只检查技能数据是否存在")
		force  = flag.Bool("force", false, "强制重新加载（删除后重新插入）")
	)
	flag.Parse()

	// 打开数据库
	db, err := sql.Open("sqlite3", *dbPath+"?_journal_mode=WAL&_busy_timeout=5000")
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	// 启用外键约束
	if _, err := db.Exec("PRAGMA foreign_keys = ON"); err != nil {
		log.Fatalf("Failed to enable foreign keys: %v", err)
	}

	// 检查技能数据
	var warriorSkillCount int
	err = db.QueryRow("SELECT COUNT(*) FROM skills WHERE class_id = 'warrior'").Scan(&warriorSkillCount)
	if err != nil {
		log.Printf("⚠️ 无法查询技能表，可能表不存在: %v", err)
		log.Println("请先运行数据库初始化")
		os.Exit(1)
	}

	fmt.Printf("当前数据库中的战士技能数量: %d\n", warriorSkillCount)

	if *check {
		if warriorSkillCount == 0 {
			fmt.Println("❌ 未找到战士技能数据")
			fmt.Println("请运行: go run main.go -load")
			os.Exit(1)
		} else {
			fmt.Println("✅ 技能数据已存在")
			os.Exit(0)
		}
	}

	// 加载技能数据
	warriorSkillsPath := filepath.Join("..", "..", "database", "warrior_skills.sql")
	if _, err := os.Stat(warriorSkillsPath); err != nil {
		log.Fatalf("❌ warrior_skills.sql 文件不存在: %v", err)
	}

	fmt.Printf("正在加载技能数据: %s\n", warriorSkillsPath)

	content, err := ioutil.ReadFile(warriorSkillsPath)
	if err != nil {
		log.Fatalf("❌ 无法读取文件: %v", err)
	}

	// 如果强制重新加载，先删除现有数据
	if *force && warriorSkillCount > 0 {
		fmt.Println("正在删除现有技能数据...")
		_, err = db.Exec("DELETE FROM skills WHERE class_id = 'warrior'")
		if err != nil {
			log.Printf("⚠️ 删除现有技能失败: %v", err)
		}
		_, err = db.Exec("DELETE FROM effects WHERE id LIKE 'eff_%'")
		if err != nil {
			log.Printf("⚠️ 删除现有效果失败: %v", err)
		}
		_, err = db.Exec("DELETE FROM passive_skills WHERE class_id = 'warrior'")
		if err != nil {
			log.Printf("⚠️ 删除现有被动技能失败: %v", err)
		}
	}

	// 执行SQL文件
	fmt.Println("正在执行SQL...")
	if _, err := db.Exec(string(content)); err != nil {
		log.Fatalf("❌ 执行SQL失败: %v", err)
	}

	// 验证加载结果
	var newCount int
	err = db.QueryRow("SELECT COUNT(*) FROM skills WHERE class_id = 'warrior'").Scan(&newCount)
	if err != nil {
		log.Fatalf("❌ 无法验证结果: %v", err)
	}

	fmt.Printf("✅ 技能数据加载完成！\n")
	fmt.Printf("   加载前: %d 个技能\n", warriorSkillCount)
	fmt.Printf("   加载后: %d 个技能\n", newCount)

	// 检查初始技能池
	initialSkillIDs := []string{
		"warrior_heroic_strike",
		"warrior_taunt",
		"warrior_shield_block",
		"warrior_cleave",
		"warrior_slam",
		"warrior_battle_shout",
		"warrior_demoralizing_shout",
		"warrior_last_stand",
		"warrior_charge",
	}

	fmt.Println("\n检查初始技能池:")
	missingCount := 0
	for _, skillID := range initialSkillIDs {
		var exists int
		err := db.QueryRow("SELECT COUNT(*) FROM skills WHERE id = ?", skillID).Scan(&exists)
		if err != nil || exists == 0 {
			fmt.Printf("  ❌ %s - 缺失\n", skillID)
			missingCount++
		} else {
			fmt.Printf("  ✅ %s\n", skillID)
		}
	}

	if missingCount > 0 {
		fmt.Printf("\n⚠️  有 %d 个初始技能缺失\n", missingCount)
		os.Exit(1)
	} else {
		fmt.Println("\n✅ 所有初始技能都已加载")
	}
}


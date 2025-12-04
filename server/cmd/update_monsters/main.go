package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// 打开数据库
	db, err := sql.Open("sqlite3", "./game.db?_journal_mode=WAL&_busy_timeout=5000")
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	// 启用外键约束
	if _, err := db.Exec("PRAGMA foreign_keys = ON"); err != nil {
		log.Fatalf("Failed to enable foreign keys: %v", err)
	}

	// 先执行 schema.sql 创建表结构
	fmt.Println("Loading schema.sql...")
	schemaPath := filepath.Join("database", "schema.sql")
	schemaContent, err := ioutil.ReadFile(schemaPath)
	if err != nil {
		log.Fatalf("Failed to read schema file: %v", err)
	}
	if _, err := db.Exec(string(schemaContent)); err != nil {
		log.Printf("Warning: Schema execution had errors (tables may already exist): %v", err)
	} else {
		fmt.Println("✓ Schema loaded")
	}

	// 先删除所有怪物数据
	fmt.Println("\nDeleting existing monster data...")
	if _, err := db.Exec("DELETE FROM monsters"); err != nil {
		log.Printf("Warning: Failed to delete monsters (may not exist yet): %v", err)
	} else {
		fmt.Println("✓ Deleted existing monsters")
	}

	// 读取 seed.sql 文件
	fmt.Println("\nLoading seed.sql...")
	seedPath := filepath.Join("database", "seed.sql")
	content, err := ioutil.ReadFile(seedPath)
	if err != nil {
		log.Fatalf("Failed to read seed file: %v", err)
	}

	// 执行整个 seed.sql（INSERT OR REPLACE 会处理重复）
	fmt.Println("Executing seed.sql...")
	if _, err := db.Exec(string(content)); err != nil {
		log.Fatalf("Failed to execute seed.sql: %v", err)
	}
	fmt.Println("✓ Seed data executed")

	fmt.Println("\n✅ Monster data update completed!")
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}


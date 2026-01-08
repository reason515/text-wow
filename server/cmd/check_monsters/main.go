package main

import (
	"database/sql"
	"fmt"
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

	// 如果从 cmd/check_monsters 运行，需要回到 server 目录
	serverDir := wd
	if filepath.Base(wd) == "check_monsters" {
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

	// 查询每个区域的怪物数量
	rows, err := db.Query(`
		SELECT 
			z.id,
			z.name,
			COUNT(m.id) as monster_count
		FROM zones z
		LEFT JOIN monsters m ON z.id = m.zone_id
		GROUP BY z.id, z.name
		ORDER BY z.id
	`)
	if err != nil {
		log.Fatalf("Failed to query: %v", err)
	}
	defer rows.Close()

	fmt.Println("区域怪物统计：")
	fmt.Println("═══════════════════════════════════════════════════════════")
	fmt.Printf("%-20s %-30s %-10s\n", "区域ID", "区域名称", "怪物数量")
	fmt.Println("───────────────────────────────────────────────────────────")

	var totalZones, zonesWithMonsters, zonesWithoutMonsters int
	var zonesWithoutMonstersList []string

	for rows.Next() {
		var zoneID, zoneName string
		var monsterCount int
		err := rows.Scan(&zoneID, &zoneName, &monsterCount)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}

		totalZones++
		if monsterCount > 0 {
			zonesWithMonsters++
			fmt.Printf("%-20s %-30s %-10d ✓\n", zoneID, zoneName, monsterCount)
		} else {
			zonesWithoutMonsters++
			zonesWithoutMonstersList = append(zonesWithoutMonstersList, fmt.Sprintf("%s (%s)", zoneID, zoneName))
			fmt.Printf("%-20s %-30s %-10d ✗\n", zoneID, zoneName, monsterCount)
		}
	}

	fmt.Println("═══════════════════════════════════════════════════════════")
	fmt.Printf("\n统计：\n")
	fmt.Printf("  总区域数: %d\n", totalZones)
	fmt.Printf("  有怪物的区域: %d\n", zonesWithMonsters)
	fmt.Printf("  无怪物的区域: %d\n", zonesWithoutMonsters)

	if len(zonesWithoutMonstersList) > 0 {
		fmt.Printf("\n缺少怪物的区域：\n")
		for _, zone := range zonesWithoutMonstersList {
			fmt.Printf("  - %s\n", zone)
		}
	}
}










































































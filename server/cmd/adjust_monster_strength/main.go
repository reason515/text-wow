package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	var (
		dbPath      = flag.String("db", "../game.db", "数据库路径")
		levelMin    = flag.Int("min", 0, "等级下限")
		levelMax    = flag.Int("max", 0, "等级上限")
		hpMult      = flag.Float64("hp", 1.0, "生命值倍数")
		attackMult  = flag.Float64("attack", 1.0, "攻击力倍数")
		defenseMult = flag.Float64("defense", 1.0, "防御力倍数")
		critBonus   = flag.Float64("crit", 0.0, "暴击率加成（绝对值，如0.02表示+2%）")
		apply       = flag.Bool("apply", false, "是否立即应用到怪物数据")
		list        = flag.Bool("list", false, "列出所有配置")
		show        = flag.Bool("show", false, "显示当前配置")
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

	// 确保配置表存在
	if err := ensureConfigTable(db); err != nil {
		log.Fatalf("Failed to ensure config table: %v", err)
	}

	// 列出配置
	if *list {
		if err := listConfigs(db); err != nil {
			log.Fatalf("Failed to list configs: %v", err)
		}
		return
	}

	// 显示配置
	if *show {
		if err := showConfig(db, *levelMin, *levelMax); err != nil {
			log.Fatalf("Failed to show config: %v", err)
		}
		return
	}

	// 设置配置
	if *levelMin > 0 && *levelMax > 0 {
		if err := setConfig(db, *levelMin, *levelMax, *hpMult, *attackMult, *defenseMult, *critBonus); err != nil {
			log.Fatalf("Failed to set config: %v", err)
		}
		fmt.Printf("✅ 配置已更新：等级 %d-%d\n", *levelMin, *levelMax)
		fmt.Printf("   HP倍数: %.2f, 攻击倍数: %.2f, 防御倍数: %.2f, 暴击率加成: %.2f%%\n",
			*hpMult, *attackMult, *defenseMult, *critBonus*100)
	} else {
		fmt.Println("请使用 -min 和 -max 指定等级范围")
		flag.Usage()
		os.Exit(1)
	}

	// 应用配置
	if *apply {
		fmt.Println("\n正在应用配置到怪物数据...")
		if err := applyConfigToMonsters(db); err != nil {
			log.Fatalf("Failed to apply config: %v", err)
		}
		fmt.Println("✅ 配置已应用到所有怪物")
	} else {
		fmt.Println("\n提示：使用 -apply 参数可以立即将配置应用到怪物数据")
	}
}

func ensureConfigTable(db *sql.DB) error {
	_, err := db.Exec(`
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
	return err
}

func setConfig(db *sql.DB, levelMin, levelMax int, hpMult, attackMult, defenseMult, critBonus float64) error {
	description := fmt.Sprintf("等级 %d-%d：HP ×%.2f, 攻击 ×%.2f, 防御 ×%.2f, 暴击率 +%.2f%%",
		levelMin, levelMax, hpMult, attackMult, defenseMult, critBonus*100)

	_, err := db.Exec(`
		INSERT OR REPLACE INTO monster_strength_config 
		(level_min, level_max, hp_multiplier, attack_multiplier, defense_multiplier, crit_rate_bonus, description, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP)
	`, levelMin, levelMax, hpMult, attackMult, defenseMult, critBonus, description)
	return err
}

func listConfigs(db *sql.DB) error {
	rows, err := db.Query(`
		SELECT level_min, level_max, hp_multiplier, attack_multiplier, 
		       defense_multiplier, crit_rate_bonus, description, is_active
		FROM monster_strength_config
		ORDER BY level_min
	`)
	if err != nil {
		return err
	}
	defer rows.Close()

	fmt.Println("当前强度配置：")
	fmt.Println("═══════════════════════════════════════════════════════════")
	for rows.Next() {
		var levelMin, levelMax int
		var hpMult, attackMult, defenseMult, critBonus float64
		var description string
		var isActive int

		if err := rows.Scan(&levelMin, &levelMax, &hpMult, &attackMult, &defenseMult, &critBonus, &description, &isActive); err != nil {
			return err
		}

		status := "启用"
		if isActive == 0 {
			status = "禁用"
		}

		fmt.Printf("等级 %2d-%2d: HP×%.2f, 攻击×%.2f, 防御×%.2f, 暴击+%.2f%% [%s]\n",
			levelMin, levelMax, hpMult, attackMult, defenseMult, critBonus*100, status)
		if description != "" {
			fmt.Printf("          %s\n", description)
		}
	}
	fmt.Println("═══════════════════════════════════════════════════════════")
	return nil
}

func showConfig(db *sql.DB, levelMin, levelMax int) error {
	var hpMult, attackMult, defenseMult, critBonus float64
	var description string
	var isActive int

	err := db.QueryRow(`
		SELECT hp_multiplier, attack_multiplier, defense_multiplier, 
		       crit_rate_bonus, description, is_active
		FROM monster_strength_config
		WHERE level_min = ? AND level_max = ?
	`, levelMin, levelMax).Scan(&hpMult, &attackMult, &defenseMult, &critBonus, &description, &isActive)

	if err == sql.ErrNoRows {
		fmt.Printf("未找到等级 %d-%d 的配置\n", levelMin, levelMax)
		return nil
	}
	if err != nil {
		return err
	}

	status := "启用"
	if isActive == 0 {
		status = "禁用"
	}

	fmt.Printf("等级 %d-%d 的配置：\n", levelMin, levelMax)
	fmt.Printf("  HP倍数: %.2f\n", hpMult)
	fmt.Printf("  攻击倍数: %.2f\n", attackMult)
	fmt.Printf("  防御倍数: %.2f\n", defenseMult)
	fmt.Printf("  暴击率加成: %.2f%%\n", critBonus*100)
	fmt.Printf("  状态: %s\n", status)
	if description != "" {
		fmt.Printf("  描述: %s\n", description)
	}
	return nil
}

func applyConfigToMonsters(db *sql.DB) error {
	// 获取所有启用的配置
	rows, err := db.Query(`
		SELECT level_min, level_max, hp_multiplier, attack_multiplier, 
		       defense_multiplier, crit_rate_bonus
		FROM monster_strength_config
		WHERE is_active = 1
		ORDER BY level_min
	`)
	if err != nil {
		return err
	}
	defer rows.Close()

	configs := []struct {
		levelMin, levelMax int
		hpMult, attackMult, defenseMult, critBonus float64
	}{}

	for rows.Next() {
		var c struct {
			levelMin, levelMax int
			hpMult, attackMult, defenseMult, critBonus float64
		}
		if err := rows.Scan(&c.levelMin, &c.levelMax, &c.hpMult, &c.attackMult, &c.defenseMult, &c.critBonus); err != nil {
			return err
		}
		configs = append(configs, c)
	}

	// 对每个配置范围应用
	for _, config := range configs {
		_, err := db.Exec(`
			UPDATE monsters
			SET 
				hp = ROUND(hp * ?),
				physical_attack = ROUND(physical_attack * ?),
				magic_attack = ROUND(magic_attack * ?),
				physical_defense = ROUND(physical_defense * ?),
				magic_defense = ROUND(magic_defense * ?),
				phys_crit_rate = MIN(0.4, phys_crit_rate + ?),
				spell_crit_rate = MIN(0.4, spell_crit_rate + ?)
			WHERE level >= ? AND level <= ?
		`, config.hpMult, config.attackMult, config.attackMult, config.defenseMult, config.defenseMult,
			config.critBonus, config.critBonus, config.levelMin, config.levelMax)
		if err != nil {
			return fmt.Errorf("failed to apply config for level %d-%d: %w", config.levelMin, config.levelMax, err)
		}
	}

	return nil
}
























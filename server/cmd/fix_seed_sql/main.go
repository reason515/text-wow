package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func main() {
	// 获取当前工作目录，然后找到 server 目录
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get working directory: %v", err)
	}

	// 如果从 cmd/fix_seed_sql 运行，需要回到 server 目录
	serverDir := wd
	if filepath.Base(wd) == "fix_seed_sql" {
		serverDir = filepath.Join(wd, "..", "..")
	} else if filepath.Base(filepath.Dir(wd)) == "cmd" {
		serverDir = filepath.Join(wd, "..")
	}

	// 读取 seed.sql
	seedPath := filepath.Join(serverDir, "database", "seed.sql")
	content, err := ioutil.ReadFile(seedPath)
	if err != nil {
		log.Fatalf("Failed to read seed.sql: %v", err)
	}

	lines := strings.Split(string(content), "\n")
	var fixedLines []string
	var inInsert bool

	for i, line := range lines {
		// 检测 INSERT 语句开始
		if strings.Contains(line, "INSERT OR REPLACE INTO monsters") {
			inInsert = true
			fixedLines = append(fixedLines, line)
			continue
		}

		// 如果在 INSERT 块中
		if inInsert {
			// 检测 VALUES 行
			if strings.Contains(line, "VALUES") {
				fixedLines = append(fixedLines, line)
				continue
			}

			// 检测数据行
			trimmed := strings.TrimSpace(line)
			if trimmed == "" {
				fixedLines = append(fixedLines, line)
				continue
			}

			// 检查是否是 INSERT 块的结束（下一个 INSERT 或注释）
			if strings.HasPrefix(trimmed, "--") || strings.Contains(trimmed, "INSERT OR REPLACE") {
				// 结束当前 INSERT 块
				inInsert = false
				// 处理上一行（最后一行数据）
				if len(fixedLines) > 0 {
					lastLine := fixedLines[len(fixedLines)-1]
					// 如果最后一行以 ,); 结尾，修复它
					if strings.HasSuffix(strings.TrimSpace(lastLine), ",);") {
						fixedLines[len(fixedLines)-1] = strings.TrimSuffix(strings.TrimSpace(lastLine), ",);") + ");"
					}
				}
				// 添加当前行
				fixedLines = append(fixedLines, line)
				continue
			}

			// 修复数据行：将 ,); 替换为 ,
			// 但如果是最后一行（下一行是空行或注释或新的INSERT），则替换为 );
			nextLineIdx := i + 1
			isLastDataLine := false
			if nextLineIdx < len(lines) {
				nextLine := strings.TrimSpace(lines[nextLineIdx])
				if nextLine == "" || strings.HasPrefix(nextLine, "--") || strings.Contains(nextLine, "INSERT OR REPLACE") {
					isLastDataLine = true
				}
			} else {
				isLastDataLine = true
			}

			if isLastDataLine {
				// 最后一行数据，应该以 ); 结尾
				fixed := strings.TrimSpace(line)
				fixed = regexp.MustCompile(`,\s*\);?\s*$`).ReplaceAllString(fixed, ");")
				fixedLines = append(fixedLines, fixed)
			} else {
				// 中间行数据，应该以 , 结尾
				fixed := strings.TrimSpace(line)
				fixed = regexp.MustCompile(`,\s*\);?\s*$`).ReplaceAllString(fixed, ",")
				fixedLines = append(fixedLines, fixed)
			}
		} else {
			// 不在 INSERT 块中，直接添加
			fixedLines = append(fixedLines, line)
		}
	}

	// 处理最后一个 INSERT 块的最后一行
	if inInsert && len(fixedLines) > 0 {
		lastLine := fixedLines[len(fixedLines)-1]
		if strings.HasSuffix(strings.TrimSpace(lastLine), ",);") {
			fixedLines[len(fixedLines)-1] = strings.TrimSuffix(strings.TrimSpace(lastLine), ",);") + ");"
		}
	}

	// 写入修复后的文件
	output := strings.Join(fixedLines, "\n")
	
	// 创建备份
	backupPath := seedPath + ".backup"
	err = ioutil.WriteFile(backupPath, content, 0644)
	if err != nil {
		log.Printf("Warning: Failed to create backup: %v", err)
	} else {
		fmt.Printf("✓ 已创建备份: %s\n", backupPath)
	}

	// 写入修复后的文件
	err = ioutil.WriteFile(seedPath, []byte(output), 0644)
	if err != nil {
		log.Fatalf("Failed to write fixed seed.sql: %v", err)
	}

	fmt.Printf("✓ 已修复 seed.sql\n")
	fmt.Printf("  修复前: %d 行\n", len(lines))
	fmt.Printf("  修复后: %d 行\n", len(fixedLines))
	
	// 统计修复的行数
	fixedCount := 0
	for _, line := range lines {
		if strings.Contains(line, ",);") {
			fixedCount++
		}
	}
	fmt.Printf("  修复了 %d 行包含 ',);' 的错误\n", fixedCount)
}


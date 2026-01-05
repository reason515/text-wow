package main

import (
	"fmt"
	"log"
	"os"

	"text-wow/internal/database"
	"text-wow/internal/test/runner"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <test_directory>")
		fmt.Println("Example: go run main.go ../../test/cases/calculator")
		os.Exit(1)
	}

	// 初始化数据库
	if err := database.Init(); err != nil {
		log.Fatalf("❌ Failed to initialize database: %v", err)
	}
	defer database.Close()

	testDir := os.Args[1]

	// 创建测试运行器
	tr := runner.NewTestRunner()

	// 运行所有测试
	results, err := tr.RunAllTests(testDir)
	if err != nil {
		fmt.Printf("Error running tests: %v\n", err)
		os.Exit(1)
	}

	// 生成报告
	reporter := runner.NewReporter()
	totalPassed := 0
	totalFailed := 0
	totalSkipped := 0

	for _, result := range results {
		report := reporter.GenerateReport(result)
		fmt.Println(report)
		fmt.Println()

		totalPassed += result.PassedTests
		totalFailed += result.FailedTests
		totalSkipped += result.SkippedTests
	}

	// 总结
	fmt.Println("=== 测试总结 ===")
	fmt.Printf("总通过: %d\n", totalPassed)
	fmt.Printf("总失败: %d\n", totalFailed)
	fmt.Printf("总跳过: %d\n", totalSkipped)

	if totalFailed > 0 {
		os.Exit(1)
	}
}


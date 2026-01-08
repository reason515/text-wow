package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"text-wow/internal/database"
	"text-wow/internal/test/runner"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stdout, "Usage: go run main.go <test_directory>\n")
		fmt.Fprintf(os.Stdout, "Example: go run main.go ../../test/cases/calculator\n")
		os.Exit(1)
	}

	// 在测试模式下，禁用log输出（避免emoji字符导致序列化错误）
	// 只有在TEST_DEBUG=1时才输出日志
	if os.Getenv("TEST_DEBUG") != "1" && os.Getenv("TEST_DEBUG") != "true" {
		log.SetOutput(io.Discard)
	}

	// 初始化数据库
	if err := database.Init(); err != nil {
		// 使用fmt.Fprintf到stderr，避免使用log.Fatalf（可能包含特殊字符）
		fmt.Fprintf(os.Stderr, "Failed to initialize database: %v\n", err)
		os.Exit(1)
	}
	defer database.Close()

	testDir := os.Args[1]

	// 创建测试运行器
	tr := runner.NewTestRunner()

	// 运行所有测试
	results, err := tr.RunAllTests(testDir)
	if err != nil {
		fmt.Fprintf(os.Stdout, "Error running tests: %v\n", err)
		os.Exit(1)
	}

	// 生成报告
	reporter := runner.NewReporter()
	totalPassed := 0
	totalFailed := 0
	totalSkipped := 0

	for _, result := range results {
		report := reporter.GenerateReport(result)
		// 使用fmt.Fprintf到stdout而不是fmt.Println，避免序列化问题
		fmt.Fprintf(os.Stdout, "%s\n\n", report)

		totalPassed += result.PassedTests
		totalFailed += result.FailedTests
		totalSkipped += result.SkippedTests
	}

	// 总结
	fmt.Fprintf(os.Stdout, "=== 测试总结 ===\n")
	fmt.Fprintf(os.Stdout, "总通过: %d\n", totalPassed)
	fmt.Fprintf(os.Stdout, "总失败: %d\n", totalFailed)
	fmt.Fprintf(os.Stdout, "总跳过: %d\n", totalSkipped)

	if totalFailed > 0 {
		os.Exit(1)
	}
}


package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"syscall"

	"text-wow/internal/database"
	"text-wow/internal/test/runner"
)

var (
	kernel32                = syscall.NewLazyDLL("kernel32.dll")
	setConsoleOutputCP      = kernel32.NewProc("SetConsoleOutputCP")
	getConsoleOutputCP      = kernel32.NewProc("GetConsoleOutputCP")
)

// setUTF8Output 设置控制台输出为 UTF-8 编码（Windows 特定）
func setUTF8Output() {
	if runtime.GOOS == "windows" {
		// UTF-8 代码页是 65001
		const CP_UTF8 = 65001
		
		// 获取当前代码页
		currentCP, _, _ := getConsoleOutputCP.Call()
		
		// 如果当前代码页不是 UTF-8，则设置为 UTF-8
		if currentCP != CP_UTF8 {
			// 设置控制台输出代码页为 UTF-8
			setConsoleOutputCP.Call(CP_UTF8)
		}
	}
}

func main() {
	// 设置 UTF-8 输出编码
	setUTF8Output()

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


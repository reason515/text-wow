package main

import (
	"fmt"
	"os"
	"text-wow/internal/database"
	"text-wow/internal/test/runner"
)

func main() {
	// 初始化数据库
	if err := database.Init(); err != nil {
		fmt.Printf("Failed to initialize database: %v\n", err)
		os.Exit(1)
	}
	defer database.Close()

	// 创建测试运行器
	tr := runner.NewTestRunner()

	// 创建一个简单的测试用例
	setup := []string{
		"创建一个1级战士角色，怒气=100/100",
	}

	// 执行setup
	for _, instruction := range setup {
		if err := tr.ParseAndExecuteSetupInstruction(instruction); err != nil {
			fmt.Printf("Setup failed: %v\n", err)
			os.Exit(1)
		}
	}

	// 创建断言执行器
	ae := tr.GetAssertionExecutor()

	// 测试获取 character.resource
	path := "character.resource"
	fmt.Printf("Testing path: %s\n", path)
	
	actual, err := ae.GetValue(path)
	fmt.Printf("getValue returned: actual=%v (type=%T), err=%v\n", actual, actual, err)
	
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else if actual == nil {
		fmt.Printf("Warning: actual is nil but err is nil\n")
	} else {
		fmt.Printf("Success: actual=%v\n", actual)
	}

	// 测试执行断言
	assertion := runner.Assertion{
		Type:     "equals",
		Target:   path,
		Expected: "100",
		Message:  "测试消息",
	}

	result := ae.Execute(assertion)
	fmt.Printf("\nAssertion Result:\n")
	fmt.Printf("  Status: %s\n", result.Status)
	fmt.Printf("  Error: %s\n", result.Error)
	fmt.Printf("  Actual: %v\n", result.Actual)
	fmt.Printf("  Expected: %s\n", result.Expected)
}



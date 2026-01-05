package runner

import (
	"fmt"
)

// Reporter 报告生成器
type Reporter struct {
}

// NewReporter 创建报告生成器
func NewReporter() *Reporter {
	return &Reporter{}
}

// GenerateReport 生成测试报告
func (r *Reporter) GenerateReport(result *TestSuiteResult) string {
	report := fmt.Sprintf("=== 测试报告 ===\n")
	report += fmt.Sprintf("测试套件: %s\n", result.TestSuite)
	report += fmt.Sprintf("总测试数: %d\n", result.TotalTests)
	report += fmt.Sprintf("通过: %d\n", result.PassedTests)
	report += fmt.Sprintf("失败: %d\n", result.FailedTests)
	report += fmt.Sprintf("跳过: %d\n", result.SkippedTests)
	report += fmt.Sprintf("耗时: %v\n\n", result.Duration)

	for _, testResult := range result.Results {
		statusIcon := "✓"
		if testResult.Status == "failed" {
			statusIcon = "✗"
		} else if testResult.Status == "skipped" {
			statusIcon = "⊘"
		}

		report += fmt.Sprintf("[%s] %s (%v)\n", statusIcon, testResult.TestName, testResult.Duration)

		if testResult.Status == "failed" {
			if testResult.Error != "" {
				report += fmt.Sprintf("  错误: %s\n", testResult.Error)
			}

			// 显示失败的断言
			for _, assertion := range testResult.Assertions {
				if assertion.Status == "failed" {
					report += fmt.Sprintf("  断言失败: %s %s %s (实际: %v)\n",
						assertion.Target, assertion.Type, assertion.Expected, assertion.Actual)
					if assertion.Error != "" {
						report += fmt.Sprintf("    错误: %s\n", assertion.Error)
					}
					if assertion.Message != "" {
						report += fmt.Sprintf("    %s\n", assertion.Message)
					}
				}
			}
		}
	}

	return report
}

// GenerateJSONReport 生成JSON格式报告
func (r *Reporter) GenerateJSONReport(result *TestSuiteResult) string {
	// TODO: 实现JSON格式报告
	return ""
}

// GenerateHTMLReport 生成HTML格式报告
func (r *Reporter) GenerateHTMLReport(result *TestSuiteResult) string {
	// TODO: 实现HTML格式报告
	return ""
}


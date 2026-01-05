package runner

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// AssertionExecutor 断言执行器
type AssertionExecutor struct {
	context map[string]interface{} // 测试上下文（存储测试数据）
}

// NewAssertionExecutor 创建断言执行器
func NewAssertionExecutor() *AssertionExecutor {
	return &AssertionExecutor{
		context: make(map[string]interface{}),
	}
}

// Execute 执行断言
func (ae *AssertionExecutor) Execute(assertion Assertion) AssertionResult {
	result := AssertionResult{
		Type:     assertion.Type,
		Target:   assertion.Target,
		Expected: assertion.Expected,
		Status:   "pending",
		Message:  assertion.Message,
	}

	// 获取实际值
	actual, err := ae.getValue(assertion.Target)
	if err != nil {
		result.Status = "failed"
		result.Error = fmt.Sprintf("failed to get value: %v", err)
		return result
	}

	result.Actual = actual

	// 根据类型执行断言
	switch assertion.Type {
	case "equals":
		result.Status = ae.assertEquals(actual, assertion.Expected)
	case "greater_than":
		result.Status = ae.assertGreaterThan(actual, assertion.Expected)
	case "less_than":
		result.Status = ae.assertLessThan(actual, assertion.Expected)
	case "greater_than_or_equal":
		result.Status = ae.assertGreaterThanOrEqual(actual, assertion.Expected)
	case "contains":
		result.Status = ae.assertContains(actual, assertion.Expected)
	case "approximately":
		result.Status = ae.assertApproximately(actual, assertion.Expected, assertion.Tolerance)
	case "range":
		result.Status = ae.assertRange(actual, assertion.Expected)
	default:
		result.Status = "failed"
		result.Error = fmt.Sprintf("unknown assertion type: %s", assertion.Type)
	}

	return result
}

// getValue 获取值（从上下文或通过路径）
func (ae *AssertionExecutor) getValue(path string) (interface{}, error) {
	// 简化实现：从上下文获取
	// TODO: 实现路径解析（如 "character.hp"）
	if value, exists := ae.context[path]; exists {
		return value, nil
	}

	// 尝试解析为数字
	if num, err := strconv.Atoi(path); err == nil {
		return num, nil
	}

	return nil, fmt.Errorf("value not found: %s", path)
}

// assertEquals 断言相等
func (ae *AssertionExecutor) assertEquals(actual interface{}, expected string) string {
	actualStr := fmt.Sprintf("%v", actual)
	if actualStr == expected {
		return "passed"
	}
	return "failed"
}

// assertGreaterThan 断言大于
func (ae *AssertionExecutor) assertGreaterThan(actual interface{}, expected string) string {
	actualNum, err := ae.toNumber(actual)
	if err != nil {
		return "failed"
	}

	expectedNum, err := strconv.ParseFloat(expected, 64)
	if err != nil {
		return "failed"
	}

	if actualNum > expectedNum {
		return "passed"
	}
	return "failed"
}

// assertLessThan 断言小于
func (ae *AssertionExecutor) assertLessThan(actual interface{}, expected string) string {
	actualNum, err := ae.toNumber(actual)
	if err != nil {
		return "failed"
	}

	expectedNum, err := strconv.ParseFloat(expected, 64)
	if err != nil {
		return "failed"
	}

	if actualNum < expectedNum {
		return "passed"
	}
	return "failed"
}

// assertGreaterThanOrEqual 断言大于等于
func (ae *AssertionExecutor) assertGreaterThanOrEqual(actual interface{}, expected string) string {
	actualNum, err := ae.toNumber(actual)
	if err != nil {
		return "failed"
	}

	expectedNum, err := strconv.ParseFloat(expected, 64)
	if err != nil {
		return "failed"
	}

	if actualNum >= expectedNum {
		return "passed"
	}
	return "failed"
}

// assertContains 断言包含
func (ae *AssertionExecutor) assertContains(actual interface{}, expected string) string {
	actualStr := fmt.Sprintf("%v", actual)
	if strings.Contains(actualStr, expected) {
		return "passed"
	}
	return "failed"
}

// assertApproximately 断言近似相等
func (ae *AssertionExecutor) assertApproximately(actual interface{}, expected string, tolerance float64) string {
	actualNum, err := ae.toNumber(actual)
	if err != nil {
		return "failed"
	}

	expectedNum, err := strconv.ParseFloat(expected, 64)
	if err != nil {
		return "failed"
	}

	diff := math.Abs(actualNum - expectedNum)
	if diff <= tolerance {
		return "passed"
	}
	return "failed"
}

// assertRange 断言范围
func (ae *AssertionExecutor) assertRange(actual interface{}, expected string) string {
	// 解析范围 "[min, max]"
	expected = strings.TrimSpace(expected)
	if !strings.HasPrefix(expected, "[") || !strings.HasSuffix(expected, "]") {
		return "failed"
	}

	expected = strings.Trim(expected, "[]")
	parts := strings.Split(expected, ",")
	if len(parts) != 2 {
		return "failed"
	}

	minStr := strings.TrimSpace(parts[0])
	maxStr := strings.TrimSpace(parts[1])

	min, err := strconv.ParseFloat(minStr, 64)
	if err != nil {
		return "failed"
	}

	max, err := strconv.ParseFloat(maxStr, 64)
	if err != nil {
		return "failed"
	}

	actualNum, err := ae.toNumber(actual)
	if err != nil {
		return "failed"
	}

	if actualNum >= min && actualNum <= max {
		return "passed"
	}
	return "failed"
}

// toNumber 转换为数字
func (ae *AssertionExecutor) toNumber(value interface{}) (float64, error) {
	switch v := value.(type) {
	case int:
		return float64(v), nil
	case int64:
		return float64(v), nil
	case float64:
		return v, nil
	case float32:
		return float64(v), nil
	case string:
		return strconv.ParseFloat(v, 64)
	default:
		return 0, fmt.Errorf("cannot convert to number: %T", value)
	}
}

// SetContext 设置测试上下文
func (ae *AssertionExecutor) SetContext(key string, value interface{}) {
	ae.context[key] = value
}

// ClearContext 清空测试上下文
func (ae *AssertionExecutor) ClearContext() {
	ae.context = make(map[string]interface{})
}



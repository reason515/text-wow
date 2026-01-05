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
	
	// 检查 (nil, nil) 的情况
	if err == nil && actual == nil {
		result.Status = "failed"
		result.Error = fmt.Sprintf("getValue returned (nil, nil) for path: %s (context keys: %v)", 
			assertion.Target, getMapKeys(ae.context))
		result.Actual = nil
		return result
	}
	
	if err != nil {
		result.Status = "failed"
		result.Error = fmt.Sprintf("failed to get value: %v", err)
		result.Actual = nil
		return result
	}

	// 检查 actual 是否为 nil（对于 "null" 或 "nil" 的断言，nil 是有效的）
	// 只有在 expected 不是 "null" 或 "nil" 时才检查 nil
	if actual == nil && assertion.Expected != "null" && assertion.Expected != "nil" {
		result.Status = "failed"
		result.Error = fmt.Sprintf("value is nil for path: %s (context keys: %v)", 
			assertion.Target, getMapKeys(ae.context))
		result.Actual = nil
		return result
	}

	result.Actual = actual

	// 根据类型执行断言
	switch assertion.Type {
	case "equals":
		result.Status = ae.assertEquals(actual, assertion.Expected)
	case "greater_than":
		// 对于greater_than断言，expected也可能是另一个路径（如turn_order[1].speed）
		expectedValue, err := ae.getValue(assertion.Expected)
		if err == nil {
			// 如果expected是路径，使用解析后的值
			result.Status = ae.assertGreaterThan(actual, fmt.Sprintf("%v", expectedValue))
		} else {
			// 否则使用原始expected值
			result.Status = ae.assertGreaterThan(actual, assertion.Expected)
		}
	case "less_than":
		result.Status = ae.assertLessThan(actual, assertion.Expected)
	case "greater_than_or_equal":
		result.Status = ae.assertGreaterThanOrEqual(actual, assertion.Expected)
	case "less_than_or_equal":
		result.Status = ae.assertLessThanOrEqual(actual, assertion.Expected)
	case "contains":
		result.Status = ae.assertContains(actual, assertion.Expected)
	case "approximately":
		result.Status = ae.assertApproximately(actual, assertion.Expected, assertion.Tolerance)
	case "range":
		result.Status = ae.assertRange(actual, assertion.Expected)
	case "not_contains":
		result.Status = ae.assertNotContains(actual, assertion.Expected)
	case "not_equals":
		result.Status = ae.assertNotEquals(actual, assertion.Expected)
	case "not_null":
		result.Status = ae.assertNotNull(actual)
	default:
		result.Status = "failed"
		result.Error = fmt.Sprintf("unknown assertion type: %s", assertion.Type)
	}

	return result
}

// getValue 获取值（从上下文或通过路径）
func (ae *AssertionExecutor) getValue(path string) (interface{}, error) {
	// 支持路径解析（如 "character.hp", "turn_order[0].type"）
	// 先尝试直接获取（包括已经设置的 turn_order[0].type 等键）
	if value, exists := ae.context[path]; exists {
		// 如果值是nil，返回nil而不是错误（这样断言可以正确处理）
		return value, nil
	}
	
	// 尝试解析数组索引路径（如 "turn_order[0].type"）
	// 注意：这个路径可能已经被设置为键（如 "turn_order[0].type"），也可能需要从数组中解析
	if strings.Contains(path, "[") && strings.Contains(path, "]") {
		// 解析数组索引
		openBracket := strings.Index(path, "[")
		closeBracket := strings.Index(path, "]")
		if openBracket > 0 && closeBracket > openBracket {
			arrayKey := path[:openBracket]
			indexStr := path[openBracket+1 : closeBracket]
			restPath := ""
			if closeBracket+1 < len(path) {
				restPath = path[closeBracket+1:]
				if strings.HasPrefix(restPath, ".") {
					restPath = restPath[1:]
				}
			}
			
			// 获取数组
			arrayValue, exists := ae.context[arrayKey]
			if !exists {
				return nil, fmt.Errorf("array not found: %s", arrayKey)
			}
			
			// 解析索引
			index, err := strconv.Atoi(indexStr)
			if err != nil {
				return nil, fmt.Errorf("invalid index: %s", indexStr)
			}
			
			// 根据类型处理
			switch arr := arrayValue.(type) {
			case []interface{}:
				if index < 0 || index >= len(arr) {
					return nil, fmt.Errorf("index out of range: %d", index)
				}
				participant := arr[index]
				if participant == nil {
					return nil, nil
				}
				// 如果有剩余路径，继续解析（简化处理）
				if restPath != "" {
					// 尝试从map中获取
					if m, ok := participant.(map[string]interface{}); ok {
						if val, exists := m[restPath]; exists {
							return val, nil
						}
						// 尝试从上下文中获取已设置的key（如 turn_order[0].character.id）
						key := fmt.Sprintf("turn_order[%d].%s", index, restPath)
						if value, exists := ae.context[key]; exists {
							return value, nil
						}
					}
				}
				return participant, nil
			default:
				return nil, fmt.Errorf("unsupported array type for path: %s", path)
			}
		}
	}
	
	// 尝试解析点号路径（如 "character.hp"）
	if strings.Contains(path, ".") {
		parts := strings.Split(path, ".")
		if len(parts) == 2 {
			baseKey := parts[0]
			fieldKey := parts[1]
			
			// 尝试从上下文获取基础对象
			if baseValue, exists := ae.context[baseKey]; exists {
				// 根据类型处理
				switch obj := baseValue.(type) {
				case map[string]interface{}:
					if fieldValue, ok := obj[fieldKey]; ok {
						return fieldValue, nil
					}
				default:
					// 尝试使用组合键
					combinedKey := fmt.Sprintf("%s.%s", baseKey, fieldKey)
					if value, exists := ae.context[combinedKey]; exists {
						return value, nil
					}
				}
			}
		}
	}

	// 尝试解析为数字
	if num, err := strconv.Atoi(path); err == nil {
		return num, nil
	}

	return nil, fmt.Errorf("value not found: %s", path)
}

// assertEquals 断言相等
func (ae *AssertionExecutor) assertEquals(actual interface{}, expected string) string {
	// 处理 "null" 或 "nil" 的特殊情况
	if expected == "null" || expected == "nil" {
		if actual == nil {
			return "passed"
		}
		// 检查是否是空字符串或零值
		actualStr := fmt.Sprintf("%v", actual)
		// 接受 nil、<nil>、空字符串、0 等表示空值的字符串
		if actualStr == "" || actualStr == "<nil>" || actualStr == "0" || actualStr == "nil" || actualStr == "null" {
			return "passed"
		}
		return "failed"
	}
	
	// 如果 expected 是一个路径（如 "character.id"），尝试从上下文获取值
	if strings.Contains(expected, ".") {
		expectedValue, err := ae.getValue(expected)
		if err == nil && expectedValue != nil {
			expected = fmt.Sprintf("%v", expectedValue)
		}
	}
	
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

// assertLessThanOrEqual 断言小于等于
func (ae *AssertionExecutor) assertLessThanOrEqual(actual interface{}, expected string) string {
	actualNum, err := ae.toNumber(actual)
	if err != nil {
		return "failed"
	}

	expectedNum, err := strconv.ParseFloat(expected, 64)
	if err != nil {
		return "failed"
	}

	if actualNum <= expectedNum {
		return "passed"
	}
	return "failed"
}

// assertContains 断言包含
func (ae *AssertionExecutor) assertContains(actual interface{}, expected string) string {
	actualStr := fmt.Sprintf("%v", actual)
	
	// 支持"or"逻辑，如"prefix_or_suffix"表示包含"prefix"或"suffix"
	if strings.Contains(expected, "_or_") {
		parts := strings.Split(expected, "_or_")
		for _, part := range parts {
			if strings.Contains(actualStr, part) {
				return "passed"
			}
		}
		return "failed"
	}
	
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

// getMapKeys 获取 map 的所有键（用于调试）
func getMapKeys(m map[string]interface{}) []string {
	if m == nil {
		return []string{"<nil>"}
	}
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// assertNotContains 断言不包含
func (ae *AssertionExecutor) assertNotContains(actual interface{}, expected string) string {
	actualStr := fmt.Sprintf("%v", actual)
	
	if strings.Contains(actualStr, expected) {
		return fmt.Sprintf("expected not to contain '%s', but got '%s'", expected, actualStr)
	}
	
	return "passed"
}

// assertNotEquals 断言不等于
func (ae *AssertionExecutor) assertNotEquals(actual interface{}, expected string) string {
	actualStr := fmt.Sprintf("%v", actual)
	
	if actualStr == expected {
		return fmt.Sprintf("expected not equal to '%s', but got '%s'", expected, actualStr)
	}
	
	return "passed"
}

// assertNotNull 断言不为空
func (ae *AssertionExecutor) assertNotNull(actual interface{}) string {
	if actual == nil {
		return "expected not null, but got nil"
	}
	
	// 检查空字符串
	if str, ok := actual.(string); ok && str == "" {
		return "expected not null, but got empty string"
	}
	
	// 检查空数组/切片
	switch v := actual.(type) {
	case []interface{}:
		if len(v) == 0 {
			return "expected not null, but got empty array"
		}
	case []string:
		if len(v) == 0 {
			return "expected not null, but got empty array"
		}
	}
	
	return "passed"
}


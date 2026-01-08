package runner

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// YAMLParser YAML解析器
type YAMLParser struct {
}

// NewYAMLParser 创建YAML解析器
func NewYAMLParser() *YAMLParser {
	return &YAMLParser{}
}

// ParseTestSuite 解析测试套件
func (yp *YAMLParser) ParseTestSuite(filePath string) (*TestSuite, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	// 移除UTF-8 BOM（如果存在）
	if len(data) >= 3 && data[0] == 0xEF && data[1] == 0xBB && data[2] == 0xBF {
		data = data[3:]
	}

	var suite TestSuite
	if err := yaml.Unmarshal(data, &suite); err != nil {
		return nil, fmt.Errorf("failed to parse YAML: %w", err)
	}

	return &suite, nil
}

// ValidateTestSuite 验证测试套件格式
func (yp *YAMLParser) ValidateTestSuite(suite *TestSuite) error {
	if suite.TestSuite == "" {
		return fmt.Errorf("test_suite is required")
	}

	if len(suite.Tests) == 0 {
		return fmt.Errorf("tests list is empty")
	}

	for i, test := range suite.Tests {
		if test.Name == "" {
			return fmt.Errorf("test[%d].name is required", i)
		}

		if test.Category == "" {
			return fmt.Errorf("test[%d].category is required", i)
		}

		if len(test.Assertions) == 0 {
			return fmt.Errorf("test[%d].assertions is empty", i)
		}
	}

	return nil
}


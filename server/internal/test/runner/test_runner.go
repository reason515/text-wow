package runner



import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"text-wow/internal/game"
	"text-wow/internal/models"

	"gopkg.in/yaml.v3"
)

// TestRunner 核心结构和主要运行逻辑



type TestRunner struct {

	parser           *YAMLParser

	assertion        *AssertionExecutor

	reporter         *Reporter

	calculator       *game.Calculator

	equipmentManager *game.EquipmentManager

	context          *TestContext

}







func NewTestRunner() *TestRunner {

	return &TestRunner{

		parser:           NewYAMLParser(),

		assertion:        NewAssertionExecutor(),

		reporter:         NewReporter(),

		calculator:       game.NewCalculator(),

		equipmentManager: game.NewEquipmentManager(),

		context: &TestContext{

			Characters: make(map[string]*models.Character),

			Monsters:   make(map[string]*models.Monster),

			Equipments: make(map[string]*models.EquipmentInstance),

			Variables:  make(map[string]interface{}),

		},

	}

}



func (tr *TestRunner) RunTestSuite(suitePath string) (*TestSuiteResult, error) {

	// 读取YAML文件

	data, err := os.ReadFile(suitePath)

	if err != nil {

		return nil, fmt.Errorf("failed to read test suite file: %w", err)

	}



	// 移除UTF-8 BOM（如果存在）

	if len(data) >= 3 && data[0] == 0xEF && data[1] == 0xBB && data[2] == 0xBF {

		data = data[3:]

	}



	// 解析YAML

	var suite TestSuite

	if err := yaml.Unmarshal(data, &suite); err != nil {

		return nil, fmt.Errorf("failed to parse test suite: %w", err)

	}



	// 运行测试用例

	result := &TestSuiteResult{

		TestSuite:  suite.TestSuite,

		TotalTests: len(suite.Tests),

		Results:    make([]TestResult, 0),

	}



	startTime := time.Now()

	for _, testCase := range suite.Tests {

		testResult := tr.RunTestCase(testCase)

		result.Results = append(result.Results, testResult)



		switch testResult.Status {

		case "passed":

			result.PassedTests++

		case "failed":

			result.FailedTests++

		case "skipped":

			result.SkippedTests++

		}

	}

	result.Duration = time.Since(startTime)



	return result, nil

}



func (tr *TestRunner) RunTestCase(testCase TestCase) TestResult {

	result := TestResult{

		TestName:   testCase.Name,

		Status:     "pending",

		Assertions: make([]AssertionResult, 0),

	}



	startTime := time.Now()

	defer func() {

		result.Duration = time.Since(startTime)

	}()



	// 在每个测试用例开始时，清空上下文（确保测试用例之间不相互影响）

	tr.context = &TestContext{

		Characters: make(map[string]*models.Character),

		Monsters:   make(map[string]*models.Monster),

		Equipments: make(map[string]*models.EquipmentInstance),

		Variables:  make(map[string]interface{}),

	}



	// 执行前置条件

	if err := tr.executeSetup(testCase.Setup); err != nil {

		result.Status = "failed"

		result.Error = fmt.Sprintf("setup failed: %v", err)

		return result

	}



	// 在setup执行后立即更新断言上下文，确保所有计算属性都被正确同步
		tr.updateAssertionContext()



	// 调试：检查setup后的上下文状态
		debugPrint("[DEBUG] RunTestCase: after setup for '%s' - characters=%d, monsters=%d, variables=%d\n", testCase.Name, len(tr.context.Characters), len(tr.context.Monsters), len(tr.context.Variables))


	if char, exists := tr.context.Characters["character"]; exists && char != nil {

		debugPrint("[DEBUG] RunTestCase: after setup, character.PhysicalAttack=%d, character pointer=%p\n", char.PhysicalAttack, char)

		// 也检查Variables中的值
			if attackVal, exists := tr.context.Variables["character_physical_attack"]; exists {

			debugPrint("[DEBUG] RunTestCase: after setup, Variables[character_physical_attack]=%v\n", attackVal)

		}

	} else if exists {

		debugPrint("[DEBUG] RunTestCase: after setup, character is nil\n")

	}

	if ratio, exists := tr.context.Variables["skill_scaling_ratio"]; exists {

		debugPrint("[DEBUG] RunTestCase: skill_scaling_ratio=%v\n", ratio)

	}



	// 执行测试步骤

	for _, step := range testCase.Steps {

		// 在执行步骤之前，检查上下文中的角色状态
			if char, exists := tr.context.Characters["character"]; exists && char != nil {

			debugPrint("[DEBUG] RunTestCase: before executeStep, character.PhysicalAttack=%d, character pointer=%p\n", char.PhysicalAttack, char)

		}

		if err := tr.executeStep(step); err != nil {

			result.Status = "failed"

			result.Error = fmt.Sprintf("step failed: %v", err)

			tr.executeTeardown(testCase.Teardown)

			return result

		}

		// 在执行步骤之后，检查上下文中的角色状态
			if char, exists := tr.context.Characters["character"]; exists && char != nil {

			debugPrint("[DEBUG] RunTestCase: after executeStep, character.PhysicalAttack=%d\n", char.PhysicalAttack)

		}

	}



	// 更新断言上下文（同步测试数据）
		tr.updateAssertionContext()



	// 执行断言

	for _, assertion := range testCase.Assertions {

		assertionResult := tr.assertion.Execute(assertion)

		result.Assertions = append(result.Assertions, assertionResult)

		if assertionResult.Status == "failed" {

			result.Status = "failed"

		}

	}



	// 执行清理

	tr.executeTeardown(testCase.Teardown)



	if result.Status == "pending" {

		result.Status = "passed"

	}



	return result

}



func (tr *TestRunner) RunAllTests(testDir string) ([]*TestSuiteResult, error) {

	var results []*TestSuiteResult



	// 遍历测试目录

	err := filepath.Walk(testDir, func(path string, info os.FileInfo, err error) error {

		if err != nil {

			return err

		}



		// 只处理YAML文件

		if !info.IsDir() && filepath.Ext(path) == ".yaml" {

			result, err := tr.RunTestSuite(path)

			if err != nil {

				return fmt.Errorf("failed to run test suite %s: %w", path, err)

			}

			results = append(results, result)

		}



		return nil

	})



	return results, err

}

// debugPrint 调试打印函数
func debugPrint(format string, args ...interface{}) {
	// 可以在这里添加日志级别控制
	fmt.Printf(format, args...)
}

// safeSetContext 安全地设置断言上下文，只设置可序列化的值
func (tr *TestRunner) safeSetContext(key string, value interface{}) {
	// 只设置基本类型和字符串
	switch v := value.(type) {
	case string, int, int64, float64, bool:
		tr.context.Variables[key] = v
	default:
		// 对于复杂类型，转换为字符串
		tr.context.Variables[key] = fmt.Sprintf("%v", v)
	}
}




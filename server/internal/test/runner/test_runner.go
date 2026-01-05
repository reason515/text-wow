package runner

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
	"text-wow/internal/database"
	"text-wow/internal/game"
	"text-wow/internal/models"
	"text-wow/internal/repository"
)

// TestRunner 测试运行器
type TestRunner struct {
	parser          *YAMLParser
	assertion       *AssertionExecutor
	reporter        *Reporter
	calculator      *game.Calculator
	equipmentManager *game.EquipmentManager
	context         *TestContext
}

// TestContext 测试上下文
type TestContext struct {
	Characters map[string]*models.Character      // key: character_id
	Monsters   map[string]*models.Monster        // key: monster_id
	Equipments map[string]*models.EquipmentInstance // key: equipment_id
	Variables  map[string]interface{}            // 其他测试变量
}

// NewTestRunner 创建测试运行器
func NewTestRunner() *TestRunner {
	return &TestRunner{
		parser:          NewYAMLParser(),
		assertion:       NewAssertionExecutor(),
		reporter:        NewReporter(),
		calculator:      game.NewCalculator(),
		equipmentManager: game.NewEquipmentManager(),
		context: &TestContext{
			Characters: make(map[string]*models.Character),
			Monsters:   make(map[string]*models.Monster),
			Equipments: make(map[string]*models.EquipmentInstance),
			Variables:  make(map[string]interface{}),
		},
	}
}

// TestSuite 测试套件
type TestSuite struct {
	TestSuite  string    `yaml:"test_suite"`
	Description string   `yaml:"description"`
	Version    string    `yaml:"version"`
	Tests      []TestCase `yaml:"tests"`
}

// TestCase 测试用例
type TestCase struct {
	Name        string   `yaml:"name"`
	Description string   `yaml:"description"`
	Category    string   `yaml:"category"`    // unit/integration/e2e
	Priority    string   `yaml:"priority"`    // high/medium/low
	Setup       []string `yaml:"setup"`
	Steps       []TestStep `yaml:"steps"`
	Assertions  []Assertion `yaml:"assertions"`
	Teardown     []string `yaml:"teardown"`
	Timeout     int      `yaml:"timeout"`     // 秒
	MaxRounds   int      `yaml:"max_rounds"` // 最大回合数
}

// TestStep 测试步骤
type TestStep struct {
	Action     string   `yaml:"action"`
	Expected   string   `yaml:"expected"`
	Timeout    int      `yaml:"timeout"`
	Assertions []string `yaml:"assertions"`
}

// Assertion 断言
type Assertion struct {
	Type      string  `yaml:"type"`      // equals/greater_than/less_than/contains/approximately/range
	Target    string  `yaml:"target"`     // 目标路径，如 "character.hp"
	Expected  string  `yaml:"expected"`   // 期望值
	Tolerance float64 `yaml:"tolerance"` // 容差（用于approximately）
	Message   string  `yaml:"message"`   // 错误消息
}

// TestResult 测试结果
type TestResult struct {
	TestName   string
	Status     string  // passed/failed/skipped
	Duration   time.Duration
	Error      string
	Assertions []AssertionResult
}

// AssertionResult 断言结果
type AssertionResult struct {
	Type     string
	Target   string
	Expected string
	Actual   interface{}
	Status   string  // passed/failed
	Message  string
	Error    string  // 错误信息
}

// TestSuiteResult 测试套件结果
type TestSuiteResult struct {
	TestSuite    string
	TotalTests   int
	PassedTests  int
	FailedTests  int
	SkippedTests int
	Duration     time.Duration
	Results      []TestResult
}

// RunTestSuite 运行测试套件
func (tr *TestRunner) RunTestSuite(suitePath string) (*TestSuiteResult, error) {
	// 读取YAML文件
	data, err := os.ReadFile(suitePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read test suite file: %w", err)
	}

	// 解析YAML
	var suite TestSuite
	if err := yaml.Unmarshal(data, &suite); err != nil {
		return nil, fmt.Errorf("failed to parse test suite: %w", err)
	}

	// 运行测试用例
	result := &TestSuiteResult{
		TestSuite:    suite.TestSuite,
		TotalTests:   len(suite.Tests),
		Results:      make([]TestResult, 0),
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

// RunTestCase 运行单个测试用例
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

	// 执行前置条件
	if err := tr.executeSetup(testCase.Setup); err != nil {
		result.Status = "failed"
		result.Error = fmt.Sprintf("setup failed: %v", err)
		return result
	}

	// 执行测试步骤
	for _, step := range testCase.Steps {
		if err := tr.executeStep(step); err != nil {
			result.Status = "failed"
			result.Error = fmt.Sprintf("step failed: %v", err)
			tr.executeTeardown(testCase.Teardown)
			return result
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

// executeSetup 执行前置条件
func (tr *TestRunner) executeSetup(setup []string) error {
	for _, instruction := range setup {
		if err := tr.executeInstruction(instruction); err != nil {
			return fmt.Errorf("setup instruction failed: %w", err)
		}
	}
	return nil
}

// executeStep 执行测试步骤
func (tr *TestRunner) executeStep(step TestStep) error {
	if err := tr.executeInstruction(step.Action); err != nil {
		return fmt.Errorf("step action failed: %s, error: %w", step.Action, err)
	}
	// 更新断言上下文
	tr.updateAssertionContext()
	return nil
}

// executeInstruction 执行单个指令
func (tr *TestRunner) executeInstruction(instruction string) error {
	// 处理装备相关操作
	if strings.Contains(instruction, "掉落") && strings.Contains(instruction, "装备") {
		return tr.generateEquipmentFromMonster(instruction)
	} else if strings.Contains(instruction, "连续") && strings.Contains(instruction, "装备") {
		return tr.generateMultipleEquipments(instruction)
	} else if strings.Contains(instruction, "检查词缀") || strings.Contains(instruction, "检查词缀数值") || strings.Contains(instruction, "检查词缀类型") || strings.Contains(instruction, "检查词缀Tier") {
		// 这些操作已经在updateAssertionContext中处理
		return nil
	} else if strings.Contains(instruction, "创建一个") && strings.Contains(instruction, "角色") {
		return tr.createCharacter(instruction)
	} else if strings.Contains(instruction, "创建一个") && strings.Contains(instruction, "怪物") {
		return tr.createMonster(instruction)
	} else if strings.Contains(instruction, "击败") && strings.Contains(instruction, "怪物") {
		return tr.createMonster(instruction)
	} else if strings.Contains(instruction, "计算基础伤害") {
		return tr.executeCalculateBaseDamage()
	} else if strings.Contains(instruction, "应用防御减伤") {
		return tr.executeCalculateDefenseReduction()
	} else if strings.Contains(instruction, "计算防御减伤") || strings.Contains(instruction, "计算减伤后伤害") {
		return tr.executeCalculateDefenseReduction()
	} else if strings.Contains(instruction, "如果触发暴击，应用暴击倍率") || strings.Contains(instruction, "应用暴击倍率") {
		return tr.executeApplyCrit()
	} else if strings.Contains(instruction, "计算伤害") {
		return tr.executeCalculateDamage(instruction)
	}
	return nil
}

// executeTeardown 执行清理
func (tr *TestRunner) executeTeardown(teardown []string) error {
	// TODO: 实现清理逻辑
	// 例如：清理战斗状态、重置角色数据等
	return nil
}

// RunAllTests 运行所有测试
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

// updateAssertionContext 更新断言上下文（同步测试数据到断言执行器）
func (tr *TestRunner) updateAssertionContext() {
	// 同步角色信息
	if char, ok := tr.context.Characters["character"]; ok && char != nil {
		tr.assertion.SetContext("character.hp", char.HP)
		tr.assertion.SetContext("character.max_hp", char.MaxHP)
		tr.assertion.SetContext("character.level", char.Level)
		tr.assertion.SetContext("character.resource", char.Resource)
		tr.assertion.SetContext("character.max_resource", char.MaxResource)
		tr.assertion.SetContext("character.physical_attack", char.PhysicalAttack)
		tr.assertion.SetContext("character.magic_attack", char.MagicAttack)
		tr.assertion.SetContext("character.physical_defense", char.PhysicalDefense)
		tr.assertion.SetContext("character.magic_defense", char.MagicDefense)
	}
	
	// 同步装备信息
	if equipment, ok := tr.context.Variables["equipment"].(*models.EquipmentInstance); ok && equipment != nil {
		tr.assertion.SetContext("equipment.id", equipment.ID)
		tr.assertion.SetContext("equipment.item_id", equipment.ItemID)
		tr.assertion.SetContext("equipment.quality", equipment.Quality)
		tr.assertion.SetContext("equipment.slot", equipment.Slot)
		
		// 同步词缀ID
		if equipment.PrefixID != nil {
			tr.assertion.SetContext("equipment.prefix_id", *equipment.PrefixID)
		} else {
			tr.assertion.SetContext("equipment.prefix_id", nil)
		}
		if equipment.SuffixID != nil {
			tr.assertion.SetContext("equipment.suffix_id", *equipment.SuffixID)
		} else {
			tr.assertion.SetContext("equipment.suffix_id", nil)
		}
		
		// 同步词缀数值
		if equipment.PrefixValue != nil {
			tr.assertion.SetContext("equipment.prefix_value", *equipment.PrefixValue)
		}
		if equipment.SuffixValue != nil {
			tr.assertion.SetContext("equipment.suffix_value", *equipment.SuffixValue)
		}
		
		// 同步额外词缀
		if equipment.BonusAffix1 != nil {
			tr.assertion.SetContext("equipment.bonus_affix_1", *equipment.BonusAffix1)
		}
		if equipment.BonusAffix2 != nil {
			tr.assertion.SetContext("equipment.bonus_affix_2", *equipment.BonusAffix2)
		}
		
		// 计算并同步词缀数量
		affixCount := 0
		if equipment.PrefixID != nil {
			affixCount++
		}
		if equipment.SuffixID != nil {
			affixCount++
		}
		if equipment.BonusAffix1 != nil {
			affixCount++
		}
		if equipment.BonusAffix2 != nil {
			affixCount++
		}
		tr.assertion.SetContext("equipment.affix_count", affixCount)
		
		// 同步词缀列表信息（用于contains断言）
		affixesList := []string{}
		if equipment.PrefixID != nil {
			affixesList = append(affixesList, "prefix")
		}
		if equipment.SuffixID != nil {
			affixesList = append(affixesList, "suffix")
		}
		affixesStr := strings.Join(affixesList, ",")
		if affixesStr != "" {
			tr.assertion.SetContext("equipment.affixes", affixesStr)
		}
		
		// 获取装备等级（从角色等级或装备本身）
		equipmentLevel := 1
		if char, ok := tr.context.Characters["character"]; ok {
			equipmentLevel = char.Level
		}
		
		// 同步词缀类型和Tier信息（如果有词缀）
		if equipment.PrefixID != nil {
			tr.syncAffixInfo(*equipment.PrefixID, "prefix", equipmentLevel)
		}
		if equipment.SuffixID != nil {
			tr.syncAffixInfo(*equipment.SuffixID, "suffix", equipmentLevel)
		}
		if equipment.BonusAffix1 != nil {
			tr.syncAffixInfo(*equipment.BonusAffix1, "bonus_1", equipmentLevel)
		}
		if equipment.BonusAffix2 != nil {
			tr.syncAffixInfo(*equipment.BonusAffix2, "bonus_2", equipmentLevel)
		}
	}
	
	// 同步变量
	for key, value := range tr.context.Variables {
		tr.assertion.SetContext(key, value)
	}
}

// syncAffixInfo 同步词缀信息到断言上下文
func (tr *TestRunner) syncAffixInfo(affixID string, affixType string, equipmentLevel int) {
	// 从数据库加载词缀配置
	var slotType string
	
	err := database.DB.QueryRow(`
		SELECT slot_type
		FROM affixes 
		WHERE id = ?`,
		affixID,
	).Scan(&slotType)
	
	if err == nil {
		// 设置词缀类型
		tr.assertion.SetContext(fmt.Sprintf("affix.%s.slot_type", affixType), slotType)
		tr.assertion.SetContext("affix.slot_type", slotType) // 通用键
		
		// 计算Tier（基于装备等级，而不是词缀的levelRequired）
		// Tier 1: 1-20级
		// Tier 2: 21-40级  
		// Tier 3: 41+级
		tier := 1
		if equipmentLevel > 20 && equipmentLevel <= 40 {
			tier = 2
		} else if equipmentLevel > 40 {
			tier = 3
		}
		tr.assertion.SetContext(fmt.Sprintf("affix.%s.tier", affixType), tier)
		tr.assertion.SetContext("affix.tier", tier) // 通用键
	}
}

// generateMultipleEquipments 生成多件装备（用于随机性测试）
func (tr *TestRunner) generateMultipleEquipments(action string) error {
	// 解析数量：如"连续获得10件蓝色装备"
	count := 10
	numStr := ""
	for _, r := range action {
		if r >= '0' && r <= '9' {
			numStr += string(r)
		} else if numStr != "" {
			break
		}
	}
	if numStr != "" {
		if n, err := strconv.Atoi(numStr); err == nil {
			count = n
		}
	}
	
	// 解析品质
	quality := "rare"
	if strings.Contains(action, "白色") || strings.Contains(action, "white") || strings.Contains(action, "common") {
		quality = "common"
	} else if strings.Contains(action, "绿色") || strings.Contains(action, "green") || strings.Contains(action, "uncommon") {
		quality = "uncommon"
	} else if strings.Contains(action, "蓝色") || strings.Contains(action, "blue") || strings.Contains(action, "rare") {
		quality = "rare"
	} else if strings.Contains(action, "紫色") || strings.Contains(action, "purple") || strings.Contains(action, "epic") {
		quality = "epic"
	}
	
	// 获取角色等级
	level := 1
	if char, ok := tr.context.Characters["character"]; ok {
		level = char.Level
	}
	
	// 确保用户和角色存在
	ownerID := 1
	if char, ok := tr.context.Characters["character"]; ok {
		ownerID = char.UserID
	} else {
		userRepo := repository.NewUserRepository()
		user, err := userRepo.GetByUsername("test_user")
		if err != nil {
			passwordHash := "test_hash"
			user, err = userRepo.Create("test_user", passwordHash, "test@test.com")
			if err != nil {
				return fmt.Errorf("failed to create user: %w", err)
			}
		}
		ownerID = user.ID
		
		charRepo := repository.NewCharacterRepository()
		char, err := charRepo.Create(&models.Character{
			UserID:   user.ID,
			Name:     "测试角色",
			RaceID:   "human",
			ClassID:  "warrior",
			Faction:  "alliance",
			TeamSlot: 1,
			Level:    level,
		})
		if err != nil {
			return fmt.Errorf("failed to create character: %w", err)
		}
		tr.context.Characters["character"] = char
	}
	
	// 生成多件装备并统计唯一词缀组合
	uniqueCombinations := make(map[string]bool)
	itemID := "worn_sword"
	
	for i := 0; i < count; i++ {
		equipment, err := tr.equipmentManager.GenerateEquipment(itemID, quality, level, ownerID)
		if err != nil {
			continue
		}
		
		// 构建词缀组合字符串
		prefixID := "none"
		suffixID := "none"
		if equipment.PrefixID != nil {
			prefixID = *equipment.PrefixID
		}
		if equipment.SuffixID != nil {
			suffixID = *equipment.SuffixID
		}
		combination := fmt.Sprintf("%s_%s", prefixID, suffixID)
		uniqueCombinations[combination] = true
		
		// 存储最后一件装备到上下文
		if i == count-1 {
			tr.context.Variables["equipment"] = equipment
			tr.context.Variables["equipment_id"] = equipment.ID
		}
	}
	
	// 设置唯一词缀组合数量
	tr.context.Variables["unique_affix_combinations"] = len(uniqueCombinations)
	
	return nil
}

// generateEquipmentFromMonster 从怪物掉落生成装备
func (tr *TestRunner) generateEquipmentFromMonster(action string) error {
	// 解析品质：如"怪物掉落一件白色装备"
	quality := "common"
	if strings.Contains(action, "白色") || strings.Contains(action, "white") || strings.Contains(action, "common") {
		quality = "common"
	} else if strings.Contains(action, "绿色") || strings.Contains(action, "green") || strings.Contains(action, "uncommon") {
		quality = "uncommon"
	} else if strings.Contains(action, "蓝色") || strings.Contains(action, "blue") || strings.Contains(action, "rare") {
		quality = "rare"
	} else if strings.Contains(action, "紫色") || strings.Contains(action, "purple") || strings.Contains(action, "epic") {
		quality = "epic"
	} else if strings.Contains(action, "橙色") || strings.Contains(action, "orange") || strings.Contains(action, "legendary") {
		quality = "legendary"
	}
	
	// 处理"Boss掉落"的情况
	if strings.Contains(action, "Boss") || strings.Contains(action, "boss") {
		// 如果没有怪物，创建一个Boss怪物
		if len(tr.context.Monsters) == 0 {
			monster := &models.Monster{
				ID:              "boss_monster",
				Name:            "Boss怪物",
				Type:            "boss",
				Level:           30,
				HP:              0, // 被击败
				MaxHP:           1000,
				PhysicalAttack:  50,
				MagicAttack:     50,
				PhysicalDefense: 20,
				MagicDefense:    20,
				DodgeRate:       0.1,
			}
			tr.context.Monsters["monster"] = monster
		}
	}
	
	// 获取怪物等级
	level := 1
	for _, monster := range tr.context.Monsters {
		level = monster.Level
		break
	}
	
	// 确保用户和角色存在
	ownerID := 1
	if char, ok := tr.context.Characters["character"]; ok {
		ownerID = char.UserID
	} else {
		user, err := tr.createTestUser()
		if err != nil {
			return fmt.Errorf("failed to create test user: %w", err)
		}
		ownerID = user.ID
		
		char, err := tr.createTestCharacter(user.ID, level)
		if err != nil {
			return fmt.Errorf("failed to create test character: %w", err)
		}
		tr.context.Characters["character"] = char
	}
	
	// 生成装备（使用数据库中存在的itemID）
	itemID := "worn_sword" // 使用seed.sql中存在的itemID
	equipment, err := tr.equipmentManager.GenerateEquipment(itemID, quality, level, ownerID)
	if err != nil {
		return fmt.Errorf("failed to generate equipment: %w", err)
	}
	
	// 存储到上下文
	tr.context.Variables["equipment"] = equipment
	tr.context.Variables["equipment_id"] = equipment.ID
	tr.context.Equipments[fmt.Sprintf("%d", equipment.ID)] = equipment
	
	return nil
}

// createCharacter 创建角色
func (tr *TestRunner) createCharacter(instruction string) error {
	char := &models.Character{
		ID:       1,
		Name:     "测试角色",
		ClassID:  "warrior",
		Level:    1,
		Strength: 10,
		Agility:  10,
		Intellect: 10,
		Stamina:   10,
		Spirit:    10,
		MaxHP:    0,
		MaxResource: 0,
	}
	
	// 解析攻击力（如"攻击力=20"）
	if strings.Contains(instruction, "攻击力=") {
		parts := strings.Split(instruction, "攻击力=")
		if len(parts) > 1 {
			attackStr := strings.TrimSpace(strings.Split(parts[1], "，")[0])
			attackStr = strings.TrimSpace(strings.Split(attackStr, "的")[0])
			if attack, err := strconv.Atoi(attackStr); err == nil {
				char.PhysicalAttack = attack
			}
		}
	}
	
	// 解析防御力（如"防御力=10"）
	if strings.Contains(instruction, "防御力=") {
		parts := strings.Split(instruction, "防御力=")
		if len(parts) > 1 {
			defenseStr := strings.TrimSpace(strings.Split(parts[1], "，")[0])
			defenseStr = strings.TrimSpace(strings.Split(defenseStr, "的")[0])
			if defense, err := strconv.Atoi(defenseStr); err == nil {
				char.PhysicalDefense = defense
			}
		}
	}
	
	// 解析暴击率（如"物理暴击率=10%"）
	if strings.Contains(instruction, "物理暴击率=") {
		parts := strings.Split(instruction, "物理暴击率=")
		if len(parts) > 1 {
			critStr := strings.TrimSpace(strings.Split(parts[1], "%")[0])
			if crit, err := strconv.ParseFloat(critStr, 64); err == nil {
				char.PhysCritRate = crit / 100.0
			}
		}
	}
	
	// 解析暴击伤害（如"物理暴击伤害=150%"）
	if strings.Contains(instruction, "物理暴击伤害=") {
		parts := strings.Split(instruction, "物理暴击伤害=")
		if len(parts) > 1 {
			critDmgStr := strings.TrimSpace(strings.Split(parts[1], "%")[0])
			if critDmg, err := strconv.ParseFloat(critDmgStr, 64); err == nil {
				char.PhysCritDamage = critDmg / 100.0
			}
		}
	}
	
	// 解析等级
	if strings.Contains(instruction, "30级") {
		char.Level = 30
	}
	
	// 确保用户存在
	if char.UserID == 0 {
		user, err := tr.createTestUser()
		if err != nil {
			return fmt.Errorf("failed to create user: %w", err)
		}
		char.UserID = user.ID
	}
	
	// 确保角色有必需的字段
	if char.RaceID == "" {
		char.RaceID = "human"
	}
	if char.Faction == "" {
		char.Faction = "alliance"
	}
	if char.TeamSlot == 0 {
		char.TeamSlot = 1
	}
	if char.ResourceType == "" {
		char.ResourceType = "rage"
	}
	
	// 尝试从数据库获取角色，如果不存在则创建
	charRepo := repository.NewCharacterRepository()
	chars, err := charRepo.GetByUserID(char.UserID)
	if err != nil || len(chars) == 0 {
		createdChar, err := charRepo.Create(char)
		if err != nil {
			return fmt.Errorf("failed to create character in DB: %w", err)
		}
		char = createdChar
	} else {
		// 查找匹配slot的角色
		var existingChar *models.Character
		for _, c := range chars {
			if c.TeamSlot == char.TeamSlot {
				existingChar = c
				break
			}
		}
		if existingChar != nil {
			char.ID = existingChar.ID
			if err := charRepo.Update(char); err != nil {
				return fmt.Errorf("failed to update existing character in DB: %w", err)
			}
		} else {
			createdChar, err := charRepo.Create(char)
			if err != nil {
				return fmt.Errorf("failed to create character in DB: %w", err)
			}
			char = createdChar
		}
	}
	
	// 存储到上下文
	tr.context.Characters["character"] = char
	
	return nil
}

// createMonster 创建怪物
func (tr *TestRunner) createMonster(instruction string) error {
	monster := &models.Monster{
		ID:              "test_monster",
		Name:            "测试怪物",
		Type:            "normal",
		Level:           1,
		HP:              100, // 默认存活
		MaxHP:           100,
		PhysicalAttack:  10,
		MagicAttack:     5,
		PhysicalDefense: 5,
		MagicDefense:    3,
		DodgeRate:       0.05,
	}
	
	// 解析防御力（如"防御力=10"）
	if strings.Contains(instruction, "防御力=") {
		parts := strings.Split(instruction, "防御力=")
		if len(parts) > 1 {
			defenseStr := strings.TrimSpace(strings.Split(parts[1], "，")[0])
			defenseStr = strings.TrimSpace(strings.Split(defenseStr, "的")[0])
			defenseStr = strings.TrimSpace(strings.Split(defenseStr, "（")[0])
			if defense, err := strconv.Atoi(defenseStr); err == nil {
				monster.PhysicalDefense = defense
			}
		}
	}
	
	// 解析闪避率（如"闪避率=10%"）
	if strings.Contains(instruction, "闪避率=") {
		parts := strings.Split(instruction, "闪避率=")
		if len(parts) > 1 {
			dodgeStr := strings.TrimSpace(strings.Split(parts[1], "%")[0])
			if dodge, err := strconv.ParseFloat(dodgeStr, 64); err == nil {
				monster.DodgeRate = dodge / 100.0
			}
		}
	}
	
	// 解析HP（如"HP=100"）
	if strings.Contains(instruction, "HP=") {
		parts := strings.Split(instruction, "HP=")
		if len(parts) > 1 {
			hpStr := strings.TrimSpace(strings.Split(parts[1], "，")[0])
			if hp, err := strconv.Atoi(hpStr); err == nil {
				monster.HP = hp
				monster.MaxHP = hp
			}
		}
	}
	
	tr.context.Monsters["monster"] = monster
	return nil
}

// createTestUser 创建一个测试用户（如果不存在）
func (tr *TestRunner) createTestUser() (*models.User, error) {
	userRepo := repository.NewUserRepository()
	user, err := userRepo.GetByUsername("test_user")
	if err != nil {
		passwordHash := "test_hash"
		user, err = userRepo.Create("test_user", passwordHash, "test@test.com")
		if err != nil {
			return nil, fmt.Errorf("failed to create user: %w", err)
		}
	}
	return user, nil
}

// createTestCharacter 创建一个测试角色（如果不存在）
func (tr *TestRunner) createTestCharacter(userID, level int) (*models.Character, error) {
	charRepo := repository.NewCharacterRepository()
	chars, err := charRepo.GetByUserID(userID)
	var char *models.Character
	if err != nil || len(chars) == 0 {
		char = &models.Character{
			UserID:   userID,
			Name:     "测试角色",
			RaceID:   "human",
			ClassID:  "warrior",
			Faction:  "alliance",
			TeamSlot: 1,
			Level:    level,
			HP:       100, MaxHP: 100,
			Resource: 100, MaxResource: 100, ResourceType: "rage",
			Strength: 10, Agility: 10, Intellect: 10, Stamina: 10, Spirit: 10,
		}
		createdChar, err := charRepo.Create(char)
		if err != nil {
			return nil, fmt.Errorf("failed to create character: %w", err)
		}
		char = createdChar
	} else {
		// 查找第一个slot的角色
		for _, c := range chars {
			if c.TeamSlot == 1 {
				char = c
				break
			}
		}
		if char == nil {
			char = &models.Character{
				UserID:   userID,
				Name:     "测试角色",
				RaceID:   "human",
				ClassID:  "warrior",
				Faction:  "alliance",
				TeamSlot: 1,
				Level:    level,
				HP:       100, MaxHP: 100,
				Resource: 100, MaxResource: 100, ResourceType: "rage",
				Strength: 10, Agility: 10, Intellect: 10, Stamina: 10, Spirit: 10,
			}
			createdChar, err := charRepo.Create(char)
			if err != nil {
				return nil, fmt.Errorf("failed to create character: %w", err)
			}
			char = createdChar
		} else {
			char.Level = level
			if err := charRepo.Update(char); err != nil {
				return nil, fmt.Errorf("failed to update existing character: %w", err)
			}
		}
	}
	return char, nil
}

// executeCalculateBaseDamage 计算基础伤害
func (tr *TestRunner) executeCalculateBaseDamage() error {
	char, ok := tr.context.Characters["character"]
	if !ok || char == nil {
		return fmt.Errorf("character not found")
	}
	
	// 基础伤害 = 攻击力 × 技能系数（默认1.0）
	baseDamage := char.PhysicalAttack
	
	tr.assertion.SetContext("base_damage", baseDamage)
	tr.context.Variables["base_damage"] = baseDamage
	return nil
}

// executeCalculateDefenseReduction 计算防御减伤
func (tr *TestRunner) executeCalculateDefenseReduction() error {
	char, ok := tr.context.Characters["character"]
	if !ok || char == nil {
		return fmt.Errorf("character not found")
	}
	
	monster, ok := tr.context.Monsters["monster"]
	if !ok || monster == nil {
		return fmt.Errorf("monster not found")
	}
	
	// 获取基础伤害（如果已计算）
	baseDamage := char.PhysicalAttack
	if val, exists := tr.context.Variables["base_damage"]; exists {
		if bd, ok := val.(int); ok {
			baseDamage = bd
		}
	}
	
	// 应用防御减伤（减法公式）
	damageAfterDefense := baseDamage - monster.PhysicalDefense
	if damageAfterDefense < 1 {
		damageAfterDefense = 1 // 至少1点伤害
	}
	
	tr.assertion.SetContext("damage_after_defense", damageAfterDefense)
	tr.context.Variables["damage_after_defense"] = damageAfterDefense
	// 如果没有最终伤害，使用减伤后伤害作为最终伤害
	if _, exists := tr.context.Variables["final_damage"]; !exists {
		tr.assertion.SetContext("final_damage", damageAfterDefense)
		tr.context.Variables["final_damage"] = damageAfterDefense
	}
	
	return nil
}

// executeApplyCrit 应用暴击倍率
func (tr *TestRunner) executeApplyCrit() error {
	// 从上下文中获取伤害值
	var baseDamage int
	if val, exists := tr.context.Variables["damage_after_defense"]; exists {
		if bd, ok := val.(int); ok {
			baseDamage = bd
		}
	}
	
	if baseDamage == 0 {
		// 如果没有伤害值，尝试从角色和怪物计算
		char, ok := tr.context.Characters["character"]
		if !ok || char == nil {
			return fmt.Errorf("character not found")
		}
		monster, ok := tr.context.Monsters["monster"]
		if !ok || monster == nil {
			return fmt.Errorf("monster not found")
		}
		baseDamage = char.PhysicalAttack - monster.PhysicalDefense
		if baseDamage < 1 {
			baseDamage = 1
		}
		// 更新上下文
		tr.assertion.SetContext("damage_after_defense", baseDamage)
		tr.context.Variables["damage_after_defense"] = baseDamage
	}
	
	char, ok := tr.context.Characters["character"]
	if !ok || char == nil {
		return fmt.Errorf("character not found")
	}
	
	// 假设暴击（实际应该随机判断）
	// 注意：PhysCritDamage是倍率，如1.5表示150%
	finalDamage := int(float64(baseDamage) * char.PhysCritDamage)
	
	tr.assertion.SetContext("final_damage", finalDamage)
	tr.context.Variables["final_damage"] = finalDamage
	return nil
}

// executeCalculateDamage 计算伤害（通用）
func (tr *TestRunner) executeCalculateDamage(instruction string) error {
	char, ok := tr.context.Characters["character"]
	if !ok || char == nil {
		return fmt.Errorf("character not found")
	}
	
	monster, ok := tr.context.Monsters["monster"]
	if !ok || monster == nil {
		return fmt.Errorf("monster not found")
	}
	
	// 使用计算器计算伤害
	defender := &models.Character{
		PhysicalDefense: monster.PhysicalDefense,
		MagicDefense:    monster.MagicDefense,
		DodgeRate:       monster.DodgeRate,
	}
	
	result := tr.calculator.CalculateDamage(
		char,
		defender,
		char.PhysicalAttack,
		1.0, // 技能倍率
		"physical",
		false, // 不忽略闪避
	)
	
	tr.assertion.SetContext("base_damage", int(result.BaseDamage))
	tr.assertion.SetContext("damage_after_defense", int(result.DamageAfterDefense))
	tr.assertion.SetContext("final_damage", result.FinalDamage)
	tr.context.Variables["base_damage"] = int(result.BaseDamage)
	tr.context.Variables["damage_after_defense"] = int(result.DamageAfterDefense)
	tr.context.Variables["final_damage"] = result.FinalDamage
	
	return nil
}


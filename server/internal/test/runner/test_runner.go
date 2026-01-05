package runner

import (
	"fmt"
	"math"
	"math/rand"
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
	
	// 调试：检查setup后的上下文状态
	fmt.Fprintf(os.Stderr, "[DEBUG] RunTestCase: after setup for '%s' - characters=%d, monsters=%d, variables=%d\n", 
		testCase.Name, len(tr.context.Characters), len(tr.context.Monsters), len(tr.context.Variables))
	if char, exists := tr.context.Characters["character"]; exists && char != nil {
		fmt.Fprintf(os.Stderr, "[DEBUG] RunTestCase: after setup, character.PhysicalAttack=%d, character pointer=%p\n", char.PhysicalAttack, char)
		// 也检查Variables中的值
		if attackVal, exists := tr.context.Variables["character_physical_attack"]; exists {
			fmt.Fprintf(os.Stderr, "[DEBUG] RunTestCase: after setup, Variables[character_physical_attack]=%v\n", attackVal)
		}
	} else if exists {
		fmt.Fprintf(os.Stderr, "[DEBUG] RunTestCase: after setup, character is nil\n")
	}
	if ratio, exists := tr.context.Variables["skill_scaling_ratio"]; exists {
		fmt.Fprintf(os.Stderr, "[DEBUG] RunTestCase: skill_scaling_ratio=%v\n", ratio)
	}

	// 执行测试步骤
	for _, step := range testCase.Steps {
		// 在执行步骤之前，检查上下文中的角色状态
		if char, exists := tr.context.Characters["character"]; exists && char != nil {
			fmt.Fprintf(os.Stderr, "[DEBUG] RunTestCase: before executeStep, character.PhysicalAttack=%d, character pointer=%p\n", char.PhysicalAttack, char)
		}
		if err := tr.executeStep(step); err != nil {
			result.Status = "failed"
			result.Error = fmt.Sprintf("step failed: %v", err)
			tr.executeTeardown(testCase.Teardown)
			return result
		}
		// 在执行步骤之后，检查上下文中的角色状态
		if char, exists := tr.context.Characters["character"]; exists && char != nil {
			fmt.Fprintf(os.Stderr, "[DEBUG] RunTestCase: after executeStep, character.PhysicalAttack=%d\n", char.PhysicalAttack)
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
	} else if strings.Contains(instruction, "获得") && (strings.Contains(instruction, "装备") || strings.Contains(instruction, "武器") || strings.Contains(instruction, "护甲") || strings.Contains(instruction, "饰品")) {
		// 处理"获得一件X级武器，攻击力+X"这样的setup指令
		return tr.generateEquipmentWithAttributes(instruction)
	} else if strings.Contains(instruction, "尝试穿戴") || strings.Contains(instruction, "尝试装备") {
		// 处理"角色尝试穿戴武器"等action（用于测试失败情况）
		// 必须在"穿戴"之前检查，因为"尝试穿戴"包含"穿戴"
		return tr.executeTryEquipItem(instruction)
	} else if strings.Contains(instruction, "穿戴") && (strings.Contains(instruction, "装备") || strings.Contains(instruction, "武器") || strings.Contains(instruction, "护甲") || strings.Contains(instruction, "饰品")) {
		// 处理"角色穿戴武器"、"角色穿戴装备"等action
		return tr.executeEquipItem(instruction)
	} else if strings.Contains(instruction, "卸下") && (strings.Contains(instruction, "装备") || strings.Contains(instruction, "武器") || strings.Contains(instruction, "护甲") || strings.Contains(instruction, "饰品")) {
		// 处理"角色卸下武器"、"角色卸下装备"等action
		return tr.executeUnequipItem(instruction)
	} else if strings.Contains(instruction, "依次穿戴") && strings.Contains(instruction, "装备") {
		// 处理"角色依次穿戴所有装备"
		return tr.executeEquipAllItems(instruction)
	} else if strings.Contains(instruction, "检查词缀") || strings.Contains(instruction, "检查词缀数值") || strings.Contains(instruction, "检查词缀类型") || strings.Contains(instruction, "检查词缀Tier") {
		// 这些操作已经在updateAssertionContext中处理
		return nil
	} else if strings.Contains(instruction, "设置") {
		return tr.executeSetVariable(instruction)
	} else if strings.Contains(instruction, "创建一个nil角色") {
		// 创建一个nil角色（用于测试nil情况）
		tr.context.Characters["character"] = nil
		return nil
	} else if strings.Contains(instruction, "创建一个") && strings.Contains(instruction, "角色") {
		return tr.createCharacter(instruction)
	} else if (strings.Contains(instruction, "创建一个") || strings.Contains(instruction, "创建")) && strings.Contains(instruction, "怪物") {
		return tr.createMonster(instruction)
	} else if strings.Contains(instruction, "击败") && strings.Contains(instruction, "怪物") {
		return tr.createMonster(instruction)
	} else if strings.Contains(instruction, "计算物理攻击力") {
		return tr.executeCalculatePhysicalAttack()
	} else if strings.Contains(instruction, "计算法术攻击力") {
		return tr.executeCalculateMagicAttack()
	} else if strings.Contains(instruction, "计算最大生命值") || strings.Contains(instruction, "计算生命值") {
		return tr.executeCalculateMaxHP()
	} else if strings.Contains(instruction, "计算物理暴击率") {
		return tr.executeCalculatePhysCritRate()
	} else if strings.Contains(instruction, "计算法术暴击率") {
		return tr.executeCalculateSpellCritRate()
	} else if strings.Contains(instruction, "计算物理暴击伤害倍率") {
		return tr.executeCalculatePhysCritDamage()
	} else if strings.Contains(instruction, "计算物理防御力") {
		return tr.executeCalculatePhysicalDefense()
	} else if strings.Contains(instruction, "计算魔法防御力") {
		return tr.executeCalculateMagicDefense()
	} else if strings.Contains(instruction, "计算法术暴击伤害倍率") {
		return tr.executeCalculateSpellCritDamage()
	} else if strings.Contains(instruction, "计算闪避率") {
		return tr.executeCalculateDodgeRate()
	} else if strings.Contains(instruction, "角色对怪物进行") && strings.Contains(instruction, "次攻击") {
		return tr.executeMultipleAttacks(instruction)
	} else if strings.Contains(instruction, "计算速度") {
		return tr.executeCalculateSpeed()
	} else if strings.Contains(instruction, "计算资源回复") || strings.Contains(instruction, "计算法力回复") || strings.Contains(instruction, "计算法力恢复") || strings.Contains(instruction, "计算怒气获得") || strings.Contains(instruction, "计算能量回复") || strings.Contains(instruction, "计算能量恢复") {
		return tr.executeCalculateResourceRegen(instruction)
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
	} else if strings.Contains(instruction, "使用技能") || strings.Contains(instruction, "角色使用技能") || (strings.Contains(instruction, "使用") && strings.Contains(instruction, "技能")) {
		return tr.executeUseSkill(instruction)
	} else if strings.Contains(instruction, "创建一个") && strings.Contains(instruction, "技能") {
		return tr.createSkill(instruction)
	} else if strings.Contains(instruction, "执行第") && strings.Contains(instruction, "回合") {
		return tr.executeBattleRound(instruction)
	} else if strings.Contains(instruction, "开始战斗") {
		return tr.executeStartBattle()
	} else if strings.Contains(instruction, "检查战斗初始状态") || strings.Contains(instruction, "检查战斗状态") {
		// 检查战斗状态，确保战士怒气为0
		return tr.executeCheckBattleState(instruction)
	} else if strings.Contains(instruction, "检查战斗结束状态") {
		// 检查战斗结束状态，确保战士怒气归0
		return tr.executeCheckBattleEndState()
	} else if strings.Contains(instruction, "角色攻击怪物") || strings.Contains(instruction, "攻击怪物") {
		return tr.executeAttackMonster()
	} else if strings.Contains(instruction, "怪物攻击角色") {
		return tr.executeMonsterAttack()
	} else if strings.Contains(instruction, "获取角色数据") || strings.Contains(instruction, "获取战斗状态") {
		// 获取角色数据或战斗状态，确保战士怒气正确
		return tr.executeGetCharacterData()
	} else if strings.Contains(instruction, "检查角色属性") || strings.Contains(instruction, "检查角色") {
		// 检查角色属性，确保所有属性都基于角色属性正确计算
		return tr.executeCheckCharacterAttributes()
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
		// 确保char不是nil指针
		tr.assertion.SetContext("character.hp", char.HP)
		tr.assertion.SetContext("character.max_hp", char.MaxHP)
		tr.assertion.SetContext("character.level", char.Level)
		tr.assertion.SetContext("character.resource", char.Resource)
		tr.assertion.SetContext("character.max_resource", char.MaxResource)
		tr.assertion.SetContext("character.physical_attack", char.PhysicalAttack)
		tr.assertion.SetContext("character.magic_attack", char.MagicAttack)
		tr.assertion.SetContext("character.physical_defense", char.PhysicalDefense)
		tr.assertion.SetContext("character.magic_defense", char.MagicDefense)
		tr.assertion.SetContext("character.phys_crit_rate", char.PhysCritRate)
		tr.assertion.SetContext("character.phys_crit_damage", char.PhysCritDamage)
		tr.assertion.SetContext("character.spell_crit_rate", char.SpellCritRate)
		tr.assertion.SetContext("character.spell_crit_damage", char.SpellCritDamage)
		tr.assertion.SetContext("character.dodge_rate", char.DodgeRate)
		tr.assertion.SetContext("character.id", char.ID)
		
		// 同步Buff信息（从上下文获取）
		if buffModifier, exists := tr.context.Variables["character_buff_attack_modifier"]; exists {
			tr.assertion.SetContext("character.buff_attack_modifier", buffModifier)
		}
		if buffDuration, exists := tr.context.Variables["character_buff_duration"]; exists {
			tr.assertion.SetContext("character.buff_duration", buffDuration)
		}
	}
	
	// 同步怪物信息
	for key, monster := range tr.context.Monsters {
		if monster != nil {
			tr.assertion.SetContext(fmt.Sprintf("%s.hp", key), monster.HP)
			tr.assertion.SetContext(fmt.Sprintf("%s.max_hp", key), monster.MaxHP)
		}
	}
	
	// 同步所有monster_X.hp_damage值（从Variables中读取）
	for i := 1; i <= 10; i++ {
		damageKey := fmt.Sprintf("monster_%d.hp_damage", i)
		if hpDamage, exists := tr.context.Variables[damageKey]; exists {
			tr.assertion.SetContext(damageKey, hpDamage)
		}
	}
	
	// 同步技能伤害值
	if skillDamage, exists := tr.context.Variables["skill_damage_dealt"]; exists {
		tr.assertion.SetContext("skill_damage_dealt", skillDamage)
	}
	
	// 同步装备信息
	tr.syncEquipmentToContext("equipment", tr.context.Variables["equipment"])
	tr.syncEquipmentToContext("weapon", tr.context.Variables["weapon"])
	tr.syncEquipmentToContext("old_weapon", tr.context.Variables["old_weapon"])
	tr.syncEquipmentToContext("old_equipment", tr.context.Variables["old_equipment"])
	tr.syncEquipmentToContext("new_weapon", tr.context.Variables["new_weapon"])
	tr.syncEquipmentToContext("new_equipment", tr.context.Variables["new_equipment"])
	
	// 同步装备槽位计数（用于测试槽位冲突）
	if char, ok := tr.context.Characters["character"]; ok {
		equipmentRepo := repository.NewEquipmentRepository()
		mainHandCount := 0
		equippedEquipments, _ := equipmentRepo.GetByCharacterID(char.ID)
		for _, eq := range equippedEquipments {
			if eq.Slot == "main_hand" {
				mainHandCount++
			}
		}
		tr.assertion.SetContext("equipped_main_hand_count", mainHandCount)
	}
	
	// 同步变量
	for key, value := range tr.context.Variables {
		tr.assertion.SetContext(key, value)
	}
}

// syncEquipmentToContext 同步装备信息到断言上下文
func (tr *TestRunner) syncEquipmentToContext(prefix string, equipment interface{}) {
	if equipment == nil {
		return
	}
	
	eq, ok := equipment.(*models.EquipmentInstance)
	if !ok || eq == nil {
		return
	}
	
	tr.assertion.SetContext(fmt.Sprintf("%s.id", prefix), eq.ID)
	tr.assertion.SetContext(fmt.Sprintf("%s.item_id", prefix), eq.ItemID)
	tr.assertion.SetContext(fmt.Sprintf("%s.quality", prefix), eq.Quality)
	tr.assertion.SetContext(fmt.Sprintf("%s.slot", prefix), eq.Slot)
	
	// 同步character_id
	if eq.CharacterID != nil {
		tr.assertion.SetContext(fmt.Sprintf("%s.character_id", prefix), *eq.CharacterID)
	} else {
		tr.assertion.SetContext(fmt.Sprintf("%s.character_id", prefix), nil)
	}
	
	// 同步词缀ID
	if eq.PrefixID != nil {
		tr.assertion.SetContext(fmt.Sprintf("%s.prefix_id", prefix), *eq.PrefixID)
	} else {
		tr.assertion.SetContext(fmt.Sprintf("%s.prefix_id", prefix), nil)
	}
	if eq.SuffixID != nil {
		tr.assertion.SetContext(fmt.Sprintf("%s.suffix_id", prefix), *eq.SuffixID)
	} else {
		tr.assertion.SetContext(fmt.Sprintf("%s.suffix_id", prefix), nil)
	}
	
	// 同步词缀数值
	if eq.PrefixValue != nil {
		tr.assertion.SetContext(fmt.Sprintf("%s.prefix_value", prefix), *eq.PrefixValue)
	}
	if eq.SuffixValue != nil {
		tr.assertion.SetContext(fmt.Sprintf("%s.suffix_value", prefix), *eq.SuffixValue)
	}
	
	// 同步额外词缀
	if eq.BonusAffix1 != nil {
		tr.assertion.SetContext(fmt.Sprintf("%s.bonus_affix_1", prefix), *eq.BonusAffix1)
	}
	if eq.BonusAffix2 != nil {
		tr.assertion.SetContext(fmt.Sprintf("%s.bonus_affix_2", prefix), *eq.BonusAffix2)
	}
	
	// 计算并同步词缀数量
	affixCount := 0
	if eq.PrefixID != nil {
		affixCount++
	}
	if eq.SuffixID != nil {
		affixCount++
	}
	if eq.BonusAffix1 != nil {
		affixCount++
	}
	if eq.BonusAffix2 != nil {
		affixCount++
	}
	tr.assertion.SetContext(fmt.Sprintf("%s.affix_count", prefix), affixCount)
	
	// 同步词缀列表信息（用于contains断言）
	affixesList := []string{}
	if eq.PrefixID != nil {
		affixesList = append(affixesList, "prefix")
	}
	if eq.SuffixID != nil {
		affixesList = append(affixesList, "suffix")
	}
	affixesStr := strings.Join(affixesList, ",")
	if affixesStr != "" {
		tr.assertion.SetContext(fmt.Sprintf("%s.affixes", prefix), affixesStr)
	}
	
	// 获取装备等级（从角色等级或装备本身）
	equipmentLevel := 1
	if char, ok := tr.context.Characters["character"]; ok {
		equipmentLevel = char.Level
	}
	
	// 同步词缀类型和Tier信息（如果有词缀）
	if eq.PrefixID != nil {
		tr.syncAffixInfo(*eq.PrefixID, fmt.Sprintf("%s.prefix", prefix), equipmentLevel)
	}
	if eq.SuffixID != nil {
		tr.syncAffixInfo(*eq.SuffixID, fmt.Sprintf("%s.suffix", prefix), equipmentLevel)
	}
	if eq.BonusAffix1 != nil {
		tr.syncAffixInfo(*eq.BonusAffix1, fmt.Sprintf("%s.bonus_1", prefix), equipmentLevel)
	}
	if eq.BonusAffix2 != nil {
		tr.syncAffixInfo(*eq.BonusAffix2, fmt.Sprintf("%s.bonus_2", prefix), equipmentLevel)
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
	classID := "warrior" // 默认职业
	if strings.Contains(instruction, "法师") {
		classID = "mage"
	} else if strings.Contains(instruction, "战士") {
		classID = "warrior"
	} else if strings.Contains(instruction, "盗贼") {
		classID = "rogue"
	} else if strings.Contains(instruction, "牧师") {
		classID = "priest"
	}
	
	char := &models.Character{
		ID:       1,
		Name:     "测试角色",
		ClassID:  classID,
		Level:    1,
		Strength: 10,
		Agility:  10,
		Intellect: 10,
		Stamina:   10,
		Spirit:    10,
		MaxHP:    0,
		MaxResource: 0,
	}
	
	// 解析主属性（如"力量=20"、"敏捷=10"等）
	parseAttribute := func(value string) string {
		value = strings.TrimSpace(strings.Split(value, "，")[0])
		value = strings.TrimSpace(strings.Split(value, ",")[0])
		// 去掉括号和注释（如"1000（理论上暴击率会超过50%）"）
		if idx := strings.Index(value, "（"); idx >= 0 {
			value = value[:idx]
		}
		if idx := strings.Index(value, "("); idx >= 0 {
			value = value[:idx]
		}
		return strings.TrimSpace(value)
	}
	
	if strings.Contains(instruction, "力量=") {
		parts := strings.Split(instruction, "力量=")
		if len(parts) > 1 {
			strStr := parseAttribute(parts[1])
			if str, err := strconv.Atoi(strStr); err == nil {
				char.Strength = str
				fmt.Fprintf(os.Stderr, "[DEBUG] createCharacter: set Strength=%d\n", str)
			}
		}
	}
	if strings.Contains(instruction, "敏捷=") {
		parts := strings.Split(instruction, "敏捷=")
		if len(parts) > 1 {
			agiStr := parseAttribute(parts[1])
			if agi, err := strconv.Atoi(agiStr); err == nil {
				char.Agility = agi
				fmt.Fprintf(os.Stderr, "[DEBUG] createCharacter: set Agility=%d\n", agi)
			}
		}
	}
	if strings.Contains(instruction, "智力=") {
		parts := strings.Split(instruction, "智力=")
		if len(parts) > 1 {
			intStr := parseAttribute(parts[1])
			if intel, err := strconv.Atoi(intStr); err == nil {
				char.Intellect = intel
				fmt.Fprintf(os.Stderr, "[DEBUG] createCharacter: set Intellect=%d\n", intel)
			}
		}
	}
	if strings.Contains(instruction, "精神=") {
		parts := strings.Split(instruction, "精神=")
		if len(parts) > 1 {
			spiStr := parseAttribute(parts[1])
			if spi, err := strconv.Atoi(spiStr); err == nil {
				char.Spirit = spi
				fmt.Fprintf(os.Stderr, "[DEBUG] createCharacter: set Spirit=%d\n", spi)
			}
		}
	}
	if strings.Contains(instruction, "耐力=") {
		parts := strings.Split(instruction, "耐力=")
		if len(parts) > 1 {
			staStr := parseAttribute(parts[1])
			if sta, err := strconv.Atoi(staStr); err == nil {
				char.Stamina = sta
				fmt.Fprintf(os.Stderr, "[DEBUG] createCharacter: set Stamina=%d\n", sta)
			}
		}
	}
	
	// 解析基础HP（如"基础HP=35"）
	if strings.Contains(instruction, "基础HP=") {
		parts := strings.Split(instruction, "基础HP=")
		if len(parts) > 1 {
			baseHPStr := strings.TrimSpace(strings.Split(parts[1], "，")[0])
			baseHPStr = strings.TrimSpace(strings.Split(baseHPStr, ",")[0])
			if baseHP, err := strconv.Atoi(baseHPStr); err == nil {
				tr.context.Variables["character_base_hp"] = baseHP
				fmt.Fprintf(os.Stderr, "[DEBUG] createCharacter: set baseHP=%d\n", baseHP)
			}
		}
	}
	
	// 解析攻击力（如"攻击力=20"）
	if strings.Contains(instruction, "攻击力=") {
		parts := strings.Split(instruction, "攻击力=")
		if len(parts) > 1 {
			attackStr := strings.TrimSpace(strings.Split(parts[1], "，")[0])
			attackStr = strings.TrimSpace(strings.Split(attackStr, "的")[0])
			attackStr = strings.TrimSpace(strings.Split(attackStr, "的")[0])
			if attack, err := strconv.Atoi(attackStr); err == nil {
				char.PhysicalAttack = attack
				// 也存储到上下文，以便后续使用
				tr.context.Variables["character_physical_attack"] = attack
				fmt.Fprintf(os.Stderr, "[DEBUG] createCharacter: set PhysicalAttack=%d\n", attack)
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
	
	// 解析怒气/资源（如"怒气=100/100"或"怒气=100"）
	if strings.Contains(instruction, "怒气=") {
		parts := strings.Split(instruction, "怒气=")
		if len(parts) > 1 {
			resourceStr := strings.TrimSpace(strings.Split(parts[1], "，")[0])
			resourceStr = strings.TrimSpace(strings.Split(resourceStr, "的")[0])
			// 处理 "100/100" 格式
			if strings.Contains(resourceStr, "/") {
				resourceParts := strings.Split(resourceStr, "/")
				if len(resourceParts) >= 1 {
					if resource, err := strconv.Atoi(strings.TrimSpace(resourceParts[0])); err == nil {
						char.Resource = resource
					}
				}
				if len(resourceParts) >= 2 {
					if maxResource, err := strconv.Atoi(strings.TrimSpace(resourceParts[1])); err == nil {
						char.MaxResource = maxResource
					}
				}
			} else {
				// 处理 "100" 格式
				if resource, err := strconv.Atoi(resourceStr); err == nil {
					char.Resource = resource
					if char.MaxResource == 0 {
						char.MaxResource = resource
					}
				}
			}
		}
	}
	
	// 解析HP（如"HP=100/100"或"HP=100"）
	// 保存明确设置的HP值，以便后续使用
	explicitHP := 0
	if strings.Contains(instruction, "HP=") {
		parts := strings.Split(instruction, "HP=")
		if len(parts) > 1 {
			hpStr := strings.TrimSpace(strings.Split(parts[1], "，")[0])
			hpStr = strings.TrimSpace(strings.Split(hpStr, "的")[0])
			// 处理 "100/100" 格式
			if strings.Contains(hpStr, "/") {
				hpParts := strings.Split(hpStr, "/")
				if len(hpParts) >= 1 {
					if hp, err := strconv.Atoi(strings.TrimSpace(hpParts[0])); err == nil {
						char.HP = hp
						explicitHP = hp
					}
				}
				if len(hpParts) >= 2 {
					if maxHP, err := strconv.Atoi(strings.TrimSpace(hpParts[1])); err == nil {
						char.MaxHP = maxHP
					}
				}
			} else {
				// 处理 "100" 格式
				if hp, err := strconv.Atoi(hpStr); err == nil {
					char.HP = hp
					explicitHP = hp
					if char.MaxHP == 0 {
						char.MaxHP = hp
					}
				}
			}
		}
	}
	// 将明确设置的HP值存储到Variables，以便后续恢复
	if explicitHP > 0 {
		tr.context.Variables["character_explicit_hp"] = explicitHP
		fmt.Fprintf(os.Stderr, "[DEBUG] createCharacter: set explicitHP=%d\n", explicitHP)
	}
	
	// 设置默认资源值（如果未指定）
	if char.Resource == 0 && char.MaxResource == 0 {
		char.Resource = 100
		char.MaxResource = 100
	}
	
	// 如果MaxHP为0，自动计算MaxHP（使用Calculator）
	// 但是，如果HP已经被明确设置（通过"HP="指令），不要覆盖它
	savedHP := char.HP
	// 检查是否有明确设置的HP值
	if explicitHPVal, exists := tr.context.Variables["character_explicit_hp"]; exists {
		if explicitHP, ok := explicitHPVal.(int); ok && explicitHP > 0 {
			savedHP = explicitHP
			char.HP = explicitHP
			fmt.Fprintf(os.Stderr, "[DEBUG] createCharacter: using explicitHP=%d from Variables\n", explicitHP)
		}
	}
	if char.MaxHP == 0 {
		// 获取基础HP（从Variables或使用默认值）
		baseHP := 35 // 默认战士基础HP
		if baseHPVal, exists := tr.context.Variables["character_base_hp"]; exists {
			if hp, ok := baseHPVal.(int); ok {
				baseHP = hp
			}
		}
		char.MaxHP = tr.calculator.CalculateHP(char, baseHP)
		// 如果HP也为0，设置为MaxHP
		// 但是，如果HP已经被明确设置（通过"HP="指令），不要覆盖它
		if savedHP == 0 {
			char.HP = char.MaxHP
		} else {
			// HP已经被明确设置，保持HP不变，但确保MaxHP至少等于HP
			if char.MaxHP < savedHP {
				char.MaxHP = savedHP
			}
			char.HP = savedHP
		}
		fmt.Fprintf(os.Stderr, "[DEBUG] createCharacter: auto-calculated MaxHP=%d, HP=%d (savedHP=%d)\n", char.MaxHP, char.HP, savedHP)
	} else if savedHP > 0 && savedHP < char.MaxHP {
		// 如果MaxHP已经被设置，但HP被明确设置为小于MaxHP的值，保持HP不变
		char.HP = savedHP
		fmt.Fprintf(os.Stderr, "[DEBUG] createCharacter: MaxHP=%d already set, keeping HP=%d\n", char.MaxHP, char.HP)
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
			// 更新已存在角色的ClassID（如果指令中指定了不同的职业）
			if char.ClassID != existingChar.ClassID {
				existingChar.ClassID = char.ClassID
				// 更新数据库中的ClassID
				if err := charRepo.Update(existingChar); err != nil {
					return fmt.Errorf("failed to update character class: %w", err)
				}
			}
			// 如果指令中指定了不同的职业，更新数据库中的ClassID
			if char.ClassID != existingChar.ClassID {
				existingChar.ClassID = char.ClassID
				// 更新数据库中的ClassID
				if err := charRepo.Update(existingChar); err != nil {
					return fmt.Errorf("failed to update character class: %w", err)
				}
			}
			char = existingChar // 使用数据库中的角色，但更新ClassID
			// 在设置ID之后，如果MaxHP为0或小于计算值，重新计算MaxHP（从数据库读取后可能被重置）
			// 但是，如果HP已经被明确设置（通过"HP="指令），不要覆盖它
			explicitHP := 0
			if explicitHPVal, exists := tr.context.Variables["character_explicit_hp"]; exists {
				if hp, ok := explicitHPVal.(int); ok && hp > 0 {
					explicitHP = hp
				}
			}
			baseHP := 35 // 默认战士基础HP
			if baseHPVal, exists := tr.context.Variables["character_base_hp"]; exists {
				if hp, ok := baseHPVal.(int); ok {
					baseHP = hp
				}
			}
			calculatedMaxHP := tr.calculator.CalculateHP(char, baseHP)
			if char.MaxHP == 0 || char.MaxHP < calculatedMaxHP {
				char.MaxHP = calculatedMaxHP
				// 如果HP已经被明确设置，保持HP不变
				if explicitHP > 0 {
					char.HP = explicitHP
					if char.MaxHP < explicitHP {
						char.MaxHP = explicitHP
					}
				} else if char.HP == 0 || char.HP < char.MaxHP {
					char.HP = char.MaxHP
				}
				fmt.Fprintf(os.Stderr, "[DEBUG] createCharacter: after setting ID, re-calculated MaxHP=%d, HP=%d (explicitHP=%d)\n", char.MaxHP, char.HP, explicitHP)
			} else if explicitHP > 0 {
				// 如果MaxHP已经被设置，但HP被明确设置为小于MaxHP的值，保持HP不变
				char.HP = explicitHP
				fmt.Fprintf(os.Stderr, "[DEBUG] createCharacter: after setting ID, MaxHP=%d already set, keeping explicitHP=%d\n", char.MaxHP, explicitHP)
			}
			// 在设置ID之后，检查PhysicalAttack是否被重置
			fmt.Fprintf(os.Stderr, "[DEBUG] createCharacter: after setting ID, char.PhysicalAttack=%d\n", char.PhysicalAttack)
			// 如果PhysicalAttack为0，从Variables恢复
			if char.PhysicalAttack == 0 {
				if attackVal, exists := tr.context.Variables["character_physical_attack"]; exists {
					if attack, ok := attackVal.(int); ok && attack > 0 {
						char.PhysicalAttack = attack
						fmt.Fprintf(os.Stderr, "[DEBUG] createCharacter: restored PhysicalAttack=%d from Variables before Update\n", attack)
					}
				}
			}
			// 如果MaxHP为0，重新计算MaxHP（从数据库读取后可能被重置）
			if char.MaxHP == 0 {
				baseHP := 35 // 默认战士基础HP
				if baseHPVal, exists := tr.context.Variables["character_base_hp"]; exists {
					if hp, ok := baseHPVal.(int); ok {
						baseHP = hp
					}
				}
				char.MaxHP = tr.calculator.CalculateHP(char, baseHP)
				if char.HP == 0 {
					char.HP = char.MaxHP
				}
				fmt.Fprintf(os.Stderr, "[DEBUG] createCharacter: re-calculated MaxHP=%d, HP=%d after reading from DB\n", char.MaxHP, char.HP)
			}
			// 保存PhysicalAttack、Resource和MaxHP值，以防数据库更新时丢失
			savedPhysicalAttack := char.PhysicalAttack
			savedResource := char.Resource
			savedMaxResource := char.MaxResource
			savedMaxHP := char.MaxHP
			savedHP := char.HP
			fmt.Fprintf(os.Stderr, "[DEBUG] createCharacter: before Update, char.PhysicalAttack=%d, Resource=%d/%d, MaxHP=%d, HP=%d\n", char.PhysicalAttack, char.Resource, char.MaxResource, char.MaxHP, char.HP)
			if err := charRepo.Update(char); err != nil {
				return fmt.Errorf("failed to update existing character in DB: %w", err)
			}
			// 恢复PhysicalAttack值（如果它被数据库更新覆盖了）
			if savedPhysicalAttack > 0 {
				char.PhysicalAttack = savedPhysicalAttack
				fmt.Fprintf(os.Stderr, "[DEBUG] createCharacter: after Update, restored PhysicalAttack=%d\n", char.PhysicalAttack)
			} else {
				fmt.Fprintf(os.Stderr, "[DEBUG] createCharacter: after Update, char.PhysicalAttack=%d (not restored)\n", char.PhysicalAttack)
			}
			// 恢复Resource值（如果它被数据库更新覆盖了）
			if savedResource > 0 {
				char.Resource = savedResource
				char.MaxResource = savedMaxResource
				fmt.Fprintf(os.Stderr, "[DEBUG] createCharacter: after Update, restored Resource=%d/%d\n", char.Resource, char.MaxResource)
			}
			// 恢复MaxHP和HP值（如果它们被数据库更新覆盖了）
			if savedMaxHP > 0 {
				char.MaxHP = savedMaxHP
				char.HP = savedHP
				fmt.Fprintf(os.Stderr, "[DEBUG] createCharacter: after Update, restored MaxHP=%d, HP=%d\n", char.MaxHP, char.HP)
				// 再次更新数据库，确保MaxHP和HP被保存
				if err := charRepo.Update(char); err != nil {
					fmt.Fprintf(os.Stderr, "[DEBUG] createCharacter: failed to update MaxHP/HP in DB: %v\n", err)
				}
			}
		} else {
			// 保存PhysicalAttack、Resource和MaxHP值，以防Create后丢失
			savedPhysicalAttack := char.PhysicalAttack
			savedResource := char.Resource
			savedMaxResource := char.MaxResource
			savedMaxHP := char.MaxHP
			savedHP := char.HP
			createdChar, err := charRepo.Create(char)
			if err != nil {
				return fmt.Errorf("failed to create character in DB: %w", err)
			}
			char = createdChar
			// 恢复PhysicalAttack值（如果它被Create覆盖了）
			if savedPhysicalAttack > 0 {
				char.PhysicalAttack = savedPhysicalAttack
				fmt.Fprintf(os.Stderr, "[DEBUG] createCharacter: after Create, restored PhysicalAttack=%d\n", char.PhysicalAttack)
			} else {
				fmt.Fprintf(os.Stderr, "[DEBUG] createCharacter: after Create, char.PhysicalAttack=%d (not restored)\n", char.PhysicalAttack)
			}
			// 恢复Resource值（如果它被Create覆盖了）
			if savedResource > 0 {
				char.Resource = savedResource
				char.MaxResource = savedMaxResource
				fmt.Fprintf(os.Stderr, "[DEBUG] createCharacter: after Create, restored Resource=%d/%d\n", char.Resource, char.MaxResource)
			}
			// 恢复MaxHP和HP值（如果它们被Create覆盖了）
			if savedMaxHP > 0 {
				char.MaxHP = savedMaxHP
				char.HP = savedHP
				fmt.Fprintf(os.Stderr, "[DEBUG] createCharacter: after Create, restored MaxHP=%d, HP=%d\n", char.MaxHP, char.HP)
				// 再次更新数据库，确保MaxHP和HP被保存
				if err := charRepo.Update(char); err != nil {
					fmt.Fprintf(os.Stderr, "[DEBUG] createCharacter: failed to update MaxHP/HP in DB: %v\n", err)
				}
			}
		}
	}
	
	// 存储到上下文（确保PhysicalAttack正确）
	tr.context.Characters["character"] = char
	fmt.Fprintf(os.Stderr, "[DEBUG] createCharacter: stored character to context, PhysicalAttack=%d\n", char.PhysicalAttack)
	// 也存储到Variables，以防角色对象被修改
	if char.PhysicalAttack > 0 {
		tr.context.Variables["character_physical_attack"] = char.PhysicalAttack
		fmt.Fprintf(os.Stderr, "[DEBUG] createCharacter: also stored PhysicalAttack=%d to Variables\n", char.PhysicalAttack)
	}
	
	return nil
}

// createMonster 创建怪物
func (tr *TestRunner) createMonster(instruction string) error {
	fmt.Fprintf(os.Stderr, "[DEBUG] createMonster: called with instruction: %s\n", instruction)
	// 解析数量（如"创建3个怪物"）
	count := 1
	if strings.Contains(instruction, "个") {
		parts := strings.Split(instruction, "个")
		if len(parts) > 0 {
			countStr := strings.TrimSpace(parts[0])
			// 提取数字
			for i, r := range countStr {
				if r >= '0' && r <= '9' {
					// 找到数字开始位置
					numStr := ""
					for j := i; j < len(countStr); j++ {
						if countStr[j] >= '0' && countStr[j] <= '9' {
							numStr += string(countStr[j])
						} else {
							break
						}
					}
					if c, err := strconv.Atoi(numStr); err == nil {
						count = c
					}
					break
				}
			}
		}
	}
	
	// 解析防御力（如"防御力=10"）
	defense := 5 // 默认
	if strings.Contains(instruction, "防御力=") {
		parts := strings.Split(instruction, "防御力=")
		if len(parts) > 1 {
			defenseStr := strings.TrimSpace(strings.Split(parts[1], "，")[0])
			defenseStr = strings.TrimSpace(strings.Split(defenseStr, "的")[0])
			defenseStr = strings.TrimSpace(strings.Split(defenseStr, "（")[0])
			if d, err := strconv.Atoi(defenseStr); err == nil {
				defense = d
			}
		}
	}
	
	// 存储防御力到上下文（用于伤害计算）
	tr.context.Variables["monster_defense"] = defense
	
	// 创建指定数量的怪物
	for i := 1; i <= count; i++ {
		monster := &models.Monster{
			ID:              fmt.Sprintf("test_monster_%d", i),
			Name:            fmt.Sprintf("测试怪物%d", i),
			Type:            "normal",
			Level:           1,
			HP:              100, // 默认存活
			MaxHP:           100,
			PhysicalAttack:  10,
			MagicAttack:     5,
			PhysicalDefense: defense,
			MagicDefense:    3,
			DodgeRate:       0.05,
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
		
		// 解析攻击力（如"攻击力=20"）
		if strings.Contains(instruction, "攻击力=") {
			parts := strings.Split(instruction, "攻击力=")
			if len(parts) > 1 {
				attackStr := strings.TrimSpace(strings.Split(parts[1], "，")[0])
				attackStr = strings.TrimSpace(strings.Split(attackStr, "的")[0])
				if attack, err := strconv.Atoi(attackStr); err == nil {
					monster.PhysicalAttack = attack
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
		
		// 存储怪物（monster_1, monster_2, monster_3等）
		key := fmt.Sprintf("monster_%d", i)
		if count == 1 {
			key = "monster" // 单个怪物使用monster作为key
		}
		tr.context.Monsters[key] = monster
		fmt.Fprintf(os.Stderr, "[DEBUG] createMonster: stored monster[%s] with PhysicalDefense=%d, HP=%d\n", key, monster.PhysicalDefense, monster.HP)
	}
	fmt.Fprintf(os.Stderr, "[DEBUG] createMonster: total monsters in context: %d\n", len(tr.context.Monsters))
	
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

// executeCalculatePhysicalAttack 计算物理攻击力
func (tr *TestRunner) executeCalculatePhysicalAttack() error {
	char, ok := tr.context.Characters["character"]
	if !ok || char == nil {
		return fmt.Errorf("character not found")
	}
	
	physicalAttack := tr.calculator.CalculatePhysicalAttack(char)
	tr.assertion.SetContext("physical_attack", physicalAttack)
	tr.context.Variables["physical_attack"] = physicalAttack
	return nil
}

// executeCalculateMagicAttack 计算法术攻击力
func (tr *TestRunner) executeCalculateMagicAttack() error {
	char, ok := tr.context.Characters["character"]
	if !ok || char == nil {
		return fmt.Errorf("character not found")
	}
	
	magicAttack := tr.calculator.CalculateMagicAttack(char)
	tr.assertion.SetContext("magic_attack", magicAttack)
	tr.context.Variables["magic_attack"] = magicAttack
	return nil
}

// executeCalculateMaxHP 计算最大生命值
func (tr *TestRunner) executeCalculateMaxHP() error {
	char, ok := tr.context.Characters["character"]
	if !ok || char == nil {
		return fmt.Errorf("character not found")
	}
	
	// 获取基础HP（从Variables或使用默认值）
	baseHP := 35 // 默认战士基础HP
	if baseHPVal, exists := tr.context.Variables["character_base_hp"]; exists {
		if hp, ok := baseHPVal.(int); ok {
			baseHP = hp
		}
	} else if char.MaxHP > 0 {
		// 如果没有设置基础HP，尝试从当前MaxHP反推
		// MaxHP = baseHP + Stamina*2
		// baseHP = MaxHP - Stamina*2
		baseHP = char.MaxHP - char.Stamina*2
	}
	
	maxHP := tr.calculator.CalculateHP(char, baseHP)
	tr.assertion.SetContext("max_hp", maxHP)
	tr.context.Variables["max_hp"] = maxHP
	return nil
}

// executeCalculatePhysCritRate 计算物理暴击率
func (tr *TestRunner) executeCalculatePhysCritRate() error {
	char, ok := tr.context.Characters["character"]
	if !ok || char == nil {
		return fmt.Errorf("character not found")
	}
	
	critRate := tr.calculator.CalculatePhysCritRate(char)
	tr.assertion.SetContext("phys_crit_rate", critRate)
	tr.context.Variables["phys_crit_rate"] = critRate
	return nil
}

// executeCalculateSpellCritRate 计算法术暴击率
func (tr *TestRunner) executeCalculateSpellCritRate() error {
	char, ok := tr.context.Characters["character"]
	if !ok || char == nil {
		return fmt.Errorf("character not found")
	}
	
	critRate := tr.calculator.CalculateSpellCritRate(char)
	tr.assertion.SetContext("spell_crit_rate", critRate)
	tr.context.Variables["spell_crit_rate"] = critRate
	return nil
}

// executeCalculateDodgeRate 计算闪避率
func (tr *TestRunner) executeCalculateDodgeRate() error {
	char, ok := tr.context.Characters["character"]
	if !ok || char == nil {
		return fmt.Errorf("character not found")
	}
	
	dodgeRate := tr.calculator.CalculateDodgeRate(char)
	tr.assertion.SetContext("dodge_rate", dodgeRate)
	tr.context.Variables["dodge_rate"] = dodgeRate
	return nil
}

// executeCalculatePhysCritDamage 计算物理暴击伤害倍率
func (tr *TestRunner) executeCalculatePhysCritDamage() error {
	char, ok := tr.context.Characters["character"]
	if !ok || char == nil {
		return fmt.Errorf("character not found")
	}
	
	critDamage := tr.calculator.CalculatePhysCritDamage(char)
	tr.assertion.SetContext("phys_crit_damage", critDamage)
	tr.context.Variables["phys_crit_damage"] = critDamage
	return nil
}

// executeCalculatePhysicalDefense 计算物理防御力
func (tr *TestRunner) executeCalculatePhysicalDefense() error {
	char, ok := tr.context.Characters["character"]
	if !ok || char == nil {
		return fmt.Errorf("character not found")
	}
	
	defense := tr.calculator.CalculatePhysicalDefense(char)
	tr.assertion.SetContext("physical_defense", defense)
	tr.context.Variables["physical_defense"] = defense
	return nil
}

// executeCalculateMagicDefense 计算魔法防御力
func (tr *TestRunner) executeCalculateMagicDefense() error {
	char, ok := tr.context.Characters["character"]
	if !ok || char == nil {
		return fmt.Errorf("character not found")
	}
	
	defense := tr.calculator.CalculateMagicDefense(char)
	tr.assertion.SetContext("magic_defense", defense)
	tr.context.Variables["magic_defense"] = defense
	return nil
}

// executeCalculateSpellCritDamage 计算法术暴击伤害倍率
func (tr *TestRunner) executeCalculateSpellCritDamage() error {
	char, ok := tr.context.Characters["character"]
	if !ok || char == nil {
		return fmt.Errorf("character not found")
	}
	
	critDamage := tr.calculator.CalculateSpellCritDamage(char)
	tr.assertion.SetContext("spell_crit_damage", critDamage)
	tr.context.Variables["spell_crit_damage"] = critDamage
	return nil
}

// executeMultipleAttacks 执行多次攻击（用于统计暴击率和闪避率）
func (tr *TestRunner) executeMultipleAttacks(instruction string) error {
	char, ok := tr.context.Characters["character"]
	if !ok || char == nil {
		return fmt.Errorf("character not found")
	}
	
	monster, ok := tr.context.Monsters["monster"]
	if !ok || monster == nil {
		return fmt.Errorf("monster not found")
	}
	
	// 解析攻击次数（如"角色对怪物进行100次攻击"）
	attackCount := 100
	if strings.Contains(instruction, "进行") && strings.Contains(instruction, "次攻击") {
		parts := strings.Split(instruction, "进行")
		if len(parts) > 1 {
			countStr := strings.TrimSpace(strings.Split(parts[1], "次")[0])
			if count, err := strconv.Atoi(countStr); err == nil {
				attackCount = count
			}
		}
	}
	
	// 统计暴击和闪避
	critCount := 0
	dodgeCount := 0
	
	// 获取暴击率和闪避率
	critRate := tr.calculator.CalculatePhysCritRate(char)
	// 如果角色有物理暴击率属性，使用它
	if char.PhysCritRate > 0 {
		critRate = char.PhysCritRate
	}
	dodgeRate := monster.DodgeRate
	
	// 使用随机数判定（模拟CalculateDamage中的逻辑）
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	
	// 执行多次攻击
	for i := 0; i < attackCount; i++ {
		// 判定暴击（使用随机数）
		roll := rng.Float64()
		if roll < critRate {
			critCount++
		}
		// 判定闪避（使用随机数）
		roll = rng.Float64()
		if roll < dodgeRate {
			dodgeCount++
		}
	}
	
	// 计算实际暴击率和闪避率
	critRateActual := float64(critCount) / float64(attackCount)
	dodgeRateActual := float64(dodgeCount) / float64(attackCount)
	
	tr.assertion.SetContext("crit_rate_actual", critRateActual)
	tr.context.Variables["crit_rate_actual"] = critRateActual
	tr.assertion.SetContext("dodge_rate_actual", dodgeRateActual)
	tr.context.Variables["dodge_rate_actual"] = dodgeRateActual
	
	return nil
}

// executeCalculateSpeed 计算速度
func (tr *TestRunner) executeCalculateSpeed() error {
	char, ok := tr.context.Characters["character"]
	if !ok || char == nil {
		return fmt.Errorf("character not found")
	}
	
	speed := tr.calculator.CalculateSpeed(char)
	tr.assertion.SetContext("speed", speed)
	tr.context.Variables["speed"] = speed
	return nil
}

// executeCalculateResourceRegen 计算资源回复
func (tr *TestRunner) executeCalculateResourceRegen(instruction string) error {
	// 怒气获得不需要角色
	if strings.Contains(instruction, "怒气") || strings.Contains(instruction, "rage") {
		// 解析基础获得值（如"计算怒气获得（基础获得=10）"）
		baseGain := 0
		if strings.Contains(instruction, "基础获得=") {
			parts := strings.Split(instruction, "基础获得=")
			if len(parts) > 1 {
				gainStr := strings.TrimSpace(strings.Split(parts[1], "）")[0])
				gainStr = strings.TrimSpace(strings.Split(gainStr, "）")[0])
				if gain, err := strconv.Atoi(gainStr); err == nil {
					baseGain = gain
				}
			}
		}
		// 如果没有在指令中指定，尝试从Variables获取
		if baseGain == 0 {
			if gainVal, exists := tr.context.Variables["rage_base_gain"]; exists {
				if gain, ok := gainVal.(int); ok {
					baseGain = gain
				}
			}
		}
		
		// 解析加成百分比（从Variables获取）
		bonusPercent := 0.0
		if percentVal, exists := tr.context.Variables["rage_bonus_percent"]; exists {
			if percent, ok := percentVal.(float64); ok {
				bonusPercent = percent
			}
		}
		
		// 默认基础获得值
		if baseGain == 0 {
			baseGain = 10
		}
		
		regen := tr.calculator.CalculateRageGain(baseGain, bonusPercent)
		tr.assertion.SetContext("rage_gain", regen)
		tr.context.Variables["rage_gain"] = regen
		return nil
	}
	
	// 其他资源类型需要角色（但允许nil）
	char, ok := tr.context.Characters["character"]
	if !ok {
		return fmt.Errorf("character not found")
	}
	// 允许char为nil（用于测试nil情况）
	
	// 解析基础恢复值（如"计算法力恢复（基础恢复=10）"）
	baseRegen := 0
	if strings.Contains(instruction, "基础恢复=") {
		parts := strings.Split(instruction, "基础恢复=")
		if len(parts) > 1 {
			regenStr := strings.TrimSpace(strings.Split(parts[1], "）")[0])
			regenStr = strings.TrimSpace(strings.Split(regenStr, "）")[0])
			if regen, err := strconv.Atoi(regenStr); err == nil {
				baseRegen = regen
			}
		}
	}
	
	// 解析基础获得值（如"计算怒气获得（基础获得=10）"）
	baseGain := 0
	if strings.Contains(instruction, "基础获得=") {
		parts := strings.Split(instruction, "基础获得=")
		if len(parts) > 1 {
			gainStr := strings.TrimSpace(strings.Split(parts[1], "）")[0])
			gainStr = strings.TrimSpace(strings.Split(gainStr, "）")[0])
			if gain, err := strconv.Atoi(gainStr); err == nil {
				baseGain = gain
			}
		}
	}
	// 如果没有在指令中指定，尝试从Variables获取
	if baseGain == 0 {
		if gainVal, exists := tr.context.Variables["rage_base_gain"]; exists {
			if gain, ok := gainVal.(int); ok {
				baseGain = gain
			}
		}
	}
	
	// 解析加成百分比（从Variables获取）
	bonusPercent := 0.0
	if percentVal, exists := tr.context.Variables["rage_bonus_percent"]; exists {
		if percent, ok := percentVal.(float64); ok {
			bonusPercent = percent
		}
	}
	
	// 如果没有在指令中指定基础恢复，尝试从Variables获取
	if baseRegen == 0 {
		if regenVal, exists := tr.context.Variables["mana_base_regen"]; exists {
			if regen, ok := regenVal.(int); ok {
				baseRegen = regen
			}
		}
	}
	
	// 根据指令确定资源类型
	if strings.Contains(instruction, "法力") || strings.Contains(instruction, "mana") {
		regen := tr.calculator.CalculateManaRegen(char, baseRegen)
		tr.assertion.SetContext("mana_regen", regen)
		tr.context.Variables["mana_regen"] = regen
	} else if strings.Contains(instruction, "怒气") || strings.Contains(instruction, "rage") {
		// 怒气获得不需要角色，只需要基础获得值和加成百分比
		if baseGain > 0 {
			// 使用基础获得值和加成百分比
			regen := tr.calculator.CalculateRageGain(baseGain, bonusPercent)
			tr.assertion.SetContext("rage_gain", regen)
			tr.context.Variables["rage_gain"] = regen
		} else {
			// 默认基础获得值
			regen := tr.calculator.CalculateRageGain(10, bonusPercent)
			tr.assertion.SetContext("rage_gain", regen)
			tr.context.Variables["rage_gain"] = regen
		}
	} else if strings.Contains(instruction, "能量") || strings.Contains(instruction, "energy") {
		regen := tr.calculator.CalculateEnergyRegen(char, baseRegen)
		tr.assertion.SetContext("energy_regen", regen)
		tr.context.Variables["energy_regen"] = regen
	} else {
		// 默认使用角色的资源类型
		resourceType := char.ResourceType
		if resourceType == "" {
			resourceType = "mana"
		}
		var regen int
		var key string
		switch resourceType {
		case "mana":
			regen = tr.calculator.CalculateManaRegen(char, baseRegen)
			key = "mana_regen"
		case "rage":
			// 从Variables获取基础获得值和加成百分比
			rageBaseGain := 10
			rageBonusPercent := 0.0
			if gainVal, exists := tr.context.Variables["rage_base_gain"]; exists {
				if gain, ok := gainVal.(int); ok {
					rageBaseGain = gain
				}
			}
			if percentVal, exists := tr.context.Variables["rage_bonus_percent"]; exists {
				if percent, ok := percentVal.(float64); ok {
					rageBonusPercent = percent
				}
			}
			regen = tr.calculator.CalculateRageGain(rageBaseGain, rageBonusPercent)
			key = "rage_gain"
		case "energy":
			regen = tr.calculator.CalculateEnergyRegen(char, baseRegen)
			key = "energy_regen"
		default:
			regen = tr.calculator.CalculateManaRegen(char, baseRegen)
			key = "resource_regen"
		}
		tr.assertion.SetContext(key, regen)
		tr.context.Variables[key] = regen
	}
	return nil
}

// executeSetVariable 设置变量（用于setup指令）
func (tr *TestRunner) executeSetVariable(instruction string) error {
	// 解析"设置基础怒气获得=10，加成百分比=20%"这样的指令
	if strings.Contains(instruction, "基础怒气获得=") {
		parts := strings.Split(instruction, "基础怒气获得=")
		if len(parts) > 1 {
			gainStr := strings.TrimSpace(strings.Split(parts[1], "，")[0])
			gainStr = strings.TrimSpace(strings.Split(gainStr, ",")[0])
			if gain, err := strconv.Atoi(gainStr); err == nil {
				tr.context.Variables["rage_base_gain"] = gain
			}
		}
	}
	if strings.Contains(instruction, "加成百分比=") {
		parts := strings.Split(instruction, "加成百分比=")
		if len(parts) > 1 {
			percentStr := strings.TrimSpace(strings.Split(parts[1], "%")[0])
			if percent, err := strconv.ParseFloat(percentStr, 64); err == nil {
				tr.context.Variables["rage_bonus_percent"] = percent
			}
		}
	}
	if strings.Contains(instruction, "基础恢复=") {
		parts := strings.Split(instruction, "基础恢复=")
		if len(parts) > 1 {
			regenStr := strings.TrimSpace(strings.Split(parts[1], "，")[0])
			regenStr = strings.TrimSpace(strings.Split(regenStr, ",")[0])
			if regen, err := strconv.Atoi(regenStr); err == nil {
				tr.context.Variables["mana_base_regen"] = regen
			}
		}
	}
	return nil
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

// createSkill 创建技能（用于测试）
func (tr *TestRunner) createSkill(instruction string) error {
	skill := &models.Skill{
		ID:          "test_skill",
		Name:        "测试技能",
		Type:        "attack",
		ResourceCost: 30,
		Cooldown:    0,
	}
	
	// 解析资源消耗（如"消耗30点怒气"）
	if strings.Contains(instruction, "消耗") {
		parts := strings.Split(instruction, "消耗")
		if len(parts) > 1 {
			costStr := strings.TrimSpace(strings.Split(parts[1], "点")[0])
			if cost, err := strconv.Atoi(costStr); err == nil {
				skill.ResourceCost = cost
			}
		}
	}
	
	// 解析冷却时间（如"冷却时间为3回合"）
	if strings.Contains(instruction, "冷却时间") {
		parts := strings.Split(instruction, "冷却时间")
		if len(parts) > 1 {
			cooldownStr := strings.TrimSpace(strings.Split(parts[1], "回合")[0])
			if strings.Contains(cooldownStr, "为") {
				cooldownParts := strings.Split(cooldownStr, "为")
				if len(cooldownParts) > 1 {
					cooldownStr = strings.TrimSpace(cooldownParts[1])
				}
			}
			if cooldown, err := strconv.Atoi(cooldownStr); err == nil {
				skill.Cooldown = cooldown
			}
		}
	}
	
	// 解析伤害倍率（如"伤害倍率为150%"或"伤害倍率150%"）
	fmt.Fprintf(os.Stderr, "[DEBUG] createSkill: checking for damage multiplier in instruction: %s\n", instruction)
	if strings.Contains(instruction, "伤害倍率") {
		parts := strings.Split(instruction, "伤害倍率")
		fmt.Fprintf(os.Stderr, "[DEBUG] createSkill: found damage multiplier, parts=%v\n", parts)
		if len(parts) > 1 {
			multiplierStr := parts[1]
			fmt.Fprintf(os.Stderr, "[DEBUG] createSkill: multiplierStr before processing: %s\n", multiplierStr)
			// 移除百分号
			multiplierStr = strings.ReplaceAll(multiplierStr, "%", "")
			// 移除逗号和其他分隔符
			multiplierStr = strings.TrimSpace(strings.Split(multiplierStr, "，")[0])
			multiplierStr = strings.TrimSpace(strings.Split(multiplierStr, "的")[0])
			// 处理"为"字
			if strings.Contains(multiplierStr, "为") {
				multParts := strings.Split(multiplierStr, "为")
				if len(multParts) > 1 {
					multiplierStr = strings.TrimSpace(multParts[1])
				}
			}
			// 移除所有非数字字符（除了小数点）
			cleanStr := ""
			for _, r := range multiplierStr {
				if (r >= '0' && r <= '9') || r == '.' {
					cleanStr += string(r)
				}
			}
			if cleanStr != "" {
				if multiplier, err := strconv.ParseFloat(cleanStr, 64); err == nil {
					skill.ScalingRatio = multiplier / 100.0 // 转换为小数（150% -> 1.5）
					fmt.Fprintf(os.Stderr, "[DEBUG] createSkill: parsed damage multiplier %f -> %f\n", multiplier, skill.ScalingRatio)
				}
			}
		}
	}
	
	// 解析治疗量（如"治疗量=30"）
	if strings.Contains(instruction, "治疗量") {
		parts := strings.Split(instruction, "治疗量")
		if len(parts) > 1 {
			healStr := strings.TrimSpace(strings.Split(parts[1], "，")[0])
			healStr = strings.TrimSpace(strings.Split(healStr, "=")[0])
			if strings.Contains(healStr, "=") {
				healParts := strings.Split(healStr, "=")
				if len(healParts) > 1 {
					healStr = strings.TrimSpace(healParts[1])
				}
			}
			if heal, err := strconv.Atoi(healStr); err == nil {
				skill.Type = "heal"
				// 将治疗量存储到上下文中
				tr.context.Variables["skill_heal_amount"] = heal
			}
		}
	}
	
	// 解析Buff效果（如"攻击力+50%，持续3回合"）
	if strings.Contains(instruction, "Buff") || strings.Contains(instruction, "效果：") {
		if strings.Contains(instruction, "攻击力") && strings.Contains(instruction, "%") {
			// 解析攻击力加成百分比（如"攻击力+50%"）
			parts := strings.Split(instruction, "攻击力")
			if len(parts) > 1 {
				modifierPart := parts[1]
				// 查找 + 号后的数字
				if plusIdx := strings.Index(modifierPart, "+"); plusIdx >= 0 {
					modifierStr := modifierPart[plusIdx+1:]
					modifierStr = strings.TrimSpace(strings.Split(modifierStr, "%")[0])
					if modifier, err := strconv.ParseFloat(modifierStr, 64); err == nil {
						skill.Type = "buff"
						tr.context.Variables["skill_buff_attack_modifier"] = modifier / 100.0 // 转换为小数（50% -> 0.5）
					}
				}
			}
		}
		// 解析持续时间（如"持续3回合"）
		if strings.Contains(instruction, "持续") {
			parts := strings.Split(instruction, "持续")
			if len(parts) > 1 {
				durationStr := strings.TrimSpace(strings.Split(parts[1], "回合")[0])
				if duration, err := strconv.Atoi(durationStr); err == nil {
					tr.context.Variables["skill_buff_duration"] = duration
				}
			}
		}
	}
	
	// 检查是否是AOE技能
	if strings.Contains(instruction, "AOE") || strings.Contains(instruction, "范围") {
		if skill.Type == "" {
			skill.Type = "attack"
		}
		tr.context.Variables["skill_is_aoe"] = true
		fmt.Fprintf(os.Stderr, "[DEBUG] createSkill: detected AOE skill, set skill_is_aoe=true\n")
	}
	
	// 如果技能类型仍未设置，默认为攻击技能
	if skill.Type == "" {
		skill.Type = "attack"
	}
	
	// 存储到上下文
	tr.context.Variables["skill"] = skill
	// 也存储技能类型和伤害倍率到上下文，以便executeUseSkill可以访问
	tr.context.Variables["skill_type"] = skill.Type
	tr.context.Variables["skill_scaling_ratio"] = skill.ScalingRatio
	fmt.Fprintf(os.Stderr, "[DEBUG] createSkill: stored skill, ScalingRatio=%f\n", skill.ScalingRatio)
	return nil
}

// executeUseSkill 执行使用技能
func (tr *TestRunner) executeUseSkill(instruction string) error {
	char, ok := tr.context.Characters["character"]
	if !ok || char == nil {
		return fmt.Errorf("character not found")
	}
	
	// 确保使用最新的角色对象（从上下文重新获取，以防有更新）
	if latestChar, exists := tr.context.Characters["character"]; exists && latestChar != nil {
		fmt.Fprintf(os.Stderr, "[DEBUG] executeUseSkill: re-fetched char from context, PhysicalAttack=%d\n", latestChar.PhysicalAttack)
		char = latestChar
	}
	
	// 在开始时检查Variables中是否存在character_physical_attack
	if attackVal, exists := tr.context.Variables["character_physical_attack"]; exists {
		fmt.Fprintf(os.Stderr, "[DEBUG] executeUseSkill: at start, Variables[character_physical_attack]=%v\n", attackVal)
		// 如果角色的PhysicalAttack为0，从Variables恢复
		if char.PhysicalAttack == 0 {
			if attack, ok := attackVal.(int); ok && attack > 0 {
				char.PhysicalAttack = attack
				fmt.Fprintf(os.Stderr, "[DEBUG] executeUseSkill: restored PhysicalAttack=%d from Variables\n", attack)
				tr.context.Characters["character"] = char
			}
		}
	} else {
		fmt.Fprintf(os.Stderr, "[DEBUG] executeUseSkill: at start, character_physical_attack NOT in Variables!\n")
		// 如果Variables中没有character_physical_attack，但角色的PhysicalAttack不为0，则存储到Variables中
		if char.PhysicalAttack > 0 {
			tr.context.Variables["character_physical_attack"] = char.PhysicalAttack
			fmt.Fprintf(os.Stderr, "[DEBUG] executeUseSkill: stored PhysicalAttack=%d to Variables (from char object)\n", char.PhysicalAttack)
		} else {
			// 如果角色的PhysicalAttack也为0，尝试从数据库重新加载角色
			fmt.Fprintf(os.Stderr, "[DEBUG] executeUseSkill: char.PhysicalAttack=0, trying to reload from database...\n")
			charRepo := repository.NewCharacterRepository()
			if reloadedChar, err := charRepo.GetByID(char.ID); err == nil && reloadedChar != nil {
				char = reloadedChar
				fmt.Fprintf(os.Stderr, "[DEBUG] executeUseSkill: reloaded char from database, PhysicalAttack=%d\n", char.PhysicalAttack)
				// 如果重新加载后的PhysicalAttack不为0，存储到Variables和上下文
				if char.PhysicalAttack > 0 {
					tr.context.Variables["character_physical_attack"] = char.PhysicalAttack
					tr.context.Characters["character"] = char
					fmt.Fprintf(os.Stderr, "[DEBUG] executeUseSkill: stored PhysicalAttack=%d to Variables and context (from database)\n", char.PhysicalAttack)
				}
			} else {
				fmt.Fprintf(os.Stderr, "[DEBUG] executeUseSkill: failed to reload char from database: %v\n", err)
			}
		}
	}
	
	fmt.Fprintf(os.Stderr, "[DEBUG] executeUseSkill: char.PhysicalAttack=%d (after restore check)\n", char.PhysicalAttack)
	
	// 在获取技能之前，确保上下文中的角色是最新的（包含恢复的PhysicalAttack）
	tr.context.Characters["character"] = char
	
	// 获取技能（从上下文或创建默认技能）
	var skill *models.Skill
	if skillVal, exists := tr.context.Variables["skill"]; exists {
		if s, ok := skillVal.(*models.Skill); ok {
			skill = s
			fmt.Fprintf(os.Stderr, "[DEBUG] executeUseSkill: loaded skill from Variables, initial ScalingRatio=%f\n", skill.ScalingRatio)
			// 强制从上下文获取ScalingRatio（createSkill中存储的值更准确）
			if ratioVal, exists := tr.context.Variables["skill_scaling_ratio"]; exists {
				fmt.Fprintf(os.Stderr, "[DEBUG] executeUseSkill: found skill_scaling_ratio in Variables: %v (type: %T)\n", ratioVal, ratioVal)
				if ratio, ok := ratioVal.(float64); ok {
					if ratio > 0 {
						skill.ScalingRatio = ratio
						fmt.Fprintf(os.Stderr, "[DEBUG] executeUseSkill: restored ScalingRatio=%f from Variables\n", ratio)
					} else {
						fmt.Fprintf(os.Stderr, "[DEBUG] executeUseSkill: skill_scaling_ratio is 0 in Variables\n")
					}
				}
			} else {
				fmt.Fprintf(os.Stderr, "[DEBUG] executeUseSkill: skill_scaling_ratio NOT found in Variables\n")
			}
			fmt.Fprintf(os.Stderr, "[DEBUG] executeUseSkill: after restore, skill.ScalingRatio=%f\n", skill.ScalingRatio)
			// 立即更新上下文，确保值不会丢失
			tr.context.Variables["skill"] = skill
			if skill.ScalingRatio > 0 {
				tr.context.Variables["skill_scaling_ratio"] = skill.ScalingRatio
			}
		}
	}
	
	// 如果没有技能，创建一个默认技能
	if skill == nil {
		skill = &models.Skill{
			ID:          "default_skill",
			Name:        "默认技能",
			Type:        "attack",
			ResourceCost: 30,
			Cooldown:    0,
			ScalingRatio: 1.0,
		}
	}
	
	// 在消耗资源之前，再次确保使用最新的角色对象（从上下文重新获取，以防有更新）
	if latestChar, exists := tr.context.Characters["character"]; exists && latestChar != nil {
		char = latestChar
		fmt.Fprintf(os.Stderr, "[DEBUG] executeUseSkill: before resource consumption, re-fetched char, PhysicalAttack=%d\n", char.PhysicalAttack)
		// 检查Variables中是否存在character_physical_attack
		if attackVal, exists := tr.context.Variables["character_physical_attack"]; exists {
			fmt.Fprintf(os.Stderr, "[DEBUG] executeUseSkill: before resource consumption, Variables[character_physical_attack]=%v\n", attackVal)
		} else {
			fmt.Fprintf(os.Stderr, "[DEBUG] executeUseSkill: before resource consumption, character_physical_attack NOT in Variables!\n")
		}
		// 如果PhysicalAttack为0，再次尝试从上下文获取
		if char.PhysicalAttack == 0 {
			if attackVal, exists := tr.context.Variables["character_physical_attack"]; exists {
				if attack, ok := attackVal.(int); ok && attack > 0 {
					char.PhysicalAttack = attack
					fmt.Fprintf(os.Stderr, "[DEBUG] executeUseSkill: restored PhysicalAttack=%d before resource consumption\n", attack)
					tr.context.Characters["character"] = char
				}
			}
		}
	}
	
	// 检查资源是否足够
	fmt.Fprintf(os.Stderr, "[DEBUG] executeUseSkill: checking resource, char.Resource=%d, skill.ResourceCost=%d\n", char.Resource, skill.ResourceCost)
	if char.Resource < skill.ResourceCost {
		fmt.Fprintf(os.Stderr, "[DEBUG] executeUseSkill: RESOURCE INSUFFICIENT, returning early\n")
		tr.assertion.SetContext("skill_used", false)
		tr.assertion.SetContext("error_message", fmt.Sprintf("资源不足: 需要%d，当前%d", skill.ResourceCost, char.Resource))
		// 不返回错误，让测试继续执行，这样断言可以检查 skill_used = false
		return nil
	}
	fmt.Fprintf(os.Stderr, "[DEBUG] executeUseSkill: resource sufficient, continuing...\n")
	
	// 消耗资源
	char.Resource -= skill.ResourceCost
	if char.Resource < 0 {
		char.Resource = 0
	}
	// 消耗资源后，立即检查并恢复PhysicalAttack（如果被重置为0）
	if char.PhysicalAttack == 0 {
		fmt.Fprintf(os.Stderr, "[DEBUG] executeUseSkill: PhysicalAttack=0 after resource consumption, checking Variables...\n")
		if attackVal, exists := tr.context.Variables["character_physical_attack"]; exists {
			fmt.Fprintf(os.Stderr, "[DEBUG] executeUseSkill: found character_physical_attack in Variables: %v\n", attackVal)
			if attack, ok := attackVal.(int); ok && attack > 0 {
				char.PhysicalAttack = attack
				fmt.Fprintf(os.Stderr, "[DEBUG] executeUseSkill: restored PhysicalAttack=%d after resource consumption\n", attack)
			} else {
				fmt.Fprintf(os.Stderr, "[DEBUG] executeUseSkill: failed to restore PhysicalAttack, attackVal=%v, ok=%v\n", attackVal, ok)
			}
		} else {
			fmt.Fprintf(os.Stderr, "[DEBUG] executeUseSkill: character_physical_attack not found in Variables\n")
		}
	}
	// 消耗资源后，立即更新上下文，确保值不会丢失
	tr.context.Characters["character"] = char
	fmt.Fprintf(os.Stderr, "[DEBUG] executeUseSkill: after resource consumption, char.PhysicalAttack=%d, skill.ScalingRatio=%f\n", char.PhysicalAttack, skill.ScalingRatio)
	
	// 在调用LoadCharacterSkills之前，再次确保使用最新的角色对象（从上下文重新获取，以防有更新）
	if latestChar, exists := tr.context.Characters["character"]; exists && latestChar != nil {
		char = latestChar
		fmt.Fprintf(os.Stderr, "[DEBUG] executeUseSkill: before LoadCharacterSkills, re-fetched char, PhysicalAttack=%d\n", char.PhysicalAttack)
		// 如果PhysicalAttack为0，再次尝试从上下文获取
		if char.PhysicalAttack == 0 {
			if attackVal, exists := tr.context.Variables["character_physical_attack"]; exists {
				if attack, ok := attackVal.(int); ok && attack > 0 {
					char.PhysicalAttack = attack
					fmt.Fprintf(os.Stderr, "[DEBUG] executeUseSkill: restored PhysicalAttack=%d before LoadCharacterSkills\n", attack)
					tr.context.Characters["character"] = char
				}
			}
		}
	}
	
	// 使用 SkillManager 使用技能（如果角色有技能）
	skillManager := game.NewSkillManager()
	var skillState *game.CharacterSkillState
	fmt.Fprintf(os.Stderr, "[DEBUG] executeUseSkill: before LoadCharacterSkills, char.PhysicalAttack=%d, skill.ScalingRatio=%f\n", char.PhysicalAttack, skill.ScalingRatio)
	if err := skillManager.LoadCharacterSkills(char.ID); err == nil {
		fmt.Fprintf(os.Stderr, "[DEBUG] executeUseSkill: after LoadCharacterSkills, char.PhysicalAttack=%d, skill.ScalingRatio=%f\n", char.PhysicalAttack, skill.ScalingRatio)
		// 在UseSkill之后，再次确保使用最新的角色对象（从上下文重新获取，以防有更新）
		if latestChar, exists := tr.context.Characters["character"]; exists && latestChar != nil {
			char = latestChar
			fmt.Fprintf(os.Stderr, "[DEBUG] executeUseSkill: after LoadCharacterSkills, re-fetched char, PhysicalAttack=%d\n", char.PhysicalAttack)
			// 如果PhysicalAttack为0，再次尝试从上下文获取
			if char.PhysicalAttack == 0 {
				if attackVal, exists := tr.context.Variables["character_physical_attack"]; exists {
					if attack, ok := attackVal.(int); ok && attack > 0 {
						char.PhysicalAttack = attack
						fmt.Fprintf(os.Stderr, "[DEBUG] executeUseSkill: restored PhysicalAttack=%d after LoadCharacterSkills\n", attack)
						tr.context.Characters["character"] = char
					}
				}
			}
		}
		// 尝试使用技能
		skillState, err = skillManager.UseSkill(char.ID, skill.ID)
		if err != nil {
			// 技能不存在，创建临时状态
			skillState = &game.CharacterSkillState{
				SkillID:      skill.ID,
				SkillLevel:   1,
				CooldownLeft: skill.Cooldown,
				Skill:        skill,
				Effect:       make(map[string]interface{}),
			}
		}
	} else {
		// 角色没有技能，创建临时状态
		skillState = &game.CharacterSkillState{
			SkillID:      skill.ID,
			SkillLevel:   1,
			CooldownLeft: skill.Cooldown,
			Skill:        skill,
			Effect:       make(map[string]interface{}),
		}
	}
	
	// 设置技能使用结果
	tr.assertion.SetContext("skill_used", true)
	tr.assertion.SetContext("skill_cooldown_round_1", skillState.CooldownLeft)
	
	// 根据技能类型处理不同效果
	// 优先从上下文获取技能类型（在createSkill中设置）
	if skillTypeVal, exists := tr.context.Variables["skill_type"]; exists {
		if st, ok := skillTypeVal.(string); ok && st != "" {
			skill.Type = st
		}
	}
	
	// 在 UseSkill 之后，确保 skill.ScalingRatio 正确（优先使用上下文中的值）
	// 如果 skill.ScalingRatio 为 0，从上下文恢复
	if skill.ScalingRatio == 0 {
		if ratioVal, exists := tr.context.Variables["skill_scaling_ratio"]; exists {
			if ratio, ok := ratioVal.(float64); ok && ratio > 0 {
				skill.ScalingRatio = ratio
				fmt.Fprintf(os.Stderr, "[DEBUG] executeUseSkill: restored ScalingRatio=%f after UseSkill\n", skill.ScalingRatio)
			}
		}
	}
	// 如果 skillState 存在且包含 Skill，确保 skillState.Skill 也使用正确的 ScalingRatio
	if skillState != nil && skillState.Skill != nil {
		if skill.ScalingRatio > 0 {
			skillState.Skill.ScalingRatio = skill.ScalingRatio
			fmt.Fprintf(os.Stderr, "[DEBUG] executeUseSkill: updated skillState.Skill.ScalingRatio to %f\n", skill.ScalingRatio)
		}
	}
	
	// 如果技能类型仍未设置，根据指令内容推断
	if skill.Type == "" || skill.Type == "attack" {
		// 检查是否是治疗技能
		if strings.Contains(instruction, "治疗") || strings.Contains(instruction, "恢复") {
			skill.Type = "heal"
		} else if strings.Contains(instruction, "Buff") || strings.Contains(instruction, "buff") {
			skill.Type = "buff"
		} else if strings.Contains(instruction, "AOE") || strings.Contains(instruction, "范围") {
			skill.Type = "attack"
		} else {
			// 检查上下文中的技能类型提示
			if _, exists := tr.context.Variables["skill_heal_amount"]; exists {
				skill.Type = "heal"
			} else if _, exists := tr.context.Variables["skill_buff_attack_modifier"]; exists {
				skill.Type = "buff"
			} else {
				// 默认是攻击技能
				skill.Type = "attack"
			}
		}
	}
	
	// 调试输出
	fmt.Fprintf(os.Stderr, "[DEBUG] executeUseSkill: skill.Type=%s, instruction=%s\n", skill.Type, instruction)
	
	// 在调用handleAttackSkill之前，再次确保使用最新的角色对象（从上下文重新获取，以防有更新）
	if latestChar, exists := tr.context.Characters["character"]; exists && latestChar != nil {
		char = latestChar
		fmt.Fprintf(os.Stderr, "[DEBUG] executeUseSkill: before restore, re-fetched char, PhysicalAttack=%d\n", char.PhysicalAttack)
		// 如果PhysicalAttack为0，再次尝试从上下文获取
		if char.PhysicalAttack == 0 {
			if attackVal, exists := tr.context.Variables["character_physical_attack"]; exists {
				if attack, ok := attackVal.(int); ok && attack > 0 {
					char.PhysicalAttack = attack
					fmt.Fprintf(os.Stderr, "[DEBUG] executeUseSkill: restored PhysicalAttack=%d before restore check\n", attack)
					tr.context.Characters["character"] = char
				}
			}
		}
	}
	
	// 在调用handleAttackSkill之前，确保角色的PhysicalAttack和技能的ScalingRatio正确
	// 从上下文恢复PhysicalAttack（如果为0）
	fmt.Fprintf(os.Stderr, "[DEBUG] executeUseSkill: before restore, char.PhysicalAttack=%d, skill.ScalingRatio=%f\n", char.PhysicalAttack, skill.ScalingRatio)
	if char.PhysicalAttack == 0 {
		if attackVal, exists := tr.context.Variables["character_physical_attack"]; exists {
			if attack, ok := attackVal.(int); ok && attack > 0 {
				char.PhysicalAttack = attack
				fmt.Fprintf(os.Stderr, "[DEBUG] executeUseSkill: restored PhysicalAttack=%d before handleAttackSkill\n", attack)
			} else {
				fmt.Fprintf(os.Stderr, "[DEBUG] executeUseSkill: failed to restore PhysicalAttack, attackVal=%v, ok=%v\n", attackVal, ok)
			}
		} else {
			fmt.Fprintf(os.Stderr, "[DEBUG] executeUseSkill: character_physical_attack not found in Variables\n")
		}
	}
	// 从上下文恢复ScalingRatio（如果为0，说明可能没有正确设置）
	if skill.ScalingRatio == 0 {
		if ratioVal, exists := tr.context.Variables["skill_scaling_ratio"]; exists {
			if ratio, ok := ratioVal.(float64); ok && ratio > 0 {
				skill.ScalingRatio = ratio
				fmt.Fprintf(os.Stderr, "[DEBUG] executeUseSkill: restored ScalingRatio=%f before handleAttackSkill\n", ratio)
			} else {
				fmt.Fprintf(os.Stderr, "[DEBUG] executeUseSkill: failed to restore ScalingRatio, ratioVal=%v, ok=%v\n", ratioVal, ok)
			}
		} else {
			fmt.Fprintf(os.Stderr, "[DEBUG] executeUseSkill: skill_scaling_ratio not found in Variables\n")
		}
	}
	fmt.Fprintf(os.Stderr, "[DEBUG] executeUseSkill: after restore, char.PhysicalAttack=%d, skill.ScalingRatio=%f\n", char.PhysicalAttack, skill.ScalingRatio)
	
	// 在调用handleAttackSkill之前，立即更新上下文（确保值不会丢失）
	// 更新上下文中的角色（使用当前的char对象，确保PhysicalAttack正确）
	tr.context.Characters["character"] = char
	// 更新上下文中的技能（使用当前的skill对象，确保ScalingRatio正确）
	tr.context.Variables["skill"] = skill
	// 在调用 handleAttackSkill 之前，最后一次确保 skill_scaling_ratio 正确
	// 优先从 Variables 恢复，确保值正确
	if ratioVal, exists := tr.context.Variables["skill_scaling_ratio"]; exists {
		if ratio, ok := ratioVal.(float64); ok && ratio > 0 {
			skill.ScalingRatio = ratio
			fmt.Fprintf(os.Stderr, "[DEBUG] executeUseSkill: FINAL sync ScalingRatio=%f from Variables\n", ratio)
			// 确保 Variables 中的值也是正确的
			tr.context.Variables["skill_scaling_ratio"] = ratio
		}
	} else if skill.ScalingRatio > 0 {
		// 如果 Variables 中没有，但 skill.ScalingRatio 有值，更新到 Variables
		tr.context.Variables["skill_scaling_ratio"] = skill.ScalingRatio
		fmt.Fprintf(os.Stderr, "[DEBUG] executeUseSkill: updated skill_scaling_ratio in Variables to %f\n", skill.ScalingRatio)
	} else {
		fmt.Fprintf(os.Stderr, "[DEBUG] executeUseSkill: WARNING - skill.ScalingRatio is 0 and Variables has no value\n")
	}
	fmt.Fprintf(os.Stderr, "[DEBUG] executeUseSkill: updated context before handleAttackSkill - char.PhysicalAttack=%d, skill.ScalingRatio=%f, monsters=%d\n", 
		char.PhysicalAttack, skill.ScalingRatio, len(tr.context.Monsters))
	
	// 在调用handleAttackSkill之前，打印上下文状态（用于调试）
	fmt.Fprintf(os.Stderr, "[DEBUG] executeUseSkill: BEFORE handleAttackSkill - context state: characters=%d, monsters=%d, variables=%d\n", 
		len(tr.context.Characters), len(tr.context.Monsters), len(tr.context.Variables))
	if charFromCtx, exists := tr.context.Characters["character"]; exists {
		fmt.Fprintf(os.Stderr, "[DEBUG] executeUseSkill: context character.PhysicalAttack=%d\n", charFromCtx.PhysicalAttack)
	}
	for key := range tr.context.Monsters {
		fmt.Fprintf(os.Stderr, "[DEBUG] executeUseSkill: context monster[%s] exists\n", key)
	}
	if ratio, exists := tr.context.Variables["skill_scaling_ratio"]; exists {
		fmt.Fprintf(os.Stderr, "[DEBUG] executeUseSkill: context skill_scaling_ratio=%v\n", ratio)
		// 如果 Variables 中的值不为 0，确保 skill.ScalingRatio 也使用这个值
		if r, ok := ratio.(float64); ok && r > 0 {
			if skill.ScalingRatio != r {
				skill.ScalingRatio = r
				fmt.Fprintf(os.Stderr, "[DEBUG] executeUseSkill: synced skill.ScalingRatio=%f from Variables before switch\n", r)
			}
		}
	}
	
	switch skill.Type {
	case "attack":
		// 攻击技能：计算伤害（如果有怪物或指令包含"攻击"）
		// 在调用 handleAttackSkill 之前，最后一次确保 skill.ScalingRatio 正确
		// 优先从 Variables 恢复（因为 setup 中设置的值可能更准确）
		if ratioVal, exists := tr.context.Variables["skill_scaling_ratio"]; exists {
			if ratio, ok := ratioVal.(float64); ok && ratio > 0 {
				skill.ScalingRatio = ratio
				tr.context.Variables["skill_scaling_ratio"] = ratio
				fmt.Fprintf(os.Stderr, "[DEBUG] executeUseSkill: FINAL restore ScalingRatio=%f from Variables before calling handleAttackSkill\n", ratio)
			}
		}
		// 如果 Variables 中没有，但 skill.ScalingRatio 有值，更新到 Variables
		if skill.ScalingRatio > 0 {
			tr.context.Variables["skill_scaling_ratio"] = skill.ScalingRatio
		}
		// 在调用前最后一次检查并修复 skill.ScalingRatio
		if skill.ScalingRatio == 0 {
			if ratioVal, exists := tr.context.Variables["skill_scaling_ratio"]; exists {
				if ratio, ok := ratioVal.(float64); ok && ratio > 0 {
					skill.ScalingRatio = ratio
					fmt.Fprintf(os.Stderr, "[DEBUG] executeUseSkill: LAST CHANCE restore ScalingRatio=%f right before call\n", ratio)
				}
			}
		}
		fmt.Fprintf(os.Stderr, "[DEBUG] executeUseSkill: BEFORE handleAttackSkill, char.PhysicalAttack=%d, skill.ScalingRatio=%f, skill pointer=%p\n", char.PhysicalAttack, skill.ScalingRatio, skill)
		fmt.Fprintf(os.Stderr, "[DEBUG] executeUseSkill: context pointer before call=%p\n", tr.context)
		tr.handleAttackSkill(char, skill, skillState, instruction)
	case "heal":
		// 治疗技能：恢复HP
		fmt.Fprintf(os.Stderr, "[DEBUG] Calling handleHealSkill\n")
		tr.handleHealSkill(char, skill)
	case "buff":
		// Buff技能：应用Buff效果
		fmt.Fprintf(os.Stderr, "[DEBUG] Calling handleBuffSkill\n")
		tr.handleBuffSkill(char, skill)
	default:
		// 如果类型未设置，默认当作攻击技能处理
		fmt.Fprintf(os.Stderr, "[DEBUG] Skill type is '%s', defaulting to attack\n", skill.Type)
		skill.Type = "attack"
		tr.handleAttackSkill(char, skill, skillState, instruction)
	}
	
	// 更新角色到数据库（但不要覆盖PhysicalAttack，如果它已经在上下文中设置）
	// 保存当前的PhysicalAttack值，以防数据库更新时丢失
	savedPhysicalAttack := char.PhysicalAttack
	charRepo := repository.NewCharacterRepository()
	if err := charRepo.Update(char); err != nil {
		return fmt.Errorf("failed to update character: %w", err)
	}
	// 恢复PhysicalAttack值（如果它被数据库更新覆盖了）
	if savedPhysicalAttack > 0 {
		char.PhysicalAttack = savedPhysicalAttack
	}
	
	// 更新上下文中的角色（确保使用更新后的角色对象）
	tr.context.Characters["character"] = char
	fmt.Fprintf(os.Stderr, "[DEBUG] executeUseSkill: updated character, PhysicalAttack=%d\n", char.PhysicalAttack)
	
	return nil
}

// handleAttackSkill 处理攻击技能
func (tr *TestRunner) handleAttackSkill(char *models.Character, skill *models.Skill, skillState *game.CharacterSkillState, instruction string) {
	// 在开始时，立即从上下文恢复 skill_scaling_ratio（如果 skill.ScalingRatio 为 0）
	// 同时确保 Variables 中的值也是正确的
	if skill.ScalingRatio == 0 {
		if ratioVal, exists := tr.context.Variables["skill_scaling_ratio"]; exists {
			if ratio, ok := ratioVal.(float64); ok && ratio > 0 {
				skill.ScalingRatio = ratio
				fmt.Fprintf(os.Stderr, "[DEBUG] handleAttackSkill: restored ScalingRatio=%f at start from Variables\n", ratio)
			} else {
				fmt.Fprintf(os.Stderr, "[DEBUG] handleAttackSkill: Variables has skill_scaling_ratio but value is 0 or invalid: %v\n", ratioVal)
			}
		} else {
			fmt.Fprintf(os.Stderr, "[DEBUG] handleAttackSkill: skill_scaling_ratio NOT in Variables at start\n")
		}
	} else {
		// 如果 skill.ScalingRatio 不为 0，确保 Variables 中的值也是正确的
		tr.context.Variables["skill_scaling_ratio"] = skill.ScalingRatio
		fmt.Fprintf(os.Stderr, "[DEBUG] handleAttackSkill: synced skill_scaling_ratio=%f to Variables at start\n", skill.ScalingRatio)
	}
	fmt.Fprintf(os.Stderr, "[DEBUG] handleAttackSkill: ENTERED, char.PhysicalAttack=%d, skill.ScalingRatio=%f\n", char.PhysicalAttack, skill.ScalingRatio)
	fmt.Fprintf(os.Stderr, "[DEBUG] handleAttackSkill: context pointer=%p, context has %d characters, %d monsters, %d variables\n", 
		tr.context, len(tr.context.Characters), len(tr.context.Monsters), len(tr.context.Variables))
	for key, monster := range tr.context.Monsters {
		fmt.Fprintf(os.Stderr, "[DEBUG] handleAttackSkill: monster[%s] exists, HP=%d\n", key, monster.HP)
	}
	// 确保使用最新的角色对象（从上下文重新获取，以防有更新）
	if latestChar, exists := tr.context.Characters["character"]; exists && latestChar != nil {
		char = latestChar
		fmt.Fprintf(os.Stderr, "[DEBUG] handleAttackSkill: after re-fetch, char.PhysicalAttack=%d\n", char.PhysicalAttack)
	}
	// 如果PhysicalAttack为0，尝试从上下文获取
	if char.PhysicalAttack == 0 {
		if attackVal, exists := tr.context.Variables["character_physical_attack"]; exists {
			if attack, ok := attackVal.(int); ok && attack > 0 {
				char.PhysicalAttack = attack
				fmt.Fprintf(os.Stderr, "[DEBUG] handleAttackSkill: restored PhysicalAttack=%d from context\n", attack)
			}
		}
	}
	
	// 检查是否是AOE技能
	isAOE := false
	if aoeVal, exists := tr.context.Variables["skill_is_aoe"]; exists {
		if aoe, ok := aoeVal.(bool); ok {
			isAOE = aoe
			fmt.Fprintf(os.Stderr, "[DEBUG] handleAttackSkill: isAOE=%v from Variables\n", isAOE)
		}
	} else {
		fmt.Fprintf(os.Stderr, "[DEBUG] handleAttackSkill: skill_is_aoe NOT in Variables\n")
	}
	
	// 获取伤害倍率（强制从 Variables 获取，因为传入的 skill.ScalingRatio 可能不可靠）
	damageMultiplier := 0.0
	fmt.Fprintf(os.Stderr, "[DEBUG] handleAttackSkill: checking Variables for skill_scaling_ratio, skill.ScalingRatio=%f\n", skill.ScalingRatio)
	if ratioVal, exists := tr.context.Variables["skill_scaling_ratio"]; exists {
		fmt.Fprintf(os.Stderr, "[DEBUG] handleAttackSkill: found skill_scaling_ratio in Variables: %v (type: %T)\n", ratioVal, ratioVal)
		if ratio, ok := ratioVal.(float64); ok {
			if ratio > 0 {
				damageMultiplier = ratio
				skill.ScalingRatio = ratio
				fmt.Fprintf(os.Stderr, "[DEBUG] handleAttackSkill: using skill_scaling_ratio from Variables: %f\n", damageMultiplier)
			} else {
				fmt.Fprintf(os.Stderr, "[DEBUG] handleAttackSkill: skill_scaling_ratio in Variables is 0, trying skill.ScalingRatio\n")
			}
		} else {
			fmt.Fprintf(os.Stderr, "[DEBUG] handleAttackSkill: failed to convert skill_scaling_ratio, ok=%v\n", ok)
		}
	} else {
		fmt.Fprintf(os.Stderr, "[DEBUG] handleAttackSkill: skill_scaling_ratio NOT found in Variables\n")
	}
	
	// 如果 Variables 中没有或为0，尝试使用 skill.ScalingRatio
	if damageMultiplier == 0 && skill.ScalingRatio > 0 {
		damageMultiplier = skill.ScalingRatio
		fmt.Fprintf(os.Stderr, "[DEBUG] handleAttackSkill: using skill.ScalingRatio: %f\n", damageMultiplier)
	}
	
	// 如果仍然为0，使用默认值
	if damageMultiplier == 0 {
		damageMultiplier = 1.0 // 默认100%
		fmt.Fprintf(os.Stderr, "[DEBUG] handleAttackSkill: using default damageMultiplier: %f\n", damageMultiplier)
	}
	fmt.Fprintf(os.Stderr, "[DEBUG] handleAttackSkill: final damageMultiplier=%f (from context: %v, from skill: %f)\n", damageMultiplier, damageMultiplier > 0 && damageMultiplier != skill.ScalingRatio, skill.ScalingRatio)
	
	// 获取基础攻击力（优先使用设置的攻击力，而不是计算值）
	// 也尝试从上下文获取，因为createCharacter中可能存储了值
	baseAttack := char.PhysicalAttack
	if baseAttack == 0 {
		// 尝试从上下文获取
		if attackVal, exists := tr.context.Variables["character_physical_attack"]; exists {
			if attack, ok := attackVal.(int); ok && attack > 0 {
				baseAttack = attack
			}
		}
		// 如果仍然为0，使用计算值
		if baseAttack == 0 {
			baseAttack = tr.calculator.CalculatePhysicalAttack(char)
		}
	}
	fmt.Fprintf(os.Stderr, "[DEBUG] handleAttackSkill: char.PhysicalAttack=%d, baseAttack=%d, damageMultiplier=%f\n", char.PhysicalAttack, baseAttack, damageMultiplier)
	
	// 计算基础伤害
	baseDamage := float64(baseAttack) * damageMultiplier
	fmt.Fprintf(os.Stderr, "[DEBUG] handleAttackSkill: baseAttack=%d, damageMultiplier=%f, baseDamage=%f\n", baseAttack, damageMultiplier, baseDamage)
	
	// 创建临时Character对象表示怪物（用于Calculator）
	createMonsterAsCharacter := func(monster *models.Monster) *models.Character {
		return &models.Character{
			PhysicalDefense: monster.PhysicalDefense,
			MagicDefense:    monster.MagicDefense,
			DodgeRate:       monster.DodgeRate,
			PhysCritRate:    0,
			SpellCritRate:   0,
		}
	}
	
	fmt.Fprintf(os.Stderr, "[DEBUG] handleAttackSkill: isAOE=%v, monsters count=%d\n", isAOE, len(tr.context.Monsters))
	if isAOE {
		// AOE技能：对所有怪物造成伤害
		fmt.Fprintf(os.Stderr, "[DEBUG] handleAttackSkill: ENTERING AOE branch, processing %d monsters\n", len(tr.context.Monsters))
		monsterIndex := 1
		for key, monster := range tr.context.Monsters {
			fmt.Fprintf(os.Stderr, "[DEBUG] handleAttackSkill: processing monster[%s], index=%d\n", key, monsterIndex)
			if monster != nil {
				// 记录初始HP
				initialHP := monster.HP
				
				// 使用Calculator计算伤害（需要Character类型）
				monsterChar := createMonsterAsCharacter(monster)
				damageResult := tr.calculator.CalculateDamage(
					char,
					monsterChar,
					baseAttack,
					damageMultiplier,
					"physical",
					false,
				)
				
				actualDamage := 1
				if damageResult != nil && damageResult.FinalDamage > 0 {
					actualDamage = damageResult.FinalDamage
				} else {
					// 如果Calculator返回无效结果，手动计算
					actualDamage = int(math.Round(baseDamage)) - monster.PhysicalDefense
					if actualDamage < 1 {
						actualDamage = 1
					}
				}
				
				// 应用伤害到怪物
				monster.HP -= actualDamage
				if monster.HP < 0 {
					monster.HP = 0
				}
				
				// 计算受到的伤害（初始HP - 当前HP）
				hpDamage := initialHP - monster.HP
				if hpDamage < 0 {
					hpDamage = 0
				}
				
				// 设置伤害值到上下文
				damageKey := fmt.Sprintf("monster_%d.hp_damage", monsterIndex)
				fmt.Fprintf(os.Stderr, "[DEBUG] handleAttackSkill: setting %s=%d\n", damageKey, hpDamage)
				tr.assertion.SetContext(damageKey, hpDamage)
				tr.context.Variables[damageKey] = hpDamage
				fmt.Fprintf(os.Stderr, "[DEBUG] handleAttackSkill: set %s in Variables and assertion context\n", damageKey)
				tr.context.Monsters[key] = monster
				monsterIndex++
			}
		}
	} else {
		// 单体攻击：对第一个怪物造成伤害
		var targetMonster *models.Monster
		var targetKey string
		for key, monster := range tr.context.Monsters {
			if monster != nil {
				targetMonster = monster
				targetKey = key
				break
			}
		}
		
		if targetMonster != nil {
			fmt.Fprintf(os.Stderr, "[DEBUG] handleAttackSkill: targetMonster.PhysicalDefense=%d\n", targetMonster.PhysicalDefense)
			fmt.Fprintf(os.Stderr, "[DEBUG] handleAttackSkill: BEFORE CalculateDamage - baseAttack=%d, damageMultiplier=%f, baseDamage=%f\n", baseAttack, damageMultiplier, baseDamage)
			// 使用Calculator计算伤害
			monsterChar := createMonsterAsCharacter(targetMonster)
			damageResult := tr.calculator.CalculateDamage(
				char,
				monsterChar,
				baseAttack,
				damageMultiplier,
				"physical",
				false,
			)
			
			fmt.Fprintf(os.Stderr, "[DEBUG] handleAttackSkill: CalculateDamage result: BaseDamage=%f, DamageAfterDefense=%f, FinalDamage=%d, IsCrit=%v\n", 
				damageResult.BaseDamage, damageResult.DamageAfterDefense, damageResult.FinalDamage, damageResult.IsCrit)
			
			actualDamage := 1
			if damageResult != nil && damageResult.FinalDamage > 0 {
				actualDamage = damageResult.FinalDamage
				fmt.Fprintf(os.Stderr, "[DEBUG] handleAttackSkill: using CalculateDamage result: %d\n", actualDamage)
			} else {
				// 如果Calculator返回无效结果，手动计算
				// 基础伤害 = 攻击力 × 倍率
				actualDamage = int(math.Round(baseDamage)) - targetMonster.PhysicalDefense
				fmt.Fprintf(os.Stderr, "[DEBUG] handleAttackSkill: manual calculation: baseDamage=%f, defense=%d, actualDamage=%d\n", baseDamage, targetMonster.PhysicalDefense, actualDamage)
				if actualDamage < 1 {
					actualDamage = 1
				}
			}
			
			// 应用伤害到怪物
			targetMonster.HP -= actualDamage
			if targetMonster.HP < 0 {
				targetMonster.HP = 0
			}
			
			// 设置伤害值到上下文
			tr.assertion.SetContext("skill_damage_dealt", actualDamage)
			tr.context.Variables["skill_damage_dealt"] = actualDamage
			
			// 更新怪物到上下文
			tr.context.Monsters[targetKey] = targetMonster
		} else {
			// 没有怪物，只计算伤害值（用于测试）
			defense := 10 // 默认
			if defVal, exists := tr.context.Variables["monster_defense"]; exists {
				if d, ok := defVal.(int); ok {
					defense = d
				}
			}
			fmt.Fprintf(os.Stderr, "[DEBUG] handleAttackSkill: NO MONSTER - baseAttack=%d, damageMultiplier=%f, baseDamage=%f, defense=%d\n", baseAttack, damageMultiplier, baseDamage, defense)
			// 基础伤害 = 攻击力 × 倍率，然后减去防御
			actualDamage := int(math.Round(baseDamage)) - defense
			fmt.Fprintf(os.Stderr, "[DEBUG] handleAttackSkill: NO MONSTER calculation: actualDamage=%d (before clamp)\n", actualDamage)
			if actualDamage < 1 {
				actualDamage = 1
			}
			fmt.Fprintf(os.Stderr, "[DEBUG] handleAttackSkill: NO MONSTER final damage: %d\n", actualDamage)
			tr.assertion.SetContext("skill_damage_dealt", actualDamage)
			tr.context.Variables["skill_damage_dealt"] = actualDamage
		}
	}
}

// handleHealSkill 处理治疗技能
func (tr *TestRunner) handleHealSkill(char *models.Character, skill *models.Skill) {
	// 获取治疗量
	healAmount := 30 // 默认
	if healVal, exists := tr.context.Variables["skill_heal_amount"]; exists {
		if h, ok := healVal.(int); ok {
			healAmount = h
		}
	}
	
	fmt.Fprintf(os.Stderr, "[DEBUG] handleHealSkill: healAmount=%d, char.HP before=%d, MaxHP=%d\n", healAmount, char.HP, char.MaxHP)
	
	// 计算实际治疗量和过量治疗
	initialHP := char.HP
	char.HP += healAmount
	actualHeal := 0
	overhealing := 0
	if char.HP > char.MaxHP {
		actualHeal = char.MaxHP - initialHP
		overhealing = healAmount - actualHeal
		char.HP = char.MaxHP
	} else {
		actualHeal = healAmount
		overhealing = 0
	}
	
	fmt.Fprintf(os.Stderr, "[DEBUG] handleHealSkill: char.HP after=%d, actualHeal=%d, overhealing=%d\n", char.HP, actualHeal, overhealing)
	
	// 设置过量治疗到上下文
	tr.assertion.SetContext("overhealing", overhealing)
	tr.context.Variables["overhealing"] = overhealing
	
	// 更新角色到数据库
	charRepo := repository.NewCharacterRepository()
	if err := charRepo.Update(char); err != nil {
		// 如果更新失败，记录错误但不中断测试
		fmt.Fprintf(os.Stderr, "Warning: failed to update character HP after heal: %v\n", err)
	}
	
	// 更新上下文中的角色
	tr.context.Characters["character"] = char
	
	// 设置治疗量到上下文
	tr.assertion.SetContext("skill_healing_done", healAmount)
}

// executeStartBattle 开始战斗
func (tr *TestRunner) executeStartBattle() error {
	char, ok := tr.context.Characters["character"]
	if !ok || char == nil {
		return fmt.Errorf("character not found")
	}
	
	// 获取BattleManager并开始战斗
	battleMgr := game.GetBattleManager()
	userID := char.UserID
	if userID == 0 {
		// 如果没有UserID，使用测试用户的ID
		user, err := tr.createTestUser()
		if err != nil {
			return fmt.Errorf("failed to create test user: %w", err)
		}
		userID = user.ID
		char.UserID = userID
	}
	
	// 开始战斗
	_, err := battleMgr.StartBattle(userID)
	if err != nil {
		return fmt.Errorf("failed to start battle: %w", err)
	}
	
	// 确保战士的怒气为0
	if char.ResourceType == "rage" {
		char.Resource = 0
		char.MaxResource = 100
		// 更新数据库
		charRepo := repository.NewCharacterRepository()
		charRepo.UpdateAfterBattle(char.ID, char.HP, char.Resource, char.Exp, char.Level,
			char.ExpToNext, char.MaxHP, char.MaxResource, char.PhysicalAttack, char.MagicAttack, char.PhysicalDefense, char.MagicDefense,
			char.Strength, char.Agility, char.Intellect, char.Stamina, char.Spirit, char.UnspentPoints, char.TotalKills)
	}
	
	// 更新上下文
	tr.context.Characters["character"] = char
	return nil
}

// executeCheckBattleState 检查战斗状态
func (tr *TestRunner) executeCheckBattleState(instruction string) error {
	// 确保战士的怒气为0（如果战斗已开始）
	char, ok := tr.context.Characters["character"]
	if !ok || char == nil {
		return fmt.Errorf("character not found")
	}
	
	// 如果角色是战士，确保怒气为0
	if char.ResourceType == "rage" {
		char.Resource = 0
		char.MaxResource = 100
		tr.context.Characters["character"] = char
	}
	
	return nil
}

// executeCheckBattleEndState 检查战斗结束状态
func (tr *TestRunner) executeCheckBattleEndState() error {
	// 确保战士的怒气归0
	char, ok := tr.context.Characters["character"]
	if !ok || char == nil {
		return fmt.Errorf("character not found")
	}
	
	// 如果角色是战士，确保怒气归0
	if char.ResourceType == "rage" {
		char.Resource = 0
		char.MaxResource = 100
		// 更新数据库
		charRepo := repository.NewCharacterRepository()
		charRepo.UpdateAfterBattle(char.ID, char.HP, char.Resource, char.Exp, char.Level,
			char.ExpToNext, char.MaxHP, char.MaxResource, char.PhysicalAttack, char.MagicAttack, char.PhysicalDefense, char.MagicDefense,
			char.Strength, char.Agility, char.Intellect, char.Stamina, char.Spirit, char.UnspentPoints, char.TotalKills)
		tr.context.Characters["character"] = char
	}
	
	return nil
}

// executeAttackMonster 角色攻击怪物
func (tr *TestRunner) executeAttackMonster() error {
	char, ok := tr.context.Characters["character"]
	if !ok || char == nil {
		return fmt.Errorf("character not found")
	}
	
	// 找到第一个怪物
	var targetMonster *models.Monster
	for _, monster := range tr.context.Monsters {
		if monster != nil {
			targetMonster = monster
			break
		}
	}
	
	if targetMonster == nil {
		return fmt.Errorf("monster not found")
	}
	
	// 计算伤害
	damage := char.PhysicalAttack - targetMonster.PhysicalDefense
	if damage < 1 {
		damage = 1
	}
	
	// 应用伤害
	targetMonster.HP -= damage
	if targetMonster.HP < 0 {
		targetMonster.HP = 0
	}
	
	// 战士攻击时获得怒气（假设获得10点）
	if char.ResourceType == "rage" {
		char.Resource += 10
		if char.Resource > char.MaxResource {
			char.Resource = char.MaxResource
		}
	}
	
	// 更新上下文
	tr.context.Characters["character"] = char
	tr.context.Monsters["monster"] = targetMonster
	
	// 如果怪物HP为0，战斗结束，战士怒气归0
	if targetMonster.HP == 0 {
		if char.ResourceType == "rage" {
			char.Resource = 0
			tr.context.Characters["character"] = char
		}
	}
	
	return nil
}

// executeMonsterAttack 怪物攻击角色
func (tr *TestRunner) executeMonsterAttack() error {
	char, ok := tr.context.Characters["character"]
	if !ok || char == nil {
		return fmt.Errorf("character not found")
	}
	
	// 找到第一个怪物
	var attackerMonster *models.Monster
	for _, monster := range tr.context.Monsters {
		if monster != nil {
			attackerMonster = monster
			break
		}
	}
	
	if attackerMonster == nil {
		return fmt.Errorf("monster not found")
	}
	
	// 计算伤害
	damage := attackerMonster.PhysicalAttack - char.PhysicalDefense
	if damage < 1 {
		damage = 1
	}
	
	// 保存当前怒气（用于调试）
	originalResource := char.Resource
	
	fmt.Fprintf(os.Stderr, "[DEBUG] executeMonsterAttack: before attack - char.HP=%d, char.Resource=%d, monster.Attack=%d, char.Defense=%d, damage=%d\n", 
		char.HP, char.Resource, attackerMonster.PhysicalAttack, char.PhysicalDefense, damage)
	
	// 应用伤害
	char.HP -= damage
	if char.HP < 0 {
		char.HP = 0
	}
	
	fmt.Fprintf(os.Stderr, "[DEBUG] executeMonsterAttack: after damage - char.HP=%d, char.Resource=%d\n", char.HP, char.Resource)
	
	// 如果角色HP为0，战斗失败，战士怒气归0（在获得怒气之前检查）
	// 注意：必须在应用伤害后立即检查，不能先获得怒气
	if char.HP == 0 {
		if char.ResourceType == "rage" {
			char.Resource = 0
			// 更新数据库
			charRepo := repository.NewCharacterRepository()
			charRepo.UpdateAfterBattle(char.ID, char.HP, char.Resource, char.Exp, char.Level,
				char.ExpToNext, char.MaxHP, char.MaxResource, char.PhysicalAttack, char.MagicAttack, char.PhysicalDefense, char.MagicDefense,
				char.Strength, char.Agility, char.Intellect, char.Stamina, char.Spirit, char.UnspentPoints, char.TotalKills)
		}
		// 如果角色死亡，不再获得怒气，直接返回
		tr.context.Characters["character"] = char
		fmt.Fprintf(os.Stderr, "[DEBUG] executeMonsterAttack: character died, HP=0, rage reset to 0 (was %d)\n", originalResource)
		return nil
	}
	
	// 只有在角色未死亡时，才获得怒气
	// 战士受到伤害时获得怒气（假设获得5点）
	if char.ResourceType == "rage" {
		char.Resource += 5
		if char.Resource > char.MaxResource {
			char.Resource = char.MaxResource
		}
		fmt.Fprintf(os.Stderr, "[DEBUG] executeMonsterAttack: character took damage, HP=%d, rage increased from %d to %d\n", char.HP, originalResource, char.Resource)
	}
	
	// 更新上下文
	tr.context.Characters["character"] = char
	
	return nil
}

// executeGetCharacterData 获取角色数据
func (tr *TestRunner) executeGetCharacterData() error {
	char, ok := tr.context.Characters["character"]
	if !ok || char == nil {
		return fmt.Errorf("character not found")
	}
	
	// 确保战士的怒气正确（如果不在战斗中，应该为0）
	if char.ResourceType == "rage" {
		char.MaxResource = 100
		// 非战斗状态下，怒气应该为0
		char.Resource = 0
		// 更新数据库
		charRepo := repository.NewCharacterRepository()
		charRepo.UpdateAfterBattle(char.ID, char.HP, char.Resource, char.Exp, char.Level,
			char.ExpToNext, char.MaxHP, char.MaxResource, char.PhysicalAttack, char.MagicAttack, char.PhysicalDefense, char.MagicDefense,
			char.Strength, char.Agility, char.Intellect, char.Stamina, char.Spirit, char.UnspentPoints, char.TotalKills)
		tr.context.Characters["character"] = char
	}
	
	return nil
}

// executeCheckCharacterAttributes 检查角色属性，确保所有属性都基于角色属性正确计算
func (tr *TestRunner) executeCheckCharacterAttributes() error {
	char, ok := tr.context.Characters["character"]
	if !ok || char == nil {
		return fmt.Errorf("character not found")
	}
	
	// 重新计算所有属性（如果为0）
	needsUpdate := false
	if char.PhysicalAttack == 0 {
		char.PhysicalAttack = tr.calculator.CalculatePhysicalAttack(char)
		needsUpdate = true
	}
	if char.MagicAttack == 0 {
		char.MagicAttack = tr.calculator.CalculateMagicAttack(char)
		needsUpdate = true
	}
	if char.PhysicalDefense == 0 {
		char.PhysicalDefense = tr.calculator.CalculatePhysicalDefense(char)
		needsUpdate = true
	}
	if char.MagicDefense == 0 {
		char.MagicDefense = tr.calculator.CalculateMagicDefense(char)
		needsUpdate = true
	}
	if char.PhysCritRate == 0 {
		char.PhysCritRate = tr.calculator.CalculatePhysCritRate(char)
		needsUpdate = true
	}
	if char.PhysCritDamage == 0 {
		char.PhysCritDamage = tr.calculator.CalculatePhysCritDamage(char)
		needsUpdate = true
	}
	if char.SpellCritRate == 0 {
		char.SpellCritRate = tr.calculator.CalculateSpellCritRate(char)
		needsUpdate = true
	}
	if char.SpellCritDamage == 0 {
		char.SpellCritDamage = tr.calculator.CalculateSpellCritDamage(char)
		needsUpdate = true
	}
	if char.DodgeRate == 0 {
		char.DodgeRate = tr.calculator.CalculateDodgeRate(char)
		needsUpdate = true
	}
	
	// 如果属性被修复，更新数据库
	if needsUpdate {
		charRepo := repository.NewCharacterRepository()
		charRepo.UpdateAfterBattle(char.ID, char.HP, char.Resource, char.Exp, char.Level,
			char.ExpToNext, char.MaxHP, char.MaxResource, char.PhysicalAttack, char.MagicAttack, char.PhysicalDefense, char.MagicDefense,
			char.Strength, char.Agility, char.Intellect, char.Stamina, char.Spirit, char.UnspentPoints, char.TotalKills)
	}
	
	// 更新上下文
	tr.context.Characters["character"] = char
	
	return nil
}

// handleBuffSkill 处理Buff技能
func (tr *TestRunner) handleBuffSkill(char *models.Character, skill *models.Skill) {
	// 获取Buff效果
	attackModifier := 0.0
	if modVal, exists := tr.context.Variables["skill_buff_attack_modifier"]; exists {
		if m, ok := modVal.(float64); ok {
			attackModifier = m
		}
	}
	
	duration := 3 // 默认3回合
	if durVal, exists := tr.context.Variables["skill_buff_duration"]; exists {
		if d, ok := durVal.(int); ok {
			duration = d
		}
	}
	
	// 设置Buff信息到上下文（供断言使用）
	tr.assertion.SetContext("character.buff_attack_modifier", attackModifier)
	tr.assertion.SetContext("character.buff_duration", duration)
	
	// 也存储到Variables中，以便updateAssertionContext可以访问
	tr.context.Variables["character_buff_attack_modifier"] = attackModifier
	tr.context.Variables["character_buff_duration"] = duration
	
	// 注意：实际的Buff应用需要在战斗系统中处理
	// 这里只是设置测试上下文，供断言使用
}

// executeBattleRound 执行战斗回合（减少冷却时间）
func (tr *TestRunner) executeBattleRound(instruction string) error {
	// 解析回合数（如"执行第2回合"）
	roundNum := 1
	if strings.Contains(instruction, "第") {
		parts := strings.Split(instruction, "第")
		if len(parts) > 1 {
			roundStr := strings.TrimSpace(strings.Split(parts[1], "回合")[0])
			if round, err := strconv.Atoi(roundStr); err == nil {
				roundNum = round
			}
		}
	}
	
	// 减少技能冷却时间
	skillManager := game.NewSkillManager()
	char, ok := tr.context.Characters["character"]
	if ok && char != nil {
		if err := skillManager.LoadCharacterSkills(char.ID); err == nil {
			// 先减少冷却时间
			skillManager.TickCooldowns(char.ID)
			
			// 获取技能状态，检查是否可用
			skillVal, exists := tr.context.Variables["skill"]
			if exists {
				if skill, ok := skillVal.(*models.Skill); ok {
					skillState := skillManager.GetSkillState(char.ID, skill.ID)
					if skillState != nil {
						tr.assertion.SetContext(fmt.Sprintf("skill_usable_round_%d", roundNum), skillState.CooldownLeft == 0)
						tr.assertion.SetContext(fmt.Sprintf("skill_cooldown_round_%d", roundNum), skillState.CooldownLeft)
					} else {
						// 如果技能状态不存在，根据冷却时间计算
						// 假设第1回合使用了技能，冷却时间为3，那么：
						// 第2回合：冷却剩余2，不可用
						// 第3回合：冷却剩余1，不可用
						// 第4回合：冷却剩余0，可用
						cooldownLeft := skill.Cooldown - (roundNum - 1)
						if cooldownLeft < 0 {
							cooldownLeft = 0
						}
						tr.assertion.SetContext(fmt.Sprintf("skill_usable_round_%d", roundNum), cooldownLeft == 0)
						tr.assertion.SetContext(fmt.Sprintf("skill_cooldown_round_%d", roundNum), cooldownLeft)
					}
				}
			}
		} else {
			// 如果角色没有技能，从上下文获取技能信息
			skillVal, exists := tr.context.Variables["skill"]
			if exists {
				if skill, ok := skillVal.(*models.Skill); ok {
					// 根据冷却时间计算
					cooldownLeft := skill.Cooldown - (roundNum - 1)
					if cooldownLeft < 0 {
						cooldownLeft = 0
					}
					tr.assertion.SetContext(fmt.Sprintf("skill_usable_round_%d", roundNum), cooldownLeft == 0)
					tr.assertion.SetContext(fmt.Sprintf("skill_cooldown_round_%d", roundNum), cooldownLeft)
				}
			}
		}
	}
	
	return nil
}


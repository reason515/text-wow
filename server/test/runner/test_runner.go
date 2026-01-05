package runner

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"text-wow/internal/game"
	"text-wow/internal/models"
	"text-wow/internal/repository"

	"gopkg.in/yaml.v3"
)

// TestContext 测试上下文 - 存储测试过程中的数据
type TestContext struct {
	UserID     int
	User       *models.User
	Characters map[string]*models.Character // key: "character", "character_1", "character_2" 等
	Monsters   map[string]*models.Monster   // key: "monster", "monster_1", "monster_2" 等
	Team       []*models.Character          // 队伍角色列表
	BattleManager *game.BattleManager
	Calculator    *game.Calculator
	LastDamage    int                        // 最后一次伤害值
	LastHealing   int                        // 最后一次治疗值
	BattleLogs    []string                   // 战斗日志
	BattleResult  *game.BattleTickResult     // 最后一次战斗结果
}

// TestRunner 测试运行器
type TestRunner struct {
	parser     *YAMLParser
	assertion  *AssertionExecutor
	reporter   *Reporter
	context    *TestContext
}

// NewTestRunner 创建测试运行器
func NewTestRunner() *TestRunner {
	tr := &TestRunner{
		parser:    NewYAMLParser(),
		assertion: NewAssertionExecutor(),
		reporter:  NewReporter(),
		context:   &TestContext{
			Characters: make(map[string]*models.Character),
			Monsters:   make(map[string]*models.Monster),
			Team:       make([]*models.Character, 0),
			BattleLogs: make([]string, 0),
		},
	}
	
	// 初始化游戏系统管理器
	tr.context.BattleManager = game.NewBattleManager()
	tr.context.Calculator = game.NewCalculator()
	
	// 设置断言执行器的测试上下文引用
	tr.assertion.SetTestContext(tr.context)
	
	return tr
}

// ResetContext 重置测试上下文
func (tr *TestRunner) ResetContext() {
	tr.context = &TestContext{
		Characters: make(map[string]*models.Character),
		Monsters:   make(map[string]*models.Monster),
		Team:       make([]*models.Character, 0),
		BattleLogs: make([]string, 0),
	}
	tr.context.BattleManager = game.NewBattleManager()
	tr.context.Calculator = game.NewCalculator()
	tr.assertion.ClearContext()
	tr.assertion.SetTestContext(tr.context)
	// 确保断言执行器可以访问Calculator
	if tr.assertion != nil && tr.context != nil {
		tr.assertion.SetTestContext(tr.context)
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
	// 重置测试上下文
	tr.ResetContext()
	
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

	// 更新上下文数据（从战斗会话中获取最新状态）
	tr.updateContextFromBattle()

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
		if err := tr.parseAndExecuteSetupInstruction(instruction); err != nil {
			return fmt.Errorf("failed to execute setup instruction '%s': %w", instruction, err)
		}
	}
	return nil
}

// parseAndExecuteSetupInstruction 解析并执行setup指令
func (tr *TestRunner) parseAndExecuteSetupInstruction(instruction string) error {
	instruction = strings.TrimSpace(instruction)
	
	// 创建角色指令
	if strings.Contains(instruction, "创建一个") && strings.Contains(instruction, "角色") {
		return tr.createCharacterFromInstruction(instruction)
	}
	
	// 创建怪物指令
	if strings.Contains(instruction, "创建一个") && strings.Contains(instruction, "怪物") {
		return tr.createMonsterFromInstruction(instruction)
	}
	
	// 创建队伍指令
	if strings.Contains(instruction, "创建一个") && strings.Contains(instruction, "队伍") {
		return tr.createTeamFromInstruction(instruction)
	}
	
	// 初始化战斗系统
	if strings.Contains(instruction, "初始化战斗系统") {
		// 战斗系统已经在 NewTestRunner 中初始化，这里只需要确保用户存在
		if tr.context.UserID == 0 {
			user, err := tr.createTestUser()
			if err != nil {
				return err
			}
			tr.context.UserID = user.ID
			tr.context.User = user
		}
		return nil
	}
	
	// 创建/准备技能指令
	if strings.Contains(instruction, "技能") && (strings.Contains(instruction, "创建") || strings.Contains(instruction, "准备")) {
		return tr.createSkillFromInstruction(instruction)
	}
	
	// 学习技能指令（在setup中）
	if strings.Contains(instruction, "学习") && strings.Contains(instruction, "技能") {
		return tr.learnSkillFromInstruction(instruction)
	}
	
	return fmt.Errorf("unknown setup instruction: %s", instruction)
}

// createCharacterFromInstruction 从指令创建角色
// 示例: "创建一个1级人类战士角色，HP=25，攻击力=8"
func (tr *TestRunner) createCharacterFromInstruction(instruction string) error {
	// 解析等级
	level := 1
	if strings.Contains(instruction, "级") {
		levelStr := ""
		for i, char := range instruction {
			if char >= '0' && char <= '9' {
				levelStr += string(char)
			} else if char == '级' && levelStr != "" {
				if parsedLevel, err := strconv.Atoi(levelStr); err == nil {
					level = parsedLevel
				}
				break
			} else if i > 0 && (char < '0' || char > '9') {
				levelStr = ""
			}
		}
	}
	
	// 解析种族和职业
	raceID := "human"
	classID := "warrior"
	if strings.Contains(instruction, "人类") {
		raceID = "human"
	} else if strings.Contains(instruction, "兽人") {
		raceID = "orc"
	}
	
	if strings.Contains(instruction, "战士") {
		classID = "warrior"
	} else if strings.Contains(instruction, "法师") {
		classID = "mage"
	} else if strings.Contains(instruction, "牧师") {
		classID = "priest"
	} else if strings.Contains(instruction, "盗贼") {
		classID = "rogue"
	}
	
	// 解析主属性（力量、敏捷、智力、耐力、精神）
	strength := 10
	agility := 10
	intellect := 10
	stamina := 10
	spirit := 10
	
	// 解析力量=xxx
	if strings.Contains(instruction, "力量=") {
		strengthStr := extractValueAfter(instruction, "力量=")
		if parsedStrength, err := strconv.Atoi(strengthStr); err == nil {
			strength = parsedStrength
		}
	}
	
	// 解析敏捷=xxx
	if strings.Contains(instruction, "敏捷=") {
		agilityStr := extractValueAfter(instruction, "敏捷=")
		if parsedAgility, err := strconv.Atoi(agilityStr); err == nil {
			agility = parsedAgility
		}
	}
	
	// 解析智力=xxx
	if strings.Contains(instruction, "智力=") {
		intellectStr := extractValueAfter(instruction, "智力=")
		if parsedIntellect, err := strconv.Atoi(intellectStr); err == nil {
			intellect = parsedIntellect
		}
	}
	
	// 解析耐力=xxx
	if strings.Contains(instruction, "耐力=") {
		staminaStr := extractValueAfter(instruction, "耐力=")
		if parsedStamina, err := strconv.Atoi(staminaStr); err == nil {
			stamina = parsedStamina
		}
	}
	
	// 解析精神=xxx
	if strings.Contains(instruction, "精神=") {
		spiritStr := extractValueAfter(instruction, "精神=")
		if parsedSpirit, err := strconv.Atoi(spiritStr); err == nil {
			spirit = parsedSpirit
		}
	}
	
	// 解析技能点=xxx
	unspentPoints := 0
	if strings.Contains(instruction, "技能点=") {
		pointsStr := extractValueAfter(instruction, "技能点=")
		if parsedPoints, err := strconv.Atoi(pointsStr); err == nil {
			unspentPoints = parsedPoints
		}
	}
	
	// 解析属性
	hp := 100
	physicalAttack := 10
	physicalDefense := 5
	magicAttack := 5
	magicDefense := 3
	
	// 解析基础HP（用于计算最大HP）
	baseHP := 35 // 默认战士基础HP
	if strings.Contains(instruction, "基础HP=") {
		baseHPStr := extractValueAfter(instruction, "基础HP=")
		if parsedBaseHP, err := strconv.Atoi(baseHPStr); err == nil {
			baseHP = parsedBaseHP
		}
	}
	
	// 解析 HP=xxx（直接设置HP，不计算）
	if strings.Contains(instruction, "HP=") {
		hpStr := extractValueAfter(instruction, "HP=")
		if parsedHP, err := strconv.Atoi(hpStr); err == nil {
			hp = parsedHP
		}
	} else {
		// 如果没有直接设置HP，根据耐力和基础HP计算
		hp = tr.context.Calculator.CalculateHP(&models.Character{Stamina: stamina}, baseHP)
	}
	
	// 解析攻击力=xxx（直接设置，不计算）
	if strings.Contains(instruction, "攻击力=") {
		attackStr := extractValueAfter(instruction, "攻击力=")
		if parsedAttack, err := strconv.Atoi(attackStr); err == nil {
			physicalAttack = parsedAttack
		}
	} else {
		// 如果没有直接设置攻击力，根据属性计算
		physicalAttack = tr.context.Calculator.CalculatePhysicalAttack(&models.Character{
			Strength: strength,
			Agility:  agility,
		})
	}
	
	// 解析防御力=xxx
	if strings.Contains(instruction, "防御力=") {
		defenseStr := extractValueAfter(instruction, "防御力=")
		if parsedDefense, err := strconv.Atoi(defenseStr); err == nil {
			physicalDefense = parsedDefense
		}
	}
	
	// 确保用户存在
	if tr.context.UserID == 0 {
		user, err := tr.createTestUser()
		if err != nil {
			return err
		}
		tr.context.UserID = user.ID
		tr.context.User = user
	}
	
	// 创建角色
	char := &models.Character{
		UserID:          tr.context.UserID,
		Name:            fmt.Sprintf("测试角色_%d", len(tr.context.Team)+1),
		RaceID:          raceID,
		ClassID:         classID,
		Faction:         "alliance",
		TeamSlot:        len(tr.context.Team) + 1,
		Level:           level,
		HP:              hp,
		MaxHP:           hp,
		Resource:        100,
		MaxResource:     100,
		ResourceType:    "rage",
		PhysicalAttack:  physicalAttack,
		MagicAttack:     magicAttack,
		PhysicalDefense: physicalDefense,
		MagicDefense:    magicDefense,
		Strength:        strength,
		Agility:         agility,
		Intellect:       intellect,
		Stamina:         stamina,
		Spirit:          spirit,
		UnspentPoints:   unspentPoints,
	}
	
	// 使用Calculator计算派生属性
	char.PhysicalAttack = tr.context.Calculator.CalculatePhysicalAttack(char)
	char.MagicAttack = tr.context.Calculator.CalculateMagicAttack(char)
	char.PhysCritRate = tr.context.Calculator.CalculatePhysCritRate(char)
	char.PhysCritDamage = tr.context.Calculator.CalculatePhysCritDamage(char)
	char.SpellCritRate = tr.context.Calculator.CalculateSpellCritRate(char)
	char.SpellCritDamage = tr.context.Calculator.CalculateSpellCritDamage(char)
	char.DodgeRate = tr.context.Calculator.CalculateDodgeRate(char)
	
	// 如果直接设置了攻击力，覆盖计算值
	if strings.Contains(instruction, "攻击力=") {
		attackStr := extractValueAfter(instruction, "攻击力=")
		if parsedAttack, err := strconv.Atoi(attackStr); err == nil {
			char.PhysicalAttack = parsedAttack
		}
	}
	
	// 确保MaxHP正确设置
	if char.MaxHP < char.HP {
		char.MaxHP = char.HP
	}
	}
	
	// 根据职业设置资源类型
	if classID == "mage" || classID == "priest" {
		char.ResourceType = "mana"
		char.MaxResource = 200
		char.Resource = 200
	}
	
	// 保存到数据库
	charRepo := repository.NewCharacterRepository()
	createdChar, err := charRepo.Create(char)
	if err != nil {
		return fmt.Errorf("failed to create character: %w", err)
	}
	
	// 加载角色的技能列表
	skillRepo := repository.NewSkillRepository()
	characterSkills, err := skillRepo.GetCharacterSkills(createdChar.ID)
	if err == nil {
		createdChar.Skills = characterSkills
	}
	
	// 存储到上下文
	key := "character"
	if len(tr.context.Team) > 0 {
		key = fmt.Sprintf("character_%d", len(tr.context.Team))
	}
	tr.context.Characters[key] = createdChar
	tr.context.Team = append(tr.context.Team, createdChar)
	
	return nil
}

// createMonsterFromInstruction 从指令创建怪物
// 示例: "创建一个1级森林狼怪物，HP=20，攻击力=5"
func (tr *TestRunner) createMonsterFromInstruction(instruction string) error {
	// 解析等级
	level := 1
	if strings.Contains(instruction, "级") {
		levelStr := ""
		for i, char := range instruction {
			if char >= '0' && char <= '9' {
				levelStr += string(char)
			} else if char == '级' && levelStr != "" {
				if parsedLevel, err := strconv.Atoi(levelStr); err == nil {
					level = parsedLevel
				}
				break
			} else if i > 0 && (char < '0' || char > '9') {
				levelStr = ""
			}
		}
	}
	
	// 解析怪物名称和类型
	monsterName := "测试怪物"
	monsterType := "normal"
	if strings.Contains(instruction, "森林狼") {
		monsterName = "森林狼"
	} else if strings.Contains(instruction, "精英") {
		monsterType = "elite"
		monsterName = "精英怪物"
	} else if strings.Contains(instruction, "Boss") || strings.Contains(instruction, "boss") {
		monsterType = "boss"
		monsterName = "Boss怪物"
	}
	
	// 解析属性
	hp := 50
	physicalAttack := 10
	physicalDefense := 5
	magicAttack := 5
	magicDefense := 3
	
	// 解析 HP=xxx
	if strings.Contains(instruction, "HP=") {
		hpStr := extractValueAfter(instruction, "HP=")
		if parsedHP, err := strconv.Atoi(hpStr); err == nil {
			hp = parsedHP
		}
	}
	
	// 解析攻击力=xxx
	if strings.Contains(instruction, "攻击力=") {
		attackStr := extractValueAfter(instruction, "攻击力=")
		if parsedAttack, err := strconv.Atoi(attackStr); err == nil {
			physicalAttack = parsedAttack
		}
	}
	
	// 解析防御力=xxx
	if strings.Contains(instruction, "防御力=") {
		defenseStr := extractValueAfter(instruction, "防御力=")
		if parsedDefense, err := strconv.Atoi(defenseStr); err == nil {
			physicalDefense = parsedDefense
		}
	}
	
	// 创建怪物
	monster := &models.Monster{
		ID:              fmt.Sprintf("test_monster_%d", len(tr.context.Monsters)+1),
		Name:            monsterName,
		Type:            monsterType,
		Level:           level,
		HP:              hp,
		MaxHP:           hp,
		PhysicalAttack:  physicalAttack,
		MagicAttack:     magicAttack,
		PhysicalDefense: physicalDefense,
		MagicDefense:    magicDefense,
		Speed:           10,
		ExpReward:       20,
		GoldMin:         1,
		GoldMax:         10,
		PhysCritRate:    0.05,
		PhysCritDamage:  1.5,
		SpellCritRate:   0.05,
		SpellCritDamage: 1.5,
		DodgeRate:       0.05,
	}
	
	// 存储到上下文
	key := "monster"
	if len(tr.context.Monsters) > 0 {
		key = fmt.Sprintf("monster_%d", len(tr.context.Monsters))
	}
	tr.context.Monsters[key] = monster
	
	return nil
}

// createTeamFromInstruction 从指令创建队伍
// 示例: "创建一个3人队伍：战士(坦克)、牧师(治疗)、法师(DPS)"
func (tr *TestRunner) createTeamFromInstruction(instruction string) error {
	// 解析队伍人数
	teamSize := 3
	if strings.Contains(instruction, "人队伍") {
		teamSizeStr := ""
		for i, char := range instruction {
			if char >= '0' && char <= '9' {
				teamSizeStr += string(char)
			} else if strings.Contains(instruction[i:], "人队伍") && teamSizeStr != "" {
				if parsedSize, err := strconv.Atoi(teamSizeStr); err == nil {
					teamSize = parsedSize
				}
				break
			}
		}
	}
	
	// 解析角色配置
	roles := []struct {
		classID string
		role    string
	}{
		{"warrior", "tank"},
		{"priest", "healer"},
		{"mage", "dps"},
	}
	
	// 根据指令解析角色
	if strings.Contains(instruction, "战士") {
		roles[0] = struct {
			classID string
			role    string
		}{"warrior", "tank"}
	}
	if strings.Contains(instruction, "牧师") {
		if teamSize >= 2 {
			roles[1] = struct {
				classID string
				role    string
			}{"priest", "healer"}
		}
	}
	if strings.Contains(instruction, "法师") {
		if teamSize >= 3 {
			roles[2] = struct {
				classID string
				role    string
			}{"mage", "dps"}
		}
	}
	
	// 创建队伍角色
	for i := 0; i < teamSize && i < len(roles); i++ {
		role := roles[i]
		instruction := fmt.Sprintf("创建一个1级人类%s角色，HP=100，攻击力=10", getClassChineseName(role.classID))
		if err := tr.createCharacterFromInstruction(instruction); err != nil {
			return err
		}
	}
	
	return nil
}

// extractValueAfter 从字符串中提取指定关键字后的值（支持负数）
func extractValueAfter(s, keyword string) string {
	idx := strings.Index(s, keyword)
	if idx == -1 {
		return ""
	}
	
	start := idx + len(keyword)
	value := ""
	// 检查是否有负号
	if start < len(s) && s[start] == '-' {
		value += "-"
		start++
	}
	for i := start; i < len(s); i++ {
		if s[i] >= '0' && s[i] <= '9' {
			value += string(s[i])
		} else if value != "" && value != "-" {
			break
		}
	}
	return value
}

// getClassChineseName 获取职业中文名称
func getClassChineseName(classID string) string {
	switch classID {
	case "warrior":
		return "战士"
	case "mage":
		return "法师"
	case "priest":
		return "牧师"
	case "rogue":
		return "盗贼"
	default:
		return "战士"
	}
}

// createTestUser 创建测试用户
func (tr *TestRunner) createTestUser() (*models.User, error) {
	userRepo := repository.NewUserRepository()
	user, err := userRepo.GetByUsername("test_user")
	if err != nil {
		// 用户不存在，创建新用户
		passwordHash := "$2a$10$test_hash_for_testing"
		user, err = userRepo.Create("test_user", passwordHash, "test@test.com")
		if err != nil {
			return nil, fmt.Errorf("failed to create test user: %w", err)
		}
	}
	return user, nil
}

// executeStep 执行测试步骤
func (tr *TestRunner) executeStep(step TestStep) error {
	action := strings.TrimSpace(step.Action)
	
	// 开始战斗
	if action == "开始战斗" || strings.Contains(action, "开始战斗") {
		return tr.executeStartBattle()
	}
	
	// 执行一个回合
	if action == "执行一个回合" || strings.Contains(action, "执行一个回合") {
		return tr.executeBattleRound()
	}
	
	// 怪物反击
	if action == "怪物反击" || strings.Contains(action, "怪物反击") {
		return tr.executeMonsterAttack()
	}
	
	// 继续战斗直到条件满足
	if strings.Contains(action, "继续战斗直到") {
		condition := extractCondition(action, "继续战斗直到")
		maxRounds := step.Timeout
		if maxRounds == 0 {
			maxRounds = 20 // 默认最大回合数
		}
		return tr.executeBattleUntil(condition, maxRounds)
	}
	
	// 角色使用技能
	if strings.Contains(action, "角色使用技能") || strings.Contains(action, "使用技能") {
		skillName := extractSkillName(action)
		return tr.executeUseSkill(skillName)
	}
	
	// 角色学习技能
	if strings.Contains(action, "学习技能") || strings.Contains(action, "角色学习技能") {
		return tr.learnSkillFromInstruction(action)
	}
	
	// 计算基础伤害
	if strings.Contains(action, "计算基础伤害") || (strings.Contains(action, "计算") && strings.Contains(action, "基础伤害")) {
		return tr.executeCalculateBaseDamage()
	}
	
	// 计算伤害（通用）
	if strings.Contains(action, "计算") && strings.Contains(action, "伤害") {
		return tr.executeCalculateDamage(action)
	}
	
	// 应用防御减伤或计算防御减伤
	if strings.Contains(action, "应用防御减伤") || strings.Contains(action, "计算防御减伤") {
		return tr.executeCalculateDefenseReduction()
	}
	
	// 计算减伤后伤害
	if strings.Contains(action, "计算减伤后伤害") {
		return tr.executeCalculateDefenseReduction()
	}
	
	// 如果暴击，应用暴击倍率
	if strings.Contains(action, "暴击") && strings.Contains(action, "暴击倍率") {
		return tr.executeApplyCrit()
	}
	
	// 计算物理攻击力
	if action == "计算物理攻击力" || strings.Contains(action, "计算物理攻击力") {
		return tr.executeCalculatePhysicalAttack()
	}
	
	// 计算法术攻击力
	if action == "计算法术攻击力" || strings.Contains(action, "计算法术攻击力") {
		return tr.executeCalculateMagicAttack()
	}
	
	// 计算最大生命值
	if action == "计算最大生命值" || strings.Contains(action, "计算最大生命值") {
		return tr.executeCalculateMaxHP()
	}
	
	// 计算物理暴击率
	if action == "计算物理暴击率" || strings.Contains(action, "计算物理暴击率") {
		return tr.executeCalculatePhysCritRate()
	}
	
	// 计算闪避率
	if action == "计算闪避率" || strings.Contains(action, "计算闪避率") {
		return tr.executeCalculateDodgeRate()
	}
	
	// 计算速度
	if action == "计算速度" || strings.Contains(action, "计算速度") {
		return tr.executeCalculateSpeed()
	}
	
	// 计算速度
	if action == "计算速度" || strings.Contains(action, "计算速度") {
		return tr.executeCalculateSpeed()
	}
	
	// 计算法力恢复
	if strings.Contains(action, "计算法力恢复") {
		return tr.executeCalculateManaRegen(action)
	}
	
	// 计算怒气获得
	if strings.Contains(action, "计算怒气获得") {
		return tr.executeCalculateRageGain(action)
	}
	
	// 计算能量恢复
	if strings.Contains(action, "计算能量恢复") {
		return tr.executeCalculateEnergyRegen(action)
	}
	
	return fmt.Errorf("unknown action: %s", action)
}

// executeStartBattle 开始战斗
func (tr *TestRunner) executeStartBattle() error {
	if tr.context.UserID == 0 {
		user, err := tr.createTestUser()
		if err != nil {
			return err
		}
		tr.context.UserID = user.ID
		tr.context.User = user
	}
	
	if len(tr.context.Team) == 0 {
		return fmt.Errorf("no characters in team")
	}
	
	// 设置战斗区域
	session := tr.context.BattleManager.GetOrCreateSession(tr.context.UserID)
	if session.CurrentZone == nil {
		gameRepo := repository.NewGameRepository()
		zone, err := gameRepo.GetZoneByID("elwynn")
		if err == nil {
			session.CurrentZone = zone
		}
	}
	
	// 如果没有怪物，创建一个
	if len(tr.context.Monsters) == 0 {
		// 使用第一个角色的等级生成怪物
		charLevel := tr.context.Team[0].Level
		monster := tr.context.Monsters["monster"]
		if monster == nil {
			// 创建一个默认怪物
			monster = &models.Monster{
				ID:              "test_monster",
				Name:            "测试怪物",
				Type:            "normal",
				Level:           charLevel,
				HP:              50,
				MaxHP:           50,
				PhysicalAttack:  10,
				MagicAttack:     5,
				PhysicalDefense: 5,
				MagicDefense:    3,
				Speed:           10,
				ExpReward:       20,
				GoldMin:         1,
				GoldMax:         10,
			}
			tr.context.Monsters["monster"] = monster
		}
		session.CurrentEnemies = []*models.Monster{monster}
	} else {
		// 使用上下文中的怪物
		enemies := make([]*models.Monster, 0)
		for _, monster := range tr.context.Monsters {
			enemies = append(enemies, monster)
		}
		session.CurrentEnemies = enemies
	}
	
	// 开始战斗
	_, err := tr.context.BattleManager.StartBattle(tr.context.UserID)
	return err
}

// executeBattleRound 执行一个战斗回合
func (tr *TestRunner) executeBattleRound() error {
	if tr.context.UserID == 0 || len(tr.context.Team) == 0 {
		return fmt.Errorf("battle not initialized")
	}
	
	// 执行战斗tick
	result, err := tr.context.BattleManager.ExecuteBattleTick(tr.context.UserID, tr.context.Team)
	if err != nil {
		return err
	}
	
	// 更新上下文
	if result != nil {
		tr.context.BattleResult = result
		if result.Character != nil {
			// 更新角色状态
			for i, char := range tr.context.Team {
				if char.ID == result.Character.ID {
					tr.context.Team[i] = result.Character
					tr.context.Characters["character"] = result.Character
					break
				}
			}
		}
		
		// 更新战斗日志
		for _, log := range result.Logs {
			tr.context.BattleLogs = append(tr.context.BattleLogs, log.Message)
		}
		
		// 更新伤害和治疗值
		if result.DamageDealt > 0 {
			tr.context.LastDamage = result.DamageDealt
		}
		if result.HealingDone > 0 {
			tr.context.LastHealing = result.HealingDone
		}
	}
	
	return nil
}

// executeMonsterAttack 执行怪物攻击
func (tr *TestRunner) executeMonsterAttack() error {
	// 怪物攻击实际上是通过执行回合来实现的
	return tr.executeBattleRound()
}

// executeBattleUntil 继续战斗直到条件满足
func (tr *TestRunner) executeBattleUntil(condition string, maxRounds int) error {
	for i := 0; i < maxRounds; i++ {
		// 检查条件
		if tr.checkCondition(condition) {
			return nil
		}
		
		// 执行一个回合
		if err := tr.executeBattleRound(); err != nil {
			return err
		}
		
		// 检查战斗是否结束
		session := tr.context.BattleManager.GetSession(tr.context.UserID)
		if session == nil || !session.IsRunning {
			break
		}
	}
	
	return nil
}

// checkCondition 检查条件是否满足
func (tr *TestRunner) checkCondition(condition string) bool {
	condition = strings.TrimSpace(condition)
	
	// 怪物死亡
	if condition == "怪物死亡" || strings.Contains(condition, "怪物死亡") {
		if len(tr.context.Monsters) > 0 {
			monster := tr.context.Monsters["monster"]
			if monster != nil && monster.HP <= 0 {
				return true
			}
		}
		// 检查战斗会话中的敌人
		session := tr.context.BattleManager.GetSession(tr.context.UserID)
		if session != nil {
			allDead := true
			for _, enemy := range session.CurrentEnemies {
				if enemy != nil && enemy.HP > 0 {
					allDead = false
					break
				}
			}
			return allDead && len(session.CurrentEnemies) > 0
		}
	}
	
	// 角色死亡
	if condition == "角色死亡" || strings.Contains(condition, "角色死亡") {
		if len(tr.context.Team) > 0 {
			char := tr.context.Team[0]
			if char.HP <= 0 {
				return true
			}
		}
	}
	
	return false
}

// executeUseSkill 执行使用技能
func (tr *TestRunner) executeUseSkill(skillName string) error {
	if len(tr.context.Team) == 0 {
		return fmt.Errorf("character not found")
	}
	
	char := tr.context.Team[0]
	
	// 获取技能ID（从技能名称或ID）
	skillID := skillName
	if skillName == "" {
		// 如果没有指定技能名称，尝试使用第一个可用技能
		skillRepo := repository.NewSkillRepository()
		characterSkills, err := skillRepo.GetCharacterSkills(char.ID)
		if err != nil || len(characterSkills) == 0 {
			return fmt.Errorf("no skills available")
		}
		skillID = characterSkills[0].SkillID
	}
	
	// 通过BattleManager使用技能
	if tr.context.BattleManager == nil {
		return fmt.Errorf("battle manager not initialized")
	}
	
	// 确保战斗会话存在
	if tr.context.UserID == 0 {
		user, err := tr.createTestUser()
		if err != nil {
			return err
		}
		tr.context.UserID = user.ID
		tr.context.User = user
	}
	
	// 确保战斗已开始
	session := tr.context.BattleManager.GetSession(tr.context.UserID)
	if session == nil {
		// 如果没有战斗会话，创建一个简单的战斗会话用于测试
		if len(tr.context.Monsters) == 0 {
			// 创建一个测试怪物
			monster := &models.Monster{
				ID:             1,
				Name:           "测试怪物",
				HP:             100,
				MaxHP:          100,
				PhysicalAttack: 10,
				PhysicalDefense: 5,
			}
			tr.context.Monsters["monster"] = monster
		}
		// 开始战斗
		if err := tr.executeStartBattle(); err != nil {
			return fmt.Errorf("failed to start battle: %w", err)
		}
		session = tr.context.BattleManager.GetSession(tr.context.UserID)
	}
	
	// 使用SkillManager直接使用技能
	// 创建SkillManager实例（用于测试）
	skillManager := game.NewSkillManager()
	
	// 加载角色技能（如果还没有加载）
	if err := skillManager.LoadCharacterSkills(char.ID); err != nil {
		tr.assertion.SetContext("skill_used", false)
		tr.assertion.SetContext("error_message", fmt.Sprintf("failed to load skills: %v", err))
		return err
	}
	
	// 检查技能是否可用
	skillState := skillManager.GetSkillState(char.ID, skillID)
	if skillState == nil {
		// 尝试带前缀
		if char.ClassID == "warrior" {
			skillState = skillManager.GetSkillState(char.ID, "warrior_"+skillID)
		}
		if skillState == nil {
			tr.assertion.SetContext("skill_used", false)
			tr.assertion.SetContext("error_message", fmt.Sprintf("skill not found: %s", skillID))
			return fmt.Errorf("skill not found: %s", skillID)
		}
	}
	
	// 检查资源是否足够
	if skillState.Skill.ResourceCost > char.Resource {
		tr.assertion.SetContext("skill_used", false)
		tr.assertion.SetContext("error_message", "资源不足")
		return fmt.Errorf("insufficient resource: need %d, have %d", skillState.Skill.ResourceCost, char.Resource)
	}
	
	// 检查冷却时间
	if skillState.CooldownLeft > 0 {
		tr.assertion.SetContext("skill_used", false)
		tr.assertion.SetContext("error_message", fmt.Sprintf("skill on cooldown: %d turns left", skillState.CooldownLeft))
		return fmt.Errorf("skill on cooldown: %d turns left", skillState.CooldownLeft)
	}
	
	// 使用技能
	usedState, err := skillManager.UseSkill(char.ID, skillID)
	if err != nil {
		tr.assertion.SetContext("skill_used", false)
		tr.assertion.SetContext("error_message", err.Error())
		return err
	}
	
	// 消耗资源
	char.Resource -= skillState.Skill.ResourceCost
	if char.Resource < 0 {
		char.Resource = 0
	}
	
	// 存储技能使用结果
	tr.assertion.SetContext("skill_used", true)
	tr.assertion.SetContext("skill_cooldown_round_1", usedState.CooldownLeft)
	
	// 如果战斗会话存在，执行技能效果
	if session != nil && len(session.CurrentEnemies) > 0 {
		target := session.CurrentEnemies[0]
		if target != nil {
			// 应用技能效果
			skillEffects := skillManager.ApplySkillEffects(usedState, char, target)
			
			// 计算伤害/治疗
			if skillState.Skill.Type == "attack" {
				// 计算伤害
				baseDamage := skillState.Skill.BaseValue
				if skillState.Skill.ScalingStat != "" {
					// 根据属性计算伤害
					var statValue int
					switch skillState.Skill.ScalingStat {
					case "strength":
						statValue = char.Strength
					case "agility":
						statValue = char.Agility
					case "intellect":
						statValue = char.Intellect
					}
					baseDamage = int(float64(baseDamage) + float64(statValue)*skillState.Skill.ScalingRatio)
				}
				
				// 应用技能等级倍率
				levelMultiplier := 1.0 + float64(skillState.SkillLevel-1)*0.15
				baseDamage = int(float64(baseDamage) * levelMultiplier)
				
				// 计算最终伤害（考虑防御）
				damage := baseDamage - target.PhysicalDefense
				if damage < 0 {
					damage = 0
				}
				
				tr.assertion.SetContext("skill_damage_dealt", damage)
				target.HP -= damage
				if target.HP < 0 {
					target.HP = 0
				}
			} else if skillState.Skill.Type == "heal" {
				// 计算治疗
				healing := skillState.Skill.BaseValue
				levelMultiplier := 1.0 + float64(skillState.SkillLevel-1)*0.15
				healing = int(float64(healing) * levelMultiplier)
				
				tr.assertion.SetContext("skill_healing_done", healing)
				char.HP += healing
				if char.HP > char.MaxHP {
					char.HP = char.MaxHP
				}
			}
			
			// 应用Buff/Debuff
			if session.BuffManager != nil {
				for effectID, effectValue := range skillEffects {
					if strings.HasPrefix(effectID, "buff_") || strings.HasPrefix(effectID, "debuff_") {
						// 这里可以应用Buff/Debuff，但需要更复杂的逻辑
						_ = effectValue
					}
				}
			}
		}
	}
	
	// 更新上下文
	tr.updateContextFromBattle()
	
	return nil
}

// executeCalculateBaseDamage 计算基础伤害
func (tr *TestRunner) executeCalculateBaseDamage() error {
	if len(tr.context.Team) == 0 {
		return fmt.Errorf("character not found")
	}
	
	char := tr.context.Team[0]
	// 基础伤害 = 攻击力 × 技能系数（默认1.0）
	baseDamage := char.PhysicalAttack
	
	// 设置到断言上下文
	tr.assertion.SetContext("base_damage", baseDamage)
	// 同时设置到测试上下文（用于后续步骤）
	tr.context.LastDamage = baseDamage
	
	return nil
}

// executeCalculateDamage 计算伤害
func (tr *TestRunner) executeCalculateDamage(action string) error {
	// 从上下文中获取角色和怪物
	if len(tr.context.Team) == 0 || len(tr.context.Monsters) == 0 {
		return fmt.Errorf("character or monster not found")
	}
	
	char := tr.context.Team[0]
	monster := tr.context.Monsters["monster"]
	if monster == nil {
		return fmt.Errorf("monster not found")
	}
	
	// 创建一个临时的角色对象作为防御者（用于伤害计算）
	defender := &models.Character{
		PhysicalDefense: monster.PhysicalDefense,
		MagicDefense:    monster.MagicDefense,
		DodgeRate:       monster.DodgeRate,
	}
	
	// 使用计算器计算伤害
	result := tr.context.Calculator.CalculateDamage(
		char,
		defender,
		char.PhysicalAttack,
		1.0, // 技能倍率
		"physical",
		false, // 不忽略闪避
	)
	
	tr.context.LastDamage = result.FinalDamage
	
	// 将伤害值存储到断言上下文中
	tr.assertion.SetContext("base_damage", int(result.BaseDamage))
	tr.assertion.SetContext("damage_after_defense", int(result.DamageAfterDefense))
	tr.assertion.SetContext("final_damage", result.FinalDamage)
	
	return nil
}

// executeCalculateDefenseReduction 计算防御减伤
func (tr *TestRunner) executeCalculateDefenseReduction() error {
	if len(tr.context.Team) == 0 || len(tr.context.Monsters) == 0 {
		return fmt.Errorf("character or monster not found")
	}
	
	char := tr.context.Team[0]
	monster := tr.context.Monsters["monster"]
	if monster == nil {
		return fmt.Errorf("monster not found")
	}
	
	// 获取基础伤害（如果已计算）
	baseDamage, ok := tr.assertion.context["base_damage"].(int)
	if !ok {
		// 如果没有基础伤害，使用攻击力
		baseDamage = char.PhysicalAttack
		tr.assertion.SetContext("base_damage", baseDamage)
	}
	
	// 应用防御减伤（减法公式）
	damageAfterDefense := baseDamage - monster.PhysicalDefense
	if damageAfterDefense < 1 {
		damageAfterDefense = 1 // 至少1点伤害
	}
	
	tr.context.LastDamage = damageAfterDefense
	tr.assertion.SetContext("damage_after_defense", damageAfterDefense)
	// 如果没有最终伤害，使用减伤后伤害作为最终伤害
	if _, exists := tr.assertion.context["final_damage"]; !exists {
		tr.assertion.SetContext("final_damage", damageAfterDefense)
	}
	
	return nil
}

// executeApplyCrit 应用暴击倍率
func (tr *TestRunner) executeApplyCrit() error {
	// 从上下文中获取伤害值
	var baseDamage int
	var ok bool
	if damage, exists := tr.assertion.context["damage_after_defense"]; exists {
		baseDamage, ok = damage.(int)
	}
	if !ok {
		baseDamage = tr.context.LastDamage
		if baseDamage == 0 {
			// 如果没有伤害值，尝试从角色和怪物计算
			if len(tr.context.Team) > 0 && len(tr.context.Monsters) > 0 {
				char := tr.context.Team[0]
				monster := tr.context.Monsters["monster"]
				if monster != nil {
					baseDamage = char.PhysicalAttack - monster.PhysicalDefense
					if baseDamage < 1 {
						baseDamage = 1
					}
				}
			}
		}
	}
	
	if len(tr.context.Team) == 0 {
		return fmt.Errorf("character not found")
	}
	
	char := tr.context.Team[0]
	// 假设暴击（实际应该随机判断）
	finalDamage := int(float64(baseDamage) * char.PhysCritDamage)
	
	tr.context.LastDamage = finalDamage
	tr.assertion.SetContext("final_damage", finalDamage)
	
	return nil
}

// extractCondition 从动作中提取条件
func extractCondition(action, prefix string) string {
	idx := strings.Index(action, prefix)
	if idx == -1 {
		return ""
	}
	return strings.TrimSpace(action[idx+len(prefix):])
}

// extractSkillName 从动作中提取技能名称
func extractSkillName(action string) string {
	// 简化实现，实际应该更智能地解析
	if strings.Contains(action, "技能") {
		parts := strings.Split(action, "技能")
		if len(parts) > 1 {
			return strings.TrimSpace(parts[1])
		}
	}
	return ""
}

// executeCalculatePhysicalAttack 计算物理攻击力
func (tr *TestRunner) executeCalculatePhysicalAttack() error {
	if len(tr.context.Team) == 0 {
		return fmt.Errorf("character not found")
	}
	if tr.assertion == nil {
		return fmt.Errorf("assertion executor is nil")
	}
	char := tr.context.Team[0]
	physicalAttack := tr.context.Calculator.CalculatePhysicalAttack(char)
	// 存储到断言上下文
	tr.assertion.SetContext("physical_attack", physicalAttack)
	// 同时更新角色的物理攻击力（用于后续断言）
	char.PhysicalAttack = physicalAttack
	return nil
}

// executeCalculateMagicAttack 计算法术攻击力
func (tr *TestRunner) executeCalculateMagicAttack() error {
	if len(tr.context.Team) == 0 {
		return fmt.Errorf("character not found")
	}
	char := tr.context.Team[0]
	magicAttack := tr.context.Calculator.CalculateMagicAttack(char)
	// 存储到断言上下文
	tr.assertion.SetContext("magic_attack", magicAttack)
	// 同时更新角色的法术攻击力（用于后续断言）
	char.MagicAttack = magicAttack
	return nil
}

// executeCalculateMaxHP 计算最大生命值
func (tr *TestRunner) executeCalculateMaxHP() error {
	if len(tr.context.Team) == 0 {
		return fmt.Errorf("character not found")
	}
	char := tr.context.Team[0]
	// 获取职业基础HP
	baseHP := 35 // 默认战士基础HP
	if char.ClassID == "mage" || char.ClassID == "warlock" {
		baseHP = 20
	} else if char.ClassID == "priest" || char.ClassID == "druid" || char.ClassID == "shaman" {
		baseHP = 22
	} else if char.ClassID == "rogue" || char.ClassID == "hunter" {
		baseHP = 25
	} else if char.ClassID == "paladin" {
		baseHP = 30
	}
	maxHP := tr.context.Calculator.CalculateHP(char, baseHP)
	// 存储到断言上下文
	tr.assertion.SetContext("max_hp", maxHP)
	// 同时更新角色的最大生命值（用于后续断言）
	char.MaxHP = maxHP
	return nil
}

// executeCalculatePhysCritRate 计算物理暴击率
func (tr *TestRunner) executeCalculatePhysCritRate() error {
	if len(tr.context.Team) == 0 {
		return fmt.Errorf("character not found")
	}
	char := tr.context.Team[0]
	critRate := tr.context.Calculator.CalculatePhysCritRate(char)
	// 存储到断言上下文
	tr.assertion.SetContext("phys_crit_rate", critRate)
	// 同时更新角色的物理暴击率（用于后续断言）
	char.PhysCritRate = critRate
	return nil
}

// executeCalculateDodgeRate 计算闪避率
func (tr *TestRunner) executeCalculateDodgeRate() error {
	if len(tr.context.Team) == 0 {
		return fmt.Errorf("character not found")
	}
	char := tr.context.Team[0]
	dodgeRate := tr.context.Calculator.CalculateDodgeRate(char)
	// 存储到断言上下文
	tr.assertion.SetContext("dodge_rate", dodgeRate)
	// 同时更新角色的闪避率（用于后续断言）
	char.DodgeRate = dodgeRate
	return nil
}

// executeCalculateSpeed 计算速度
func (tr *TestRunner) executeCalculateSpeed() error {
	if len(tr.context.Team) == 0 {
		return fmt.Errorf("character not found")
	}
	char := tr.context.Team[0]
	speed := tr.context.Calculator.CalculateSpeed(char)
	// 存储到断言上下文
	tr.assertion.SetContext("speed", speed)
	return nil
}

// executeCalculateManaRegen 计算法力恢复
func (tr *TestRunner) executeCalculateManaRegen(action string) error {
	if len(tr.context.Team) == 0 {
		return fmt.Errorf("character not found")
	}
	char := tr.context.Team[0]
	
	// 从action中提取基础恢复值，默认10
	baseRegen := 10
	if strings.Contains(action, "基础恢复=") {
		baseRegenStr := extractValueAfter(action, "基础恢复=")
		if parsedBaseRegen, err := strconv.Atoi(baseRegenStr); err == nil {
			baseRegen = parsedBaseRegen
		}
	}
	
	manaRegen := tr.context.Calculator.CalculateManaRegen(char, baseRegen)
	tr.assertion.SetContext("mana_regen", manaRegen)
	return nil
}

// executeCalculateRageGain 计算怒气获得
func (tr *TestRunner) executeCalculateRageGain(action string) error {
	// 从action中提取基础获得和加成百分比
	baseGain := 10
	bonusPercent := 0.0
	
	// 解析基础获得
	if strings.Contains(action, "基础获得=") || strings.Contains(action, "基础怒气获得=") {
		var baseGainStr string
		if strings.Contains(action, "基础获得=") {
			baseGainStr = extractValueAfter(action, "基础获得=")
		} else {
			baseGainStr = extractValueAfter(action, "基础怒气获得=")
		}
		if parsedBaseGain, err := strconv.Atoi(baseGainStr); err == nil {
			baseGain = parsedBaseGain
		}
	}
	
	// 解析加成百分比
	if strings.Contains(action, "加成百分比=") || strings.Contains(action, "加成=") {
		var bonusStr string
		if strings.Contains(action, "加成百分比=") {
			bonusStr = extractValueAfter(action, "加成百分比=")
		} else {
			bonusStr = extractValueAfter(action, "加成=")
		}
		// 移除%符号
		bonusStr = strings.TrimSuffix(bonusStr, "%")
		if parsedBonus, err := strconv.ParseFloat(bonusStr, 64); err == nil {
			bonusPercent = parsedBonus
		}
	}
	
	rageGain := tr.context.Calculator.CalculateRageGain(baseGain, bonusPercent)
	tr.assertion.SetContext("rage_gain", rageGain)
	return nil
}

// executeCalculateEnergyRegen 计算能量恢复
func (tr *TestRunner) executeCalculateEnergyRegen(action string) error {
	if len(tr.context.Team) == 0 {
		return fmt.Errorf("character not found")
	}
	char := tr.context.Team[0]
	
	// 从action中提取基础恢复值，默认10
	baseRegen := 10
	if strings.Contains(action, "基础恢复=") {
		baseRegenStr := extractValueAfter(action, "基础恢复=")
		if parsedBaseRegen, err := strconv.Atoi(baseRegenStr); err == nil {
			baseRegen = parsedBaseRegen
		}
	}
	
	energyRegen := tr.context.Calculator.CalculateEnergyRegen(char, baseRegen)
	tr.assertion.SetContext("energy_regen", energyRegen)
	return nil
}

// updateContextFromBattle 从战斗会话更新上下文数据
func (tr *TestRunner) updateContextFromBattle() {
	if tr.context.UserID == 0 {
		return
	}
	
	session := tr.context.BattleManager.GetSession(tr.context.UserID)
	if session == nil {
		return
	}
	
	// 更新怪物状态
	if len(session.CurrentEnemies) > 0 {
		for i, enemy := range session.CurrentEnemies {
			if enemy != nil {
				key := "monster"
				if i > 0 {
					key = fmt.Sprintf("monster_%d", i)
				}
				tr.context.Monsters[key] = enemy
			}
		}
	}
	
	// 更新角色状态（从数据库重新加载以确保最新）
	if len(tr.context.Team) > 0 {
		charRepo := repository.NewCharacterRepository()
		for i, char := range tr.context.Team {
			reloaded, err := charRepo.GetByID(char.ID)
			if err == nil {
				tr.context.Team[i] = reloaded
				if i == 0 {
					tr.context.Characters["character"] = reloaded
				} else {
					tr.context.Characters[fmt.Sprintf("character_%d", i)] = reloaded
				}
			}
		}
	}
	
	// 更新战斗日志
	if len(session.BattleLogs) > 0 {
		tr.context.BattleLogs = make([]string, 0, len(session.BattleLogs))
		for _, log := range session.BattleLogs {
			tr.context.BattleLogs = append(tr.context.BattleLogs, log.Message)
		}
	}
	
	// 更新一些统计值到断言上下文
	tr.assertion.SetContext("team_alive_count", tr.countAliveCharacters())
	tr.assertion.SetContext("team_total_exp", tr.calculateTeamTotalExp())
}

// countAliveCharacters 计算存活角色数量
func (tr *TestRunner) countAliveCharacters() int {
	count := 0
	for _, char := range tr.context.Team {
		if char != nil && char.HP > 0 {
			count++
		}
	}
	return count
}

// calculateTeamTotalExp 计算队伍总经验值
func (tr *TestRunner) calculateTeamTotalExp() int {
	total := 0
	for _, char := range tr.context.Team {
		if char != nil {
			total += char.Exp
		}
	}
	return total
}

// executeTeardown 执行清理
func (tr *TestRunner) executeTeardown(teardown []string) error {
	for _, instruction := range teardown {
		instruction = strings.TrimSpace(instruction)
		
		// 清理战斗状态
		if strings.Contains(instruction, "清理战斗状态") || strings.Contains(instruction, "清理战斗") {
			if tr.context.UserID > 0 {
				tr.context.BattleManager.StopBattle(tr.context.UserID)
			}
		}
		
		// 重置角色数据
		if strings.Contains(instruction, "重置角色数据") || strings.Contains(instruction, "重置角色") {
			// 从数据库重新加载角色数据
			if len(tr.context.Team) > 0 {
				charRepo := repository.NewCharacterRepository()
				for i, char := range tr.context.Team {
					reloaded, err := charRepo.GetByID(char.ID)
					if err == nil {
						tr.context.Team[i] = reloaded
						tr.context.Characters["character"] = reloaded
					}
				}
			}
		}
	}
	
	// 清理上下文数据（但保留角色和怪物引用，因为可能还需要用于断言）
	tr.context.LastDamage = 0
	tr.context.LastHealing = 0
	tr.context.BattleLogs = make([]string, 0)
	tr.context.BattleResult = nil
	
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

// createSkillFromInstruction 从指令创建/准备技能
// 示例: "准备一个可学习的技能：冲锋（Charge）"
// 示例: "创建一个消耗30点怒气的技能"
// 示例: "创建一个冷却时间为3回合的技能"
func (tr *TestRunner) createSkillFromInstruction(instruction string) error {
	// 解析技能名称
	skillName := ""
	skillID := ""
	
	// 提取技能名称（如"冲锋"、"Charge"）
	if strings.Contains(instruction, "：") {
		parts := strings.Split(instruction, "：")
		if len(parts) > 1 {
			skillName = strings.TrimSpace(parts[1])
			// 提取括号中的ID
			if strings.Contains(skillName, "（") && strings.Contains(skillName, "）") {
				start := strings.Index(skillName, "（")
				end := strings.Index(skillName, "）")
				if start < end {
					skillID = strings.ToLower(strings.TrimSpace(skillName[start+3 : end]))
					skillName = strings.TrimSpace(skillName[:start])
				}
			}
		}
	}
	
	// 如果没有指定名称，生成一个默认的
	if skillID == "" {
		skillID = "test_skill_" + fmt.Sprintf("%d", len(tr.context.Team)+1)
	}
	if skillName == "" {
		skillName = "测试技能"
	}
	
	// 解析技能属性
	resourceCost := 0
	cooldown := 0
	baseValue := 0
	skillType := "attack"
	targetType := "enemy"
	
	// 解析资源消耗
	if strings.Contains(instruction, "消耗") && strings.Contains(instruction, "点") {
		costStr := extractValueAfter(instruction, "消耗")
		if parsedCost, err := strconv.Atoi(costStr); err == nil {
			resourceCost = parsedCost
		}
	}
	
	// 解析冷却时间
	if strings.Contains(instruction, "冷却时间") {
		cooldownStr := extractValueAfter(instruction, "冷却时间")
		if parsedCooldown, err := strconv.Atoi(cooldownStr); err == nil {
			cooldown = parsedCooldown
		}
	}
	
	// 解析伤害倍率
	if strings.Contains(instruction, "伤害倍率") {
		multiplierStr := extractValueAfter(instruction, "伤害倍率")
		multiplierStr = strings.TrimSuffix(multiplierStr, "%")
		if parsedMultiplier, err := strconv.ParseFloat(multiplierStr, 64); err == nil {
			baseValue = int(parsedMultiplier)
		}
	}
	
	// 解析治疗量
	if strings.Contains(instruction, "治疗量") {
		healStr := extractValueAfter(instruction, "治疗量")
		if parsedHeal, err := strconv.Atoi(healStr); err == nil {
			baseValue = parsedHeal
			skillType = "heal"
			targetType = "self"
		}
	}
	
	// 解析AOE技能
	if strings.Contains(instruction, "AOE") {
		targetType = "enemy_all"
	}
	
	// 解析Buff技能
	if strings.Contains(instruction, "Buff") {
		skillType = "buff"
		targetType = "self"
	}
	
	// 创建技能（存储在测试上下文中，不写入数据库）
	// 这里我们只是准备技能，实际学习需要调用learnSkillFromInstruction
	// 将技能信息存储到上下文中，供后续使用
	skillKey := "prepared_skill_" + skillID
	tr.assertion.SetContext(skillKey, map[string]interface{}{
		"id":            skillID,
		"name":          skillName,
		"resource_cost": resourceCost,
		"cooldown":      cooldown,
		"base_value":    baseValue,
		"type":          skillType,
		"target_type":   targetType,
	})
	
	return nil
}

// learnSkillFromInstruction 从指令学习技能
// 示例: "角色学习技能：冲锋"
func (tr *TestRunner) learnSkillFromInstruction(instruction string) error {
	if len(tr.context.Team) == 0 {
		return fmt.Errorf("character not found")
	}
	
	char := tr.context.Team[0]
	
	// 提取技能名称
	skillName := ""
	if strings.Contains(instruction, "：") {
		parts := strings.Split(instruction, "：")
		if len(parts) > 1 {
			skillName = strings.TrimSpace(parts[1])
		}
	}
	
	// 如果没有指定技能名称，尝试从上下文中获取准备的技能
	skillID := ""
	if skillName != "" {
		// 尝试从技能名称转换为ID
		skillID = strings.ToLower(skillName)
		// 如果是中文名称，尝试映射到ID
		skillNameMap := map[string]string{
			"冲锋":   "charge",
			"charge": "charge",
		}
		if mappedID, exists := skillNameMap[skillName]; exists {
			skillID = mappedID
		}
	}
	
	// 如果还是没有ID，尝试从准备的技能中获取
	if skillID == "" {
		// 查找第一个准备的技能
		for key, value := range tr.assertion.context {
			if strings.HasPrefix(key, "prepared_skill_") {
				if skillMap, ok := value.(map[string]interface{}); ok {
					if id, exists := skillMap["id"].(string); exists {
						skillID = id
						break
					}
				}
			}
		}
	}
	
	if skillID == "" {
		return fmt.Errorf("skill ID not found")
	}
	
	// 检查技能是否存在
	skillRepo := repository.NewSkillRepository()
	skill, err := skillRepo.GetSkillByID(skillID)
	if err != nil {
		// 如果技能不存在，创建一个临时技能用于测试
		skill = &models.Skill{
			ID:           skillID,
			Name:         skillName,
			Type:         "attack",
			TargetType:   "enemy",
			BaseValue:    0,
			ResourceCost: 0,
			Cooldown:     0,
			ClassID:      char.ClassID,
		}
	}
	
	// 学习技能
	err = skillRepo.AddCharacterSkill(char.ID, skillID, 1)
	if err != nil {
		tr.assertion.SetContext("skill_learned", false)
		tr.assertion.SetContext("error_message", err.Error())
		return err
	}
	
	tr.assertion.SetContext("skill_learned", true)
	
	// 重新加载角色技能
	characterSkills, _ := skillRepo.GetCharacterSkills(char.ID)
	char.Skills = characterSkills
	
	// 更新上下文中的角色对象
	for i, teamChar := range tr.context.Team {
		if teamChar.ID == char.ID {
			tr.context.Team[i].Skills = characterSkills
			break
		}
	}
	for key, contextChar := range tr.context.Characters {
		if contextChar.ID == char.ID {
			tr.context.Characters[key].Skills = characterSkills
			break
		}
	}
	
	return nil
}


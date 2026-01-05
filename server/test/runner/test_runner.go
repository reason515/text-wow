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
	
	// 解析属性
	hp := 100
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
		PhysCritRate:    0.05,
		PhysCritDamage:  1.5,
		SpellCritRate:   0.05,
		SpellCritDamage: 1.5,
		DodgeRate:       0.05,
		Strength:        10,
		Agility:         10,
		Intellect:       10,
		Stamina:         10,
		Spirit:          10,
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

// extractValueAfter 从字符串中提取指定关键字后的值
func extractValueAfter(s, keyword string) string {
	idx := strings.Index(s, keyword)
	if idx == -1 {
		return ""
	}
	
	start := idx + len(keyword)
	value := ""
	for i := start; i < len(s); i++ {
		if s[i] >= '0' && s[i] <= '9' {
			value += string(s[i])
		} else if value != "" {
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
	
	// 计算伤害
	if strings.Contains(action, "计算") && strings.Contains(action, "伤害") {
		return tr.executeCalculateDamage(action)
	}
	
	// 应用防御减伤
	if strings.Contains(action, "应用防御减伤") || strings.Contains(action, "计算防御减伤") {
		return tr.executeCalculateDefenseReduction()
	}
	
	// 如果暴击，应用暴击倍率
	if strings.Contains(action, "暴击") && strings.Contains(action, "暴击倍率") {
		return tr.executeApplyCrit()
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
	// 这里需要实现技能使用逻辑
	// 暂时返回nil，因为技能系统可能需要更复杂的集成
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
	
	// 计算基础伤害
	baseDamage := char.PhysicalAttack
	// 应用防御减伤（减法公式）
	damageAfterDefense := baseDamage - monster.PhysicalDefense
	if damageAfterDefense < 1 {
		damageAfterDefense = 1 // 至少1点伤害
	}
	
	tr.context.LastDamage = damageAfterDefense
	tr.assertion.SetContext("damage_after_defense", damageAfterDefense)
	
	return nil
}

// executeApplyCrit 应用暴击倍率
func (tr *TestRunner) executeApplyCrit() error {
	// 从上下文中获取伤害值
	baseDamage, ok := tr.assertion.context["damage_after_defense"].(int)
	if !ok {
		baseDamage = tr.context.LastDamage
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


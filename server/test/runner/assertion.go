package runner

import (
	"fmt"
	"math"
	"reflect"
	"strconv"
	"strings"
)

// AssertionExecutor 断言执行器
type AssertionExecutor struct {
	context      map[string]interface{} // 测试上下文（存储测试数据）
	testContext  *TestContext            // 测试上下文引用
}

// NewAssertionExecutor 创建断言执行器
func NewAssertionExecutor() *AssertionExecutor {
	return &AssertionExecutor{
		context: make(map[string]interface{}),
	}
}

// SetTestContext 设置测试上下文引用
func (ae *AssertionExecutor) SetTestContext(ctx *TestContext) {
	ae.testContext = ctx
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
	
	// 强制检查：如果err为nil但actual也为nil，这是不应该发生的
	if err == nil && actual == nil {
		result.Status = "failed"
		teamLen := 0
		hasCalculator := false
		if ae.testContext != nil {
			teamLen = len(ae.testContext.Team)
			hasCalculator = ae.testContext.Calculator != nil
		}
		result.Error = fmt.Sprintf("getValue returned (nil, nil) for path: %s (context keys: %v, testContext: %v, team: %d, calculator: %v)", 
			assertion.Target, getMapKeys(ae.context), ae.testContext != nil, teamLen, hasCalculator)
		result.Actual = nil
		return result
	}
	
	if err != nil {
		result.Status = "failed"
		result.Error = fmt.Sprintf("failed to get value: %v", err)
		result.Actual = nil
		return result
	}

	// 检查值是否为nil（使用反射来正确处理interface{}类型的nil）
	// 首先检查简单的nil比较
	if actual == nil {
		result.Status = "failed"
		teamLen := 0
		hasCalculator := false
		if ae.testContext != nil {
			teamLen = len(ae.testContext.Team)
			hasCalculator = ae.testContext.Calculator != nil
		}
		result.Error = fmt.Sprintf("value is nil for path: %s (context keys: %v, testContext: %v, team: %d, calculator: %v, actual type: %T)", 
			assertion.Target, getMapKeys(ae.context), ae.testContext != nil, teamLen, hasCalculator, actual)
		result.Actual = nil
		return result
	}
	
	// 使用反射检查interface{}类型的nil（即使在actual != nil时也要检查，因为可能是interface{}(nil)）
	rv := reflect.ValueOf(actual)
	if !rv.IsValid() || (rv.Kind() == reflect.Interface && rv.IsNil()) {
		result.Status = "failed"
		teamLen := 0
		hasCalculator := false
		if ae.testContext != nil {
			teamLen = len(ae.testContext.Team)
			hasCalculator = ae.testContext.Calculator != nil
		}
		result.Error = fmt.Sprintf("value is nil (interface{}) for path: %s (context keys: %v, testContext: %v, team: %d, calculator: %v)", 
			assertion.Target, getMapKeys(ae.context), ae.testContext != nil, teamLen, hasCalculator)
		result.Actual = nil
		return result
	}
	
	// 检查其他类型的nil
	if rv.IsValid() && (rv.Kind() == reflect.Ptr || rv.Kind() == reflect.Slice || rv.Kind() == reflect.Map || rv.Kind() == reflect.Chan || rv.Kind() == reflect.Func) {
		if rv.IsNil() {
			result.Status = "failed"
			teamLen := 0
			hasCalculator := false
			if ae.testContext != nil {
				teamLen = len(ae.testContext.Team)
				hasCalculator = ae.testContext.Calculator != nil
			}
			result.Error = fmt.Sprintf("value is nil (pointer/slice/map) for path: %s (context keys: %v, testContext: %v, team: %d, calculator: %v)", 
				assertion.Target, getMapKeys(ae.context), ae.testContext != nil, teamLen, hasCalculator)
			result.Actual = nil
			return result
		}
	}

	// 在设置 result.Actual 之前，再次检查 actual 是否为 nil
	// 这应该不会发生，因为前面的检查应该已经捕获了这种情况
	if actual == nil {
		result.Status = "failed"
		teamLen := 0
		hasCalculator := false
		if ae.testContext != nil {
			teamLen = len(ae.testContext.Team)
			hasCalculator = ae.testContext.Calculator != nil
		}
		result.Error = fmt.Sprintf("actual is nil before setting result.Actual for path: %s (context keys: %v, testContext: %v, team: %d, calculator: %v)", 
			assertion.Target, getMapKeys(ae.context), ae.testContext != nil, teamLen, hasCalculator)
		result.Actual = nil
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
	// 检查 context 是否为 nil
	if ae.context == nil {
		return nil, fmt.Errorf("context is nil for path: %s", path)
	}
	
	// 首先尝试从简单上下文获取
	if value, exists := ae.context[path]; exists {
		if value == nil {
			return nil, fmt.Errorf("value is nil for path: %s (key exists but value is nil)", path)
		}
		return value, nil
	}

	// 尝试解析为数字
	if num, err := strconv.Atoi(path); err == nil {
		return num, nil
	}

	// 如果有测试上下文，尝试路径解析
	if ae.testContext != nil {
		value, err := ae.resolvePath(path)
		if err == nil {
			// 如果解析成功但值为nil，返回错误
			if value == nil {
				return nil, fmt.Errorf("resolvePath returned nil for path: %s", path)
			}
			// 将计算得到的值存储到上下文中，以便后续使用
			if !strings.Contains(path, ".") {
				ae.context[path] = value
			}
			return value, nil
		}
		// 如果路径解析失败，尝试作为简单键再次查找
		if val, exists := ae.context[path]; exists {
			if val == nil {
				return nil, fmt.Errorf("value is nil for path: %s (key exists but value is nil)", path)
			}
			return val, nil
		}
		// 确保总是返回错误，而不是(nil, nil)
		if err == nil {
			return nil, fmt.Errorf("resolvePath returned (nil, nil) for path: %s (context keys: %v)", path, getMapKeys(ae.context))
		}
		return nil, fmt.Errorf("value not found: %s (context keys: %v, resolvePath error: %v)", path, getMapKeys(ae.context), err)
	}

	// 最终检查：如果到达这里，说明没有找到值
	// 确保总是返回错误，而不是(nil, nil)
	err := fmt.Errorf("value not found: %s (no test context, context keys: %v)", path, getMapKeys(ae.context))
	return nil, err
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

// resolvePath 解析路径，支持点号分隔和数组索引
func (ae *AssertionExecutor) resolvePath(path string) (interface{}, error) {
	// 如果没有测试上下文，无法解析路径
	if ae.testContext == nil {
		return nil, fmt.Errorf("test context not set")
	}
	
	parts := strings.Split(path, ".")
	if len(parts) == 0 {
		return nil, fmt.Errorf("invalid path: %s", path)
	}

	// 处理根对象
	var current interface{}
	root := parts[0]

	switch root {
	case "character", "character_0":
		if len(ae.testContext.Characters) > 0 {
			if char, exists := ae.testContext.Characters["character"]; exists {
				current = char
			} else if len(ae.testContext.Team) > 0 {
				current = ae.testContext.Team[0]
			}
		}
	case "monster", "monster_0":
		if len(ae.testContext.Monsters) > 0 {
			if monster, exists := ae.testContext.Monsters["monster"]; exists {
				current = monster
			}
		}
	case "last_damage":
		return ae.testContext.LastDamage, nil
	case "last_healing":
		return ae.testContext.LastHealing, nil
	case "battle_logs":
		return strings.Join(ae.testContext.BattleLogs, "\n"), nil
	case "team_alive_count":
		return ae.countAliveCharacters(), nil
	case "team_total_exp":
		return ae.calculateTeamTotalExp(), nil
	case "enemy_alive_count":
		return ae.countAliveMonsters(), nil
	case "battle_state":
		return ae.getBattleState(), nil
	case "battle_result":
		return ae.getBattleResult(), nil
	// 计算属性（从上下文获取，如果上下文中没有则尝试计算）
	case "physical_attack", "magic_attack", "max_hp", "phys_crit_rate", "dodge_rate", "speed",
		"mana_regen", "rage_gain", "energy_regen", "base_damage", "damage_after_defense", "final_damage",
		"overhealing", "skill_damage_dealt", "skill_used", "equip_success", "error_message":
		// 首先尝试从上下文获取
		if val, exists := ae.context[root]; exists {
			if val == nil {
				return nil, fmt.Errorf("value is nil for path: %s (key exists but value is nil)", root)
			}
			return val, nil
		}
		// 如果上下文中没有，尝试从角色计算
		if ae.testContext != nil && len(ae.testContext.Team) > 0 && ae.testContext.Calculator != nil {
			char := ae.testContext.Team[0]
			var calculatedValue interface{}
			switch root {
			case "physical_attack":
				calculatedValue = ae.testContext.Calculator.CalculatePhysicalAttack(char)
			case "magic_attack":
				calculatedValue = ae.testContext.Calculator.CalculateMagicAttack(char)
			case "max_hp":
				baseHP := 35
				if char.ClassID == "mage" || char.ClassID == "warlock" {
					baseHP = 20
				} else if char.ClassID == "priest" || char.ClassID == "druid" || char.ClassID == "shaman" {
					baseHP = 22
				} else if char.ClassID == "rogue" || char.ClassID == "hunter" {
					baseHP = 25
				} else if char.ClassID == "paladin" {
					baseHP = 30
				}
				calculatedValue = ae.testContext.Calculator.CalculateHP(char, baseHP)
			case "phys_crit_rate":
				calculatedValue = ae.testContext.Calculator.CalculatePhysCritRate(char)
			case "dodge_rate":
				calculatedValue = ae.testContext.Calculator.CalculateDodgeRate(char)
			case "speed":
				calculatedValue = ae.testContext.Calculator.CalculateSpeed(char)
			case "mana_regen":
				calculatedValue = ae.testContext.Calculator.CalculateManaRegen(char)
			case "rage_gain":
				calculatedValue = ae.testContext.Calculator.CalculateRageGain(char)
			case "energy_regen":
				calculatedValue = ae.testContext.Calculator.CalculateEnergyRegen(char)
			default:
				return nil, fmt.Errorf("calculated attribute not handled: %s", root)
			}
			// 将计算得到的值存储到上下文中
			if calculatedValue != nil {
				ae.context[root] = calculatedValue
				return calculatedValue, nil
			}
			return nil, fmt.Errorf("calculated value is nil for path: %s", root)
		}
		return nil, fmt.Errorf("value not found: %s (testContext: %v, team: %d, calculator: %v)", root, ae.testContext != nil, len(ae.testContext.Team), ae.testContext != nil && ae.testContext.Calculator != nil)
	case "warrior", "priest", "mage", "rogue":
		// 通过职业名称查找角色
		current = ae.findCharacterByClass(root)
	case "character_1", "character_2", "character_3", "character_4":
		// 通过索引查找角色
		idx := 0
		if strings.HasPrefix(root, "character_") {
			if parsedIdx, err := strconv.Atoi(strings.TrimPrefix(root, "character_")); err == nil {
				idx = parsedIdx
			}
		}
		if idx < len(ae.testContext.Team) {
			current = ae.testContext.Team[idx]
		}
	case "monster_1", "monster_2", "monster_3":
		// 通过索引查找怪物
		idx := 0
		if strings.HasPrefix(root, "monster_") {
			if parsedIdx, err := strconv.Atoi(strings.TrimPrefix(root, "monster_")); err == nil {
				idx = parsedIdx
			}
		}
		keys := make([]string, 0, len(ae.testContext.Monsters))
		for k := range ae.testContext.Monsters {
			keys = append(keys, k)
		}
		if idx < len(keys) {
			current = ae.testContext.Monsters[keys[idx]]
		}
	default:
		// 尝试从简单上下文获取
		if value, exists := ae.context[root]; exists {
			current = value
		} else {
			// 尝试解析为数字
			if num, err := strconv.Atoi(root); err == nil {
				return num, nil
			}
			return nil, fmt.Errorf("unknown root object: %s", root)
		}
	}

	if current == nil {
		return nil, fmt.Errorf("object not found: %s", root)
	}

	// 处理嵌套路径
	for i := 1; i < len(parts); i++ {
		part := parts[i]
		
		// 检查是否是数组索引
		if strings.Contains(part, "[") && strings.Contains(part, "]") {
			// 解析数组索引，如 "characters[0]"
			idxStart := strings.Index(part, "[")
			idxEnd := strings.Index(part, "]")
			if idxStart == -1 || idxEnd == -1 {
				return nil, fmt.Errorf("invalid array index: %s", part)
			}
			
			arrayName := part[:idxStart]
			idxStr := part[idxStart+1 : idxEnd]
			idx, err := strconv.Atoi(idxStr)
			if err != nil {
				return nil, fmt.Errorf("invalid array index: %s", idxStr)
			}
			
			// 获取数组
			var arr interface{}
			switch arrayName {
			case "characters":
				if idx < len(ae.testContext.Team) {
					arr = ae.testContext.Team[idx]
				}
			case "monsters":
				// 从monsters map中获取
				keys := make([]string, 0, len(ae.testContext.Monsters))
				for k := range ae.testContext.Monsters {
					keys = append(keys, k)
				}
				if idx < len(keys) {
					arr = ae.testContext.Monsters[keys[idx]]
				}
			default:
				return nil, fmt.Errorf("unknown array: %s", arrayName)
			}
			
			if arr == nil {
				return nil, fmt.Errorf("array index out of range: %s[%d]", arrayName, idx)
			}
			
			current = arr
		} else {
			// 普通属性访问
			// 如果是map类型，直接访问
			if mapValue, ok := current.(map[string]interface{}); ok {
				if val, exists := mapValue[part]; exists {
					current = val
				} else {
					return nil, fmt.Errorf("field not found: %s in %s", part, strings.Join(parts[:i], "."))
				}
			} else {
				current = ae.getFieldValue(current, part)
				if current == nil {
					return nil, fmt.Errorf("field not found: %s in %s", part, strings.Join(parts[:i], "."))
				}
			}
		}
	}

	// 确保不会返回 (nil, nil)
	if current == nil {
		return nil, fmt.Errorf("object is nil at end of resolvePath for path: %s", path)
	}

	return current, nil
}

// getFieldValue 获取结构体字段值（使用反射或类型断言）
func (ae *AssertionExecutor) getFieldValue(obj interface{}, fieldName string) interface{} {
	switch v := obj.(type) {
	case *models.Character:
		return ae.getCharacterField(v, fieldName)
	case *models.Monster:
		return ae.getMonsterField(v, fieldName)
	case *models.CharacterSkill:
		return ae.getCharacterSkillField(v, fieldName)
	case map[string]interface{}:
		return v[fieldName]
	default:
		return nil
	}
}

// getCharacterSkillField 获取角色技能字段值
func (ae *AssertionExecutor) getCharacterSkillField(skill *models.CharacterSkill, fieldName string) interface{} {
	switch fieldName {
	case "level", "Level", "skill_level", "skillLevel":
		return skill.SkillLevel
	case "id", "ID", "skill_id", "skillId":
		return skill.SkillID
	case "exp", "Exp", "skill_exp", "skillExp":
		return skill.SkillExp
	case "exp_to_next", "expToNext", "ExpToNext":
		return skill.ExpToNext
	case "damage_multiplier", "damageMultiplier":
		// 根据技能等级计算伤害倍率
		// 1级: 100%, 2级: 115%, 3级: 130%, 4级: 145%, 5级: 160%
		multipliers := map[int]float64{
			1: 1.0,
			2: 1.15,
			3: 1.30,
			4: 1.45,
			5: 1.60,
		}
		if multiplier, exists := multipliers[skill.SkillLevel]; exists {
			return multiplier
		}
		return 1.0
	default:
		return nil
	}
}

// getCharacterField 获取角色字段值
func (ae *AssertionExecutor) getCharacterField(char *models.Character, fieldName string) interface{} {
	switch fieldName {
	case "hp", "HP":
		return char.HP
	case "max_hp", "maxHp", "MaxHP":
		return char.MaxHP
	case "resource", "Resource":
		return char.Resource
	case "max_resource", "maxResource", "MaxResource":
		return char.MaxResource
	case "level", "Level":
		return char.Level
	case "exp", "Exp":
		return char.Exp
	case "physical_attack", "physicalAttack", "PhysicalAttack":
		return char.PhysicalAttack
	case "magic_attack", "magicAttack", "MagicAttack":
		return char.MagicAttack
	case "physical_defense", "physicalDefense", "PhysicalDefense":
		return char.PhysicalDefense
	case "magic_defense", "magicDefense", "MagicDefense":
		return char.MagicDefense
	case "strength", "Strength":
		return char.Strength
	case "agility", "Agility":
		return char.Agility
	case "intellect", "Intellect":
		return char.Intellect
	case "stamina", "Stamina":
		return char.Stamina
	case "spirit", "Spirit":
		return char.Spirit
	case "phys_crit_rate", "physCritRate":
		return char.PhysCritRate
	case "phys_crit_damage", "physCritDamage":
		return char.PhysCritDamage
	case "spell_crit_rate", "spellCritRate":
		return char.SpellCritRate
	case "spell_crit_damage", "spellCritDamage":
		return char.SpellCritDamage
	case "dodge_rate", "dodgeRate":
		return char.DodgeRate
	case "is_dead", "isDead", "IsDead":
		return char.IsDead
	case "id", "ID":
		return char.ID
	case "name", "Name":
		return char.Name
	case "class_id", "classId", "ClassID":
		return char.ClassID
	case "threat", "Threat":
		// 从战斗会话获取威胁值
		if ae.testContext != nil && ae.testContext.UserID > 0 {
			session := ae.testContext.BattleManager.GetSession(ae.testContext.UserID)
			if session != nil && len(session.ThreatTable) > 0 {
				// 查找该角色的威胁值
				for _, threatMap := range session.ThreatTable {
					if threat, exists := threatMap[char.ID]; exists {
						return threat
					}
				}
			}
		}
		return 0
	case "speed", "Speed":
		// 从Calculator计算速度
		if ae.testContext != nil && ae.testContext.Calculator != nil {
			return ae.testContext.Calculator.CalculateSpeed(char)
		}
		return char.Agility // 默认返回敏捷值
	case "skills", "Skills":
		// 返回技能列表（技能ID列表）
		if char.Skills != nil {
			skillIDs := make([]string, 0, len(char.Skills))
			for _, skill := range char.Skills {
				skillIDs = append(skillIDs, skill.SkillID)
			}
			return skillIDs
		}
		return []string{}
	default:
		// 处理skill_xxx格式的字段（如skill_charge）
		if strings.HasPrefix(fieldName, "skill_") {
			skillID := strings.TrimPrefix(fieldName, "skill_")
			// 查找对应的技能
			if char.Skills != nil {
				for _, skill := range char.Skills {
					if skill.SkillID == skillID {
						// 返回技能对象（用于后续访问.level等属性）
						return skill
					}
				}
			}
			return nil
		}
		// 尝试从上下文获取计算值
		if ae.context != nil {
			if val, exists := ae.context[fieldName]; exists {
				return val
			}
		}
		return nil
	}
}

// getMonsterField 获取怪物字段值
func (ae *AssertionExecutor) getMonsterField(monster *models.Monster, fieldName string) interface{} {
	switch fieldName {
	case "hp", "HP":
		return monster.HP
	case "max_hp", "maxHp", "MaxHP":
		return monster.MaxHP
	case "level", "Level":
		return monster.Level
	case "physical_attack", "physicalAttack", "PhysicalAttack":
		return monster.PhysicalAttack
	case "magic_attack", "magicAttack", "MagicAttack":
		return monster.MagicAttack
	case "physical_defense", "physicalDefense", "PhysicalDefense":
		return monster.PhysicalDefense
	case "magic_defense", "magicDefense", "MagicDefense":
		return monster.MagicDefense
	case "speed", "Speed":
		return monster.Speed
	case "exp_reward", "expReward":
		return monster.ExpReward
	case "id", "ID":
		return monster.ID
	case "name", "Name":
		return monster.Name
	case "type", "Type":
		return monster.Type
	case "debuff_defense_modifier":
		// 从Buff系统获取防御力Debuff
		if ae.testContext != nil && ae.testContext.UserID > 0 {
			session := ae.testContext.BattleManager.GetSession(ae.testContext.UserID)
			if session != nil {
				// 这里需要从BuffManager获取，暂时返回0
				// TODO: 实现从BuffManager获取Debuff值
			}
		}
		return 0.0
	case "actual_defense":
		// 实际防御力 = 基础防御力 * (1 + debuff_modifier)
		baseDefense := monster.PhysicalDefense
		debuffMod := 0.0
		if ae.testContext != nil && ae.testContext.UserID > 0 {
			// TODO: 从BuffManager获取debuff值
		}
		return int(float64(baseDefense) * (1.0 + debuffMod))
	default:
		return nil
	}
}

// assertEquals 断言相等（支持比较运算符如 "> 0", "< 100" 等）
func (ae *AssertionExecutor) assertEquals(actual interface{}, expected string) string {
	actualStr := fmt.Sprintf("%v", actual)
	
	// 如果期望值是纯字符串，直接比较
	if actualStr == expected {
		return "passed"
	}
	
	// 检查是否包含比较运算符
	expected = strings.TrimSpace(expected)
	if strings.HasPrefix(expected, ">") {
		// 大于比较
		expectedNum, err := strconv.ParseFloat(strings.TrimSpace(expected[1:]), 64)
		if err != nil {
			return "failed"
		}
		actualNum, err := ae.toNumber(actual)
		if err != nil {
			return "failed"
		}
		if actualNum > expectedNum {
			return "passed"
		}
		return "failed"
	} else if strings.HasPrefix(expected, "<") {
		// 小于比较
		expectedNum, err := strconv.ParseFloat(strings.TrimSpace(expected[1:]), 64)
		if err != nil {
			return "failed"
		}
		actualNum, err := ae.toNumber(actual)
		if err != nil {
			return "failed"
		}
		if actualNum < expectedNum {
			return "passed"
		}
		return "failed"
	} else if strings.HasPrefix(expected, ">=") {
		// 大于等于比较
		expectedNum, err := strconv.ParseFloat(strings.TrimSpace(expected[2:]), 64)
		if err != nil {
			return "failed"
		}
		actualNum, err := ae.toNumber(actual)
		if err != nil {
			return "failed"
		}
		if actualNum >= expectedNum {
			return "passed"
		}
		return "failed"
	} else if strings.HasPrefix(expected, "<=") {
		// 小于等于比较
		expectedNum, err := strconv.ParseFloat(strings.TrimSpace(expected[2:]), 64)
		if err != nil {
			return "failed"
		}
		actualNum, err := ae.toNumber(actual)
		if err != nil {
			return "failed"
		}
		if actualNum <= expectedNum {
			return "passed"
		}
		return "failed"
	} else if strings.HasPrefix(expected, "!=") {
		// 不等于比较
		expectedNum, err := strconv.ParseFloat(strings.TrimSpace(expected[2:]), 64)
		if err != nil {
			// 如果不是数字，按字符串比较
			if actualStr != strings.TrimSpace(expected[2:]) {
				return "passed"
			}
			return "failed"
		}
		actualNum, err := ae.toNumber(actual)
		if err != nil {
			return "failed"
		}
		if actualNum != expectedNum {
			return "passed"
		}
		return "failed"
	}
	
	// 尝试数字比较
	expectedNum, err := strconv.ParseFloat(expected, 64)
	if err == nil {
		actualNum, err := ae.toNumber(actual)
		if err == nil && actualNum == expectedNum {
			return "passed"
		}
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
	if ae.context == nil {
		ae.context = make(map[string]interface{})
	}
	ae.context[key] = value
}

// ClearContext 清空测试上下文
func (ae *AssertionExecutor) ClearContext() {
	ae.context = make(map[string]interface{})
}

// countAliveCharacters 计算存活角色数量
func (ae *AssertionExecutor) countAliveCharacters() int {
	if ae.testContext == nil {
		return 0
	}
	count := 0
	for _, char := range ae.testContext.Team {
		if char != nil && char.HP > 0 {
			count++
		}
	}
	return count
}

// calculateTeamTotalExp 计算队伍总经验值
func (ae *AssertionExecutor) calculateTeamTotalExp() int {
	if ae.testContext == nil {
		return 0
	}
	total := 0
	for _, char := range ae.testContext.Team {
		if char != nil {
			total += char.Exp
		}
	}
	return total
}

// countAliveMonsters 计算存活怪物数量
func (ae *AssertionExecutor) countAliveMonsters() int {
	if ae.testContext == nil {
		return 0
	}
	count := 0
	for _, monster := range ae.testContext.Monsters {
		if monster != nil && monster.HP > 0 {
			count++
		}
	}
	return count
}

// getBattleState 获取战斗状态
func (ae *AssertionExecutor) getBattleState() string {
	if ae.testContext == nil || ae.testContext.UserID == 0 {
		return "idle"
	}
	session := ae.testContext.BattleManager.GetSession(ae.testContext.UserID)
	if session == nil {
		return "idle"
	}
	if session.IsRunning {
		return "in_progress"
	}
	if session.IsResting {
		return "resting"
	}
	return "idle"
}

// getBattleResult 获取战斗结果
func (ae *AssertionExecutor) getBattleResult() map[string]interface{} {
	result := make(map[string]interface{})
	if ae.testContext == nil || ae.testContext.BattleResult == nil {
		result["is_victory"] = false
		return result
	}
	// 检查是否有存活的角色和怪物
	aliveChars := ae.countAliveCharacters()
	aliveMonsters := ae.countAliveMonsters()
	result["is_victory"] = aliveChars > 0 && aliveMonsters == 0
	return result
}

// findCharacterByClass 通过职业名称查找角色
func (ae *AssertionExecutor) findCharacterByClass(className string) *models.Character {
	if ae.testContext == nil {
		return nil
	}
	classMap := map[string]string{
		"warrior": "warrior",
		"priest":  "priest",
		"mage":    "mage",
		"rogue":   "rogue",
	}
	targetClassID := classMap[className]
	if targetClassID == "" {
		return nil
	}
	for _, char := range ae.testContext.Team {
		if char != nil && char.ClassID == targetClassID {
			return char
		}
	}
	return nil
}



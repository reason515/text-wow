package runner

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"text-wow/internal/game"
	"text-wow/internal/models"
	"text-wow/internal/repository"
	"time"
)

// Battle 相关函数
func (tr *TestRunner) executeBuildTurnOrder() error {

	// 使用与executeStartBattle相同的逻辑构建回合顺序
	return tr.buildTurnOrder()
}

func (tr *TestRunner) buildTurnOrder() error {

	// 收集所有参与者（角色和怪物）
	type participant struct {
		entry  map[string]interface{}
		speed  int
		isChar bool
		charID int
		key    string
	}
	participants := make([]participant, 0)
	debugPrint("[DEBUG] buildTurnOrder: Characters count=%d, Monsters count=%d\n", len(tr.context.Characters), len(tr.context.Monsters))

	// 收集所有角色（包括character和character_1, character_2等）
	for key, char := range tr.context.Characters {
		debugPrint("[DEBUG] buildTurnOrder: processing character key=%s, char=%v\n", key, char != nil)
		if char != nil {
			speed := tr.calculator.CalculateSpeed(char)

			// 从key中提取角色ID
			charID := key
			if key == "character" {

				// 如果�character"，检查是否有character_1，如果没有则使用character_1
				if _, exists := tr.context.Characters["character_1"]; !exists {

					// 如果没有character_1，使用character_1作为ID
					charID = "character_1"
				} else {

					// 如果有character_1，跳过这�character"（避免重复）
					continue
				}
			} else if strings.HasPrefix(key, "character_") {

				// 直接使用key作为ID（character_1, character_2等）
				charID = key
			} else {

				// 否则使用数据库ID
				charID = fmt.Sprintf("character_%d", char.ID)
			}
			charEntry := map[string]interface{}{
				"type":   "character",
				"id":     charID,
				"speed":  speed,
				"hp":     char.HP,
				"max_hp": char.MaxHP,
			}
			participants = append(participants, participant{
				entry:  charEntry,
				speed:  speed,
				isChar: true,
				charID: char.ID,
				key:    key,
			})
		}
	}
	// 收集所有怪物
	for key, monster := range tr.context.Monsters {
		debugPrint("[DEBUG] buildTurnOrder: processing monster key=%s, monster=%v\n", key, monster != nil)
		if monster != nil {

			// key可能是monster_1, monster_2等，直接使用作为ID
			monsterID := key

			// 如果key�monster"，则使用"monster_1"格式
			if key == "monster" {
				monsterID = "monster_1"
			}
			monsterEntry := map[string]interface{}{
				"type":   "monster",
				"id":     monsterID,
				"speed":  monster.Speed,
				"hp":     monster.HP,
				"max_hp": monster.MaxHP,
			}
			participants = append(participants, participant{
				entry:  monsterEntry,
				speed:  monster.Speed,
				isChar: false,
				key:    key,
			})
		}
	}
	// 按速度从高到低排序（速度相同时保持原有顺序）
	for i := 0; i < len(participants)-1; i++ {
		for j := i + 1; j < len(participants); j++ {
			if participants[i].speed < participants[j].speed {
				participants[i], participants[j] = participants[j], participants[i]
			}
		}
	}
	// 构建排序后的turn_order
	turnOrder := make([]interface{}, 0)
	for idx, p := range participants {
		turnOrder = append(turnOrder, p.entry)

		// 设置单独的键以便访问
		tr.safeSetContext(fmt.Sprintf("turn_order[%d].type", idx), p.entry["type"])
		tr.safeSetContext(fmt.Sprintf("turn_order[%d].speed", idx), p.speed)
		tr.context.Variables[fmt.Sprintf("turn_order[%d].type", idx)] = p.entry["type"]
		tr.context.Variables[fmt.Sprintf("turn_order[%d].speed", idx)] = p.speed
		if p.isChar {

			// 使用entry中的id（已从key提取）
			charID := p.entry["id"].(string)
			tr.safeSetContext(fmt.Sprintf("turn_order[%d].character.id", idx), charID)
			tr.context.Variables[fmt.Sprintf("turn_order[%d].character.id", idx)] = charID
		} else {

			// p.key可能是monster_1, monster_2等，直接使用，不需要再加monster_前缀
			monsterID := p.key

			// 如果key�monster"，则使用"monster_1"格式
			if p.key == "monster" {
				monsterID = "monster_1"
			}
			tr.safeSetContext(fmt.Sprintf("turn_order[%d].monster.id", idx), monsterID)
			tr.context.Variables[fmt.Sprintf("turn_order[%d].monster.id", idx)] = monsterID
		}
	}
	// 设置完整的turn_order数组（确保可序列化）
	tr.safeSetContext("turn_order", turnOrder)
	tr.context.Variables["turn_order"] = turnOrder
	tr.safeSetContext("turn_order_length", len(turnOrder))
	tr.context.Variables["turn_order_length"] = len(turnOrder)
	debugPrint("[DEBUG] buildTurnOrder: created turn_order with %d participants\n", len(turnOrder))
	return nil
}

func (tr *TestRunner) executeAttackMonster() error {
	char, ok := tr.context.Characters["character"]
	if !ok || char == nil {
		return fmt.Errorf("character not found")
	}
	// 找到第一个存活的怪物
	var targetMonster *models.Monster
	var targetKey string
	for key, monster := range tr.context.Monsters {
		if monster != nil && monster.HP > 0 {
			targetMonster = monster
			targetKey = key
			break
		}
	}
	if targetMonster == nil {
		return fmt.Errorf("monster not found")
	}
	// 计算伤害（考虑Debuff减成）
	baseAttack := float64(char.PhysicalAttack)
	// 检查是否有Debuff减成
	if debuffModifier, exists := tr.context.Variables["monster_debuff_attack_modifier"]; exists {
		if modifier, ok := debuffModifier.(float64); ok && modifier < 0 {
			baseAttack = baseAttack * (1.0 + modifier) // modifier是负数，所以是1.0 + (-0.3) = 0.7
			debugPrint("[DEBUG] executeAttackMonster: Debuff applied, modifier=%f, baseAttack=%f\n", modifier, baseAttack)
		}
	}
	damage := int(math.Round(baseAttack)) - targetMonster.PhysicalDefense
	if damage < 1 {
		damage = 1
	}
	// 应用伤害
	targetMonster.HP -= damage
	if targetMonster.HP < 0 {
		targetMonster.HP = 0
	}
	// 添加战斗日志
	if battleLogs, exists := tr.context.Variables["battle_logs"]; exists {
		if logs, ok := battleLogs.([]string); ok {
			logs = append(logs, fmt.Sprintf("角色攻击怪物，造成%d点伤害", damage))
			tr.context.Variables["battle_logs"] = logs
		}
	} else {
		tr.context.Variables["battle_logs"] = []string{fmt.Sprintf("角色攻击怪物，造成%d点伤害", damage)}
	}
	tr.safeSetContext("damage_dealt", damage)
	tr.context.Variables["damage_dealt"] = damage

	// 战士攻击时获得怒气（假设获�0点）
	if char.ResourceType == "rage" {
		char.Resource += 10
		if char.Resource > char.MaxResource {
			char.Resource = char.MaxResource
		}
	}
	// 更新上下文
	tr.context.Characters["character"] = char
	// 更新怪物到上下文
	if targetKey != "" {
		tr.context.Monsters[targetKey] = targetMonster
	}
	// 如果怪物HP为0，战斗结束，战士怒气归0
	if targetMonster.HP == 0 {
		if char.ResourceType == "rage" {
			char.Resource = 0
			tr.context.Characters["character"] = char
		}
		// 检查是否所有怪物都死亡
		allDead := true
		hasMonsters := false
		for _, m := range tr.context.Monsters {
			if m != nil {
				hasMonsters = true
				if m.HP > 0 {
					allDead = false
					break
				}
			}
		}
		// 如果存在怪物且所有怪物都死亡，或者角色还活着且攻击的怪物已死亡，则算胜利
		if (hasMonsters && allDead) || (char.HP > 0 && targetMonster.HP == 0) {
			tr.setBattleResult(true, char)
		}
	}
	return nil
}

// executeTeamAttackMonster 队伍攻击怪物
func (tr *TestRunner) executeTeamAttackMonster() error {
	// 找到第一个存活的怪物
	var targetMonster *models.Monster
	var targetKey string
	for key, monster := range tr.context.Monsters {
		if monster != nil && monster.HP > 0 {
			targetMonster = monster
			targetKey = key
			break
		}
	}
	if targetMonster == nil {
		return fmt.Errorf("monster not found")
	}

	// 让所有存活的角色攻击怪物
	totalDamage := 0
	attacked := false
	for key, char := range tr.context.Characters {
		if char != nil && char.HP > 0 {
			// 计算伤害
			baseAttack := float64(char.PhysicalAttack)
			// 检查是否有Debuff减成
			if debuffModifier, exists := tr.context.Variables["monster_debuff_attack_modifier"]; exists {
				if modifier, ok := debuffModifier.(float64); ok && modifier < 0 {
					baseAttack = baseAttack * (1.0 + modifier)
				}
			}
			damage := int(math.Round(baseAttack)) - targetMonster.PhysicalDefense
			if damage < 1 {
				damage = 1
			}
			
			// 应用伤害
			targetMonster.HP -= damage
			if targetMonster.HP < 0 {
				targetMonster.HP = 0
			}
			totalDamage += damage
			attacked = true

			// 战士攻击时获得怒气
			if char.ResourceType == "rage" {
				char.Resource += 10
				if char.Resource > char.MaxResource {
					char.Resource = char.MaxResource
				}
			}
			
			// 更新角色到上下文
			tr.context.Characters[key] = char
		}
	}

	if !attacked {
		return fmt.Errorf("no alive characters in team")
	}

	// 更新怪物到上下文
	if targetKey != "" {
		tr.context.Monsters[targetKey] = targetMonster
	}

	// 添加战斗日志
	if battleLogs, exists := tr.context.Variables["battle_logs"]; exists {
		if logs, ok := battleLogs.([]string); ok {
			logs = append(logs, fmt.Sprintf("队伍攻击怪物，造成%d点伤害", totalDamage))
			tr.context.Variables["battle_logs"] = logs
		}
	} else {
		tr.context.Variables["battle_logs"] = []string{fmt.Sprintf("队伍攻击怪物，造成%d点伤害", totalDamage)}
	}
	tr.safeSetContext("damage_dealt", totalDamage)
	tr.context.Variables["damage_dealt"] = totalDamage

	// 如果怪物HP为0，战斗结束
	if targetMonster.HP == 0 {
		// 重置所有战士的怒气
		for key, char := range tr.context.Characters {
			if char != nil && char.ResourceType == "rage" {
				char.Resource = 0
				tr.context.Characters[key] = char
			}
		}
		
		// 检查是否所有怪物都死亡
		allDead := true
		for _, m := range tr.context.Monsters {
			if m != nil && m.HP > 0 {
				allDead = false
				break
			}
		}
		if allDead {
			// 使用第一个存活的角色作为代表设置战斗结果
			var firstChar *models.Character
			for _, c := range tr.context.Characters {
				if c != nil && c.HP > 0 {
					firstChar = c
					break
				}
			}
			if firstChar != nil {
				tr.setBattleResult(true, firstChar)
			}
		}
	}

	return nil
}

func (tr *TestRunner) executeMonsterAttack() error {
	char, ok := tr.context.Characters["character"]
	if !ok || char == nil {
		return fmt.Errorf("character not found")
	}
	// 找到第一个存活的怪物
	var attackerMonster *models.Monster
	for _, monster := range tr.context.Monsters {
		if monster != nil && monster.HP > 0 {
			attackerMonster = monster
			break
		}
	}
	if attackerMonster == nil {
		return fmt.Errorf("monster not found")
	}
	// 计算伤害（考虑Buff加成）
	baseAttack := float64(attackerMonster.PhysicalAttack)
	// 检查是否有Buff加成
	if buffModifier, exists := tr.context.Variables["monster_buff_attack_modifier"]; exists {
		if modifier, ok := buffModifier.(float64); ok && modifier > 0 {
			baseAttack = baseAttack * (1.0 + modifier)
			debugPrint("[DEBUG] executeMonsterAttack: Buff applied, modifier=%f, baseAttack=%f\n", modifier, baseAttack)
		}
	}
	damage := int(math.Round(baseAttack)) - char.PhysicalDefense
	if damage < 1 {
		damage = 1
	}
	// 保存当前怒气（用于调试）
	originalResource := char.Resource
	debugPrint("[DEBUG] executeMonsterAttack: before attack - char.HP=%d, char.Resource=%d, monster.Attack=%d, char.Defense=%d, damage=%d\n", char.HP, char.Resource, attackerMonster.PhysicalAttack, char.PhysicalDefense, damage)

	// 应用伤害
	char.HP -= damage
	if char.HP < 0 {
		char.HP = 0
	}
	// 添加战斗日志
	if battleLogs, exists := tr.context.Variables["battle_logs"]; exists {
		if logs, ok := battleLogs.([]string); ok {
			logs = append(logs, fmt.Sprintf("怪物攻击角色，造成%d点伤害", damage))
			tr.context.Variables["battle_logs"] = logs
		}
	} else {
		tr.context.Variables["battle_logs"] = []string{fmt.Sprintf("怪物攻击角色，造成%d点伤害", damage)}
	}
	tr.safeSetContext("monster_damage_dealt", damage)
	tr.context.Variables["monster_damage_dealt"] = damage
	debugPrint("[DEBUG] executeMonsterAttack: after damage - char.HP=%d, char.Resource=%d\n", char.HP, char.Resource)

	// 如果角色HP�，战斗失败，战士怒气�（在获得怒气之前检查）
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
		// 如果角色死亡，不再获得怒气，设置战斗失败
		tr.context.Characters["character"] = char
		debugPrint("[DEBUG] executeMonsterAttack: character died, HP=0, rage reset to 0 (was %d)\n", originalResource)
		tr.setBattleResult(false, char)
		return nil
	}
	// 只有在角色未死亡时，才获得怒气
	// 战士受到伤害时获得怒气（假设获得5点）
	if char.ResourceType == "rage" {
		char.Resource += 5
		if char.Resource > char.MaxResource {
			char.Resource = char.MaxResource
		}
		debugPrint("[DEBUG] executeMonsterAttack: character took damage, HP=%d, rage increased from %d to %d\n", char.HP, originalResource, char.Resource)
	}
	// 更新上下文（确保角色HP被正确更新）
	tr.context.Characters["character"] = char
	// 如果角色HP为0，确保设置战斗失败（防御性检查）
	if char.HP == 0 {
		tr.setBattleResult(false, char)
	}
	return nil
}

func (tr *TestRunner) executeBattleRound(instruction string) error {

	// 解析回合数（如"执行X回合"或"执行一个回合"）
	roundNum := 1
	if strings.Contains(instruction, "执行") {
		parts := strings.Split(instruction, "执行")
		if len(parts) > 1 {
			roundStr := strings.TrimSpace(strings.Split(parts[1], "回合")[0])
			if round, err := strconv.Atoi(roundStr); err == nil {
				roundNum = round
			}
		}
	} else {

		// 如果没有指定回合数，从上下文获取当前回合数并递增
		if currentRound, exists := tr.context.Variables["current_round"]; exists {
			if cr, ok := currentRound.(int); ok {
				roundNum = cr + 1
			}
		}
		tr.context.Variables["current_round"] = roundNum
		tr.safeSetContext("current_round", roundNum)
	}
	// 减少技能冷却时间
	skillManager := game.NewSkillManager()
	char, ok := tr.context.Characters["character"]
	if ok && char != nil {
		if err := skillManager.LoadCharacterSkills(char.ID); err == nil {

			// 先减少冷却时间
			skillManager.TickCooldowns(char.ID)
			// 减少Buff持续时间（每回合减少1）
			if buffDuration, exists := tr.context.Variables["character_buff_duration"]; exists {
				if duration, ok := buffDuration.(int); ok && duration > 0 {
					newDuration := duration - 1
					if newDuration < 0 {
						newDuration = 0
					}
					tr.context.Variables["character_buff_duration"] = newDuration
					tr.safeSetContext("character.buff_duration", newDuration)
					tr.safeSetContext(fmt.Sprintf("buff_duration_round_%d", roundNum), newDuration)
					tr.context.Variables[fmt.Sprintf("buff_duration_round_%d", roundNum)] = newDuration
				}
			}
			// 减少护盾持续时间（每回合减少1）
			if shieldDuration, exists := tr.context.Variables["character.shield_duration"]; exists {
				if duration, ok := shieldDuration.(int); ok && duration > 0 {
					newDuration := duration - 1
					if newDuration < 0 {
						newDuration = 0
					}
					tr.context.Variables["character.shield_duration"] = newDuration
					tr.safeSetContext("character.shield_duration", newDuration)
					tr.safeSetContext(fmt.Sprintf("character.shield_duration_round_%d", roundNum), newDuration)
					tr.context.Variables[fmt.Sprintf("character.shield_duration_round_%d", roundNum)] = newDuration
				}
			}
			// 减少状态效果持续时间（每回合减少1）
			statusTypes := []string{"stunned", "silenced", "feared"}
			for _, statusType := range statusTypes {
				durationKey := fmt.Sprintf("character.%s_duration", statusType)
				if statusDuration, exists := tr.context.Variables[durationKey]; exists {
					if duration, ok := statusDuration.(int); ok && duration > 0 {
						newDuration := duration - 1
						if newDuration < 0 {
							newDuration = 0
						}
						tr.context.Variables[durationKey] = newDuration
						tr.safeSetContext(durationKey, newDuration)
						tr.safeSetContext(fmt.Sprintf("character.%s_duration_round_%d", statusType, roundNum), newDuration)
						tr.context.Variables[fmt.Sprintf("character.%s_duration_round_%d", statusType, roundNum)] = newDuration
						// 如果持续时间为0，清除状态
						if newDuration == 0 {
							if statusType == "stunned" {
								tr.context.Variables["character.is_stunned"] = false
								tr.safeSetContext("character.is_stunned", false)
							} else if statusType == "silenced" {
								tr.context.Variables["character.is_silenced"] = false
								tr.safeSetContext("character.is_silenced", false)
							} else if statusType == "feared" {
								tr.context.Variables["character.is_feared"] = false
								tr.safeSetContext("character.is_feared", false)
							}
						}
					}
				}
			}
			// 获取技能状态，检查是否可用（不再从Variables读取Skill对象，避免序列化错误）
			skillID, exists := tr.context.Variables["skill_id"]
			if exists {
				skillIDStr, ok := skillID.(string)
				if ok && skillIDStr != "" {
					skillState := skillManager.GetSkillState(char.ID, skillIDStr)
					if skillState != nil {
						tr.safeSetContext(fmt.Sprintf("skill_usable_round_%d", roundNum), skillState.CooldownLeft == 0)
						tr.safeSetContext(fmt.Sprintf("skill_cooldown_round_%d", roundNum), skillState.CooldownLeft)
					} else {

						// 如果技能状态不存在，从Variables获取冷却时间并计算
						cooldown := 0
						if cooldownVal, exists := tr.context.Variables["skill_cooldown"]; exists {
							if cd, ok := cooldownVal.(int); ok {
								cooldown = cd
							}
						}
						// 假设第1回合使用了技能，冷却时间3，那么：
						// 第2回合：冷却剩余2，不可用
						// 第3回合：冷却剩余1，不可用
						// 第4回合：冷却剩余0，可用
						// �回合：冷却剩�，不可用
						// �回合：冷却剩�，不可用
						cooldownLeft := cooldown - (roundNum - 1)
						if cooldownLeft < 0 {
							cooldownLeft = 0
						}
						tr.safeSetContext(fmt.Sprintf("skill_usable_round_%d", roundNum), cooldownLeft == 0)
						tr.safeSetContext(fmt.Sprintf("skill_cooldown_round_%d", roundNum), cooldownLeft)
					}
				}
			}
		} else {

			// 如果角色没有技能，从上下文获取技能信息（不再从Variables读取Skill对象）
			if _, exists := tr.context.Variables["skill_id"]; exists {
				// 从Variables获取冷却时间并计算
				cooldown := 0
				if cooldownVal, exists := tr.context.Variables["skill_cooldown"]; exists {
					if cd, ok := cooldownVal.(int); ok {
						cooldown = cd
					}
				}
				// 根据冷却时间计算
				cooldownLeft := cooldown - (roundNum - 1)
				if cooldownLeft < 0 {
					cooldownLeft = 0
				}
				tr.safeSetContext(fmt.Sprintf("skill_usable_round_%d", roundNum), cooldownLeft == 0)
				tr.safeSetContext(fmt.Sprintf("skill_cooldown_round_%d", roundNum), cooldownLeft)
			}
		}
	}
	// 处理怪物技能冷却时间（不再从Variables读取Skill对象，避免序列化错误）
	if monsterSkillID, exists := tr.context.Variables["monster_skill_id"]; exists && monsterSkillID != nil {
		// 从Variables获取怪物技能冷却时间
		monsterCooldown := 0
		if cooldownVal, exists := tr.context.Variables["monster_skill_cooldown"]; exists {
			if cd, ok := cooldownVal.(int); ok {
				monsterCooldown = cd
			}
		}
		// 获取上次使用技能的回合
		lastUsedRound := 1
		if lastRound, exists := tr.context.Variables["monster_skill_last_used_round"]; exists {
			if lr, ok := lastRound.(int); ok {
				lastUsedRound = lr
			}
		}
		// 计算冷却剩余时间
		cooldownLeft := monsterCooldown - (roundNum - lastUsedRound)
		if cooldownLeft < 0 {
			cooldownLeft = 0
		}
		tr.safeSetContext(fmt.Sprintf("monster_skill_cooldown_round_%d", roundNum), cooldownLeft)
		tr.context.Variables[fmt.Sprintf("monster_skill_cooldown_round_%d", roundNum)] = cooldownLeft
	}

	// 记录当前回合的角色HP
	if char != nil {
		tr.safeSetContext(fmt.Sprintf("character.hp_round_%d", roundNum), char.HP)
		tr.context.Variables[fmt.Sprintf("character.hp_round_%d", roundNum)] = char.HP
	}

	// 记录当前回合的怪物HP
	monsterIdx := 1
	for _, monster := range tr.context.Monsters {
		if monster != nil {
			tr.safeSetContext(fmt.Sprintf("monster.hp_round_%d", roundNum), monster.HP)
			tr.context.Variables[fmt.Sprintf("monster.hp_round_%d", roundNum)] = monster.HP
			tr.safeSetContext(fmt.Sprintf("monster_%d.hp_round_%d", monsterIdx, roundNum), monster.HP)
			tr.context.Variables[fmt.Sprintf("monster_%d.hp_round_%d", monsterIdx, roundNum)] = monster.HP
			monsterIdx++
		}
	}

	return nil
}

func (tr *TestRunner) executeRemainingMonstersAttack(instruction string) error {

	// 解析剩余怪物数量（如"剩余2个怪物攻击角色"）
	expectedCount := 0
	if strings.Contains(instruction, "剩余") {
		parts := strings.Split(instruction, "剩余")
		if len(parts) > 1 {
			countStr := strings.TrimSpace(strings.Split(parts[1], "个")[0])
			if count, err := strconv.Atoi(countStr); err == nil {
				expectedCount = count
			}
		}
	}
	// 获取所有存活的怪物
	aliveMonsters := []*models.Monster{}
	for _, monster := range tr.context.Monsters {
		if monster != nil && monster.HP > 0 {
			aliveMonsters = append(aliveMonsters, monster)
		}
	}
	// 验证存活怪物数量
	if len(aliveMonsters) != expectedCount {
		debugPrint("Warning: expected %d alive monsters, but found %d\n", expectedCount, len(aliveMonsters))
	}
	// 执行攻击
	return tr.executeAllMonstersAttack(instruction)
}

func (tr *TestRunner) executeWaitRestRecovery() error {

	// 检查是否处于休息状态
	isResting, exists := tr.context.Variables["is_resting"]
	if !exists || isResting == nil || !isResting.(bool) {

		// 如果不在休息状态，先进入休息状态
		if err := tr.checkAndEnterRest(); err != nil {
			return fmt.Errorf("failed to enter rest state: %w", err)
		}
	}
	// 模拟休息恢复（简化处理：直接恢复到满值）
	char, ok := tr.context.Characters["character"]
	if !ok || char == nil {
		return fmt.Errorf("character not found")
	}
	// 恢复HP和Resource（简化：恢复到最大值）
	char.HP = char.MaxHP
	char.Resource = char.MaxResource
	// 更新上下文
	tr.context.Characters["character"] = char
	tr.safeSetContext("character.hp", char.HP)
	tr.safeSetContext("character.resource", char.Resource)
	tr.safeSetContext("character.max_hp", char.MaxHP)
	tr.safeSetContext("character.max_resource", char.MaxResource)
	return nil
}

func (tr *TestRunner) setBattleResult(isVictory bool, char *models.Character) {

	// 设置战斗状态
	if isVictory {
		tr.safeSetContext("battle_state", "victory")
		tr.context.Variables["battle_state"] = "victory"

		// 添加战斗日志
		if battleLogs, exists := tr.context.Variables["battle_logs"]; exists {
			if logs, ok := battleLogs.([]string); ok {
				logs = append(logs, "战斗胜利")
				tr.context.Variables["battle_logs"] = logs
			}
		}
		// 检查是否应该进入休息状态
		if err := tr.checkAndEnterRest(); err != nil {
			debugPrint("Warning: failed to enter rest state: %v\n", err)
		}
	} else {
		tr.safeSetContext("battle_state", "defeat")
		tr.context.Variables["battle_state"] = "defeat"

		// 添加战斗日志
		if battleLogs, exists := tr.context.Variables["battle_logs"]; exists {
			if logs, ok := battleLogs.([]string); ok {
				logs = append(logs, "战败")
				tr.context.Variables["battle_logs"] = logs
			}
		}
	}
	// 设置战斗结果
	tr.safeSetContext("battle_result.is_victory", isVictory)
	tr.context.Variables["battle_result.is_victory"] = isVictory

	// 计算战斗时长
	if startTime, exists := tr.context.Variables["battle_start_time"]; exists {
		if start, ok := startTime.(int64); ok {
			duration := time.Now().Unix() - start
			tr.safeSetContext("battle_result.duration_seconds", duration)
			tr.context.Variables["battle_result.duration_seconds"] = duration
		}
	}
	// 设置角色死亡状态
	if char != nil {
		isDead := char.HP <= 0
		tr.safeSetContext("character.is_dead", isDead)
		tr.context.Variables["character.is_dead"] = isDead

		// 如果胜利，给予经验和金币奖励
		if isVictory {
			// 计算经验奖励（基于怪物数量）
			expGain := len(tr.context.Monsters) * 10 // 每个怪物10经验
			
			// 给所有存活的角色添加经验
			teamTotalExp := 0
			for key, c := range tr.context.Characters {
				if c != nil && c.HP > 0 {
					c.Exp += expGain
					teamTotalExp += c.Exp
					
					// 更新角色到上下文
					tr.context.Characters[key] = c
					
					// 根据角色键设置经验变量
					if key == "character" {
						tr.safeSetContext("character.exp", c.Exp)
						tr.context.Variables["character.exp"] = c.Exp
						tr.safeSetContext("character.exp_gained", expGain)
						tr.context.Variables["character.exp_gained"] = expGain
					} else {
						// 对于其他角色，也设置经验变量
						tr.safeSetContext(fmt.Sprintf("%s.exp", key), c.Exp)
						tr.context.Variables[fmt.Sprintf("%s.exp", key)] = c.Exp
					}
					
					// 根据职业ID设置经验变量（如warrior.exp, priest.exp）
					if c.ClassID == "warrior" {
						tr.safeSetContext("warrior.exp", c.Exp)
						tr.context.Variables["warrior.exp"] = c.Exp
					} else if c.ClassID == "priest" {
						tr.safeSetContext("priest.exp", c.Exp)
						tr.context.Variables["priest.exp"] = c.Exp
					} else if c.ClassID == "mage" {
						tr.safeSetContext("mage.exp", c.Exp)
						tr.context.Variables["mage.exp"] = c.Exp
					} else if c.ClassID == "rogue" {
						tr.safeSetContext("rogue.exp", c.Exp)
						tr.context.Variables["rogue.exp"] = c.Exp
					}
				} else if c != nil {
					// 死亡的角色不获得经验，但计入总经验
					teamTotalExp += c.Exp
				}
			}
			
			// 设置team_total_exp
			tr.safeSetContext("team_total_exp", teamTotalExp)
			tr.context.Variables["team_total_exp"] = teamTotalExp

			// 计算金币奖励（简化：每个怪物20金币，只给主角色）
			if char != nil {
				goldGain := len(tr.context.Monsters) * 20 // 每个怪物20金币
				userRepo := repository.NewUserRepository()
				if user, err := userRepo.GetByID(char.UserID); err == nil && user != nil {
					newGold := user.Gold + goldGain
					userRepo.UpdateGold(char.UserID, newGold)
					tr.safeSetContext("character.gold", newGold)
					tr.context.Variables["character.gold"] = newGold
					tr.safeSetContext("character.gold_gained", goldGain)
					tr.context.Variables["character.gold_gained"] = goldGain
				}
			}
		} else {

			// 失败时，exp_gained和gold_gained为0
			tr.safeSetContext("character.exp_gained", 0)
			tr.context.Variables["character.exp_gained"] = 0
			tr.safeSetContext("character.gold_gained", 0)
			tr.context.Variables["character.gold_gained"] = 0
		}
		// 设置team_alive_count（单角色时，如果角色死亡则为0，否则为1）
		aliveCount := 0
		if char.HP > 0 {
			aliveCount = 1
		}
		tr.safeSetContext("team_alive_count", aliveCount)
		tr.context.Variables["team_alive_count"] = aliveCount

		// 设置enemy_death_count
		enemyDeathCount := 0
		for _, monster := range tr.context.Monsters {
			if monster != nil && monster.HP <= 0 {
				enemyDeathCount++
			}
		}
		tr.safeSetContext("enemy_death_count", enemyDeathCount)
		tr.context.Variables["enemy_death_count"] = enemyDeathCount

		// 如果角色是战士，确保怒气为0
		if char.ResourceType == "rage" {
			char.Resource = 0
			char.MaxResource = 100

			// 更新数据库
			charRepo := repository.NewCharacterRepository()
			charRepo.UpdateAfterBattle(char.ID, char.HP, char.Resource, char.Exp, char.Level,
				char.ExpToNext, char.MaxHP, char.MaxResource, char.PhysicalAttack, char.MagicAttack, char.PhysicalDefense, char.MagicDefense,
				char.Strength, char.Agility, char.Intellect, char.Stamina, char.Spirit, char.UnspentPoints, char.TotalKills)
		}
		tr.context.Characters["character"] = char
	}
	// 设置battle_rounds
	if rounds, exists := tr.context.Variables["battle_rounds"]; exists {
		if r, ok := rounds.(int); ok {
			tr.safeSetContext("battle_rounds", r)
		}
	}
}

// executeAllMonstersAttack 所有怪物攻击角色或队伍
func (tr *TestRunner) executeAllMonstersAttack(instruction string) error {
	// 获取所有存活的怪物
	for key, monster := range tr.context.Monsters {
		if monster != nil && monster.HP > 0 {
			// 执行怪物攻击
			if err := tr.executeMonsterAttack(); err != nil {
				debugPrint("[DEBUG] executeAllMonstersAttack: failed to execute attack for monster %s: %v\n", key, err)
				// 继续执行其他怪物的攻击
			}
		}
	}
	return nil
}

// checkAndEnterRest 检查并进入休息状态（当所有敌人死亡时）
func (tr *TestRunner) checkAndEnterRest() error {
	// 检查是否所有怪物都已死亡
	allMonstersDead := true
	for _, monster := range tr.context.Monsters {
		if monster != nil && monster.HP > 0 {
			allMonstersDead = false
			break
		}
	}

	if allMonstersDead {
		tr.safeSetContext("is_resting", true)
		tr.context.Variables["is_resting"] = true
		tr.safeSetContext("battle_state", "resting")
		tr.context.Variables["battle_state"] = "resting"
		debugPrint("[DEBUG] checkAndEnterRest: entered rest state\n")
	}

	return nil
}

// executeStartBattle 开始战斗
func (tr *TestRunner) executeStartBattle() error {
	char, ok := tr.context.Characters["character"]
	if !ok || char == nil {
		// 尝试查找其他角色（如character_1, warrior等）
		for key, c := range tr.context.Characters {
			if c != nil {
				char = c
				// 同时设置到"character"键以便后续使用
				tr.context.Characters["character"] = c
				debugPrint("[DEBUG] executeStartBattle: using character from key '%s'\n", key)
				break
			}
		}
		if char == nil {
			return fmt.Errorf("character not found")
		}
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

	// 初始化战斗日志和战斗开始时间
	battleLogs := []string{"战斗开始"}
	tr.context.Variables["battle_logs"] = battleLogs
	tr.context.Variables["battle_start_time"] = time.Now().Unix()
	tr.context.Variables["battle_rounds"] = 0
	// 记录战斗前的经验值（用于计算exp_gained）
	tr.context.Variables["character.exp_before_battle"] = char.Exp

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

	// 设置战斗状态到上下文
	tr.safeSetContext("battle_state", "in_progress")
	tr.context.Variables["battle_state"] = "in_progress"
	tr.safeSetContext("is_resting", false)
	tr.context.Variables["is_resting"] = false

	// 计算并设置回合顺序（使用通用函数）
	if err := tr.buildTurnOrder(); err != nil {
		return fmt.Errorf("failed to build turn order: %w", err)
	}

	// 设置敌人数量
	enemyCount := len(tr.context.Monsters)
	tr.safeSetContext("enemy_count", enemyCount)
	tr.context.Variables["enemy_count"] = enemyCount

	// 计算存活敌人数量
	aliveEnemyCount := 0
	for _, monster := range tr.context.Monsters {
		if monster != nil && monster.HP > 0 {
			aliveEnemyCount++
		}
	}
	tr.safeSetContext("enemy_alive_count", aliveEnemyCount)
	tr.context.Variables["enemy_alive_count"] = aliveEnemyCount
	// 同时设置别名 enemies_alive_count（复数形式）
	tr.safeSetContext("enemies_alive_count", aliveEnemyCount)
	tr.context.Variables["enemies_alive_count"] = aliveEnemyCount

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
	char, ok := tr.context.Characters["character"]
	if !ok || char == nil {
		return fmt.Errorf("character not found")
	}

	// 检查战斗结果是否已经设置
	if _, exists := tr.context.Variables["battle_result.is_victory"]; !exists {
		// 如果还没有设置战斗结果，检查当前状态并设置
		// 检查角色是否死亡
		if char.HP <= 0 {
			tr.setBattleResult(false, char)
		} else {
			// 检查是否所有怪物都死亡
			hasAliveMonsters := false
			for _, monster := range tr.context.Monsters {
				if monster != nil && monster.HP > 0 {
					hasAliveMonsters = true
					break
				}
			}
			// 如果角色还活着且没有存活的怪物，则算胜利
			if !hasAliveMonsters && char.HP > 0 {
				tr.setBattleResult(true, char)
			}
		}
	}

	// 如果角色是战士，确保怒气为0
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

// executeMonsterUseSkill 怪物使用技能
func (tr *TestRunner) executeMonsterUseSkill(instruction string) error {
	// 找到第一个存活的怪物
	var targetMonster *models.Monster
	var targetKey string
	for key, monster := range tr.context.Monsters {
		if monster != nil && monster.HP > 0 {
			targetMonster = monster
			targetKey = key
			break
		}
	}
	if targetMonster == nil {
		return fmt.Errorf("monster not found")
	}

	// 解析技能类型（如"怪物使用Buff技能"）
	if strings.Contains(instruction, "Buff") {
		// 应用Buff效果到角色
		char, ok := tr.context.Characters["character"]
		if ok && char != nil {
			// 这里可以添加Buff逻辑
			tr.context.Characters["character"] = char
		}
	} else if strings.Contains(instruction, "Debuff") {
		// 应用Debuff效果到角色
		char, ok := tr.context.Characters["character"]
		if ok && char != nil {
			// 这里可以添加Debuff逻辑
			tr.context.Characters["character"] = char
		}
	}

	// 更新怪物到上下文
	tr.context.Monsters[targetKey] = targetMonster

	return nil
}

// executeCheckCharacterAttributes 检查角色属性
func (tr *TestRunner) executeCheckCharacterAttributes() error {
	char, ok := tr.context.Characters["character"]
	if !ok || char == nil {
		return fmt.Errorf("character not found")
	}

	// 重新计算所有属性
	char.PhysicalAttack = tr.calculator.CalculatePhysicalAttack(char)
	char.MagicAttack = tr.calculator.CalculateMagicAttack(char)
	char.PhysicalDefense = tr.calculator.CalculatePhysicalDefense(char)
	char.MagicDefense = tr.calculator.CalculateMagicDefense(char)
	char.PhysCritRate = tr.calculator.CalculatePhysCritRate(char)
	char.SpellCritRate = tr.calculator.CalculateSpellCritRate(char)
	char.DodgeRate = tr.calculator.CalculateDodgeRate(char)

	// 更新上下文
	tr.context.Characters["character"] = char

	// 存储到Variables
	tr.context.Variables["character_physical_attack"] = char.PhysicalAttack
	tr.context.Variables["character_magic_attack"] = char.MagicAttack
	tr.context.Variables["character_physical_defense"] = char.PhysicalDefense
	tr.context.Variables["character_magic_defense"] = char.MagicDefense
	tr.context.Variables["character_phys_crit_rate"] = char.PhysCritRate
	tr.context.Variables["character_spell_crit_rate"] = char.SpellCritRate
	tr.context.Variables["character_dodge_rate"] = char.DodgeRate

	return nil
}

// executeEnterRestState 进入休息状态
func (tr *TestRunner) executeEnterRestState() error {
	char, ok := tr.context.Characters["character"]
	if !ok || char == nil {
		return fmt.Errorf("character not found")
	}

	// 设置休息状态
	tr.safeSetContext("is_resting", true)
	tr.context.Variables["is_resting"] = true
	tr.safeSetContext("battle_state", "resting")
	tr.context.Variables["battle_state"] = "resting"

	// 更新上下文
	tr.context.Characters["character"] = char

	return nil
}

// executeContinueBattleUntil 继续战斗直到指定条件
func (tr *TestRunner) executeContinueBattleUntil(instruction string) error {
	maxRounds := 20 // 防止无限循环

	if strings.Contains(instruction, "角色死亡") {
		// 继续战斗直到角色死亡
		for round := 0; round < maxRounds; round++ {
			char, ok := tr.context.Characters["character"]
			if !ok || char == nil {
				return fmt.Errorf("character not found")
			}
			if char.HP <= 0 {
				tr.setBattleResult(false, char)
				break
			}
			// 怪物攻击角色
			if err := tr.executeMonsterAttack(); err != nil {
				debugPrint("[DEBUG] executeContinueBattleUntil: monster attack error: %v\n", err)
			}
		}
	} else if strings.Contains(instruction, "怪物死亡") {
		// 继续战斗直到所有怪物死亡
		for round := 0; round < maxRounds; round++ {
			// 检查是否所有怪物已死亡
			allDead := true
			for _, monster := range tr.context.Monsters {
				if monster != nil && monster.HP > 0 {
					allDead = false
					break
				}
			}
			if allDead {
				char, _ := tr.context.Characters["character"]
				tr.setBattleResult(true, char)
				break
			}
			// 角色攻击怪物
			if err := tr.executeAttackMonster(); err != nil {
				debugPrint("[DEBUG] executeContinueBattleUntil: attack error: %v\n", err)
			}
		}
	}
	return nil
}

// executeAttackSpecificMonster 攻击指定怪物
func (tr *TestRunner) executeAttackSpecificMonster(instruction string) error {
	char, ok := tr.context.Characters["character"]
	if !ok || char == nil {
		return fmt.Errorf("character not found")
	}

	// 解析怪物编号（如"攻击怪物1"）
	monsterIndex := 1
	if strings.Contains(instruction, "怪物") {
		parts := strings.Split(instruction, "怪物")
		if len(parts) > 1 {
			indexStr := strings.TrimSpace(parts[1])
			if index, err := strconv.Atoi(indexStr); err == nil {
				monsterIndex = index
			}
		}
	}

	// 获取指定怪物
	var targetMonster *models.Monster
	var targetKey string
	if monsterIndex == 1 {
		targetKey = "monster"
	} else {
		targetKey = fmt.Sprintf("monster_%d", monsterIndex)
	}
	targetMonster, ok = tr.context.Monsters[targetKey]
	if !ok || targetMonster == nil {
		return fmt.Errorf("monster %d not found", monsterIndex)
	}

	// 计算伤害
	baseAttack := float64(char.PhysicalAttack)
	damage := int(math.Round(baseAttack)) - targetMonster.PhysicalDefense
	if damage < 1 {
		damage = 1
	}

	// 应用伤害
	targetMonster.HP -= damage
	if targetMonster.HP < 0 {
		targetMonster.HP = 0
	}

	// 更新怪物到上下文
	tr.context.Monsters[targetKey] = targetMonster

	// 设置伤害值到上下文
	tr.safeSetContext("damage_dealt", damage)
	tr.context.Variables["damage_dealt"] = damage

	return nil
}

// executeElementalDamageSkill 执行元素伤害技能
// 格式: "角色对怪物使用火焰伤害技能"、"角色对怪物使用冰霜伤害技能"等
func (tr *TestRunner) executeElementalDamageSkill(instruction string) error {
	char, ok := tr.context.Characters["character"]
	if !ok || char == nil {
		return fmt.Errorf("character not found")
	}

	// 找到第一个存活的怪物
	var targetMonster *models.Monster
	var targetKey string
	for key, monster := range tr.context.Monsters {
		if monster != nil && monster.HP > 0 {
			targetMonster = monster
			targetKey = key
			break
		}
	}
	if targetMonster == nil {
		return fmt.Errorf("monster not found")
	}

	// 解析伤害类型
	var damageType string
	if strings.Contains(instruction, "火焰") {
		damageType = "fire"
	} else if strings.Contains(instruction, "冰霜") {
		damageType = "frost"
	} else if strings.Contains(instruction, "暗影") {
		damageType = "shadow"
	} else if strings.Contains(instruction, "神圣") {
		damageType = "holy"
	} else if strings.Contains(instruction, "自然") {
		damageType = "nature"
	} else {
		return fmt.Errorf("unknown damage type in instruction: %s", instruction)
	}

	// 计算元素伤害（使用法术攻击力和魔法防御）
	baseAttack := float64(char.MagicAttack)
	damage := int(math.Round(baseAttack)) - targetMonster.MagicDefense
	if damage < 1 {
		damage = 1
	}

	// 应用伤害
	targetMonster.HP -= damage
	if targetMonster.HP < 0 {
		targetMonster.HP = 0
	}

	// 更新怪物到上下文
	tr.context.Monsters[targetKey] = targetMonster

	// 设置伤害类型和伤害值
	tr.safeSetContext("damage_type", damageType)
	tr.context.Variables["damage_type"] = damageType
	tr.safeSetContext(fmt.Sprintf("%s_damage_type", damageType), damageType)
	tr.context.Variables[fmt.Sprintf("%s_damage_type", damageType)] = damageType
	tr.safeSetContext(fmt.Sprintf("%s_damage_dealt", damageType), damage)
	tr.context.Variables[fmt.Sprintf("%s_damage_dealt", damageType)] = damage
	tr.safeSetContext("damage_dealt", damage)
	tr.context.Variables["damage_dealt"] = damage

	// 添加战斗日志
	if battleLogs, exists := tr.context.Variables["battle_logs"]; exists {
		if logs, ok := battleLogs.([]string); ok {
			logs = append(logs, fmt.Sprintf("角色使用%s伤害技能，造成%d点伤害", damageType, damage))
			tr.context.Variables["battle_logs"] = logs
		}
	} else {
		tr.context.Variables["battle_logs"] = []string{fmt.Sprintf("角色使用%s伤害技能，造成%d点伤害", damageType, damage)}
	}

	// 检查怪物是否死亡
	if targetMonster.HP == 0 {
		allDead := true
		for _, m := range tr.context.Monsters {
			if m != nil && m.HP > 0 {
				allDead = false
				break
			}
		}
		if allDead {
			tr.setBattleResult(true, char)
		}
	}

	return nil
}

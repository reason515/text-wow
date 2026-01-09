package runner

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"text-wow/internal/database"
	"text-wow/internal/game"
	"text-wow/internal/models"
	"text-wow/internal/repository"

	"gopkg.in/yaml.v3"
)
// Battle 相关函数

func (tr *TestRunner) executeBuildTurnOrder() error {
	// 使用与executeStartBattle相同的逻辑构建回合顺序
	return tr.buildTurnOrder()
}

func (tr *TestRunner) buildTurnOrder() error {
	// 收集所有参与者（角色和怪物�?	type participant struct {
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
				// 如果�?character"，检查是否有character_1，如果没有则使用character_1
				if _, exists := tr.context.Characters["character_1"]; !exists {
					// 如果没有character_1，使用character_1作为ID
					charID = "character_1"
				} else {
					// 如果有character_1，跳过这�?character"（避免重复）
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
			// 如果key�?monster"，则使用"monster_1"格式
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
			// 使用entry中的id（已经从key提取�?			charID := p.entry["id"].(string)
			tr.safeSetContext(fmt.Sprintf("turn_order[%d].character.id", idx), charID)
			tr.context.Variables[fmt.Sprintf("turn_order[%d].character.id", idx)] = charID
		} else {
			// p.key可能是monster_1, monster_2等，直接使用，不需要再加monster_前缀
			monsterID := p.key
			// 如果key�?monster"，则使用"monster_1"格式
			if p.key == "monster" {
				monsterID = "monster_1"
			}
			tr.safeSetContext(fmt.Sprintf("turn_order[%d].monster.id", idx), monsterID)
			tr.context.Variables[fmt.Sprintf("turn_order[%d].monster.id", idx)] = monsterID
		}
	}

	// 设置完整的turn_order数组（确保可序列化）
	if isSerializable(turnOrder) {
		tr.safeSetContext("turn_order", turnOrder)
	tr.context.Variables["turn_order"] = turnOrder
	} else {
		debugPrint("[DEBUG] buildTurnOrder: turn_order is not serializable, skipping\n")
	}
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

	// 计算伤害（考虑Debuff减成�?	baseAttack := float64(char.PhysicalAttack)
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
			logs = append(logs, fmt.Sprintf("角色攻击怪物，造成%d点伤�?, damage))
			tr.context.Variables["battle_logs"] = logs
		}
	} else {
		tr.context.Variables["battle_logs"] = []string{fmt.Sprintf("角色攻击怪物，造成%d点伤�?, damage)}
	}

	// 设置伤害值到上下�?	tr.safeSetContext("damage_dealt", damage)
	tr.context.Variables["damage_dealt"] = damage

	// 战士攻击时获得怒气（假设获�?0点）
	if char.ResourceType == "rage" {
		char.Resource += 10
		if char.Resource > char.MaxResource {
			char.Resource = char.MaxResource
		}
	}

	// 更新上下�?	tr.context.Characters["character"] = char
	// 更新怪物到上下文
	if targetKey != "" {
		tr.context.Monsters[targetKey] = targetMonster
	}

	// 如果怪物HP�?，战斗结束，战士怒气�?
	if targetMonster.HP == 0 {
		if char.ResourceType == "rage" {
			char.Resource = 0
			tr.context.Characters["character"] = char
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

	// 计算伤害（考虑Buff加成�?	baseAttack := float64(attackerMonster.PhysicalAttack)
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
			logs = append(logs, fmt.Sprintf("怪物攻击角色，造成%d点伤�?, damage))
			tr.context.Variables["battle_logs"] = logs
		}
	} else {
		tr.context.Variables["battle_logs"] = []string{fmt.Sprintf("怪物攻击角色，造成%d点伤�?, damage)}
	}

	// 设置伤害值到上下�?	tr.safeSetContext("monster_damage_dealt", damage)
	tr.context.Variables["monster_damage_dealt"] = damage

	debugPrint("[DEBUG] executeMonsterAttack: after damage - char.HP=%d, char.Resource=%d\n", char.HP, char.Resource)

	// 如果角色HP�?，战斗失败，战士怒气�?（在获得怒气之前检查）
	// 注意：必须在应用伤害后立即检查，不能先获得怒气
	if char.HP == 0 {
		if char.ResourceType == "rage" {
			char.Resource = 0
			// 更新数据�?			charRepo := repository.NewCharacterRepository()
			charRepo.UpdateAfterBattle(char.ID, char.HP, char.Resource, char.Exp, char.Level,
				char.ExpToNext, char.MaxHP, char.MaxResource, char.PhysicalAttack, char.MagicAttack, char.PhysicalDefense, char.MagicDefense,
				char.Strength, char.Agility, char.Intellect, char.Stamina, char.Spirit, char.UnspentPoints, char.TotalKills)
		}
		// 如果角色死亡，不再获得怒气，直接返�?		tr.context.Characters["character"] = char
		debugPrint("[DEBUG] executeMonsterAttack: character died, HP=0, rage reset to 0 (was %d)\n", originalResource)
		return nil
	}

	// 只有在角色未死亡时，才获得怒气
	// 战士受到伤害时获得怒气（假设获�?点）
	if char.ResourceType == "rage" {
		char.Resource += 5
		if char.Resource > char.MaxResource {
			char.Resource = char.MaxResource
		}
		debugPrint("[DEBUG] executeMonsterAttack: character took damage, HP=%d, rage increased from %d to %d\n", char.HP, originalResource, char.Resource)
	}

	// 更新上下�?	tr.context.Characters["character"] = char

	return nil
}

func (tr *TestRunner) executeBattleRound(instruction string) error {
	// 解析回合数（�?执行�?回合"�?执行一个回�?�?	roundNum := 1
	if strings.Contains(instruction, "�?) {
		parts := strings.Split(instruction, "�?)
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

		// 减少技能冷却时�?		skillManager := game.NewSkillManager()
	char, ok := tr.context.Characters["character"]
	if ok && char != nil {
		if err := skillManager.LoadCharacterSkills(char.ID); err == nil {
				// 先减少冷却时�?				skillManager.TickCooldowns(char.ID)
				
				// 减少Buff持续时间（每回合�?�?				if buffDuration, exists := tr.context.Variables["character_buff_duration"]; exists {
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
				
				// 减少护盾持续时间（每回合�?�?				if shieldDuration, exists := tr.context.Variables["character.shield_duration"]; exists {
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

			// 获取技能状态，检查是否可用（不再从Variables读取Skill对象，避免序列化错误�?			skillID, exists := tr.context.Variables["skill_id"]
			if exists {
				skillIDStr, ok := skillID.(string)
				if ok && skillIDStr != "" {
					skillState := skillManager.GetSkillState(char.ID, skillIDStr)
					if skillState != nil {
						tr.safeSetContext(fmt.Sprintf("skill_usable_round_%d", roundNum), skillState.CooldownLeft == 0)
						tr.safeSetContext(fmt.Sprintf("skill_cooldown_round_%d", roundNum), skillState.CooldownLeft)
					} else {
						// 如果技能状态不存在，从Variables获取冷却时间并计�?						cooldown := 0
						if cooldownVal, exists := tr.context.Variables["skill_cooldown"]; exists {
							if cd, ok := cooldownVal.(int); ok {
								cooldown = cd
							}
						}
						// 假设�?回合使用了技能，冷却时间�?，那么：
						// �?回合：冷却剩�?，不可用
						// �?回合：冷却剩�?，不可用
						// �?回合：冷却剩�?，可�?						cooldownLeft := cooldown - (roundNum - 1)
						if cooldownLeft < 0 {
							cooldownLeft = 0
						}
						tr.safeSetContext(fmt.Sprintf("skill_usable_round_%d", roundNum), cooldownLeft == 0)
						tr.safeSetContext(fmt.Sprintf("skill_cooldown_round_%d", roundNum), cooldownLeft)
					}
				}
			}
		} else {
			// 如果角色没有技能，从上下文获取技能信息（不再从Variables读取Skill对象�?			if _, exists := tr.context.Variables["skill_id"]; exists {
				// 从Variables获取冷却时间并计�?				cooldown := 0
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

	// 处理怪物技能冷却时间（不再从Variables读取Skill对象，避免序列化错误�?	if monsterSkillID, exists := tr.context.Variables["monster_skill_id"]; exists && monsterSkillID != nil {
		// 从Variables获取怪物技能冷却时�?		monsterCooldown := 0
		if cooldownVal, exists := tr.context.Variables["monster_skill_cooldown"]; exists {
			if cd, ok := cooldownVal.(int); ok {
				monsterCooldown = cd
			}
		}
		// 获取上次使用技能的回合�?		lastUsedRound := 1
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

	return nil
}

func (tr *TestRunner) executeRemainingMonstersAttack(instruction string) error {
	// 解析剩余怪物数量（如"剩余2个怪物攻击角色"�?	expectedCount := 0
	if strings.Contains(instruction, "剩余") {
		parts := strings.Split(instruction, "剩余")
		if len(parts) > 1 {
			countStr := strings.TrimSpace(strings.Split(parts[1], "�?)[0])
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
	// 检查是否处于休息状�?	isResting, exists := tr.context.Variables["is_resting"]
	if !exists || isResting == nil || !isResting.(bool) {
		// 如果不在休息状态，先进入休息状�?		if err := tr.checkAndEnterRest(); err != nil {
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

	// 更新上下�?	tr.context.Characters["character"] = char
	tr.safeSetContext("character.hp", char.HP)
	tr.safeSetContext("character.resource", char.Resource)
	tr.safeSetContext("character.max_hp", char.MaxHP)
	tr.safeSetContext("character.max_resource", char.MaxResource)

	return nil
}

func (tr *TestRunner) setBattleResult(isVictory bool, char *models.Character) {
	// 设置战斗状�?	if isVictory {
		tr.safeSetContext("battle_state", "victory")
		tr.context.Variables["battle_state"] = "victory"
		// 添加战斗日志
		if battleLogs, exists := tr.context.Variables["battle_logs"]; exists {
			if logs, ok := battleLogs.([]string); ok {
				logs = append(logs, "战斗胜利")
				tr.context.Variables["battle_logs"] = logs
			}
		}
		// 检查是否应该进入休息状�?		if err := tr.checkAndEnterRest(); err != nil {
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

	// 设置角色死亡状�?	if char != nil {
		isDead := char.HP <= 0
		tr.safeSetContext("character.is_dead", isDead)
		tr.context.Variables["character.is_dead"] = isDead

		// 如果胜利，给予经验和金币奖励
		if isVictory {
			// 计算经验奖励（基于怪物数量�?			expGain := len(tr.context.Monsters) * 10 // 简化：每个怪物10经验
			char.Exp += expGain
			tr.safeSetContext("character.exp", char.Exp)
			tr.context.Variables["character.exp"] = char.Exp
			tr.safeSetContext("character.exp_gained", expGain)
			tr.context.Variables["character.exp_gained"] = expGain

			// 计算金币奖励（简化：每个怪物10-30金币�?			goldGain := len(tr.context.Monsters) * 15 // 简化：每个怪物15金币
			userRepo := repository.NewUserRepository()
			if user, err := userRepo.GetByID(char.UserID); err == nil && user != nil {
				newGold := user.Gold + goldGain
				userRepo.UpdateGold(char.UserID, newGold)
				tr.safeSetContext("character.gold", newGold)
				tr.context.Variables["character.gold"] = newGold
				tr.safeSetContext("character.gold_gained", goldGain)
				tr.context.Variables["character.gold_gained"] = goldGain
			}

			// 设置team_total_exp（单角色时等于character.exp�?			tr.safeSetContext("team_total_exp", char.Exp)
			tr.context.Variables["team_total_exp"] = char.Exp
		} else {
			// 失败时，exp_gained和gold_gained�?
			tr.safeSetContext("character.exp_gained", 0)
			tr.context.Variables["character.exp_gained"] = 0
			tr.safeSetContext("character.gold_gained", 0)
			tr.context.Variables["character.gold_gained"] = 0
		}

		// 设置team_alive_count（单角色时，如果角色死亡则为0，否则为1�?		aliveCount := 0
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

		// 如果角色是战士，确保怒气�?
		if char.ResourceType == "rage" {
			char.Resource = 0
			char.MaxResource = 100
			// 更新数据�?			charRepo := repository.NewCharacterRepository()
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


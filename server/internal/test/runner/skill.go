package runner

import (
	"fmt"
	"strconv"
	"strings"

	"text-wow/internal/models"
	"text-wow/internal/repository"
)

// createSkill 创建技能（用于测试）
func (tr *TestRunner) createSkill(instruction string) error {
	// 默认资源消耗：如果是治疗技能，设为0（测试环境）；否则设为30
	defaultResourceCost := 30
	if strings.Contains(instruction, "治疗") || strings.Contains(instruction, "恢复") {
		defaultResourceCost = 0 // 治疗技能在测试中默认不消耗资源
	}

	skill := &models.Skill{
		ID:           "test_skill",
		Name:         "测试技能",
		Type:         "attack",
		ResourceCost: defaultResourceCost,
		Cooldown:     0,
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

	// 解析冷却时间（如"冷却时间3回合"）
	if strings.Contains(instruction, "冷却时间") {
		parts := strings.Split(instruction, "冷却时间")
		if len(parts) > 1 {
			cooldownStr := strings.TrimSpace(strings.Split(parts[1], "回合")[0])
			if strings.Contains(cooldownStr, "=") {
				cooldownParts := strings.Split(cooldownStr, "=")
				if len(cooldownParts) > 1 {
					cooldownStr = strings.TrimSpace(cooldownParts[1])
				}
			}
			if cooldown, err := strconv.Atoi(cooldownStr); err == nil {
				skill.Cooldown = cooldown
			}
		}
	}

	// 解析伤害倍率（如"伤害倍率=50%"或"伤害倍率150%"）
	debugPrint("[DEBUG] createSkill: checking for damage multiplier in instruction: %s\n", instruction)
	if strings.Contains(instruction, "伤害倍率") {
		parts := strings.Split(instruction, "伤害倍率")
		debugPrint("[DEBUG] createSkill: found damage multiplier, parts=%v\n", parts)
		if len(parts) > 1 {
			multiplierStr := parts[1]
			debugPrint("[DEBUG] createSkill: multiplierStr before processing: %s\n", multiplierStr)
			// 移除百分号
			multiplierStr = strings.ReplaceAll(multiplierStr, "%", "")
			// 移除逗号和其他分隔符
			multiplierStr = strings.TrimSpace(strings.Split(multiplierStr, "）")[0])
			multiplierStr = strings.TrimSpace(strings.Split(multiplierStr, "，")[0])
			// 处理"="
			if strings.Contains(multiplierStr, "=") {
				multParts := strings.Split(multiplierStr, "=")
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
					debugPrint("[DEBUG] createSkill: parsed damage multiplier %f -> %f\n", multiplier, skill.ScalingRatio)
				}
			}
		}
	}

	// 解析治疗量（如"治疗=30"或"治疗=20"）
	if strings.Contains(instruction, "治疗=") {
		parts := strings.Split(instruction, "治疗=")
		if len(parts) > 1 {
			healStr := strings.TrimSpace(strings.Split(parts[1], "）")[0])
			healStr = strings.TrimSpace(strings.Split(healStr, ",")[0])
			// 解析"=20"格式
			if strings.Contains(healStr, "=") {
				healParts := strings.Split(healStr, "=")
				if len(healParts) > 1 {
					healStr = strings.TrimSpace(healParts[1])
				}
			}
			if heal, err := strconv.Atoi(healStr); err == nil {
				skill.Type = "heal"
				// 将治疗量存储到上下文
				tr.context.Variables["skill_heal_amount"] = heal
				// 如果是治疗技能且没有明确指定资源消耗，设置为0（测试环境）
				if !strings.Contains(instruction, "消耗") {
					skill.ResourceCost = 0
					debugPrint("[DEBUG] createSkill: set ResourceCost=0 for heal skill (test environment)\n")
				}
				debugPrint("[DEBUG] createSkill: parsed heal amount=%d\n", heal)
			}
		}
	}

	// 解析Buff效果（如"攻击力+50%，持续3回合"或"效果：攻击力+50%，持续3回合"）
	if strings.Contains(instruction, "Buff") || strings.Contains(instruction, "效果=") || strings.Contains(instruction, "效果:") {
		skill.Type = "buff" // 设置为Buff技能类型
		if strings.Contains(instruction, "攻击力") && strings.Contains(instruction, "%") {
			// 解析攻击力加成百分比（如"攻击力+50%"或"效果：攻击力+50%"）
			parts := strings.Split(instruction, "攻击力")
			if len(parts) > 1 {
				modifierPart := parts[1]
				// 查找 + 号后的数字
				if plusIdx := strings.Index(modifierPart, "+"); plusIdx >= 0 {
					modifierStr := modifierPart[plusIdx+1:]
					modifierStr = strings.TrimSpace(strings.Split(modifierStr, "%")[0])
					if modifier, err := strconv.ParseFloat(modifierStr, 64); err == nil {
						tr.context.Variables["skill_buff_attack_modifier"] = modifier / 100.0 // 转换为小数（50% -> 0.5）
						debugPrint("[DEBUG] createSkill: parsed buff attack modifier=%f (from %s%%)\n", modifier/100.0, modifierStr)
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
					debugPrint("[DEBUG] createSkill: parsed buff duration=%d\n", duration)
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
		debugPrint("[DEBUG] createSkill: detected AOE skill, set skill_is_aoe=true\n")
	}

	// 如果技能类型仍未设置，默认为攻击技能
	if skill.Type == "" {
		skill.Type = "attack"
	}

	// 存储到上下文（只存储基本字段，不存储整个对象）
	tr.context.Variables["skill_id"] = skill.ID
	tr.context.Variables["skill_type"] = skill.Type
	tr.context.Variables["skill_name"] = skill.Name
	// 确保skill_scaling_ratio被正确存储（如果为0，使用默认1.0）
	if skill.ScalingRatio > 0 {
		tr.context.Variables["skill_scaling_ratio"] = skill.ScalingRatio
	} else {
		// 如果ScalingRatio为0，使用默认1.0
		skill.ScalingRatio = 1.0
		tr.context.Variables["skill_scaling_ratio"] = 1.0
		debugPrint("[DEBUG] createSkill: ScalingRatio was 0, using default 1.0\n")
	}
	debugPrint("[DEBUG] createSkill: stored skill, ScalingRatio=%f, skill_scaling_ratio=%v\n", skill.ScalingRatio, tr.context.Variables["skill_scaling_ratio"])
	return nil
}

// executeLearnSkill 执行学习技能
func (tr *TestRunner) executeLearnSkill(instruction string) error {
	char, ok := tr.context.Characters["character"]
	if !ok || char == nil {
		tr.safeSetContext("skill_learned", false)
		tr.safeSetContext("error_message", "角色不存在")
		return fmt.Errorf("character not found")
	}

	// 从上下文获取技能ID（不再从Variables读取Skill对象，避免序列化错误）
	skillID, exists := tr.context.Variables["skill_id"]
	if !exists {
		tr.safeSetContext("skill_learned", false)
		tr.safeSetContext("error_message", "技能不存在，请先创建技能")
		return fmt.Errorf("skill not found in context, please create a skill first")
	}

	skillIDStr, ok := skillID.(string)
	if !ok {
		tr.safeSetContext("skill_learned", false)
		tr.safeSetContext("error_message", "技能ID无效")
		return fmt.Errorf("skill_id is not a valid string")
	}

	// 从数据库加载技能对象
	skillRepo := repository.NewSkillRepository()
	skill, err := skillRepo.GetSkillByID(skillIDStr)
	if err != nil || skill == nil {
		// 如果数据库中没有，从Variables中的基本字段重新构建Skill对象
		skill = &models.Skill{
			ID: skillIDStr,
		}
		if skillName, exists := tr.context.Variables["skill_name"]; exists {
			if name, ok := skillName.(string); ok {
				skill.Name = name
			}
		}
		if skillType, exists := tr.context.Variables["skill_type"]; exists {
			if st, ok := skillType.(string); ok {
				skill.Type = st
			}
		}
		if scalingRatio, exists := tr.context.Variables["skill_scaling_ratio"]; exists {
			if ratio, ok := scalingRatio.(float64); ok {
				skill.ScalingRatio = ratio
			}
		}
		// 设置默认值
		if skill.Type == "" {
			skill.Type = "attack"
		}
		if skill.ScalingRatio == 0 {
			skill.ScalingRatio = 1.0
		}
		if skill.ResourceCost == 0 {
			skill.ResourceCost = 30
		}
	}

	// 使用skillRepo让角色学习技能
	err = skillRepo.AddCharacterSkill(char.ID, skill.ID, 1)
	if err != nil {
		tr.safeSetContext("skill_learned", false)
		tr.safeSetContext("error_message", err.Error())
		return fmt.Errorf("failed to learn skill: %w", err)
	}

	// 设置学习成功标志
	tr.safeSetContext("skill_learned", true)
	tr.context.Variables["skill_learned"] = true
	debugPrint("[DEBUG] executeLearnSkill: character %d learned skill %s\n", char.ID, skill.ID)
	return nil
}

// executeUseSkill 执行使用技能（简化版本）
func (tr *TestRunner) executeUseSkill(instruction string) error {
	char, ok := tr.context.Characters["character"]
	if !ok || char == nil {
		return fmt.Errorf("character not found")
	}

	// 获取技能ID
	skillID, exists := tr.context.Variables["skill_id"]
	if !exists {
		return fmt.Errorf("skill not found in context, please create a skill first")
	}

	skillIDStr, ok := skillID.(string)
	if !ok {
		return fmt.Errorf("skill_id is not a valid string")
	}

	// 从数据库加载技能对象
	skillRepo := repository.NewSkillRepository()
	skill, err := skillRepo.GetSkillByID(skillIDStr)
	if err != nil || skill == nil {
		// 如果数据库中没有，从Variables中的基本字段重新构建Skill对象
		skill = &models.Skill{
			ID: skillIDStr,
		}
		if skillName, exists := tr.context.Variables["skill_name"]; exists {
			if name, ok := skillName.(string); ok {
				skill.Name = name
			}
		}
		if skillType, exists := tr.context.Variables["skill_type"]; exists {
			if st, ok := skillType.(string); ok && st != "" {
				skill.Type = st
			}
		}
		if scalingRatio, exists := tr.context.Variables["skill_scaling_ratio"]; exists {
			if ratio, ok := scalingRatio.(float64); ok && ratio > 0 {
				skill.ScalingRatio = ratio
			}
		}
		// 设置默认值
		if skill.Type == "" {
			skill.Type = "attack"
		}
		if skill.ScalingRatio == 0 {
			skill.ScalingRatio = 1.0
		}
		if skill.ResourceCost == 0 {
			skill.ResourceCost = 30
		}
	}

	// 检查资源是否足够
	if char.Resource < skill.ResourceCost {
		tr.safeSetContext("skill_used", false)
		tr.safeSetContext("error_message", fmt.Sprintf("资源不足: 需要%d，当前%d", skill.ResourceCost, char.Resource))
		return nil
	}

	// 消耗资源
	char.Resource -= skill.ResourceCost
	if char.Resource < 0 {
		char.Resource = 0
	}

	// 设置技能使用结果
	tr.safeSetContext("skill_used", true)
	tr.context.Variables["skill_used"] = true

	// 根据技能类型处理不同效果
	if skill.Type == "" {
		skill.Type = "attack"
	}

	// 更新上下文
	tr.context.Characters["character"] = char

	// 更新数据库
	charRepo := repository.NewCharacterRepository()
	if err := charRepo.Update(char); err != nil {
		return fmt.Errorf("failed to update character: %w", err)
	}

	return nil
}

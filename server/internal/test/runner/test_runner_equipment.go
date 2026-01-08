package runner

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"text-wow/internal/models"
	"text-wow/internal/repository"
)

// generateEquipmentWithAttributes 生成带指定属性的装备
func (tr *TestRunner) generateEquipmentWithAttributes(instruction string) error {
	// 先解析装备要求（等级、职业、属性），再确定level
	levelRequired := 0
	classRequired := ""
	strengthRequired := 0
	agilityRequired := 0
	intellectRequired := 0
	
	// 解析等级要求：如"需要10级才能装备"或"需要10级才能装备的武器"
	// 注意：必须同时包含"需要"和"级才能装备"，且不包含"只有"（避免与职业要求混淆）
	if strings.Contains(instruction, "需要") && strings.Contains(instruction, "级才能装备") && !strings.Contains(instruction, "只有") {
		re := regexp.MustCompile(`需要(\d+)级才能装备`)
		matches := re.FindStringSubmatch(instruction)
		if len(matches) > 1 {
			if l, err := strconv.Atoi(matches[1]); err == nil {
				levelRequired = l
				debugPrint("[DEBUG] generateEquipmentWithAttributes: parsed level_required=%d\n", levelRequired)
			}
		}
	}
	
	// 解析职业要求：如"只有战士才能装备"
	if strings.Contains(instruction, "只有") && strings.Contains(instruction, "才能装备") {
		// 使用更宽泛的正则表达式匹配中文职业名
		re := regexp.MustCompile(`只有([^才]+)才能装备`)
		matches := re.FindStringSubmatch(instruction)
		if len(matches) > 1 {
			classRequired = strings.TrimSpace(matches[1])
			// 转换中文职业名称为ID
			if classRequired == "战士" {
				classRequired = "warrior"
			} else if classRequired == "法师" {
				classRequired = "mage"
			} else if classRequired == "盗贼" {
				classRequired = "rogue"
			} else if classRequired == "牧师" {
				classRequired = "priest"
			}
		}
	}
	
	// 解析属性要求：如"需要力量5才能装备"
	// 注意：必须包含"需要"和属性名称（力量/敏捷/智力）和"才能装备"，且不包含"级"（避免与等级要求混淆）
	if strings.Contains(instruction, "需要") && (strings.Contains(instruction, "力量") || strings.Contains(instruction, "敏捷") || strings.Contains(instruction, "智力")) && strings.Contains(instruction, "才能装备") && !strings.Contains(instruction, "级") {
		// 力量要求
		if strings.Contains(instruction, "力量") {
			re := regexp.MustCompile(`需要力量(\d+)才能装备`)
			matches := re.FindStringSubmatch(instruction)
			if len(matches) > 1 {
				if s, err := strconv.Atoi(matches[1]); err == nil {
					strengthRequired = s
				}
			}
		}
		// 敏捷要求
		if strings.Contains(instruction, "敏捷") {
			re := regexp.MustCompile(`需要敏捷(\d+)才能装备`)
			matches := re.FindStringSubmatch(instruction)
			if len(matches) > 1 {
				if a, err := strconv.Atoi(matches[1]); err == nil {
					agilityRequired = a
				}
			}
		}
		// 智力要求
		if strings.Contains(instruction, "智力") {
			re := regexp.MustCompile(`需要智力(\d+)才能装备`)
			matches := re.FindStringSubmatch(instruction)
			if len(matches) > 1 {
				if i, err := strconv.Atoi(matches[1]); err == nil {
					intellectRequired = i
				}
			}
		}
	}
	
	// 解析等级（用于生成装备，如果有职业或属性要求但没有等级要求，使用角色当前等级）
	level := 1
	if strings.Contains(instruction, "级") && !strings.Contains(instruction, "需要") {
		// 提取数字，如"10级"（但不是"需要10级才能装备"）
		re := regexp.MustCompile(`(\d+)级`)
		matches := re.FindStringSubmatch(instruction)
		if len(matches) > 1 {
			if l, err := strconv.Atoi(matches[1]); err == nil {
				level = l
			}
		}
	} else if char, ok := tr.context.Characters["character"]; ok {
		level = char.Level
	}
	
	// 如果有职业要求或属性要求但没有等级要求，使用角色当前等级（避免等级检查失败）
	if (classRequired != "" || strengthRequired > 0 || agilityRequired > 0 || intellectRequired > 0) && levelRequired == 0 {
		if char, ok := tr.context.Characters["character"]; ok {
			level = char.Level
		}
	}
	
	// 解析装备类型和槽
	itemID := "worn_sword"
	slot := "main_hand"
	if strings.Contains(instruction, "武器") {
		itemID = "worn_sword"
		slot = "main_hand"
	} else if strings.Contains(instruction, "护甲") || strings.Contains(instruction, "盔甲") {
		itemID = "cloth_robe"
		slot = "armor"
	} else if strings.Contains(instruction, "饰品") {
		itemID = "ring"
		slot = "accessory"
	} else if strings.Contains(instruction, "盾牌") {
		itemID = "wooden_shield"
		slot = "off_hand"
	}
	
	// 解析品质（默认蓝色）
	quality := "rare"
	if strings.Contains(instruction, "白色") || strings.Contains(instruction, "common") {
		quality = "common"
	} else if strings.Contains(instruction, "绿色") || strings.Contains(instruction, "uncommon") {
		quality = "uncommon"
	} else if strings.Contains(instruction, "蓝色") || strings.Contains(instruction, "rare") {
		quality = "rare"
	} else if strings.Contains(instruction, "紫色") || strings.Contains(instruction, "epic") {
		quality = "epic"
	} else if strings.Contains(instruction, "橙色") || strings.Contains(instruction, "legendary") {
		quality = "legendary"
	}
	
	// 如果有特殊要求，创建或更新item配置
	debugPrint("[DEBUG] generateEquipmentWithAttributes: checking requirements - levelRequired=%d, classRequired=%s, strengthRequired=%d, agilityRequired=%d, intellectRequired=%d\n", 
		levelRequired, classRequired, strengthRequired, agilityRequired, intellectRequired)
	if levelRequired > 0 || classRequired != "" || strengthRequired > 0 || agilityRequired > 0 || intellectRequired > 0 {
		// 生成唯一的item ID（基于要求）
		testItemID := fmt.Sprintf("test_item_%s_%d_%s_%d_%d_%d", slot, levelRequired, classRequired, strengthRequired, agilityRequired, intellectRequired)
		debugPrint("[DEBUG] generateEquipmentWithAttributes: will create test item with ID=%s\n", testItemID)
		
		// 获取基础item配置
		gameRepo := repository.NewGameRepository()
		baseItem, err := gameRepo.GetItemByID(itemID)
		if err != nil {
			// 如果基础item不存在，创建一个默认配
			baseItem = map[string]interface{}{
				"id":            testItemID,
				"name":          "测试装备",
				"type":          "equipment",
				"slot":          slot,
				"quality":       quality,
				"level_required": levelRequired,
				"class_required": classRequired,
				"strength_required": strengthRequired,
				"agility_required": agilityRequired,
				"intellect_required": intellectRequired,
			}
		} else {
			// 复制基础配置并更新要
			baseItem["id"] = testItemID
			baseItem["level_required"] = levelRequired
			if classRequired != "" {
				baseItem["class_required"] = classRequired
			}
			baseItem["strength_required"] = strengthRequired
			baseItem["agility_required"] = agilityRequired
			baseItem["intellect_required"] = intellectRequired
		}
		
		// 在数据库中插入或更新item配置
		debugPrint("[DEBUG] generateEquipmentWithAttributes: calling createOrUpdateTestItem for itemID=%s\n", testItemID)
		err = tr.createOrUpdateTestItem(baseItem)
		if err != nil {
			return fmt.Errorf("failed to create test item: %w", err)
		}
		debugPrint("[DEBUG] generateEquipmentWithAttributes: createOrUpdateTestItem succeeded for itemID=%s\n", testItemID)
		
		// 使用测试item ID
		itemID = testItemID
		debugPrint("[DEBUG] generateEquipmentWithAttributes: updated itemID to %s with requirements (level=%d, class=%s, str=%d, agi=%d, int=%d)\n", 
			itemID, levelRequired, classRequired, strengthRequired, agilityRequired, intellectRequired)
	}
	
	debugPrint("[DEBUG] generateEquipmentWithAttributes: using itemID=%s to generate equipment\n", itemID)
	
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
	
	// 生成装备
	equipment, err := tr.equipmentManager.GenerateEquipment(itemID, quality, level, ownerID)
	if err != nil {
		return fmt.Errorf("failed to generate equipment: %w", err)
	}
	
	// 解析并设置属性（如果指令中指定了）
	if strings.Contains(instruction, "攻击力") || strings.Contains(instruction, "物理攻击") {
		re := regexp.MustCompile(`攻击力\+(\d+)|物理攻击\+(\d+)`)
		matches := re.FindStringSubmatch(instruction)
		if len(matches) > 1 {
			attackValue := 0
			if matches[1] != "" {
				attackValue, _ = strconv.Atoi(matches[1])
			} else if matches[2] != "" {
				attackValue, _ = strconv.Atoi(matches[2])
			}
			if attackValue > 0 {
				// 通过修改词缀值来设置攻击力加成
				// 这里我们需要创建一个攻击力词缀
				// 暂时先存储到变量中，后续在穿戴时应用
				tr.context.Variables[fmt.Sprintf("equipment_%d_physical_attack", equipment.ID)] = attackValue
			}
		}
	}
	if strings.Contains(instruction, "防御") || strings.Contains(instruction, "物理防御") {
		re := regexp.MustCompile(`防御力\+(\d+)|物理防御\+(\d+)`)
		matches := re.FindStringSubmatch(instruction)
		if len(matches) > 1 {
			defenseValue := 0
			if matches[1] != "" {
				defenseValue, _ = strconv.Atoi(matches[1])
			} else if matches[2] != "" {
				defenseValue, _ = strconv.Atoi(matches[2])
			}
			if defenseValue > 0 {
				tr.context.Variables[fmt.Sprintf("equipment_%d_physical_defense", equipment.ID)] = defenseValue
			}
		}
	}
	if strings.Contains(instruction, "力量") {
		re := regexp.MustCompile(`力量\+(\d+)`)
		matches := re.FindStringSubmatch(instruction)
		if len(matches) > 1 {
			strengthValue, _ := strconv.Atoi(matches[1])
			if strengthValue > 0 {
				tr.context.Variables[fmt.Sprintf("equipment_%d_strength", equipment.ID)] = strengthValue
			}
		}
	}
	
	// 存储到上下文（只存储基本字段，不存储整个对象）
	tr.context.Variables["equipment_id"] = equipment.ID
	tr.context.Variables["equipment_item_id"] = equipment.ItemID
	tr.context.Variables["equipment_quality"] = equipment.Quality
	tr.context.Variables["equipment_slot"] = equipment.Slot
	tr.context.Equipments[fmt.Sprintf("%d", equipment.ID)] = equipment
	
	// 如果指令中提到"已装备"，直接装备
	if strings.Contains(instruction, "已装备") {
		char, ok := tr.context.Characters["character"]
		if ok {
			if err := tr.equipmentManager.EquipItem(char.ID, equipment.ID); err != nil {
				return fmt.Errorf("failed to equip item: %w", err)
			}
			// 重新加载角色以获取更新后的属性
			charRepo := repository.NewCharacterRepository()
			updatedChar, err := charRepo.GetByID(char.ID)
			if err == nil {
				tr.context.Characters["character"] = updatedChar
			}
		}
	}
	
	return nil
}

// executeEquipItem 执行穿戴装备操作
func (tr *TestRunner) executeEquipItem(instruction string) error {
	char, ok := tr.context.Characters["character"]
	if !ok || char == nil {
		return fmt.Errorf("character not found")
	}
	
	// 获取要穿戴的装备
	var equipment *models.EquipmentInstance
	if eq, ok := tr.context.Variables["equipment"].(*models.EquipmentInstance); ok {
		equipment = eq
	} else if eqID, ok := tr.context.Variables["equipment_id"].(int); ok {
		equipmentRepo := repository.NewEquipmentRepository()
		eq, err := equipmentRepo.GetByID(eqID)
		if err != nil {
			return fmt.Errorf("failed to get equipment: %w", err)
		}
		equipment = eq
	} else {
		// 尝试从Equipments map中获
		for _, eq := range tr.context.Equipments {
			equipment = eq
			break
		}
	}
	
	if equipment == nil {
		return fmt.Errorf("equipment not found")
	}
	
	// 检查是否有旧装备（用于替换测试）
	equipmentRepo := repository.NewEquipmentRepository()
	oldEquipment, err := equipmentRepo.GetByCharacterAndSlot(char.ID, equipment.Slot)
	if err == nil && oldEquipment != nil {
		tr.context.Variables["old_weapon"] = oldEquipment
		tr.context.Variables["old_equipment"] = oldEquipment
	}
	
	// 记录穿戴前的攻击力（用于断言）
	tr.context.Variables["previous_physical_attack"] = char.PhysicalAttack
	tr.context.Variables["base_physical_attack"] = char.PhysicalAttack
	
	// 穿戴装备
	if err := tr.equipmentManager.EquipItem(char.ID, equipment.ID); err != nil {
		return fmt.Errorf("failed to equip item: %w", err)
	}
	
	// 重新加载角色以获取更新后的属
	charRepo := repository.NewCharacterRepository()
	updatedChar, err := charRepo.GetByID(char.ID)
	if err != nil {
		return fmt.Errorf("failed to reload character: %w", err)
	}
	tr.context.Characters["character"] = updatedChar
	
	// 重新加载装备以获取更新后的character_id
	updatedEquipment, err := equipmentRepo.GetByID(equipment.ID)
	if err == nil {
		tr.context.Variables["equipment"] = updatedEquipment
		tr.context.Variables["new_weapon"] = updatedEquipment
		tr.context.Variables["new_equipment"] = updatedEquipment
		tr.context.Equipments[fmt.Sprintf("%d", updatedEquipment.ID)] = updatedEquipment
	}
	
	// 如果有旧装备，重新加载以获取更新后的character_id（应该是nil）
	if oldEquipment != nil {
		updatedOldEquipment, err := equipmentRepo.GetByID(oldEquipment.ID)
		if err == nil {
			tr.context.Variables["old_weapon"] = updatedOldEquipment
			tr.context.Variables["old_equipment"] = updatedOldEquipment
		}
	}
	
	return nil
}

// executeUnequipItem 执行卸下装备操作
func (tr *TestRunner) executeUnequipItem(instruction string) error {
	char, ok := tr.context.Characters["character"]
	if !ok || char == nil {
		return fmt.Errorf("character not found")
	}
	
	// 获取要卸下的装备
	var equipment *models.EquipmentInstance
	if eq, ok := tr.context.Variables["equipment"].(*models.EquipmentInstance); ok {
		equipment = eq
	} else if eqID, ok := tr.context.Variables["equipment_id"].(int); ok {
		equipmentRepo := repository.NewEquipmentRepository()
		eq, err := equipmentRepo.GetByID(eqID)
		if err != nil {
			return fmt.Errorf("failed to get equipment: %w", err)
		}
		equipment = eq
	} else {
		// 尝试从Equipments map中获
		for _, eq := range tr.context.Equipments {
			if eq.CharacterID != nil && *eq.CharacterID == char.ID {
				equipment = eq
				break
			}
		}
	}
	
	if equipment == nil {
		return fmt.Errorf("equipment not found")
	}
	
	// 记录卸下前的攻击力（用于断言）
	tr.context.Variables["previous_physical_attack"] = char.PhysicalAttack
	
	// 卸下装备
	if err := tr.equipmentManager.UnequipItem(char.ID, equipment.ID); err != nil {
		return fmt.Errorf("failed to unequip item: %w", err)
	}
	
	// 重新加载角色以获取更新后的属
	charRepo := repository.NewCharacterRepository()
	updatedChar, err := charRepo.GetByID(char.ID)
	if err != nil {
		return fmt.Errorf("failed to reload character: %w", err)
	}
	tr.context.Characters["character"] = updatedChar
	
	// 重新加载装备以获取更新后的character_id
	equipmentRepo := repository.NewEquipmentRepository()
	updatedEquipment, err := equipmentRepo.GetByID(equipment.ID)
	if err == nil {
		tr.context.Variables["equipment"] = updatedEquipment
		tr.context.Variables["weapon"] = updatedEquipment
		tr.context.Equipments[fmt.Sprintf("%d", updatedEquipment.ID)] = updatedEquipment
	}
	
	return nil
}

// executeTryEquipItem 尝试穿戴装备（用于测试失败情况）
func (tr *TestRunner) executeTryEquipItem(instruction string) error {
	debugPrint("[DEBUG] executeTryEquipItem: called with instruction=%s\n", instruction)
	
	char, ok := tr.context.Characters["character"]
	if !ok || char == nil {
		return fmt.Errorf("character not found")
	}
	
	// 获取要穿戴的装备
	var equipment *models.EquipmentInstance
	if eq, ok := tr.context.Variables["equipment"].(*models.EquipmentInstance); ok {
		equipment = eq
		debugPrint("[DEBUG] executeTryEquipItem: found equipment from Variables, ID=%d\n", equipment.ID)
	} else if eqID, ok := tr.context.Variables["equipment_id"].(int); ok {
		equipmentRepo := repository.NewEquipmentRepository()
		eq, err := equipmentRepo.GetByID(eqID)
		if err != nil {
			return fmt.Errorf("failed to get equipment: %w", err)
		}
		equipment = eq
		debugPrint("[DEBUG] executeTryEquipItem: found equipment from equipment_id, ID=%d\n", equipment.ID)
	} else {
		// 尝试从Equipments map中获
		for _, eq := range tr.context.Equipments {
			equipment = eq
			debugPrint("[DEBUG] executeTryEquipItem: found equipment from Equipments map, ID=%d\n", equipment.ID)
			break
		}
	}
	
	if equipment == nil {
		return fmt.Errorf("equipment not found")
	}
	
	// 尝试穿戴装备
	debugPrint("[DEBUG] executeTryEquipItem: attempting to equip equipment ID=%d for character ID=%d\n", equipment.ID, char.ID)
	err := tr.equipmentManager.EquipItem(char.ID, equipment.ID)
	if err != nil {
		// 装备失败，记录错误信息
		tr.context.Variables["equip_success"] = false
		errorMsg := err.Error()
		// 确保错误消息包含中文关键词（用于contains断言）
		tr.context.Variables["error_message"] = errorMsg
		debugPrint("[DEBUG] executeTryEquipItem: equip failed, error=%s, set equip_success=false, error_message=%s\n", errorMsg, errorMsg)
		return nil // 不返回错误，因为这是预期的失败
	}
	
	// 装备成功
	tr.context.Variables["equip_success"] = true
	debugPrint("[DEBUG] executeTryEquipItem: equip succeeded, set equip_success=true\n")
	
	// 重新加载角色
	charRepo := repository.NewCharacterRepository()
	updatedChar, err := charRepo.GetByID(char.ID)
	if err == nil {
		tr.context.Characters["character"] = updatedChar
	}
	
	return nil
}

// executeEquipAllItems 依次穿戴所有装备
func (tr *TestRunner) executeEquipAllItems(instruction string) error {
	char, ok := tr.context.Characters["character"]
	if !ok || char == nil {
		return fmt.Errorf("character not found")
	}
	
	// 记录穿戴前的属性（用于断言）
	tr.context.Variables["base_physical_attack"] = char.PhysicalAttack
	tr.context.Variables["base_physical_defense"] = char.PhysicalDefense
	tr.context.Variables["base_strength"] = char.Strength
	
	// 遍历所有装备并穿戴
	for _, equipment := range tr.context.Equipments {
		if equipment.CharacterID == nil {
			// 装备未穿戴，尝试穿戴
			if err := tr.equipmentManager.EquipItem(char.ID, equipment.ID); err != nil {
				// 忽略错误，继续穿戴其他装备
				continue
			}
		}
	}
	
	// 重新加载角色以获取更新后的属
	charRepo := repository.NewCharacterRepository()
	updatedChar, err := charRepo.GetByID(char.ID)
	if err != nil {
		return fmt.Errorf("failed to reload character: %w", err)
	}
	tr.context.Characters["character"] = updatedChar
	
	return nil
}


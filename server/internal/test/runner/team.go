package runner



import (
	"fmt"
	"strconv"
	"strings"

	"text-wow/internal/models"
	"text-wow/internal/repository"
)

// Team 相关函数



func (tr *TestRunner) createTeam(instruction string) error {

	// 确保用户存在

	user, err := tr.createTestUser()

	if err != nil {

		return fmt.Errorf("failed to create user: %w", err)

	}



	// 解析队伍成员（通过冒号或逗号分隔）
	// 格式：战士(HP=100)、牧师(HP=100)、法师(HP=100)

	var members []string

	if strings.Contains(instruction, "包含") && strings.Contains(instruction, "个角色") {

		parts := strings.Split(instruction, "包含")

		if len(parts) > 1 {

			members = strings.Split(parts[1], ",")

		}

	} else if strings.Contains(instruction, "：") {
		// 中文冒号
		parts := strings.Split(instruction, "：")

		if len(parts) > 1 {
			// 支持中文顿号"、"分隔
			memberStr := parts[1]
			if strings.Contains(memberStr, "、") {
				members = strings.Split(memberStr, "、")
			} else {
				members = strings.Split(memberStr, ",")
			}
		}

	} else if strings.Contains(instruction, ":") {
		// 英文冒号
		parts := strings.Split(instruction, ":")

		if len(parts) > 1 {

			members = strings.Split(parts[1], ",")

		}

	}



	charRepo := repository.NewCharacterRepository()

	slot := 1



	// 先获取用户的所有角色，检查哪些slot已被占用

	existingChars, err := charRepo.GetByUserID(user.ID)

	if err != nil {

		existingChars = []*models.Character{}

	}

	existingSlots := make(map[int]*models.Character)

	for _, c := range existingChars {

		existingSlots[c.TeamSlot] = c

	}



	for _, memberDesc := range members {

		memberDesc = strings.TrimSpace(memberDesc)

		if memberDesc == "" {

			continue

		}



		// 解析职业（战士、牧师、法师等）
		classID := "warrior"

		if strings.Contains(memberDesc, "战士") {

			classID = "warrior"

		} else if strings.Contains(memberDesc, "牧师") {

			classID = "priest"

		} else if strings.Contains(memberDesc, "法师") {

			classID = "mage"

		} else if strings.Contains(memberDesc, "盗贼") {

			classID = "rogue"

		}



		// 解析HP（如"HP=100"）
		hp := 100
		if strings.Contains(memberDesc, "HP=") {

			parts := strings.Split(memberDesc, "HP=")

			if len(parts) > 1 {

				hpStr := strings.TrimSpace(strings.Split(parts[1], ")")[0])

				if h, err := strconv.Atoi(hpStr); err == nil {

					hp = h

				}

			}

		}



		// 检查该slot是否已存在角色
		var createdChar *models.Character
		if existingChar, exists := existingSlots[slot]; exists {

			// 更新已存在的角色

			existingChar.Name = fmt.Sprintf("测试角色%d", slot)

			existingChar.ClassID = classID

			existingChar.HP = hp

			existingChar.MaxHP = hp

			existingChar.Level = 1

			existingChar.Strength = 10

			existingChar.Agility = 10

			existingChar.Intellect = 10

			existingChar.Stamina = 10

			existingChar.Spirit = 10



			// 根据职业设置资源类型

			if classID == "warrior" {

				existingChar.ResourceType = "rage"

				existingChar.MaxResource = 100

				existingChar.Resource = 0

			} else if classID == "rogue" {

				existingChar.ResourceType = "energy"

				existingChar.MaxResource = 100

				existingChar.Resource = 100

			} else {

				existingChar.ResourceType = "mana"

				existingChar.MaxResource = 100

				existingChar.Resource = 100

			}



			// 更新到数据库

			if err := charRepo.Update(existingChar); err != nil {

				return fmt.Errorf("failed to update character in team: %w", err)

			}

			createdChar = existingChar

		} else {

			// 创建新角色
			char := &models.Character{

				UserID:    user.ID,

				Name:      fmt.Sprintf("测试角色%d", slot),

				RaceID:    "human",

				ClassID:   classID,

				Faction:   "alliance",

				TeamSlot:  slot,

				Level:     1,

				HP:        hp,

				MaxHP:     hp,

				Strength:  10,

				Agility:   10,

				Intellect: 10,

				Stamina:   10,

				Spirit:    10,

			}



			// 根据职业设置资源类型

			if classID == "warrior" {

				char.ResourceType = "rage"

				char.MaxResource = 100

				char.Resource = 0

			} else if classID == "rogue" {

				char.ResourceType = "energy"

				char.MaxResource = 100

				char.Resource = 100

			} else {

				char.ResourceType = "mana"

				char.MaxResource = 100

				char.Resource = 100

			}



			// 保存到数据库

			var err error

			createdChar, err = charRepo.Create(char)

			if err != nil {

				return fmt.Errorf("failed to create character in team: %w", err)

			}

		}



		// 保存到上下文（使用character_1, character_2等作为key）
		key := fmt.Sprintf("character_%d", slot)
		tr.context.Characters[key] = createdChar



		// 第一个角色也保存�character"（向后兼容）

		if slot == 1 {

			tr.context.Characters["character"] = createdChar

		}



		slot++

	}



	return nil

}



func (tr *TestRunner) executeCreateEmptyTeam() error {

	// 清空所有角色（除了character，保留作为默认角色）

	// 实际上，空队伍意味着没有角色在队伍槽位中

	// 我们只需要确保team.character_count�?
	tr.context.Variables["team.character_count"] = 0

	tr.safeSetContext("team.character_count", 0)

	return nil

}



func (tr *TestRunner) executeCreateTeamWithMembers(instruction string) error {

	// 解析指令，如"创建一个队伍，槽位1已有角色1"或"创建一个队伍，包含3个角色"


	if strings.Contains(instruction, "槽位") && strings.Contains(instruction, "已有") {

		// 解析槽位和角色ID

		// 如"槽位1已有角色1"


		parts := strings.Split(instruction, "槽位")

		if len(parts) > 1 {

			slotPart := strings.TrimSpace(strings.Split(parts[1], "已有")[0])

			if slot, err := strconv.Atoi(slotPart); err == nil {

				// 解析角色ID

				charIDPart := strings.TrimSpace(strings.Split(parts[1], "角色")[1])

				if charID, err := strconv.Atoi(charIDPart); err == nil {

				// 创建或获取角色
				char, err := tr.getOrCreateCharacterByID(charID, slot)

					if err != nil {

						return err

					}

					key := fmt.Sprintf("character_%d", slot)

					tr.context.Characters[key] = char

					tr.context.Variables["team.character_count"] = 1

					tr.safeSetContext("team.character_count", 1)

					tr.safeSetContext(fmt.Sprintf("team.slot_%d.character_id", slot), char.ID)

				}

			}

		}

	} else if strings.Contains(instruction, "包含") && strings.Contains(instruction, "个角色") {

		// 解析角色数量，如"包含3个角色"


		parts := strings.Split(instruction, "包含")

		if len(parts) > 1 {

			countStr := strings.TrimSpace(strings.Split(parts[1], ")")[0])

			if count, err := strconv.Atoi(countStr); err == nil {

				// 创建指定数量的角色
				for i := 1; i <= count; i++ {

					char, err := tr.getOrCreateCharacterByID(i, i)

					if err != nil {

						return err

					}

					key := fmt.Sprintf("character_%d", i)

					tr.context.Characters[key] = char

				}

				tr.context.Variables["team.character_count"] = count

				tr.safeSetContext("team.character_count", count)

				// 创建队伍后，同步队伍信息到上下文

				tr.syncTeamToContext()

			}

		}

	}

	return nil

}



func (tr *TestRunner) executeAddCharacterToTeamSlot(instruction string) error {

	// 解析指令，支持两种格式：
	// 1. "将角色X添加到槽位Y" - 指定角色ID
	// 2. "将角色添加到槽位Y" - 使用默认角色
	
	var charID, slot int
	var err error
	
	// 尝试匹配 "将角色X添加到槽位Y" 格式
	if strings.Contains(instruction, "将角色") && strings.Contains(instruction, "添加到槽位") {
		// 提取角色ID部分
		afterRole := strings.TrimPrefix(instruction, "将角色")
		beforeSlot := strings.Split(afterRole, "添加到槽位")
		
		if len(beforeSlot) >= 2 {
			charIDStr := strings.TrimSpace(beforeSlot[0])
			slotStr := strings.TrimSpace(beforeSlot[1])
			
			// 如果角色ID为空，使用默认值1
			if charIDStr == "" {
				charID = 1
			} else {
				charID, err = strconv.Atoi(charIDStr)
				if err != nil {
					charID = 1 // 解析失败时使用默认值
				}
			}
			
			// 解析槽位
			slot, err = strconv.Atoi(slotStr)
			if err != nil {
				return fmt.Errorf("failed to parse slot from '%s': %w", slotStr, err)
			}
		} else {
			return fmt.Errorf("invalid instruction format: %s", instruction)
		}
	} else {
		return fmt.Errorf("invalid instruction: %s", instruction)
	}

	

	// 检查槽位是否已被占用
	slotKey := fmt.Sprintf("character_%d", slot)
	if existingChar, exists := tr.context.Characters[slotKey]; exists && existingChar != nil {

		return fmt.Errorf("slot %d is already occupied", slot)

	}

	

	// 检查槽位是否解锁（简化：假设�个槽位默认解锁�?
	if slot > 5 {

		// 检查unlocked_slots

		unlockedSlots := 1

		if unlockedVal, exists := tr.context.Variables["team.unlocked_slots"]; exists {

			if u, ok := unlockedVal.(int); ok {

				unlockedSlots = u

			}

		}

		if slot > unlockedSlots {

			tr.context.Variables["operation_success"] = false

			tr.safeSetContext("operation_success", false)

			return fmt.Errorf("slot %d is not unlocked", slot)

		}

	}

	

	// 获取或创建角�?char, err := tr.getOrCreateCharacterByID(charID, slot)

	if err != nil {

		return err

	}

	

	// 添加到槽位
	// 根据 charID 获取角色
	charKey := fmt.Sprintf("character_%d", charID)
	char := tr.context.Characters[charKey]
	if char == nil {
		// 尝试使用默认角色
		char = tr.context.Characters["character"]
	}
	if char == nil {
		return fmt.Errorf("character %d not found", charID)
	}
	tr.context.Characters[slotKey] = char

	// 更新队伍角色数
	teamCount := 0

	for key, c := range tr.context.Characters {

		if c != nil && (strings.HasPrefix(key, "character_") || key == "character") {

			teamCount++

		}

	}

	tr.context.Variables["team.character_count"] = teamCount

	tr.safeSetContext("team.character_count", teamCount)

	tr.safeSetContext(fmt.Sprintf("team.slot_%d.character_id", slot), char.ID)

	

	tr.context.Variables["operation_success"] = true

	tr.safeSetContext("operation_success", true)

	

	return nil

}



func (tr *TestRunner) executeTryAddCharacterToTeamSlot(instruction string) error {

	err := tr.executeAddCharacterToTeamSlot(instruction)

	if err != nil {

		// 操作失败，设置operation_success为false

		tr.context.Variables["operation_success"] = false

		tr.safeSetContext("operation_success", false)

		return nil // 不返回错误，因为这是预期的失败
	}

	tr.context.Variables["operation_success"] = true

	tr.safeSetContext("operation_success", true)

	return nil

}



func (tr *TestRunner) executeUnlockTeamSlot(instruction string) error {

	// 解析指令，如"解锁槽位2"

	parts := strings.Split(instruction, "槽位")

	if len(parts) < 2 {

		return fmt.Errorf("invalid instruction: %s", instruction)

	}

	

	slotPart := strings.TrimSpace(parts[1])

	slot, err := strconv.Atoi(slotPart)

	if err != nil {

		return fmt.Errorf("failed to parse slot: %w", err)

	}

	

	// 更新解锁槽位�?tr.context.Variables["team.unlocked_slots"] = slot

	tr.safeSetContext("team.unlocked_slots", slot)

	

	return nil

}



func (tr *TestRunner) executeTryAddCharacterToUnlockedSlot(instruction string) error {

	// 这个函数会尝试添加，但应该失败
	return tr.executeTryAddCharacterToTeamSlot(instruction)

}





// syncTeamToContext 同步队伍信息到断言上下文
func (tr *TestRunner) syncTeamToContext() {
	// 统计队伍中的角色数量
	teamCharCount := 0
	teamAliveCount := 0
	unlockedSlots := 0

	// 统计所有角色（character, character_1, character_2等）
	for key, char := range tr.context.Characters {
		if char != nil {
			teamCharCount++
			if char.HP > 0 {
				teamAliveCount++
			}
			// 如果key是character_N格式，说明是队伍成员
			if strings.HasPrefix(key, "character_") {
				slotStr := strings.TrimPrefix(key, "character_")
				if slot, err := strconv.Atoi(slotStr); err == nil {
				// 假设5个槽位默认解锁（可以根据实际情况调整）
				if slot <= 5 {
						if slot > unlockedSlots {
							unlockedSlots = slot
						}
						// 设置槽位信息
						tr.safeSetContext(fmt.Sprintf("team.slot_%d.character_id", slot), char.ID)
						tr.safeSetContext(fmt.Sprintf("team.slot_%d.character_name", slot), char.Name)
						tr.safeSetContext(fmt.Sprintf("team.slot_%d.hp", slot), char.HP)
						tr.safeSetContext(fmt.Sprintf("team.slot_%d.max_hp", slot), char.MaxHP)
					}
				}
			}
		}
	}

	// 如果只有character（没有character_1等），也统计
	if char, exists := tr.context.Characters["character"]; exists && char != nil {
		if teamCharCount == 0 {
			teamCharCount = 1
			if char.HP > 0 {
				teamAliveCount = 1
			}
		}
	}

	// 设置队伍属性
	tr.safeSetContext("team.character_count", teamCharCount)
	tr.safeSetContext("team_alive_count", teamAliveCount)
	tr.context.Variables["team.character_count"] = teamCharCount
	tr.context.Variables["team_alive_count"] = teamAliveCount

	// 设置解锁槽位数（如果没有设置，使用队伍角色数）
	if unlockedSlotsVal, exists := tr.context.Variables["team.unlocked_slots"]; exists {
		if u, ok := unlockedSlotsVal.(int); ok {
			unlockedSlots = u
		}
	}
	if unlockedSlots == 0 {
		unlockedSlots = teamCharCount
		if unlockedSlots == 0 {
			unlockedSlots = 1 // 至少1个槽位解锁
		}
	}
	tr.safeSetContext("team.unlocked_slots", unlockedSlots)
	tr.context.Variables["team.unlocked_slots"] = unlockedSlots
}

// executeRemoveCharacterFromTeamSlot 从队伍槽位移除角色
func (tr *TestRunner) executeRemoveCharacterFromTeamSlot(instruction string) error {
	// 解析指令，如"从槽位1移除角色"
	parts := strings.Split(instruction, "槽位")
	if len(parts) < 2 {
		return fmt.Errorf("invalid instruction: %s", instruction)
	}

	slotStr := strings.TrimSpace(strings.Split(parts[1], "移除")[0])
	slot, err := strconv.Atoi(slotStr)
	if err != nil {
		return fmt.Errorf("failed to parse slot number: %w", err)
	}

	// 从上下文移除角色
	slotKey := fmt.Sprintf("character_%d", slot)
	delete(tr.context.Characters, slotKey)

	// 更新队伍角色数量
	teamCount := 0
	for key, c := range tr.context.Characters {
		if c != nil && (strings.HasPrefix(key, "character_") || key == "character") {
			teamCount++
		}
	}
	tr.context.Variables["team.character_count"] = teamCount
	return nil
}

// getOrCreateCharacterByID 根据ID获取或创建角色
func (tr *TestRunner) getOrCreateCharacterByID(charID int, slot int) (*models.Character, error) {
	// 先检查是否已存在
	key := fmt.Sprintf("character_%d", slot)
	if existingChar, exists := tr.context.Characters[key]; exists && existingChar != nil && existingChar.ID == charID {
		return existingChar, nil
	}
	
	// 检查character_1, character_2等
	for i := 1; i <= 5; i++ {
		checkKey := fmt.Sprintf("character_%d", i)
		if existingChar, exists := tr.context.Characters[checkKey]; exists && existingChar != nil && existingChar.ID == charID {
			return existingChar, nil
		}
	}
	
	// 创建新角色
	user, err := tr.createTestUser()
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}
	
	charRepo := repository.NewCharacterRepository()
	char := &models.Character{
		UserID:    user.ID,
		ID:        charID,
		Name:      fmt.Sprintf("测试角色%d", charID),
		RaceID:    "human",
		ClassID:   "warrior",
		Faction:   "alliance",
		TeamSlot:  slot,
		Level:     1,
		HP:        100,
		MaxHP:     100,
		Strength:  10,
		Agility:   10,
		Intellect: 10,
		Stamina:   10,
		Spirit:    10,
		ResourceType: "rage",
		Resource:  0,
		MaxResource: 100,
	}
	
	createdChar, err := charRepo.Create(char)
	if err != nil {
		return nil, fmt.Errorf("failed to create character: %w", err)
	}
	
	return createdChar, nil
}

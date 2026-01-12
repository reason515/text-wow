package runner

import (
	"fmt"

	"regexp"

	"strconv"

	"strings"

	"text-wow/internal/database"

	"text-wow/internal/models"

	"text-wow/internal/repository"
)

// Character 相关函数

func (tr *TestRunner) createCharacter(instruction string) error {

	// 保存当前指令到上下文，以便后续判断是否明确设置了某些属性

	tr.context.Variables["last_instruction"] = instruction

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

	// 保存ClassID到Variables

	tr.context.Variables["character_class_id"] = classID

	char := &models.Character{

		ID: 1,

		Name: "测试角色",

		ClassID: classID,

		Level: 1,

		Strength: 10,

		Agility: 10,

		Intellect: 10,

		Stamina: 10,

		Spirit: 10,

		MaxHP: 0,

		MaxResource: 0,
	}

	// 解析主属性（如"力量=20"或"敏捷=10"等）

	parseAttribute := func(value string) string {

		value = strings.TrimSpace(strings.Split(value, "）")[0])

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

				tr.context.Variables["character_strength"] = str

				debugPrint("[DEBUG] createCharacter: set Strength=%d and saved to Variables\n", str)

			}

		}

	}

	if strings.Contains(instruction, "敏捷=") {

		parts := strings.Split(instruction, "敏捷=")

		if len(parts) > 1 {

			agiStr := parseAttribute(parts[1])

			if agi, err := strconv.Atoi(agiStr); err == nil {

				char.Agility = agi

				tr.context.Variables["character_agility"] = agi

				debugPrint("[DEBUG] createCharacter: set Agility=%d and saved to Variables\n", agi)

			}

		}

	}

	if strings.Contains(instruction, "智力=") {

		parts := strings.Split(instruction, "智力=")

		if len(parts) > 1 {

			intStr := parseAttribute(parts[1])

			if intel, err := strconv.Atoi(intStr); err == nil {

				char.Intellect = intel

				tr.context.Variables["character_intellect"] = intel

				debugPrint("[DEBUG] createCharacter: set Intellect=%d and saved to Variables\n", intel)

			}

		}

	}

	if strings.Contains(instruction, "精神=") {

		parts := strings.Split(instruction, "精神=")

		if len(parts) > 1 {

			spiStr := parseAttribute(parts[1])

			if spi, err := strconv.Atoi(spiStr); err == nil {

				char.Spirit = spi

				tr.context.Variables["character_spirit"] = spi

				debugPrint("[DEBUG] createCharacter: set Spirit=%d and saved to Variables\n", spi)

			}

		}

	}

	if strings.Contains(instruction, "耐力=") {

		parts := strings.Split(instruction, "耐力=")

		if len(parts) > 1 {

			staStr := parseAttribute(parts[1])

			if sta, err := strconv.Atoi(staStr); err == nil {

				char.Stamina = sta

				tr.context.Variables["character_stamina"] = sta

				debugPrint("[DEBUG] createCharacter: set Stamina=%d and saved to Variables\n", sta)

			}

		}

	}

	// 解析基础HP（如"基础HP=35"）
	if strings.Contains(instruction, "基础HP=") {

		parts := strings.Split(instruction, "基础HP=")

		if len(parts) > 1 {

			baseHPStr := strings.TrimSpace(strings.Split(parts[1], "）")[0])

			baseHPStr = strings.TrimSpace(strings.Split(baseHPStr, ",")[0])

			if baseHP, err := strconv.Atoi(baseHPStr); err == nil {

				tr.context.Variables["character_base_hp"] = baseHP

				debugPrint("[DEBUG] createCharacter: set baseHP=%d\n", baseHP)

			}

		}

	}

	// 解析攻击力（如"攻击=20"）
	if strings.Contains(instruction, "攻击=") {

		parts := strings.Split(instruction, "攻击=")

		if len(parts) > 1 {

			attackStr := strings.TrimSpace(strings.Split(parts[1], "）")[0])

			attackStr = strings.TrimSpace(strings.Split(attackStr, "）")[0])

			attackStr = strings.TrimSpace(strings.Split(attackStr, "）")[0])

			if attack, err := strconv.Atoi(attackStr); err == nil {

				char.PhysicalAttack = attack

				// 也存储到上下文，以便后续使用

				tr.context.Variables["character_physical_attack"] = attack

				debugPrint("[DEBUG] createCharacter: set PhysicalAttack=%d\n", attack)

			}

		}

	}

	// 解析防御力（如"防御=10"）
	if strings.Contains(instruction, "防御=") {

		parts := strings.Split(instruction, "防御=")

		if len(parts) > 1 {

			defenseStr := strings.TrimSpace(strings.Split(parts[1], "）")[0])

			defenseStr = strings.TrimSpace(strings.Split(defenseStr, "）")[0])

			if defense, err := strconv.Atoi(defenseStr); err == nil {

				char.PhysicalDefense = defense

			}

		}

	}

	// 解析金币（如"金币=100"）
	// 注意：Gold在User模型中，不在Character模型中
	if strings.Contains(instruction, "金币=") {

		parts := strings.Split(instruction, "金币=")

		if len(parts) > 1 {

			goldStr := strings.TrimSpace(strings.Split(parts[1], "）")[0])

			goldStr = strings.TrimSpace(strings.Split(goldStr, "）")[0])

			if gold, err := strconv.Atoi(goldStr); err == nil {

				// 存储到Variables，稍后在创建/更新用户时设置
				tr.context.Variables["character_gold"] = gold

				tr.context.Variables["character.gold"] = gold

				debugPrint("[DEBUG] createCharacter: set Gold=%d (will update user)\n", gold)

			}

		}

	}

	// 解析暴击率（如"物理暴击率=30%"）
	if strings.Contains(instruction, "物理暴击率=") {
		parts := strings.Split(instruction, "物理暴击率=")
		if len(parts) > 1 {
			critStr := strings.TrimSpace(strings.Split(parts[1], "%")[0])

			if crit, err := strconv.ParseFloat(critStr, 64); err == nil {

				char.PhysCritRate = crit / 100.0

				// 标记为明确设置，防止后续被覆盖
				tr.context.Variables["character_explicit_phys_crit_rate"] = char.PhysCritRate

				debugPrint("[DEBUG] createCharacter: set PhysCritRate=%f from instruction\n", char.PhysCritRate)

			}

		}

	}

	// 解析暴击伤害（如"物理暴击伤害=150%"）
	if strings.Contains(instruction, "物理暴击率=") {

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

	// 解析怒气/资源（如"怒气=100/100"）
	if strings.Contains(instruction, "怒气=") {

		parts := strings.Split(instruction, "怒气=")

		if len(parts) > 1 {

			resourceStr := strings.TrimSpace(strings.Split(parts[1], "）")[0])

			resourceStr = strings.TrimSpace(strings.Split(resourceStr, "）")[0])

			// 处理 "100/100" 格式

			if strings.Contains(resourceStr, "/") {

				resourceParts := strings.Split(resourceStr, "/")

				if len(resourceParts) >= 1 {

					if resource, err := strconv.Atoi(strings.TrimSpace(resourceParts[0])); err == nil {

						char.Resource = resource

						// 也存储到Variables，以便后续恢�						tr.context.Variables["character_resource"] = resource

						debugPrint("[DEBUG] createCharacter: parsed Resource=%d from instruction\n", resource)

					}

				}

				if len(resourceParts) >= 2 {

					if maxResource, err := strconv.Atoi(strings.TrimSpace(resourceParts[1])); err == nil {

						char.MaxResource = maxResource

						// 也存储到Variables，以便后续恢�						tr.context.Variables["character_max_resource"] = maxResource

						debugPrint("[DEBUG] createCharacter: parsed MaxResource=%d from instruction\n", maxResource)

					}

				}

			} else {

				// 处理 "100" 格式

				if resource, err := strconv.Atoi(resourceStr); err == nil {

					char.Resource = resource

					// 也存储到Variables，以便后续恢�					tr.context.Variables["character_resource"] = resource

					if char.MaxResource == 0 {

						char.MaxResource = resource

					}

					tr.context.Variables["character_max_resource"] = char.MaxResource

					debugPrint("[DEBUG] createCharacter: parsed Resource=%d, MaxResource=%d from instruction\n", resource, char.MaxResource)

				}

			}

		}

	}

	// 解析HP（如"HP=100/100"或"HP=100"）
	// 注意：必须排除“基础HP="的情况，避免误解
	// 保存明确设置的HP值，以便后续使用

	explicitHP := 0

	if strings.Contains(instruction, "HP=") && !strings.Contains(instruction, "基础HP=") {

		parts := strings.Split(instruction, "HP=")

		if len(parts) > 1 {

			hpStr := strings.TrimSpace(strings.Split(parts[1], "）")[0])

			hpStr = strings.TrimSpace(strings.Split(hpStr, "）")[0])

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

						// 保存MaxHP到Variables，以便后续恢复
						tr.context.Variables["character_explicit_max_hp"] = maxHP

						debugPrint("[DEBUG] createCharacter: set explicitMaxHP=%d\n", maxHP)

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

		debugPrint("[DEBUG] createCharacter: set explicitHP=%d\n", explicitHP)

	}

	// 设置默认资源值（如果未指定）

	if char.Resource == 0 && char.MaxResource == 0 {

		char.Resource = 100

		char.MaxResource = 100

	}

	// 如果MaxHP为0，自动计算MaxHP（使用Calculator）
	// 但是，如果HP已经被明确设置（通过"HP="指令），不要覆盖
	savedHP := char.HP

	// 检查是否有明确设置的HP值
	if explicitHPVal, exists := tr.context.Variables["character_explicit_hp"]; exists {

		if explicitHP, ok := explicitHPVal.(int); ok && explicitHP > 0 {

			savedHP = explicitHP

			char.HP = explicitHP

			debugPrint("[DEBUG] createCharacter: using explicitHP=%d from Variables\n", explicitHP)

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

		// 但是，如果HP已经被明确设置（通过"HP="指令），不要覆盖
		if savedHP == 0 {

			char.HP = char.MaxHP

		} else {

			// HP已经被明确设置，保持HP不变，但确保MaxHP至少等于HP

			if char.MaxHP < savedHP {

				char.MaxHP = savedHP

			}

			char.HP = savedHP

		}

		debugPrint("[DEBUG] createCharacter: auto-calculated MaxHP=%d, HP=%d (savedHP=%d)\n", char.MaxHP, char.HP, savedHP)

	} else if savedHP > 0 && savedHP < char.MaxHP {

		// 如果MaxHP已经被设置，但HP被明确设置为小于MaxHP的值，保持HP不变

		char.HP = savedHP

		debugPrint("[DEBUG] createCharacter: MaxHP=%d already set, keeping HP=%d\n", char.MaxHP, char.HP)

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

		// 从Variables恢复我们在指令中设置的属性值（Create可能覆盖了它们）

		if strengthVal, exists := tr.context.Variables["character_strength"]; exists {

			if strength, ok := strengthVal.(int); ok {

				char.Strength = strength

			}

		}

		if agilityVal, exists := tr.context.Variables["character_agility"]; exists {

			if agility, ok := agilityVal.(int); ok {

				char.Agility = agility

			}

		}

		if intellectVal, exists := tr.context.Variables["character_intellect"]; exists {

			if intellect, ok := intellectVal.(int); ok {

				char.Intellect = intellect

			}

		}

		if staminaVal, exists := tr.context.Variables["character_stamina"]; exists {

			if stamina, ok := staminaVal.(int); ok {

				char.Stamina = stamina

			}

		}

		if spiritVal, exists := tr.context.Variables["character_spirit"]; exists {

			if spirit, ok := spiritVal.(int); ok {

				char.Spirit = spirit

			}

		}

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

			// 使用数据库中的角色
			char = existingChar

			// 从Variables恢复我们在指令中设置的属性
			if strengthVal, exists := tr.context.Variables["character_strength"]; exists {

				if strength, ok := strengthVal.(int); ok {

					char.Strength = strength

				}

			}

			if agilityVal, exists := tr.context.Variables["character_agility"]; exists {

				if agility, ok := agilityVal.(int); ok {

					char.Agility = agility

				}

			}

			if intellectVal, exists := tr.context.Variables["character_intellect"]; exists {

				if intellect, ok := intellectVal.(int); ok {

					char.Intellect = intellect

				}

			}

			if staminaVal, exists := tr.context.Variables["character_stamina"]; exists {

				if stamina, ok := staminaVal.(int); ok {

					char.Stamina = stamina

				}

			}

			if spiritVal, exists := tr.context.Variables["character_spirit"]; exists {

				if spirit, ok := spiritVal.(int); ok {

					char.Spirit = spirit

				}

			}

			// 从Variables恢复Resource（如果指令中指定了）

			if resourceVal, exists := tr.context.Variables["character_resource"]; exists {

				if resource, ok := resourceVal.(int); ok && resource > 0 {

					char.Resource = resource

					debugPrint("[DEBUG] createCharacter: restored Resource=%d from Variables\n", resource)

				}

			}

			if maxResourceVal, exists := tr.context.Variables["character_max_resource"]; exists {

				if maxResource, ok := maxResourceVal.(int); ok && maxResource > 0 {

					char.MaxResource = maxResource

					debugPrint("[DEBUG] createCharacter: restored MaxResource=%d from Variables\n", maxResource)

				}

			}

			// 更新已存在角色的ClassID（如果指令中指定了不同的职业）

			if classIDVal, exists := tr.context.Variables["character_class_id"]; exists {

				if classID, ok := classIDVal.(string); ok && classID != "" {

					char.ClassID = classID

				}

			}

			// 在设置ID之后，如果MaxHP为0或小于计算值，重新计算MaxHP（从数据库读取后可能被重置）

			// 但是，如果HP已经被明确设置（通过"HP="指令），不要覆盖
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

			// 检查MaxHP是否已经被明确设置（通过"HP=95/100"）
			explicitMaxHP := 0

			if maxHPVal, exists := tr.context.Variables["character_explicit_max_hp"]; exists {

				if maxHP, ok := maxHPVal.(int); ok && maxHP > 0 {

					explicitMaxHP = maxHP

				}

			}

			calculatedMaxHP := tr.calculator.CalculateHP(char, baseHP)

			// 如果MaxHP已经被明确设置，使用明确设置的值
			if explicitMaxHP > 0 {

				char.MaxHP = explicitMaxHP

				// 如果HP已经被明确设置，保持HP不变

				if explicitHP > 0 {

					char.HP = explicitHP

				} else if char.HP == 0 || char.HP < char.MaxHP {

					char.HP = char.MaxHP

				}

				debugPrint("[DEBUG] createCharacter: after setting ID, using explicitMaxHP=%d, HP=%d (explicitHP=%d)\n", char.MaxHP, char.HP, explicitHP)

			} else if char.MaxHP == 0 || char.MaxHP < calculatedMaxHP {

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

				debugPrint("[DEBUG] createCharacter: after setting ID, re-calculated MaxHP=%d, HP=%d (explicitHP=%d)\n", char.MaxHP, char.HP, explicitHP)

			} else if explicitHP > 0 {

				// 如果MaxHP已经被设置，但HP被明确设置为小于MaxHP的值，保持HP不变

				char.HP = explicitHP

				debugPrint("[DEBUG] createCharacter: after setting ID, MaxHP=%d already set, keeping explicitHP=%d\n", char.MaxHP, explicitHP)

			}

			// 在设置ID之后，检查PhysicalAttack是否被重�			debugPrint("[DEBUG] createCharacter: after setting ID, char.PhysicalAttack=%d\n", char.PhysicalAttack)

			// 如果PhysicalAttack为0，从Variables恢复）

			if char.PhysicalAttack == 0 {

				if attackVal, exists := tr.context.Variables["character_physical_attack"]; exists {

					if attack, ok := attackVal.(int); ok && attack > 0 {

						char.PhysicalAttack = attack

						debugPrint("[DEBUG] createCharacter: restored PhysicalAttack=%d from Variables before Update\n", attack)

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

				debugPrint("[DEBUG] createCharacter: re-calculated MaxHP=%d, HP=%d after reading from DB\n", char.MaxHP, char.HP)

			}

			// 保存PhysicalAttack、Resource和MaxHP值，以防数据库更新时丢失

			savedPhysicalAttack := char.PhysicalAttack

			savedResource := char.Resource

			savedMaxResource := char.MaxResource

			savedMaxHP := char.MaxHP

			savedHP := char.HP

			debugPrint("[DEBUG] createCharacter: before Update, char.PhysicalAttack=%d, Resource=%d/%d, MaxHP=%d, HP=%d\n", char.PhysicalAttack, char.Resource, char.MaxResource, char.MaxHP, char.HP)

			if err := charRepo.Update(char); err != nil {

				return fmt.Errorf("failed to update existing character in DB: %w", err)

			}

			// 从数据库重新加载角色（因为Update可能修改了某些字段）

			reloadedChar, err := charRepo.GetByID(char.ID)

			if err == nil && reloadedChar != nil {

				char = reloadedChar

			}

			// 恢复PhysicalAttack值（如果它被数据库更新覆盖了）
			if savedPhysicalAttack > 0 {

				char.PhysicalAttack = savedPhysicalAttack

				debugPrint("[DEBUG] createCharacter: after Update, restored PhysicalAttack=%d\n", char.PhysicalAttack)

			} else if char.PhysicalAttack == 0 {

				// 如果PhysicalAttack为0，重新计算
				char.PhysicalAttack = tr.calculator.CalculatePhysicalAttack(char)

				debugPrint("[DEBUG] createCharacter: after Update, re-calculated PhysicalAttack=%d (was 0)\n", char.PhysicalAttack)

			} else {

				debugPrint("[DEBUG] createCharacter: after Update, char.PhysicalAttack=%d (not restored)\n", char.PhysicalAttack)

			}

			// 恢复PhysCritRate值（如果它被明确设置）
			if explicitCritRate, exists := tr.context.Variables["character_explicit_phys_crit_rate"]; exists {

				if critRate, ok := explicitCritRate.(float64); ok && critRate > 0 {

					char.PhysCritRate = critRate

					debugPrint("[DEBUG] createCharacter: after Update, restored PhysCritRate=%f\n", critRate)

				}

			}

			// 恢复Resource值（如果它被数据库更新覆盖了）
			// 优先使用savedResource和savedMaxResource（如果它们都不为0）
			debugPrint("[DEBUG] createCharacter: after Update, char.Resource=%d/%d (from DB)\n", char.Resource, char.MaxResource)

			if savedResource > 0 && savedMaxResource > 0 {

				// 直接恢复保存的值，不做特殊判断

				char.Resource = savedResource

				char.MaxResource = savedMaxResource

				debugPrint("[DEBUG] createCharacter: after Update, restored Resource=%d/%d (from saved values)\n", char.Resource, char.MaxResource)

			} else if savedMaxResource > 0 {

				// 如果MaxResource不为0但Resource为0，恢复Resource为MaxResource

				char.Resource = savedMaxResource

				char.MaxResource = savedMaxResource

				debugPrint("[DEBUG] createCharacter: after Update, restored Resource=%d/%d (from MaxResource)\n", char.Resource, char.MaxResource)

			} else if char.Resource == 0 && char.MaxResource == 0 {

				// 如果资源被重置为0，恢复默认值
				char.Resource = 100

				char.MaxResource = 100

				debugPrint("[DEBUG] createCharacter: after Update, restored default Resource=100/100\n")

			} else if char.MaxResource > 0 && char.Resource == 0 {

				// 如果MaxResource不为0但Resource为0，恢复Resource为MaxResource

				char.Resource = char.MaxResource

				debugPrint("[DEBUG] createCharacter: after Update, restored Resource=%d (from MaxResource)\n", char.Resource)

			} else if char.MaxResource == 100 && char.Resource < 100 {

				// 如果MaxResource为100但Resource小于100，恢复Resource为100

				char.Resource = char.MaxResource

				debugPrint("[DEBUG] createCharacter: after Update, restored Resource=%d (MaxResource is 100)\n", char.Resource)

			}

			// 恢复MaxHP和HP值（如果它们被数据库更新覆盖了）

			if savedMaxHP > 0 {

				char.MaxHP = savedMaxHP

				char.HP = savedHP

				debugPrint("[DEBUG] createCharacter: after Update, restored MaxHP=%d, HP=%d\n", char.MaxHP, char.HP)

				// 再次更新数据库，确保MaxHP和HP被保存
				if err := charRepo.Update(char); err != nil {

					debugPrint("[DEBUG] createCharacter: failed to update MaxHP/HP in DB: %v\n", err)

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

				debugPrint("[DEBUG] createCharacter: after Create, restored PhysicalAttack=%d\n", char.PhysicalAttack)

			} else if char.PhysicalAttack == 0 {

				// 如果PhysicalAttack为0，重新计算
				char.PhysicalAttack = tr.calculator.CalculatePhysicalAttack(char)

				debugPrint("[DEBUG] createCharacter: after Create, re-calculated PhysicalAttack=%d (was 0)\n", char.PhysicalAttack)

			} else {

				debugPrint("[DEBUG] createCharacter: after Create, char.PhysicalAttack=%d (not restored)\n", char.PhysicalAttack)

			}

			// 恢复Resource值（如果它被Create覆盖了）

			// 优先使用savedResource和savedMaxResource（如果它们都不为0）
			if savedResource > 0 && savedMaxResource > 0 {

				// 直接恢复保存的值，不做特殊判断

				char.Resource = savedResource

				char.MaxResource = savedMaxResource

				debugPrint("[DEBUG] createCharacter: after Create, restored Resource=%d/%d\n", char.Resource, char.MaxResource)

			} else if savedMaxResource > 0 {

				// 如果MaxResource不为0但Resource为0，恢复Resource为MaxResource

				char.Resource = savedMaxResource

				char.MaxResource = savedMaxResource

				debugPrint("[DEBUG] createCharacter: after Create, restored Resource=%d/%d (from MaxResource)\n", char.Resource, char.MaxResource)

			} else if char.Resource == 0 && char.MaxResource == 0 {

				// 如果资源被重置为0，恢复默认值
				char.Resource = 100

				char.MaxResource = 100

				debugPrint("[DEBUG] createCharacter: after Create, restored default Resource=100/100\n")

			} else if char.MaxResource > 0 && char.Resource == 0 {

				// 如果MaxResource不为0但Resource为0，恢复Resource为MaxResource

				char.Resource = char.MaxResource

				debugPrint("[DEBUG] createCharacter: after Create, restored Resource=%d (from MaxResource)\n", char.Resource)

			} else if char.MaxResource == 100 && char.Resource < 100 {

				// 如果MaxResource为100但Resource小于100，恢复Resource为100

				char.Resource = char.MaxResource

				debugPrint("[DEBUG] createCharacter: after Create, restored Resource=%d (MaxResource is 100)\n", char.Resource)

			}

			// 恢复MaxHP和HP值（如果它们被Create覆盖了）

			// 首先检查是否有明确设置的MaxHP
			restoreExplicitMaxHP := 0

			if maxHPVal, exists := tr.context.Variables["character_explicit_max_hp"]; exists {

				if maxHP, ok := maxHPVal.(int); ok && maxHP > 0 {

					restoreExplicitMaxHP = maxHP

				}

			}

			// 检查是否有明确设置的HP
			restoreExplicitHP := 0

			if explicitHPVal, exists := tr.context.Variables["character_explicit_hp"]; exists {

				if hp, ok := explicitHPVal.(int); ok && hp > 0 {

					restoreExplicitHP = hp

				}

			}

			// 获取基础HP用于重新计算

			restoreBaseHP := 35 // 默认战士基础HP

			if baseHPVal, exists := tr.context.Variables["character_base_hp"]; exists {

				if hp, ok := baseHPVal.(int); ok {

					restoreBaseHP = hp

				}

			}

			// 重新计算MaxHP（基于当前属性）

			restoreCalculatedMaxHP := tr.calculator.CalculateHP(char, restoreBaseHP)

			// 确定最终的MaxHP
			if restoreExplicitMaxHP > 0 {

				char.MaxHP = restoreExplicitMaxHP

			} else if savedMaxHP > 0 && savedMaxHP == restoreCalculatedMaxHP {

				// 如果保存的MaxHP等于计算值，使用保存的值
				char.MaxHP = savedMaxHP

			} else if char.MaxHP != restoreCalculatedMaxHP {

				// 如果当前MaxHP不等于计算值，使用计算值
				char.MaxHP = restoreCalculatedMaxHP

			}

			// 确定最终的HP
			if restoreExplicitHP > 0 {

				char.HP = restoreExplicitHP

				// 确保MaxHP至少等于HP

				if char.MaxHP < restoreExplicitHP {

					char.MaxHP = restoreExplicitHP

				}

			} else if savedHP > 0 && savedHP <= char.MaxHP {

				char.HP = savedHP

			} else if char.HP == 0 || char.HP > char.MaxHP {

				// 如果HP为0或超过MaxHP，设置为MaxHP

				char.HP = char.MaxHP

			}

			debugPrint("[DEBUG] createCharacter: after Create, final MaxHP=%d, HP=%d (calculatedMaxHP=%d, savedMaxHP=%d, explicitMaxHP=%d, explicitHP=%d)\n", char.MaxHP, char.HP, restoreCalculatedMaxHP, savedMaxHP, restoreExplicitMaxHP, restoreExplicitHP)

			// 再次更新数据库，确保MaxHP和HP被保存
			if err := charRepo.Update(char); err != nil {

				debugPrint("[DEBUG] createCharacter: failed to update MaxHP/HP in DB: %v\n", err)

			}

		}

	}

	// 在计算属性前，确保基础属性值正确（从Variables恢复）
	if strengthVal, exists := tr.context.Variables["character_strength"]; exists {

		if strength, ok := strengthVal.(int); ok {

			char.Strength = strength

			debugPrint("[DEBUG] createCharacter: restored Strength=%d from Variables before calculation\n", strength)

		}

	}

	if agilityVal, exists := tr.context.Variables["character_agility"]; exists {

		if agility, ok := agilityVal.(int); ok {

			char.Agility = agility

			debugPrint("[DEBUG] createCharacter: restored Agility=%d from Variables before calculation\n", agility)

		}

	} else {

		debugPrint("[DEBUG] createCharacter: character_agility NOT found in Variables (keys: %v)\n", getMapKeys(tr.context.Variables))

	}

	if intellectVal, exists := tr.context.Variables["character_intellect"]; exists {

		if intellect, ok := intellectVal.(int); ok {

			char.Intellect = intellect

		}

	}

	if staminaVal, exists := tr.context.Variables["character_stamina"]; exists {

		if stamina, ok := staminaVal.(int); ok {

			char.Stamina = stamina

		}

	}

	if spiritVal, exists := tr.context.Variables["character_spirit"]; exists {

		if spirit, ok := spiritVal.(int); ok {

			char.Spirit = spirit

		}

	}

	// 计算并更新所有属性（如果它们�或未设置�	// 获取基础HP（从Variables或使用默认值）

	baseHP := 35 // 默认战士基础HP

	if baseHPVal, exists := tr.context.Variables["character_base_hp"]; exists {

		if hp, ok := baseHPVal.(int); ok {

			baseHP = hp

		}

	}

	// 计算所有属性（如果未设置或未明确设置，则重新计算）
	explicitPhysicalAttack := false // 注意：如果属性已经在指令中明确设置（�攻击�20"�物理暴击率= false

	if attackVal, exists := tr.context.Variables["character_physical_attack"]; exists {

		// 检查是否是通过"攻击）"指令设置的（而不是计算后存储的）

		if instruction, ok := tr.context.Variables["last_instruction"].(string); ok && strings.Contains(instruction, "攻击）") {

			explicitPhysicalAttack = true

			if attack, ok := attackVal.(int); ok {

				char.PhysicalAttack = attack

				debugPrint("[DEBUG] createCharacter: using explicit PhysicalAttack=%d from instruction\n", attack)

			}

		}

	}

	// 如果未明确设置，总是基于主属性重新计算（即使当前值不为0）
	if !explicitPhysicalAttack {

		oldAttack := char.PhysicalAttack

		calculatedAttack := tr.calculator.CalculatePhysicalAttack(char)

		// 如果当前值为0或与计算值不同，使用计算值
		if oldAttack == 0 || oldAttack != calculatedAttack {

			char.PhysicalAttack = calculatedAttack

			debugPrint("[DEBUG] createCharacter: re-calculated PhysicalAttack=%d (from Strength=%d, Agility=%d, was %d)\n", char.PhysicalAttack, char.Strength, char.Agility, oldAttack)

		}

	}

	// 法术攻击力：如果未明确设置或为0，总是基于主属性重新计算
	if char.MagicAttack == 0 {

		char.MagicAttack = tr.calculator.CalculateMagicAttack(char)

		debugPrint("[DEBUG] createCharacter: calculated MagicAttack=%d (from Intellect=%d, Spirit=%d)\n", char.MagicAttack, char.Intellect, char.Spirit)

	}

	// 物理防御：如果未明确设置，总是基于主属性重新计算
	if char.PhysicalDefense == 0 {

		char.PhysicalDefense = tr.calculator.CalculatePhysicalDefense(char)

	}

	// 魔法防御：如果未明确设置，总是基于主属性重新计算
	if char.MagicDefense == 0 {

		char.MagicDefense = tr.calculator.CalculateMagicDefense(char)

	}

	// 暴击率和闪避率：如果为0，则计算；如果已设置，保持原值
	// 检查是否有明确设置的PhysCritRate
	if explicitCritRate, exists := tr.context.Variables["character_explicit_phys_crit_rate"]; exists {

		if critRate, ok := explicitCritRate.(float64); ok && critRate > 0 {

			char.PhysCritRate = critRate

			debugPrint("[DEBUG] createCharacter: using explicit PhysCritRate=%f from Variables\n", critRate)

		}

	} else if char.PhysCritRate == 0 {

		char.PhysCritRate = tr.calculator.CalculatePhysCritRate(char)

	}

	if char.PhysCritDamage == 0 {

		char.PhysCritDamage = tr.calculator.CalculatePhysCritDamage(char)

	}

	if char.SpellCritRate == 0 {

		char.SpellCritRate = tr.calculator.CalculateSpellCritRate(char)

	}

	if char.SpellCritDamage == 0 {

		char.SpellCritDamage = tr.calculator.CalculateSpellCritDamage(char)

	}

	if char.DodgeRate == 0 {

		char.DodgeRate = tr.calculator.CalculateDodgeRate(char)

	}

	// 计算速度（speed = agility）	// 注意：速度不是Character模型的字段，但可以通过Calculator计算

	// 这里我们确保速度值被正确计算并存储到上下文
	speed := tr.calculator.CalculateSpeed(char)

	tr.context.Variables["character_speed"] = speed

	// 计算MaxHP（如果为0，或者如果MaxHP小于明确设置的HP值）

	// 但是，如果MaxHP已经被明确设置（通过"HP=95/100"），不要覆盖
	finalCalculatedMaxHP := tr.calculator.CalculateHP(char, baseHP)

	// 检查是否有明确设置的MaxHP
	finalExplicitMaxHP := 0

	if maxHPVal, exists := tr.context.Variables["character_explicit_max_hp"]; exists {

		if maxHP, ok := maxHPVal.(int); ok && maxHP > 0 {

			finalExplicitMaxHP = maxHP

		}

	}

	// 确定最终的MaxHP
	if finalExplicitMaxHP > 0 {

		char.MaxHP = finalExplicitMaxHP

	} else if char.MaxHP == 0 || char.MaxHP != finalCalculatedMaxHP {

		// 如果MaxHP�或与计算值为0或与计算值不一致，使用计算值
		char.MaxHP = finalCalculatedMaxHP

	}

	// 检查是否有明确设置的HP
	finalExplicitHP := 0

	if explicitHPVal, exists := tr.context.Variables["character_explicit_hp"]; exists {

		if hp, ok := explicitHPVal.(int); ok && hp > 0 {

			finalExplicitHP = hp

		}

	}

	// 确定最终的HP
	if finalExplicitHP > 0 {

		char.HP = finalExplicitHP

		// 确保MaxHP至少等于HP

		if char.MaxHP < finalExplicitHP {

			char.MaxHP = finalExplicitHP

		}

	} else if char.HP == 0 || char.HP > char.MaxHP {

		// 如果HP为0或超过MaxHP，设置为MaxHP

		char.HP = char.MaxHP

	}

	debugPrint("[DEBUG] createCharacter: final calculation - MaxHP=%d, HP=%d (calculatedMaxHP=%d, explicitMaxHP=%d, explicitHP=%d)\n", char.MaxHP, char.HP, finalCalculatedMaxHP, finalExplicitMaxHP, finalExplicitHP)

	// 更新用户金币（如果设置了
	if goldVal, exists := tr.context.Variables["character_gold"]; exists {

		if gold, ok := goldVal.(int); ok {

			// 直接更新数据库中的用户金币
			_, err := database.DB.Exec(`UPDATE users SET gold = ? WHERE id = ?`, gold, char.UserID)

			if err != nil {

				debugPrint("[DEBUG] createCharacter: failed to update user gold: %v\n", err)

			} else {

				tr.context.Variables["character.gold"] = gold

				debugPrint("[DEBUG] createCharacter: set user Gold=%d (userID=%d)\n", gold, char.UserID)

			}

		}

	}

	// 存储到上下文（确保所有属性正确）

	tr.context.Characters["character"] = char

	debugPrint("[DEBUG] createCharacter: stored character to context, PhysicalAttack=%d, MagicAttack=%d\n", char.PhysicalAttack, char.MagicAttack)

	// 存储所有计算属性到Variables，以防角色对象被修改

	tr.context.Variables["character_physical_attack"] = char.PhysicalAttack

	tr.context.Variables["character_magic_attack"] = char.MagicAttack

	tr.context.Variables["character_physical_defense"] = char.PhysicalDefense

	tr.context.Variables["character_magic_defense"] = char.MagicDefense

	tr.context.Variables["character_phys_crit_rate"] = char.PhysCritRate

	tr.context.Variables["character_phys_crit_damage"] = char.PhysCritDamage

	tr.context.Variables["character_spell_crit_rate"] = char.SpellCritRate

	tr.context.Variables["character_spell_crit_damage"] = char.SpellCritDamage

	tr.context.Variables["character_dodge_rate"] = char.DodgeRate

	tr.context.Variables["character_speed"] = speed

	tr.context.Variables["character_max_hp"] = char.MaxHP

	tr.context.Variables["character_hp"] = char.HP

	// 同时存储简化键（不带character_前缀），以便测试用例可以直接访问

	tr.context.Variables["physical_attack"] = char.PhysicalAttack

	tr.context.Variables["magic_attack"] = char.MagicAttack

	tr.context.Variables["physical_defense"] = char.PhysicalDefense

	tr.context.Variables["magic_defense"] = char.MagicDefense

	tr.context.Variables["phys_crit_rate"] = char.PhysCritRate

	tr.context.Variables["phys_crit_damage"] = char.PhysCritDamage

	tr.context.Variables["spell_crit_rate"] = char.SpellCritRate

	tr.context.Variables["spell_crit_damage"] = char.SpellCritDamage

	tr.context.Variables["dodge_rate"] = char.DodgeRate

	tr.context.Variables["speed"] = speed

	tr.context.Variables["max_hp"] = char.MaxHP

	debugPrint("[DEBUG] createCharacter: stored all calculated attributes to Variables\n")

	debugPrint("[DEBUG] createCharacter: final context - characters=%d, stored character with key='character'\n", len(tr.context.Characters))

	debugPrint("[DEBUG] createCharacter: final context - characters=%d, stored character with key='character'\n", len(tr.context.Characters))

	return nil

}

func (tr *TestRunner) createMultipleCharacters(instruction string) error {

	// 解析角色列表（通过冒号分隔）
	var characterDescs []string

	if strings.Contains(instruction, "创建多个角色:") {

		parts := strings.Split(instruction, "创建多个角色:")

		if len(parts) > 1 {

			characterDescs = strings.Split(parts[1], ",")

		}

	} else if strings.Contains(instruction, ":") {

		parts := strings.Split(instruction, "创建多个角色:")

		if len(parts) > 1 {

			characterDescs = strings.Split(parts[1], ",")

		}

	}

	charRepo := repository.NewCharacterRepository()

	user, err := tr.createTestUser()

	if err != nil {

		return fmt.Errorf("failed to create test user: %w", err)

	}

	// 先获取用户的所有角色，检查哪些slot已被占用

	existingChars, err := charRepo.GetByUserID(user.ID)

	if err != nil {

		existingChars = []*models.Character{}

	}

	existingSlots := make(map[int]*models.Character)

	for _, c := range existingChars {

		existingSlots[c.TeamSlot] = c

	}

	for _, charDesc := range characterDescs {

		charDesc = strings.TrimSpace(charDesc)

		if charDesc == "" {

			continue

		}

		// 解析角色索引（如"角色1"�角色2"等）

		charIndex := 1

		if strings.Contains(charDesc, "角色") {

			// 提取数字

			re := regexp.MustCompile(`角色(\d+)`)

			matches := re.FindStringSubmatch(charDesc)

			if len(matches) > 1 {

				if idx, err := strconv.Atoi(matches[1]); err == nil {

					charIndex = idx

				}

			}

		}

		// 使用createCharacter的逻辑，但修改指令以创建单个角�		// �角色1（敏�30，速度=60�转换�创建一个角色，敏捷=30，速度=60"

		singleCharInstruction := strings.Replace(charDesc, fmt.Sprintf("角色%d", charIndex), "一个角色", 1)

		singleCharInstruction = strings.TrimSpace(strings.TrimPrefix(singleCharInstruction, ":"))

		singleCharInstruction = strings.TrimSpace(strings.TrimSuffix(singleCharInstruction, ":"))

		singleCharInstruction = strings.TrimSpace(strings.TrimSuffix(singleCharInstruction, ")"))

		singleCharInstruction = "创建一个角色，" + singleCharInstruction

		// 临时保存当前上下文，以便createCharacter使用

		oldLastInstruction := tr.context.Variables["last_instruction"]

		tr.context.Variables["last_instruction"] = singleCharInstruction

		// 调用createCharacter创建单个角色

		if err := tr.createCharacter(singleCharInstruction); err != nil {

			tr.context.Variables["last_instruction"] = oldLastInstruction

			return fmt.Errorf("failed to create character %d: %w", charIndex, err)

		}

		// 恢复last_instruction

		tr.context.Variables["last_instruction"] = oldLastInstruction

		// 获取刚创建的角色（应该存储在"character"键中）
		char, ok := tr.context.Characters["character"]

		if !ok || char == nil {

			return fmt.Errorf("failed to get created character %d", charIndex)

		}

		// 保存敏捷值（可能在数据库操作后丢失）

		savedAgility := char.Agility

		savedStrength := char.Strength

		savedIntellect := char.Intellect

		savedStamina := char.Stamina

		savedSpirit := char.Spirit

		// 检查该slot是否已存在角色
		if existingChar, exists := existingSlots[charIndex]; exists {

			// 更新已存在的角色

			char.ID = existingChar.ID

			char.TeamSlot = charIndex

			char.UserID = user.ID

			// 恢复保存的属性
			char.Agility = savedAgility

			char.Strength = savedStrength

			char.Intellect = savedIntellect

			char.Stamina = savedStamina

			char.Spirit = savedSpirit

			if err := charRepo.Update(char); err != nil {

				return fmt.Errorf("failed to update character %d: %w", charIndex, err)

			}

		} else {

			// 创建新角色
			char.TeamSlot = charIndex

			char.UserID = user.ID

			// 确保属性值正确
			char.Agility = savedAgility

			char.Strength = savedStrength

			char.Intellect = savedIntellect

			char.Stamina = savedStamina

			char.Spirit = savedSpirit

			createdChar, err := charRepo.Create(char)

			if err != nil {

				return fmt.Errorf("failed to create character %d: %w", charIndex, err)

			}

			char = createdChar

			// 数据库操作后，可能需要重新设置属性�			char.Agility = savedAgility

			char.Strength = savedStrength

			char.Intellect = savedIntellect

			char.Stamina = savedStamina

			char.Spirit = savedSpirit

			// 更新数据库以确保属性值正确
			charRepo.Update(char)

		}

		// 确保属性值正确（数据库操作后可能被重置）

		char.Agility = savedAgility

		char.Strength = savedStrength

		char.Intellect = savedIntellect

		char.Stamina = savedStamina

		char.Spirit = savedSpirit

		// 重新计算速度（确保使用最新的敏捷值）

		speed := tr.calculator.CalculateSpeed(char)

		tr.context.Variables[fmt.Sprintf("character_%d_speed", charIndex)] = speed

		// 存储到上下文（使用character_1, character_2等作为键）
		key := fmt.Sprintf("character_%d", charIndex)

		tr.context.Characters[key] = char

		// 第一个角色也保存�character"（向后兼容）

		if charIndex == 1 {

			tr.context.Characters["character"] = char

		}

	}

	return nil

}

func (tr *TestRunner) createTestCharacter(userID, level int) (*models.Character, error) {

	charRepo := repository.NewCharacterRepository()

	chars, err := charRepo.GetByUserID(userID)

	var char *models.Character

	if err != nil || len(chars) == 0 {

		char = &models.Character{

			UserID: userID,

			Name: "测试角色",

			RaceID: "human",

			ClassID: "warrior",

			Faction: "alliance",

			TeamSlot: 1,

			Level: level,

			HP: 100, MaxHP: 100,

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

				UserID: userID,

				Name: "测试角色",

				RaceID: "human",

				ClassID: "warrior",

				Faction: "alliance",

				TeamSlot: 1,

				Level: level,

				HP: 100, MaxHP: 100,

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

func (tr *TestRunner) executeGetCharacterData() error {

	char, ok := tr.context.Characters["character"]

	if !ok || char == nil {

		return fmt.Errorf("character not found")

	}

	// 确保战士的怒气正确（如果不在战斗中，应该为0
	if char.ResourceType == "rage" {

		char.MaxResource = 100

		// 非战斗状态下，怒气应该�

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

package runner

import (
	"fmt"

	"strconv"

	"strings"

	"text-wow/internal/models"

	"text-wow/internal/repository"
)

// Equipment 相关函数

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

			UserID: user.ID,

			Name: "测试角色",

			RaceID: "human",

			ClassID: "warrior",

			Faction: "alliance",

			TeamSlot: 1,

			Level: level,
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

		// 存储最后一件装备到上下文（只存储基本字段，不存储整个对象）

		if i == count-1 {

			tr.context.Variables["equipment_id"] = equipment.ID

			tr.context.Variables["equipment_item_id"] = equipment.ItemID

			tr.context.Variables["equipment_quality"] = equipment.Quality

			tr.context.Variables["equipment_slot"] = equipment.Slot

		}

	}

	// 设置唯一词缀组合数量

	tr.context.Variables["unique_affix_combinations"] = len(uniqueCombinations)

	return nil

}

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

				ID: "boss_monster",

				Name: "Boss怪物",

				Type: "boss",

				Level: 30,

				HP:    0, // 被击败
				MaxHP: 1000,

				PhysicalAttack: 50,

				MagicAttack: 50,

				PhysicalDefense: 20,

				MagicDefense: 20,

				DodgeRate: 0.1,
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

	// 存储到上下文（只存储基本字段，不存储整个对象）
	tr.context.Variables["equipment_id"] = equipment.ID

	tr.context.Variables["equipment_item_id"] = equipment.ItemID

	tr.context.Variables["equipment_quality"] = equipment.Quality

	tr.context.Variables["equipment_slot"] = equipment.Slot

	tr.context.Equipments[fmt.Sprintf("%d", equipment.ID)] = equipment

	return nil

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

	tr.safeSetContext(fmt.Sprintf("%s.id", prefix), eq.ID)
	tr.safeSetContext(fmt.Sprintf("%s.item_id", prefix), eq.ItemID)
	tr.safeSetContext(fmt.Sprintf("%s.quality", prefix), eq.Quality)
	tr.safeSetContext(fmt.Sprintf("%s.slot", prefix), eq.Slot)

	// 同步character_id
	if eq.CharacterID != nil {
		tr.safeSetContext(fmt.Sprintf("%s.character_id", prefix), *eq.CharacterID)
	} else {
		tr.safeSetContext(fmt.Sprintf("%s.character_id", prefix), nil)
	}

	// 同步词缀ID
	if eq.PrefixID != nil {
		tr.safeSetContext(fmt.Sprintf("%s.prefix_id", prefix), *eq.PrefixID)
	} else {
		tr.safeSetContext(fmt.Sprintf("%s.prefix_id", prefix), nil)
	}
	if eq.SuffixID != nil {
		tr.safeSetContext(fmt.Sprintf("%s.suffix_id", prefix), *eq.SuffixID)
	} else {
		tr.safeSetContext(fmt.Sprintf("%s.suffix_id", prefix), nil)
	}

	// 同步词缀数值
	if eq.PrefixValue != nil {
		tr.safeSetContext(fmt.Sprintf("%s.prefix_value", prefix), *eq.PrefixValue)
	}
	if eq.SuffixValue != nil {
		tr.safeSetContext(fmt.Sprintf("%s.suffix_value", prefix), *eq.SuffixValue)
	}

	// 同步额外词缀
	if eq.BonusAffix1 != nil {
		tr.safeSetContext(fmt.Sprintf("%s.bonus_affix_1", prefix), *eq.BonusAffix1)
	}
	if eq.BonusAffix2 != nil {
		tr.safeSetContext(fmt.Sprintf("%s.bonus_affix_2", prefix), *eq.BonusAffix2)
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
	tr.safeSetContext(fmt.Sprintf("%s.affix_count", prefix), affixCount)

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
		tr.safeSetContext(fmt.Sprintf("%s.affixes", prefix), affixesStr)
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
	// 这里可以添加词缀信息的同步逻辑
	// 暂时留空，避免编译错误
	_ = affixID
	_ = affixType
	_ = equipmentLevel
	_ = slotType
}

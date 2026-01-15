package game

import (
	"fmt"
	"math/rand"
	"os"
	"sync"

	"text-wow/internal/database"
	"text-wow/internal/models"
	"text-wow/internal/repository"
)

// EquipmentManager 装备管理器 - 管理装备穿戴、词缀生成、装备强化
// 
// 注意：装备系统的完整实现需要 EquipmentRepository 支持。
// EquipmentRepository 需要实现以下方法：
//   - GetEquipmentByID(equipmentID int) (*EquipmentInstance, error)
//   - GetEquipmentsByCharacter(characterID int) ([]*EquipmentInstance, error)
//   - GetEquipmentByCharacterAndSlot(characterID int, slot string) (*EquipmentInstance, error)
//   - CreateEquipment(equipment *EquipmentInstance) error
//   - UpdateEquipment(equipment *EquipmentInstance) error
//   - DeleteEquipment(equipmentID int) error
//
// 装备系统的核心功能（词缀生成、装备强化等）已经实现，但穿戴/卸下功能需要 EquipmentRepository。
type EquipmentManager struct {
	mu          sync.RWMutex
	charRepo    *repository.CharacterRepository
	equipmentRepo *repository.EquipmentRepository
	gameRepo    *repository.GameRepository
	affixGenerator *AffixGenerator
	calculator  *Calculator
}

// AffixGenerator 词缀生成器
type AffixGenerator struct {
	mu          sync.RWMutex
	affixConfigs map[string]*AffixConfig // 词缀配置缓存
	gameRepo    *repository.GameRepository
	rng         *rand.Rand
}

// AffixConfig 词缀配置
type AffixConfig struct {
	ID          string
	Name        string
	Type        string  // prefix/suffix
	SlotType    string  // weapon/armor/accessory/all
	Rarity      string  // common/uncommon/rare/epic
	EffectType  string
	EffectStat  string
	MinValue    float64
	MaxValue    float64
	ValueType   string  // flat/percent
	LevelRequired int
	Tier        int
}

// NewEquipmentManager 创建装备管理器
func NewEquipmentManager() *EquipmentManager {
	return &EquipmentManager{
		charRepo:      repository.NewCharacterRepository(),
		equipmentRepo: repository.NewEquipmentRepository(),
		gameRepo:      repository.NewGameRepository(),
		affixGenerator: NewAffixGenerator(),
		calculator:    NewCalculator(),
	}
}

// NewAffixGenerator 创建词缀生成器
func NewAffixGenerator() *AffixGenerator {
	return &AffixGenerator{
		affixConfigs: make(map[string]*AffixConfig),
		gameRepo:     repository.NewGameRepository(),
		rng:          rand.New(rand.NewSource(rand.Int63())),
	}
}

// GenerateEquipment 生成装备（掉落时）
// 根据品质生成对应数量的词缀
func (em *EquipmentManager) GenerateEquipment(itemID string, quality string, level int, ownerID int) (*models.EquipmentInstance, error) {
	em.mu.Lock()
	defer em.mu.Unlock()

	// 确定词缀数量
	affixCount := em.getAffixCountByQuality(quality)

	// 生成词缀
	affixes := make([]*GeneratedAffix, 0, affixCount)
	// 从items表查询实际slot
	slot := "weapon" // 默认值
	itemData, err := em.gameRepo.GetItemByID(itemID)
	if err == nil {
		if slotValue, ok := itemData["slot"].(string); ok && slotValue != "" {
			slot = slotValue
		}
	}
	
	for i := 0; i < affixCount; i++ {
		affix, err := em.affixGenerator.GenerateAffix(slot, level, i == 0) // 第一个是前缀
		if err != nil {
			return nil, fmt.Errorf("failed to generate affix: %w", err)
		}
		affixes = append(affixes, affix)
	}

	// 创建装备实例
	equipment := &models.EquipmentInstance{
		ItemID:         itemID,
		OwnerID:        ownerID,
		CharacterID:    nil, // 初始在背包中
		Slot:           slot,
		Quality:        quality,
		EvolutionStage: 1,
		EvolutionPath:  nil,
		PrefixID:       nil,
		PrefixValue:    nil,
		SuffixID:       nil,
		SuffixValue:    nil,
		BonusAffix1:    nil,
		BonusAffix1Value: nil,
		BonusAffix2:    nil,
		BonusAffix2Value: nil,
		LegendaryEffectID: nil,
		IsLocked:       false,
	}

	// 设置词缀
	// 第一个词缀：前缀
	// 第二个词缀：后缀
	// 第三个词缀：额外词缀1（紫色+）
	// 第四个词缀：额外词缀2（橙色+）
	prefixSet := false
	suffixSet := false
	
	for i, affix := range affixes {
		if affix.Type == "prefix" && !prefixSet {
			equipment.PrefixID = &affix.ID
			equipment.PrefixValue = &affix.Value
			prefixSet = true
		} else if affix.Type == "suffix" && !suffixSet {
			equipment.SuffixID = &affix.ID
			equipment.SuffixValue = &affix.Value
			suffixSet = true
		} else if i == 2 {
			// 第三个词缀作为额外词缀1
			equipment.BonusAffix1 = &affix.ID
			equipment.BonusAffix1Value = &affix.Value
		} else if i == 3 {
			// 第四个词缀作为额外词缀2
			equipment.BonusAffix2 = &affix.ID
			equipment.BonusAffix2Value = &affix.Value
		}
	}
	
	// 如果前缀或后缀还没有设置，但还有词缀可用，使用它们
	if !prefixSet && len(affixes) > 0 {
		equipment.PrefixID = &affixes[0].ID
		equipment.PrefixValue = &affixes[0].Value
		prefixSet = true
	}
	if !suffixSet && len(affixes) > 1 {
		equipment.SuffixID = &affixes[1].ID
		equipment.SuffixValue = &affixes[1].Value
		suffixSet = true
	}

	// 保存到数据库
	return em.equipmentRepo.Create(equipment)
	/*
	// 获取基础物品配置
	item, err := em.gameRepo.GetItemByID(itemID)
	if err != nil {
		return nil, fmt.Errorf("failed to get item: %w", err)
	}

	// 确定词缀数量
	affixCount := em.getAffixCountByQuality(quality)

	// 生成词缀
	affixes := make([]*GeneratedAffix, 0, affixCount)
	for i := 0; i < affixCount; i++ {
		affix, err := em.affixGenerator.GenerateAffix(item.Slot, level, i == 0) // 第一个是前缀
		if err != nil {
			return nil, fmt.Errorf("failed to generate affix: %w", err)
		}
		affixes = append(affixes, affix)
	}

	// 创建装备实例
	equipment := &models.EquipmentInstance{
		ItemID:     itemID,
		Quality:    quality,
		Level:      level,
		PrefixID:   nil,
		PrefixValue: nil,
		SuffixID:   nil,
		SuffixValue: nil,
	}

	// 设置词缀
	if len(affixes) > 0 {
		if affixes[0].Type == "prefix" {
			equipment.PrefixID = &affixes[0].ID
			equipment.PrefixValue = &affixes[0].Value
		} else {
			equipment.SuffixID = &affixes[0].ID
			equipment.SuffixValue = &affixes[0].Value
		}
	}

	if len(affixes) > 1 {
		if affixes[1].Type == "suffix" {
			equipment.SuffixID = &affixes[1].ID
			equipment.SuffixValue = &affixes[1].Value
		}
	}

	return equipment, nil
	*/
}

// getAffixCountByQuality 根据品质获取词缀数量
func (em *EquipmentManager) getAffixCountByQuality(quality string) int {
	switch quality {
	case "common", "white":
		return 0
	case "uncommon", "green":
		return 1
	case "rare", "blue":
		return 2
	case "epic", "purple":
		return 3
	case "legendary", "orange":
		return 4
	case "unique":
		return 0 // 独特装备使用固定词缀
	default:
		return 0
	}
}

// GeneratedAffix 生成的词缀
type GeneratedAffix struct {
	ID    string
	Name  string
	Type  string  // prefix/suffix
	Value float64
}

// GenerateAffix 生成词缀
func (ag *AffixGenerator) GenerateAffix(slot string, level int, isPrefix bool) (*GeneratedAffix, error) {
	// 确定词缀类型
	affixType := "suffix"
	if isPrefix {
		affixType = "prefix"
	}

	// 确定词缀Tier
	tier := ag.determineAffixTier(level)

	// 从词缀池中随机选择
	affixConfig, err := ag.selectRandomAffix(slot, affixType, tier, level)
	if err != nil {
		return nil, err
	}

	// 生成词缀数值
	value := ag.generateAffixValue(affixConfig, tier)

	return &GeneratedAffix{
		ID:    affixConfig.ID,
		Name:  affixConfig.Name,
		Type:  affixType,
		Value: value,
	}, nil
}

// GenerateAffixByID 根据ID生成词缀（用于精华材料）
func (ag *AffixGenerator) GenerateAffixByID(affixID string, slot string, level int, isPrefix bool) (*GeneratedAffix, error) {
	// 从数据库加载词缀配置
	affixConfig, err := ag.loadAffixConfigByID(affixID)
	if err != nil {
		return nil, fmt.Errorf("failed to load affix config: %w", err)
	}

	// 确定词缀Tier
	tier := ag.determineAffixTier(level)

	// 生成词缀数值
	value := ag.generateAffixValue(affixConfig, tier)

	affixType := "suffix"
	if isPrefix {
		affixType = "prefix"
	}

	return &GeneratedAffix{
		ID:    affixConfig.ID,
		Name:  affixConfig.Name,
		Type:  affixType,
		Value: value,
	}, nil
}

// loadAffixConfigByID 根据ID加载词缀配置
func (ag *AffixGenerator) loadAffixConfigByID(affixID string) (*AffixConfig, error) {
	// 检查缓存
	ag.mu.RLock()
	if config, exists := ag.affixConfigs[affixID]; exists {
		ag.mu.RUnlock()
		return config, nil
	}
	ag.mu.RUnlock()

	// 从数据库加载
	var id, name, affixType, slotType, rarity, effectType, effectStat, valueType string
	var minValue, maxValue float64
	var levelRequired, tier int

	err := database.DB.QueryRow(`
		SELECT id, name, type, slot_type, rarity, effect_type, effect_stat,
		       min_value, max_value, value_type, level_required
		FROM affixes WHERE id = ?`, affixID,
	).Scan(
		&id, &name, &affixType, &slotType, &rarity, &effectType, &effectStat,
		&minValue, &maxValue, &valueType, &levelRequired,
	)
	if err != nil {
		return nil, err
	}

	// 确定Tier（简化处理，根据level_required）
	if levelRequired <= 20 {
		tier = 1
	} else if levelRequired <= 40 {
		tier = 2
	} else {
		tier = 3
	}

	config := &AffixConfig{
		ID:            id,
		Name:          name,
		Type:          affixType,
		SlotType:      slotType,
		Rarity:        rarity,
		EffectType:    effectType,
		EffectStat:    effectStat,
		MinValue:      minValue,
		MaxValue:      maxValue,
		ValueType:     valueType,
		LevelRequired: levelRequired,
		Tier:          tier,
	}

	// 缓存
	ag.mu.Lock()
	ag.affixConfigs[affixID] = config
	ag.mu.Unlock()

	return config, nil
}

// determineAffixTier 确定词缀Tier
func (ag *AffixGenerator) determineAffixTier(level int) int {
	if level <= 20 {
		return 1
	} else if level <= 40 {
		return 2
	} else {
		return 3
	}
}

// selectRandomAffix 从词缀池中随机选择词缀
func (ag *AffixGenerator) selectRandomAffix(slot, affixType string, tier, level int) (*AffixConfig, error) {
	// 从数据库中选择符合条件的词缀
	// 根据slot、affixType、tier筛选合适的词缀
	// slot_type可以是"all"或匹配的slot
	// 将slot转换为对应的slot_type（main_hand -> weapon, chest -> armor等）
	slotType := slot
	if slot == "main_hand" || slot == "off_hand" {
		slotType = "weapon"
	} else if slot == "chest" || slot == "legs" || slot == "head" || slot == "feet" || slot == "hands" {
		slotType = "armor"
	} else if slot == "ring" || slot == "neck" || slot == "trinket" {
		slotType = "accessory"
	}
	
	rows, err := database.DB.Query(`
		SELECT id, name, type, slot_type, rarity, effect_type, effect_stat,
		       min_value, max_value, value_type, level_required
		FROM affixes 
		WHERE type = ? 
		  AND (slot_type = ? OR slot_type = 'all')
		  AND level_required <= ?
		ORDER BY RANDOM()
		LIMIT 1`,
		affixType, slotType, level,
	)
	if err != nil {
		// 如果查询失败，返回一个默认词缀（但ID需要存在于数据库中）
		// 尝试使用一个简单的ID，如果不存在，则返回错误
		return nil, fmt.Errorf("failed to query affixes: %w", err)
	}
	defer rows.Close()

	if rows.Next() {
		var id, name, dbAffixType, slotType, rarity, effectType, effectStat, valueType string
		var minValue, maxValue float64
		var levelRequired int

		err := rows.Scan(&id, &name, &dbAffixType, &slotType, &rarity, &effectType, &effectStat,
			&minValue, &maxValue, &valueType, &levelRequired)
		if err != nil {
			return nil, fmt.Errorf("failed to scan affix: %w", err)
		}

		// 确定Tier
		dbTier := tier
		if levelRequired <= 20 {
			dbTier = 1
		} else if levelRequired <= 40 {
			dbTier = 2
		} else {
			dbTier = 3
		}

		return &AffixConfig{
			ID:            id,
			Name:          name,
			Type:          dbAffixType,
			SlotType:      slotType,
			Rarity:        rarity,
			EffectType:    effectType,
			EffectStat:    effectStat,
			MinValue:      minValue,
			MaxValue:      maxValue,
			ValueType:     valueType,
			LevelRequired: levelRequired,
			Tier:          dbTier,
		}, nil
	}

	// 如果没有找到词缀，返回错误
	return nil, fmt.Errorf("no affix found for slot=%s, type=%s, level=%d", slot, affixType, level)
}

// generateAffixValue 生成词缀数值
func (ag *AffixGenerator) generateAffixValue(config *AffixConfig, tier int) float64 {
	// 计算Tier倍率
	tierMultiplier := 1.0
	switch tier {
	case 2:
		tierMultiplier = 1.8
	case 3:
		tierMultiplier = 2.5
	}

	// 在范围内随机
	minValue := config.MinValue * tierMultiplier
	maxValue := config.MaxValue * tierMultiplier

	return minValue + ag.rng.Float64()*(maxValue-minValue)
}

// EquipItem 穿戴装备
func (em *EquipmentManager) EquipItem(characterID int, equipmentID int) error {
	em.mu.Lock()
	defer em.mu.Unlock()

	// 获取角色
	char, err := em.charRepo.GetByID(characterID)
	if err != nil {
		return fmt.Errorf("failed to get character: %w", err)
	}

	// 获取装备
	equipment, err := em.equipmentRepo.GetByID(equipmentID)
	if err != nil {
		return fmt.Errorf("failed to get equipment: %w", err)
	}

	// 验证装备所有权
	if equipment.OwnerID != char.UserID {
		return fmt.Errorf("equipment does not belong to user")
	}

	// 验证装备要求（等级、职业等）
	if err := em.validateEquipmentRequirements(char, equipment); err != nil {
		return err
	}

	// 卸下同槽位的旧装备（如果有）
	// 注意：不能调用UnequipItem，因为已经持有锁，会导致死锁
	oldEquipment, err := em.equipmentRepo.GetByCharacterAndSlot(characterID, equipment.Slot)
	if err != nil {
		return fmt.Errorf("failed to check existing equipment: %w", err)
	}
	if oldEquipment != nil {
		// 直接卸下旧装备（不调用UnequipItem，避免死锁）
		oldEquipment.CharacterID = nil
		if err := em.equipmentRepo.Update(oldEquipment); err != nil {
			return fmt.Errorf("failed to unequip old equipment: %w", err)
		}
	}

	// 穿戴装备
	equipment.CharacterID = &characterID
	if err := em.equipmentRepo.Update(equipment); err != nil {
		return fmt.Errorf("failed to equip item: %w", err)
	}

	// 更新角色属性（应用装备加成）
	return em.updateCharacterAttributes(char)
}

// UnequipItem 卸下装备
func (em *EquipmentManager) UnequipItem(characterID int, equipmentID int) error {
	em.mu.Lock()
	defer em.mu.Unlock()

	equipment, err := em.equipmentRepo.GetByID(equipmentID)
	if err != nil {
		return fmt.Errorf("failed to get equipment: %w", err)
	}

	if equipment.CharacterID == nil || *equipment.CharacterID != characterID {
		return fmt.Errorf("equipment not equipped by character")
	}

	equipment.CharacterID = nil
	if err := em.equipmentRepo.Update(equipment); err != nil {
		return fmt.Errorf("failed to unequip item: %w", err)
	}

	// 更新角色属性（移除装备加成）
	char, err := em.charRepo.GetByID(characterID)
	if err != nil {
		return err
	}

	return em.updateCharacterAttributes(char)
}

// validateEquipmentRequirements 验证装备要求
func (em *EquipmentManager) validateEquipmentRequirements(char *models.Character, equipment *models.EquipmentInstance) error {
	// 获取基础物品配置
	item, err := em.gameRepo.GetItemByID(equipment.ItemID)
	if err != nil {
		return fmt.Errorf("failed to get item: %w", err)
	}

	// 调试：打印item配置（仅在TEST_DEBUG=1时输出）
	if os.Getenv("TEST_DEBUG") == "1" {
		fmt.Printf("[DEBUG] validateEquipmentRequirements: itemID=%s, charLevel=%d, charClass=%s, item level_required=%v, item class_required=%v\n", 
			equipment.ItemID, char.Level, char.ClassID, item["level_required"], item["class_required"])
	}

	// 检查等级要求
	levelRequired, ok := item["level_required"].(int)
	if !ok {
		// 尝试从数据库读取的int类型
		if levelRequiredFloat, ok := item["level_required"].(float64); ok {
			levelRequired = int(levelRequiredFloat)
		} else {
			levelRequired = 0
		}
	}
	if levelRequired > 0 && char.Level < levelRequired {
		return fmt.Errorf("等级不足：角色等级 %d 低于需求等级 %d", char.Level, levelRequired)
	}

	// 检查职业要求
	classRequired, _ := item["class_required"].(string)
	if classRequired != "" && char.ClassID != classRequired {
		return fmt.Errorf("职业不匹配：角色职业 %s 不符合需求职业 %s", char.ClassID, classRequired)
	}

	// 检查属性要求
	strengthRequired := getIntFromMap(item, "strength_required")
	if strengthRequired > 0 && char.Strength < strengthRequired {
		return fmt.Errorf("属性不足：角色力量 %d 低于需求 %d", char.Strength, strengthRequired)
	}
	agilityRequired := getIntFromMap(item, "agility_required")
	if agilityRequired > 0 && char.Agility < agilityRequired {
		return fmt.Errorf("character agility %d is below required %d", char.Agility, agilityRequired)
	}
	intellectRequired := getIntFromMap(item, "intellect_required")
	if intellectRequired > 0 && char.Intellect < intellectRequired {
		return fmt.Errorf("character intellect %d is below required %d", char.Intellect, intellectRequired)
	}

	return nil
}

// getIntFromMap 从map中安全获取int值
func getIntFromMap(m map[string]interface{}, key string) int {
	val, exists := m[key]
	if !exists {
		return 0
	}
	switch v := val.(type) {
	case int:
		return v
	case float64:
		return int(v)
	case int64:
		return int(v)
	default:
		return 0
	}
}

// updateCharacterAttributes 更新角色属性（应用装备加成）
func (em *EquipmentManager) updateCharacterAttributes(char *models.Character) error {
	// 获取角色所有装备
	equipments, err := em.equipmentRepo.GetByCharacterID(char.ID)
	if err != nil {
		return fmt.Errorf("failed to get equipments: %w", err)
	}

	// 计算装备属性加成总和
	equipmentBonuses := em.calculateEquipmentBonuses(equipments)

	// 获取职业配置以获取基础HP和MP
	class, err := em.gameRepo.GetClassByID(char.ClassID)
	if err != nil {
		// 如果获取失败，使用默认值
		char.PhysicalAttack = em.calculator.CalculatePhysicalAttack(char) + equipmentBonuses.PhysicalAttack
		char.MagicAttack = em.calculator.CalculateMagicAttack(char) + equipmentBonuses.MagicAttack
		char.MaxHP = em.calculator.CalculateHP(char, 100) + equipmentBonuses.HPBonus // 默认基础HP
		char.MaxResource = em.calculator.CalculateMP(char, 50) + equipmentBonuses.MPBonus // 默认基础MP
	} else {
		// 计算基础HP和MP（考虑等级成长）
		baseHP := class.BaseHP + class.HPPerLevel*(char.Level-1)
		baseMP := class.BaseResource + class.ResourcePerLevel*(char.Level-1)
		
		// 重新计算派生属性并应用装备加成
		char.PhysicalAttack = em.calculator.CalculatePhysicalAttack(char) + equipmentBonuses.PhysicalAttack
		char.MagicAttack = em.calculator.CalculateMagicAttack(char) + equipmentBonuses.MagicAttack
		char.MaxHP = em.calculator.CalculateHP(char, baseHP) + equipmentBonuses.HPBonus
		char.MaxResource = em.calculator.CalculateMP(char, baseMP) + equipmentBonuses.MPBonus
	}

	// 应用装备的基础属性加成
	char.Strength += equipmentBonuses.Strength
	char.Agility += equipmentBonuses.Agility
	char.Intellect += equipmentBonuses.Intellect
	char.Stamina += equipmentBonuses.Stamina
	char.Spirit += equipmentBonuses.Spirit

	// 更新角色到数据库
	return em.charRepo.Update(char)
}

// EquipmentBonuses 装备属性加成总和
type EquipmentBonuses struct {
	Strength        int
	Agility         int
	Intellect       int
	Stamina         int
	Spirit          int
	PhysicalAttack  int
	MagicAttack     int
	PhysicalDefense int
	MagicDefense    int
	HPBonus         int
	MPBonus         int
	CritRate        float64
	CritDamage      float64
	DodgeRate       float64
}

// calculateEquipmentBonuses 计算所有装备的属性加成总和
func (em *EquipmentManager) calculateEquipmentBonuses(equipments []*models.EquipmentInstance) EquipmentBonuses {
	bonuses := EquipmentBonuses{}

	for _, equipment := range equipments {
		// 获取基础物品配置
		item, err := em.gameRepo.GetItemByID(equipment.ItemID)
		if err != nil {
			continue
		}

		// 累加基础物品属性
		if strength, ok := item["strength"].(int); ok {
			bonuses.Strength += strength
		}
		if agility, ok := item["agility"].(int); ok {
			bonuses.Agility += agility
		}
		if intellect, ok := item["intellect"].(int); ok {
			bonuses.Intellect += intellect
		}
		if stamina, ok := item["stamina"].(int); ok {
			bonuses.Stamina += stamina
		}
		if spirit, ok := item["spirit"].(int); ok {
			bonuses.Spirit += spirit
		}
		if attack, ok := item["attack"].(int); ok {
			bonuses.PhysicalAttack += attack
		}
		if defense, ok := item["defense"].(int); ok {
			bonuses.PhysicalDefense += defense
		}
		if hpBonus, ok := item["hp_bonus"].(int); ok {
			bonuses.HPBonus += hpBonus
		}
		if mpBonus, ok := item["mp_bonus"].(int); ok {
			bonuses.MPBonus += mpBonus
		}

		// 累加词缀属性
		em.applyAffixBonuses(equipment, &bonuses)
	}

	return bonuses
}

// applyAffixBonuses 应用词缀属性加成
func (em *EquipmentManager) applyAffixBonuses(equipment *models.EquipmentInstance, bonuses *EquipmentBonuses) {
	// 处理前缀
	if equipment.PrefixID != nil && equipment.PrefixValue != nil {
		em.applyAffixBonus(*equipment.PrefixID, *equipment.PrefixValue, bonuses)
	}

	// 处理后缀
	if equipment.SuffixID != nil && equipment.SuffixValue != nil {
		em.applyAffixBonus(*equipment.SuffixID, *equipment.SuffixValue, bonuses)
	}

	// 处理额外词缀1
	if equipment.BonusAffix1 != nil && equipment.BonusAffix1Value != nil {
		em.applyAffixBonus(*equipment.BonusAffix1, *equipment.BonusAffix1Value, bonuses)
	}

	// 处理额外词缀2
	if equipment.BonusAffix2 != nil && equipment.BonusAffix2Value != nil {
		em.applyAffixBonus(*equipment.BonusAffix2, *equipment.BonusAffix2Value, bonuses)
	}
}

// applyAffixBonus 应用单个词缀的属性加成
func (em *EquipmentManager) applyAffixBonus(affixID string, affixValue float64, bonuses *EquipmentBonuses) {
	// 加载词缀配置
	affixConfig, err := em.affixGenerator.loadAffixConfigByID(affixID)
	if err != nil {
		// 如果加载失败，跳过该词缀
		return
	}

	// 根据effect_stat和value_type应用加成
	effectStat := affixConfig.EffectStat
	valueType := affixConfig.ValueType

	// 如果是百分比类型，需要在实际计算时应用（这里只处理固定值）
	if valueType == "percent" {
		// 百分比加成在计算时应用，这里不处理
		return
	}

	// 应用固定值加成
	switch effectStat {
	case "strength":
		bonuses.Strength += int(affixValue)
	case "agility":
		bonuses.Agility += int(affixValue)
	case "intellect":
		bonuses.Intellect += int(affixValue)
	case "stamina":
		bonuses.Stamina += int(affixValue)
	case "spirit":
		bonuses.Spirit += int(affixValue)
	case "attack", "physical_attack":
		bonuses.PhysicalAttack += int(affixValue)
	case "magic_attack":
		bonuses.MagicAttack += int(affixValue)
	case "defense", "physical_defense":
		bonuses.PhysicalDefense += int(affixValue)
	case "magic_defense":
		bonuses.MagicDefense += int(affixValue)
	case "hp", "hp_bonus":
		bonuses.HPBonus += int(affixValue)
	case "mp", "mp_bonus", "resource_bonus":
		bonuses.MPBonus += int(affixValue)
	case "crit_rate":
		bonuses.CritRate += float64(affixValue)
	case "crit_damage":
		bonuses.CritDamage += float64(affixValue)
	case "dodge_rate":
		bonuses.DodgeRate += float64(affixValue)
	}
}

// EnhanceEquipment 强化装备
func (em *EquipmentManager) EnhanceEquipment(characterID int, equipmentID int, materials []Material) error {
	em.mu.Lock()
	defer em.mu.Unlock()

	equipment, err := em.equipmentRepo.GetByID(equipmentID)
	if err != nil {
		return fmt.Errorf("failed to get equipment: %w", err)
	}

	// 验证装备所有权
	if equipment.OwnerID != characterID {
		return fmt.Errorf("equipment does not belong to character")
	}

	// 处理强化材料
	for _, material := range materials {
		switch material.Type {
		case "catalyst":
			// 催化剂：强化已有词缀数值
			if err := em.enhanceAffixValue(equipment, material); err != nil {
				return err
			}
		case "essence":
			// 精华材料：保证特定词缀
			if err := em.guaranteeAffix(equipment, material); err != nil {
				return err
			}
		case "base":
			// 基础材料：重铸/添加词缀
			if err := em.reforgeAffix(equipment, material); err != nil {
				return err
			}
		case "protection":
			// 保护材料：锁定词缀
			if err := em.lockAffix(equipment, material); err != nil {
				return err
			}
		}
	}

	return em.equipmentRepo.Update(equipment)
}

// Material 强化材料
type Material struct {
	ID     string
	Type   string  // catalyst/essence/base/protection
	AffixID string // 目标词缀ID（用于精华和保护材料）
}

// enhanceAffixValue 强化词缀数值
func (em *EquipmentManager) enhanceAffixValue(equipment *models.EquipmentInstance, material Material) error {
	// 提升词缀数值（提升10-20%）
	enhanceMultiplier := 1.1 + em.affixGenerator.rng.Float64()*0.1 // 1.1-1.2倍

	// 提升前缀数值
	if equipment.PrefixID != nil && equipment.PrefixValue != nil {
		newValue := *equipment.PrefixValue * enhanceMultiplier
		equipment.PrefixValue = &newValue
	}

	// 提升后缀数值
	if equipment.SuffixID != nil && equipment.SuffixValue != nil {
		newValue := *equipment.SuffixValue * enhanceMultiplier
		equipment.SuffixValue = &newValue
	}

	// 提升额外词缀数值
	if equipment.BonusAffix1 != nil && equipment.BonusAffix1Value != nil {
		newValue := *equipment.BonusAffix1Value * enhanceMultiplier
		equipment.BonusAffix1Value = &newValue
	}

	if equipment.BonusAffix2 != nil && equipment.BonusAffix2Value != nil {
		newValue := *equipment.BonusAffix2Value * enhanceMultiplier
		equipment.BonusAffix2Value = &newValue
	}

	return nil
}

// guaranteeAffix 保证特定词缀
func (em *EquipmentManager) guaranteeAffix(equipment *models.EquipmentInstance, material Material) error {
	// 如果装备没有该词缀，添加该词缀
	// 如果已有该词缀，提升数值
	if material.AffixID == "" {
		return fmt.Errorf("affix ID required for essence material")
	}

	// 检查是否已有该词缀
	hasAffix := false
	if equipment.PrefixID != nil && *equipment.PrefixID == material.AffixID {
		hasAffix = true
		// 提升前缀数值
		if equipment.PrefixValue != nil {
			newValue := *equipment.PrefixValue * 1.15
			equipment.PrefixValue = &newValue
		}
	} else if equipment.SuffixID != nil && *equipment.SuffixID == material.AffixID {
		hasAffix = true
		// 提升后缀数值
		if equipment.SuffixValue != nil {
			newValue := *equipment.SuffixValue * 1.15
			equipment.SuffixValue = &newValue
		}
	}

	// 如果没有该词缀，尝试添加到空槽位
	if !hasAffix {
		if equipment.PrefixID == nil {
			// 生成前缀
			affix, err := em.affixGenerator.GenerateAffixByID(material.AffixID, equipment.Slot, 1, true)
			if err == nil {
				equipment.PrefixID = &affix.ID
				equipment.PrefixValue = &affix.Value
			}
		} else if equipment.SuffixID == nil {
			// 生成后缀
			affix, err := em.affixGenerator.GenerateAffixByID(material.AffixID, equipment.Slot, 1, false)
			if err == nil {
				equipment.SuffixID = &affix.ID
				equipment.SuffixValue = &affix.Value
			}
		}
	}

	return nil
}

// reforgeAffix 重铸词缀
func (em *EquipmentManager) reforgeAffix(equipment *models.EquipmentInstance, material Material) error {
	// 随机重铸词缀（随机选择前缀或后缀重铸）
	if em.affixGenerator.rng.Float64() < 0.5 {
		// 重铸前缀
		if equipment.PrefixID != nil {
			affix, err := em.affixGenerator.GenerateAffix(equipment.Slot, 1, true)
			if err == nil {
				equipment.PrefixID = &affix.ID
				equipment.PrefixValue = &affix.Value
			}
		}
	} else {
		// 重铸后缀
		if equipment.SuffixID != nil {
			affix, err := em.affixGenerator.GenerateAffix(equipment.Slot, 1, false)
			if err == nil {
				equipment.SuffixID = &affix.ID
				equipment.SuffixValue = &affix.Value
			}
		}
	}
	return nil
}

// lockAffix 锁定词缀
func (em *EquipmentManager) lockAffix(equipment *models.EquipmentInstance, material Material) error {
	// 锁定词缀，防止被重铸改变
	equipment.IsLocked = true
	return nil
}


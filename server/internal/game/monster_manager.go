package game

import (
	"fmt"
	"math"
	"math/rand"
	"sync"

	"text-wow/internal/models"
	"text-wow/internal/repository"
)

// MonsterManager 怪物管理器 - 管理怪物生成、配置和AI
type MonsterManager struct {
	mu          sync.RWMutex
	monsterConfigs map[string]*MonsterConfig // 怪物配置缓存
	gameRepo    *repository.GameRepository
	calculator  *Calculator
}

// MonsterConfig 怪物配置（从配置表加载）
type MonsterConfig struct {
	ID                string
	Name              string
	Type              string  // normal/elite/boss/special
	Level             int
	BaseHP            int
	BasePhysicalAttack int
	BaseMagicAttack   int
	BasePhysicalDefense int
	BaseMagicDefense  int
	HPMultiplier      float64
	AttackMultiplier  float64
	DefenseMultiplier float64
	CritRate          float64
	CritDamage        float64
	DodgeRate         float64
	Speed             int
	AIType            string
	AIBehavior        string  // JSON格式
	SkillIDs          []string
	DropTableID       string
	BalanceVersion    int
}

// NewMonsterManager 创建怪物管理器
func NewMonsterManager() *MonsterManager {
	return &MonsterManager{
		monsterConfigs: make(map[string]*MonsterConfig),
		gameRepo:       repository.NewGameRepository(),
		calculator:     NewCalculator(),
	}
}

// LoadMonsterConfig 加载怪物配置
func (mm *MonsterManager) LoadMonsterConfig(monsterID string) (*MonsterConfig, error) {
	mm.mu.RLock()
	if config, exists := mm.monsterConfigs[monsterID]; exists {
		mm.mu.RUnlock()
		return config, nil
	}
	mm.mu.RUnlock()

	// 从数据库加载配置
	monster, err := mm.gameRepo.GetMonsterByID(monsterID)
	if err != nil {
		return nil, fmt.Errorf("failed to load monster config: %w", err)
	}

	config := &MonsterConfig{
		ID:                monster.ID,
		Name:              monster.Name,
		Type:              monster.Type,
		Level:             monster.Level,
		BaseHP:            monster.HP,
		BasePhysicalAttack: monster.PhysicalAttack,
		BaseMagicAttack:    monster.MagicAttack,
		BasePhysicalDefense: monster.PhysicalDefense,
		BaseMagicDefense:   monster.MagicDefense,
		CritRate:           monster.PhysCritRate,
		CritDamage:          monster.PhysCritDamage,
		DodgeRate:          monster.DodgeRate,
		Speed:              monster.Speed,
		AIType:             monster.AIType,
		AIBehavior:         monster.AIBehavior,
		HPMultiplier:        1.0,
		AttackMultiplier:   1.0,
		DefenseMultiplier:  1.0,
	}

	// 根据怪物类型应用系数
	switch monster.Type {
	case "elite":
		config.HPMultiplier = 1.8
		config.AttackMultiplier = 1.5
		config.DefenseMultiplier = 1.3
	case "boss":
		config.HPMultiplier = 4.0
		config.AttackMultiplier = 2.5
		config.DefenseMultiplier = 2.0
	case "special":
		// 特殊怪物根据配置调整
		config.HPMultiplier = 1.5
		config.AttackMultiplier = 1.5
		config.DefenseMultiplier = 1.5
	}

	mm.mu.Lock()
	mm.monsterConfigs[monsterID] = config
	mm.mu.Unlock()

	return config, nil
}

// GenerateMonster 生成怪物实例
func (mm *MonsterManager) GenerateMonster(zoneID string, level int) (*models.Monster, error) {
	// 从区域获取怪物列表
	zone, err := mm.gameRepo.GetZoneByID(zoneID)
	if err != nil {
		return nil, fmt.Errorf("failed to get zone: %w", err)
	}

	if len(zone.Monsters) == 0 {
		return nil, fmt.Errorf("zone has no monsters")
	}

	// 根据权重随机选择怪物
	monster := mm.selectMonsterByWeight(zone.Monsters, level)
	if monster == nil {
		return nil, fmt.Errorf("no suitable monster found")
	}

	// 加载怪物配置（从缓存或数据库）
	config, err := mm.LoadMonsterConfig(monster.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to load monster config: %w", err)
	}

	// 根据配置生成怪物实例
	generatedMonster := mm.createMonsterInstance(config, level)

	return generatedMonster, nil
}

// selectMonsterByWeight 根据权重随机选择怪物
func (mm *MonsterManager) selectMonsterByWeight(monsters []models.Monster, level int) *models.Monster {
	totalWeight := 0
	suitableMonsters := make([]models.Monster, 0)

	// 筛选适合等级的怪物
	for _, m := range monsters {
		if m.Level >= level-2 && m.Level <= level+2 {
			totalWeight += m.SpawnWeight
			suitableMonsters = append(suitableMonsters, m)
		}
	}

	if len(suitableMonsters) == 0 {
		// 如果没有适合等级的，选择所有怪物
		for _, m := range monsters {
			totalWeight += m.SpawnWeight
			suitableMonsters = append(suitableMonsters, m)
		}
	}

	if totalWeight == 0 {
		return nil
	}

	// 随机选择
	roll := rand.Intn(totalWeight)
	currentWeight := 0
	for i := range suitableMonsters {
		currentWeight += suitableMonsters[i].SpawnWeight
		if roll < currentWeight {
			return &suitableMonsters[i]
		}
	}

	return &suitableMonsters[len(suitableMonsters)-1]
}

// createMonsterInstance 根据配置创建怪物实例
func (mm *MonsterManager) createMonsterInstance(config *MonsterConfig, level int) *models.Monster {
	// 计算等级系数（每级增长约5%）
	levelMultiplier := 1.0 + float64(level-config.Level)*0.05
	if levelMultiplier < 0.5 {
		levelMultiplier = 0.5
	}
	if levelMultiplier > 2.0 {
		levelMultiplier = 2.0
	}

	// 计算最终属性
	hp := int(math.Round(float64(config.BaseHP) * config.HPMultiplier * levelMultiplier))
	physicalAttack := int(math.Round(float64(config.BasePhysicalAttack) * config.AttackMultiplier * levelMultiplier))
	magicAttack := int(math.Round(float64(config.BaseMagicAttack) * config.AttackMultiplier * levelMultiplier))
	physicalDefense := int(math.Round(float64(config.BasePhysicalDefense) * config.DefenseMultiplier * levelMultiplier))
	magicDefense := int(math.Round(float64(config.BaseMagicDefense) * config.DefenseMultiplier * levelMultiplier))

	monster := &models.Monster{
		ID:              config.ID,
		ZoneID:          "",
		Name:            config.Name,
		Level:           level,
		Type:            config.Type,
		HP:              hp,
		MaxHP:           hp,
		MP:              100, // 默认MP
		MaxMP:           100,
		PhysicalAttack:  physicalAttack,
		MagicAttack:     magicAttack,
		PhysicalDefense: physicalDefense,
		MagicDefense:    magicDefense,
		AttackType:      "physical", // 默认物理攻击
		PhysCritRate:    config.CritRate,
		PhysCritDamage:  config.CritDamage,
		SpellCritRate:   config.CritRate,
		SpellCritDamage: config.CritDamage,
		DodgeRate:       config.DodgeRate,
		Speed:           config.Speed,
		ExpReward:       level * 5, // 基础经验 = 等级 × 5
		GoldMin:         level,
		GoldMax:         level * 3,
		SpawnWeight:     1,
		AIType:          config.AIType,
		AIBehavior:      config.AIBehavior,
	}

	// 加载怪物技能
	skills, err := mm.gameRepo.GetMonsterSkills(config.ID)
	if err == nil {
		monster.MonsterSkills = skills
	}

	return monster
}

// DropResult 掉落结果
type DropResult struct {
	ItemID   string
	Quantity int
}

// CalculateDrops 计算怪物掉落
// 根据怪物的掉落表，计算实际掉落的物品
func (mm *MonsterManager) CalculateDrops(monsterID string, monsterType string) ([]DropResult, error) {
	// 获取掉落表
	drops, err := mm.gameRepo.GetMonsterDrops(monsterID)
	if err != nil {
		// 如果没有掉落表，返回空结果
		return []DropResult{}, nil
	}

	var results []DropResult
	
	// 根据怪物类型应用掉落率修正
	dropRateMultiplier := 1.0
	switch monsterType {
	case "elite":
		dropRateMultiplier = 1.2 // 精英怪物掉落率提升20%
	case "boss":
		dropRateMultiplier = 1.5 // Boss掉落率提升50%
	case "special":
		dropRateMultiplier = 1.3 // 特殊怪物掉落率提升30%
	}

	// 遍历掉落表，根据概率计算掉落
	for _, drop := range drops {
		// 应用掉落率修正
		adjustedRate := drop.DropRate * dropRateMultiplier
		if adjustedRate > 1.0 {
			adjustedRate = 1.0 // 确保不超过100%
		}

		// 随机判断是否掉落
		if rand.Float64() < adjustedRate {
			// 计算掉落数量
			quantity := drop.MinQuantity
			if drop.MaxQuantity > drop.MinQuantity {
				quantity = drop.MinQuantity + rand.Intn(drop.MaxQuantity-drop.MinQuantity+1)
			}

			results = append(results, DropResult{
				ItemID:   drop.ItemID,
				Quantity: quantity,
			})
		}
	}

	return results, nil
}

// ReloadMonsterConfig 重新加载怪物配置（热更新）
func (mm *MonsterManager) ReloadMonsterConfig(monsterID string) error {
	mm.mu.Lock()
	delete(mm.monsterConfigs, monsterID)
	mm.mu.Unlock()

	_, err := mm.LoadMonsterConfig(monsterID)
	return err
}

// GetMonsterByID 根据ID获取怪物（便捷方法）
func (mm *MonsterManager) GetMonsterByID(monsterID string) (*models.Monster, error) {
	return mm.gameRepo.GetMonsterByID(monsterID)
}

// GetMonsterConfig 获取怪物配置
func (mm *MonsterManager) GetMonsterConfig(monsterID string) (*MonsterConfig, error) {
	return mm.LoadMonsterConfig(monsterID)
}

// LoadAllMonsterConfigs 批量加载所有怪物配置
func (mm *MonsterManager) LoadAllMonsterConfigs() error {
	// 从区域配置中获取所有怪物ID
	zones, err := mm.gameRepo.GetZones()
	if err != nil {
		return fmt.Errorf("failed to get zones: %w", err)
	}

	monsterIDs := make(map[string]bool)
	for _, zone := range zones {
		monsters, err := mm.gameRepo.GetMonstersByZone(zone.ID)
		if err != nil {
			continue // 跳过无法加载的区域
		}
		for _, monster := range monsters {
			monsterIDs[monster.ID] = true
		}
	}

	// 加载所有怪物配置
	for monsterID := range monsterIDs {
		_, err := mm.LoadMonsterConfig(monsterID)
		if err != nil {
			// 记录错误但继续加载其他怪物
			fmt.Printf("Warning: failed to load monster config %s: %v\n", monsterID, err)
		}
	}

	return nil
}

// ReloadAllMonsterConfigs 重新加载所有怪物配置（热更新）
func (mm *MonsterManager) ReloadAllMonsterConfigs() error {
	mm.mu.Lock()
	mm.monsterConfigs = make(map[string]*MonsterConfig)
	mm.mu.Unlock()

	return mm.LoadAllMonsterConfigs()
}
package config

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"text-wow/internal/database"
)

// ConfigManager 配置管理器 - 管理所有游戏配置数据
type ConfigManager struct {
	mu              sync.RWMutex
	configs         map[string]*ConfigCache // 配置缓存
	listeners       []ConfigChangeListener  // 配置变更监听器
	versionManager  *VersionManager
}

// ConfigCache 配置缓存
type ConfigCache struct {
	Data      interface{}
	Version   int
	UpdatedAt time.Time
}

// ConfigChangeListener 配置变更监听器接口
type ConfigChangeListener interface {
	OnConfigChange(configType string, version int)
}

// VersionManager 版本管理器
type VersionManager struct {
	mu sync.RWMutex
}

// NewConfigManager 创建配置管理器
func NewConfigManager() *ConfigManager {
	return &ConfigManager{
		configs:        make(map[string]*ConfigCache),
		listeners:      make([]ConfigChangeListener, 0),
		versionManager: &VersionManager{},
	}
}

// LoadConfig 加载配置
// configType: monster/skill/item/economy/zone
func (cm *ConfigManager) LoadConfig(configType string) error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	switch configType {
	case "monster":
		return cm.loadMonsterConfigs()
	case "skill":
		return cm.loadSkillConfigs()
	case "item":
		return cm.loadItemConfigs()
	case "economy":
		return cm.loadEconomyConfigs()
	case "zone":
		return cm.loadZoneConfigs()
	default:
		return fmt.Errorf("unknown config type: %s", configType)
	}
}

// loadMonsterConfigs 加载怪物配置
func (cm *ConfigManager) loadMonsterConfigs() error {
	// 从数据库加载怪物配置
	rows, err := database.DB.Query(`
		SELECT id, name, type, level, hp, physical_attack, magic_attack,
		       physical_defense, magic_defense, crit_rate, crit_damage,
		       dodge_rate, speed, ai_type, ai_behavior, spawn_weight
		FROM monsters
	`)
	if err != nil {
		return fmt.Errorf("failed to load monster configs: %w", err)
	}
	defer rows.Close()

	monsters := make(map[string]interface{})
	for rows.Next() {
		var id, name, monsterType, aiType string
		var level, hp, physicalAttack, magicAttack, physicalDefense, magicDefense, speed, spawnWeight int
		var critRate, critDamage, dodgeRate float64
		var aiBehavior sql.NullString

		err := rows.Scan(&id, &name, &monsterType, &level, &hp, &physicalAttack, &magicAttack,
			&physicalDefense, &magicDefense, &critRate, &critDamage, &dodgeRate, &speed,
			&aiType, &aiBehavior, &spawnWeight)
		if err != nil {
			log.Printf("Failed to scan monster config: %v", err)
			continue
		}

		monster := map[string]interface{}{
			"id":               id,
			"name":             name,
			"type":             monsterType,
			"level":            level,
			"hp":               hp,
			"physical_attack":  physicalAttack,
			"magic_attack":     magicAttack,
			"physical_defense": physicalDefense,
			"magic_defense":    magicDefense,
			"crit_rate":        critRate,
			"crit_damage":      critDamage,
			"dodge_rate":       dodgeRate,
			"speed":            speed,
			"ai_type":          aiType,
			"spawn_weight":    spawnWeight,
		}
		if aiBehavior.Valid {
			monster["ai_behavior"] = aiBehavior.String
		}

		monsters[id] = monster
	}

	version := cm.getConfigVersion("monster")
	if version == 0 {
		version = 1 // 默认版本
	}
	cm.configs["monster"] = &ConfigCache{
		Data:      monsters,
		Version:   version,
		UpdatedAt: time.Now(),
	}

	log.Printf("✅ Loaded %d monster configs", len(monsters))
	return nil
}

// loadSkillConfigs 加载技能配置
func (cm *ConfigManager) loadSkillConfigs() error {
	rows, err := database.DB.Query(`
		SELECT id, name, class_id, skill_type, resource_cost, cooldown,
		       damage_multiplier, healing_multiplier, effect_type, effect_value
		FROM skills
	`)
	if err != nil {
		return fmt.Errorf("failed to load skill configs: %w", err)
	}
	defer rows.Close()

	skills := make(map[string]interface{})
	for rows.Next() {
		var id, name, classID, skillType, effectType string
		var resourceCost, cooldown, effectValue int
		var damageMultiplier, healingMultiplier float64

		err := rows.Scan(&id, &name, &classID, &skillType, &resourceCost, &cooldown,
			&damageMultiplier, &healingMultiplier, &effectType, &effectValue)
		if err != nil {
			log.Printf("Failed to scan skill config: %v", err)
			continue
		}

		skill := map[string]interface{}{
			"id":                 id,
			"name":               name,
			"class_id":           classID,
			"skill_type":         skillType,
			"resource_cost":      resourceCost,
			"cooldown":           cooldown,
			"damage_multiplier":  damageMultiplier,
			"healing_multiplier": healingMultiplier,
			"effect_type":        effectType,
			"effect_value":       effectValue,
		}

		skills[id] = skill
	}

	version := cm.getConfigVersion("skill")
	if version == 0 {
		version = 1
	}
	cm.configs["skill"] = &ConfigCache{
		Data:      skills,
		Version:   version,
		UpdatedAt: time.Now(),
	}

	log.Printf("✅ Loaded %d skill configs", len(skills))
	return nil
}

// loadItemConfigs 加载物品配置
func (cm *ConfigManager) loadItemConfigs() error {
	rows, err := database.DB.Query(`
		SELECT id, name, type, slot, quality, base_tier, level_required,
		       class_required, strength_required, agility_required,
		       intellect_required, stamina_required, spirit_required
		FROM items
	`)
	if err != nil {
		return fmt.Errorf("failed to load item configs: %w", err)
	}
	defer rows.Close()

	items := make(map[string]interface{})
	for rows.Next() {
		var id, name, itemType, slot, quality, baseTier string
		var levelRequired int
		var classRequired sql.NullString
		var strengthRequired, agilityRequired, intellectRequired, staminaRequired, spiritRequired int

		err := rows.Scan(&id, &name, &itemType, &slot, &quality, &baseTier, &levelRequired,
			&classRequired, &strengthRequired, &agilityRequired,
			&intellectRequired, &staminaRequired, &spiritRequired)
		if err != nil {
			log.Printf("Failed to scan item config: %v", err)
			continue
		}

		item := map[string]interface{}{
			"id":                 id,
			"name":               name,
			"type":               itemType,
			"slot":               slot,
			"quality":            quality,
			"base_tier":          baseTier,
			"level_required":     levelRequired,
			"strength_required":  strengthRequired,
			"agility_required":   agilityRequired,
			"intellect_required": intellectRequired,
			"stamina_required":   staminaRequired,
			"spirit_required":    spiritRequired,
		}
		if classRequired.Valid {
			item["class_required"] = classRequired.String
		}

		items[id] = item
	}

	version := cm.getConfigVersion("item")
	if version == 0 {
		version = 1
	}
	cm.configs["item"] = &ConfigCache{
		Data:      items,
		Version:   version,
		UpdatedAt: time.Now(),
	}

	log.Printf("✅ Loaded %d item configs", len(items))
	return nil
}

// loadEconomyConfigs 加载经济配置
func (cm *ConfigManager) loadEconomyConfigs() error {
	// 从数据库加载经济配置（如果有economy_config表）
	// 目前简化处理，返回默认配置
	economy := map[string]interface{}{
		"gold_base_multiplier":    1.0,
		"exp_base_multiplier":     1.0,
		"drop_base_multiplier":    1.0,
		"enhance_cost_base":       100,
		"enhance_cost_multiplier": 1.5,
	}
	
	version := cm.getConfigVersion("economy")
	if version == 0 {
		version = 1
	}
	cm.configs["economy"] = &ConfigCache{
		Data:      economy,
		Version:   version,
		UpdatedAt: time.Now(),
	}
	
	log.Printf("✅ Loaded economy configs")
	return nil
}

// loadZoneConfigs 加载区域配置
func (cm *ConfigManager) loadZoneConfigs() error {
	rows, err := database.DB.Query(`
		SELECT id, name, description, min_level, max_level, faction,
		       exp_modifier, gold_modifier, drop_modifier,
		       unlock_zone_id, required_exploration
		FROM zones
	`)
	if err != nil {
		return fmt.Errorf("failed to load zone configs: %w", err)
	}
	defer rows.Close()

	zones := make(map[string]interface{})
	for rows.Next() {
		var id, name, description, faction string
		var minLevel, maxLevel, requiredExploration int
		var expModifier, goldModifier, dropModifier float64
		var unlockZoneID sql.NullString

		err := rows.Scan(&id, &name, &description, &minLevel, &maxLevel, &faction,
			&expModifier, &goldModifier, &dropModifier, &unlockZoneID, &requiredExploration)
		if err != nil {
			log.Printf("Failed to scan zone config: %v", err)
			continue
		}

		zone := map[string]interface{}{
			"id":                   id,
			"name":                 name,
			"description":          description,
			"min_level":            minLevel,
			"max_level":            maxLevel,
			"faction":              faction,
			"exp_modifier":         expModifier,
			"gold_modifier":        goldModifier,
			"drop_modifier":        dropModifier,
			"required_exploration": requiredExploration,
		}
		if unlockZoneID.Valid {
			zone["unlock_zone_id"] = unlockZoneID.String
		}

		zones[id] = zone
	}

	version := cm.getConfigVersion("zone")
	if version == 0 {
		version = 1
	}
	cm.configs["zone"] = &ConfigCache{
		Data:      zones,
		Version:   version,
		UpdatedAt: time.Now(),
	}

	log.Printf("✅ Loaded %d zone configs", len(zones))
	return nil
}

// GetMonsterConfig 获取怪物配置
func (cm *ConfigManager) GetMonsterConfig(monsterID string) (map[string]interface{}, error) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	cache, exists := cm.configs["monster"]
	if !exists {
		return nil, fmt.Errorf("monster configs not loaded")
	}

	monsters, ok := cache.Data.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid monster config cache")
	}

	monster, exists := monsters[monsterID]
	if !exists {
		return nil, fmt.Errorf("monster %s not found", monsterID)
	}

	monsterMap, ok := monster.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid monster config format")
	}

	return monsterMap, nil
}

// GetSkillConfig 获取技能配置
func (cm *ConfigManager) GetSkillConfig(skillID string) (map[string]interface{}, error) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	cache, exists := cm.configs["skill"]
	if !exists {
		return nil, fmt.Errorf("skill configs not loaded")
	}

	skills, ok := cache.Data.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid skill config cache")
	}

	skill, exists := skills[skillID]
	if !exists {
		return nil, fmt.Errorf("skill %s not found", skillID)
	}

	skillMap, ok := skill.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid skill config format")
	}

	return skillMap, nil
}

// GetItemConfig 获取物品配置
func (cm *ConfigManager) GetItemConfig(itemID string) (map[string]interface{}, error) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	cache, exists := cm.configs["item"]
	if !exists {
		return nil, fmt.Errorf("item configs not loaded")
	}

	items, ok := cache.Data.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid item config cache")
	}

	item, exists := items[itemID]
	if !exists {
		return nil, fmt.Errorf("item %s not found", itemID)
	}

	itemMap, ok := item.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid item config format")
	}

	return itemMap, nil
}

// GetZoneConfig 获取区域配置
func (cm *ConfigManager) GetZoneConfig(zoneID string) (map[string]interface{}, error) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	cache, exists := cm.configs["zone"]
	if !exists {
		return nil, fmt.Errorf("zone configs not loaded")
	}

	zones, ok := cache.Data.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid zone config cache")
	}

	zone, exists := zones[zoneID]
	if !exists {
		return nil, fmt.Errorf("zone %s not found", zoneID)
	}

	zoneMap, ok := zone.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid zone config format")
	}

	return zoneMap, nil
}

// GetEconomyConfig 获取经济配置
func (cm *ConfigManager) GetEconomyConfig() (map[string]interface{}, error) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	cache, exists := cm.configs["economy"]
	if !exists {
		return nil, fmt.Errorf("economy configs not loaded")
	}

	economy, ok := cache.Data.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid economy config cache")
	}

	return economy, nil
}

// GetAllConfigs 获取所有已加载的配置类型
func (cm *ConfigManager) GetAllConfigs() []string {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	configTypes := make([]string, 0, len(cm.configs))
	for configType := range cm.configs {
		configTypes = append(configTypes, configType)
	}
	return configTypes
}

// ReloadConfig 热更新配置
func (cm *ConfigManager) ReloadConfig(configType string) error {
	oldVersion := cm.getConfigVersion(configType)
	
	if err := cm.LoadConfig(configType); err != nil {
		return err
	}

	newVersion := cm.getConfigVersion(configType)
	if newVersion != oldVersion {
		cm.notifyConfigChange(configType, newVersion)
	}

	return nil
}

// RegisterConfigChangeListener 注册配置变更监听器
func (cm *ConfigManager) RegisterConfigChangeListener(listener ConfigChangeListener) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cm.listeners = append(cm.listeners, listener)
}

// notifyConfigChange 通知配置变更
func (cm *ConfigManager) notifyConfigChange(configType string, version int) {
	cm.mu.RLock()
	listeners := make([]ConfigChangeListener, len(cm.listeners))
	copy(listeners, cm.listeners)
	cm.mu.RUnlock()

	for _, listener := range listeners {
		listener.OnConfigChange(configType, version)
	}
}

// getConfigVersion 获取配置版本（从数据库读取最新版本）
func (cm *ConfigManager) getConfigVersion(configType string) int {
	var version int
	err := database.DB.QueryRow(`
		SELECT COALESCE(MAX(version), 0) FROM config_versions WHERE config_type = ?
	`, configType).Scan(&version)
	if err != nil {
		// 如果表不存在或查询失败，返回1作为默认版本
		log.Printf("Warning: failed to get config version for %s: %v", configType, err)
		return 1
	}
	return version
}

// SaveConfigVersion 保存配置版本
func (cm *ConfigManager) SaveConfigVersion(configType string, version int, configData interface{}, description string) error {
	dataJSON, err := json.Marshal(configData)
	if err != nil {
		return fmt.Errorf("failed to marshal config data: %w", err)
	}

	_, err = database.DB.Exec(`
		INSERT INTO config_versions (config_type, version, config_data, description)
		VALUES (?, ?, ?, ?)
	`, configType, version, string(dataJSON), description)
	if err != nil {
		return fmt.Errorf("failed to save config version: %w", err)
	}

	// 更新缓存版本
	cm.mu.Lock()
	if cache, exists := cm.configs[configType]; exists {
		cache.Version = version
		cache.UpdatedAt = time.Now()
	}
	cm.mu.Unlock()

	return nil
}






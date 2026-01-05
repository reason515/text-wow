package game

import (
	"fmt"
	"sync"

	"text-wow/internal/models"
	"text-wow/internal/repository"
)

// ZoneManager 地图管理器 - 管理区域加载、缓存、解锁检查等
type ZoneManager struct {
	mu          sync.RWMutex
	zones       map[string]*models.Zone // 区域缓存
	gameRepo    *repository.GameRepository
	explorationRepo *repository.ExplorationRepository
}

// NewZoneManager 创建地图管理器
func NewZoneManager() *ZoneManager {
	return &ZoneManager{
		zones:          make(map[string]*models.Zone),
		gameRepo:       repository.NewGameRepository(),
		explorationRepo: repository.NewExplorationRepository(),
	}
}

// LoadZone 加载区域（从数据库加载并缓存）
func (zm *ZoneManager) LoadZone(zoneID string) (*models.Zone, error) {
	zm.mu.RLock()
	if zone, exists := zm.zones[zoneID]; exists {
		zm.mu.RUnlock()
		return zone, nil
	}
	zm.mu.RUnlock()

	// 从数据库加载
	zone, err := zm.gameRepo.GetZoneByID(zoneID)
	if err != nil {
		return nil, err
	}

	// 缓存
	zm.mu.Lock()
	zm.zones[zoneID] = zone
	zm.mu.Unlock()

	return zone, nil
}

// GetZone 获取区域（从缓存或数据库）
func (zm *ZoneManager) GetZone(zoneID string) (*models.Zone, error) {
	zm.mu.RLock()
	if zone, exists := zm.zones[zoneID]; exists {
		zm.mu.RUnlock()
		return zone, nil
	}
	zm.mu.RUnlock()

	return zm.LoadZone(zoneID)
}

// GetAllZones 获取所有区域
func (zm *ZoneManager) GetAllZones() ([]*models.Zone, error) {
	zones, err := zm.gameRepo.GetZones()
	if err != nil {
		return nil, err
	}

	// 转换为指针数组并更新缓存
	zonePtrs := make([]*models.Zone, 0, len(zones))
	zm.mu.Lock()
	for i := range zones {
		zonePtrs = append(zonePtrs, &zones[i])
		zm.zones[zones[i].ID] = &zones[i]
	}
	zm.mu.Unlock()

	return zonePtrs, nil
}

// GetZonesByLevel 根据等级获取可用区域
func (zm *ZoneManager) GetZonesByLevel(level int, faction string) ([]*models.Zone, error) {
	allZones, err := zm.GetAllZones()
	if err != nil {
		return nil, err
	}

	available := make([]*models.Zone, 0)
	for _, zone := range allZones {
		// 检查等级范围
		if level >= zone.MinLevel && level <= zone.MaxLevel {
			// 检查阵营
			if zone.Faction == "" || zone.Faction == faction {
				available = append(available, zone)
			}
		}
	}

	return available, nil
}

// IsZoneUnlocked 检查区域是否已解锁
func (zm *ZoneManager) IsZoneUnlocked(userID int, zoneID string) (bool, error) {
	zone, err := zm.GetZone(zoneID)
	if err != nil {
		return false, err
	}

	// 如果没有解锁要求，直接返回true
	if zone.RequiredExploration == 0 && zone.UnlockZoneID == nil {
		return true, nil
	}

	// 检查探索度
	if zm.explorationRepo == nil {
		return false, fmt.Errorf("exploration repository not initialized")
	}

	unlocked, err := zm.explorationRepo.IsZoneUnlocked(userID, zone)
	if err != nil {
		return false, err
	}

	return unlocked, nil
}

// CheckZoneAccess 检查区域访问条件（等级、阵营、解锁状态）
func (zm *ZoneManager) CheckZoneAccess(userID int, zoneID string, playerLevel int, playerFaction string) error {
	zone, err := zm.GetZone(zoneID)
	if err != nil {
		return fmt.Errorf("zone not found: %s", zoneID)
	}

	// 检查等级限制
	if playerLevel < zone.MinLevel {
		return fmt.Errorf("level too low, need level %d", zone.MinLevel)
	}

	if playerLevel > zone.MaxLevel {
		return fmt.Errorf("level too high, max level for this zone is %d", zone.MaxLevel)
	}

	// 检查阵营限制
	if zone.Faction != "" && zone.Faction != playerFaction {
		return fmt.Errorf("faction mismatch, this zone is for %s only", zone.Faction)
	}

	// 检查解锁状态
	unlocked, err := zm.IsZoneUnlocked(userID, zoneID)
	if err != nil {
		return fmt.Errorf("failed to check zone unlock status: %v", err)
	}
	if !unlocked {
		if zone.UnlockZoneID != nil {
			unlockZone, _ := zm.GetZone(*zone.UnlockZoneID)
			if unlockZone != nil {
				exploration, _ := zm.explorationRepo.GetExploration(userID, *zone.UnlockZoneID)
				return fmt.Errorf("zone locked: need %d exploration in %s (current: %d)", zone.RequiredExploration, unlockZone.Name, exploration.Exploration)
			}
		}
		return fmt.Errorf("zone locked: need %d exploration", zone.RequiredExploration)
	}

	return nil
}

// GetZoneInfo 获取区域信息（包括解锁状态）
func (zm *ZoneManager) GetZoneInfo(userID int, zoneID string) (*ZoneInfo, error) {
	zone, err := zm.GetZone(zoneID)
	if err != nil {
		return nil, err
	}

	unlocked, _ := zm.IsZoneUnlocked(userID, zoneID)

	info := &ZoneInfo{
		Zone:    zone,
		Unlocked: unlocked,
	}

	// 如果未解锁，获取解锁进度
	if !unlocked && zm.explorationRepo != nil {
		if zone.UnlockZoneID != nil {
			exploration, _ := zm.explorationRepo.GetExploration(userID, *zone.UnlockZoneID)
			info.RequiredExploration = zone.RequiredExploration
			info.CurrentExploration = exploration.Exploration
			unlockZone, _ := zm.GetZone(*zone.UnlockZoneID)
			if unlockZone != nil {
				info.UnlockZoneName = unlockZone.Name
			}
		}
	}

	return info, nil
}

// ZoneInfo 区域信息（包含解锁状态）
type ZoneInfo struct {
	Zone                *models.Zone `json:"zone"`
	Unlocked            bool          `json:"unlocked"`
	RequiredExploration int           `json:"requiredExploration,omitempty"`
	CurrentExploration  int           `json:"currentExploration,omitempty"`
	UnlockZoneName      string        `json:"unlockZoneName,omitempty"`
}

// CalculateExpMultiplier 计算经验倍率
func (zm *ZoneManager) CalculateExpMultiplier(zoneID string) float64 {
	zone, err := zm.GetZone(zoneID)
	if err != nil {
		return 1.0 // 默认倍率
	}
	if zone.ExpMulti <= 0 {
		return 1.0
	}
	return zone.ExpMulti
}

// CalculateGoldMultiplier 计算金币倍率
func (zm *ZoneManager) CalculateGoldMultiplier(zoneID string) float64 {
	zone, err := zm.GetZone(zoneID)
	if err != nil {
		return 1.0 // 默认倍率
	}
	if zone.GoldMulti <= 0 {
		return 1.0
	}
	return zone.GoldMulti
}

// CalculateDropMultiplier 计算掉落倍率（与金币倍率相同）
func (zm *ZoneManager) CalculateDropMultiplier(zoneID string) float64 {
	// 掉落倍率通常与金币倍率相同
	return zm.CalculateGoldMultiplier(zoneID)
}

// ReloadZone 重新加载区域（用于热更新）
func (zm *ZoneManager) ReloadZone(zoneID string) error {
	zm.mu.Lock()
	delete(zm.zones, zoneID)
	zm.mu.Unlock()

	_, err := zm.LoadZone(zoneID)
	return err
}

// ReloadAllZones 重新加载所有区域（用于热更新）
func (zm *ZoneManager) ReloadAllZones() error {
	zm.mu.Lock()
	zm.zones = make(map[string]*models.Zone)
	zm.mu.Unlock()

	_, err := zm.GetAllZones()
	return err
}


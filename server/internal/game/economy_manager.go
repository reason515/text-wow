package game

import (
	"fmt"
	"sync"

	"text-wow/internal/repository"
)

// EconomyManager 经济管理器 - 管理金币获取与消耗
type EconomyManager struct {
	mu       sync.RWMutex
	userRepo *repository.UserRepository
	config   *EconomyConfig
}

// EconomyConfig 经济配置
type EconomyConfig struct {
	GoldMultiplier           float64 // 金币获取倍率
	MaterialPriceMultiplier  float64 // 材料价格倍率
	EquipmentPriceMultiplier float64 // 装备价格倍率
}

// NewEconomyManager 创建经济管理器
func NewEconomyManager() *EconomyManager {
	return &EconomyManager{
		userRepo: repository.NewUserRepository(),
		config: &EconomyConfig{
			GoldMultiplier:           1.0,
			MaterialPriceMultiplier:  1.0,
			EquipmentPriceMultiplier: 1.0,
		},
	}
}

// CalculateGoldReward 计算金币奖励
// 公式: 基础金币 × 区域倍率 × 难度倍率 × 全局倍率
func (em *EconomyManager) CalculateGoldReward(baseGold int, zoneMultiplier float64, difficultyMultiplier float64) int {
	totalMultiplier := zoneMultiplier * difficultyMultiplier * em.config.GoldMultiplier
	return int(float64(baseGold) * totalMultiplier)
}

// CalculateMaterialPrice 计算材料价格
func (em *EconomyManager) CalculateMaterialPrice(basePrice int, rarity string) int {
	rarityMultiplier := 1.0
	switch rarity {
	case "common":
		rarityMultiplier = 1.0
	case "uncommon":
		rarityMultiplier = 1.5
	case "rare":
		rarityMultiplier = 2.0
	case "epic":
		rarityMultiplier = 3.0
	}

	return int(float64(basePrice) * rarityMultiplier * em.config.MaterialPriceMultiplier)
}

// AddGold 增加金币
func (em *EconomyManager) AddGold(userID int, amount int) error {
	user, err := em.userRepo.GetByID(userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	user.Gold += amount
	user.TotalGoldGained += amount

	// 使用 UpdateGold 方法更新金币
	return em.userRepo.UpdateGold(userID, amount)
}

// SpendGold 消耗金币
func (em *EconomyManager) SpendGold(userID int, amount int) error {
	user, err := em.userRepo.GetByID(userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	if user.Gold < amount {
		return fmt.Errorf("insufficient gold: have %d, need %d", user.Gold, amount)
	}

	// 使用 UpdateGold 方法更新金币（传入负数表示减少）
	return em.userRepo.UpdateGold(userID, -amount)
}

// GetGold 获取金币
func (em *EconomyManager) GetGold(userID int) (int, error) {
	user, err := em.userRepo.GetByID(userID)
	if err != nil {
		return 0, fmt.Errorf("failed to get user: %w", err)
	}

	return user.Gold, nil
}

// UpdateConfig 更新经济配置（热更新）
func (em *EconomyManager) UpdateConfig(config *EconomyConfig) {
	em.mu.Lock()
	defer em.mu.Unlock()

	em.config = config
}

// GetConfig 获取经济配置
func (em *EconomyManager) GetConfig() *EconomyConfig {
	em.mu.RLock()
	defer em.mu.RUnlock()

	return em.config
}


package game

import (
	"fmt"
	"sync"
	"time"

	"text-wow/internal/repository"
)

// TradingManager 交易管理器 - 管理玩家间装备交易和拍卖行
type TradingManager struct {
	mu          sync.RWMutex
	userRepo    *repository.UserRepository
	gameRepo    *repository.GameRepository
	economyMgr  *EconomyManager
	equipmentMgr *EquipmentManager
}

// AuctionListing 拍卖行上架信息
type AuctionListing struct {
	ID          int
	SellerID    int
	EquipmentID int
	ItemID      string
	Price       int
	ListedAt    time.Time
	ExpiresAt   time.Time
	Status      string // active/sold/expired/cancelled
	BuyerID     *int
	SoldAt      *time.Time
}

// TradeOffer 交易报价
type TradeOffer struct {
	ID          int
	FromUserID  int
	ToUserID    int
	EquipmentID int
	Price       int
	Status      string // pending/accepted/rejected/cancelled
	CreatedAt   time.Time
	ExpiresAt   time.Time
}

// NewTradingManager 创建交易管理器
func NewTradingManager() *TradingManager {
	return &TradingManager{
		userRepo:     repository.NewUserRepository(),
		gameRepo:     repository.NewGameRepository(),
		economyMgr:   NewEconomyManager(),
		equipmentMgr: NewEquipmentManager(),
	}
}

// ListItem 上架装备到拍卖行
// 功能：玩家将装备上架到拍卖行，设置价格，其他玩家可以购买
func (tm *TradingManager) ListItem(sellerID int, equipmentID int, price int, durationDays int) error {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	// 验证价格
	if price <= 0 {
		return fmt.Errorf("price must be greater than 0")
	}

	// 验证装备所有权
	equipmentRepo := repository.NewEquipmentRepository()
	equipment, err := equipmentRepo.GetByID(equipmentID)
	if err != nil {
		return fmt.Errorf("failed to get equipment: %w", err)
	}
	if equipment.OwnerID != sellerID {
		return fmt.Errorf("equipment does not belong to seller")
	}
	if equipment.CharacterID != nil {
		return fmt.Errorf("equipment is currently equipped, please unequip first")
	}

	// 计算上架费（价格的5%）
	listingFee := int(float64(price) * 0.05)
	if listingFee < 10 {
		listingFee = 10 // 最低上架费10金币
	}

	// 检查卖家是否有足够金币支付上架费
	seller, err := tm.userRepo.GetByID(sellerID)
	if err != nil {
		return fmt.Errorf("failed to get seller: %w", err)
	}
	if seller.Gold < listingFee {
		return fmt.Errorf("insufficient gold for listing fee: need %d, have %d", listingFee, seller.Gold)
	}

	// 扣除上架费
	if err := tm.economyMgr.SpendGold(sellerID, listingFee); err != nil {
		return fmt.Errorf("failed to pay listing fee: %w", err)
	}

	// 创建上架记录
	// TODO: 需要实现 AuctionRepository
	// listing := &AuctionListing{
	// 	SellerID:    sellerID,
	// 	EquipmentID: equipmentID,
	// 	Price:       price,
	// 	ListedAt:    time.Now(),
	// 	ExpiresAt:   time.Now().AddDate(0, 0, durationDays),
	// 	Status:      "active",
	// }
	// if err := tm.auctionRepo.CreateListing(listing); err != nil {
	// 	// 如果创建失败，退还上架费
	// 	tm.economyMgr.AddGold(sellerID, listingFee)
	// 	return fmt.Errorf("failed to create listing: %w", err)
	// }

	return fmt.Errorf("AuctionRepository not implemented yet")
}

// BuyItem 从拍卖行购买装备
// 功能：玩家从拍卖行购买装备，扣除金币，转移装备所有权
func (tm *TradingManager) BuyItem(buyerID int, listingID int) error {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	// 获取上架信息
	// TODO: 需要实现 AuctionRepository
	// listing, err := tm.auctionRepo.GetListingByID(listingID)
	// if err != nil {
	// 	return fmt.Errorf("failed to get listing: %w", err)
	// }
	// if listing.Status != "active" {
	// 	return fmt.Errorf("listing is not active")
	// }
	// if listing.ExpiresAt.Before(time.Now()) {
	// 	return fmt.Errorf("listing has expired")
	// }
	// if listing.SellerID == buyerID {
	// 	return fmt.Errorf("cannot buy your own item")
	// }

	// TODO: 实现完整的购买逻辑
	// 检查买家是否有足够金币
	_, err := tm.userRepo.GetByID(buyerID)
	if err != nil {
		return fmt.Errorf("failed to get buyer: %w", err)
	}
	
	// 获取拍卖列表
	// listing, err := tm.auctionRepo.GetListingByID(listingID)
	// if err != nil {
	// 	return fmt.Errorf("failed to get listing: %w", err)
	// }
	
	// if buyer.Gold < listing.Price {
	// 	return fmt.Errorf("insufficient gold: need %d, have %d", listing.Price, buyer.Gold)
	// }

	// 扣除买家金币
	// if err := tm.economyMgr.SpendGold(buyerID, listing.Price); err != nil {
	// 	return fmt.Errorf("failed to deduct gold: %w", err)
	// }

	// 计算交易手续费（价格的5%）
	// transactionFee := int(float64(listing.Price) * 0.05)
	// sellerGold := listing.Price - transactionFee

	// 转移装备所有权
	// equipmentRepo := repository.NewEquipmentRepository()
	// equipment, err := equipmentRepo.GetByID(listing.EquipmentID)
	// if err != nil {
	// 	// 如果转移失败，退还金币
	// 	tm.economyMgr.AddGold(buyerID, listing.Price)
	// 	return fmt.Errorf("failed to get equipment: %w", err)
	// }
	// equipment.OwnerID = buyerID
	// equipment.CharacterID = nil // 确保装备未装备
	// if err := equipmentRepo.Update(equipment); err != nil {
	// 	// 如果更新失败，退还金币
	// 	tm.economyMgr.AddGold(buyerID, listing.Price)
	// 	return fmt.Errorf("failed to transfer equipment: %w", err)
	// }
	
	return fmt.Errorf("BuyFromAuction not fully implemented yet")

	// 支付卖家金币（扣除手续费）
	// if err := tm.economyMgr.AddGold(listing.SellerID, sellerGold); err != nil {
	// 	// 如果支付失败，记录错误但不回滚（装备已转移）
	// 	fmt.Printf("[ERROR] Failed to pay seller: %v\n", err)
	// }

	// 更新上架状态
	// listing.Status = "sold"
	// listing.BuyerID = &buyerID
	// now := time.Now()
	// listing.SoldAt = &now
	// if err := tm.auctionRepo.UpdateListing(listing); err != nil {
	// 	fmt.Printf("[ERROR] Failed to update listing status: %v\n", err)
	// }

	return fmt.Errorf("AuctionRepository and EquipmentRepository not implemented yet")
}

// CancelListing 取消上架
// 功能：卖家取消上架，装备归还，上架费不退还
func (tm *TradingManager) CancelListing(sellerID int, listingID int) error {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	// 获取上架信息
	// TODO: 需要实现 AuctionRepository
	// listing, err := tm.auctionRepo.GetListingByID(listingID)
	// if err != nil {
	// 	return fmt.Errorf("failed to get listing: %w", err)
	// }
	// if listing.SellerID != sellerID {
	// 	return fmt.Errorf("listing does not belong to seller")
	// }
	// if listing.Status != "active" {
	// 	return fmt.Errorf("listing is not active")
	// }

	// 更新上架状态
	// listing.Status = "cancelled"
	// if err := tm.auctionRepo.UpdateListing(listing); err != nil {
	// 	return fmt.Errorf("failed to cancel listing: %w", err)
	// }

	// 注意：上架费不退还，这是为了防止玩家频繁上架/取消

	return fmt.Errorf("AuctionRepository not implemented yet")
}

// CreateDirectTrade 创建直接交易
// 功能：两个玩家之间的直接交易，需要双方确认
func (tm *TradingManager) CreateDirectTrade(fromUserID int, toUserID int, equipmentID int, price int) error {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	// 验证交易双方
	if fromUserID == toUserID {
		return fmt.Errorf("cannot trade with yourself")
	}

	// 验证价格
	if price <= 0 {
		return fmt.Errorf("price must be greater than 0")
	}

	// 验证装备所有权
	// TODO: 需要实现 GetEquipmentByID
	// equipment, err := tm.equipmentMgr.GetEquipmentByID(equipmentID)
	// if err != nil {
	// 	return fmt.Errorf("failed to get equipment: %w", err)
	// }
	// if equipment.OwnerID != fromUserID {
	// 	return fmt.Errorf("equipment does not belong to seller")
	// }

	// 创建交易报价
	// offer := &TradeOffer{
	// 	FromUserID:  fromUserID,
	// 	ToUserID:    toUserID,
	// 	EquipmentID: equipmentID,
	// 	Price:       price,
	// 	Status:      "pending",
	// 	CreatedAt:   time.Now(),
	// 	ExpiresAt:   time.Now().AddDate(0, 0, 7), // 7天过期
	// }
	// if err := tm.tradeRepo.CreateOffer(offer); err != nil {
	// 	return fmt.Errorf("failed to create trade offer: %w", err)
	// }

	return fmt.Errorf("TradeRepository not implemented yet")
}

// AcceptTrade 接受交易
// 功能：买家接受交易报价，完成交易
func (tm *TradingManager) AcceptTrade(buyerID int, offerID int) error {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	// 获取交易报价
	// TODO: 需要实现 TradeRepository
	// offer, err := tm.tradeRepo.GetOfferByID(offerID)
	// if err != nil {
	// 	return fmt.Errorf("failed to get trade offer: %w", err)
	// }
	// if offer.ToUserID != buyerID {
	// 	return fmt.Errorf("trade offer does not belong to buyer")
	// }
	// if offer.Status != "pending" {
	// 	return fmt.Errorf("trade offer is not pending")
	// }
	// if offer.ExpiresAt.Before(time.Now()) {
	// 	return fmt.Errorf("trade offer has expired")
	// }

	// 检查买家是否有足够金币
	// buyer, err := tm.userRepo.GetByID(buyerID)
	// if err != nil {
	// 	return fmt.Errorf("failed to get buyer: %w", err)
	// }
	// if buyer.Gold < offer.Price {
	// 	return fmt.Errorf("insufficient gold: need %d, have %d", offer.Price, buyer.Gold)
	// }

	// 扣除买家金币
	// if err := tm.economyMgr.SpendGold(buyerID, offer.Price); err != nil {
	// 	return fmt.Errorf("failed to deduct gold: %w", err)
	// }

	// 计算交易手续费（价格的5%）
	// transactionFee := int(float64(offer.Price) * 0.05)
	// sellerGold := offer.Price - transactionFee

	// 转移装备所有权
	// equipment, err := tm.equipmentMgr.GetEquipmentByID(offer.EquipmentID)
	// if err != nil {
	// 	// 如果转移失败，退还金币
	// 	tm.economyMgr.AddGold(buyerID, offer.Price)
	// 	return fmt.Errorf("failed to get equipment: %w", err)
	// }
	// equipment.OwnerID = buyerID
	// equipment.CharacterID = nil
	// if err := tm.equipmentRepo.UpdateEquipment(equipment); err != nil {
	// 	// 如果更新失败，退还金币
	// 	tm.economyMgr.AddGold(buyerID, offer.Price)
	// 	return fmt.Errorf("failed to transfer equipment: %w", err)
	// }

	// 支付卖家金币（扣除手续费）
	// if err := tm.economyMgr.AddGold(offer.FromUserID, sellerGold); err != nil {
	// 	fmt.Printf("[ERROR] Failed to pay seller: %v\n", err)
	// }

	// 更新交易状态
	// offer.Status = "accepted"
	// if err := tm.tradeRepo.UpdateOffer(offer); err != nil {
	// 	fmt.Printf("[ERROR] Failed to update trade offer: %v\n", err)
	// }

	return fmt.Errorf("TradeRepository and EquipmentRepository not implemented yet")
}

// RejectTrade 拒绝交易
// 功能：买家拒绝交易报价
func (tm *TradingManager) RejectTrade(buyerID int, offerID int) error {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	// 获取交易报价
	// TODO: 需要实现 TradeRepository
	// offer, err := tm.tradeRepo.GetOfferByID(offerID)
	// if err != nil {
	// 	return fmt.Errorf("failed to get trade offer: %w", err)
	// }
	// if offer.ToUserID != buyerID {
	// 	return fmt.Errorf("trade offer does not belong to buyer")
	// }
	// if offer.Status != "pending" {
	// 	return fmt.Errorf("trade offer is not pending")
	// }

	// 更新交易状态
	// offer.Status = "rejected"
	// if err := tm.tradeRepo.UpdateOffer(offer); err != nil {
	// 	return fmt.Errorf("failed to reject trade offer: %w", err)
	// }

	return fmt.Errorf("TradeRepository not implemented yet")
}

// GetActiveListings 获取活跃的上架列表
// 功能：获取拍卖行中所有活跃的上架装备
func (tm *TradingManager) GetActiveListings(filters map[string]interface{}) ([]*AuctionListing, error) {
	tm.mu.RLock()
	defer tm.mu.RUnlock()

	// TODO: 需要实现 AuctionRepository
	// listings, err := tm.auctionRepo.GetActiveListings(filters)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to get listings: %w", err)
	// }

	// return listings, nil
	return nil, fmt.Errorf("AuctionRepository not implemented yet")
}

// GetUserListings 获取用户的上架列表
// 功能：获取指定用户的所有上架装备
func (tm *TradingManager) GetUserListings(userID int) ([]*AuctionListing, error) {
	tm.mu.RLock()
	defer tm.mu.RUnlock()

	// TODO: 需要实现 AuctionRepository
	// listings, err := tm.auctionRepo.GetListingsBySellerID(userID)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to get user listings: %w", err)
	// }

	// return listings, nil
	return nil, fmt.Errorf("AuctionRepository not implemented yet")
}

// GetUserTradeOffers 获取用户的交易报价
// 功能：获取指定用户收到和发出的交易报价
func (tm *TradingManager) GetUserTradeOffers(userID int) (sent []*TradeOffer, received []*TradeOffer, err error) {
	tm.mu.RLock()
	defer tm.mu.RUnlock()

	// TODO: 需要实现 TradeRepository
	// sent, err = tm.tradeRepo.GetOffersByFromUserID(userID)
	// if err != nil {
	// 	return nil, nil, fmt.Errorf("failed to get sent offers: %w", err)
	// }

	// received, err = tm.tradeRepo.GetOffersByToUserID(userID)
	// if err != nil {
	// 	return nil, nil, fmt.Errorf("failed to get received offers: %w", err)
	// }

	// return sent, received, nil
	return nil, nil, fmt.Errorf("TradeRepository not implemented yet")
}

// CalculateTransactionFee 计算交易手续费
// 功能：计算交易手续费（价格的5%）
func (tm *TradingManager) CalculateTransactionFee(price int) int {
	return int(float64(price) * 0.05)
}

// CalculateListingFee 计算上架费
// 功能：计算上架费（价格的5%，最低10金币）
func (tm *TradingManager) CalculateListingFee(price int) int {
	fee := int(float64(price) * 0.05)
	if fee < 10 {
		fee = 10
	}
	return fee
}


package game

import (
	"fmt"
	"sync"

	"text-wow/internal/models"
	"text-wow/internal/repository"
)

// TeamManager 队伍管理器 - 管理玩家队伍（1-5人）
type TeamManager struct {
	mu       sync.RWMutex
	charRepo *repository.CharacterRepository
	calculator *Calculator
}

// Team 队伍信息
type Team struct {
	UserID        int
	MaxSize       int          // 最大队伍人数 (1-5)
	UnlockedSlots int          // 已解锁槽位数
	Characters    []*models.Character // 队伍中的角色（按槽位排序）
}

// NewTeamManager 创建队伍管理器
func NewTeamManager() *TeamManager {
	return &TeamManager{
		charRepo:   repository.NewCharacterRepository(),
		calculator: NewCalculator(),
	}
}

// GetTeam 获取玩家队伍
func (tm *TeamManager) GetTeam(userID int) (*Team, error) {
	characters, err := tm.charRepo.GetByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get characters: %w", err)
	}

	// 按槽位排序
	sortedChars := make([]*models.Character, 5)
	for _, char := range characters {
		if char.TeamSlot >= 1 && char.TeamSlot <= 5 {
			sortedChars[char.TeamSlot-1] = char
		}
	}

	// 移除空槽位
	teamChars := make([]*models.Character, 0, 5)
	for _, char := range sortedChars {
		if char != nil {
			teamChars = append(teamChars, char)
		}
	}

	// 获取用户信息以确定最大队伍人数
	userRepo := repository.NewUserRepository()
	user, err := userRepo.GetByID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	
	maxSize := user.MaxTeamSize
	if maxSize == 0 {
		maxSize = 5 // 默认值
	}
	
	unlockedSlots := user.UnlockedSlots
	if unlockedSlots == 0 {
		unlockedSlots = len(teamChars)
		if unlockedSlots == 0 {
			unlockedSlots = 1 // 至少解锁1个槽位
		}
	}

	return &Team{
		UserID:        userID,
		MaxSize:         maxSize,
		UnlockedSlots:  unlockedSlots,
		Characters:     teamChars,
	}, nil
}

// AddCharacterToTeam 添加角色到队伍
func (tm *TeamManager) AddCharacterToTeam(userID int, characterID int, slot int) error {
	if slot < 1 || slot > 5 {
		return fmt.Errorf("invalid slot number: %d (must be 1-5)", slot)
	}

	// 检查槽位是否已解锁
	team, err := tm.GetTeam(userID)
	if err != nil {
		return err
	}

	if slot > team.UnlockedSlots {
		return fmt.Errorf("slot %d is not unlocked (unlocked: %d)", slot, team.UnlockedSlots)
	}

	// 检查槽位是否已被占用
	for _, char := range team.Characters {
		if char.TeamSlot == slot {
			return fmt.Errorf("slot %d is already occupied by character %d", slot, char.ID)
		}
	}

	// 获取角色
	char, err := tm.charRepo.GetByID(characterID)
	if err != nil {
		return fmt.Errorf("failed to get character: %w", err)
	}

	// 验证角色所有权
	if char.UserID != userID {
		return fmt.Errorf("character does not belong to user")
	}

	// 更新角色的槽位
	char.TeamSlot = slot
	char.IsActive = true

	return tm.charRepo.Update(char)
}

// RemoveCharacterFromTeam 从队伍中移除角色
func (tm *TeamManager) RemoveCharacterFromTeam(userID int, slot int) error {
	team, err := tm.GetTeam(userID)
	if err != nil {
		return err
	}

	var charToRemove *models.Character
	for _, char := range team.Characters {
		if char.TeamSlot == slot {
			charToRemove = char
			break
		}
	}

	if charToRemove == nil {
		return fmt.Errorf("no character in slot %d", slot)
	}

	charToRemove.TeamSlot = 0
	charToRemove.IsActive = false

	return tm.charRepo.Update(charToRemove)
}

// GetActiveCharacters 获取活跃角色列表（用于战斗）
func (tm *TeamManager) GetActiveCharacters(userID int) ([]*models.Character, error) {
	team, err := tm.GetTeam(userID)
	if err != nil {
		return nil, err
	}

	activeChars := make([]*models.Character, 0)
	for _, char := range team.Characters {
		if char.IsActive && !char.IsDead {
			activeChars = append(activeChars, char)
		}
	}

	return activeChars, nil
}

// CalculateTeamAttributes 计算队伍总属性
func (tm *TeamManager) CalculateTeamAttributes(team *Team) *TeamAttributes {
	attrs := &TeamAttributes{
		TotalHP:           0,
		TotalPhysicalAttack: 0,
		TotalMagicAttack:  0,
		TotalPhysicalDefense: 0,
		TotalMagicDefense: 0,
	}

	for _, char := range team.Characters {
		if char.IsActive && !char.IsDead {
			attrs.TotalHP += char.MaxHP
			attrs.TotalPhysicalAttack += tm.calculator.CalculatePhysicalAttack(char)
			attrs.TotalMagicAttack += tm.calculator.CalculateMagicAttack(char)
			attrs.TotalPhysicalDefense += char.PhysicalDefense
			attrs.TotalMagicDefense += char.MagicDefense
		}
	}

	return attrs
}

// TeamAttributes 队伍属性
type TeamAttributes struct {
	TotalHP            int
	TotalPhysicalAttack int
	TotalMagicAttack   int
	TotalPhysicalDefense int
	TotalMagicDefense  int
}


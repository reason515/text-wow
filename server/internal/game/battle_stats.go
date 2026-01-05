package game

import (
	"fmt"
	"sync"
	"time"

	"text-wow/internal/models"
	"text-wow/internal/repository"
)

// BattleStatsCollector 战斗统计收集器
type BattleStatsCollector struct {
	mu              sync.RWMutex
	battleStatsRepo *repository.BattleStatsRepository
	stats           map[int]*CharacterBattleStats // key: characterID
}

// CharacterBattleStats 角色战斗统计
type CharacterBattleStats struct {
	CharacterID int
	TeamSlot    int

	// 伤害统计
	DamageDealt    int
	PhysicalDamage int
	MagicDamage    int
	FireDamage     int
	FrostDamage    int
	ShadowDamage   int
	HolyDamage     int
	NatureDamage   int
	DotDamage      int

	// 暴击统计
	CritCount  int
	CritDamage int
	MaxCrit    int

	// 承伤统计
	DamageTaken    int
	PhysicalTaken  int
	MagicTaken     int
	DamageBlocked  int
	DamageAbsorbed int

	// 闪避统计
	DodgeCount int
	BlockCount int
	HitCount   int

	// 治疗统计
	HealingDone     int
	HealingReceived int
	Overhealing     int
	SelfHealing     int
	HotHealing      int

	// 技能统计
	SkillUses   int
	SkillHits   int
	SkillMisses int

	// 控制统计
	CcApplied  int
	CcReceived int
	Dispels    int
	Interrupts int

	// 其他统计
	Kills             int
	Deaths            int
	Resurrects        int
	ResourceUsed      int
	ResourceGenerated int
}

// NewBattleStatsCollector 创建战斗统计收集器
func NewBattleStatsCollector() *BattleStatsCollector {
	return &BattleStatsCollector{
		battleStatsRepo: repository.NewBattleStatsRepository(),
		stats:           make(map[int]*CharacterBattleStats),
	}
}

// InitializeBattle 初始化战斗统计
func (bsc *BattleStatsCollector) InitializeBattle(characterIDs []int) {
	bsc.mu.Lock()
	defer bsc.mu.Unlock()

	for _, charID := range characterIDs {
		bsc.stats[charID] = &CharacterBattleStats{
			CharacterID: charID,
		}
	}
}

// RecordDamage 记录伤害
func (bsc *BattleStatsCollector) RecordDamage(characterID int, damage int, damageType string, isCrit bool) {
	bsc.mu.Lock()
	defer bsc.mu.Unlock()

	stats, exists := bsc.stats[characterID]
	if !exists {
		return
	}

	stats.DamageDealt += damage
	stats.HitCount++

	switch damageType {
	case "physical":
		stats.PhysicalDamage += damage
	case "magic":
		stats.MagicDamage += damage
	case "fire":
		stats.FireDamage += damage
	case "frost":
		stats.FrostDamage += damage
	case "shadow":
		stats.ShadowDamage += damage
	case "holy":
		stats.HolyDamage += damage
	case "nature":
		stats.NatureDamage += damage
	}

	if isCrit {
		stats.CritCount++
		stats.CritDamage += damage
		if damage > stats.MaxCrit {
			stats.MaxCrit = damage
		}
	}
}

// RecordDamageTaken 记录承伤
func (bsc *BattleStatsCollector) RecordDamageTaken(characterID int, damage int, damageType string, isDodged bool, isBlocked bool) {
	bsc.mu.Lock()
	defer bsc.mu.Unlock()

	stats, exists := bsc.stats[characterID]
	if !exists {
		return
	}

	if isDodged {
		stats.DodgeCount++
		return
	}

	if isBlocked {
		stats.BlockCount++
		stats.DamageBlocked += damage
		return
	}

	stats.DamageTaken += damage
	stats.HitCount++

	switch damageType {
	case "physical":
		stats.PhysicalTaken += damage
	case "magic":
		stats.MagicTaken += damage
	}
}

// RecordHealing 记录治疗
func (bsc *BattleStatsCollector) RecordHealing(healerID int, targetID int, healing int, overhealing int, isSelf bool) {
	bsc.mu.Lock()
	defer bsc.mu.Unlock()

	stats, exists := bsc.stats[healerID]
	if !exists {
		return
	}

	stats.HealingDone += healing
	stats.Overhealing += overhealing

	if isSelf {
		stats.SelfHealing += healing
	}

	// 记录目标受到的治疗
	if targetStats, exists := bsc.stats[targetID]; exists {
		targetStats.HealingReceived += healing
	}
}

// RecordSkillUse 记录技能使用
func (bsc *BattleStatsCollector) RecordSkillUse(characterID int, skillID string, hit bool, resourceCost int) {
	bsc.mu.Lock()
	defer bsc.mu.Unlock()

	stats, exists := bsc.stats[characterID]
	if !exists {
		return
	}

	stats.SkillUses++
	stats.ResourceUsed += resourceCost

	if hit {
		stats.SkillHits++
	} else {
		stats.SkillMisses++
	}
}

// RecordKill 记录击杀
func (bsc *BattleStatsCollector) RecordKill(characterID int) {
	bsc.mu.Lock()
	defer bsc.mu.Unlock()

	if stats, exists := bsc.stats[characterID]; exists {
		stats.Kills++
	}
}

// RecordDeath 记录死亡
func (bsc *BattleStatsCollector) RecordDeath(characterID int) {
	bsc.mu.Lock()
	defer bsc.mu.Unlock()

	if stats, exists := bsc.stats[characterID]; exists {
		stats.Deaths++
	}
}

// GetStats 获取角色统计
func (bsc *BattleStatsCollector) GetStats(characterID int) *CharacterBattleStats {
	bsc.mu.RLock()
	defer bsc.mu.RUnlock()

	return bsc.stats[characterID]
}

// GetAllStats 获取所有角色统计
func (bsc *BattleStatsCollector) GetAllStats() map[int]*CharacterBattleStats {
	bsc.mu.RLock()
	defer bsc.mu.RUnlock()

	result := make(map[int]*CharacterBattleStats)
	for k, v := range bsc.stats {
		result[k] = v
	}
	return result
}

// SaveBattleStats 保存战斗统计到数据库
func (bsc *BattleStatsCollector) SaveBattleStats(battleID int, duration time.Duration, totalRounds int) error {
	bsc.mu.RLock()
	defer bsc.mu.RUnlock()

	for _, stats := range bsc.stats {
		// 保存角色战斗统计
		battleCharStats := &models.BattleCharacterStats{
			BattleID:    battleID,
			CharacterID: stats.CharacterID,
			TeamSlot:    stats.TeamSlot,
			DamageDealt: stats.DamageDealt,
			PhysicalDamage: stats.PhysicalDamage,
			MagicDamage:    stats.MagicDamage,
			FireDamage:     stats.FireDamage,
			FrostDamage:    stats.FrostDamage,
			ShadowDamage:   stats.ShadowDamage,
			HolyDamage:     stats.HolyDamage,
			NatureDamage:   stats.NatureDamage,
			DotDamage:      stats.DotDamage,
			CritCount:      stats.CritCount,
			CritDamage:     stats.CritDamage,
			MaxCrit:        stats.MaxCrit,
			DamageTaken:    stats.DamageTaken,
			PhysicalTaken:  stats.PhysicalTaken,
			MagicTaken:     stats.MagicTaken,
			DamageBlocked:  stats.DamageBlocked,
			DamageAbsorbed: stats.DamageAbsorbed,
			DodgeCount:     stats.DodgeCount,
			BlockCount:     stats.BlockCount,
			HitCount:       stats.HitCount,
			HealingDone:    stats.HealingDone,
			HealingReceived: stats.HealingReceived,
			Overhealing:    stats.Overhealing,
			SelfHealing:    stats.SelfHealing,
			HotHealing:     stats.HotHealing,
			SkillUses:      stats.SkillUses,
			SkillHits:      stats.SkillHits,
			SkillMisses:    stats.SkillMisses,
			CcApplied:      stats.CcApplied,
			CcReceived:    stats.CcReceived,
			Dispels:        stats.Dispels,
			Interrupts:     stats.Interrupts,
			Kills:          stats.Kills,
			Deaths:         stats.Deaths,
			Resurrects:     stats.Resurrects,
			ResourceUsed:   stats.ResourceUsed,
			ResourceGenerated: stats.ResourceGenerated,
		}

		if _, err := bsc.battleStatsRepo.CreateBattleCharacterStats(battleCharStats); err != nil {
			return fmt.Errorf("failed to save character stats: %w", err)
		}
	}

	return nil
}

// ClearStats 清空统计（战斗结束后）
func (bsc *BattleStatsCollector) ClearStats() {
	bsc.mu.Lock()
	defer bsc.mu.Unlock()

	bsc.stats = make(map[int]*CharacterBattleStats)
}

// StatsAnalyzer 统计分析器
type StatsAnalyzer struct {
	collector *BattleStatsCollector
}

// NewStatsAnalyzer 创建统计分析器
func NewStatsAnalyzer(collector *BattleStatsCollector) *StatsAnalyzer {
	return &StatsAnalyzer{
		collector: collector,
	}
}

// CalculateDPS 计算DPS
func (sa *StatsAnalyzer) CalculateDPS(characterID int, duration time.Duration) float64 {
	stats := sa.collector.GetStats(characterID)
	if stats == nil || duration.Seconds() == 0 {
		return 0
	}

	return float64(stats.DamageDealt) / duration.Seconds()
}

// CalculateHPS 计算HPS
func (sa *StatsAnalyzer) CalculateHPS(characterID int, duration time.Duration) float64 {
	stats := sa.collector.GetStats(characterID)
	if stats == nil || duration.Seconds() == 0 {
		return 0
	}

	return float64(stats.HealingDone) / duration.Seconds()
}

// GetDamageDistribution 获取伤害分布
func (sa *StatsAnalyzer) GetDamageDistribution(characterID int) map[string]int {
	stats := sa.collector.GetStats(characterID)
	if stats == nil {
		return nil
	}

	return map[string]int{
		"physical": stats.PhysicalDamage,
		"magic":    stats.MagicDamage,
		"fire":     stats.FireDamage,
		"frost":    stats.FrostDamage,
		"shadow":   stats.ShadowDamage,
		"holy":     stats.HolyDamage,
		"nature":   stats.NatureDamage,
		"dot":      stats.DotDamage,
		"total":    stats.DamageDealt,
	}
}

// GetActualCritRate 获取实际暴击率
func (sa *StatsAnalyzer) GetActualCritRate(characterID int) float64 {
	stats := sa.collector.GetStats(characterID)
	if stats == nil || stats.HitCount == 0 {
		return 0
	}

	return float64(stats.CritCount) / float64(stats.HitCount)
}

// GetActualDodgeRate 获取实际闪避率
func (sa *StatsAnalyzer) GetActualDodgeRate(characterID int) float64 {
	stats := sa.collector.GetStats(characterID)
	if stats == nil {
		return 0
	}

	totalAttacks := stats.DodgeCount + stats.HitCount
	if totalAttacks == 0 {
		return 0
	}

	return float64(stats.DodgeCount) / float64(totalAttacks)
}


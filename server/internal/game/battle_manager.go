package game

import (
	"fmt"
	"math"
	"math/rand"
	"sort"
	"strings"
	"sync"
	"time"

	"text-wow/internal/models"
	"text-wow/internal/repository"
)

// BattleManager 战斗管理器 - 管理所有用户的战斗状态
type BattleManager struct {
	mu                  sync.RWMutex
	sessions            map[int]*BattleSession // key: userID
	gameRepo            *repository.GameRepository
	charRepo            *repository.CharacterRepository
	explorationRepo     *repository.ExplorationRepository // 探索度仓库
	inventoryRepo       *repository.InventoryRepository   // 背包仓库
	skillManager        *SkillManager
	buffManager         *BuffManager
	passiveSkillManager *PassiveSkillManager
	strategyExecutor    *StrategyExecutor
	battleStatsRepo     *repository.BattleStatsRepository // 战斗统计仓库

	// 新增系统集成
	calculator           *Calculator           // 数值计算系统
	monsterManager       *MonsterManager       // 怪物管理系统
	teamManager          *TeamManager          // 队伍管理系统
	zoneManager          *ZoneManager          // 地图管理系统
	equipmentManager     *EquipmentManager     // 装备管理系统
	battleStatsCollector *BattleStatsCollector // 战斗统计收集器

	// 用户自定义统计会话管理
	statsSessions   map[int]*StatsSession // key: userID, 用户自定义的统计会话
	statsSessionsMu sync.RWMutex          // 统计会话的锁
}

// StatsSession 用户自定义统计会话
type StatsSession struct {
	UserID    int
	StartTime time.Time
	IsActive  bool
}

// BattleSession 用户战斗会话
type BattleSession struct {
	UserID             int
	IsRunning          bool
	CurrentZone        *models.Zone
	CurrentEnemy       *models.Monster   // 保留用于向后兼容
	CurrentEnemies     []*models.Monster // 多个敌人支持
	BattleLogs         []models.BattleLog
	BattleCount        int
	SessionKills       int
	SessionGold        int
	SessionExp         int
	StartedAt          time.Time
	LastTick           time.Time
	IsResting          bool       // 是否在休息
	RestUntil          *time.Time // 休息结束时间
	RestStartedAt      *time.Time // 休息开始时间
	LastRestTick       *time.Time // 上次恢复处理的时间
	RestSpeed          float64    // 恢复速度倍率
	CurrentBattleExp   int        // 本场战斗获得的经验
	CurrentBattleGold  int        // 本场战斗获得的金币
	CurrentBattleKills int        // 本场战斗击杀数
	CurrentTurnIndex   int        // 回合控制：-1=玩家回合，>=0=敌人索引
	JustEncountered    bool       // 刚遭遇敌人，需要等待1个tick再开始战斗

	// 战斗统计收集
	BattleStartTime    time.Time                              // 本场战斗开始时间
	CurrentBattleRound int                                    // 本场战斗回合数
	CharacterStats     map[int]*CharacterBattleStatsCollector // 角色战斗统计收集器
	SkillBreakdown     map[int]map[string]*SkillUsageStats    // 角色->技能ID->技能使用统计

	// 威胁值系统
	ThreatTable map[string]map[int]int // 怪物ID -> 角色ID -> 威胁值

	// 速度排序回合系统
	TurnOrder             []*TurnParticipant // 回合顺序队列（按速度排序）
	CurrentTurnOrderIndex int                // 当前回合队列索引
}

// TurnParticipant 回合参与者
type TurnParticipant struct {
	Type      string            // "character" 或 "monster"
	Character *models.Character // 如果是角色
	Monster   *models.Monster   // 如果是怪物
	Speed     int               // 速度值
	Index     int               // 原始索引（用于角色或怪物数组）
}

// CharacterBattleStatsCollector 角色战斗统计收集器（内存中收集，战斗结束时保存）
type CharacterBattleStatsCollector struct {
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

// SkillUsageStats 技能使用统计
type SkillUsageStats struct {
	SkillID      string
	UseCount     int
	HitCount     int
	CritCount    int
	TotalDamage  int
	TotalHealing int
	ResourceCost int
}

// NewBattleManager 创建战斗管理器
func NewBattleManager() *BattleManager {
	return &BattleManager{
		sessions:             make(map[int]*BattleSession),
		gameRepo:             repository.NewGameRepository(),
		charRepo:             repository.NewCharacterRepository(),
		explorationRepo:      repository.NewExplorationRepository(),
		inventoryRepo:        repository.NewInventoryRepository(),
		skillManager:         NewSkillManager(),
		buffManager:          NewBuffManager(),
		passiveSkillManager:  NewPassiveSkillManager(),
		strategyExecutor:     NewStrategyExecutor(),
		battleStatsRepo:      repository.NewBattleStatsRepository(),
		calculator:           NewCalculator(),
		monsterManager:       NewMonsterManager(),
		teamManager:          NewTeamManager(),
		zoneManager:          NewZoneManager(),
		equipmentManager:     NewEquipmentManager(),
		battleStatsCollector: NewBattleStatsCollector(),
		statsSessions:        make(map[int]*StatsSession),
	}
}

// 全局战斗管理器实例
var battleManager *BattleManager
var once sync.Once

// GetBattleManager 获取战斗管理器单例
func GetBattleManager() *BattleManager {
	once.Do(func() {
		battleManager = NewBattleManager()
	})
	return battleManager
}

// GetOrCreateSession 获取或创建战斗会话
func (m *BattleManager) GetOrCreateSession(userID int) *BattleSession {
	m.mu.Lock()
	defer m.mu.Unlock()

	if session, exists := m.sessions[userID]; exists {
		// 如果会话存在但没有地图，不在这里设置默认地图
		// 让 GetBattleStatus 根据角色阵营来设置正确的默认地图
		return session
	}

	session := &BattleSession{
		UserID:                userID,
		BattleLogs:            make([]models.BattleLog, 0),
		StartedAt:             time.Now(),
		CurrentEnemies:        make([]*models.Monster, 0),
		CurrentTurnIndex:      -1, // 初始化为玩家回合
		RestSpeed:             1.0,
		CurrentZone:           nil,                          // 不在这里设置默认地图，让 GetBattleStatus 根据角色阵营设置
		ThreatTable:           make(map[string]map[int]int), // 初始化威胁值表
		CharacterStats:        make(map[int]*CharacterBattleStatsCollector),
		SkillBreakdown:        make(map[int]map[string]*SkillUsageStats),
		TurnOrder:             make([]*TurnParticipant, 0), // 初始化回合队列
		CurrentTurnOrderIndex: -1,                          // 初始化为-1，表示需要重新排序
	}

	m.sessions[userID] = session
	return session
}

// GetSession 获取战斗会话
func (m *BattleManager) GetSession(userID int) *BattleSession {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.sessions[userID]
}

// ToggleBattle 切换战斗状态
func (m *BattleManager) ToggleBattle(userID int) (bool, error) {
	session := m.GetOrCreateSession(userID)

	m.mu.Lock()
	defer m.mu.Unlock()

	session.IsRunning = !session.IsRunning
	session.LastTick = time.Now()

	if session.IsRunning {
		// 如果没有设置区域，设置默认区域
		if session.CurrentZone == nil {
			zone, err := m.gameRepo.GetZoneByID("elwynn")
			if err == nil {
				session.CurrentZone = zone
			}
		}
		session.CurrentTurnIndex = -1 // 重置为玩家回合
		m.addLog(session, "system", ">> 开始自动战斗...", "#33ff33")
	} else {
		m.addLog(session, "system", ">> 暂停自动战斗", "#ffff00")
	}

	return session.IsRunning, nil
}

// StartBattle 开始战斗
func (m *BattleManager) StartBattle(userID int) (bool, error) {
	session := m.GetOrCreateSession(userID)

	m.mu.Lock()
	defer m.mu.Unlock()

	if session.IsRunning {
		return true, nil
	}

	session.IsRunning = true
	session.LastTick = time.Now()
	session.CurrentTurnIndex = -1 // 重置为玩家回合

	// 设置默认区域
	if session.CurrentZone == nil {
		zone, err := m.gameRepo.GetZoneByID("elwynn")
		if err == nil {
			session.CurrentZone = zone
		}
	}

	// 战斗开始时，重置所有战士角色的怒气为0
	characters, err := m.charRepo.GetByUserID(userID)
	if err == nil {
		for _, char := range characters {
			if char != nil && char.ResourceType == "rage" {
				char.Resource = 0
				char.MaxResource = 100
				// 立即保存到数据库
				m.charRepo.UpdateAfterBattle(char.ID, char.HP, char.Resource, char.Exp, char.Level,
					char.ExpToNext, char.MaxHP, char.MaxResource, char.PhysicalAttack, char.MagicAttack, char.PhysicalDefense, char.MagicDefense,
					char.Strength, char.Agility, char.Intellect, char.Stamina, char.Spirit, char.UnspentPoints, char.TotalKills)
			}
		}
	}

	m.addLog(session, "system", ">> 开始自动战斗...", "#33ff33")
	return true, nil
}

// StopBattle 停止战斗
func (m *BattleManager) StopBattle(userID int) error {
	session := m.GetSession(userID)
	if session == nil {
		return nil
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	session.IsRunning = false
	m.addLog(session, "system", ">> 暂停自动战斗", "#ffff00")
	return nil
}

// ExecuteBattleTick 执行战斗回合（回合制：每tick只执行一个动作）
func (m *BattleManager) ExecuteBattleTick(userID int, characters []*models.Character) (*BattleTickResult, error) {
	session := m.GetOrCreateSession(userID)

	m.mu.Lock()
	defer m.mu.Unlock()

	// 如果没有角色，返回nil
	if len(characters) == 0 {
		return nil, nil
	}

	// 使用第一个角色进行战斗
	char := characters[0]

	// 确保战士的怒气上限为100（每次tick都检查，防止被覆盖）
	if char.ResourceType == "rage" {
		char.MaxResource = 100
	}

	// 加载角色的技能（如果还没有加载）
	if m.skillManager != nil {
		if err := m.skillManager.LoadCharacterSkills(char.ID); err != nil {
			// 如果加载失败，记录日志但不中断战斗
			m.addLog(session, "system", fmt.Sprintf("警告：无法加载角色技能: %v", err), "#ffaa00")
		}
	}

	// 加载角色的被动技能（如果还没有加载）
	if m.passiveSkillManager != nil {
		if err := m.passiveSkillManager.LoadCharacterPassiveSkills(char.ID); err != nil {
			// 如果加载失败，记录日志但不中断战斗
			m.addLog(session, "system", fmt.Sprintf("警告：无法加载角色被动技能: %v", err), "#ffaa00")
		}
	}

	// 如果战斗未运行且不在休息状态，检查是否需要返回角色数据
	// 如果角色刚复活（之前死亡但现在不死亡），需要返回一次数据让前端更新
	if !session.IsRunning && !session.IsResting {
		// 从数据库重新加载角色数据以确保状态正确
		updatedChar, err := m.charRepo.GetByID(char.ID)
		if err == nil && updatedChar != nil {
			char = updatedChar
			// 确保战士的怒气上限为100
			if char.ResourceType == "rage" {
				char.MaxResource = 100
			}
			// 如果角色已经复活（之前死亡但现在不死亡），返回角色数据
			if !char.IsDead {
				// 返回角色数据，让前端知道角色已经复活
				return &BattleTickResult{
					Character:    char,
					Enemy:        nil,
					Enemies:      nil,
					Logs:         []models.BattleLog{},
					IsRunning:    false,
					IsResting:    false,
					RestUntil:    nil,
					SessionKills: session.SessionKills,
					SessionGold:  session.SessionGold,
					SessionExp:   session.SessionExp,
					BattleCount:  session.BattleCount,
				}, nil
			}
		}
		// 如果无法获取角色数据或角色仍然死亡，返回nil
		return nil, nil
	}

	session.LastTick = time.Now()
	logs := make([]models.BattleLog, 0)

	// 检查角色是否死亡且还没到复活时间
	now := time.Now()
	if char.IsDead && char.ReviveAt != nil && now.Before(*char.ReviveAt) {
		// 角色死亡但还没到复活时间，进入休息状态
		if !session.IsResting {
			// 计算休息时间（复活时间 + 恢复时间）
			reviveRemaining := char.ReviveAt.Sub(now)
			recoveryTime := 25 * time.Second // 恢复一半HP需要的时间
			restDuration := reviveRemaining + recoveryTime
			restUntil := now.Add(restDuration)
			session.IsResting = true
			session.RestUntil = &restUntil
			session.RestStartedAt = &now
			session.LastRestTick = &now
			session.RestSpeed = 1.0
			// 保持 isRunning = true，这样按钮会显示"停止挂机"，休息状态可以自动处理

			remainingSeconds := int(reviveRemaining.Seconds()) + 1
			m.addLog(session, "death", fmt.Sprintf("%s 正在复活中... (剩余 %d 秒)", char.Name, remainingSeconds), "#ff0000")
			logs = append(logs, session.BattleLogs[len(session.BattleLogs)-1])
		}
	}

	// 如果正在休息，处理休息
	if session.IsResting && session.RestUntil != nil {
		initialHP := char.HP
		initialMP := char.Resource
		now := time.Now()
		m.processRest(session, char)

		// 更新LastTick，用于下次计算时间差
		session.LastTick = now

		if !session.IsResting {
			// 休息结束，保存角色数据
			m.charRepo.UpdateAfterBattle(char.ID, char.HP, char.Resource, char.Exp, char.Level,
				char.ExpToNext, char.MaxHP, char.MaxResource, char.PhysicalAttack, char.MagicAttack, char.PhysicalDefense, char.MagicDefense,
				char.Strength, char.Agility, char.Intellect, char.Stamina, char.Spirit, char.UnspentPoints, char.TotalKills)

			// 休息结束后，确保返回角色数据，让前端知道休息已结束
			// 从数据库重新加载角色数据以确保状态正确
			updatedChar, err := m.charRepo.GetByID(char.ID)
			if err == nil && updatedChar != nil {
				char = updatedChar
				// 确保战士的怒气上限为100
				if char.ResourceType == "rage" {
					char.MaxResource = 100
				}
			}

			// 如果角色已经复活（不再死亡），自动恢复战斗
			if !char.IsDead {
				session.IsRunning = true
				m.addLog(session, "system", ">> 休息结束，自动恢复战斗", "#33ff33")
			} else {
				m.addLog(session, "system", ">> 休息结束，准备下一场战斗", "#33ff33")
			}
			logs = append(logs, session.BattleLogs[len(session.BattleLogs)-1])
		} else {
			// 仍在休息中
			remaining := time.Until(*session.RestUntil)
			if remaining > 0 {
				m.addLog(session, "system", fmt.Sprintf(">> 休息中... (剩余 %d 秒)", int(remaining.Seconds())+1), "#888888")
				logs = append(logs, session.BattleLogs[len(session.BattleLogs)-1])
			}
		}

		// 保存角色数据更新
		if char.HP != initialHP || char.Resource != initialMP {
			m.charRepo.UpdateAfterBattle(char.ID, char.HP, char.Resource, char.Exp, char.Level,
				char.ExpToNext, char.MaxHP, char.MaxResource, char.PhysicalAttack, char.MagicAttack, char.PhysicalDefense, char.MagicDefense,
				char.Strength, char.Agility, char.Intellect, char.Stamina, char.Spirit, char.UnspentPoints, char.TotalKills)
		}

		return &BattleTickResult{
			Character:    char,
			Enemy:        nil,
			Enemies:      session.CurrentEnemies,
			Logs:         logs,
			IsRunning:    session.IsRunning,
			IsResting:    session.IsResting,
			RestUntil:    session.RestUntil,
			SessionKills: session.SessionKills,
			SessionGold:  session.SessionGold,
			SessionExp:   session.SessionExp,
			BattleCount:  session.BattleCount,
		}, nil
	}

	// 获取存活的敌人
	aliveEnemies := make([]*models.Monster, 0)
	if session.CurrentEnemies != nil {
		for _, enemy := range session.CurrentEnemies {
			if enemy != nil && enemy.HP > 0 {
				aliveEnemies = append(aliveEnemies, enemy)
			}
		}
	}

	// 如果没有敌人，生成新的
	if len(aliveEnemies) == 0 {
		// 重置本场战斗统计
		session.CurrentBattleExp = 0
		session.CurrentBattleGold = 0
		session.CurrentBattleKills = 0
		session.CurrentTurnIndex = -1 // 玩家回合

		// 初始化战斗统计收集器
		m.initBattleStats(session, characters)

		// 战斗开始时，确保战士的怒气为0，最大怒气为100
		if char.ResourceType == "rage" {
			char.Resource = 0
			char.MaxResource = 100
		}

		err := m.spawnEnemies(session, char.Level, len(characters))
		if err != nil {
			// 如果生成敌人失败，记录错误并返回
			m.addLog(session, "error", fmt.Sprintf("生成敌人失败: %v", err), "#ff0000")
			return nil, fmt.Errorf("failed to spawn enemies: %v", err)
		}
		logs = append(logs, session.BattleLogs[len(session.BattleLogs)-1])

		// 构建回合顺序队列（按速度排序）
		m.buildTurnOrder(session, characters, session.CurrentEnemies)

		// 初始化战斗回合数和开始时间
		session.CurrentBattleRound = 1
		session.BattleStartTime = time.Now()

		// 添加战斗开始日志
		enemyNames := make([]string, 0, len(session.CurrentEnemies))
		for _, enemy := range session.CurrentEnemies {
			if enemy != nil {
				enemyNames = append(enemyNames, enemy.Name)
			}
		}
		if len(enemyNames) > 0 {
			enemyList := strings.Join(enemyNames, "、")
			m.addLog(session, "system", fmt.Sprintf("━━━ 遭遇敌人：%s ━━━", enemyList), "#ffaa00")
			logs = append(logs, session.BattleLogs[len(session.BattleLogs)-1])
		}

		// 标记刚遭遇敌人，需要等待1个tick再开始战斗
		session.JustEncountered = true

		// 刚遭遇敌人，这个tick只显示信息，不执行战斗
		return &BattleTickResult{
			Character:    char,
			Enemy:        session.CurrentEnemy,
			Enemies:      session.CurrentEnemies,
			Logs:         logs,
			IsRunning:    session.IsRunning,
			IsResting:    session.IsResting,
			RestUntil:    session.RestUntil,
			SessionKills: session.SessionKills,
			SessionGold:  session.SessionGold,
			SessionExp:   session.SessionExp,
			BattleCount:  session.BattleCount,
		}, nil
	}

	// 如果刚遭遇敌人，这个tick只显示信息，不执行战斗
	if session.JustEncountered {
		session.JustEncountered = false // 清除标志，下一个tick开始战斗
		return &BattleTickResult{
			Character:    char,
			Enemy:        session.CurrentEnemy,
			Enemies:      session.CurrentEnemies,
			Logs:         logs,
			IsRunning:    session.IsRunning,
			IsResting:    session.IsResting,
			RestUntil:    session.RestUntil,
			SessionKills: session.SessionKills,
			SessionGold:  session.SessionGold,
			SessionExp:   session.SessionExp,
			BattleCount:  session.BattleCount,
		}, nil
	}

	// 确保TurnOrder已初始化
	if session.TurnOrder == nil || len(session.TurnOrder) == 0 || session.CurrentTurnOrderIndex < 0 {
		m.buildTurnOrder(session, characters, session.CurrentEnemies)
	}

	// 使用速度排序的回合系统
	currentParticipant := m.getCurrentTurnParticipant(session)
	if currentParticipant == nil {
		// 如果没有参与者，重新构建队列
		m.buildTurnOrder(session, characters, session.CurrentEnemies)
		currentParticipant = m.getCurrentTurnParticipant(session)
		if currentParticipant == nil {
			// 仍然没有参与者，可能是所有角色和敌人都死亡了
			return &BattleTickResult{
				Character:    char,
				Enemy:        session.CurrentEnemy,
				Enemies:      session.CurrentEnemies,
				Logs:         logs,
				IsRunning:    session.IsRunning,
				IsResting:    session.IsResting,
				RestUntil:    session.RestUntil,
				SessionKills: session.SessionKills,
				SessionGold:  session.SessionGold,
				SessionExp:   session.SessionExp,
				BattleCount:  session.BattleCount,
			}, nil
		}
	}

	// 根据参与者类型设置CurrentTurnIndex以保持向后兼容
	// 然后使用原有的回合逻辑执行行动
	if currentParticipant.Type == "character" {
		// 角色回合：设置CurrentTurnIndex为-1以保持兼容
		actingChar := currentParticipant.Character
		if actingChar == nil || actingChar.HP <= 0 {
			// 角色已死亡，跳过
			m.moveToNextTurn(session, characters, aliveEnemies)
			return &BattleTickResult{
				Character:    char,
				Enemy:        session.CurrentEnemy,
				Enemies:      session.CurrentEnemies,
				Logs:         logs,
				IsRunning:    session.IsRunning,
				IsResting:    session.IsResting,
				RestUntil:    session.RestUntil,
				SessionKills: session.SessionKills,
				SessionGold:  session.SessionGold,
				SessionExp:   session.SessionExp,
				BattleCount:  session.BattleCount,
			}, nil
		}
		// 只处理主要角色（当前实现只支持单角色）
		if actingChar.ID != char.ID {
			// 其他角色，暂时跳过（多角色系统后续实现）
			m.moveToNextTurn(session, characters, aliveEnemies)
			return &BattleTickResult{
				Character:    char,
				Enemy:        session.CurrentEnemy,
				Enemies:      session.CurrentEnemies,
				Logs:         logs,
				IsRunning:    session.IsRunning,
				IsResting:    session.IsResting,
				RestUntil:    session.RestUntil,
				SessionKills: session.SessionKills,
				SessionGold:  session.SessionGold,
				SessionExp:   session.SessionExp,
				BattleCount:  session.BattleCount,
			}, nil
		}
		session.CurrentTurnIndex = -1
	} else {
		// 怪物回合：找到怪物在aliveEnemies中的索引
		actingEnemy := currentParticipant.Monster
		if actingEnemy == nil || actingEnemy.HP <= 0 {
			// 怪物已死亡，跳过
			m.moveToNextTurn(session, characters, aliveEnemies)
			return &BattleTickResult{
				Character:    char,
				Enemy:        session.CurrentEnemy,
				Enemies:      session.CurrentEnemies,
				Logs:         logs,
				IsRunning:    session.IsRunning,
				IsResting:    session.IsResting,
				RestUntil:    session.RestUntil,
				SessionKills: session.SessionKills,
				SessionGold:  session.SessionGold,
				SessionExp:   session.SessionExp,
				BattleCount:  session.BattleCount,
			}, nil
		}
		// 找到怪物在aliveEnemies中的索引
		enemyIndex := -1
		for i, enemy := range aliveEnemies {
			if enemy != nil && enemy.ID == actingEnemy.ID {
				enemyIndex = i
				break
			}
		}
		if enemyIndex >= 0 {
			session.CurrentTurnIndex = enemyIndex
		} else {
			// 找不到怪物，跳过这个回合
			m.moveToNextTurn(session, characters, aliveEnemies)
			return &BattleTickResult{
				Character:    char,
				Enemy:        session.CurrentEnemy,
				Enemies:      session.CurrentEnemies,
				Logs:         logs,
				IsRunning:    session.IsRunning,
				IsResting:    session.IsResting,
				RestUntil:    session.RestUntil,
				SessionKills: session.SessionKills,
				SessionGold:  session.SessionGold,
				SessionExp:   session.SessionExp,
				BattleCount:  session.BattleCount,
			}, nil
		}
	}

	// 原有的回合制逻辑：CurrentTurnIndex == -1 表示玩家回合，>=0 表示敌人索引
	// 现在这个逻辑会根据TurnOrder系统设置的CurrentTurnIndex来执行
	if session.CurrentTurnIndex == -1 {
		// 玩家回合：攻击第一个存活的敌人
		if len(aliveEnemies) > 0 {
			target := aliveEnemies[0]
			targetHPPercent := float64(target.HP) / float64(target.MaxHP)
			hasMultipleEnemies := len(aliveEnemies) > 1
			targetIndex := 0

			// 使用技能管理器选择技能
			var skillState *CharacterSkillState
			var strategyDecision *SkillDecision
			var strategy *models.BattleStrategy
			var battleCtx *BattleContext

			// 优先使用策略执行器
			hasStrategy := false
			if m.strategyExecutor != nil {
				strategy = m.strategyExecutor.GetActiveStrategy(char.ID)
				if strategy != nil {
					hasStrategy = true
					// 构建战斗上下文
					battleCtx = &BattleContext{
						Character:    char,
						Enemies:      aliveEnemies,
						Allies:       characters,
						Target:       target,
						CurrentRound: session.BattleCount,
						SkillManager: m.skillManager,
						BuffManager:  m.buffManager,
					}
					strategyDecision = m.strategyExecutor.ExecuteStrategy(strategy, battleCtx)
				}
			}

			// 根据策略决策或默认逻辑选择技能
			if strategyDecision != nil {
				// 更新目标（无论是普通攻击还是技能，都应该使用策略选择的目标）
				if strategyDecision.TargetIndex >= 0 && strategyDecision.TargetIndex < len(aliveEnemies) {
					targetIndex = strategyDecision.TargetIndex
					target = aliveEnemies[targetIndex]
					targetHPPercent = float64(target.HP) / float64(target.MaxHP)
				}

				if strategyDecision.IsNormalAttack {
					// 策略决定使用普通攻击
					skillState = nil
				} else if strategyDecision.SkillID != "" {
					// 策略决定使用特定技能
					skillState = m.skillManager.GetSkillState(char.ID, strategyDecision.SkillID)
					if skillState == nil {
						// 尝试带 warrior_ 前缀
						skillState = m.skillManager.GetSkillState(char.ID, "warrior_"+strategyDecision.SkillID)
					}
				}
			} else if hasStrategy {
				// 有策略但返回 nil，表示没有可用技能或应该使用普通攻击
				// 不使用 SelectBestSkill，因为它不检查条件规则限制
				skillState = nil
				// 即使策略返回nil，也应该根据策略的目标优先级选择目标
				if strategy != nil {
					targetIndex = m.strategyExecutor.SelectTargetByStrategy(strategy, battleCtx, "")
					if targetIndex >= 0 && targetIndex < len(aliveEnemies) {
						target = aliveEnemies[targetIndex]
						targetHPPercent = float64(target.HP) / float64(target.MaxHP)
					}
				}
			} else if m.skillManager != nil {
				// 没有策略，使用默认逻辑
				skillState = m.skillManager.SelectBestSkill(char.ID, char.Resource, targetHPPercent, hasMultipleEnemies, m.buffManager)
			}
			_ = targetIndex // 避免未使用警告

			var skillName string
			var playerDamage int
			var resourceCost int
			var usedSkill bool
			var skillEffects map[string]interface{}
			var isCrit bool
			var damageDetails *DamageCalculationDetails
			var shouldDealDamage bool // 是否应该造成伤害（只有attack类型的技能才造成伤害）
			var isDodged bool         // 是否被闪避
			var ignoresDodge bool     // 技能是否无视闪避
			var originalResource int  // 资源变化前的值（用于日志显示）

			// 保存资源变化前的值
			originalResource = char.Resource

			if skillState != nil && skillState.Skill != nil {
				// 使用技能
				skillName = skillState.Skill.Name
				resourceCost = m.skillManager.GetSkillResourceCost(skillState)

				// 判断技能是否应该造成伤害（只有attack类型的技能才造成伤害）
				shouldDealDamage = skillState.Skill.Type == "attack"

				// 检查技能是否无视闪避
				ignoresDodge = m.skillIgnoresDodge(skillState.Skill)

				// 检查资源是否足够
				if resourceCost <= char.Resource {

					var baseDamage int
					// playerDamage, isCrit, and damageDetails are already declared in outer scope
					// Do not redeclare them here to avoid shadowing outer scope variables

					if shouldDealDamage {
						// 计算技能伤害（基础伤害，暴击在后面处理）
						baseDamage = m.skillManager.CalculateSkillDamage(skillState, char, target, m.passiveSkillManager, m.buffManager)

						// 计算实际攻击力（用于公式显示，需要包含Buff加成）
						skillRatio := skillState.Skill.ScalingRatio
						actualAttackForFormula := float64(char.PhysicalAttack)
						attackModifiers := []string{}

						// 检查被动技能的攻击力加成
						if m.passiveSkillManager != nil {
							attackModifier := m.passiveSkillManager.GetPassiveModifier(char.ID, "attack")
							if attackModifier > 0 {
								actualAttackForFormula = actualAttackForFormula * (1.0 + attackModifier/100.0)
								attackModifiers = append(attackModifiers, fmt.Sprintf("被动攻击+%.0f%%", attackModifier))
							}
						}

						// 检查Buff的攻击力加成（战斗怒吼等）
						if m.buffManager != nil {
							attackBuffValue := m.buffManager.GetBuffValue(char.ID, "attack")
							if attackBuffValue > 0 {
								actualAttackForFormula = actualAttackForFormula * (1.0 + attackBuffValue/100.0)
								attackModifiers = append(attackModifiers, fmt.Sprintf("Buff攻击+%.0f%%", attackBuffValue))
							}
						}

						scaledDamage := actualAttackForFormula * skillRatio

						// 创建技能伤害详情
						damageDetails = &DamageCalculationDetails{
							BaseAttack:       char.PhysicalAttack,
							ActualAttack:     actualAttackForFormula,
							BaseDefense:      target.PhysicalDefense,
							BaseDamage:       float64(baseDamage),
							AttackModifiers:  attackModifiers,
							DefenseModifiers: []string{},
							ActualCritRate:   -1, // -1 表示未设置
							RandomRoll:       -1, // -1 表示未设置
							SkillRatio:       skillRatio,
							ScaledDamage:     scaledDamage,
						}

						// 计算暴击（技能也可以暴击，应用被动技能和Buff加成）
						// 根据伤害类型选择使用物理暴击率还是法术暴击率
						var baseCritRate, baseCritDamage float64
						var critType string
						if skillState.Skill.DamageType == "physical" {
							baseCritRate = char.PhysCritRate
							baseCritDamage = char.PhysCritDamage
							critType = "phys_crit_rate"
						} else {
							// 法术伤害（magic/fire/frost/shadow/holy/nature）
							baseCritRate = char.SpellCritRate
							baseCritDamage = char.SpellCritDamage
							critType = "spell_crit_rate"
						}

						actualCritRate := baseCritRate
						damageDetails.BaseCritRate = baseCritRate
						damageDetails.CritModifiers = []string{}

						if m.passiveSkillManager != nil {
							// 检查特定类型暴击率加成
							critModifier := m.passiveSkillManager.GetPassiveModifier(char.ID, critType)
							if critModifier > 0 {
								actualCritRate = baseCritRate + critModifier/100.0
								damageDetails.CritModifiers = append(damageDetails.CritModifiers,
									fmt.Sprintf("被动暴击+%.0f%%", critModifier))
							}
							// 检查通用暴击率加成（同时影响物理和法术）
							generalCritModifier := m.passiveSkillManager.GetPassiveModifier(char.ID, "crit_rate")
							if generalCritModifier > 0 {
								actualCritRate = actualCritRate + generalCritModifier/100.0
								damageDetails.CritModifiers = append(damageDetails.CritModifiers,
									fmt.Sprintf("被动暴击+%.0f%%", generalCritModifier))
							}
						}
						// 应用Buff的暴击率加成（鲁莽等）
						if m.buffManager != nil {
							// 检查特定类型暴击率加成
							critBuffValue := m.buffManager.GetBuffValue(char.ID, critType)
							if critBuffValue > 0 {
								actualCritRate = actualCritRate + critBuffValue/100.0
								damageDetails.CritModifiers = append(damageDetails.CritModifiers,
									fmt.Sprintf("Buff暴击+%.0f%%", critBuffValue))
							}
							// 检查通用暴击率加成（同时影响物理和法术）
							generalCritBuffValue := m.buffManager.GetBuffValue(char.ID, "crit_rate")
							if generalCritBuffValue > 0 {
								actualCritRate = actualCritRate + generalCritBuffValue/100.0
								damageDetails.CritModifiers = append(damageDetails.CritModifiers,
									fmt.Sprintf("Buff暴击+%.0f%%", generalCritBuffValue))
							}
						}
						if actualCritRate > 1.0 {
							actualCritRate = 1.0
						}
						damageDetails.ActualCritRate = actualCritRate
						randomRoll := rand.Float64()
						damageDetails.RandomRoll = randomRoll
						isCrit = randomRoll < actualCritRate
						damageDetails.IsCrit = isCrit
						damageDetails.CritMultiplier = baseCritDamage

						if isCrit {
							playerDamage = int(float64(baseDamage) * baseCritDamage)
						} else {
							playerDamage = baseDamage
						}
						damageDetails.FinalDamage = playerDamage
					}

					// 应用技能效果
					skillEffects = m.skillManager.ApplySkillEffects(skillState, char, target)

					// 应用Buff/Debuff效果
					m.applySkillBuffs(skillState, char, target, skillEffects)

					// 应用Debuff到敌人（挫志怒吼、旋风斩等）
					m.applySkillDebuffs(skillState, char, target, aliveEnemies, skillEffects)

					// 保存资源变化前的值
					originalResource := char.Resource

					// 消耗资源
					char.Resource -= resourceCost
					if char.Resource < 0 {
						char.Resource = 0
					}

					// 使用技能（设置冷却）
					m.skillManager.UseSkill(char.ID, skillState.SkillID)
					usedSkill = true

					// 处理被动技能的使用技能时效果
					m.handlePassiveOnSkillUseEffects(char, skillState.SkillID, session, &logs)

					// 处理技能特殊效果（怒气获得等）
					if rageGain, ok := skillEffects["rageGain"].(int); ok {
						// 应用被动技能的怒气获得加成（愤怒掌握等）
						actualRageGain := m.applyRageGenerationModifiers(char.ID, rageGain)
						char.Resource += actualRageGain
						if char.Resource > char.MaxResource {
							char.Resource = char.MaxResource
						}
					}

					// 只有attack类型的技能才造成伤害
					if shouldDealDamage {
						// 【闪避判定】检查主目标是否闪避（非AOE技能）
						if skillState.Skill.TargetType != "enemy_all" {
							if m.checkDodge(target.DodgeRate, ignoresDodge) {
								isDodged = true
							}
						}

						// 处理AOE技能（旋风斩等）
						if skillState.Skill.TargetType == "enemy_all" {
							// 根据技能伤害类型获取暴击伤害
							var aoeCritDamage float64
							if skillState.Skill.DamageType == "physical" {
								aoeCritDamage = char.PhysCritDamage
							} else {
								aoeCritDamage = char.SpellCritDamage
							}
							// 对所有敌人造成伤害（AOE技能每个敌人单独判定闪避）
							for _, enemy := range aliveEnemies {
								if enemy.HP > 0 {
									// AOE 技能每个敌人单独判定闪避
									if m.checkDodge(enemy.DodgeRate, ignoresDodge) {
										m.addLog(session, "dodge", fmt.Sprintf("%s 闪避了 %s 的攻击！", enemy.Name, char.Name), "#00ffff")
										logs = append(logs, session.BattleLogs[len(session.BattleLogs)-1])
										continue
									}
									damage := m.skillManager.CalculateSkillDamage(skillState, char, enemy, m.passiveSkillManager, m.buffManager)
									if isCrit {
										// 根据技能伤害类型选择暴击伤害
										damage = int(float64(damage) * aoeCritDamage)
									}
									enemy.HP -= damage
									if enemy.HP < 0 {
										enemy.HP = 0
									}
									// 更新威胁值（AOE技能对每个目标都产生威胁）
									m.updateThreat(session, enemy.ID, char.ID, damage)
								}
							}
							// playerDamage用于日志显示（主目标伤害）
						} else if skillState.SkillID == "warrior_cleave" {
							// 顺劈斩：主目标+相邻目标
							// 主目标闪避检查已在上方完成，如果未闪避则造成伤害
							if !isDodged {
								target.HP -= playerDamage
							}

							// 对相邻目标造成伤害（最多2个）
							// 收集相邻目标的日志信息，稍后记录（在主目标日志之后）
							adjacentLogs := make([]models.BattleLog, 0)
							adjacentTotalDamage := 0 // 累计波及伤害总和，用于统计
							adjacentCount := 0
							processedEnemies := make(map[*models.Monster]bool) // 记录已处理的敌人，避免重复
							for _, enemy := range aliveEnemies {
								// 确保不是主目标，且未处理过，且还有空位
								if enemy != target && enemy.HP > 0 && adjacentCount < 2 && !processedEnemies[enemy] {
									processedEnemies[enemy] = true // 标记为已处理
									// 相邻目标单独判定闪避
									if m.checkDodge(enemy.DodgeRate, ignoresDodge) {
										// 先创建日志但不立即添加到session，稍后统一添加
										adjacentLog := models.BattleLog{
											LogType: "dodge",
											Message: fmt.Sprintf("%s 闪避了 %s 的攻击！", enemy.Name, char.Name),
											Color:   "#00ffff",
										}
										adjacentLogs = append(adjacentLogs, adjacentLog)
										adjacentCount++
										continue
									}
									// 计算相邻目标伤害
									if effect, ok := skillState.Effect["adjacentMultiplier"].(float64); ok {
										adjacentDamage := int(float64(char.PhysicalAttack) * effect)
										// 基础伤害 = 实际攻击力 - 目标防御力（不再除以2）
										adjacentDamage = adjacentDamage - enemy.PhysicalDefense
										if adjacentDamage < 1 {
											adjacentDamage = 1
										}
										if isCrit {
											// 顺劈斩是物理技能，使用物理暴击伤害
											adjacentDamage = int(float64(adjacentDamage) * char.PhysCritDamage)
										}
										adjacentOldHP := enemy.HP
										enemy.HP -= adjacentDamage
										if enemy.HP < 0 {
											enemy.HP = 0
										}
										// 更新威胁值（顺劈斩对相邻目标也产生威胁）
										m.updateThreat(session, enemy.ID, char.ID, adjacentDamage)
										adjacentCount++
										adjacentTotalDamage += adjacentDamage // 累计伤害用于统计
										adjacentHPChange := m.formatHPChange(enemy.Name, adjacentOldHP, enemy.HP, enemy.MaxHP)
										// 先创建日志但不立即添加到session，稍后统一添加
										adjacentLog := models.BattleLog{
											LogType:    "combat",
											Message:    fmt.Sprintf("%s 的顺劈斩波及到 %s，造成 %d 点伤害%s", char.Name, enemy.Name, adjacentDamage, adjacentHPChange),
											Color:      "#ffaa00",
											DamageType: "physical",
										}
										adjacentLogs = append(adjacentLogs, adjacentLog)
									}
								}
							}
							// 将相邻目标日志信息和总伤害存储到skillState中，稍后记录
							if skillState.Effect == nil {
								skillState.Effect = make(map[string]interface{})
							}
							skillState.Effect["_adjacentLogs"] = adjacentLogs
							skillState.Effect["_adjacentTotalDamage"] = adjacentTotalDamage
						} else {
							// 单体技能 - 如果未闪避则造成伤害
							if !isDodged {
								target.HP -= playerDamage
								// 更新威胁值（威胁值等于伤害值）
								m.updateThreat(session, target.ID, char.ID, playerDamage)
							}
						}
					} else {
						// buff技能使用后，还需要进行普通攻击
						// 先记录buff技能使用日志
						buffResourceChangeText := m.formatResourceChange(char.ResourceType, originalResource, char.Resource)
						m.addLog(session, "combat", fmt.Sprintf("%s 使用 [%s]%s", char.Name, skillName, buffResourceChangeText), "#8888ff")
						logs = append(logs, session.BattleLogs[len(session.BattleLogs)-1])
						// 重置资源消耗，避免普通攻击日志重复显示
						resourceCost = 0
						// 设置skillState为nil，让后续代码进行普通攻击
						skillState = nil
					}
				} else {
					// 资源不足，使用普通攻击
					skillState = nil
				}
			}

			// 如果没有使用技能或资源不足，或使用了buff技能，使用普通攻击
			if skillState == nil {
				skillName = "普通攻击"
				shouldDealDamage = true // 普通攻击造成伤害
				ignoresDodge = false    // 普通攻击不无视闪避

				// 【闪避判定】检查目标是否闪避普通攻击
				if m.checkDodge(target.DodgeRate, ignoresDodge) {
					isDodged = true
				}
				// 计算实际物理攻击力（应用被动技能加成）
				actualAttack := float64(char.PhysicalAttack)
				damageDetails = &DamageCalculationDetails{
					BaseAttack:       char.PhysicalAttack,
					BaseDefense:      target.PhysicalDefense,
					AttackModifiers:  []string{},
					DefenseModifiers: []string{},
					ActualCritRate:   -1, // -1 表示未设置
					RandomRoll:       -1, // -1 表示未设置
				}

				if m.passiveSkillManager != nil {
					attackModifier := m.passiveSkillManager.GetPassiveModifier(char.ID, "attack")
					if attackModifier > 0 {
						actualAttack = actualAttack * (1.0 + attackModifier/100.0)
						damageDetails.AttackModifiers = append(damageDetails.AttackModifiers,
							fmt.Sprintf("被动攻击+%.0f%%", attackModifier))
					}
					// 应用被动技能的伤害加成
					damageModifier := m.passiveSkillManager.GetPassiveModifier(char.ID, "damage")
					if damageModifier > 0 {
						actualAttack = actualAttack * (1.0 + damageModifier/100.0)
						damageDetails.AttackModifiers = append(damageDetails.AttackModifiers,
							fmt.Sprintf("被动伤害+%.0f%%", damageModifier))
					}

					// 处理低血量时的攻击力加成（狂暴之心）
					hpPercent := float64(char.HP) / float64(char.MaxHP)
					passives := m.passiveSkillManager.GetPassiveSkills(char.ID)
					for _, passive := range passives {
						if passive.Passive.EffectType == "stat_mod" && passive.Passive.ID == "warrior_passive_berserker_heart" {
							// 根据等级计算触发阈值（1级50%，5级30%）
							threshold := 0.50 - float64(passive.Level-1)*0.05
							if hpPercent < threshold {
								// 根据等级计算攻击力加成（1级20%，5级60%）
								attackBonus := 20.0 + float64(passive.Level-1)*10.0
								actualAttack = actualAttack * (1.0 + attackBonus/100.0)
								damageDetails.AttackModifiers = append(damageDetails.AttackModifiers,
									fmt.Sprintf("狂暴之心+%.0f%%", attackBonus))
							}
						}
					}
				}
				// 应用Buff的攻击力加成（战斗怒吼、狂暴之怒、天神下凡等）
				if m.buffManager != nil {
					attackBuffValue := m.buffManager.GetBuffValue(char.ID, "attack")
					if attackBuffValue > 0 {
						actualAttack = actualAttack * (1.0 + attackBuffValue/100.0)
						damageDetails.AttackModifiers = append(damageDetails.AttackModifiers,
							fmt.Sprintf("Buff攻击+%.0f%%", attackBuffValue))
					}
				}

				damageDetails.ActualAttack = actualAttack
				damageDetails.ActualDefense = float64(target.PhysicalDefense)

				// 计算实际用于伤害计算的攻击力（四舍五入）
				attackUsedInCalc := int(math.Round(actualAttack))
				baseDamage, calcDetails := m.calculatePhysicalDamageWithDetails(attackUsedInCalc, target.PhysicalDefense)
				damageDetails.BaseDamage = calcDetails.BaseDamage
				damageDetails.BaseAttack = attackUsedInCalc // 确保公式显示的是实际用于计算的值
				damageDetails.Variance = calcDetails.Variance

				// 计算暴击率（普通攻击使用物理暴击率，应用被动技能和Buff加成）
				actualCritRate := char.PhysCritRate
				damageDetails.BaseCritRate = char.PhysCritRate
				damageDetails.CritModifiers = []string{}

				if m.passiveSkillManager != nil {
					// 检查物理暴击率加成
					critModifier := m.passiveSkillManager.GetPassiveModifier(char.ID, "phys_crit_rate")
					if critModifier > 0 {
						actualCritRate = char.PhysCritRate + critModifier/100.0
						damageDetails.CritModifiers = append(damageDetails.CritModifiers,
							fmt.Sprintf("被动暴击+%.0f%%", critModifier))
					}
					// 检查通用暴击率加成（同时影响物理和法术）
					generalCritModifier := m.passiveSkillManager.GetPassiveModifier(char.ID, "crit_rate")
					if generalCritModifier > 0 {
						actualCritRate = actualCritRate + generalCritModifier/100.0
						damageDetails.CritModifiers = append(damageDetails.CritModifiers,
							fmt.Sprintf("被动暴击+%.0f%%", generalCritModifier))
					}
				}
				// 应用Buff的暴击率加成（鲁莽等）
				if m.buffManager != nil {
					// 检查物理暴击率加成
					critBuffValue := m.buffManager.GetBuffValue(char.ID, "phys_crit_rate")
					if critBuffValue > 0 {
						actualCritRate = actualCritRate + critBuffValue/100.0
						damageDetails.CritModifiers = append(damageDetails.CritModifiers,
							fmt.Sprintf("Buff暴击+%.0f%%", critBuffValue))
					}
					// 检查通用暴击率加成（同时影响物理和法术）
					generalCritBuffValue := m.buffManager.GetBuffValue(char.ID, "crit_rate")
					if generalCritBuffValue > 0 {
						actualCritRate = actualCritRate + generalCritBuffValue/100.0
						damageDetails.CritModifiers = append(damageDetails.CritModifiers,
							fmt.Sprintf("Buff暴击+%.0f%%", generalCritBuffValue))
					}
				}
				if actualCritRate > 1.0 {
					actualCritRate = 1.0
				}
				damageDetails.ActualCritRate = actualCritRate
				// 使用 Calculator 进行暴击判定（内部会处理上限）
				isCrit = m.calculator.ShouldCrit(actualCritRate)
				damageDetails.IsCrit = isCrit
				damageDetails.CritMultiplier = char.PhysCritDamage
				damageDetails.RandomRoll = 0 // Calculator内部处理随机数

				if isCrit {
					playerDamage = int(float64(baseDamage) * char.PhysCritDamage)
				} else {
					playerDamage = baseDamage
				}
				damageDetails.FinalDamage = playerDamage

				// 如果未闪避，造成伤害
				if !isDodged {
					target.HP -= playerDamage
					// 更新威胁值（威胁值等于伤害值）
					m.updateThreat(session, target.ID, char.ID, playerDamage)
					// 记录伤害统计
					if m.battleStatsCollector != nil {
						m.battleStatsCollector.RecordDamage(char.ID, playerDamage, "physical", isCrit)
					}
				}
				// 注意：闪避统计只记录角色的闪避，怪物的闪避不记录到角色统计中
				resourceCost = 0
				usedSkill = false
			}
			// 如果使用了技能，isCrit已经在上面计算了

			// 普通攻击获得怒气（只有普通攻击才获得怒气，使用技能时不获得，闪避时不获得）
			if char.ResourceType == "rage" && !usedSkill && !isDodged {
				var baseRageGain int
				if isCrit {
					baseRageGain = 10 // 暴击获得10点怒气
				} else {
					baseRageGain = 5 // 普通攻击获得5点怒气
				}

				// 应用被动技能的怒气获得加成（愤怒掌握等）
				rageGain := m.applyRageGenerationModifiers(char.ID, baseRageGain)

				char.Resource += rageGain
				// 确保不超过最大值
				if char.Resource > char.MaxResource {
					char.Resource = char.MaxResource
				}
			}

			// 处理被动技能的特殊效果（攻击时触发）- 闪避时不触发
			if !isDodged {
				m.handlePassiveOnHitEffects(char, playerDamage, usedSkill, session, &logs)

				// 处理被动技能的暴击时效果（如果暴击）
				if isCrit {
					m.handlePassiveOnCritEffects(char, playerDamage, usedSkill, session, &logs)
				}
			}

			// 构建战斗日志消息，包含资源变化（带颜色）
			resourceChangeText := m.formatResourceChange(char.ResourceType, originalResource, char.Resource)

			// 格式化伤害公式
			formulaText := ""
			if damageDetails != nil {
				formulaText = m.formatDamageFormula(damageDetails)
			}

			// 记录技能使用日志
			if shouldDealDamage {
				if isDodged {
					// 被闪避时显示闪避日志
					m.addLog(session, "dodge", fmt.Sprintf("%s 闪避了 %s 使用的 [%s]！%s", target.Name, char.Name, skillName, resourceChangeText), "#00ffff")
				} else {
					// 计算目标HP变化（需要在造成伤害前记录原始HP）
					// 注意：此时伤害已经造成，target.HP已经是伤害后的值
					// 所以我们需要在造成伤害前记录原始HP，这里使用伤害值反推
					targetOldHP := target.HP + playerDamage
					if targetOldHP > target.MaxHP {
						targetOldHP = target.MaxHP
					}
					hpChangeText := m.formatHPChange(target.Name, targetOldHP, target.HP, target.MaxHP)
					playerDamageType := "physical"
					if skillState != nil && skillState.Skill != nil {
						if dt := normalizeDamageType(skillState.Skill.DamageType); dt != "" {
							playerDamageType = dt
						}
					}

					// 攻击类技能：记录伤害
					if isCrit {
						m.addLog(session, "combat", fmt.Sprintf("%s 使用 [%s] 💥暴击！对 %s 造成 %d 点伤害%s%s%s", char.Name, skillName, target.Name, playerDamage, formulaText, hpChangeText, resourceChangeText), "#ff6b6b", withDamageType(playerDamageType))
					} else {
						m.addLog(session, "combat", fmt.Sprintf("%s 使用 [%s] 对 %s 造成 %d 点伤害%s%s%s", char.Name, skillName, target.Name, playerDamage, formulaText, hpChangeText, resourceChangeText), "#ffaa00", withDamageType(playerDamageType))
					}

					// 如果是顺劈斩，在主目标日志后记录相邻目标的日志
					if skillState != nil && skillState.SkillID == "warrior_cleave" {
						// 先添加主目标日志到logs
						logs = append(logs, session.BattleLogs[len(session.BattleLogs)-1])

						// 然后添加相邻目标的日志
						if adjacentLogsRaw, ok := skillState.Effect["_adjacentLogs"]; ok {
							if adjacentLogs, ok := adjacentLogsRaw.([]models.BattleLog); ok {
								for _, adjacentLog := range adjacentLogs {
									// 将日志添加到session并记录到logs
									if adjacentLog.LogType == "dodge" {
										// 闪避日志不需要伤害类型
										m.addLog(session, adjacentLog.LogType, adjacentLog.Message, adjacentLog.Color)
									} else {
										// 伤害日志需要伤害类型（使用日志中存储的DamageType，如果没有则使用physical）
										damageType := adjacentLog.DamageType
										if damageType == "" {
											damageType = "physical"
										}
										m.addLog(session, adjacentLog.LogType, adjacentLog.Message, adjacentLog.Color, withDamageType(damageType))
									}
									logs = append(logs, session.BattleLogs[len(session.BattleLogs)-1])
								}
								// 清理临时数据
								delete(skillState.Effect, "_adjacentLogs")
							}
						}
					}

					// 记录造成伤害的统计
					m.recordDamageDealt(session, char.ID, char.TeamSlot, playerDamage, playerDamageType, isCrit)

					// 如果是顺劈斩，记录波及伤害到统计
					totalSkillDamage := playerDamage // 技能总伤害（用于技能使用统计）
					if skillState != nil && skillState.SkillID == "warrior_cleave" {
						if adjacentTotalDamageRaw, ok := skillState.Effect["_adjacentTotalDamage"]; ok {
							if adjacentTotalDamage, ok := adjacentTotalDamageRaw.(int); ok && adjacentTotalDamage > 0 {
								// 波及伤害也计入统计（物理伤害，是否暴击取决于主目标是否暴击）
								m.recordDamageDealt(session, char.ID, char.TeamSlot, adjacentTotalDamage, "physical", isCrit)
								totalSkillDamage += adjacentTotalDamage // 累计到技能总伤害
								// 清理临时数据
								delete(skillState.Effect, "_adjacentTotalDamage")
							}
						}
					}

					// 记录技能使用统计（包含主目标和波及伤害的总和）
					skillID := ""
					if skillState != nil {
						skillID = skillState.SkillID
					}
					m.recordSkillUsage(session, char.ID, char.TeamSlot, skillID, totalSkillDamage, 0, resourceCost, true, isCrit)
				}
			} else {
				// 非攻击类技能（buff/debuff/control等）：只记录使用，不记录伤害
				m.addLog(session, "combat", fmt.Sprintf("%s 使用 [%s]%s", char.Name, skillName, resourceChangeText), "#8888ff")

				// 记录非伤害技能使用统计
				skillID := ""
				if skillState != nil {
					skillID = skillState.SkillID
				}
				m.recordSkillUsage(session, char.ID, char.TeamSlot, skillID, 0, 0, resourceCost, true, false)
			}
			// 对于顺劈斩，主目标和相邻目标的日志都已经在上面添加了，这里跳过避免重复
			// 对于其他技能，添加技能使用日志
			if skillState == nil || skillState.SkillID != "warrior_cleave" {
				logs = append(logs, session.BattleLogs[len(session.BattleLogs)-1])
			}
			// 顺劈斩的日志已经在上面处理完毕，不需要再添加

			// 处理技能特殊效果日志（在技能使用日志之后，闪避时不触发伤害相关效果）
			if skillEffects != nil && !isDodged {
				if stun, ok := skillEffects["stun"].(bool); ok && stun {
					m.addLog(session, "combat", fmt.Sprintf("%s 被眩晕了！", target.Name), "#ff00ff")
					logs = append(logs, session.BattleLogs[len(session.BattleLogs)-1])
				}
				// 处理基于伤害的恢复（嗜血等）
				if healPercent, ok := skillEffects["healPercent"].(float64); ok && usedSkill {
					healAmount := int(float64(playerDamage) * healPercent / 100.0)
					char.HP += healAmount
					if char.HP > char.MaxHP {
						char.HP = char.MaxHP
					}
					m.addLog(session, "heal", fmt.Sprintf("%s 恢复了 %d 点生命值", char.Name, healAmount), "#00ff00")
					logs = append(logs, session.BattleLogs[len(session.BattleLogs)-1])
				}
				// 处理破釜沉舟的立即恢复（基于最大HP）
				if healMaxHpPercent, ok := skillEffects["healMaxHpPercent"].(float64); ok && usedSkill {
					healAmount := int(float64(char.MaxHP) * healMaxHpPercent / 100.0)
					char.HP += healAmount
					if char.HP > char.MaxHP {
						char.HP = char.MaxHP
					}
					m.addLog(session, "heal", fmt.Sprintf("%s 的破釜沉舟恢复了 %d 点生命值", char.Name, healAmount), "#00ff00")
					logs = append(logs, session.BattleLogs[len(session.BattleLogs)-1])
				}
			}

			// 减少技能冷却时间
			m.skillManager.TickCooldowns(char.ID)

			// 减少Buff/Debuff持续时间
			expiredBuffs := m.buffManager.TickBuffs(char.ID)
			for _, expired := range expiredBuffs {
				m.addLog(session, "buff", fmt.Sprintf("%s 的 %s 效果消失了", char.Name, expired.Name), "#888888")
				logs = append(logs, session.BattleLogs[len(session.BattleLogs)-1])
			}

			// 处理DOT/HOT效果（在Buff持续时间减少之后）
			dotDamage, hotHealing := m.buffManager.ProcessDOTEffects(char.ID, session.CurrentBattleRound)
			if dotDamage > 0 {
				char.HP -= dotDamage
				if char.HP < 0 {
					char.HP = 0
				}
				m.addLog(session, "dot", fmt.Sprintf("%s 受到持续伤害，损失 %d 点生命值", char.Name, dotDamage), "#ff6666")
				logs = append(logs, session.BattleLogs[len(session.BattleLogs)-1])
			}
			if hotHealing > 0 {
				originalHP := char.HP
				char.HP += hotHealing
				if char.HP > char.MaxHP {
					char.HP = char.MaxHP
				}
				actualHealing := char.HP - originalHP
				if actualHealing > 0 {
					m.addLog(session, "hot", fmt.Sprintf("%s 的持续恢复效果恢复了 %d 点生命值", char.Name, actualHealing), "#00ff00")
					logs = append(logs, session.BattleLogs[len(session.BattleLogs)-1])
				}
			}

			// 检查目标是否死亡
			if target.HP <= 0 {
				// 确保HP归零
				target.HP = 0

				// 处理战争机器的击杀回怒效果
				m.handleWarMachineRageGain(char, session, &logs)

				// 处理被动技能的击杀时效果
				m.handlePassiveOnKillEffects(char, target, session, &logs)

				// 敌人死亡
				expGain := target.ExpReward
				goldGain := target.GoldMin + rand.Intn(target.GoldMax-target.GoldMin+1)

				// 应用区域收益倍率
				if session.CurrentZone != nil && m.zoneManager != nil {
					expMulti := m.zoneManager.CalculateExpMultiplier(session.CurrentZone.ID)
					goldMulti := m.zoneManager.CalculateGoldMultiplier(session.CurrentZone.ID)
					expGain = int(float64(expGain) * expMulti)
					goldGain = int(float64(goldGain) * goldMulti)
				}

				// 记录敌人死亡日志（敌人名字用红色，避免前端错误着色）
				m.addLog(session, "kill", fmt.Sprintf("💀 <span style=\"color: #ff7777\">%s</span> 被击杀！获得 <span style=\"color: #3d85c6\">%d</span> 经验、<span style=\"color: #ffd700\">%d</span> 金币", target.Name, expGain, goldGain), "#ff6b6b")
				logs = append(logs, session.BattleLogs[len(session.BattleLogs)-1])

				// 记录击杀统计
				m.recordKill(session, char.ID, char.TeamSlot)

				// 增加探索度（每击杀一个怪物增加1点探索度）
				if session.CurrentZone != nil && m.explorationRepo != nil {
					err := m.explorationRepo.AddExploration(session.UserID, session.CurrentZone.ID, 1)
					if err != nil {
						fmt.Printf("[WARN] Failed to add exploration: %v\n", err)
					}
				}

				session.CurrentBattleExp += expGain
				session.CurrentBattleGold += goldGain
				session.CurrentBattleKills++
				session.SessionExp += expGain
				session.SessionGold += goldGain
				session.SessionKills++

				char.Exp += expGain
				char.TotalKills++

				// 检查升级
				for char.Exp >= char.ExpToNext {
					char.Exp -= char.ExpToNext
					char.Level++
					char.ExpToNext = int(float64(char.ExpToNext) * 1.5)

					// 获得可分配属性点（不再自动增加主属性）
					char.UnspentPoints += 5

					// 升级时回满生命与资源（不改变上限）
					char.HP = char.MaxHP
					if char.ResourceType == "rage" {
						// 战士怒气上限固定为100，不重置怒气值
						char.MaxResource = 100
					} else if char.ResourceType == "energy" {
						// 盗贼等能量职业上限固定100，升级回满
						char.MaxResource = 100
						char.Resource = char.MaxResource
					} else {
						char.Resource = char.MaxResource
					}

					m.addLog(session, "levelup", fmt.Sprintf("🎉【升级】恭喜！%s 升到了 %d 级！", char.Name, char.Level), "#ffd700")
					logs = append(logs, session.BattleLogs[len(session.BattleLogs)-1])
				}
			}

			// 移动到下一个回合（使用TurnOrder系统）
			m.moveToNextTurn(session, characters, aliveEnemies)
		}
	} else {
		// 敌人回合：当前索引的敌人攻击玩家
		if session.CurrentTurnIndex < len(aliveEnemies) {
			enemy := aliveEnemies[session.CurrentTurnIndex]

			// 检查敌人是否处于眩晕状态
			enemyDebuffs := m.buffManager.GetEnemyDebuffs(enemy.ID)
			isStunned := false
			for _, debuff := range enemyDebuffs {
				if debuff.Type == "stun" {
					isStunned = true
					break
				}
			}

			if isStunned {
				// 敌人被眩晕，无法行动
				m.addLog(session, "combat", fmt.Sprintf("%s 处于眩晕状态，无法行动！", enemy.Name), "#ff00ff")
				logs = append(logs, session.BattleLogs[len(session.BattleLogs)-1])

				// 减少敌人debuff持续时间
				expiredDebuffs := m.buffManager.TickEnemyDebuffs(enemy.ID)
				for _, expiredID := range expiredDebuffs {
					// 检查是否是眩晕debuff
					if expiredID == "charge_stun" {
						m.addLog(session, "buff", fmt.Sprintf("%s 的眩晕效果消失了", enemy.Name), "#888888")
						logs = append(logs, session.BattleLogs[len(session.BattleLogs)-1])
					}
				}

				// 移动到下一个回合（使用TurnOrder系统）
				m.moveToNextTurn(session, characters, aliveEnemies)

				// 返回结果（眩晕，无行动）
				return &BattleTickResult{
					Character:    char,
					Enemy:        session.CurrentEnemy,
					Enemies:      session.CurrentEnemies,
					Logs:         logs,
					IsRunning:    session.IsRunning,
					IsResting:    session.IsResting,
					RestUntil:    session.RestUntil,
					SessionKills: session.SessionKills,
					SessionGold:  session.SessionGold,
					SessionExp:   session.SessionExp,
					BattleCount:  session.BattleCount,
				}, nil
			}

			// 【闪避判定】玩家尝试闪避敌人攻击
			playerDodgeRate := m.calculateCharacterDodgeRate(char)
			if m.checkDodge(playerDodgeRate, false) {
				// 闪避成功！
				m.addLog(session, "dodge", fmt.Sprintf("%s 闪避了 %s 的攻击！", char.Name, enemy.Name), "#00ffff")
				logs = append(logs, session.BattleLogs[len(session.BattleLogs)-1])

				// 记录闪避统计
				m.recordDodge(session, char.ID, char.TeamSlot)

				// 移动到下一个回合（使用TurnOrder系统）
				m.moveToNextTurn(session, characters, aliveEnemies)

				// 返回结果（闪避成功，无伤害）
				return &BattleTickResult{
					Character:    char,
					Enemy:        session.CurrentEnemy,
					Enemies:      session.CurrentEnemies,
					Logs:         logs,
					IsRunning:    session.IsRunning,
					IsResting:    session.IsResting,
					RestUntil:    session.RestUntil,
					SessionKills: session.SessionKills,
					SessionGold:  session.SessionGold,
					SessionExp:   session.SessionExp,
					BattleCount:  session.BattleCount,
				}, nil
			}

			// 决定敌人的攻击类型（物理/魔法）
			attackType := m.resolveEnemyAttackType(enemy)

			// 基础伤害计算（根据攻击类型选择不同的防御）
			var baseEnemyDamage int
			var enemyDamageDetails *DamageCalculationDetails
			if attackType == "magic" {
				baseEnemyDamage, enemyDamageDetails = m.calculateMagicDamageWithDetails(enemy.MagicAttack, char.MagicDefense)
			} else {
				baseEnemyDamage, enemyDamageDetails = m.calculatePhysicalDamageWithDetails(enemy.PhysicalAttack, char.PhysicalDefense)
			}

			enemyDamage := baseEnemyDamage

			// 敌人暴击判定
			var baseCritRate, baseCritDamage float64
			if attackType == "magic" {
				baseCritRate = enemy.SpellCritRate
				baseCritDamage = enemy.SpellCritDamage
			} else {
				baseCritRate = enemy.PhysCritRate
				baseCritDamage = enemy.PhysCritDamage
			}
			actualCritRate := baseCritRate
			if actualCritRate > 1.0 {
				actualCritRate = 1.0
			}
			critRoll := rand.Float64()
			isEnemyCrit := critRoll < actualCritRate
			if enemyDamageDetails != nil {
				enemyDamageDetails.BaseCritRate = baseCritRate
				enemyDamageDetails.ActualCritRate = actualCritRate
				enemyDamageDetails.RandomRoll = critRoll
				enemyDamageDetails.IsCrit = isEnemyCrit
				enemyDamageDetails.CritMultiplier = baseCritDamage
			}
			if isEnemyCrit {
				enemyDamage = int(float64(enemyDamage) * baseCritDamage)
			}

			// 应用buff/debuff效果（如盾牌格挡的减伤等）
			originalDamage := enemyDamage
			enemyDamage = m.buffManager.CalculateDamageTakenWithBuffs(enemyDamage, char.ID, true)
			if enemyDamage != originalDamage {
				reduction := float64(originalDamage-enemyDamage) / float64(originalDamage) * 100.0
				enemyDamageDetails.DefenseModifiers = append(enemyDamageDetails.DefenseModifiers,
					fmt.Sprintf("减伤Buff -%.0f%%", reduction))
			}

			// 处理被动技能的减伤效果（不灭意志等）
			originalDamage2 := enemyDamage
			enemyDamage = m.handlePassiveDamageReduction(char, enemyDamage)
			if enemyDamage != originalDamage2 {
				reduction := float64(originalDamage2-enemyDamage) / float64(originalDamage2) * 100.0
				enemyDamageDetails.DefenseModifiers = append(enemyDamageDetails.DefenseModifiers,
					fmt.Sprintf("被动减伤 -%.0f%%", reduction))
			}

			// 处理被动技能的受到伤害时效果
			m.handlePassiveOnDamageTakenEffects(char, enemyDamage, session, &logs)

			// 处理被动技能的受到伤害时效果
			m.handlePassiveOnDamageTakenEffects(char, enemyDamage, session, &logs)

			// 处理护盾效果（不灭壁垒等）
			shieldAmount := m.buffManager.GetBuffValue(char.ID, "shield")
			if shieldAmount > 0 {
				// 有护盾，先消耗护盾
				shieldInt := int(shieldAmount)
				if enemyDamage <= shieldInt {
					// 伤害完全被护盾吸收
					shieldInt -= enemyDamage
					absorbedDamage := enemyDamage
					enemyDamage = 0
					m.addLog(session, "shield", fmt.Sprintf("%s 的护盾吸收了 %d 点伤害", char.Name, absorbedDamage), "#00ffff")
					logs = append(logs, session.BattleLogs[len(session.BattleLogs)-1])
					// 更新护盾值（通过更新Buff的value）
					m.updateShieldValue(char.ID, float64(shieldInt))
				} else {
					// 护盾被击破，剩余伤害继续
					absorbedDamage := shieldInt
					enemyDamage -= shieldInt
					m.addLog(session, "shield", fmt.Sprintf("%s 的护盾吸收了 %d 点伤害后被击破", char.Name, absorbedDamage), "#00ffff")
					logs = append(logs, session.BattleLogs[len(session.BattleLogs)-1])
					m.updateShieldValue(char.ID, 0)
				}
			}

			if enemyDamageDetails != nil {
				enemyDamageDetails.FinalDamage = enemyDamage
			}

			// 处理被动技能的生存效果（坚韧不拔等）- 在受到伤害前检查
			originalHP := char.HP
			char.HP -= enemyDamage

			// 如果受到致命伤害，检查坚韧不拔效果
			if originalHP > 0 && char.HP <= 0 {
				if m.passiveSkillManager != nil {
					passives := m.passiveSkillManager.GetPassiveSkills(char.ID)
					for _, passive := range passives {
						if passive.Passive.EffectType == "survival" && passive.Passive.ID == "warrior_passive_unbreakable" {
							// 坚韧不拔：受到致命伤害时保留1点HP
							char.HP = 1
							m.addLog(session, "survival", fmt.Sprintf("%s 的坚韧不拔效果触发，保留了1点生命值！", char.Name), "#ff00ff")
							logs = append(logs, session.BattleLogs[len(session.BattleLogs)-1])
							break // 只触发一次
						}
					}
				}
			}

			// 处理反击效果（反击风暴、复仇被动等）
			m.handleCounterAttacks(char, enemy, enemyDamage, session, &logs)

			// 处理被动技能的反射效果（盾牌反射被动等）
			m.handlePassiveReflectEffects(char, enemy, enemyDamage, session, &logs)

			// 处理主动技能的反射效果（盾牌反射技能等）
			m.handleActiveReflectEffects(char, enemy, enemyDamage, session, &logs)

			// 记录受到伤害的统计
			m.recordDamageTaken(session, char.ID, char.TeamSlot, enemyDamage, attackType, 0, 0)

			// 保存资源变化前的值（用于日志显示）
			originalResource := char.Resource

			// 战士受到伤害时获得怒气
			if char.ResourceType == "rage" && enemyDamage > 0 {
				// 受到伤害获得怒气: 伤害/最大HP × 50，至少1点
				baseRageGain := int(float64(enemyDamage) / float64(char.MaxHP) * 50)
				if baseRageGain < 1 {
					baseRageGain = 1
				}

				// 应用被动技能的怒气获得加成（愤怒掌握等）
				rageGain := m.applyRageGenerationModifiers(char.ID, baseRageGain)

				char.Resource += rageGain
				if char.Resource > char.MaxResource {
					char.Resource = char.MaxResource
				}

				// 记录资源获得统计
				m.recordResourceGenerated(session, char.ID, char.TeamSlot, rageGain)
			}

			// 构建战斗日志消息，包含资源变化（带颜色）
			resourceChangeText := m.formatResourceChange(char.ResourceType, originalResource, char.Resource)

			// 格式化伤害公式
			enemyFormulaText := ""
			if enemyDamageDetails != nil {
				enemyFormulaText = m.formatDamageFormula(enemyDamageDetails)
			}

			// 格式化HP变化（使用已保存的originalHP）
			playerHPChangeText := m.formatHPChange(char.Name, originalHP, char.HP, char.MaxHP)

			damageColor := "#ff4444"
			if isEnemyCrit {
				m.addLog(session, "combat", fmt.Sprintf("%s 进行了💥暴击，对 %s 造成 %d 点伤害%s%s%s", enemy.Name, char.Name, enemyDamage, enemyFormulaText, playerHPChangeText, resourceChangeText), damageColor, withDamageType(attackType))
			} else {
				m.addLog(session, "combat", fmt.Sprintf("%s 攻击命中 %s，造成 %d 点伤害%s%s%s", enemy.Name, char.Name, enemyDamage, enemyFormulaText, playerHPChangeText, resourceChangeText), damageColor, withDamageType(attackType))
			}
			logs = append(logs, session.BattleLogs[len(session.BattleLogs)-1])

			// 减少敌人debuff持续时间
			expiredDebuffs := m.buffManager.TickEnemyDebuffs(enemy.ID)
			for _, expiredID := range expiredDebuffs {
				if expiredID == "charge_stun" {
					m.addLog(session, "buff", fmt.Sprintf("%s 的眩晕效果消失了", enemy.Name), "#888888")
					logs = append(logs, session.BattleLogs[len(session.BattleLogs)-1])
				}
			}

			// 检查玩家是否死亡
			if char.HP <= 0 {
				char.TotalDeaths++
				// 角色死亡时不停止战斗，保持 isRunning = true，这样休息状态可以自动处理
				// 用户已经开启了自动战斗，死亡只是暂时进入休息状态，休息结束后应该自动恢复战斗
				session.CurrentEnemies = nil
				session.CurrentEnemy = nil
				session.CurrentTurnIndex = -1

				// 角色死亡时，战士的怒气归0
				if char.ResourceType == "rage" {
					char.Resource = 0
				}

				// 计算复活时间
				reviveDuration := m.calculateReviveTime(userID)
				now := time.Now()
				reviveAt := now.Add(reviveDuration)

				// 设置角色HP为0（死亡状态）
				char.HP = 0
				char.IsDead = true
				char.ReviveAt = &reviveAt

				// 角色死亡时，立即清除所有buff和debuff，技能冷却重置
				if m.buffManager != nil {
					m.buffManager.ClearBuffs(char.ID)
				}
				// 清除技能状态（包括冷却时间）
				if m.skillManager != nil {
					m.skillManager.ClearCharacterSkills(char.ID)
				}

				m.addLog(session, "death", fmt.Sprintf("%s 被击败了... 需要 %d 秒复活", char.Name, int(reviveDuration.Seconds())), "#ff0000")
				logs = append(logs, session.BattleLogs[len(session.BattleLogs)-1])

				// 战斗失败总结
				m.addBattleSummary(session, false, &logs)

				// 记录角色死亡统计
				m.recordDeath(session, char.ID, char.TeamSlot)

				// 保存战斗统计到数据库（战斗失败）
				monsterID := ""
				if len(session.CurrentEnemies) > 0 && session.CurrentEnemies[0] != nil {
					monsterID = session.CurrentEnemies[0].ID
				}
				zoneID := ""
				if session.CurrentZone != nil {
					zoneID = session.CurrentZone.ID
				}
				m.saveBattleStats(session, session.UserID, zoneID, monsterID, false, characters)

				// 战斗失败时，战士的怒气归0
				if char.ResourceType == "rage" {
					char.Resource = 0
				}
				// 保存死亡数据（包括死亡标记、复活时间和怒气归0）
				m.charRepo.UpdateAfterDeath(char.ID, char.HP, char.Resource, char.TotalDeaths, &reviveAt)

				// 进入休息状态，休息时间 = 复活时间 + 恢复时间（恢复一半HP需要的时间）
				// 恢复时间：从0恢复到50% HP，每秒恢复2%，需要25秒
				recoveryTime := 25 * time.Second
				restDuration := reviveDuration + recoveryTime
				restUntil := now.Add(restDuration)
				session.IsResting = true
				session.RestUntil = &restUntil
				session.RestStartedAt = &now
				session.LastRestTick = &now
				session.RestSpeed = 1.0

				m.addLog(session, "system", fmt.Sprintf(">> 进入休息恢复状态 (预计 %d 秒)", int(restDuration.Seconds())+1), "#33ff33")
				logs = append(logs, session.BattleLogs[len(session.BattleLogs)-1])

				// 清除战斗统计收集器
				m.clearBattleStats(session)

				// 重置本场战斗统计
				session.CurrentBattleExp = 0
				session.CurrentBattleGold = 0
				session.CurrentBattleKills = 0
				session.CurrentTurnIndex = -1

				// 角色死亡时，立即返回，确保前端清除敌人显示
				// 保持 isRunning = true，这样按钮会显示"停止挂机"，休息状态可以自动处理
				return &BattleTickResult{
					Character:    char,
					Enemy:        nil,
					Enemies:      nil, // 明确返回 nil，让前端清除敌人显示
					Logs:         logs,
					IsRunning:    session.IsRunning, // 保持运行状态，不停止
					IsResting:    session.IsResting,
					RestUntil:    session.RestUntil,
					SessionKills: session.SessionKills,
					SessionGold:  session.SessionGold,
					SessionExp:   session.SessionExp,
					BattleCount:  session.BattleCount,
				}, nil
			} else {
				// 移动到下一个回合（使用TurnOrder系统）
				m.moveToNextTurn(session, characters, aliveEnemies)
			}
		} else {
			// 索引超出范围，重新构建TurnOrder
			m.buildTurnOrder(session, characters, aliveEnemies)
		}
	}

	// 更新存活敌人列表
	aliveEnemies = make([]*models.Monster, 0)
	for _, enemy := range session.CurrentEnemies {
		if enemy != nil && enemy.HP > 0 {
			aliveEnemies = append(aliveEnemies, enemy)
		}
	}

	// 如果所有敌人都被击败，处理战斗结束
	if len(aliveEnemies) == 0 && len(session.CurrentEnemies) > 0 {
		// 确保所有敌人的HP都归零
		for _, enemy := range session.CurrentEnemies {
			if enemy != nil && enemy.HP <= 0 {
				enemy.HP = 0
			}
		}

		// 处理怪物掉落
		m.processMonsterDrops(session, session.CurrentEnemies, &logs, characters)

		// 战斗胜利总结
		m.addBattleSummary(session, true, &logs)

		// 保存战斗统计到数据库
		monsterID := ""
		if len(session.CurrentEnemies) > 0 && session.CurrentEnemies[0] != nil {
			monsterID = session.CurrentEnemies[0].ID
		}
		zoneID := ""
		if session.CurrentZone != nil {
			zoneID = session.CurrentZone.ID
		}
		m.saveBattleStats(session, session.UserID, zoneID, monsterID, true, characters)

		// 战斗结束后，清除所有角色的buff和debuff，怒气归0，技能冷却重置
		for _, c := range characters {
			// 清除所有buff和debuff
			if m.buffManager != nil {
				m.buffManager.ClearBuffs(c.ID)
			}
			// 清除技能状态（包括冷却时间）
			if m.skillManager != nil {
				m.skillManager.ClearCharacterSkills(c.ID)
			}
			// 战士的怒气归0
			if c.ResourceType == "rage" {
				c.Resource = 0
			}
			// 保存所有角色的数据（包括战士的怒气归0）
			m.charRepo.UpdateAfterBattle(c.ID, c.HP, c.Resource, c.Exp, c.Level,
				c.ExpToNext, c.MaxHP, c.MaxResource, c.PhysicalAttack, c.MagicAttack, c.PhysicalDefense, c.MagicDefense,
				c.Strength, c.Agility, c.Intellect, c.Stamina, c.Spirit, c.UnspentPoints, c.TotalKills)
		}

		// 计算并开始休息
		restDuration := m.calculateRestTime(char)
		now := time.Now()
		restUntil := now.Add(restDuration)
		session.IsResting = true
		session.RestUntil = &restUntil
		session.RestStartedAt = &now
		session.LastRestTick = &now
		session.RestSpeed = 1.0 // 默认恢复速度

		m.addLog(session, "system", fmt.Sprintf(">> 开始休息恢复 (预计 %d 秒)", int(restDuration.Seconds())+1), "#33ff33")
		logs = append(logs, session.BattleLogs[len(session.BattleLogs)-1])

		// 清除战斗统计收集器
		m.clearBattleStats(session)

		// 重置本场战斗统计
		session.CurrentBattleExp = 0
		session.CurrentBattleGold = 0
		session.CurrentBattleKills = 0
		session.CurrentTurnIndex = -1

		// 先返回一次带有HP=0的敌人状态，让前端更新显示
		// 创建敌人副本，确保HP为0
		defeatedEnemies := make([]*models.Monster, len(session.CurrentEnemies))
		for i, enemy := range session.CurrentEnemies {
			if enemy != nil {
				defeatedEnemy := *enemy
				defeatedEnemy.HP = 0
				defeatedEnemies[i] = &defeatedEnemy
			}
		}

		// 清除敌人（在返回结果之后）
		session.CurrentEnemies = nil
		session.CurrentEnemy = nil

		// 返回带有HP=0的敌人状态
		return &BattleTickResult{
			Character:    char,
			Enemy:        nil,
			Enemies:      defeatedEnemies, // 返回HP=0的敌人副本
			Logs:         logs,
			IsRunning:    session.IsRunning,
			IsResting:    session.IsResting,
			RestUntil:    session.RestUntil,
			SessionKills: session.SessionKills,
			SessionGold:  session.SessionGold,
			SessionExp:   session.SessionExp,
			BattleCount:  session.BattleCount,
		}, nil
	}

	// 保存角色数据更新
	m.charRepo.UpdateAfterBattle(char.ID, char.HP, char.Resource, char.Exp, char.Level,
		char.ExpToNext, char.MaxHP, char.MaxResource, char.PhysicalAttack, char.MagicAttack, char.PhysicalDefense, char.MagicDefense,
		char.Strength, char.Agility, char.Intellect, char.Stamina, char.Spirit, char.UnspentPoints, char.TotalKills)

	return &BattleTickResult{
		Character:    char,
		Enemy:        session.CurrentEnemy,
		Enemies:      session.CurrentEnemies,
		Logs:         logs,
		IsRunning:    session.IsRunning,
		IsResting:    session.IsResting,
		RestUntil:    session.RestUntil,
		SessionKills: session.SessionKills,
		SessionGold:  session.SessionGold,
		SessionExp:   session.SessionExp,
		BattleCount:  session.BattleCount,
	}, nil
}

// spawnEnemies 生成多个敌人
// 敌人数量基于玩家角色数量：最高概率出现在等于玩家数量的敌人，最多相差不超过2
func (m *BattleManager) spawnEnemies(session *BattleSession, playerLevel int, playerCount int) error {
	if session.CurrentZone == nil {
		// 加载默认区域
		zone, err := m.gameRepo.GetZoneByID("elwynn")
		if err != nil {
			fmt.Printf("[ERROR] Failed to get zone: %v\n", err)
			return fmt.Errorf("failed to get zone: %v", err)
		}
		session.CurrentZone = zone
		// DEBUG输出仅在TEST_DEBUG=1时启用
	}

	// 获取区域怪物
	monsters, err := m.gameRepo.GetMonstersByZone(session.CurrentZone.ID)
	if err != nil {
		fmt.Printf("[ERROR] Failed to get monsters for zone %s: %v\n", session.CurrentZone.ID, err)
		return fmt.Errorf("failed to get monsters for zone %s: %v", session.CurrentZone.ID, err)
	}
	if len(monsters) == 0 {
		fmt.Printf("[ERROR] No monsters in zone %s (ID: %s)\n", session.CurrentZone.Name, session.CurrentZone.ID)
		// 如果当前区域没有怪物，尝试使用默认区域
		if session.CurrentZone.ID != "elwynn" {
			fmt.Printf("[WARN] Trying fallback to elwynn zone\n")
			fallbackZone, err := m.gameRepo.GetZoneByID("elwynn")
			if err == nil {
				monsters, err = m.gameRepo.GetMonstersByZone("elwynn")
				if err == nil && len(monsters) > 0 {
					session.CurrentZone = fallbackZone
					fmt.Printf("[INFO] Using fallback zone: elwynn\n")
				}
			}
		}
		if len(monsters) == 0 {
			return fmt.Errorf("no monsters available in zone %s", session.CurrentZone.ID)
		}
	}
	// fmt.Printf("[DEBUG] Found %d monsters in zone %s\n", len(monsters), session.CurrentZone.ID)

	// 基于玩家角色数量生成敌人数量（加权随机）
	// 敌人数量范围：max(1, playerCount-2) 到 playerCount+2
	// 权重：等于玩家数量的权重最高（5），相差1的权重为2，相差2的权重为1
	minEnemyCount := 1
	if playerCount > 2 {
		minEnemyCount = playerCount - 2
	}
	maxEnemyCount := playerCount + 2

	// 构建加权随机：每个可能的敌人数量及其权重
	type enemyCountWeight struct {
		count  int
		weight int
	}
	weights := make([]enemyCountWeight, 0)
	for count := minEnemyCount; count <= maxEnemyCount; count++ {
		diff := int(math.Abs(float64(count - playerCount)))
		// 权重：相差0（相等）=5（提高概率），相差1=2，相差2=1
		weight := 5 - diff*2
		if weight < 1 {
			weight = 1
		}
		weights = append(weights, enemyCountWeight{count: count, weight: weight})
	}

	// 计算总权重
	totalWeight := 0
	for _, w := range weights {
		totalWeight += w.weight
	}

	// 加权随机选择
	var enemyCount int
	if totalWeight <= 0 {
		// 如果总权重为0（理论上不应该发生），使用默认值
		enemyCount = playerCount
	} else {
		randomValue := rand.Intn(totalWeight)
		currentWeight := 0
		enemyCount = minEnemyCount // 默认值
		for _, w := range weights {
			currentWeight += w.weight
			if currentWeight > randomValue {
				enemyCount = w.count
				break
			}
		}
	}

	// 调试日志：记录生成的敌人数量
	// fmt.Printf("[DEBUG] Enemy count generation: playerCount=%d, enemyCount=%d, range=[%d,%d]\n",
	// 	playerCount, enemyCount, minEnemyCount, maxEnemyCount)

	// 重置威胁表（新战斗开始）
	m.resetThreatTable(session)

	session.CurrentEnemies = make([]*models.Monster, 0, enemyCount)

	var enemyNames []string
	for i := 0; i < enemyCount; i++ {
		// 优先使用 MonsterManager 生成怪物（配置化，支持平衡调整）
		var enemy *models.Monster
		var err error
		if m.monsterManager != nil {
			enemy, err = m.monsterManager.GenerateMonster(session.CurrentZone.ID, playerLevel)
		}

		// 如果生成失败，回退到旧方法
		if enemy == nil || err != nil {
			template := m.selectMonsterByWeight(monsters)
			enemy = &models.Monster{
				ID:              template.ID,
				ZoneID:          template.ZoneID,
				Name:            template.Name,
				Level:           template.Level,
				Type:            template.Type,
				HP:              template.HP,
				MaxHP:           template.HP,
				PhysicalAttack:  template.PhysicalAttack,
				MagicAttack:     template.MagicAttack,
				PhysicalDefense: template.PhysicalDefense,
				MagicDefense:    template.MagicDefense,
				ExpReward:       template.ExpReward,
				GoldMin:         template.GoldMin,
				GoldMax:         template.GoldMax,
			}
		}

		session.CurrentEnemies = append(session.CurrentEnemies, enemy)
		enemyNames = append(enemyNames, fmt.Sprintf("%s (Lv.%d)", enemy.Name, enemy.Level))
	}

	// 保留 CurrentEnemy 用于向后兼容（指向第一个敌人）
	if len(session.CurrentEnemies) > 0 {
		session.CurrentEnemy = session.CurrentEnemies[0]
	}

	session.BattleCount++
	if len(enemyNames) == 0 {
		return fmt.Errorf("failed to generate enemies")
	}
	enemyList := enemyNames[0]
	if len(enemyNames) > 1 {
		enemyList = fmt.Sprintf("%s 等 %d 个敌人", enemyNames[0], len(enemyNames))
	}
	m.addLog(session, "encounter", fmt.Sprintf("━━━ 战斗 #%d ━━━ 遭遇: %s", session.BattleCount, enemyList), "#00ff00")

	return nil
}

// Note: buildTurnOrder will be called after enemies are spawned in ExecuteBattleTick

// selectMonsterByWeight 根据权重随机选择怪物
// 使用加权随机算法：权重越高，被选中的概率越大
// 稀有怪物（elite/boss）的权重较低，普通怪物（normal）的权重较高
func (m *BattleManager) selectMonsterByWeight(monsters []models.Monster) models.Monster {
	if len(monsters) == 0 {
		return models.Monster{}
	}

	// 计算总权重
	totalWeight := 0
	for _, monster := range monsters {
		weight := monster.SpawnWeight
		if weight <= 0 {
			weight = 1 // 确保权重至少为1，避免除零错误
		}
		totalWeight += weight
	}

	if totalWeight == 0 {
		// 如果所有怪物权重都是0，使用简单随机选择
		return monsters[rand.Intn(len(monsters))]
	}

	// 生成 0 到 totalWeight 之间的随机数
	randomValue := rand.Intn(totalWeight)

	// 遍历怪物列表，累加权重，找到对应的怪物
	currentWeight := 0
	for _, monster := range monsters {
		weight := monster.SpawnWeight
		if weight <= 0 {
			weight = 1
		}
		currentWeight += weight
		if currentWeight > randomValue {
			return monster
		}
	}

	// 如果循环结束还没找到（理论上不应该发生），返回最后一个
	return monsters[len(monsters)-1]
}

// ChangeZone 切换区域
func (m *BattleManager) ChangeZone(userID int, zoneID string, playerLevel int, playerFaction string) error {
	session := m.GetOrCreateSession(userID)

	// 使用ZoneManager检查区域访问条件
	if m.zoneManager != nil {
		err := m.zoneManager.CheckZoneAccess(userID, zoneID, playerLevel, playerFaction)
		if err != nil {
			return err
		}
	}

	// 加载区域
	var zone *models.Zone
	var err error
	if m.zoneManager != nil {
		zone, err = m.zoneManager.GetZone(zoneID)
	} else {
		zone, err = m.gameRepo.GetZoneByID(zoneID)
	}
	if err != nil {
		return fmt.Errorf("zone not found: %s", zoneID)
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	session.CurrentZone = zone
	session.CurrentEnemy = nil
	session.CurrentEnemies = make([]*models.Monster, 0) // 清空所有敌人
	session.JustEncountered = false                     // 重置遭遇标志

	m.addLog(session, "zone", fmt.Sprintf(">> 你来到了 [%s]", zone.Name), "#00ffff")
	m.addLog(session, "zone", zone.Description, "#888888")

	return nil
}

// GetBattleStatus 获取战斗状态
func (m *BattleManager) GetBattleStatus(userID int) *models.BattleStatus {
	session := m.GetSession(userID)
	if session == nil {
		return &models.BattleStatus{
			IsRunning: false,
		}
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	status := &models.BattleStatus{
		IsRunning:      session.IsRunning,
		CurrentMonster: session.CurrentEnemy,
		CurrentEnemies: session.CurrentEnemies,
		BattleCount:    session.BattleCount,
		TotalKills:     session.SessionKills,
		TotalExp:       session.SessionExp,
		TotalGold:      session.SessionGold,
		IsResting:      session.IsResting,
		RestUntil:      session.RestUntil,
	}

	if session.CurrentZone != nil {
		status.CurrentZoneID = session.CurrentZone.ID
	}

	return status
}

// GetCharacterBuffs 获取角色的所有Buff/Debuff信息（用于API返回）
func (m *BattleManager) GetCharacterBuffs(characterID int) []*models.BuffInfo {
	if m.buffManager == nil {
		return []*models.BuffInfo{}
	}

	buffInstances := m.buffManager.GetBuffs(characterID)
	buffs := make([]*models.BuffInfo, 0, len(buffInstances))

	for _, buff := range buffInstances {
		description := m.getBuffDescription(buff)
		buffInfo := &models.BuffInfo{
			EffectID:     buff.EffectID,
			Name:         buff.Name,
			Type:         buff.Type,
			IsBuff:       buff.IsBuff,
			Duration:     buff.Duration,
			Value:        buff.Value,
			StatAffected: buff.StatAffected,
			Description:  description,
		}
		buffs = append(buffs, buffInfo)
	}

	return buffs
}

// getBuffDescription 获取Buff的描述文本
func (m *BattleManager) getBuffDescription(buff *BuffInstance) string {
	switch buff.StatAffected {
	case "attack":
		if buff.IsBuff {
			return fmt.Sprintf("提升%.0f%%物理攻击力", buff.Value)
		}
		return fmt.Sprintf("降低%.0f%%物理攻击力", -buff.Value)
	case "defense":
		if buff.IsBuff {
			return fmt.Sprintf("提升%.0f%%物理防御", buff.Value)
		}
		return fmt.Sprintf("降低%.0f%%物理防御", -buff.Value)
	case "physical_damage_taken":
		return fmt.Sprintf("减少%.0f%%受到的物理伤害", -buff.Value)
	case "damage_taken":
		return fmt.Sprintf("减少%.0f%%受到的伤害", -buff.Value)
	case "crit_rate":
		// 通用暴击率（同时影响物理和法术）
		if buff.IsBuff {
			return fmt.Sprintf("提升%.0f%%暴击率", buff.Value)
		}
		return fmt.Sprintf("降低%.0f%%暴击率", -buff.Value)
	case "phys_crit_rate":
		if buff.IsBuff {
			return fmt.Sprintf("提升%.0f%%物理暴击率", buff.Value)
		}
		return fmt.Sprintf("降低%.0f%%物理暴击率", -buff.Value)
	case "spell_crit_rate":
		if buff.IsBuff {
			return fmt.Sprintf("提升%.0f%%法术暴击率", buff.Value)
		}
		return fmt.Sprintf("降低%.0f%%法术暴击率", -buff.Value)
	case "healing_received":
		return fmt.Sprintf("降低%.0f%%治疗效果", buff.Value)
	case "shield":
		return fmt.Sprintf("获得相当于最大HP %.0f%%的护盾", buff.Value/float64(100))
	case "reflect":
		return fmt.Sprintf("反射%.0f%%受到的伤害", buff.Value)
	case "counter_attack":
		return fmt.Sprintf("受到攻击时反击，造成%.0f%%物理攻击力伤害", buff.Value)
	case "cc_immune":
		return "免疫控制效果"
	default:
		// 如果没有匹配的类型，返回buff名称
		return buff.Name
	}
}

// GetBattleLogs 获取战斗日志
func (m *BattleManager) GetBattleLogs(userID int, limit int) []models.BattleLog {
	session := m.GetSession(userID)
	if session == nil {
		return []models.BattleLog{}
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	logs := session.BattleLogs
	if limit > 0 && len(logs) > limit {
		logs = logs[len(logs)-limit:]
	}
	return logs
}

// checkDodge 检查闪避（返回 true 表示闪避成功）
// dodgeRate: 闪避率（0.0-1.0）
// ignoresDodge: 技能是否无视闪避
// 使用新的 Calculator 系统进行统一判定
func (m *BattleManager) checkDodge(dodgeRate float64, ignoresDodge bool) bool {
	// 如果技能无视闪避，直接返回 false（未闪避）
	if ignoresDodge {
		return false
	}

	// 使用 Calculator 进行闪避判定（内部会处理上限）
	return m.calculator.ShouldDodge(dodgeRate)
}

// skillIgnoresDodge 检查技能是否无视闪避
func (m *BattleManager) skillIgnoresDodge(skill *models.Skill) bool {
	if skill == nil || skill.Tags == "" {
		return false
	}
	// Tags 是 JSON 数组字符串，检查是否包含 "ignores_dodge"
	return strings.Contains(skill.Tags, "ignores_dodge")
}

// calculateCharacterDodgeRate 计算角色实际闪避率（包含被动和Buff加成）
func (m *BattleManager) calculateCharacterDodgeRate(char *models.Character) float64 {
	baseDodgeRate := char.DodgeRate

	// 应用被动技能的闪避率加成
	if m.passiveSkillManager != nil {
		dodgeModifier := m.passiveSkillManager.GetPassiveModifier(char.ID, "dodge_rate")
		if dodgeModifier > 0 {
			baseDodgeRate = baseDodgeRate + dodgeModifier/100.0
		}
	}

	// 应用Buff的闪避率加成
	if m.buffManager != nil {
		dodgeBuffValue := m.buffManager.GetBuffValue(char.ID, "dodge_rate")
		if dodgeBuffValue > 0 {
			baseDodgeRate = baseDodgeRate + dodgeBuffValue/100.0
		}
	}

	// 闪避率上限50%
	if baseDodgeRate > 0.5 {
		baseDodgeRate = 0.5
	}

	return baseDodgeRate
}

// DamageCalculationDetails 伤害计算详情
type DamageCalculationDetails struct {
	BaseAttack       int      // 基础攻击力
	ActualAttack     float64  // 实际攻击力（应用加成后）
	BaseDefense      int      // 基础防御力
	ActualDefense    float64  // 实际防御力（应用Debuff后）
	BaseDamage       float64  // 基础伤害（攻击-防御/2）
	FinalDamage      int      // 最终伤害（应用随机波动后）
	Variance         float64  // 随机波动值
	IsCrit           bool     // 是否暴击
	CritMultiplier   float64  // 暴击倍率
	BaseCritRate     float64  // 基础暴击率
	ActualCritRate   float64  // 实际暴击率（应用加成后）
	RandomRoll       float64  // 随机数（用于暴击判定）
	AttackModifiers  []string // 攻击力加成说明
	DefenseModifiers []string // 防御力修改说明
	CritModifiers    []string // 暴击率加成说明
	SkillRatio       float64  // 技能倍率（0表示普通攻击）
	ScaledDamage     float64  // 技能倍率后的伤害（攻击×倍率）
}

// calculateMagicDamageWithDetails 计算魔法伤害（返回详情）
// 使用新的 Calculator 系统进行统一计算
func (m *BattleManager) calculateMagicDamageWithDetails(attack, defense int) (int, *DamageCalculationDetails) {
	details := &DamageCalculationDetails{
		BaseAttack:       attack,
		ActualAttack:     float64(attack),
		BaseDefense:      defense,
		ActualDefense:    float64(defense),
		AttackModifiers:  []string{},
		DefenseModifiers: []string{},
	}

	// 基础伤害 = 攻击力
	baseDamage := float64(attack)
	details.BaseDamage = baseDamage

	// 应用防御减伤（减法公式：伤害 = 攻击 - 防御）
	damageAfterDefense := baseDamage - float64(defense)
	if damageAfterDefense < 1 {
		damageAfterDefense = 1
	}
	details.BaseDamage = damageAfterDefense

	details.Variance = 0
	details.FinalDamage = int(math.Round(damageAfterDefense))

	return details.FinalDamage, details
}

// calculatePhysicalDamage 计算物理伤害（返回详情）
// 使用新的 Calculator 系统进行统一计算
func (m *BattleManager) calculatePhysicalDamageWithDetails(attack, defense int) (int, *DamageCalculationDetails) {
	details := &DamageCalculationDetails{
		BaseAttack:       attack,
		ActualAttack:     float64(attack),
		BaseDefense:      defense,
		ActualDefense:    float64(defense),
		AttackModifiers:  []string{},
		DefenseModifiers: []string{},
	}

	// 基础伤害 = 攻击力
	baseDamage := float64(attack)
	details.BaseDamage = baseDamage

	// 应用防御减伤（减法公式：伤害 = 攻击 - 防御）
	damageAfterDefense := baseDamage - float64(defense)
	if damageAfterDefense < 1 {
		damageAfterDefense = 1
	}
	details.BaseDamage = damageAfterDefense

	details.Variance = 0 // 不再使用随机波动，未来通过装备的攻击力上下限实现
	details.FinalDamage = int(math.Round(damageAfterDefense))

	return details.FinalDamage, details
}

// resolveEnemyAttackType 决定敌人的攻击类型，如果未配置则按数值推断
func (m *BattleManager) resolveEnemyAttackType(enemy *models.Monster) string {
	if enemy.AttackType != "" {
		return enemy.AttackType
	}
	// 简单推断：如果魔法攻击更高且大于0，则使用魔法，否则物理
	if enemy.MagicAttack > enemy.PhysicalAttack && enemy.MagicAttack > 0 {
		return "magic"
	}
	return "physical"
}

// addLog 添加日志
type logOption func(*models.BattleLog)

func withDamageType(damageType string) logOption {
	return func(log *models.BattleLog) {
		if damageType != "" {
			log.DamageType = damageType
		}
	}
}

func normalizeDamageType(damageType string) string {
	switch strings.ToLower(damageType) {
	case "physical":
		return "physical"
	case "magic", "fire", "frost", "shadow", "holy", "nature", "arcane":
		return "magic"
	default:
		return ""
	}
}

func (m *BattleManager) addLog(session *BattleSession, logType, message, color string, opts ...logOption) {
	log := models.BattleLog{
		Message:   message,
		LogType:   logType,
		Color:     color,
		CreatedAt: time.Now(),
	}

	for _, opt := range opts {
		opt(&log)
	}
	session.BattleLogs = append(session.BattleLogs, log)

	// 保持日志数量在合理范围
	if len(session.BattleLogs) > 200 {
		session.BattleLogs = session.BattleLogs[len(session.BattleLogs)-200:]
	}
}

// addBattleSummary 添加战斗总结和分割线
func (m *BattleManager) addBattleSummary(session *BattleSession, isVictory bool, logs *[]models.BattleLog) {
	// 生成战斗总结，使用不同颜色标记不同指标
	var summaryMsg string
	battleDuration := time.Since(session.BattleStartTime)
	battleDurationSeconds := int(battleDuration.Seconds())

	if isVictory {
		if session.CurrentBattleKills > 0 {
			// 使用HTML标签为不同部分添加颜色
			// 结果：绿色 #00ff00，击杀：红色 #ff4444，经验：蓝色 #3d85c6，金币：金色 #ffd700，回合：紫色 #aa00ff
			summaryMsg = fmt.Sprintf("━━━ 战斗总结 ━━━ 结果: <span style=\"color: #00ff00\">✓ 胜利</span> | 击杀: <span style=\"color: #ff4444\">%d</span> | 经验: <span style=\"color: #3d85c6\">%d</span> | 金币: <span style=\"color: #ffd700\">%d</span> | 回合: <span style=\"color: #aa00ff\">%d</span> | 耗时: <span style=\"color: #888888\">%d秒</span>",
				session.CurrentBattleKills, session.CurrentBattleExp, session.CurrentBattleGold, session.CurrentBattleRound, battleDurationSeconds)
		} else {
			summaryMsg = fmt.Sprintf("━━━ 战斗总结 ━━━ 结果: <span style=\"color: #00ff00\">✓ 胜利</span> | 回合: <span style=\"color: #aa00ff\">%d</span> | 耗时: <span style=\"color: #888888\">%d秒</span>",
				session.CurrentBattleRound, battleDurationSeconds)
		}
		m.addLog(session, "battle_summary", summaryMsg, "#00ff00")
		*logs = append(*logs, session.BattleLogs[len(session.BattleLogs)-1])
	} else {
		// 失败时的总结
		if session.CurrentBattleKills > 0 {
			// 结果：红色 #ff6666，击杀：橙色 #ffaa00，经验：蓝色 #3d85c6，金币：金色 #ffd700，回合：紫色 #aa00ff
			summaryMsg = fmt.Sprintf("━━━ 战斗总结 ━━━ 结果: <span style=\"color: #ff6666\">✗ 失败</span> | 击杀: <span style=\"color: #ffaa00\">%d</span> | 经验: <span style=\"color: #3d85c6\">%d</span> | 金币: <span style=\"color: #ffd700\">%d</span> | 回合: <span style=\"color: #aa00ff\">%d</span> | 耗时: <span style=\"color: #888888\">%d秒</span>",
				session.CurrentBattleKills, session.CurrentBattleExp, session.CurrentBattleGold, session.CurrentBattleRound, battleDurationSeconds)
		} else {
			summaryMsg = fmt.Sprintf("━━━ 战斗总结 ━━━ 结果: <span style=\"color: #ff6666\">✗ 失败</span> | 回合: <span style=\"color: #aa00ff\">%d</span> | 耗时: <span style=\"color: #888888\">%d秒</span>",
				session.CurrentBattleRound, battleDurationSeconds)
		}
		m.addLog(session, "battle_summary", summaryMsg, "#ff6666")
		*logs = append(*logs, session.BattleLogs[len(session.BattleLogs)-1])
	}

	// 添加分割线
	m.addLog(session, "battle_separator", "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━", "#666666")
	*logs = append(*logs, session.BattleLogs[len(session.BattleLogs)-1])
}

// getResourceName 获取资源的中文名称
func (m *BattleManager) getResourceName(resourceType string) string {
	switch resourceType {
	case "rage":
		return "怒气"
	case "mana":
		return "MP"
	case "energy":
		return "能量"
	default:
		return "资源"
	}
}

// getResourceColor 获取资源的颜色（参考魔兽世界，但区别于伤害红色）
func (m *BattleManager) getResourceColor(resourceType string) string {
	switch resourceType {
	case "rage":
		return "#e25822" // 橙红色 - 怒气（区别于伤害的红色）
	case "mana":
		return "#3d85c6" // 蓝色 - 法力
	case "energy":
		return "#ffd700" // 金色/黄色 - 能量
	default:
		return "#ffffff" // 白色 - 默认
	}
}

// formatDamageFormula 格式化伤害计算公式文本（简洁版）
func (m *BattleManager) formatDamageFormula(details *DamageCalculationDetails) string {
	if details == nil {
		return ""
	}

	var parts []string

	// 检查是否为技能伤害（有技能倍率）
	isSkillDamage := details.SkillRatio > 0

	if isSkillDamage {
		// 技能伤害公式：攻击 × 倍率 - 防御 = 伤害
		// 使用四舍五入后的实际攻击力
		attackDisplay := int(math.Round(details.ActualAttack))
		if attackDisplay == 0 {
			attackDisplay = details.BaseAttack
		}

		// 计算实际数学结果（攻击×倍率-防御）
		scaledAttack := float64(attackDisplay) * details.SkillRatio
		rawDamage := scaledAttack - float64(details.BaseDefense)

		if rawDamage < 1 {
			// 如果计算结果小于1，显示实际计算和最低伤害说明
			baseFormula := fmt.Sprintf("%d攻 × %.1f - %d防 = %.0f → 最低1",
				attackDisplay, details.SkillRatio, details.BaseDefense, rawDamage)
			parts = append(parts, baseFormula)
		} else {
			baseFormula := fmt.Sprintf("%d攻 × %.1f - %d防 = %.0f",
				attackDisplay, details.SkillRatio, details.BaseDefense, details.BaseDamage)
			parts = append(parts, baseFormula)
		}

		// 如果有攻击力加成，显示加成说明
		if len(details.AttackModifiers) > 0 {
			modText := strings.Join(details.AttackModifiers, ", ")
			parts = append(parts, modText)
		}
	} else {
		// 普通攻击公式：攻击 - 防御 = 伤害
		// BaseAttack 已经是四舍五入后的实际计算值，直接使用
		attackUsed := details.BaseAttack

		// 计算实际数学结果
		rawDamage := attackUsed - details.BaseDefense
		if rawDamage < 1 {
			// 如果计算结果小于1，显示实际计算和最低伤害说明
			baseFormula := fmt.Sprintf("%d攻 - %d防 = %d → 最低1", attackUsed, details.BaseDefense, rawDamage)
			parts = append(parts, baseFormula)
		} else {
			baseFormula := fmt.Sprintf("%d攻 - %d防 = %d", attackUsed, details.BaseDefense, rawDamage)
			parts = append(parts, baseFormula)
		}
	}

	// 如果暴击，显示暴击计算
	if details.IsCrit && details.CritMultiplier > 0 {
		critFormula := fmt.Sprintf("%.0f × %.1f暴击 = %d", details.BaseDamage, details.CritMultiplier, details.FinalDamage)
		parts = append(parts, critFormula)
	}

	// 如果有防御修改（减伤等），简洁显示
	if len(details.DefenseModifiers) > 0 {
		modText := strings.Join(details.DefenseModifiers, ", ")
		parts = append(parts, modText)
	}

	if len(parts) == 0 {
		return ""
	}

	// 使用暗灰色显示公式（不抢眼，作为补充信息）
	// 注意：使用圆括号而非方括号，避免前端将其误识别为技能名
	formulaText := strings.Join(parts, " → ")
	return fmt.Sprintf(" <span style=\"color: #888888\">(%s)</span>", formulaText)
}

// formatHPChange 格式化HP变化显示
func (m *BattleManager) formatHPChange(name string, oldHP, newHP, maxHP int) string {
	// 计算HP百分比
	newPercent := float64(newHP) / float64(maxHP) * 100
	// 根据HP百分比选择颜色（使用青色系，区别于伤害红色）
	var color string
	if newPercent > 50 {
		color = "#4ecdc4" // 青绿色 - 健康
	} else if newPercent > 25 {
		color = "#ffe66d" // 淡黄色 - 警告
	} else {
		color = "#ff6b6b" // 珊瑚红 - 危险
	}
	// 使用尖括号避免与前端技能名识别冲突
	return fmt.Sprintf(" <span style=\"color: %s\">〈%s: %d→%d〉</span>", color, name, oldHP, newHP)
}

// formatResourceChange 格式化资源变化文本（带颜色），显示为 A->B 格式
func (m *BattleManager) formatResourceChange(resourceType string, originalValue int, finalValue int) string {
	if originalValue == finalValue {
		return ""
	}

	resourceName := m.getResourceName(resourceType)
	color := m.getResourceColor(resourceType)

	// 显示为 A->B 格式
	return fmt.Sprintf(" <span style=\"color: %s\">(%s %d->%d)</span>", color, resourceName, originalValue, finalValue)
}

// calculateReviveTime 计算复活时间（根据死亡人数）
func (m *BattleManager) calculateReviveTime(userID int) time.Duration {
	// 获取所有角色（所有角色都参与战斗）
	characters, err := m.charRepo.GetByUserID(userID)
	if err != nil {
		return 30 * time.Second // 默认30秒
	}

	// 统计死亡角色的数量
	deadCount := 0
	for _, char := range characters {
		if char.IsDead {
			deadCount++
		}
	}

	// 如果没有死亡角色，返回默认值
	if deadCount == 0 {
		deadCount = 1 // 至少有一个角色死亡才会调用这个函数
	}

	// 根据死亡人数计算复活时间（秒）
	var reviveSeconds int
	switch deadCount {
	case 1:
		reviveSeconds = 30
	case 2:
		reviveSeconds = 60
	case 3:
		reviveSeconds = 90
	case 4:
		reviveSeconds = 120
	default: // 5人或更多
		reviveSeconds = 180
	}

	return time.Duration(reviveSeconds) * time.Second
}

// calculateRestTime 计算休息时间（基于HP/MP损失）
// 注意：战士的怒气不需要恢复，战斗结束后直接归0，每场战斗从0开始
func (m *BattleManager) calculateRestTime(char *models.Character) time.Duration {
	hpLoss := float64(char.MaxHP - char.HP)

	// 战士的怒气不需要恢复，只计算HP损失
	// 其他职业需要计算MP损失
	var mpLoss float64
	if char.ResourceType != "rage" {
		mpLoss = float64(char.MaxResource - char.Resource)
	} else {
		// 战士的怒气不参与休息时间计算
		mpLoss = 0
	}

	// 如果已经满血满蓝（或满血），不需要休息
	if hpLoss <= 0 && mpLoss <= 0 {
		return 0
	}

	// 分别计算HP和MP的恢复时间
	// 每秒恢复2%，所以需要的时间 = 损失百分比 / 0.02 = 损失百分比 * 50
	hpLossPercent := hpLoss / float64(char.MaxHP)

	hpRestSeconds := hpLossPercent * 50.0
	var mpRestSeconds float64
	if char.ResourceType != "rage" && char.MaxResource > 0 {
		mpLossPercent := mpLoss / float64(char.MaxResource)
		mpRestSeconds = mpLossPercent * 50.0
	} else {
		mpRestSeconds = 0
	}

	// 取两者中的最大值，因为HP和MP是同时恢复的
	restSeconds := hpRestSeconds
	if mpRestSeconds > restSeconds {
		restSeconds = mpRestSeconds
	}

	// 最少1秒
	if restSeconds < 1.0 {
		restSeconds = 1.0
	}

	// 应用恢复速度倍率（未来可以从装备获取）
	restSpeed := 1.0 // 默认恢复速度
	if restSpeed > 0 {
		restSeconds = restSeconds / restSpeed
	}

	return time.Duration(restSeconds) * time.Second
}

// processRest 处理休息期间的恢复
func (m *BattleManager) processRest(session *BattleSession, char *models.Character) {
	if !session.IsResting || session.RestUntil == nil || session.RestStartedAt == nil {
		return
	}

	now := time.Now()

	// 检查角色是否已经复活（如果角色死亡且有复活时间）
	if char.IsDead && char.ReviveAt != nil {
		if now.After(*char.ReviveAt) || now.Equal(*char.ReviveAt) {
			// 如果MaxHP为0，先重新计算MaxHP
			if char.MaxHP == 0 {
				// 获取职业信息以获取BaseHP
				class, err := m.gameRepo.GetClassByID(char.ClassID)
				if err == nil && class != nil {
					char.MaxHP = m.calculator.CalculateHP(char, class.BaseHP)
				}
				// 如果仍然为0，使用默认值
				if char.MaxHP == 0 {
					char.MaxHP = 100 // 默认值
				}
			}
			// 复活时间到了，恢复角色到一半HP
			char.HP = char.MaxHP / 2
			if char.HP < 1 {
				char.HP = 1 // 至少1点HP
			}
			char.IsDead = false
			char.ReviveAt = nil

			// 更新数据库，清除死亡标记
			m.charRepo.SetDead(char.ID, false, nil)

			// 更新角色HP
			m.charRepo.UpdateAfterBattle(char.ID, char.HP, char.Resource, char.Exp, char.Level,
				char.ExpToNext, char.MaxHP, char.MaxResource, char.PhysicalAttack, char.MagicAttack, char.PhysicalDefense, char.MagicDefense,
				char.Strength, char.Agility, char.Intellect, char.Stamina, char.Spirit, char.UnspentPoints, char.TotalKills)

			// 记录复活日志
			m.addLog(session, "revive", fmt.Sprintf("%s 已复活，HP恢复至 %d/%d", char.Name, char.HP, char.MaxHP), "#00ff00")
		}
	}

	// 检查是否已经恢复满血满蓝，如果是则提前结束休息
	if char.HP >= char.MaxHP && char.Resource >= char.MaxResource {
		session.IsResting = false
		session.RestUntil = nil
		session.RestStartedAt = nil
		session.LastRestTick = nil
		return
	}

	if now.Before(*session.RestUntil) {
		// 计算从上次恢复到现在经过的时间
		var elapsed time.Duration
		if session.LastRestTick == nil {
			// 第一次恢复，从休息开始时间计算
			elapsed = now.Sub(*session.RestStartedAt)
		} else {
			// 从上次恢复时间计算
			elapsed = now.Sub(*session.LastRestTick)
		}

		// 如果时间间隔太长（超过1秒），限制为1秒，避免一次性恢复过多
		if elapsed > time.Second {
			elapsed = time.Second
		}

		// 如果时间间隔太短（小于0.1秒），跳过恢复，避免频繁计算
		if elapsed < 100*time.Millisecond {
			return
		}

		// 计算恢复速度（每秒恢复最大值的2%）
		restSpeed := session.RestSpeed
		if restSpeed <= 0 {
			restSpeed = 1.0
		}

		// 计算经过的秒数
		elapsedSeconds := elapsed.Seconds()

		// 如果MaxHP为0，先重新计算MaxHP
		if char.MaxHP == 0 {
			// 获取职业信息以获取BaseHP
			class, err := m.gameRepo.GetClassByID(char.ClassID)
			if err == nil && class != nil {
				char.MaxHP = m.calculator.CalculateHP(char, class.BaseHP)
			}
			// 如果仍然为0，使用默认值
			if char.MaxHP == 0 {
				char.MaxHP = 100 // 默认值
			}
		}

		// 如果角色已经死亡但还没到复活时间，不恢复HP
		if char.IsDead && char.ReviveAt != nil && now.Before(*char.ReviveAt) {
			// 只恢复资源（如果适用），不恢复HP
		} else {
			// 根据实际经过的时间计算恢复量
			hpRegenPercent := 0.02 * restSpeed * elapsedSeconds // 每秒2%

			hpRegen := int(float64(char.MaxHP) * hpRegenPercent)

			// 确保至少恢复1点（如果还没满）
			if hpRegen < 1 && char.HP < char.MaxHP {
				hpRegen = 1
			}

			char.HP += hpRegen
			if char.HP > char.MaxHP {
				char.HP = char.MaxHP
			}
		}

		// 战士的怒气不在休息时恢复，只在战斗中通过攻击/受击获得
		if char.ResourceType != "rage" {
			mpRegenPercent := 0.02 * restSpeed * elapsedSeconds
			mpRegen := int(float64(char.MaxResource) * mpRegenPercent)

			if mpRegen < 1 && char.Resource < char.MaxResource {
				mpRegen = 1
			}

			char.Resource += mpRegen
			if char.Resource > char.MaxResource {
				char.Resource = char.MaxResource
			}
		}

		// 更新上次恢复时间
		session.LastRestTick = &now

		// 恢复后再次检查是否已满，如果满了则提前结束休息
		if char.HP >= char.MaxHP && char.Resource >= char.MaxResource {
			session.IsResting = false
			session.RestUntil = nil
			session.RestStartedAt = nil
			session.LastRestTick = nil
		}
	} else {
		// 休息时间到了，结束休息
		// 确保HP已满
		if char.HP < char.MaxHP {
			char.HP = char.MaxHP
		}
		// 战士的怒气不在休息时恢复，只在战斗中通过攻击/受击获得
		if char.ResourceType != "rage" {
			if char.Resource < char.MaxResource {
				char.Resource = char.MaxResource
			}
		}
		session.IsResting = false
		session.RestUntil = nil
		session.RestStartedAt = nil
		session.LastRestTick = nil
	}
}

// BattleTickResult 战斗回合结果
type BattleTickResult struct {
	Character    *models.Character  `json:"character"`
	Enemy        *models.Monster    `json:"enemy,omitempty"`
	Enemies      []*models.Monster  `json:"enemies,omitempty"` // 多个敌人支持
	Logs         []models.BattleLog `json:"logs"`
	IsRunning    bool               `json:"isRunning"`
	IsResting    bool               `json:"isResting"`           // 是否在休息
	RestUntil    *time.Time         `json:"restUntil,omitempty"` // 休息结束时间
	SessionKills int                `json:"sessionKills"`
	SessionGold  int                `json:"sessionGold"`
	SessionExp   int                `json:"sessionExp"`
	BattleCount  int                `json:"battleCount"`
}

// applySkillBuffs 应用技能的Buff/Debuff效果
func (m *BattleManager) applySkillBuffs(skillState *CharacterSkillState, character *models.Character, target *models.Monster, skillEffects map[string]interface{}) {
	skill := skillState.Skill
	effect := skillState.Effect

	switch skill.ID {
	case "warrior_shield_block":
		// 盾牌格挡：减少受到的物理伤害
		if damageReduction, ok := effect["damageReduction"].(float64); ok {
			duration := 2
			if d, ok := effect["duration"].(int); ok {
				duration = d
			}
			m.buffManager.ApplyBuff(character.ID, "shield_block", "盾牌格挡", "buff", true, duration, -damageReduction, "physical_damage_taken", "")
		}
	case "warrior_battle_shout":
		// 战斗怒吼：提升攻击力
		if attackBonus, ok := effect["attackBonus"].(float64); ok {
			duration := 5
			if d, ok := effect["duration"].(int); ok {
				duration = d
			}
			m.buffManager.ApplyBuff(character.ID, "battle_shout", "战斗怒吼", "buff", true, duration, attackBonus, "attack", "")
		}
	case "warrior_demoralizing_shout":
		// 挫志怒吼：降低所有敌人攻击力（在applySkillDebuffs中处理）
	case "warrior_whirlwind":
		// 旋风斩：降低所有敌人防御（在applySkillDebuffs中处理）
	case "warrior_mortal_strike":
		// 致死打击：降低目标治疗效果
		if healingReduction, ok := effect["healingReduction"].(float64); ok {
			duration := 3
			if d, ok := effect["debuffDuration"].(float64); ok {
				duration = int(d)
			}
			// 应用到目标敌人
			if target != nil {
				m.buffManager.ApplyEnemyDebuff(target.ID, "mortal_strike", "致死打击", "debuff", duration, healingReduction, "healing_received", "")
			}
		}
	case "warrior_last_stand":
		// 破釜沉舟：立即恢复最大HP的百分比
		if healPercent, ok := effect["healPercent"].(float64); ok {
			// 立即恢复
			healAmount := int(float64(character.MaxHP) * healPercent / 100.0)
			character.HP += healAmount
			if character.HP > character.MaxHP {
				character.HP = character.MaxHP
			}
			// 通过skillEffects传递，在战斗日志中显示
			skillEffects["healMaxHpPercent"] = healPercent
		}
	case "warrior_unbreakable_barrier":
		// 不灭壁垒：获得护盾
		if shieldPercent, ok := effect["shieldPercent"].(float64); ok {
			duration := 4
			if d, ok := effect["duration"].(float64); ok {
				duration = int(d)
			}
			shieldAmount := int(float64(character.MaxHP) * shieldPercent / 100.0)
			// 使用Buff存储护盾值，statAffected为"shield"，value为护盾值
			m.buffManager.ApplyBuff(character.ID, "unbreakable_barrier", "不灭壁垒", "buff", true, duration, float64(shieldAmount), "shield", "")
		}
	case "warrior_shield_reflection":
		// 盾牌反射：反射受到的伤害
		if reflectPercent, ok := effect["reflectPercent"].(float64); ok {
			duration := 2
			if d, ok := effect["duration"].(float64); ok {
				duration = int(d)
			}
			// 使用Buff存储反射比例，statAffected为"reflect"，value为反射百分比
			m.buffManager.ApplyBuff(character.ID, "shield_reflection", "盾牌反射", "buff", true, duration, reflectPercent, "reflect", "")
		}
	case "warrior_shield_wall":
		// 盾墙：大幅减少受到的伤害
		if damageReduction, ok := effect["damageReduction"].(float64); ok {
			duration := 2
			if d, ok := effect["duration"].(float64); ok {
				duration = int(d)
			}
			m.buffManager.ApplyBuff(character.ID, "shield_wall", "盾墙", "buff", true, duration, -damageReduction, "damage_taken", "")
		}
	case "warrior_recklessness":
		// 鲁莽：提升暴击率，但受到伤害增加
		if critBonus, ok := effect["critBonus"].(float64); ok {
			duration := 3
			if d, ok := effect["duration"].(float64); ok {
				duration = int(d)
			}
			m.buffManager.ApplyBuff(character.ID, "recklessness_crit", "鲁莽", "buff", true, duration, critBonus, "crit_rate", "")
		}
		if damageIncrease, ok := effect["damageTakenIncrease"].(float64); ok {
			duration := 3
			if d, ok := effect["duration"].(float64); ok {
				duration = int(d)
			}
			m.buffManager.ApplyBuff(character.ID, "recklessness_damage", "鲁莽", "debuff", false, duration, damageIncrease, "damage_taken", "")
		}
	case "warrior_retaliation":
		// 反击风暴：受到攻击时反击
		if counterDamage, ok := effect["counterDamagePercent"].(float64); ok {
			duration := 3
			if d, ok := effect["duration"].(float64); ok {
				duration = int(d)
			}
			m.buffManager.ApplyBuff(character.ID, "retaliation", "反击风暴", "buff", true, duration, counterDamage, "counter_attack", "")
		}
	case "warrior_berserker_rage":
		// 狂暴之怒：提升攻击力和怒气获取
		if attackBonus, ok := effect["attackBonus"].(float64); ok {
			duration := 4
			if d, ok := effect["duration"].(float64); ok {
				duration = int(d)
			}
			m.buffManager.ApplyBuff(character.ID, "berserker_rage", "狂暴之怒", "buff", true, duration, attackBonus, "attack", "")
		}
	case "warrior_avatar":
		// 天神下凡：大幅提升攻击力，免疫控制
		if attackBonus, ok := effect["attackBonus"].(float64); ok {
			duration := 3
			if d, ok := effect["duration"].(float64); ok {
				duration = int(d)
			}
			m.buffManager.ApplyBuff(character.ID, "avatar", "天神下凡", "buff", true, duration, attackBonus, "attack", "")
		}
		if immuneCC, ok := effect["immuneCC"].(bool); ok && immuneCC {
			duration := 3
			if d, ok := effect["duration"].(float64); ok {
				duration = int(d)
			}
			m.buffManager.ApplyBuff(character.ID, "avatar_cc_immune", "天神下凡", "buff", true, duration, 1.0, "cc_immune", "")
		}
	}
}

// handleCounterAttacks 处理反击效果
func (m *BattleManager) handleCounterAttacks(character *models.Character, attacker *models.Monster, damageTaken int, session *BattleSession, logs *[]models.BattleLog) {
	// 处理Buff的反击效果（反击风暴）
	buffs := m.buffManager.GetBuffs(character.ID)
	for _, buff := range buffs {
		if buff.StatAffected == "counter_attack" && buff.IsBuff {
			// 反击风暴：对攻击者造成反击伤害
			counterDamage := int(float64(character.PhysicalAttack) * buff.Value / 100.0)
			attackerOldHP := attacker.HP
			attacker.HP -= counterDamage
			if attacker.HP < 0 {
				attacker.HP = 0
			}
			// 更新威胁值（反击也产生威胁）
			m.updateThreat(session, attacker.ID, character.ID, counterDamage)
			counterHPChange := m.formatHPChange(attacker.Name, attackerOldHP, attacker.HP, attacker.MaxHP)
			m.addLog(session, "combat", fmt.Sprintf("%s 的反击风暴对 %s 造成 %d 点反击伤害%s", character.Name, attacker.Name, counterDamage, counterHPChange), "#ff8800", withDamageType("physical"))
			*logs = append(*logs, session.BattleLogs[len(session.BattleLogs)-1])
		}
	}

	// 处理被动技能的反击效果（复仇）
	if m.passiveSkillManager != nil {
		passives := m.passiveSkillManager.GetPassiveSkills(character.ID)
		for _, passive := range passives {
			if passive.Passive.EffectType == "counter_attack" {
				// 复仇：受到攻击时概率反击
				// effectValue是触发概率（百分比），需要根据等级计算实际概率和伤害
				triggerChance := passive.EffectValue / 100.0
				if rand.Float64() < triggerChance {
					// 计算反击伤害（根据等级：1级100%，5级180%）
					counterDamagePercent := 100.0 + float64(passive.Level-1)*20.0
					// 计算实际攻击力（应用被动技能和Buff加成）
					actualAttack := float64(character.PhysicalAttack)
					if m.passiveSkillManager != nil {
						attackModifier := m.passiveSkillManager.GetPassiveModifier(character.ID, "attack")
						actualAttack = actualAttack * (1.0 + attackModifier/100.0)
					}
					if m.buffManager != nil {
						attackBuffValue := m.buffManager.GetBuffValue(character.ID, "attack")
						if attackBuffValue > 0 {
							actualAttack = actualAttack * (1.0 + attackBuffValue/100.0)
						}
					}
					counterDamage := int(actualAttack * counterDamagePercent / 100.0)
					counterDamage = counterDamage - attacker.PhysicalDefense/2
					if counterDamage < 1 {
						counterDamage = 1
					}
					revengeOldHP := attacker.HP
					attacker.HP -= counterDamage
					if attacker.HP < 0 {
						attacker.HP = 0
					}
					// 更新威胁值（复仇反击也产生威胁）
					m.updateThreat(session, attacker.ID, character.ID, counterDamage)
					revengeHPChange := m.formatHPChange(attacker.Name, revengeOldHP, attacker.HP, attacker.MaxHP)
					m.addLog(session, "combat", fmt.Sprintf("%s 的复仇对 %s 造成 %d 点反击伤害%s", character.Name, attacker.Name, counterDamage, revengeHPChange), "#ff8800", withDamageType("physical"))
					*logs = append(*logs, session.BattleLogs[len(session.BattleLogs)-1])
				}
			}
		}
	}
}

// handlePassiveOnHitEffects 处理被动技能的攻击时效果
func (m *BattleManager) handlePassiveOnHitEffects(character *models.Character, damageDealt int, usedSkill bool, session *BattleSession, logs *[]models.BattleLog) {
	if m.passiveSkillManager == nil {
		return
	}

	passives := m.passiveSkillManager.GetPassiveSkills(character.ID)
	for _, passive := range passives {
		if passive.Passive == nil {
			continue // 跳过Passive为nil的情况
		}
		switch passive.Passive.EffectType {
		case "on_hit_heal":
			// 血之狂热：每次攻击恢复生命值
			healPercent := passive.EffectValue // 百分比值（如1.0表示1%）
			// 使用浮点数计算，然后四舍五入
			healAmountFloat := float64(character.MaxHP) * healPercent / 100.0
			healAmount := int(healAmountFloat + 0.5) // 四舍五入
			// 确保至少恢复1点（如果计算为0但EffectValue>0）
			if healPercent > 0 && healAmount == 0 && character.MaxHP > 0 {
				healAmount = 1
			}
			if healAmount > 0 {
				oldHP := character.HP
				character.HP += healAmount
				if character.HP > character.MaxHP {
					character.HP = character.MaxHP
				}
				// 只有在HP实际增加时才记录日志
				if character.HP > oldHP {
					m.addLog(session, "heal", fmt.Sprintf("%s 的血之狂热恢复了 %d 点生命值", character.Name, healAmount), "#00ff00")
					*logs = append(*logs, session.BattleLogs[len(session.BattleLogs)-1])
				}
			}
		case "on_hit_resource":
			// 攻击时获得资源（如怒气、法力）
			resourceGain := int(passive.EffectValue)
			if resourceGain > 0 {
				character.Resource += resourceGain
				if character.Resource > character.MaxResource {
					character.Resource = character.MaxResource
				}
				resourceName := m.getResourceName(character.ResourceType)
				m.addLog(session, "resource", fmt.Sprintf("%s 获得了 %d 点%s", character.Name, resourceGain, resourceName), "#8888ff")
				*logs = append(*logs, session.BattleLogs[len(session.BattleLogs)-1])
			}
		}
	}
}

// handlePassiveOnCritEffects 处理被动技能的暴击时效果
func (m *BattleManager) handlePassiveOnCritEffects(character *models.Character, critDamage int, usedSkill bool, session *BattleSession, logs *[]models.BattleLog) {
	if m.passiveSkillManager == nil {
		return
	}

	passives := m.passiveSkillManager.GetPassiveSkills(character.ID)
	for _, passive := range passives {
		if passive.Passive == nil {
			continue
		}
		switch passive.Passive.EffectType {
		case "on_crit_heal":
			// 暴击时恢复生命值
			healPercent := passive.EffectValue // 百分比值（基于暴击伤害）
			healAmount := int(float64(critDamage) * healPercent / 100.0)
			if healAmount > 0 {
				character.HP += healAmount
				if character.HP > character.MaxHP {
					character.HP = character.MaxHP
				}
				m.addLog(session, "heal", fmt.Sprintf("%s 的暴击恢复效果恢复了 %d 点生命值", character.Name, healAmount), "#00ff00")
				*logs = append(*logs, session.BattleLogs[len(session.BattleLogs)-1])
			}
		case "on_crit_resource":
			// 暴击时获得额外资源
			resourceGain := int(passive.EffectValue)
			if resourceGain > 0 {
				character.Resource += resourceGain
				if character.Resource > character.MaxResource {
					character.Resource = character.MaxResource
				}
				resourceName := m.getResourceName(character.ResourceType)
				m.addLog(session, "resource", fmt.Sprintf("%s 的暴击获得了额外 %d 点%s", character.Name, resourceGain, resourceName), "#8888ff")
				*logs = append(*logs, session.BattleLogs[len(session.BattleLogs)-1])
			}
		}
	}
}

// handlePassiveOnKillEffects 处理被动技能的击杀时效果
func (m *BattleManager) handlePassiveOnKillEffects(character *models.Character, killedEnemy *models.Monster, session *BattleSession, logs *[]models.BattleLog) {
	if m.passiveSkillManager == nil {
		return
	}

	passives := m.passiveSkillManager.GetPassiveSkills(character.ID)
	for _, passive := range passives {
		if passive.Passive == nil {
			continue
		}
		switch passive.Passive.EffectType {
		case "on_kill_heal":
			// 击杀时恢复生命值
			healPercent := passive.EffectValue // 百分比值（基于最大HP）
			healAmount := int(float64(character.MaxHP) * healPercent / 100.0)
			if healAmount > 0 {
				character.HP += healAmount
				if character.HP > character.MaxHP {
					character.HP = character.MaxHP
				}
				m.addLog(session, "heal", fmt.Sprintf("%s 的击杀恢复效果恢复了 %d 点生命值", character.Name, healAmount), "#00ff00")
				*logs = append(*logs, session.BattleLogs[len(session.BattleLogs)-1])
			}
		case "on_kill_resource":
			// 击杀时获得资源（战争机器等）
			// 这个已经在handleWarMachineRageGain中处理，这里可以添加其他被动技能
			resourceGain := int(passive.EffectValue)
			if resourceGain > 0 && passive.Passive.ID != "warrior_passive_war_machine" {
				character.Resource += resourceGain
				if character.Resource > character.MaxResource {
					character.Resource = character.MaxResource
				}
				resourceName := m.getResourceName(character.ResourceType)
				m.addLog(session, "resource", fmt.Sprintf("%s 的击杀获得了 %d 点%s", character.Name, resourceGain, resourceName), "#8888ff")
				*logs = append(*logs, session.BattleLogs[len(session.BattleLogs)-1])
			}
		}
	}
}

// handlePassiveDamageReduction 处理被动技能的减伤效果
func (m *BattleManager) handlePassiveDamageReduction(character *models.Character, damage int) int {
	if m.passiveSkillManager == nil {
		return damage
	}

	passives := m.passiveSkillManager.GetPassiveSkills(character.ID)
	for _, passive := range passives {
		if passive.Passive.EffectType == "survival" && passive.Passive.ID == "warrior_passive_unbreakable_will" {
			// 不灭意志：HP低于阈值时减伤
			hpPercent := float64(character.HP) / float64(character.MaxHP)
			// 根据等级计算触发阈值（1级30%，5级10%）
			threshold := 0.30 - float64(passive.Level-1)*0.05
			if hpPercent < threshold {
				// 根据等级计算减伤比例（1级25%，5级65%）
				reductionPercent := 25.0 + float64(passive.Level-1)*10.0
				damage = int(float64(damage) * (1.0 - reductionPercent/100.0))
				if damage < 1 {
					damage = 1
				}
			}
		}
	}

	return damage
}

// handleActiveReflectEffects 处理主动技能的反射效果
func (m *BattleManager) handleActiveReflectEffects(character *models.Character, attacker *models.Monster, damageTaken int, session *BattleSession, logs *[]models.BattleLog) {
	if m.buffManager == nil {
		return
	}

	buffs := m.buffManager.GetBuffs(character.ID)
	for _, buff := range buffs {
		if buff.StatAffected == "reflect" && buff.IsBuff && buff.EffectID == "shield_reflection" {
			// 盾牌反射（主动技能）：反射受到的伤害
			reflectPercent := buff.Value // 百分比值（如50.0表示50%）
			reflectDamage := int(float64(damageTaken) * reflectPercent / 100.0)
			if reflectDamage > 0 {
				reflectOldHP := attacker.HP
				attacker.HP -= reflectDamage
				if attacker.HP < 0 {
					attacker.HP = 0
				}
				// 更新威胁值（反射伤害也产生威胁）
				m.updateThreat(session, attacker.ID, character.ID, reflectDamage)
				reflectHPChange := m.formatHPChange(attacker.Name, reflectOldHP, attacker.HP, attacker.MaxHP)
				m.addLog(session, "combat", fmt.Sprintf("%s 的盾牌反射对 %s 造成 %d 点反射伤害%s", character.Name, attacker.Name, reflectDamage, reflectHPChange), "#ff8800", withDamageType("magic"))
				*logs = append(*logs, session.BattleLogs[len(session.BattleLogs)-1])
			}
		}
	}
}

// updateShieldValue 更新护盾值
func (m *BattleManager) updateShieldValue(characterID int, newShieldValue float64) {
	if m.buffManager == nil {
		return
	}

	buffs := m.buffManager.GetBuffs(characterID)
	if buff, exists := buffs["unbreakable_barrier"]; exists {
		buff.Value = newShieldValue
	}
}

// applySkillDebuffs 应用技能的Debuff效果到敌人
func (m *BattleManager) applySkillDebuffs(skillState *CharacterSkillState, character *models.Character, target *models.Monster, allEnemies []*models.Monster, skillEffects map[string]interface{}) {
	skill := skillState.Skill
	effect := skillState.Effect

	// 处理眩晕效果（冲锋等技能）
	if stun, ok := skillEffects["stun"].(bool); ok && stun {
		stunDuration := 1 // 默认1回合
		if duration, ok := skillEffects["stunDuration"].(int); ok {
			stunDuration = duration
		} else if duration, ok := skillEffects["stunDuration"].(float64); ok {
			stunDuration = int(duration)
		}
		// 应用到目标敌人
		if target != nil && target.HP > 0 {
			m.buffManager.ApplyEnemyDebuff(target.ID, "charge_stun", "冲锋眩晕", "stun", stunDuration, 0, "", "")
		}
	}

	switch skill.ID {
	case "warrior_demoralizing_shout":
		// 挫志怒吼：降低所有敌人攻击力
		if attackReduction, ok := effect["attackReduction"].(float64); ok {
			duration := 3
			if d, ok := effect["duration"].(float64); ok {
				duration = int(d)
			}
			// 应用到所有存活的敌人
			for _, enemy := range allEnemies {
				if enemy.HP > 0 {
					m.buffManager.ApplyEnemyDebuff(enemy.ID, "demoralizing_shout", "挫志怒吼", "debuff", duration, attackReduction, "attack", "")
				}
			}
		}
	case "warrior_whirlwind":
		// 旋风斩：降低所有敌人防御
		if defenseReduction, ok := effect["defenseReduction"].(float64); ok {
			duration := 2
			if d, ok := effect["debuffDuration"].(float64); ok {
				duration = int(d)
			}
			// 应用到所有存活的敌人
			for _, enemy := range allEnemies {
				if enemy.HP > 0 {
					m.buffManager.ApplyEnemyDebuff(enemy.ID, "whirlwind", "旋风斩", "debuff", duration, defenseReduction, "defense", "")
				}
			}
		}
	case "warrior_mortal_strike":
		// 致死打击：降低目标治疗效果
		if healingReduction, ok := effect["healingReduction"].(float64); ok {
			duration := 3
			if d, ok := effect["debuffDuration"].(float64); ok {
				duration = int(d)
			}
			// 应用到目标敌人
			if target != nil && target.HP > 0 {
				m.buffManager.ApplyEnemyDebuff(target.ID, "mortal_strike", "致死打击", "debuff", duration, healingReduction, "healing_received", "")
			}
		}
	}
}

// handlePassiveOnDamageTakenEffects 处理被动技能的受到伤害时效果
func (m *BattleManager) handlePassiveOnDamageTakenEffects(character *models.Character, damageTaken int, session *BattleSession, logs *[]models.BattleLog) {
	if m.passiveSkillManager == nil {
		return
	}

	passives := m.passiveSkillManager.GetPassiveSkills(character.ID)
	for _, passive := range passives {
		switch passive.Passive.EffectType {
		case "on_damage_taken_resource":
			// 受到伤害时获得资源（如战士的受击回怒）
			// 这个已经在其他地方处理，这里可以添加其他被动技能
			resourceGain := int(passive.EffectValue * float64(damageTaken) / 100.0)
			if resourceGain > 0 {
				character.Resource += resourceGain
				if character.Resource > character.MaxResource {
					character.Resource = character.MaxResource
				}
				resourceName := m.getResourceName(character.ResourceType)
				m.addLog(session, "resource", fmt.Sprintf("%s 受到伤害获得了 %d 点%s", character.Name, resourceGain, resourceName), "#8888ff")
				*logs = append(*logs, session.BattleLogs[len(session.BattleLogs)-1])
			}
		case "on_damage_taken_heal":
			// 受到伤害时恢复生命值（如吸血效果）
			healPercent := passive.EffectValue // 百分比值（基于受到的伤害）
			healAmount := int(float64(damageTaken) * healPercent / 100.0)
			if healAmount > 0 {
				character.HP += healAmount
				if character.HP > character.MaxHP {
					character.HP = character.MaxHP
				}
				m.addLog(session, "heal", fmt.Sprintf("%s 的伤害恢复效果恢复了 %d 点生命值", character.Name, healAmount), "#00ff00")
				*logs = append(*logs, session.BattleLogs[len(session.BattleLogs)-1])
			}
		}
	}
}

// handlePassiveOnSkillUseEffects 处理被动技能的使用技能时效果
func (m *BattleManager) handlePassiveOnSkillUseEffects(character *models.Character, skillID string, session *BattleSession, logs *[]models.BattleLog) {
	if m.passiveSkillManager == nil {
		return
	}

	passives := m.passiveSkillManager.GetPassiveSkills(character.ID)
	for _, passive := range passives {
		switch passive.Passive.EffectType {
		case "on_skill_use_resource":
			// 使用技能时获得资源
			resourceGain := int(passive.EffectValue)
			if resourceGain > 0 {
				character.Resource += resourceGain
				if character.Resource > character.MaxResource {
					character.Resource = character.MaxResource
				}
				resourceName := m.getResourceName(character.ResourceType)
				m.addLog(session, "resource", fmt.Sprintf("%s 使用技能获得了 %d 点%s", character.Name, resourceGain, resourceName), "#8888ff")
				*logs = append(*logs, session.BattleLogs[len(session.BattleLogs)-1])
			}
		case "on_skill_use_heal":
			// 使用技能时恢复生命值
			healPercent := passive.EffectValue // 百分比值（基于最大HP）
			healAmount := int(float64(character.MaxHP) * healPercent / 100.0)
			if healAmount > 0 {
				character.HP += healAmount
				if character.HP > character.MaxHP {
					character.HP = character.MaxHP
				}
				m.addLog(session, "heal", fmt.Sprintf("%s 使用技能恢复了 %d 点生命值", character.Name, healAmount), "#00ff00")
				*logs = append(*logs, session.BattleLogs[len(session.BattleLogs)-1])
			}
		}
	}
}

// handlePassiveReflectEffects 处理被动技能的反射效果
func (m *BattleManager) handlePassiveReflectEffects(character *models.Character, attacker *models.Monster, damageTaken int, session *BattleSession, logs *[]models.BattleLog) {
	if m.passiveSkillManager == nil {
		return
	}

	passives := m.passiveSkillManager.GetPassiveSkills(character.ID)
	for _, passive := range passives {
		if passive.Passive.EffectType == "reflect" && passive.Passive.ID == "warrior_passive_shield_reflection" {
			// 盾牌反射（被动）：受到物理攻击时反射伤害
			reflectPercent := passive.EffectValue // 百分比值（如10.0表示10%）
			reflectDamage := int(float64(damageTaken) * reflectPercent / 100.0)
			if reflectDamage > 0 {
				passiveReflectOldHP := attacker.HP
				attacker.HP -= reflectDamage
				if attacker.HP < 0 {
					attacker.HP = 0
				}
				// 更新威胁值（被动反射伤害也产生威胁）
				m.updateThreat(session, attacker.ID, character.ID, reflectDamage)
				passiveReflectHPChange := m.formatHPChange(attacker.Name, passiveReflectOldHP, attacker.HP, attacker.MaxHP)
				m.addLog(session, "combat", fmt.Sprintf("%s 的盾牌反射对 %s 造成 %d 点反射伤害%s", character.Name, attacker.Name, reflectDamage, passiveReflectHPChange), "#ff8800", withDamageType("magic"))
				*logs = append(*logs, session.BattleLogs[len(session.BattleLogs)-1])
			}
		}
	}
}

// ═══════════════════════════════════════════════════════════
// 战斗统计收集方法
// ═══════════════════════════════════════════════════════════

// initBattleStats 初始化本场战斗的统计收集器
func (m *BattleManager) initBattleStats(session *BattleSession, characters []*models.Character) {
	session.BattleStartTime = time.Now()
	session.CurrentBattleRound = 0
	session.CharacterStats = make(map[int]*CharacterBattleStatsCollector)
	session.SkillBreakdown = make(map[int]map[string]*SkillUsageStats)

	// 使用新的 BattleStatsCollector 初始化
	if m.battleStatsCollector != nil {
		characterIDs := make([]int, len(characters))
		for i, char := range characters {
			characterIDs[i] = char.ID
		}
		m.battleStatsCollector.InitializeBattle(characterIDs)
	}

	// 保留旧的统计收集器（向后兼容）
	for _, char := range characters {
		session.CharacterStats[char.ID] = &CharacterBattleStatsCollector{
			CharacterID: char.ID,
			TeamSlot:    char.TeamSlot,
		}
		session.SkillBreakdown[char.ID] = make(map[string]*SkillUsageStats)
	}
}

// getOrCreateCharacterStats 获取或创建角色统计收集器
func (m *BattleManager) getOrCreateCharacterStats(session *BattleSession, characterID int, teamSlot int) *CharacterBattleStatsCollector {
	if session.CharacterStats == nil {
		session.CharacterStats = make(map[int]*CharacterBattleStatsCollector)
	}
	if stats, exists := session.CharacterStats[characterID]; exists {
		return stats
	}
	stats := &CharacterBattleStatsCollector{
		CharacterID: characterID,
		TeamSlot:    teamSlot,
	}
	session.CharacterStats[characterID] = stats
	return stats
}

// recordDamageDealt 记录角色造成的伤害
func (m *BattleManager) recordDamageDealt(session *BattleSession, characterID int, teamSlot int, damage int, damageType string, isCrit bool) {
	stats := m.getOrCreateCharacterStats(session, characterID, teamSlot)
	stats.DamageDealt += damage

	// 按伤害类型分类记录
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
	default:
		stats.PhysicalDamage += damage // 默认为物理
	}

	// 暴击统计
	if isCrit {
		stats.CritCount++
		stats.CritDamage += damage
		if damage > stats.MaxCrit {
			stats.MaxCrit = damage
		}
	}
}

// recordDamageTaken 记录角色受到的伤害
func (m *BattleManager) recordDamageTaken(session *BattleSession, characterID int, teamSlot int, damage int, damageType string, blocked int, absorbed int) {
	stats := m.getOrCreateCharacterStats(session, characterID, teamSlot)
	stats.DamageTaken += damage
	stats.DamageBlocked += blocked
	stats.DamageAbsorbed += absorbed
	stats.HitCount++

	// 按伤害类型分类记录
	switch damageType {
	case "physical":
		stats.PhysicalTaken += damage
	case "magic", "fire", "frost", "shadow", "holy", "nature":
		stats.MagicTaken += damage
	default:
		stats.PhysicalTaken += damage
	}
}

// recordHealing 记录治疗
func (m *BattleManager) recordHealing(session *BattleSession, healerID int, healerSlot int, targetID int, targetSlot int, healing int, overhealing int, isSelfHeal bool, isHot bool) {
	// 记录治疗者的输出
	healerStats := m.getOrCreateCharacterStats(session, healerID, healerSlot)
	healerStats.HealingDone += healing
	healerStats.Overhealing += overhealing
	if isSelfHeal {
		healerStats.SelfHealing += healing
	}
	if isHot {
		healerStats.HotHealing += healing
	}

	// 记录目标的受到治疗（如果不是自我治疗）
	if targetID != healerID {
		targetStats := m.getOrCreateCharacterStats(session, targetID, targetSlot)
		targetStats.HealingReceived += healing
	}
}

// recordSkillUsage 记录技能使用
func (m *BattleManager) recordSkillUsage(session *BattleSession, characterID int, teamSlot int, skillID string, damage int, healing int, resourceCost int, isHit bool, isCrit bool) {
	stats := m.getOrCreateCharacterStats(session, characterID, teamSlot)
	stats.SkillUses++
	if isHit {
		stats.SkillHits++
	} else {
		stats.SkillMisses++
	}
	stats.ResourceUsed += resourceCost

	// 记录技能明细
	if session.SkillBreakdown == nil {
		session.SkillBreakdown = make(map[int]map[string]*SkillUsageStats)
	}
	if session.SkillBreakdown[characterID] == nil {
		session.SkillBreakdown[characterID] = make(map[string]*SkillUsageStats)
	}

	skillStats, exists := session.SkillBreakdown[characterID][skillID]
	if !exists {
		skillStats = &SkillUsageStats{SkillID: skillID}
		session.SkillBreakdown[characterID][skillID] = skillStats
	}

	skillStats.UseCount++
	if isHit {
		skillStats.HitCount++
	}
	if isCrit {
		skillStats.CritCount++
	}
	skillStats.TotalDamage += damage
	skillStats.TotalHealing += healing
	skillStats.ResourceCost += resourceCost
}

// recordResourceGenerated 记录资源获得
func (m *BattleManager) recordResourceGenerated(session *BattleSession, characterID int, teamSlot int, amount int) {
	stats := m.getOrCreateCharacterStats(session, characterID, teamSlot)
	stats.ResourceGenerated += amount
}

// recordDodge 记录闪避
func (m *BattleManager) recordDodge(session *BattleSession, characterID int, teamSlot int) {
	stats := m.getOrCreateCharacterStats(session, characterID, teamSlot)
	stats.DodgeCount++
}

// recordKill 记录击杀
func (m *BattleManager) recordKill(session *BattleSession, characterID int, teamSlot int) {
	stats := m.getOrCreateCharacterStats(session, characterID, teamSlot)
	stats.Kills++
}

// recordDeath 记录死亡
func (m *BattleManager) recordDeath(session *BattleSession, characterID int, teamSlot int) {
	stats := m.getOrCreateCharacterStats(session, characterID, teamSlot)
	stats.Deaths++
}

// recordCcApplied 记录施加控制
func (m *BattleManager) recordCcApplied(session *BattleSession, characterID int, teamSlot int) {
	stats := m.getOrCreateCharacterStats(session, characterID, teamSlot)
	stats.CcApplied++
}

// recordCcReceived 记录受到控制
func (m *BattleManager) recordCcReceived(session *BattleSession, characterID int, teamSlot int) {
	stats := m.getOrCreateCharacterStats(session, characterID, teamSlot)
	stats.CcReceived++
}

// incrementBattleRound 增加战斗回合数
func (m *BattleManager) incrementBattleRound(session *BattleSession) {
	session.CurrentBattleRound++
}

// processMonsterDrops 处理怪物掉落
func (m *BattleManager) processMonsterDrops(session *BattleSession, enemies []*models.Monster, logs *[]models.BattleLog, characters []*models.Character) {
	if m.monsterManager == nil || enemies == nil || len(enemies) == 0 {
		return
	}

	// 如果没有角色，无法分配物品
	if len(characters) == 0 {
		return
	}

	// 使用第一个角色作为掉落接收者（未来可以支持多角色分配）
	character := characters[0]

	// 获取区域掉落倍率（如果区域管理器可用）
	dropMultiplier := 1.0
	if m.zoneManager != nil && session.CurrentZone != nil {
		dropMultiplier = m.zoneManager.CalculateDropMultiplier(session.CurrentZone.ID)
	}

	// 遍历所有被击败的敌人，计算掉落
	for _, enemy := range enemies {
		if enemy == nil || enemy.HP > 0 {
			continue
		}

		// 计算掉落
		drops, err := m.monsterManager.CalculateDrops(enemy.ID, enemy.Type)
		if err != nil {
			fmt.Printf("[WARN] Failed to calculate drops for monster %s: %v\n", enemy.ID, err)
			continue
		}

		// 如果有掉落，分配物品并记录日志
		if len(drops) > 0 {
			dropMessages := make([]string, 0)
			for _, drop := range drops {
				// 检查物品类型
				itemData, err := m.gameRepo.GetItemByID(drop.ItemID)
				if err != nil {
					fmt.Printf("[WARN] Failed to get item data for %s: %v\n", drop.ItemID, err)
					// 如果不是装备，直接添加到背包
					if m.inventoryRepo != nil {
						m.inventoryRepo.AddItem(character.ID, drop.ItemID, drop.Quantity)
					}
					dropMessages = append(dropMessages, fmt.Sprintf("%s x%d", drop.ItemID, drop.Quantity))
					continue
				}

				itemType, _ := itemData["type"].(string)
				// 如果是装备，生成装备实例
				if itemType == "equipment" && m.equipmentManager != nil {
					// 确定装备品质（根据怪物类型）
					quality := m.determineEquipmentQuality(enemy.Type, dropMultiplier)

					// 生成装备实例
					_, err := m.equipmentManager.GenerateEquipment(drop.ItemID, quality, enemy.Level, character.UserID)
					if err != nil {
						fmt.Printf("[WARN] Failed to generate equipment %s: %v\n", drop.ItemID, err)
						// 如果生成失败，仍然尝试添加到背包
						if m.inventoryRepo != nil {
							m.inventoryRepo.AddItem(character.ID, drop.ItemID, drop.Quantity)
						}
						dropMessages = append(dropMessages, fmt.Sprintf("%s x%d", drop.ItemID, drop.Quantity))
					} else {
						// 装备生成成功，添加到背包
						if m.inventoryRepo != nil {
							// 将装备实例ID添加到背包（需要InventoryRepository支持装备实例）
							// 暂时使用ItemID
							m.inventoryRepo.AddItem(character.ID, drop.ItemID, 1)
						}
						qualityName := m.getQualityDisplayName(quality)
						dropMessages = append(dropMessages, fmt.Sprintf("<span style=\"color: %s\">%s</span> x1",
							m.getQualityColor(quality), qualityName))
					}
				} else {
					// 非装备物品，直接添加到背包
					if m.inventoryRepo != nil {
						err := m.inventoryRepo.AddItem(character.ID, drop.ItemID, drop.Quantity)
						if err != nil {
							fmt.Printf("[WARN] Failed to add item %s to inventory: %v\n", drop.ItemID, err)
						}
					}
					dropMessages = append(dropMessages, fmt.Sprintf("%s x%d", drop.ItemID, drop.Quantity))
				}
			}

			if len(dropMessages) > 0 {
				dropText := fmt.Sprintf("🎁 击败 <span style=\"color: #ff7777\">%s</span> 获得: %s",
					enemy.Name, strings.Join(dropMessages, ", "))
				m.addLog(session, "loot", dropText, "#4ecdc4")
				*logs = append(*logs, session.BattleLogs[len(session.BattleLogs)-1])
			}
		}
	}
}

// determineEquipmentQuality 根据怪物类型确定装备品质
func (m *BattleManager) determineEquipmentQuality(monsterType string, dropMultiplier float64) string {
	// 基础品质分布（根据文档）
	var qualityWeights map[string]float64

	switch monsterType {
	case "normal":
		// 普通怪物: 30%白, 35%绿, 25%蓝, 8%紫, 1.8%橙, 0.2%传说
		qualityWeights = map[string]float64{
			"common":    30.0,
			"uncommon":  35.0,
			"rare":      25.0,
			"epic":      8.0,
			"legendary": 1.8,
			"mythic":    0.2,
		}
	case "elite":
		// 精英怪物: 40%白, 35%绿, 20%蓝, 4%紫, 0.9%橙, 0.1%传说
		qualityWeights = map[string]float64{
			"common":    40.0,
			"uncommon":  35.0,
			"rare":      20.0,
			"epic":      4.0,
			"legendary": 0.9,
			"mythic":    0.1,
		}
	case "boss":
		// Boss: 10%白, 25%绿, 35%蓝, 20%紫, 8%橙, 2%传说
		qualityWeights = map[string]float64{
			"common":    10.0,
			"uncommon":  25.0,
			"rare":      35.0,
			"epic":      20.0,
			"legendary": 8.0,
			"mythic":    2.0,
		}
	default:
		// 默认使用普通怪物分布
		qualityWeights = map[string]float64{
			"common":    30.0,
			"uncommon":  35.0,
			"rare":      25.0,
			"epic":      8.0,
			"legendary": 1.8,
			"mythic":    0.2,
		}
	}

	// 应用区域掉落倍率（提升高品质装备概率）
	if dropMultiplier > 1.0 {
		// 提升高品质装备的权重
		qualityWeights["epic"] *= dropMultiplier
		qualityWeights["legendary"] *= dropMultiplier
		qualityWeights["mythic"] *= dropMultiplier
	}

	// 归一化权重
	totalWeight := 0.0
	for _, weight := range qualityWeights {
		totalWeight += weight
	}

	// 随机选择品质
	randValue := rand.Float64() * totalWeight
	currentWeight := 0.0

	qualityOrder := []string{"common", "uncommon", "rare", "epic", "legendary", "mythic"}
	for _, quality := range qualityOrder {
		currentWeight += qualityWeights[quality]
		if randValue <= currentWeight {
			return quality
		}
	}

	// 默认返回普通品质
	return "common"
}

// getQualityDisplayName 获取品质显示名称
func (m *BattleManager) getQualityDisplayName(quality string) string {
	names := map[string]string{
		"common":    "普通",
		"uncommon":  "优秀",
		"rare":      "精良",
		"epic":      "稀有",
		"legendary": "史诗",
		"mythic":    "传说",
	}
	if name, ok := names[quality]; ok {
		return name
	}
	return "普通"
}

// getQualityColor 获取品质颜色
func (m *BattleManager) getQualityColor(quality string) string {
	colors := map[string]string{
		"common":    "#ffffff", // 白色
		"uncommon":  "#1eff00", // 绿色
		"rare":      "#0070dd", // 蓝色
		"epic":      "#a335ee", // 紫色
		"legendary": "#ff8000", // 橙色
		"mythic":    "#ffd700", // 金色
	}
	if color, ok := colors[quality]; ok {
		return color
	}
	return "#ffffff"
}

// saveBattleStats 保存战斗统计到数据库
func (m *BattleManager) saveBattleStats(session *BattleSession, userID int, zoneID string, monsterID string, isVictory bool, characters []*models.Character) {
	if m.battleStatsRepo == nil {
		return
	}

	// 如果没有统计数据，跳过保存
	if session.CharacterStats == nil || len(session.CharacterStats) == 0 {
		return
	}

	// 计算战斗时长
	duration := int(time.Since(session.BattleStartTime).Seconds())

	// 计算团队总伤害和治疗
	var teamDamageDealt, teamDamageTaken, teamHealingDone int
	for _, stats := range session.CharacterStats {
		teamDamageDealt += stats.DamageDealt
		teamDamageTaken += stats.DamageTaken
		teamHealingDone += stats.HealingDone
	}

	// 创建战斗记录
	result := "victory"
	if !isVictory {
		result = "defeat"
	}

	battleRecord := &models.BattleRecord{
		UserID:          userID,
		ZoneID:          zoneID,
		BattleType:      "pve",
		MonsterID:       monsterID,
		TotalRounds:     session.CurrentBattleRound,
		DurationSeconds: duration,
		Result:          result,
		TeamDamageDealt: teamDamageDealt,
		TeamDamageTaken: teamDamageTaken,
		TeamHealingDone: teamHealingDone,
		ExpGained:       session.CurrentBattleExp,
		GoldGained:      session.CurrentBattleGold,
	}

	// 保存战斗记录
	battleID, err := m.battleStatsRepo.CreateBattleRecord(battleRecord)
	if err != nil {
		fmt.Printf("[ERROR] Failed to save battle record: %v\n", err)
		return
	}

	// 保存每个角色的统计数据
	today := time.Now().Format("2006-01-02")
	for characterID, collector := range session.CharacterStats {
		// 创建角色战斗统计
		charStats := &models.BattleCharacterStats{
			BattleID:          int(battleID),
			CharacterID:       characterID,
			TeamSlot:          collector.TeamSlot,
			DamageDealt:       collector.DamageDealt,
			PhysicalDamage:    collector.PhysicalDamage,
			MagicDamage:       collector.MagicDamage,
			FireDamage:        collector.FireDamage,
			FrostDamage:       collector.FrostDamage,
			ShadowDamage:      collector.ShadowDamage,
			HolyDamage:        collector.HolyDamage,
			NatureDamage:      collector.NatureDamage,
			DotDamage:         collector.DotDamage,
			CritCount:         collector.CritCount,
			CritDamage:        collector.CritDamage,
			MaxCrit:           collector.MaxCrit,
			DamageTaken:       collector.DamageTaken,
			PhysicalTaken:     collector.PhysicalTaken,
			MagicTaken:        collector.MagicTaken,
			DamageBlocked:     collector.DamageBlocked,
			DamageAbsorbed:    collector.DamageAbsorbed,
			DodgeCount:        collector.DodgeCount,
			BlockCount:        collector.BlockCount,
			HitCount:          collector.HitCount,
			HealingDone:       collector.HealingDone,
			HealingReceived:   collector.HealingReceived,
			Overhealing:       collector.Overhealing,
			SelfHealing:       collector.SelfHealing,
			HotHealing:        collector.HotHealing,
			SkillUses:         collector.SkillUses,
			SkillHits:         collector.SkillHits,
			SkillMisses:       collector.SkillMisses,
			CcApplied:         collector.CcApplied,
			CcReceived:        collector.CcReceived,
			Dispels:           collector.Dispels,
			Interrupts:        collector.Interrupts,
			Kills:             collector.Kills,
			Deaths:            collector.Deaths,
			Resurrects:        collector.Resurrects,
			ResourceUsed:      collector.ResourceUsed,
			ResourceGenerated: collector.ResourceGenerated,
		}

		_, err := m.battleStatsRepo.CreateBattleCharacterStats(charStats)
		if err != nil {
			fmt.Printf("[ERROR] Failed to save character battle stats: %v\n", err)
		}

		// 更新角色生涯统计
		err = m.battleStatsRepo.UpdateLifetimeStats(characterID, charStats, isVictory, "pve", session.CurrentBattleRound)
		if err != nil {
			fmt.Printf("[ERROR] Failed to update lifetime stats: %v\n", err)
		}

		// 保存技能明细
		if skillBreakdown, exists := session.SkillBreakdown[characterID]; exists {
			for skillID, skillStats := range skillBreakdown {
				breakdown := &models.BattleSkillBreakdown{
					BattleID:     int(battleID),
					CharacterID:  characterID,
					SkillID:      skillID,
					UseCount:     skillStats.UseCount,
					HitCount:     skillStats.HitCount,
					CritCount:    skillStats.CritCount,
					TotalDamage:  skillStats.TotalDamage,
					TotalHealing: skillStats.TotalHealing,
					ResourceCost: skillStats.ResourceCost,
				}
				_, err := m.battleStatsRepo.CreateBattleSkillBreakdown(breakdown)
				if err != nil {
					fmt.Printf("[ERROR] Failed to save skill breakdown: %v\n", err)
				}
			}
		}
	}

	// 更新每日统计
	err = m.battleStatsRepo.UpdateDailyStats(
		userID, today, isVictory,
		teamDamageDealt, teamHealingDone, teamDamageTaken,
		session.CurrentBattleExp, session.CurrentBattleGold,
		session.CurrentBattleKills, 0, // deaths 需要从角色统计中计算
	)
	if err != nil {
		fmt.Printf("[ERROR] Failed to update daily stats: %v\n", err)
	}
}

// clearBattleStats 清除本场战斗的统计数据
func (m *BattleManager) clearBattleStats(session *BattleSession) {
	session.CharacterStats = nil
	session.SkillBreakdown = nil
	session.CurrentBattleRound = 0
}

// ═══════════════════════════════════════════════════════════
// 用户自定义统计会话管理
// ═══════════════════════════════════════════════════════════

// StartStatsSession 开始统计会话
func (m *BattleManager) StartStatsSession(userID int) {
	m.statsSessionsMu.Lock()
	defer m.statsSessionsMu.Unlock()

	m.statsSessions[userID] = &StatsSession{
		UserID:    userID,
		StartTime: time.Now(),
		IsActive:  true,
	}
}

// ResetStatsSession 重置统计会话
func (m *BattleManager) ResetStatsSession(userID int) {
	m.statsSessionsMu.Lock()
	defer m.statsSessionsMu.Unlock()

	delete(m.statsSessions, userID)
}

// GetStatsSession 获取统计会话
func (m *BattleManager) GetStatsSession(userID int) *StatsSession {
	m.statsSessionsMu.RLock()
	defer m.statsSessionsMu.RUnlock()

	return m.statsSessions[userID]
}

// updateThreat 更新威胁值
// 当角色对怪物造成伤害时，增加该角色对该怪物的威胁值
// threatGain: 威胁值增加量（通常等于伤害值，但可以根据技能类型调整）
func (m *BattleManager) updateThreat(session *BattleSession, monsterID string, characterID int, threatGain int) {
	if session == nil || session.ThreatTable == nil {
		return
	}

	// 初始化该怪物的威胁表（如果不存在）
	if session.ThreatTable[monsterID] == nil {
		session.ThreatTable[monsterID] = make(map[int]int)
	}

	// 增加威胁值
	session.ThreatTable[monsterID][characterID] += threatGain

	// 确保威胁值不为负数
	if session.ThreatTable[monsterID][characterID] < 0 {
		session.ThreatTable[monsterID][characterID] = 0
	}
}

// getThreatTableForMonster 获取特定怪物的威胁表
func (m *BattleManager) getThreatTableForMonster(session *BattleSession, monsterID string) map[int]int {
	if session == nil || session.ThreatTable == nil {
		return make(map[int]int)
	}

	if threatTable, exists := session.ThreatTable[monsterID]; exists {
		return threatTable
	}

	return make(map[int]int)
}

// resetThreatTable 重置威胁表（新战斗开始时调用）
func (m *BattleManager) resetThreatTable(session *BattleSession) {
	if session == nil {
		return
	}

	// 清空所有威胁表
	session.ThreatTable = make(map[string]map[int]int)
}

// buildTurnOrder 构建回合顺序队列（按速度排序）
// 包含所有角色和敌人，按速度从高到低排序
func (m *BattleManager) buildTurnOrder(session *BattleSession, characters []*models.Character, enemies []*models.Monster) {
	if session == nil {
		return
	}

	turnOrder := make([]*TurnParticipant, 0)

	// 添加所有角色到队列
	for i, char := range characters {
		if char == nil || char.HP <= 0 {
			continue
		}
		speed := m.calculator.CalculateSpeed(char)
		turnOrder = append(turnOrder, &TurnParticipant{
			Type:      "character",
			Character: char,
			Speed:     speed,
			Index:     i,
		})
	}

	// 添加所有敌人到队列
	for i, enemy := range enemies {
		if enemy == nil || enemy.HP <= 0 {
			continue
		}
		speed := enemy.Speed
		if speed <= 0 {
			speed = 10 // 默认速度
		}
		turnOrder = append(turnOrder, &TurnParticipant{
			Type:    "monster",
			Monster: enemy,
			Speed:   speed,
			Index:   i,
		})
	}

	// 按速度从高到低排序（速度相同则随机）
	sort.Slice(turnOrder, func(i, j int) bool {
		if turnOrder[i].Speed != turnOrder[j].Speed {
			return turnOrder[i].Speed > turnOrder[j].Speed
		}
		// 速度相同时，随机排序（使用索引作为随机种子）
		return rand.Intn(2) == 0
	})

	session.TurnOrder = turnOrder
	session.CurrentTurnOrderIndex = 0
}

// getCurrentTurnParticipant 获取当前回合的参与者
func (m *BattleManager) getCurrentTurnParticipant(session *BattleSession) *TurnParticipant {
	if session == nil || session.TurnOrder == nil || len(session.TurnOrder) == 0 {
		return nil
	}
	if session.CurrentTurnOrderIndex < 0 || session.CurrentTurnOrderIndex >= len(session.TurnOrder) {
		return nil
	}
	return session.TurnOrder[session.CurrentTurnOrderIndex]
}

// moveToNextTurn 移动到下一个回合
func (m *BattleManager) moveToNextTurn(session *BattleSession, characters []*models.Character, enemies []*models.Monster) {
	if session == nil {
		return
	}

	// 移动到下一个参与者
	session.CurrentTurnOrderIndex++

	// 如果所有参与者都行动完毕，开始新的一轮
	if session.CurrentTurnOrderIndex >= len(session.TurnOrder) {
		// 重新构建回合队列（因为可能有角色/敌人死亡，速度可能变化）
		m.buildTurnOrder(session, characters, enemies)
		// 增加回合数
		session.CurrentBattleRound++
		m.incrementBattleRound(session)

		// 添加回合开始日志（每5回合显示一次，避免日志过多）
		// 注意：这个日志不会在moveToNextTurn中直接添加，而是在需要时由调用者添加
		// 避免在每次移动回合时都产生日志
	}
}

// checkBattleEnd 检查战斗是否结束
// 返回: (isEnded, isVictory, allCharactersDead)
// isEnded: 战斗是否结束
// isVictory: 是否胜利（仅当isEnded为true时有效）
// allCharactersDead: 所有角色是否都死亡（仅当isEnded为true时有效）
func (m *BattleManager) checkBattleEnd(session *BattleSession, characters []*models.Character, enemies []*models.Monster) (isEnded bool, isVictory bool, allCharactersDead bool) {
	if session == nil {
		return false, false, false
	}

	// 检查所有角色是否都死亡
	allDead := true
	hasAliveCharacter := false
	for _, char := range characters {
		if char != nil && char.HP > 0 && !char.IsDead {
			allDead = false
			hasAliveCharacter = true
			break
		}
	}

	// 检查所有敌人是否都被击败
	allEnemiesDefeated := true
	hasAliveEnemy := false
	for _, enemy := range enemies {
		if enemy != nil && enemy.HP > 0 {
			allEnemiesDefeated = false
			hasAliveEnemy = true
			break
		}
	}

	// 如果所有角色都死亡，战斗失败
	if allDead && !hasAliveCharacter {
		return true, false, true
	}

	// 如果所有敌人都被击败，战斗胜利
	if allEnemiesDefeated && !hasAliveEnemy && len(enemies) > 0 {
		return true, true, false
	}

	// 战斗继续
	return false, false, false
}

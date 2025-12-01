package game

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"text-wow/internal/models"
	"text-wow/internal/repository"
)

// BattleManager 战斗管理器 - 管理所有用户的战斗状态
type BattleManager struct {
	mu       sync.RWMutex
	sessions map[int]*BattleSession // key: userID
	gameRepo *repository.GameRepository
	charRepo *repository.CharacterRepository
}

// BattleSession 用户战斗会话
type BattleSession struct {
	UserID            int
	IsRunning         bool
	CurrentZone       *models.Zone
	CurrentEnemy      *models.Monster      // 保留用于向后兼容
	CurrentEnemies    []*models.Monster    // 多个敌人支持
	BattleLogs       []models.BattleLog
	BattleCount       int
	SessionKills      int
	SessionGold       int
	SessionExp        int
	StartedAt         time.Time
	LastTick          time.Time
	IsResting         bool       // 是否在休息
	RestUntil         *time.Time // 休息结束时间
	RestStartedAt     *time.Time // 休息开始时间
	LastRestTick      *time.Time // 上次恢复处理的时间
	RestSpeed         float64    // 恢复速度倍率
	CurrentBattleExp  int        // 本场战斗获得的经验
	CurrentBattleGold int        // 本场战斗获得的金币
	CurrentBattleKills int       // 本场战斗击杀数
	CurrentTurnIndex  int        // 回合控制：-1=玩家回合，>=0=敌人索引
	JustEncountered   bool       // 刚遭遇敌人，需要等待1个tick再开始战斗
}

// NewBattleManager 创建战斗管理器
func NewBattleManager() *BattleManager {
	return &BattleManager{
		sessions: make(map[int]*BattleSession),
		gameRepo: repository.NewGameRepository(),
		charRepo: repository.NewCharacterRepository(),
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
		return session
	}

	session := &BattleSession{
		UserID:           userID,
		BattleLogs:       make([]models.BattleLog, 0),
		StartedAt:        time.Now(),
		CurrentEnemies:    make([]*models.Monster, 0),
		CurrentTurnIndex:  -1, // 初始化为玩家回合
		RestSpeed:        1.0,
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

	if !session.IsRunning || len(characters) == 0 {
		return nil, nil
	}

	session.LastTick = time.Now()
	logs := make([]models.BattleLog, 0)

	// 使用第一个角色进行战斗
	char := characters[0]
	
	// 确保战士的怒气上限为100（每次tick都检查，防止被覆盖）
	if char.ResourceType == "rage" {
		char.MaxResource = 100
	}
	
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
			session.IsRunning = false
			
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
			// 休息结束
			m.addLog(session, "system", ">> 休息结束，准备下一场战斗", "#33ff33")
			logs = append(logs, session.BattleLogs[len(session.BattleLogs)-1])
		} else {
			// 仍在休息中
			remaining := session.RestUntil.Sub(time.Now())
			if remaining > 0 {
				m.addLog(session, "system", fmt.Sprintf(">> 休息中... (剩余 %d 秒)", int(remaining.Seconds())+1), "#888888")
				logs = append(logs, session.BattleLogs[len(session.BattleLogs)-1])
			}
		}
		
		// 保存角色数据更新
		if char.HP != initialHP || char.Resource != initialMP {
			m.charRepo.UpdateAfterBattle(char.ID, char.HP, char.Resource, char.Exp, char.Level,
				char.ExpToNext, char.MaxHP, char.MaxResource, char.Attack, char.Defense,
				char.Strength, char.Agility, char.Stamina, char.TotalKills)
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
	for _, enemy := range session.CurrentEnemies {
		if enemy != nil && enemy.HP > 0 {
			aliveEnemies = append(aliveEnemies, enemy)
		}
	}

	// 如果没有敌人，生成新的
	if len(aliveEnemies) == 0 {
		// 重置本场战斗统计
		session.CurrentBattleExp = 0
		session.CurrentBattleGold = 0
		session.CurrentBattleKills = 0
		session.CurrentTurnIndex = -1 // 玩家回合
		
		// 战斗开始时，确保战士的怒气为0，最大怒气为100
		if char.ResourceType == "rage" {
			char.Resource = 0
			char.MaxResource = 100
		}
		
		err := m.spawnEnemies(session, char.Level)
		if err != nil {
			return nil, err
		}
		logs = append(logs, session.BattleLogs[len(session.BattleLogs)-1])
		
		// 标记刚遭遇敌人，需要等待1个tick再开始战斗
		session.JustEncountered = true
		
		// 更新存活敌人列表
		aliveEnemies = session.CurrentEnemies
		
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

	// 回合制逻辑：CurrentTurnIndex == -1 表示玩家回合，>=0 表示敌人索引
	if session.CurrentTurnIndex == -1 {
		// 玩家回合：攻击第一个存活的敌人
		if len(aliveEnemies) > 0 {
			target := aliveEnemies[0]
			
			// 确定使用的技能和消耗
			skillName, skillCost := m.getSkillForAttack(char)
			
			// 如果是战士，检查怒气是否足够使用技能
			if char.ResourceType == "rage" {
				if skillCost > 0 && char.Resource < skillCost {
					// 怒气不足，只能使用普通攻击
					skillName = "普通攻击"
					skillCost = 0
				}
			}
			
			playerDamage := m.calculateDamage(char.Attack, target.Defense)
			isCrit := rand.Float64() < char.CritRate
			if isCrit {
				playerDamage = int(float64(playerDamage) * char.CritDamage)
			}
			target.HP -= playerDamage
			
			// 消耗资源（如果是战士，消耗怒气）
			if char.ResourceType == "rage" && skillCost > 0 {
				char.Resource -= skillCost
				if char.Resource < 0 {
					char.Resource = 0
				}
			}
			
			// 战士攻击获得怒气
			if char.ResourceType == "rage" {
				if isCrit {
					char.Resource += 10 // 暴击获得10点怒气
				} else {
					char.Resource += 5 // 普通攻击获得5点怒气
				}
				// 确保不超过最大值
				if char.Resource > char.MaxResource {
					char.Resource = char.MaxResource
				}
			}

			if isCrit {
				m.addLog(session, "combat", fmt.Sprintf("%s 使用 [%s] 💥暴击！对 %s 造成 %d 点伤害", char.Name, skillName, target.Name, playerDamage), "#ff6b6b")
			} else {
				m.addLog(session, "combat", fmt.Sprintf("%s 使用 [%s] 对 %s 造成 %d 点伤害", char.Name, skillName, target.Name, playerDamage), "#ffaa00")
			}
			logs = append(logs, session.BattleLogs[len(session.BattleLogs)-1])

			// 检查目标是否死亡
			if target.HP <= 0 {
				// 确保HP归零
				target.HP = 0
				// 敌人死亡
				expGain := target.ExpReward
				goldGain := target.GoldMin + rand.Intn(target.GoldMax-target.GoldMin+1)

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

					// 升级属性提升
					char.MaxHP += 15
					char.HP = char.MaxHP
					
					// 战士的怒气最大值固定为100，不随升级改变
					if char.ResourceType == "rage" {
						char.MaxResource = 100
						// 升级时怒气保持不变，不重置为最大值
					} else {
						char.MaxResource += 8
						char.Resource = char.MaxResource
					}
					
					char.Strength += 2
					char.Agility += 1
					char.Stamina += 2
					char.Attack = char.Strength / 2
					char.Defense = char.Stamina / 3

					m.addLog(session, "levelup", fmt.Sprintf("🎉【升级】恭喜！%s 升到了 %d 级！", char.Name, char.Level), "#ffd700")
					logs = append(logs, session.BattleLogs[len(session.BattleLogs)-1])
				}
			}

			// 移动到下一个敌人回合
			session.CurrentTurnIndex = 0
		}
	} else {
		// 敌人回合：当前索引的敌人攻击玩家
		if session.CurrentTurnIndex < len(aliveEnemies) {
			enemy := aliveEnemies[session.CurrentTurnIndex]
			enemyDamage := m.calculateDamage(enemy.Attack, char.Defense)
			char.HP -= enemyDamage
			
			// 战士受到伤害时获得怒气
			if char.ResourceType == "rage" && enemyDamage > 0 {
				// 受到伤害获得怒气: 伤害/最大HP × 50，至少1点
				rageGain := int(float64(enemyDamage) / float64(char.MaxHP) * 50)
				if rageGain < 1 {
					rageGain = 1
				}
				char.Resource += rageGain
				if char.Resource > char.MaxResource {
					char.Resource = char.MaxResource
				}
			}

			m.addLog(session, "combat", fmt.Sprintf("%s 攻击了 %s，造成 %d 点伤害", enemy.Name, char.Name, enemyDamage), "#ff4444")
			logs = append(logs, session.BattleLogs[len(session.BattleLogs)-1])

			// 检查玩家是否死亡
			if char.HP <= 0 {
				char.TotalDeaths++
				session.IsRunning = false
				session.CurrentEnemies = nil
				session.CurrentEnemy = nil
				session.CurrentTurnIndex = -1

				// 计算复活时间
				reviveDuration := m.calculateReviveTime(userID)
				now := time.Now()
				reviveAt := now.Add(reviveDuration)
				
				// 设置角色HP为0（死亡状态）
				char.HP = 0
				char.IsDead = true
				char.ReviveAt = &reviveAt

				m.addLog(session, "death", fmt.Sprintf("%s 被击败了... 需要 %d 秒复活", char.Name, int(reviveDuration.Seconds())), "#ff0000")
				logs = append(logs, session.BattleLogs[len(session.BattleLogs)-1])

				// 保存死亡数据（包括死亡标记和复活时间）
				m.charRepo.UpdateAfterDeath(char.ID, char.HP, char.TotalDeaths, &reviveAt)
				
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
			} else {
				// 移动到下一个敌人或回到玩家回合
				session.CurrentTurnIndex++
				if session.CurrentTurnIndex >= len(aliveEnemies) {
					session.CurrentTurnIndex = -1 // 回到玩家回合
				}
			}
		} else {
			// 索引超出范围，回到玩家回合
			session.CurrentTurnIndex = -1
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
		
		// 战斗胜利总结
		if session.CurrentBattleKills > 0 {
			summaryMsg := fmt.Sprintf("━━━ 战斗总结 ━━━ 击杀: %d | 经验: %d | 金币: %d", 
				session.CurrentBattleKills, session.CurrentBattleExp, session.CurrentBattleGold)
			m.addLog(session, "battle_summary", summaryMsg, "#ffd700")
			logs = append(logs, session.BattleLogs[len(session.BattleLogs)-1])
		}

		// 添加分割线
		m.addLog(session, "battle_separator", "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━", "#666666")
		logs = append(logs, session.BattleLogs[len(session.BattleLogs)-1])

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
		char.ExpToNext, char.MaxHP, char.MaxResource, char.Attack, char.Defense,
		char.Strength, char.Agility, char.Stamina, char.TotalKills)

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

// spawnEnemy 生成敌人（向后兼容）
func (m *BattleManager) spawnEnemy(session *BattleSession, playerLevel int) error {
	return m.spawnEnemies(session, playerLevel)
}

// spawnEnemies 生成多个敌人（1-3个随机）
func (m *BattleManager) spawnEnemies(session *BattleSession, playerLevel int) error {
	if session.CurrentZone == nil {
		// 加载默认区域
		zone, err := m.gameRepo.GetZoneByID("elwynn")
		if err != nil {
			fmt.Printf("[ERROR] Failed to get zone: %v\n", err)
			return fmt.Errorf("failed to get zone: %v", err)
		}
		session.CurrentZone = zone
		fmt.Printf("[DEBUG] Loaded zone: %s\n", zone.Name)
	}

	// 获取区域怪物
	monsters, err := m.gameRepo.GetMonstersByZone(session.CurrentZone.ID)
	if err != nil {
		fmt.Printf("[ERROR] Failed to get monsters: %v\n", err)
		return fmt.Errorf("failed to get monsters: %v", err)
	}
	if len(monsters) == 0 {
		fmt.Printf("[ERROR] No monsters in zone %s\n", session.CurrentZone.ID)
		return fmt.Errorf("no monsters in zone %s", session.CurrentZone.ID)
	}
	fmt.Printf("[DEBUG] Found %d monsters in zone\n", len(monsters))

	// 随机生成1-3个敌人
	enemyCount := 1 + rand.Intn(3) // 1-3个
	session.CurrentEnemies = make([]*models.Monster, 0, enemyCount)
	
	var enemyNames []string
	for i := 0; i < enemyCount; i++ {
		// 随机选择一个怪物模板
		template := monsters[rand.Intn(len(monsters))]
		
		enemy := &models.Monster{
			ID:        template.ID,
			ZoneID:    template.ZoneID,
			Name:      template.Name,
			Level:     template.Level,
			Type:      template.Type,
			HP:        template.HP,
			MaxHP:     template.HP,
			Attack:    template.Attack,
			Defense:   template.Defense,
			ExpReward: template.ExpReward,
			GoldMin:   template.GoldMin,
			GoldMax:   template.GoldMax,
		}
		session.CurrentEnemies = append(session.CurrentEnemies, enemy)
		enemyNames = append(enemyNames, fmt.Sprintf("%s (Lv.%d)", enemy.Name, enemy.Level))
	}

	// 保留 CurrentEnemy 用于向后兼容（指向第一个敌人）
	if len(session.CurrentEnemies) > 0 {
		session.CurrentEnemy = session.CurrentEnemies[0]
	}

	session.BattleCount++
	enemyList := fmt.Sprintf("%s", enemyNames[0])
	if len(enemyNames) > 1 {
		enemyList = fmt.Sprintf("%s 等 %d 个敌人", enemyNames[0], len(enemyNames))
	}
	m.addLog(session, "encounter", fmt.Sprintf("━━━ 战斗 #%d ━━━ 遭遇: %s", session.BattleCount, enemyList), "#ffff00")

	return nil
}

// ChangeZone 切换区域
func (m *BattleManager) ChangeZone(userID int, zoneID string, playerLevel int) error {
	session := m.GetOrCreateSession(userID)

	zone, err := m.gameRepo.GetZoneByID(zoneID)
	if err != nil {
		return fmt.Errorf("zone not found: %s", zoneID)
	}

	if playerLevel < zone.MinLevel {
		return fmt.Errorf("level too low, need level %d", zone.MinLevel)
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	session.CurrentZone = zone
	session.CurrentEnemy = nil

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

// calculateDamage 计算伤害
func (m *BattleManager) calculateDamage(attack, defense int) int {
	baseDamage := attack - defense/2
	if baseDamage < 1 {
		baseDamage = 1
	}
	// 添加随机波动 ±20%
	variance := float64(baseDamage) * 0.2
	damage := float64(baseDamage) + (rand.Float64()*2-1)*variance
	return int(damage)
}

// addLog 添加日志
func (m *BattleManager) addLog(session *BattleSession, logType, message, color string) {
	log := models.BattleLog{
		Message:   message,
		LogType:   logType,
		CreatedAt: time.Now(),
	}
	session.BattleLogs = append(session.BattleLogs, log)

	// 保持日志数量在合理范围
	if len(session.BattleLogs) > 200 {
		session.BattleLogs = session.BattleLogs[len(session.BattleLogs)-200:]
	}
}

// getRandomSkillName 获取随机技能名称
func (m *BattleManager) getRandomSkillName(classID string) string {
	skills := map[string][]string{
		"warrior": {"英勇打击", "雷霆一击", "顺劈斩", "致死打击"},
		"paladin": {"圣光术", "十字军打击", "正义之锤", "审判"},
		"hunter":  {"奥术射击", "多重射击", "瞄准射击", "稳固射击"},
		"rogue":   {"邪恶攻击", "剔骨", "背刺", "毒刃"},
		"priest":  {"惩击", "暗言术:痛", "神圣之火", "心灵震爆"},
		"mage":    {"火球术", "寒冰箭", "奥术飞弹", "炎爆术"},
		"warlock": {"暗影箭", "腐蚀术", "献祭", "混乱箭"},
		"druid":   {"月火术", "愤怒", "挥击", "横扫"},
		"shaman":  {"闪电箭", "闪电链", "熔岩爆裂", "烈焰震击"},
	}

	if classSkills, ok := skills[classID]; ok {
		return classSkills[rand.Intn(len(classSkills))]
	}
	return "普通攻击"
}

// getSkillForAttack 获取攻击技能名称和消耗
func (m *BattleManager) getSkillForAttack(char *models.Character) (string, int) {
	// 战士技能及其怒气消耗
	warriorSkills := []struct {
		name string
		cost int
	}{
		{"英勇打击", 10},
		{"雷霆一击", 15},
		{"顺劈斩", 12},
		{"致死打击", 20},
	}
	
	// 如果是战士，返回随机技能和消耗
	if char.ResourceType == "rage" {
		skill := warriorSkills[rand.Intn(len(warriorSkills))]
		return skill.name, skill.cost
	}
	
	// 其他职业使用普通技能，不消耗资源（或消耗法力，但这里简化处理）
	skillName := m.getRandomSkillName(char.ClassID)
	return skillName, 0
}

// calculateReviveTime 计算复活时间（根据死亡人数）
func (m *BattleManager) calculateReviveTime(userID int) time.Duration {
	deadCount, err := m.charRepo.CountDeadByUserID(userID)
	if err != nil {
		deadCount = 1 // 默认值
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
func (m *BattleManager) calculateRestTime(char *models.Character) time.Duration {
	hpLoss := float64(char.MaxHP - char.HP)
	mpLoss := float64(char.MaxResource - char.Resource)
	
	// 如果已经满血满蓝，不需要休息
	if hpLoss <= 0 && mpLoss <= 0 {
		return 0
	}
	
	// 基础休息时间：每损失1% HP/MP = 0.1秒，最少1秒
	// 使用HP和MP损失的平均值
	hpLossPercent := hpLoss / float64(char.MaxHP)
	mpLossPercent := mpLoss / float64(char.MaxResource)
	totalLoss := (hpLossPercent + mpLossPercent) / 2.0
	
	// 每秒恢复2%，所以需要的时间 = 损失百分比 / 0.02
	// 但为了更合理，我们使用：每损失1% = 0.5秒（因为每秒恢复2%，所以50%损失需要25秒）
	restSeconds := totalLoss * 50.0
	
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
			// 复活时间到了，恢复角色到一半HP
			char.HP = char.MaxHP / 2
			char.IsDead = false
			char.ReviveAt = nil
			
			// 更新数据库，清除死亡标记
			m.charRepo.SetDead(char.ID, false, nil)
			
			// 更新角色HP
			m.charRepo.UpdateAfterBattle(char.ID, char.HP, char.Resource, char.Exp, char.Level,
				char.ExpToNext, char.MaxHP, char.MaxResource, char.Attack, char.Defense,
				char.Strength, char.Agility, char.Stamina, char.TotalKills)
			
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
	Character    *models.Character `json:"character"`
	Enemy        *models.Monster   `json:"enemy,omitempty"`
	Enemies      []*models.Monster `json:"enemies,omitempty"` // 多个敌人支持
	Logs         []models.BattleLog `json:"logs"`
	IsRunning    bool              `json:"isRunning"`
	IsResting    bool              `json:"isResting"`    // 是否在休息
	RestUntil    *time.Time        `json:"restUntil,omitempty"` // 休息结束时间
	SessionKills int               `json:"sessionKills"`
	SessionGold  int               `json:"sessionGold"`
	SessionExp   int               `json:"sessionExp"`
	BattleCount  int               `json:"battleCount"`
}


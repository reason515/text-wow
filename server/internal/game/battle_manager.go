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
	RestSpeed         float64    // 恢复速度倍率
	CurrentBattleExp  int        // 本场战斗获得的经验
	CurrentBattleGold int        // 本场战斗获得的金币
	CurrentBattleKills int       // 本场战斗击杀数
	CurrentTurnIndex  int        // 回合控制：-1=玩家回合，>=0=敌人索引
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

	// 如果正在休息，处理休息
	if session.IsResting && session.RestUntil != nil {
		initialHP := char.HP
		initialMP := char.Resource
		m.processRest(session, char)
		
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
		
		err := m.spawnEnemies(session, char.Level)
		if err != nil {
			return nil, err
		}
		logs = append(logs, session.BattleLogs[len(session.BattleLogs)-1])
		
		// 更新存活敌人列表
		aliveEnemies = session.CurrentEnemies
	}

	// 回合制逻辑：CurrentTurnIndex == -1 表示玩家回合，>=0 表示敌人索引
	if session.CurrentTurnIndex == -1 {
		// 玩家回合：攻击第一个存活的敌人
		if len(aliveEnemies) > 0 {
			target := aliveEnemies[0]
			playerDamage := m.calculateDamage(char.Attack, target.Defense)
			isCrit := rand.Float64() < char.CritRate
			if isCrit {
				playerDamage = int(float64(playerDamage) * char.CritDamage)
			}
			target.HP -= playerDamage

			skillName := m.getRandomSkillName(char.ClassID)
			if isCrit {
				m.addLog(session, "combat", fmt.Sprintf("%s 使用 [%s] 💥暴击！对 %s 造成 %d 点伤害", char.Name, skillName, target.Name, playerDamage), "#ff6b6b")
			} else {
				m.addLog(session, "combat", fmt.Sprintf("%s 使用 [%s] 对 %s 造成 %d 点伤害", char.Name, skillName, target.Name, playerDamage), "#ffaa00")
			}
			logs = append(logs, session.BattleLogs[len(session.BattleLogs)-1])

			// 检查目标是否死亡
			if target.HP <= 0 {
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
					char.MaxResource += 8
					char.Resource = char.MaxResource
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

			m.addLog(session, "combat", fmt.Sprintf("%s 攻击了 %s，造成 %d 点伤害", enemy.Name, char.Name, enemyDamage), "#ff4444")
			logs = append(logs, session.BattleLogs[len(session.BattleLogs)-1])

			// 检查玩家是否死亡
			if char.HP <= 0 {
				char.HP = char.MaxHP / 2
				char.TotalDeaths++
				session.IsRunning = false
				session.CurrentEnemies = nil
				session.CurrentEnemy = nil
				session.CurrentTurnIndex = -1

				m.addLog(session, "death", fmt.Sprintf("%s 被击败了... 正在复活", char.Name), "#ff0000")
				logs = append(logs, session.BattleLogs[len(session.BattleLogs)-1])

				// 保存死亡数据
				m.charRepo.UpdateAfterDeath(char.ID, char.HP, char.TotalDeaths)
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
		initialHP := char.MaxHP
		initialMP := char.MaxResource
		restDuration := m.calculateRestTime(char, initialHP, initialMP)
		restUntil := time.Now().Add(restDuration)
		session.IsResting = true
		session.RestUntil = &restUntil
		session.RestSpeed = 1.0 // 默认恢复速度

		m.addLog(session, "system", fmt.Sprintf(">> 开始休息恢复 (预计 %d 秒)", int(restDuration.Seconds())+1), "#33ff33")
		logs = append(logs, session.BattleLogs[len(session.BattleLogs)-1])

		// 重置本场战斗统计
		session.CurrentBattleExp = 0
		session.CurrentBattleGold = 0
		session.CurrentBattleKills = 0
		session.CurrentTurnIndex = -1

		// 清除敌人
		session.CurrentEnemies = nil
		session.CurrentEnemy = nil
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

// calculateRestTime 计算休息时间（基于HP/MP损失）
func (m *BattleManager) calculateRestTime(char *models.Character, initialHP, initialMP int) time.Duration {
	hpLoss := float64(char.MaxHP - char.HP)
	mpLoss := float64(char.MaxResource - char.Resource)
	
	// 基础休息时间：每损失1% HP/MP = 0.1秒，最少1秒
	totalLoss := (hpLoss/float64(char.MaxHP) + mpLoss/float64(char.MaxResource)) / 2.0
	restSeconds := totalLoss * 10.0
	
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
	if !session.IsResting || session.RestUntil == nil {
		return
	}
	
	now := time.Now()
	if now.Before(*session.RestUntil) {
		// 计算恢复速度（每秒恢复一定百分比）
		restSpeed := session.RestSpeed
		if restSpeed <= 0 {
			restSpeed = 1.0
		}
		
		// 每秒恢复最大值的2%
		hpRegen := int(float64(char.MaxHP) * 0.02 * restSpeed)
		mpRegen := int(float64(char.MaxResource) * 0.02 * restSpeed)
		
		char.HP += hpRegen
		if char.HP > char.MaxHP {
			char.HP = char.MaxHP
		}
		
		char.Resource += mpRegen
		if char.Resource > char.MaxResource {
			char.Resource = char.MaxResource
		}
	} else {
		// 休息结束
		session.IsResting = false
		session.RestUntil = nil
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


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
	UserID       int
	IsRunning    bool
	CurrentZone  *models.Zone
	CurrentEnemy *models.Monster
	BattleLogs   []models.BattleLog
	BattleCount  int
	SessionKills int
	SessionGold  int
	SessionExp   int
	StartedAt    time.Time
	LastTick     time.Time
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
		UserID:     userID,
		BattleLogs: make([]models.BattleLog, 0),
		StartedAt:  time.Now(),
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

// ExecuteBattleTick 执行战斗回合
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

	// 如果没有当前敌人，生成一个
	if session.CurrentEnemy == nil {
		err := m.spawnEnemy(session, char.Level)
		if err != nil {
			return nil, err
		}
		logs = append(logs, session.BattleLogs[len(session.BattleLogs)-1])
	}

	enemy := session.CurrentEnemy

	// 执行战斗回合
	// 玩家攻击
	playerDamage := m.calculateDamage(char.Attack, enemy.Defense)
	isCrit := rand.Float64() < char.CritRate
	if isCrit {
		playerDamage = int(float64(playerDamage) * char.CritDamage)
	}
	enemy.HP -= playerDamage

	skillName := m.getRandomSkillName(char.ClassID)
	if isCrit {
		m.addLog(session, "combat", fmt.Sprintf("你使用 [%s] 💥暴击！对 %s 造成 %d 点伤害", skillName, enemy.Name, playerDamage), "#ff6b6b")
	} else {
		m.addLog(session, "combat", fmt.Sprintf("你使用 [%s] 对 %s 造成 %d 点伤害", skillName, enemy.Name, playerDamage), "#ffaa00")
	}
	logs = append(logs, session.BattleLogs[len(session.BattleLogs)-1])

	// 检查敌人是否死亡
	if enemy.HP <= 0 {
		// 胜利！
		expGain := enemy.ExpReward
		goldGain := enemy.GoldMin + rand.Intn(enemy.GoldMax-enemy.GoldMin+1)

		session.SessionExp += expGain
		session.SessionGold += goldGain
		session.SessionKills++

		m.addLog(session, "victory", fmt.Sprintf(">> %s 被击败！获得 %d 经验值", enemy.Name, expGain), "#33ff33")
		logs = append(logs, session.BattleLogs[len(session.BattleLogs)-1])
		m.addLog(session, "loot", fmt.Sprintf(">> 拾取 %d 金币", goldGain), "#ffd700")
		logs = append(logs, session.BattleLogs[len(session.BattleLogs)-1])

		// 更新角色数据
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

			m.addLog(session, "levelup", fmt.Sprintf("🎉【升级】恭喜！你升到了 %d 级！", char.Level), "#ffd700")
			logs = append(logs, session.BattleLogs[len(session.BattleLogs)-1])
		}

		// 保存角色数据更新
		m.charRepo.UpdateAfterBattle(char.ID, char.HP, char.Resource, char.Exp, char.Level,
			char.ExpToNext, char.MaxHP, char.MaxResource, char.Attack, char.Defense,
			char.Strength, char.Agility, char.Stamina, char.TotalKills)

		// 清除敌人，下回合生成新的
		session.CurrentEnemy = nil

		// 恢复一些HP
		healAmount := char.MaxHP / 10
		char.HP += healAmount
		if char.HP > char.MaxHP {
			char.HP = char.MaxHP
		}
	} else {
		// 敌人反击
		enemyDamage := m.calculateDamage(enemy.Attack, char.Defense)
		char.HP -= enemyDamage

		m.addLog(session, "combat", fmt.Sprintf("%s 攻击了你，造成 %d 点伤害", enemy.Name, enemyDamage), "#ff4444")
		logs = append(logs, session.BattleLogs[len(session.BattleLogs)-1])

		// 检查玩家是否死亡
		if char.HP <= 0 {
			char.HP = char.MaxHP / 2
			char.TotalDeaths++
			session.IsRunning = false
			session.CurrentEnemy = nil

			m.addLog(session, "death", fmt.Sprintf("你被 %s 击败了... 正在复活", enemy.Name), "#ff0000")
			logs = append(logs, session.BattleLogs[len(session.BattleLogs)-1])

			// 保存死亡数据
			m.charRepo.UpdateAfterDeath(char.ID, char.HP, char.TotalDeaths)
		}
	}

	return &BattleTickResult{
		Character:    char,
		Enemy:        session.CurrentEnemy,
		Logs:         logs,
		IsRunning:    session.IsRunning,
		SessionKills: session.SessionKills,
		SessionGold:  session.SessionGold,
		SessionExp:   session.SessionExp,
		BattleCount:  session.BattleCount,
	}, nil
}

// spawnEnemy 生成敌人
func (m *BattleManager) spawnEnemy(session *BattleSession, playerLevel int) error {
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

	// 随机选择一个怪物
	template := monsters[rand.Intn(len(monsters))]

	session.CurrentEnemy = &models.Monster{
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

	session.BattleCount++
	m.addLog(session, "encounter", fmt.Sprintf("━━━ 战斗 #%d ━━━ 遭遇: %s (Lv.%d)", session.BattleCount, template.Name, template.Level), "#ffff00")

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
		BattleCount:    session.BattleCount,
		TotalKills:     session.SessionKills,
		TotalExp:       session.SessionExp,
		TotalGold:      session.SessionGold,
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

// BattleTickResult 战斗回合结果
type BattleTickResult struct {
	Character    *models.Character `json:"character"`
	Enemy        *models.Monster   `json:"enemy,omitempty"`
	Logs         []models.BattleLog `json:"logs"`
	IsRunning    bool              `json:"isRunning"`
	SessionKills int               `json:"sessionKills"`
	SessionGold  int               `json:"sessionGold"`
	SessionExp   int               `json:"sessionExp"`
	BattleCount  int               `json:"battleCount"`
}


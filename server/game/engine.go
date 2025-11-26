package game

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// Engine 游戏引擎
type Engine struct {
	mu           sync.RWMutex
	character    *Character
	currentZone  *Zone
	currentEnemy *Monster
	battleLogs   []BattleLog
	isRunning    bool

	// 统计
	battleCount  int
	sessionKills int
	sessionGold  int
	sessionExp   int
}

// 预定义区域
var zones = map[string]*Zone{
	"elwynn": {
		ID:          "elwynn",
		Name:        "艾尔文森林",
		Description: "人类王国暴风城外的宁静森林，适合新手冒险者。",
		MinLevel:    1,
		Monsters: []Monster{
			{ID: "wolf", Name: "森林狼", Level: 1, MaxHP: 30, Attack: 5, Defense: 1, ExpDrop: 15, GoldMin: 1, GoldMax: 5},
			{ID: "kobold", Name: "狗头人矿工", Level: 2, MaxHP: 40, Attack: 6, Defense: 2, ExpDrop: 20, GoldMin: 2, GoldMax: 8},
			{ID: "defias", Name: "迪菲亚劫匪", Level: 3, MaxHP: 55, Attack: 8, Defense: 3, ExpDrop: 30, GoldMin: 5, GoldMax: 12},
		},
	},
	"westfall": {
		ID:          "westfall",
		Name:        "西部荒野",
		Description: "曾经肥沃的农田，如今被迪菲亚兄弟会占领。",
		MinLevel:    5,
		Monsters: []Monster{
			{ID: "harvest", Name: "收割傀儡", Level: 5, MaxHP: 80, Attack: 12, Defense: 5, ExpDrop: 50, GoldMin: 8, GoldMax: 20},
			{ID: "defias_rogue", Name: "迪菲亚盗贼", Level: 6, MaxHP: 90, Attack: 14, Defense: 6, ExpDrop: 65, GoldMin: 10, GoldMax: 25},
			{ID: "gnoll", Name: "豺狼人", Level: 7, MaxHP: 110, Attack: 16, Defense: 7, ExpDrop: 80, GoldMin: 12, GoldMax: 30},
		},
	},
	"duskwood": {
		ID:          "duskwood",
		Name:        "暮色森林",
		Description: "被永恒黑暗笼罩的诡异森林，亡灵与狼人出没。",
		MinLevel:    10,
		Monsters: []Monster{
			{ID: "skeleton", Name: "骷髅战士", Level: 10, MaxHP: 150, Attack: 22, Defense: 10, ExpDrop: 120, GoldMin: 15, GoldMax: 40},
			{ID: "ghoul", Name: "食尸鬼", Level: 12, MaxHP: 180, Attack: 26, Defense: 12, ExpDrop: 150, GoldMin: 20, GoldMax: 50},
			{ID: "worgen", Name: "狼人", Level: 14, MaxHP: 220, Attack: 32, Defense: 15, ExpDrop: 200, GoldMin: 25, GoldMax: 60},
		},
	},
}

// NewEngine 创建新游戏引擎
func NewEngine() *Engine {
	return &Engine{
		battleLogs: make([]BattleLog, 0),
	}
}

// CreateCharacter 创建角色
func (e *Engine) CreateCharacter(name, race, class string) *Character {
	e.mu.Lock()
	defer e.mu.Unlock()

	// 根据职业设置基础属性
	var maxHP, maxMP, attack, defense int
	switch class {
	case "warrior":
		maxHP, maxMP, attack, defense = 120, 20, 15, 8
	case "mage":
		maxHP, maxMP, attack, defense = 70, 100, 20, 3
	case "rogue":
		maxHP, maxMP, attack, defense = 90, 50, 18, 5
	case "priest":
		maxHP, maxMP, attack, defense = 80, 80, 12, 4
	default:
		maxHP, maxMP, attack, defense = 100, 50, 12, 5
	}

	e.character = &Character{
		Name:      name,
		Race:      race,
		ClassName: class,
		Level:     1,
		Exp:       0,
		ExpToNext: 100,
		MaxHP:     maxHP,
		HP:        maxHP,
		MaxMP:     maxMP,
		MP:        maxMP,
		Attack:    attack,
		Defense:   defense,
		Gold:      0,
	}

	// 默认区域
	e.currentZone = zones["elwynn"]

	// 初始日志
	e.addLog("info", fmt.Sprintf("欢迎来到艾泽拉斯，%s！", name), "#00FF00")
	e.addLog("info", fmt.Sprintf("你是一名 %s %s", getRaceName(race), getClassName(class)), "#00FF00")
	e.addLog("info", "输入命令开始你的冒险...", "#888888")

	return e.character
}

// GetCharacter 获取角色
func (e *Engine) GetCharacter() *Character {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.character
}

// ToggleBattle 切换战斗状态
func (e *Engine) ToggleBattle() bool {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.isRunning = !e.isRunning
	if e.isRunning {
		e.addLog("info", ">> 开始自动战斗...", "#00FF00")
	} else {
		e.addLog("info", ">> 暂停自动战斗", "#FFFF00")
	}
	return e.isRunning
}

// BattleTick 执行一次战斗回合
func (e *Engine) BattleTick() *BattleResult {
	e.mu.Lock()
	defer e.mu.Unlock()

	if e.character == nil {
		return nil
	}

	logs := make([]BattleLog, 0)

	// 如果没有当前敌人，生成一个
	if e.currentEnemy == nil || e.currentEnemy.HP <= 0 {
		e.spawnEnemy()
		logs = append(logs, e.battleLogs[len(e.battleLogs)-1])
	}

	// 玩家攻击
	playerDamage := e.calculateDamage(e.character.Attack, e.currentEnemy.Defense)
	e.currentEnemy.HP -= playerDamage
	e.addLog("damage", fmt.Sprintf("你使用 [英勇打击] 对 %s 造成 %d 点伤害", e.currentEnemy.Name, playerDamage), "#FF6B6B")
	logs = append(logs, e.battleLogs[len(e.battleLogs)-1])

	// 检查敌人是否死亡
	if e.currentEnemy.HP <= 0 {
		e.onEnemyDefeated()
		logs = append(logs, e.battleLogs[len(e.battleLogs)-3:]...)
	} else {
		// 敌人攻击
		enemyDamage := e.calculateDamage(e.currentEnemy.Attack, e.character.Defense)
		e.character.HP -= enemyDamage
		e.addLog("damage", fmt.Sprintf("%s 攻击你造成 %d 点伤害", e.currentEnemy.Name, enemyDamage), "#FF4444")
		logs = append(logs, e.battleLogs[len(e.battleLogs)-1])

		// 检查玩家是否死亡
		if e.character.HP <= 0 {
			e.onPlayerDeath()
			logs = append(logs, e.battleLogs[len(e.battleLogs)-1])
		}
	}

	return &BattleResult{
		Character: e.character,
		Enemy:     e.currentEnemy,
		Logs:      logs,
		Status:    e.getBattleStatus(),
	}
}

// spawnEnemy 生成敌人
func (e *Engine) spawnEnemy() {
	if e.currentZone == nil || len(e.currentZone.Monsters) == 0 {
		return
	}

	// 随机选择怪物
	template := e.currentZone.Monsters[rand.Intn(len(e.currentZone.Monsters))]

	e.currentEnemy = &Monster{
		ID:      template.ID,
		Name:    template.Name,
		Level:   template.Level,
		MaxHP:   template.MaxHP,
		HP:      template.MaxHP,
		Attack:  template.Attack,
		Defense: template.Defense,
		ExpDrop: template.ExpDrop,
		GoldMin: template.GoldMin,
		GoldMax: template.GoldMax,
	}

	e.battleCount++
	e.addLog("info", fmt.Sprintf("━━━ 战斗 #%d ━━━ 遭遇: %s (Lv.%d)", e.battleCount, e.currentEnemy.Name, e.currentEnemy.Level), "#FFFF00")
}

// onEnemyDefeated 敌人被击败
func (e *Engine) onEnemyDefeated() {
	// 获得经验
	expGain := e.currentEnemy.ExpDrop
	e.character.Exp += expGain
	e.sessionExp += expGain

	// 获得金币
	goldGain := e.currentEnemy.GoldMin + rand.Intn(e.currentEnemy.GoldMax-e.currentEnemy.GoldMin+1)
	e.character.Gold += goldGain
	e.sessionGold += goldGain

	// 更新统计
	e.character.TotalKills++
	e.sessionKills++

	e.addLog("exp", fmt.Sprintf(">> %s 被击败！获得 %d 经验值", e.currentEnemy.Name, expGain), "#00FF00")
	e.addLog("loot", fmt.Sprintf(">> 拾取 %d 金币", goldGain), "#FFD700")

	// 检查升级
	for e.character.Exp >= e.character.ExpToNext {
		e.character.Exp -= e.character.ExpToNext
		e.character.Level++
		e.character.ExpToNext = e.character.Level * 100

		// 升级属性提升
		e.character.MaxHP += 10
		e.character.HP = e.character.MaxHP
		e.character.MaxMP += 5
		e.character.MP = e.character.MaxMP
		e.character.Attack += 2
		e.character.Defense += 1

		e.addLog("levelup", fmt.Sprintf("★★★ 恭喜升级！你现在是 %d 级 ★★★", e.character.Level), "#FFD700")
	}

	e.currentEnemy = nil
}

// onPlayerDeath 玩家死亡
func (e *Engine) onPlayerDeath() {
	e.addLog("info", "你被击败了...正在复活", "#FF0000")
	e.character.HP = e.character.MaxHP
	e.character.MP = e.character.MaxMP
	e.currentEnemy = nil
	e.isRunning = false
}

// calculateDamage 计算伤害
func (e *Engine) calculateDamage(attack, defense int) int {
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
func (e *Engine) addLog(logType, message, color string) {
	log := BattleLog{
		Time:    time.Now().Format("15:04:05"),
		Type:    logType,
		Message: message,
		Color:   color,
	}
	e.battleLogs = append(e.battleLogs, log)

	// 保持日志数量在合理范围
	if len(e.battleLogs) > 100 {
		e.battleLogs = e.battleLogs[len(e.battleLogs)-100:]
	}
}

// GetBattleLogs 获取战斗日志
func (e *Engine) GetBattleLogs() []BattleLog {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.battleLogs
}

// GetBattleStatus 获取战斗状态
func (e *Engine) GetBattleStatus() BattleStatus {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.getBattleStatus()
}

func (e *Engine) getBattleStatus() BattleStatus {
	zoneName := ""
	if e.currentZone != nil {
		zoneName = e.currentZone.Name
	}
	return BattleStatus{
		IsRunning:    e.isRunning,
		CurrentZone:  zoneName,
		CurrentEnemy: e.currentEnemy,
		BattleCount:  e.battleCount,
		SessionKills: e.sessionKills,
		SessionGold:  e.sessionGold,
		SessionExp:   e.sessionExp,
	}
}

// GetZones 获取所有区域
func (e *Engine) GetZones() []*Zone {
	result := make([]*Zone, 0, len(zones))
	for _, zone := range zones {
		result = append(result, zone)
	}
	return result
}

// ChangeZone 切换区域
func (e *Engine) ChangeZone(zoneID string) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	zone, ok := zones[zoneID]
	if !ok {
		return fmt.Errorf("未知区域: %s", zoneID)
	}

	if e.character != nil && e.character.Level < zone.MinLevel {
		return fmt.Errorf("等级不足，需要 %d 级", zone.MinLevel)
	}

	e.currentZone = zone
	e.currentEnemy = nil
	e.addLog("info", fmt.Sprintf(">> 你来到了 [%s]", zone.Name), "#00FFFF")
	e.addLog("info", zone.Description, "#888888")

	return nil
}

// 辅助函数
func getRaceName(race string) string {
	names := map[string]string{
		"human":    "人类",
		"dwarf":    "矮人",
		"nightelf": "暗夜精灵",
		"gnome":    "侏儒",
		"orc":      "兽人",
		"undead":   "亡灵",
		"tauren":   "牛头人",
		"troll":    "巨魔",
	}
	if name, ok := names[race]; ok {
		return name
	}
	return race
}

func getClassName(class string) string {
	names := map[string]string{
		"warrior": "战士",
		"mage":    "法师",
		"rogue":   "盗贼",
		"priest":  "牧师",
		"paladin": "圣骑士",
		"hunter":  "猎人",
		"warlock": "术士",
		"druid":   "德鲁伊",
		"shaman":  "萨满",
	}
	if name, ok := names[class]; ok {
		return name
	}
	return class
}



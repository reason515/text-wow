package game

import (
	"fmt"
	"math/rand"
	"sync"
	"text-wow/internal/models"
	"time"
)

// Engine 游戏引擎
type Engine struct {
	mu            sync.RWMutex
	character     *models.Character
	strategy      *models.Strategy
	battleStatus  *models.BattleStatus
	battleLogs    []models.BattleLog
	currentZone   *models.Zone
	stopChan      chan struct{}
	skills        []models.Skill
	skillCooldowns map[string]int
}

var engine *Engine

func InitEngine() {
	engine = &Engine{
		battleStatus: &models.BattleStatus{
			IsRunning: false,
		},
		battleLogs:     make([]models.BattleLog, 0),
		skillCooldowns: make(map[string]int),
	}
}

func GetEngine() *Engine {
	return engine
}

// SetCharacter 设置角色
func (e *Engine) SetCharacter(c *models.Character) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.character = c

	// 设置职业技能
	if skills, ok := models.ClassSkills[c.Class]; ok {
		e.skills = skills
	} else {
		e.skills = models.ClassSkills["warrior"]
	}
}

// GetCharacter 获取角色
func (e *Engine) GetCharacter() *models.Character {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.character
}

// SetStrategy 设置策略
func (e *Engine) SetStrategy(s *models.Strategy) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.strategy = s
}

// GetStrategy 获取策略
func (e *Engine) GetStrategy() *models.Strategy {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.strategy
}

// SetZone 设置区域
func (e *Engine) SetZone(zoneID string) bool {
	e.mu.Lock()
	defer e.mu.Unlock()

	for _, z := range models.Zones {
		if z.ID == zoneID {
			e.currentZone = &z
			return true
		}
	}
	return false
}

// GetCurrentZone 获取当前区域
func (e *Engine) GetCurrentZone() *models.Zone {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.currentZone
}

// GetBattleStatus 获取战斗状态
func (e *Engine) GetBattleStatus() *models.BattleStatus {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.battleStatus
}

// GetBattleLogs 获取战斗日志
func (e *Engine) GetBattleLogs(limit int) []models.BattleLog {
	e.mu.RLock()
	defer e.mu.RUnlock()

	if limit <= 0 || limit > len(e.battleLogs) {
		limit = len(e.battleLogs)
	}

	start := len(e.battleLogs) - limit
	if start < 0 {
		start = 0
	}
	return e.battleLogs[start:]
}

// AddLog 添加日志
func (e *Engine) AddLog(message string, logType string) {
	e.mu.Lock()
	defer e.mu.Unlock()

	log := models.BattleLog{
		ID:        len(e.battleLogs) + 1,
		Message:   message,
		LogType:   logType,
		CreatedAt: time.Now(),
	}
	e.battleLogs = append(e.battleLogs, log)

	// 保留最近500条日志
	if len(e.battleLogs) > 500 {
		e.battleLogs = e.battleLogs[len(e.battleLogs)-500:]
	}
}

// StartBattle 开始自动战斗
func (e *Engine) StartBattle() bool {
	e.mu.Lock()
	if e.battleStatus.IsRunning {
		e.mu.Unlock()
		return false
	}

	if e.character == nil || e.currentZone == nil {
		e.mu.Unlock()
		return false
	}

	e.battleStatus.IsRunning = true
	now := time.Now()
	e.battleStatus.SessionStart = &now
	e.battleStatus.BattleCount = 0
	e.battleStatus.TotalKills = 0
	e.battleStatus.TotalExp = 0
	e.battleStatus.TotalGold = 0
	e.stopChan = make(chan struct{})
	e.mu.Unlock()

	e.AddLog(fmt.Sprintf("⚔️ 开始在 [%s] 自动战斗...", e.currentZone.Name), "system")

	go e.battleLoop()
	return true
}

// StopBattle 停止自动战斗
func (e *Engine) StopBattle() {
	e.mu.Lock()
	if !e.battleStatus.IsRunning {
		e.mu.Unlock()
		return
	}

	close(e.stopChan)
	e.battleStatus.IsRunning = false
	e.battleStatus.CurrentMonster = nil
	e.mu.Unlock()

	e.AddLog("🛑 停止自动战斗", "system")
}

// battleLoop 战斗循环
func (e *Engine) battleLoop() {
	for {
		select {
		case <-e.stopChan:
			return
		default:
			e.runSingleBattle()
			time.Sleep(500 * time.Millisecond) // 战斗间隔
		}
	}
}

// runSingleBattle 进行一场战斗
func (e *Engine) runSingleBattle() {
	e.mu.Lock()
	if e.currentZone == nil || len(e.currentZone.Monsters) == 0 {
		e.mu.Unlock()
		return
	}

	// 随机选择怪物
	monsterTemplate := e.currentZone.Monsters[rand.Intn(len(e.currentZone.Monsters))]
	monster := monsterTemplate // 复制一份
	e.battleStatus.CurrentMonster = &monster
	e.battleStatus.BattleCount++
	battleNum := e.battleStatus.BattleCount
	e.mu.Unlock()

	e.AddLog(fmt.Sprintf("━━━ 战斗 #%d ━━━", battleNum), "system")
	e.AddLog(fmt.Sprintf("🐺 遭遇: %s Lv.%d (HP: %d)", monster.Name, monster.Level, monster.HP), "combat")

	// 重置技能冷却
	e.mu.Lock()
	e.skillCooldowns = make(map[string]int)
	e.mu.Unlock()

	// 战斗回合
	round := 0
	for monster.HP > 0 && e.character.HP > 0 {
		select {
		case <-e.stopChan:
			return
		default:
		}

		round++
		e.executeBattleRound(&monster, round)
		time.Sleep(300 * time.Millisecond) // 回合间隔
	}

	// 战斗结束
	if monster.HP <= 0 {
		e.handleVictory(&monster)
	} else {
		e.handleDefeat()
	}
}

// executeBattleRound 执行战斗回合
func (e *Engine) executeBattleRound(monster *models.Monster, round int) {
	e.mu.Lock()
	defer e.mu.Unlock()

	// 玩家回合 - 根据策略选择技能
	skill := e.selectSkill()
	damage := e.calculateDamage(skill, monster)

	if e.character.MP >= skill.MPCost {
		e.character.MP -= skill.MPCost
		monster.HP -= damage
		if monster.HP < 0 {
			monster.HP = 0
		}

		e.mu.Unlock()
		e.AddLog(fmt.Sprintf("⚔️ 你使用 [%s] 造成 %d 点伤害 (怪物HP: %d)", skill.Name, damage, monster.HP), "combat")
		e.mu.Lock()

		// 设置技能冷却
		if skill.Cooldown > 0 {
			e.skillCooldowns[skill.ID] = skill.Cooldown
		}
	} else {
		// MP不足，使用普通攻击
		basicDamage := e.character.Strength + rand.Intn(5)
		monster.HP -= basicDamage

		e.mu.Unlock()
		e.AddLog(fmt.Sprintf("⚔️ 你进行普通攻击造成 %d 点伤害", basicDamage), "combat")
		e.mu.Lock()
	}

	// 怪物死亡检查
	if monster.HP <= 0 {
		return
	}

	// 怪物回合
	monsterDamage := monster.Attack - (e.character.Stamina / 2) + rand.Intn(5)
	if monsterDamage < 1 {
		monsterDamage = 1
	}
	e.character.HP -= monsterDamage

	e.mu.Unlock()
	e.AddLog(fmt.Sprintf("💥 %s 攻击你造成 %d 点伤害 (你的HP: %d/%d)", monster.Name, monsterDamage, e.character.HP, e.character.MaxHP), "combat")
	e.mu.Lock()

	// 减少所有技能冷却
	for id := range e.skillCooldowns {
		if e.skillCooldowns[id] > 0 {
			e.skillCooldowns[id]--
		}
	}
}

// selectSkill 根据策略选择技能
func (e *Engine) selectSkill() models.Skill {
	if e.strategy != nil && len(e.strategy.SkillPriority) > 0 {
		for _, skillID := range e.strategy.SkillPriority {
			// 检查冷却
			if cd, ok := e.skillCooldowns[skillID]; ok && cd > 0 {
				continue
			}

			// 查找技能
			for _, skill := range e.skills {
				if skill.ID == skillID && e.character.MP >= skill.MPCost {
					return skill
				}
			}
		}
	}

	// 默认返回普通攻击
	for _, skill := range e.skills {
		if skill.ID == "attack" {
			return skill
		}
	}

	return models.Skill{ID: "attack", Name: "普通攻击", Damage: 10, MPCost: 0}
}

// calculateDamage 计算伤害
func (e *Engine) calculateDamage(skill models.Skill, monster *models.Monster) int {
	baseDamage := skill.Damage + (e.character.Strength / 2)
	variance := rand.Intn(10) - 5
	damage := baseDamage + variance - (monster.Defense / 2)
	if damage < 1 {
		damage = 1
	}
	return damage
}

// handleVictory 处理胜利
func (e *Engine) handleVictory(monster *models.Monster) {
	e.mu.Lock()
	defer e.mu.Unlock()

	// 获得经验
	exp := monster.ExpReward
	e.character.Exp += exp
	e.battleStatus.TotalExp += exp

	// 获得金币
	gold := monster.GoldMin + rand.Intn(monster.GoldMax-monster.GoldMin+1)
	e.character.Gold += gold
	e.battleStatus.TotalGold += gold

	e.battleStatus.TotalKills++
	e.battleStatus.CurrentMonster = nil

	e.mu.Unlock()
	e.AddLog(fmt.Sprintf("🏆 击败 %s! 获得 %d 经验, %d 金币", monster.Name, exp, gold), "loot")
	e.mu.Lock()

	// 检查升级
	if e.character.Exp >= e.character.ExpToNext {
		e.levelUp()
	}

	// 回复一些HP和MP
	e.character.HP += e.character.MaxHP / 10
	if e.character.HP > e.character.MaxHP {
		e.character.HP = e.character.MaxHP
	}
	e.character.MP += e.character.MaxMP / 5
	if e.character.MP > e.character.MaxMP {
		e.character.MP = e.character.MaxMP
	}
}

// handleDefeat 处理失败
func (e *Engine) handleDefeat() {
	e.mu.Lock()
	e.character.HP = e.character.MaxHP / 2
	e.character.MP = e.character.MaxMP / 2
	e.battleStatus.CurrentMonster = nil
	e.mu.Unlock()

	e.AddLog("💀 你被击败了! 复活中...", "system")

	time.Sleep(2 * time.Second)

	e.mu.Lock()
	e.character.HP = e.character.MaxHP
	e.character.MP = e.character.MaxMP
	e.mu.Unlock()

	e.AddLog("✨ 你已复活，继续战斗!", "system")
}

// levelUp 升级
func (e *Engine) levelUp() {
	e.character.Level++
	e.character.Exp -= e.character.ExpToNext
	e.character.ExpToNext = int(float64(e.character.ExpToNext) * 1.5)

	// 属性提升
	e.character.MaxHP += 20
	e.character.HP = e.character.MaxHP
	e.character.MaxMP += 10
	e.character.MP = e.character.MaxMP
	e.character.Strength += 2
	e.character.Agility += 2
	e.character.Intellect += 2
	e.character.Stamina += 2
	e.character.Spirit += 2

	e.mu.Unlock()
	e.AddLog(fmt.Sprintf("🎉 升级! 你现在是 Lv.%d!", e.character.Level), "levelup")
	e.mu.Lock()
}


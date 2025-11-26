package game

import (
	"math/rand"
	"sync"
	"time"
)

// GameEngine 游戏引擎
type GameEngine struct {
	State     *GameState
	Zones     map[string]*Zone
	mu        sync.RWMutex
	battleID  int
	listeners []chan BattleResult
}

// NewGameEngine 创建新的游戏引擎
func NewGameEngine() *GameEngine {
	engine := &GameEngine{
		Zones:     make(map[string]*Zone),
		listeners: make([]chan BattleResult, 0),
	}

	// 初始化区域
	engine.initZones()

	// 创建默认角色
	engine.State = &GameState{
		Character:   engine.createDefaultCharacter(),
		CurrentZone: "elwynn_forest",
		IsAutoFight: false,
		BattleCount: 0,
	}

	return engine
}

// 初始化游戏区域
func (e *GameEngine) initZones() {
	e.Zones["elwynn_forest"] = &Zone{
		ID:          "elwynn_forest",
		Name:        "艾尔文森林",
		Description: "联盟的新手村，阳光透过树叶洒落，偶尔能听到狼嚎声...",
		MinLevel:    1,
		Monsters: []Monster{
			{ID: "wolf", Name: "森林狼", Level: 2, MaxHP: 45, Attack: 8, Defense: 2, ExpReward: 20, GoldDrop: 5,
				LootTable: []Loot{{Name: "狼皮", Chance: 0.6}, {Name: "狼牙", Chance: 0.3}}},
			{ID: "boar", Name: "野猪", Level: 3, MaxHP: 60, Attack: 10, Defense: 4, ExpReward: 30, GoldDrop: 8,
				LootTable: []Loot{{Name: "野猪肉", Chance: 0.7}, {Name: "野猪蹄", Chance: 0.2}}},
			{ID: "kobold", Name: "狗头人", Level: 4, MaxHP: 55, Attack: 12, Defense: 3, ExpReward: 35, GoldDrop: 12,
				LootTable: []Loot{{Name: "蜡烛", Chance: 0.8}, {Name: "破旧矿镐", Chance: 0.1}}},
		},
	}

	e.Zones["westfall"] = &Zone{
		ID:          "westfall",
		Name:        "西部荒野",
		Description: "一片荒芜的农田，迪菲亚盗贼在此横行...",
		MinLevel:    10,
		Monsters: []Monster{
			{ID: "defias", Name: "迪菲亚盗贼", Level: 11, MaxHP: 120, Attack: 22, Defense: 8, ExpReward: 80, GoldDrop: 25,
				LootTable: []Loot{{Name: "红色面罩", Chance: 0.3}, {Name: "盗贼匕首", Chance: 0.1}}},
			{ID: "harvest_golem", Name: "收割傀儡", Level: 12, MaxHP: 150, Attack: 25, Defense: 12, ExpReward: 100, GoldDrop: 30,
				LootTable: []Loot{{Name: "金属零件", Chance: 0.5}, {Name: "傀儡核心", Chance: 0.05}}},
		},
	}

	e.Zones["duskwood"] = &Zone{
		ID:          "duskwood",
		Name:        "暮色森林",
		Description: "永恒的黑夜笼罩着这片森林，亡灵和狼人在阴影中游荡...",
		MinLevel:    20,
		Monsters: []Monster{
			{ID: "skeleton", Name: "腐化骷髅", Level: 21, MaxHP: 200, Attack: 35, Defense: 15, ExpReward: 150, GoldDrop: 40,
				LootTable: []Loot{{Name: "骨片", Chance: 0.6}, {Name: "暗影精华", Chance: 0.1}}},
			{ID: "worgen", Name: "夜色镇狼人", Level: 23, MaxHP: 280, Attack: 45, Defense: 18, ExpReward: 200, GoldDrop: 55,
				LootTable: []Loot{{Name: "狼人之爪", Chance: 0.4}, {Name: "月亮护符", Chance: 0.05}}},
		},
	}
}

// 创建默认角色
func (e *GameEngine) createDefaultCharacter() *Character {
	return &Character{
		ID:        "player_1",
		Name:      "勇士",
		Race:      "人类",
		Class:     "战士",
		Level:     1,
		Exp:       0,
		ExpToNext: 100,
		MaxHP:     100,
		CurrentHP: 100,
		MaxMP:     50,
		CurrentMP: 50,
		Stats: Stats{
			Strength:  12,
			Agility:   8,
			Intellect: 5,
			Stamina:   10,
			Spirit:    6,
		},
		Skills: []Skill{
			{ID: "heroic_strike", Name: "英勇打击", Description: "一次强力的武器攻击", Damage: 25, ManaCost: 0, Cooldown: 0, Type: "physical"},
			{ID: "thunder_clap", Name: "雷霆一击", Description: "震荡周围敌人", Damage: 35, ManaCost: 10, Cooldown: 2, Type: "physical"},
			{ID: "execute", Name: "斩杀", Description: "对低血量敌人造成巨额伤害", Damage: 60, ManaCost: 15, Cooldown: 4, Type: "physical"},
		},
		Gold:       0,
		TotalKills: 0,
	}
}

// GetState 获取当前游戏状态
func (e *GameEngine) GetState() *GameState {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.State
}

// GetZones 获取所有区域
func (e *GameEngine) GetZones() map[string]*Zone {
	return e.Zones
}

// SetZone 设置当前区域
func (e *GameEngine) SetZone(zoneID string) bool {
	e.mu.Lock()
	defer e.mu.Unlock()

	zone, exists := e.Zones[zoneID]
	if !exists {
		return false
	}

	if e.State.Character.Level < zone.MinLevel {
		return false
	}

	e.State.CurrentZone = zoneID
	return true
}

// StartAutoFight 开始自动战斗
func (e *GameEngine) StartAutoFight() {
	e.mu.Lock()
	e.State.IsAutoFight = true
	e.mu.Unlock()
}

// StopAutoFight 停止自动战斗
func (e *GameEngine) StopAutoFight() {
	e.mu.Lock()
	e.State.IsAutoFight = false
	e.mu.Unlock()
}

// DoBattle 执行一次战斗
func (e *GameEngine) DoBattle() *BattleResult {
	e.mu.Lock()
	defer e.mu.Unlock()

	zone := e.Zones[e.State.CurrentZone]
	if zone == nil {
		return nil
	}

	// 随机选择一个怪物
	monsterTemplate := zone.Monsters[rand.Intn(len(zone.Monsters))]
	monster := Monster{
		ID:        monsterTemplate.ID,
		Name:      monsterTemplate.Name,
		Level:     monsterTemplate.Level,
		MaxHP:     monsterTemplate.MaxHP,
		CurrentHP: monsterTemplate.MaxHP,
		Attack:    monsterTemplate.Attack,
		Defense:   monsterTemplate.Defense,
		ExpReward: monsterTemplate.ExpReward,
		GoldDrop:  monsterTemplate.GoldDrop,
		LootTable: monsterTemplate.LootTable,
	}

	e.battleID++
	result := e.executeBattle(&monster)
	result.BattleCount = e.battleID

	// 更新统计
	e.State.BattleCount = e.battleID
	if result.Victory {
		e.State.TodayKills++
		e.State.TodayGold += result.GoldGained
		e.State.TodayExp += result.ExpGained
		e.State.Character.TotalKills++
	}

	// 战斗后恢复一些HP/MP
	e.State.Character.CurrentHP += e.State.Character.MaxHP / 10
	if e.State.Character.CurrentHP > e.State.Character.MaxHP {
		e.State.Character.CurrentHP = e.State.Character.MaxHP
	}
	e.State.Character.CurrentMP += e.State.Character.MaxMP / 5
	if e.State.Character.CurrentMP > e.State.Character.MaxMP {
		e.State.Character.CurrentMP = e.State.Character.MaxMP
	}

	// 重置技能冷却
	for i := range e.State.Character.Skills {
		e.State.Character.Skills[i].CurrentCD = 0
	}

	return result
}

// executeBattle 执行战斗逻辑
func (e *GameEngine) executeBattle(monster *Monster) *BattleResult {
	result := &BattleResult{
		Logs:       make([]BattleLog, 0),
		LootGained: make([]string, 0),
	}

	char := e.State.Character
	round := 0

	// 战斗开始日志
	result.Logs = append(result.Logs, BattleLog{
		Round:     0,
		Actor:     "system",
		Action:    "encounter",
		Message:   "【战斗开始】遭遇: " + monster.Name + " (Lv." + itoa(monster.Level) + ")",
		Timestamp: time.Now().UnixMilli(),
	})

	// 回合制战斗
	for char.CurrentHP > 0 && monster.CurrentHP > 0 {
		round++

		// 玩家回合
		skill := e.selectSkill()
		damage := e.calculateDamage(char.Stats.Strength, skill.Damage, monster.Defense)

		// 暴击判定 (基于敏捷)
		isCrit := rand.Float64() < float64(char.Stats.Agility)/100.0
		if isCrit {
			damage = int(float64(damage) * 1.5)
		}

		monster.CurrentHP -= damage
		char.CurrentMP -= skill.ManaCost

		result.Logs = append(result.Logs, BattleLog{
			Round:     round,
			Actor:     char.Name,
			Action:    skill.Name,
			Target:    monster.Name,
			Damage:    damage,
			IsCrit:    isCrit,
			Message:   e.formatPlayerAttackMessage(char.Name, skill.Name, monster.Name, damage, isCrit),
			Timestamp: time.Now().UnixMilli(),
		})

		// 更新技能冷却
		for i := range char.Skills {
			if char.Skills[i].CurrentCD > 0 {
				char.Skills[i].CurrentCD--
			}
		}
		if skill.Cooldown > 0 {
			for i := range char.Skills {
				if char.Skills[i].ID == skill.ID {
					char.Skills[i].CurrentCD = skill.Cooldown
				}
			}
		}

		// 检查怪物是否死亡
		if monster.CurrentHP <= 0 {
			break
		}

		// 怪物回合
		monsterDamage := e.calculateDamage(monster.Attack, 0, char.Stats.Stamina/2)
		char.CurrentHP -= monsterDamage

		result.Logs = append(result.Logs, BattleLog{
			Round:     round,
			Actor:     monster.Name,
			Action:    "攻击",
			Target:    char.Name,
			Damage:    monsterDamage,
			Message:   monster.Name + " 攻击了你，造成 " + itoa(monsterDamage) + " 点伤害",
			Timestamp: time.Now().UnixMilli(),
		})
	}

	// 战斗结果
	if monster.CurrentHP <= 0 {
		result.Victory = true
		result.ExpGained = monster.ExpReward
		result.GoldGained = monster.GoldDrop + rand.Intn(monster.GoldDrop/2)

		// 经验和金币
		char.Exp += result.ExpGained
		char.Gold += result.GoldGained

		result.Logs = append(result.Logs, BattleLog{
			Round:     round,
			Actor:     "system",
			Action:    "victory",
			Message:   "【胜利】" + monster.Name + " 被击败！获得 " + itoa(result.ExpGained) + " 经验, " + itoa(result.GoldGained) + " 金币",
			Timestamp: time.Now().UnixMilli(),
		})

		// 掉落判定
		for _, loot := range monster.LootTable {
			if rand.Float64() < loot.Chance {
				result.LootGained = append(result.LootGained, loot.Name)
				result.Logs = append(result.Logs, BattleLog{
					Round:     round,
					Actor:     "system",
					Action:    "loot",
					Message:   "获得物品: [" + loot.Name + "]",
					Timestamp: time.Now().UnixMilli(),
				})
			}
		}

		// 检查升级
		if char.Exp >= char.ExpToNext {
			char.Level++
			char.Exp -= char.ExpToNext
			char.ExpToNext = int(float64(char.ExpToNext) * 1.5)

			// 提升属性
			char.Stats.Strength += 2
			char.Stats.Agility += 1
			char.Stats.Stamina += 2
			char.Stats.Intellect += 1
			char.Stats.Spirit += 1

			// 提升HP/MP上限
			char.MaxHP += 15
			char.MaxMP += 8
			char.CurrentHP = char.MaxHP
			char.CurrentMP = char.MaxMP

			result.LevelUp = true
			result.Logs = append(result.Logs, BattleLog{
				Round:     round,
				Actor:     "system",
				Action:    "levelup",
				Message:   "🎉【升级】恭喜！你升到了 " + itoa(char.Level) + " 级！",
				Timestamp: time.Now().UnixMilli(),
			})
		}
	} else {
		result.Victory = false
		result.Logs = append(result.Logs, BattleLog{
			Round:     round,
			Actor:     "system",
			Action:    "defeat",
			Message:   "【战败】你被 " + monster.Name + " 击败了...",
			Timestamp: time.Now().UnixMilli(),
		})

		// 复活并恢复一半HP
		char.CurrentHP = char.MaxHP / 2
	}

	return result
}

// selectSkill 选择技能（简单AI策略）
func (e *GameEngine) selectSkill() *Skill {
	char := e.State.Character

	// 优先使用可用的高伤害技能
	for i := range char.Skills {
		skill := &char.Skills[i]
		if skill.CurrentCD == 0 && char.CurrentMP >= skill.ManaCost && skill.Damage > 30 {
			return skill
		}
	}

	// 默认使用普通攻击
	return &char.Skills[0]
}

// calculateDamage 计算伤害
func (e *GameEngine) calculateDamage(attack, skillDamage, defense int) int {
	baseDamage := attack + skillDamage - defense/2
	// 添加一些随机波动
	variance := rand.Intn(baseDamage/5+1) - baseDamage/10
	damage := baseDamage + variance
	if damage < 1 {
		damage = 1
	}
	return damage
}

// formatPlayerAttackMessage 格式化玩家攻击消息
func (e *GameEngine) formatPlayerAttackMessage(player, skill, target string, damage int, isCrit bool) string {
	msg := "你使用了 [" + skill + "]"
	if isCrit {
		msg += " 💥暴击！"
	}
	msg += " 对 " + target + " 造成 " + itoa(damage) + " 点伤害"
	return msg
}

// 简单的int转string
func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	if n < 0 {
		return "-" + itoa(-n)
	}
	var digits []byte
	for n > 0 {
		digits = append([]byte{byte('0' + n%10)}, digits...)
		n /= 10
	}
	return string(digits)
}

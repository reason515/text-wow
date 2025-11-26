package game

// 属性类型
type Stats struct {
	Strength  int `json:"strength"`  // 力量 - 物理攻击
	Agility   int `json:"agility"`   // 敏捷 - 暴击/闪避
	Intellect int `json:"intellect"` // 智力 - 法术攻击
	Stamina   int `json:"stamina"`   // 耐力 - 生命值
	Spirit    int `json:"spirit"`    // 精神 - 法力回复
}

// 角色
type Character struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Race      string `json:"race"`      // 种族
	Class     string `json:"class"`     // 职业
	Level     int    `json:"level"`     // 等级
	Exp       int    `json:"exp"`       // 经验值
	ExpToNext int    `json:"expToNext"` // 升级所需经验

	MaxHP     int `json:"maxHp"`     // 最大生命值
	CurrentHP int `json:"currentHp"` // 当前生命值
	MaxMP     int `json:"maxMp"`     // 最大法力值
	CurrentMP int `json:"currentMp"` // 当前法力值

	Stats  Stats   `json:"stats"`  // 基础属性
	Skills []Skill `json:"skills"` // 技能列表

	Gold       int `json:"gold"`       // 金币
	TotalKills int `json:"totalKills"` // 总击杀数
}

// 技能
type Skill struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Damage      int    `json:"damage"`   // 伤害值
	ManaCost    int    `json:"manaCost"` // 法力消耗
	Cooldown    int    `json:"cooldown"` // 冷却回合数
	CurrentCD   int    `json:"currentCd"`
	Type        string `json:"type"` // physical/magical
}

// 怪物
type Monster struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Level     int    `json:"level"`
	MaxHP     int    `json:"maxHp"`
	CurrentHP int    `json:"currentHp"`
	Attack    int    `json:"attack"`
	Defense   int    `json:"defense"`
	ExpReward int    `json:"expReward"` // 经验奖励
	GoldDrop  int    `json:"goldDrop"`  // 金币掉落
	LootTable []Loot `json:"lootTable"` // 掉落表
}

// 掉落物品
type Loot struct {
	Name   string  `json:"name"`
	Chance float64 `json:"chance"` // 掉落概率 0-1
}

// 战斗日志条目
type BattleLog struct {
	Round     int    `json:"round"`
	Actor     string `json:"actor"`
	Action    string `json:"action"`
	Target    string `json:"target"`
	Damage    int    `json:"damage,omitempty"`
	Healing   int    `json:"healing,omitempty"`
	IsCrit    bool   `json:"isCrit,omitempty"`
	Message   string `json:"message"`
	Timestamp int64  `json:"timestamp"`
}

// 战斗结果
type BattleResult struct {
	Victory     bool        `json:"victory"`
	Logs        []BattleLog `json:"logs"`
	ExpGained   int         `json:"expGained"`
	GoldGained  int         `json:"goldGained"`
	LootGained  []string    `json:"lootGained"`
	LevelUp     bool        `json:"levelUp"`
	BattleCount int         `json:"battleCount"`
}

// 游戏状态
type GameState struct {
	Character   *Character `json:"character"`
	CurrentZone string     `json:"currentZone"`
	IsAutoFight bool       `json:"isAutoFight"`
	BattleCount int        `json:"battleCount"`
	TodayKills  int        `json:"todayKills"`
	TodayGold   int        `json:"todayGold"`
	TodayExp    int        `json:"todayExp"`
}

// 区域
type Zone struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	MinLevel    int       `json:"minLevel"`
	Monsters    []Monster `json:"monsters"`
}

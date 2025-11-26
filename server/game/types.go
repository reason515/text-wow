package game

// Character 玩家角色
type Character struct {
	Name      string `json:"name"`
	Race      string `json:"race"`
	ClassName string `json:"class"`
	Level     int    `json:"level"`
	Exp       int    `json:"exp"`
	ExpToNext int    `json:"exp_to_next"`

	// 属性
	MaxHP   int `json:"max_hp"`
	HP      int `json:"hp"`
	MaxMP   int `json:"max_mp"`
	MP      int `json:"mp"`
	Attack  int `json:"attack"`
	Defense int `json:"defense"`

	// 统计
	Gold       int `json:"gold"`
	TotalKills int `json:"total_kills"`
}

// Monster 怪物
type Monster struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Level   int    `json:"level"`
	MaxHP   int    `json:"max_hp"`
	HP      int    `json:"hp"`
	Attack  int    `json:"attack"`
	Defense int    `json:"defense"`
	ExpDrop int    `json:"exp_drop"`
	GoldMin int    `json:"gold_min"`
	GoldMax int    `json:"gold_max"`
}

// Zone 区域
type Zone struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	MinLevel    int       `json:"min_level"`
	Monsters    []Monster `json:"monsters"`
}

// BattleLog 战斗日志
type BattleLog struct {
	Time    string `json:"time"`
	Type    string `json:"type"` // info, damage, heal, loot, exp, levelup
	Message string `json:"message"`
	Color   string `json:"color"` // 颜色代码
}

// BattleStatus 战斗状态
type BattleStatus struct {
	IsRunning    bool     `json:"is_running"`
	CurrentZone  string   `json:"current_zone"`
	CurrentEnemy *Monster `json:"current_enemy"`
	BattleCount  int      `json:"battle_count"`
	SessionKills int      `json:"session_kills"`
	SessionGold  int      `json:"session_gold"`
	SessionExp   int      `json:"session_exp"`
}

// BattleResult 战斗回合结果
type BattleResult struct {
	Character *Character   `json:"character"`
	Enemy     *Monster     `json:"enemy"`
	Logs      []BattleLog  `json:"logs"`
	Status    BattleStatus `json:"status"`
}


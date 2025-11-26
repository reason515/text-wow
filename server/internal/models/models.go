package models

import "time"

// Character 角色
type Character struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Faction     string    `json:"faction"`     // alliance, horde
	Race        string    `json:"race"`        // human, orc, etc.
	Class       string    `json:"class"`       // warrior, mage, etc.
	Level       int       `json:"level"`
	Exp         int       `json:"exp"`
	ExpToNext   int       `json:"expToNext"`
	HP          int       `json:"hp"`
	MaxHP       int       `json:"maxHp"`
	MP          int       `json:"mp"`
	MaxMP       int       `json:"maxMp"`
	Strength    int       `json:"strength"`
	Agility     int       `json:"agility"`
	Intellect   int       `json:"intellect"`
	Stamina     int       `json:"stamina"`
	Spirit      int       `json:"spirit"`
	Gold        int       `json:"gold"`
	CurrentZone string    `json:"currentZone"`
	CreatedAt   time.Time `json:"createdAt"`
}

// Strategy 战斗策略
type Strategy struct {
	ID               int      `json:"id"`
	CharacterID      int      `json:"characterId"`
	SkillPriority    []string `json:"skillPriority"`    // 技能优先级
	HPPotionThreshold int     `json:"hpPotionThreshold"` // HP低于此百分比使用药水
	MPPotionThreshold int     `json:"mpPotionThreshold"` // MP低于此百分比使用药水
	TargetPriority   string   `json:"targetPriority"`   // lowest_hp, highest_hp, random
	AutoLoot         bool     `json:"autoLoot"`         // 自动拾取
}

// Monster 怪物
type Monster struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Level    int    `json:"level"`
	HP       int    `json:"hp"`
	MaxHP    int    `json:"maxHp"`
	Attack   int    `json:"attack"`
	Defense  int    `json:"defense"`
	ExpReward int   `json:"expReward"`
	GoldMin  int    `json:"goldMin"`
	GoldMax  int    `json:"goldMax"`
}

// Zone 区域
type Zone struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	MinLevel    int       `json:"minLevel"`
	MaxLevel    int       `json:"maxLevel"`
	Monsters    []Monster `json:"monsters"`
	Faction     string    `json:"faction"` // alliance, horde, neutral
}

// BattleLog 战斗日志
type BattleLog struct {
	ID        int       `json:"id"`
	Message   string    `json:"message"`
	LogType   string    `json:"logType"` // combat, loot, system, levelup
	CreatedAt time.Time `json:"createdAt"`
}

// BattleStatus 战斗状态
type BattleStatus struct {
	IsRunning     bool       `json:"isRunning"`
	CurrentMonster *Monster  `json:"currentMonster"`
	BattleCount   int        `json:"battleCount"`
	TotalKills    int        `json:"totalKills"`
	TotalExp      int        `json:"totalExp"`
	TotalGold     int        `json:"totalGold"`
	SessionStart  *time.Time `json:"sessionStart"`
}

// Skill 技能
type Skill struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Damage      int    `json:"damage"`
	MPCost      int    `json:"mpCost"`
	Cooldown    int    `json:"cooldown"` // 回合冷却
}

// ClassSkills 职业技能映射
var ClassSkills = map[string][]Skill{
	"warrior": {
		{ID: "attack", Name: "普通攻击", Description: "基础物理攻击", Damage: 10, MPCost: 0, Cooldown: 0},
		{ID: "heroic_strike", Name: "英勇打击", Description: "强力一击", Damage: 25, MPCost: 10, Cooldown: 0},
		{ID: "thunder_clap", Name: "雷霆一击", Description: "范围伤害", Damage: 15, MPCost: 15, Cooldown: 2},
		{ID: "execute", Name: "斩杀", Description: "对低血量目标伤害翻倍", Damage: 40, MPCost: 20, Cooldown: 3},
	},
	"mage": {
		{ID: "attack", Name: "普通攻击", Description: "基础攻击", Damage: 5, MPCost: 0, Cooldown: 0},
		{ID: "fireball", Name: "火球术", Description: "发射火球", Damage: 30, MPCost: 15, Cooldown: 0},
		{ID: "frostbolt", Name: "寒冰箭", Description: "冰系攻击", Damage: 25, MPCost: 12, Cooldown: 0},
		{ID: "arcane_missiles", Name: "奥术飞弹", Description: "连续攻击", Damage: 35, MPCost: 20, Cooldown: 2},
	},
	"hunter": {
		{ID: "attack", Name: "普通攻击", Description: "弓箭射击", Damage: 12, MPCost: 0, Cooldown: 0},
		{ID: "aimed_shot", Name: "瞄准射击", Description: "精准射击", Damage: 28, MPCost: 12, Cooldown: 1},
		{ID: "multi_shot", Name: "多重射击", Description: "多目标攻击", Damage: 18, MPCost: 15, Cooldown: 2},
		{ID: "kill_shot", Name: "杀戮射击", Description: "致命一击", Damage: 45, MPCost: 25, Cooldown: 3},
	},
	"rogue": {
		{ID: "attack", Name: "普通攻击", Description: "匕首攻击", Damage: 12, MPCost: 0, Cooldown: 0},
		{ID: "sinister_strike", Name: "邪恶攻击", Description: "积累连击点", Damage: 20, MPCost: 8, Cooldown: 0},
		{ID: "eviscerate", Name: "剔骨", Description: "消耗连击点", Damage: 35, MPCost: 15, Cooldown: 1},
		{ID: "backstab", Name: "背刺", Description: "暴击伤害", Damage: 40, MPCost: 20, Cooldown: 2},
	},
	"priest": {
		{ID: "attack", Name: "普通攻击", Description: "法杖攻击", Damage: 5, MPCost: 0, Cooldown: 0},
		{ID: "smite", Name: "惩击", Description: "神圣伤害", Damage: 22, MPCost: 12, Cooldown: 0},
		{ID: "shadow_word_pain", Name: "暗言术:痛", Description: "持续伤害", Damage: 15, MPCost: 10, Cooldown: 0},
		{ID: "mind_blast", Name: "心灵震爆", Description: "精神攻击", Damage: 35, MPCost: 20, Cooldown: 2},
	},
}

// Zones 区域数据
var Zones = []Zone{
	{
		ID:          "elwynn_forest",
		Name:        "艾尔文森林",
		Description: "联盟新手村，风景优美的森林地带",
		MinLevel:    1,
		MaxLevel:    10,
		Faction:     "alliance",
		Monsters: []Monster{
			{ID: "wolf", Name: "森林狼", Level: 2, HP: 50, MaxHP: 50, Attack: 8, Defense: 2, ExpReward: 15, GoldMin: 1, GoldMax: 3},
			{ID: "kobold", Name: "狗头人", Level: 3, HP: 65, MaxHP: 65, Attack: 10, Defense: 3, ExpReward: 20, GoldMin: 2, GoldMax: 5},
			{ID: "defias_thug", Name: "迪菲亚暴徒", Level: 5, HP: 90, MaxHP: 90, Attack: 14, Defense: 5, ExpReward: 35, GoldMin: 3, GoldMax: 8},
			{ID: "murloc", Name: "鱼人", Level: 4, HP: 70, MaxHP: 70, Attack: 12, Defense: 3, ExpReward: 25, GoldMin: 2, GoldMax: 6},
		},
	},
	{
		ID:          "durotar",
		Name:        "杜隆塔尔",
		Description: "部落新手村，炎热的沙漠地带",
		MinLevel:    1,
		MaxLevel:    10,
		Faction:     "horde",
		Monsters: []Monster{
			{ID: "scorpion", Name: "蝎子", Level: 2, HP: 45, MaxHP: 45, Attack: 9, Defense: 4, ExpReward: 15, GoldMin: 1, GoldMax: 3},
			{ID: "boar", Name: "野猪", Level: 3, HP: 70, MaxHP: 70, Attack: 11, Defense: 3, ExpReward: 20, GoldMin: 2, GoldMax: 5},
			{ID: "vile_familiar", Name: "邪恶小鬼", Level: 4, HP: 55, MaxHP: 55, Attack: 15, Defense: 2, ExpReward: 28, GoldMin: 3, GoldMax: 7},
			{ID: "razormane", Name: "钢鬃豪猪人", Level: 5, HP: 95, MaxHP: 95, Attack: 13, Defense: 6, ExpReward: 35, GoldMin: 3, GoldMax: 8},
		},
	},
	{
		ID:          "westfall",
		Name:        "西部荒野",
		Description: "曾经繁荣的农田，如今被迪菲亚占据",
		MinLevel:    10,
		MaxLevel:    20,
		Faction:     "alliance",
		Monsters: []Monster{
			{ID: "harvest_golem", Name: "收割傀儡", Level: 12, HP: 180, MaxHP: 180, Attack: 22, Defense: 12, ExpReward: 80, GoldMin: 8, GoldMax: 15},
			{ID: "defias_pillager", Name: "迪菲亚掠夺者", Level: 14, HP: 220, MaxHP: 220, Attack: 28, Defense: 10, ExpReward: 100, GoldMin: 10, GoldMax: 20},
			{ID: "gnoll", Name: "豺狼人", Level: 11, HP: 160, MaxHP: 160, Attack: 20, Defense: 8, ExpReward: 65, GoldMin: 6, GoldMax: 12},
		},
	},
	{
		ID:          "barrens",
		Name:        "贫瘠之地",
		Description: "广阔的草原，部落的必经之路",
		MinLevel:    10,
		MaxLevel:    25,
		Faction:     "horde",
		Monsters: []Monster{
			{ID: "zhevra", Name: "斑马", Level: 12, HP: 150, MaxHP: 150, Attack: 18, Defense: 8, ExpReward: 70, GoldMin: 5, GoldMax: 12},
			{ID: "centaur", Name: "半人马", Level: 15, HP: 250, MaxHP: 250, Attack: 30, Defense: 12, ExpReward: 120, GoldMin: 12, GoldMax: 25},
			{ID: "harpy", Name: "鹰身人", Level: 13, HP: 180, MaxHP: 180, Attack: 25, Defense: 6, ExpReward: 85, GoldMin: 8, GoldMax: 18},
		},
	},
}


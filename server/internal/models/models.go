package models

import "time"

// ═══════════════════════════════════════════════════════════
// 用户相关
// ═══════════════════════════════════════════════════════════

// User 用户
type User struct {
	ID              int        `json:"id"`
	Username        string     `json:"username"`
	Email           string     `json:"email,omitempty"`
	MaxTeamSize     int        `json:"maxTeamSize"`
	UnlockedSlots   int        `json:"unlockedSlots"`
	Gold            int        `json:"gold"`
	CurrentZoneID   string     `json:"currentZoneId"`
	TotalKills      int        `json:"totalKills"`
	TotalGoldGained int        `json:"totalGoldGained"`
	PlayTime        int        `json:"playTime"` // 秒
	CreatedAt       time.Time  `json:"createdAt"`
	LastLoginAt     *time.Time `json:"lastLoginAt,omitempty"`
}

// UserCredentials 用户登录凭据
type UserCredentials struct {
	Username string `json:"username" binding:"required,min=2,max=32"`
	Password string `json:"password" binding:"required,min=6,max=64"`
}

// UserRegister 用户注册请求
type UserRegister struct {
	Username string `json:"username" binding:"required,min=2,max=32"`
	Password string `json:"password" binding:"required,min=6,max=64"`
	Email    string `json:"email" binding:"omitempty,email"`
}

// AuthResponse 认证响应
type AuthResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

// ═══════════════════════════════════════════════════════════
// 角色相关
// ═══════════════════════════════════════════════════════════

// Character 角色
type Character struct {
	ID              int         `json:"id"`
	UserID          int         `json:"userId"`
	Name            string      `json:"name"`
	RaceID          string      `json:"raceId"`
	ClassID         string      `json:"classId"`
	Faction         string      `json:"faction"`  // alliance, horde
	TeamSlot        int         `json:"teamSlot"` // 1-5
	IsActive        bool        `json:"isActive"`
	IsDead          bool        `json:"isDead"`
	ReviveAt        *time.Time  `json:"reviveAt,omitempty"`
	Level           int         `json:"level"`
	Exp             int         `json:"exp"`
	ExpToNext       int         `json:"expToNext"`
	HP              int         `json:"hp"`
	MaxHP           int         `json:"maxHp"`
	Resource        int         `json:"resource"`     // 当前能量
	MaxResource     int         `json:"maxResource"`  // 最大能量
	ResourceType    string      `json:"resourceType"` // mana/rage/energy
	Strength        int         `json:"strength"`
	Agility         int         `json:"agility"`
	Intellect       int         `json:"intellect"`
	Stamina         int         `json:"stamina"`
	Spirit          int         `json:"spirit"`
	UnspentPoints   int         `json:"unspentPoints"`
	PhysicalAttack  int         `json:"physicalAttack"`
	MagicAttack     int         `json:"magicAttack"`
	PhysicalDefense int         `json:"physicalDefense"`
	MagicDefense    int         `json:"magicDefense"`
	PhysCritRate    float64     `json:"physCritRate"`    // 物理暴击率
	PhysCritDamage  float64     `json:"physCritDamage"`  // 物理暴击伤害
	SpellCritRate   float64     `json:"spellCritRate"`   // 法术暴击率
	SpellCritDamage float64     `json:"spellCritDamage"` // 法术暴击伤害
	DodgeRate       float64     `json:"dodgeRate"`       // 闪避率
	TotalKills      int         `json:"totalKills"`
	TotalDeaths     int         `json:"totalDeaths"`
	CreatedAt       time.Time   `json:"createdAt"`
	Buffs           []*BuffInfo `json:"buffs,omitempty"` // Buff/Debuff信息（不存储在数据库）
}

// BuffInfo Buff/Debuff信息（用于API返回）
type BuffInfo struct {
	EffectID     string  `json:"effectId"`
	Name         string  `json:"name"`
	Type         string  `json:"type"`
	IsBuff       bool    `json:"isBuff"`
	Duration     int     `json:"duration"`
	Value        float64 `json:"value"`
	StatAffected string  `json:"statAffected"`
	Description  string  `json:"description,omitempty"` // 效果描述
}

// CharacterCreate 创建角色请求
type CharacterCreate struct {
	Name    string `json:"name" binding:"required,min=2,max=32"`
	RaceID  string `json:"raceId" binding:"required"`
	ClassID string `json:"classId" binding:"required"`
}

// Team 小队信息
type Team struct {
	UserID        int          `json:"userId"`
	MaxSize       int          `json:"maxSize"`
	UnlockedSlots int          `json:"unlockedSlots"`
	Characters    []*Character `json:"characters"`
}

// ═══════════════════════════════════════════════════════════
// 种族和职业配置
// ═══════════════════════════════════════════════════════════

// Race 种族配置
type Race struct {
	ID            string  `json:"id"`
	Name          string  `json:"name"`
	Faction       string  `json:"faction"`
	Description   string  `json:"description"`
	StrengthBase  int     `json:"strengthBase"`
	AgilityBase   int     `json:"agilityBase"`
	IntellectBase int     `json:"intellectBase"`
	StaminaBase   int     `json:"staminaBase"`
	SpiritBase    int     `json:"spiritBase"`
	StrengthPct   float64 `json:"strengthPct"`
	AgilityPct    float64 `json:"agilityPct"`
	IntellectPct  float64 `json:"intellectPct"`
	StaminaPct    float64 `json:"staminaPct"`
	SpiritPct     float64 `json:"spiritPct"`
}

// Class 职业配置
type Class struct {
	ID                 string  `json:"id"`
	Name               string  `json:"name"`
	Description        string  `json:"description"`
	Role               string  `json:"role"` // tank/healer/dps/hybrid
	PrimaryStat        string  `json:"primaryStat"`
	ResourceType       string  `json:"resourceType"` // mana/rage/energy
	BaseHP             int     `json:"baseHp"`
	BaseResource       int     `json:"baseResource"`
	HPPerLevel         int     `json:"hpPerLevel"`
	ResourcePerLevel   int     `json:"resourcePerLevel"`
	ResourceRegen      float64 `json:"resourceRegen"`
	ResourceRegenPct   float64 `json:"resourceRegenPct"`
	BaseStrength       int     `json:"baseStrength"`
	BaseAgility        int     `json:"baseAgility"`
	BaseIntellect      int     `json:"baseIntellect"`
	BaseStamina        int     `json:"baseStamina"`
	BaseSpirit         int     `json:"baseSpirit"`
	BaseThreatModifier float64 `json:"baseThreatModifier"`
	CombatRole         string  `json:"combatRole"`
	IsRanged           bool    `json:"isRanged"`
}

// ═══════════════════════════════════════════════════════════
// 区域和怪物
// ═══════════════════════════════════════════════════════════

// Zone 区域
type Zone struct {
	ID                  string    `json:"id"`
	Name                string    `json:"name"`
	Description         string    `json:"description"`
	MinLevel            int       `json:"minLevel"`
	MaxLevel            int       `json:"maxLevel"`
	Faction             string    `json:"faction"` // alliance/horde/neutral
	ExpMulti            float64   `json:"expMulti"`
	GoldMulti           float64   `json:"goldMulti"`
	UnlockZoneID        *string   `json:"unlockZoneId,omitempty"`        // 需要探索的前置地图ID
	RequiredExploration int       `json:"requiredExploration"`            // 解锁所需探索度
	Monsters            []Monster `json:"monsters,omitempty"`
}

// ZoneExploration 玩家地图探索度
type ZoneExploration struct {
	UserID      int    `json:"userId"`
	ZoneID      string `json:"zoneId"`
	Exploration int    `json:"exploration"` // 当前探索度
	Kills       int    `json:"kills"`       // 在该地图的击杀数
}

// Monster 怪物
type Monster struct {
	ID              string  `json:"id"`
	ZoneID          string  `json:"zoneId"`
	Name            string  `json:"name"`
	Level           int     `json:"level"`
	Type            string  `json:"type"` // normal/elite/boss
	HP              int     `json:"hp"`
	MaxHP           int     `json:"maxHp"`
	PhysicalAttack  int     `json:"physicalAttack"`
	MagicAttack     int     `json:"magicAttack"`
	PhysicalDefense int     `json:"physicalDefense"`
	MagicDefense    int     `json:"magicDefense"`
	AttackType      string  `json:"attackType"`      // physical/magic
	PhysCritRate    float64 `json:"physCritRate"`    // 物理暴击率
	PhysCritDamage  float64 `json:"physCritDamage"`  // 物理暴击伤害
	SpellCritRate   float64 `json:"spellCritRate"`   // 法术暴击率
	SpellCritDamage float64 `json:"spellCritDamage"` // 法术暴击伤害
	DodgeRate       float64 `json:"dodgeRate"`       // 闪避率
	ExpReward       int     `json:"expReward"`
	GoldMin         int     `json:"goldMin"`
	GoldMax         int     `json:"goldMax"`
	SpawnWeight     int     `json:"spawnWeight"`
}

// ═══════════════════════════════════════════════════════════
// 战斗相关
// ═══════════════════════════════════════════════════════════

// BattleLog 战斗日志
type BattleLog struct {
	ID         int       `json:"id"`
	Message    string    `json:"message"`
	LogType    string    `json:"logType"` // combat/loot/system/levelup/damage/heal/buff
	Source     string    `json:"source,omitempty"`
	Target     string    `json:"target,omitempty"`
	Value      int       `json:"value,omitempty"`
	Color      string    `json:"color,omitempty"`
	DamageType string    `json:"damageType,omitempty"` // physical/magic
	CreatedAt  time.Time `json:"createdAt"`
}

// BattleStatus 战斗状态
type BattleStatus struct {
	IsRunning      bool         `json:"isRunning"`
	CurrentMonster *Monster     `json:"currentMonster,omitempty"`
	CurrentEnemies []*Monster   `json:"currentEnemies,omitempty"` // 多个敌人支持
	CurrentZoneID  string       `json:"currentZoneId,omitempty"`
	Team           []*Character `json:"team,omitempty"`
	BattleCount    int          `json:"battleCount"`
	TotalKills     int          `json:"totalKills"`
	TotalExp       int          `json:"totalExp"`
	TotalGold      int          `json:"totalGold"`
	SessionStart   *time.Time   `json:"sessionStart,omitempty"`
	IsResting      bool         `json:"isResting"`           // 是否在休息
	RestUntil      *time.Time   `json:"restUntil,omitempty"` // 休息结束时间
}

// ═══════════════════════════════════════════════════════════
// 战斗统计
// ═══════════════════════════════════════════════════════════

// BattleRecord 战斗记录 - 单场战斗的基本信息
type BattleRecord struct {
	ID              int       `json:"id"`
	UserID          int       `json:"userId"`
	ZoneID          string    `json:"zoneId"`
	BattleType      string    `json:"battleType"` // pve/pvp/boss/abyss
	MonsterID       string    `json:"monsterId,omitempty"`
	OpponentUserID  *int      `json:"opponentUserId,omitempty"`
	TotalRounds     int       `json:"totalRounds"`
	DurationSeconds int       `json:"durationSeconds"`
	Result          string    `json:"result"` // victory/defeat/draw/flee
	TeamDamageDealt int       `json:"teamDamageDealt"`
	TeamDamageTaken int       `json:"teamDamageTaken"`
	TeamHealingDone int       `json:"teamHealingDone"`
	ExpGained       int       `json:"expGained"`
	GoldGained      int       `json:"goldGained"`
	CreatedAt       time.Time `json:"createdAt"`

	// 关联数据（不存储在本表）
	CharacterStats []*BattleCharacterStats `json:"characterStats,omitempty"`
	SkillBreakdown []*BattleSkillBreakdown `json:"skillBreakdown,omitempty"`
}

// BattleCharacterStats 单场战斗角色统计 - 记录每场战斗中每个角色的详细表现
type BattleCharacterStats struct {
	ID          int `json:"id"`
	BattleID    int `json:"battleId"`
	CharacterID int `json:"characterId"`
	TeamSlot    int `json:"teamSlot"`

	// 伤害统计
	DamageDealt    int `json:"damageDealt"`    // 造成总伤害
	PhysicalDamage int `json:"physicalDamage"` // 物理伤害
	MagicDamage    int `json:"magicDamage"`    // 魔法伤害
	FireDamage     int `json:"fireDamage"`     // 火焰伤害
	FrostDamage    int `json:"frostDamage"`    // 冰霜伤害
	ShadowDamage   int `json:"shadowDamage"`   // 暗影伤害
	HolyDamage     int `json:"holyDamage"`     // 神圣伤害
	NatureDamage   int `json:"natureDamage"`   // 自然伤害
	DotDamage      int `json:"dotDamage"`      // DOT伤害

	// 暴击统计
	CritCount  int `json:"critCount"`  // 暴击次数
	CritDamage int `json:"critDamage"` // 暴击总伤害
	MaxCrit    int `json:"maxCrit"`    // 最高单次暴击

	// 承伤统计
	DamageTaken    int `json:"damageTaken"`    // 受到总伤害
	PhysicalTaken  int `json:"physicalTaken"`  // 物理承伤
	MagicTaken     int `json:"magicTaken"`     // 魔法承伤
	DamageBlocked  int `json:"damageBlocked"`  // 格挡伤害
	DamageAbsorbed int `json:"damageAbsorbed"` // 护盾吸收

	// 闪避统计
	DodgeCount int `json:"dodgeCount"` // 闪避次数
	BlockCount int `json:"blockCount"` // 格挡次数
	HitCount   int `json:"hitCount"`   // 被命中次数

	// 治疗统计
	HealingDone     int `json:"healingDone"`     // 造成治疗
	HealingReceived int `json:"healingReceived"` // 受到治疗
	Overhealing     int `json:"overhealing"`     // 过量治疗
	SelfHealing     int `json:"selfHealing"`     // 自我治疗
	HotHealing      int `json:"hotHealing"`      // HOT治疗

	// 技能统计
	SkillUses   int `json:"skillUses"`   // 技能使用次数
	SkillHits   int `json:"skillHits"`   // 技能命中次数
	SkillMisses int `json:"skillMisses"` // 技能未命中

	// 控制统计
	CcApplied  int `json:"ccApplied"`  // 施加控制次数
	CcReceived int `json:"ccReceived"` // 受到控制次数
	Dispels    int `json:"dispels"`    // 驱散次数
	Interrupts int `json:"interrupts"` // 打断次数

	// 其他统计
	Kills             int `json:"kills"`             // 击杀数(最后一击)
	Deaths            int `json:"deaths"`            // 死亡次数
	Resurrects        int `json:"resurrects"`        // 复活次数
	ResourceUsed      int `json:"resourceUsed"`      // 消耗能量
	ResourceGenerated int `json:"resourceGenerated"` // 获得能量
}

// CharacterLifetimeStats 角色生涯统计 - 累计统计数据
type CharacterLifetimeStats struct {
	CharacterID int `json:"characterId"`

	// 战斗场次
	TotalBattles int `json:"totalBattles"` // 总战斗场数
	Victories    int `json:"victories"`    // 胜利场数
	Defeats      int `json:"defeats"`      // 失败场数
	PveBattles   int `json:"pveBattles"`   // PVE战斗数
	PvpBattles   int `json:"pvpBattles"`   // PVP战斗数
	BossKills    int `json:"bossKills"`    // Boss击杀数

	// 累计伤害
	TotalDamageDealt    int `json:"totalDamageDealt"`    // 总造成伤害
	TotalPhysicalDamage int `json:"totalPhysicalDamage"` // 物理总伤害
	TotalMagicDamage    int `json:"totalMagicDamage"`    // 魔法总伤害
	TotalCritDamage     int `json:"totalCritDamage"`     // 暴击总伤害
	TotalCritCount      int `json:"totalCritCount"`      // 总暴击次数
	HighestDamageSingle int `json:"highestDamageSingle"` // 单次最高伤害
	HighestDamageBattle int `json:"highestDamageBattle"` // 单场最高伤害

	// 累计承伤
	TotalDamageTaken    int `json:"totalDamageTaken"`    // 总承受伤害
	TotalDamageBlocked  int `json:"totalDamageBlocked"`  // 总格挡伤害
	TotalDamageAbsorbed int `json:"totalDamageAbsorbed"` // 总吸收伤害
	TotalDodgeCount     int `json:"totalDodgeCount"`     // 总闪避次数

	// 累计治疗
	TotalHealingDone     int `json:"totalHealingDone"`     // 总治疗量
	TotalHealingReceived int `json:"totalHealingReceived"` // 总受到治疗
	TotalOverhealing     int `json:"totalOverhealing"`     // 总过量治疗
	HighestHealingSingle int `json:"highestHealingSingle"` // 单次最高治疗
	HighestHealingBattle int `json:"highestHealingBattle"` // 单场最高治疗

	// 击杀与死亡
	TotalKills      int `json:"totalKills"`        // 总击杀数
	TotalDeaths     int `json:"totalDeaths"`       // 总死亡数
	KillStreakBest  int `json:"killStreakBest"`    // 最长连杀
	CurrentKillStrk int `json:"currentKillStreak"` // 当前连杀

	// 技能使用
	TotalSkillUses int `json:"totalSkillUses"` // 技能总使用次数
	TotalSkillHits int `json:"totalSkillHits"` // 技能总命中数

	// 资源统计
	TotalResourceUsed int `json:"totalResourceUsed"` // 总消耗能量
	TotalRounds       int `json:"totalRounds"`       // 总战斗回合数
	TotalBattleTime   int `json:"totalBattleTime"`   // 总战斗时间(秒)

	// 最后更新
	LastBattleAt *time.Time `json:"lastBattleAt,omitempty"`
	UpdatedAt    time.Time  `json:"updatedAt"`
}

// BattleSkillBreakdown 战斗技能明细 - 每场战斗中各技能的使用和效果
type BattleSkillBreakdown struct {
	ID           int    `json:"id"`
	BattleID     int    `json:"battleId"`
	CharacterID  int    `json:"characterId"`
	SkillID      string `json:"skillId"`
	UseCount     int    `json:"useCount"`     // 使用次数
	HitCount     int    `json:"hitCount"`     // 命中次数
	CritCount    int    `json:"critCount"`    // 暴击次数
	TotalDamage  int    `json:"totalDamage"`  // 造成总伤害
	TotalHealing int    `json:"totalHealing"` // 造成总治疗
	ResourceCost int    `json:"resourceCost"` // 总消耗能量

	// 关联数据（不存储在本表）
	SkillName string `json:"skillName,omitempty"`
}

// DailyStatistics 每日统计汇总 - 每日战斗数据快照
type DailyStatistics struct {
	ID               int       `json:"id"`
	UserID           int       `json:"userId"`
	StatDate         string    `json:"statDate"` // YYYY-MM-DD 格式
	BattlesCount     int       `json:"battlesCount"`
	Victories        int       `json:"victories"`
	Defeats          int       `json:"defeats"`
	TotalDamage      int       `json:"totalDamage"`
	TotalHealing     int       `json:"totalHealing"`
	TotalDamageTaken int       `json:"totalDamageTaken"`
	ExpGained        int       `json:"expGained"`
	GoldGained       int       `json:"goldGained"`
	PlayTime         int       `json:"playTime"` // 游戏时长(秒)
	Kills            int       `json:"kills"`
	Deaths           int       `json:"deaths"`
	CreatedAt        time.Time `json:"createdAt,omitempty"`
}

// ═══════════════════════════════════════════════════════════
// 战斗统计 API 响应
// ═══════════════════════════════════════════════════════════

// BattleStatsOverview 战斗统计概览 - 用于前端面板展示
type BattleStatsOverview struct {
	// 会话统计
	SessionStats *SessionStats `json:"sessionStats,omitempty"`

	// 角色生涯统计
	LifetimeStats []*CharacterLifetimeStats `json:"lifetimeStats,omitempty"`

	// 今日统计
	TodayStats *DailyStatistics `json:"todayStats,omitempty"`

	// 最近战斗
	RecentBattles []*BattleRecord `json:"recentBattles,omitempty"`
}

// SessionStats 会话统计 - 当前游戏会话的统计数据
type SessionStats struct {
	TotalBattles    int       `json:"totalBattles"`
	TotalKills      int       `json:"totalKills"`
	TotalExp        int       `json:"totalExp"`
	TotalGold       int       `json:"totalGold"`
	TotalDamage     int       `json:"totalDamage"`
	TotalHealing    int       `json:"totalHealing"`
	SessionStart    time.Time `json:"sessionStart"`
	DurationSeconds int       `json:"durationSeconds"`
}

// CharacterBattleSummary 角色战斗摘要 - 用于角色详情页展示
type CharacterBattleSummary struct {
	CharacterID   int     `json:"characterId"`
	CharacterName string  `json:"characterName"`
	TotalBattles  int     `json:"totalBattles"`
	Victories     int     `json:"victories"`
	WinRate       float64 `json:"winRate"`
	TotalDamage   int     `json:"totalDamage"`
	TotalHealing  int     `json:"totalHealing"`
	TotalKills    int     `json:"totalKills"`
	TotalDeaths   int     `json:"totalDeaths"`
	KDRatio       float64 `json:"kdRatio"` // 击杀/死亡比
	AvgDPS        float64 `json:"avgDps"`  // 平均每回合伤害
	AvgHPS        float64 `json:"avgHps"`  // 平均每回合治疗
}

// ═══════════════════════════════════════════════════════════
// DPS分析相关
// ═══════════════════════════════════════════════════════════

// SkillDPSAnalysis 技能DPS分析 - 单个技能的详细DPS数据
type SkillDPSAnalysis struct {
	SkillID           string  `json:"skillId"`
	SkillName         string  `json:"skillName"`
	TotalDamage       int     `json:"totalDamage"`       // 总伤害
	UseCount          int     `json:"useCount"`          // 使用次数
	HitCount          int     `json:"hitCount"`          // 命中次数
	CritCount         int     `json:"critCount"`         // 暴击次数
	AvgDamage         float64 `json:"avgDamage"`         // 平均伤害
	MaxDamage         int     `json:"maxDamage"`         // 最高伤害
	DPS               float64 `json:"dps"`               // 每秒伤害
	DamagePercent     float64 `json:"damagePercent"`     // 伤害占比(%)
	ResourceCost      int     `json:"resourceCost"`      // 总消耗能量
	DamagePerResource float64 `json:"damagePerResource"` // 每点能量伤害
	HitRate           float64 `json:"hitRate"`           // 命中率(%)
	CritRate          float64 `json:"critRate"`          // 暴击率(%)
}

// CharacterDPSAnalysis 角色DPS分析 - 单个角色的完整DPS数据
type CharacterDPSAnalysis struct {
	CharacterID       int                 `json:"characterId"`
	CharacterName     string              `json:"characterName"`
	TotalDamage       int                 `json:"totalDamage"`       // 总伤害
	TotalHealing      int                 `json:"totalHealing"`      // 总治疗
	Duration          int                 `json:"duration"`          // 战斗时长(秒)
	TotalDPS          float64             `json:"totalDps"`          // 总DPS
	TotalHPS          float64             `json:"totalHps"`          // 总HPS
	SkillBreakdown    []*SkillDPSAnalysis `json:"skillBreakdown"`    // 技能明细
	DamageComposition *DamageComposition  `json:"damageComposition"` // 伤害构成
}

// DamageComposition 伤害构成 - 按类型分类的伤害统计
type DamageComposition struct {
	Physical    int                `json:"physical"`    // 物理伤害
	Magic       int                `json:"magic"`       // 魔法伤害
	Fire        int                `json:"fire"`        // 火焰伤害
	Frost       int                `json:"frost"`       // 冰霜伤害
	Shadow      int                `json:"shadow"`      // 暗影伤害
	Holy        int                `json:"holy"`        // 神圣伤害
	Nature      int                `json:"nature"`      // 自然伤害
	Dot         int                `json:"dot"`         // DOT伤害
	Total       int                `json:"total"`       // 总伤害
	Percentages map[string]float64 `json:"percentages"` // 各类型占比(%)
}

// BattleDPSAnalysis 战斗DPS分析 - 单场战斗的完整DPS数据
type BattleDPSAnalysis struct {
	BattleID              int                     `json:"battleId"`
	Duration              int                     `json:"duration"`              // 战斗时长(秒)
	TotalRounds           int                     `json:"totalRounds"`           // 总回合数
	BattleCount           int                     `json:"battleCount,omitempty"` // 战斗场次（累计统计时使用）
	TeamDPS               float64                 `json:"teamDps"`               // 队伍总DPS
	TeamHPS               float64                 `json:"teamHps"`               // 队伍总HPS
	Characters            []*CharacterDPSAnalysis `json:"characters"`            // 各角色DPS分析
	TeamDamageComposition *DamageComposition      `json:"teamDamageComposition"` // 队伍伤害构成
}

// ═══════════════════════════════════════════════════════════
// 技能
// ═══════════════════════════════════════════════════════════

// Skill 技能
type Skill struct {
	ID             string  `json:"id"`
	Name           string  `json:"name"`
	Description    string  `json:"description"`
	ClassID        string  `json:"classId"`
	Type           string  `json:"type"`       // attack/heal/buff/debuff/dot/hot/shield/control
	TargetType     string  `json:"targetType"` // self/ally/enemy/ally_all/enemy_all
	DamageType     string  `json:"damageType,omitempty"`
	BaseValue      int     `json:"baseValue"`
	ScalingStat    string  `json:"scalingStat,omitempty"`
	ScalingRatio   float64 `json:"scalingRatio"`
	ResourceCost   int     `json:"resourceCost"`
	Cooldown       int     `json:"cooldown"`
	LevelRequired  int     `json:"levelRequired"`
	ThreatModifier float64 `json:"threatModifier"`
	ThreatType     string  `json:"threatType"`     // normal/high/taunt/reduce/clear
	Tags           string  `json:"tags,omitempty"` // JSON数组字符串
}

// CharacterSkill 角色技能（已学会的技能）
type CharacterSkill struct {
	ID          int    `json:"id"`
	CharacterID int    `json:"characterId"`
	SkillID     string `json:"skillId"`
	SkillLevel  int    `json:"skillLevel"`
	Slot        *int   `json:"slot,omitempty"`
	IsAuto      bool   `json:"isAuto"`
	Skill       *Skill `json:"skill,omitempty"` // 关联的技能详情
}

// PassiveSkill 被动技能
type PassiveSkill struct {
	ID           string  `json:"id"`
	Name         string  `json:"name"`
	Description  string  `json:"description"`
	ClassID      string  `json:"classId"`
	Rarity       string  `json:"rarity"` // common/rare/epic/legendary
	Tier         int     `json:"tier"`
	EffectType   string  `json:"effectType"`
	EffectValue  float64 `json:"effectValue"`
	EffectStat   string  `json:"effectStat,omitempty"`
	MaxLevel     int     `json:"maxLevel"`
	LevelScaling float64 `json:"levelScaling"`
}

// CharacterPassiveSkill 角色被动技能
type CharacterPassiveSkill struct {
	ID          int           `json:"id"`
	CharacterID int           `json:"characterId"`
	PassiveID   string        `json:"passiveId"`
	Level       int           `json:"level"`
	AcquiredAt  time.Time     `json:"acquiredAt"`
	Passive     *PassiveSkill `json:"passive,omitempty"` // 关联的被动技能详情
}

// SkillSelection 技能选择机会
type SkillSelection struct {
	CharacterID     int                      `json:"characterId"`
	Level           int                      `json:"level"`
	SelectionType   string                   `json:"selectionType"`             // "initial_active" / "active" / "passive"
	CanUpgrade      bool                     `json:"canUpgrade"`                // 是否可以升级现有技能
	UpgradeSkills   []*CharacterSkill        `json:"upgradeSkills,omitempty"`   // 可升级的主动技能列表
	UpgradePassives []*CharacterPassiveSkill `json:"upgradePassives,omitempty"` // 可升级的被动技能列表
	NewSkills       []*Skill                 `json:"newSkills,omitempty"`       // 新技能选项（随机4个）
	NewPassives     []*PassiveSkill          `json:"newPassives,omitempty"`     // 新被动技能选项（随机4个）
}

// SkillSelectionRequest 技能选择请求
type SkillSelectionRequest struct {
	CharacterID int    `json:"characterId" binding:"required"`
	SkillID     string `json:"skillId,omitempty"`   // 选择的技能ID（新技能或升级）
	PassiveID   string `json:"passiveId,omitempty"` // 选择的被动技能ID（新被动或升级）
	IsUpgrade   bool   `json:"isUpgrade"`           // 是否为升级
}

// InitialSkillSelectionRequest 初始技能选择请求
type InitialSkillSelectionRequest struct {
	CharacterID int      `json:"characterId" binding:"required"`
	SkillIDs    []string `json:"skillIds" binding:"required,len=2"` // 必须选择2个技能
}

// ═══════════════════════════════════════════════════════════
// 战斗策略
// ═══════════════════════════════════════════════════════════

// BattleStrategy 战斗策略
type BattleStrategy struct {
	ID                   int                `json:"id"`
	CharacterID          int                `json:"characterId"`
	Name                 string             `json:"name"`
	IsActive             bool               `json:"isActive"`
	SkillPriority        []string           `json:"skillPriority"`        // 技能优先级列表
	ConditionalRules     []ConditionalRule  `json:"conditionalRules"`     // 条件规则
	TargetPriority       string             `json:"targetPriority"`       // 默认目标选择策略
	SkillTargetOverrides map[string]string  `json:"skillTargetOverrides"` // 技能目标覆盖
	ResourceThreshold    int                `json:"resourceThreshold"`    // 资源阈值
	ReservedSkills       []ReservedSkill    `json:"reservedSkills"`       // 保留技能
	AutoTargetSettings   AutoTargetSettings `json:"autoTargetSettings"`   // 智能目标设置
	CreatedAt            time.Time          `json:"createdAt"`
	UpdatedAt            *time.Time         `json:"updatedAt,omitempty"`
}

// ConditionalRule 条件规则
type ConditionalRule struct {
	ID        string        `json:"id"`
	Priority  int           `json:"priority"`
	Enabled   bool          `json:"enabled"`
	Condition RuleCondition `json:"condition"`
	Action    RuleAction    `json:"action"`
}

// RuleCondition 规则条件
type RuleCondition struct {
	Type     string  `json:"type"`              // 条件类型: self_hp_percent, alive_enemy_count, target_hp_percent, etc.
	Operator string  `json:"operator"`          // 比较运算符: <, >, <=, >=, =, !=
	Value    float64 `json:"value"`             // 条件值
	SkillID  string  `json:"skillId,omitempty"` // 技能ID (用于 skill_ready 条件)
	BuffID   string  `json:"buffId,omitempty"`  // Buff ID (用于 has_buff 条件)
}

// RuleAction 规则动作
type RuleAction struct {
	Type    string `json:"type"`              // 动作类型: use_skill, normal_attack
	SkillID string `json:"skillId,omitempty"` // 使用的技能ID
	Comment string `json:"comment,omitempty"` // 备注
}

// ReservedSkill 保留技能
type ReservedSkill struct {
	SkillID   string        `json:"skillId"`
	Condition RuleCondition `json:"condition"`
}

// AutoTargetSettings 智能目标设置
type AutoTargetSettings struct {
	PositionalAutoOptimize bool `json:"positionalAutoOptimize"` // 位置技能自动优化
	ExecuteAutoTarget      bool `json:"executeAutoTarget"`      // 斩杀技能自动选择低血量
	HealAutoTarget         bool `json:"healAutoTarget"`         // 治疗技能自动选择低血量队友
}

// StrategyCreateRequest 创建策略请求
type StrategyCreateRequest struct {
	CharacterID  int    `json:"characterId" binding:"required"`
	Name         string `json:"name" binding:"required,min=1,max=32"`
	FromTemplate string `json:"fromTemplate,omitempty"` // 从模板创建
}

// StrategyUpdateRequest 更新策略请求
type StrategyUpdateRequest struct {
	Name                 *string             `json:"name,omitempty"`
	IsActive             *bool               `json:"isActive,omitempty"`
	SkillPriority        []string            `json:"skillPriority,omitempty"`
	ConditionalRules     []ConditionalRule   `json:"conditionalRules,omitempty"`
	TargetPriority       *string             `json:"targetPriority,omitempty"`
	SkillTargetOverrides map[string]string   `json:"skillTargetOverrides,omitempty"`
	ResourceThreshold    *int                `json:"resourceThreshold,omitempty"`
	ReservedSkills       []ReservedSkill     `json:"reservedSkills,omitempty"`
	AutoTargetSettings   *AutoTargetSettings `json:"autoTargetSettings,omitempty"`
}

// ═══════════════════════════════════════════════════════════
// API 响应
// ═══════════════════════════════════════════════════════════

// APIResponse 通用API响应
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

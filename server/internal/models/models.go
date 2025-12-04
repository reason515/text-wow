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
	ID           int        `json:"id"`
	UserID       int        `json:"userId"`
	Name         string     `json:"name"`
	RaceID       string     `json:"raceId"`
	ClassID      string     `json:"classId"`
	Faction      string     `json:"faction"`  // alliance, horde
	TeamSlot     int        `json:"teamSlot"` // 1-5
	IsActive     bool       `json:"isActive"`
	IsDead       bool       `json:"isDead"`
	ReviveAt     *time.Time `json:"reviveAt,omitempty"`
	Level        int        `json:"level"`
	Exp          int        `json:"exp"`
	ExpToNext    int        `json:"expToNext"`
	HP           int        `json:"hp"`
	MaxHP        int        `json:"maxHp"`
	Resource     int        `json:"resource"`     // 当前能量
	MaxResource  int        `json:"maxResource"`  // 最大能量
	ResourceType string     `json:"resourceType"` // mana/rage/energy
	Strength     int        `json:"strength"`
	Agility      int        `json:"agility"`
	Intellect    int        `json:"intellect"`
	Stamina      int        `json:"stamina"`
	Spirit       int        `json:"spirit"`
	Attack       int        `json:"attack"`
	Defense      int        `json:"defense"`
	CritRate     float64    `json:"critRate"`
	CritDamage   float64    `json:"critDamage"`
	TotalKills   int        `json:"totalKills"`
	TotalDeaths  int        `json:"totalDeaths"`
	CreatedAt    time.Time  `json:"createdAt"`
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
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	MinLevel    int       `json:"minLevel"`
	MaxLevel    int       `json:"maxLevel"`
	Faction     string    `json:"faction"` // alliance/horde/neutral
	ExpMulti    float64   `json:"expMulti"`
	GoldMulti   float64   `json:"goldMulti"`
	Monsters    []Monster `json:"monsters,omitempty"`
}

// Monster 怪物
type Monster struct {
	ID          string `json:"id"`
	ZoneID      string `json:"zoneId"`
	Name        string `json:"name"`
	Level       int    `json:"level"`
	Type        string `json:"type"` // normal/elite/boss
	HP          int    `json:"hp"`
	MaxHP       int    `json:"maxHp"`
	Attack      int    `json:"attack"`
	Defense     int    `json:"defense"`
	ExpReward   int    `json:"expReward"`
	GoldMin     int    `json:"goldMin"`
	GoldMax     int    `json:"goldMax"`
	SpawnWeight int    `json:"spawnWeight"`
}

// ═══════════════════════════════════════════════════════════
// 战斗相关
// ═══════════════════════════════════════════════════════════

// BattleLog 战斗日志
type BattleLog struct {
	ID        int       `json:"id"`
	Message   string    `json:"message"`
	LogType   string    `json:"logType"` // combat/loot/system/levelup/damage/heal/buff
	Source    string    `json:"source,omitempty"`
	Target    string    `json:"target,omitempty"`
	Value     int       `json:"value,omitempty"`
	CreatedAt time.Time `json:"createdAt"`
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
// API 响应
// ═══════════════════════════════════════════════════════════

// APIResponse 通用API响应
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

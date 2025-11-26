// 游戏核心类型定义

// 阵营
export type Faction = 'alliance' | 'horde'

// 种族
export type Race = 
  | 'human' | 'dwarf' | 'nightelf' | 'gnome' | 'draenei'  // 联盟
  | 'orc' | 'troll' | 'undead' | 'tauren' | 'bloodelf'    // 部落

// 职业
export type CharacterClass = 
  | 'warrior' | 'paladin' | 'hunter' | 'rogue' 
  | 'priest' | 'mage' | 'warlock' | 'druid' | 'shaman'

// 属性
export interface Stats {
  strength: number      // 力量 - 物理攻击
  agility: number       // 敏捷 - 暴击/闪避
  intellect: number     // 智力 - 法术攻击/法力值
  stamina: number       // 耐力 - 生命值
  spirit: number        // 精神 - 法力回复/生命回复
}

// 战斗属性（计算后）
export interface CombatStats {
  maxHp: number
  currentHp: number
  maxMp: number
  currentMp: number
  attack: number
  defense: number
  critRate: number      // 暴击率 0-100
  dodgeRate: number     // 闪避率 0-100
}

// 技能
export interface Skill {
  id: string
  name: string
  description: string
  damage: number        // 伤害倍率
  mpCost: number        // 法力消耗
  cooldown: number      // 冷却回合
  currentCooldown: number
  type: 'physical' | 'magical' | 'heal'
  icon?: string
}

// 装备槽位
export type EquipSlot = 'weapon' | 'head' | 'chest' | 'legs' | 'feet' | 'hands' | 'trinket'

// 装备品质
export type ItemQuality = 'common' | 'uncommon' | 'rare' | 'epic' | 'legendary'

// 装备
export interface Equipment {
  id: string
  name: string
  slot: EquipSlot
  quality: ItemQuality
  level: number
  stats: Partial<Stats>
  requiredLevel: number
}

// 背包物品
export interface InventoryItem {
  id: string
  name: string
  type: 'equipment' | 'consumable' | 'material' | 'quest'
  quantity: number
  data?: Equipment
}

// 玩家角色
export interface Character {
  id: string
  name: string
  faction: Faction
  race: Race
  class: CharacterClass
  level: number
  exp: number
  expToNextLevel: number
  gold: number
  stats: Stats
  combatStats: CombatStats
  skills: Skill[]
  equipment: Partial<Record<EquipSlot, Equipment>>
  inventory: InventoryItem[]
}

// 怪物
export interface Monster {
  id: string
  name: string
  level: number
  maxHp: number
  currentHp: number
  attack: number
  defense: number
  expReward: number
  goldReward: [number, number]  // [min, max]
  lootTable: LootItem[]
}

// 掉落物
export interface LootItem {
  itemId: string
  name: string
  dropRate: number  // 0-100
  quantity: [number, number]
}

// 战斗区域
export interface Zone {
  id: string
  name: string
  description: string
  levelRange: [number, number]
  monsters: string[]  // monster ids
  unlockLevel: number
}

// 战斗日志条目
export interface LogEntry {
  id: number
  timestamp: number
  message: string
  type: 'system' | 'combat-start' | 'damage' | 'heal' | 'loot' | 'exp' | 'levelup' | 'death'
}

// 战斗策略
export interface BattleStrategy {
  skillPriority: string[]           // 技能使用优先级
  useHealAt: number                 // HP低于多少%时使用治疗
  targetPriority: 'lowest_hp' | 'highest_hp' | 'random'
  autoPotionAt: number              // HP低于多少%时使用药水
}

// 游戏状态
export interface GameState {
  character: Character | null
  currentZone: Zone | null
  isInCombat: boolean
  currentMonster: Monster | null
  battleSpeed: 1 | 2 | 5 | 10
  isPaused: boolean
  logs: LogEntry[]
  strategy: BattleStrategy
  statistics: GameStatistics
}

// 游戏统计
export interface GameStatistics {
  totalKills: number
  totalExp: number
  totalGold: number
  totalPlayTime: number  // seconds
  highestDamage: number
  deathCount: number
}

// 种族信息
export interface RaceInfo {
  id: Race
  name: string
  faction: Faction
  description: string
  bonusStats: Partial<Stats>
}

// 职业信息
export interface ClassInfo {
  id: CharacterClass
  name: string
  description: string
  primaryStat: keyof Stats
  startingSkills: string[]
  availableRaces: Race[]
}


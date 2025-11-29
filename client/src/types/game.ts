// ═══════════════════════════════════════════════════════════
// 用户相关
// ═══════════════════════════════════════════════════════════

export interface User {
  id: number
  username: string
  email?: string
  maxTeamSize: number
  unlockedSlots: number
  gold: number
  currentZoneId: string
  totalKills: number
  totalGoldGained: number
  playTime: number
  createdAt: string
  lastLoginAt?: string
}

export interface AuthResponse {
  token: string
  user: User
}

export interface LoginRequest {
  username: string
  password: string
}

export interface RegisterRequest {
  username: string
  password: string
  email?: string
}

// ═══════════════════════════════════════════════════════════
// 角色相关
// ═══════════════════════════════════════════════════════════

export interface Character {
  id: number
  userId: number
  name: string
  raceId: string
  classId: string
  faction: string
  teamSlot: number
  isActive: boolean
  isDead: boolean
  reviveAt?: string
  level: number
  exp: number
  expToNext: number
  hp: number
  maxHp: number
  resource: number
  maxResource: number
  resourceType: string
  strength: number
  agility: number
  intellect: number
  stamina: number
  spirit: number
  attack: number
  defense: number
  critRate: number
  critDamage: number
  totalKills: number
  totalDeaths: number
  createdAt: string
}

export interface CharacterCreate {
  name: string
  raceId: string
  classId: string
}

export interface Team {
  userId: number
  maxSize: number
  unlockedSlots: number
  characters: Character[]
}

// ═══════════════════════════════════════════════════════════
// 种族和职业
// ═══════════════════════════════════════════════════════════

export interface Race {
  id: string
  name: string
  faction: string
  description: string
  baseStrengthBonus: number
  baseAgilityBonus: number
  baseIntellectBonus: number
  baseStaminaBonus: number
  baseSpiritBonus: number
}

export interface Class {
  id: string
  name: string
  description: string
  role: string
  primaryStat: string
  resourceType: string
  baseHp: number
  baseResource: number
  combatRole: string
  isRanged: boolean
}

// ═══════════════════════════════════════════════════════════
// 区域和怪物
// ═══════════════════════════════════════════════════════════

export interface Zone {
  id: string
  name: string
  description: string
  minLevel: number
  maxLevel: number
  faction: string
  expMulti: number
  goldMulti: number
}

export interface Monster {
  id: string
  zoneId: string
  name: string
  level: number
  type: string
  hp: number
  maxHp: number
  attack: number
  defense: number
  expReward: number
  goldMin: number
  goldMax: number
}

// ═══════════════════════════════════════════════════════════
// 战斗相关
// ═══════════════════════════════════════════════════════════

export interface BattleLog {
  id?: number
  message: string
  logType: string
  source?: string
  target?: string
  value?: number
  createdAt: string
}

export interface BattleStatus {
  isRunning: boolean
  currentMonster?: Monster
  team: Character[]
  battleCount: number
  totalKills: number
  totalExp: number
  totalGold: number
  sessionStart?: string
}

// ═══════════════════════════════════════════════════════════
// API 响应
// ═══════════════════════════════════════════════════════════

export interface APIResponse<T = any> {
  success: boolean
  message?: string
  data?: T
  error?: string
}

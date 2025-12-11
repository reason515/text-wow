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
  unspentPoints?: number
  physicalAttack: number
  magicAttack: number
  physicalDefense: number
  magicDefense: number
  physCritRate: number      // 物理暴击率
  physCritDamage: number    // 物理暴击伤害
  spellCritRate: number     // 法术暴击率
  spellCritDamage: number   // 法术暴击伤害
  dodgeRate: number         // 闪避率
  totalKills: number
  totalDeaths: number
  createdAt: string
  buffs?: BuffInfo[]
}

export interface BuffInfo {
  effectId: string
  name: string
  type: string
  isBuff: boolean
  duration: number
  value: number
  statAffected: string
  description?: string
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
  strengthBase: number
  agilityBase: number
  intellectBase: number
  staminaBase: number
  spiritBase: number
  strengthPct: number
  agilityPct: number
  intellectPct: number
  staminaPct: number
  spiritPct: number
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
  physicalAttack: number
  magicAttack: number
  physicalDefense: number
  magicDefense: number
  attackType: string         // physical/magic
  physCritRate: number       // 物理暴击率
  physCritDamage: number     // 物理暴击伤害
  spellCritRate: number      // 法术暴击率
  spellCritDamage: number    // 法术暴击伤害
  dodgeRate: number         // 闪避率
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
  color?: string
  damageType?: 'physical' | 'magic'
  actor?: string
  action?: string
  time?: string
  type?: string
  createdAt?: string
}

export interface BattleStatus {
  isRunning: boolean
  currentMonster?: Monster
  currentEnemies?: Monster[]  // 多个敌人支持
  currentZoneId?: string
  team?: Character[]
  battleCount: number
  totalKills: number
  totalExp: number
  totalGold: number
  sessionStart?: string
}

export interface BattleResult {
  character?: Character
  enemy?: Monster
  enemies?: Monster[]  // 多个敌人支持
  logs?: BattleLog[]
  status?: BattleStatus
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

// ═══════════════════════════════════════════════════════════
// 职业颜色映射 (魔兽世界经典配色)
// ═══════════════════════════════════════════════════════════

export const CLASS_COLORS: Record<string, string> = {
  warrior: '#C79C6E',   // 战士 - 棕褐色
  paladin: '#F58CBA',   // 圣骑士 - 粉色
  hunter: '#ABD473',    // 猎人 - 草绿色
  rogue: '#FFF569',     // 盗贼 - 黄色
  priest: '#FFFFFF',    // 牧师 - 白色
  mage: '#69CCF0',      // 法师 - 天蓝色
  warlock: '#9482C9',   // 术士 - 紫色
  druid: '#FF7D0A',     // 德鲁伊 - 橙色
  shaman: '#0070DE',    // 萨满 - 蓝色
}

export const CLASS_NAMES: Record<string, string> = {
  warrior: '战士',
  paladin: '圣骑士',
  hunter: '猎人',
  rogue: '盗贼',
  priest: '牧师',
  mage: '法师',
  warlock: '术士',
  druid: '德鲁伊',
  shaman: '萨满',
}

// 获取职业CSS类名
export function getClassColorClass(classId: string): string {
  return `class-${classId}`
}

// 获取职业颜色值
export function getClassColor(classId: string): string {
  return CLASS_COLORS[classId] || '#33ff33'
}

// ═══════════════════════════════════════════════════════════
// 物品品质颜色
// ═══════════════════════════════════════════════════════════

export const QUALITY_COLORS: Record<string, string> = {
  common: '#9d9d9d',    // 普通 - 灰色
  uncommon: '#1eff00',  // 优秀 - 绿色
  rare: '#0070dd',      // 稀有 - 蓝色
  epic: '#a335ee',      // 史诗 - 紫色
  legendary: '#ff8000', // 传说 - 橙色
  mythic: '#e6cc80',    // 神话 - 金色
}

export const QUALITY_NAMES: Record<string, string> = {
  common: '普通',
  uncommon: '优秀',
  rare: '稀有',
  epic: '史诗',
  legendary: '传说',
  mythic: '神话',
}

// 获取品质CSS类名
export function getQualityColorClass(quality: string): string {
  return `quality-${quality}`
}

// 获取品质颜色值
export function getQualityColor(quality: string): string {
  return QUALITY_COLORS[quality] || '#9d9d9d'
}

// ═══════════════════════════════════════════════════════════
// 资源颜色映射 (参考魔兽世界)
// ═══════════════════════════════════════════════════════════

export const RESOURCE_COLORS: Record<string, string> = {
  rage: '#ff4444',    // 红色 - 怒气
  mana: '#3d85c6',    // 蓝色 - 法力
  energy: '#ffd700',  // 金色/黄色 - 能量
}

// 获取资源颜色值
export function getResourceColor(resourceType: string): string {
  return RESOURCE_COLORS[resourceType] || '#ffffff'
}

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
  unlockZoneId?: string | null
  requiredExploration: number
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

// ═══════════════════════════════════════════════════════════
// 战斗策略
// ═══════════════════════════════════════════════════════════

export interface BattleStrategy {
  id: number
  characterId: number
  name: string
  isActive: boolean
  skillPriority: string[]
  conditionalRules: ConditionalRule[]
  targetPriority: string
  skillTargetOverrides: Record<string, string>
  resourceThreshold: number
  reservedSkills: ReservedSkill[]
  autoTargetSettings: AutoTargetSettings
  createdAt: string
  updatedAt?: string
}

export interface ConditionalRule {
  id: string
  priority: number
  enabled: boolean
  condition: RuleCondition
  action: RuleAction
}

export interface RuleCondition {
  type: string
  operator: string
  value: number
  skillId?: string
  buffId?: string
}

export interface RuleAction {
  type: string
  skillId?: string
  comment?: string
}

export interface ReservedSkill {
  skillId: string
  condition: RuleCondition
}

export interface AutoTargetSettings {
  positionalAutoOptimize: boolean
  executeAutoTarget: boolean
  healAutoTarget: boolean
}

export interface StrategyCreateRequest {
  characterId: number
  name: string
  fromTemplate?: string
}

export interface StrategyUpdateRequest {
  name?: string
  isActive?: boolean
  skillPriority?: string[]
  conditionalRules?: ConditionalRule[]
  targetPriority?: string
  skillTargetOverrides?: Record<string, string>
  resourceThreshold?: number
  reservedSkills?: ReservedSkill[]
  autoTargetSettings?: AutoTargetSettings
}

export interface ConditionTypeInfo {
  type: string
  name: string
  category: string
  operators: string[]
  valueType: string
}

export interface TargetPriorityInfo {
  value: string
  label: string
}

export interface StrategyTemplate {
  id: string
  name: string
  description: string
}

// ═══════════════════════════════════════════════════════════
// 战斗统计
// ═══════════════════════════════════════════════════════════

export interface BattleRecord {
  id: number
  userId: number
  zoneId: string
  battleType: string
  monsterId?: string
  opponentUserId?: number
  totalRounds: number
  durationSeconds: number
  result: string
  teamDamageDealt: number
  teamDamageTaken: number
  teamHealingDone: number
  expGained: number
  goldGained: number
  createdAt: string
  characterStats?: BattleCharacterStats[]
  skillBreakdown?: BattleSkillBreakdown[]
}

export interface BattleCharacterStats {
  id: number
  battleId: number
  characterId: number
  teamSlot: number
  damageDealt: number
  physicalDamage: number
  magicDamage: number
  fireDamage: number
  frostDamage: number
  shadowDamage: number
  holyDamage: number
  natureDamage: number
  dotDamage: number
  critCount: number
  critDamage: number
  maxCrit: number
  damageTaken: number
  physicalDamageTaken: number
  magicDamageTaken: number
  damageBlocked: number
  damageAbsorbed: number
  dodgeCount: number
  blockCount: number
  hitCount: number
  healingDone: number
  healingReceived: number
  overhealing: number
  selfHealing: number
  hotHealing: number
  skillUses: number
  skillHits: number
  skillMisses: number
  ccApplied: number
  ccReceived: number
  dispels: number
  interrupts: number
  kills: number
  deaths: number
  resurrects: number
  resourceUsed: number
  resourceGenerated: number
}

export interface CharacterLifetimeStats {
  characterId: number
  totalBattles: number
  victories: number
  defeats: number
  pveBattles: number
  pvpBattles: number
  bossKills: number
  totalDamageDealt: number
  totalPhysicalDamage: number
  totalMagicDamage: number
  totalCritDamage: number
  totalCritCount: number
  highestDamageSingle: number
  highestDamageBattle: number
  totalDamageTaken: number
  totalDamageBlocked: number
  totalDamageAbsorbed: number
  totalDodgeCount: number
  totalHealingDone: number
  totalHealingReceived: number
  totalOverhealing: number
  highestHealingSingle: number
  highestHealingBattle: number
  totalKills: number
  totalDeaths: number
  killStreakBest: number
  currentKillStreak: number
  totalSkillUses: number
  totalSkillHits: number
  totalResourceUsed: number
  totalRounds: number
  totalBattleTime: number
  lastBattleAt?: string
  updatedAt: string
}

export interface BattleSkillBreakdown {
  id: number
  battleId: number
  characterId: number
  skillId: string
  useCount: number
  hitCount: number
  critCount: number
  totalDamage: number
  totalHealing: number
  resourceCost: number
  skillName?: string
}

export interface DailyStatistics {
  id: number
  userId: number
  statDate: string
  battlesCount: number
  victories: number
  defeats: number
  totalDamage: number
  totalHealing: number
  totalDamageTaken: number
  expGained: number
  goldGained: number
  playTime: number
  kills: number
  deaths: number
  createdAt?: string
}

export interface SessionStats {
  totalBattles: number
  totalKills: number
  totalExp: number
  totalGold: number
  totalDamage: number
  totalHealing: number
  sessionStart: string
  durationSeconds: number
}

export interface BattleStatsOverview {
  sessionStats?: SessionStats
  lifetimeStats?: CharacterLifetimeStats[]
  todayStats?: DailyStatistics
  recentBattles?: BattleRecord[]
}

export interface CharacterBattleSummary {
  characterId: number
  characterName: string
  totalBattles: number
  victories: number
  winRate: number
  totalDamage: number
  totalHealing: number
  totalKills: number
  totalDeaths: number
  kdRatio: number
  avgDps: number
  avgHps: number
}

// ═══════════════════════════════════════════════════════════
// DPS分析
// ═══════════════════════════════════════════════════════════

export interface SkillDPSAnalysis {
  skillId: string
  skillName: string
  totalDamage: number
  useCount: number
  hitCount: number
  critCount: number
  avgDamage: number
  maxDamage: number
  dps: number
  damagePercent: number
  resourceCost: number
  damagePerResource: number
  hitRate: number
  critRate: number
}

export interface DamageComposition {
  physical: number
  magic: number
  fire: number
  frost: number
  shadow: number
  holy: number
  nature: number
  dot: number
  total: number
  percentages: Record<string, number>
}

export interface CharacterDPSAnalysis {
  characterId: number
  characterName: string
  totalDamage: number
  totalHealing: number
  duration: number
  totalDps: number
  totalHps: number
  skillBreakdown: SkillDPSAnalysis[]
  damageComposition: DamageComposition
}

export interface BattleDPSAnalysis {
  battleId: number
  duration: number
  totalRounds: number
  battleCount?: number
  teamDps: number
  teamHps: number
  characters: CharacterDPSAnalysis[]
  teamDamageComposition: DamageComposition
}

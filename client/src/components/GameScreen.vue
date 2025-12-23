<script setup lang="ts">
import { ref, computed, watch, nextTick, onMounted } from 'vue'
import { useGameStore } from '../stores/game'
import { useCharacterStore } from '../stores/character'
import { useAuthStore } from '../stores/auth'
import { getClassColor, getResourceColor } from '../types/game'
import ChatPanel from './ChatPanel.vue'
import StrategyEditor from './StrategyEditor.vue'
import BattleStatsPanel from './BattleStatsPanel.vue'

const emit = defineEmits<{
  logout: []
  'create-character': []
}>()

const game = useGameStore()
const charStore = useCharacterStore()
const authStore = useAuthStore()
const logContainer = ref<HTMLElement | null>(null)
const skillSelection = computed(() => game.skillSelection)
const skillSelectionLoading = computed(() => game.skillSelectionLoading)
const selectingSkill = ref(false)
const selectionError = ref('')

// 角色详情弹窗
const showCharacterDetail = ref(false)
const selectedCharacter = ref<any>(null)
const characterSkills = ref<any[]>([])
const passiveSkills = ref<any[]>([])
const loadingSkills = ref(false)
const allocating = ref(false)

// 策略编辑器
const showStrategyEditor = ref(false)
const strategyCharacterId = ref<number | null>(null)
const strategyCharacterSkills = ref<any[]>([])

// 战斗统计面板
const showStatsPanel = ref(false)

// 地图选择面板
const showZoneSelector = ref(false)
const availableZones = ref<any[]>([])
const loadingZones = ref(false)
const zoneError = ref('')

// 敌人名称到攻击类型的映射，用于保持敌人颜色一致
const enemyAttackTypeMap = ref<Record<string, string>>({})

// 当前地图名称
const currentZoneName = computed(() => {
  const zoneId = game.battleStatus?.currentZoneId || game.battleStatus?.current_zone
  if (!zoneId) return '未知地图'
  
  // 从地图列表中查找
  const zone = game.zones.find(z => z.id === zoneId)
  if (zone) return zone.name
  
  // 如果还没加载地图列表，显示ID
  return zoneId
})

// 打开策略编辑器
async function openStrategyEditor() {
  // 获取当前活跃角色
  const activeChar = charStore.activeCharacter
  if (!activeChar) {
    console.warn('No active character')
    return
  }
  strategyCharacterId.value = activeChar.id
  // 获取角色技能用于策略配置
  try {
    const response = await fetch(`/api/characters/${activeChar.id}/skills`, {
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      }
    })
    const data = await response.json()
    if (data.success && data.data) {
      strategyCharacterSkills.value = data.data.activeSkills || []
    }
  } catch (e) {
    console.error('Failed to fetch skills for strategy:', e)
  }
  showStrategyEditor.value = true
}

// 关闭策略编辑器
function closeStrategyEditor() {
  showStrategyEditor.value = false
}

// 打开战斗统计面板
function openStatsPanel() {
  showStatsPanel.value = true
}

// 关闭战斗统计面板
function closeStatsPanel() {
  showStatsPanel.value = false
}

// 打开地图选择器
async function openZoneSelector() {
  console.log('Opening zone selector...')
  showZoneSelector.value = true
  await loadZones()
}

// 关闭地图选择器
function closeZoneSelector() {
  showZoneSelector.value = false
}

// 加载地图列表
async function loadZones() {
  loadingZones.value = true
  zoneError.value = ''
  try {
    console.log('Loading zones...')
    await Promise.all([game.fetchZones(), game.fetchExplorations()])
    console.log('Zones loaded:', game.zones)
    console.log('Explorations loaded:', game.explorations)
    availableZones.value = game.zones
    if (availableZones.value.length === 0) {
      zoneError.value = '没有可用的地图'
    }
  } catch (e) {
    zoneError.value = '加载地图失败: ' + (e instanceof Error ? e.message : String(e))
    console.error('Failed to load zones:', e)
  } finally {
    loadingZones.value = false
  }
}

// 检查地图是否已解锁
function isZoneUnlocked(zone: any): boolean {
  // 如果没有解锁条件，直接返回true
  if (!zone.unlockZoneId || zone.requiredExploration === 0 || zone.requiredExploration === undefined) {
    return true // 初始地图或无需解锁
  }
  const exploration = game.explorations[zone.unlockZoneId]
  if (!exploration) {
    return false // 没有探索度数据，视为未解锁
  }
  return exploration.exploration >= zone.requiredExploration
}

// 获取解锁进度信息
function getUnlockProgress(zone: any): { current: number; required: number; unlocked: boolean; unlockZoneName: string } | null {
  // 如果没有解锁条件，不显示进度
  if (!zone.unlockZoneId || zone.requiredExploration === 0 || zone.requiredExploration === undefined) {
    return null // 初始地图
  }
  const exploration = game.explorations[zone.unlockZoneId]
  const unlockZone = game.zones.find(z => z.id === zone.unlockZoneId)
  return {
    current: exploration?.exploration || 0,
    required: zone.requiredExploration,
    unlocked: isZoneUnlocked(zone),
    unlockZoneName: unlockZone?.name || zone.unlockZoneId
  }
}

// 切换地图
async function selectZone(zoneId: string) {
  console.log('Selecting zone:', zoneId)
  zoneError.value = ''
  const success = await game.changeZone(zoneId)
  if (success) {
    console.log('Zone changed successfully')
    closeZoneSelector()
    await game.fetchBattleStatus()
    await game.fetchBattleLogs()
  } else {
    zoneError.value = '切换地图失败，请检查等级和阵营限制'
    console.error('Failed to change zone')
  }
}

// 显示角色详情
async function showDetail(char: any) {
  selectedCharacter.value = char
  showCharacterDetail.value = true
  // 获取角色技能
  await fetchCharacterSkills(char.id)
}

// 获取角色技能
async function fetchCharacterSkills(characterId: number) {
  loadingSkills.value = true
  try {
    const response = await fetch(`/api/characters/${characterId}/skills`, {
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      }
    })
    const data = await response.json()
    if (data.success && data.data) {
      // API返回格式: { activeSkills: [...], passiveSkills: [...] }
      characterSkills.value = data.data.activeSkills || []
      passiveSkills.value = data.data.passiveSkills || []
    } else {
      characterSkills.value = []
      passiveSkills.value = []
    }
  } catch (e) {
    console.error('Failed to fetch character skills:', e)
    characterSkills.value = []
    passiveSkills.value = []
  } finally {
    loadingSkills.value = false
  }
}

// 检查是否有技能选择机会（被动/主动）
async function refreshSkillSelection(force = false) {
  await game.checkSkillSelection(force)
}

// 提交技能选择
async function submitSelection(payload: { skillId?: string; passiveId?: string; isUpgrade: boolean }) {
  selectionError.value = ''
  selectingSkill.value = true
  try {
    const ok = await game.submitSkillSelection(payload)
    if (!ok) {
      selectionError.value = '技能选择提交失败，请重试'
    }
  } catch (e) {
    console.error('submit selection failed', e)
    selectionError.value = '技能选择提交失败，请稍后重试'
  } finally {
    selectingSkill.value = false
  }
}

function getPassiveName(passive: any): string {
  if (!passive) return '未知被动'
  return passive.passive?.name || passive.name || passive.passiveId || '未知被动'
}

function getSkillName(skill: any): string {
  if (!skill) return '未知技能'
  return skill.skill?.name || skill.name || skill.skillId || '未知技能'
}

// 关闭角色详情
function closeDetail() {
  showCharacterDetail.value = false
  selectedCharacter.value = null
  characterSkills.value = []
  passiveSkills.value = []
  hideSkillTooltip() // 关闭时也隐藏tooltip
  hideBuffTooltip() // 关闭时也隐藏buff tooltip
}

// 获取buff tooltip文本
function getBuffTooltip(buff: any): string {
  if (!buff) return ''
  
  const parts: string[] = []
  // Buff名称
  const buffName = buff.name || '未知效果'
  parts.push(`【${buffName}】`)
  
  // Buff描述 - 优先使用后端返回的description
  const description = buff.description || ''
  if (description && description.trim()) {
    parts.push(description)
  } else {
    // 如果没有描述，根据类型生成描述
    const statAffected = buff.statAffected || ''
    const isBuff = buff.isBuff !== undefined ? buff.isBuff : true
    const value = buff.value !== undefined ? buff.value : 0
    
    if (statAffected === 'attack' && isBuff) {
      parts.push(`提升 ${Math.round(value)}% 物理攻击力`)
    } else if (statAffected === 'attack' && !isBuff) {
      parts.push(`降低 ${Math.round(Math.abs(value))}% 物理攻击力`)
    } else if (statAffected === 'defense' && isBuff) {
      parts.push(`提升 ${Math.round(value)}% 物理防御`)
    } else if (statAffected === 'physical_damage_taken' || statAffected === 'damage_taken') {
      parts.push(`减少 ${Math.round(Math.abs(value))}% 受到的物理伤害`)
    } else if (statAffected === 'crit_rate' && isBuff) {
      parts.push(`提升 ${Math.round(value)}% 暴击率`)
    } else if (statAffected === 'shield') {
      parts.push(`获得相当于最大HP ${Math.round(value)} 点的护盾`)
    } else if (statAffected === 'reflect') {
      parts.push(`反射 ${Math.round(value)}% 受到的伤害`)
    } else if (statAffected === 'counter_attack') {
      parts.push(`受到攻击时反击，造成 ${Math.round(value)}% 物理攻击力伤害`)
    } else if (statAffected === 'cc_immune') {
      parts.push('免疫控制效果')
    } else if (statAffected === 'healing_received') {
      parts.push(`降低 ${Math.round(value)}% 治疗效果`)
    } else {
      // 如果无法识别类型，至少显示buff名称
      parts.push(buffName)
    }
  }
  
  // 持续时间
  const duration = buff.duration !== undefined ? buff.duration : null
  if (duration !== undefined && duration !== null && duration > 0) {
    parts.push(`━━━━━━━━━━━━━━━━`)
    parts.push(`剩余: ${duration} 回合`)
  }
  
  const result = parts.join('\n')
  return result || `【${buffName}】` // 至少返回buff名称
}

// 获取技能tooltip文本
function getSkillTooltip(skill: any): string {
  if (!skill) return ''
  
  // 如果没有skill详情，至少返回skillId
  if (!skill.skill) {
    return skill.skillId || '未知技能'
  }
  
  const parts: string[] = []
  // 技能名称
  const skillName = skill.skill.name || skill.skillId || '未知技能'
  parts.push(skillName)
  
  // 技能描述
  if (skill.skill.description) {
    parts.push(skill.skill.description)
  }
  
  // 技能详情
  const details: string[] = []
  if (skill.skillLevel) {
    details.push(`等级: ${skill.skillLevel}`)
  }
  if (skill.skill.resourceCost !== undefined && skill.skill.resourceCost !== null) {
    const resourceName = getResourceTypeName(selectedCharacter.value)
    const resourceShort = resourceName === '怒气' ? '怒' : resourceName === '能量' ? '能' : 'MP'
    details.push(`消耗: ${skill.skill.resourceCost}${resourceShort}`)
  }
  if (skill.skill.cooldown !== undefined && skill.skill.cooldown !== null && skill.skill.cooldown > 0) {
    details.push(`冷却: ${skill.skill.cooldown}回合`)
  }
  
  if (details.length > 0) {
    parts.push(details.join(' | '))
  }
  
  const result = parts.join('\n')
  return result || skillName // 至少返回技能名称
}

// 计算被动技能的实际效果值（含等级成长）
function calculatePassiveEffectValue(passive: any, level: number): number {
  const base = passive?.effectValue ?? passive?.EffectValue ?? 0
  const scaling = passive?.levelScaling ?? passive?.LevelScaling ?? 0
  const lvl = Math.max(1, level || passive?.level || passive?.Level || 1)
  return base + (lvl - 1) * scaling
}

// 汇总所有被动的属性加成（百分比数值）
const passiveStatModifiers = computed<Record<string, number>>(() => {
  const mods: Record<string, number> = {}
  passiveSkills.value.forEach(ps => {
    const passive = ps.passive || ps.Passive || ps
    if (!passive || passive.effectType !== 'stat_mod') return
    const stat = passive.effectStat || passive.effect_stat
    if (!stat) return
    const level = ps.level || ps.Level || ps.passiveLevel || ps.skillLevel || 1
    const value = calculatePassiveEffectValue(passive, level)
    mods[stat] = (mods[stat] || 0) + value
  })
  return mods
})

// 计算包含被动加成后的展示属性
const displayedStats = computed(() => {
  const zero = {
    physicalAttack: 0,
    magicAttack: 0,
    physicalDefense: 0,
    magicDefense: 0,
    physCritRate: 0,
    physCritDamage: 0,
    spellCritRate: 0,
    spellCritDamage: 0,
    dodgeRate: 0,
    physicalAttackBonusPct: 0,
    magicAttackBonusPct: 0,
    physicalDefenseBonusPct: 0,
    magicDefenseBonusPct: 0,
    physCritBonusPct: 0,
    spellCritBonusPct: 0,
    physCritDamageBonusPct: 0,
    spellCritDamageBonusPct: 0,
    dodgeBonusPct: 0,
  }
  const char = selectedCharacter.value
  if (!char) return zero

  const str = char.strength || 0
  const agi = char.agility || 0
  const intl = char.intellect || 0
  const spr = char.spirit || 0

  const basePhysicalAttack = char.physicalAttack ?? (str * 0.6 + agi * 0.2)
  const baseMagicAttack = char.magicAttack ?? (intl * 1 + spr * 0.2)
  const basePhysicalDefense = char.physicalDefense ?? (str * 0.2 + (char.stamina || 0) * 0.3)
  const baseMagicDefense = char.magicDefense ?? (intl * 0.2 + spr * 0.3)
  const basePhysCritRate = char.physCritRate ?? (0.05 + agi / 20)
  const basePhysCritDamage = char.physCritDamage ?? (1.5 + str * 0.003)
  const baseSpellCritRate = char.spellCritRate ?? (0.05 + spr / 20)
  const baseSpellCritDamage = char.spellCritDamage ?? (1.5 + intl * 0.003)
  const baseDodgeRate = char.dodgeRate ?? (0.05 + agi / 20)

  const physicalAttackBonusPct = passiveStatModifiers.value['physical_attack'] || 0
  const magicAttackBonusPct = passiveStatModifiers.value['magic_attack'] || 0
  const physicalDefenseBonusPct = passiveStatModifiers.value['physical_defense'] || 0
  const magicDefenseBonusPct = passiveStatModifiers.value['magic_defense'] || 0
  const physCritBonusPct = (passiveStatModifiers.value['phys_crit_rate'] || 0) + (passiveStatModifiers.value['crit_rate'] || 0)
  const spellCritBonusPct = (passiveStatModifiers.value['spell_crit_rate'] || 0) + (passiveStatModifiers.value['crit_rate'] || 0)
  const physCritDamageBonusPct = passiveStatModifiers.value['phys_crit_damage'] || 0
  const spellCritDamageBonusPct = passiveStatModifiers.value['spell_crit_damage'] || 0
  const dodgeBonusPct = passiveStatModifiers.value['dodge_rate'] || 0

  return {
    physicalAttack: Math.round(basePhysicalAttack * (1 + physicalAttackBonusPct / 100)),
    magicAttack: Math.round(baseMagicAttack * (1 + magicAttackBonusPct / 100)),
    physicalDefense: Math.round(basePhysicalDefense * (1 + physicalDefenseBonusPct / 100)),
    magicDefense: Math.round(baseMagicDefense * (1 + magicDefenseBonusPct / 100)),
    physCritRate: basePhysCritRate + physCritBonusPct / 100,
    physCritDamage: basePhysCritDamage + physCritDamageBonusPct / 100,
    spellCritRate: baseSpellCritRate + spellCritBonusPct / 100,
    spellCritDamage: baseSpellCritDamage + spellCritDamageBonusPct / 100,
    dodgeRate: Math.min(0.5, baseDodgeRate + dodgeBonusPct / 100),
    physicalAttackBonusPct,
    magicAttackBonusPct,
    physicalDefenseBonusPct,
    magicDefenseBonusPct,
    physCritBonusPct,
    spellCritBonusPct,
    physCritDamageBonusPct,
    spellCritDamageBonusPct,
    dodgeBonusPct,
  }
})

function getEffectStatLabel(stat: string | undefined): string {
  const map: Record<string, string> = {
    physical_attack: '物理攻击',
    magic_attack: '魔法攻击',
    physical_defense: '物理防御',
    magic_defense: '魔法防御',
    phys_crit_rate: '物理暴击率',
    spell_crit_rate: '法术暴击率',
    crit_rate: '通用暴击率',
    phys_crit_damage: '物理暴击伤害',
    spell_crit_damage: '法术暴击伤害',
    dodge_rate: '闪避率',
    hp: '生命值',
    max_hp: '最大生命值',
    attack: '攻击力',
    defense: '防御',
    damage: '伤害',
    threat: '仇恨',
    resistance: '控制抗性',
    // 多属性被动技能
    threat_and_defense: '仇恨和防御',
    attack_and_crit: '攻击和暴击',
    damage_and_threat: '伤害和仇恨',
    hp_and_defense: '生命值和防御',
    damage_and_crit: '伤害和暴击',
    hp_defense_resistance: '生命值、防御和抗性',
    damage_crit_threat: '伤害、暴击和仇恨',
    hp_defense_resistance_immune: '生命值、防御、抗性和免疫',
  }
  return map[stat || ''] || (stat || '未知属性')
}

function formatPassiveEffect(passiveSkill: any): string {
  const passive = passiveSkill?.passive || passiveSkill?.Passive || passiveSkill
  if (!passive) return ''
  const level = passiveSkill.level || passiveSkill.Level || passiveSkill.passiveLevel || passiveSkill.skillLevel || 1
  const value = calculatePassiveEffectValue(passive, level)
  if (passive.effectType === 'stat_mod') {
    return `+${value}% ${getEffectStatLabel(passive.effectStat || passive.effect_stat)}`
  }
  return passive.description || ''
}

// 获取被动技能tooltip
function getPassiveTooltip(passiveSkill: any): string {
  if (!passiveSkill) return ''
  const passive = passiveSkill.passive || passiveSkill.Passive || passiveSkill
  const level = passiveSkill.level || passiveSkill.Level || passiveSkill.passiveLevel || passiveSkill.skillLevel || 1
  const parts: string[] = []
  const name = passive?.name || passiveSkill.passiveId || '未知被动'
  parts.push(name)
  if (passive?.description) {
    parts.push(passive.description)
  }
  if (passive?.effectType === 'stat_mod') {
    const value = calculatePassiveEffectValue(passive, level)
    parts.push(`效果: +${value}% ${getEffectStatLabel(passive.effectStat || passive.effect_stat)}`)
  }
  parts.push(`等级: ${level}/${passive?.maxLevel ?? passive?.MaxLevel ?? 5}`)
  return parts.filter(Boolean).join('\n')
}

// 被动技能tooltip展示
function handlePassiveTooltip(event: MouseEvent, passiveSkill: any) {
  const tooltipText = getPassiveTooltip(passiveSkill)
  if (!tooltipText) return
  if (skillTooltipEl) {
    skillTooltipEl.remove()
  }
  skillTooltipEl = document.createElement('div')
  skillTooltipEl.className = 'skill-tooltip-fixed'
  skillTooltipEl.textContent = tooltipText
  document.body.appendChild(skillTooltipEl)

  const rect = (event.currentTarget as HTMLElement).getBoundingClientRect()
  const tooltipRect = skillTooltipEl.getBoundingClientRect()
  let left = rect.left + (rect.width / 2) - (tooltipRect.width / 2)
  let top = rect.top - tooltipRect.height - 8
  if (left < 10) left = 10
  if (left + tooltipRect.width > window.innerWidth - 10) {
    left = window.innerWidth - tooltipRect.width - 10
  }
  if (top < 10) {
    top = rect.bottom + 8
  }
  skillTooltipEl.style.left = left + 'px'
  skillTooltipEl.style.top = top + 'px'
}

// 战斗属性 tooltip，包含被动加成说明
function getCombatStatTooltip(statKey: string): string {
  const char = selectedCharacter.value
  if (!char) return ''
  const stats = displayedStats.value
  const agi = char.agility || 0
  const str = char.strength || 0
  const intl = char.intellect || 0
  const spr = char.spirit || 0
  const sta = char.stamina || 0

  switch (statKey) {
    case 'physicalAttack': {
      const base = (str * 0.6 + agi * 0.2).toFixed(1)
      const bonus = stats.physicalAttackBonusPct
      return [
        `物理攻击 = 力量×0.6 + 敏捷×0.2${bonus ? `，再乘(1 + 被动${bonus.toFixed(1)}%)` : ''}`,
        `力量: ${str}, 敏捷: ${agi}`,
        `基础计算: ${str}×0.6 + ${agi}×0.2 = ${base}`,
        bonus ? `被动加成: +${bonus.toFixed(1)}%` : '',
        `最终: ${stats.physicalAttack}`
      ].filter(Boolean).join('\n')
    }
    case 'magicAttack': {
      const base = (intl * 1 + spr * 0.2).toFixed(1)
      const bonus = stats.magicAttackBonusPct
      return [
        `魔法攻击 = 智力×1.0 + 精神×0.2${bonus ? `，再乘(1 + 被动${bonus.toFixed(1)}%)` : ''}`,
        `智力: ${intl}, 精神: ${spr}`,
        `基础计算: ${intl}×1.0 + ${spr}×0.2 = ${base}`,
        bonus ? `被动加成: +${bonus.toFixed(1)}%` : '',
        `最终: ${stats.magicAttack}`
      ].filter(Boolean).join('\n')
    }
    case 'physicalDefense': {
      const base = (str * 0.2 + sta * 0.3).toFixed(1)
      const bonus = stats.physicalDefenseBonusPct
      return [
        `物理防御 = 力量×0.2 + 耐力×0.3${bonus ? `，再乘(1 + 被动${bonus.toFixed(1)}%)` : ''}`,
        `力量: ${str}, 耐力: ${sta}`,
        `基础计算: ${str}×0.2 + ${sta}×0.3 = ${base}`,
        bonus ? `被动加成: +${bonus.toFixed(1)}%` : '',
        `最终: ${stats.physicalDefense}`
      ].filter(Boolean).join('\n')
    }
    case 'magicDefense': {
      const base = (intl * 0.2 + spr * 0.3).toFixed(1)
      const bonus = stats.magicDefenseBonusPct
      return [
        `魔法防御 = 智力×0.2 + 精神×0.3${bonus ? `，再乘(1 + 被动${bonus.toFixed(1)}%)` : ''}`,
        `智力: ${intl}, 精神: ${spr}`,
        `基础计算: ${intl}×0.2 + ${spr}×0.3 = ${base}`,
        bonus ? `被动加成: +${bonus.toFixed(1)}%` : '',
        `最终: ${stats.magicDefense}`
      ].filter(Boolean).join('\n')
    }
    case 'physCritRate': {
      const bonus = stats.physCritBonusPct
      return [
        `物理暴击率 = 5% + 敏捷/20${bonus ? ` + 被动(${bonus.toFixed(1)}%)` : ''}`,
        `当前敏捷: ${agi}`,
        `计算: 5% + ${agi}/20${bonus ? ` + ${bonus.toFixed(1)}%` : ''} = ${(stats.physCritRate * 100).toFixed(1)}%`
      ].join('\n')
    }
    case 'physCritDamage': {
      const bonus = stats.physCritDamageBonusPct
      return [
        `物理暴击伤害 = 150% + 力量×0.3%${bonus ? ` + 被动(${bonus.toFixed(1)}%)` : ''}`,
        `当前力量: ${str}`,
        `计算: 150% + ${str}×0.3%${bonus ? ` + ${bonus.toFixed(1)}%` : ''} = ${(stats.physCritDamage * 100).toFixed(1)}%`
      ].join('\n')
    }
    case 'spellCritRate': {
      const bonus = stats.spellCritBonusPct
      return [
        `法术暴击率 = 5% + 精神/20${bonus ? ` + 被动(${bonus.toFixed(1)}%)` : ''}`,
        `当前精神: ${spr}`,
        `计算: 5% + ${spr}/20${bonus ? ` + ${bonus.toFixed(1)}%` : ''} = ${(stats.spellCritRate * 100).toFixed(1)}%`
      ].join('\n')
    }
    case 'spellCritDamage': {
      const bonus = stats.spellCritDamageBonusPct
      return [
        `法术暴击伤害 = 150% + 智力×0.3%${bonus ? ` + 被动(${bonus.toFixed(1)}%)` : ''}`,
        `当前智力: ${intl}`,
        `计算: 150% + ${intl}×0.3%${bonus ? ` + ${bonus.toFixed(1)}%` : ''} = ${(stats.spellCritDamage * 100).toFixed(1)}%`
      ].join('\n')
    }
    case 'dodgeRate': {
      const bonus = stats.dodgeBonusPct
      return [
        `闪避率 = 5% + 敏捷/20${bonus ? ` + 被动(${bonus.toFixed(1)}%)` : ''}`,
        `当前敏捷: ${agi}`,
        `计算: 5% + ${agi}/20${bonus ? ` + ${bonus.toFixed(1)}%` : ''} = ${(stats.dodgeRate * 100).toFixed(1)}%`,
        '闪避可以躲避物理和魔法攻击'
      ].join('\n')
    }
    default:
      return ''
  }
}

type PrimaryStatKey = 'strength' | 'agility' | 'intellect' | 'stamina' | 'spirit'

// 主属性 tooltip：说明该属性对派生数值的影响
function getPrimaryStatTooltip(statKey: PrimaryStatKey): string {
  const char = selectedCharacter.value
  if (!char) return ''

  const str = char.strength || 0
  const agi = char.agility || 0
  const intl = char.intellect || 0
  const sta = char.stamina || 0
  const spr = char.spirit || 0

  switch (statKey) {
    case 'strength': {
      const physAtk = (str * 0.6).toFixed(1)
      const physDef = (str * 0.2).toFixed(1)
      const critDmg = ((char.physCritDamage ?? (1.5 + str * 0.003)) * 100).toFixed(1)
      return [
        '力量',
        `- 每点提供 0.6 物理攻击 (当前贡献: +${physAtk})`,
        `- 每点提供 0.2 物理防御 (当前贡献: +${physDef})`,
        `- 物理暴击伤害 = 150% + 力量×0.3% (当前: ${critDmg}%)`
      ].join('\n')
    }
    case 'agility': {
      const physAtk = (agi * 0.2).toFixed(1)
      const critRate = ((char.physCritRate ?? (0.05 + agi / 20)) * 100).toFixed(1)
      const dodge = ((char.dodgeRate ?? (0.05 + agi / 20)) * 100).toFixed(1)
      return [
        '敏捷',
        `- 每点提供 0.2 物理攻击 (当前贡献: +${physAtk})`,
        `- 物理暴击率 = 5% + 敏捷/20 (当前: ${critRate}%)`,
        `- 闪避率 = 5% + 敏捷/20 (当前: ${dodge}%)`
      ].join('\n')
    }
    case 'intellect': {
      const magicAtk = (intl * 1).toFixed(1)
      const magicDef = (intl * 0.2).toFixed(1)
      const spellCritDmg = ((char.spellCritDamage ?? (1.5 + intl * 0.003)) * 100).toFixed(1)
      return [
        '智力',
        `- 每点提供 1.0 魔法攻击 (当前贡献: +${magicAtk})`,
        `- 每点提供 0.2 魔法防御 (当前贡献: +${magicDef})`,
        `- 法术暴击伤害 = 150% + 智力×0.3% (当前: ${spellCritDmg}%)`
      ].join('\n')
    }
    case 'stamina': {
      const physDef = (sta * 0.3).toFixed(1)
      return [
        '耐力',
        `- 每点提供 0.3 物理防御 (当前贡献: +${physDef})`,
        '- 提升生存能力，适合坦克/前排'
      ].join('\n')
    }
    case 'spirit': {
      const magicAtk = (spr * 0.2).toFixed(1)
      const magicDef = (spr * 0.3).toFixed(1)
      const spellCritRate = ((char.spellCritRate ?? (0.05 + spr / 20)) * 100).toFixed(1)
      return [
        '精神',
        `- 每点提供 0.2 魔法攻击 (当前贡献: +${magicAtk})`,
        `- 每点提供 0.3 魔法防御 (当前贡献: +${magicDef})`,
        `- 法术暴击率 = 5% + 精神/20 (当前: ${spellCritRate}%)`
      ].join('\n')
    }
    default:
      return ''
  }
}

// 分配单点主属性
async function allocateStat(statKey: PrimaryStatKey) {
  if (!selectedCharacter.value) return
  if (allocating.value) return
  if (!selectedCharacter.value.unspentPoints || selectedCharacter.value.unspentPoints <= 0) return

  allocating.value = true
  try {
    const body: Record<string, number> = {
      strength: 0,
      agility: 0,
      intellect: 0,
      stamina: 0,
      spirit: 0
    }
    body[statKey] = 1

    const resp = await fetch(`/api/characters/${selectedCharacter.value.id}/allocate`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${localStorage.getItem('token') || ''}`
      },
      body: JSON.stringify(body)
    })

    const data = await resp.json()
    if (data.success && data.data) {
      selectedCharacter.value = data.data
      // 同步到 charStore 列表
      const idx = charStore.characters.findIndex((c: any) => c.id === data.data.id)
      if (idx >= 0) {
        charStore.characters[idx] = data.data
      }
    } else {
      console.error('Allocate stat failed:', data.error || data.message)
    }
  } catch (e) {
    console.error('Allocate stat failed:', e)
  } finally {
    allocating.value = false
  }
}

// 处理buff tooltip显示（使用fixed定位，智能调整位置避免超出屏幕）
let buffTooltipEl: HTMLElement | null = null

function handleBuffTooltip(event: MouseEvent, buff: any) {
  const tooltipText = getBuffTooltip(buff)
  if (!tooltipText) return
  
  // 移除旧的tooltip
  if (buffTooltipEl) {
    buffTooltipEl.remove()
  }
  
  // 创建新的tooltip元素
  buffTooltipEl = document.createElement('div')
  buffTooltipEl.className = 'buff-tooltip-fixed'
  buffTooltipEl.textContent = tooltipText
  document.body.appendChild(buffTooltipEl)
  
  // 计算位置
  const rect = (event.currentTarget as HTMLElement).getBoundingClientRect()
  const tooltipRect = buffTooltipEl.getBoundingClientRect()
  
  // 默认显示在右侧
  let left = rect.right + 8
  let top = rect.top + (rect.height / 2) - (tooltipRect.height / 2)
  
  // 如果右侧空间不够，显示在左侧
  if (left + tooltipRect.width > window.innerWidth - 10) {
    left = rect.left - tooltipRect.width - 8
  }
  
  // 如果左侧也不够，显示在上方
  if (left < 10) {
    left = rect.left + (rect.width / 2) - (tooltipRect.width / 2)
    top = rect.top - tooltipRect.height - 8
  }
  
  // 确保不超出视口
  if (left < 10) left = 10
  if (left + tooltipRect.width > window.innerWidth - 10) {
    left = window.innerWidth - tooltipRect.width - 10
  }
  if (top < 10) {
    top = rect.bottom + 8
  }
  if (top + tooltipRect.height > window.innerHeight - 10) {
    top = window.innerHeight - tooltipRect.height - 10
  }
  
  buffTooltipEl.style.left = left + 'px'
  buffTooltipEl.style.top = top + 'px'
}

function hideBuffTooltip() {
  if (buffTooltipEl) {
    buffTooltipEl.remove()
    buffTooltipEl = null
  }
}

// 处理技能/属性/战斗统计 tooltip（使用 fixed，避免被面板裁剪）
let skillTooltipEl: HTMLElement | null = null
let attrTooltipEl: HTMLElement | null = null
let combatTooltipEl: HTMLElement | null = null

function handleSkillTooltip(event: MouseEvent, skill: any) {
  const tooltipText = getSkillTooltip(skill)
  if (!tooltipText) return
  
  // 移除旧的tooltip
  if (skillTooltipEl) {
    skillTooltipEl.remove()
  }
  
  // 创建新的tooltip元素
  skillTooltipEl = document.createElement('div')
  skillTooltipEl.className = 'skill-tooltip-fixed'
  skillTooltipEl.textContent = tooltipText
  document.body.appendChild(skillTooltipEl)
  
  // 计算位置
  const rect = (event.currentTarget as HTMLElement).getBoundingClientRect()
  const tooltipRect = skillTooltipEl.getBoundingClientRect()
  
  let left = rect.left + (rect.width / 2) - (tooltipRect.width / 2)
  let top = rect.top - tooltipRect.height - 8
  
  // 确保不超出视口
  if (left < 10) left = 10
  if (left + tooltipRect.width > window.innerWidth - 10) {
    left = window.innerWidth - tooltipRect.width - 10
  }
  if (top < 10) {
    top = rect.bottom + 8
  }
  
  skillTooltipEl.style.left = left + 'px'
  skillTooltipEl.style.top = top + 'px'
}

function hideSkillTooltip() {
  if (skillTooltipEl) {
    skillTooltipEl.remove()
    skillTooltipEl = null
  }
}

// 属性tooltip（使用fixed，避免被面板裁剪）
function handleAttrTooltip(event: MouseEvent, statKey: PrimaryStatKey) {
  const tooltipText = getPrimaryStatTooltip(statKey)
  if (!tooltipText) return

  if (attrTooltipEl) {
    attrTooltipEl.remove()
  }

  attrTooltipEl = document.createElement('div')
  attrTooltipEl.className = 'attr-tooltip-fixed'
  attrTooltipEl.textContent = tooltipText
  document.body.appendChild(attrTooltipEl)

  const rect = (event.currentTarget as HTMLElement).getBoundingClientRect()
  const tooltipRect = attrTooltipEl.getBoundingClientRect()

  // 默认在右侧居中
  let left = rect.right + 12
  let top = rect.top + rect.height / 2 - tooltipRect.height / 2

  // 如果右侧超出，放到左侧
  if (left + tooltipRect.width > window.innerWidth - 10) {
    left = rect.left - tooltipRect.width - 12
  }

  // 如果左右都不够，居中放上方
  if (left < 10) {
    left = rect.left + rect.width / 2 - tooltipRect.width / 2
    top = rect.top - tooltipRect.height - 10
  }

  // 视口保护
  if (left < 10) left = 10
  if (left + tooltipRect.width > window.innerWidth - 10) {
    left = window.innerWidth - tooltipRect.width - 10
  }
  if (top < 10) {
    top = rect.bottom + 10
  }
  if (top + tooltipRect.height > window.innerHeight - 10) {
    top = window.innerHeight - tooltipRect.height - 10
  }

  attrTooltipEl.style.left = `${left}px`
  attrTooltipEl.style.top = `${top}px`
}

function hideAttrTooltip() {
  if (attrTooltipEl) {
    attrTooltipEl.remove()
    attrTooltipEl = null
  }
}

// 战斗统计 tooltip（fixed，避免被面板 overflow 裁剪）
function handleCombatTooltip(event: MouseEvent, tooltipText?: string) {
  if (!tooltipText) return

  if (combatTooltipEl) {
    combatTooltipEl.remove()
  }

  combatTooltipEl = document.createElement('div')
  combatTooltipEl.className = 'attr-tooltip-fixed'
  combatTooltipEl.textContent = tooltipText
  document.body.appendChild(combatTooltipEl)

  const rect = (event.currentTarget as HTMLElement).getBoundingClientRect()
  const tooltipRect = combatTooltipEl.getBoundingClientRect()

  // 默认显示在元素上方居中
  let left = rect.left + rect.width / 2 - tooltipRect.width / 2
  let top = rect.top - tooltipRect.height - 10

  // 若上方空间不足，放到底部
  if (top < 10) {
    top = rect.bottom + 10
  }

  // 视口保护
  if (left < 10) left = 10
  if (left + tooltipRect.width > window.innerWidth - 10) {
    left = window.innerWidth - tooltipRect.width - 10
  }
  if (top + tooltipRect.height > window.innerHeight - 10) {
    top = window.innerHeight - tooltipRect.height - 10
  }

  combatTooltipEl.style.left = `${left}px`
  combatTooltipEl.style.top = `${top}px`
}

function hideCombatTooltip() {
  if (combatTooltipEl) {
    combatTooltipEl.remove()
    combatTooltipEl = null
  }
}


// 初始化：从 characterStore 获取角色数据并同步到 gameStore
onMounted(async () => {
  console.log('GameScreen mounted')
  console.log('charStore.characters:', charStore.characters)
  
  // 如果没有角色，先尝试获取
  if (charStore.characters.length === 0) {
    await charStore.fetchCharacters()
  }
  
  // 获取第一个角色（所有角色都参与战斗）
  const activeChar = charStore.characters[0]
  
  console.log('activeChar:', activeChar)
  
  // 优先从 API 获取最新的角色数据（包含死亡/复活状态）
  await game.fetchCharacter()
  
  if (game.character) {
    console.log('Character loaded from API:', game.character)
  } else if (activeChar) {
    // 如果 API 没有返回，使用 characterStore 中的数据作为后备
    game.character = activeChar
    console.log('Character synced from characterStore:', game.character)
  }
  
  // 获取战斗状态和日志
  await game.fetchBattleStatus()
  await game.fetchBattleLogs()
  await game.fetchZones() // 确保地图列表已加载
  await refreshSkillSelection(true)
  
  // 检查是否需要选择地图（延迟一下，确保数据已加载）
  await nextTick()
  const zoneId = game.battleStatus?.currentZoneId || game.battleStatus?.current_zone
  if (!zoneId || zoneId === '' || currentZoneName.value === '未知地图') {
    // 如果没有地图，自动打开地图选择器引导玩家选择
    console.log('No zone selected, opening zone selector...')
    setTimeout(() => {
      openZoneSelector()
    }, 500) // 延迟500ms，让界面先渲染完成
  }
  
  // 如果战斗状态中有队伍数据，使用第一个角色作为当前显示角色
  // Team 是一个数组，包含所有角色（所有角色都参与战斗）
  if (game.battleStatus?.team && Array.isArray(game.battleStatus.team) && game.battleStatus.team.length > 0) {
    game.character = game.battleStatus.team[0]
    console.log('Character updated from battle status team:', game.character)
    console.log('Team size:', game.battleStatus.team.length)
  } else if (charStore.characters.length > 0) {
    // 如果没有队伍数据，使用第一个角色
    game.character = charStore.characters[0]
    console.log('Character set from characters:', game.character)
  }
  
  if (!game.character) {
    console.warn('No character found after all attempts!')
  }
})

// 自动滚动到底部
watch(() => game.battleLogs.length, async () => {
  await nextTick()
  if (logContainer.value) {
    logContainer.value.scrollTop = logContainer.value.scrollHeight
  }
})

// 计算角色HP/MP/EXP百分比（用于详情弹窗）
function getHpPercent(char: any): number {
  if (!char) return 0
  const maxHp = char.maxHp || char.max_hp || 100
  const hp = char.hp || 0
  return maxHp > 0 ? (hp / maxHp) * 100 : 0
}

function getMpPercent(char: any): number {
  if (!char) return 0
  const maxResource = char.maxResource || char.max_resource || char.max_mp || 100
  const resource = char.resource || char.mp || 0
  return maxResource > 0 ? (resource / maxResource) * 100 : 0
}

function getExpPercent(char: any): number {
  if (!char) return 0
  const expToNext = char.expToNext || char.exp_to_next || 100
  const exp = char.exp || 0
  return expToNext > 0 ? (exp / expToNext) * 100 : 0
}

function getResourceTypeName(char: any): string {
  if (!char) return 'MP'
  const type = char.resourceType || 'mana'
  const names: Record<string, string> = {
    mana: '法力',
    rage: '怒气',
    energy: '能量'
  }
  return names[type] || 'MP'
}

const enemyHpPercent = computed(() => {
  if (!game.currentEnemy) return 0
  const enemy = game.currentEnemy as any
  const maxHp = enemy.maxHp || enemy.max_hp || 100
  const hp = Math.max(0, enemy.hp || 0) // 确保HP不会小于0
  return (hp / maxHp) * 100
})

// 计算每个敌人的HP百分比
function getEnemyHpPercent(enemy: any): number {
  if (!enemy) return 0
  const maxHp = enemy.maxHp || enemy.max_hp || enemy.hp || 100
  const hp = Math.max(0, enemy.hp || 0) // 确保HP不会小于0
  return maxHp > 0 ? (hp / maxHp) * 100 : 0
}

// 获取怪物名称颜色（根据攻击类型，与战斗日志保持一致）
function getEnemyNameColor(enemy: any): string {
  if (!enemy) return '#ff7777' // 默认红色（物理攻击）
  
  // 优先使用 attackType 字段
  let attackType = enemy.attackType || enemy.attack_type
  if (attackType) {
    attackType = attackType.toLowerCase()
    if (attackType === 'magic') {
      return '#7777ff' // 魔法攻击用蓝色
    }
    return '#ff7777' // 物理攻击用红色
  }
  
  // 如果没有 attackType，根据魔法攻击和物理攻击的数值推断
  const magicAttack = enemy.magicAttack || enemy.magic_attack || 0
  const physicalAttack = enemy.physicalAttack || enemy.physical_attack || 0
  
  if (magicAttack > physicalAttack && magicAttack > 0) {
    return '#7777ff' // 魔法攻击用蓝色
  }
  
  return '#ff7777' // 默认物理攻击用红色
}

// 获取资源类型名称

// 获取日志类型的CSS类
function getLogClass(type: string) {
  return `log-type-${type}`
}

// 获取种族名称
function getRaceName(race: string) {
  const names: Record<string, string> = {
    human: '人类', dwarf: '矮人', nightelf: '暗夜精灵', gnome: '侏儒',
    orc: '兽人', undead: '亡灵', tauren: '牛头人', troll: '巨魔'
  }
  return names[race] || race
}

// 获取职业名称
function getClassName(cls: string) {
  const names: Record<string, string> = {
    warrior: '战士', mage: '法师', rogue: '盗贼', priest: '牧师',
    paladin: '圣骑士', hunter: '猎人', warlock: '术士', druid: '德鲁伊', shaman: '萨满'
  }
  return names[cls] || cls
}

// 格式化战斗日志时间
function formatLogTime(log: any): string {
  if (log.time) return log.time
  if (log.createdAt) {
    const date = new Date(log.createdAt)
    return date.toLocaleTimeString('zh-CN', { hour: '2-digit', minute: '2-digit', second: '2-digit' })
  }
  return new Date().toLocaleTimeString('zh-CN', { hour: '2-digit', minute: '2-digit', second: '2-digit' })
}

function getDamageBadge(damageType?: string): string {
  const type = (damageType || '').toLowerCase()
  if (type === 'physical') {
    return '<span class="damage-badge damage-badge-physical">物理</span>'
  }
  if (type === 'magic') {
    return '<span class="damage-badge damage-badge-magic">魔法</span>'
  }
  return ''
}

// 格式化日志消息，添加颜色标记
function formatLogMessage(log: any): string {
  let message = ''
  if (log.message) {
    message = log.message
  } else if (log.logType && log.value) {
    message = `${log.logType}: ${log.value}`
  } else {
    message = log.logType || '未知'
  }
  
  // 如果没有消息，直接返回
  if (!message) return ''
  
  // 伤害类型徽标
  const badge = getDamageBadge(log.damageType || log.damage_type)
  if (badge) {
    message = `${badge} ${message}`
  }
  
  // 获取角色名（我方）
  const playerName = game.character?.name || '你'
  const playerNameVariants = [playerName, '你', '勇士'] // 可能的变体
  
  // 获取角色职业颜色（根据职业ID）
  const character = game.character as any
  const classId = character?.classId || character?.class || ''
  const playerColor = getClassColor(classId) // 使用职业颜色，如果没有职业则使用默认绿色
  
  // 获取资源类型和资源颜色（用于技能颜色）
  const resourceType = character?.resourceType || 'mana'
  const resourceColor = getResourceColor(resourceType)
  
  // 获取敌方角色名（从当前敌人或日志中的target/actor字段，或从消息文本中提取）
  let enemyName = ''
  // 优先使用target字段（如果actor是我方，target就是敌方）
  if (log.target && log.target !== playerName && !playerNameVariants.includes(log.target)) {
    enemyName = log.target
  } 
  // 如果actor不是我方角色，则actor是敌方
  else if (log.actor && log.actor !== playerName && !playerNameVariants.includes(log.actor) && log.actor !== 'system') {
    enemyName = log.actor
  } 
  // 如果还没有找到，尝试从当前敌人列表中查找（检查消息中是否包含敌人名称）
  if (!enemyName && game.currentEnemies && game.currentEnemies.length > 0) {
    // 从消息文本中查找匹配的敌人名称
    for (const enemy of game.currentEnemies) {
      const enemyNameToCheck = (enemy as any)?.name || ''
      if (enemyNameToCheck && message.includes(enemyNameToCheck)) {
        enemyName = enemyNameToCheck
        break
      }
    }
    // 如果仍然没有找到，使用第一个敌人作为默认值
    if (!enemyName) {
      const currentEnemy = game.currentEnemies[0] as any
      enemyName = currentEnemy?.name || ''
    }
  }
  
  // 如果仍然没有找到，尝试从已记录的敌人映射中查找（即使currentEnemies为空）
  if (!enemyName) {
    // 从消息文本中查找所有已记录的敌人名称
    for (const recordedEnemyName of Object.keys(enemyAttackTypeMap.value)) {
      if (recordedEnemyName && message.includes(recordedEnemyName)) {
        enemyName = recordedEnemyName
        break
      }
    }
  }
  
  // 获取技能名（从日志的action字段或消息中的方括号内容）
  let skillName = ''
  if (log.action && log.action !== '攻击' && log.action !== 'encounter' && log.action !== 'victory' && log.action !== 'defeat' && log.action !== 'loot' && log.action !== 'levelup') {
    skillName = log.action
  }
  
  // 根据敌人名称确定敌人颜色：物理攻击用红色，魔法攻击用蓝色
  // 使用敌人名称到攻击类型的映射，确保同一敌人颜色一致
  let enemyColor = '#ff7777' // 默认红色（物理攻击）
  
  if (enemyName) {
    // 首先检查映射中是否已有该敌人的攻击类型
    let attackType = enemyAttackTypeMap.value[enemyName]
    
    // 如果映射中没有，尝试从当前敌人列表中查找
    if (!attackType && game.currentEnemies && game.currentEnemies.length > 0) {
      const matchedEnemy = game.currentEnemies.find((e: any) => e?.name === enemyName) as any
      if (matchedEnemy) {
        // 优先使用敌人的attackType
        if (matchedEnemy.attackType) {
          attackType = matchedEnemy.attackType.toLowerCase()
        } else if (matchedEnemy.magicAttack > matchedEnemy.physicalAttack && matchedEnemy.magicAttack > 0) {
          // 如果魔法攻击更高，推断为魔法类型
          attackType = 'magic'
        } else {
          attackType = 'physical'
        }
        // 保存到映射中
        enemyAttackTypeMap.value[enemyName] = attackType
      }
    }
    
    // 如果仍然没有找到，尝试从当前日志的damageType推断（仅用于首次遇到）
    if (!attackType && (log.damageType || log.damage_type)) {
      const damageType = (log.damageType || log.damage_type || '').toLowerCase()
      if (damageType === 'magic' || damageType === 'physical') {
        attackType = damageType
        enemyAttackTypeMap.value[enemyName] = attackType
      }
    }
    
    // 根据攻击类型设置颜色
    if (attackType === 'magic') {
      enemyColor = '#7777ff' // 魔法攻击用蓝色
    }
  }
  
  // 当遇到新敌人时，更新映射（从遭遇日志中）
  if (log.logType === 'encounter' && game.currentEnemies && game.currentEnemies.length > 0) {
    game.currentEnemies.forEach((enemy: any) => {
      if (enemy?.name && !enemyAttackTypeMap.value[enemy.name]) {
        if (enemy.attackType) {
          enemyAttackTypeMap.value[enemy.name] = enemy.attackType.toLowerCase()
        } else if (enemy.magicAttack > enemy.physicalAttack && enemy.magicAttack > 0) {
          enemyAttackTypeMap.value[enemy.name] = 'magic'
        } else {
          enemyAttackTypeMap.value[enemy.name] = 'physical'
        }
      }
    })
  }
  
  // 解析消息并添加颜色标记（传入资源颜色用于技能颜色，传入敌人颜色）
  return formatMessageWithColors(message, playerName, playerNameVariants, enemyName, skillName, playerColor, resourceColor, enemyColor)
}

// 格式化消息，为角色名和技能名添加颜色
function formatMessageWithColors(
  message: string,
  playerName: string,
  playerNameVariants: string[],
  enemyName: string,
  skillName: string,
  playerColor: string = '#ffff55', // 默认金色，如果未传入则使用默认值
  resourceColor: string = '#ffffff', // 资源颜色，用于技能颜色
  enemyColor: string = '#ff7777' // 敌人颜色，根据伤害类型传入
): string {
  // 转义HTML特殊字符
  const escapeHtml = (text: string) => {
    const div = document.createElement('div')
    div.textContent = text
    return div.innerHTML
  }
  
  // 检查位置是否在HTML标签内
  const isInHtmlTag = (text: string, pos: number): boolean => {
    const before = text.substring(0, pos)
    const lastOpen = before.lastIndexOf('<')
    const lastClose = before.lastIndexOf('>')
    return lastOpen > lastClose
  }
  
  // 检查位置是否已经在span标签内
  const isInSpanTag = (text: string, pos: number): boolean => {
    const before = text.substring(0, pos)
    const lastSpanOpen = before.lastIndexOf('<span')
    const lastSpanClose = before.lastIndexOf('</span>')
    if (lastSpanOpen === -1) return false
    return lastSpanOpen > lastSpanClose
  }
  
  // 定义颜色（使用传入的职业颜色，敌方使用传入的颜色）
  const normalAttackColor = '#ffffff' // 普通攻击使用白色
  const skillColor = resourceColor // 技能使用资源颜色（与消耗的资源颜色一致）
  
  // 处理消息：保护已有的 HTML 标签，转义纯文本部分
  // 使用占位符保护 HTML 标签
  const htmlPlaceholders: string[] = []
  let processedMessage = message
  let placeholderIndex = 0
  
  // 提取所有 HTML 标签并用占位符替换
  processedMessage = processedMessage.replace(/<[^>]+>/g, (match) => {
    const placeholder = `__HTML_PLACEHOLDER_${placeholderIndex}__`
    htmlPlaceholders[placeholderIndex] = match
    placeholderIndex++
    return placeholder
  })
  
  // 转义纯文本部分
  let formatted = escapeHtml(processedMessage)
  
  // 恢复 HTML 标签
  htmlPlaceholders.forEach((html, index) => {
    formatted = formatted.replace(`__HTML_PLACEHOLDER_${index}__`, html)
  })
  
  // 标记技能名（方括号内的内容）- 优先处理，避免与其他标记冲突
  // 普通攻击使用白色，其他技能使用资源颜色（与消耗的资源颜色一致）
  formatted = formatted.replace(/\[([^\]]+)\]/g, (match, skill) => {
    const isNormalAttack = skill === '普通攻击'
    const color = isNormalAttack ? normalAttackColor : skillColor
    return `<span style="color: ${color}">[${escapeHtml(skill)}]</span>`
  })
  
  // 标记我方角色名（按长度从长到短排序，避免短名称覆盖长名称）
  const sortedPlayerNames = [...playerNameVariants].filter(n => n).sort((a, b) => b.length - a.length)
  sortedPlayerNames.forEach(name => {
    if (name) {
      const regex = new RegExp(escapeRegex(name), 'g')
      // 收集所有匹配位置（从后往前处理，避免索引变化）
      const matches: Array<{ match: string; index: number }> = []
      let match
      while ((match = regex.exec(formatted)) !== null) {
        matches.push({ match: match[0], index: match.index })
      }
      // 从后往前替换
      for (let i = matches.length - 1; i >= 0; i--) {
        const { match: matchText, index } = matches[i]
        if (!isInHtmlTag(formatted, index) && !isInSpanTag(formatted, index)) {
          formatted = formatted.substring(0, index) + 
                      `<span style="color: ${playerColor}">${matchText}</span>` + 
                      formatted.substring(index + matchText.length)
        }
      }
    }
  })
  
  // 标记敌方角色名（避免与已标记的内容冲突）
  // 如果enemyName为空，尝试从消息中提取所有可能的敌人名称
  let enemiesToMark: string[] = []
  if (enemyName) {
    enemiesToMark = [enemyName]
  } else {
    // 首先从当前敌人列表中查找
    if (game.currentEnemies && game.currentEnemies.length > 0) {
      game.currentEnemies.forEach((enemy: any) => {
        const name = enemy?.name || ''
        if (name && formatted.includes(name) && !enemiesToMark.includes(name)) {
          enemiesToMark.push(name)
        }
      })
    }
    // 如果仍然没有找到，从已记录的敌人映射中查找（即使currentEnemies为空）
    if (enemiesToMark.length === 0) {
      for (const recordedEnemyName of Object.keys(enemyAttackTypeMap.value)) {
        if (recordedEnemyName && formatted.includes(recordedEnemyName) && !enemiesToMark.includes(recordedEnemyName)) {
          enemiesToMark.push(recordedEnemyName)
        }
      }
    }
  }
  
  // 为每个敌人名称设置颜色（按长度从长到短排序，避免短名称覆盖长名称）
  const sortedEnemyNames = enemiesToMark.sort((a, b) => b.length - a.length)
  sortedEnemyNames.forEach(name => {
    if (name) {
      // 获取该敌人的攻击类型和颜色
      let color = '#ff7777' // 默认红色
      let attackType = enemyAttackTypeMap.value[name]
      // 如果映射中没有，尝试从当前敌人列表中查找
      if (!attackType && game.currentEnemies && game.currentEnemies.length > 0) {
        const matchedEnemy = game.currentEnemies.find((e: any) => e?.name === name) as any
        if (matchedEnemy) {
          if (matchedEnemy.attackType) {
            attackType = matchedEnemy.attackType.toLowerCase()
          } else if (matchedEnemy.magicAttack > matchedEnemy.physicalAttack && matchedEnemy.magicAttack > 0) {
            attackType = 'magic'
          } else {
            attackType = 'physical'
          }
          enemyAttackTypeMap.value[name] = attackType
        }
      }
      // 如果仍然没有找到，使用默认物理攻击类型（避免颜色丢失）
      if (!attackType) {
        attackType = 'physical'
        // 如果映射中没有，也保存到映射中，避免下次查找
        if (!enemyAttackTypeMap.value[name]) {
          enemyAttackTypeMap.value[name] = attackType
        }
      }
      if (attackType === 'magic') {
        color = '#7777ff' // 魔法攻击用蓝色
      }
      
      const regex = new RegExp(escapeRegex(name), 'g')
      // 收集所有匹配位置（从后往前处理，避免索引变化）
      const matches: Array<{ match: string; index: number }> = []
      let match
      while ((match = regex.exec(formatted)) !== null) {
        matches.push({ match: match[0], index: match.index })
      }
      // 从后往前替换
      for (let i = matches.length - 1; i >= 0; i--) {
        const { match: matchText, index } = matches[i]
        if (!isInHtmlTag(formatted, index) && !isInSpanTag(formatted, index)) {
          const beforeReplace = formatted.substring(index, index + matchText.length)
          // 使用 !important 确保敌人名称颜色不被父元素的颜色覆盖
          formatted = formatted.substring(0, index) + 
                      `<span style="color: ${color} !important">${matchText}</span>` + 
                      formatted.substring(index + matchText.length)
        }
      }
    }
  })
  
  return formatted
}

// 转义正则表达式特殊字符
function escapeRegex(str: string): string {
  return str.replace(/[.*+?^${}()|[\]\\]/g, '\\$&')
}
</script>

<template>
  <div class="game-screen">
    <!-- 如果没有角色数据，显示提示 -->
    <div v-if="!game.character" class="no-character">
      <div class="no-character-message">
        <h2>未找到角色数据</h2>
        <p>正在加载角色信息...</p>
        <p style="color: #888; font-size: 12px; margin-top: 20px;">
          如果长时间未加载，请刷新页面或检查网络连接
        </p>
      </div>
    </div>
    
    <template v-else>
      <!-- 顶部栏 -->
      <div class="game-header">
        <div class="header-left">
          <span class="username">{{ authStore.user?.username || '玩家' }}</span>
          <span class="user-id">{{ authStore.user?.id || '00' }}</span>
        </div>
        <div class="header-right">
          <button class="header-btn" @click="$emit('create-character')">新建角色</button>
          <button class="header-btn" @click="$emit('logout')">登出</button>
        </div>
      </div>

      <!-- 状态栏 -->
      <div class="status-line">
        <span class="stat-zone" @click="openZoneSelector" title="点击切换地图">
          🗺️ {{ currentZoneName }}
        </span>
        <span class="stat-separator">|</span>
        <span class="stat-battle">战斗: {{ (game.battleStatus as any)?.battleCount || (game.battleStatus as any)?.battle_count || 0 }}</span>
        <span class="stat-separator">|</span>
        <span class="stat-kills">击杀: {{ (game.battleStatus as any)?.totalKills || (game.battleStatus as any)?.session_kills || 0 }}</span>
        <span class="stat-separator">|</span>
        <span class="stat-exp">+{{ (game.battleStatus as any)?.totalExp || (game.battleStatus as any)?.session_exp || 0 }} EXP</span>
        <span class="stat-separator">|</span>
        <span class="stat-gold">+{{ (game.battleStatus as any)?.totalGold || (game.battleStatus as any)?.session_gold || 0 }} G</span>
        <span class="battle-status" :class="{ active: game.isRunning }">
          {{ game.isRunning ? '× 战斗中' : '○ 待机' }}
        </span>
      </div>

      <!-- 主内容区 -->
      <div class="game-main">
        <!-- 左侧角色信息面板 -->
        <div class="game-sidebar">
          <!-- 队伍成员列表（显示所有角色，点击查看详情） -->
          <div v-if="charStore.characters.length > 0" class="team-panel">
            <div class="team-panel-title">
              队伍成员 ({{ charStore.characters.length }}/5)
              <span v-if="charStore.characters.length >= 5" class="team-panel-full">已满</span>
            </div>
            <div class="team-characters">
              <div
                v-for="char in charStore.characters"
                :key="char.id"
                class="team-character-card"
                :class="{ dead: char.isDead }"
                @click="showDetail(char)"
              >
                <div 
                  class="team-character-name"
                  :style="{ 
                    color: getClassColor(char.classId || ''),
                    textShadow: `0 0 8px ${getClassColor(char.classId || '')}`
                  }"
                >
                  {{ char.name }}
                </div>
                <div class="team-character-level">
                  Lv.{{ char.level }} {{ getClassName(char.classId || '') }}
                </div>
                <div class="team-character-hp">
                  <div class="team-character-hp-label">HP:</div>
                  <div class="team-character-hp-bar">
                    <div 
                      class="team-character-hp-fill" 
                      :style="{ 
                        width: getHpPercent(char) + '%' 
                      }"
                    ></div>
                  </div>
                  <div class="team-character-hp-value">
                    {{ char.hp || 0 }}/{{ char.maxHp || 100 }}
                  </div>
                </div>
                <div class="team-character-resource">
                  <div class="team-character-resource-label">{{ getResourceTypeName(char) }}:</div>
                  <div class="team-character-resource-bar">
                    <div 
                      class="team-character-resource-fill" 
                      :style="{ 
                        width: getMpPercent(char) + '%',
                        background: getResourceTypeName(char) === '怒气' ? 'linear-gradient(90deg, #ff4444, #ff6666)' : 
                                    getResourceTypeName(char) === '能量' ? 'linear-gradient(90deg, #ffd700, #ffed4e)' :
                                    'linear-gradient(90deg, #3d85c6, #5ba3d6)'
                      }"
                    ></div>
                  </div>
                  <div class="team-character-resource-value">
                    {{ (char as any).resource || (char as any).mp || 0 }}/{{ (char as any).maxResource || (char as any).max_resource || (char as any).max_mp || 100 }}
                  </div>
                </div>
                <!-- Buff/Debuff显示 -->
                <div v-if="char.buffs && char.buffs.length > 0" class="team-character-buffs">
                  <div
                    v-for="buff in char.buffs"
                    :key="buff.effectId"
                    class="buff-icon"
                    :class="{ 'buff-positive': buff.isBuff, 'buff-negative': !buff.isBuff }"
                    :data-tooltip="getBuffTooltip(buff)"
                    data-tooltip-right
                    @mouseenter="handleBuffTooltip($event, buff)"
                    @mouseleave="hideBuffTooltip"
                  >
                    {{ (buff.name || '?').charAt(0) }}
                  </div>
                </div>
                <div v-if="char.isDead" class="team-character-dead">
                  死亡中...
                </div>
              </div>
            </div>
          </div>
          
          <!-- 空状态提示 -->
          <div v-else class="no-characters-hint">
            <div class="hint-text">还没有角色</div>
            <button class="hint-btn" @click="$emit('create-character')">
              创建角色
            </button>
          </div>
        </div>

        <!-- 中间战斗日志区域 -->
        <div class="game-content">
          <!-- 敌人信息面板（固定在顶部） -->
          <div v-if="game.currentEnemies && game.currentEnemies.length > 0" class="enemies-panel">
            <div 
              v-for="(enemy, index) in game.currentEnemies" 
              :key="index"
              class="enemy-info"
              :class="{ 
                'enemy-dead': (enemy as any)?.hp <= 0,
                'enemy-normal': !(enemy as any)?.type || (enemy as any)?.type === 'normal',
                'enemy-elite': (enemy as any)?.type === 'elite',
                'enemy-boss': (enemy as any)?.type === 'boss'
              }"
            >
              <span 
                class="enemy-name"
                :style="{
                  color: getEnemyNameColor(enemy)
                }"
              >
                <span v-if="(enemy as any)?.type === 'elite'" class="enemy-rarity-icon">⭐</span>
                <span v-else-if="(enemy as any)?.type === 'boss'" class="enemy-rarity-icon">👑</span>
                <span v-else class="enemy-rarity-icon">⚔</span>
                {{ (enemy as any)?.name || '未知敌人' }} (Lv.{{ (enemy as any)?.level || 1 }})
              </span>
              <div class="enemy-hp">
                <span class="enemy-hp-label">HP:</span>
                <div class="enemy-bar">
                  <div class="enemy-bar-fill" :style="{ width: getEnemyHpPercent(enemy) + '%' }"></div>
                </div>
                <span class="enemy-hp-value">
                  {{ Math.max(0, (enemy as any)?.hp || 0) }}/{{ (enemy as any)?.maxHp || (enemy as any)?.max_hp || (enemy as any)?.hp || 100 }}
                </span>
              </div>
            </div>
          </div>
          
          <div class="terminal-content" ref="logContainer">
            <!-- 战斗日志 -->
            <div class="battle-log">
              <div 
                v-for="(log, index) in game.battleLogs" 
                :key="index"
                class="log-line"
              >
                <span class="log-time">[{{ formatLogTime(log) }}]</span>
                <span 
                  class="log-message"
                  :class="getLogClass(log.type || log.logType || 'info')"
                  :style="{ color: log.color || '#00ff00' }"
                  v-html="formatLogMessage(log)"
                ></span>
              </div>
              <div class="log-line" v-if="game.isRunning">
                <span class="log-time"></span>
                <span class="log-message" style="color: #00ff00">
                  等待下一回合...<span class="cursor"></span>
                </span>
              </div>
            </div>
          </div>

          <!-- 控制按钮 -->
          <div class="control-bar">
            <button 
              class="cmd-btn" 
              :class="{ active: game.isRunning }"
              @click="game.toggleBattle"
            >
              {{ game.isRunning ? '[停止挂机]' : '[开始挂机]' }}
            </button>
            <button 
              class="cmd-btn" 
              @click="openStrategyEditor"
              :disabled="!charStore.activeCharacter"
            >
              [S] 策略
            </button>
            <button 
              class="cmd-btn" 
              @click="openStatsPanel"
            >
              [T] 统计
            </button>
            <button class="cmd-btn" disabled>
              [E] 装备
            </button>
            <button 
              class="cmd-btn" 
              @click="openZoneSelector"
            >
              [M] 地图
            </button>
          </div>
        </div>
      </div>

      <!-- 底部聊天面板 -->
      <ChatPanel />

      <!-- 策略编辑器 -->
      <StrategyEditor 
        v-if="showStrategyEditor && strategyCharacterId"
        :character-id="strategyCharacterId"
        :character-skills="strategyCharacterSkills"
        @close="closeStrategyEditor"
      />

      <!-- 战斗统计面板 -->
      <BattleStatsPanel 
        v-if="showStatsPanel"
        @close="closeStatsPanel"
      />

      <!-- 地图选择面板 -->
      <div v-if="showZoneSelector" class="zone-selector-overlay" @click.self="closeZoneSelector">
        <div class="zone-selector">
          <div class="zone-selector-header">
            <h2>选择地图</h2>
            <button class="close-btn" @click="closeZoneSelector">×</button>
          </div>
          
          <div v-if="zoneError" class="zone-error">{{ zoneError }}</div>
          
          <div v-if="loadingZones" class="zone-loading">加载中...</div>
          
          <div v-else-if="availableZones.length === 0" class="zone-loading">
            没有可用的地图，请检查你的等级和阵营
          </div>
          
          <div v-else class="zone-list">
            <div 
              v-for="zone in availableZones" 
              :key="zone.id"
              class="zone-item"
              :class="{ 
                'zone-current': game.battleStatus?.currentZoneId === zone.id,
                'zone-locked': (charStore.activeCharacter && charStore.activeCharacter.level < zone.minLevel) || !isZoneUnlocked(zone)
              }"
              @click="(charStore.activeCharacter && charStore.activeCharacter.level >= zone.minLevel && isZoneUnlocked(zone)) ? selectZone(zone.id) : null"
            >
              <div class="zone-name">{{ zone.name }}</div>
              <div class="zone-info">
                <span class="zone-level">等级: {{ zone.minLevel }}-{{ zone.maxLevel }}</span>
                <span class="zone-faction" :class="`faction-${zone.faction || 'neutral'}`">
                  {{ zone.faction === 'alliance' ? '联盟' : zone.faction === 'horde' ? '部落' : 'PVP' }}
                </span>
                <span class="zone-multiplier">倍率: {{ zone.expMulti }}x</span>
              </div>
              <div class="zone-description">{{ zone.description }}</div>
              
              <!-- 探索度进度显示 -->
              <div v-if="getUnlockProgress(zone)" class="zone-exploration-progress">
                <div class="exploration-label">
                  解锁条件: 在 <strong>{{ getUnlockProgress(zone)?.unlockZoneName }}</strong> 探索度达到 {{ getUnlockProgress(zone)?.required }}
                </div>
                <div class="exploration-bar">
                  <div 
                    class="exploration-fill" 
                    :style="{ width: `${Math.min(100, (getUnlockProgress(zone)?.current || 0) / getUnlockProgress(zone)?.required * 100)}%` }"
                  ></div>
                  <span class="exploration-text">
                    {{ getUnlockProgress(zone)?.current }} / {{ getUnlockProgress(zone)?.required }}
                  </span>
                </div>
              </div>
              
              <!-- 当前地图探索度显示 -->
              <div v-if="game.explorations[zone.id]" class="zone-current-exploration">
                当前探索度: {{ game.explorations[zone.id].exploration }} (击杀: {{ game.explorations[zone.id].kills }})
              </div>
              
              <!-- 锁定提示 -->
              <div v-if="charStore.activeCharacter && charStore.activeCharacter.level < zone.minLevel" class="zone-locked-hint">
                需要等级 {{ zone.minLevel }}
              </div>
              <div v-if="!isZoneUnlocked(zone) && zone.unlockZoneId" class="zone-locked-hint">
                未解锁（探索度不足）
              </div>
            </div>
          </div>
        </div>
      </div>
    </template>
    
    <!-- 技能选择弹窗：被动/主动选择 -->
    <div v-if="skillSelection" class="skill-select-overlay">
      <div class="skill-select-modal">
        <div class="skill-select-header">
          <div class="skill-select-title">
            {{ skillSelection.selectionType === 'passive' ? '被动技能选择' : '主动技能选择' }}
            <span class="skill-select-level">Lv.{{ skillSelection.level }}</span>
          </div>
          <div class="skill-select-sub">请选择强化已有或学习新技能</div>
        </div>

        <div class="skill-select-columns">
          <div class="skill-select-column">
            <div class="skill-select-column-title">强化已有</div>
            <div v-if="skillSelection.selectionType === 'passive'">
              <div v-if="skillSelection.upgradePassives && skillSelection.upgradePassives.length > 0" class="skill-select-list">
                <button
                  v-for="item in skillSelection.upgradePassives"
                  :key="item.passiveId"
                  class="skill-select-item"
                  :disabled="selectingSkill"
                  @click="submitSelection({ passiveId: item.passiveId, isUpgrade: true })"
                >
                  <div class="skill-select-name">{{ getPassiveName(item) }}</div>
                  <div class="skill-select-desc">当前等级: {{ item.level }}/5</div>
                </button>
              </div>
              <div v-else class="skill-select-empty">暂无可强化的被动</div>
            </div>

            <div v-else>
              <div v-if="skillSelection.upgradeSkills && skillSelection.upgradeSkills.length > 0" class="skill-select-list">
                <button
                  v-for="item in skillSelection.upgradeSkills"
                  :key="item.skillId"
                  class="skill-select-item"
                  :disabled="selectingSkill"
                  @click="submitSelection({ skillId: item.skillId, isUpgrade: true })"
                >
                  <div class="skill-select-name">{{ getSkillName(item) }}</div>
                  <div class="skill-select-desc">当前等级: {{ item.skillLevel }}/5</div>
                </button>
              </div>
              <div v-else class="skill-select-empty">暂无可强化的主动技能</div>
            </div>
          </div>

          <div class="skill-select-column">
            <div class="skill-select-column-title">学习新技能</div>
            <div v-if="skillSelection.selectionType === 'passive'">
              <div v-if="skillSelection.newPassives && skillSelection.newPassives.length > 0" class="skill-select-list">
                <button
                  v-for="item in skillSelection.newPassives"
                  :key="item.id"
                  class="skill-select-item"
                  :disabled="selectingSkill"
                  @click="submitSelection({ passiveId: item.id, isUpgrade: false })"
                >
                  <div class="skill-select-name">{{ item.name }}</div>
                  <div class="skill-select-desc">{{ item.description }}</div>
                </button>
              </div>
              <div v-else class="skill-select-empty">暂无新的被动可学习</div>
            </div>

            <div v-else>
              <div v-if="skillSelection.newSkills && skillSelection.newSkills.length > 0" class="skill-select-list">
                <button
                  v-for="item in skillSelection.newSkills"
                  :key="item.id"
                  class="skill-select-item"
                  :disabled="selectingSkill"
                  @click="submitSelection({ skillId: item.id, isUpgrade: false })"
                >
                  <div class="skill-select-name">{{ item.name }}</div>
                  <div class="skill-select-desc">{{ item.description }}</div>
                </button>
              </div>
              <div v-else class="skill-select-empty">暂无新的主动技能可学习</div>
            </div>
          </div>
        </div>

        <div class="skill-select-footer">
          <span class="skill-select-error" v-if="selectionError">{{ selectionError }}</span>
          <span class="skill-select-hint" v-else>本次选择后，下一次机会会在相应等级里程碑出现</span>
          <div class="skill-select-actions">
            <button class="skill-select-btn" :disabled="selectingSkill || skillSelectionLoading" @click="refreshSkillSelection(true)">
              刷新选项
            </button>
          </div>
        </div>
      </div>
    </div>
    
    <!-- 角色详情弹窗 -->
    <div v-if="showCharacterDetail && selectedCharacter" class="character-detail-modal" @click.self="closeDetail">
      <div class="character-detail-content">
        <div class="character-detail-header">
          <div 
            class="character-detail-name"
            :style="{ 
              color: getClassColor(selectedCharacter.classId || ''),
              textShadow: `0 0 10px ${getClassColor(selectedCharacter.classId || '')}`
            }"
          >
            {{ selectedCharacter.name }}
          </div>
          <button class="character-detail-close" @click="closeDetail">×</button>
        </div>
        
        <div class="character-detail-level">
          Lv.{{ selectedCharacter.level }} {{ getClassName(selectedCharacter.classId || '') }}
        </div>
        
        <!-- 进度条 -->
        <div class="character-detail-progress">
          <div class="character-detail-progress-item">
            <div class="character-detail-progress-label">生命值</div>
            <div class="character-detail-progress-bar hp-bar">
              <div class="character-detail-progress-fill" :style="{ width: getHpPercent(selectedCharacter) + '%' }"></div>
            </div>
            <div class="character-detail-progress-text">
              {{ selectedCharacter.hp || 0 }}/{{ selectedCharacter.maxHp || 100 }}
            </div>
          </div>
          
          <div class="character-detail-progress-item">
            <div class="character-detail-progress-label">{{ getResourceTypeName(selectedCharacter) }}</div>
            <div class="character-detail-progress-bar mp-bar">
              <div class="character-detail-progress-fill" :style="{ width: getMpPercent(selectedCharacter) + '%' }"></div>
            </div>
            <div class="character-detail-progress-text">
              {{ (selectedCharacter as any).resource || (selectedCharacter as any).mp || 0 }}/{{ (selectedCharacter as any).maxResource || (selectedCharacter as any).max_resource || (selectedCharacter as any).max_mp || 100 }}
            </div>
          </div>
          
          <div class="character-detail-progress-item">
            <div class="character-detail-progress-label">经验值</div>
            <div class="character-detail-progress-bar exp-bar">
              <div class="character-detail-progress-fill" :style="{ width: getExpPercent(selectedCharacter) + '%' }"></div>
            </div>
            <div class="character-detail-progress-text">
              {{ selectedCharacter.exp || 0 }}/{{ selectedCharacter.expToNext || selectedCharacter.exp_to_next || 100 }}
            </div>
          </div>
        </div>

        <!-- 属性 -->
        <div class="character-detail-unspent" v-if="selectedCharacter.unspentPoints !== undefined">
          剩余点数: {{ selectedCharacter.unspentPoints || 0 }}
        </div>
        <div class="character-detail-stats">
          <div
            class="character-detail-stat"
            @mouseenter="handleAttrTooltip($event, 'strength')"
            @mouseleave="hideAttrTooltip"
          >
            <span class="character-detail-stat-label">力量</span>
            <span class="character-detail-stat-value">{{ selectedCharacter.strength || 0 }}</span>
            <button 
              class="stat-allocate-btn" 
              @click.stop="allocateStat('strength')" 
              :disabled="!selectedCharacter.unspentPoints || allocating"
            >+</button>
          </div>
          <div
            class="character-detail-stat"
            @mouseenter="handleAttrTooltip($event, 'agility')"
            @mouseleave="hideAttrTooltip"
          >
            <span class="character-detail-stat-label">敏捷</span>
            <span class="character-detail-stat-value">{{ selectedCharacter.agility || 0 }}</span>
            <button 
              class="stat-allocate-btn" 
              @click.stop="allocateStat('agility')" 
              :disabled="!selectedCharacter.unspentPoints || allocating"
            >+</button>
          </div>
          <div
            class="character-detail-stat"
            @mouseenter="handleAttrTooltip($event, 'intellect')"
            @mouseleave="hideAttrTooltip"
          >
            <span class="character-detail-stat-label">智力</span>
            <span class="character-detail-stat-value">{{ selectedCharacter.intellect || 0 }}</span>
            <button 
              class="stat-allocate-btn" 
              @click.stop="allocateStat('intellect')" 
              :disabled="!selectedCharacter.unspentPoints || allocating"
            >+</button>
          </div>
          <div
            class="character-detail-stat"
            @mouseenter="handleAttrTooltip($event, 'stamina')"
            @mouseleave="hideAttrTooltip"
          >
            <span class="character-detail-stat-label">耐力</span>
            <span class="character-detail-stat-value">{{ selectedCharacter.stamina || 0 }}</span>
            <button 
              class="stat-allocate-btn" 
              @click.stop="allocateStat('stamina')" 
              :disabled="!selectedCharacter.unspentPoints || allocating"
            >+</button>
          </div>
          <div
            class="character-detail-stat"
            @mouseenter="handleAttrTooltip($event, 'spirit')"
            @mouseleave="hideAttrTooltip"
          >
            <span class="character-detail-stat-label">精神</span>
            <span class="character-detail-stat-value">{{ selectedCharacter.spirit || 0 }}</span>
            <button 
              class="stat-allocate-btn" 
              @click.stop="allocateStat('spirit')" 
              :disabled="!selectedCharacter.unspentPoints || allocating"
            >+</button>
          </div>
        </div>

        <!-- 战斗统计 -->
        <div class="character-detail-combat-stats">
          <div
            class="character-detail-combat-stat"
            @mouseenter="handleCombatTooltip($event, getCombatStatTooltip('physicalAttack'))"
            @mouseleave="hideCombatTooltip"
          >
            <span class="character-detail-combat-stat-label">物理攻击</span>
            <span class="character-detail-combat-stat-value">{{ displayedStats.physicalAttack }}</span>
          </div>
          <div
            class="character-detail-combat-stat"
            @mouseenter="handleCombatTooltip($event, getCombatStatTooltip('magicAttack'))"
            @mouseleave="hideCombatTooltip"
          >
            <span class="character-detail-combat-stat-label">魔法攻击</span>
            <span class="character-detail-combat-stat-value">{{ displayedStats.magicAttack }}</span>
          </div>
          <div
            class="character-detail-combat-stat"
            @mouseenter="handleCombatTooltip($event, getCombatStatTooltip('physicalDefense'))"
            @mouseleave="hideCombatTooltip"
          >
            <span class="character-detail-combat-stat-label">物理防御</span>
            <span class="character-detail-combat-stat-value">{{ displayedStats.physicalDefense }}</span>
          </div>
          <div
            class="character-detail-combat-stat"
            @mouseenter="handleCombatTooltip($event, getCombatStatTooltip('magicDefense'))"
            @mouseleave="hideCombatTooltip"
          >
            <span class="character-detail-combat-stat-label">魔法防御</span>
            <span class="character-detail-combat-stat-value">{{ displayedStats.magicDefense }}</span>
          </div>
          <div
            class="character-detail-combat-stat"
            @mouseenter="handleCombatTooltip($event, getCombatStatTooltip('physCritRate'))"
            @mouseleave="hideCombatTooltip"
          >
            <span class="character-detail-combat-stat-label">物理暴击</span>
            <span class="character-detail-combat-stat-value">{{ (displayedStats.physCritRate * 100).toFixed(1) }}%</span>
          </div>
          <div
            class="character-detail-combat-stat"
            @mouseenter="handleCombatTooltip($event, getCombatStatTooltip('physCritDamage'))"
            @mouseleave="hideCombatTooltip"
          >
            <span class="character-detail-combat-stat-label">物暴伤害</span>
            <span class="character-detail-combat-stat-value">{{ (displayedStats.physCritDamage * 100).toFixed(0) }}%</span>
          </div>
          <div
            class="character-detail-combat-stat"
            @mouseenter="handleCombatTooltip($event, getCombatStatTooltip('spellCritRate'))"
            @mouseleave="hideCombatTooltip"
          >
            <span class="character-detail-combat-stat-label">法术暴击</span>
            <span class="character-detail-combat-stat-value">{{ (displayedStats.spellCritRate * 100).toFixed(1) }}%</span>
          </div>
          <div
            class="character-detail-combat-stat"
            @mouseenter="handleCombatTooltip($event, getCombatStatTooltip('spellCritDamage'))"
            @mouseleave="hideCombatTooltip"
          >
            <span class="character-detail-combat-stat-label">法暴伤害</span>
            <span class="character-detail-combat-stat-value">{{ (displayedStats.spellCritDamage * 100).toFixed(0) }}%</span>
          </div>
          <div
            class="character-detail-combat-stat"
            @mouseenter="handleCombatTooltip($event, getCombatStatTooltip('dodgeRate'))"
            @mouseleave="hideCombatTooltip"
          >
            <span class="character-detail-combat-stat-label">闪避率</span>
            <span class="character-detail-combat-stat-value">{{ (displayedStats.dodgeRate * 100).toFixed(1) }}%</span>
          </div>
        </div>

        <!-- Buff/Debuff显示 -->
        <div v-if="selectedCharacter.buffs && selectedCharacter.buffs.length > 0" class="character-detail-buffs">
          <div class="character-detail-section-title">Buff/Debuff</div>
          <div class="character-detail-buffs-list">
            <div
              v-for="buff in selectedCharacter.buffs"
              :key="buff.effectId"
              class="character-detail-buff-item"
              :class="{ 'buff-positive': buff.isBuff, 'buff-negative': !buff.isBuff }"
            >
              <div class="buff-item-name">{{ buff.name }}</div>
              <div class="buff-item-desc">{{ buff.description || '' }}</div>
              <div class="buff-item-duration">剩余 {{ buff.duration }} 回合</div>
            </div>
          </div>
        </div>

        <!-- 被动技能列表 -->
        <div class="character-detail-skills">
          <div class="character-detail-section-title">被动技能 ({{ passiveSkills.length }})</div>
          <div v-if="loadingSkills" class="character-detail-loading">加载中...</div>
          <div v-else-if="passiveSkills.length === 0" class="character-detail-no-skills">暂无被动技能</div>
          <div v-else class="character-detail-skills-list">
            <div
              v-for="passive in passiveSkills"
              :key="passive.id || passive.passiveId"
              class="character-detail-skill-item"
              :data-tooltip="getPassiveTooltip(passive)"
              @mouseenter="handlePassiveTooltip($event, passive)"
              @mouseleave="hideSkillTooltip"
            >
              <div class="skill-item-main">
                <span class="skill-item-name">{{ passive.passive?.name || passive.passiveId }}</span>
                <span class="skill-item-level">Lv.{{ passive.level || 1 }}</span>
              </div>
              <div class="skill-item-meta">
                <span class="skill-item-cost">{{ formatPassiveEffect(passive) }}</span>
              </div>
            </div>
          </div>
        </div>

        <!-- 技能列表 -->
        <div class="character-detail-skills">
          <div class="character-detail-section-title">已掌握的技能 ({{ characterSkills.length }})</div>
          <div v-if="loadingSkills" class="character-detail-loading">加载中...</div>
          <div v-else-if="characterSkills.length === 0" class="character-detail-no-skills">暂无技能</div>
          <div v-else class="character-detail-skills-list">
            <div
              v-for="skill in characterSkills"
              :key="skill.id"
              class="character-detail-skill-item"
              :data-tooltip="getSkillTooltip(skill)"
              @mouseenter="handleSkillTooltip($event, skill)"
              @mouseleave="hideSkillTooltip"
            >
              <div class="skill-item-main">
                <span class="skill-item-name">{{ skill.skill?.name || skill.skillId }}</span>
                <span class="skill-item-level">Lv.{{ skill.skillLevel }}</span>
              </div>
              <div class="skill-item-meta">
                <span v-if="skill.skill?.resourceCost" class="skill-item-cost">
                  {{ skill.skill.resourceCost }}{{ getResourceTypeName(selectedCharacter) === '怒气' ? '怒' : getResourceTypeName(selectedCharacter) === '能量' ? '能' : 'MP' }}
                </span>
                <span v-if="skill.skill?.cooldown" class="skill-item-cooldown">
                  CD:{{ skill.skill.cooldown }}
                </span>
              </div>
            </div>
          </div>
        </div>

        <!-- 总结统计 -->
        <div class="character-detail-summary">
          <div class="character-detail-summary-kills">击杀: {{ selectedCharacter.totalKills || 0 }}</div>
          <div class="character-detail-summary-deaths">死亡: {{ selectedCharacter.totalDeaths || 0 }}</div>
        </div>
      </div>
    </div>
  </div>
</template>

<style>
/* 使用全局样式，terminal.css 中已定义大部分样式 */
.game-screen {
  display: flex;
  flex-direction: column;
  height: 100vh;
  overflow: hidden;
}

/* 顶部栏 */
.game-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 8px 16px;
  border-bottom: 1px solid var(--terminal-green);
  background: rgba(0, 50, 0, 0.3);
}

.header-left {
  display: flex;
  align-items: center;
  gap: 8px;
}

.username {
  color: var(--terminal-green);
  font-weight: bold;
}

.user-id {
  color: var(--terminal-gray);
  font-size: 12px;
}

.header-right {
  display: flex;
  gap: 12px;
}

.header-btn {
  background: transparent;
  border: 1px solid var(--terminal-green);
  color: var(--terminal-green);
  padding: 4px 12px;
  font-family: inherit;
  font-size: 12px;
  cursor: pointer;
  transition: all 0.2s;
}

.header-btn:hover {
  background: var(--terminal-green);
  color: var(--terminal-bg);
}

/* 状态栏 */
.status-line {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 6px 16px;
  border-bottom: 1px solid var(--terminal-gray);
  background: rgba(0, 0, 0, 0.2);
  font-size: 12px;
}

.stat-zone {
  color: var(--terminal-cyan);
  cursor: pointer;
  text-decoration: underline;
  text-decoration-color: var(--terminal-cyan);
  transition: all 0.2s;
}

.stat-zone:hover {
  color: var(--terminal-green);
  text-decoration-color: var(--terminal-green);
}

.stat-battle {
  color: var(--terminal-green);
}

.stat-kills {
  color: var(--terminal-red);
}

.stat-exp {
  color: var(--terminal-cyan);
}

.stat-gold {
  color: var(--terminal-gold);
}

.stat-separator {
  color: var(--terminal-gray);
  opacity: 0.5;
}

.battle-status {
  margin-left: auto;
  color: var(--terminal-gray);
}

.battle-status.active {
  color: var(--terminal-red);
}

/* 主内容区 */
.game-main {
  display: flex;
  flex: 1;
  overflow: hidden;
}

.game-sidebar {
  width: 280px;
  border-right: 2px solid var(--border-color);
  padding: 15px;
  background: rgba(0, 0, 0, 0.3);
  overflow-y: auto;
  display: flex;
  flex-direction: column;
  gap: 15px;
}

/* 空状态提示 */
.no-characters-hint {
  text-align: center;
  padding: 40px 20px;
  color: var(--terminal-gray);
}

.hint-text {
  font-size: 14px;
  margin-bottom: 20px;
}

.hint-btn {
  background: transparent;
  border: 1px solid var(--terminal-green);
  color: var(--terminal-green);
  padding: 8px 16px;
  font-family: inherit;
  font-size: 12px;
  cursor: pointer;
  transition: all 0.2s;
}

.hint-btn:hover {
  background: var(--terminal-green);
  color: var(--terminal-bg);
}

/* 队伍面板 */
.team-panel {
  border: 1px solid var(--border-color);
  background: rgba(0, 0, 0, 0.5);
  padding: 10px;
  flex: 1;
  min-height: 0;
  display: flex;
  flex-direction: column;
}

.team-panel-title {
  color: var(--terminal-cyan);
  font-size: 12px;
  margin-bottom: 8px;
  text-align: center;
  padding-bottom: 8px;
  border-bottom: 1px solid var(--text-dim);
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 8px;
}

.team-panel-full {
  color: var(--terminal-red);
  font-size: 10px;
  padding: 2px 6px;
  border: 1px solid var(--terminal-red);
  border-radius: 2px;
}

.team-characters {
  display: flex;
  flex-direction: column;
  gap: 8px;
  flex: 1;
  overflow-y: auto;
  overflow-x: hidden;
}

/* 队伍成员滚动条样式 */
.team-characters::-webkit-scrollbar {
  width: 6px;
}

.team-characters::-webkit-scrollbar-track {
  background: rgba(0, 0, 0, 0.3);
}

.team-characters::-webkit-scrollbar-thumb {
  background: var(--terminal-gray);
  border-radius: 3px;
}

.team-characters::-webkit-scrollbar-thumb:hover {
  background: var(--terminal-green);
}

.team-character-card {
  border: 1px solid var(--text-dim);
  padding: 8px;
  background: rgba(0, 0, 0, 0.3);
  transition: all 0.2s;
  cursor: pointer;
}

.team-character-card:hover {
  border-color: var(--terminal-green);
  background: rgba(0, 255, 0, 0.05);
  box-shadow: 0 0 8px rgba(0, 255, 0, 0.2);
}

.team-character-card.dead {
  opacity: 0.6;
  border-color: var(--terminal-red);
}

.team-character-name {
  font-size: 13px;
  font-weight: bold;
  margin-bottom: 4px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.team-character-level {
  color: var(--text-secondary);
  font-size: 11px;
  margin-bottom: 6px;
}

.team-character-hp {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 10px;
}

.team-character-hp-label {
  color: var(--text-gray);
  min-width: 24px;
}

.team-character-hp-bar {
  flex: 1;
  height: 8px;
  background: var(--bg-color);
  border: 1px solid var(--text-green);
  overflow: hidden;
}

.team-character-hp-fill {
  height: 100%;
  background: linear-gradient(90deg, #00ff00, #44ff44);
  transition: width 0.3s ease;
}

.team-character-hp-value {
  color: var(--text-green);
  font-size: 10px;
  white-space: nowrap;
  min-width: 50px;
  text-align: right;
}

.team-character-resource {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 10px;
  margin-top: 4px;
}

.team-character-resource-label {
  color: var(--text-gray);
  font-size: 10px;
  min-width: 32px;
}

.team-character-resource-bar {
  flex: 1;
  height: 6px;
  background: var(--bg-color);
  border: 1px solid var(--text-dim);
  overflow: hidden;
}

.team-character-resource-fill {
  height: 100%;
  transition: width 0.3s ease;
}

.team-character-resource-value {
  color: var(--text-cyan);
  font-size: 10px;
  white-space: nowrap;
  min-width: 50px;
  text-align: right;
}

.team-character-buffs {
  display: flex;
  flex-wrap: wrap;
  gap: 4px;
  margin-top: 6px;
  min-height: 18px;
}

.buff-icon {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 18px;
  height: 18px;
  font-size: 10px;
  font-weight: bold;
  border: 1px solid;
  border-radius: 2px;
  cursor: help;
  transition: all 0.2s;
}

.buff-icon:hover {
  transform: scale(1.2);
  z-index: 10;
}

.buff-positive {
  background: rgba(0, 255, 0, 0.2);
  border-color: #00ff00;
  color: #00ff00;
}

.buff-negative {
  background: rgba(255, 0, 0, 0.2);
  border-color: #ff4444;
  color: #ff4444;
}

.team-character-dead {
  color: var(--terminal-red);
  font-size: 10px;
  margin-top: 4px;
  text-align: center;
}

/* 角色详情弹窗 */
.character-detail-modal {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.8);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 10000;
  padding: 20px;
}

.skill-select-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.75);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 12000;
  padding: 20px;
}

.skill-select-modal {
  background: rgba(0, 30, 0, 0.95);
  border: 2px solid var(--terminal-green);
  padding: 20px;
  max-width: 900px;
  width: 100%;
  color: var(--text-primary);
  box-shadow: 0 0 25px rgba(0, 255, 0, 0.3);
}

.skill-select-header {
  margin-bottom: 16px;
}

.skill-select-title {
  font-size: 18px;
  color: var(--terminal-gold);
  display: flex;
  align-items: center;
  gap: 10px;
}

.skill-select-level {
  font-size: 12px;
  color: var(--text-secondary);
}

.skill-select-sub {
  color: var(--text-secondary);
  font-size: 13px;
  margin-top: 4px;
}

.skill-select-columns {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 12px;
}

.skill-select-column {
  background: rgba(0, 0, 0, 0.35);
  border: 1px solid var(--text-dim);
  padding: 12px;
  min-height: 220px;
}

.skill-select-column-title {
  font-size: 14px;
  color: var(--terminal-cyan);
  margin-bottom: 8px;
}

.skill-select-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.skill-select-item {
  width: 100%;
  text-align: left;
  background: rgba(0, 0, 0, 0.5);
  border: 1px solid var(--terminal-green);
  color: var(--text-primary);
  padding: 10px;
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.2s;
}

.skill-select-item:hover:not(:disabled) {
  background: rgba(0, 255, 0, 0.08);
  transform: translateY(-1px);
}

.skill-select-item:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.skill-select-name {
  font-weight: bold;
  margin-bottom: 4px;
  color: var(--terminal-gold);
}

.skill-select-desc {
  font-size: 12px;
  color: var(--text-secondary);
}

.skill-select-empty {
  color: var(--text-secondary);
  font-size: 13px;
  padding: 8px 0;
}

.skill-select-footer {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-top: 12px;
  gap: 12px;
}

.skill-select-error {
  color: var(--terminal-red);
  font-size: 13px;
}

.skill-select-hint {
  color: var(--text-secondary);
  font-size: 12px;
}

.skill-select-actions {
  display: flex;
  gap: 8px;
}

.skill-select-btn {
  background: var(--terminal-green);
  color: var(--terminal-bg);
  border: none;
  padding: 8px 14px;
  border-radius: 4px;
  cursor: pointer;
  font-weight: bold;
  transition: all 0.2s;
}

.skill-select-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.skill-select-btn:not(:disabled):hover {
  background: #00ff9a;
}

.character-detail-content {
  background: rgba(0, 20, 0, 0.95);
  border: 2px solid var(--terminal-green);
  padding: 20px;
  max-width: 500px;
  width: 100%;
  max-height: 90vh;
  overflow-y: auto;
  box-shadow: 0 0 30px rgba(0, 255, 0, 0.3);
  position: relative;
}

.character-detail-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 15px;
  padding-bottom: 10px;
  border-bottom: 1px solid var(--text-dim);
}

.character-detail-name {
  font-family: var(--font-pixel);
  font-size: 20px;
  font-weight: bold;
}

.character-detail-close {
  background: transparent;
  border: 1px solid var(--terminal-gray);
  color: var(--terminal-gray);
  width: 28px;
  height: 28px;
  font-size: 20px;
  cursor: pointer;
  transition: all 0.2s;
  display: flex;
  align-items: center;
  justify-content: center;
}

.character-detail-close:hover {
  border-color: var(--terminal-red);
  color: var(--terminal-red);
}

.character-detail-level {
  color: var(--text-cyan);
  margin-bottom: 15px;
  font-size: 14px;
  text-align: center;
}

.character-detail-progress {
  margin-bottom: 15px;
}

.character-detail-progress-item {
  margin-bottom: 12px;
}

.character-detail-progress-label {
  color: var(--text-secondary);
  font-size: 12px;
  margin-bottom: 4px;
}

.character-detail-progress-bar {
  width: 100%;
  height: 14px;
  background: var(--bg-color);
  border: 1px solid var(--text-dim);
  position: relative;
  overflow: hidden;
  margin-bottom: 4px;
}

.character-detail-progress-fill {
  height: 100%;
  transition: width 0.3s ease;
}

.character-detail-progress-bar.hp-bar .character-detail-progress-fill {
  background: linear-gradient(90deg, #00ff00, #44ff44);
  box-shadow: 0 0 10px rgba(0, 255, 0, 0.5);
}

.character-detail-progress-bar.mp-bar .character-detail-progress-fill {
  background: linear-gradient(90deg, #ff4444, #ff6666);
  box-shadow: 0 0 10px rgba(255, 68, 68, 0.5);
}

.character-detail-progress-bar.exp-bar .character-detail-progress-fill {
  background: linear-gradient(90deg, #ffd700, #ffed4e);
  box-shadow: 0 0 10px rgba(255, 215, 0, 0.5);
}

.character-detail-progress-text {
  color: var(--text-primary);
  font-size: 12px;
}

.character-detail-stats {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 8px;
  margin-bottom: 15px;
  font-size: 14px;
  padding-top: 15px;
  border-top: 1px solid var(--text-dim);
}

.character-detail-stat {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.character-detail-stat-label {
  color: var(--text-secondary);
}

.character-detail-stat-value {
  color: var(--text-white);
}

.character-detail-unspent {
  color: var(--terminal-cyan);
  font-size: 12px;
  margin-bottom: 6px;
  text-align: right;
}

.stat-allocate-btn {
  margin-left: 6px;
  width: 20px;
  height: 20px;
  border: 1px solid var(--terminal-green);
  background: rgba(0, 0, 0, 0.4);
  color: var(--terminal-green);
  cursor: pointer;
  font-size: 12px;
  line-height: 1;
  transition: all 0.15s ease;
}

.stat-allocate-btn:disabled {
  border-color: var(--text-dim);
  color: var(--text-dim);
  cursor: not-allowed;
  opacity: 0.7;
}

.stat-allocate-btn:not(:disabled):hover {
  background: var(--terminal-green);
  color: var(--terminal-bg);
}

.character-detail-combat-stats {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 8px;
  margin-bottom: 15px;
  font-size: 14px;
  padding-top: 15px;
  border-top: 1px solid var(--text-dim);
}

.character-detail-combat-stat {
  display: flex;
  justify-content: space-between;
  cursor: help;
  position: relative;
}

.character-detail-combat-stat:hover {
  color: var(--terminal-cyan);
}

.character-detail-combat-stat-label {
  color: var(--text-secondary);
}

.character-detail-combat-stat-value {
  color: var(--text-cyan);
}

.character-detail-summary {
  display: flex;
  justify-content: space-between;
  font-size: 12px;
  padding-top: 10px;
  border-top: 1px solid var(--text-dim);
}

.character-detail-summary-kills {
  color: var(--terminal-red);
}

.character-detail-summary-deaths {
  color: var(--terminal-gray);
  opacity: 0.7;
}

.character-detail-section-title {
  color: var(--terminal-cyan);
  font-size: 14px;
  font-weight: bold;
  margin-bottom: 10px;
  padding-bottom: 6px;
  border-bottom: 1px solid var(--text-dim);
}

.character-detail-buffs {
  margin-bottom: 15px;
  padding-top: 15px;
  border-top: 1px solid var(--text-dim);
}

.character-detail-buffs-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.character-detail-buff-item {
  padding: 8px;
  border: 1px solid;
  border-radius: 4px;
  font-size: 12px;
}

.character-detail-buff-item.buff-positive {
  background: rgba(0, 255, 0, 0.1);
  border-color: #00ff00;
}

.character-detail-buff-item.buff-negative {
  background: rgba(255, 0, 0, 0.1);
  border-color: #ff4444;
}

.buff-item-name {
  font-weight: bold;
  margin-bottom: 4px;
}

.buff-item-desc {
  color: var(--text-secondary);
  font-size: 11px;
  margin-bottom: 4px;
}

.buff-item-duration {
  color: var(--text-gray);
  font-size: 10px;
}

.character-detail-skills {
  margin-bottom: 15px;
  padding-top: 15px;
  border-top: 1px solid var(--text-dim);
}

.character-detail-loading,
.character-detail-no-skills {
  color: var(--text-gray);
  font-size: 12px;
  text-align: center;
  padding: 10px;
}

.character-detail-skills-list {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 6px;
  max-height: 250px;
  overflow-y: auto;
  overflow-x: hidden;
  padding-right: 2px;
  box-sizing: border-box;
  width: 100%;
}

.character-detail-skill-item {
  padding: 6px 8px;
  border: 1px solid var(--text-dim);
  border-radius: 3px;
  background: rgba(0, 0, 0, 0.3);
  font-size: 11px;
  cursor: help;
  transition: all 0.2s;
  min-width: 0;
  box-sizing: border-box;
  overflow: hidden;
}

.character-detail-skill-item:hover {
  border-color: var(--terminal-green);
  background: rgba(0, 255, 0, 0.05);
}

.skill-item-main {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 4px;
}

.skill-item-name {
  font-weight: bold;
  color: var(--terminal-cyan);
  font-size: 12px;
  flex: 1;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.skill-item-level {
  color: var(--terminal-gold);
  font-size: 10px;
  margin-left: 6px;
  flex-shrink: 0;
}

.skill-item-meta {
  display: flex;
  gap: 8px;
  font-size: 9px;
  color: var(--text-gray);
}

.skill-item-cost {
  color: var(--terminal-cyan);
}

.skill-item-cooldown {
  color: var(--terminal-yellow);
}

/* 弹窗滚动条样式 */
.character-detail-content::-webkit-scrollbar {
  width: 6px;
}

.character-detail-content::-webkit-scrollbar-track {
  background: rgba(0, 0, 0, 0.3);
}

.character-detail-content::-webkit-scrollbar-thumb {
  background: var(--terminal-gray);
  border-radius: 3px;
}

.character-detail-content::-webkit-scrollbar-thumb:hover {
  background: var(--terminal-green);
}

.character-name {
  font-family: var(--font-pixel);
  font-size: 18px;
  margin-bottom: 8px;
  /* 颜色通过内联样式动态设置 */
}

.character-level {
  color: var(--text-cyan);
  margin-bottom: 15px;
  font-size: 14px;
}

/* 进度条区域 */
.progress-section {
  margin-bottom: 15px;
}

.progress-item {
  margin-bottom: 12px;
}

.progress-label {
  color: var(--text-secondary);
  font-size: 12px;
  margin-bottom: 4px;
}

.progress-bar {
  width: 100%;
  height: 14px;
  background: var(--bg-color);
  border: 1px solid var(--text-dim);
  position: relative;
  overflow: hidden;
  margin-bottom: 4px;
}

.progress-fill {
  height: 100%;
  transition: width 0.3s ease;
}

.hp-bar .progress-fill {
  background: linear-gradient(90deg, #00ff00, #44ff44);
  box-shadow: 0 0 10px rgba(0, 255, 0, 0.5);
}

.mp-bar .progress-fill {
  background: linear-gradient(90deg, #ff4444, #ff6666);
  box-shadow: 0 0 10px rgba(255, 68, 68, 0.5);
}

.exp-bar .progress-fill {
  background: linear-gradient(90deg, #ffd700, #ffed4e);
  box-shadow: 0 0 10px rgba(255, 215, 0, 0.5);
}

.progress-text {
  color: var(--text-primary);
  font-size: 12px;
}

/* 属性网格 */
.character-stats {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 8px;
  margin-bottom: 15px;
  font-size: 14px;
}

.character-stat {
  display: flex;
  justify-content: space-between;
}

.character-stat-label {
  color: var(--text-secondary);
}

.character-stat-value {
  color: var(--text-white);
}

/* 战斗统计 */
.combat-stats {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 8px;
  margin-bottom: 15px;
  font-size: 14px;
  padding-top: 15px;
  border-top: 1px solid var(--text-dim);
}

.combat-stat {
  display: flex;
  justify-content: space-between;
}

.combat-stat-label {
  color: var(--text-secondary);
}

.combat-stat-value {
  color: var(--text-cyan);
}

/* 总结统计 */
.summary-stats {
  display: flex;
  justify-content: space-between;
  font-size: 12px;
  padding-top: 10px;
  border-top: 1px solid var(--text-dim);
}

.summary-kills {
  color: var(--terminal-red);
}

.summary-deaths {
  color: var(--terminal-gray);
  opacity: 0.7;
}

/* 游戏内容区 */
.game-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

/* 敌人信息面板（固定在顶部，横向排列） */
.enemies-panel {
  position: sticky;
  top: 0;
  z-index: 100;
  border-bottom: 2px solid var(--border-color);
  background: rgba(0, 0, 0, 0.95);
  padding: 8px 12px;
  display: flex;
  gap: 12px;
  flex-wrap: wrap;
  flex-shrink: 0;
  overflow-x: auto;
  overflow-y: hidden;
}

.terminal-content {
  flex: 1;
  min-height: 0;
  overflow-y: auto;
  padding: 16px;
  position: relative;
  z-index: 1;
}

.battle-log {
  position: relative;
  z-index: 1;
}

.no-character {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 100vh;
  padding: 40px;
}

.no-character-message {
  text-align: center;
  border: 2px solid var(--terminal-green);
  padding: 40px;
  background: rgba(0, 50, 0, 0.3);
  box-shadow: 0 0 20px rgba(0, 255, 0, 0.1);
}

.no-character-message h2 {
  color: var(--terminal-gold);
  margin-bottom: 20px;
  font-size: 18px;
}

.no-character-message p {
  color: var(--terminal-green);
  font-size: 14px;
  margin: 10px 0;
}

/* 敌人信息样式覆盖（横向排列） */
.enemy-info {
  position: relative;
  z-index: 11;
  border: 1px solid var(--text-dim);
  padding: 6px 10px;
  display: flex;
  flex-direction: column;
  gap: 4px;
  min-width: 160px;
  flex: 1;
  max-width: 280px;
  background: rgba(50, 0, 0, 0.5);
  transition: all 0.3s;
  border-radius: 4px;
}

/* 普通怪物样式 */
.enemy-info.enemy-normal {
  border: 1px solid var(--text-dim);
  background: rgba(50, 0, 0, 0.5);
  position: relative;
}

.enemy-info.enemy-normal::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  border: 1px solid rgba(100, 100, 100, 0.3);
  border-radius: 4px;
  pointer-events: none;
}

/* 精英怪物样式 - 蓝色边框，发光效果 */
.enemy-info.enemy-elite {
  border: 2px solid #4a90d9 !important;
  background: linear-gradient(135deg, rgba(20, 30, 60, 0.8), rgba(30, 50, 90, 0.6)) !important;
  box-shadow: 0 0 12px rgba(74, 144, 217, 0.6), 
              0 0 20px rgba(74, 144, 217, 0.3),
              inset 0 0 12px rgba(74, 144, 217, 0.15),
              0 0 0 1px rgba(74, 144, 217, 0.5) !important;
  position: relative;
  animation: elite-glow-border 2s ease-in-out infinite alternate;
}

.enemy-info.enemy-elite::before {
  content: '';
  position: absolute;
  top: -2px;
  left: -2px;
  right: -2px;
  bottom: -2px;
  background: linear-gradient(45deg, 
    transparent 30%, 
    rgba(74, 144, 217, 0.3) 50%, 
    transparent 70%);
  border-radius: 6px;
  pointer-events: none;
  animation: elite-shine 3s linear infinite;
}

.enemy-info.enemy-elite .enemy-name {
  /* 颜色由 getEnemyNameColor 函数动态设置，这里只设置阴影效果 */
  text-shadow: 0 0 8px currentColor,
               0 0 12px currentColor,
               0 0 4px currentColor !important;
  font-weight: 700;
}

.enemy-info.enemy-elite .enemy-rarity-icon {
  color: #6bb3ff !important;
  text-shadow: 0 0 10px rgba(74, 144, 217, 1),
               0 0 15px rgba(74, 144, 217, 0.8),
               0 0 20px rgba(74, 144, 217, 0.6) !important;
  animation: elite-glow 2s ease-in-out infinite alternate;
  filter: drop-shadow(0 0 4px rgba(74, 144, 217, 0.8));
}

.enemy-info.enemy-elite .enemy-bar {
  border-color: #4a90d9 !important;
  box-shadow: 0 0 4px rgba(74, 144, 217, 0.4);
}

.enemy-info.enemy-elite .enemy-bar-fill {
  background: linear-gradient(90deg, #4a90d9, #6bb3ff) !important;
  box-shadow: 0 0 6px rgba(74, 144, 217, 0.6);
}

/* Boss怪物样式 - 橙色/金色边框，强烈发光效果 */
.enemy-info.enemy-boss {
  border: 3px solid #ff6b35 !important;
  background: linear-gradient(135deg, rgba(60, 20, 20, 0.9), rgba(80, 30, 30, 0.7)) !important;
  box-shadow: 0 0 16px rgba(255, 107, 53, 0.8), 
              0 0 28px rgba(255, 107, 53, 0.5),
              0 0 40px rgba(255, 107, 53, 0.3),
              inset 0 0 16px rgba(255, 107, 53, 0.25),
              0 0 0 2px rgba(255, 215, 0, 0.6) !important;
  position: relative;
  animation: boss-pulse 2s ease-in-out infinite;
  transform: scale(1.02);
}

.enemy-info.enemy-boss::before {
  content: '';
  position: absolute;
  top: -3px;
  left: -3px;
  right: -3px;
  bottom: -3px;
  background: linear-gradient(45deg, 
    transparent 30%, 
    rgba(255, 215, 0, 0.4) 50%, 
    transparent 70%);
  border-radius: 8px;
  pointer-events: none;
  animation: boss-shine 2s linear infinite;
  z-index: -1;
}

.enemy-info.enemy-boss::after {
  content: '⚡';
  position: absolute;
  top: -8px;
  right: -8px;
  font-size: 20px;
  color: #ffd700;
  text-shadow: 0 0 10px rgba(255, 215, 0, 1),
               0 0 20px rgba(255, 215, 0, 0.8);
  animation: boss-sparkle 1.5s ease-in-out infinite;
  pointer-events: none;
}

.enemy-info.enemy-boss .enemy-name {
  /* 颜色由 getEnemyNameColor 函数动态设置，这里只设置阴影效果 */
  text-shadow: 0 0 10px currentColor,
               0 0 16px currentColor,
               0 0 24px currentColor,
               0 0 6px rgba(255, 215, 0, 0.8) !important;
  font-weight: 900;
  font-size: 15px;
  letter-spacing: 0.5px;
}

.enemy-info.enemy-boss .enemy-rarity-icon {
  color: #ffd700 !important;
  text-shadow: 0 0 12px rgba(255, 215, 0, 1),
               0 0 18px rgba(255, 215, 0, 1),
               0 0 24px rgba(255, 215, 0, 0.8),
               0 0 30px rgba(255, 215, 0, 0.6) !important;
  animation: boss-crown 1.5s ease-in-out infinite;
  filter: drop-shadow(0 0 6px rgba(255, 215, 0, 1));
  font-size: 18px;
}

.enemy-info.enemy-boss .enemy-bar {
  border: 2px solid #ff6b35 !important;
  box-shadow: 0 0 8px rgba(255, 107, 53, 0.6),
              inset 0 0 4px rgba(255, 107, 53, 0.3);
}

.enemy-info.enemy-boss .enemy-bar-fill {
  background: linear-gradient(90deg, #ff6b35, #ff8c5a, #ffd700) !important;
  box-shadow: 0 0 10px rgba(255, 107, 53, 0.8),
              inset 0 0 4px rgba(255, 215, 0, 0.4);
  animation: boss-bar-glow 1.5s ease-in-out infinite alternate;
}

.enemy-info.enemy-dead {
  opacity: 0.5;
  border-color: var(--text-gray);
  box-shadow: none;
  animation: none;
}

.enemy-info .enemy-name {
  font-size: 14px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  /* 颜色由 getEnemyNameColor 函数动态设置 */
  color: #ff7777; /* 默认红色（物理攻击），会被内联样式覆盖 */
  font-weight: bold;
  display: flex;
  align-items: center;
  gap: 4px;
}

.enemy-rarity-icon {
  font-size: 16px;
  display: inline-block;
}

.enemy-info .enemy-hp {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 10px;
}

.enemy-info .enemy-hp-label {
  color: var(--text-gray);
  font-size: 10px;
  min-width: 24px;
}

/* 地图选择器样式 */
.zone-selector-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.8);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.zone-selector {
  background: var(--bg-dark);
  border: 2px solid var(--text-dim);
  border-radius: 4px;
  width: 90%;
  max-width: 800px;
  max-height: 80vh;
  display: flex;
  flex-direction: column;
  color: var(--text-bright);
}

.zone-selector-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 15px 20px;
  border-bottom: 1px solid var(--text-dim);
}

.zone-selector-header h2 {
  margin: 0;
  color: var(--text-bright);
  font-size: 18px;
}

.close-btn {
  background: none;
  border: none;
  color: var(--text-bright);
  font-size: 24px;
  cursor: pointer;
  padding: 0;
  width: 30px;
  height: 30px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.close-btn:hover {
  color: var(--text-red);
}

.zone-error {
  padding: 15px 20px;
  color: var(--text-red);
  background: rgba(255, 0, 0, 0.1);
  border-bottom: 1px solid var(--text-dim);
}

.zone-loading {
  padding: 40px;
  text-align: center;
  color: var(--text-dim);
}

.zone-list {
  padding: 15px;
  overflow-y: auto;
  flex: 1;
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 15px;
}

.zone-item {
  background: rgba(50, 0, 0, 0.3);
  border: 1px solid var(--text-dim);
  border-radius: 4px;
  padding: 15px;
  cursor: pointer;
  transition: all 0.2s;
}

.zone-item:hover {
  border-color: var(--text-bright);
  background: rgba(50, 0, 0, 0.5);
}

.zone-item.zone-current {
  border-color: var(--terminal-green);
  background: rgba(0, 50, 0, 0.3);
}

.zone-item.zone-locked {
  opacity: 0.5;
  cursor: not-allowed;
}

.zone-name {
  font-size: 16px;
  font-weight: bold;
  color: var(--text-bright);
  margin-bottom: 8px;
}

.zone-info {
  display: flex;
  gap: 10px;
  flex-wrap: wrap;
  margin-bottom: 8px;
  font-size: 12px;
}

.zone-level {
  color: var(--text-dim);
}

.zone-faction {
  padding: 2px 6px;
  border-radius: 2px;
  font-size: 11px;
}

.faction-alliance {
  background: rgba(74, 144, 217, 0.3);
  color: #4a90d9;
}

.faction-horde {
  background: rgba(196, 30, 58, 0.3);
  color: #c41e3a;
}

.faction-neutral {
  background: rgba(255, 215, 0, 0.3);
  color: #ffd700;
}

.zone-multiplier {
  color: var(--terminal-green);
}

.zone-description {
  font-size: 12px;
  color: var(--text-dim);
  line-height: 1.4;
  margin-top: 8px;
}

.zone-locked-hint {
  margin-top: 8px;
  font-size: 11px;
  color: var(--text-red);
}

.zone-exploration-progress {
  margin-top: 10px;
  padding: 8px;
  background: rgba(0, 0, 0, 0.3);
  border-radius: 4px;
}

.exploration-label {
  font-size: 11px;
  color: var(--text-dim);
  margin-bottom: 6px;
}

.exploration-label strong {
  color: var(--text-bright);
}

.exploration-bar {
  position: relative;
  height: 20px;
  background: rgba(0, 0, 0, 0.5);
  border: 1px solid var(--text-dim);
  border-radius: 2px;
  overflow: hidden;
}

.exploration-fill {
  height: 100%;
  background: linear-gradient(90deg, var(--terminal-green), #4ade80);
  transition: width 0.3s;
}

.exploration-text {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  font-size: 10px;
  color: var(--text-bright);
  font-weight: bold;
  text-shadow: 0 0 2px rgba(0, 0, 0, 0.8);
}

.zone-current-exploration {
  margin-top: 8px;
  font-size: 11px;
  color: var(--terminal-green);
  font-weight: bold;
}

.enemy-info .enemy-bar {
  flex: 1;
  min-width: 60px;
  height: 10px;
  background: var(--bg-color);
  border: 1px solid var(--text-red);
}

.enemy-info .enemy-hp-value {
  color: #ff4444;
  font-size: 10px;
  white-space: nowrap;
  min-width: 50px;
  text-align: right;
}

/* 稀有度动画效果 */
@keyframes elite-glow {
  0% {
    text-shadow: 0 0 10px rgba(74, 144, 217, 1),
                 0 0 15px rgba(74, 144, 217, 0.8),
                 0 0 20px rgba(74, 144, 217, 0.6);
  }
  100% {
    text-shadow: 0 0 15px rgba(74, 144, 217, 1),
                 0 0 20px rgba(74, 144, 217, 1),
                 0 0 30px rgba(74, 144, 217, 0.8);
  }
}

@keyframes elite-glow-border {
  0% {
    box-shadow: 0 0 12px rgba(74, 144, 217, 0.6), 
                0 0 20px rgba(74, 144, 217, 0.3),
                inset 0 0 12px rgba(74, 144, 217, 0.15),
                0 0 0 1px rgba(74, 144, 217, 0.5);
  }
  100% {
    box-shadow: 0 0 16px rgba(74, 144, 217, 0.8), 
                0 0 28px rgba(74, 144, 217, 0.5),
                inset 0 0 16px rgba(74, 144, 217, 0.25),
                0 0 0 2px rgba(74, 144, 217, 0.7);
  }
}

@keyframes elite-shine {
  0% {
    transform: translateX(-100%) translateY(-100%);
    opacity: 0;
  }
  50% {
    opacity: 1;
  }
  100% {
    transform: translateX(100%) translateY(100%);
    opacity: 0;
  }
}

@keyframes boss-pulse {
  0%, 100% {
    box-shadow: 0 0 16px rgba(255, 107, 53, 0.8), 
                0 0 28px rgba(255, 107, 53, 0.5),
                0 0 40px rgba(255, 107, 53, 0.3),
                inset 0 0 16px rgba(255, 107, 53, 0.25),
                0 0 0 2px rgba(255, 215, 0, 0.6);
    transform: scale(1.02);
  }
  50% {
    box-shadow: 0 0 20px rgba(255, 107, 53, 1), 
                0 0 36px rgba(255, 107, 53, 0.7),
                0 0 50px rgba(255, 107, 53, 0.5),
                inset 0 0 20px rgba(255, 107, 53, 0.35),
                0 0 0 3px rgba(255, 215, 0, 0.8);
    transform: scale(1.03);
  }
}

@keyframes boss-crown {
  0%, 100% {
    transform: scale(1) rotate(0deg);
    text-shadow: 0 0 12px rgba(255, 215, 0, 1),
                 0 0 18px rgba(255, 215, 0, 1),
                 0 0 24px rgba(255, 215, 0, 0.8);
  }
  25% {
    transform: scale(1.1) rotate(-3deg);
    text-shadow: 0 0 15px rgba(255, 215, 0, 1),
                 0 0 22px rgba(255, 215, 0, 1),
                 0 0 30px rgba(255, 215, 0, 1);
  }
  50% {
    transform: scale(1.15) rotate(0deg);
    text-shadow: 0 0 18px rgba(255, 215, 0, 1),
                 0 0 26px rgba(255, 215, 0, 1),
                 0 0 36px rgba(255, 215, 0, 1);
  }
  75% {
    transform: scale(1.1) rotate(3deg);
    text-shadow: 0 0 15px rgba(255, 215, 0, 1),
                 0 0 22px rgba(255, 215, 0, 1),
                 0 0 30px rgba(255, 215, 0, 1);
  }
}

@keyframes boss-shine {
  0% {
    transform: translateX(-100%) translateY(-100%) rotate(45deg);
    opacity: 0;
  }
  50% {
    opacity: 1;
  }
  100% {
    transform: translateX(100%) translateY(100%) rotate(45deg);
    opacity: 0;
  }
}

@keyframes boss-sparkle {
  0%, 100% {
    transform: scale(1) rotate(0deg);
    opacity: 0.8;
  }
  50% {
    transform: scale(1.2) rotate(180deg);
    opacity: 1;
  }
}

@keyframes boss-bar-glow {
  0% {
    box-shadow: 0 0 10px rgba(255, 107, 53, 0.8),
                inset 0 0 4px rgba(255, 215, 0, 0.4);
  }
  100% {
    box-shadow: 0 0 15px rgba(255, 107, 53, 1),
                inset 0 0 6px rgba(255, 215, 0, 0.6);
  }
}
</style>
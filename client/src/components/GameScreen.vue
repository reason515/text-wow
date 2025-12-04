<script setup lang="ts">
import { ref, computed, watch, nextTick, onMounted } from 'vue'
import { useGameStore } from '../stores/game'
import { useCharacterStore } from '../stores/character'
import { useAuthStore } from '../stores/auth'
import { getClassColor, getResourceColor } from '../types/game'
import ChatPanel from './ChatPanel.vue'

const emit = defineEmits<{
  logout: []
  'create-character': []
}>()

const game = useGameStore()
const charStore = useCharacterStore()
const authStore = useAuthStore()
const logContainer = ref<HTMLElement | null>(null)

// 角色详情弹窗
const showCharacterDetail = ref(false)
const selectedCharacter = ref<any>(null)
const characterSkills = ref<any[]>([])
const loadingSkills = ref(false)

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
    } else {
      characterSkills.value = []
    }
  } catch (e) {
    console.error('Failed to fetch character skills:', e)
    characterSkills.value = []
  } finally {
    loadingSkills.value = false
  }
}

// 关闭角色详情
function closeDetail() {
  showCharacterDetail.value = false
  selectedCharacter.value = null
  characterSkills.value = []
  hideSkillTooltip() // 关闭时也隐藏tooltip
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

// 处理技能tooltip显示（使用fixed定位避免被overflow裁剪）
let skillTooltipEl: HTMLElement | null = null

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
  
  // 获取敌方角色名（从当前敌人或日志中的target/actor字段）
  let enemyName = ''
  // 优先使用target字段（如果actor是我方，target就是敌方）
  if (log.target && log.target !== playerName && !playerNameVariants.includes(log.target)) {
    enemyName = log.target
  } 
  // 如果actor不是我方角色，则actor是敌方
  else if (log.actor && log.actor !== playerName && !playerNameVariants.includes(log.actor) && log.actor !== 'system') {
    enemyName = log.actor
  } 
  // 最后尝试从当前敌人列表获取
  else if (game.currentEnemies && game.currentEnemies.length > 0) {
    const currentEnemy = game.currentEnemies[0] as any
    enemyName = currentEnemy?.name || ''
  }
  
  // 获取技能名（从日志的action字段或消息中的方括号内容）
  let skillName = ''
  if (log.action && log.action !== '攻击' && log.action !== 'encounter' && log.action !== 'victory' && log.action !== 'defeat' && log.action !== 'loot' && log.action !== 'levelup') {
    skillName = log.action
  }
  
  // 解析消息并添加颜色标记（传入资源颜色用于技能颜色）
  return formatMessageWithColors(message, playerName, playerNameVariants, enemyName, skillName, playerColor, resourceColor)
}

// 格式化消息，为角色名和技能名添加颜色
function formatMessageWithColors(
  message: string,
  playerName: string,
  playerNameVariants: string[],
  enemyName: string,
  skillName: string,
  playerColor: string = '#ffff55', // 默认金色，如果未传入则使用默认值
  resourceColor: string = '#ffffff' // 资源颜色，用于技能颜色
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
  
  // 定义颜色（使用传入的职业颜色，敌方使用固定颜色）
  const enemyColor = '#ff7777' // var(--text-red)
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
  if (enemyName) {
    const regex = new RegExp(escapeRegex(enemyName), 'g')
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
                    `<span style="color: ${enemyColor}">${matchText}</span>` + 
                    formatted.substring(index + matchText.length)
      }
    }
  }
  
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
                    {{ char.resource || char.mp || 0 }}/{{ char.maxResource || char.max_resource || char.max_mp || 100 }}
                  </div>
                </div>
                <!-- Buff/Debuff显示 -->
                <div v-if="char.buffs && char.buffs.length > 0" class="team-character-buffs">
                  <div
                    v-for="buff in char.buffs"
                    :key="buff.effectId"
                    class="buff-icon"
                    :class="{ 'buff-positive': buff.isBuff, 'buff-negative': !buff.isBuff }"
                    :data-tooltip="buff.name + '\n' + (buff.description || '') + '\n剩余 ' + buff.duration + ' 回合'"
                  >
                    {{ buff.name.charAt(0) }}
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
              :class="{ 'enemy-dead': (enemy as any)?.hp <= 0 }"
            >
              <span class="enemy-name">
                ⚔ {{ (enemy as any)?.name || '未知敌人' }} (Lv.{{ (enemy as any)?.level || 1 }})
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
            <button class="cmd-btn" @click="game.battleTick" :disabled="game.isRunning">
              [→] 单步战斗
            </button>
            <button class="cmd-btn" disabled>
              [S] 策略
            </button>
            <button class="cmd-btn" disabled>
              [E] 装备
            </button>
            <button class="cmd-btn" disabled>
              [M] 地图
            </button>
          </div>
        </div>
      </div>

      <!-- 底部聊天面板 -->
      <ChatPanel />
    </template>
    
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
              {{ selectedCharacter.resource || selectedCharacter.mp || 0 }}/{{ selectedCharacter.maxResource || selectedCharacter.max_resource || selectedCharacter.max_mp || 100 }}
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
        <div class="character-detail-stats">
          <div class="character-detail-stat">
            <span class="character-detail-stat-label">力量</span>
            <span class="character-detail-stat-value">{{ selectedCharacter.strength || 0 }}</span>
          </div>
          <div class="character-detail-stat">
            <span class="character-detail-stat-label">敏捷</span>
            <span class="character-detail-stat-value">{{ selectedCharacter.agility || 0 }}</span>
          </div>
          <div class="character-detail-stat">
            <span class="character-detail-stat-label">智力</span>
            <span class="character-detail-stat-value">{{ selectedCharacter.intellect || 0 }}</span>
          </div>
          <div class="character-detail-stat">
            <span class="character-detail-stat-label">耐力</span>
            <span class="character-detail-stat-value">{{ selectedCharacter.stamina || 0 }}</span>
          </div>
          <div class="character-detail-stat">
            <span class="character-detail-stat-label">精神</span>
            <span class="character-detail-stat-value">{{ selectedCharacter.spirit || 0 }}</span>
          </div>
        </div>

        <!-- 战斗统计 -->
        <div class="character-detail-combat-stats">
          <div class="character-detail-combat-stat" :data-tooltip="`物理攻击力 = 力量 / 2\n当前: ${selectedCharacter.strength || 0} / 2 = ${selectedCharacter.physicalAttack || 0}`">
            <span class="character-detail-combat-stat-label">物理攻击</span>
            <span class="character-detail-combat-stat-value">{{ selectedCharacter.physicalAttack || 0 }}</span>
          </div>
          <div class="character-detail-combat-stat" :data-tooltip="`魔法攻击力 = 智力 / 2\n当前: ${selectedCharacter.intellect || 0} / 2 = ${selectedCharacter.magicAttack || 0}`">
            <span class="character-detail-combat-stat-label">魔法攻击</span>
            <span class="character-detail-combat-stat-value">{{ selectedCharacter.magicAttack || 0 }}</span>
          </div>
          <div class="character-detail-combat-stat" :data-tooltip="`物理防御力 = 耐力 / 3\n当前: ${selectedCharacter.stamina || 0} / 3 = ${selectedCharacter.physicalDefense || 0}`">
            <span class="character-detail-combat-stat-label">物理防御</span>
            <span class="character-detail-combat-stat-value">{{ selectedCharacter.physicalDefense || 0 }}</span>
          </div>
          <div class="character-detail-combat-stat" :data-tooltip="`魔法防御力 = (智力 + 精神) / 4\n当前: (${selectedCharacter.intellect || 0} + ${selectedCharacter.spirit || 0}) / 4 = ${selectedCharacter.magicDefense || 0}`">
            <span class="character-detail-combat-stat-label">魔法防御</span>
            <span class="character-detail-combat-stat-value">{{ selectedCharacter.magicDefense || 0 }}</span>
          </div>
          <div class="character-detail-combat-stat" :data-tooltip="`暴击率 = 基础值 + 被动技能加成 + Buff加成\n基础值: 5%\n当前: ${((selectedCharacter.critRate || 0.05) * 100).toFixed(1)}%`">
            <span class="character-detail-combat-stat-label">暴击率</span>
            <span class="character-detail-combat-stat-value">{{ ((selectedCharacter.critRate || 0.05) * 100).toFixed(1) }}%</span>
          </div>
          <div class="character-detail-combat-stat" :data-tooltip="`暴击伤害 = 基础值 × 暴击伤害倍率\n基础值: 150%\n当前: ${((selectedCharacter.critDamage || 1.5) * 100).toFixed(0)}%`">
            <span class="character-detail-combat-stat-label">暴击伤害</span>
            <span class="character-detail-combat-stat-value">{{ ((selectedCharacter.critDamage || 1.5) * 100).toFixed(0) }}%</span>
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
}

.character-detail-stat-label {
  color: var(--text-secondary);
}

.character-detail-stat-value {
  color: var(--text-white);
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
}

.character-detail-skill-item {
  padding: 6px 8px;
  border: 1px solid var(--text-dim);
  border-radius: 3px;
  background: rgba(0, 0, 0, 0.3);
  font-size: 11px;
  cursor: help;
  transition: all 0.2s;
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
  position: relative;
  z-index: 10;
  border-bottom: 2px solid var(--border-color);
  background: rgba(0, 0, 0, 0.8);
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
  transition: opacity 0.3s;
}

.enemy-info.enemy-dead {
  opacity: 0.5;
  border-color: var(--text-gray);
}

.enemy-info .enemy-name {
  font-size: 14px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  color: var(--text-red);
  font-weight: bold;
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
</style>

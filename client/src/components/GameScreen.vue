<script setup lang="ts">
import { ref, computed, watch, nextTick, onMounted } from 'vue'
import { useGameStore } from '../stores/game'
import { useCharacterStore } from '../stores/character'
import { useAuthStore } from '../stores/auth'
import { getClassColor } from '../types/game'
import ChatPanel from './ChatPanel.vue'

const emit = defineEmits<{
  logout: []
  'create-character': []
}>()

const game = useGameStore()
const charStore = useCharacterStore()
const authStore = useAuthStore()
const logContainer = ref<HTMLElement | null>(null)

// 初始化：从 characterStore 获取角色数据并同步到 gameStore
onMounted(async () => {
  console.log('GameScreen mounted')
  console.log('charStore.characters:', charStore.characters)
  console.log('charStore.activeCharacters:', charStore.activeCharacters)
  
  // 如果没有角色，先尝试获取
  if (charStore.characters.length === 0) {
    await charStore.fetchCharacters()
  }
  
  // 获取当前激活的角色
  const activeChar = charStore.activeCharacters[0] || charStore.characters[0]
  
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
  
  // 如果战斗状态中有角色数据，使用它（可能包含最新的死亡/复活状态）
  // Team 是一个数组，不是包含 characters 的对象
  if (game.battleStatus?.team && Array.isArray(game.battleStatus.team) && game.battleStatus.team.length > 0) {
    game.character = game.battleStatus.team[0]
    console.log('Character updated from battle status:', game.character)
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

// 计算HP/MP/EXP百分比
const hpPercent = computed(() => {
  if (!game.character) return 0
  const char = game.character as any
  const maxHp = char.maxHp || char.max_hp || 100
  const hp = char.hp || 0
  return (hp / maxHp) * 100
})

const mpPercent = computed(() => {
  if (!game.character) return 0
  const char = game.character as any
  // 支持多种字段名：resource/maxResource 或 mp/max_mp
  const maxResource = char.maxResource || char.max_resource || char.max_mp || 100
  const resource = char.resource || char.mp || 0
  return (resource / maxResource) * 100
})

const expPercent = computed(() => {
  if (!game.character) return 0
  const char = game.character as any
  // 支持多种字段名：expToNext 或 exp_to_next
  const expToNext = char.expToNext || char.exp_to_next || 100
  const exp = char.exp || 0
  return (exp / expToNext) * 100
})

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
const resourceTypeName = computed(() => {
  if (!game.character) return 'MP'
  const char = game.character as any
  const type = char.resourceType || 'mana'
  const names: Record<string, string> = {
    mana: '法力',
    rage: '怒气',
    energy: '能量'
  }
  return names[type] || 'MP'
})

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
  
  // 解析消息并添加颜色标记
  return formatMessageWithColors(message, playerName, playerNameVariants, enemyName, skillName, playerColor)
}

// 格式化消息，为角色名和技能名添加颜色
function formatMessageWithColors(
  message: string,
  playerName: string,
  playerNameVariants: string[],
  enemyName: string,
  skillName: string,
  playerColor: string = '#ffff55' // 默认金色，如果未传入则使用默认值
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
  
  // 定义颜色（使用传入的职业颜色，敌方和技能使用固定颜色）
  const enemyColor = '#ff7777' // var(--text-red)
  const skillColor = '#77ffff' // var(--text-cyan)
  
  // 先转义整个消息
  let formatted = escapeHtml(message)
  
  // 标记技能名（方括号内的内容）- 优先处理，避免与其他标记冲突
  formatted = formatted.replace(/\[([^\]]+)\]/g, (match, skill) => {
    return `<span style="color: ${skillColor}">[${escapeHtml(skill)}]</span>`
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
        <span>战斗: {{ (game.battleStatus as any)?.battleCount || (game.battleStatus as any)?.battle_count || 0 }}</span>
        <span>|</span>
        <span>击杀: {{ (game.battleStatus as any)?.totalKills || (game.battleStatus as any)?.session_kills || 0 }}</span>
        <span>|</span>
        <span>+{{ (game.battleStatus as any)?.totalExp || (game.battleStatus as any)?.session_exp || 0 }} EXP</span>
        <span>|</span>
        <span>+{{ (game.battleStatus as any)?.totalGold || (game.battleStatus as any)?.session_gold || 0 }} G</span>
        <span class="battle-status" :class="{ active: game.isRunning }">
          {{ game.isRunning ? '× 战斗中' : '○ 待机' }}
        </span>
      </div>

      <!-- 主内容区 -->
      <div class="game-main">
        <!-- 左侧角色信息面板 -->
        <div class="game-sidebar">
          <div class="character-card">
            <div class="character-name">
              {{ game.character?.name }}
            </div>
            <div class="character-level">
              Lv.{{ game.character?.level }} {{ getClassName((game.character as any)?.classId || (game.character as any)?.class || '') }}
            </div>
            
            <!-- 进度条 -->
            <div class="progress-section">
              <div class="progress-item">
                <div class="progress-label">生命值</div>
                <div class="progress-bar hp-bar">
                  <div class="progress-fill" :style="{ width: hpPercent + '%' }"></div>
                </div>
                <div class="progress-text">
                  {{ (game.character as any)?.hp || 0 }}/{{ (game.character as any)?.maxHp || (game.character as any)?.max_hp || 100 }}
                </div>
              </div>
              
              <div class="progress-item">
                <div class="progress-label">{{ resourceTypeName }}</div>
                <div class="progress-bar mp-bar">
                  <div class="progress-fill" :style="{ width: mpPercent + '%' }"></div>
                </div>
                <div class="progress-text">
                  {{ (game.character as any)?.resource || (game.character as any)?.mp || 0 }}/{{ (game.character as any)?.maxResource || (game.character as any)?.max_resource || (game.character as any)?.max_mp || 100 }}
                </div>
              </div>
              
              <div class="progress-item">
                <div class="progress-label">经验值</div>
                <div class="progress-bar exp-bar">
                  <div class="progress-fill" :style="{ width: expPercent + '%' }"></div>
                </div>
                <div class="progress-text">
                  {{ (game.character as any)?.exp || 0 }}/{{ (game.character as any)?.expToNext || (game.character as any)?.exp_to_next || 100 }}
                </div>
              </div>
            </div>

            <!-- 属性 -->
            <div class="character-stats">
              <div class="character-stat">
                <span class="character-stat-label">力量</span>
                <span class="character-stat-value">{{ (game.character as any)?.strength || 0 }}</span>
              </div>
              <div class="character-stat">
                <span class="character-stat-label">敏捷</span>
                <span class="character-stat-value">{{ (game.character as any)?.agility || 0 }}</span>
              </div>
              <div class="character-stat">
                <span class="character-stat-label">智力</span>
                <span class="character-stat-value">{{ (game.character as any)?.intellect || 0 }}</span>
              </div>
              <div class="character-stat">
                <span class="character-stat-label">耐力</span>
                <span class="character-stat-value">{{ (game.character as any)?.stamina || 0 }}</span>
              </div>
              <div class="character-stat">
                <span class="character-stat-label">精神</span>
                <span class="character-stat-value">{{ (game.character as any)?.spirit || 0 }}</span>
              </div>
            </div>

            <!-- 战斗统计 -->
            <div class="combat-stats">
              <div class="combat-stat">
                <span class="combat-stat-label">攻击力</span>
                <span class="combat-stat-value">{{ (game.character as any)?.attack || 0 }}</span>
              </div>
              <div class="combat-stat">
                <span class="combat-stat-label">防御力</span>
                <span class="combat-stat-value">{{ (game.character as any)?.defense || 0 }}</span>
              </div>
              <div class="combat-stat">
                <span class="combat-stat-label">暴击率</span>
                <span class="combat-stat-value">{{ ((game.character as any)?.critRate || 0).toFixed(1) }}%</span>
              </div>
              <div class="combat-stat">
                <span class="combat-stat-label">暴击伤害</span>
                <span class="combat-stat-value">{{ ((game.character as any)?.critDamage || 150).toFixed(0) }}%</span>
              </div>
            </div>

            <!-- 总结统计 -->
            <div class="summary-stats">
              <div>击杀: {{ (game.character as any)?.totalKills || 0 }}</div>
              <div>死亡: {{ (game.character as any)?.totalDeaths || 0 }}</div>
            </div>
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
}

.character-card {
  border: 2px solid var(--border-color);
  padding: 15px;
  background: rgba(51, 255, 51, 0.03);
}

.character-name {
  font-family: var(--font-pixel);
  font-size: 14px;
  color: var(--text-gold);
  text-shadow: 0 0 10px var(--text-gold);
  margin-bottom: 8px;
}

.character-level {
  color: var(--text-cyan);
  margin-bottom: 15px;
  font-size: 12px;
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
  font-size: 11px;
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
  background: linear-gradient(90deg, #00aaff, #00ddff);
  box-shadow: 0 0 10px rgba(0, 170, 255, 0.5);
}

.progress-text {
  color: var(--text-primary);
  font-size: 11px;
}

/* 属性网格 */
.character-stats {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 8px;
  margin-bottom: 15px;
  font-size: 12px;
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
  font-size: 12px;
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
  font-size: 11px;
  color: var(--text-gray);
  padding-top: 10px;
  border-top: 1px solid var(--text-dim);
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
  border-bottom: 2px solid var(--border-color);
  background: rgba(0, 0, 0, 0.5);
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
  font-size: 24px;
}

.no-character-message p {
  color: var(--terminal-green);
  font-size: 16px;
  margin: 10px 0;
}

/* 敌人信息样式覆盖（横向排列） */
.enemy-info {
  border: 1px solid var(--text-dim);
  padding: 6px 10px;
  display: flex;
  flex-direction: column;
  gap: 4px;
  min-width: 160px;
  flex: 1;
  max-width: 280px;
  background: rgba(50, 0, 0, 0.3);
  transition: opacity 0.3s;
}

.enemy-info.enemy-dead {
  opacity: 0.5;
  border-color: var(--text-gray);
}

.enemy-info .enemy-name {
  font-size: 12px;
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

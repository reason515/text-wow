<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted, watch, nextTick } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { useCharacterStore } from '@/stores/character'
import { useGameStore } from '@/stores/game'
import ChatPanel from '@/components/ChatPanel.vue'
import { CLASS_COLORS, CLASS_NAMES } from '@/types/game'

const emit = defineEmits<{
  logout: []
  'create-character': []
}>()

const authStore = useAuthStore()
const charStore = useCharacterStore()
const gameStore = useGameStore()

// 日志容器引用
const logContainer = ref<HTMLElement | null>(null)

// 当前选中的角色
const selectedCharacter = computed(() => {
  return charStore.characters[0] || null
})

// 资源类型名称
const resourceTypeName = computed(() => {
  if (!selectedCharacter.value) return '能量'
  const types: Record<string, string> = {
    mana: '法力',
    rage: '怒气',
    energy: '能量'
  }
  return types[selectedCharacter.value.resourceType] || '能量'
})

// 资源条颜色
const resourceBarColor = computed(() => {
  if (!selectedCharacter.value) return '#4a90d9'
  const colors: Record<string, string> = {
    mana: '#4a90d9',
    rage: '#c41e3a',
    energy: '#f0b90b'
  }
  return colors[selectedCharacter.value.resourceType] || '#4a90d9'
})

// 百分比计算
const hpPercent = computed(() => {
  if (!selectedCharacter.value) return 0
  return (selectedCharacter.value.hp / selectedCharacter.value.maxHp) * 100
})

const resourcePercent = computed(() => {
  if (!selectedCharacter.value) return 0
  return (selectedCharacter.value.resource / selectedCharacter.value.maxResource) * 100
})

const expPercent = computed(() => {
  if (!selectedCharacter.value) return 0
  return (selectedCharacter.value.exp / selectedCharacter.value.expToNext) * 100
})

// 阵营颜色
const factionColor = computed(() => {
  if (!selectedCharacter.value) return '#888'
  return selectedCharacter.value.faction === 'alliance' ? '#4a90d9' : '#c41e3a'
})

// 职业颜色
const classColor = computed(() => {
  if (!selectedCharacter.value) return '#33ff33'
  return CLASS_COLORS[selectedCharacter.value.classId] || '#33ff33'
})

// 职业名称
const className = computed(() => {
  if (!selectedCharacter.value) return ''
  return CLASS_NAMES[selectedCharacter.value.classId] || selectedCharacter.value.classId
})

// 敌人HP百分比
const enemyHpPercent = computed(() => {
  if (!gameStore.currentEnemy) return 0
  return (gameStore.currentEnemy.hp / gameStore.currentEnemy.maxHp) * 100
})

// 获取日志颜色
function getLogColor(logType: string): string {
  const colors: Record<string, string> = {
    system: '#33ff33',
    combat: '#ffaa00',
    victory: '#33ff33',
    loot: '#ffd700',
    levelup: '#ffd700',
    death: '#ff0000',
    encounter: '#ffff00',
    zone: '#00ffff'
  }
  return colors[logType] || '#33ff33'
}

// 格式化时间
function formatTime(date: string | Date): string {
  const d = typeof date === 'string' ? new Date(date) : date
  return d.toLocaleTimeString('zh-CN', { hour12: false, hour: '2-digit', minute: '2-digit', second: '2-digit' })
}

// 自动滚动到底部
function scrollToBottom() {
  nextTick(() => {
    if (logContainer.value) {
      logContainer.value.scrollTop = logContainer.value.scrollHeight
    }
  })
}

// 监听日志变化，自动滚动
watch(() => gameStore.battleLogs.length, () => {
  scrollToBottom()
})

// 加载数据
onMounted(async () => {
  await charStore.fetchCharacters()
  await gameStore.init()
  scrollToBottom()
})

onUnmounted(() => {
  gameStore.cleanup()
})

// 处理登出
function handleLogout() {
  gameStore.cleanup()
  emit('logout')
}

// 创建新角色
function createNewCharacter() {
  emit('create-character')
}

// 切换战斗
async function handleToggleBattle() {
  await gameStore.toggleBattle()
  // 刷新角色数据
  if (!gameStore.isRunning) {
    await charStore.fetchCharacters()
  }
}
</script>

<template>
  <div class="game-screen">
    <!-- 顶部导航 -->
    <div class="top-bar">
      <div class="user-info">
        <span class="username">{{ authStore.username }}</span>
        <span class="separator">|</span>
        <span class="gold">💰 {{ authStore.user?.gold || 0 }} G</span>
      </div>
      <div class="battle-stats" v-if="gameStore.battleStatus.battleCount > 0">
        <span>战斗: {{ gameStore.battleStatus.battleCount }}</span>
        <span class="separator">|</span>
        <span>击杀: {{ gameStore.battleStatus.totalKills }}</span>
        <span class="separator">|</span>
        <span class="exp">+{{ gameStore.battleStatus.totalExp }} EXP</span>
        <span class="separator">|</span>
        <span class="gold-gain">+{{ gameStore.battleStatus.totalGold }} G</span>
      </div>
      <div class="actions">
        <button class="text-btn" @click="createNewCharacter">新建角色</button>
        <button class="text-btn logout" @click="handleLogout">登出</button>
      </div>
    </div>

    <!-- 主内容区 -->
    <div class="main-content">
      <!-- 左侧：角色信息 -->
      <div class="character-panel" v-if="selectedCharacter">
        <div class="panel-header">
          <span class="faction-badge" :style="{ backgroundColor: factionColor }">
            {{ selectedCharacter.faction === 'alliance' ? '联盟' : '部落' }}
          </span>
          <h2 class="character-name">{{ selectedCharacter.name }}</h2>
        </div>
        
        <div class="character-class" :style="{ color: classColor, textShadow: `0 0 8px ${classColor}` }">
          Lv.{{ selectedCharacter.level }} {{ className }}
        </div>

        <!-- 状态条 -->
        <div class="stat-bars">
          <div class="stat-bar">
            <div class="bar-header">
              <span class="bar-label">生命值</span>
              <span class="bar-value">{{ selectedCharacter.hp }}/{{ selectedCharacter.maxHp }}</span>
            </div>
            <div class="bar-track hp">
              <div class="bar-fill" :style="{ width: hpPercent + '%' }"></div>
            </div>
          </div>

          <div class="stat-bar">
            <div class="bar-header">
              <span class="bar-label">{{ resourceTypeName }}</span>
              <span class="bar-value">{{ selectedCharacter.resource }}/{{ selectedCharacter.maxResource }}</span>
            </div>
            <div class="bar-track" :style="{ '--bar-color': resourceBarColor }">
              <div class="bar-fill" :style="{ width: resourcePercent + '%', backgroundColor: resourceBarColor }"></div>
            </div>
          </div>

          <div class="stat-bar">
            <div class="bar-header">
              <span class="bar-label">经验值</span>
              <span class="bar-value">{{ selectedCharacter.exp }}/{{ selectedCharacter.expToNext }}</span>
            </div>
            <div class="bar-track exp">
              <div class="bar-fill" :style="{ width: expPercent + '%' }"></div>
            </div>
          </div>
        </div>

        <!-- 属性 -->
        <div class="attributes">
          <div class="attr">
            <span class="attr-label">力量</span>
            <span class="attr-value str">{{ selectedCharacter.strength }}</span>
          </div>
          <div class="attr">
            <span class="attr-label">敏捷</span>
            <span class="attr-value agi">{{ selectedCharacter.agility }}</span>
          </div>
          <div class="attr">
            <span class="attr-label">智力</span>
            <span class="attr-value int">{{ selectedCharacter.intellect }}</span>
          </div>
          <div class="attr">
            <span class="attr-label">耐力</span>
            <span class="attr-value sta">{{ selectedCharacter.stamina }}</span>
          </div>
          <div class="attr">
            <span class="attr-label">精神</span>
            <span class="attr-value spi">{{ selectedCharacter.spirit }}</span>
          </div>
        </div>

        <!-- 战斗属性 -->
        <div class="combat-stats">
          <div class="combat-stat">
            <span>攻击力</span>
            <span class="value">{{ selectedCharacter.attack }}</span>
          </div>
          <div class="combat-stat">
            <span>防御力</span>
            <span class="value">{{ selectedCharacter.defense }}</span>
          </div>
          <div class="combat-stat">
            <span>暴击率</span>
            <span class="value">{{ (selectedCharacter.critRate * 100).toFixed(1) }}%</span>
          </div>
          <div class="combat-stat">
            <span>暴击伤害</span>
            <span class="value">{{ (selectedCharacter.critDamage * 100).toFixed(0) }}%</span>
          </div>
        </div>

        <!-- 统计 -->
        <div class="stats-row">
          <span>击杀: {{ selectedCharacter.totalKills }}</span>
          <span>死亡: {{ selectedCharacter.totalDeaths }}</span>
        </div>
      </div>

      <!-- 右侧：战斗日志区 -->
      <div class="battle-panel">
        <div class="panel-header">
          <h2>战斗日志</h2>
          <span class="battle-indicator" :class="{ running: gameStore.isRunning }">
            {{ gameStore.isRunning ? '⚔️ 战斗中' : '⏸️ 已暂停' }}
          </span>
        </div>

        <!-- 当前敌人 -->
        <div class="enemy-bar" v-if="gameStore.currentEnemy">
          <div class="enemy-info">
            <span class="enemy-name">{{ gameStore.currentEnemy.name }}</span>
            <span class="enemy-level">Lv.{{ gameStore.currentEnemy.level }}</span>
          </div>
          <div class="enemy-hp-track">
            <div class="enemy-hp-fill" :style="{ width: enemyHpPercent + '%' }"></div>
          </div>
          <div class="enemy-hp-text">{{ gameStore.currentEnemy.hp }}/{{ gameStore.currentEnemy.maxHp }}</div>
        </div>
        
        <div class="battle-log" ref="logContainer">
          <div 
            v-for="(log, index) in gameStore.battleLogs" 
            :key="index"
            class="log-entry"
            :style="{ color: getLogColor(log.logType) }"
          >
            <span class="log-time">[{{ formatTime(log.createdAt) }}]</span>
            <span class="log-message">{{ log.message }}</span>
          </div>
          
          <div v-if="gameStore.battleLogs.length === 0" class="log-placeholder">
            <p>🎮 欢迎来到艾泽拉斯！</p>
            <p>点击「开始挂机」按钮开始自动战斗</p>
          </div>
        </div>

        <!-- 控制按钮 -->
        <div class="control-bar">
          <button 
            class="cmd-btn primary" 
            @click="handleToggleBattle"
            :disabled="!selectedCharacter || gameStore.isLoading"
          >
            {{ gameStore.isRunning ? '[⏹] 停止挂机' : '[▶] 开始挂机' }}
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

    <!-- 聊天面板 -->
    <ChatPanel />

    <!-- 小队列表 -->
    <div class="team-bar" v-if="charStore.characters.length > 0">
      <div class="team-label">小队成员:</div>
      <div class="team-list">
        <div 
          v-for="char in charStore.characters" 
          :key="char.id"
          class="team-member"
          :class="{ 
            selected: char.id === selectedCharacter?.id,
            dead: char.isDead 
          }"
          :style="{ '--member-class-color': CLASS_COLORS[char.classId] || '#33ff33' }"
        >
          <span class="member-name" :style="{ color: CLASS_COLORS[char.classId] }">{{ char.name }}</span>
          <span class="member-class">{{ CLASS_NAMES[char.classId] || char.classId }}</span>
          <span class="member-level">Lv.{{ char.level }}</span>
          <div class="member-hp">
            <div class="hp-fill" :style="{ width: (char.hp / char.maxHp * 100) + '%' }"></div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.game-screen {
  display: flex;
  flex-direction: column;
  height: 100vh;
  padding: 15px;
  gap: 15px;
}

/* 顶部栏 */
.top-bar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 10px 15px;
  border: 1px solid var(--terminal-gray);
  background: rgba(0, 0, 0, 0.3);
}

.user-info {
  display: flex;
  align-items: center;
  gap: 15px;
}

.username {
  color: var(--terminal-green);
  font-size: 14px;
}

.separator {
  color: var(--terminal-gray);
}

.gold {
  color: var(--terminal-gold);
  font-size: 13px;
}

.battle-stats {
  display: flex;
  align-items: center;
  gap: 10px;
  font-size: 12px;
  color: var(--terminal-gray);
}

.battle-stats .exp {
  color: #9b59b6;
}

.battle-stats .gold-gain {
  color: var(--terminal-gold);
}

.actions {
  display: flex;
  gap: 15px;
}

.text-btn {
  background: none;
  border: none;
  color: var(--terminal-cyan);
  font-family: inherit;
  font-size: 12px;
  cursor: pointer;
  transition: color 0.3s;
}

.text-btn:hover {
  color: var(--terminal-green);
}

.text-btn.logout:hover {
  color: var(--terminal-red);
}

/* 主内容区 */
.main-content {
  flex: 1;
  display: grid;
  grid-template-columns: 320px 1fr;
  gap: 15px;
  min-height: 0;
}

/* 面板通用样式 */
.character-panel,
.battle-panel {
  border: 1px solid var(--terminal-green);
  background: rgba(0, 50, 0, 0.2);
  display: flex;
  flex-direction: column;
}

.panel-header {
  padding: 15px;
  border-bottom: 1px solid var(--terminal-gray);
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 10px;
}

.panel-header h2 {
  color: var(--terminal-green);
  font-size: 14px;
  margin: 0;
  text-transform: uppercase;
  letter-spacing: 2px;
}

.battle-indicator {
  font-size: 11px;
  color: var(--terminal-gray);
}

.battle-indicator.running {
  color: var(--terminal-green);
  animation: pulse 1s infinite;
}

@keyframes pulse {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.5; }
}

/* 角色面板 */
.faction-badge {
  padding: 2px 8px;
  font-size: 10px;
  text-transform: uppercase;
  color: white;
}

.character-name {
  color: var(--terminal-gold) !important;
  font-size: 16px !important;
}

.character-class {
  padding: 10px 15px;
  color: var(--terminal-cyan);
  font-size: 12px;
  border-bottom: 1px solid var(--terminal-gray);
}

/* 状态条 */
.stat-bars {
  padding: 15px;
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.stat-bar {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.bar-header {
  display: flex;
  justify-content: space-between;
  font-size: 12px;
  white-space: nowrap;
}

.bar-label {
  color: var(--terminal-gray);
  white-space: nowrap;
}

.bar-value {
  color: var(--terminal-green);
  white-space: nowrap;
}

.bar-track {
  height: 8px;
  background: rgba(255, 255, 255, 0.1);
  position: relative;
}

.bar-track.hp .bar-fill {
  background: linear-gradient(90deg, #2d5016, #4a8c2a);
}

.bar-track.exp .bar-fill {
  background: linear-gradient(90deg, #6b21a8, #9333ea);
}

.bar-fill {
  height: 100%;
  transition: width 0.3s ease;
}

/* 属性 */
.attributes {
  display: grid;
  grid-template-columns: repeat(5, 1fr);
  gap: 3px;
  padding: 0 10px 15px;
  border-bottom: 1px solid var(--terminal-gray);
}

.attr {
  text-align: center;
  padding: 6px 2px;
  background: rgba(0, 0, 0, 0.3);
  white-space: nowrap;
}

.attr-label {
  display: block;
  font-size: 11px;
  color: var(--terminal-gray);
  margin-bottom: 3px;
  white-space: nowrap;
}

.attr-value {
  font-size: 14px;
  font-weight: bold;
  white-space: nowrap;
}

.attr-value.str { color: #ff6b6b; }
.attr-value.agi { color: #69db7c; }
.attr-value.int { color: #74c0fc; }
.attr-value.sta { color: #ffd43b; }
.attr-value.spi { color: #da77f2; }

/* 战斗属性 */
.combat-stats {
  padding: 15px;
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 6px;
}

.combat-stat {
  display: flex;
  justify-content: space-between;
  font-size: 12px;
  padding: 5px 8px;
  background: rgba(0, 0, 0, 0.3);
  white-space: nowrap;
}

.combat-stat span:first-child {
  color: var(--terminal-gray);
}

.combat-stat .value {
  color: var(--terminal-green);
}

.stats-row {
  padding: 10px 15px;
  display: flex;
  justify-content: space-around;
  font-size: 11px;
  color: var(--terminal-gray);
  border-top: 1px solid var(--terminal-gray);
  margin-top: auto;
}

/* 敌人状态条 */
.enemy-bar {
  padding: 10px 15px;
  border-bottom: 1px solid var(--terminal-gray);
  background: rgba(255, 0, 0, 0.05);
}

.enemy-info {
  display: flex;
  justify-content: space-between;
  margin-bottom: 5px;
}

.enemy-name {
  color: var(--terminal-red);
  font-size: 13px;
}

.enemy-level {
  color: var(--terminal-gray);
  font-size: 11px;
}

.enemy-hp-track {
  height: 6px;
  background: rgba(255, 255, 255, 0.1);
  margin-bottom: 3px;
}

.enemy-hp-fill {
  height: 100%;
  background: linear-gradient(90deg, #8b0000, #c41e3a);
  transition: width 0.3s ease;
}

.enemy-hp-text {
  font-size: 10px;
  color: var(--terminal-gray);
  text-align: right;
}

/* 战斗面板 */
.battle-log {
  flex: 1;
  padding: 15px;
  overflow-y: auto;
  min-height: 0;
  font-size: 12px;
  line-height: 1.6;
}

.log-entry {
  margin-bottom: 4px;
  word-break: break-word;
}

.log-time {
  color: var(--terminal-gray);
  margin-right: 8px;
  font-size: 11px;
}

.log-placeholder {
  color: var(--terminal-gray);
  text-align: center;
  padding: 40px 20px;
}

.log-placeholder p {
  margin: 10px 0;
}

/* 控制栏 */
.control-bar {
  padding: 15px;
  display: flex;
  gap: 10px;
  border-top: 1px solid var(--terminal-gray);
}

.cmd-btn {
  flex: 1;
  padding: 10px;
  background: transparent;
  border: 1px solid var(--terminal-gray);
  color: var(--terminal-gray);
  font-family: inherit;
  font-size: 11px;
  cursor: pointer;
  transition: all 0.3s;
}

.cmd-btn.primary {
  border-color: var(--terminal-green);
  color: var(--terminal-green);
}

.cmd-btn:not(:disabled):hover {
  border-color: var(--terminal-green);
  color: var(--terminal-green);
  background: rgba(0, 255, 0, 0.1);
}

.cmd-btn.primary:not(:disabled):hover {
  background: var(--terminal-green);
  color: var(--terminal-bg);
}

.cmd-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

/* 小队栏 */
.team-bar {
  display: flex;
  align-items: center;
  gap: 15px;
  padding: 10px 15px;
  border: 1px solid var(--terminal-gray);
  background: rgba(0, 0, 0, 0.3);
}

.team-label {
  color: var(--terminal-gray);
  font-size: 11px;
  white-space: nowrap;
}

.team-list {
  display: flex;
  gap: 10px;
  flex: 1;
  overflow-x: auto;
}

.team-member {
  display: flex;
  flex-direction: column;
  gap: 3px;
  padding: 8px 12px;
  border: 1px solid var(--terminal-gray);
  min-width: 100px;
  cursor: pointer;
  transition: all 0.3s;
}

.team-member:hover,
.team-member.selected {
  border-color: var(--terminal-green);
  background: rgba(0, 255, 0, 0.05);
}

.team-member.dead {
  opacity: 0.5;
  border-color: var(--terminal-red);
}

.member-name {
  font-size: 12px;
  text-shadow: 0 0 5px currentColor;
}

.member-class {
  color: var(--terminal-gray);
  font-size: 10px;
}

.member-level {
  color: var(--terminal-gray);
  font-size: 10px;
}

.member-hp {
  height: 3px;
  background: rgba(255, 255, 255, 0.1);
  margin-top: 3px;
}

.hp-fill {
  height: 100%;
  background: var(--terminal-green);
}

/* 响应式 */
@media (max-width: 768px) {
  .main-content {
    grid-template-columns: 1fr;
  }
  
  .character-panel {
    max-height: 300px;
  }
}
</style>

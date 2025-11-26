<script setup lang="ts">
import { ref, computed, watch, nextTick } from 'vue'
import { useGameStore } from '../stores/game'

const game = useGameStore()
const logContainer = ref<HTMLElement | null>(null)

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
  return (game.character.hp / game.character.max_hp) * 100
})

const mpPercent = computed(() => {
  if (!game.character) return 0
  return (game.character.mp / game.character.max_mp) * 100
})

const expPercent = computed(() => {
  if (!game.character) return 0
  return (game.character.exp / game.character.exp_to_next) * 100
})

const enemyHpPercent = computed(() => {
  if (!game.currentEnemy) return 0
  return (game.currentEnemy.hp / game.currentEnemy.max_hp) * 100
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
</script>

<template>
  <div class="game-screen">
    <!-- 战斗日志区 -->
    <div class="terminal-content" ref="logContainer">
      <!-- 当前敌人信息 -->
      <div v-if="game.currentEnemy" class="enemy-info">
        <span class="enemy-name">
          ⚔ {{ game.currentEnemy.name }} (Lv.{{ game.currentEnemy.level }})
        </span>
        <div class="enemy-hp">
          <span style="color: #888">HP:</span>
          <div class="enemy-bar">
            <div class="enemy-bar-fill" :style="{ width: enemyHpPercent + '%' }"></div>
          </div>
          <span style="color: #ff4444">{{ game.currentEnemy.hp }}/{{ game.currentEnemy.max_hp }}</span>
        </div>
      </div>

      <!-- 战斗日志 -->
      <div class="battle-log">
        <div 
          v-for="(log, index) in game.battleLogs" 
          :key="index"
          class="log-line"
        >
          <span class="log-time">[{{ log.time }}]</span>
          <span 
            class="log-message"
            :class="getLogClass(log.type)"
            :style="{ color: log.color }"
          >
            {{ log.message }}
          </span>
        </div>
        <div class="log-line" v-if="game.isRunning">
          <span class="log-time"></span>
          <span class="log-message" style="color: #00ff00">
            等待下一回合...<span class="cursor"></span>
          </span>
        </div>
      </div>
    </div>

    <!-- 状态栏 -->
    <div class="status-bar">
      <!-- 角色信息行 -->
      <div class="status-row">
        <div class="character-info">
          <span>
            <span class="stat-label">角色: </span>
            <span class="stat-value">{{ game.character?.name }}</span>
          </span>
          <span>
            <span class="stat-value">{{ getRaceName(game.character?.race || '') }} {{ getClassName(game.character?.class || '') }}</span>
          </span>
          <span>
            <span class="stat-label">Lv.</span>
            <span class="stat-value">{{ game.character?.level }}</span>
          </span>
        </div>
        <div class="character-info">
          <span>
            <span class="stat-label">区域: </span>
            <span class="stat-value" style="color: #00ffff">{{ game.battleStatus.current_zone || '未知' }}</span>
          </span>
          <span>
            <span class="stat-label">击杀: </span>
            <span class="stat-value">{{ game.battleStatus.session_kills }}</span>
          </span>
          <span>
            <span class="stat-label">金币: </span>
            <span class="stat-gold">{{ game.character?.gold }} G</span>
          </span>
        </div>
      </div>

      <!-- HP/MP/EXP条 -->
      <div class="status-row">
        <div class="bar-container">
          <span class="bar-label">HP:</span>
          <div class="bar-wrapper">
            <div class="bar bar-hp">
              <div class="bar-fill" :style="{ width: hpPercent + '%' }"></div>
            </div>
            <span class="bar-text">{{ game.character?.hp }}/{{ game.character?.max_hp }}</span>
          </div>
        </div>
        <div class="bar-container">
          <span class="bar-label">MP:</span>
          <div class="bar-wrapper">
            <div class="bar bar-mp">
              <div class="bar-fill" :style="{ width: mpPercent + '%' }"></div>
            </div>
            <span class="bar-text">{{ game.character?.mp }}/{{ game.character?.max_mp }}</span>
          </div>
        </div>
        <div class="bar-container">
          <span class="bar-label">EXP:</span>
          <div class="bar-wrapper">
            <div class="bar bar-exp">
              <div class="bar-fill" :style="{ width: expPercent + '%' }"></div>
            </div>
            <span class="bar-text">{{ game.character?.exp }}/{{ game.character?.exp_to_next }}</span>
          </div>
        </div>
      </div>
    </div>

    <!-- 控制栏 -->
    <div class="control-bar">
      <button 
        class="cmd-btn" 
        :class="{ active: game.isRunning }"
        @click="game.toggleBattle"
      >
        {{ game.isRunning ? '[■] 停止' : '[▶] 开始挂机' }}
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
</template>

<style scoped>
.game-screen {
  display: flex;
  flex-direction: column;
  height: 100%;
}

.terminal-content {
  flex: 1;
  min-height: 0;
}
</style>




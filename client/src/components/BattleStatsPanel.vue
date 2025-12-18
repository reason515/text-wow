<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { get, post } from '../api/client'
import type { 
  BattleStatsOverview, 
  CharacterLifetimeStats, 
  BattleRecord,
  DailyStatistics,
  SessionStats,
  BattleDPSAnalysis
} from '../types/game'

const props = defineProps<{
  characterId?: number
}>()

const emit = defineEmits<{
  close: []
}>()

const loading = ref(false)
const error = ref<string | null>(null)
const overview = ref<BattleStatsOverview | null>(null)
const activeTab = ref<'overview' | 'lifetime' | 'daily' | 'recent' | 'cumulative'>('overview')
const cumulativeDpsAnalysis = ref<BattleDPSAnalysis | null>(null)
const loadingCumulative = ref(false)
const statsSessionActive = ref(false)
const statsSessionStartTime = ref<string | null>(null)
const loadingSessionStatus = ref(false)
const battleDpsAnalysis = ref<BattleDPSAnalysis | null>(null)
const loadingBattleDps = ref(false)
const selectedBattleId = ref<number | null>(null)

// 加载统计数据
async function loadStats() {
  loading.value = true
  error.value = null
  
  try {
    const response = await get<BattleStatsOverview>('/stats/overview')
    if (response.success && response.data) {
      overview.value = response.data
    } else {
      error.value = response.error || '加载统计数据失败'
    }
  } catch (e) {
    console.error('Failed to load stats:', e)
    error.value = '加载统计数据时发生错误'
  } finally {
    loading.value = false
  }
}

// 格式化数字
function formatNumber(num: number): string {
  if (num >= 1000000) {
    return (num / 1000000).toFixed(2) + 'M'
  }
  if (num >= 1000) {
    return (num / 1000).toFixed(2) + 'K'
  }
  return num.toString()
}

// 格式化时间（秒）
function formatTime(seconds: number): string {
  const hours = Math.floor(seconds / 3600)
  const minutes = Math.floor((seconds % 3600) / 60)
  const secs = seconds % 60
  
  if (hours > 0) {
    return `${hours}小时${minutes}分钟`
  }
  if (minutes > 0) {
    return `${minutes}分钟${secs}秒`
  }
  return `${secs}秒`
}

// 格式化日期
function formatDate(dateStr: string): string {
  const date = new Date(dateStr)
  return date.toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  })
}

// 计算胜率
function getWinRate(victories: number, total: number): number {
  if (total === 0) return 0
  return (victories / total) * 100
}

// 计算K/D比
function getKDRatio(kills: number, deaths: number): number {
  if (deaths === 0) return kills > 0 ? kills : 0
  return kills / deaths
}

// 格式化DPS
function formatDPS(dps: number): string {
  if (dps >= 1000) {
    return dps.toFixed(1) + 'K'
  }
  return dps.toFixed(1)
}

// 获取伤害类型名称
function getDamageTypeName(type: string): string {
  const names: Record<string, string> = {
    physical: '物理',
    magic: '魔法',
    fire: '火焰',
    frost: '冰霜',
    shadow: '暗影',
    holy: '神圣',
    nature: '自然',
    dot: '持续伤害'
  }
  return names[type] || type
}

// 获取伤害类型颜色
function getDamageTypeColor(type: string): string {
  const colors: Record<string, string> = {
    physical: '#ff6b6b',
    magic: '#4ecdc4',
    fire: '#ff6b35',
    frost: '#95e1d3',
    shadow: '#6c5ce7',
    holy: '#feca57',
    nature: '#48dbfb',
    dot: '#a29bfe'
  }
  return colors[type] || '#ffffff'
}

// 获取伤害类型数值
function getDamageTypeValue(type: string, composition: any): number {
  const values: Record<string, number> = {
    physical: composition.physical || 0,
    magic: composition.magic || 0,
    fire: composition.fire || 0,
    frost: composition.frost || 0,
    shadow: composition.shadow || 0,
    holy: composition.holy || 0,
    nature: composition.nature || 0,
    dot: composition.dot || 0
  }
  return values[type] || 0
}

// 获取技能显示名称
function getSkillDisplayName(skill: any): string {
  // 如果有技能名称，直接使用
  if (skill.skillName && skill.skillName.trim() !== '') {
    return skill.skillName
  }
  // 如果skillId为空或者是普通攻击相关的，显示"普通攻击"
  if (!skill.skillId || skill.skillId === '' || skill.skillId === 'normal_attack') {
    return '普通攻击'
  }
  // 否则使用skillId
  return skill.skillId
}

// 获取累计战斗场次
function getCumulativeBattleCount(): number {
  if (cumulativeDpsAnalysis.value && cumulativeDpsAnalysis.value.battleCount) {
    return cumulativeDpsAnalysis.value.battleCount
  }
  return 0
}

// 检查是否有任何伤害类型数据
function hasAnyDamageType(composition: any): boolean {
  if (!composition || !composition.percentages) {
    return false
  }
  for (const value of Object.values(composition.percentages) as number[]) {
    if (value > 0) {
      return true
    }
  }
  return false
}

// 加载统计会话状态
async function loadSessionStatus() {
  loadingSessionStatus.value = true
  try {
    const response = await get<{ isActive: boolean; startTime?: string }>('/stats/session/status')
    if (response.success && response.data) {
      statsSessionActive.value = response.data.isActive
      statsSessionStartTime.value = response.data.startTime || null
    }
  } catch (e) {
    console.error('Failed to load session status:', e)
  } finally {
    loadingSessionStatus.value = false
  }
}

// 开始统计会话
async function startStatsSession() {
  try {
    const response = await post('/stats/session/start')
    if (response.success) {
      statsSessionActive.value = true
      statsSessionStartTime.value = new Date().toISOString()
      await loadCumulativeDPS()
      activeTab.value = 'cumulative'
    } else {
      error.value = response.error || '开始统计失败'
    }
  } catch (e) {
    console.error('Failed to start stats session:', e)
    error.value = '开始统计时发生错误'
  }
}

// 重置统计会话
async function resetStatsSession() {
  if (!confirm('确定要重置统计吗？这将清除当前的累计统计数据。')) {
    return
  }
  
  try {
    const response = await post('/stats/session/reset')
    if (response.success) {
      statsSessionActive.value = false
      statsSessionStartTime.value = null
      cumulativeDpsAnalysis.value = null
      // 保持在 DPS统计 页，显示“未开始统计”的状态
    } else {
      error.value = response.error || '重置统计失败'
    }
  } catch (e) {
    console.error('Failed to reset stats session:', e)
    error.value = '重置统计时发生错误'
  }
}

// 加载累计DPS分析
async function loadCumulativeDPS() {
  if (!statsSessionActive.value) {
    return
  }
  
  loadingCumulative.value = true
  error.value = null
  
  try {
    const response = await get<BattleDPSAnalysis>('/stats/cumulative/dps')
    if (response.success && response.data) {
      cumulativeDpsAnalysis.value = response.data
    } else {
      error.value = response.error || '加载累计DPS分析失败'
    }
  } catch (e) {
    console.error('Failed to load cumulative DPS:', e)
    error.value = '加载累计DPS分析时发生错误'
  } finally {
    loadingCumulative.value = false
  }
}

// 加载单场战斗DPS分析
async function loadBattleDPS(battleId: number) {
  loadingBattleDps.value = true
  error.value = null
  selectedBattleId.value = battleId
  
  try {
    const response = await get<BattleDPSAnalysis>(`/stats/battles/${battleId}/dps`)
    if (response.success && response.data) {
      battleDpsAnalysis.value = response.data
      // 切换到累计统计标签页以显示DPS分析
      activeTab.value = 'cumulative'
    } else {
      error.value = response.error || '加载战斗DPS分析失败'
    }
  } catch (e) {
    console.error('Failed to load battle DPS:', e)
    error.value = '加载战斗DPS分析时发生错误'
  } finally {
    loadingBattleDps.value = false
  }
}

// 关闭单场战斗DPS分析
function closeBattleDPS() {
  battleDpsAnalysis.value = null
  selectedBattleId.value = null
}

onMounted(() => {
  loadStats()
  loadSessionStatus()
})
</script>

<template>
  <div class="stats-panel-overlay" @click.self="emit('close')">
    <div class="stats-panel">
      <div class="stats-panel-header">
        <h2 class="stats-panel-title">战斗统计</h2>
        <button class="stats-panel-close" @click="emit('close')">×</button>
      </div>

      <!-- 标签页 -->
      <div class="stats-tabs">
        <button 
          class="stats-tab" 
          :class="{ active: activeTab === 'overview' }"
          @click="activeTab = 'overview'"
        >
          概览
        </button>
        <button 
          class="stats-tab" 
          :class="{ active: activeTab === 'lifetime' }"
          @click="activeTab = 'lifetime'"
        >
          生涯统计
        </button>
        <button 
          class="stats-tab" 
          :class="{ active: activeTab === 'daily' }"
          @click="activeTab = 'daily'"
        >
          每日统计
        </button>
        <button 
          class="stats-tab" 
          :class="{ active: activeTab === 'recent' }"
          @click="activeTab = 'recent'"
        >
          最近战斗
        </button>
        <button 
          class="stats-tab" 
          :class="{ active: activeTab === 'cumulative' }"
          @click="activeTab = 'cumulative'; loadCumulativeDPS()"
        >
          DPS统计
        </button>
      </div>

      <!-- 内容区域 -->
      <div class="stats-content">
        <div v-if="loading" class="stats-loading">加载中...</div>
        <div v-else-if="error" class="stats-error">{{ error }}</div>
        
        <!-- 概览 -->
        <div v-else-if="activeTab === 'overview' && overview" class="stats-overview">
          <!-- 会话统计 -->
          <div v-if="overview.sessionStats" class="stats-section">
            <h3 class="stats-section-title">当前会话</h3>
            <div class="stats-grid">
              <div class="stats-item">
                <span class="stats-label">战斗场次</span>
                <span class="stats-value">{{ overview.sessionStats.totalBattles }}</span>
              </div>
              <div class="stats-item">
                <span class="stats-label">击杀数</span>
                <span class="stats-value">{{ formatNumber(overview.sessionStats.totalKills) }}</span>
              </div>
              <div class="stats-item">
                <span class="stats-label">获得经验</span>
                <span class="stats-value">{{ formatNumber(overview.sessionStats.totalExp) }}</span>
              </div>
              <div class="stats-item">
                <span class="stats-label">获得金币</span>
                <span class="stats-value">{{ formatNumber(overview.sessionStats.totalGold) }}</span>
              </div>
              <div class="stats-item">
                <span class="stats-label">游戏时长</span>
                <span class="stats-value">{{ formatTime(overview.sessionStats.durationSeconds) }}</span>
              </div>
            </div>
          </div>

          <!-- 今日统计 -->
          <div v-if="overview.todayStats" class="stats-section">
            <h3 class="stats-section-title">今日统计</h3>
            <div class="stats-grid">
              <div class="stats-item">
                <span class="stats-label">战斗场次</span>
                <span class="stats-value">{{ overview.todayStats.battlesCount }}</span>
              </div>
              <div class="stats-item">
                <span class="stats-label">胜利</span>
                <span class="stats-value stats-success">{{ overview.todayStats.victories }}</span>
              </div>
              <div class="stats-item">
                <span class="stats-label">失败</span>
                <span class="stats-value stats-danger">{{ overview.todayStats.defeats }}</span>
              </div>
              <div class="stats-item">
                <span class="stats-label">总伤害</span>
                <span class="stats-value">{{ formatNumber(overview.todayStats.totalDamage) }}</span>
              </div>
              <div class="stats-item">
                <span class="stats-label">总治疗</span>
                <span class="stats-value">{{ formatNumber(overview.todayStats.totalHealing) }}</span>
              </div>
              <div class="stats-item">
                <span class="stats-label">击杀数</span>
                <span class="stats-value">{{ formatNumber(overview.todayStats.kills) }}</span>
              </div>
            </div>
          </div>

          <!-- 角色生涯统计 -->
          <div v-if="overview.lifetimeStats && overview.lifetimeStats.length > 0" class="stats-section">
            <h3 class="stats-section-title">角色生涯统计</h3>
            <div 
              v-for="stats in overview.lifetimeStats" 
              :key="stats.characterId"
              class="lifetime-stats-card"
            >
              <div class="lifetime-stats-header">
                <span class="lifetime-stats-title">角色 #{{ stats.characterId }}</span>
                <span class="lifetime-stats-battles">{{ stats.totalBattles }} 场战斗</span>
              </div>
              <div class="stats-grid">
                <div class="stats-item">
                  <span class="stats-label">胜率</span>
                  <span class="stats-value">{{ getWinRate(stats.victories, stats.totalBattles).toFixed(1) }}%</span>
                </div>
                <div class="stats-item">
                  <span class="stats-label">总伤害</span>
                  <span class="stats-value">{{ formatNumber(stats.totalDamageDealt) }}</span>
                </div>
                <div class="stats-item">
                  <span class="stats-label">总治疗</span>
                  <span class="stats-value">{{ formatNumber(stats.totalHealingDone) }}</span>
                </div>
                <div class="stats-item">
                  <span class="stats-label">击杀/死亡</span>
                  <span class="stats-value">{{ stats.totalKills }} / {{ stats.totalDeaths }}</span>
                </div>
                <div class="stats-item">
                  <span class="stats-label">K/D比</span>
                  <span class="stats-value">{{ getKDRatio(stats.totalKills, stats.totalDeaths).toFixed(2) }}</span>
                </div>
                <div class="stats-item">
                  <span class="stats-label">最高单次伤害</span>
                  <span class="stats-value">{{ formatNumber(stats.highestDamageSingle) }}</span>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- 生涯统计 -->
        <div v-else-if="activeTab === 'lifetime' && overview?.lifetimeStats" class="stats-lifetime">
          <div 
            v-for="stats in overview.lifetimeStats" 
            :key="stats.characterId"
            class="lifetime-detail-card"
          >
            <h3 class="lifetime-detail-title">角色 #{{ stats.characterId }}</h3>
            
            <div class="stats-section">
              <h4 class="stats-subtitle">战斗记录</h4>
              <div class="stats-grid">
                <div class="stats-item">
                  <span class="stats-label">总战斗</span>
                  <span class="stats-value">{{ stats.totalBattles }}</span>
                </div>
                <div class="stats-item">
                  <span class="stats-label">胜利</span>
                  <span class="stats-value stats-success">{{ stats.victories }}</span>
                </div>
                <div class="stats-item">
                  <span class="stats-label">失败</span>
                  <span class="stats-value stats-danger">{{ stats.defeats }}</span>
                </div>
                <div class="stats-item">
                  <span class="stats-label">胜率</span>
                  <span class="stats-value">{{ getWinRate(stats.victories, stats.totalBattles).toFixed(1) }}%</span>
                </div>
                <div class="stats-item">
                  <span class="stats-label">PVE战斗</span>
                  <span class="stats-value">{{ stats.pveBattles }}</span>
                </div>
                <div class="stats-item">
                  <span class="stats-label">PVP战斗</span>
                  <span class="stats-value">{{ stats.pvpBattles }}</span>
                </div>
              </div>
            </div>

            <div class="stats-section">
              <h4 class="stats-subtitle">伤害统计</h4>
              <div class="stats-grid">
                <div class="stats-item">
                  <span class="stats-label">总伤害</span>
                  <span class="stats-value">{{ formatNumber(stats.totalDamageDealt) }}</span>
                </div>
                <div class="stats-item">
                  <span class="stats-label">物理伤害</span>
                  <span class="stats-value">{{ formatNumber(stats.totalPhysicalDamage) }}</span>
                </div>
                <div class="stats-item">
                  <span class="stats-label">魔法伤害</span>
                  <span class="stats-value">{{ formatNumber(stats.totalMagicDamage) }}</span>
                </div>
                <div class="stats-item">
                  <span class="stats-label">暴击次数</span>
                  <span class="stats-value">{{ formatNumber(stats.totalCritCount) }}</span>
                </div>
                <div class="stats-item">
                  <span class="stats-label">暴击伤害</span>
                  <span class="stats-value">{{ formatNumber(stats.totalCritDamage) }}</span>
                </div>
                <div class="stats-item">
                  <span class="stats-label">最高单次伤害</span>
                  <span class="stats-value stats-highlight">{{ formatNumber(stats.highestDamageSingle) }}</span>
                </div>
                <div class="stats-item">
                  <span class="stats-label">最高单场伤害</span>
                  <span class="stats-value stats-highlight">{{ formatNumber(stats.highestDamageBattle) }}</span>
                </div>
              </div>
            </div>

            <div class="stats-section">
              <h4 class="stats-subtitle">承伤统计</h4>
              <div class="stats-grid">
                <div class="stats-item">
                  <span class="stats-label">总承伤</span>
                  <span class="stats-value">{{ formatNumber(stats.totalDamageTaken) }}</span>
                </div>
                <div class="stats-item">
                  <span class="stats-label">格挡伤害</span>
                  <span class="stats-value">{{ formatNumber(stats.totalDamageBlocked) }}</span>
                </div>
                <div class="stats-item">
                  <span class="stats-label">吸收伤害</span>
                  <span class="stats-value">{{ formatNumber(stats.totalDamageAbsorbed) }}</span>
                </div>
                <div class="stats-item">
                  <span class="stats-label">闪避次数</span>
                  <span class="stats-value">{{ formatNumber(stats.totalDodgeCount) }}</span>
                </div>
              </div>
            </div>

            <div class="stats-section">
              <h4 class="stats-subtitle">治疗统计</h4>
              <div class="stats-grid">
                <div class="stats-item">
                  <span class="stats-label">总治疗</span>
                  <span class="stats-value">{{ formatNumber(stats.totalHealingDone) }}</span>
                </div>
                <div class="stats-item">
                  <span class="stats-label">受到治疗</span>
                  <span class="stats-value">{{ formatNumber(stats.totalHealingReceived) }}</span>
                </div>
                <div class="stats-item">
                  <span class="stats-label">过量治疗</span>
                  <span class="stats-value">{{ formatNumber(stats.totalOverhealing) }}</span>
                </div>
                <div class="stats-item">
                  <span class="stats-label">最高单次治疗</span>
                  <span class="stats-value stats-highlight">{{ formatNumber(stats.highestHealingSingle) }}</span>
                </div>
              </div>
            </div>

            <div class="stats-section">
              <h4 class="stats-subtitle">击杀与死亡</h4>
              <div class="stats-grid">
                <div class="stats-item">
                  <span class="stats-label">总击杀</span>
                  <span class="stats-value stats-success">{{ formatNumber(stats.totalKills) }}</span>
                </div>
                <div class="stats-item">
                  <span class="stats-label">总死亡</span>
                  <span class="stats-value stats-danger">{{ formatNumber(stats.totalDeaths) }}</span>
                </div>
                <div class="stats-item">
                  <span class="stats-label">K/D比</span>
                  <span class="stats-value">{{ getKDRatio(stats.totalKills, stats.totalDeaths).toFixed(2) }}</span>
                </div>
                <div class="stats-item">
                  <span class="stats-label">最长连杀</span>
                  <span class="stats-value stats-highlight">{{ stats.killStreakBest }}</span>
                </div>
                <div class="stats-item">
                  <span class="stats-label">当前连杀</span>
                  <span class="stats-value">{{ stats.currentKillStreak }}</span>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- 每日统计 -->
        <div v-else-if="activeTab === 'daily' && overview?.todayStats" class="stats-daily">
          <div class="stats-section">
            <h3 class="stats-section-title">今日统计 ({{ overview.todayStats.statDate }})</h3>
            <div class="stats-grid">
              <div class="stats-item">
                <span class="stats-label">战斗场次</span>
                <span class="stats-value">{{ overview.todayStats.battlesCount }}</span>
              </div>
              <div class="stats-item">
                <span class="stats-label">胜利</span>
                <span class="stats-value stats-success">{{ overview.todayStats.victories }}</span>
              </div>
              <div class="stats-item">
                <span class="stats-label">失败</span>
                <span class="stats-value stats-danger">{{ overview.todayStats.defeats }}</span>
              </div>
              <div class="stats-item">
                <span class="stats-label">总伤害</span>
                <span class="stats-value">{{ formatNumber(overview.todayStats.totalDamage) }}</span>
              </div>
              <div class="stats-item">
                <span class="stats-label">总治疗</span>
                <span class="stats-value">{{ formatNumber(overview.todayStats.totalHealing) }}</span>
              </div>
              <div class="stats-item">
                <span class="stats-label">总承伤</span>
                <span class="stats-value">{{ formatNumber(overview.todayStats.totalDamageTaken) }}</span>
              </div>
              <div class="stats-item">
                <span class="stats-label">获得经验</span>
                <span class="stats-value">{{ formatNumber(overview.todayStats.expGained) }}</span>
              </div>
              <div class="stats-item">
                <span class="stats-label">获得金币</span>
                <span class="stats-value">{{ formatNumber(overview.todayStats.goldGained) }}</span>
              </div>
              <div class="stats-item">
                <span class="stats-label">游戏时长</span>
                <span class="stats-value">{{ formatTime(overview.todayStats.playTime) }}</span>
              </div>
              <div class="stats-item">
                <span class="stats-label">击杀数</span>
                <span class="stats-value stats-success">{{ formatNumber(overview.todayStats.kills) }}</span>
              </div>
              <div class="stats-item">
                <span class="stats-label">死亡数</span>
                <span class="stats-value stats-danger">{{ formatNumber(overview.todayStats.deaths) }}</span>
              </div>
            </div>
          </div>
        </div>

        <!-- 最近战斗 -->
        <div v-else-if="activeTab === 'recent' && overview?.recentBattles" class="stats-recent">
          <div v-if="overview.recentBattles.length === 0" class="stats-empty">
            暂无战斗记录
          </div>
          <div v-else class="battle-list">
            <div 
              v-for="battle in overview.recentBattles" 
              :key="battle.id"
              class="battle-item"
            >
              <div class="battle-header">
                <span class="battle-result" :class="battle.result">
                  {{ battle.result === 'victory' ? '胜利' : battle.result === 'defeat' ? '失败' : battle.result }}
                </span>
                <span class="battle-date">{{ formatDate(battle.createdAt) }}</span>
              </div>
              <div class="battle-info">
                <div class="battle-info-item">
                  <span class="battle-info-label">区域:</span>
                  <span class="battle-info-value">{{ battle.zoneId }}</span>
                </div>
                <div class="battle-info-item">
                  <span class="battle-info-label">类型:</span>
                  <span class="battle-info-value">{{ battle.battleType }}</span>
                </div>
                <div class="battle-info-item">
                  <span class="battle-info-label">回合数:</span>
                  <span class="battle-info-value">{{ battle.totalRounds }}</span>
                </div>
                <div class="battle-info-item">
                  <span class="battle-info-label">时长:</span>
                  <span class="battle-info-value">{{ formatTime(battle.durationSeconds) }}</span>
                </div>
              </div>
              <div class="battle-stats">
                <div class="battle-stat">
                  <span class="battle-stat-label">伤害:</span>
                  <span class="battle-stat-value">{{ formatNumber(battle.teamDamageDealt) }}</span>
                </div>
                <div class="battle-stat">
                  <span class="battle-stat-label">承伤:</span>
                  <span class="battle-stat-value">{{ formatNumber(battle.teamDamageTaken) }}</span>
                </div>
                <div class="battle-stat">
                  <span class="battle-stat-label">治疗:</span>
                  <span class="battle-stat-value">{{ formatNumber(battle.teamHealingDone) }}</span>
                </div>
                <div class="battle-stat">
                  <span class="battle-stat-label">经验:</span>
                  <span class="battle-stat-value">{{ formatNumber(battle.expGained) }}</span>
                </div>
                <div class="battle-stat">
                  <span class="battle-stat-label">金币:</span>
                  <span class="battle-stat-value">{{ formatNumber(battle.goldGained) }}</span>
                </div>
              </div>
              <div class="battle-actions">
                <button 
                  class="battle-action-btn"
                  @click="loadBattleDPS(battle.id)"
                  :disabled="loadingBattleDps"
                >
                  {{ loadingBattleDps && selectedBattleId === battle.id ? '加载中...' : '查看DPS分析' }}
                </button>
              </div>
            </div>
          </div>
        </div>

        <!-- 累计统计 / 单场战斗DPS分析 -->
        <div v-else-if="activeTab === 'cumulative'" class="stats-cumulative">
          <!-- 单场战斗DPS分析 -->
          <div v-if="battleDpsAnalysis" class="battle-dps-header">
            <h3 class="battle-dps-title">单场战斗DPS分析</h3>
            <button class="battle-dps-close" @click="closeBattleDPS">×</button>
          </div>
          
          <!-- DPS统计页内的会话控制（仅累计统计时显示） -->
          <div v-if="!battleDpsAnalysis" class="stats-session-controls">
            <div class="stats-session-status">
              <span v-if="statsSessionActive" class="session-active">
                ● 统计中 ({{ statsSessionStartTime ? formatDate(statsSessionStartTime) : '已开始' }})
              </span>
              <span v-else class="session-inactive">○ 未开始统计</span>
            </div>
            <div class="stats-session-actions">
              <button 
                v-if="!statsSessionActive"
                class="session-btn session-btn-start"
                @click="startStatsSession"
              >
                开始统计
              </button>
              <button 
                v-else
                class="session-btn session-btn-reset"
                @click="resetStatsSession"
              >
                重置统计
              </button>
            </div>
          </div>

          <div v-if="loadingCumulative || loadingBattleDps" class="stats-loading">
            {{ battleDpsAnalysis ? '加载战斗DPS分析中...' : '加载累计统计中...' }}
          </div>
          <div v-else-if="battleDpsAnalysis && (!battleDpsAnalysis.characters || battleDpsAnalysis.characters.length === 0)" class="stats-empty">
            <div class="dps-empty-message">
              <p>该战斗暂无DPS数据</p>
            </div>
          </div>
          <div v-else-if="!battleDpsAnalysis && !statsSessionActive" class="stats-empty">
            <div class="dps-empty-message">
              <p>请先开始统计会话</p>
              <p class="dps-empty-hint">点击上方的"开始统计"按钮开始累计统计</p>
            </div>
          </div>
          <div v-else-if="!battleDpsAnalysis && (!cumulativeDpsAnalysis || (cumulativeDpsAnalysis.characters && cumulativeDpsAnalysis.characters.length === 0))" class="stats-empty">
            <div class="dps-empty-message">
              <p>暂无累计统计数据</p>
              <p class="dps-empty-hint">开始战斗后，统计数据将在这里显示</p>
            </div>
          </div>
          <div v-else class="dps-analysis">
            <!-- DPS总览 -->
            <div class="stats-section">
              <h3 class="stats-section-title">
                {{ battleDpsAnalysis ? '单场战斗DPS统计' : '累计DPS统计' }}
                <span v-if="!battleDpsAnalysis && statsSessionStartTime" class="session-time-badge">
                  自 {{ formatDate(statsSessionStartTime) }}
                </span>
              </h3>
              <div class="dps-overview">
                <div v-if="!battleDpsAnalysis" class="dps-overview-item">
                  <span class="dps-overview-label">战斗场次</span>
                  <span class="dps-overview-value">{{ getCumulativeBattleCount() }}</span>
                </div>
                <div class="dps-overview-item">
                  <span class="dps-overview-label">{{ battleDpsAnalysis ? '战斗时长' : '累计时长' }}</span>
                  <span class="dps-overview-value">{{ formatTime((battleDpsAnalysis || cumulativeDpsAnalysis).duration) }}</span>
                </div>
                <div class="dps-overview-item">
                  <span class="dps-overview-label">总回合数</span>
                  <span class="dps-overview-value">{{ (battleDpsAnalysis || cumulativeDpsAnalysis).totalRounds }}</span>
                </div>
                <div class="dps-overview-item">
                  <span class="dps-overview-label">队伍DPS</span>
                  <span class="dps-overview-value stats-highlight">{{ formatDPS((battleDpsAnalysis || cumulativeDpsAnalysis).teamDps) }}</span>
                </div>
                <div class="dps-overview-item">
                  <span class="dps-overview-label">队伍HPS</span>
                  <span class="dps-overview-value">{{ formatDPS((battleDpsAnalysis || cumulativeDpsAnalysis).teamHps) }}</span>
                </div>
              </div>
            </div>

            <!-- 伤害构成 -->
            <div v-if="(battleDpsAnalysis || cumulativeDpsAnalysis)?.teamDamageComposition && hasAnyDamageType((battleDpsAnalysis || cumulativeDpsAnalysis).teamDamageComposition)" class="stats-section">
              <h3 class="stats-section-title">{{ battleDpsAnalysis ? '战斗伤害构成' : '累计伤害构成' }}</h3>
              <div class="damage-composition">
                <template v-for="(value, key) in (battleDpsAnalysis || cumulativeDpsAnalysis).teamDamageComposition.percentages" :key="key">
                  <div 
                    v-if="value > 0"
                    class="damage-type-item"
                  >
                    <div class="damage-type-header">
                      <span class="damage-type-name">{{ getDamageTypeName(key) }}</span>
                      <span class="damage-type-percent">{{ value.toFixed(1) }}%</span>
                    </div>
                    <div class="damage-type-bar">
                      <div 
                        class="damage-type-fill" 
                        :style="{ width: value + '%', backgroundColor: getDamageTypeColor(key) }"
                      ></div>
                    </div>
                    <span class="damage-type-value">{{ formatNumber(getDamageTypeValue(key, (battleDpsAnalysis || cumulativeDpsAnalysis).teamDamageComposition)) }}</span>
                  </div>
                </template>
              </div>
            </div>

            <!-- 角色DPS分析 -->
            <div 
              v-for="character in (battleDpsAnalysis || cumulativeDpsAnalysis).characters" 
              :key="character.characterId"
              class="stats-section character-dps-section"
            >
              <h3 class="stats-section-title">
                {{ character.characterName || `角色 #${character.characterId}` }}
                <span class="character-dps-badge">{{ formatDPS(character.totalDps) }} DPS</span>
              </h3>

              <!-- 角色总览 -->
              <div class="character-dps-overview">
                <div class="dps-overview-item">
                  <span class="dps-overview-label">总伤害</span>
                  <span class="dps-overview-value">{{ formatNumber(character.totalDamage) }}</span>
                </div>
                <div class="dps-overview-item">
                  <span class="dps-overview-label">总治疗</span>
                  <span class="dps-overview-value">{{ formatNumber(character.totalHealing) }}</span>
                </div>
                <div class="dps-overview-item">
                  <span class="dps-overview-label">DPS</span>
                  <span class="dps-overview-value stats-highlight">{{ formatDPS(character.totalDps) }}</span>
                </div>
                <div class="dps-overview-item">
                  <span class="dps-overview-label">HPS</span>
                  <span class="dps-overview-value">{{ formatDPS(character.totalHps) }}</span>
                </div>
              </div>

              <!-- 技能DPS明细 -->
              <div class="skill-dps-list">
                <h4 class="stats-subtitle">{{ battleDpsAnalysis ? '技能DPS明细' : '累计技能DPS明细' }}</h4>
                <div class="skill-dps-table">
                  <div class="skill-dps-header">
                    <div class="skill-dps-col-name">技能</div>
                    <div class="skill-dps-col-dps">DPS</div>
                    <div class="skill-dps-col-damage">总伤害</div>
                    <div class="skill-dps-col-percent">占比</div>
                    <div class="skill-dps-col-uses">使用</div>
                    <div class="skill-dps-col-avg">平均</div>
                  </div>
                  <div 
                    v-for="skill in character.skillBreakdown" 
                    :key="skill.skillId"
                    class="skill-dps-row"
                  >
                    <div class="skill-dps-col-name">{{ getSkillDisplayName(skill) }}</div>
                    <div class="skill-dps-col-dps stats-highlight">{{ formatDPS(skill.dps) }}</div>
                    <div class="skill-dps-col-damage">{{ formatNumber(skill.totalDamage) }}</div>
                    <div class="skill-dps-col-percent">{{ skill.damagePercent.toFixed(1) }}%</div>
                    <div class="skill-dps-col-uses">{{ skill.useCount }}</div>
                    <div class="skill-dps-col-avg">{{ formatNumber(skill.avgDamage) }}</div>
                  </div>
                </div>
              </div>

              <!-- 角色伤害构成 -->
              <div v-if="character.damageComposition && hasAnyDamageType(character.damageComposition)" class="character-damage-composition">
                <h4 class="stats-subtitle">{{ battleDpsAnalysis ? '伤害类型构成' : '累计伤害类型构成' }}</h4>
                <div class="damage-composition">
                  <template v-for="(value, key) in character.damageComposition.percentages" :key="key">
                    <div 
                      v-if="value > 0"
                      class="damage-type-item"
                    >
                      <div class="damage-type-header">
                        <span class="damage-type-name">{{ getDamageTypeName(key) }}</span>
                        <span class="damage-type-percent">{{ value.toFixed(1) }}%</span>
                      </div>
                      <div class="damage-type-bar">
                        <div 
                          class="damage-type-fill" 
                          :style="{ width: value + '%', backgroundColor: getDamageTypeColor(key) }"
                        ></div>
                      </div>
                      <span class="damage-type-value">{{ formatNumber(getDamageTypeValue(key, character.damageComposition)) }}</span>
                    </div>
                  </template>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.stats-panel-overlay {
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

.stats-panel {
  background: rgba(0, 20, 0, 0.95);
  border: 2px solid var(--terminal-green);
  width: 100%;
  max-width: 900px;
  max-height: 90vh;
  display: flex;
  flex-direction: column;
  box-shadow: 0 0 30px rgba(0, 255, 0, 0.3);
}

.stats-panel-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 20px;
  border-bottom: 1px solid var(--terminal-green);
}

.stats-panel-title {
  color: var(--terminal-gold);
  font-size: 20px;
  margin: 0;
}

.stats-panel-close {
  background: transparent;
  border: 1px solid var(--terminal-gray);
  color: var(--terminal-gray);
  width: 32px;
  height: 32px;
  font-size: 24px;
  cursor: pointer;
  transition: all 0.2s;
  display: flex;
  align-items: center;
  justify-content: center;
}

.stats-panel-close:hover {
  border-color: var(--terminal-red);
  color: var(--terminal-red);
}

.stats-tabs {
  display: flex;
  border-bottom: 1px solid var(--text-dim);
  background: rgba(0, 0, 0, 0.3);
}

.stats-tab {
  flex: 1;
  padding: 12px;
  background: transparent;
  border: none;
  border-bottom: 2px solid transparent;
  color: var(--text-secondary);
  cursor: pointer;
  transition: all 0.2s;
  font-family: inherit;
  font-size: 14px;
}

.stats-tab:hover {
  color: var(--terminal-green);
  background: rgba(0, 255, 0, 0.05);
}

.stats-tab.active {
  color: var(--terminal-cyan);
  border-bottom-color: var(--terminal-cyan);
  background: rgba(0, 255, 255, 0.05);
}

.stats-content {
  flex: 1;
  overflow-y: auto;
  padding: 20px;
}

.stats-loading,
.stats-error {
  text-align: center;
  padding: 40px;
  color: var(--text-secondary);
}

.stats-error {
  color: var(--terminal-red);
}

.stats-section {
  margin-bottom: 24px;
}

.stats-section-title {
  color: var(--terminal-cyan);
  font-size: 16px;
  margin-bottom: 12px;
  padding-bottom: 8px;
  border-bottom: 1px solid var(--text-dim);
}

.stats-subtitle {
  color: var(--terminal-green);
  font-size: 14px;
  margin-bottom: 8px;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(150px, 1fr));
  gap: 12px;
}

.stats-item {
  display: flex;
  flex-direction: column;
  padding: 8px;
  background: rgba(0, 0, 0, 0.3);
  border: 1px solid var(--text-dim);
  border-radius: 4px;
}

.stats-label {
  color: var(--text-secondary);
  font-size: 12px;
  margin-bottom: 4px;
}

.stats-value {
  color: var(--terminal-green);
  font-size: 16px;
  font-weight: bold;
}

.stats-value.stats-success {
  color: var(--terminal-green);
}

.stats-value.stats-danger {
  color: var(--terminal-red);
}

.stats-value.stats-highlight {
  color: var(--terminal-gold);
}

.lifetime-stats-card {
  background: rgba(0, 0, 0, 0.4);
  border: 1px solid var(--text-dim);
  border-radius: 4px;
  padding: 12px;
  margin-bottom: 12px;
}

.lifetime-stats-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
  padding-bottom: 8px;
  border-bottom: 1px solid var(--text-dim);
}

.lifetime-stats-title {
  color: var(--terminal-cyan);
  font-weight: bold;
}

.lifetime-stats-battles {
  color: var(--text-secondary);
  font-size: 12px;
}

.lifetime-detail-card {
  background: rgba(0, 0, 0, 0.4);
  border: 1px solid var(--text-dim);
  border-radius: 4px;
  padding: 16px;
  margin-bottom: 16px;
}

.lifetime-detail-title {
  color: var(--terminal-cyan);
  font-size: 18px;
  margin-bottom: 16px;
  padding-bottom: 8px;
  border-bottom: 1px solid var(--terminal-cyan);
}

.battle-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.battle-item {
  background: rgba(0, 0, 0, 0.4);
  border: 1px solid var(--text-dim);
  border-radius: 4px;
  padding: 12px;
}

.battle-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
  padding-bottom: 8px;
  border-bottom: 1px solid var(--text-dim);
}

.battle-result {
  font-weight: bold;
  font-size: 14px;
}

.battle-result.victory {
  color: var(--terminal-green);
}

.battle-result.defeat {
  color: var(--terminal-red);
}

.battle-date {
  color: var(--text-secondary);
  font-size: 12px;
}

.battle-info {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
  margin-bottom: 8px;
}

.battle-info-item {
  display: flex;
  gap: 4px;
  font-size: 12px;
}

.battle-info-label {
  color: var(--text-secondary);
}

.battle-info-value {
  color: var(--terminal-green);
}

.battle-stats {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
  padding-top: 8px;
  border-top: 1px solid var(--text-dim);
}

.battle-stat {
  display: flex;
  gap: 4px;
  font-size: 12px;
}

.battle-stat-label {
  color: var(--text-secondary);
}

.battle-stat-value {
  color: var(--terminal-cyan);
  font-weight: bold;
}

.stats-empty {
  text-align: center;
  padding: 40px;
  color: var(--text-secondary);
}

.battle-actions {
  margin-top: 12px;
  padding-top: 12px;
  border-top: 1px solid var(--text-dim);
}

.battle-actions {
  margin-top: 12px;
  padding-top: 12px;
  border-top: 1px solid var(--text-dim);
  display: flex;
  justify-content: flex-end;
}

.battle-action-btn {
  width: 100%;
  padding: 8px;
  background: rgba(0, 255, 0, 0.1);
  border: 1px solid var(--terminal-green);
  color: var(--terminal-green);
  cursor: pointer;
  transition: all 0.2s;
  font-family: inherit;
  font-size: 12px;
}

.battle-action-btn:hover {
  background: rgba(0, 255, 0, 0.2);
  transform: translateY(-1px);
}

/* DPS分析样式 */
.dps-analysis {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.dps-overview {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(150px, 1fr));
  gap: 12px;
}

.dps-overview-item {
  display: flex;
  flex-direction: column;
  padding: 12px;
  background: rgba(0, 0, 0, 0.3);
  border: 1px solid var(--text-dim);
  border-radius: 4px;
}

.dps-overview-label {
  color: var(--text-secondary);
  font-size: 12px;
  margin-bottom: 4px;
}

.dps-overview-value {
  color: var(--terminal-cyan);
  font-size: 18px;
  font-weight: bold;
}

.character-dps-section {
  background: rgba(0, 0, 0, 0.2);
  border: 1px solid var(--text-dim);
  border-radius: 4px;
  padding: 16px;
}

.character-dps-badge {
  margin-left: 12px;
  padding: 4px 8px;
  background: rgba(0, 255, 255, 0.2);
  border: 1px solid var(--terminal-cyan);
  border-radius: 4px;
  color: var(--terminal-cyan);
  font-size: 12px;
  font-weight: bold;
}

.character-dps-overview {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(120px, 1fr));
  gap: 8px;
  margin-bottom: 16px;
}

.skill-dps-list {
  margin-top: 16px;
}

.skill-dps-table {
  margin-top: 8px;
}

.skill-dps-header,
.skill-dps-row {
  display: grid;
  grid-template-columns: 2fr 1fr 1fr 1fr 0.8fr 1fr;
  gap: 8px;
  padding: 8px;
  border-bottom: 1px solid var(--text-dim);
}

.skill-dps-header {
  background: rgba(0, 0, 0, 0.4);
  font-weight: bold;
  color: var(--terminal-cyan);
  font-size: 12px;
}

.skill-dps-row {
  transition: background 0.2s;
}

.skill-dps-row:hover {
  background: rgba(0, 255, 0, 0.05);
}

.skill-dps-col-name {
  color: var(--text-primary);
  font-weight: 500;
}

.skill-dps-col-dps {
  color: var(--terminal-gold);
  font-weight: bold;
}

.skill-dps-col-damage {
  color: var(--terminal-green);
}

.skill-dps-col-percent {
  color: var(--terminal-cyan);
}

.skill-dps-col-uses {
  color: var(--text-secondary);
  text-align: center;
}

.skill-dps-col-avg {
  color: var(--text-secondary);
}

.damage-composition {
  display: flex;
  flex-direction: column;
  gap: 12px;
  margin-top: 12px;
}

.damage-type-item {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.damage-type-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 12px;
}

.damage-type-name {
  color: var(--text-primary);
  font-weight: 500;
}

.damage-type-percent {
  color: var(--terminal-cyan);
  font-weight: bold;
}

.damage-type-bar {
  width: 100%;
  height: 20px;
  background: rgba(0, 0, 0, 0.3);
  border: 1px solid var(--text-dim);
  border-radius: 2px;
  overflow: hidden;
  position: relative;
}

.damage-type-fill {
  height: 100%;
  transition: width 0.3s ease;
  border-radius: 2px;
}

.damage-type-value {
  color: var(--text-secondary);
  font-size: 11px;
}

.character-damage-composition {
  margin-top: 16px;
  padding-top: 16px;
  border-top: 1px solid var(--text-dim);
}

.dps-empty-message {
  text-align: center;
  padding: 40px 20px;
}

.dps-empty-message p {
  margin: 8px 0;
  color: var(--text-secondary);
}

.dps-empty-hint {
  font-size: 12px;
  color: var(--text-gray);
  margin-top: 16px;
}

/* 统计会话控制 */
.stats-session-controls {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 20px;
  border-bottom: 1px solid var(--text-dim);
  background: rgba(0, 0, 0, 0.3);
}

.stats-session-status {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 13px;
}

.session-active {
  color: var(--terminal-green);
}

.session-inactive {
  color: var(--text-secondary);
}

.stats-session-actions {
  display: flex;
  gap: 8px;
}

.session-btn {
  padding: 6px 12px;
  border: 1px solid;
  background: transparent;
  cursor: pointer;
  transition: all 0.2s;
  font-family: inherit;
  font-size: 12px;
}

.session-btn-start {
  border-color: var(--terminal-green);
  color: var(--terminal-green);
}

.session-btn-start:hover {
  background: rgba(0, 255, 0, 0.1);
}

.session-btn-reset {
  border-color: var(--terminal-red);
  color: var(--terminal-red);
}

.session-btn-reset:hover {
  background: rgba(255, 0, 0, 0.1);
}

.battle-dps-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 20px;
  border-bottom: 1px solid var(--text-dim);
  background: rgba(0, 0, 0, 0.3);
  margin-bottom: 16px;
}

.battle-dps-title {
  margin: 0;
  font-size: 16px;
  color: var(--text-primary);
}

.battle-dps-close {
  background: transparent;
  border: 1px solid var(--text-dim);
  color: var(--text-secondary);
  width: 28px;
  height: 28px;
  border-radius: 4px;
  cursor: pointer;
  font-size: 18px;
  line-height: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s;
}

.battle-dps-close:hover {
  background: rgba(255, 0, 0, 0.1);
  border-color: var(--terminal-red);
  color: var(--terminal-red);
}

.session-time-badge {
  margin-left: 12px;
  padding: 2px 8px;
  background: rgba(0, 255, 255, 0.1);
  border: 1px solid var(--terminal-cyan);
  border-radius: 4px;
  color: var(--terminal-cyan);
  font-size: 11px;
  font-weight: normal;
}

.stats-tab:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

/* 滚动条样式 */
.stats-content::-webkit-scrollbar {
  width: 8px;
}

.stats-content::-webkit-scrollbar-track {
  background: rgba(0, 0, 0, 0.3);
}

.stats-content::-webkit-scrollbar-thumb {
  background: var(--terminal-gray);
  border-radius: 4px;
}

.stats-content::-webkit-scrollbar-thumb:hover {
  background: var(--terminal-green);
}
</style>









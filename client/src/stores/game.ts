import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { get, post } from '@/api/client'
import type { Character, BattleLog, Monster, Zone, BattleStatus as BattleStatusType } from '@/types/game'

// 战斗回合结果
interface BattleTickResult {
  character: Character
  enemy?: Monster
  logs: BattleLog[]
  isRunning: boolean
  sessionKills: number
  sessionGold: number
  sessionExp: number
  battleCount: number
}

export const useGameStore = defineStore('game', () => {
  // 状态
  const battleLogs = ref<BattleLog[]>([])
  const battleStatus = ref<BattleStatusType>({
    isRunning: false,
    battleCount: 0,
    totalKills: 0,
    totalExp: 0,
    totalGold: 0,
    team: []
  })
  const zones = ref<Zone[]>([])
  const currentZone = ref<Zone | null>(null)
  const currentEnemy = ref<Monster | null>(null)
  const isLoading = ref(false)
  const battleInterval = ref<number | null>(null)

  // 计算属性
  const isRunning = computed(() => battleStatus.value.isRunning)

  // ═══════════════════════════════════════════════════════════
  // API 调用
  // ═══════════════════════════════════════════════════════════

  // 获取战斗状态
  async function fetchBattleStatus() {
    const res = await get<BattleStatusType>('/battle/status')
    if (res.success && res.data) {
      battleStatus.value = res.data
      currentEnemy.value = res.data.currentMonster || null
    }
  }

  // 获取战斗日志
  async function fetchBattleLogs() {
    const res = await get<{ logs: BattleLog[] }>('/battle/logs')
    if (res.success && res.data) {
      battleLogs.value = res.data.logs || []
    }
  }

  // 获取区域列表
  async function fetchZones() {
    const res = await get<Zone[]>('/zones')
    if (res.success && res.data) {
      zones.value = res.data
    }
  }

  // 开始战斗
  async function startBattle() {
    isLoading.value = true
    try {
      const res = await post<{ isRunning: boolean }>('/battle/start')
      if (res.success && res.data) {
        battleStatus.value.isRunning = res.data.isRunning
        if (res.data.isRunning) {
          startBattleLoop()
        }
      }
    } finally {
      isLoading.value = false
    }
  }

  // 停止战斗
  async function stopBattle() {
    isLoading.value = true
    try {
      const res = await post<{ isRunning: boolean }>('/battle/stop')
      if (res.success) {
        battleStatus.value.isRunning = false
        stopBattleLoop()
      }
    } finally {
      isLoading.value = false
    }
  }

  // 切换战斗状态
  async function toggleBattle() {
    isLoading.value = true
    try {
      console.log('[ToggleBattle] Calling...')
      const res = await post<{ isRunning: boolean }>('/battle/toggle')
      console.log('[ToggleBattle] Response:', res)
      
      if (res.success && res.data) {
        battleStatus.value.isRunning = res.data.isRunning
        console.log('[ToggleBattle] isRunning:', res.data.isRunning)
        
        if (res.data.isRunning) {
          startBattleLoop()
          console.log('[ToggleBattle] Started battle loop')
        } else {
          stopBattleLoop()
          console.log('[ToggleBattle] Stopped battle loop')
        }
      } else {
        console.warn('[ToggleBattle] Failed or no data:', res)
      }
    } finally {
      isLoading.value = false
    }
  }

  // 执行战斗回合
  async function battleTick() {
    try {
      const res = await post<BattleTickResult>('/battle/tick')
      console.log('[BattleTick] Response:', res)
      
      if (res.success && res.data) {
        const result = res.data
        console.log('[BattleTick] Result:', result)
        
        // 更新状态
        battleStatus.value.isRunning = result.isRunning
        battleStatus.value.battleCount = result.battleCount
        battleStatus.value.totalKills = result.sessionKills
        battleStatus.value.totalExp = result.sessionExp
        battleStatus.value.totalGold = result.sessionGold
        
        // 更新当前敌人
        currentEnemy.value = result.enemy || null

        // 添加新日志
        if (result.logs && result.logs.length > 0) {
          console.log('[BattleTick] Adding logs:', result.logs.length)
          for (const log of result.logs) {
            addLog(log)
          }
        }

        // 如果战斗停止，停止循环
        if (!result.isRunning) {
          stopBattleLoop()
        }

        return result.character
      } else {
        console.warn('[BattleTick] No data or not successful:', res)
      }
    } catch (e) {
      console.error('Battle tick failed:', e)
      stopBattleLoop()
    }
    return null
  }

  // 切换区域
  async function changeZone(zoneId: string) {
    isLoading.value = true
    try {
      const res = await post<{ status: BattleStatusType; logs: BattleLog[] }>('/battle/zone', { zoneId })
      if (res.success && res.data) {
        if (res.data.status) {
          battleStatus.value = res.data.status
        }
        if (res.data.logs) {
          for (const log of res.data.logs) {
            addLog(log)
          }
        }
        // 更新当前区域
        currentZone.value = zones.value.find(z => z.id === zoneId) || null
        return true
      }
      return false
    } finally {
      isLoading.value = false
    }
  }

  // ═══════════════════════════════════════════════════════════
  // 战斗循环
  // ═══════════════════════════════════════════════════════════

  function startBattleLoop() {
    if (battleInterval.value) return
    battleInterval.value = window.setInterval(() => {
      if (battleStatus.value.isRunning) {
        battleTick()
      }
    }, 1500) // 每1.5秒一个回合
  }

  function stopBattleLoop() {
    if (battleInterval.value) {
      clearInterval(battleInterval.value)
      battleInterval.value = null
    }
  }

  // ═══════════════════════════════════════════════════════════
  // 日志管理
  // ═══════════════════════════════════════════════════════════

  function addLog(log: BattleLog) {
    battleLogs.value.push(log)
    // 保持日志数量
    if (battleLogs.value.length > 100) {
      battleLogs.value = battleLogs.value.slice(-100)
    }
  }

  function clearLogs() {
    battleLogs.value = []
  }

  // ═══════════════════════════════════════════════════════════
  // 初始化
  // ═══════════════════════════════════════════════════════════

  async function init() {
    await Promise.all([
      fetchZones(),
      fetchBattleStatus(),
      fetchBattleLogs()
    ])

    // 如果战斗正在运行，启动循环
    if (battleStatus.value.isRunning) {
      startBattleLoop()
    }
  }

  // 清理
  function cleanup() {
    stopBattleLoop()
  }

  return {
    // 状态
    battleLogs,
    battleStatus,
    zones,
    currentZone,
    currentEnemy,
    isLoading,
    // 计算属性
    isRunning,
    // 方法
    fetchBattleStatus,
    fetchBattleLogs,
    fetchZones,
    startBattle,
    stopBattle,
    toggleBattle,
    battleTick,
    changeZone,
    addLog,
    clearLogs,
    init,
    cleanup,
    startBattleLoop,
    stopBattleLoop
  }
})


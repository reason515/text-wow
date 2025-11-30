import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { Character, BattleLog, BattleStatus, Monster, Zone, BattleResult } from '../types/game'

const API_BASE = '/api'

export const useGameStore = defineStore('game', () => {
  // 状态
  const character = ref<Character | null>(null)
  const battleLogs = ref<BattleLog[]>([])
  const battleStatus = ref<BattleStatus>({
    is_running: false,
    current_zone: '',
    current_enemy: null,
    battle_count: 0,
    session_kills: 0,
    session_gold: 0,
    session_exp: 0
  })
  const zones = ref<Zone[]>([])
  const isLoading = ref(false)
  const battleInterval = ref<number | null>(null)

  // 计算属性
  const hasCharacter = computed(() => character.value !== null)
  const isRunning = computed(() => battleStatus.value.is_running)
  const currentEnemy = computed(() => battleStatus.value.current_enemy)

  // API 调用
  async function createCharacter(name: string, race: string, className: string) {
    isLoading.value = true
    try {
      const res = await fetch(`${API_BASE}/character`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ name, race, class: className })
      })
      character.value = await res.json()
      await fetchBattleLogs()
      await fetchZones()
    } finally {
      isLoading.value = false
    }
  }

  async function fetchCharacter() {
    try {
      const res = await fetch(`${API_BASE}/character`)
      if (res.ok) {
        character.value = await res.json()
      }
    } catch (e) {
      console.error('Failed to fetch character:', e)
    }
  }

  async function fetchBattleLogs() {
    try {
      const res = await fetch(`${API_BASE}/battle/logs`)
      const data = await res.json()
      battleLogs.value = data.logs || []
    } catch (e) {
      console.error('Failed to fetch logs:', e)
    }
  }

  async function fetchBattleStatus() {
    try {
      const res = await fetch(`${API_BASE}/battle/status`)
      battleStatus.value = await res.json()
    } catch (e) {
      console.error('Failed to fetch status:', e)
    }
  }

  async function fetchZones() {
    try {
      const res = await fetch(`${API_BASE}/zones`)
      const data = await res.json()
      zones.value = data.zones || []
    } catch (e) {
      console.error('Failed to fetch zones:', e)
    }
  }

  async function toggleBattle() {
    try {
      const res = await fetch(`${API_BASE}/battle/toggle`, { method: 'POST' })
      const data = await res.json()
      battleStatus.value.is_running = data.running

      if (data.running) {
        startBattleLoop()
      } else {
        stopBattleLoop()
      }
    } catch (e) {
      console.error('Failed to toggle battle:', e)
    }
  }

  async function battleTick() {
    try {
      const res = await fetch(`${API_BASE}/battle/tick`, { method: 'POST' })
      const result: BattleResult = await res.json()
      
      if (result) {
        character.value = result.character
        battleStatus.value = result.status
        
        // 添加新日志
        if (result.logs && result.logs.length > 0) {
          battleLogs.value.push(...result.logs)
          // 保持日志数量
          if (battleLogs.value.length > 100) {
            battleLogs.value = battleLogs.value.slice(-100)
          }

        }
      }
    } catch (e) {
      console.error('Battle tick failed:', e)
    }
  }

  async function changeZone(zoneId: string) {
    try {
      const res = await fetch(`${API_BASE}/zone/change`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ zone_id: zoneId })
      })
      if (res.ok) {
        await fetchBattleStatus()
        await fetchBattleLogs()
      }
    } catch (e) {
      console.error('Failed to change zone:', e)
    }
  }

  function startBattleLoop() {
    if (battleInterval.value) return
    battleInterval.value = window.setInterval(() => {
      if (battleStatus.value.is_running) {
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

  // 添加本地日志（不发送到服务器）
  function addLocalLog(type: BattleLog['type'], message: string, color: string = '#00FF00') {
    const log: BattleLog = {
      time: new Date().toLocaleTimeString('zh-CN', { hour12: false }),
      type,
      message,
      color
    }
    battleLogs.value.push(log)
  }

  return {
    // 状态
    character,
    battleLogs,
    battleStatus,
    zones,
    isLoading,
    // 计算属性
    hasCharacter,
    isRunning,
    currentEnemy,
    // 方法
    createCharacter,
    fetchCharacter,
    fetchBattleLogs,
    fetchBattleStatus,
    fetchZones,
    toggleBattle,
    battleTick,
    changeZone,
    addLocalLog,
    startBattleLoop,
    stopBattleLoop
  }
})








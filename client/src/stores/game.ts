import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { Character, BattleLog, BattleStatus, Monster, Zone, BattleResult } from '../types/game'
import { get, post } from '@/api/client'

export const useGameStore = defineStore('game', () => {
  // 状态
  const character = ref<Character | null>(null)
  const battleLogs = ref<BattleLog[]>([])
  // 扩展 BattleStatus 以支持两种命名格式
  const battleStatus = ref<any>({
    is_running: false,
    isRunning: false,
    current_zone: '',
    currentZoneId: '',
    current_enemy: null,
    currentMonster: null,
    current_enemies: null,
    currentEnemies: null,
    battle_count: 0,
    battleCount: 0,
    session_kills: 0,
    totalKills: 0,
    session_gold: 0,
    totalGold: 0,
    session_exp: 0,
    totalExp: 0
  })
  const zones = ref<Zone[]>([])
  const isLoading = ref(false)
  const battleInterval = ref<number | null>(null)

  // 计算属性
  const hasCharacter = computed(() => character.value !== null)
  const isRunning = computed(() => battleStatus.value?.is_running ?? false)
  const currentEnemy = computed(() => battleStatus.value?.current_enemy ?? null)
  const currentEnemies = computed(() => {
    const enemies = battleStatus.value?.current_enemies ?? battleStatus.value?.currentEnemies ?? null
    return enemies && Array.isArray(enemies) ? enemies : (enemies ? [enemies] : [])
  })

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
      const response = await get<{ logs: BattleLog[] }>('/battle/logs')
      if (response.success && response.data) {
        battleLogs.value = response.data.logs || []
      }
    } catch (e) {
      console.error('Failed to fetch logs:', e)
    }
  }

  async function fetchBattleStatus() {
    try {
      const response = await get<any>('/battle/status')
      if (response.success && response.data) {
        const data = response.data
        // 同步两种命名格式，确保所有字段都有默认值
        battleStatus.value = {
          is_running: data.isRunning ?? data.is_running ?? false,
          isRunning: data.isRunning ?? data.is_running ?? false,
          current_zone: data.currentZoneId ?? data.current_zone ?? '',
          currentZoneId: data.currentZoneId ?? data.current_zone ?? '',
          current_enemy: data.currentMonster ?? data.current_enemy ?? null,
          currentMonster: data.currentMonster ?? data.current_enemy ?? null,
          current_enemies: data.currentEnemies ?? data.current_enemies ?? null,
          currentEnemies: data.currentEnemies ?? data.current_enemies ?? null,
          battle_count: data.battleCount ?? data.battle_count ?? 0,
          battleCount: data.battleCount ?? data.battle_count ?? 0,
          session_kills: data.totalKills ?? data.session_kills ?? 0,
          totalKills: data.totalKills ?? data.session_kills ?? 0,
          session_gold: data.totalGold ?? data.session_gold ?? 0,
          totalGold: data.totalGold ?? data.session_gold ?? 0,
          session_exp: data.totalExp ?? data.session_exp ?? 0,
          totalExp: data.totalExp ?? data.session_exp ?? 0,
          ...data // 保留其他字段
        }
      } else {
        // 如果获取失败，确保 battleStatus.value 有默认值
        if (!battleStatus.value) {
          battleStatus.value = {
            is_running: false,
            isRunning: false,
            current_zone: '',
            currentZoneId: '',
            current_enemy: null,
            currentMonster: null,
            battle_count: 0,
            battleCount: 0,
            session_kills: 0,
            totalKills: 0,
            session_gold: 0,
            totalGold: 0,
            session_exp: 0,
            totalExp: 0
          }
        }
      }
    } catch (e) {
      console.error('Failed to fetch status:', e)
      // 出错时也确保有默认值
      if (!battleStatus.value) {
        battleStatus.value = {
          is_running: false,
          isRunning: false,
          current_zone: '',
          currentZoneId: '',
          current_enemy: null,
          currentMonster: null,
          battle_count: 0,
          battleCount: 0,
          session_kills: 0,
          totalKills: 0,
          session_gold: 0,
          totalGold: 0,
          session_exp: 0,
          totalExp: 0
        }
      }
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
      console.log('Toggle battle called')
      console.log('battleStatus.value:', battleStatus.value)
      
      // 确保 battleStatus.value 存在
      if (!battleStatus.value) {
        battleStatus.value = {
          is_running: false,
          isRunning: false,
          current_zone: '',
          currentZoneId: '',
          current_enemy: null,
          currentMonster: null,
          battle_count: 0,
          battleCount: 0,
          session_kills: 0,
          totalKills: 0,
          session_gold: 0,
          totalGold: 0,
          session_exp: 0,
          totalExp: 0
        }
      }
      
      const response = await post<{ isRunning?: boolean; running?: boolean }>('/battle/toggle')
      console.log('Toggle battle response:', response)
      
      if (response.success && response.data) {
        // 后端返回 isRunning，前端兼容两种格式
        const isRunning = response.data.isRunning ?? response.data.running ?? false
        console.log('Battle isRunning:', isRunning)
        
        // 确保 battleStatus.value 存在后再设置
        if (battleStatus.value) {
          battleStatus.value.is_running = isRunning
          battleStatus.value.isRunning = isRunning
        }

        if (isRunning) {
          console.log('Starting battle loop')
          startBattleLoop()
        } else {
          console.log('Stopping battle loop')
          stopBattleLoop()
        }
      } else {
        console.error('Toggle battle failed:', response.error)
        alert('开始战斗失败: ' + (response.error || '未知错误'))
      }
    } catch (e) {
      console.error('Failed to toggle battle:', e)
      alert('开始战斗失败: ' + (e instanceof Error ? e.message : '未知错误'))
    }
  }

  async function battleTick() {
    try {
      const response = await post<BattleResult>('/battle/tick')
      
      if (response.success && response.data) {
        const result = response.data
        character.value = result.character
        
        // 确保 battleStatus.value 存在
        if (!battleStatus.value) {
          battleStatus.value = {
            is_running: false,
            isRunning: false,
            current_zone: '',
            currentZoneId: '',
            current_enemy: null,
            currentMonster: null,
            battle_count: 0,
            battleCount: 0,
            session_kills: 0,
            totalKills: 0,
            session_gold: 0,
            totalGold: 0,
            session_exp: 0,
            totalExp: 0
          }
        }
        
        // 合并状态数据，而不是直接替换
        if (result.status) {
          const status = result.status as any
          battleStatus.value = {
            ...battleStatus.value,
            is_running: status.isRunning ?? status.is_running ?? battleStatus.value.is_running,
            isRunning: status.isRunning ?? status.is_running ?? battleStatus.value.isRunning,
            current_zone: status.currentZoneId ?? status.current_zone ?? battleStatus.value.current_zone,
            currentZoneId: status.currentZoneId ?? status.current_zone ?? battleStatus.value.currentZoneId,
            current_enemy: status.currentMonster ?? status.current_enemy ?? battleStatus.value.current_enemy,
            currentMonster: status.currentMonster ?? status.current_enemy ?? battleStatus.value.currentMonster,
            current_enemies: status.currentEnemies ?? status.current_enemies ?? battleStatus.value.current_enemies,
            currentEnemies: status.currentEnemies ?? status.current_enemies ?? battleStatus.value.currentEnemies,
            battle_count: status.battleCount ?? status.battle_count ?? battleStatus.value.battle_count,
            battleCount: status.battleCount ?? status.battle_count ?? battleStatus.value.battleCount,
            session_kills: status.totalKills ?? status.session_kills ?? battleStatus.value.session_kills,
            totalKills: status.totalKills ?? status.session_kills ?? battleStatus.value.totalKills,
            session_gold: status.totalGold ?? status.session_gold ?? battleStatus.value.session_gold,
            totalGold: status.totalGold ?? status.session_gold ?? battleStatus.value.totalGold,
            session_exp: status.totalExp ?? status.session_exp ?? battleStatus.value.session_exp,
            totalExp: status.totalExp ?? status.session_exp ?? battleStatus.value.totalExp,
            ...status // 保留其他字段
          }
        }
        
        // 如果 result 中有 enemies 字段，也更新
        if (result.enemies) {
          battleStatus.value.current_enemies = result.enemies
          battleStatus.value.currentEnemies = result.enemies
        }
        
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
    
    // 确保 battleStatus.value 存在
    if (!battleStatus.value) {
      battleStatus.value = {
        is_running: false,
        isRunning: false,
        current_zone: '',
        currentZoneId: '',
        current_enemy: null,
        currentMonster: null,
        battle_count: 0,
        battleCount: 0,
        session_kills: 0,
        totalKills: 0,
        session_gold: 0,
        totalGold: 0,
        session_exp: 0,
        totalExp: 0
      }
    }
    
    battleInterval.value = window.setInterval(() => {
      if (battleStatus.value?.is_running) {
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
    currentEnemies,
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








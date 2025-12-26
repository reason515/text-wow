import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { Character, BattleLog, BattleStatus, Monster, Zone, BattleResult } from '../types/game'
import { get, post } from '@/api/client'
import { useCharacterStore } from './character'

const API_BASE = '/api'

export const useGameStore = defineStore('game', () => {
  // 状态
  const character = ref<Character | null>(null)
  const battleLogs = ref<BattleLog[]>([])
  // 扩展 BattleStatus 以支持两种命名格式
  const battleStatus = ref<any>({
    is_running: false,
    isRunning: false,
    is_resting: false,
    isResting: false,
    restUntil: null,
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
  const explorations = ref<Record<string, { exploration: number; kills: number }>>({})
  const isLoading = ref(false)
  const battleInterval = ref<number | null>(null)
  const currentZone = ref<Zone | null>(null)
  const skillSelection = ref<any | null>(null)            // 待处理的技能选择机会
  const skillSelectionLoading = ref(false)
  const lastSelectionCheckLevel = ref<number | null>(null)

  // 计算属性
  const hasCharacter = computed(() => character.value !== null)
  const isRunning = computed(() => battleStatus.value?.is_running ?? false)
  const currentEnemy = computed(() => battleStatus.value?.current_enemy ?? null)
  const currentEnemies = computed(() => {
    const enemies = battleStatus.value?.current_enemies ?? battleStatus.value?.currentEnemies ?? null
    // 如果是数组，直接返回；如果是null/undefined，返回空数组；如果是单个对象，包装成数组
    if (enemies === null || enemies === undefined) {
      return []
    }
    if (Array.isArray(enemies)) {
      return enemies
    }
    // 如果是单个对象，包装成数组
    return [enemies]
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
      const response = await get<Character>('/character')
      if (response.success && response.data) {
        character.value = response.data
      } else {
        character.value = null
      }
    } catch (e) {
      console.error('Failed to fetch character:', e)
      character.value = null
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
          is_resting: data.isResting ?? data.is_resting ?? false,
          isResting: data.isResting ?? data.is_resting ?? false,
          restUntil: data.restUntil ?? data.rest_until ?? null,
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
        
        // 如果战斗状态中包含角色数据（Team），更新角色数据
        // Team 是一个数组，不是包含 characters 的对象
        if (data.team && Array.isArray(data.team) && data.team.length > 0) {
          // 同步更新 characterStore 中的角色列表，确保界面显示最新的HP/MP等数据
          const charStore = useCharacterStore()
          charStore.characters = data.team
          
          // 使用第一个角色更新当前显示角色
          character.value = data.team[0]
        } else if (data.character) {
          // 如果直接返回了角色数据，也更新
          character.value = data.character
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
      const response = await get<any>('/battle/zones')
      if (response.success && response.data) {
        zones.value = response.data.zones || response.data || []
        // 保存探索度数据
        if (response.data.explorations) {
          explorations.value = response.data.explorations
        }
        // 设置当前区域
        if (battleStatus.value?.currentZoneId) {
          currentZone.value = zones.value.find(z => z.id === battleStatus.value.currentZoneId) || null
        }
      }
    } catch (e) {
      console.error('Failed to fetch zones:', e)
    }
  }

  async function fetchExplorations() {
    try {
      const response = await get<Record<string, { exploration: number; kills: number }>>('/battle/explorations')
      if (response.success && response.data) {
        explorations.value = response.data
      }
    } catch (e) {
      console.error('Failed to fetch explorations:', e)
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
          // 先启动循环，然后立即执行一次 battleTick 以开始战斗
          startBattleLoop()
          // 立即执行一次战斗回合，不要等待定时器
          battleTick().catch(console.error)
        } else {
          console.log('Stopping battle loop')
          stopBattleLoop()
        }
      } else {
        console.error('Toggle battle failed:', response.error)
        if (typeof alert !== 'undefined') {
          alert('开始战斗失败: ' + (response.error || '未知错误'))
        }
      }
    } catch (e) {
      console.error('Failed to toggle battle:', e)
      if (typeof alert !== 'undefined') {
        alert('开始战斗失败: ' + (e instanceof Error ? e.message : '未知错误'))
      }
    }
  }

  async function battleTick(): Promise<Character | null> {
    try {
      const response = await post<BattleResult>('/battle/tick')
      
      if (response.success && response.data) {
        const result = response.data
        const prevLevel = character.value?.level ?? null
        character.value = result.character
        
        // 战斗后刷新战斗状态以获取最新的角色数据（包括所有角色的HP/MP等）
        // fetchBattleStatus 会更新 charStore.characters
        await fetchBattleStatus()
        
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
        
        // 更新战斗状态：后端直接返回 BattleTickResult，包含 IsRunning、IsResting 等字段
        // 同时也支持通过 status 对象传递（向后兼容）
        const status = (result as any).status || result
        
        // 优先处理 enemies 字段（后端返回的是小写 "enemies"）
        let enemiesToSet: any = null
        if ('enemies' in result && result.enemies !== undefined) {
          enemiesToSet = result.enemies
          console.log('[DEBUG] Found enemies in result.enemies:', enemiesToSet)
        } else if (status && 'enemies' in status && status.enemies !== undefined) {
          enemiesToSet = status.enemies
          console.log('[DEBUG] Found enemies in status.enemies:', enemiesToSet)
        } else if (status && 'currentEnemies' in status && status.currentEnemies !== undefined) {
          enemiesToSet = status.currentEnemies
          console.log('[DEBUG] Found enemies in status.currentEnemies:', enemiesToSet)
        } else if (status && 'current_enemies' in status && status.current_enemies !== undefined) {
          enemiesToSet = status.current_enemies
          console.log('[DEBUG] Found enemies in status.current_enemies:', enemiesToSet)
        } else {
          console.log('[DEBUG] No enemies found in result. result keys:', Object.keys(result || {}), 'status keys:', status ? Object.keys(status) : 'no status')
        }
        
        if (status) {
          battleStatus.value = {
            ...battleStatus.value,
            is_running: status.isRunning ?? status.is_running ?? battleStatus.value.is_running,
            isRunning: status.isRunning ?? status.is_running ?? battleStatus.value.isRunning,
            is_resting: status.isResting ?? status.is_resting ?? battleStatus.value.is_resting ?? false,
            isResting: status.isResting ?? status.is_resting ?? battleStatus.value.isResting ?? false,
            restUntil: status.restUntil ?? status.rest_until ?? battleStatus.value.restUntil,
            current_zone: status.currentZoneId ?? status.current_zone ?? battleStatus.value.current_zone,
            currentZoneId: status.currentZoneId ?? status.current_zone ?? battleStatus.value.currentZoneId,
            current_enemy: status.currentMonster ?? status.current_enemy ?? battleStatus.value.current_enemy,
            currentMonster: status.currentMonster ?? status.current_enemy ?? battleStatus.value.currentMonster,
            // 使用预先处理的 enemies 值，如果找到了就使用，否则保留旧值
            current_enemies: enemiesToSet !== undefined ? enemiesToSet : battleStatus.value.current_enemies,
            currentEnemies: enemiesToSet !== undefined ? enemiesToSet : battleStatus.value.currentEnemies,
            battle_count: status.battleCount ?? status.battle_count ?? battleStatus.value.battle_count,
            battleCount: status.battleCount ?? status.battle_count ?? battleStatus.value.battleCount,
            session_kills: status.sessionKills ?? status.totalKills ?? status.session_kills ?? battleStatus.value.session_kills,
            totalKills: status.sessionKills ?? status.totalKills ?? status.session_kills ?? battleStatus.value.totalKills,
            session_gold: status.sessionGold ?? status.totalGold ?? status.session_gold ?? battleStatus.value.session_gold,
            totalGold: status.sessionGold ?? status.totalGold ?? status.session_gold ?? battleStatus.value.totalGold,
            session_exp: status.sessionExp ?? status.totalExp ?? status.session_exp ?? battleStatus.value.session_exp,
            totalExp: status.sessionExp ?? status.totalExp ?? status.session_exp ?? battleStatus.value.totalExp,
            ...status // 保留其他字段
          }
        }
        
        // 如果找到了 enemies 值，确保更新（防止上面的逻辑没有正确设置）
        if (enemiesToSet !== undefined && enemiesToSet !== null) {
          battleStatus.value.current_enemies = enemiesToSet
          battleStatus.value.currentEnemies = enemiesToSet
        }

        // 检查是否出现技能选择机会（被动/主动）
        const newLevel = character.value?.level ?? null
        if (character.value) {
          // 若等级提升或当前等级为里程碑（3或5的倍数）且还没有待选机会，则尝试查询
          const milestone = newLevel !== null && (newLevel % 3 === 0 || newLevel % 5 === 0)
          if (skillSelection.value) {
            // 保留已有的待选机会，避免被覆盖
            lastSelectionCheckLevel.value = newLevel
          } else if (newLevel !== null && (newLevel !== prevLevel || milestone)) {
            await checkSkillSelection(true)
          }
        }
        
        // 添加新日志
        if (result.logs && result.logs.length > 0) {
          battleLogs.value.push(...result.logs)
          // 保持日志数量
          if (battleLogs.value.length > 100) {
            battleLogs.value = battleLogs.value.slice(-100)
          }
        }
        
        return character.value
      }
      return null
    } catch (e) {
      console.error('Battle tick failed:', e)
      return null
    }
  }

  // 获取技能选择机会（被动/主动）
  async function checkSkillSelection(force = false) {
    if (!character.value) return null

    const level = character.value.level ?? 0
    if (!force && skillSelection.value) {
      return skillSelection.value
    }
    if (!force && lastSelectionCheckLevel.value === level) {
      return skillSelection.value
    }

    skillSelectionLoading.value = true
    try {
      const response = await get<any>(`/characters/${character.value.id}/skills/selection`)
      if (response.success && response.data) {
        skillSelection.value = response.data
      } else {
        skillSelection.value = null
      }
      lastSelectionCheckLevel.value = level
      return skillSelection.value
    } catch (e) {
      console.error('Failed to check skill selection:', e)
      return null
    } finally {
      skillSelectionLoading.value = false
    }
  }

  // 提交技能选择
  async function submitSkillSelection(payload: { skillId?: string; passiveId?: string; isUpgrade: boolean }): Promise<boolean> {
    if (!character.value) return false
    try {
      const requestBody = {
        ...payload,
        characterId: character.value.id,
      }

      const response = await post(`/characters/${character.value.id}/skills/select`, requestBody)
      if (response.success) {
        // 选择成功后，清空待选机会并刷新角色数据
        skillSelection.value = null
        await fetchCharacter()
        return true
      } else {
        console.error('Select skill failed:', response.error)
        return false
      }
    } catch (e) {
      console.error('Select skill failed:', e)
      return false
    }
  }

  async function changeZone(zoneId: string): Promise<boolean> {
    try {
      const response = await post<any>('/battle/change-zone', { zoneId })
      if (response.success) {
        await fetchBattleStatus()
        await fetchBattleLogs()
        // 更新当前区域
        const zone = zones.value.find(z => z.id === zoneId)
        if (zone) {
          currentZone.value = zone
        }
        return true
      }
      return false
    } catch (e) {
      console.error('Failed to change zone:', e)
      return false
    }
  }

  function startBattleLoop() {
    // 如果循环已经存在，先停止它
    if (battleInterval.value) {
      stopBattleLoop()
    }
    
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
    
    console.log('Starting battle loop, isRunning:', battleStatus.value?.is_running)
    
    battleInterval.value = (typeof window !== 'undefined' ? window.setInterval : setInterval)(() => {
      // 如果战斗正在运行，或者正在休息，都需要调用 battleTick 来处理
      const isRunning = battleStatus.value?.is_running ?? false
      const isResting = battleStatus.value?.isResting ?? battleStatus.value?.is_resting ?? false
      
      console.log('Battle loop tick, isRunning:', isRunning, 'isResting:', isResting)
      
      if (isRunning || isResting) {
        battleTick()
      } else {
        // 即使战斗停止且不在休息，也要定期获取角色数据和战斗状态（用于显示死亡/复活状态）
        fetchCharacter().catch(console.error)
        fetchBattleStatus().catch(console.error)
      }
    }, 1500) // 每1.5秒一个回合
  }

  function stopBattleLoop() {
    if (battleInterval.value) {
      (typeof window !== 'undefined' ? window.clearInterval : clearInterval)(battleInterval.value)
      battleInterval.value = null
    }
  }

  // 添加本地日志（不发送到服务器）
  function addLocalLog(type: BattleLog['type'], message: string, color: string = '#00FF00') {
    const log: BattleLog = {
      time: new Date().toLocaleTimeString('zh-CN', { hour12: false }),
      logType: type || 'info',
      type,
      message,
      color
    }
    battleLogs.value.push(log)
  }

  // 添加日志（兼容测试）
  function addLog(log: BattleLog) {
    battleLogs.value.push(log)
    // 保持日志数量
    if (battleLogs.value.length > 100) {
      battleLogs.value = battleLogs.value.slice(-100)
    }
  }

  // 清除日志
  function clearLogs() {
    battleLogs.value = []
  }

  // 开始战斗
  async function startBattle(): Promise<boolean> {
    try {
      const response = await post<{ isRunning: boolean }>('/battle/start')
      if (response.success && response.data) {
        battleStatus.value.is_running = response.data.isRunning
        battleStatus.value.isRunning = response.data.isRunning
        if (response.data.isRunning) {
          startBattleLoop()
        }
        return true
      }
      return false
    } catch (e) {
      console.error('Failed to start battle:', e)
      return false
    }
  }

  // 停止战斗
  async function stopBattle(): Promise<boolean> {
    try {
      const response = await post<{ isRunning: boolean }>('/battle/stop')
      if (response.success && response.data) {
        battleStatus.value.is_running = response.data.isRunning
        battleStatus.value.isRunning = response.data.isRunning
        stopBattleLoop()
        return true
      }
      return false
    } catch (e) {
      console.error('Failed to stop battle:', e)
      return false
    }
  }

  // 初始化
  async function init() {
    isLoading.value = true
    try {
      await Promise.all([
        fetchCharacter(),
        fetchBattleStatus(),
        fetchBattleLogs(),
        fetchZones()
      ])
      
      // 如果战斗正在运行，启动战斗循环
      if (battleStatus.value?.is_running || battleStatus.value?.isRunning) {
        startBattleLoop()
      }
    } catch (e) {
      console.error('Failed to initialize:', e)
    } finally {
      isLoading.value = false
    }
  }

  // 清理
  function cleanup() {
    stopBattleLoop()
    battleLogs.value = []
    character.value = null
    battleStatus.value = {
      is_running: false,
      isRunning: false,
      is_resting: false,
      isResting: false,
      restUntil: null,
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
    }
  }

  return {
    // 状态
    character,
    battleLogs,
    battleStatus,
    zones,
    explorations,
    currentZone,
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
    fetchExplorations,
    toggleBattle,
    battleTick,
    changeZone,
    skillSelection,
    skillSelectionLoading,
    checkSkillSelection,
    submitSkillSelection,
    addLocalLog,
    addLog,
    clearLogs,
    startBattle,
    stopBattle,
    init,
    cleanup,
    startBattleLoop,
    stopBattleLoop
  }
})








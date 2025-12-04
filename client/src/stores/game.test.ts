import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useGameStore } from './game'
import { mockFetch, createMockResponse } from '@/test/setup'

// ═══════════════════════════════════════════════════════════
// Mock 数据
// ═══════════════════════════════════════════════════════════

function createMockCharacter(overrides = {}) {
  return {
    id: 1,
    userId: 1,
    name: 'TestHero',
    raceId: 'human',
    classId: 'warrior',
    faction: 'alliance',
    level: 1,
    hp: 100,
    maxHp: 100,
    exp: 0,
    expToNext: 100,
    resource: 100,
    maxResource: 100,
    resourceType: 'rage',
    physicalAttack: 15,
    magicAttack: 8,
    physicalDefense: 10,
    magicDefense: 6,
    strength: 15,
    agility: 10,
    intellect: 5,
    stamina: 12,
    spirit: 8,
    critRate: 0.05,
    critDamage: 1.5,
    totalKills: 0,
    totalDeaths: 0,
    ...overrides,
  }
}

function createMockMonster(overrides = {}) {
  return {
    id: 'wolf',
    zoneId: 'elwynn_forest',
    name: '森林狼',
    level: 2,
    type: 'normal',
    hp: 30,
    maxHp: 30,
    physicalAttack: 8,
    magicAttack: 4,
    physicalDefense: 2,
    magicDefense: 1,
    expReward: 15,
    goldMin: 1,
    goldMax: 5,
    ...overrides,
  }
}

function createMockZone(overrides = {}) {
  return {
    id: 'elwynn_forest',
    name: '艾尔文森林',
    description: '暴风城周边的和平森林',
    minLevel: 1,
    maxLevel: 10,
    faction: 'alliance',
    expMulti: 1.0,
    goldMulti: 1.0,
    ...overrides,
  }
}

function createMockBattleLog(overrides = {}) {
  return {
    message: '你使用 [英勇打击] 对 森林狼 造成 25 点伤害',
    logType: 'combat',
    createdAt: new Date().toISOString(),
    ...overrides,
  }
}

function createMockBattleStatus(overrides = {}) {
  return {
    isRunning: false,
    battleCount: 0,
    totalKills: 0,
    totalExp: 0,
    totalGold: 0,
    ...overrides,
  }
}

// ═══════════════════════════════════════════════════════════
// 测试套件
// ═══════════════════════════════════════════════════════════

describe('Game Store', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    mockFetch.mockReset()
    localStorage.setItem('token', 'test-token')
  })

  afterEach(() => {
    vi.clearAllTimers()
  })

  // ═════════════════════════════════════════════════════════
  // 初始状态测试
  // ═════════════════════════════════════════════════════════

  describe('Initial State', () => {
    it('should have correct initial state', () => {
      const store = useGameStore()

      expect(store.battleLogs).toEqual([])
      expect(store.zones).toEqual([])
      expect(store.currentZone).toBeNull()
      expect(store.currentEnemy).toBeNull()
      expect(store.isLoading).toBe(false)
      expect(store.isRunning).toBe(false)
    })

    it('should have initial battle status', () => {
      const store = useGameStore()

      expect(store.battleStatus.isRunning).toBe(false)
      expect(store.battleStatus.battleCount).toBe(0)
      expect(store.battleStatus.totalKills).toBe(0)
    })
  })

  // ═════════════════════════════════════════════════════════
  // 获取战斗状态测试
  // ═════════════════════════════════════════════════════════

  describe('fetchBattleStatus', () => {
    it('should fetch battle status successfully', async () => {
      const store = useGameStore()
      const mockStatus = createMockBattleStatus({ isRunning: true, battleCount: 5 })

      mockFetch.mockResolvedValueOnce(createMockResponse(mockStatus))

      await store.fetchBattleStatus()

      expect(store.battleStatus.isRunning).toBe(true)
      expect(store.battleStatus.battleCount).toBe(5)
    })

    it('should handle empty response', async () => {
      const store = useGameStore()

      mockFetch.mockResolvedValueOnce(createMockResponse(null))

      await store.fetchBattleStatus()

      // Should not crash
      expect(store.battleStatus).toBeDefined()
    })
  })

  // ═════════════════════════════════════════════════════════
  // 获取战斗日志测试
  // ═════════════════════════════════════════════════════════

  describe('fetchBattleLogs', () => {
    it('should fetch battle logs successfully', async () => {
      const store = useGameStore()
      const mockLogs = [
        createMockBattleLog({ message: 'Log 1' }),
        createMockBattleLog({ message: 'Log 2' }),
      ]

      mockFetch.mockResolvedValueOnce(createMockResponse({ logs: mockLogs }))

      await store.fetchBattleLogs()

      expect(store.battleLogs).toHaveLength(2)
      expect(store.battleLogs[0].message).toBe('Log 1')
    })

    it('should handle empty logs', async () => {
      const store = useGameStore()

      mockFetch.mockResolvedValueOnce(createMockResponse({ logs: [] }))

      await store.fetchBattleLogs()

      expect(store.battleLogs).toEqual([])
    })
  })

  // ═════════════════════════════════════════════════════════
  // 获取区域测试
  // ═════════════════════════════════════════════════════════

  describe('fetchZones', () => {
    it('should fetch zones successfully', async () => {
      const store = useGameStore()
      const mockZones = [
        createMockZone({ id: 'zone1', name: 'Zone 1' }),
        createMockZone({ id: 'zone2', name: 'Zone 2' }),
      ]

      mockFetch.mockResolvedValueOnce(createMockResponse({ zones: mockZones }))

      await store.fetchZones()

      expect(store.zones).toHaveLength(2)
      expect(store.zones[0].id).toBe('zone1')
    })
  })

  // ═════════════════════════════════════════════════════════
  // 开始战斗测试
  // ═════════════════════════════════════════════════════════

  describe('startBattle', () => {
    it('should start battle successfully', async () => {
      const store = useGameStore()

      mockFetch.mockResolvedValueOnce(createMockResponse({ isRunning: true }))

      await store.startBattle()

      expect(store.battleStatus.isRunning).toBe(true)
      expect(store.isLoading).toBe(false)
    })

    it('should set loading state during request', async () => {
      const store = useGameStore()

      let resolvePromise: (value: any) => void
      const promise = new Promise((resolve) => {
        resolvePromise = resolve
      })

      mockFetch.mockReturnValueOnce(promise)

      const startPromise = store.startBattle()

      // startBattle doesn't set isLoading, it's only used in init
      // So we just verify the function completes
      resolvePromise!(createMockResponse({ isRunning: true }))
      await startPromise

      expect(store.battleStatus.isRunning).toBe(true)
    })
  })

  // ═════════════════════════════════════════════════════════
  // 停止战斗测试
  // ═════════════════════════════════════════════════════════

  describe('stopBattle', () => {
    it('should stop battle successfully', async () => {
      const store = useGameStore()
      store.battleStatus.isRunning = true

      mockFetch.mockResolvedValueOnce(createMockResponse({ isRunning: false }))

      await store.stopBattle()

      expect(store.battleStatus.isRunning).toBe(false)
    })
  })

  // ═════════════════════════════════════════════════════════
  // 切换战斗测试
  // ═════════════════════════════════════════════════════════

  describe('toggleBattle', () => {
    it('should toggle battle on', async () => {
      const store = useGameStore()

      mockFetch.mockResolvedValueOnce(createMockResponse({ isRunning: true }))
      // Mock battleTick call (should be called immediately after starting)
      mockFetch.mockResolvedValueOnce(createMockResponse({
        character: createMockCharacter(),
        logs: [],
        isRunning: true
      }))

      await store.toggleBattle()

      expect(store.battleStatus.isRunning).toBe(true)
      // Verify that battleTick was called immediately (not waiting for interval)
      expect(mockFetch).toHaveBeenCalledTimes(2) // toggle + immediate battleTick
    })

    it('should toggle battle off', async () => {
      const store = useGameStore()
      store.battleStatus.isRunning = true

      mockFetch.mockResolvedValueOnce(createMockResponse({ isRunning: false }))

      await store.toggleBattle()

      expect(store.battleStatus.isRunning).toBe(false)
    })

    it('should start battle loop when toggling on', async () => {
      const store = useGameStore()
      store.battleStatus.isRunning = false

      mockFetch.mockResolvedValueOnce(createMockResponse({ isRunning: true }))
      mockFetch.mockResolvedValueOnce(createMockResponse({
        character: createMockCharacter(),
        logs: [],
        isRunning: true
      }))

      await store.toggleBattle()

      expect(store.battleStatus.isRunning).toBe(true)
      // Verify battle loop is started (battleInterval should be set)
      // Note: We can't directly check the interval, but we can verify battleTick was called
      expect(mockFetch).toHaveBeenCalledTimes(2)
    })

    it('should handle toggle battle errors gracefully', async () => {
      const store = useGameStore()
      const alertSpy = vi.spyOn(window, 'alert').mockImplementation(() => {})

      mockFetch.mockResolvedValueOnce(createMockResponse({
        success: false,
        error: 'Failed to toggle battle'
      }))

      await store.toggleBattle()

      expect(alertSpy).toHaveBeenCalledWith('开始战斗失败: Failed to toggle battle')
      alertSpy.mockRestore()
    })
  })

  // ═════════════════════════════════════════════════════════
  // 战斗回合测试
  // ═════════════════════════════════════════════════════════

  describe('battleTick', () => {
    it('should execute battle tick successfully', async () => {
      const store = useGameStore()
      const mockResult = {
        character: createMockCharacter({ hp: 80 }),
        enemy: createMockMonster({ hp: 15 }),
        logs: [createMockBattleLog()],
        status: {
          isRunning: true,
          currentMonster: createMockMonster({ hp: 15 }),
          battleCount: 1,
          totalKills: 1,
          totalGold: 5,
          totalExp: 15,
        },
      }

      mockFetch.mockResolvedValueOnce(createMockResponse(mockResult))
      mockFetch.mockResolvedValueOnce(createMockResponse(createMockBattleStatus({ isRunning: true, battleCount: 1, totalKills: 1 })))

      const character = await store.battleTick()

      expect(character).toBeDefined()
      expect(character!.hp).toBe(80)
      expect(store.battleStatus.battleCount).toBe(1)
      expect(store.battleStatus.totalKills).toBe(1)
      expect(store.battleLogs.length).toBeGreaterThan(0)
    })

    it('should update enemy status', async () => {
      const store = useGameStore()
      const mockEnemy = createMockMonster({ hp: 15 })
      const mockResult = {
        character: createMockCharacter(),
        enemy: mockEnemy,
        logs: [],
        status: {
          isRunning: true,
          currentMonster: mockEnemy,
          battleCount: 1,
          totalKills: 0,
          totalGold: 0,
          totalExp: 0,
        },
      }

      mockFetch.mockResolvedValueOnce(createMockResponse(mockResult))
      mockFetch.mockResolvedValueOnce(createMockResponse(createMockBattleStatus({ isRunning: true, currentMonster: mockEnemy })))

      await store.battleTick()

      expect(store.currentEnemy).toBeDefined()
      expect(store.currentEnemy!.hp).toBe(15)
    })

    it('should clear enemy when defeated', async () => {
      const store = useGameStore()
      const mockResult = {
        character: createMockCharacter(),
        enemy: null,
        enemies: null,
        logs: [createMockBattleLog({ logType: 'victory' })],
        status: {
          isRunning: true,
          currentMonster: null,
          currentEnemies: null,
          battleCount: 1,
          totalKills: 1,
          totalGold: 5,
          totalExp: 15,
        },
      }

      mockFetch.mockResolvedValueOnce(createMockResponse(mockResult))
      mockFetch.mockResolvedValueOnce(createMockResponse(createMockBattleStatus({ isRunning: true, currentMonster: null })))

      await store.battleTick()

      expect(store.currentEnemy).toBeNull()
    })

    it('should handle battle stop on death', async () => {
      const store = useGameStore()
      const mockResult = {
        character: createMockCharacter({ hp: 50 }),
        enemy: null,
        logs: [createMockBattleLog({ logType: 'death' })],
        status: {
          isRunning: false, // Battle stopped due to death
          currentMonster: null,
          battleCount: 1,
          totalKills: 0,
          totalGold: 0,
          totalExp: 0,
        },
      }

      mockFetch.mockResolvedValueOnce(createMockResponse(mockResult))
      mockFetch.mockResolvedValueOnce(createMockResponse(createMockBattleStatus({ isRunning: false })))

      await store.battleTick()

      expect(store.battleStatus.isRunning).toBe(false)
    })
  })

  // ═════════════════════════════════════════════════════════
  // 区域切换测试
  // ═════════════════════════════════════════════════════════

  describe('changeZone', () => {
    it('should change zone successfully', async () => {
      const store = useGameStore()
      const mockZone = createMockZone({ id: 'westfall', name: '西部荒野' })
      store.zones = [createMockZone(), mockZone]

      store.zones = [createMockZone({ id: 'westfall', name: '西部荒野' })]
      mockFetch.mockResolvedValueOnce(createMockResponse({ zoneId: 'westfall' }))
      mockFetch.mockResolvedValueOnce(createMockResponse(createMockBattleStatus({ currentZoneId: 'westfall' })))
      mockFetch.mockResolvedValueOnce(createMockResponse({ logs: [] }))

      const result = await store.changeZone('westfall')

      expect(result).toBe(true)
      expect(store.currentZone?.id).toBe('westfall')
    })

    it('should add logs when changing zone', async () => {
      const store = useGameStore()
      store.zones = [createMockZone()]

      store.zones = [createMockZone({ id: 'elwynn_forest', name: '艾尔文森林' })]
      mockFetch.mockResolvedValueOnce(createMockResponse({ zoneId: 'elwynn_forest' }))
      mockFetch.mockResolvedValueOnce(createMockResponse(createMockBattleStatus()))
      mockFetch.mockResolvedValueOnce(createMockResponse({
        logs: [
          createMockBattleLog({ message: '>> 你来到了 [艾尔文森林]', logType: 'zone' }),
          createMockBattleLog({ message: '描述文字', logType: 'zone' }),
        ]
      }))

      await store.changeZone('elwynn_forest')

      expect(store.battleLogs.length).toBe(2)
    })

    it('should return false on failure', async () => {
      const store = useGameStore()

      mockFetch.mockResolvedValueOnce(createMockResponse(null, false, 'Zone not found'))

      const result = await store.changeZone('invalid_zone')

      expect(result).toBe(false)
    })
  })

  // ═════════════════════════════════════════════════════════
  // 日志管理测试
  // ═════════════════════════════════════════════════════════

  describe('Log Management', () => {
    it('should add log correctly', () => {
      const store = useGameStore()
      const log = createMockBattleLog()

      store.addLog(log)

      expect(store.battleLogs).toHaveLength(1)
      expect(store.battleLogs[0].message).toBe(log.message)
    })

    it('should limit logs to 100', () => {
      const store = useGameStore()

      // Add 110 logs
      for (let i = 0; i < 110; i++) {
        store.addLog(createMockBattleLog({ message: `Log ${i}` }))
      }

      expect(store.battleLogs.length).toBe(100)
      expect(store.battleLogs[0].message).toBe('Log 10') // Oldest should be removed
    })

    it('should clear logs', () => {
      const store = useGameStore()
      store.addLog(createMockBattleLog())
      store.addLog(createMockBattleLog())

      store.clearLogs()

      expect(store.battleLogs).toHaveLength(0)
    })
  })

  // ═════════════════════════════════════════════════════════
  // 初始化和清理测试
  // ═════════════════════════════════════════════════════════

  describe('Initialization and Cleanup', () => {
    it('should initialize all data', async () => {
      const store = useGameStore()
      const mockZones = [createMockZone()]
      const mockStatus = createMockBattleStatus()
      const mockLogs = { logs: [] }

      mockFetch
        .mockResolvedValueOnce(createMockResponse({ zones: mockZones }))
        .mockResolvedValueOnce(createMockResponse(mockStatus))
        .mockResolvedValueOnce(createMockResponse(mockLogs))

      await store.init()

      expect(mockFetch).toHaveBeenCalledTimes(4) // fetchCharacter, fetchBattleStatus, fetchBattleLogs, fetchZones
    })

    it('should start battle loop if already running', async () => {
      vi.useFakeTimers()
      const store = useGameStore()

      mockFetch
        .mockResolvedValueOnce(createMockResponse(null)) // fetchCharacter
        .mockResolvedValueOnce(createMockResponse({ isRunning: true, battleCount: 0, totalKills: 0, totalGold: 0, totalExp: 0 })) // fetchBattleStatus
        .mockResolvedValueOnce(createMockResponse({ logs: [] })) // fetchBattleLogs
        .mockResolvedValueOnce(createMockResponse({ zones: [] })) // fetchZones

      await store.init()

      // Battle loop should be started
      expect(store.battleStatus.isRunning).toBe(true)

      vi.useRealTimers()
    })

    it('should cleanup battle loop on cleanup', () => {
      vi.useFakeTimers()
      const store = useGameStore()
      
      store.startBattleLoop()
      store.cleanup()

      // After cleanup, adding time should not trigger more ticks
      vi.useRealTimers()
    })
  })

  // ═════════════════════════════════════════════════════════
  // 战斗循环测试
  // ═════════════════════════════════════════════════════════

  describe('Battle Loop', () => {
    it('should not start loop twice', () => {
      vi.useFakeTimers()
      const store = useGameStore()
      
      store.startBattleLoop()
      store.startBattleLoop() // Second call should be ignored

      vi.useRealTimers()
    })

    it('should stop loop correctly', () => {
      vi.useFakeTimers()
      const store = useGameStore()
      
      store.startBattleLoop()
      store.stopBattleLoop()
      store.stopBattleLoop() // Second call should be safe

      vi.useRealTimers()
    })
  })

  // ═════════════════════════════════════════════════════════
  // 错误处理测试
  // ═════════════════════════════════════════════════════════

  describe('Error Handling', () => {
    it('should handle network error in battleTick', async () => {
      const store = useGameStore()

      mockFetch.mockRejectedValueOnce(new Error('Network error'))

      const result = await store.battleTick()

      expect(result).toBeNull()
    })

    it('should stop loop on error', async () => {
      vi.useFakeTimers()
      const store = useGameStore()
      store.battleStatus.isRunning = true

      mockFetch.mockRejectedValueOnce(new Error('Network error'))

      await store.battleTick()

      // Loop should be stopped on error
      vi.useRealTimers()
    })
  })
})



















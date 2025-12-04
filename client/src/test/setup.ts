import { beforeEach, vi } from 'vitest'
import { config } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'

// ═══════════════════════════════════════════════════════════
// 全局测试配置
// ═══════════════════════════════════════════════════════════

// 每个测试前重置 Pinia
beforeEach(() => {
  setActivePinia(createPinia())
})

// ═══════════════════════════════════════════════════════════
// Mock localStorage
// ═══════════════════════════════════════════════════════════

const localStorageMock = (() => {
  let store: Record<string, string> = {}
  
  return {
    getItem: vi.fn((key: string) => store[key] || null),
    setItem: vi.fn((key: string, value: string) => {
      store[key] = value
    }),
    removeItem: vi.fn((key: string) => {
      delete store[key]
    }),
    clear: vi.fn(() => {
      store = {}
    }),
    get length() {
      return Object.keys(store).length
    },
    key: vi.fn((index: number) => {
      return Object.keys(store)[index] || null
    }),
  }
})()

Object.defineProperty(global, 'localStorage', {
  value: localStorageMock,
})

// 每个测试前清除 localStorage
beforeEach(() => {
  localStorageMock.clear()
  vi.clearAllMocks()
})

// ═══════════════════════════════════════════════════════════
// Mock fetch
// ═══════════════════════════════════════════════════════════

export const mockFetch = vi.fn()
global.fetch = mockFetch

beforeEach(() => {
  mockFetch.mockReset()
})

// ═══════════════════════════════════════════════════════════
// 辅助函数
// ═══════════════════════════════════════════════════════════

export function createMockResponse<T>(data: T, success = true, error?: string) {
  return Promise.resolve({
    json: () => Promise.resolve({
      success,
      data,
      error,
    }),
    ok: success,
    status: success ? 200 : 400,
  })
}

export function createMockUser(overrides = {}) {
  return {
    id: 1,
    username: 'testuser',
    email: 'test@example.com',
    maxTeamSize: 3,
    unlockedSlots: 3,
    gold: 0,
    currentZoneId: 'elwynn_forest',
    totalKills: 0,
    totalGoldGained: 0,
    playTime: 0,
    createdAt: new Date().toISOString(),
    lastLoginAt: null,
    ...overrides,
  }
}

export function createMockAuthResponse(user = createMockUser()) {
  return {
    token: 'mock-jwt-token-12345',
    user,
  }
}



















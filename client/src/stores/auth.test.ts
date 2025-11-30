import { describe, it, expect, vi, beforeEach } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useAuthStore } from './auth'
import { mockFetch, createMockResponse, createMockUser, createMockAuthResponse } from '@/test/setup'

describe('Auth Store', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    mockFetch.mockReset()
    vi.mocked(localStorage.getItem).mockReturnValue(null)
  })

  // ═══════════════════════════════════════════════════════════
  // 初始状态测试
  // ═══════════════════════════════════════════════════════════

  describe('Initial State', () => {
    it('should have null user initially', () => {
      const store = useAuthStore()
      expect(store.user).toBeNull()
    })

    it('should not be authenticated initially', () => {
      const store = useAuthStore()
      expect(store.isAuthenticated).toBe(false)
    })

    it('should not be loading initially', () => {
      const store = useAuthStore()
      expect(store.loading).toBe(false)
    })

    it('should have no error initially', () => {
      const store = useAuthStore()
      expect(store.error).toBeNull()
    })

    it('should have empty username when not logged in', () => {
      const store = useAuthStore()
      expect(store.username).toBe('')
    })
  })

  // ═══════════════════════════════════════════════════════════
  // 注册测试
  // ═══════════════════════════════════════════════════════════

  describe('Register', () => {
    it('should register successfully', async () => {
      const mockUser = createMockUser({ username: 'newuser' })
      const mockAuth = createMockAuthResponse(mockUser)
      mockFetch.mockResolvedValue(createMockResponse(mockAuth))

      const store = useAuthStore()
      const result = await store.register({
        username: 'newuser',
        password: 'password123',
        email: 'new@example.com',
      })

      expect(result).toBe(true)
      expect(store.user).toEqual(mockUser)
      expect(store.isAuthenticated).toBe(true)
      expect(store.error).toBeNull()
    })

    it('should set loading state during registration', async () => {
      let resolvePromise: Function
      const promise = new Promise((resolve) => {
        resolvePromise = resolve
      })
      mockFetch.mockReturnValue(promise as any)

      const store = useAuthStore()
      const registerPromise = store.register({
        username: 'test',
        password: 'password',
      })

      expect(store.loading).toBe(true)

      resolvePromise!(createMockResponse(createMockAuthResponse()))
      await registerPromise

      expect(store.loading).toBe(false)
    })

    it('should handle registration failure', async () => {
      mockFetch.mockResolvedValue(createMockResponse(null, false, '用户名已存在'))

      const store = useAuthStore()
      const result = await store.register({
        username: 'existing',
        password: 'password',
      })

      expect(result).toBe(false)
      expect(store.user).toBeNull()
      expect(store.error).toBe('用户名已存在')
    })

    it('should handle network error during registration', async () => {
      mockFetch.mockRejectedValue(new Error('Network failure'))

      const store = useAuthStore()
      const result = await store.register({
        username: 'test',
        password: 'password',
      })

      expect(result).toBe(false)
      // 错误消息来自 api/client.ts 的网络错误处理
      expect(store.error).toMatch(/网络错误|请求异常/)
    })

    it('should save token on successful registration', async () => {
      const mockAuth = createMockAuthResponse()
      mockFetch.mockResolvedValue(createMockResponse(mockAuth))

      const store = useAuthStore()
      await store.register({
        username: 'test',
        password: 'password',
      })

      expect(localStorage.setItem).toHaveBeenCalledWith('token', mockAuth.token)
    })
  })

  // ═══════════════════════════════════════════════════════════
  // 登录测试
  // ═══════════════════════════════════════════════════════════

  describe('Login', () => {
    it('should login successfully', async () => {
      const mockUser = createMockUser({ username: 'loginuser' })
      const mockAuth = createMockAuthResponse(mockUser)
      mockFetch.mockResolvedValue(createMockResponse(mockAuth))

      const store = useAuthStore()
      const result = await store.login({
        username: 'loginuser',
        password: 'password123',
      })

      expect(result).toBe(true)
      expect(store.user).toEqual(mockUser)
      expect(store.username).toBe('loginuser')
      expect(store.isAuthenticated).toBe(true)
    })

    it('should handle invalid credentials', async () => {
      mockFetch.mockResolvedValue(createMockResponse(null, false, '用户名或密码错误'))

      const store = useAuthStore()
      const result = await store.login({
        username: 'wrong',
        password: 'wrong',
      })

      expect(result).toBe(false)
      expect(store.error).toBe('用户名或密码错误')
    })

    it('should clear previous error on new login attempt', async () => {
      const store = useAuthStore()
      store.error = '之前的错误'

      mockFetch.mockResolvedValue(createMockResponse(createMockAuthResponse()))

      await store.login({
        username: 'test',
        password: 'password',
      })

      expect(store.error).toBeNull()
    })

    it('should save token on successful login', async () => {
      const mockAuth = createMockAuthResponse()
      mockFetch.mockResolvedValue(createMockResponse(mockAuth))

      const store = useAuthStore()
      await store.login({
        username: 'test',
        password: 'password',
      })

      expect(localStorage.setItem).toHaveBeenCalledWith('token', mockAuth.token)
    })
  })

  // ═══════════════════════════════════════════════════════════
  // 登出测试
  // ═══════════════════════════════════════════════════════════

  describe('Logout', () => {
    it('should clear user state on logout', async () => {
      // First login
      const mockAuth = createMockAuthResponse()
      mockFetch.mockResolvedValue(createMockResponse(mockAuth))

      const store = useAuthStore()
      await store.login({ username: 'test', password: 'pass' })
      expect(store.isAuthenticated).toBe(true)

      // Then logout
      store.logout()

      expect(store.user).toBeNull()
      expect(store.isAuthenticated).toBe(false)
    })

    it('should clear token on logout', () => {
      const store = useAuthStore()
      store.logout()

      expect(localStorage.removeItem).toHaveBeenCalledWith('token')
    })
  })

  // ═══════════════════════════════════════════════════════════
  // 初始化测试
  // ═══════════════════════════════════════════════════════════

  describe('Init', () => {
    it('should fetch current user if token exists', async () => {
      vi.mocked(localStorage.getItem).mockReturnValue('valid-token')
      const mockUser = createMockUser()
      mockFetch.mockResolvedValue(createMockResponse(mockUser))

      const store = useAuthStore()
      await store.init()

      expect(mockFetch).toHaveBeenCalled()
      expect(store.user).toEqual(mockUser)
    })

    it('should not fetch if no token', async () => {
      vi.mocked(localStorage.getItem).mockReturnValue(null)

      const store = useAuthStore()
      await store.init()

      expect(mockFetch).not.toHaveBeenCalled()
    })

    it('should clear invalid token', async () => {
      vi.mocked(localStorage.getItem).mockReturnValue('invalid-token')
      mockFetch.mockResolvedValue(createMockResponse(null, false, 'Invalid token'))

      const store = useAuthStore()
      await store.init()

      expect(localStorage.removeItem).toHaveBeenCalledWith('token')
      expect(store.user).toBeNull()
    })
  })

  // ═══════════════════════════════════════════════════════════
  // fetchCurrentUser 测试
  // ═══════════════════════════════════════════════════════════

  describe('fetchCurrentUser', () => {
    it('should update user on success', async () => {
      const mockUser = createMockUser({ username: 'fetched' })
      mockFetch.mockResolvedValue(createMockResponse(mockUser))

      const store = useAuthStore()
      const result = await store.fetchCurrentUser()

      expect(result).toBe(true)
      expect(store.user).toEqual(mockUser)
    })

    it('should clear token on failure', async () => {
      mockFetch.mockResolvedValue(createMockResponse(null, false, 'Unauthorized'))

      const store = useAuthStore()
      const result = await store.fetchCurrentUser()

      expect(result).toBe(false)
      expect(store.user).toBeNull()
      expect(localStorage.removeItem).toHaveBeenCalledWith('token')
    })

    it('should handle network error', async () => {
      mockFetch.mockRejectedValue(new Error('Network error'))

      const store = useAuthStore()
      const result = await store.fetchCurrentUser()

      expect(result).toBe(false)
      expect(store.user).toBeNull()
    })
  })

  // ═══════════════════════════════════════════════════════════
  // 并发操作测试
  // ═══════════════════════════════════════════════════════════

  describe('Concurrent Operations', () => {
    it('should handle multiple login attempts', async () => {
      const mockAuth = createMockAuthResponse()
      mockFetch.mockResolvedValue(createMockResponse(mockAuth))

      const store = useAuthStore()
      
      // 同时发起多个登录请求
      const promises = [
        store.login({ username: 'user1', password: 'pass' }),
        store.login({ username: 'user2', password: 'pass' }),
      ]

      const results = await Promise.all(promises)
      
      // 两个请求都应该成功
      expect(results).toEqual([true, true])
    })
  })
})


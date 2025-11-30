import { describe, it, expect, vi, beforeEach } from 'vitest'
import { setToken, clearToken, isLoggedIn, get, post, put, del } from './client'
import { mockFetch, createMockResponse } from '@/test/setup'

// ═══════════════════════════════════════════════════════════
// Token 管理测试
// ═══════════════════════════════════════════════════════════

describe('Token Management', () => {
  describe('setToken', () => {
    it('should store token in localStorage', () => {
      setToken('test-token-123')
      expect(localStorage.setItem).toHaveBeenCalledWith('token', 'test-token-123')
    })
  })

  describe('clearToken', () => {
    it('should remove token from localStorage', () => {
      clearToken()
      expect(localStorage.removeItem).toHaveBeenCalledWith('token')
    })
  })

  describe('isLoggedIn', () => {
    it('should return true when token exists', () => {
      vi.mocked(localStorage.getItem).mockReturnValue('valid-token')
      expect(isLoggedIn()).toBe(true)
    })

    it('should return false when token is null', () => {
      vi.mocked(localStorage.getItem).mockReturnValue(null)
      expect(isLoggedIn()).toBe(false)
    })

    it('should return false when token is empty string', () => {
      vi.mocked(localStorage.getItem).mockReturnValue('')
      expect(isLoggedIn()).toBe(false)
    })
  })
})

// ═══════════════════════════════════════════════════════════
// API 请求测试
// ═══════════════════════════════════════════════════════════

describe('API Requests', () => {
  beforeEach(() => {
    mockFetch.mockReset()
    vi.mocked(localStorage.getItem).mockReturnValue(null)
  })

  describe('GET request', () => {
    it('should make GET request to correct endpoint', async () => {
      mockFetch.mockResolvedValue(createMockResponse({ data: 'test' }))

      await get('/test-endpoint')

      expect(mockFetch).toHaveBeenCalledWith(
        '/api/test-endpoint',
        expect.objectContaining({
          method: 'GET',
          headers: expect.objectContaining({
            'Content-Type': 'application/json',
          }),
        })
      )
    })

    it('should include auth header when token exists', async () => {
      vi.mocked(localStorage.getItem).mockReturnValue('my-token')
      mockFetch.mockResolvedValue(createMockResponse({ data: 'test' }))

      await get('/protected')

      expect(mockFetch).toHaveBeenCalledWith(
        '/api/protected',
        expect.objectContaining({
          headers: expect.objectContaining({
            Authorization: 'Bearer my-token',
          }),
        })
      )
    })

    it('should return parsed response data', async () => {
      const mockData = { id: 1, name: 'Test' }
      mockFetch.mockResolvedValue(createMockResponse(mockData))

      const response = await get<typeof mockData>('/data')

      expect(response.success).toBe(true)
      expect(response.data).toEqual(mockData)
    })

    it('should handle network errors gracefully', async () => {
      mockFetch.mockRejectedValue(new Error('Network error'))

      const response = await get('/failing-endpoint')

      expect(response.success).toBe(false)
      expect(response.error).toContain('网络错误')
    })
  })

  describe('POST request', () => {
    it('should make POST request with body', async () => {
      mockFetch.mockResolvedValue(createMockResponse({ success: true }))
      const body = { username: 'test', password: 'pass' }

      await post('/auth/login', body)

      expect(mockFetch).toHaveBeenCalledWith(
        '/api/auth/login',
        expect.objectContaining({
          method: 'POST',
          body: JSON.stringify(body),
        })
      )
    })

    it('should handle POST without body', async () => {
      mockFetch.mockResolvedValue(createMockResponse({ success: true }))

      await post('/simple')

      expect(mockFetch).toHaveBeenCalledWith(
        '/api/simple',
        expect.objectContaining({
          method: 'POST',
          body: undefined,
        })
      )
    })
  })

  describe('PUT request', () => {
    it('should make PUT request with body', async () => {
      mockFetch.mockResolvedValue(createMockResponse({ success: true }))
      const body = { name: 'Updated' }

      await put('/resource/1', body)

      expect(mockFetch).toHaveBeenCalledWith(
        '/api/resource/1',
        expect.objectContaining({
          method: 'PUT',
          body: JSON.stringify(body),
        })
      )
    })
  })

  describe('DELETE request', () => {
    it('should make DELETE request', async () => {
      mockFetch.mockResolvedValue(createMockResponse({ success: true }))

      await del('/resource/1')

      expect(mockFetch).toHaveBeenCalledWith(
        '/api/resource/1',
        expect.objectContaining({
          method: 'DELETE',
        })
      )
    })
  })
})

// ═══════════════════════════════════════════════════════════
// 错误处理测试
// ═══════════════════════════════════════════════════════════

describe('Error Handling', () => {
  it('should handle JSON parse errors', async () => {
    mockFetch.mockResolvedValue({
      json: () => Promise.reject(new Error('Invalid JSON')),
    })

    const response = await get('/bad-json')

    expect(response.success).toBe(false)
    expect(response.error).toContain('网络错误')
  })

  it('should handle server error responses', async () => {
    mockFetch.mockResolvedValue(createMockResponse(null, false, 'Server error'))

    const response = await get('/error-endpoint')

    expect(response.success).toBe(false)
    expect(response.error).toBe('Server error')
  })

  it('should handle timeout errors', async () => {
    mockFetch.mockRejectedValue(new Error('Request timeout'))

    const response = await post('/slow-endpoint', {})

    expect(response.success).toBe(false)
    expect(response.error).toContain('网络错误')
    expect(response.error).toContain('timeout')
  })
})


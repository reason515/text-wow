import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { User, AuthResponse, LoginRequest, RegisterRequest } from '@/types/game'
import { post, get, setToken, clearToken, isLoggedIn } from '@/api/client'

export const useAuthStore = defineStore('auth', () => {
  // 状态
  const user = ref<User | null>(null)
  const loading = ref(false)
  const error = ref<string | null>(null)

  // 计算属性
  const isAuthenticated = computed(() => !!user.value)
  const username = computed(() => user.value?.username || '')

  // 初始化 - 检查本地存储的登录状态
  async function init() {
    if (isLoggedIn()) {
      await fetchCurrentUser()
    }
  }

  // 注册
  async function register(data: RegisterRequest): Promise<boolean> {
    loading.value = true
    error.value = null

    try {
      const response = await post<AuthResponse>('/auth/register', data)
      
      if (response.success && response.data) {
        setToken(response.data.token)
        user.value = response.data.user
        return true
      } else {
        error.value = response.error || '注册失败'
        return false
      }
    } catch (e) {
      const errorMsg = e instanceof Error ? e.message : '未知错误'
      error.value = `请求异常: ${errorMsg}`
      return false
    } finally {
      loading.value = false
    }
  }

  // 登录
  async function login(data: LoginRequest): Promise<boolean> {
    loading.value = true
    error.value = null

    try {
      const response = await post<AuthResponse>('/auth/login', data)
      
      if (response.success && response.data) {
        setToken(response.data.token)
        user.value = response.data.user
        return true
      } else {
        error.value = response.error || '登录失败'
        return false
      }
    } catch (e) {
      const errorMsg = e instanceof Error ? e.message : '未知错误'
      error.value = `请求异常: ${errorMsg}`
      return false
    } finally {
      loading.value = false
    }
  }

  // 登出
  function logout() {
    clearToken()
    user.value = null
  }

  // 获取当前用户信息
  async function fetchCurrentUser(): Promise<boolean> {
    loading.value = true

    try {
      const response = await get<User>('/user')
      
      if (response.success && response.data) {
        user.value = response.data
        return true
      } else {
        // Token可能过期
        clearToken()
        user.value = null
        return false
      }
    } catch (e) {
      clearToken()
      user.value = null
      return false
    } finally {
      loading.value = false
    }
  }

  return {
    // 状态
    user,
    loading,
    error,
    // 计算属性
    isAuthenticated,
    username,
    // 方法
    init,
    register,
    login,
    logout,
    fetchCurrentUser,
  }
})


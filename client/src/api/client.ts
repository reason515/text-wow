import type { APIResponse } from '@/types/game'

const API_BASE = '/api'

// 获取存储的token
function getToken(): string | null {
  return localStorage.getItem('token')
}

// 设置token
export function setToken(token: string): void {
  localStorage.setItem('token', token)
}

// 清除token
export function clearToken(): void {
  localStorage.removeItem('token')
}

// 检查是否已登录
export function isLoggedIn(): boolean {
  return !!getToken()
}

// 通用请求函数
async function request<T>(
  endpoint: string,
  options: RequestInit = {}
): Promise<APIResponse<T>> {
  const token = getToken()
  
  const headers: HeadersInit = {
    'Content-Type': 'application/json',
    ...options.headers,
  }
  
  if (token) {
    headers['Authorization'] = `Bearer ${token}`
  }

  try {
    const response = await fetch(`${API_BASE}${endpoint}`, {
      ...options,
      headers,
    })

    const data = await response.json()
    return data
  } catch (error) {
    console.error('API Error:', error)
    const errorMessage = error instanceof Error ? error.message : 'Unknown error'
    return {
      success: false,
      error: `网络错误: ${errorMessage}`,
    }
  }
}

// GET 请求
export async function get<T>(endpoint: string): Promise<APIResponse<T>> {
  return request<T>(endpoint, { method: 'GET' })
}

// POST 请求
export async function post<T>(endpoint: string, body?: any): Promise<APIResponse<T>> {
  return request<T>(endpoint, {
    method: 'POST',
    body: body ? JSON.stringify(body) : undefined,
  })
}

// PUT 请求
export async function put<T>(endpoint: string, body?: any): Promise<APIResponse<T>> {
  return request<T>(endpoint, {
    method: 'PUT',
    body: body ? JSON.stringify(body) : undefined,
  })
}

// DELETE 请求
export async function del<T>(endpoint: string): Promise<APIResponse<T>> {
  return request<T>(endpoint, { method: 'DELETE' })
}


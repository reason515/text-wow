import type { APIResponse } from '@/types/game'

const API_BASE = '/api'

// 获取存储的token
export function getToken(): string | null {
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

    // 检查响应状态
    if (!response.ok) {
      // 尝试解析错误响应
      let errorMessage = `HTTP ${response.status}: ${response.statusText}`
      try {
        const errorData = await response.json()
        if (errorData.error) {
          errorMessage = errorData.error
        }
      } catch {
        // 如果无法解析JSON，使用状态文本
      }
      return {
        success: false,
        error: errorMessage,
      }
    }

    // 检查响应内容类型
    const contentType = response.headers.get('content-type')
    if (!contentType || !contentType.includes('application/json')) {
      const text = await response.text()
      if (!text) {
        return {
          success: false,
          error: '服务器返回空响应',
        }
      }
      return {
        success: false,
        error: `服务器返回非JSON响应: ${text.substring(0, 100)}`,
      }
    }

    // 检查响应体是否为空
    const text = await response.text()
    if (!text || text.trim() === '') {
      return {
        success: false,
        error: '服务器返回空响应',
      }
    }

    try {
      const data = JSON.parse(text)
      return data
    } catch (parseError) {
      console.error('JSON Parse Error:', parseError, 'Response text:', text)
      return {
        success: false,
        error: `网络错误: JSON解析错误 - ${parseError instanceof Error ? parseError.message : 'Unknown error'}`,
      }
    }
  } catch (error) {
    console.error('API Error:', error)
    const errorMessage = error instanceof Error ? error.message : 'Unknown error'
    
    // 检查是否是网络错误
    if (errorMessage.includes('Failed to fetch') || errorMessage.includes('NetworkError')) {
      return {
        success: false,
        error: '无法连接到服务器，请检查后端服务是否已启动 (http://localhost:8080)',
      }
    }
    
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


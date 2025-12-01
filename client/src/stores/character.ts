import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { Character, Team, Race, Class, CharacterCreate } from '@/types/game'
import { get, post, put } from '@/api/client'

export const useCharacterStore = defineStore('character', () => {
  // 状态
  const characters = ref<Character[]>([])
  const team = ref<Team | null>(null)
  const races = ref<Race[]>([])
  const classes = ref<Class[]>([])
  const loading = ref(false)
  const error = ref<string | null>(null)

  // 计算属性
  const hasCharacters = computed(() => characters.value.length > 0)
  const activeCharacters = computed(() => 
    characters.value.filter(c => c.isActive)
  )
  const allianceRaces = computed(() => 
    races.value.filter(r => r.faction === 'alliance')
  )
  const hordeRaces = computed(() => 
    races.value.filter(r => r.faction === 'horde')
  )

  // 获取种族列表
  async function fetchRaces(): Promise<boolean> {
    try {
      const response = await get<Race[]>('/races')
      if (response.success && response.data) {
        races.value = response.data
        return true
      }
      return false
    } catch (e) {
      return false
    }
  }

  // 获取职业列表
  async function fetchClasses(): Promise<boolean> {
    try {
      const response = await get<Class[]>('/classes')
      if (response.success && response.data) {
        classes.value = response.data
        return true
      }
      return false
    } catch (e) {
      return false
    }
  }

  // 获取角色列表
  async function fetchCharacters(): Promise<boolean> {
    loading.value = true
    error.value = null

    try {
      const response = await get<Character[]>('/characters')
      if (response.success && response.data) {
        characters.value = response.data
        return true
      }
      error.value = response.error || 'Failed to fetch characters'
      return false
    } catch (e) {
      error.value = 'Network error'
      return false
    } finally {
      loading.value = false
    }
  }

  // 获取小队信息
  async function fetchTeam(): Promise<boolean> {
    loading.value = true

    try {
      const response = await get<Team>('/team')
      if (response.success && response.data) {
        team.value = response.data
        return true
      }
      return false
    } catch (e) {
      return false
    } finally {
      loading.value = false
    }
  }

  // 创建角色
  async function createCharacter(data: CharacterCreate): Promise<Character | null> {
    loading.value = true
    error.value = null

    try {
      const response = await post<Character>('/characters', data)
      if (response.success && response.data) {
        characters.value.push(response.data)
        return response.data
      }
      error.value = response.error || 'Failed to create character'
      return null
    } catch (e) {
      error.value = 'Network error'
      return null
    } finally {
      loading.value = false
    }
  }

  // 设置角色激活状态
  async function setCharacterActive(characterId: number, active: boolean): Promise<boolean> {
    try {
      const response = await put('/characters/active', { characterId, active })
      if (response.success) {
        const char = characters.value.find(c => c.id === characterId)
        if (char) {
          char.isActive = active
        }
        return true
      }
      return false
    } catch (e) {
      return false
    }
  }

  // 初始化游戏数据
  async function init() {
    await Promise.all([
      fetchRaces(),
      fetchClasses(),
    ])
  }

  // 清空数据（登出时调用）
  function clear() {
    characters.value = []
    team.value = null
  }

  return {
    // 状态
    characters,
    team,
    races,
    classes,
    loading,
    error,
    // 计算属性
    hasCharacters,
    activeCharacters,
    allianceRaces,
    hordeRaces,
    // 方法
    fetchRaces,
    fetchClasses,
    fetchCharacters,
    fetchTeam,
    createCharacter,
    setCharacterActive,
    init,
    clear,
  }
})







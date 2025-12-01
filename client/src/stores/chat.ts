import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { get, post } from '@/api/client'

export interface ChatMessage {
  id: number
  channel: string
  faction?: string
  zoneId?: string
  senderId: number
  senderName: string
  senderClass?: string
  receiverId?: number
  content: string
  createdAt: string
}

export interface OnlineUser {
  userId: number
  characterId?: number
  characterName: string
  faction: string
  zoneId?: string
  isOnline: boolean
}

export const useChatStore = defineStore('chat', () => {
  // 状态
  const messages = ref<ChatMessage[]>([])
  const onlineUsers = ref<OnlineUser[]>([])
  const currentChannel = ref('world')
  const loading = ref(false)
  const error = ref<string | null>(null)
  const onlineCount = ref({ alliance: 0, horde: 0 })
  const lastWhisperFrom = ref<string | null>(null)

  // 计算属性
  const filteredMessages = computed(() => {
    if (currentChannel.value === 'all') {
      return messages.value
    }
    return messages.value.filter(m => m.channel === currentChannel.value)
  })

  // 频道颜色
  const channelColors: Record<string, string> = {
    world: '#ffd700',    // 金色
    zone: '#ff8c00',     // 橙色
    trade: '#32cd32',    // 绿色
    lfg: '#4169e1',      // 蓝色
    whisper: '#da70d6',  // 紫色
    system: '#ffffff',   // 白色
  }

  // 获取消息
  async function fetchMessages(channel: string = 'recent', zoneId?: string) {
    loading.value = true
    error.value = null

    try {
      let url = `/chat/messages?channel=${channel}`
      if (zoneId) {
        url += `&zoneId=${zoneId}`
      }

      const response = await get<ChatMessage[]>(url)
      if (response.success && response.data) {
        messages.value = response.data
      } else {
        error.value = response.error || 'Failed to fetch messages'
      }
    } catch (e) {
      error.value = 'Network error'
    } finally {
      loading.value = false
    }
  }

  // 发送消息
  async function sendMessage(content: string, channel: string = 'world', receiver?: string) {
    error.value = null

    try {
      const response = await post<ChatMessage>('/chat/send', {
        channel,
        content,
        receiver,
      })

      if (response.success && response.data) {
        messages.value.push(response.data)
        return true
      } else {
        error.value = response.error || 'Failed to send message'
        return false
      }
    } catch (e) {
      error.value = 'Network error'
      return false
    }
  }

  // 获取在线用户
  async function fetchOnlineUsers() {
    try {
      const response = await get<{
        users: OnlineUser[]
        allianceCount: number
        hordeCount: number
      }>('/chat/online')

      if (response.success && response.data) {
        onlineUsers.value = response.data.users
        onlineCount.value = {
          alliance: response.data.allianceCount,
          horde: response.data.hordeCount,
        }
      }
    } catch (e) {
      // 忽略错误
    }
  }

  // 屏蔽玩家
  async function blockPlayer(playerName: string) {
    try {
      const response = await post('/chat/block', { playerName })
      return response.success
    } catch (e) {
      return false
    }
  }

  // 取消屏蔽
  async function unblockPlayer(playerName: string) {
    try {
      const response = await post('/chat/unblock', { playerName })
      return response.success
    } catch (e) {
      return false
    }
  }

  // 设置在线状态
  async function setOnline(characterId: number, zoneId?: string) {
    try {
      await post('/chat/online', { characterId, zoneId })
    } catch (e) {
      // 忽略错误
    }
  }

  // 设置离线状态
  async function setOffline() {
    try {
      await post('/chat/offline', {})
    } catch (e) {
      // 忽略错误
    }
  }

  // 心跳
  async function heartbeat() {
    try {
      await post('/chat/heartbeat', {})
    } catch (e) {
      // 忽略错误
    }
  }

  // 切换频道
  function setChannel(channel: string) {
    currentChannel.value = channel
  }

  // 添加消息 (用于实时推送)
  function addMessage(msg: ChatMessage) {
    messages.value.push(msg)
    
    // 如果是私聊，记录发送者
    if (msg.channel === 'whisper') {
      lastWhisperFrom.value = msg.senderName
    }

    // 限制消息数量
    if (messages.value.length > 200) {
      messages.value = messages.value.slice(-100)
    }
  }

  // 解析命令
  function parseCommand(input: string): { command: string; args: string[] } | null {
    if (!input.startsWith('/')) {
      return null
    }

    const parts = input.slice(1).split(' ')
    const command = parts[0].toLowerCase()
    const args = parts.slice(1)

    return { command, args }
  }

  // 处理输入
  async function handleInput(input: string): Promise<boolean> {
    const cmd = parseCommand(input)

    if (!cmd) {
      // 普通消息，发送到当前频道
      return sendMessage(input, currentChannel.value)
    }

    switch (cmd.command) {
      case 's':
      case 'say':
      case 'world':
        return sendMessage(cmd.args.join(' '), 'world')

      case 'z':
      case 'zone':
        return sendMessage(cmd.args.join(' '), 'zone')

      case 't':
      case 'trade':
        return sendMessage(cmd.args.join(' '), 'trade')

      case 'lfg':
        return sendMessage(cmd.args.join(' '), 'lfg')

      case 'w':
      case 'whisper':
        if (cmd.args.length < 2) {
          error.value = 'Usage: /w <player> <message>'
          return false
        }
        return sendMessage(cmd.args.slice(1).join(' '), 'whisper', cmd.args[0])

      case 'r':
      case 'reply':
        if (!lastWhisperFrom.value) {
          error.value = 'No one to reply to'
          return false
        }
        return sendMessage(cmd.args.join(' '), 'whisper', lastWhisperFrom.value)

      case 'block':
        if (cmd.args.length < 1) {
          error.value = 'Usage: /block <player>'
          return false
        }
        return blockPlayer(cmd.args[0])

      case 'unblock':
        if (cmd.args.length < 1) {
          error.value = 'Usage: /unblock <player>'
          return false
        }
        return unblockPlayer(cmd.args[0])

      default:
        error.value = `Unknown command: /${cmd.command}`
        return false
    }
  }

  // 获取频道颜色
  function getChannelColor(channel: string): string {
    return channelColors[channel] || '#888888'
  }

  // 清空
  function clear() {
    messages.value = []
    onlineUsers.value = []
    currentChannel.value = 'world'
  }

  return {
    // 状态
    messages,
    onlineUsers,
    currentChannel,
    loading,
    error,
    onlineCount,
    lastWhisperFrom,
    // 计算属性
    filteredMessages,
    // 方法
    fetchMessages,
    sendMessage,
    fetchOnlineUsers,
    blockPlayer,
    unblockPlayer,
    setOnline,
    setOffline,
    heartbeat,
    setChannel,
    addMessage,
    handleInput,
    getChannelColor,
    clear,
  }
})





<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted, nextTick, watch } from 'vue'
import { useChatStore } from '@/stores/chat'
import { useCharacterStore } from '@/stores/character'

const chatStore = useChatStore()
const charStore = useCharacterStore()

// 状态
const inputMessage = ref('')
const chatContainer = ref<HTMLElement | null>(null)
const isMinimized = ref(false)
const showOnline = ref(false)

// 频道配置
const channels = [
  { id: 'world', name: '世界', shortcut: '/s' },
  { id: 'zone', name: '区域', shortcut: '/z' },
  { id: 'trade', name: '交易', shortcut: '/t' },
  { id: 'lfg', name: '组队', shortcut: '/lfg' },
]

// 自动滚动
watch(() => chatStore.messages.length, async () => {
  await nextTick()
  scrollToBottom()
})

function scrollToBottom() {
  if (chatContainer.value) {
    chatContainer.value.scrollTop = chatContainer.value.scrollHeight
  }
}

// 发送消息
async function handleSend() {
  if (!inputMessage.value.trim()) return

  const success = await chatStore.handleInput(inputMessage.value)
  if (success) {
    inputMessage.value = ''
  }
}

// 处理按键
function handleKeydown(e: KeyboardEvent) {
  if (e.key === 'Enter' && !e.shiftKey) {
    e.preventDefault()
    handleSend()
  }
}

// 格式化时间
function formatTime(dateStr: string): string {
  const date = new Date(dateStr)
  return date.toLocaleTimeString('zh-CN', { hour: '2-digit', minute: '2-digit' })
}

// 获取频道标签
function getChannelTag(channel: string): string {
  const tags: Record<string, string> = {
    world: '世界',
    zone: '区域',
    trade: '交易',
    lfg: '组队',
    whisper: '私聊',
    system: '系统',
  }
  return tags[channel] || channel
}

// 获取职业颜色
function getClassColor(classId?: string): string {
  const colors: Record<string, string> = {
    warrior: '#c79c6e',
    paladin: '#f58cba',
    hunter: '#abd473',
    rogue: '#fff569',
    priest: '#ffffff',
    mage: '#69ccf0',
    warlock: '#9482c9',
    druid: '#ff7d0a',
    shaman: '#0070de',
  }
  return colors[classId || ''] || '#cccccc'
}

// 心跳定时器
let heartbeatTimer: number | null = null
let onlineTimer: number | null = null

onMounted(async () => {
  // 加载最近消息
  await chatStore.fetchMessages('recent')
  
  // 设置在线状态
  if (charStore.characters.length > 0) {
    const char = charStore.characters[0]
    await chatStore.setOnline(char.id)
  }

  // 获取在线用户
  await chatStore.fetchOnlineUsers()

  // 启动心跳
  heartbeatTimer = window.setInterval(() => {
    chatStore.heartbeat()
  }, 60000) // 每分钟

  // 定期刷新在线用户
  onlineTimer = window.setInterval(() => {
    chatStore.fetchOnlineUsers()
  }, 30000) // 每30秒

  scrollToBottom()
})

onUnmounted(() => {
  // 设置离线
  chatStore.setOffline()

  // 清理定时器
  if (heartbeatTimer) {
    clearInterval(heartbeatTimer)
  }
  if (onlineTimer) {
    clearInterval(onlineTimer)
  }
})
</script>

<template>
  <div class="chat-panel" :class="{ minimized: isMinimized }">
    <!-- 标题栏 -->
    <div class="chat-header">
      <span class="chat-title">💬 聊天</span>
      <div class="chat-controls">
        <span class="online-count" @click="showOnline = !showOnline">
          👥 {{ chatStore.onlineCount.alliance + chatStore.onlineCount.horde }}
        </span>
        <button class="minimize-btn" @click="isMinimized = !isMinimized">
          {{ isMinimized ? '▲' : '▼' }}
        </button>
      </div>
    </div>

    <template v-if="!isMinimized">
      <!-- 频道选择 -->
      <div class="channel-tabs">
        <button 
          v-for="ch in channels" 
          :key="ch.id"
          class="channel-tab"
          :class="{ active: chatStore.currentChannel === ch.id }"
          :style="{ '--channel-color': chatStore.getChannelColor(ch.id) }"
          @click="chatStore.setChannel(ch.id)"
        >
          {{ ch.name }}
        </button>
      </div>

      <!-- 消息列表 -->
      <div class="chat-messages" ref="chatContainer">
        <div 
          v-for="msg in chatStore.filteredMessages" 
          :key="msg.id"
          class="chat-message"
        >
          <span class="msg-time">[{{ formatTime(msg.createdAt) }}]</span>
          <span 
            class="msg-channel"
            :style="{ color: chatStore.getChannelColor(msg.channel) }"
          >
            [{{ getChannelTag(msg.channel) }}]
          </span>
          <span 
            class="msg-sender"
            :style="{ color: getClassColor(msg.senderClass) }"
          >
            {{ msg.senderName }}:
          </span>
          <span class="msg-content">{{ msg.content }}</span>
        </div>

        <div v-if="chatStore.messages.length === 0" class="no-messages">
          暂无消息，来说点什么吧！
        </div>
      </div>

      <!-- 在线用户列表 -->
      <div v-if="showOnline" class="online-panel">
        <div class="online-header">
          在线玩家 ({{ chatStore.onlineUsers.length }})
          <button class="close-btn" @click="showOnline = false">×</button>
        </div>
        <div class="online-list">
          <div 
            v-for="user in chatStore.onlineUsers" 
            :key="user.userId"
            class="online-user"
          >
            {{ user.characterName }}
          </div>
          <div v-if="chatStore.onlineUsers.length === 0" class="no-online">
            暂无同阵营在线玩家
          </div>
        </div>
      </div>

      <!-- 输入框 -->
      <div class="chat-input-container">
        <input 
          v-model="inputMessage"
          class="chat-input"
          type="text"
          placeholder="输入消息... (/w 玩家名 私聊)"
          maxlength="200"
          @keydown="handleKeydown"
        />
        <button class="send-btn" @click="handleSend" :disabled="!inputMessage.trim()">
          发送
        </button>
      </div>

      <!-- 错误提示 -->
      <div v-if="chatStore.error" class="chat-error">
        {{ chatStore.error }}
      </div>
    </template>
  </div>
</template>

<style scoped>
.chat-panel {
  display: flex;
  flex-direction: column;
  border: 1px solid var(--terminal-green);
  background: rgba(0, 20, 0, 0.9);
  width: 100%;
  max-height: 300px;
  font-size: 12px;
}

.chat-panel.minimized {
  max-height: auto;
}

.chat-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 8px 12px;
  border-bottom: 1px solid var(--terminal-gray);
  background: rgba(0, 50, 0, 0.5);
}

.chat-title {
  color: var(--terminal-green);
  font-weight: bold;
}

.chat-controls {
  display: flex;
  align-items: center;
  gap: 10px;
}

.online-count {
  color: var(--terminal-cyan);
  font-size: 11px;
  cursor: pointer;
}

.online-count:hover {
  color: var(--terminal-green);
}

.minimize-btn {
  background: none;
  border: none;
  color: var(--terminal-gray);
  cursor: pointer;
  font-size: 10px;
}

.minimize-btn:hover {
  color: var(--terminal-green);
}

/* 频道选项卡 */
.channel-tabs {
  display: flex;
  padding: 5px;
  gap: 5px;
  border-bottom: 1px solid var(--terminal-gray);
}

.channel-tab {
  background: transparent;
  border: 1px solid var(--terminal-gray);
  color: var(--terminal-gray);
  padding: 3px 8px;
  font-family: inherit;
  font-size: 11px;
  cursor: pointer;
  transition: all 0.2s;
}

.channel-tab:hover {
  border-color: var(--channel-color, var(--terminal-green));
  color: var(--channel-color, var(--terminal-green));
}

.channel-tab.active {
  border-color: var(--channel-color, var(--terminal-green));
  color: var(--channel-color, var(--terminal-green));
  background: rgba(255, 255, 255, 0.05);
}

/* 消息列表 */
.chat-messages {
  flex: 1;
  overflow-y: auto;
  padding: 8px;
  min-height: 120px;
  max-height: 180px;
}

.chat-message {
  margin-bottom: 4px;
  line-height: 1.4;
  word-break: break-word;
}

.msg-time {
  color: var(--terminal-gray);
  font-size: 10px;
}

.msg-channel {
  margin-left: 4px;
  font-size: 10px;
}

.msg-sender {
  margin-left: 4px;
  font-weight: bold;
}

.msg-content {
  color: var(--terminal-green);
  margin-left: 4px;
}

.no-messages {
  color: var(--terminal-gray);
  text-align: center;
  padding: 20px;
}

/* 在线用户面板 */
.online-panel {
  position: absolute;
  right: 0;
  top: 100%;
  width: 200px;
  background: rgba(0, 20, 0, 0.95);
  border: 1px solid var(--terminal-green);
  z-index: 100;
}

.online-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 8px;
  border-bottom: 1px solid var(--terminal-gray);
  color: var(--terminal-cyan);
  font-size: 11px;
}

.close-btn {
  background: none;
  border: none;
  color: var(--terminal-gray);
  cursor: pointer;
  font-size: 14px;
}

.online-list {
  max-height: 200px;
  overflow-y: auto;
  padding: 5px;
}

.online-user {
  padding: 3px 8px;
  color: var(--terminal-green);
  font-size: 11px;
}

.online-user:hover {
  background: rgba(0, 255, 0, 0.1);
}

.no-online {
  color: var(--terminal-gray);
  text-align: center;
  padding: 10px;
  font-size: 11px;
}

/* 输入框 */
.chat-input-container {
  display: flex;
  padding: 8px;
  gap: 8px;
  border-top: 1px solid var(--terminal-gray);
}

.chat-input {
  flex: 1;
  background: rgba(0, 0, 0, 0.5);
  border: 1px solid var(--terminal-gray);
  color: var(--terminal-green);
  padding: 6px 10px;
  font-family: inherit;
  font-size: 12px;
}

.chat-input:focus {
  outline: none;
  border-color: var(--terminal-green);
}

.chat-input::placeholder {
  color: var(--terminal-gray);
}

.send-btn {
  background: transparent;
  border: 1px solid var(--terminal-green);
  color: var(--terminal-green);
  padding: 6px 15px;
  font-family: inherit;
  font-size: 11px;
  cursor: pointer;
  transition: all 0.2s;
}

.send-btn:hover:not(:disabled) {
  background: var(--terminal-green);
  color: var(--terminal-bg);
}

.send-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

/* 错误提示 */
.chat-error {
  padding: 5px 10px;
  background: rgba(255, 0, 0, 0.1);
  color: var(--terminal-red);
  font-size: 11px;
  text-align: center;
}

/* 滚动条 */
.chat-messages::-webkit-scrollbar,
.online-list::-webkit-scrollbar {
  width: 6px;
}

.chat-messages::-webkit-scrollbar-track,
.online-list::-webkit-scrollbar-track {
  background: transparent;
}

.chat-messages::-webkit-scrollbar-thumb,
.online-list::-webkit-scrollbar-thumb {
  background: var(--terminal-gray);
}

.chat-messages::-webkit-scrollbar-thumb:hover,
.online-list::-webkit-scrollbar-thumb:hover {
  background: var(--terminal-green);
}
</style>







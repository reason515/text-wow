<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { useCharacterStore } from '@/stores/character'
import AuthScreen from '@/components/AuthScreen.vue'
import CreateCharacter from '@/components/CreateCharacter.vue'
import GameScreen from '@/components/GameScreen.vue'

// 页面状态
type PageState = 'loading' | 'auth' | 'create-character' | 'game'

const authStore = useAuthStore()
const charStore = useCharacterStore()

const currentPage = ref<PageState>('loading')
const initialized = ref(false)

// 初始化
onMounted(async () => {
  // 初始化认证状态
  await authStore.init()
  
  // 加载游戏配置
  await charStore.init()

  if (authStore.isAuthenticated) {
    // 已登录，加载角色
    await charStore.fetchCharacters()
    
    if (charStore.hasCharacters) {
      currentPage.value = 'game'
    } else {
      currentPage.value = 'create-character'
    }
  } else {
    currentPage.value = 'auth'
  }
  
  initialized.value = true
})

// 事件处理
function onAuthSuccess() {
  // 登录成功后检查是否有角色
  charStore.fetchCharacters().then(() => {
    if (charStore.hasCharacters) {
      currentPage.value = 'game'
    } else {
      currentPage.value = 'create-character'
    }
  })
}

function onCharacterCreated() {
  currentPage.value = 'game'
}

function onLogout() {
  authStore.logout()
  charStore.clear()
  currentPage.value = 'auth'
}

function goToCreateCharacter() {
  currentPage.value = 'create-character'
}

function backToGame() {
  if (charStore.hasCharacters) {
    currentPage.value = 'game'
  }
}
</script>

<template>
  <div class="app">
    <!-- 加载中 -->
    <div v-if="currentPage === 'loading'" class="loading-screen">
      <div class="loading-text">正在连接艾泽拉斯...</div>
      <div class="loading-bar">
        <div class="loading-progress"></div>
      </div>
    </div>

    <!-- 登录/注册 -->
    <AuthScreen 
      v-else-if="currentPage === 'auth'" 
      @success="onAuthSuccess"
    />

    <!-- 创建角色 -->
    <CreateCharacter 
      v-else-if="currentPage === 'create-character'"
      @created="onCharacterCreated"
      @back="backToGame"
    />

    <!-- 游戏主界面 -->
    <GameScreen 
      v-else-if="currentPage === 'game'"
      @logout="onLogout"
      @create-character="goToCreateCharacter"
    />

    <!-- CRT 效果 -->
    <div class="crt-overlay"></div>
    <div class="scanlines"></div>
  </div>
</template>

<style>
/* 终端样式已在 main.ts 中全局导入 */

.app {
  min-height: 100vh;
  position: relative;
}

/* 加载界面 */
.loading-screen {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 100vh;
}

.loading-text {
  color: var(--terminal-green);
  font-size: 14px;
  margin-bottom: 20px;
  animation: blink 1s infinite;
}

.loading-bar {
  width: 300px;
  height: 4px;
  background: var(--terminal-gray);
  overflow: hidden;
}

.loading-progress {
  width: 0%;
  height: 100%;
  background: var(--terminal-green);
  animation: loading 2s ease-in-out infinite;
}

@keyframes loading {
  0% { width: 0%; }
  50% { width: 100%; }
  100% { width: 0%; }
}

@keyframes blink {
  0%, 50% { opacity: 1; }
  51%, 100% { opacity: 0.5; }
}

/* CRT效果 */
.crt-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  pointer-events: none;
  background: radial-gradient(ellipse at center, transparent 0%, rgba(0,0,0,0.2) 100%);
  z-index: 9998;
}

.scanlines {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  pointer-events: none;
  background: repeating-linear-gradient(
    0deg,
    rgba(0, 0, 0, 0.1),
    rgba(0, 0, 0, 0.1) 1px,
    transparent 1px,
    transparent 2px
  );
  z-index: 9999;
}
</style>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useAuthStore } from '@/stores/auth'

const emit = defineEmits<{
  success: []
}>()

const authStore = useAuthStore()

// 状态
const mode = ref<'login' | 'register'>('login')
const username = ref('')
const password = ref('')
const confirmPassword = ref('')
const email = ref('')

// 计算属性
const isLogin = computed(() => mode.value === 'login')
const title = computed(() => isLogin.value ? '登录' : '注册')
const submitText = computed(() => isLogin.value ? '进入艾泽拉斯' : '创建账号')
const switchText = computed(() => isLogin.value ? '没有账号？点击注册' : '已有账号？点击登录')

const canSubmit = computed(() => {
  if (!username.value || !password.value) return false
  if (!isLogin.value && password.value !== confirmPassword.value) return false
  if (!isLogin.value && password.value.length < 6) return false
  return true
})

const passwordError = computed(() => {
  if (!isLogin.value && confirmPassword.value && password.value !== confirmPassword.value) {
    return '两次输入的密码不一致'
  }
  if (!isLogin.value && password.value && password.value.length < 6) {
    return '密码至少需要6个字符'
  }
  return null
})

// 方法
function toggleMode() {
  mode.value = isLogin.value ? 'register' : 'login'
  authStore.error = null
}

async function handleSubmit() {
  if (!canSubmit.value) return

  let success: boolean
  if (isLogin.value) {
    success = await authStore.login({
      username: username.value,
      password: password.value,
    })
  } else {
    success = await authStore.register({
      username: username.value,
      password: password.value,
      email: email.value || undefined,
    })
  }

  if (success) {
    emit('success')
  }
}
</script>

<template>
  <div class="auth-screen">
    <!-- 游戏标题 -->
    <div class="game-title">
      <pre class="ascii-art">
╔════════════════════════════════════════════════════════════╗
║  ████████╗███████╗██╗  ██╗████████╗    ██╗    ██╗ ██████╗ ██╗    ██╗  ║
║  ╚══██╔══╝██╔════╝╚██╗██╔╝╚══██╔══╝    ██║    ██║██╔═══██╗██║    ██║  ║
║     ██║   █████╗   ╚███╔╝    ██║       ██║ █╗ ██║██║   ██║██║ █╗ ██║  ║
║     ██║   ██╔══╝   ██╔██╗    ██║       ██║███╗██║██║   ██║██║███╗██║  ║
║     ██║   ███████╗██╔╝ ██╗   ██║       ╚███╔███╔╝╚██████╔╝╚███╔███╔╝  ║
║     ╚═╝   ╚══════╝╚═╝  ╚═╝   ╚═╝        ╚══╝╚══╝  ╚═════╝  ╚══╝╚══╝   ║
╚════════════════════════════════════════════════════════════╝
      </pre>
      <div class="subtitle">— 艾泽拉斯放置冒险 —</div>
    </div>

    <!-- 登录/注册表单 -->
    <div class="auth-form-container">
      <div class="auth-form">
        <h2 class="form-title">{{ title }}</h2>
        
        <form @submit.prevent="handleSubmit">
          <div class="form-group">
            <label>用户名</label>
            <input 
              v-model="username" 
              type="text" 
              placeholder="请输入用户名"
              maxlength="32"
              autofocus
            />
          </div>

          <div class="form-group" v-if="!isLogin">
            <label>邮箱 <span class="optional">(可选)</span></label>
            <input 
              v-model="email" 
              type="email" 
              placeholder="请输入邮箱"
            />
          </div>

          <div class="form-group">
            <label>密码</label>
            <input 
              v-model="password" 
              type="password" 
              placeholder="请输入密码"
            />
          </div>

          <div class="form-group" v-if="!isLogin">
            <label>确认密码</label>
            <input 
              v-model="confirmPassword" 
              type="password" 
              placeholder="请再次输入密码"
            />
            <span class="field-error" v-if="passwordError">{{ passwordError }}</span>
          </div>

          <div class="error-message" v-if="authStore.error">
            <span class="error-icon">⚠</span> {{ authStore.error }}
          </div>

          <button 
            type="submit" 
            class="submit-btn"
            :disabled="!canSubmit || authStore.loading"
          >
            <span v-if="authStore.loading" class="loading-dots">...</span>
            <span v-else>{{ submitText }}</span>
          </button>
        </form>

        <div class="switch-mode" @click="toggleMode">
          {{ switchText }}
        </div>
      </div>
    </div>

    <!-- 版本信息 -->
    <div class="version-info">
      v0.1.0 | 联盟与部落，为了艾泽拉斯！
    </div>
  </div>
</template>

<style scoped>
.auth-screen {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 100vh;
  padding: 20px;
}

.game-title {
  text-align: center;
  margin-bottom: 30px;
}

.ascii-art {
  color: var(--terminal-gold);
  font-size: 8px;
  line-height: 1.2;
  text-shadow: 0 0 10px var(--terminal-gold);
  margin: 0;
}

.subtitle {
  color: var(--terminal-cyan);
  font-size: 12px;
  margin-top: 10px;
  letter-spacing: 4px;
}

.auth-form-container {
  width: 100%;
  max-width: 400px;
}

.auth-form {
  border: 1px solid var(--terminal-green);
  padding: 30px;
  background: rgba(0, 50, 0, 0.3);
  box-shadow: 0 0 20px rgba(0, 255, 0, 0.1);
}

.form-title {
  color: var(--terminal-green);
  text-align: center;
  margin-bottom: 25px;
  font-size: 18px;
  text-transform: uppercase;
  letter-spacing: 3px;
}

.form-group {
  margin-bottom: 20px;
}

.form-group label {
  display: block;
  color: var(--terminal-cyan);
  margin-bottom: 8px;
  font-size: 12px;
  text-transform: uppercase;
}

.form-group .optional {
  color: var(--terminal-gray);
  font-size: 10px;
  text-transform: lowercase;
}

.form-group input {
  width: 100%;
  padding: 12px 15px;
  background: rgba(0, 0, 0, 0.5);
  border: 1px solid var(--terminal-gray);
  color: var(--terminal-green);
  font-family: inherit;
  font-size: 14px;
  transition: all 0.3s ease;
}

.form-group input:focus {
  outline: none;
  border-color: var(--terminal-green);
  box-shadow: 0 0 10px rgba(0, 255, 0, 0.2);
}

.form-group input::placeholder {
  color: var(--terminal-gray);
}

.field-error {
  display: block;
  color: var(--terminal-red);
  font-size: 10px;
  margin-top: 5px;
}

.error-message {
  background: rgba(255, 0, 0, 0.1);
  border: 1px solid var(--terminal-red);
  color: var(--terminal-red);
  padding: 10px 15px;
  margin-bottom: 20px;
  font-size: 12px;
}

.error-icon {
  margin-right: 5px;
}

.submit-btn {
  width: 100%;
  padding: 15px;
  background: transparent;
  border: 2px solid var(--terminal-gold);
  color: var(--terminal-gold);
  font-family: inherit;
  font-size: 14px;
  text-transform: uppercase;
  letter-spacing: 2px;
  cursor: pointer;
  transition: all 0.3s ease;
}

.submit-btn:hover:not(:disabled) {
  background: var(--terminal-gold);
  color: var(--terminal-bg);
  box-shadow: 0 0 20px rgba(255, 215, 0, 0.3);
}

.submit-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.loading-dots {
  animation: blink 1s infinite;
}

@keyframes blink {
  0%, 50% { opacity: 1; }
  51%, 100% { opacity: 0; }
}

.switch-mode {
  text-align: center;
  margin-top: 20px;
  color: var(--terminal-cyan);
  font-size: 12px;
  cursor: pointer;
  transition: color 0.3s ease;
}

.switch-mode:hover {
  color: var(--terminal-green);
  text-decoration: underline;
}

.version-info {
  position: fixed;
  bottom: 20px;
  color: var(--terminal-gray);
  font-size: 10px;
}

/* 响应式 */
@media (max-width: 480px) {
  .ascii-art {
    font-size: 5px;
  }
  
  .auth-form {
    padding: 20px;
  }
}
</style>











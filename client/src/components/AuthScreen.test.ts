import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import AuthScreen from './AuthScreen.vue'
import { useAuthStore } from '@/stores/auth'

describe('AuthScreen Component', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  // ═══════════════════════════════════════════════════════════
  // 渲染测试
  // ═══════════════════════════════════════════════════════════

  describe('Rendering', () => {
    it('should render login form by default', () => {
      const wrapper = mount(AuthScreen)
      
      expect(wrapper.find('.form-title').text()).toBe('登录')
      expect(wrapper.find('input[type="text"]').exists()).toBe(true)
      expect(wrapper.find('input[type="password"]').exists()).toBe(true)
    })

    it('should render game title', () => {
      const wrapper = mount(AuthScreen)
      
      expect(wrapper.find('.game-title').exists()).toBe(true)
      expect(wrapper.find('.ascii-art').exists()).toBe(true)
    })

    it('should show register fields when in register mode', async () => {
      const wrapper = mount(AuthScreen)
      
      // Use component method instead of DOM event
      await wrapper.vm.$nextTick()
      wrapper.vm.toggleMode()
      await wrapper.vm.$nextTick()
      
      expect(wrapper.find('.form-title').text()).toBe('注册')
      expect(wrapper.findAll('input[type="password"]').length).toBe(2) // password + confirm
      expect(wrapper.find('input[type="email"]').exists()).toBe(true)
    })

    it('should show version info', () => {
      const wrapper = mount(AuthScreen)
      
      expect(wrapper.find('.version-info').exists()).toBe(true)
      expect(wrapper.find('.version-info').text()).toContain('v0.1.0')
    })
  })

  // ═══════════════════════════════════════════════════════════
  // 表单交互测试
  // ═══════════════════════════════════════════════════════════

  describe('Form Interaction', () => {
    it('should update username on input', async () => {
      const wrapper = mount(AuthScreen)
      
      // Directly set the component's reactive property
      wrapper.vm.username = 'testuser'
      await wrapper.vm.$nextTick()
      
      expect(wrapper.vm.username).toBe('testuser')
      // Verify it's reflected in the input element
      const input = wrapper.find('input[type="text"]')
      expect((input.element as HTMLInputElement).value).toBe('testuser')
    })

    it('should update password on input', async () => {
      const wrapper = mount(AuthScreen)
      
      // Directly set the component's reactive property
      wrapper.vm.password = 'mypassword'
      await wrapper.vm.$nextTick()
      
      expect(wrapper.vm.password).toBe('mypassword')
      // Verify it's reflected in the input element
      const input = wrapper.find('input[type="password"]')
      expect((input.element as HTMLInputElement).value).toBe('mypassword')
    })

    it('should toggle between login and register mode', async () => {
      const wrapper = mount(AuthScreen)
      
      expect(wrapper.find('.form-title').text()).toBe('登录')
      
      wrapper.vm.toggleMode()
      await wrapper.vm.$nextTick()
      expect(wrapper.find('.form-title').text()).toBe('注册')
      
      wrapper.vm.toggleMode()
      await wrapper.vm.$nextTick()
      expect(wrapper.find('.form-title').text()).toBe('登录')
    })

    it('should clear error when toggling mode', async () => {
      const wrapper = mount(AuthScreen)
      const authStore = useAuthStore()
      
      authStore.error = '测试错误'
      await wrapper.vm.$nextTick()
      
      wrapper.vm.toggleMode()
      await wrapper.vm.$nextTick()
      
      expect(authStore.error).toBeNull()
    })
  })

  // ═══════════════════════════════════════════════════════════
  // 表单验证测试
  // ═══════════════════════════════════════════════════════════

  describe('Form Validation', () => {
    it('should disable submit button when username is empty', async () => {
      const wrapper = mount(AuthScreen)
      
      wrapper.vm.password = 'password'
      await wrapper.vm.$nextTick()
      
      const submitBtn = wrapper.find('.submit-btn')
      expect(submitBtn.attributes('disabled')).toBeDefined()
    })

    it('should disable submit button when password is empty', async () => {
      const wrapper = mount(AuthScreen)
      
      wrapper.vm.username = 'username'
      await wrapper.vm.$nextTick()
      
      const submitBtn = wrapper.find('.submit-btn')
      expect(submitBtn.attributes('disabled')).toBeDefined()
    })

    it('should enable submit button when form is valid', async () => {
      const wrapper = mount(AuthScreen)
      
      wrapper.vm.username = 'username'
      wrapper.vm.password = 'password'
      await wrapper.vm.$nextTick()
      
      const submitBtn = wrapper.find('.submit-btn')
      expect(submitBtn.attributes('disabled')).toBeUndefined()
    })

    it('should show password mismatch error in register mode', async () => {
      const wrapper = mount(AuthScreen)
      
      wrapper.vm.toggleMode()
      wrapper.vm.username = 'username'
      wrapper.vm.password = 'password1'
      wrapper.vm.confirmPassword = 'password2'
      await wrapper.vm.$nextTick()
      
      expect(wrapper.find('.field-error').exists()).toBe(true)
      expect(wrapper.find('.field-error').text()).toContain('两次输入的密码不一致')
    })

    it('should show password length error in register mode', async () => {
      const wrapper = mount(AuthScreen)
      
      wrapper.vm.toggleMode()
      wrapper.vm.username = 'username'
      wrapper.vm.password = 'short'
      await wrapper.vm.$nextTick()
      
      expect(wrapper.find('.field-error').exists()).toBe(true)
      expect(wrapper.find('.field-error').text()).toContain('密码至少需要6个字符')
    })

    it('should disable submit in register mode when passwords do not match', async () => {
      const wrapper = mount(AuthScreen)
      
      wrapper.vm.toggleMode()
      wrapper.vm.username = 'username'
      wrapper.vm.password = 'password1'
      wrapper.vm.confirmPassword = 'password2'
      await wrapper.vm.$nextTick()
      
      const submitBtn = wrapper.find('.submit-btn')
      expect(submitBtn.attributes('disabled')).toBeDefined()
    })
  })

  // ═══════════════════════════════════════════════════════════
  // 表单提交测试
  // ═══════════════════════════════════════════════════════════

  describe('Form Submission', () => {
    it('should call login on submit in login mode', async () => {
      const wrapper = mount(AuthScreen)
      const authStore = useAuthStore()
      const loginSpy = vi.spyOn(authStore, 'login').mockResolvedValue(true)
      
      wrapper.vm.username = 'testuser'
      wrapper.vm.password = 'password'
      await wrapper.vm.$nextTick()
      await wrapper.vm.handleSubmit()
      await flushPromises()
      
      expect(loginSpy).toHaveBeenCalledWith({
        username: 'testuser',
        password: 'password',
      })
    })

    it('should call register on submit in register mode', async () => {
      const wrapper = mount(AuthScreen)
      const authStore = useAuthStore()
      const registerSpy = vi.spyOn(authStore, 'register').mockResolvedValue(true)
      
      wrapper.vm.toggleMode()
      wrapper.vm.username = 'newuser'
      wrapper.vm.email = 'test@example.com'
      wrapper.vm.password = 'password123'
      wrapper.vm.confirmPassword = 'password123'
      await wrapper.vm.$nextTick()
      await wrapper.vm.handleSubmit()
      await flushPromises()
      
      expect(registerSpy).toHaveBeenCalledWith({
        username: 'newuser',
        password: 'password123',
        email: 'test@example.com',
      })
    })

    it('should emit success event on successful login', async () => {
      const wrapper = mount(AuthScreen)
      const authStore = useAuthStore()
      vi.spyOn(authStore, 'login').mockResolvedValue(true)
      
      wrapper.vm.username = 'user'
      wrapper.vm.password = 'pass'
      await wrapper.vm.$nextTick()
      await wrapper.vm.handleSubmit()
      await flushPromises()
      
      expect(wrapper.emitted('success')).toBeTruthy()
    })

    it('should not emit success on failed login', async () => {
      const wrapper = mount(AuthScreen)
      const authStore = useAuthStore()
      vi.spyOn(authStore, 'login').mockResolvedValue(false)
      
      wrapper.vm.username = 'user'
      wrapper.vm.password = 'wrong'
      await wrapper.vm.$nextTick()
      await wrapper.vm.handleSubmit()
      await flushPromises()
      
      expect(wrapper.emitted('success')).toBeFalsy()
    })
  })

  // ═══════════════════════════════════════════════════════════
  // 加载状态测试
  // ═══════════════════════════════════════════════════════════

  describe('Loading State', () => {
    it('should show loading indicator when loading', async () => {
      const wrapper = mount(AuthScreen)
      const authStore = useAuthStore()
      authStore.loading = true
      await wrapper.vm.$nextTick()
      
      expect(wrapper.find('.loading-dots').exists()).toBe(true)
    })

    it('should disable submit button when loading', async () => {
      const wrapper = mount(AuthScreen)
      const authStore = useAuthStore()
      
      wrapper.vm.username = 'user'
      wrapper.vm.password = 'pass'
      authStore.loading = true
      await wrapper.vm.$nextTick()
      
      expect(wrapper.find('.submit-btn').attributes('disabled')).toBeDefined()
    })
  })

  // ═══════════════════════════════════════════════════════════
  // 错误显示测试
  // ═══════════════════════════════════════════════════════════

  describe('Error Display', () => {
    it('should display error message from store', async () => {
      const wrapper = mount(AuthScreen)
      const authStore = useAuthStore()
      
      authStore.error = '登录失败：密码错误'
      await wrapper.vm.$nextTick()
      
      expect(wrapper.find('.error-message').exists()).toBe(true)
      expect(wrapper.find('.error-message').text()).toContain('登录失败：密码错误')
    })

    it('should not show error message when no error', () => {
      const wrapper = mount(AuthScreen)
      
      expect(wrapper.find('.error-message').exists()).toBe(false)
    })
  })
})



















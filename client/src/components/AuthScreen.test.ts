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
      
      await wrapper.find('.switch-mode').trigger('click')
      
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
      
      const input = wrapper.find('input[type="text"]')
      await input.setValue('testuser')
      
      expect((input.element as HTMLInputElement).value).toBe('testuser')
    })

    it('should update password on input', async () => {
      const wrapper = mount(AuthScreen)
      
      const input = wrapper.find('input[type="password"]')
      await input.setValue('mypassword')
      
      expect((input.element as HTMLInputElement).value).toBe('mypassword')
    })

    it('should toggle between login and register mode', async () => {
      const wrapper = mount(AuthScreen)
      
      expect(wrapper.find('.form-title').text()).toBe('登录')
      
      await wrapper.find('.switch-mode').trigger('click')
      expect(wrapper.find('.form-title').text()).toBe('注册')
      
      await wrapper.find('.switch-mode').trigger('click')
      expect(wrapper.find('.form-title').text()).toBe('登录')
    })

    it('should clear error when toggling mode', async () => {
      const wrapper = mount(AuthScreen)
      const authStore = useAuthStore()
      
      authStore.error = '测试错误'
      await wrapper.vm.$nextTick()
      
      await wrapper.find('.switch-mode').trigger('click')
      
      expect(authStore.error).toBeNull()
    })
  })

  // ═══════════════════════════════════════════════════════════
  // 表单验证测试
  // ═══════════════════════════════════════════════════════════

  describe('Form Validation', () => {
    it('should disable submit button when username is empty', async () => {
      const wrapper = mount(AuthScreen)
      
      await wrapper.find('input[type="password"]').setValue('password')
      
      const submitBtn = wrapper.find('.submit-btn')
      expect(submitBtn.attributes('disabled')).toBeDefined()
    })

    it('should disable submit button when password is empty', async () => {
      const wrapper = mount(AuthScreen)
      
      await wrapper.find('input[type="text"]').setValue('username')
      
      const submitBtn = wrapper.find('.submit-btn')
      expect(submitBtn.attributes('disabled')).toBeDefined()
    })

    it('should enable submit button when form is valid', async () => {
      const wrapper = mount(AuthScreen)
      
      await wrapper.find('input[type="text"]').setValue('username')
      await wrapper.find('input[type="password"]').setValue('password')
      
      const submitBtn = wrapper.find('.submit-btn')
      expect(submitBtn.attributes('disabled')).toBeUndefined()
    })

    it('should show password mismatch error in register mode', async () => {
      const wrapper = mount(AuthScreen)
      
      await wrapper.find('.switch-mode').trigger('click')
      await wrapper.find('input[type="text"]').setValue('username')
      
      const passwordInputs = wrapper.findAll('input[type="password"]')
      await passwordInputs[0].setValue('password1')
      await passwordInputs[1].setValue('password2')
      
      expect(wrapper.find('.field-error').exists()).toBe(true)
      expect(wrapper.find('.field-error').text()).toContain('两次输入的密码不一致')
    })

    it('should show password length error in register mode', async () => {
      const wrapper = mount(AuthScreen)
      
      await wrapper.find('.switch-mode').trigger('click')
      await wrapper.find('input[type="text"]').setValue('username')
      
      const passwordInputs = wrapper.findAll('input[type="password"]')
      await passwordInputs[0].setValue('short')
      
      expect(wrapper.find('.field-error').exists()).toBe(true)
      expect(wrapper.find('.field-error').text()).toContain('密码至少需要6个字符')
    })

    it('should disable submit in register mode when passwords do not match', async () => {
      const wrapper = mount(AuthScreen)
      
      await wrapper.find('.switch-mode').trigger('click')
      await wrapper.find('input[type="text"]').setValue('username')
      
      const passwordInputs = wrapper.findAll('input[type="password"]')
      await passwordInputs[0].setValue('password1')
      await passwordInputs[1].setValue('password2')
      
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
      
      await wrapper.find('input[type="text"]').setValue('testuser')
      await wrapper.find('input[type="password"]').setValue('password')
      await wrapper.find('form').trigger('submit')
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
      
      await wrapper.find('.switch-mode').trigger('click')
      await wrapper.find('input[type="text"]').setValue('newuser')
      await wrapper.find('input[type="email"]').setValue('test@example.com')
      
      const passwordInputs = wrapper.findAll('input[type="password"]')
      await passwordInputs[0].setValue('password123')
      await passwordInputs[1].setValue('password123')
      
      await wrapper.find('form').trigger('submit')
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
      
      await wrapper.find('input[type="text"]').setValue('user')
      await wrapper.find('input[type="password"]').setValue('pass')
      await wrapper.find('form').trigger('submit')
      await flushPromises()
      
      expect(wrapper.emitted('success')).toBeTruthy()
    })

    it('should not emit success on failed login', async () => {
      const wrapper = mount(AuthScreen)
      const authStore = useAuthStore()
      vi.spyOn(authStore, 'login').mockResolvedValue(false)
      
      await wrapper.find('input[type="text"]').setValue('user')
      await wrapper.find('input[type="password"]').setValue('wrong')
      await wrapper.find('form').trigger('submit')
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
      
      await wrapper.find('input[type="text"]').setValue('user')
      await wrapper.find('input[type="password"]').setValue('pass')
      
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





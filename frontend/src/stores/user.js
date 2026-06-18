import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import * as authApi from '@/api/auth'

const TOKEN_KEY = 'paddle_token'
const USER_KEY = 'paddle_user'

/**
 * 用户全局状态(Composition API 风格)
 * - token / userInfo 持久化到 localStorage,刷新页面后自动恢复
 */
export const useUserStore = defineStore('user', () => {
  // ---- state ----
  const token = ref(localStorage.getItem(TOKEN_KEY) || '')
  const userInfo = ref(JSON.parse(localStorage.getItem(USER_KEY) || 'null'))

  // ---- getters ----
  const isLoggedIn = computed(() => !!token.value)
  const role = computed(() => userInfo.value?.role || '')
  const username = computed(() => userInfo.value?.username || '')
  const companyName = computed(() => userInfo.value?.company_name || '')

  // ---- actions ----
  /**
   * 登录: 成功后将 token 与用户信息写入 localStorage
   */
  async function login(credentials) {
    const data = await authApi.login(credentials)
    token.value = data.token
    userInfo.value = data.user
    localStorage.setItem(TOKEN_KEY, data.token)
    localStorage.setItem(USER_KEY, JSON.stringify(data.user))
    return data
  }

  /**
   * 注册: 仅完成注册动作,不自动登录
   */
  async function register(payload) {
    return authApi.register(payload)
  }

  /**
   * 拉取当前登录用户最新信息
   */
  async function fetchProfile() {
    const data = await authApi.getProfile()
    userInfo.value = data
    localStorage.setItem(USER_KEY, JSON.stringify(data))
    return data
  }

  /**
   * 退出登录: 清空内存状态与 localStorage
   */
  function logout() {
    token.value = ''
    userInfo.value = null
    localStorage.removeItem(TOKEN_KEY)
    localStorage.removeItem(USER_KEY)
  }

  return {
    // state
    token,
    userInfo,
    // getters
    isLoggedIn,
    role,
    username,
    companyName,
    // actions
    login,
    register,
    fetchProfile,
    logout
  }
})

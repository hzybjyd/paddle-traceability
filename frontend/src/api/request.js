import axios from 'axios'
import { ElMessage } from 'element-plus'
import { useUserStore } from '@/stores/user'
import router from '@/router'
import { translateApiMessage } from '@/utils/i18n'

/**
 * Axios instance wrapper
 * - Base URL: import.meta.env.VITE_API_BASE
 * - Request interceptor: automatically inject JWT token
 * - Response interceptor: unified handling of business codes, error prompts, 401 redirect
 */
const request = axios.create({
  baseURL: import.meta.env.VITE_API_BASE || '/api/v1',
  timeout: 15000,
  headers: {
    'Content-Type': 'application/json'
  }
})

// Request interceptor: inject Authorization header
request.interceptors.request.use(
  (config) => {
    const userStore = useUserStore()
    if (userStore.token) {
      config.headers.Authorization = `Bearer ${userStore.token}`
    }
    return config
  },
  (error) => Promise.reject(error)
)

// Response interceptor: unified handling of { code, message, data }
// Backend message is in English -> translated to Chinese via i18n table
request.interceptors.response.use(
  (response) => {
    const res = response.data
    // Compatible with responses without business code (e.g. non /api/v1 requests)
    if (res && typeof res === 'object' && 'code' in res) {
      // Translate message into Chinese for consistent display
      res.message = translateApiMessage(res.message)
      if (res.code === 200 || res.code === 201) {
        return res.data
      }
      // Business error
      ElMessage.error(res.message || 'request failed')
      return Promise.reject(new Error(res.message || 'request failed'))
    }
    return res
  },
  (error) => {
    const status = error.response?.status
    const rawMsg = error.response?.data?.message || error.message || 'network error'
    const msg = translateApiMessage(rawMsg) || rawMsg
    if (status === 401) {
      // Token expired: clear user state and redirect to login
      const userStore = useUserStore()
      userStore.logout()
      ElMessage.error('login expired, please log in again')
      router.push('/login')
    } else {
      ElMessage.error(msg)
    }
    return Promise.reject(error)
  }
)

export default request

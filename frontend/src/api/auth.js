import request from './request'

/**
 * 用户注册
 * @param {object} data - { username, password, role, company_name, phone, public_key }
 */
export function register(data) {
  return request({ url: '/auth/register', method: 'post', data })
}

/**
 * 用户登录
 * @param {object} data - { username, password }
 * @returns {Promise<{token: string, expires_at: string, user: object}>}
 */
export function login(data) {
  return request({ url: '/auth/login', method: 'post', data })
}

/**
 * 获取当前登录用户的个人信息
 */
export function getProfile() {
  return request({ url: '/auth/profile', method: 'get' })
}

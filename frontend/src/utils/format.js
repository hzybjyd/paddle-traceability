/**
 * 通用格式化工具函数
 */

/**
 * 补零
 * @param {number} n
 * @returns {string}
 */
function pad(n) {
  return n < 10 ? '0' + n : String(n)
}

/**
 * 将日期对象/字符串/时间戳格式化为指定格式字符串
 * 支持的占位符: YYYY MM DD HH mm ss
 * @param {Date|string|number} date
 * @param {string} fmt 默认 'YYYY-MM-DD HH:mm:ss'
 * @returns {string}
 */
export function formatDate(date, fmt = 'YYYY-MM-DD HH:mm:ss') {
  if (!date) return ''
  const d = date instanceof Date ? date : new Date(date)
  if (isNaN(d.getTime())) return ''

  return fmt
    .replace('YYYY', String(d.getFullYear()))
    .replace('MM', pad(d.getMonth() + 1))
    .replace('DD', pad(d.getDate()))
    .replace('HH', pad(d.getHours()))
    .replace('mm', pad(d.getMinutes()))
    .replace('ss', pad(d.getSeconds()))
}

/**
 * 短日期格式 YYYY-MM-DD
 * @param {Date|string|number} date
 * @returns {string}
 */
export function formatDateShort(date) {
  return formatDate(date, 'YYYY-MM-DD')
}

/**
 * 哈希字符串截断显示
 * 例: 'abcdefghijklmnop' -> 'abcdefghij...'
 * @param {string} hash
 * @param {number} len 前后保留的字符数,默认 10
 * @returns {string}
 */
export function truncateHash(hash, len = 10) {
  if (!hash) return ''
  const s = String(hash)
  if (s.length <= len * 2 + 3) return s
  return `${s.slice(0, len)}...${s.slice(-len)}`
}

/**
 * 复制文本到剪贴板
 * @param {string} text
 * @returns {Promise<boolean>} 成功返回 true,失败返回 false
 */
export async function copyToClipboard(text) {
  if (text === undefined || text === null) return false
  const value = String(text)

  // 优先使用 Clipboard API
  if (navigator.clipboard && window.isSecureContext) {
    try {
      await navigator.clipboard.writeText(value)
      return true
    } catch (e) {
      // 降级到 fallback
    }
  }

  // Fallback: 使用 textarea + execCommand
  try {
    const textarea = document.createElement('textarea')
    textarea.value = value
    textarea.setAttribute('readonly', '')
    textarea.style.position = 'fixed'
    textarea.style.left = '-9999px'
    textarea.style.top = '0'
    document.body.appendChild(textarea)
    textarea.select()
    const ok = document.execCommand('copy')
    document.body.removeChild(textarea)
    return ok
  } catch (e) {
    return false
  }
}

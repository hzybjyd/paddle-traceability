import axios from 'axios'
import { ElMessage } from 'element-plus'

/**
 * 防伪验证(公开接口,无需登录)
 * 返回后端完整响应对象 { code, message, verified, data }
 *  - verified: true  表示链上已查询到该产品,data 为产品溯源信息
 *  - verified: false 表示未查询到该产品,data 为 null
 * @param {string} uid - 产品 UID
 */
export async function verifyProduct(uid) {
  try {
    const response = await axios.get(
      `${import.meta.env.VITE_API_BASE || '/api/v1'}/verify/${uid}`,
      { timeout: 15000 }
    )
    return response.data
  } catch (error) {
    const status = error.response?.status
    const msg = error.response?.data?.message || error.message || '网络异常'
    ElMessage.error(msg)
    if (status === 404) {
      // 产品不存在时,后端仍可能返回 200 + verified=false;
      // 此处为极端网络错误兜底
    }
    return Promise.reject(error)
  }
}

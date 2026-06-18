import request from './request'

/**
 * 创建产品(自动上链存证)
 * @param {object} data - { brand, model, material, rubber_type, batch_no, production_date, quality_report_hash }
 */
export function createProduct(data) {
  return request({ url: '/products', method: 'post', data })
}

/**
 * 获取产品列表
 * @param {object} params - { page, page_size, status }
 */
export function listProducts(params) {
  return request({ url: '/products', method: 'get', params })
}

/**
 * 获取产品详情
 * @param {string} id - 产品 UID(雪花算法生成的 19 位 ID)
 */
export function getProduct(id) {
  return request({ url: `/products/${id}`, method: 'get' })
}

/**
 * 更新产品信息
 * @param {string} id - 产品 UID
 * @param {object} data - 待更新的字段
 */
export function updateProduct(id, data) {
  return request({ url: `/products/${id}`, method: 'put', data })
}

/**
 * 获取产品全生命周期溯源信息
 * @param {string} id - 产品 UID
 */
export function getProductTrace(id) {
  return request({ url: `/products/${id}/trace`, method: 'get' })
}

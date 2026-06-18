import request from './request'

/**
 * 添加物流记录(自动上链存证并更新产品状态)
 * @param {object} data - { product_uid, action, warehouse_name, location, carrier, remark }
 *   - action: 'INBOUND' 入库 / 'OUTBOUND' 出库
 */
export function addLogistics(data) {
  return request({ url: '/logistics', method: 'post', data })
}

/**
 * 查询物流记录列表
 * @param {object} params - { product_uid }
 */
export function getLogistics(params) {
  return request({ url: '/logistics', method: 'get', params })
}

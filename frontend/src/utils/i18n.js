// English -> Chinese mapping for backend identifiers
// Backend stores all values in English; frontend displays in Chinese.
// When the backend adds a new enum value, add it here.

export const PRODUCT_STATUS_MAP = {
  PRODUCED: '已生产',
  IN_TRANSIT: '运输中',
  IN_STOCK: '在库',
  SOLD: '已售出',
}

export const USER_ROLE_MAP = {
  FACTORY: '厂商',
  LOGISTICS: '物流',
  RETAILER: '经销商',
}

export const TX_TYPE_MAP = {
  CREATE: '创建',
  TRANSFER: '流转',
  CONFIRM: '确认',
  UPDATE: '更新',
}

export const CHAIN_STATUS_MAP = {
  CONFIRMED: '已确认',
  PENDING: '待确认',
  FAILED: '失败',
}

export const STAGE_MAP = {
  PRODUCTION: '生产出厂',
  LOGISTICS_TRANSFER: '物流流转',
  WAREHOUSE_INBOUND: '物流入库',
  WAREHOUSE_OUTBOUND: '物流出库',
  SALE_CONFIRM: '销售确认',
}

export const LOGISTICS_ACTION_MAP = {
  INBOUND: '入库',
  OUTBOUND: '出库',
}

// Generic helper: look up the Chinese label, fallback to the raw key
export function translate(map, key, fallback) {
  if (key === null || key === undefined || key === '') {
    return fallback ?? ''
  }
  return map[key] ?? key
}

// API response message -> Chinese
// The backend message is in English; map to Chinese for direct display.
// If a message is not in the map, the original English is returned.
export const API_MESSAGE_MAP = {
  'register success': '注册成功',
  'login success': '登录成功',
  'product created and recorded on chain': '产品创建成功，已上链存证',
  'product status updated and recorded on chain': '产品状态更新成功，已上链存证',
  'logistics record added and recorded on chain': '物流记录已添加并上链存证',
  'verified: this product is authentic': '验证通过，该产品为正品',
  'verification failed: no product record found on chain, may be counterfeit':
    '验证失败，链上未查询到该产品记录，可能为假冒产品',
  'username or password incorrect': '用户名或密码错误',
  'invalid request': '请求参数错误',
  'username already exists': '用户名已存在',
  'invalid role type': '无效的角色类型',
  'missing product_uid': '请提供产品唯一标识',
  'missing product_uid parameter': '请提供 product_uid 参数',
  'product not found': '产品不存在',
  'error occurred during verification': '验证过程中发生错误',
  'only FACTORY role can create products': '仅厂商角色可创建产品',
  'only LOGISTICS and RETAILER can add logistics records': '仅物流和经销商可添加物流记录',
  'permission denied': '无权执行此操作',
  'invalid product id': '无效的产品ID',
  'invalid production date format, expected YYYY-MM-DD': '生产日期格式错误，应为 YYYY-MM-DD',
  'chain attestation failed': '上链存证失败',
  'create product failed': '创建产品失败',
  'save tx record failed': '保存交易记录失败',
  'update product failed': '更新产品失败',
  'update product status failed': '更新产品状态失败',
  'create logistics record failed': '创建物流记录失败',
  'query product failed': '查询产品失败',
  'query tx records failed': '查询交易记录失败',
  'query logistics records failed': '查询物流记录失败',
  'query product list failed': '查询产品列表失败',
  'password encryption failed': '密码加密失败',
  'create user failed': '创建用户失败',
  'query user failed': '查询用户失败',
  'generate token failed': '生成令牌失败',
  'user not found': '用户不存在',
  'authentication failed': '认证失败',
  'missing authorization header': '未提供认证信息',
  'invalid authorization format': '认证格式错误',
}

export function translateApiMessage(msg) {
  if (!msg) return ''
  return API_MESSAGE_MAP[msg] ?? msg
}

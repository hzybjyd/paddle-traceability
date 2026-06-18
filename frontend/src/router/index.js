import { createRouter, createWebHistory } from 'vue-router'
import { useUserStore } from '@/stores/user'

/**
 * 路由表
 * - meta.requiresAuth: 是否需要登录
 * - meta.roles: 允许访问的角色列表(空数组表示仅需登录,不限制角色)
 * - views 目录尚未创建,统一使用动态 import,后续创建文件即可生效
 */
const routes = [
  // ============ 独立路由(不套 MainLayout) ============
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/auth/Login.vue'),
    meta: { requiresAuth: false, title: '登录' }
  },

  // ============ MainLayout 包裹的路由 ============
  {
    path: '/',
    component: () => import('@/layouts/MainLayout.vue'),
    redirect: '/verify',
    children: [
      // 公开：防伪验证
      {
        path: 'verify',
        name: 'Verify',
        component: () => import('@/views/verify/VerifyPage.vue'),
        meta: { requiresAuth: false, title: '防伪验证' }
      },
      {
        path: 'verify/:uid',
        name: 'VerifyDetail',
        component: () => import('@/views/verify/VerifyPage.vue'),
        meta: { requiresAuth: false, title: '防伪验证' }
      },
      // 需登录路由
      {
        path: 'factory',
        name: 'FactoryProductList',
        // 厂家/经销商产品列表
        component: () => import('@/views/factory/ProductList.vue'),
        meta: {
          requiresAuth: true,
          roles: ['FACTORY', 'RETAILER'],
          title: '产品管理'
        }
      },
      {
        path: 'factory/create',
        name: 'FactoryProductCreate',
        // 厂家创建产品(上链存证，仅厂家可用)
        component: () => import('@/views/factory/ProductCreate.vue'),
        meta: {
          requiresAuth: true,
          roles: ['FACTORY'],
          title: '创建产品'
        }
      },
      {
        path: 'logistics',
        name: 'LogisticsPanel',
        // 物流操作面板(仅物流公司)
        component: () => import('@/views/logistics/LogisticsPanel.vue'),
        meta: {
          requiresAuth: true,
          roles: ['LOGISTICS'],
          title: '物流管理'
        }
      },
    ]
  },

  // ============ 兜底 404 ============
  {
    path: '/:pathMatch(.*)*',
    name: 'NotFound',
    // 未匹配到的路由跳转到防伪验证页
    redirect: '/verify'
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

/**
 * 全局前置守卫
 * 1. 需要登录但未登录 -> 跳转 /login
 * 2. 已登录访问 /login -> 跳转 /verify
 * 3. 角色不匹配 -> 跳转 /verify
 */
router.beforeEach((to, from, next) => {
  const userStore = useUserStore()
  const requiresAuth = to.matched.some((r) => r.meta.requiresAuth)
  const allowedRoles = to.meta.roles

  // 需要登录但未登录
  if (requiresAuth && !userStore.isLoggedIn) {
    return next({ path: '/login', query: { redirect: to.fullPath } })
  }

  // 已登录访问登录页,直接进入防伪验证
  if (to.path === '/login' && userStore.isLoggedIn) {
    return next('/verify')
  }

  // 角色权限校验
  if (
    requiresAuth &&
    Array.isArray(allowedRoles) &&
    allowedRoles.length > 0 &&
    !allowedRoles.includes(userStore.role)
  ) {
    return next('/verify')
  }

  next()
})

export default router

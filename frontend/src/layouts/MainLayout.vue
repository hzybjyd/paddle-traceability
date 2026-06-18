<template>
  <el-container class="main-layout">
    <!-- 顶部 -->
    <el-header class="main-header" height="60px">
      <div class="main-header__brand">
        <el-icon :size="22" class="brand-icon"><Histogram /></el-icon>
        <span class="brand-title">乒乓球拍防伪溯源系统</span>
      </div>
      <div class="main-header__user">
        <!-- 未登录：显示登录/注册 -->
        <template v-if="!isLoggedIn">
          <el-button
            type="primary"
            plain
            size="small"
            @click="$router.push('/login')"
          >
            登录
          </el-button>
          <el-button
            type="primary"
            size="small"
            @click="$router.push({ path: '/login', query: { tab: 'register' } })"
          >
            注册
          </el-button>
        </template>
        <!-- 已登录 -->
        <template v-else>
          <span class="user-company">{{ companyName || username || '未登录' }}</span>
          <el-tag
            v-if="role"
            :type="roleTagType"
            size="small"
            effect="light"
            round
            class="user-role"
          >
            {{ roleLabel }}
          </el-tag>
          <el-button
            type="primary"
            plain
            :icon="SwitchButton"
            size="small"
            class="logout-btn"
            @click="handleLogout"
          >
            退出登录
          </el-button>
        </template>
      </div>
    </el-header>

    <el-container class="main-body">
      <!-- 左侧菜单 -->
      <el-aside width="220px" class="main-aside">
        <el-menu
          :default-active="activeMenu"
          :router="true"
          class="main-menu"
          background-color="#001529"
          text-color="#cfd8dc"
          active-text-color="#ffffff"
        >
          <!-- 通用:防伪验证 -->
          <el-menu-item index="/verify">
            <el-icon><CircleCheck /></el-icon>
            <template #title>防伪验证</template>
          </el-menu-item>

          <!-- FACTORY -->
          <template v-if="role === 'FACTORY'">
            <el-menu-item index="/factory">
              <el-icon><Goods /></el-icon>
              <template #title>产品列表</template>
            </el-menu-item>
            <el-menu-item index="/factory/create">
              <el-icon><Plus /></el-icon>
              <template #title>创建产品</template>
            </el-menu-item>
          </template>

          <!-- RETAILER -->
          <template v-if="role === 'RETAILER'">
            <el-menu-item index="/factory">
              <el-icon><Goods /></el-icon>
              <template #title>产品管理</template>
            </el-menu-item>
          </template>

          <!-- LOGISTICS -->
          <template v-if="role === 'LOGISTICS'">
            <el-menu-item index="/logistics">
              <el-icon><Van /></el-icon>
              <template #title>物流管理</template>
            </el-menu-item>
          </template>

        </el-menu>
      </el-aside>

      <!-- 内容区 -->
      <el-main class="main-content">
        <router-view v-slot="{ Component }">
          <transition name="fade" mode="out-in">
            <component :is="Component" />
          </transition>
        </router-view>
      </el-main>
    </el-container>
  </el-container>
</template>

<script setup>
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  CircleCheck,
  Goods,
  Plus,
  Van,
  Histogram,
  SwitchButton
} from '@element-plus/icons-vue'
import { useUserStore } from '@/stores/user'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()

// 用户信息
const isLoggedIn = computed(() => userStore.isLoggedIn)
const role = computed(() => userStore.role)
const username = computed(() => userStore.username)
const companyName = computed(() => userStore.companyName)

// 角色 -> 中文标签
const ROLE_LABELS = {
  FACTORY: '厂家',
  LOGISTICS: '物流',
  RETAILER: '经销商'
}

const roleLabel = computed(() => ROLE_LABELS[role.value] || role.value || '')

// 角色 -> 标签颜色
const ROLE_TAG_TYPES = {
  FACTORY: 'success',
  LOGISTICS: 'warning',
  RETAILER: 'primary'
}

const roleTagType = computed(() => ROLE_TAG_TYPES[role.value] || 'info')

// 菜单高亮:根据当前路由的顶级路径
const activeMenu = computed(() => {
  const path = route.path
  if (path.startsWith('/verify')) return '/verify'
  if (path.startsWith('/factory/create')) return '/factory/create'
  if (path.startsWith('/factory')) return '/factory'
  if (path.startsWith('/logistics')) return '/logistics'
  return path
})

// 退出登录
async function handleLogout() {
  try {
    await ElMessageBox.confirm('确认退出登录吗?', '提示', {
      type: 'warning',
      confirmButtonText: '退出',
      cancelButtonText: '取消'
    })
    userStore.logout()
    ElMessage.success('已退出登录')
    router.push('/login')
  } catch (e) {
    // 用户取消
  }
}
</script>

<style scoped>
.main-layout {
  height: 100vh;
  width: 100%;
}

.main-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  background: linear-gradient(90deg, #1d4ed8 0%, #3b82f6 100%);
  color: #ffffff;
  padding: 0 20px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
}

.main-header__brand {
  display: flex;
  align-items: center;
  gap: 8px;
}

.brand-icon {
  color: #ffffff;
}

.brand-title {
  font-size: 18px;
  font-weight: 600;
  color: #ffffff;
  letter-spacing: 1px;
}

.main-header__user {
  display: flex;
  align-items: center;
  gap: 12px;
}

.user-company {
  color: #ffffff;
  font-size: 14px;
  max-width: 240px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.user-role {
  background: rgba(255, 255, 255, 0.2);
  border-color: rgba(255, 255, 255, 0.3);
  color: #ffffff;
}

.logout-btn {
  background: rgba(255, 255, 255, 0.15);
  border-color: rgba(255, 255, 255, 0.4);
  color: #ffffff;
}

.logout-btn:hover {
  background: rgba(255, 255, 255, 0.25);
  border-color: #ffffff;
  color: #ffffff;
}

.main-body {
  height: calc(100vh - 60px);
}

.main-aside {
  background-color: #001529;
  overflow: hidden;
}

.main-menu {
  border-right: none;
  height: 100%;
}

.main-content {
  background: #f5f7fa;
  padding: 20px;
  overflow: auto;
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.2s ease;
}
.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>

<template>
  <div class="auth-page">
    <div class="auth-card">
      <!-- 顶部 Logo + 标题 -->
      <div class="auth-header">
        <div class="auth-logo">
          <el-icon :size="40" color="#fff">
            <TrophyBase />
          </el-icon>
        </div>
        <h1 class="auth-title">乒乓球拍防伪溯源系统</h1>
        <p class="auth-subtitle">基于区块链的防伪溯源解决方案</p>
      </div>

      <!-- 登录 / 注册 Tab 切换 -->
      <el-tabs v-model="activeTab" class="auth-tabs" stretch>
        <!-- ===== 登录 ===== -->
        <el-tab-pane label="登录" name="login">
          <el-form
            ref="loginFormRef"
            :model="loginForm"
            :rules="loginRules"
            size="large"
            @keyup.enter="handleLogin"
          >
            <el-form-item prop="username">
              <el-input
                v-model="loginForm.username"
                placeholder="请输入用户名"
                :prefix-icon="User"
                clearable
              />
            </el-form-item>
            <el-form-item prop="password">
              <el-input
                v-model="loginForm.password"
                type="password"
                placeholder="请输入密码"
                :prefix-icon="Lock"
                show-password
                clearable
              />
            </el-form-item>
            <el-form-item>
              <el-button
                type="primary"
                :loading="loginLoading"
                class="auth-submit"
                @click="handleLogin"
              >
                登 录
              </el-button>
            </el-form-item>
          </el-form>
        </el-tab-pane>

        <!-- ===== 注册 ===== -->
        <el-tab-pane label="注册" name="register">
          <el-form
            ref="registerFormRef"
            :model="registerForm"
            :rules="registerRules"
            size="large"
            @keyup.enter="handleRegister"
          >
            <el-form-item prop="username">
              <el-input
                v-model="registerForm.username"
                placeholder="用户名 (3-50字符)"
                :prefix-icon="User"
                clearable
              />
            </el-form-item>
            <el-form-item prop="password">
              <el-input
                v-model="registerForm.password"
                type="password"
                placeholder="密码 (至少 6 位)"
                :prefix-icon="Lock"
                show-password
                clearable
              />
            </el-form-item>
            <el-form-item prop="company_name">
              <el-input
                v-model="registerForm.company_name"
                placeholder="企业/组织名称"
                :prefix-icon="OfficeBuilding"
                clearable
              />
            </el-form-item>
            <el-form-item prop="phone">
              <el-input
                v-model="registerForm.phone"
                placeholder="手机号 (选填)"
                :prefix-icon="Phone"
                clearable
              />
            </el-form-item>
            <el-form-item prop="role">
              <el-select
                v-model="registerForm.role"
                placeholder="请选择角色"
                style="width: 100%"
              >
                <el-option
                  v-for="opt in roleOptions"
                  :key="opt.value"
                  :label="opt.label"
                  :value="opt.value"
                />
              </el-select>
            </el-form-item>
            <el-form-item>
              <el-button
                type="primary"
                :loading="registerLoading"
                class="auth-submit"
                @click="handleRegister"
              >
                注 册
              </el-button>
            </el-form-item>
          </el-form>
        </el-tab-pane>
      </el-tabs>

      <!-- 底部 -->
      <div class="auth-footer">
        基于 Vue 3 + Element Plus + XuperChain
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import {
  User,
  Lock,
  Phone,
  OfficeBuilding
} from '@element-plus/icons-vue'
import { useUserStore } from '@/stores/user'

const router = useRouter()
const route = useRoute()
const userStore = useUserStore()

// 当前激活的 Tab
const activeTab = ref(route.query.tab === 'register' ? 'register' : 'login')

// ===== 登录表单 =====
const loginFormRef = ref(null)
const loginLoading = ref(false)
const loginForm = reactive({
  username: '',
  password: ''
})
const loginRules = {
  username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
  password: [{ required: true, message: '请输入密码', trigger: 'blur' }]
}

async function handleLogin() {
  if (!loginFormRef.value) return
  await loginFormRef.value.validate(async (valid) => {
    if (!valid) return
    loginLoading.value = true
    try {
      const data = await userStore.login({
        username: loginForm.username,
        password: loginForm.password
      })
      ElMessage.success('登录成功')
      // 根据角色跳转;如果 URL 带 redirect 优先使用
      const redirect = route.query.redirect
      if (redirect) {
        router.push(redirect)
      } else {
        router.push(resolveHomeByRole(data?.user?.role))
      }
    } catch (err) {
      // 错误提示由 request 拦截器统一处理
    } finally {
      loginLoading.value = false
    }
  })
}

// ===== 注册表单 =====
const registerFormRef = ref(null)
const registerLoading = ref(false)
const registerForm = reactive({
  username: '',
  password: '',
  company_name: '',
  phone: '',
  role: ''
})
const registerRules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 3, max: 50, message: '用户名长度 3-50 字符', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, message: '密码至少 6 位', trigger: 'blur' }
  ],
  role: [{ required: true, message: '请选择角色', trigger: 'change' }],
  phone: [
    {
      pattern: /^$|^1[3-9]\d{9}$/,
      message: '请输入正确的手机号',
      trigger: 'blur'
    }
  ]
}

const roleOptions = [
  { label: '厂家', value: 'FACTORY' },
  { label: '物流', value: 'LOGISTICS' },
  { label: '经销商', value: 'RETAILER' }
]

async function handleRegister() {
  if (!registerFormRef.value) return
  await registerFormRef.value.validate(async (valid) => {
    if (!valid) return
    registerLoading.value = true
    try {
      await userStore.register({ ...registerForm })
      ElMessage.success('注册成功,请登录')
      // 重置注册表单并切换到登录 Tab
      resetRegisterForm()
      activeTab.value = 'login'
    } catch (err) {
      // 错误提示由 request 拦截器统一处理
    } finally {
      registerLoading.value = false
    }
  })
}

function resetRegisterForm() {
  registerForm.username = ''
  registerForm.password = ''
  registerForm.company_name = ''
  registerForm.phone = ''
  registerForm.role = ''
  registerFormRef.value?.clearValidate()
}

// 根据角色返回首页
function resolveHomeByRole(role) {
  switch (role) {
    case 'FACTORY':
      return '/factory'
    case 'LOGISTICS':
    case 'RETAILER':
      return '/logistics'
    default:
      return '/verify'
  }
}
</script>

<style scoped>
.auth-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #c2e0ff 0%, #d9c8ff 100%);
  padding: 24px;
}

.auth-card {
  width: 100%;
  max-width: 460px;
  background: #ffffff;
  border-radius: 16px;
  box-shadow: 0 20px 50px rgba(60, 80, 160, 0.18);
  padding: 36px 40px 28px;
}

.auth-header {
  text-align: center;
  margin-bottom: 24px;
}

.auth-logo {
  width: 72px;
  height: 72px;
  margin: 0 auto 14px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #5b8def 0%, #8b6fe8 100%);
  box-shadow: 0 8px 20px rgba(91, 141, 239, 0.4);
}

.auth-title {
  font-size: 22px;
  font-weight: 700;
  color: #303133;
  margin: 0 0 6px;
}

.auth-subtitle {
  font-size: 13px;
  color: #909399;
  margin: 0;
}

.auth-tabs {
  margin-top: 4px;
}

.auth-tabs :deep(.el-tabs__nav-wrap)::after {
  background-color: #ebeef5;
}

.auth-submit {
  width: 100%;
  height: 44px;
  font-size: 15px;
  letter-spacing: 4px;
  background: linear-gradient(135deg, #5b8def 0%, #8b6fe8 100%);
  border: none;
}

.auth-submit:hover {
  background: linear-gradient(135deg, #4a7ade 0%, #7a5ed8 100%);
}

.auth-footer {
  margin-top: 20px;
  text-align: center;
  font-size: 12px;
  color: #c0c4cc;
  letter-spacing: 0.5px;
}
</style>

<template>
  <div class="logistics-panel">
    <PageHeader title="物流管理" subtitle="添加产品流转记录 / 查询产品历史轨迹" />

    <el-row :gutter="16">
      <!-- 上半部分:添加流转记录 -->
      <el-col :xs="24">
        <el-card shadow="hover" class="panel-card">
          <template #header>
            <div class="card-header">
              <el-icon :size="18" color="#5b8def"><Box /></el-icon>
              <span class="card-title">添加流转记录</span>
              <el-tag
                :type="actionTagType"
                size="small"
                effect="light"
                round
                class="card-tag"
              >
                当前动作:{{ actionLabel }}
              </el-tag>
            </div>
          </template>

          <el-form
            ref="addFormRef"
            :model="addForm"
            :rules="addRules"
            label-width="100px"
            size="default"
            class="add-form"
            @submit.prevent
          >
            <el-row :gutter="16">
              <el-col :xs="24" :sm="12">
                <el-form-item label="产品UID" prop="product_uid">
                  <el-input
                    v-model="addForm.product_uid"
                    placeholder="请输入19位数字产品UID"
                    :maxlength="32"
                    clearable
                  />
                </el-form-item>
              </el-col>
              <el-col :xs="24" :sm="12">
                <el-form-item label="动作类型" prop="action">
                  <el-radio-group v-model="addForm.action">
                    <el-radio-button
                      v-for="opt in ACTION_OPTIONS"
                      :key="opt.value"
                      :value="opt.value"
                    >
                      {{ opt.label }}
                    </el-radio-button>
                  </el-radio-group>
                </el-form-item>
              </el-col>
              <el-col :xs="24" :sm="12">
                <el-form-item label="仓库名称" prop="warehouse_name">
                  <el-input
                    v-model="addForm.warehouse_name"
                    placeholder="如:北京中转仓"
                    maxlength="100"
                    clearable
                  />
                </el-form-item>
              </el-col>
              <el-col :xs="24" :sm="12">
                <el-form-item label="地理位置" prop="location">
                  <el-input
                    v-model="addForm.location"
                    placeholder="如:北京市朝阳区"
                    maxlength="200"
                    clearable
                  />
                </el-form-item>
              </el-col>
              <el-col :xs="24" :sm="12">
                <el-form-item label="承运方" prop="carrier">
                  <el-input
                    v-model="addForm.carrier"
                    placeholder="如:顺丰速运"
                    maxlength="100"
                    clearable
                  />
                </el-form-item>
              </el-col>
              <el-col :xs="24" :sm="12">
                <el-form-item label="备注" prop="remark">
                  <el-input
                    v-model="addForm.remark"
                    placeholder="可选,补充说明"
                    maxlength="200"
                    clearable
                  />
                </el-form-item>
              </el-col>
            </el-row>
            <el-form-item>
              <el-button
                type="primary"
                :loading="addSubmitting"
                :icon="Position"
                @click="handleAdd"
              >
                提交并上链
              </el-button>
              <el-button :icon="RefreshLeft" @click="resetAddForm">重置</el-button>
            </el-form-item>
          </el-form>
        </el-card>
      </el-col>

      <!-- 下半部分:查询流转历史 -->
      <el-col :xs="24">
        <el-card shadow="hover" class="panel-card">
          <template #header>
            <div class="card-header">
              <el-icon :size="18" color="#e6a23c"><Search /></el-icon>
              <span class="card-title">查询流转历史</span>
            </div>
          </template>

          <div class="query-row">
            <el-input
              v-model="queryUid"
              placeholder="请输入要查询的产品UID"
              :maxlength="32"
              clearable
              class="query-input"
              @keyup.enter="handleQuery"
            />
            <el-button
              type="primary"
              :icon="Search"
              :loading="queryLoading"
              @click="handleQuery"
            >
              查询
            </el-button>
            <el-button :icon="Refresh" @click="handleQueryReset">清空</el-button>
          </div>

          <el-table
            v-loading="queryLoading"
            :data="records"
            border
            stripe
            style="width: 100%"
            empty-text="暂无流转记录"
            class="records-table"
          >
            <el-table-column label="动作" width="100" align="center">
              <template #default="{ row }">
                <el-tag
                  :type="row.action === 'INBOUND' ? 'success' : 'warning'"
                  effect="light"
                  round
                  size="small"
                >
                  {{ translate(LOGISTICS_ACTION_MAP, row.action, row.action) }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="warehouse_name" label="仓库" min-width="140" show-overflow-tooltip>
              <template #default="{ row }">
                {{ row.warehouse_name || '-' }}
              </template>
            </el-table-column>
            <el-table-column prop="location" label="位置" min-width="160" show-overflow-tooltip>
              <template #default="{ row }">
                {{ row.location || '-' }}
              </template>
            </el-table-column>
            <el-table-column prop="carrier" label="承运方" min-width="120" show-overflow-tooltip>
              <template #default="{ row }">
                {{ row.carrier || '-' }}
              </template>
            </el-table-column>
            <el-table-column prop="remark" label="备注" min-width="160" show-overflow-tooltip>
              <template #default="{ row }">
                {{ row.remark || '-' }}
              </template>
            </el-table-column>
            <el-table-column label="创建时间" width="170">
              <template #default="{ row }">
                <span class="muted">{{ formatDate(row.created_at) }}</span>
              </template>
            </el-table-column>
          </el-table>

          <el-empty
            v-if="!queryLoading && hasQueried && records.length === 0"
            description="未查询到该产品的流转记录"
            :image-size="100"
          />
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { ref, reactive, computed } from 'vue'
import { ElMessage } from 'element-plus'
import {
  Box,
  Search,
  Position,
  RefreshLeft,
  Refresh
} from '@element-plus/icons-vue'
import { useUserStore } from '@/stores/user'
import { addLogistics, getLogistics } from '@/api/logistics'
import { formatDate } from '@/utils/format'
import { LOGISTICS_ACTION_MAP, translate } from '@/utils/i18n'
import PageHeader from '@/components/PageHeader.vue'

const userStore = useUserStore()

// Action options: labels come from LOGISTICS_ACTION_MAP (i18n)
const ACTION_OPTIONS = [
  { label: LOGISTICS_ACTION_MAP.INBOUND, value: 'INBOUND' },
  { label: LOGISTICS_ACTION_MAP.OUTBOUND, value: 'OUTBOUND' }
]

// 19位数字校验
const UID_REGEX = /^\d{19}$/

// ===== 添加流转记录 =====
const addFormRef = ref(null)
const addSubmitting = ref(false)
const addForm = reactive({
  product_uid: '',
  action: 'INBOUND',
  warehouse_name: '',
  location: '',
  carrier: '',
  remark: ''
})

const addRules = {
  product_uid: [
    { required: true, message: '请输入产品UID', trigger: 'blur' },
    {
      pattern: UID_REGEX,
      message: '产品UID必须为19位数字',
      trigger: 'blur'
    }
  ],
  action: [{ required: true, message: '请选择动作类型', trigger: 'change' }]
}

const actionLabel = computed(() => {
  const opt = ACTION_OPTIONS.find((it) => it.value === addForm.action)
  return opt ? opt.label : '-'
})

const actionTagType = computed(() => {
  return addForm.action === 'INBOUND' ? 'success' : 'warning'
})

function resetAddForm() {
  addFormRef.value?.resetFields()
  // resetFields 会重置为初始值(包含 action: 'INBOUND'),无需额外处理
}

async function handleAdd() {
  if (!addFormRef.value) return
  await addFormRef.value.validate(async (valid) => {
    if (!valid) return
    addSubmitting.value = true
    try {
      // 清理空字符串字段
      const payload = {
        product_uid: addForm.product_uid,
        action: addForm.action
      }
      if (addForm.warehouse_name) payload.warehouse_name = addForm.warehouse_name
      if (addForm.location) payload.location = addForm.location
      if (addForm.carrier) payload.carrier = addForm.carrier
      if (addForm.remark) payload.remark = addForm.remark
      await addLogistics(payload)
      ElMessage.success('物流记录已添加并上链存证')
      resetAddForm()
    } catch (err) {
      // 错误提示由 request 拦截器统一处理
    } finally {
      addSubmitting.value = false
    }
  })
}

// ===== 查询流转历史 =====
const queryUid = ref('')
const queryLoading = ref(false)
const hasQueried = ref(false)
const records = ref([])

async function handleQuery() {
  const uid = (queryUid.value || '').trim()
  if (!uid) {
    ElMessage.warning('请输入产品UID')
    return
  }
  if (!UID_REGEX.test(uid)) {
    ElMessage.warning('产品UID必须为19位数字')
    return
  }
  queryLoading.value = true
  hasQueried.value = false
  records.value = []
  try {
    const res = await getLogistics({ product_uid: uid })
    records.value = res?.records || []
    hasQueried.value = true
  } catch (err) {
    hasQueried.value = true
  } finally {
    queryLoading.value = false
  }
}

function handleQueryReset() {
  queryUid.value = ''
  records.value = []
  hasQueried.value = false
}
</script>

<style scoped>
.logistics-panel {
  padding: 0 4px;
}

.panel-card {
  border-radius: 10px;
  margin-bottom: 16px;
}

.card-header {
  display: flex;
  align-items: center;
  gap: 8px;
}

.card-title {
  font-size: 16px;
  font-weight: 600;
  color: #303133;
}

.card-tag {
  margin-left: 4px;
}

.add-form {
  margin-top: 4px;
}

.query-row {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 16px;
  flex-wrap: wrap;
}

.query-input {
  width: 320px;
  max-width: 100%;
}

.records-table {
  margin-top: 4px;
}

.muted {
  color: #909399;
  font-size: 12px;
}

@media (max-width: 600px) {
  .query-input {
    width: 100%;
  }
}
</style>

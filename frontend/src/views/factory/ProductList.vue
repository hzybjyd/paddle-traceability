<template>
  <div class="product-list">
    <PageHeader :title="pageTitle" :subtitle="pageSubtitle">
      <el-button
        v-if="userStore.role === 'FACTORY'"
        type="primary"
        :icon="Plus"
        @click="goCreate"
      >
        创建产品
      </el-button>
    </PageHeader>

    <!-- 统计卡片 -->
    <el-row :gutter="16" class="stat-row">
      <el-col :xs="12" :sm="6">
        <el-card shadow="hover" class="stat-card stat-card--total">
          <div class="stat-card__label">产品总数</div>
          <div class="stat-card__value">{{ stats.total }}</div>
        </el-card>
      </el-col>
      <el-col :xs="12" :sm="6">
        <el-card shadow="hover" class="stat-card stat-card--success">
          <div class="stat-card__label">已生产</div>
          <div class="stat-card__value">{{ stats.produced }}</div>
        </el-card>
      </el-col>
      <el-col :xs="12" :sm="6">
        <el-card shadow="hover" class="stat-card stat-card--warning">
          <div class="stat-card__label">运输中</div>
          <div class="stat-card__value">{{ stats.inTransit }}</div>
        </el-card>
      </el-col>
      <el-col :xs="12" :sm="6">
        <el-card shadow="hover" class="stat-card stat-card--info">
          <div class="stat-card__label">已售出</div>
          <div class="stat-card__value">{{ stats.sold }}</div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 筛选栏 -->
    <el-card shadow="never" class="filter-card">
      <el-form :inline="true" :model="filterForm" class="filter-form">
        <el-form-item label="状态">
          <el-select
            v-model="filterForm.status"
            placeholder="全部状态"
            clearable
            style="width: 160px"
            @change="handleSearch"
          >
            <el-option
              v-for="opt in STATUS_OPTIONS"
              :key="opt.value"
              :label="opt.label"
              :value="opt.value"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="关键字">
          <el-input
            v-model="filterForm.keyword"
            placeholder="品牌 / 型号 / UID"
            clearable
            style="width: 240px"
            @keyup.enter="handleSearch"
            @clear="handleSearch"
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" :icon="Search" @click="handleSearch">查询</el-button>
          <el-button :icon="Refresh" @click="handleReset">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 表格 -->
    <el-card shadow="never" class="table-card">
      <el-table
        v-loading="loading"
        :data="filteredList"
        border
        stripe
        :empty-text="loading ? '加载中…' : '暂无产品数据'"
        style="width: 100%"
      >
        <el-table-column label="产品UID" min-width="200">
          <template #default="{ row }">
            <el-tooltip :content="row.product_uid" placement="top">
              <code class="uid-cell">{{ truncateHash(row.product_uid, 8) }}</code>
            </el-tooltip>
          </template>
        </el-table-column>
        <el-table-column prop="brand" label="品牌" min-width="120" show-overflow-tooltip />
        <el-table-column prop="model" label="型号" min-width="140" show-overflow-tooltip />
        <el-table-column prop="material" label="底板材质" min-width="120" show-overflow-tooltip />
        <el-table-column prop="batch_no" label="批次号" min-width="140" show-overflow-tooltip />
        <el-table-column prop="production_date" label="生产日期" width="120" />
        <el-table-column label="状态" width="100" align="center">
          <template #default="{ row }">
            <ProductStatusTag :status="row.status" />
          </template>
        </el-table-column>
        <el-table-column label="创建时间" width="170">
          <template #default="{ row }">
            <span class="muted">{{ formatDate(row.created_at) }}</span>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="240" fixed="right" align="center">
          <template #default="{ row }">
            <el-button type="primary" link :icon="View" @click="openTrace(row)">
              查看溯源
            </el-button>
            <el-dropdown
              trigger="click"
              @command="(cmd) => handleStatusChange(row, cmd)"
            >
              <el-button type="warning" link :icon="Edit">
                修改状态
              </el-button>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item
                    v-for="opt in STATUS_OPTIONS"
                    :key="opt.value"
                    :command="opt.value"
                    :disabled="row.status === opt.value"
                  >
                    {{ opt.label }}
                  </el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
      <div class="pagination-wrap">
        <el-pagination
          v-model:current-page="page"
          v-model:page-size="pageSize"
          :total="total"
          :page-sizes="[10, 20, 50, 100]"
          layout="total, sizes, prev, pager, next, jumper"
          background
          @size-change="fetchList"
          @current-change="fetchList"
        />
      </div>
    </el-card>

    <!-- 溯源弹窗 -->
    <el-dialog
      v-model="traceVisible"
      title="产品溯源"
      width="760px"
      :close-on-click-modal="false"
      destroy-on-close
    >
      <div v-loading="traceLoading">
        <el-descriptions
          v-if="currentRow"
          class="trace-info"
          :column="2"
          border
          size="small"
        >
          <el-descriptions-item label="产品UID">
            <code class="hash-text">{{ currentRow.product_uid }}</code>
          </el-descriptions-item>
          <el-descriptions-item label="状态">
            <ProductStatusTag :status="currentRow.status" />
          </el-descriptions-item>
          <el-descriptions-item label="品牌">{{ currentRow.brand || '-' }}</el-descriptions-item>
          <el-descriptions-item label="型号">{{ currentRow.model || '-' }}</el-descriptions-item>
        </el-descriptions>
        <el-divider />
        <TraceTimeline :chain="traceChain" />
      </div>
      <template #footer>
        <el-button @click="traceVisible = false">关闭</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  Plus,
  Search,
  Refresh,
  View,
  Edit
} from '@element-plus/icons-vue'
import { useUserStore } from '@/stores/user'
import { listProducts, updateProduct, getProductTrace } from '@/api/product'
import { formatDate, truncateHash } from '@/utils/format'
import { PRODUCT_STATUS_MAP, translate } from '@/utils/i18n'
import PageHeader from '@/components/PageHeader.vue'
import ProductStatusTag from '@/components/ProductStatusTag.vue'
import TraceTimeline from '@/components/TraceTimeline.vue'

const router = useRouter()
const userStore = useUserStore()

// Status options: backend stores English values; labels are translated from PRODUCT_STATUS_MAP
const STATUS_OPTIONS = Object.keys(PRODUCT_STATUS_MAP).map((key) => ({
  label: PRODUCT_STATUS_MAP[key],
  value: key
}))

// Page title adapts to role
const pageTitle = computed(() => {
  return userStore.role === 'FACTORY' ? '我的产品' : '产品管理'
})
const pageSubtitle = computed(() => {
  return userStore.role === 'FACTORY'
    ? '管理本厂家生产的所有乒乓球拍产品'
    : '查看所有已上链产品，确认售出后可标记为已售出'
})

// 列表
const loading = ref(false)
const list = ref([])
const page = ref(1)
const pageSize = ref(10)
const total = ref(0)

// 筛选
const filterForm = reactive({
  status: '',
  keyword: ''
})

// 客户端关键字过滤(状态由后端过滤)
const filteredList = computed(() => {
  const kw = (filterForm.keyword || '').trim().toLowerCase()
  if (!kw) return list.value
  return list.value.filter((row) => {
    return (
      (row.brand && String(row.brand).toLowerCase().includes(kw)) ||
      (row.model && String(row.model).toLowerCase().includes(kw)) ||
      (row.product_uid && String(row.product_uid).includes(kw))
    )
  })
})

// 统计
const stats = computed(() => {
  const all = list.value
  return {
    total: all.length,
    produced: all.filter((it) => it.status === 'PRODUCED').length,
    inTransit: all.filter((it) => it.status === 'IN_TRANSIT').length,
    sold: all.filter((it) => it.status === 'SOLD').length
  }
})

// 溯源弹窗
const traceVisible = ref(false)
const traceLoading = ref(false)
const currentRow = ref(null)
const traceChain = ref([])

async function fetchList() {
  loading.value = true
  try {
    const params = {
      page: page.value,
      page_size: pageSize.value
    }
    if (filterForm.status) params.status = filterForm.status
    const res = await listProducts(params)
    list.value = res?.items || []
    total.value = res?.total ?? list.value.length
  } catch (err) {
    // request 拦截器已统一提示
  } finally {
    loading.value = false
  }
}

function handleSearch() {
  page.value = 1
  fetchList()
}

function handleReset() {
  filterForm.status = ''
  filterForm.keyword = ''
  page.value = 1
  fetchList()
}

function goCreate() {
  router.push('/factory/create')
}

async function openTrace(row) {
  currentRow.value = row
  traceVisible.value = true
  traceLoading.value = true
  traceChain.value = []
  try {
    const res = await getProductTrace(row.product_uid)
    traceChain.value = res?.trace_chain || []
  } catch (err) {
    // 已统一提示
  } finally {
    traceLoading.value = false
  }
}

async function handleStatusChange(row, newStatus) {
  if (!newStatus || row.status === newStatus) return
  try {
    await ElMessageBox.confirm(
      `确认将产品 ${truncateHash(row.product_uid, 8)} 状态修改为 "${
        STATUS_OPTIONS.find((s) => s.value === newStatus)?.label || newStatus
      }" 吗?`,
      '修改状态',
      {
        type: 'warning',
        confirmButtonText: '确认修改',
        cancelButtonText: '取消'
      }
    )
  } catch (e) {
    return // 用户取消
  }
  try {
    await updateProduct(row.product_uid, { status: newStatus })
    ElMessage.success('状态修改成功')
    fetchList()
  } catch (err) {
    // 拦截器已提示
  }
}

onMounted(() => {
  fetchList()
})
</script>

<style scoped>
.product-list {
  padding: 0 4px;
}

.stat-row {
  margin-bottom: 16px;
}

.stat-card {
  border-radius: 10px;
  text-align: center;
  padding: 4px 0;
}

.stat-card__label {
  font-size: 13px;
  color: #909399;
  margin-bottom: 6px;
}

.stat-card__value {
  font-size: 26px;
  font-weight: 600;
  line-height: 1.2;
}

.stat-card--total .stat-card__value { color: #303133; }
.stat-card--success .stat-card__value { color: #67c23a; }
.stat-card--warning .stat-card__value { color: #e6a23c; }
.stat-card--info .stat-card__value { color: #909399; }

.filter-card,
.table-card {
  border-radius: 10px;
  margin-bottom: 16px;
}

.filter-form {
  margin: 0;
}

.uid-cell {
  font-family: 'Consolas', 'Monaco', 'Courier New', monospace;
  font-size: 12px;
  color: #5b8def;
  background: #f0f5ff;
  padding: 2px 6px;
  border-radius: 4px;
}

.muted {
  color: #909399;
  font-size: 12px;
}

.pagination-wrap {
  display: flex;
  justify-content: flex-end;
  margin-top: 16px;
}

.hash-text {
  font-family: 'Consolas', 'Monaco', 'Courier New', monospace;
  font-size: 12px;
  color: #606266;
  background: #f5f7fa;
  padding: 2px 6px;
  border-radius: 4px;
}

.trace-info {
  margin-bottom: 8px;
}
</style>

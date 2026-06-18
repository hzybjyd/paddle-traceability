<template>
  <div class="verify-page">
    <!-- 顶部标题 -->
    <div class="verify-header">
      <h1 class="verify-title">防伪验证</h1>
      <p class="verify-subtitle">输入产品 UID 查询真伪与完整溯源</p>
    </div>

    <!-- 验证卡片 -->
    <el-card class="verify-card" shadow="hover">
      <!-- 输入区 -->
      <div class="verify-input-row">
        <el-input
          v-model="uid"
          size="large"
          placeholder="请输入产品UID（19位数字）"
          :prefix-icon="Search"
          :maxlength="32"
          clearable
          class="verify-input"
          @keyup.enter="handleVerify"
        />
        <el-button
          type="primary"
          size="large"
          :loading="loading"
          :icon="Search"
          class="verify-btn"
          @click="handleVerify"
        >
          立即验证
        </el-button>
      </div>
      <div class="verify-tip">
        示例: 2066928303960231936
      </div>
    </el-card>

    <!-- 验证中 -->
    <el-card v-if="loading" class="verify-result" shadow="hover">
      <el-skeleton :rows="4" animated />
    </el-card>

    <!-- 验证失败: 假冒产品 -->
    <el-card v-else-if="hasQueried && !result.verified" class="verify-result verify-result--fail" shadow="hover">
      <el-result icon="error" :title="failTitle" :sub-title="failSubtitle">
        <template #icon>
          <div class="fail-icon">
            <el-icon :size="80" color="#f56c6c">
              <CircleClose />
            </el-icon>
          </div>
        </template>
      </el-result>
      <el-alert
        type="error"
        :closable="false"
        show-icon
        title="温馨提示"
        :description="failAlert"
      />
    </el-card>

    <!-- 验证成功: 正品 -->
    <template v-else-if="hasQueried && result.verified && result.data">
      <el-card class="verify-result verify-result--success" shadow="hover">
        <el-result title="✓ 正品" sub-title="区块链溯源验证通过,产品为正品">
          <template #icon>
            <div class="success-icon">
              <el-icon :size="80" color="#67c23a">
                <CircleCheckFilled />
              </el-icon>
            </div>
          </template>
        </el-result>

        <!-- 链验证状态 -->
        <div class="chain-tags">
          <el-tag
            v-if="result.data.chain_verified !== undefined"
            :type="result.data.chain_verified ? 'success' : 'danger'"
            effect="dark"
            size="large"
            round
          >
            <el-icon class="tag-icon"><Check /></el-icon>
            链验证:{{ result.data.chain_verified ? '已验证' : '未通过' }}
          </el-tag>
          <el-tag
            v-if="result.data.data_hash_matched !== undefined"
            :type="result.data.data_hash_matched ? 'success' : 'warning'"
            effect="light"
            size="large"
            round
          >
            <el-icon class="tag-icon"><DocumentChecked /></el-icon>
            数据哈希:{{ result.data.data_hash_matched ? '匹配' : '不匹配' }}
          </el-tag>
        </div>

        <!-- 详细信息 -->
        <el-descriptions
          class="product-info"
          title="产品信息"
          :column="2"
          border
          size="default"
        >
          <el-descriptions-item label="产品 UID">
            <code class="hash-text">{{ result.data.product_uid }}</code>
          </el-descriptions-item>
          <el-descriptions-item label="状态">
            <ProductStatusTag :status="result.data.status" effect="dark" />
          </el-descriptions-item>
          <el-descriptions-item label="品牌">{{ result.data.brand || '-' }}</el-descriptions-item>
          <el-descriptions-item label="型号">{{ result.data.model || '-' }}</el-descriptions-item>
          <el-descriptions-item label="底板材质">{{ result.data.material || '-' }}</el-descriptions-item>
          <el-descriptions-item label="胶皮类型">{{ result.data.rubber_type || '-' }}</el-descriptions-item>
          <el-descriptions-item label="生产日期">{{ result.data.production_date || '-' }}</el-descriptions-item>
          <el-descriptions-item label="状态描述">
            {{ translate(PRODUCT_STATUS_MAP, result.data.status, '-') }}
          </el-descriptions-item>
        </el-descriptions>
      </el-card>

      <!-- 简化版溯源时间线 -->
      <el-card v-if="chainList.length" class="verify-result verify-trace" shadow="hover">
        <template #header>
          <div class="trace-header">
            <el-icon :size="18" color="#5b8def"><Histogram /></el-icon>
            <span class="trace-title">溯源时间线</span>
          </div>
        </template>
        <TraceTimeline :chain="chainList" />
      </el-card>

      <!-- 链上交易哈希(可选) -->
      <el-card v-if="txHash" class="verify-result verify-txhash" shadow="hover">
        <template #header>
          <div class="trace-header">
            <el-icon :size="18" color="#8b6fe8"><Link /></el-icon>
            <span class="trace-title">链上交易哈希</span>
          </div>
        </template>
        <div class="txhash-row">
          <el-tooltip :content="txHash" placement="top">
            <code class="hash-text hash-text--lg">{{ truncateHash(txHash) }}</code>
          </el-tooltip>
          <el-button
            type="primary"
            size="small"
            :icon="CopyDocument"
            @click="handleCopy(txHash)"
          >
            复制
          </el-button>
        </div>
      </el-card>
    </template>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import {
  Search,
  Check,
  CopyDocument,
  Link,
  Histogram,
  DocumentChecked
} from '@element-plus/icons-vue'
import { verifyProduct } from '@/api/verify'
import { truncateHash, copyToClipboard, formatDate } from '@/utils/format'
import { PRODUCT_STATUS_MAP, translate, translateApiMessage } from '@/utils/i18n'
import ProductStatusTag from '@/components/ProductStatusTag.vue'
import TraceTimeline from '@/components/TraceTimeline.vue'

const route = useRoute()

const uid = ref('')
const loading = ref(false)
const hasQueried = ref(false)
const result = ref({ verified: false, data: null, message: '' })

// Failure state texts
const failTitle = computed(() => {
  return result.value.message ? '未查询到该产品记录' : '验证失败'
})
const failSubtitle = computed(() => {
  return result.value.message
    ? '该产品未在区块链上查询到溯源记录,可能为假冒产品'
    : ''
})
const failAlert = computed(() => {
  return '该产品未在区块链上查询到溯源记录,可能为假冒产品。购买时请认准官方渠道,谨防上当受骗。'
})

// Convert backend trace_summary into the format expected by TraceTimeline.
// The backend returns English stage values; the timeline component translates them to Chinese.
const chainList = computed(() => {
  const list = result.value.data?.trace_summary
  if (!Array.isArray(list)) return []
  return list.map((item) => ({
    stage: item.stage,
    operator: item.operator,
    timestamp: item.time || item.timestamp,
    chain_status: result.value.data?.chain_verified ? 'CONFIRMED' : 'PENDING',
    data_hash: result.value.data?.data_hash || '',
    tx_hash: item.tx_hash || '',
    detail: item.detail || null
  }))
})

// On-chain transaction hash: prefer the first trace entry with a hash, otherwise from the top-level data
const txHash = computed(() => {
  const firstWithTx = chainList.value.find((it) => it.tx_hash)
  return (
    firstWithTx?.tx_hash ||
    result.value.data?.tx_hash ||
    ''
  )
})

async function handleVerify() {
  const value = (uid.value || '').trim()
  if (!value) {
    ElMessage.warning('请输入产品 UID')
    return
  }
  loading.value = true
  hasQueried.value = false
  try {
    const res = await verifyProduct(value)
    result.value = {
      verified: !!res.verified,
      data: res.data || null,
      message: translateApiMessage(res.message)
    }
    hasQueried.value = true
  } catch (err) {
    hasQueried.value = false
  } finally {
    loading.value = false
  }
}

async function handleCopy(text) {
  const ok = await copyToClipboard(text)
  if (ok) {
    ElMessage.success('已复制到剪贴板')
  } else {
    ElMessage.error('复制失败,请手动复制')
  }
}

onMounted(() => {
  // 如果 URL 带 :uid 参数,自动填入并触发验证
  const uidParam = route.params.uid
  if (uidParam) {
    uid.value = String(uidParam)
    handleVerify()
  }
})
</script>

<style scoped>
.verify-page {
  max-width: 880px;
  margin: 0 auto;
  padding: 32px 16px 48px;
}

.verify-header {
  text-align: center;
  margin-bottom: 28px;
}

.verify-title {
  font-size: 32px;
  font-weight: 700;
  color: #303133;
  margin: 0 0 8px;
  letter-spacing: 2px;
}

.verify-subtitle {
  font-size: 14px;
  color: #909399;
  margin: 0;
}

.verify-card {
  border-radius: 12px;
  margin-bottom: 20px;
  background: linear-gradient(135deg, #ffffff 0%, #f5f8ff 100%);
}

.verify-input-row {
  display: flex;
  gap: 12px;
  align-items: center;
}

.verify-input {
  flex: 1;
}

.verify-input :deep(.el-input__wrapper) {
  padding: 4px 12px;
  border-radius: 8px;
  box-shadow: 0 0 0 1px #dcdfe6 inset;
}

.verify-btn {
  min-width: 140px;
  height: 44px;
  font-size: 15px;
  letter-spacing: 2px;
  background: linear-gradient(135deg, #5b8def 0%, #8b6fe8 100%);
  border: none;
}

.verify-btn:hover {
  background: linear-gradient(135deg, #4a7ade 0%, #7a5ed8 100%);
}

.verify-tip {
  margin-top: 10px;
  font-size: 12px;
  color: #c0c4cc;
  text-align: left;
}

.verify-result {
  border-radius: 12px;
  margin-bottom: 20px;
}

.verify-result--fail :deep(.el-result__title) {
  color: #f56c6c;
}

.verify-result--success :deep(.el-result__title) {
  color: #67c23a;
  font-size: 28px;
  font-weight: 700;
}

.fail-icon,
.success-icon {
  display: flex;
  justify-content: center;
  align-items: center;
  padding: 16px 0;
}

.chain-tags {
  display: flex;
  gap: 12px;
  justify-content: center;
  margin: 0 0 24px;
  flex-wrap: wrap;
}

.chain-tags .el-tag {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 0 14px;
  height: 32px;
  font-size: 14px;
}

.tag-icon {
  margin-right: 2px;
}

.product-info {
  margin-top: 8px;
}

.hash-text {
  font-family: 'Consolas', 'Monaco', 'Courier New', monospace;
  font-size: 13px;
  color: #606266;
  background: #f5f7fa;
  padding: 2px 8px;
  border-radius: 4px;
  word-break: break-all;
}

.hash-text--lg {
  font-size: 14px;
  padding: 6px 12px;
}

.trace-header {
  display: flex;
  align-items: center;
  gap: 8px;
}

.trace-title {
  font-size: 16px;
  font-weight: 600;
  color: #303133;
}

.txhash-row {
  display: flex;
  align-items: center;
  gap: 12px;
  flex-wrap: wrap;
}

@media (max-width: 600px) {
  .verify-input-row {
    flex-direction: column;
    align-items: stretch;
  }
  .verify-btn {
    width: 100%;
  }
  .verify-title {
    font-size: 24px;
  }
}
</style>

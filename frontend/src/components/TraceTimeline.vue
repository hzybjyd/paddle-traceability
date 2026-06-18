<template>
  <div class="trace-timeline">
    <el-empty v-if="!chain || chain.length === 0" description="暂无溯源记录" />
    <el-timeline v-else>
      <el-timeline-item
        v-for="(item, index) in chain"
        :key="index"
        :timestamp="formatDate(item.timestamp)"
        placement="top"
        :hollow="item.chain_status !== 'CONFIRMED'"
        :type="timelineType(item.chain_status)"
        size="large"
      >
        <el-card shadow="hover" class="trace-card">
          <div class="trace-card__header">
            <span class="trace-card__stage">{{ translateStage(item.stage) }}</span>
            <el-tag
              :type="chainStatusType(item.chain_status)"
              size="small"
              effect="light"
              round
            >
              {{ chainStatusLabel(item.chain_status) }}
            </el-tag>
          </div>

          <el-descriptions :column="1" size="small" class="trace-card__desc">
            <el-descriptions-item label="操作方">
              {{ item.operator || '-' }}
            </el-descriptions-item>
            <el-descriptions-item label="数据哈希">
              <el-tooltip :content="item.data_hash || ''" placement="top">
                <code class="hash-text">{{ truncateHash(item.data_hash) }}</code>
              </el-tooltip>
            </el-descriptions-item>
            <el-descriptions-item label="链上交易哈希">
              <el-tooltip :content="item.tx_hash || ''" placement="top">
                <code class="hash-text">{{ truncateHash(item.tx_hash) }}</code>
              </el-tooltip>
            </el-descriptions-item>
            <el-descriptions-item v-if="item.detail && item.detail.action" label="操作类型">
              {{ translateLogisticsAction(item.detail.action) }}
            </el-descriptions-item>
          </el-descriptions>

          <div v-if="item.detail" class="trace-card__detail">
            <el-divider class="detail-divider" />
            <pre class="detail-json">{{ JSON.stringify(item.detail, null, 2) }}</pre>
          </div>
        </el-card>
      </el-timeline-item>
    </el-timeline>
  </div>
</template>

<script setup>
import { formatDate, truncateHash } from '@/utils/format'
import { STAGE_MAP, CHAIN_STATUS_MAP, LOGISTICS_ACTION_MAP, translate } from '@/utils/i18n'

defineProps({
  chain: {
    type: Array,
    default: () => []
  }
})

// Chain status -> tag color
function chainStatusType(status) {
  switch (status) {
    case 'CONFIRMED':
      return 'success'
    case 'PENDING':
      return 'warning'
    case 'FAILED':
      return 'danger'
    default:
      return 'info'
  }
}

// Chain status -> Chinese label
function chainStatusLabel(status) {
  return translate(CHAIN_STATUS_MAP, status, status || '未知')
}

// Timeline node type (affects left dot color)
function timelineType(status) {
  switch (status) {
    case 'CONFIRMED':
      return 'success'
    case 'PENDING':
      return 'warning'
    case 'FAILED':
      return 'danger'
    default:
      return 'primary'
  }
}

// Stage -> Chinese label
function translateStage(stage) {
  return translate(STAGE_MAP, stage, stage || '-')
}

// Logistics action -> Chinese label
function translateLogisticsAction(action) {
  return translate(LOGISTICS_ACTION_MAP, action, action || '-')
}
</script>

<style scoped>
.trace-timeline {
  width: 100%;
}

.trace-card {
  border-radius: 8px;
}

.trace-card__header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 12px;
}

.trace-card__stage {
  font-size: 16px;
  font-weight: 600;
  color: #303133;
}

.trace-card__desc {
  margin-top: 4px;
}

.hash-text {
  font-family: 'Consolas', 'Monaco', 'Courier New', monospace;
  font-size: 12px;
  color: #606266;
  background: #f5f7fa;
  padding: 2px 6px;
  border-radius: 4px;
}

.detail-divider {
  margin: 8px 0;
}

.detail-json {
  font-family: 'Consolas', 'Monaco', 'Courier New', monospace;
  font-size: 12px;
  color: #303133;
  background: #fafafa;
  border: 1px solid #ebeef5;
  border-radius: 4px;
  padding: 8px 12px;
  margin: 0;
  white-space: pre-wrap;
  word-break: break-all;
  max-height: 240px;
  overflow: auto;
}
</style>

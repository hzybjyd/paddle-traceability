<template>
  <div class="product-create">
    <PageHeader title="创建产品" subtitle="填写产品信息,提交后将自动上链存证">
      <el-button :icon="ArrowLeft" @click="goBack">返回列表</el-button>
    </PageHeader>

    <!-- 成功结果页 -->
    <el-card
      v-if="result"
      shadow="hover"
      class="result-card"
    >
      <el-result icon="success" title="产品创建成功,已上链存证">
        <template #icon>
          <div class="result-icon">
            <el-icon :size="72" color="#67c23a"><CircleCheckFilled /></el-icon>
          </div>
        </template>
        <template #sub-title>
          <div class="result-sub">请妥善保存产品 UID,可用于后续溯源查询与流转登记</div>
        </template>
      </el-result>

      <el-descriptions class="result-info" :column="1" border size="default">
        <el-descriptions-item label="产品 UID">
          <span class="uid-text">{{ result.product_uid }}</span>
          <el-button
            type="primary"
            link
            :icon="CopyDocument"
            class="copy-btn"
            @click="handleCopy(result.product_uid, '产品UID')"
          >
            复制
          </el-button>
        </el-descriptions-item>
        <el-descriptions-item label="链上交易哈希">
          <code class="hash-text">{{ result.tx_hash }}</code>
          <el-button
            type="primary"
            link
            :icon="CopyDocument"
            class="copy-btn"
            @click="handleCopy(result.tx_hash, '交易哈希')"
          >
            复制
          </el-button>
        </el-descriptions-item>
        <el-descriptions-item label="区块高度">
          <el-tag effect="plain" round>#{{ result.block_height }}</el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="创建时间">
          {{ formatDate(result.created_at) }}
        </el-descriptions-item>
      </el-descriptions>

      <div class="result-actions">
        <el-button type="primary" :icon="List" @click="goBack">返回列表</el-button>
        <el-button :icon="Plus" @click="handleContinue">继续创建</el-button>
      </div>
    </el-card>

    <!-- 创建表单 -->
    <el-row v-else :gutter="16">
      <el-col :xs="24" :md="14">
        <el-card shadow="hover" class="form-card">
          <template #header>
            <div class="card-header">
              <el-icon :size="18" color="#5b8def"><DocumentAdd /></el-icon>
              <span class="card-title">产品信息</span>
            </div>
          </template>
          <el-form
            ref="formRef"
            :model="form"
            :rules="rules"
            label-width="100px"
            label-position="right"
            size="default"
            @submit.prevent
          >
            <el-form-item label="品牌" prop="brand">
              <el-input
                v-model="form.brand"
                placeholder="请输入品牌,如:红双喜"
                maxlength="50"
                clearable
              />
            </el-form-item>
            <el-form-item label="型号" prop="model">
              <el-input
                v-model="form.model"
                placeholder="请输入产品型号"
                maxlength="100"
                clearable
              />
            </el-form-item>
            <el-form-item label="底板材质" prop="material">
              <el-input
                v-model="form.material"
                placeholder="如:碳素纤维 / 纯木 / 芳碳"
                maxlength="100"
                clearable
              />
            </el-form-item>
            <el-form-item label="胶皮类型" prop="rubber_type">
              <el-input
                v-model="form.rubber_type"
                placeholder="如:反胶 / 正胶 / 长胶"
                maxlength="100"
                clearable
              />
            </el-form-item>
            <el-form-item label="批次号" prop="batch_no">
              <el-input
                v-model="form.batch_no"
                placeholder="如:20260617001"
                maxlength="50"
                clearable
              />
            </el-form-item>
            <el-form-item label="生产日期" prop="production_date">
              <el-date-picker
                v-model="form.production_date"
                type="date"
                placeholder="请选择生产日期"
                value-format="YYYY-MM-DD"
                style="width: 100%"
              />
            </el-form-item>
            <el-form-item label="质检报告哈希" prop="quality_report_hash">
              <el-input
                v-model="form.quality_report_hash"
                placeholder="可选,留空则自动生成"
                maxlength="64"
                clearable
              />
            </el-form-item>
            <el-form-item>
              <el-button
                type="primary"
                size="large"
                :loading="submitting"
                :icon="Position"
                class="submit-btn"
                @click="handleSubmit"
              >
                提交并上链
              </el-button>
              <el-button size="large" :icon="RefreshLeft" @click="handleReset">重置</el-button>
            </el-form-item>
          </el-form>
        </el-card>
      </el-col>

      <el-col :xs="24" :md="10">
        <el-card shadow="hover" class="tip-card">
          <template #header>
            <div class="card-header">
              <el-icon :size="18" color="#e6a23c"><InfoFilled /></el-icon>
              <span class="card-title">设计说明</span>
            </div>
          </template>
          <ul class="tip-list">
            <li>
              <strong>产品 UID</strong>
              <p>系统将使用雪花算法自动生成 19 位数字 ID,作为产品的唯一标识。</p>
            </li>
            <li>
              <strong>上链存证</strong>
              <p>提交后产品信息将写入百度超级链,生成不可篡改的交易哈希。</p>
            </li>
            <li>
              <strong>质检报告哈希</strong>
              <p>留空时系统自动生成 SHA-256 哈希;也可由质检系统预先计算后填入。</p>
            </li>
            <li>
              <strong>生产日期</strong>
              <p>必须为有效日期,系统会将其写入链上存证数据,影响后续溯源展示。</p>
            </li>
            <li>
              <strong>字段说明</strong>
              <p>带 <span class="required">*</span> 的字段为必填,其他字段建议尽量填写完整。</p>
            </li>
          </ul>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import {
  ArrowLeft,
  DocumentAdd,
  CircleCheckFilled,
  CopyDocument,
  List,
  Plus,
  Position,
  RefreshLeft,
  InfoFilled
} from '@element-plus/icons-vue'
import { useUserStore } from '@/stores/user'
import { createProduct } from '@/api/product'
import { formatDate, copyToClipboard } from '@/utils/format'
import PageHeader from '@/components/PageHeader.vue'

const router = useRouter()
const userStore = useUserStore()

// 提交结果
const result = ref(null)
const submitting = ref(false)

// 表单
const formRef = ref(null)
const form = reactive({
  brand: '',
  model: '',
  material: '',
  rubber_type: '',
  batch_no: '',
  production_date: '',
  quality_report_hash: ''
})

const rules = {
  brand: [{ required: true, message: '请输入品牌', trigger: 'blur' }],
  model: [{ required: true, message: '请输入产品型号', trigger: 'blur' }],
  material: [{ required: true, message: '请输入底板材质', trigger: 'blur' }],
  batch_no: [{ required: true, message: '请输入生产批次号', trigger: 'blur' }],
  production_date: [
    { required: true, message: '请选择生产日期', trigger: 'change' }
  ]
}

function goBack() {
  router.push('/factory')
}

function handleReset() {
  formRef.value?.resetFields()
}

function handleContinue() {
  result.value = null
  handleReset()
}

async function handleCopy(text, label) {
  if (!text) return
  const ok = await copyToClipboard(text)
  if (ok) {
    ElMessage.success(`${label} 已复制到剪贴板`)
  } else {
    ElMessage.error('复制失败,请手动复制')
  }
}

async function handleSubmit() {
  if (!formRef.value) return
  await formRef.value.validate(async (valid) => {
    if (!valid) return
    submitting.value = true
    try {
      // 清理空字符串字段,避免后端对空值校验失败
      const payload = {}
      Object.keys(form).forEach((k) => {
        const v = form[k]
        if (v !== '' && v !== null && v !== undefined) {
          payload[k] = v
        }
      })
      const res = await createProduct(payload)
      result.value = res
      ElMessage.success('产品创建成功')
    } catch (err) {
      // 错误提示已由 request 拦截器统一处理
    } finally {
      submitting.value = false
    }
  })
}
</script>

<style scoped>
.product-create {
  padding: 0 4px;
}

.form-card,
.tip-card,
.result-card {
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

.submit-btn {
  min-width: 160px;
  background: linear-gradient(135deg, #5b8def 0%, #8b6fe8 100%);
  border: none;
}

.submit-btn:hover {
  background: linear-gradient(135deg, #4a7ade 0%, #7a5ed8 100%);
}

.tip-list {
  list-style: none;
  padding: 0;
  margin: 0;
}

.tip-list li {
  padding: 10px 0;
  border-bottom: 1px dashed #ebeef5;
  font-size: 13px;
  color: #606266;
  line-height: 1.6;
}

.tip-list li:last-child {
  border-bottom: none;
}

.tip-list strong {
  display: block;
  color: #303133;
  font-size: 14px;
  margin-bottom: 4px;
}

.tip-list p {
  margin: 0;
  color: #909399;
}

.required {
  color: #f56c6c;
  margin: 0 2px;
}

.result-card {
  text-align: center;
  padding: 8px 0;
}

.result-icon {
  display: flex;
  justify-content: center;
  align-items: center;
  padding: 12px 0 4px;
}

.result-sub {
  color: #909399;
  font-size: 13px;
  margin-top: 4px;
}

.result-info {
  max-width: 720px;
  margin: 16px auto 8px;
  text-align: left;
}

.uid-text {
  font-family: 'Consolas', 'Monaco', 'Courier New', monospace;
  font-size: 16px;
  font-weight: 600;
  color: #5b8def;
  background: #f0f5ff;
  padding: 4px 10px;
  border-radius: 6px;
  letter-spacing: 1px;
}

.hash-text {
  font-family: 'Consolas', 'Monaco', 'Courier New', monospace;
  font-size: 12px;
  color: #606266;
  background: #f5f7fa;
  padding: 2px 8px;
  border-radius: 4px;
  word-break: break-all;
  display: inline-block;
  max-width: 100%;
}

.copy-btn {
  margin-left: 8px;
}

.result-actions {
  margin-top: 20px;
  display: flex;
  justify-content: center;
  gap: 12px;
}
</style>

<template>
  <div class="page-header">
    <div class="page-header__left">
      <div class="page-header__title">{{ title }}</div>
      <div v-if="subtitle" class="page-header__subtitle">{{ subtitle }}</div>
    </div>
    <div class="page-header__right">
      <slot />
    </div>
  </div>
  <el-divider class="page-header__divider" />
</template>

<script setup>
defineProps({
  title: {
    type: String,
    required: true,
    default: ''
  },
  subtitle: {
    type: String,
    default: ''
  }
})
</script>

<style scoped>
.page-header {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 16px;
  padding: 4px 0 0 0;
}

.page-header__left {
  flex: 1;
  min-width: 0;
}

.page-header__title {
  font-size: 20px;
  font-weight: 600;
  color: #303133;
  line-height: 1.4;
}

.page-header__subtitle {
  margin-top: 4px;
  font-size: 13px;
  color: #909399;
  line-height: 1.5;
}

.page-header__right {
  flex-shrink: 0;
  display: flex;
  align-items: center;
  gap: 8px;
}

.page-header__divider {
  margin: 16px 0 20px 0;
}
</style>

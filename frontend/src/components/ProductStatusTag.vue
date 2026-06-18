<template>
  <el-tag :type="tagType" :effect="effect" round>
    {{ label }}
  </el-tag>
</template>

<script setup>
import { computed } from 'vue'
import { PRODUCT_STATUS_MAP, translate } from '@/utils/i18n'

const props = defineProps({
  status: {
    type: String,
    default: ''
  },
  effect: {
    type: String,
    default: 'light'
  }
})

// Status -> { Chinese label, Element Plus tag type }
// Backend stores English values; look up the Chinese label here.
const STATUS_TYPES = {
  PRODUCED: 'success',
  IN_TRANSIT: 'warning',
  IN_STOCK: 'primary',
  SOLD: 'info'
}

const label = computed(() => {
  return translate(PRODUCT_STATUS_MAP, props.status, props.status ?? '')
})

const tagType = computed(() => {
  return STATUS_TYPES[props.status] ?? 'info'
})
</script>

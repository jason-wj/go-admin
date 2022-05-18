<template>
  <div :class="{'hidden':hidden}" class="pagination-container">
    <el-pagination
      v-model:current-page="currentPage"
      v-model:page-size="pageSize"
      :background="background"
      :layout="layout"
      :page-sizes="pageSizes"
      :total="total"
      v-bind="$attrs"
      @size-change="handleSizeChange"
      @current-change="handleCurrentChange"
    />
  </div>
</template>

<script>
export default {
  name: 'Pagination',
}
</script>

<script setup>
import { scrollTo } from '@/utils/scroll-to'
import { computed, getCurrentInstance } from 'vue'

const context = getCurrentInstance()
const props = defineProps({
  total: {
    required: true,
    type: Number
  },
  page: {
    type: Number,
    default: 1
  },
  limit: {
    type: Number,
    default: 20
  },
  pageSizes: {
    type: Array,
    default() {
      return [10, 20, 30, 50]
    }
  },
  layout: {
    type: String,
    default: 'total, sizes, prev, pager, next, jumper'
  },
  background: {
    type: Boolean,
    default: true
  },
  autoScroll: {
    type: Boolean,
    default: false
  },
  hidden: {
    type: Boolean,
    default: false
  }
})

const currentPage = computed({
  get() {
    return props.page
  },
  set(val) {
    context.emit('update:page', val)
  }
})

const pageSize = computed({
  get() {
    return props.limit
  },
  set(val) {
    context.emit('update:limit', val)
  }
})

const handleSizeChange = (val) => {
  context.emit('pagination', { page: currentPage, limit: val })
  if (props.autoScroll) {
    scrollTo(0, 800)
  }
}
const handleCurrentChange = (val) => {
  context.emit('pagination', { page: val, limit: pageSize })
  if (props.autoScroll) {
    scrollTo(0, 800)
  }
}
</script>

<style scoped>
.pagination-container {
  background: #fff;
  padding: 32px 16px;
}
.pagination-container.hidden {
  display: none;
}
</style>

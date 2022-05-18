<template>
  <el-dialog v-model="open" :close-on-click-modal="false" title="选择文章分类" width="800px" top="5vh" append-to-body destroy-on-close>
    <el-form ref="searchForm" :model="queryParams" :inline="true" label-width="68px">
      <el-form-item label-width="100" label="分类编号" prop="id">
        <el-input
          v-model="queryParams.id"
          placeholder="请输入分类编号"
          clearable
          size="small"
          @keyup.enter.native="handleQuery"
        />
      </el-form-item>
      <el-form-item label-width="100" label="分类名称" prop="name">
        <el-input
          v-model="queryParams.name"
          placeholder="请输入分类名称"
          clearable
          size="small"
          @keyup.enter.native="handleQuery"
        />
      </el-form-item>

      <el-form-item>
        <el-button type="primary" icon="search" size="small" @click="handleQuery">搜索</el-button>
        <el-button icon="refresh" size="small" @click="resetQuery">重置</el-button>
      </el-form-item>
    </el-form>

    <el-table v-loading="loading" stripe border :data="categoryList" height="260" @row-click="clickRow">
      <el-table-column label="序号" type="index" align="center" width="80">
        <template #default="scope">
          <span>{{ (queryParams.pageIndex - 1) * queryParams.pageSize + scope.$index + 1 }}</span>
        </template>
      </el-table-column>
      <el-table-column width="100" label="分类编号" align="center" prop="id" :show-overflow-tooltip="true" />
      <el-table-column width="200" label="分类名称" align="center" prop="name" :show-overflow-tooltip="true" />
      <el-table-column width="100" label="更新人编号" align="center" prop="updateBy" :show-overflow-tooltip="true" />
      <el-table-column label="创建时间" align="center" prop="createdAt" :show-overflow-tooltip="true">
        <template #default="scope">
          <span>{{ formatDate(scope.row.createdAt) }}</span>
        </template>
      </el-table-column>
    </el-table>

    <Pagination
      v-show="total>0"
      v-model:page="queryParams.pageIndex"
      v-model:limit="queryParams.pageSize"
      :total="total"
      @pagination="getList"
    />
  </el-dialog>
</template>

<script>
export default {
  name: 'SelectCategory',
}
</script>

<script setup>

import { listCategory } from '@/api/plugins/content/category'
import { ref } from 'vue'
const emit = defineEmits(['ok'])

// 遮罩层
const loading = ref(true)
// 总条数
const total = ref(0)
// 是否显示弹出层
const open = ref(false)
// 类型数据字典
const categoryList = ref([])

const searchForm = ref(null)

// 查询参数
const queryParams = ref(
  {
    pageIndex: 1,
    pageSize: 10,
    name: undefined
  })

/** 查询参数列表 */
const getList = () => {
  loading.value = true
  listCategory(queryParams.value)
    .then(response => {
      categoryList.value = response.data.list
      total.value = response.data.count
      loading.value = false
    })
}

// 显示弹框
const show = () => {
  queryParams.value = {
    pageIndex: 1,
    pageSize: 20,
    name: undefined
  }
  getList()
  open.value = true
}

const clickRow = (row) => {
  open.value = false
  emit('ok', row)
}

/** 搜索按钮操作 */
const handleQuery = () => {
  queryParams.pageIndex = 1
  getList()
}

/** 重置按钮操作 */
const resetQuery = () => {
  searchForm.value.resetFields()
  handleQuery()
}

getList()

defineExpose({ show })
</script>

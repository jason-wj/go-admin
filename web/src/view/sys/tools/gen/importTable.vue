<template>
  <!-- 导入表 -->
  <el-dialog v-model="visible" title="导入表" width="800px" top="5vh" append-to-body destroy-on-close>
    <el-form ref="searchForm" :model="queryParams" :inline="true">
      <el-form-item label="表名称" prop="tableName">
        <el-input
          v-model="queryParams.tableName"
          placeholder="请输入表名称"
          clearable
          size="small"
          @keyup.enter.native="handleQuery"
        />
      </el-form-item>
      <!--      <el-form-item label="表描述" prop="tableComment">
        <el-input
          v-model="queryParams.tableComment"
          placeholder="请输入表描述"
          clearable
          size="small"
          @keyup.enter.native="handleQuery"
        />
      </el-form-item>-->
      <el-form-item>
        <el-button type="primary" icon="search" size="small" @click="handleQuery">搜索</el-button>
        <el-button icon="refresh" size="small" @click="resetQuery">重置</el-button>
      </el-form-item>
    </el-form>
    <el-row>
      <el-table
        ref="table"
        :data="dbTableList"
        height="260px"
        @row-click="clickRow"
        @selection-change="handleSelectionChange"
      >
        <el-table-column type="selection" width="55" />
        <el-table-column prop="tableName" label="表名称" />
        <el-table-column prop="tableComment" label="表描述" />
        <el-table-column prop="createTime" label="创建时间">
          <template #default="scope">
            <span>{{ formatDate(scope.row.createTime) }}</span>
          </template>
        </el-table-column>
        <!--        <el-table-column prop="updateTime" label="更新时间" />-->
      </el-table>
      <pagination
        v-show="total>0"
        v-model:page="queryParams.pageIndex"
        v-model:limit="queryParams.pageSize"
        :total="total"
        @pagination="getList"
      />
    </el-row>
    <div slot="footer" class="dialog-footer">
      <el-button type="primary" :loading="loading" @click="handleImportTable">确 定</el-button>
      <el-button @click="visible = false">取 消</el-button>
    </div>
  </el-dialog>
</template>

<script setup>
import { listDbTable, importTable } from '@/api/sys/tools/gen'
import { getCurrentInstance, ref } from 'vue'
import { ElMessage } from 'element-plus'

// 遮罩层
const loading = ref(false)
// 遮罩层
const visible = ref(false)
// 选中数组值
const tables = ref([])
// 总条数
const total = ref(0)
// 选中数组值
const dbTableList = ref([])

// ref form
const searchForm = ref(null)

// 查询参数
const queryParams = ref({
  pageIndex: 1,
  pageSize: 10,
  tableName: undefined,
  tableComment: undefined
})

const { refs } = getCurrentInstance()
const emit = defineEmits(['ok'])

// 显示弹框
const show = () => {
  getList()
  visible.value = true
}

const clickRow = (row) => {
  refs.table.toggleRowSelection(row)
}

// 多选框选中数据
const handleSelectionChange = (selection) => {
  tables.value = selection.map(item => item.tableName)
}

// 查询表数据
const getList = async() => {
  const res = await listDbTable(queryParams.value)
  if (res.code === 200) {
    dbTableList.value = res.data.list
    total.value = res.data.count
  }
}

/** 搜索按钮操作 */
const handleQuery = () => {
  queryParams.value.pageIndex = 1
  getList()
}

/** 重置按钮操作 */
const resetQuery = () => {
  searchForm.value.resetFields()
  handleQuery()
}

/** 导入按钮操作 */
const handleImportTable = async() => {
  loading.value = true
  const res = await importTable({ tables: tables.value.join(',') })
  let resType = 'error'
  if (res.code === 200) {
    visible.value = false
    resType = 'success'
    emit('ok')
  }
  loading.value = false
  ElMessage({
    type: resType,
    message: res.msg
  })
}

defineExpose({ show })
</script>

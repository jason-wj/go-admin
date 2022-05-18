
<template>
  <div>
    <div class="gva-search-box">

      <el-form ref="searchForm" :model="queryParams" :inline="true" label-position="left">
        <el-form-item label="主键编码" prop="id">
          <el-input v-model="queryParams.id" placeholder="请输入主键编码" clearable size="small" @keyup.enter.native="handleQuery" />
        </el-form-item>
        <el-form-item label="用户名" prop="username">
          <el-input v-model="queryParams.username" placeholder="请输入用户名" clearable size="small" @keyup.enter.native="handleQuery" />
        </el-form-item>
        <el-form-item label="状态" prop="status">
          <el-select v-model="queryParams.status" placeholder="登录日志状态" clearable size="small">
            <el-option
              v-for="dict in statusOptions"
              :key="dict.dictValue"
              :label="dict.dictLabel"
              :value="dict.dictValue"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="ip地址" prop="ipaddr">
          <el-input v-model="queryParams.ipaddr" placeholder="请输入ip地址" clearable size="small" @keyup.enter.native="handleQuery" />
        </el-form-item>
        <el-form-item label="创建时间" prop="createdAt">
          <el-date-picker v-model="dateRange" size="small" type="datetimerange" range-separator="至" start-placeholder="开始日期" end-placeholder="结束日期" align="right" value-format="yyyy-MM-dd HH:mm:ss" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" icon="search" size="small" @click="handleQuery">搜索</el-button>
          <el-button icon="refresh" size="small" @click="resetQuery">重置</el-button>
        </el-form-item>
      </el-form>

      <el-row :gutter="10" class="mb8">
        <el-col :span="1.5">
          <el-button type="success" icon="download" size="small" @click="handleExport">Excel导出</el-button>
        </el-col>
      </el-row>
    </div>
    <el-table v-loading="loading" stripe border :data="loginlogList">
      <el-table-column label="序号" type="index" align="center" width="60">
        <template #default="scope">
          <span>{{ (queryParams.pageIndex - 1) * queryParams.pageSize + scope.$index + 1 }}</span>
        </template>
      </el-table-column>
      <el-table-column min-width="100" label="主键编码" align="center" prop="id" :show-overflow-tooltip="true" />
      <el-table-column min-width="100" label="用户名" align="center" prop="username" :show-overflow-tooltip="true" />
      <el-table-column min-width="100" label="状态" align="center" prop="status" :formatter="statusFormat">
        <template #default="scope">
          {{ statusFormat(scope.row) }}
        </template>
      </el-table-column>
      <el-table-column min-width="100" label="信息" align="center" prop="msg" :show-overflow-tooltip="true" />
      <el-table-column min-width="180" label="ip地址" align="center" prop="ipaddr" :show-overflow-tooltip="true" />
      <el-table-column min-width="100" label="归属地" align="center" prop="loginLocation" :show-overflow-tooltip="true" />
      <el-table-column min-width="180" label="浏览器" align="center" prop="browser" :show-overflow-tooltip="true" />
      <el-table-column min-width="180" label="系统" align="center" prop="os" :show-overflow-tooltip="true" />
      <el-table-column min-width="100" label="固件" align="center" prop="platform" :show-overflow-tooltip="true" />
      <el-table-column min-width="180" label="登录时间" align="center" prop="loginTime" :show-overflow-tooltip="true">
        <template #default="scope">
          <span>{{ formatDate(scope.row.loginTime) }}</span>
        </template>
      </el-table-column>
      <el-table-column min-width="180" label="备注" align="center" prop="remark" :show-overflow-tooltip="true" />
      <el-table-column min-width="180" label="创建时间" align="center" prop="createdAt" :show-overflow-tooltip="true">
        <template #default="scope">
          <span>{{ formatDate(scope.row.createdAt) }}</span>
        </template>
      </el-table-column>
      <el-table-column min-width="180" label="最后更新时间" align="center" prop="updatedAt" :show-overflow-tooltip="true">
        <template #default="scope">
          <span>{{ formatDate(scope.row.updatedAt) }}</span>
        </template>
      </el-table-column>
      <!--      <el-table-column min-width="100" label="创建者" align="center" prop="createBy" :show-overflow-tooltip="true" />-->
      <!--      <el-table-column min-width="100" label="更新者" align="center" prop="updateBy" :show-overflow-tooltip="true" />-->
      <el-table-column min-width="130" fixed="right" label="操作" align="center" class-name="small-padding fixed-width">
        <template #default="scope">
          <el-button size="small" type="text" icon="delete" @click="handleDelete(scope.row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>
    <Pagination v-show="total>0" v-model:page="queryParams.pageIndex" v-model:limit="queryParams.pageSize" :total="total" @pagination="getList" />
  </div>
</template>

<script>
export default {
  name: 'LoginLog',
}
</script>

<script setup>
import { getCurrentInstance, ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { resolveBlob } from '@/utils/zipdownload'
import { exportLoginLog, delLoginLog, listLoginLog } from '@/api/sys/loginlog'

const { proxy } = getCurrentInstance()

// 遮罩层
const loading = ref(true)
// 总条数
const total = ref(0)
// 是否显示弹出层
const open = ref(false)
// 日期范围
const dateRange = ref([])
// 数据列表
const loginlogList = ref([])
// 类型数据字典
const statusOptions = ref([])
// 查询参数
const queryParams = ref({
  pageIndex: 1,
  pageSize: 10,
  id: undefined,
  username: undefined,
  status: undefined,
  ipaddr: undefined,
})

const searchForm = ref(null)

const init = () => {
  getList()
  proxy.getDicts('sys_status').then(response => {
    statusOptions.value = response.data
  })
}

/** 查询参数列表 */
const getList = () => {
  loading.value = true
  listLoginLog(proxy.addDateRange(queryParams.value, dateRange.value)).then(response => {
    loginlogList.value = response.data.list
    total.value = response.data.count
    loading.value = false
  })
}

// 字典
const statusFormat = (row) => {
  return proxy.selectDictLabel(statusOptions.value, row.status)
}

/** 搜索按钮操作 */
const handleQuery = () => {
  queryParams.value.pageIndex = 1
  getList()
}

/** 重置按钮操作 */
const resetQuery = () => {
  dateRange.value = []
  searchForm.value.resetFields()
  handleQuery()
}

/** 删除按钮操作 */
const handleDelete = (row) => {
  const ids = [row.id]

  ElMessageBox.confirm('是否确认删除编号为"' + ids + '"的数据项?', '警告', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(() => {
    return delLoginLog({ 'ids': ids })
  }).then((response) => {
    if (response.code === 200) {
      open.value = false
      getList()
      ElMessage({
        type: 'success',
        message: response.msg,
        showClose: true,
      })
    }
  }).catch(() => {
    ElMessage({
      type: 'info',
      message: '取消操作'
    })
  })
}

/** 下载excel */
const handleExport = () => {
  ElMessageBox.confirm('是否确认导出所选数据？', '提示', {
    confirmButtonText: '确认',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(() => {
    exportLoginLog(proxy.addDateRange(queryParams.value, dateRange.value)).then(response => {
      resolveBlob(response, 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet', '登录日志')
    })
  }).catch(() => {
  })
}

init()
</script>

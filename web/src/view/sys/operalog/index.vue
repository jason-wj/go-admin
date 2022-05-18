
<template>
  <div>
    <div class="gva-search-box">

      <el-form ref="searchForm" :model="queryParams" :inline="true" label-position="left">
        <el-form-item label="状态" prop="status">
          <el-select v-model="queryParams.status" placeholder="操作状态" clearable size="small" style="width: 160px">
            <el-option
              v-for="dict in statusOptions"
              :key="dict.dictValue"
              :label="dict.dictLabel"
              :value="dict.dictValue"
            />
          </el-select>
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
    <el-table v-loading="loading" stripe border :data="operalogList">
      <el-table-column label="序号" type="index" align="center" min-width="60">
        <template #default="scope">
          <span>{{ (queryParams.pageIndex - 1) * queryParams.pageSize + scope.$index + 1 }}</span>
        </template>
      </el-table-column>
      <el-table-column label="编号" align="center" min-width="70" prop="id" />
      <el-table-column label="Request Method" align="center" min-width="150" prop="requestMethod" />
      <el-table-column label="Request Host" align="center" min-width="150" prop="operIp" />
      <el-table-column label="Request Location" align="center" min-width="150" prop="operLocation" />
      <el-table-column label="Request 耗时" align="center" min-width="150" prop="latencyTime" />
      <el-table-column label="Request Url" align="center" min-width="300" prop="operUrl" />
      <el-table-column label="操作人员" prop="operName" align="center" min-width="100" :show-overflow-tooltip="true" />
      <el-table-column min-width="100" label="状态" align="center" prop="status" :formatter="statusFormat">
        <template #default="scope">
          {{ statusFormat(scope.row) }}
        </template>
      </el-table-column>
      <el-table-column label="操作日期" prop="operTime" align="center" min-width="160">
        <template #default="scope">
          <span>{{ formatDate(scope.row.operTime) }}</span>
        </template>
      </el-table-column>
      <el-table-column label="操作" min-width="150" fixed="right" align="center" class-name="small-padding fixed-width">
        <template #default="scope">
          <el-button
            size="small"
            type="text"
            icon="view"
            @click="handleView(scope.row,scope.index)"
          >详细</el-button>
          <el-button size="small" type="text" icon="delete" @click="handleDelete(scope.row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>
    <Pagination v-show="total>0" v-model:page="queryParams.pageIndex" v-model:limit="queryParams.pageSize" :total="total" @pagination="getList" />
    <!-- 添加或修改对话框 -->
    <el-dialog v-model="open" title="操作日志详细" :close-on-click-modal="false" width="700px" append-to-body destroy-on-close>
      <el-form ref="inputform" :model="form" label-width="100px" size="small">
        <el-row>
          <el-col :span="24">
            <el-form-item label="请求地址：">{{ form.operUrl }}</el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item
              label="登录信息："
            >{{ form.operName }} / {{ form.operIp }} / {{ form.operLocation }}</el-form-item>
          </el-col>

          <el-col :span="12">
            <el-form-item label="请求方式：">{{ form.requestMethod }}</el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="耗时：">{{ form.latencyTime }}</el-form-item>
          </el-col>
          <el-col :span="24">
            <el-form-item label="请求参数：">{{ form.operParam }}</el-form-item>
          </el-col>
          <el-col :span="24">
            <el-form-item label="返回参数：">{{ form.jsonResult }}</el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="操作状态：">
              <el-select v-model="form.status" :disabled="true">
                <el-option
                  v-for="dict in statusOptions"
                  :key="dict.dictValue"
                  :label="dict.dictLabel"
                  :value="dict.dictValue"
                />
              </el-select>

            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="操作时间：">{{ formatDate(form.operTime) }}</el-form-item>
          </el-col>
          <el-col :span="24">
            <el-form-item v-if="form.status === 1" label="异常信息：">{{ form.errorMsg }}</el-form-item>
          </el-col>
        </el-row>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button @click="open = false">关 闭</el-button>
      </div>
    </el-dialog>
  </div>
</template>

<script>
export default {
  name: 'OperaLog',
}
</script>

<script setup>
import { getCurrentInstance, ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { resolveBlob } from '@/utils/zipdownload'
import { listOperaLog, delOperalog, exportOperaLog } from '@/api/sys/operalog'

const { proxy } = getCurrentInstance()

// 遮罩层
const loading = ref(true)
// 总条数
const total = ref(0)
// 是否显示弹出层
const open = ref(false)
// 类型数据字典
const statusOptions = ref([])
// 日期范围
const dateRange = ref([])
// 表单参数
const form = ref({})
// 数据列表
const operalogList = ref([])
// 查询参数
const queryParams = ref({
  pageIndex: 1,
  pageSize: 10,
  title: undefined,
  operName: undefined,
  businessType: undefined,
  status: undefined,
  createdAtOrder: 'desc'
})
// ref form
const inputform = ref(null)
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
  listOperaLog(proxy.addDateRange(queryParams.value, dateRange.value)).then(response => {
    operalogList.value = response.data.list
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

/** 详细按钮操作 */
const handleView = (row) => {
  open.value = true
  form.value = row
}

/** 删除按钮操作 */
const handleDelete = (row) => {
  const ids = [row.id]
  ElMessageBox.confirm('是否确认删除编号为"' + ids + '"的数据项?', '警告', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(() => {
    return delOperalog({ 'ids': ids })
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
    exportOperaLog(proxy.addDateRange(queryParams.value, dateRange.value)).then(response => {
      resolveBlob(response, 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet', '操作日志')
    })
  }).catch(() => {
  })
}

init()
</script>


<template>
  <div>
    <div class="gva-search-box">

      <el-form ref="searchForm" :model="queryParams" :inline="true" label-position="left">
        <el-form-item label="配置名称" prop="configName">
          <el-input v-model="queryParams.configName" placeholder="请输入配置名称" clearable size="small" @keyup.enter.native="handleQuery" />
        </el-form-item>
        <el-form-item label="配置键名" prop="configKey">
          <el-input v-model="queryParams.configKey" placeholder="请输入配置键名" clearable size="small" @keyup.enter.native="handleQuery" />
        </el-form-item>
        <el-form-item label="系统内置" prop="configType">
          <el-select v-model="queryParams.configType" placeholder="系统配置" clearable size="small">
            <el-option
                v-for="dict in configTypeOptions"
                :key="dict.dictValue"
                :label="dict.dictLabel"
                :value="dict.dictValue"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="是否前台" prop="isFrontend">
          <el-select v-model="queryParams.isFrontend" placeholder="系统配置是否前台" clearable size="small">
            <el-option
                v-for="dict in isFrontendOptions"
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
          <el-button type="primary" icon="plus" size="small" @click="handleAdd">新增</el-button>
        </el-col>
        <el-col :span="1.5">
          <el-button type="success" icon="download" size="small" @click="handleExport">Excel导出</el-button>
        </el-col>
      </el-row>
    </div>
    <el-table v-loading="loading" stripe border :data="configList">
      <el-table-column label="序号" type="index" align="center" width="60">
        <template #default="scope">
          <span>{{ (queryParams.pageIndex - 1) * queryParams.pageSize + scope.$index + 1 }}</span>
        </template>
      </el-table-column>
      <el-table-column min-width="100" label="主键编码" align="center" prop="id" :show-overflow-tooltip="true" />
      <el-table-column min-width="150" label="配置名称" align="center" prop="configName" :show-overflow-tooltip="true" />
      <el-table-column min-width="150" label="配置键名" align="center" prop="configKey" :show-overflow-tooltip="true" />
      <el-table-column min-width="180" label="配置键值" align="center" prop="configValue" :show-overflow-tooltip="true" />
      <el-table-column min-width="100" label="系统内置" align="center" prop="configType" :formatter="configTypeFormat">
        <template #default="scope">
          {{ configTypeFormat(scope.row) }}
        </template>
      </el-table-column>
      <el-table-column min-width="100" label="是否前台" align="center" prop="isFrontend" :formatter="isFrontendFormat">
        <template #default="scope">
          {{ isFrontendFormat(scope.row) }}
        </template>
      </el-table-column>
      <el-table-column min-width="200" label="备注" align="center" prop="remark" :show-overflow-tooltip="true" />
      <el-table-column min-width="100" label="创建者" align="center" prop="createBy" :show-overflow-tooltip="true" />
      <el-table-column min-width="100" label="更新者" align="center" prop="updateBy" :show-overflow-tooltip="true" />
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
      <el-table-column min-width="130" fixed="right" label="操作" align="center" class-name="small-padding fixed-width">
        <template #default="scope">
          <el-button size="small" type="text" icon="edit" @click="handleUpdate(scope.row)">修改</el-button>
          <el-button size="small" type="text" icon="delete" @click="handleDelete(scope.row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>
    <Pagination v-show="total>0" v-model:page="queryParams.pageIndex" v-model:limit="queryParams.pageSize" :total="total" @pagination="getList" />
    <!-- 添加或修改对话框 -->
    <el-dialog v-model="open" :close-on-click-modal="false" :title="title" width="500px" append-to-body destroy-on-close>
      <el-form ref="inputform" :model="form" :rules="rules" label-width="80px">
        <el-form-item label="配置名称" prop="configName">
          <el-input v-model="form.configName" placeholder="配置名称" />
        </el-form-item>
        <el-form-item label="配置键名" prop="configKey">
          <el-input v-model="form.configKey" placeholder="配置键名" :disabled="isEdit" />
        </el-form-item>
        <el-form-item label="配置键值" prop="configValue">
          <el-input v-model="form.configValue" placeholder="配置键值" />
        </el-form-item>
        <el-form-item label="系统内置" prop="configType">
          <el-select v-model="form.configType" placeholder="请选择">
            <el-option
                v-for="dict in configTypeOptions"
                :key="dict.dictValue"
                :label="dict.dictLabel"
                :value="dict.dictValue"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="是否前台" prop="isFrontend">
          <el-select v-model="form.isFrontend" placeholder="请选择">
            <el-option
                v-for="dict in isFrontendOptions"
                :key="dict.dictValue"
                :label="dict.dictLabel"
                :value="dict.dictValue"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="备注" prop="remark">
          <el-input v-model="form.remark" placeholder="备注" />
        </el-form-item>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button type="primary" @click="submitForm">确 定</el-button>
        <el-button @click="cancel">取 消</el-button>
      </div>
    </el-dialog>
  </div>
</template>

<script>
export default {
  name: 'Config',
}
</script>

<script setup>
import { getCurrentInstance, ref, nextTick } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { resolveBlob } from '@/utils/zipdownload'
import { exportConfig, addConfig, delConfig, getConfig, listConfig, updateConfig } from '@/api/sys/config'

const { proxy } = getCurrentInstance()

// 遮罩层
const loading = ref(true)
// 总条数
const total = ref(0)
// 弹出层标题
const title = ref('')
// 是否显示弹出层
const open = ref(false)
// 是否编辑
const isEdit = ref(false)
// 日期范围
const dateRange = ref([])
// 数据列表
const configList = ref([])
// 类型数据字典
const configTypeOptions = ref([])
// 类型数据字典
const isFrontendOptions = ref([])
// 查询参数
const queryParams = ref({
  pageIndex: 1,
  pageSize: 10,
  id: undefined,
  configName: undefined,
  configKey: undefined,
  configType: undefined,
  isFrontend: undefined,
})
// 表单参数
const form = ref({})
// 表单校验
const rules = ref({
  configName: [{ required: true, message: '配置名称不能为空', trigger: 'blur' }],
  configKey: [{ required: true, message: '配置键名不能为空', trigger: 'blur' }],
  configType: [{ required: true, message: '系统内置不能为空', trigger: 'blur' }],
  isFrontend: [{ required: true, message: '是否前台不能为空', trigger: 'blur' }],
})

// ref form
const inputform = ref(null)
const searchForm = ref(null)

const init = () => {
  getList()
  proxy.getDicts('sys_yes_no').then(response => {
    configTypeOptions.value = response.data
    isFrontendOptions.value = response.data
    console.log(isFrontendOptions.value)
  })
}

/** 查询参数列表 */
const getList = () => {
  loading.value = true
  listConfig(proxy.addDateRange(queryParams.value, dateRange.value)).then(response => {
    configList.value = response.data.list
    total.value = response.data.count
    loading.value = false
  })
}

// 取消按钮
const cancel = () => {
  open.value = false
  reset()
}

// 表单重置
const reset = () => {
  form.value = {
    id: undefined,
    configName: undefined,
    configKey: undefined,
    configValue: undefined,
    configType: undefined,
    isFrontend: undefined,
    remark: undefined,
  }
  nextTick(() => {
    if (inputform.value != null) {
      inputform.value.resetFields()
    }
  })
}

// 字典
const configTypeFormat = (row) => {
  return proxy.selectDictLabel(configTypeOptions.value, row.configType)
}
const isFrontendFormat = (row) => {
  return proxy.selectDictLabel(isFrontendOptions.value, row.isFrontend)
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

/** 新增按钮操作 */
const handleAdd = () => {
  open.value = true
  title.value = '添加系统配置'
  isEdit.value = false
  reset()
}

/** 修改按钮操作 */
const handleUpdate = (row) => {
  const id = row.id
  getConfig(id).then(response => {
    open.value = true
    title.value = '修改系统配置'
    isEdit.value = true
    reset()
    form.value = response.data
  })
}

/** 提交按钮 */
const submitForm = () => {
  inputform.value.validate(valid => {
    if (valid) {
      if (form.value.id !== undefined) {
        updateConfig(form.value).then(response => {
          if (response.code === 200) {
            open.value = false
            getList()
            ElMessage({
              type: 'success',
              message: response.msg,
              showClose: true,
            })
          }
        })
      } else {
        addConfig(form.value).then(response => {
          if (response.code === 200) {
            open.value = false
            getList()
            ElMessage({
              type: 'success',
              message: response.msg,
              showClose: true,
            })
          }
        })
      }
    }
  })
}

/** 删除按钮操作 */
const handleDelete = (row) => {
  const ids = [row.id]

  ElMessageBox.confirm('是否确认删除编号为"' + ids + '"的数据项?', '警告', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(() => {
    return delConfig({ 'ids': ids })
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
    exportConfig(proxy.addDateRange(queryParams.value, dateRange.value)).then(response => {
      resolveBlob(response, 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet', '系统配置')
    })
  }).catch(() => {
  })
}

init()
</script>

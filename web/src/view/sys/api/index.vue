
<template>
  <div>
    <div class="gva-search-box">

      <el-form ref="searchForm" :model="queryParams" :inline="true" label-position="left">
        <el-form-item label="标题" prop="title">
          <el-input v-model="queryParams.title" placeholder="请输入标题" clearable size="small" @keyup.enter.native="handleQuery" />
        </el-form-item>
        <el-form-item label="地址" prop="path">
          <el-input v-model="queryParams.path" placeholder="请输入地址" clearable size="small" @keyup.enter.native="handleQuery" />
        </el-form-item>
        <el-form-item label="接口类型" prop="type">
          <el-select v-model="queryParams.type" placeholder="接口类型" clearable size="small">
            <el-option
              v-for="dict in typeOptions"
              :key="dict.dictValue"
              :label="dict.dictLabel"
              :value="dict.dictValue"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="请求类型" prop="action">
          <el-select v-model="queryParams.action" placeholder=" 接口请求类型" clearable size="small">
            <el-option
              v-for="dict in actionOptions"
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
    <el-table v-loading="loading" stripe border :data="apiList">
      <el-table-column label="序号" type="index" align="center" width="60">
        <template #default="scope">
          <span>{{ (queryParams.pageIndex - 1) * queryParams.pageSize + scope.$index + 1 }}</span>
        </template>
      </el-table-column>
      <el-table-column min-width="100" label="主键编码" align="center" prop="id" :show-overflow-tooltip="true" />
      <el-table-column min-width="150" label="标题" align="center" prop="title" :show-overflow-tooltip="true" />
      <el-table-column min-width="300" label="地址" align="center" prop="path" :show-overflow-tooltip="true" />
      <el-table-column min-width="120" label="接口类型" align="center" prop="type" :formatter="typeFormat">
        <template #default="scope">
          {{ typeFormat(scope.row) }}
        </template>
      </el-table-column>
      <el-table-column min-width="120" label="请求类型" align="center" prop="action" :formatter="actionFormat">
        <template #default="scope">
          {{ actionFormat(scope.row) }}
        </template>
      </el-table-column>
      <el-table-column min-width="300" label="handle" align="center" prop="handle" :show-overflow-tooltip="true" />
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
      <el-table-column min-width="100" label="创建者" align="center" prop="createBy" :show-overflow-tooltip="true" />
      <el-table-column min-width="100" label="更新者" align="center" prop="updateBy" :show-overflow-tooltip="true" />
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
        <el-form-item label="handle" prop="handle">
          <el-input v-model="form.handle" placeholder="handle" />
        </el-form-item>
        <el-form-item label="标题" prop="title">
          <el-input v-model="form.title" placeholder="标题" />
        </el-form-item>
        <el-form-item label="接口类型" prop="type">
          <el-select v-model="form.type" placeholder="请选择">
            <el-option
              v-for="dict in typeOptions"
              :key="dict.dictValue"
              :label="dict.dictLabel"
              :value="dict.dictValue"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="请求类型" prop="action">
          <el-select v-model="form.action" placeholder="请选择">
            <el-option
              v-for="dict in actionOptions"
              :key="dict.dictValue"
              :label="dict.dictLabel"
              :value="dict.dictValue"
            />
          </el-select>
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
  name: 'Api',
}
</script>

<script setup>
import { getCurrentInstance, ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { resolveBlob } from '@/utils/zipdownload'
import { exportApi, delApi, getApi, listApi, updateApi } from '@/api/sys/api'
import { nextTick } from 'vue'

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
const apiList = ref([])
// 类型数据字典
const typeOptions = ref([])
// 类型数据字典
const actionOptions = ref([])
// 查询参数
const queryParams = ref({
  pageIndex: 1,
  pageSize: 10,
  id: undefined,
  title: undefined,
  path: undefined,
  type: undefined,
  action: undefined,
})
// 表单参数
const form = ref({})
// 表单校验
const rules = ref({
  title: [{ required: true, message: '标题不能为空', trigger: 'blur' }],
  path: [{ required: true, message: '地址不能为空', trigger: 'blur' }],
  type: [{ required: true, message: '接口类型不能为空', trigger: 'blur' }],
  action: [{ required: true, message: '请求类型不能为空', trigger: 'blur' }],
})

// ref form
const inputform = ref(null)
const searchForm = ref(null)

const init = () => {
  getList()
  proxy.getDicts('sys_api_type').then(response => {
    typeOptions.value = response.data
  })
  proxy.getDicts('sys_request_type').then(response => {
    actionOptions.value = response.data
  })
}

/** 查询参数列表 */
const getList = () => {
  loading.value = true
  listApi(proxy.addDateRange(queryParams.value, dateRange.value)).then(response => {
    apiList.value = response.data.list
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
    handle: undefined,
    title: undefined,
    type: undefined,
    action: undefined,
  }
  nextTick(() => {
    if (inputform.value != null) {
      inputform.value.resetFields()
    }
  })
}

// 字典
const typeFormat = (row) => {
  return proxy.selectDictLabel(typeOptions.value, row.type)
}
const actionFormat = (row) => {
  return proxy.selectDictLabel(actionOptions.value, row.action)
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
  title.value = '添加 接口'
  isEdit.value = false
  reset()
}

/** 修改按钮操作 */
const handleUpdate = (row) => {
  const id = row.id
  getApi(id).then(response => {
    open.value = true
    title.value = '修改 接口'
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
        updateApi(form.value).then(response => {
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
    return delApi({ 'ids': ids })
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
    exportApi(proxy.addDateRange(queryParams.value, dateRange.value)).then(response => {
      resolveBlob(response, 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet', ' 接口')
    })
  }).catch(() => {
  })
}

init()
</script>


<template>
  <div>
    <div class="gva-search-box">

      <el-form ref="searchForm" :model="queryParams" :inline="true" label-position="left">
        <el-form-item label="字典名称" prop="dictType">
          <el-select v-model="queryParams.dictType" placeholder="字典名称选择" size="small">
            <el-option
              v-for="item in typeOptions"
              :key="item.dictId"
              :label="item.dictName"
              :value="item.dictType"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="字典标签" prop="dictLabel">
          <el-input v-model="queryParams.dictLabel" placeholder="请输入字典标签" clearable size="small" @keyup.enter.native="handleQuery" />
        </el-form-item>
        <el-form-item label="状态" prop="status">
          <el-select v-model="queryParams.status" placeholder="字典数据" clearable size="small">
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
          <el-button type="primary" icon="plus" size="small" @click="handleAdd">新增</el-button>
        </el-col>
        <el-col :span="1.5">
          <el-button type="success" icon="download" size="small" @click="handleExport">Excel导出</el-button>
        </el-col>
      </el-row>
    </div>
    <el-table v-loading="loading" stripe border :data="dictdataList">
      <el-table-column label="序号" type="index" align="center" width="60">
        <template #default="scope">
          <span>{{ (queryParams.pageIndex - 1) * queryParams.pageSize + scope.$index + 1 }}</span>
        </template>
      </el-table-column>
      <el-table-column min-width="100" label="字典编码" align="center" prop="dictCode" :show-overflow-tooltip="true" />
      <el-table-column min-width="100" label="字典标签" align="center" prop="dictLabel" :show-overflow-tooltip="true" />
      <el-table-column min-width="100" label="字典键值" align="center" prop="dictValue" :show-overflow-tooltip="true" />
      <el-table-column min-width="100" label="字典类型" align="center" prop="dictType" :show-overflow-tooltip="true" />
      <el-table-column min-width="100" label="字典排序" align="center" prop="dictSort" :show-overflow-tooltip="true" />
      <el-table-column min-width="100" label="状态" align="center" prop="status" :formatter="statusFormat">
        <template #default="scope">
          {{ statusFormat(scope.row) }}
        </template>
      </el-table-column>
      <el-table-column min-width="100" label="备注" align="center" prop="remark" :show-overflow-tooltip="true" />
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
        <el-form-item label="字典类型" prop="dictType">
          <el-input v-model="form.dictType" placeholder="字典类型" :disabled="true" />
        </el-form-item>
        <el-form-item label="字典标签" prop="dictLabel">
          <el-input v-model="form.dictLabel" placeholder="字典标签" :disabled="isEdit" />
        </el-form-item>
        <el-form-item label="字典键值" prop="dictValue">
          <el-input v-model="form.dictValue" placeholder="字典键值" :disabled="isEdit" />
        </el-form-item>
        <el-form-item label="字典排序" prop="dictSort">
          <el-input-number v-model="form.dictSort" />
        </el-form-item>
        <el-form-item label="状态" prop="status">
          <el-select v-model="form.status" placeholder="请选择">
            <el-option
              v-for="dict in statusOptions"
              :key="dict.dictValue"
              :label="dict.dictLabel"
              :value="dict.dictValue"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="备注" prop="remark">
          <el-input v-model="form.remark" placeholder="" />
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
  name: 'SysDictData',
}
</script>

<script setup>
import { getCurrentInstance, ref } from 'vue'
import { useRoute } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { resolveBlob } from '@/utils/zipdownload'
import { listSysDictType, getSysDictType } from '@/api/sys/dicttype'
import { exportSysDictData, addSysDictData, delSysDictData, getSysDictData, listSysDictData, updateSysDictData } from '@/api/sys/dictdata'
import { nextTick } from 'vue'

const route = useRoute()

const dictId = ref(route.params.dictId)

const defaultDictType = ref('')

const { proxy } = getCurrentInstance()

// 类型数据字典
const typeOptions = ref([])

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
const dictdataList = ref([])
// 类型数据字典
const statusOptions = ref([])
// 查询参数
const queryParams = ref({
  pageIndex: 1,
  pageSize: 10,
  dictType: undefined,
  dictLabel: undefined,
  status: undefined,
})
// 表单参数
const form = ref({})
// 表单校验
const rules = ref({
  dictLabel: [{ required: true, message: '字典标签不能为空', trigger: 'blur' }],
  status: [{ required: true, message: '不能为空', trigger: 'blur' }],
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

/** 查询字典类型详细 */
const sysDictType = async(dictId) => {
  const response = await getSysDictType(dictId)
  queryParams.value.dictType = response.data.dictType
  defaultDictType.value = response.data.dictType
  init()
}

/** 查询字典类型列表 */
const getListSysDictType = async() => {
  const response = await listSysDictType({ pageSize: 1000 })
  typeOptions.value = response.data.list
}

/** 查询参数列表 */
const getList = () => {
  console.log(queryParams.value)
  loading.value = true
  listSysDictData(proxy.addDateRange(queryParams.value, dateRange.value)).then(response => {
    dictdataList.value = response.data.list
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
    dictSort: undefined,
    dictLabel: undefined,
    dictValue: undefined,
    dictType: undefined,
    status: undefined,
    remark: undefined,
  }
  nextTick(() => {
    if (inputform.value != null) {
      inputform.value.resetFields()
    }
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
  queryParams.value.dictType = defaultDictType.value
  handleQuery()
}

/** 新增按钮操作 */
const handleAdd = () => {
  open.value = true
  title.value = '添加字典数据'
  isEdit.value = false
  reset()
  form.value.dictType = queryParams.value.dictType
}

/** 修改按钮操作 */
const handleUpdate = (row) => {
  const dictCode = row.dictCode
  getSysDictData(dictCode).then(response => {
    open.value = true
    title.value = '修改字典数据'
    isEdit.value = true
    reset()
    form.value = response.data
  })
}

/** 提交按钮 */
const submitForm = () => {
  inputform.value.validate(valid => {
    if (valid) {
      if (form.value.dictCode !== undefined) {
        updateSysDictData(form.value).then(response => {
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
        addSysDictData(form.value).then(response => {
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
  const ids = [row.dictCode]

  ElMessageBox.confirm('是否确认删除编号为"' + ids + '"的数据项?', '警告', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(() => {
    return delSysDictData({ 'ids': ids })
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
    exportSysDictData(proxy.addDateRange(queryParams.value, dateRange.value)).then(response => {
      resolveBlob(response, 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet', '字典数据')
    })
  }).catch(() => {
  })
}

sysDictType(dictId.value)
getListSysDictType()
</script>


<template>
  <div>
    <div class="gva-search-box">

      <el-form ref="searchForm" :model="queryParams" :inline="true" label-position="left">
        <el-form-item label="主键" prop="id">
          <el-input v-model="queryParams.id" placeholder="请输入主键" clearable size="small" @keyup.enter.native="handleQuery" />
        </el-form-item>
        <el-form-item label="版本号" prop="version">
          <el-input v-model="queryParams.version" placeholder="请输入版本号" clearable size="small" @keyup.enter.native="handleQuery" />
        </el-form-item>
        <el-form-item label="系统平台" prop="platform"><el-select v-model="queryParams.platform" placeholder="系统平台" clearable size="small">
          <el-option
            v-for="dict in platformOptions"
            :key="dict.dictValue"
            :label="dict.dictLabel"
            :value="dict.dictValue"
          />
        </el-select>
        </el-form-item>
        <el-form-item label="版本类型" prop="type"><el-select v-model="queryParams.type" placeholder="版本类型" clearable size="small">
          <el-option
            v-for="dict in typeOptions"
            :key="dict.dictValue"
            :label="dict.dictLabel"
            :value="dict.dictValue"
          />
        </el-select>
        </el-form-item>
        <el-form-item label="下载类型" prop="downloadType"><el-select v-model="queryParams.downloadType" placeholder="下载类型" clearable size="small">
          <el-option
            v-for="dict in downloadTypeOptions"
            :key="dict.dictValue"
            :label="dict.dictLabel"
            :value="dict.dictValue"
          />
        </el-select>
        </el-form-item>
        <el-form-item label="发布状态" prop="status"><el-select v-model="queryParams.status" placeholder="发布状态" clearable size="small">
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
    <el-table v-loading="loading" stripe border :data="appList">
      <el-table-column label="序号" type="index" align="center" width="60">
        <template #default="scope">
          <span>{{ (queryParams.pageIndex - 1) * queryParams.pageSize + scope.$index + 1 }}</span>
        </template>
      </el-table-column>
      <el-table-column min-width="100" label="主键" align="center" prop="id" :show-overflow-tooltip="true" />
      <el-table-column min-width="100" label="版本号" align="center" prop="version" :show-overflow-tooltip="true" />
      <el-table-column min-width="100" label="发布状态" align="center" prop="status" :formatter="statusFormat">
        <template #default="scope">
          {{ statusFormat(scope.row) }}
        </template>
      </el-table-column>
      <el-table-column min-width="100" label="平台" align="center" prop="platform" :formatter="platformFormat">
        <template #default="scope">
          {{ platformFormat(scope.row) }}
        </template>
      </el-table-column>
      <el-table-column min-width="100" label="版本" align="center" prop="type" :formatter="typeFormat">
        <template #default="scope">
          {{ typeFormat(scope.row) }}
        </template>
      </el-table-column>
      <el-table-column min-width="100" label="下载类型" align="center" prop="downloadType" :formatter="downloadTypeFormat">
        <template #default="scope">
          {{ downloadTypeFormat(scope.row) }}
        </template>
      </el-table-column>
      <el-table-column min-width="300" label="本地地址" align="center" prop="localAddress" :show-overflow-tooltip="true" />
      <el-table-column min-width="300" label="下载地址" align="center" prop="downloadUrl" :show-overflow-tooltip="true" />
      <el-table-column min-width="300" label="备注信息" align="center" prop="remark" :show-overflow-tooltip="true" />
      <el-table-column min-width="100" label="创建者" align="center" prop="createBy" :show-overflow-tooltip="true" />
      <el-table-column min-width="100" label="更新者" align="center" prop="updateBy" :show-overflow-tooltip="true" />
      <el-table-column min-width="130" label="更新时间" align="center" prop="updatedAt" :show-overflow-tooltip="true">
        <template #default="scope">
          <span>{{ formatDate(scope.row.updatedAt) }}</span>
        </template>
      </el-table-column>
      <el-table-column min-width="130" label="创建时间" align="center" prop="createdAt" :show-overflow-tooltip="true">
        <template #default="scope">
          <span>{{ formatDate(scope.row.createdAt) }}</span>
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
        <el-row v-if="!isEdit">
          <el-col :span="12">
            <el-form-item label="版本号" prop="version">
              <el-input v-model="form.version" placeholder="版本号" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="平台" prop="platform">
              <el-select v-model="form.platform" placeholder="请选择">
                <el-option
                  v-for="dict in platformOptions"
                  :key="dict.dictValue"
                  :label="dict.dictLabel"
                  :value="dict.dictValue"
                />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="下载类型" prop="downloadType">
              <el-select v-model="form.downloadType" placeholder="请选择">
                <el-option
                  v-for="dict in downloadTypeOptions"
                  :key="dict.dictValue"
                  :label="dict.dictLabel"
                  :value="dict.dictValue"
                />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="App类型" prop="type">
              <el-select v-model="form.type" placeholder="请选择">
                <el-option
                  v-for="dict in typeOptions"
                  :key="dict.dictValue"
                  :label="dict.dictLabel"
                  :value="dict.dictValue"
                />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col v-if="form.downloadType==='1'" :span="24">
            <el-form-item label="下载地址" prop="downloadUrl">
              <el-input v-model="form.downloadUrl" placeholder="下载地址" />
            </el-form-item>
          </el-col>
          <el-col v-if="form.downloadType==='0'" :span="24">
            <el-form-item label="上传App">
              <el-upload
                ref="upload"
                :limit="1"
                accept=".apk, .ipa"
                :headers="upload.headers"
                :action="url"
                :disabled="upload.isUploading"
                :on-progress="handleFileUploadProgress"
                :on-success="handleFileSuccess"
                :show-file-list="true"
                :auto-upload="true"
                drag
              >
                <i class="el-icon-upload" />
                <div class="el-upload__text">
                  将文件拖到此处，或<em>点击上传</em>
                </div>
                <div slot="tip" class="el-upload__tip" style="color:red">提示：仅允许导入“apk”或“ipa”格式文件！</div>
              </el-upload>
            </el-form-item>
          </el-col>
          <el-col :span="24">
            <el-form-item label="更新内容" prop="remark">
              <el-input
                v-model="form.remark"
                type="textarea"
                :rows="4"
                placeholder="请输入更新内容"
              />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row v-if="isEdit">
          <el-col :span="12">
            <el-form-item label="App状态" prop="status">
              <el-select v-model="form.status" placeholder="请选择">
                <el-option
                  v-for="dict in statusOptions"
                  :key="dict.dictValue"
                  :label="dict.dictLabel"
                  :value="dict.dictValue"
                />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>
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
  name: 'App',
}
</script>

<script setup>
import { getCurrentInstance, ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { resolveBlob } from '@/utils/zipdownload'
import { useUserStore } from '@/pinia/modules/user'
import { exportApp, addApp, delApp, getApp, listApp, updateApp } from '@/api/plugins/filemgr/app'
import { nextTick } from 'vue'

const { proxy } = getCurrentInstance()

const userStore = useUserStore()

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
const appList = ref([])
// 类型数据字典
const platformOptions = ref([])
// 类型数据字典
const typeOptions = ref([])
// 类型数据字典
const downloadTypeOptions = ref([])
// 类型数据字典
const statusOptions = ref([])
// 查询参数
const queryParams = ref({
  pageIndex: 1,
  pageSize: 10,
  version: undefined,
  platform: undefined,
  type: undefined,
  downloadType: undefined,
  downloadUrl: undefined
})
// 表单参数
const form = ref({})
// 表单校验
const rules = ref({
  version: [{ required: true, message: '版本号不能为空', trigger: 'blur' }],
  platform: [{ required: true, message: '平台不能为空', trigger: 'blur' }],
  type: [{ required: true, message: '版本不能为空', trigger: 'blur' }],
  downloadType: [{ required: true, message: '下载类型不能为空', trigger: 'blur' }],
  status: [{ required: true, message: '发布状态不能为空', trigger: 'blur' }],
})

const baseUrl = ref(import.meta.env.VITE_BASE_API)
const ip = ref(import.meta.env.VITE_BASE_PATH)
const port = ref(import.meta.env.VITE_SERVER_PORT)

const url = ref(ip.value + ':' + port.value + baseUrl.value + '/v1/filemgr/app/upload')

const upload = ref({
  // 是否禁用上传
  isUploading: false,
  // 设置上传的请求头部
  headers: { Authorization: 'Bearer ' + userStore.token },
})

// ref form
const inputform = ref(null)
const searchForm = ref(null)

const init = () => {
  getList()
  proxy.getDicts('app_platform').then(response => {
    platformOptions.value = response.data
  })
  proxy.getDicts('app_type').then(response => {
    typeOptions.value = response.data
  })
  proxy.getDicts('app_download_type').then(response => {
    downloadTypeOptions.value = response.data
  })
  proxy.getDicts('app_status').then(response => {
    statusOptions.value = response.data
  })
}

/** 查询参数列表 */
const getList = () => {
  loading.value = true
  listApp(proxy.addDateRange(queryParams.value, dateRange.value)).then(response => {
    appList.value = response.data.list
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
    version: undefined,
    platform: undefined,
    type: undefined,
    localAddress: undefined,
    bucketName: undefined,
    ossKey: undefined,
    downloadNum: undefined,
    downloadType: undefined,
    downloadUrl: undefined,
    remark: undefined,
    status: undefined,
  }
  nextTick(() => {
    if (inputform.value != null) {
      inputform.value.resetFields()
    }
  })
}

// 字典
const platformFormat = (row) => {
  return proxy.selectDictLabel(platformOptions.value, row.platform)
}
const typeFormat = (row) => {
  return proxy.selectDictLabel(typeOptions.value, row.type)
}
const downloadTypeFormat = (row) => {
  return proxy.selectDictLabel(downloadTypeOptions.value, row.downloadType)
}
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

/** 新增按钮操作 */
const handleAdd = () => {
  open.value = true
  title.value = '添加app升级管理'
  isEdit.value = false
  reset()
}

/** 修改按钮操作 */
const handleUpdate = (row) => {
  const id = row.id
  getApp(id).then(response => {
    open.value = true
    title.value = '修改app升级管理'
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
        updateApp(form.value).then(response => {
          if (response.code === 200) {
            ElMessage({
              type: 'success',
              message: response.msg,
              showClose: true,
            })
            open.value = false
            getList()
          }
        })
      } else {
        addApp(form.value).then(response => {
          if (response.code === 200) {
            ElMessage({
              type: 'success',
              message: response.msg,
              showClose: true,
            })
            open.value = false
            getList()
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
    return delApp({ 'ids': ids })
  }).then((response) => {
    let respType = 'error'
    if (response.code === 200) {
      respType = 'success'
      open.value = false
      getList()
    }
    ElMessage({
      type: respType,
      message: response.msg,
      showClose: true,
    })
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
    exportApp(proxy.addDateRange(queryParams.value, dateRange.value)).then(response => {
      resolveBlob(response, 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet', 'app升级管理')
    })
  }).catch(() => {
  })
}

// 文件上传中处理
const handleFileUploadProgress = (event, file, fileList) => {
  upload.value.isUploading = true
}
// 文件上传成功处理
const handleFileSuccess = (response, file, fileList) => {
  upload.value.isUploading = false
  if (response.code !== 200 || response.data === undefined || response.data === '') {
    return
  }
  form.value.localAddress = response.data
}

init()
</script>

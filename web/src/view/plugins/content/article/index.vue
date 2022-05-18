
<template>
  <div>
    <div class="gva-search-box">

      <el-form ref="searchForm" :model="queryParams" :inline="true" label-position="left">
        <el-form-item label="主键编码" prop="id">
          <el-input v-model="queryParams.id" placeholder="请输入主键编码" clearable size="small" @keyup.enter.native="handleQuery" />
        </el-form-item>
        <el-form-item label="分类编号" prop="cateId">
          <el-input v-model="queryParams.cateId" placeholder="请输入分类编号" clearable size="small" @keyup.enter.native="handleQuery" />
        </el-form-item>
        <el-form-item label="名称" prop="name">
          <el-input v-model="queryParams.name" placeholder="请输入名称" clearable size="small" @keyup.enter.native="handleQuery" />
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
    <el-table v-loading="loading" stripe border :data="articleList">
      <el-table-column label="序号" type="index" align="center" width="60">
        <template #default="scope">
          <span>{{ (queryParams.pageIndex - 1) * queryParams.pageSize + scope.$index + 1 }}</span>
        </template>
      </el-table-column>
      <el-table-column min-width="100" label="主键编码" align="center" prop="id" :show-overflow-tooltip="true" />
      <el-table-column min-width="100" label="分类编号" align="center" prop="cateId" :show-overflow-tooltip="true" />
      <el-table-column min-width="100" label="名称" align="center" prop="name" :show-overflow-tooltip="true" />
      <el-table-column min-width="100" label="内容" align="center" prop="content" :show-overflow-tooltip="true" />
      <el-table-column min-width="100" label="备注信息" align="center" prop="remark" :show-overflow-tooltip="true" />
      <el-table-column min-width="100" label="状态" align="center" prop="status" :formatter="statusFormat">
        <template #default="scope">
          {{ statusFormat(scope.row) }}
        </template>
      </el-table-column>
      <el-table-column min-width="100" label="创建者" align="center" prop="createBy" :show-overflow-tooltip="true" />
      <el-table-column min-width="100" label="更新者" align="center" prop="updateBy" :show-overflow-tooltip="true" />
      <el-table-column min-width="170" label="更新时间" align="center" prop="updatedAt" :show-overflow-tooltip="true">
        <template #default="scope">
          <span>{{ formatDate(scope.row.updatedAt) }}</span>
        </template>
      </el-table-column>
      <el-table-column min-width="170" label="创建时间" align="center" prop="createdAt" :show-overflow-tooltip="true">
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
    <el-dialog v-model="open" :close-on-click-modal="false" :title="title" width="800px" append-to-body destroy-on-close>
      <el-form ref="inputform" :model="form" :rules="rules" label-width="80px">
        <el-form-item label="分类编号" prop="cateId">
          <el-input v-model="form.cateId" placeholder="分类编号" readonly @click.native="openSelectCategory">
            <template slot="prepend">{{ categoryName }}</template>
          </el-input>
        </el-form-item>
        <el-form-item label="名称" prop="name">
          <el-input v-model="form.name" placeholder="名称" />
        </el-form-item>
        <el-form-item label="内容" prop="content">
          <Tinymce ref="editorRef" v-model:modelValue="form.content" :disabled="disabled" :height="200" :width="660" />
        </el-form-item>
        <el-form-item label="备注信息" prop="remark">
          <el-input v-model="form.remark" placeholder="备注信息" />
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
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button type="primary" @click="submitForm">确 定</el-button>
        <el-button @click="cancel">取 消</el-button>
      </div>
      <SelectCategory ref="selectCategory" @ok="getCategorySelectInfo" />
    </el-dialog>
  </div>
</template>

<script>
export default {
  name: 'Article',
}
</script>

<script setup>
import { getCurrentInstance, ref } from 'vue'
import Tinymce from '@/components/TEditor/index.vue'
import SelectCategory from '@/view/plugins/content/select/selectCategory.vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { resolveBlob } from '@/utils/zipdownload'
import { exportArticle, addArticle, delArticle, getArticle, listArticle, updateArticle } from '@/api/plugins/content/article'
import { nextTick } from 'vue'

const { proxy } = getCurrentInstance()

const categoryName = ref('请选择')
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
const articleList = ref([])
// 类型数据字典
const statusOptions = ref([])
// 查询参数
const queryParams = ref({
  pageIndex: 1,
  pageSize: 10,
  id: undefined,
  cateId: undefined,
  name: undefined,
  content: undefined,
  remark: undefined,
  status: undefined,
  createBy: undefined,
  updateBy: undefined,
  updatedAt: undefined,
  createdAt: undefined,
})
// 表单参数
const form = ref({})

const disabled = ref(true)

// 表单校验
const rules = ref({
  cateId: [{ required: true, message: '分类编号不能为空', trigger: 'blur' }],
  name: [{ required: true, message: '名称不能为空', trigger: 'blur' }],
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
  listArticle(proxy.addDateRange(queryParams.value, dateRange.value)).then(response => {
    articleList.value = response.data.list
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
  disabled.value = true
  form.value = {
    id: undefined,
    cateId: undefined,
    name: undefined,
    content: undefined,
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
  title.value = '添加文章管理'
  isEdit.value = false
  disabled.value = false
  reset()
}

/** 修改按钮操作 */
const handleUpdate = (row) => {
  const id = row.id
  getArticle(id).then(response => {
    open.value = true
    disabled.value = false
    title.value = '修改文章管理'
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
        updateArticle(form.value).then(response => {
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
        addArticle(form.value).then(response => {
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
    return delArticle({ 'ids': ids })
  }).then((response) => {
    if (response.code === 200) {
      ElMessage({
        type: 'success',
        message: response.msg,
        showClose: true,
      })
      open.value = false
      getList()
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
    exportArticle(proxy.addDateRange(queryParams.value, dateRange.value)).then(response => {
      resolveBlob(response, 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet', '文章管理')
    })
  }).catch(() => {
  })
}

/** 显示选项框selectUserLevel */
const selectCategory = ref(null)
const openSelectCategory = () => {
  selectCategory.value.show()
}

/** 获取返回的信息 */
const getCategorySelectInfo = (row) => {
  form.value.cateId = row.id
  categoryName.value = row.name
  inputform.value.validateField('cateId')
}

init()
</script>

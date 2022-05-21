<template>
  <div>
    <div class="gva-search-box">
      <el-form ref="searchForm" :model="queryParams" :inline="true" label-position="left">
        <el-form-item label="表名称" prop="tableName">
          <el-input
              v-model="queryParams.tableName"
              placeholder="请输入表名称"
              clearable
              size="small"
              @keyup.enter.native="handleQuery"
          />
        </el-form-item>
        <el-form-item label="菜单名称" prop="tableComment">
          <el-input
              v-model="queryParams.tableComment"
              placeholder="请输入菜单名称"
              clearable
              size="small"
              @keyup.enter.native="handleQuery"
          />
        </el-form-item>

        <el-form-item>
          <el-button size="small" type="primary" icon="search" @click="handleQuery">查询</el-button>
          <el-button size="small" icon="refresh" @click="resetQuery">重置</el-button>
        </el-form-item>
      </el-form>
    </div>
    <div class="gva-table-box">
      <div class="gva-btn-list">
        <el-button size="small" type="primary" icon="upload" @click="openImportTable">导入</el-button>
      </div>

      <el-table v-loading="loading" stripe border :data="tableList">
        <!--        <el-table-column label="序号" align="center" prop="tableId" min-width="60" />-->
        <el-table-column label="序号" type="index" align="center" min-width="30">
          <template #default="scope">
            <span>{{ (queryParams.pageIndex - 1) * queryParams.pageSize + scope.$index + 1 }}</span>
          </template>
        </el-table-column>
        <el-table-column label="表编号" align="center" prop="tableId" min-width="30" />
        <el-table-column label="表名称" align="center" prop="tableName" :show-overflow-tooltip="true" min-width="80" />
        <el-table-column label="菜单名称" align="center" prop="tableComment" :show-overflow-tooltip="true" min-width="60" />
        <el-table-column label="模型名称" align="center" prop="className" :show-overflow-tooltip="true" min-width="60" />
        <el-table-column label="创建时间" align="center" prop="createdAt" min-width="80">
          <template #default="scope">
            <span>{{ formatDate(scope.row.createdAt) }}</span>
          </template>
        </el-table-column>

        <el-table-column label="操作" align="center" class-name="small-padding fixed-width" min-width="100">
          <template #default="scope">
            <el-button type="text" size="small" icon="el-icon-edit" @click="openEditTable(scope.row)">编辑</el-button>
            <el-button type="text" size="small" icon="el-icon-view" @click="handlePreview(scope.row)">预览</el-button>
            <el-button slot="reference" type="text" size="small" icon="el-icon-view" @click="handleDownloadCode(scope.row)">代码下载</el-button>
            <el-button slot="reference" type="text" size="small" icon="el-icon-view" @click="handleToProject(scope.row)">代码生成</el-button>
            <el-button slot="reference" type="text" size="small" icon="el-icon-view" @click="handleToDB(scope.row)">配置生成</el-button>
            <el-button slot="reference" type="text" size="small" icon="el-icon-delete" @click="handleSingleDelete(scope.row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
      <pagination
          v-show="total>0"
          v-model:page="queryParams.pageIndex"
          v-model:limit="queryParams.pageSize"
          :total="total"
          @pagination="getList"
      />
    </div>

    <el-dialog v-model="preview.open" :title="preview.title" fullscreen append-to-body destroy-on-close>
      <div class="el-dialog-container">
        <!--        <div class="tag-group">
          &lt;!&ndash; eslint-disable-next-line vue/valid-v-for &ndash;&gt;
          <el-tag v-for="(value, key) in preview.data" @click="codeChange(key)">
            &lt;!&ndash;            {{ key.substring(key.lastIndexOf('/') + 1, key.indexOf('.go.template')) }}&ndash;&gt;
            <template #default="">
              {{ key.substring(key.lastIndexOf('/') + 1, key.indexOf('.go.template')) }}
            </template>
          </el-tag>
        </div>
        <div>
          <Codemirror ref="cmEditor" v-model:value="codestr" :options="cmOptions" border />
        </div>-->
        <el-tabs v-model="preview.activeName" tab-position="left" @tab-click="codeChange">
          <el-tab-pane
              v-for="(value, key) in preview.data"
              :key="key"
              :label="key.substring(key.lastIndexOf('/')+1,key.indexOf('.template'))"
              :name="key"
          />
          <pre class="pre">
          <Codemirror v-model:value="codestr" :options="cmOptions" border />
          </pre>
        </el-tabs>
      </div>

    </el-dialog>
    <ImportTable ref="importTable" @ok="handleQuery" />
    <EditTable ref="editTable" @ok="handleQuery" />
  </div>
</template>

<script>
export default {
  name: 'Gen',
}
</script>

<script setup>
import { listTable, previewTable, delTable, toDBTable, toProjectTable, downloadCode } from '@/api/sys/tools/gen'
import ImportTable from '@/view/sys/tools/gen/importTable.vue'
import EditTable from '@/view/sys/tools/gen/editTable.vue'
import { resolveBlob } from '@/utils/zipdownload'
import Codemirror from 'codemirror-editor-vue3'
import 'codemirror/mode/nginx/nginx.js'
import 'codemirror/theme/dracula.css'
import { ElMessage, ElMessageBox } from 'element-plus'

import 'codemirror/mode/javascript/javascript'
import 'codemirror/mode/go/go'
import 'codemirror/mode/vue/vue'
import { getCurrentInstance, ref } from 'vue'

const { proxy } = getCurrentInstance()

const cmOptions = ref({
  mode: 'text/x-nginx-conf', // Language mode
  theme: 'dracula', // Theme
  lineNumbers: true, // Show line number
  smartIndent: true, // Smart indent
  indentUnit: 4, // The smart indent unit is 2 spaces in length
  foldGutter: true, // Code folding
  matchBrackets: true,
  autoCloseBrackets: true,
  styleActiveLine: true, // Display the style of the selected row
  readOnly: true,
})

const queryParams = ref({
  pageIndex: 1,
  pageSize: 10,
  tableName: undefined,
  tableComment: undefined
})

// 预览参数
const activateName = 'static/template/api.go.template'
const preview = ref({
  open: false,
  title: '代码预览',
  data: {},
  activeName: activateName
})

const codestr = ref('')
// 遮罩层
const loading = ref(true)
// 选中数组
const ids = ref([])
// 总条数
const total = ref(0)
// 表数据
const tableList = ref([])
// 日期范围
const dateRange = ref([])

/** 重置按钮操作 */
const resetQuery = () => {
  dateRange.value = []
  proxy.resetForm('searchForm')
  handleQuery()
}

/** 查询表集合 */
const getList = () => {
  loading.value = true
  listTable(proxy.addDateRange(queryParams.value, dateRange.value)).then(response => {
        tableList.value = response.data.list
        total.value = response.data.count
        loading.value = false
      }
  )
}

getList()

const codeChange = (tab, event) => {
  let key = tab
  if (tab.props !== undefined) {
    key = tab.props.name
  }
  if (key.indexOf('js') > -1) {
    cmOptions.value.mode = 'text/javascript'
  }
  if (key.indexOf('model') > -1 || key.indexOf('router') > -1 || key.indexOf('api') > -1 || key.indexOf('service') > -1 || key.indexOf('dto') > -1) {
    cmOptions.value.mode = 'text/x-go'
  }
  if (key.indexOf('vue') > -1) {
    cmOptions.value.mode = 'text/x-vue'
  }
  codestr.value = preview.value.data[key]
}

/** 搜索按钮操作 */
const handleQuery = (e) => {
  queryParams.value.pageIndex = 1
  queryParams.value.pageSize = 10
  getList()
}

/** 生成代码操作 */
const handleDownloadCode = (row) => {
  ElMessageBox.confirm('确认下载？', '提示', {
    confirmButtonText: '确认',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(() => {
    downloadCode(row.tableId).then(response => {
      resolveBlob(response, 'application/zip', 'code.zip')
    })
  }).catch(() => {
  })
}

const importTable = ref(null)
/** 打开导入表弹窗 */
const openImportTable = () => {
  importTable.value.show()
}

const editTable = ref(null)
/** 打开编辑弹窗 */
const openEditTable = (row) => {
  const tableId = row.tableId
  editTable.value.show(tableId)
}

/** 修改按钮操作 */
/* const handleEditTable = (row) => {

  router.push({
    name: 'editTable', params: {
      tableId: tableId
    }
  })
}*/

/** 预览按钮 */
const handlePreview = async(row) => {
  const res = await previewTable(row.tableId)
  if (res.code === 200) {
    preview.value.data = res.data
    preview.value.open = true
    preview.value.activeName = activateName
    codeChange(activateName)
  } else {
    ElMessage({
      type: 'error',
      message: res.msg
    })
  }
}

const handleToProject = async(row) => {
  ElMessageBox.confirm('正在使用代码生成请确认', '提示', {
    confirmButtonText: '生成',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async() => {
    const res = await toProjectTable(row.tableId, false)
    ElMessage({
      type: 'success',
      message: res.msg
    })
  }).catch(() => {
    ElMessage({
      type: 'info',
      message: '取消操作'
    })
  })
}

const handleToDB = async(row) => {
  ElMessageBox.confirm('正在使用【菜单以及API生成到数据库】请确认?', '提示', {
    confirmButtonText: '写入DB',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async() => {
    const res = await toDBTable(row.tableId)
    if (res.code === 200) {
      ElMessage({
        type: 'success',
        message: res.msg
      })
    }
  }).catch(() => {
    ElMessage({
      type: 'info',
      message: '取消操作'
    })
  })
}

const handleSingleDelete = async(row) => {
  ElMessageBox.confirm('确认删除数据项？', '提示', {
    confirmButtonText: '删除',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async() => {
    const myids = row.tableId || ids.value
    const res = await delTable(myids)
    let resType = 'error'
    if (res.code === 200) {
      resType = 'success'
      open.value = false
      getList()
    }
    ElMessage({ type: resType, message: res.msg })
  }).catch(() => {
    ElMessage({
      type: 'info',
      message: '取消操作'
    })
  })
}
</script>

<style lang="scss" scoped>
.el-dialog-container :v-deep(.el-dialog__body) {
  height: 600px;
  overflow: hidden;

.el-scrollbar__view {
  height: 100%;
}

.pre {
  height: 546px;
  overflow: hidden;
}

.el-scrollbar {
  height: 100%;
}

.el-scrollbar__wrap::-webkit-scrollbar {
  display: none;
}

}
:v-deep(.el-dialog__body) {
  padding: 0 20px;
  margin: 0;
}

.tag-group .el-tag {
  margin-left: 10px;
}

</style>

<style lang="scss">
#codemirror {
  height: auto;
  margin: 0;
  overflow: auto;
}

.CodeMirror {
  overflow: auto;
  border: 1px solid #eee;
  height: 600px;
}
</style>

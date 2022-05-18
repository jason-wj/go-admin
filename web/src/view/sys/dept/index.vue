<template>
  <div>
    <div class="gva-search-box">
      <el-form ref="searchForm" :model="queryParams" :inline="true" label-position="left">
        <el-form-item label="部门名称" prop="deptName">
          <el-input v-model="queryParams.deptName" placeholder="请输入部门名称" clearable size="small" @keyup.enter.native="handleQuery" />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="queryParams.status" placeholder="部门状态" clearable size="small">
            <el-option
              v-for="dict in statusOptions"
              :key="dict.dictValue"
              :label="dict.dictLabel"
              :value="dict.dictValue"
            />
          </el-select>
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
      </el-row>
    </div>
    <el-table v-loading="loading" :data="deptList" row-key="deptId" default-expand-all border :tree-props="{children: 'children', hasChildren: 'hasChildren'}">
      <el-table-column prop="deptName" label="部门名称" />
      <el-table-column prop="sort" align="center" label="排序" width="200" />
      <el-table-column prop="status" label="状态" :formatter="statusFormat" width="100">
        <template #default="scope">
          <el-tag
            :type="scope.row.status === 1 ? 'danger' : 'success'"
            disable-transitions
          >{{ statusFormat(scope.row) }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="创建时间" align="center" prop="createdAt" width="200">
        <template #default="scope">
          <span>{{ formatDate(scope.row.createdAt) }}</span>
        </template>
      </el-table-column>
      <el-table-column label="操作" align="center" class-name="small-padding fixed-width">
        <template #default="scope">
          <el-button
            size="mini"
            type="text"
            icon="edit"
            @click="handleUpdate(scope.row)"
          >修改
          </el-button>
          <el-button
            size="mini"
            type="text"
            icon="plus"
            @click="handleAdd(scope.row)"
          >新增
          </el-button>
          <el-button
            v-if="scope.row.p_id !== 0"
            size="mini"
            type="text"
            icon="delete"
            @click="handleDelete(scope.row)"
          >删除
          </el-button>
        </template>
      </el-table-column>
    </el-table>

    <!-- 添加或修改部门对话框 -->
    <el-dialog v-model="open" :title="title" width="600px" append-to-body destroy-on-close>
      <el-form ref="inputform" :model="form" :rules="rules" label-width="80px">
        <el-row>
          <el-col :span="24">
            <el-form-item label="上级部门" prop="parentId">
              <treeselect
                v-model="form.parentId"
                :options="deptOptions"
                :normalizer="normalizer"
                :show-count="true"
                placeholder="选择上级部门"
                :is-disabled="isEdit"
              />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="部门名称" prop="deptName">
              <el-input v-model="form.deptName" placeholder="请输入部门名称" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="显示排序" prop="orderNum">
              <el-input-number v-model="form.sort" controls-position="right" :min="0" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="负责人" prop="leader">
              <el-input v-model="form.leader" placeholder="请输入负责人" maxlength="20" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="联系电话" prop="phone">
              <el-input v-model="form.phone" placeholder="请输入联系电话" maxlength="11" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="邮箱" prop="email">
              <el-input v-model="form.email" placeholder="请输入邮箱" maxlength="50" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="部门状态">
              <el-radio-group v-model="form.status">
                <el-radio
                  v-for="dict in statusOptions"
                  :key="dict.dictValue"
                  :label="dict.dictValue"
                >{{ dict.dictLabel }}
                </el-radio>
              </el-radio-group>
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
  name: 'Dept',
}
</script>

<script setup>
import { getDeptList, getDept, delDept, addDept, updateDept } from '@/api/sys/dept'
import Treeselect from 'vue3-treeselect'
import 'vue3-treeselect/dist/vue3-treeselect.css'
import { getCurrentInstance, nextTick, ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'

const { proxy } = getCurrentInstance()

// 遮罩层
const loading = ref(true)
// 表格树数据
const deptList = ref([])
// 部门树选项
const deptOptions = ref([])
// 弹出层标题
const title = ref('')
// 是否编辑
const isEdit = ref(false)
// 是否显示弹出层
const open = ref(false)
// 类型数据字典
const statusOptions = ref([])
// 查询参数
const queryParams = ref({
  deptName: undefined,
  status: undefined
})
// 表单参数
const form = ref({})
// 表单校验
const rules = ref({
  parentId: [{ required: true, message: '上级部门不能为空', trigger: 'blur' }],
  deptName: [{ required: true, message: '部门名称不能为空', trigger: 'blur' }],
  sort: [{ required: true, message: '菜单顺序不能为空', trigger: 'blur' }],
  email: [{ type: 'email', message: "'请输入正确的邮箱地址", trigger: ['blur', 'change'] }],
  phone: [{ pattern: /^1[3|4|5|6|7|8|9][0-9]\d{8}$/, message: '请输入正确的手机号码', trigger: 'blur' }]
})

// ref form
const inputform = ref(null)
const searchForm = ref(null)

const init = () => {
  getList()
  proxy.getDicts('sys_normal_disable').then(response => {
    statusOptions.value = response.data
  })
}

/** 查询参数列表 */
const getList = () => {
  loading.value = true
  getDeptList(queryParams.value).then(response => {
    deptList.value = response.data
    loading.value = false
  })
}

/** 转换部门数据结构 */
const normalizer = (node) => {
  if (node.children && !node.children.length) {
    delete node.children
  }
  return {
    id: node.deptId,
    label: node.deptName,
    children: node.children
  }
}

/** 查询部门下拉树结构 */
const getTreeselect = (e) => {
  getDeptList().then(response => {
    deptOptions.value = []

    if (e === 'update') {
      const dept = { deptId: 0, deptName: '主类目', children: [], isDisabled: true }
      dept.children = response.data
      deptOptions.value.push(dept)
    } else {
      const dept = { deptId: 0, deptName: '主类目', children: [] }
      dept.children = response.data
      deptOptions.value.push(dept)
    }
  })
}

// 字典
const statusFormat = (row) => {
  return proxy.selectDictLabel(statusOptions.value, row.status)
}

// 取消按钮
const cancel = () => {
  open.value = false
  reset()
}

// 表单重置
const reset = () => {
  form.value = {
    deptId: undefined,
    parentId: undefined,
    deptName: undefined,
    sort: 10,
    leader: undefined,
    phone: undefined,
    email: undefined,
    status: '0'
  }
  nextTick(() => {
    if (inputform.value != null) {
      inputform.value.resetFields()
    }
  })
}

/** 重置按钮操作 */
const resetQuery = () => {
  searchForm.value.resetFields()
  handleQuery()
}

/** 搜索按钮操作 */
const handleQuery = () => {
  getList()
}

/** 新增按钮操作 */
const handleAdd = (row) => {
  open.value = true
  title.value = '添加部门'
  isEdit.value = false
  reset()

  getTreeselect('add')
  if (row !== undefined) {
    form.value.parentId = row.deptId
  }
}

/** 修改按钮操作 */
const handleUpdate = (row) => {
  getTreeselect('update')

  getDept(row.deptId).then(response => {
    open.value = true
    title.value = '修改部门'
    isEdit.value = true
    reset()

    form.value = response.data
    form.value.status = String(form.value.status)
    form.value.sort = String(form.value.sort)
  })
}

/** 提交按钮 */
const submitForm = () => {
  inputform.value.validate(valid => {
    if (valid) {
      form.value.status = parseInt(form.value.status)
      // form.value.sort = parseInt(form.value.sort)
      if (form.value.deptId !== undefined) {
        updateDept(form.value, form.value.deptId).then(response => {
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
        addDept(form.value).then(response => {
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
  const ids = row.deptId

  ElMessageBox.confirm('是否确认删除编号为"' + ids + '"的数据项?', '警告', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(() => {
    return delDept({ 'ids': ids })
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

init()
</script>

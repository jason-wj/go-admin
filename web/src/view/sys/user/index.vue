<template>
  <div>
    <!--部门数据-->
    <!--    <el-col :span="4" :xs="24">
      <div class="head-container">
        <el-input
          v-model="deptName"
          placeholder="请输入部门名称"
          clearable
          size="small"
          prefix-icon="el-icon-search"
          style="margin-bottom: 20px"
        />
      </div>
      <div class="head-container">
        <el-tree
          ref="tree"
          :data="deptOptions"
          :props="defaultProps"
          :expand-on-click-node="false"
          :filter-node-method="filterNode"
          default-expand-all
          @node-click="handleNodeClick"
        />
      </div>
    </el-col>-->

    <!--用户数据-->
    <div class="gva-search-box">
      <el-form ref="searchForm" :model="queryParams" :inline="true" label-width="68px">
        <el-form-item label="用户名称" prop="username">
          <el-input
            v-model="queryParams.username"
            placeholder="请输入用户名称"
            clearable
            size="small"
            style="width: 160px"
            @keyup.enter.native="handleQuery"
          />
        </el-form-item>
        <el-form-item label="手机号码" prop="phone">
          <el-input
            v-model="queryParams.phone"
            placeholder="请输入手机号码"
            clearable
            size="small"
            style="width: 160px"
            @keyup.enter.native="handleQuery"
          />
        </el-form-item>
        <el-form-item label="状态" prop="status">
          <el-select
            v-model="queryParams.status"
            placeholder="用户状态"
            clearable
            size="small"
            style="width: 160px"
          >
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
          <el-button
            type="primary"
            icon="plus"
            size="small"
            @click="handleAdd"
          >新增
          </el-button>
        </el-col>
      </el-row>
    </div>
    <el-table v-loading="loading" :data="userList" border>
      <el-table-column label="序号" type="index" align="center" min-width="60">
        <template #default="scope">
          <span>{{ (queryParams.pageIndex - 1) * queryParams.pageSize + scope.$index + 1 }}</span>
        </template>
      </el-table-column>
      <el-table-column label="编号" width="75" align="center" prop="userId" />
      <el-table-column label="登录名" width="105" align="center" prop="username" :show-overflow-tooltip="true" />
      <el-table-column label="昵称" prop="nickName" align="center" :show-overflow-tooltip="true" />
      <el-table-column label="部门" prop="dept.deptName" align="center" :show-overflow-tooltip="true" />
      <el-table-column label="手机号" prop="phone" min-width="100" align="center" />
      <el-table-column min-width="100" label="状态" align="center" prop="status" :formatter="statusFormat">
        <template #default="scope">
          {{ statusFormat(scope.row) }}
        </template>
      </el-table-column>
      <el-table-column label="创建时间" prop="createdAt" min-width="130" align="center">
        <template #default="scope">
          <span>{{ formatDate(scope.row.createdAt) }}</span>
        </template>
      </el-table-column>
      <el-table-column label="操作" min-width="180" fix="right" class-name="small-padding fixed-width">
        <template #default="scope">
          <el-button size="small" type="text" icon="edit" @click="handleUpdate(scope.row)">修改</el-button>
          <el-button
            v-if="scope.row.userId !== 1"
            size="small"
            type="text"
            icon="delete"
            @click="handleDelete(scope.row)"
          >删除</el-button>
          <el-button size="small" type="text" icon="el-icon-key" @click="handleResetPwd(scope.row)">重置密码</el-button>
        </template>
      </el-table-column>
    </el-table>

    <Pagination v-show="total>0" v-model:page="queryParams.pageIndex" v-model:limit="queryParams.pageSize" :total="total" @pagination="getList" />

    <!-- 添加或修改参数配置对话框 -->
    <el-dialog v-model="open" :title="title" width="600px" append-to-body destroy-on-close>
      <el-form ref="inputform" :model="form" :rules="rules" label-width="80px">
        <el-row>
          <el-col :span="12">
            <el-form-item label="用户昵称" prop="nickName">
              <el-input v-model="form.nickName" placeholder="请输入用户昵称" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="归属部门" prop="deptId">
              <treeselect
                v-model="form.deptId"
                :options="deptOptions"
                placeholder="请选择归属部门"
              />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="手机号码" prop="phone">
              <el-input v-model="form.phone" placeholder="请输入手机号码" maxlength="11" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="邮箱" prop="email">
              <el-input v-model="form.email" placeholder="请输入邮箱" maxlength="50" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="用户名称" prop="username">
              <el-input v-model="form.username" placeholder="请输入用户名称" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item v-if="form.userId === undefined" label="用户密码" prop="password">
              <el-input v-model="form.password" placeholder="请输入用户密码" type="password" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="用户性别">
              <el-select v-model="form.sex" placeholder="请选择">
                <el-option
                  v-for="dict in sexOptions"
                  :key="dict.dictValue"
                  :label="dict.dictLabel"
                  :value="dict.dictValue"
                />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="状态">
              <el-radio-group v-model="form.status">
                <el-radio
                  v-for="dict in statusOptions"
                  :key="dict.dictValue"
                  :label="dict.dictValue"
                >{{ dict.dictLabel }}</el-radio>
              </el-radio-group>
            </el-form-item>
          </el-col>

          <el-col :span="12">
            <el-form-item label="岗位">
              <el-select v-model="form.postId" placeholder="请选择" @change="$forceUpdate()">
                <el-option
                  v-for="item in postOptions"
                  :key="item.postId"
                  :label="item.postName"
                  :value="item.postId"
                  :disabled="item.status === 1"
                />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="角色">
              <el-select v-model="form.roleId" placeholder="请选择" @change="$forceUpdate()">
                <el-option
                  v-for="item in roleOptions"
                  :key="item.roleId"
                  :label="item.roleName"
                  :value="item.roleId"
                  :disabled="item.status === 1"
                />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="24">
            <el-form-item label="备注">
              <el-input v-model="form.remark" type="textarea" placeholder="请输入内容" />
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
  name: 'User',
}
</script>

<script setup>
import { listUser, getUser, delUser, addUser, updateUser, resetUserPwd } from '@/api/sys/user'
import { listPost } from '@/api/sys/post'
import Treeselect from 'vue3-treeselect'
import 'vue3-treeselect/dist/vue3-treeselect.css'
import { listRole } from '@/api/sys/role'
import { treeselect as deptSelect } from '@/api/sys/dept'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ref, watch } from 'vue'
import { getCurrentInstance, nextTick } from 'vue'

// 遮罩层
const loading = ref(true)
// 总条数
const total = ref(0)
// 数据列表
const userList = ref([])
// 弹出层标题
const title = ref('')
// 部门树选项
const deptOptions = ref([])
// 是否显示弹出层
const open = ref(false)
// 部门名称
const deptName = ref('')
// 默认密码
const initPassword = ref('')
// 日期范围
const dateRange = ref([])
// 类型数据字典
const statusOptions = ref([])
// 性别状态字典
const sexOptions = ref([])
// 岗位选项
const postOptions = ref([])
// 角色选项
const roleOptions = ref([])
// 表单参数
const form = ref({})

const defaultProps = ref({
  children: 'children',
  label: 'label'
})

// 查询参数
const queryParams = ref({
  pageIndex: 1,
  pageSize: 10,
  username: undefined,
  phone: undefined,
  status: undefined,
  deptId: undefined
})

// 表单校验
const rules = ref({
  username: [{ required: true, message: '用户名称不能为空', trigger: 'blur' }],
  nickName: [{ required: true, message: '用户昵称不能为空', trigger: 'blur' }],
  deptId: [{ required: true, message: '归属部门不能为空', trigger: 'blur' }],
  password: [{ required: true, message: '用户密码不能为空', trigger: 'blur' }],
  email: [
    { required: true, message: '邮箱地址不能为空', trigger: 'blur' },
    { type: 'email', message: "'请输入正确的邮箱地址", trigger: ['blur', 'change'] }
  ],
  phone: [
    { required: true, message: '手机号码不能为空', trigger: 'blur' },
    { pattern: /^1[3|4|5|6|7|8|9][0-9]\d{8}$/, message: '请输入正确的手机号码', trigger: 'blur' }
  ]
})

// ref form
const inputform = ref(null)
const searchForm = ref(null)
const tree = ref(null)

const { proxy } = getCurrentInstance()

watch(deptName, () => {
  // 根据名称筛选部门树
  tree.value.filter(deptName.value)
})

const init = () => {
  getList()
  getTreeselect()
  proxy.getDicts('sys_normal_disable').then(response => {
    statusOptions.value = response.data
  })
  proxy.getDicts('sys_user_sex').then(response => {
    sexOptions.value = response.data
  })
  proxy.getConfigKey('sys_user_initPassword').then(response => {
    initPassword.value = response.data.configValue
  })
}

// 字典
const statusFormat = (row) => {
  return proxy.selectDictLabel(statusOptions.value, row.status)
}

/** 查询参数列表 */
const getList = () => {
  loading.value = true
  listUser(proxy.addDateRange(queryParams.value, dateRange.value)).then(response => {
    userList.value = response.data.list
    total.value = response.data.count
    loading.value = false
  })
}

/** 查询部门下拉树结构 */
const getTreeselect = () => {
  deptSelect().then(response => {
    deptOptions.value = response.data
  })
}

// 筛选节点
const filterNode = (value, data) => {
  if (!value) return true
  return data.label.indexOf(value) !== -1
}

// 节点单击事件
const handleNodeClick = (data) => {
  queryParams.value.deptId = '/' + data.id + '/'
  getList()
}

// 取消按钮
const cancel = () => {
  open.value = false
  reset()
}

// 表单重置
const reset = () => {
  form.value = {
    userId: undefined,
    deptId: undefined,
    username: undefined,
    nickName: undefined,
    password: undefined,
    phone: undefined,
    email: undefined,
    sex: undefined,
    status: '2',
    remark: undefined,
    postIds: undefined,
    roleIds: undefined
  }
  nextTick(() => {
    if (inputform.value != null) {
      inputform.value.resetFields()
    }
  })
}

/** 搜索按钮操作 */
const handleQuery = () => {
  queryParams.value.pageIndex = 1
  getList()
}

/** 重置按钮操作 */
const resetQuery = () => {
  queryParams.value.deptId = undefined
  dateRange.value = []
  searchForm.value.resetFields()
  handleQuery()
}

/** 新增按钮操作 */
const handleAdd = () => {
  open.value = true
  title.value = '添加用户'
  reset()
  form.value.password = initPassword

  getTreeselect()

  listPost({ pageSize: 1000 }).then(response => {
    postOptions.value = response.data.list
  })
  listRole({ pageSize: 1000 }).then(response => {
    roleOptions.value = response.data.list
  })
}

/** 修改按钮操作 */
const handleUpdate = (row) => {
  const userId = row.userId
  getUser(userId).then(response => {
    open.value = true
    title.value = '修改用户'
    reset()
    form.value = response.data
    form.value.password = ''
  })
  listPost({ pageSize: 1000 }).then(response => {
    postOptions.value = response.data.list
  })
  listRole({ pageSize: 1000 }).then(response => {
    roleOptions.value = response.data.list
  })
}

/** 重置密码按钮操作 */
const handleResetPwd = (row) => {
  ElMessageBox.prompt('请输入"' + row.username + '"的新密码', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消'
  }).then(({ value }) => {
    resetUserPwd(row.userId, value).then(response => {
      if (response.code === 200) {
        ElMessage({
          type: 'success',
          message: '修改成功，新密码是：' + value,
          showClose: true,
        })
      }
    })
  }).catch(() => {})
}

/** 提交按钮 */
const submitForm = () => {
  inputform.value.validate(valid => {
    if (valid) {
      if (form.value.userId !== undefined) {
        updateUser(form.value).then(response => {
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
        addUser(form.value).then(response => {
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
  const ids = [row.userId]
  ElMessageBox.confirm('是否确认删除编号为"' + ids + '"的数据项?', '警告', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(() => {
    return delUser({ 'ids': ids })
  }).then((response) => {
    if (response.code === 200) {
      open.value = false
      ElMessage({
        type: 'success',
        message: response.msg,
        showClose: true,
      })
      getList()
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

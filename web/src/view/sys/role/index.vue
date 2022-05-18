<template>
  <div>
    <div class="gva-search-box">
      <el-form ref="searchForm" :model="queryParams" :inline="true" label-position="left">
        <el-form-item label="名称" prop="roleName">
          <el-input
            v-model="queryParams.roleName"
            placeholder="请输入角色名称"
            clearable
            size="small"
            style="width: 160px"
            @keyup.enter.native="handleQuery"
          />
        </el-form-item>
        <el-form-item label="权限字符" prop="roleKey">
          <el-input
            v-model="queryParams.roleKey"
            placeholder="请输入权限字符"
            clearable
            size="small"
            style="width: 160px"
            @keyup.enter.native="handleQuery"
          />
        </el-form-item>
        <el-form-item label="状态" prop="status">
          <el-select
            v-model="queryParams.status"
            placeholder="角色状态"
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
    <el-table v-loading="loading" :data="roleList" border>
      <el-table-column label="序号" type="index" align="center" min-width="60">
        <template #default="scope">
          <span>{{ (queryParams.pageIndex - 1) * queryParams.pageSize + scope.$index + 1 }}</span>
        </template>
      </el-table-column>
      <el-table-column label="编码" sortable="custom" prop="roleId" min-width="80" align="center" />
      <el-table-column label="名称" sortable="custom" prop="roleName" align="center" :show-overflow-tooltip="true" />
      <el-table-column label="权限字符" prop="roleKey" :show-overflow-tooltip="true" align="center" min-width="100" />
      <el-table-column label="排序" sortable="custom" prop="roleSort" align="center" min-width="100" />

      <el-table-column prop="status" label="状态" :formatter="statusFormat" align="center" min-width="100">
        <template #default="scope">
          <el-tag
            :type="scope.row.status === 1 ? 'danger' : 'success'"
            disable-transitions
          >{{ statusFormat(scope.row) }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="创建时间" align="center" prop="createdAt" min-width="130">
        <template #default="scope">
          <span>{{ formatDate(scope.row.createdAt) }}</span>
        </template>
      </el-table-column>
      <el-table-column label="操作" align="center" class-name="small-padding fixed-width" min-width="130">
        <template #default="scope">
          <el-button
            size="small"
            type="text"
            icon="edit"
            @click="handleUpdate(scope.row)"
          >修改
          </el-button>
          <el-button
            size="small"
            type="text"
            icon="check"
            @click="handleDataScope(scope.row)"
          >数据权限
          </el-button>
          <el-button
            v-if="scope.row.roleKey!=='admin'"
            size="small"
            type="text"
            icon="delete"
            @click="handleDelete(scope.row)"
          >删除
          </el-button>
        </template>
      </el-table-column>
    </el-table>

    <pagination
      v-show="total>0"
      :total="total"
      :page.sync="queryParams.pageIndex"
      :limit.sync="queryParams.pageSize"
      @pagination="getList"
    />
    <!-- 添加或修改角色配置对话框 -->
    <el-dialog v-if="open" v-model="open" :title="title" width="500px" append-to-body destroy-on-close>
      <el-form ref="inputform" :model="form" :rules="rules" label-width="80px">
        <el-form-item label="角色名称" prop="roleName">
          <el-input v-model="form.roleName" placeholder="请输入角色名称" :disabled="isEdit" />
        </el-form-item>
        <el-form-item label="权限字符" prop="roleKey">
          <el-input v-model="form.roleKey" placeholder="请输入权限字符" :disabled="isEdit" />
        </el-form-item>
        <el-form-item label="角色顺序" prop="roleSort">
          <el-input-number v-model="form.roleSort" controls-position="right" :min="0" />
        </el-form-item>
        <el-form-item label="状态">
          <el-radio-group v-model="form.status">
            <el-radio
              v-for="dict in statusOptions"
              :key="dict.dictValue"
              :label="dict.dictValue"
            >{{ dict.dictLabel }}
            </el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="菜单权限">
          <el-tree v-if="menuOptions.length>0" ref="menuTree" :empty-text="menuOptionsAlert" :data="menuOptions" show-checkbox node-key="id" check-strictly />
          <div v-if="menuOptions.length<=0">{{ menuOptionsAlert }}</div>
        </el-form-item>
        <el-form-item label="备注">
          <el-input v-model="form.remark" type="textarea" placeholder="请输入内容" />
        </el-form-item>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button type="primary" @click="submitForm">确 定</el-button>
        <el-button @click="cancel">取 消</el-button>
      </div>
    </el-dialog>

    <!-- 分配角色数据权限对话框 -->
    <el-dialog v-if="openDataScope" v-model="openDataScope" :title="title" width="500px" append-to-body destroy-on-close>
      <el-form :model="form" label-width="80px">
        <el-form-item label="角色名称">
          <el-input v-model="form.roleName" :disabled="true" />
        </el-form-item>
        <el-form-item label="权限字符">
          <el-input v-model="form.roleKey" :disabled="true" />
        </el-form-item>
        <el-form-item label="权限范围">
          <el-select v-model="form.dataScope">
            <el-option
              v-for="item in dataScopeOptions"
              :key="item.value"
              :label="item.label"
              :value="item.value"
            />
          </el-select>
        </el-form-item>
        <el-form-item v-show="form.dataScope === 2" label="数据权限">
          <el-tree
            ref="dept"
            :data="deptOptions"
            show-checkbox
            default-expand-all
            node-key="id"
            empty-text="加载中，请稍后"
            :props="defaultProps"
          />
        </el-form-item>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button type="primary" @click="submitDataScope">确 定</el-button>
        <el-button @click="cancelDataScope">取 消</el-button>
      </div>
    </el-dialog>
  </div>
</template>

<script>
export default {
  name: 'Role',
}
</script>

<script setup>
import { listRole, getRole, delRole, addRole, updateRole, dataScope } from '@/api/sys/role'
import { roleMenuTreeselect } from '@/api/sys/menu'
import { treeselect as deptTreeselect, roleDeptTreeselect } from '@/api/sys/dept'

import { nextTick, getCurrentInstance, ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'

const { proxy } = getCurrentInstance()

// 遮罩层
const loading = ref(true)
// 总条数
const total = ref(0)
// 弹出层标题
const title = ref('')
// 数据列表
const roleList = ref([])
// 是否显示弹出层
const open = ref(false)
// 是否显示弹出层（数据权限）
const openDataScope = ref(false)
// 是否编辑
const isEdit = ref(false)
// 日期范围
const dateRange = ref([])
// 类型数据字典
const statusOptions = ref([])
// 菜单列表
const menuOptions = ref([])
// 数据列表
const menuList = ref([])
// 部门列表
const deptOptions = ref([])
// 提示
const menuOptionsAlert = ref('加载中，请稍后')

// 查询参数
const queryParams = ref({
  pageIndex: 1,
  pageSize: 10,
  postName: undefined,
  postCode: undefined,
  status: undefined,
})

// 表单参数
const form = ref({
  sysMenu: []
})

const defaultProps = ref({
  children: 'children',
  label: 'label'
})

// 表单校验
const rules = ref({
  roleName: [{ required: true, message: '角色名称不能为空', trigger: 'blur' }],
  roleKey: [{ required: true, message: '权限字符不能为空', trigger: 'blur' }],
  roleSort: [{ required: true, message: '角色顺序不能为空', trigger: 'blur' }]
})

// ref form
const menuTree = ref(null)
const inputform = ref(null)
const searchForm = ref(null)
const dept = ref(null)

// 数据范围选项
const dataScopeOptions = ref([
  {
    value: '1',
    label: '全部数据权限'
  },
  {
    value: '2',
    label: '自定数据权限'
  },
  {
    value: '3',
    label: '本部门数据权限'
  },
  {
    value: '4',
    label: '本部门及以下数据权限'
  },
  {
    value: '5',
    label: '仅本人数据权限'
  }
])

const init = () => {
  getList()
  getMenuTreeselect()
  proxy.getDicts('sys_normal_disable').then(response => {
    statusOptions.value = response.data
  })
}

/** 查询参数列表 */
const getList = () => {
  loading.value = true
  listRole(proxy.addDateRange(queryParams.value, dateRange.value)).then(response => {
    roleList.value = response.data.list
    total.value = response.data.count
    loading.value = false
  })
}

/** 查询菜单树结构 */
const getMenuTreeselect = () => {
  roleMenuTreeselect(0).then(response => {
    menuOptions.value = response.data.menus
    menuList.value = menuOptions.value
  })
}

// 字典
const statusFormat = (row) => {
  return proxy.selectDictLabel(statusOptions.value, row.status)
}

/** 查询部门树结构 */
const getDeptTreeselect = () => {
  deptTreeselect().then(response => {
    deptOptions.value = response.data.list
  })
}

// 所有菜单节点数据
const getMenuAllCheckedKeys = () => {
  // 目前被选中的菜单节点
  const checkedKeys = menuTree.value.getHalfCheckedKeys()
  console.log('目前被选中的菜单节点', checkedKeys)
  // 半选中的菜单节点
  const halfCheckedKeys = menuTree.value.getCheckedKeys()
  console.log('半选中的菜单节点', halfCheckedKeys)
  checkedKeys.unshift.apply(checkedKeys, halfCheckedKeys)
  console.log('所有选中节点汇总', checkedKeys)
  return checkedKeys
}

// 所有部门节点数据
const getDeptAllCheckedKeys = () => {
  // 目前被选中的部门节点
  const checkedKeys = dept.value.getCheckedKeys()
  // 半选中的部门节点
  const halfCheckedKeys = dept.value.getCheckedKeys()
  checkedKeys.unshift.apply(checkedKeys, halfCheckedKeys)
  return checkedKeys
}

/** 根据角色ID查询菜单树结构 */
const getRoleMenuTreeselect = (row, checkedKeys) => {
  if (row.roleKey === 'admin') {
    menuOptionsAlert.value = '系统超级管理员无需此操作'
    menuOptions.value = []
  } else {
    nextTick(() => {
      menuTree.value.setCheckedKeys(checkedKeys)
    })
  }
}

/** 根据角色ID查询部门树结构 */
const getRoleDeptTreeselect = (roleId) => {
  roleDeptTreeselect(roleId).then(response => {
    deptOptions.value = response.data.depts
    nextTick(() => {
      dept.value.setCheckedKeys(response.data.checkedKeys)
    })
  })
}

// 取消按钮
const cancel = () => {
  open.value = false
  reset()
}

// 取消按钮（数据权限）
const cancelDataScope = () => {
  openDataScope.value = false
  reset()
}

// 表单重置
const reset = () => {
  menuOptions.value = menuList.value
  form.value = {
    roleId: undefined,
    roleName: undefined,
    roleKey: undefined,
    roleSort: 0,
    status: '2',
    menuIds: [],
    deptIds: [],
    sysMenu: [],
    remark: undefined
  }
  nextTick(() => {
    if (menuTree.value !== null) {
      menuTree.value.setCheckedKeys([])
    }
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
  dateRange.value = []
  searchForm.value.resetFields()
  handleQuery()
}

/** 新增按钮操作 */
const handleAdd = () => {
  open.value = true
  title.value = '添加角色'
  isEdit.value = false
  reset()
}

/** 修改按钮操作 */
const handleUpdate = (row) => {
  const roleId = row.roleId
  getRole(roleId).then(response => {
    title.value = '修改角色'
    isEdit.value = true
    open.value = true
    reset()
    form.value = response.data
    getRoleMenuTreeselect(row, response.data.menuIds)
  })
}

/** 分配数据权限操作 */
const handleDataScope = (row) => {
  getRole(row.roleId).then(response => {
    openDataScope.value = true
    title.value = '分配数据权限'
    reset()
    form.value = response.data
    getRoleDeptTreeselect(row.roleId)
  })
}

/** 提交按钮 */
const submitForm = () => {
  inputform.value.validate(valid => {
    if (valid) {
      if (form.value.roleId !== undefined) {
        form.value.menuIds = getMenuAllCheckedKeys()
        updateRole(form.value, form.value.roleId).then(response => {
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
        form.value.menuIds = getMenuAllCheckedKeys()
        addRole(form.value).then(response => {
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

/** 提交按钮（数据权限） */
const submitDataScope = () => {
  if (form.value.roleId !== undefined) {
    form.value.deptIds = getDeptAllCheckedKeys()
    // console.log(this.getDeptAllCheckedKeys())
    dataScope(form.value).then(response => {
      if (response.code === 200) {
        ElMessage({
          type: 'success',
          message: response.msg,
          showClose: true,
        })
        openDataScope.value = false
        getList()
      }
    })
  }
}

/** 删除按钮操作 */
const handleDelete = (row) => {
  const ids = [row.roleId]

  ElMessageBox.confirm('是否确认删除角色编号为"' + ids + '"的数据项?', '警告', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(() => {
    return delRole({ 'ids': ids })
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

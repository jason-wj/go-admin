<template>
  <div>
    <div class="gva-search-box">
      <el-form ref="searchForm" :model="queryParams" :inline="true" label-position="left">
        <el-form-item label="菜单名称">
          <el-input v-model="queryParams.title" placeholder="请输入菜单名称" clearable size="small" @keyup.enter.native="handleQuery" />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="queryParams.visible" placeholder="菜单状态" clearable size="small">
            <el-option
              v-for="dict in visibleOptions"
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
    <el-table v-loading="loading" :data="menuList" border row-key="menuId" :tree-props="{children: 'children', hasChildren: 'hasChildren'}">
      <el-table-column prop="title" label="菜单名称" :show-overflow-tooltip="true" min-width="80" align="center" />
      <el-table-column prop="icon" label="图标" align="center" min-width="40">
        <template #default="scope">
          <svg-icon :icon-class="scope.row.icon" />
        </template>
      </el-table-column>
      <el-table-column prop="sort" label="排序" min-width="40" />
      <el-table-column prop="permission" label="权限标识" :show-overflow-tooltip="true" min-width="80" align="center" />
      <el-table-column prop="component" label="组建路径" :show-overflow-tooltip="true" min-width="180" align="center" />
      <el-table-column prop="hidden" label="可见" :formatter="visibleFormat" min-width="80" align="center">
        <template #default="scope">
          {{ visibleFormat(scope.row) }}
        </template>
      </el-table-column>
      <el-table-column label="创建时间" align="center" prop="createdAt" width="180">
        <template #default="scope">
          <span>{{ formatDate(scope.row.createdAt) }}</span>
        </template>
      </el-table-column>
      <el-table-column label="操作" align="center" class-name="small-padding fixed-width" width="180">
        <template #default="scope">
          <el-button
            size="small"
            type="text"
            icon="el-icon-edit"
            @click="handleUpdate(scope.row)"
          >修改</el-button>
          <el-button
            size="small"
            type="text"
            icon="el-icon-plus"
            @click="handleAdd(scope.row)"
          >新增</el-button>
          <el-button
            size="small"
            type="text"
            icon="el-icon-delete"
            @click="handleDelete(scope.row)"
          >删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <!-- 添加或修改菜单对话框 -->
    <el-dialog v-model="open" :close-on-click-modal="false" :title="title" width="1000px" append-to-body destroy-on-close>
      <el-form ref="inputform" :model="form" :rules="rules" label-position="top" label-width="150px">
        <el-row>
          <el-col :span="24">
            <span slot="label">
              上级菜单
              <el-tooltip content="指当前菜单停靠的菜单归属" placement="top">
                <el-icon><question-filled /></el-icon>
              </el-tooltip>
            </span>
            <el-form-item prop="parentId">
              <treeselect
                v-model="form.parentId"
                :options="menuOptions"
                :normalizer="normalizer"
                :show-count="true"
                placeholder="选择上级菜单"
              />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <span slot="label">
              菜单标题
              <el-tooltip content="菜单位置显示的说明信息" placement="top">
                <el-icon><question-filled /></el-icon>
              </el-tooltip>
            </span>
            <el-form-item prop="title">
              <el-input v-model="form.title" placeholder="请输入菜单标题" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <span slot="label">
              显示排序
              <el-tooltip content="根据序号升序排列" placement="top">
                <el-icon><question-filled /></el-icon>
              </el-tooltip>
            </span>
            <el-form-item prop="sort">
              <el-input-number v-model="form.sort" controls-position="right" :min="0" />
            </el-form-item>
          </el-col>
          <el-col :span="24">
            <span slot="label">
              菜单类型
              <el-tooltip content="包含目录：以及菜单或者菜单组，菜单：具体对应某一个页面，按钮：功能才做按钮；" placement="top">
                <el-icon><question-filled /></el-icon>
              </el-tooltip>
            </span>
            <el-form-item prop="menuType">
              <el-radio-group v-model="form.menuType">
                <el-radio label="M">目录</el-radio>
                <el-radio label="C">菜单</el-radio>
                <el-radio label="F">按钮</el-radio>
              </el-radio-group>
            </el-form-item>
          </el-col>
          <el-col :span="24">
            <el-form-item label="菜单图标">
              <el-popover title="选择图标" placement="top-start" :width="500" trigger="click" @show="iconSelectReset">
                <template #reference>
                  <el-input slot="reference" v-model="form.icon" placeholder="点击选择图标" readonly>
                    <template #prepend>
                      <svg-icon v-if="form.icon" slot="prefix" :icon-class="form.icon" style="height: 32px;width: 16px;" />
                      <svg-icon v-else slot="prefix" icon-class="search" style="height: 32px;width: 16px;" />
                    </template>
                  </el-input>
                </template>
                <IconSelect ref="iconSelect" @selected="selected" />
              </el-popover>
            </el-form-item>
          </el-col>
          <el-col v-if="form.menuType === 'M' || form.menuType === 'C'" :span="12">
            <span slot="label">
              路由名称
              <el-tooltip content="需要和页面name保持一致，对应页面即可选择缓存" placement="top">
                <el-icon><question-filled /></el-icon>
              </el-tooltip>
            </span>
            <el-form-item prop="name">
              <el-input v-model="form.name" placeholder="请输入路由名称" />
            </el-form-item>
          </el-col>

          <el-col v-if="form.menuType === 'M' || form.menuType === 'C'" :span="12">
            <span slot="label">
              组件路径
              <el-tooltip content="菜单对应的具体vue页面文件路径views的下级路径/admin/sys-api/index；目录类型：填写Layout，如何有二级目录请参照日志目录填写；" placement="top">
                <el-icon><question-filled /></el-icon>
              </el-tooltip>
            </span>
            <el-form-item prop="component">
              <el-input v-model="form.component" placeholder="请输入组件路径" />
            </el-form-item>
          </el-col>

          <el-col v-if="form.menuType === 'M' || form.menuType === 'C'" :span="12">
            <span slot="label">
              是否外链
              <el-tooltip content="可以通过iframe打开指定地址" placement="top">
                <el-icon><question-filled /></el-icon>
              </el-tooltip>
            </span>
            <el-form-item>
              <el-radio-group v-model="form.isFrame">
                <el-radio label="0">是</el-radio>
                <el-radio label="1">否</el-radio>
              </el-radio-group>
            </el-form-item>
          </el-col>

          <el-col v-if="form.menuType !== 'F'" :span="12">
            <span slot="label">
              路由地址
              <el-tooltip content="访问此页面自定义的url地址，建议/开头书写，例如 /app-name/menu-name" placement="top">
                <el-icon><question-filled /></el-icon>
              </el-tooltip>
            </span>
            <el-form-item prop="path">
              <el-input v-model="form.path" placeholder="请输入路由地址" />
            </el-form-item>
          </el-col>

          <el-col v-if="form.menuType === 'F' || form.menuType === 'C'" :span="12">
            <span slot="label">
              权限标识
              <el-tooltip content="前端权限控制按钮是否显示" placement="top">
                <el-icon><question-filled /></el-icon>
              </el-tooltip>
            </span>
            <el-form-item>
              <el-input v-model="form.permission" placeholder="请权限标识" maxlength="50" />
            </el-form-item>
          </el-col>
          <el-col v-if="form.menuType !== 'F'" :span="12">
            <span slot="label">
              菜单状态
              <el-tooltip content="需要显示在菜单列表的菜单设置为显示，否则设置为隐藏" placement="top">
                <el-icon><question-filled /></el-icon>
              </el-tooltip>
            </span>
            <el-form-item>
              <el-radio-group v-model="form.visible">
                <el-radio
                  v-for="dict in visibleOptions"
                  :key="dict.dictValue"
                  :label="dict.dictValue"
                >{{ dict.dictLabel }}</el-radio>
              </el-radio-group>
            </el-form-item>
          </el-col>
          <el-col v-if="form.menuType === 'F' || form.menuType === 'C'" :span="24">
            <span slot="label">
              api权限
              <el-tooltip content="配置在这个才做上需要使用到的接口，否则在设置用户角色时，接口将无权访问。" placement="top">
                <el-icon><question-filled /></el-icon>
              </el-tooltip>
            </span>
            <el-form-item>
              <el-transfer
                v-model="form.apis"
                style="text-align: left; display: inline-block"
                filterable
                :props="{
                  key: 'id',
                  label: 'title'
                }"
                :titles="['未授权', '已授权']"
                :button-texts="['收回', '授权 ']"
                :format="{
                  noChecked: '${total}',
                  hasChecked: '${checked}/${total}'
                }"
                class="panel"
                :data="sysapiList"
                @change="handleChange"
              />
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
  name: 'Menu',
}
</script>

<script setup>
import { listMenu, getMenu, delMenu, addMenu, updateMenu } from '@/api/sys/menu'
import { listApi } from '@/api/sys/api'

import Treeselect from 'vue3-treeselect'
import 'vue3-treeselect/dist/vue3-treeselect.css'
import IconSelect from '@/components/IconSelect/index.vue'
import { getCurrentInstance, ref, nextTick } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'

const { proxy } = getCurrentInstance()

// 遮罩层
const loading = ref(true)
// 菜单表格树数据
const menuList = ref([])
// 菜单表格树数据 接口
const sysapiList = ref([])
// 菜单树选项
const menuOptions = ref([])
// 弹出层标题
const title = ref('')
// 是否显示弹出层
const open = ref(false)
// 菜单状态数据字典
const visibleOptions = ref([])
// 查询参数
const queryParams = ref({
  title: undefined,
  visible: undefined
})
// 表单参数
const form = ref({
  apis: [],
  sysApi: []
})

// 表单校验
const rules = ref({
  title: [{ required: true, message: '菜单标题不能为空', trigger: 'blur' }],
  sort: [{ required: true, message: '菜单顺序不能为空', trigger: 'blur' }]
})

// ref form
const inputform = ref(null)
const searchForm = ref(null)
const iconSelect = ref(null)

const init = () => {
  getList()
  getApiList()
  proxy.getDicts('sys_show_hide').then(response => {
    visibleOptions.value = response.data
  })
}

/** 重置按钮操作 */
const resetQuery = () => {
  searchForm.value.resetFields()
  handleQuery()
}

const handleChange = (value, direction, movedKeys) => {
  const list = form.value.sysApi
  form.value.apis = value
  if (direction === 'right') {
    for (let x = 0; x < movedKeys.length; x++) {
      for (let index = 0; index < sysapiList.value.length; index++) {
        const element = sysapiList.value[index]
        if (element.id === movedKeys[x]) {
          list.push(element)
          break
        }
      }
    }
    form.value.sysApi = list
  } else if (direction === 'left') {
    const l = []
    for (let index = 0; index < movedKeys.length; index++) {
      const element = movedKeys[index]
      for (let x = 0; x < list.length; x++) {
        const e = list[x]
        if (element !== e.id) {
          l.push()
          break
        }
      }
    }
    form.value.sysApi = l
  }
}

const getApiList = () => {
  loading.value = true
  listApi({ 'pageSize': 10000, 'type': 'BUS' }).then(response => {
    sysapiList.value = response.data.list
    loading.value = false
  }
  )
}

// 选择图标
const selected = (name) => {
  form.value.icon = name
}

/** 查询菜单列表 */
const getList = () => {
  loading.value = true
  listMenu(queryParams.value).then(response => {
    menuList.value = response.data
    loading.value = false
  })
}

/** 转换菜单数据结构 */
const normalizer = (node) => {
  if (node.children && !node.children.length) {
    delete node.children
  }
  return {
    id: node.menuId,
    label: node.title,
    children: node.children
  }
}

/** 查询菜单下拉树结构 */
const getTreeselect = () => {
  listMenu().then(response => {
    menuOptions.value = []
    const menu = { menuId: 0, title: '主类目', children: [] }
    menu.children = response.data
    menuOptions.value.push(menu)
  })
}

// 菜单显示状态字典翻译
const visibleFormat = (row) => {
  let hidden = 0
  if (row.hidden === true) {
    hidden = 1
  }
  return proxy.selectDictLabel(visibleOptions.value, hidden)
}

const iconSelectReset = () => {
  nextTick(() => {
    if (iconSelect.value != null) {
      iconSelect.value.reset()
    }
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
    menuId: undefined,
    parentId: 0,
    name: undefined,
    icon: undefined,
    menuType: 'M',
    apis: [],
    sort: 0,
    action: undefined,
    isFrame: '1',
    visible: '0'
  }
  nextTick(() => {
    if (inputform.value != null) {
      inputform.value.resetFields()
    }
  })
}

/** 搜索按钮操作 */
const handleQuery = () => {
  getList()
}

/** 新增按钮操作 */
const handleAdd = (row) => {
  open.value = true
  title.value = '添加菜单'
  reset()
  getTreeselect()
  if (row != null) {
    form.value.parentId = row.menuId
  }
}

/** 修改按钮操作 */
const handleUpdate = (row) => {
  getTreeselect()
  console.log(row.menuId)
  getMenu(row.menuId).then(response => {
    open.value = true
    title.value = '修改菜单'
    reset()
    form.value = response.data
  })
}

const setApis = (apiArray) => {
  const l = []
  for (let index = 0; index < apiArray.length; index++) {
    const element = apiArray[index]
    l.push(element.id)
  }
  form.value.apis = l
}

/** 提交按钮 */
const submitForm = () => {
  inputform.value.validate(valid => {
    if (valid) {
      if (form.value.menuId !== undefined) {
        updateMenu(form.value, form.value.menuId).then(response => {
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
        addMenu(form.value).then(response => {
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
  const ids = [row.menuId]

  ElMessageBox.confirm('是否确认删除名称为"' + row.name + '"的数据项?', '警告', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(() => {
    return delMenu({ 'ids': ids })
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

<style lang="css">

.panel {
  margin-left:0;
}
.panel .el-transfer-panel{
  width: 300px;
}

.el-col {
padding: 0 5px;
}
.el-drawer__header{
margin-bottom: 0;
}
</style>

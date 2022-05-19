<template>
  <el-dialog v-model="visible" title="编辑" fullscreen top="5vh" append-to-body destroy-on-close>
    <el-card>
      <el-tabs v-model="activeName">
        <el-tab-pane label="基本信息" name="basic">
          <BasicInfoForm ref="basicInfo" :info="info" />
        </el-tab-pane>
        <el-tab-pane label="字段信息" name="cloum">
          <el-alert
              title="⚠️表字段中的id、create_by、update_by、created_at、updated_at的字段在此列表中已经隐藏"
              type="warning"
              show-icon
          />
          <el-table v-loading="dataLoading" :data="columns" :max-height="tableHeight" style="width: 100%">
            <el-table-column fixed label="序号" type="index" width="50" />
            <el-table-column
                fixed
                label="字段列名"
                prop="columnName"
                width="150"
                :show-overflow-tooltip="true"
            />
            <el-table-column fixed label="字段描述" width="150">
              <template #default="scope">
                <el-input v-model="scope.row.columnComment" />
              </template>
            </el-table-column>
            <el-table-column
                label="物理类型"
                prop="columnType"
                width="120"
                :show-overflow-tooltip="true"
            />
            <el-table-column label="go类型" width="120">
              <template #default="scope">
                <el-select v-model="scope.row.goType">
                  <el-option label="int64" value="int64" />
                  <el-option label="string" value="string" />
                  <el-option label="decimal" value="decimal.Decimal" />
                  <el-option label="time" value="time.Time" />
                  <!-- <el-option label="int" value="int" />
                  <el-option label="bool" value="bool" /> -->
                </el-select>
              </template>
            </el-table-column>
            <el-table-column label="go属性" width="150">
              <template #default="scope">
                <el-input v-model="scope.row.goField" />
              </template>
            </el-table-column>
            <el-table-column label="json属性" width="150">
              <template #default="scope">
                <el-input v-model="scope.row.jsonField" />
              </template>
            </el-table-column>

            <el-table-column label="编辑" width="50">
              <template #default="scope">
                <el-checkbox v-model="scope.row.isInsert" true-label="1" false-label="0" />
              </template>
            </el-table-column>
            <!-- <el-table-column label="编辑" width="70" :render-header="renderHeadeUpdate" :cell-style="{'text-align':'center'}">
            <template slot-scope="scope">
              <el-checkbox v-model="scope.row.isEdit" true-label="1" false-label="0" />
            </template>
          </el-table-column> -->
            <el-table-column label="列表" width="70" :scoped-slot="renderHeadeList" :cell-style="{'text-align':'center'}">
              <template #default="scope">
                <el-checkbox v-model="scope.row.isList" true-label="1" false-label="0" />
              </template>
            </el-table-column>
            <el-table-column label="查询" width="70" :scoped-slot="renderHeadeSearch" :cell-style="{'text-align':'center'}">
              <template #default="scope">
                <el-checkbox v-model="scope.row.isQuery" true-label="1" false-label="0" />
              </template>
            </el-table-column>
            <el-table-column label="查询方式" width="120">
              <template #default="scope">
                <el-select v-model="scope.row.queryType">
                  <el-option label="=" value="EQ" />
                  <el-option label="!=" value="NE" />
                  <el-option label=">" value="GT" />
                  <el-option label=">=" value="GTE" />
                  <el-option label="<" value="LT" />
                  <el-option label="<=" value="LTE" />
                  <el-option label="LIKE" value="LIKE" />
                  <!-- <el-option label="BETWEEN" value="BETWEEN" /> -->
                </el-select>
              </template>
            </el-table-column>
            <el-table-column label="必填" width="50">
              <template #default="scope">
                <el-checkbox v-model="scope.row.isRequired" true-label="1" false-label="0" />
              </template>
            </el-table-column>
            <el-table-column label="显示类型" width="140">
              <template #default="scope">
                <el-select v-model="scope.row.htmlType">
                  <el-option label="文本框" value="input" />
                  <el-option label="下拉框" value="select" />
                  <el-option label="单选框" value="radio" />
                  <!-- <el-option label="文件选择" value="file" /> -->
                  <!-- <el-option label="复选框" value="checkbox" />
                <el-option label="日期控件" value="datetime" />-->
                  <el-option label="文本域" value="textarea" />

                </el-select>
              </template>
            </el-table-column>
            <el-table-column label="字典类型" width="160">
              <template #default="scope">
                <el-select v-model="scope.row.dictType" clearable filterable placeholder="请选择">
                  <el-option
                      v-for="dict in dictOptions"
                      :key="dict.dictType"
                      :label="dict.dictName"
                      :value="dict.dictType"
                  >
                    <span style="float: left">{{ dict.dictName }}</span>
                    <span style="float: right; color: #8492a6; font-size: 13px">{{ dict.dictType }}</span>
                  </el-option>
                </el-select>
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>
        <el-tab-pane label="生成信息" name="genInfo">
          <GenInfoForm ref="genInfo" :info="info" />
        </el-tab-pane>
      </el-tabs>
      <div slot="footer" class="dialog-footer">
        <el-button type="primary" :loading="loading" @click="submitForm">提 交</el-button>
        <el-button @click="visible = false">取 消</el-button>
      </div>
    </el-card>

  </el-dialog>
</template>

<script>
export default {
  name: 'EditTable',
}
</script>

<script setup>
import { ref } from 'vue'
import { getGenTable, updateGenTable, getTableTree } from '@/api/sys/tools/gen'
import { optionselect as getDictOptionselect } from '@/api/sys/dicttype'
import BasicInfoForm from '@/view/sys/tools/gen/basicInfoForm.vue'
import GenInfoForm from '@/view/sys/tools/gen/genInfoForm.vue'
import { ElMessage } from 'element-plus'
const emit = defineEmits(['ok'])

const dataLoading = ref(false)
// 遮罩层
const loading = ref(false)
// dialog可视
const visible = ref(false)

// 选中选项卡的 name
const activeName = ref('cloum')
// 表格的高度
const tableHeight = ref(document.documentElement.scrollHeight - 245 + 'px')
// 表列信息
const columns = ref([])
const tableTree = ref([])
// 字典信息
const dictOptions = ref([])
// 表详细信息
const info = ref({})

const show = (tableId) => {
  visible.value = true
  dataLoading.value = true
  getTableTree().then(response => {
    tableTree.value = response.data
    tableTree.value.unshift({ tableId: 0, className: '请选择' })
  })
  if (tableId) {
    // 获取表详细信息
    getGenTable(tableId).then(res => {
      columns.value = res.data.list
      info.value = res.data.info

      columns.value.forEach(item => {
        tableTree.value.filter(function(e) {
          if (e.tableId === item.fkTableNameClass) {
            item.fkCol = e.columns || [{ columnId: 0, columnName: '请选择' }]
            // item.fkCol.unshift({ columnId: 0, columnName: '请选择' })
          }
        })
      })
      dataLoading.value = false
    })

    /** 查询字典下拉列表 */
    getDictOptionselect().then(response => {
      dictOptions.value = response.data
    })
  }
}

/* const renderHeadeUpdate = (h, { column, $index }) => {
  // h 是一个渲染函数       column 是一个对象表示当前列      $index 第几列
  return h('div', [
    h('span', column.label + '  ', { align: 'center', marginTop: '0px' }),
    h(
      'el-popover',
      { props: { placement: 'top-start', width: '270', trigger: 'hover' }},
      [
        h('p', '是否在表单编辑时能够编辑，打√表示需要', { class: 'text-align: center; margin: 0' }),
        // 生成 i 标签 ，添加icon 设置 样式，slot 必填
        h('i', { class: 'el-icon-question', style: 'color:#ccc,padding-top:5px', slot: 'reference' })
      ]
    )
  ])
}*/

const renderHeadeList = (h, { column, $index }) => {
  // h 是一个渲染函数       column 是一个对象表示当前列      $index 第几列
  return h('div', [
    h('span', column.label + '  ', { align: 'center', marginTop: '0px' }),
    h(
        'el-popover',
        { props: { placement: 'top-start', width: '260', trigger: 'hover' }},
        [
          h('p', '是否在列表中展示，打√表示需要展示', { class: 'text-align: center; margin: 0' }),
          h('i', { class: 'el-icon-question', style: 'color:#ccc,padding-top:5px', slot: 'reference' })
        ]
    )
  ])
}

const renderHeadeSearch = (h, { column, $index }) => {
  return h('div', [
    h('span', column.label + '  ', { align: 'center', marginTop: '0px' }),
    h(
        'el-popover',
        { props: { placement: 'top-start', width: '270', trigger: 'hover' }},
        [
          h('p', '是都当做搜索条件，打√表示做为搜索条件', { class: 'text-align: center; margin: 0' }),
          h('i', { class: 'el-icon-question', style: 'color:#ccc,padding-top:5px', slot: 'reference' })
        ]
    )
  ])
}

/* const handleChangeConfig = (row, index) => {
  tableTree.value.filter(function(item) {
    if (item.tableName === row.fkTableName) {
      row.fkCol = item.columns
      // row.fkCol.unshift({ columnId: 0, columnName: '请选择' })
    }
  })
}*/

const basicInfo = ref(null)
const genInfo = ref(null)

/** 提交按钮 */
const submitForm = () => {
  loading.value = true
  const basicInfoForm = basicInfo.value.basicInfoForm
  const genInfoForm = genInfo.value.genInfoForm
  Promise.all([basicInfoForm, genInfoForm].map(getFormPromise)).then(res => {
    const validateResult = res.every(item => !!item)
    if (validateResult) {
      const genTable = Object.assign({}, basicInfoForm.model, genInfoForm.model)
      genTable.columns = columns.value
      updateGenTable(genTable).then(res => {
        loading.value = false
        let resType = 'error'
        if (res.code === 200) {
          resType = 'success'
          visible.value = false
          emit('ok')
        }
        ElMessage({
          type: resType,
          message: res.msg
        })
      })
    } else {
      loading.value = false
      ElMessage({
        type: 'error',
        message: '表单校验未通过，请重新检查提交内容'
      })
    }
  })
}

/* const getTables = () => {
  getTableTree().then(response => {
    tableTree.value = response.data
    tableTree.value.unshift({ tableId: 0, className: '请选择' })
  })
}*/

/* const getTablesCol = (tableName) => {
  return tableTree.value.filter(function(item) {
    if (item.tableName === tableName) {
      return item.columns
    }
  })
}*/

const getFormPromise = (form) => {
  return new Promise(resolve => {
    form.validate(res => {
      resolve(res)
    })
  })
}

defineExpose({ show })
</script>

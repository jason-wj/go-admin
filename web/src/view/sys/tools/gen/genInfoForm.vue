<template>
  <el-form ref="genInfoForm" :model="info" :rules="rules" label-width="150px">
    <el-row>
      <el-col :span="12">
        <el-form-item prop="tplCategory">
          <template #label=""> 生成模板</template>
          <el-select v-model="info.tplCategory">
            <el-option label="关系表（增删改查）" value="crud" />
          </el-select>
        </el-form-item>
      </el-col>

      <el-col :span="12">
        <el-form-item prop="tplCategory">
          <template #label=""> 是否插件</template>
          <el-select v-model="info.isPlugin">
            <el-option label="否" value="0" />
            <el-option label="是" value="1" />
          </el-select>
        </el-form-item>
      </el-col>

      <el-col :span="12">
        <el-form-item prop="packageName">
          <template #label="">
            应用名
            <el-tooltip content="应用名，例如：在app文件夹下将该功能发到那个应用中，默认：admin" placement="top">
              <el-icon><question-filled /></el-icon>
            </el-tooltip>
          </template>
          <el-input v-model="info.packageName" />
        </el-form-item>
      </el-col>

      <el-col :span="12">
        <el-form-item prop="businessName">
          <template #label="">
            业务名
            <el-tooltip content="可理解为功能英文名，例如 user" placement="top">
              <el-icon><question-filled /></el-icon>
            </el-tooltip>
          </template>
          <el-input v-model="info.businessName" />
        </el-form-item>
      </el-col>

      <el-col :span="12">
        <el-form-item prop="functionName">
          <template #label="">
            功能描述
            <el-tooltip content="同步的数据库表备注，用作类描述，例如：用户" placement="top">
              <el-icon><question-filled /></el-icon>
            </el-tooltip>
          </template>
          <el-input v-model="info.functionName" />
        </el-form-item>
      </el-col>
      <el-col :span="12">
        <el-form-item prop="moduleName">
          <template #label="">
            接口路径
            <el-tooltip content="接口路径，例如：api/v1/{sys-user}" placement="top">
              <el-icon><question-filled /></el-icon>
            </el-tooltip>
          </template>
          <el-input v-model="info.moduleName">
            <template slot="prepend">api/{version}/</template>
            <template slot="append">...</template>
          </el-input>
        </el-form-item>
      </el-col>
    </el-row>
  </el-form>
</template>

<script>
export default {
  name: 'GenInfoForm',
}
</script>

<script setup>
import { ref } from 'vue'

const props = defineProps({
  info: {
    type: Object,
    default: null
  }
})

const rules = ref({
  tplCategory: [
    { required: true, message: '请选择生成模板', trigger: 'blur' }
  ],
  isPlugin: [
    { required: true, message: '请选择是否插件', trigger: 'blur' }
  ],
  packageName: [
    { required: true, message: '请输入生成包路径', trigger: 'blur' },
    { pattern: /^[a-z]*$/g, trigger: 'blur', message: '只允许小写字母,例如 system 格式' }
  ],
  moduleName: [
    { required: true, message: '请输入生成模块名', trigger: 'blur' },
    // { pattern: /^[a-z\-]*[a-z]$/g, trigger: 'blur', message: '只允许小写字母,例如 sys-demo 格式' }
  ],
  businessName: [
    { required: true, message: '请输入生成业务名', trigger: 'blur' },
    { pattern: /^[a-z][A-Za-z]+$/, trigger: 'blur', message: '校验规则:  只允许输入字母 a-z 或大写 A-Z ，并且小写字母开头' }
  ],
  functionName: [
    { required: true, message: '请输入生成功能名', trigger: 'blur' }
  ]
})

const genInfoForm = ref(null)

defineExpose({ genInfoForm })
</script>

<template>
  <div class="content">
    <Codemirror v-model:value="codestr" :options="cmOptions" border />
  </div>
</template>

<script>
export default {
  name: 'Config',
}
</script>

<script setup>
import Codemirror from 'codemirror-editor-vue3'
import 'codemirror/mode/nginx/nginx.js'
import 'codemirror/theme/darcula.css'
import { ElMessage } from 'element-plus'
const codestr = ref('')
import { ref } from 'vue'
import { getConfig } from '@/api/sys/tools/config'
import 'codemirror/mode/yaml/yaml.js'

const cmOptions = ref({
  mode: 'text/yaml', // Language mode
  theme: 'darcula', // Theme
  lineNumbers: true, // Show line number
  smartIndent: true, // Smart indent
  indentUnit: 4, // The smart indent unit is 2 spaces in length
  foldGutter: true, // Code folding
  matchBrackets: true,
  autoCloseBrackets: true,
  styleActiveLine: true, // Display the style of the selected row
  readOnly: true,
})

/** 预览 */
const handlePreview = async() => {
  const res = await getConfig()
  if (res.code === 200) {
    codestr.value = res.data
  } else {
    ElMessage({
      type: 'error',
      message: res.msg
    })
  }
}
handlePreview()
</script>

<style lang="scss">
.content {
  height: 600px;
  overflow: hidden;
}
</style>

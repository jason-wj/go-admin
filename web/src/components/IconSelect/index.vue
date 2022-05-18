<template>
  <div class="icon-body">
    <el-input v-model="name" style="position: relative;" clearable placeholder="请输入图标名称" suffix-icon="search" @clear="filterIcons" @input.native="filterIcons" />
    <div class="icon-list">
      <div v-for="(item, index) in iconList" :key="index" @click="selectedIcon(item)">
        <svg-icon :icon-class="item" style="height: 30px;width: 16px;" />
        <span>{{ item }}</span>
      </div>
    </div>
  </div>
</template>
<script>
export default {
  name: 'IconSelect',
}
</script>

<script setup>
import icons from './requireIcons'
import { ref } from 'vue'

import { getCurrentInstance } from 'vue'

const context = getCurrentInstance()

// 弹出层标题
const name = ref('')
// 数据
const iconList = ref(icons)

const filterIcons = () => {
  if (name.value) {
    iconList.value = iconList.value.filter(item => item.includes(name.value))
  } else {
    iconList.value = icons
  }
}

const selectedIcon = (name) => {
  context.emit('selected', name)
  document.body.click()
}

const reset = () => {
  name.value = ''
  iconList.value = icons
}

defineExpose({ reset })
</script>

<style rel="stylesheet/scss" lang="scss" scoped>
  .icon-body {
    width: 100%;
    padding: 10px;
    .icon-list {
      height: 200px;
      overflow-y: scroll;
      div {
        height: 30px;
        line-height: 30px;
        margin-bottom: -5px;
        cursor: pointer;
        width: 33%;
        float: left;
      }
      span {
        display: inline-block;
        vertical-align: -0.15em;
        fill: currentColor;
        overflow: hidden;
      }
    }
  }
</style>

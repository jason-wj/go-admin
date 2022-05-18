<template>
  <div>
    <el-row :gutter="10" class="system_state">
      <el-col sm="24" :md="8" :span="12">
        <el-card v-if="info.cpu" class="box-card" :body-style="bodyStyle">
          <template #header>
            <span>CPU使用率</span>
            <div class="monitor-header">
              <el-progress :color="userStore.activeColor" type="circle" :percentage="info.cpu.Percent" />
            </div>
          </template>
          <div class="monitor-footer">
            <Cell label="CPU主频" :value="info.cpu.cpuInfo[0].modelName.split('@ ')[1]" border />
            <Cell label="核心数" :value="`${info.cpu.cpuInfo[0].cores}`" />
            <el-row v-for="(item, index) in info.cpu.cpus" :key="index" :gutter="10">
              <el-col :span="12">core {{ index }}:</el-col>
              <el-col
                :span="12"
              ><el-progress
                type="line"
                :percentage="+item.toFixed(0)"
                :color="colors"
              /></el-col>
            </el-row>
          </div>
        </el-card>
      </el-col>
      <el-col :sm="24" :md="8" :span="12">
        <el-card v-if="info.mem" class="box-card" :body-style="bodyStyle">
          <template #header>
            <span>内存使用率</span>
            <div class="monitor-header">
              <el-progress :color="userStore.activeColor" type="circle" :percentage="info.mem.usage" />
            </div>
          </template>
          <div class="monitor-footer">
            <Cell label="总内存" :value="info.mem.total+'G'" border />
            <Cell label="已用内存" :value="info.mem.used+'G'" />
          </div>
        </el-card>
      </el-col>
      <el-col :sm="24" :md="8" :span="12">
        <el-card v-if="info.mem" class="box-card" :body-style="bodyStyle">
          <template #header>
            <span>磁盘信息</span>
            <div class="monitor-header">
              <el-progress :color="userStore.activeColor" type="circle" :percentage=" Number(( (info.disk.total-info.disk.free) / info.disk.total * 100).toFixed(2))" />
            </div>
          </template>
          <div class="monitor-footer">
            <Cell label="总磁盘" :value="info.disk.total+'G'" border />
            <Cell label="已用磁盘" :value="info.disk.total-info.disk.free+'G'" />
          </div>
        </el-card>
      </el-col>
    </el-row>

    <el-card v-if="info.os">
      <template #header>
        <span>go运行环境</span>
      </template>
      <div>
        <Cell label="GO 版本" :value="info.os.version" border />
        <Cell label="Goroutine" :value="`${info.os.numGoroutine}`" border />
        <Cell label="项目地址" :value="info.os.projectDir" />
      </div>
    </el-card>

    <el-card v-if="info.os">
      <template #header>
        <span>服务器信息</span>
      </template>
      <div>
        <Cell label="主机名称" :value="info.os.hostName" border />
        <Cell label="操作系统" :value="info.os.goOs" border />
        <Cell label="服务器IP" :value="info.os.ip" border />
        <Cell label="系统架构" :value="info.os.arch" border />
        <Cell label="CPU" :value="info.cpu.cpuInfo[0].modelName" border />
        <Cell label="当前时间" :value="info.os.time" />
      </div>
    </el-card>

    <el-card>
      <template #header>
        <span>磁盘状态</span>
      </template>
      <div class="el-table el-table--enable-row-hover el-table--medium">
        <table cellspacing="0" style="width: 100%;">
          <thead>
            <tr>
              <th class="is-leaf"><div class="cell">盘符路径</div></th>
              <th class="is-leaf"><div class="cell">文件系统</div></th>
              <th class="is-leaf"><div class="cell">总大小</div></th>

              <th class="is-leaf"><div class="cell">可用大小</div></th>
              <th class="is-leaf"><div class="cell">已用大小</div></th>
              <th class="is-leaf"><div class="cell">已用百分比</div></th>
            </tr>
          </thead>
          <tbody v-if="info.diskList">
            <tr v-for="(forList,index) in info.diskList" :key="index">
              <td><div class="cell">{{ forList.path }}</div></td>
              <td><div class="cell">{{ forList.fstype }}</div></td>
              <td><div class="cell">{{ forList.total }}M</div></td>
              <td><div class="cell">{{ forList.free }}M</div></td>
              <td><div class="cell">{{ forList.used }}M</div></td>
              <td><div class="cell" :class="{'text-danger': forList.usedPercent > 80}">{{ forList.usedPercent }}%</div></td>
            </tr>
          </tbody>
        </table>
      </div>
    </el-card>
  </div>
</template>

<script>
export default {
  name: 'Monitor',
}
</script>

<script setup>
import Cell from '@/components/Cell/index.vue'
import { getServer } from '@/api/sys/tools/monitor'
import { onUnmounted, ref } from 'vue'
import { useUserStore } from '@/pinia/modules/user'

const bodyStyle = ref(
  { height: '180px', 'overflow-y': 'scroll' }
)

const colors = ref([
  { color: '#5cb87a', percentage: 20 },
  { color: '#e6a23c', percentage: 40 },
  { color: '#f56c6c', percentage: 80 }
])

const userStore = useUserStore()

const info = ref({})
const timer = ref(null)

const getServerInfo = () => {
  getServer().then(ret => {
    if (ret.code === 200) {
      info.value = ret
    }
  })
}
getServerInfo()

onUnmounted(() => {
  clearInterval(timer.value)
  timer.value = null
})

timer.value = setInterval(() => {
  getServerInfo()
}, 1000 * 10)

</script>

<style lang="scss" scoped>
  .monitor-header{
  display: flex;
  justify-content: center;
  align-items: center;
  }
  .system_state {
  padding: 10px;
  }

  .el-card__body{
  padding: 20px 20px 0 20px!important;
  }

  .box-card {
  height: 400px;
  }
</style>

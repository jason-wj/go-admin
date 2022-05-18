import { createApp } from 'vue'
import 'element-plus/dist/index.css'
import './style/element_visiable.scss'
import ElementPlus from 'element-plus'
// import zhCn from 'element-plus/es/locale/lang/zh-cn'
// 引入gin-vue-admin前端初始化相关内容
import './core/go-admin-web'
// 引入封装的router
import router from '@/router/index'
import '@/permission'
import run from '@/core/go-admin-web'
import auth from '@/directive/auth'
import { store } from '@/pinia'
import App from './App.vue'
import i18n from './i18n' // added by mohamed hassan to multilangauge
import Pagination from '@/components/Pagination/index.vue'
import { getDicts } from '@/api/sys/dictdata'
import { getConfigKey } from '@/api/sys/config'

import SvgIcon from '@/components/SvgIcon/index.vue'// svg component
import 'virtual:svg-icons-register'

import { formatDate } from '@/utils/format'
import { addDateRange, selectDictLabel } from '@/utils/costum'

const app = createApp(App)
app.config.productionTip = false
// 全局组件挂载
app.component('Pagination', Pagination)
app.component('SvgIcon', SvgIcon)
app.config.globalProperties.addDateRange = addDateRange
app.config.globalProperties.formatDate = formatDate
app.config.globalProperties.getDicts = getDicts
app.config.globalProperties.getConfigKey = getConfigKey
app.config.globalProperties.selectDictLabel = selectDictLabel

app
  .use(run)
  .use(store)
  .use(auth)
  .use(router)
  .use(i18n)
  .use(ElementPlus, { i18n: (key, value) => i18n.t(key, value) })
  .mount('#app')

export default app

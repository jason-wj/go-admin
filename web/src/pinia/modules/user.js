import { login, getUserInfo } from '@/api/sys/user'
// import { jsonInBlacklist } from '@/api/jwt'
import router from '@/router/index'
import { ElMessage } from 'element-plus'
import { defineStore } from 'pinia'
import { ref, computed, watch } from 'vue'
import { useRouterStore } from './router'
import settings from '@/settings'

export const useUserStore = defineStore('user', () => {
  const userInfo = ref({})
  const mode = ref(settings.sideMode)
  const defaultRouter = ref(settings.defaultRouter)
  const token = ref(window.localStorage.getItem('token') || '')
  const language = ref(window.localStorage.getItem('langauge') || 'en') // added by mohamed hassan to allow store selected language for multilanguage support.

  const setUserInfo = (val) => {
    userInfo.value = val
  }

  const setToken = (val) => {
    token.value = val
  }

  // added by mohame hassan to allow store selected language for multilanguage support.
  const setLanguage = (val) => {
    console.log('setLanguage called with value: ' + val)
    language.value = val
  }

  const getLanguage = () => {
    return language.value
  }

  const NeedInit = () => {
    token.value = ''
    window.localStorage.removeItem('token')
    localStorage.clear()
    router.push({ name: 'Init', replace: true })
  }

  const ResetUserInfo = (value = {}) => {
    userInfo.value = {
      ...userInfo.value,
      ...value
    }
  }

  /* 获取用户信息*/
  const GetUserInfo = async() => {
    const res = await getUserInfo()
    if (res.code === 200) {
      setUserInfo(res.data)
    }
    return res
  }
  /* 登录*/
  const LoginIn = async(loginInfo) => {
    const res = await login(loginInfo)
    if (res.code === 200) {
      setUserInfo(res.data.userInfo)
      setToken(res.data.token)
      const routerStore = useRouterStore()
      await routerStore.SetAsyncRouter()
      const asyncRouters = routerStore.asyncRouters
      asyncRouters.forEach(asyncRouter => {
        router.addRoute(asyncRouter)
      })
      router.push({ name: defaultRouter.value })
      return true
    }
  }
  /* 登出*/
  const LoginOut = async() => {
    /* const res = await jsonInBlacklist()
            if (res.code === 0) {
              token.value = ''
              sessionStorage.clear()
              localStorage.clear()
              router.push({ name: 'Login', replace: true })
              window.location.reload()
            }*/

    token.value = ''
    sessionStorage.clear()
    localStorage.clear()
    router.push({ name: 'Login', replace: true })
    window.location.reload()
  }
  /* 设置侧边栏模式*/
  const changeSideMode = async(data) => {
    mode.value = data
    ElMessage({
      type: 'success',
      message: '设置成功'
    })
    /* const res = await setSelfInfo({ conf.sideMode: data })
            if (res.code === 0) {
              userInfo.value.sideMode = data
              ElMessage({
                type: 'success',
                message: '设置成功'
              })
            }*/
  }
  const sideMode = computed(() => {
    if (mode.value === 'dark') {
      return settings.backgroundColor
    } else if (mode.value === 'light') {
      return settings.baseColor
    } else {
      return settings.sideMode
    }
  })
  const baseColor = computed(() => {
    if (mode.value === 'dark') {
      return settings.baseColor
    } else if (mode.value === 'light') {
      return settings.backgroundColor
    } else {
      return settings.baseColor
    }
  })

  const backgroundColor = computed(() => {
    if (mode.value === 'dark') {
      return settings.backgroundColor
    } else if (mode.value === 'light') {
      return settings.baseColor
    } else {
      return settings.baseColor
    }
  })
  const activeColor = computed(() => {
    if (mode.value === 'dark' || mode.value === 'light') {
      return settings.activeColor
    }
    return settings.activeColor
  })

  watch(token, () => {
    window.localStorage.setItem('token', token.value)
  })

  return {
    userInfo,
    token,
    language, // added by mohame hassan to allow store selected language for multilanguage support.
    NeedInit,
    ResetUserInfo,
    GetUserInfo,
    LoginIn,
    LoginOut,
    setLanguage, // added by mohame hassan to allow store selected language for multilanguage support.
    getLanguage, // added by mohame hassan to allow store selected language for multilanguage support.
    changeSideMode,
    mode,
    defaultRouter,
    sideMode,
    setToken,
    baseColor,
    backgroundColor,
    activeColor
  }
})

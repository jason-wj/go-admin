import { createRouter, createWebHashHistory } from 'vue-router'

const routes = [{
  path: '/',
  redirect: '/login'
},
{
  path: '/login',
  name: 'Login',
  component: () => import('@/view/sys/login/index.vue')
}
]

const router = createRouter({
  history: createWebHashHistory(),
  routes
})

export default router

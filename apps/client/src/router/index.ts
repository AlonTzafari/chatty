import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '../views/HomeView.vue'
import { useAuthStore } from '@/stores/auth'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      
      component: HomeView,
    },
    {
      path: '/about',
      name: 'about',
      // route level code-splitting
      // this generates a separate chunk (About.[hash].js) for this route
      // which is lazy-loaded when the route is visited.
      component: () => import('../views/AboutView.vue')
    },
    {
      path: '/login',
      name: 'login',
      // route level code-splitting
      // this generates a separate chunk (About.[hash].js) for this route
      // which is lazy-loaded when the route is visited.
      component: () => import('../views/login-view.vue')
    },
    {
      path: '/create-channel',
      name: 'createChannel',
      // route level code-splitting
      // this generates a separate chunk (About.[hash].js) for this route
      // which is lazy-loaded when the route is visited.
      component: () => import('../views/create-channel-view.vue')
    },
  ]
})
router.beforeEach(async (to, from, next) => {
    if(to.path === '/login') {
        return next() 
    }
    const authStore = useAuthStore()
    const user = await authStore.fetchUser()
    if(user == null) {
        next({path: '/login', query: {from: to.fullPath}})
    } 
    return next()
})
export default router

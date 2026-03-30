import { createRouter, createWebHashHistory } from 'vue-router'
import { useAuth } from '../composables/useAuth'

import LoginView from '../views/LoginView.vue'
import DashboardView from '../views/DashboardView.vue'
import UrlsListView from '../views/UrlsListView.vue'
import CreateUrlView from '../views/CreateUrlView.vue'
import UrlDetailView from '../views/UrlDetailView.vue'
import AccountView from '../views/AccountView.vue'

const routes = [
  {
    path: '/login',
    name: 'login',
    component: LoginView,
    meta: { public: true },
  },
  {
    path: '/',
    name: 'dashboard',
    component: DashboardView,
  },
  {
    path: '/urls',
    name: 'urls',
    component: UrlsListView,
  },
  {
    path: '/urls/new',
    name: 'create-url',
    component: CreateUrlView,
  },
  {
    path: '/urls/:shortcode',
    name: 'url-detail',
    component: UrlDetailView,
    props: true,
  },
  {
    path: '/account',
    name: 'account',
    component: AccountView,
  },
]

const router = createRouter({
  history: createWebHashHistory(),
  routes,
})

router.beforeEach((to) => {
  const { isAuthenticated } = useAuth()
  if (!to.meta.public && !isAuthenticated.value) {
    return { name: 'login' }
  }
  if (to.name === 'login' && isAuthenticated.value) {
    return { name: 'dashboard' }
  }
})

export default router

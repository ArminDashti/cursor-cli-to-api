import { createRouter, createWebHistory } from 'vue-router'
import { loadAuth } from '@/lib/utils'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/login',
      name: 'login',
      component: () => import('@/views/auth/LoginView.vue'),
      meta: { public: true },
    },
    {
      path: '/',
      redirect: '/playground/chat',
    },
    {
      path: '/playground',
      redirect: '/playground/chat',
    },
    {
      path: '/playground/chat',
      name: 'playground-chat',
      component: () => import('@/views/playground/ChatCompletionsView.vue'),
    },
    {
      path: '/playground/responses',
      name: 'playground-responses',
      component: () => import('@/views/playground/ResponsesView.vue'),
    },
    {
      path: '/docs',
      redirect: '/docs/overview',
    },
    {
      path: '/docs/overview',
      name: 'docs-overview',
      component: () => import('@/views/docs/OverviewView.vue'),
    },
    {
      path: '/docs/chat-completions',
      name: 'docs-chat',
      component: () => import('@/views/docs/ChatCompletionsDocView.vue'),
    },
    {
      path: '/docs/responses',
      name: 'docs-responses',
      component: () => import('@/views/docs/ResponsesDocView.vue'),
    },
  ],
})

router.beforeEach((to) => {
  if (to.meta.public) return true
  if (!loadAuth()) return { name: 'login', query: { redirect: to.fullPath } }
  return true
})

export default router

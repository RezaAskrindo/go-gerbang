import { createRouter, createWebHistory } from 'vue-router'
import MainView from '@/views/MainView.vue'
import HomeView from '@/views/HomeView.vue'
import AuthView from '@/views/AuthView.vue'
import NotFoundView from '@/views/NotFoundView.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/auth',
      component: AuthView,
      // component: () => import('@/views/AuthView.vue'),
      children: [
        {
          path: '',
          redirect: { name: 'Login' }
        },
        {
          path: 'login',
          name: 'Login',
          component: () => import('@/components/auth/LoginAuth.vue'),
        },
        {
          path: 'sign-up',
          name: 'Sign Up',
          component: () => import('@/components/auth/SignUpUser.vue'),
        },
        {
          path: 'reset-password',
          name: 'Reset Password',
          component: () => import('@/components/auth/ResetPassword.vue'),
        },
        {
          path: 'forget-password',
          name: 'Forget Password',
          component: () => import('@/components/auth/ForgetPassword.vue'),
        },
      ]
    },
    {
      path: '/',
      component: MainView,
      children: [
        {
          path: '/',
          name: 'home',
          component: HomeView
        }
      ]
    },
    {
      path: '/:pathMatch(.*)*',
      name: 'NotFound',
      component: NotFoundView,
    },
  ],
})

export default router

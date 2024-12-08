/**
 * router/index.ts
 *
 * Automatic routes for `./src/pages/*.vue`
 */

// Composables
import { createRouter, createWebHistory } from 'vue-router/auto'
import { setupLayouts } from 'virtual:generated-layouts'
import { routes } from 'vue-router/auto-routes'
import { useAuthStore } from '@/stores/auth'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: setupLayouts(routes),
})

// Workaround for https://github.com/vitejs/vite/issues/11804
router.onError((err, to) => {
  if (err?.message?.includes?.('Failed to fetch dynamically imported module')) {
    if (!localStorage.getItem('vuetify:dynamic-reload')) {
      console.log('Reloading page to fix dynamic import error')
      localStorage.setItem('vuetify:dynamic-reload', 'true')
      location.assign(to.fullPath)
    } else {
      console.error('Dynamic import error, reloading page did not fix it', err)
    }
  } else {
    console.error(err)
  }
})

router.isReady().then(() => {
  localStorage.removeItem('vuetify:dynamic-reload')
})

router.beforeEach(async (to, from, next) => {
  const authStore = useAuthStore()

  if (to.meta.layout === 'admin') {
    await authStore.checkAuth()

    if (!authStore.isAuth) {
      next({
        path: '/login',
        query: { redirect: to.fullPath },
      })
      return
    }

    let isAdmin = false
    try {
      const me = await authStore.me()
      if (me.role) isAdmin = me.role.admin
    } catch {
      isAdmin = false
    }

    if (!isAdmin) {
      next({ path: '/' })
      return
    }
  } else if (to.meta.layout === 'default') {
    await authStore.checkAuth()

    if (!authStore.isAuth) {
      next({
        path: '/login',
        query: { redirect: to.fullPath },
      })
      return
    }
  }

  next()
})

export default router

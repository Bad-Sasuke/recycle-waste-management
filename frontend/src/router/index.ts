import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '../views/HomeView.vue'
import { useUsersStore } from '@/stores/users'
import { useShopStore } from '@/stores/shop'

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
      component: () => import('../views/AboutView.vue'),
    },
    {
      path: '/marketplace',
      name: 'marketplace',
      component: () => import('../views/MarketplaceView.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/auth',
      name: 'auth',
      component: () => import('../views/AuthView.vue'),
    },
    {
      path: '/profile',
      name: 'profile',
      component: () => import('../views/ProfileView.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/shop/create',
      name: 'create-shop',
      component: () => import('../views/CreateShopView.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/manage-shop',
      name: 'manage-shop',
      component: () => import('../views/ManageShopView.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/settings',
      name: 'settings',
      component: () => import('../views/SettingsView.vue'),
      meta: { requiresAuth: true },
    },
  ],
})

// Global navigation guard to handle authentication and shop checking
router.beforeEach(async (to, from, next) => {
  // Check if this is the create shop route - allow access without shop check
  if (to.name === 'create-shop') {
    next()
    return
  }

  // Always check user login state
  const usersStore = useUsersStore()
  await usersStore.checkLogin() // Refresh authentication state

  // Check if user is logged in
  if (usersStore.isLogin) {
    // For all authenticated users, check if they have a shop
    const shopStore = useShopStore()

    // Only check user shop if we haven't checked yet
    if (!shopStore.checked) {
      await shopStore.checkUserShop()
    }

    // Check if the user doesn't have a shop and redirect to create shop page
    // This ensures users who don't have shops are always redirected to create one
    if (!shopStore.hasShop) {
      // For all authenticated routes, redirect to create shop if no shop exists
      // Store the intended destination so we can redirect back after creating a shop
      next({
        name: 'create-shop',
        query: { redirect: to.fullPath },
      })
      return
    }
  }

  if (to.matched.some((record) => record.meta.requiresAuth)) {
    // This route requires auth, check if logged in
    // If not, redirect to login page
    if (!usersStore.isLogin) {
      next({
        path: '/auth',
        query: { redirect: to.fullPath },
      })
    } else {
      next()
    }
  } else {
    next() // Make sure to always call next()!
  }
})

export default router

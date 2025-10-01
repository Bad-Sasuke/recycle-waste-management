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
      meta: { requiresAuth: true }
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
      meta: { requiresAuth: true }
    },
    {
      path: '/shop/create',
      name: 'create-shop',
      component: () => import('../views/CreateShopView.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/shop/manage',
      name: 'manage-shop',
      component: () => import('../views/ManageShopView.vue'),
      meta: { requiresAuth: true }
    }
  ],
})

// Global navigation guard to handle authentication and shop checking
router.beforeEach(async (to, from, next) => {
  if (to.matched.some(record => record.meta.requiresAuth)) {
    // This route requires auth, check if logged in
    // If not, redirect to login page
    const usersStore = useUsersStore();
    await usersStore.checkLogin(); // Refresh authentication state
    
    if (!usersStore.isLogin) {
      next({
        path: '/auth',
        query: { redirect: to.fullPath }
      })
    } else {
      // Check if this is the create shop route - allow access without shop check
      if (to.name === 'create-shop') {
        next();
        return;
      }
      
      // For other authenticated routes, check if user has a shop
      const shopStore = useShopStore();
      if (!shopStore.checked) {
        await shopStore.checkUserShop();
      }
      
      // If user doesn't have a shop and is trying to access shop management or other protected areas,
      // redirect to create shop page (except for the create shop page itself)
      if (!shopStore.hasShop && to.name !== 'create-shop') {
        next({ name: 'create-shop' });
      } else {
        next();
      }
    }
  } else {
    next() // Make sure to always call next()!
  }
})

export default router

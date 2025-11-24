import { defineStore } from 'pinia'
import { getCookie, deleteCookie } from '../stores/cookie'
import { useShopStore } from './shop'

interface User {
  user_id: string
  username: string
  email: string
  image_url: string
  created_at: string
  last_login: string
  role?: string
}

// Interface for waste data to edit
export interface WasteToEdit {
  id: string
  name: string
  price: number
  category: string
  url: string
  last_update: string
}

export const useUsersStore = defineStore('users', {
  state: () => ({
    isLogin: false,
    jwt: '',
    user: null as User | null,
  }),
  getters: {
    profileImage: (state) => state.user?.image_url || '',
    username: (state) => state.user?.username || '',
  },
  actions: {
    async checkLogin() {
      const token = getCookie('token')
      this.isLogin = !!token
      this.jwt = token || ''

      // If logged in, fetch user profile
      if (this.isLogin) {
        await this.fetchUserProfile()
      } else {
        // If not logged in, reset the shop store as well
        const shopStore = useShopStore()
        shopStore.resetShop()
      }
    },

    async fetchUserProfile() {
      if (!this.isLogin || !this.jwt) return

      try {
        const apiUrl = import.meta.env.VITE_WEB_API
        if (!apiUrl) {
          console.error('API configuration error. Please check environment settings.')
          return
        }

        const response = await fetch(`${apiUrl}/api/user/profile`, {
          headers: {
            Authorization: `Bearer ${this.jwt}`,
            'Content-Type': 'application/json',
          },
        })

        if (response.ok) {
          const result = await response.json()
          this.user = result.data

          // After fetching user profile, check if user has a shop
          // This will be handled by the router guard, but we can still trigger the shop check here
          const shopStore = useShopStore()
          await shopStore.checkUserShop()
        } else {
          console.error(`Failed to fetch profile data: ${response.status} ${response.statusText}`)
          // If profile fetch fails, set user to null to avoid stale data
          this.user = null
        }
      } catch (error) {
        console.error('Error fetching user profile:', error)
        this.user = null
      }
    },

    logout() {
      deleteCookie('token')
      this.isLogin = false
      this.jwt = ''
      this.user = null
    },

    async updateUserRole(role: string) {
      if (!this.isLogin || !this.jwt) return false

      try {
        const apiUrl = import.meta.env.VITE_WEB_API
        const response = await fetch(`${apiUrl}/api/user/update-role-user`, {
          method: 'PUT',
          headers: {
            Authorization: `Bearer ${this.jwt}`,
            'Content-Type': 'application/json',
          },
          body: JSON.stringify({ role }),
        })

        if (response.ok) {
          await this.fetchUserProfile()
          return true
        }
        return false
      } catch (error) {
        console.error('Error updating user role:', error)
        return false
      }
    },
  },
})

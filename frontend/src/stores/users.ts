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
  shop_id?: string
}

interface Employee {
  employee_id: string
  shop_id: string
  first_name: string
  last_name: string
  username: string
  created_at?: string
  updated_at?: string
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
    employee: null as Employee | null,
    isEmployee: false,
  }),
  getters: {
    profileImage: (state) => state.user?.image_url || '',
    username: (state) => {
      if (state.isEmployee && state.employee) {
        return state.employee.username
      }
      return state.user?.username || ''
    },
    displayName: (state) => {
      if (state.isEmployee && state.employee) {
        return `${state.employee.first_name} ${state.employee.last_name}`
      }
      return state.user?.username || 'ไม่มีชื่อ'
    },
  },
  actions: {
    async checkLogin() {
      const token = getCookie('token')
      this.isLogin = !!token
      this.jwt = token || ''

      // Check if this is an employee or user
      if (this.isLogin && token) {
        // Check if employee data is stored in localStorage
        const employeeData = localStorage.getItem('employee_data')
        if (employeeData) {
          try {
            this.employee = JSON.parse(employeeData)
            this.isEmployee = true
            this.user = null // Employee doesn't have user profile
          } catch (error) {
            console.error('Error parsing employee data:', error)
          }
        } else {
          // Regular user
          this.isEmployee = false
          await this.fetchUserProfile()
        }
      } else {
        // If not logged in, reset the shop store as well
        const shopStore = useShopStore()
        shopStore.resetShop()
      }
    },

    async fetchUserProfile() {
      if (!this.isLogin || !this.jwt || this.isEmployee) return

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
        } else if (response.status === 403) {
          // 403 likely means this is an employee token, not a user token
          // Check if employee data exists in localStorage
          const employeeData = localStorage.getItem('employee_data')
          if (employeeData) {
            try {
              this.employee = JSON.parse(employeeData)
              this.isEmployee = true
              this.user = null
            } catch (error) {
              console.error('Error parsing employee data:', error)
            }
          }
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

    setEmployeeData(employee: Employee) {
      this.employee = employee
      this.isEmployee = true
      this.user = null
      // Store employee data in localStorage
      localStorage.setItem('employee_data', JSON.stringify(employee))
    },

    logout() {
      deleteCookie('token')
      localStorage.removeItem('employee_data')
      this.isLogin = false
      this.jwt = ''
      this.user = null
      this.employee = null
      this.isEmployee = false
    },

    async updateUserRole(role: string) {
      if (!this.isLogin || !this.jwt || this.isEmployee) return false

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

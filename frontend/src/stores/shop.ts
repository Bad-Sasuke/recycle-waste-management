import { defineStore } from 'pinia'
import type { Shop } from '@/types/shop'
import { fetchMyShop, createShop, updateShop, deleteShop, fetchShops } from '@/services/shop'
import { getCookie } from '@/stores/cookie'

export const useShopStore = defineStore('shop', {
  state: () => ({
    shop: null as Shop | null,
    allShops: [] as Shop[],
    isLoading: false,
    hasShop: false,
    checked: false, // To track whether we've checked for shop existence
    pagination: {
      page: 1,
      limit: 12,
      total_pages: 1,
      total_items: 0,
    },
  }),

  getters: {
    shopName: (state) => state.shop?.name || '',
    shopImage: (state) => state.shop?.image_url || '',
    shopAddress: (state) => state.shop?.address || '',
  },

  actions: {
    async fetchAllShops(page: number = 1, limit: number = 100) {
      this.isLoading = true
      try {
        const result = await fetchShops(page, limit)
        if (result.data) {
          this.allShops = result.data as Shop[]
          this.pagination = {
            page: result.page || 1,
            limit: result.limit || limit,
            total_pages: result.total_pages || 1,
            total_items: result.total_items || 0,
          }
        } else {
          this.allShops = []
        }
      } catch (error) {
        console.error('Error fetching all shops:', error)
        this.allShops = []
      } finally {
        this.isLoading = false
      }
    },

    async checkUserShop() {
      if (!getCookie('token')) {
        this.hasShop = false
        this.checked = true
        return false
      }

      this.isLoading = true
      try {
        const result = await fetchMyShop()
        if (result.data) {
          this.shop = result.data as Shop
          this.hasShop = true
        } else {
          this.shop = null
          this.hasShop = false
        }
      } catch (error) {
        console.error('Error checking user shop:', error)
        this.hasShop = false
        this.shop = null
      } finally {
        this.isLoading = false
        this.checked = true
      }
      return this.hasShop
    },

    async createShop(formData: FormData) {
      this.isLoading = true
      try {
        const result = await createShop(formData)
        if (result.message.includes('successfully')) {
          // Refresh the shop data
          await this.checkUserShop()
          return { success: true, message: result.message }
        } else {
          return { success: false, message: result.message }
        }
      } catch (error) {
        console.error('Error creating shop:', error)
        return { success: false, message: 'Error creating shop' }
      } finally {
        this.isLoading = false
      }
    },

    async updateShop(shopId: string, formData: FormData) {
      this.isLoading = true
      try {
        const result = await updateShop(shopId, formData)
        if (result.message.includes('successfully')) {
          // Refresh the shop data
          await this.checkUserShop()
          return { success: true, message: result.message }
        } else {
          return { success: false, message: result.message }
        }
      } catch (error) {
        console.error('Error updating shop:', error)
        return { success: false, message: 'Error updating shop' }
      } finally {
        this.isLoading = false
      }
    },

    async deleteShop(shopId: string) {
      this.isLoading = true
      try {
        const result = await deleteShop(shopId)
        if (result.message.includes('successfully')) {
          this.shop = null
          this.hasShop = false
          return { success: true, message: result.message }
        } else {
          return { success: false, message: result.message }
        }
      } catch (error) {
        console.error('Error deleting shop:', error)
        return { success: false, message: 'Error deleting shop' }
      } finally {
        this.isLoading = false
      }
    },

    // Reset the shop state
    resetShop() {
      this.shop = null
      this.hasShop = false
      this.checked = false
    },
  },
})

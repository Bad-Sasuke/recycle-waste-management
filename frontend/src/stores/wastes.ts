import { defineStore } from 'pinia'
import config from '../config'
import { useUsersStore } from './users'
import { useShopStore } from './shop'
import type RecycleWaste from '../types/recycle_waste'
import type { GroupedRecyclableItem } from '../types/recycle_waste'

const webAPI = config.webAPI

// Interface for waste data to edit
export interface WasteToEdit {
  id: string
  name: string
  price: number
  category: string
  url: string
  last_update: string
}

export const useWastesStore = defineStore('wastes', {
  state: () => ({
    wastes: [] as RecycleWaste[], // ใช้ interface ที่นำเข้า
    groupedWastes: [] as GroupedRecyclableItem[], // สำหรับข้อมูลที่จัดกลุ่มตามชื่อสินค้า
    wasteToEdit: null as WasteToEdit | null,
    pagination: {
      page: 1,
      limit: 12,
      total_pages: 1,
      total_items: 0,
    },
    groupedPagination: {
      page: 1,
      limit: 12,
      total_pages: 1,
      total_items: 0,
    },
  }),
  actions: {
    async fetchWastes(page: number = 1, limit: number = 12) {
      try {
        const response = await fetch(
          webAPI + `/api/recycle-waste/get-wastes?page=${page}&limit=${limit}`,
        )
        const dataJson = await response.json()

        if (!response.ok) {
          throw new Error(`HTTP error! Status: ${response.status}`)
        }

        // Handle paginated response
        if (dataJson.data && Array.isArray(dataJson.data)) {
          // Check if the data is already grouped (contains shops array) or individual items
          if (dataJson.data.length > 0 && dataJson.data[0].shops !== undefined) {
            // Data is grouped
            this.groupedWastes = dataJson.data as GroupedRecyclableItem[]

            // Update grouped pagination info
            if (dataJson.page !== undefined) this.groupedPagination.page = dataJson.page
            if (dataJson.limit !== undefined) this.groupedPagination.limit = dataJson.limit
            if (dataJson.total_pages !== undefined)
              this.groupedPagination.total_pages = dataJson.total_pages
            if (dataJson.total_items !== undefined)
              this.groupedPagination.total_items = dataJson.total_items
          } else {
            // Data is individual items
            const wastes = dataJson.data as RecycleWaste[]

            wastes.forEach((item: RecycleWaste) => {
              // Check if last_update is a valid string or date, and only convert if it is
              if (item.last_update) {
                item.last_update = new Date(item.last_update).toLocaleDateString('th-TH', {
                  day: 'numeric',
                  month: 'numeric',
                  year: 'numeric',
                })
              } else {
                item.last_update = '' // Set to an empty string or a default value
              }
            })

            this.wastes = wastes

            // Update pagination info if available in response
            if (dataJson.page !== undefined) this.pagination.page = dataJson.page
            if (dataJson.limit !== undefined) this.pagination.limit = dataJson.limit
            if (dataJson.total_pages !== undefined)
              this.pagination.total_pages = dataJson.total_pages
            if (dataJson.total_items !== undefined)
              this.pagination.total_items = dataJson.total_items
          }
        } else {
          console.error('Data is not an array or is undefined')
          this.wastes = [] // Optionally set to an empty array
          this.groupedWastes = [] // Also reset grouped data
        }
      } catch (error) {
        console.error('Error fetching waste data:', error)
      }
    },
    async addWaste(waste: unknown) {
      try {
        const formData = new FormData()
        formData.append('name', (waste as { name: string }).name)
        formData.append('price', (waste as { price: number }).price.toString())
        formData.append('category', (waste as { category: string }).category)
        if (waste && typeof waste === 'object' && 'imageFile' in waste) {
          // ตรวจสอบว่า imageFile มีอยู่ใน waste
          formData.append('image_file', (waste as { imageFile: File }).imageFile)
        }

        const usersStore = useUsersStore()
        const shopStore = useShopStore()

        // Add shop_id to the form data
        if (shopStore.shop && shopStore.shop.shop_id) {
          formData.append('shop_id', shopStore.shop.shop_id)
        }

        const response = await fetch(webAPI + '/api/recycle-waste/add-waste', {
          method: 'POST',
          headers: {
            Authorization: `Bearer ${usersStore.jwt}`,
          },
          body: formData, // ใช้ formData เป็น body
        })
        if (response.status !== 200) {
          return response.status
        }
        await this.fetchWastes()
        return response.status
      } catch (error) {
        console.error('err:', error)
      }
    },
    async deleteWaste(id: string) {
      try {
        const usersStore = useUsersStore()
        const shopStore = useShopStore()

        let deleteUrl = webAPI + '/api/recycle-waste/delete-waste/' + id

        // Add shop_id as query parameter if available
        if (shopStore.shop && shopStore.shop.shop_id) {
          deleteUrl += `?shop_id=${shopStore.shop.shop_id}`
        }

        const response = await fetch(deleteUrl, {
          method: 'DELETE',
          headers: {
            Authorization: `Bearer ${usersStore.jwt}`,
          },
        })
        if (response.status !== 200) {
          return response.status
        }
        await this.fetchWastes()
        return response.status
      } catch (error) {
        console.error('err:', error)
      }
    },
    async updateWaste(id: string, waste: unknown) {
      try {
        const formData = new FormData()
        formData.append('name', (waste as { name: string }).name)
        formData.append('price', (waste as { price: number }).price.toString())
        formData.append('category', (waste as { category: string }).category)
        if (
          waste &&
          typeof waste === 'object' &&
          'imageFile' in waste &&
          (waste as { imageFile: File }).imageFile
        ) {
          // ตรวจสอบว่า imageFile มีอยู่ใน waste
          formData.append('image_file', (waste as { imageFile: File }).imageFile)
        }

        const usersStore = useUsersStore()
        const shopStore = useShopStore()

        // Add shop_id to the form data
        if (shopStore.shop && shopStore.shop.shop_id) {
          formData.append('shop_id', shopStore.shop.shop_id)
        }

        const response = await fetch(webAPI + '/api/recycle-waste/edit-waste/' + id, {
          method: 'PUT', // ใช้ PUT method สำหรับการอัปเดต
          headers: {
            Authorization: `Bearer ${usersStore.jwt}`,
          },
          body: formData,
        })
        if (response.status !== 200) {
          console.error('Error updating waste:', response.status)
          return response.status
        }
        await this.fetchWastes() // Refresh the waste list
        return response.status
      } catch (error) {
        console.error('err:', error)
      }
    },

    // Set the waste to be edited
    setWasteToEdit(
      id: string,
      wasteData: {
        name: string
        price: number
        category: string
        url: string
        last_update: string
      },
    ) {
      this.wasteToEdit = {
        id,
        name: wasteData.name,
        price: wasteData.price,
        category: wasteData.category,
        url: wasteData.url,
        last_update: wasteData.last_update,
      }
    },

    // Clear the waste to be edited
    clearWasteToEdit() {
      this.wasteToEdit = null
    },
  },
})

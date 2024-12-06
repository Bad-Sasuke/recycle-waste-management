import { defineStore } from 'pinia'
import config from '../config'
import { fetchData } from '../services/api'
import type CategoryWaste from '../types/category_waste'

const webAPI = config.webAPI

export const useCategoryWasteStore = defineStore('categoryWaste', {
    state: () => ({
        category: [] as CategoryWaste[], // ใช้ interface ที่นำเข้า
    }),
    actions: {
        async fetchCategoryWaste() {
            try {
                const dataJson = await fetchData(webAPI + '/api/category-waste/get-category')
                if (!dataJson.data) {
                    throw new Error('error')
                }
                this.category = dataJson.data as { id?: string; name?: string }[];
            } catch (error) {
                console.error('err:', error)
            }
        },
    },
})

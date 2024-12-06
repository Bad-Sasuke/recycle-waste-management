import { defineStore } from 'pinia'
import config from '../config'
import { fetchData } from '../services/api'
import * as entities from '../types/category_waste'

const webAPI = config.webAPI

export const useCategoryWasteStore = defineStore('categoryWaste', {
    state: () => ({
        category: [] as entities.CategoryWaste[], // ใช้ interface ที่นำเข้า
    }),
    actions: {
        async fetchCategoryWaste() {
            try {
                const dataJson = await fetchData(webAPI + '/api/category-waste/get-category')
                if (!dataJson.data) {
                    throw new Error('error')
                }
                this.category = dataJson.data;
            } catch (error) {
                console.error('err:', error)
            }
        },
    },
})

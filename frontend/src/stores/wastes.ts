import { defineStore } from 'pinia'
import config from '../config'
import { fetchData, postData } from '../services/api'
import * as entities from '../types/recycle_waste'

const webAPI = config.webAPI

export const useWastesStore = defineStore('wastes', {
  state: () => ({
    wastes: [] as entities.RecycleWaste[], // ใช้ interface ที่นำเข้า
  }),
  actions: {
    async fetchWastes() {
      try {
        const dataJson = await fetchData(webAPI + '/api/recycle-waste/get-wastes')
        dataJson.data.forEach((item: entities.RecycleWaste) => {
          item.lastUpdate = new Date(item.lastUpdate).toLocaleDateString('th-TH', {
            day: 'numeric',
            month: 'numeric',
            year: 'numeric',
          })
        })
        if (!dataJson.data) {
          throw new Error('error')
        }
        this.wastes = dataJson.data
      } catch (error) {
        console.error('err:', error)
      }
    },
    async addWaste(waste: any) {
      try {
        const formData = new FormData();
        formData.append('name', waste.name);
        formData.append('price', waste.price);
        formData.append('category', waste.category);
        if (waste.imageFile) {
          formData.append('image_file', waste.imageFile);  // 'image' คือตัวแปรใน API ของคุณ
        }
        const response = await fetch(webAPI + '/api/recycle-waste/add-waste', {
          method: 'POST',
          body: formData,  // ใช้ formData เป็น body
        });
        if (!response.status == 200) {
          return response.status
        }
        await this.fetchWastes()
        return response.status
      } catch (error) {
        console.error('err:', error)
      }
    }
  },
})

import { defineStore } from 'pinia'
import config from '../config'
import { fetchData } from '../services/api'
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
  },
})

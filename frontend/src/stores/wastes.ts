import { defineStore } from 'pinia'
import config from '../config'
import { fetchData } from '../services/api'
import type RecycleWaste from '../types/recycle_waste'

const webAPI = config.webAPI



export const useWastesStore = defineStore('wastes', {
  state: () => ({
    wastes: [] as RecycleWaste[], // ใช้ interface ที่นำเข้า
  }),
  actions: {
    async fetchWastes() {
      try {
        const dataJson = await fetchData(webAPI + '/api/recycle-waste/get-wastes');

        // Ensure dataJson.data is properly typed as RecycleWaste[]
        if (dataJson.data && Array.isArray(dataJson.data)) {
          const wastes = dataJson.data as RecycleWaste[];

          wastes.forEach((item: RecycleWaste) => {
            // Check if lastUpdate is a valid string or date, and only convert if it is
            if (item.lastUpdate) {
              item.lastUpdate = new Date(item.lastUpdate).toLocaleDateString('th-TH', {
                day: 'numeric',
                month: 'numeric',
                year: 'numeric',
              });
            } else {
              item.lastUpdate = ''; // Set to an empty string or a default value
            }
          });

          this.wastes = wastes;
        } else {
          console.error('Data is not an array or is undefined');
          this.wastes = []; // Optionally set to an empty array
        }
      } catch (error) {
        console.error('Error fetching waste data:', error);
      }
    },
    async addWaste(waste: unknown) {
      try {
        const formData = new FormData();
        formData.append('name', (waste as { name: string }).name);
        formData.append('price', (waste as { price: number }).price.toString());
        formData.append('category', (waste as { category: string }).category);
        if (waste && typeof waste === 'object' && 'imageFile' in waste) {
          // ตรวจสอบว่า imageFile มีอยู่ใน waste
          formData.append('image_file', (waste as { imageFile: File }).imageFile);
        }
        const response = await fetch(webAPI + '/api/recycle-waste/add-waste', {
          method: 'POST',
          body: formData,  // ใช้ formData เป็น body
        });
        if (response.status !== 200) {
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

import { defineStore } from 'pinia'
import config from '../config'
import { fetchData } from '../services/api'
import type RecycleWaste from '../types/recycle_waste'

const webAPI = config.webAPI



// Interface for waste data to edit
export interface WasteToEdit {
  id: string;
  name: string;
  price: number;
  category: string;
  url: string;
  lastUpdate: string;
}

export const useWastesStore = defineStore('wastes', {
  state: () => ({
    wastes: [] as RecycleWaste[], // ใช้ interface ที่นำเข้า
    wasteToEdit: null as WasteToEdit | null,
    pagination: {
      page: 1,
      limit: 12,
      total_pages: 1,
      total_items: 0
    }
  }),
  actions: {
    async fetchWastes(page: number = 1, limit: number = 12) {
      try {
        const dataJson = await fetchData(webAPI + `/api/recycle-waste/get-wastes?page=${page}&limit=${limit}`);

        // Handle paginated response
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
          
          // Update pagination info if available in response
          if (dataJson.page !== undefined) this.pagination.page = dataJson.page;
          if (dataJson.limit !== undefined) this.pagination.limit = dataJson.limit;
          if (dataJson.total_pages !== undefined) this.pagination.total_pages = dataJson.total_pages;
          if (dataJson.total_items !== undefined) this.pagination.total_items = dataJson.total_items;
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
    },
    async deleteWaste(id: string) {
      try {
        const response = await fetch(webAPI + '/api/recycle-waste/delete-waste/' + id, {
          method: 'DELETE',
        });
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
        const formData = new FormData();
        formData.append('name', (waste as { name: string }).name);
        formData.append('price', (waste as { price: number }).price.toString());
        formData.append('category', (waste as { category: string }).category);
        if (waste && typeof waste === 'object' && 'imageFile' in waste && (waste as { imageFile: File }).imageFile) {
          // ตรวจสอบว่า imageFile มีอยู่ใน waste
          formData.append('image_file', (waste as { imageFile: File }).imageFile);
        }
        const response = await fetch(webAPI + '/api/recycle-waste/edit-waste/' + id, {
          method: 'PUT', // ใช้ PUT method สำหรับการอัปเดต
          body: formData,
        });
        if (response.status !== 200) {
          console.error('Error updating waste:', response.status);
          return response.status;
        }
        await this.fetchWastes(); // Refresh the waste list
        return response.status;
      } catch (error) {
        console.error('err:', error);
      }
    },
    
    // Set the waste to be edited
    setWasteToEdit(id: string, wasteData: { name: string; price: number; category: string; url: string; lastUpdate: string }) {
      this.wasteToEdit = {
        id,
        name: wasteData.name,
        price: wasteData.price,
        category: wasteData.category,
        url: wasteData.url,
        lastUpdate: wasteData.lastUpdate
      };
    },
    
    // Clear the waste to be edited
    clearWasteToEdit() {
      this.wasteToEdit = null;
    }
  },
})

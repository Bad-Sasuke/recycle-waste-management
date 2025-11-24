import {
  addCustomerRequest,
  fetchCustomerMyRequests,
  fetchCustomerRequestAll,
  acceptCustomerRequest,
  cancelCustomerRequest,
  completeCustomerRequest,
} from '@/services/customer_request'
import type { CustomerRequestBody, CustomerRequest } from '@/types/customer_request'
import { defineStore } from 'pinia'

export const useCustomerRequestStore = defineStore('customer_request', {
  state: () => ({
    customerRequestsMyRequests: [] as CustomerRequest[],
    customerRequestsAll: [] as CustomerRequest[],
    errorCustomerRequest: '' as string,
    // Pagination metadata
    total: 0,
    page: 1,
    limit: 10,
    totalPages: 0,
    // Distance filter
    maxDistance: 20, // Default 20 km
  }),
  actions: {
    async addCustomerRequest(customerRequest: CustomerRequestBody) {
      const response = await addCustomerRequest(customerRequest)
      if (response.data) {
        this.customerRequestsMyRequests.push({
          // random id for lib standard
          customer_request_id: Math.random().toString(36).substring(2, 9),
          latitude: customerRequest.latitude,
          longitude: customerRequest.longitude,
          description: customerRequest.description,
          status: 'pending',
          distance: 0, // Will be calculated by backend
          created_at: new Date().toISOString(),
          updated_at: new Date().toISOString(),
        })
        this.errorCustomerRequest = ''
      } else {
        this.errorCustomerRequest = response.message
      }
    },
    async fetchCustomerMyRequests() {
      try {
        const response = await fetchCustomerMyRequests()
        this.customerRequestsMyRequests = response.data ?? []
      } catch (error) {
        throw error
      }
    },
    async fetchCustomerRequestsAll(page = 1, limit = 10, maxDistance = 20) {
      try {
        const response = await fetchCustomerRequestAll(page, limit, maxDistance)
        this.customerRequestsAll = response.data ?? []
        this.total = response.total
        this.page = response.page
        this.limit = response.limit
        this.totalPages = response.total_pages
        this.maxDistance = maxDistance
      } catch (error) {
        throw error
      }
    },
    // Method to update max distance and refresh data
    async updateMaxDistance(newMaxDistance: number) {
      this.maxDistance = newMaxDistance
      await this.fetchCustomerRequestsAll(this.page, this.limit, newMaxDistance)
    },
    // Method to accept a customer request
    async acceptRequest(customerRequestId: string) {
      try {
        await acceptCustomerRequest(customerRequestId)
        // Update the status in the local state
        const request = this.customerRequestsAll.find(
          (req) => req.customer_request_id === customerRequestId,
        )
        if (request) {
          request.status = 'accepted'
          request.updated_at = new Date().toISOString()
        }
      } catch (error) {
        console.error('Failed to accept customer request:', error)
        throw error
      }
    },
    // Method to cancel a customer request
    async cancelRequest(customerRequestId: string, reason: string) {
      // Validate reason length
      if (reason.length > 100) {
        throw new Error('Reason must not exceed 100 characters')
      }
      if (!reason.trim()) {
        throw new Error('Reason is required')
      }

      try {
        await cancelCustomerRequest(customerRequestId, reason)
        // Update the status in the local state
        const request = this.customerRequestsAll.find(
          (req) => req.customer_request_id === customerRequestId,
        )
        if (request) {
          request.status = 'cancelled'
          request.updated_at = new Date().toISOString()
        }
      } catch (error) {
        console.error('Failed to cancel customer request:', error)
        throw error
      }
    },
    // Method to complete a customer request
    async completeRequest(customerRequestId: string) {
      try {
        await completeCustomerRequest(customerRequestId)
        // Update the status in the local state
        const request = this.customerRequestsAll.find(
          (req) => req.customer_request_id === customerRequestId,
        )
        if (request) {
          request.status = 'done'
          request.updated_at = new Date().toISOString()
        }
      } catch (error) {
        console.error('Failed to complete customer request:', error)
        throw error
      }
    },
  },
})

export interface CustomerRequestBody {
  latitude: number
  longitude: number
  description: string
}

export interface CustomerRequest {
  customer_request_id: string
  latitude: number
  longitude: number
  description: string
  status: string
  distance: number // Distance in kilometers
  created_at: string
  updated_at: string
}

export interface CustomerRequestResponse {
  data: CustomerRequest[]
  message: string
}

export interface PaginatedCustomerRequestResponse {
  data: CustomerRequest[]
  total: number
  page: number
  limit: number
  total_pages: number
}

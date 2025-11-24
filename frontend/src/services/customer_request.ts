import { getCookie } from '@/stores/cookie'
import type { CustomerRequestBody, CustomerRequestResponse } from '@/types/customer_request'
import type { ResponseAPI } from '@/types/response'

const addCustomerRequest = async (
  customerRequest: CustomerRequestBody,
): Promise<ResponseAPI<CustomerRequestResponse>> => {
  const apiUrl = import.meta.env.VITE_WEB_API
  const token = getCookie('token')
  if (!apiUrl || !token) {
    console.error('API configuration error or missing token. Please check environment settings.')
    return { message: 'API configuration error or missing token' }
  }
  const response = await fetch(`${apiUrl}/api/customer-request`, {
    method: 'POST',
    headers: {
      Authorization: `Bearer ${token}`,
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(customerRequest),
  })
  const data = await response.json()
  return data
}

const fetchCustomerMyRequests = async (): Promise<CustomerRequestResponse> => {
  const apiUrl = import.meta.env.VITE_WEB_API
  const token = getCookie('token')
  if (!apiUrl || !token) {
    console.error('API configuration error or missing token. Please check environment settings.')
  }
  const response = await fetch(`${apiUrl}/api/customer-request/my-request`, {
    method: 'GET',
    headers: {
      Authorization: `Bearer ${token}`,
      'Content-Type': 'application/json',
    },
  })
  if (response.ok) {
    const data = await response.json()
    const customerRequestsMyRequests = data as CustomerRequestResponse
    return customerRequestsMyRequests
  } else {
    throw new Error(`HTTP error! Status: ${response.status}`)
  }
}

const fetchCustomerRequestAll = async (
  page: number = 1,
  limit: number = 10,
  maxDistance: number = 20,
): Promise<import('@/types/customer_request').PaginatedCustomerRequestResponse> => {
  const apiUrl = import.meta.env.VITE_WEB_API
  const token = getCookie('token')
  if (!apiUrl || !token) {
    console.error('API configuration error or missing token. Please check environment settings.')
    throw new Error('API configuration error or missing token')
  }
  const response = await fetch(
    `${apiUrl}/api/customer-request/all?page=${page}&limit=${limit}&maxDistance=${maxDistance}`,
    {
      method: 'GET',
      headers: {
        Authorization: `Bearer ${token}`,
        'Content-Type': 'application/json',
      },
    },
  )
  if (response.ok) {
    const data = await response.json()
    return data as import('@/types/customer_request').PaginatedCustomerRequestResponse
  } else {
    throw new Error(`HTTP error! Status: ${response.status}`)
  }
}

const acceptCustomerRequest = async (customerRequestId: string): Promise<ResponseAPI<null>> => {
  const apiUrl = import.meta.env.VITE_WEB_API
  const token = getCookie('token')
  if (!apiUrl || !token) {
    console.error('API configuration error or missing token. Please check environment settings.')
    throw new Error('API configuration error or missing token')
  }

  const response = await fetch(`${apiUrl}/api/customer-request/accept/${customerRequestId}`, {
    method: 'PUT',
    headers: {
      Authorization: `Bearer ${token}`,
      'Content-Type': 'application/json',
    },
  })

  if (response.ok) {
    const data = await response.json()
    return data as ResponseAPI<null>
  } else {
    const errorData = await response.json()
    throw new Error(errorData.message || `HTTP error! Status: ${response.status}`)
  }
}

const cancelCustomerRequest = async (
  customerRequestId: string,
  reason: string,
): Promise<ResponseAPI<null>> => {
  const apiUrl = import.meta.env.VITE_WEB_API
  const token = getCookie('token')
  if (!apiUrl || !token) {
    console.error('API configuration error or missing token. Please check environment settings.')
    throw new Error('API configuration error or missing token')
  }

  const response = await fetch(`${apiUrl}/api/customer-request/cancel/${customerRequestId}`, {
    method: 'PUT',
    headers: {
      Authorization: `Bearer ${token}`,
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({ reason }),
  })

  if (response.ok) {
    const data = await response.json()
    return data as ResponseAPI<null>
  } else {
    const errorData = await response.json()
    throw new Error(errorData.message || `HTTP error! Status: ${response.status}`)
  }
}

const completeCustomerRequest = async (customerRequestId: string): Promise<ResponseAPI<null>> => {
  const apiUrl = import.meta.env.VITE_WEB_API
  const token = getCookie('token')
  if (!apiUrl || !token) {
    console.error('API configuration error or missing token. Please check environment settings.')
    throw new Error('API configuration error or missing token')
  }

  const response = await fetch(`${apiUrl}/api/customer-request/complete/${customerRequestId}`, {
    method: 'PUT',
    headers: {
      Authorization: `Bearer ${token}`,
      'Content-Type': 'application/json',
    },
  })

  if (response.ok) {
    const data = await response.json()
    return data as ResponseAPI<null>
  } else {
    const errorData = await response.json()
    throw new Error(errorData.message || `HTTP error! Status: ${response.status}`)
  }
}

export {
  addCustomerRequest,
  fetchCustomerRequestAll,
  fetchCustomerMyRequests,
  acceptCustomerRequest,
  cancelCustomerRequest,
  completeCustomerRequest,
}

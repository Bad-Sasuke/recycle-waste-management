import config from '@/config'

// Get stocks by shop ID with pagination
export const fetchStocksByShopID = async (
  shopId: string,
  page: number = 1,
  pageSize: number = 10,
) => {
  try {
    const apiUrl = config.webAPI
    if (!apiUrl) {
      console.error('API configuration error')
      return { success: false, message: 'API configuration error', data: [] }
    }

    const response = await fetch(
      `${apiUrl}/api/stocks/shop/${shopId}?page=${page}&page_size=${pageSize}`,
    )

    if (!response.ok) {
      throw new Error(`HTTP error! Status: ${response.status}`)
    }

    const data = await response.json()
    console.log('Stock API Response:', data) // Debug log
    return data
  } catch (error) {
    console.error('Error fetching stocks:', error)
    return { success: false, message: 'Failed to fetch stocks', data: [] }
  }
}

// Get receipts by shop ID with pagination
export const fetchReceiptsByShopID = async (
  shopId: string,
  page: number = 1,
  pageSize: number = 10,
) => {
  try {
    const apiUrl = config.webAPI
    if (!apiUrl) {
      console.error('API configuration error')
      return { success: false, message: 'API configuration error', data: [] }
    }

    const response = await fetch(
      `${apiUrl}/api/receipts/shop/${shopId}?page=${page}&page_size=${pageSize}`,
    )

    if (!response.ok) {
      throw new Error(`HTTP error! Status: ${response.status}`)
    }

    const data = await response.json()
    console.log('Receipt API Response:', data) // Debug log
    return data
  } catch (error) {
    console.error('Error fetching receipts:', error)
    return { success: false, message: 'Failed to fetch receipts', data: [] }
  }
}

export default {
  fetchStocksByShopID,
  fetchReceiptsByShopID,
}

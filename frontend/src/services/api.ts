interface ResponseAPI {
  message: string
  data?: any
  page?: number
  limit?: number
  total_pages?: number
  total_items?: number
}

const fetchData = async (url: string): Promise<ResponseAPI> => {
  try {
    // Make the GET request using Fetch API
    const response = await fetch(url)

    if (!response.ok) {
      // Throw an error if the response is not okay
      throw new Error(`HTTP error! Status: ${response.status}`)
    }

    // Parse the response as JSON
    const data: ResponseAPI = await response.json()

    return data
  } catch (error) {
    // Handle errors (e.g., network issues, invalid JSON, etc.)
    console.error('Error fetching data:', error)

    // Return a default ResponseAPI object in case of an error
    return {
      message: 'An error occurred while fetching data.',
      data: [], // Empty array as a fallback value
    }
  }
}

const postData = async (url: string, data: unknown): Promise<ResponseAPI | string> => {
  try {
    // Make the POST request using Fetch API
    const response = await fetch(url, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(data),
    })

    if (!response.ok) {
      // Throw an error if the response is not okay
      throw new Error(`HTTP error! Status: ${response.status}`)
    }

    // Parse the response as JSON
    const responseData: ResponseAPI = await response.json()

    return responseData
  } catch (error: unknown) {
    // Handle errors (e.g., network issues, invalid JSON, etc.)
    console.error('Error fetching data:', error)

    // Check if the error is an instance of Error and return the message
    if (error instanceof Error) {
      return error.message
    }

    // Return a default error message if the error type is unknown
    return 'An unknown error occurred.'
  }
}

// Helper function to get cookie value
const getCookie = (name: string): string | null => {
  const value = `; ${document.cookie}`
  const parts = value.split(`; ${name}=`)
  if (parts.length === 2) {
    return parts.pop()?.split(';').shift() || null
  }
  return null
}

const postDataWithAuth = async (url: string, data: unknown): Promise<ResponseAPI | string> => {
  try {
    const token = getCookie('token')

    if (!token) {
      throw new Error('No authentication token found')
    }

    // Get API base URL from environment variable
    const apiUrl = import.meta.env.VITE_WEB_API || 'http://localhost:8080'
    const fullUrl = `${apiUrl}${url}`

    // Make the POST request using Fetch API with JWT
    const response = await fetch(fullUrl, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${token}`,
      },
      body: JSON.stringify(data),
    })

    if (!response.ok) {
      // Try to parse error message from response
      const errorData = await response.json().catch(() => ({}))
      throw new Error(errorData.message || `HTTP error! Status: ${response.status}`)
    }

    // Parse the response as JSON
    const responseData: ResponseAPI = await response.json()

    return responseData
  } catch (error: unknown) {
    // Handle errors (e.g., network issues, invalid JSON, etc.)
    console.error('Error posting data with auth:', error)

    // Check if the error is an instance of Error and return the message
    if (error instanceof Error) {
      return error.message
    }

    // Return a default error message if the error type is unknown
    return 'An unknown error occurred.'
  }
}

export { fetchData, postData, postDataWithAuth }
